/* instrument — поведение. Опционально.

   Кит рисует состояние и объявляет роли. Роль — это обещание: `role="listbox"`
   говорит вспомогательной технологии, что стрелки работают. Пока обещание не
   выполнено, компонент не «не доделан», он ОБМАНЫВАЕТ — и это хуже, чем не
   объявлять роль вовсе.

   Этот файл выполняет обещание и больше ничего не делает.

   ── Что он НЕ делает ──────────────────────────────────────────────────────

   Не рисует. Ни одного присвоения style, ни одного добавления класса
   оформления. Он ставит атрибуты, которые УЖЕ есть в разметочном контракте —
   tabindex, aria-selected, aria-checked, aria-expanded, — а как они выглядят,
   решает CSS. Как только скрипт начнёт красить, кит потеряет главное своё
   свойство: приложение перестанет иметь возможность переопределить вид, не
   трогая поведение.

   Не обязателен. Без него всё, что не требует клавиатуры, работает как
   работало: раскрытие на <details>, поповер на Popover API, модалка на
   <dialog>, валидация на :user-invalid.

   ── Подключение ───────────────────────────────────────────────────────────

     <script type="module" src="src/kit.js"></script>

   Модуль, а не классический скрипт: у него есть именованные экспорты, и на
   них потом сядут обёртки Svelte и React. Сборки по-прежнему нет.

   ── Почему делегирование, а не инициализация ──────────────────────────────

   Слушатели висят на документе, а элементы ищутся в момент нажатия. Для
   агентного интерфейса это не оптимизация, а единственный рабочий вариант:
   строки очереди, шаги и узлы дерева ПРИБЫВАЮТ во время работы, и любой
   init(el) при старте промахнулся бы мимо всего, что появилось позже.

   Единственное, что нельзя сделать лениво, — начальный tabindex: группа без
   него не достижима по Tab вообще. Поэтому за появлением групп следит
   MutationObserver, а не человек.
   ========================================================================= */

/* Роли, у которых есть контракт клавиатуры, и всё, что о них надо знать.

   Список закрыт. Он совпадает с таблицей разметочного контракта в
   конституции — если там появится строка, она обязана появиться и здесь,
   иначе кит снова начнёт обещать невыполнимое. */
const GROUPS = {
  listbox: {
    item: '[role="option"]',
    axis: 'vertical',
    // Выделение следует за фокусом: одиночный listbox по APG ведёт себя так,
    // и очередь задач — ровно этот случай.
    follows: 'aria-selected',
  },
  tree: {
    item: '[role="treeitem"]',
    axis: 'vertical',
    follows: 'aria-selected',
  },
  menu: {
    item: '[role^="menuitem"]',
    axis: 'vertical',
    // У пункта меню выделения нет: он действие, а не выбор.
    follows: null,
  },
  radiogroup: {
    item: '[role="radio"]',
    axis: 'horizontal',
    follows: 'aria-checked',
  },
  tablist: {
    item: '[role="tab"]',
    axis: 'horizontal',
    follows: 'aria-selected',
  },
};

const GROUP_SELECTOR = Object.keys(GROUPS).map((r) => `[role="${r}"]`).join(',');

/* Ось берётся из aria-orientation, если он есть: горизонтальное дерево и
   вертикальные вкладки существуют, и разметка вправе так сказать. */
function axisOf(group, spec) {
  return group.getAttribute('aria-orientation') || spec.axis;
}

/* Элементы группы.
 *
 * closest() обязателен: вложенные группы существуют (подменю внутри меню,
 * группа внутри дерева), и без него стрелка в родителе прыгала бы по чужим
 * пунктам. Отключённые и невидимые выбрасываются — фокус на них не ставится
 * ни платформой, ни нами. */
function itemsOf(group, spec) {
  return [...group.querySelectorAll(spec.item)].filter(
    (el) =>
      el.closest(GROUP_SELECTOR) === group &&
      el.getAttribute('aria-disabled') !== 'true' &&
      !el.disabled &&
      el.offsetParent !== null,
  );
}

/* Бегущий tabindex: ровно один элемент группы достижим по Tab.
 *
 * Это и есть весь смысл. Список из двухсот строк, у каждой tabindex="0",
 * заставляет двести раз нажать Tab, чтобы уйти дальше; список, где ни у одной
 * его нет, недостижим вовсе — так кит и жил до сих пор. */
function roving(group, spec, current) {
  const items = itemsOf(group, spec);
  if (!items.length) return;
  const active = current && items.includes(current)
    ? current
    : items.find((el) => el.getAttribute(spec.follows || 'aria-selected') === 'true') || items[0];
  for (const el of items) el.tabIndex = el === active ? 0 : -1;
  return active;
}

/* Перевести фокус и, если роль этого требует, состояние выбора.
 *
 * Выбор снимается со ВСЕЙ группы, а не переключается на цели: два выделенных
 * пункта в одиночном списке — состояние, из которого разметка уже не выйдет. */
function move(group, spec, to) {
  const items = itemsOf(group, spec);
  for (const el of items) {
    el.tabIndex = el === to ? 0 : -1;
    if (spec.follows) el.setAttribute(spec.follows, String(el === to));
  }
  to.focus();
}

function step(items, from, delta) {
  const i = items.indexOf(from);
  if (i < 0) return items[0];
  // По кругу: у списка нет «конца», после которого некуда деться. APG
  // допускает оба поведения, но в очереди задач заворот избавляет от
  // тупика на последней строке.
  return items[(i + delta + items.length) % items.length];
}

/* Разворот и сворачивание узла дерева.
 *
 * Стрелка «вперёд» на свёрнутом узле раскрывает его, на раскрытом — уходит к
 * первому потомку. «Назад» симметрично. Это контракт роли treeitem, и без
 * него aria-expanded в разметке — украшение. */
function treeArrow(group, spec, item, forward) {
  const expandable = item.hasAttribute('aria-expanded');
  const open = item.getAttribute('aria-expanded') === 'true';

  if (forward) {
    if (expandable && !open) {
      item.setAttribute('aria-expanded', 'true');
      return true;
    }
    if (expandable && open) {
      const items = itemsOf(group, spec);
      const next = items[items.indexOf(item) + 1];
      if (next) move(group, spec, next);
      return true;
    }
    return false;
  }

  if (expandable && open) {
    item.setAttribute('aria-expanded', 'false');
    return true;
  }
  // Свёрнутый узел уводит фокус к родителю: уровень берётся из aria-level,
  // потому что вложенность в разметке может быть плоской.
  const level = Number(item.getAttribute('aria-level') || 1);
  const items = itemsOf(group, spec);
  for (let i = items.indexOf(item) - 1; i >= 0; i--) {
    if (Number(items[i].getAttribute('aria-level') || 1) < level) {
      move(group, spec, items[i]);
      return true;
    }
  }
  return false;
}

const NEXT = { vertical: 'ArrowDown', horizontal: 'ArrowRight' };
const PREV = { vertical: 'ArrowUp', horizontal: 'ArrowLeft' };

function onKeydown(e) {
  if (e.defaultPrevented || e.ctrlKey || e.metaKey || e.altKey) return;

  const item = e.target.closest?.(
    '[role="option"],[role="treeitem"],[role^="menuitem"],[role="radio"],[role="tab"]',
  );
  if (!item) return;

  const group = item.closest(GROUP_SELECTOR);
  if (!group) return;
  const spec = GROUPS[group.getAttribute('role')];
  if (!spec) return;

  const axis = axisOf(group, spec);
  const items = itemsOf(group, spec);
  if (!items.length) return;

  // Дерево забирает поперечные стрелки под раскрытие — но только их.
  if (group.getAttribute('role') === 'tree') {
    const forward = e.key === 'ArrowRight';
    if (forward || e.key === 'ArrowLeft') {
      if (treeArrow(group, spec, item, forward)) e.preventDefault();
      return;
    }
  }

  switch (e.key) {
    case NEXT[axis]:
      e.preventDefault();
      move(group, spec, step(items, item, 1));
      return;
    case PREV[axis]:
      e.preventDefault();
      move(group, spec, step(items, item, -1));
      return;
    case 'Home':
      e.preventDefault();
      move(group, spec, items[0]);
      return;
    case 'End':
      e.preventDefault();
      move(group, spec, items[items.length - 1]);
      return;
    case ' ':
    case 'Enter':
      // Радио и вкладка уже выбраны стрелкой — нажимать нечего. Остальным
      // клик отдаётся платформе: обработчик приложения висит на нём, и
      // выдумывать второе событие незачем.
      if (spec.follows && group.getAttribute('role') !== 'radiogroup') {
        move(group, spec, item);
      }
      if (e.key === ' ') e.preventDefault();
      item.click();
      return;
    default:
  }
}

/* Фокус, пришедший мышью, тоже обязан перевесить tabindex: иначе следующий
 * Tab уйдёт не оттуда, куда человек только что ткнул. */
function onFocusin(e) {
  const group = e.target.closest?.(GROUP_SELECTOR);
  if (!group) return;
  const spec = GROUPS[group.getAttribute('role')];
  if (!spec) return;
  const item = e.target.closest(spec.item);
  if (item) roving(group, spec, item);
}

/* ── Тосты ──────────────────────────────────────────────────────────────────

   Здесь скрипт делает то, чего CSS не может в принципе: считает время и
   держит очередь. Вид по-прежнему целиком в CSS — ни одного присвоения style.

   Область открывается popover'ом ОДИН раз и не закрывается. Причина не в
   удобстве: закрытый popover это display: none, а живой регион в display:
   none не озвучивается — скринридер не услышал бы ни одного тоста.
   ───────────────────────────────────────────────────────────────────────── */

const TOAST_LIMIT = 4;
const TOAST_DURATION = 5000;

function toastRegion(doc) {
  let region = doc.querySelector('.inst-toasts');
  if (region) {
    ensureOpen(region);
    return region;
  }
  region = doc.createElement('div');
  region.className = 'inst-toasts';
  region.setAttribute('popover', 'manual');
  // Живой регион вежливый: тост сообщает результат, а не прерывает работу.
  // Ошибка перебивает — на ней самой стоит role="alert".
  region.setAttribute('aria-live', 'polite');
  region.setAttribute('aria-label', 'Уведомления');
  doc.body.append(region);
  ensureOpen(region);
  return region;
}

/* Открыть область один раз.
 *
 * Перевести её НАВЕРХ верхнего слоя нельзя, и это проверено: hidePopover() с
 * последующим showPopover() — и в одной задаче, и через задачу — не поднимает
 * область над модалкой, открытой showModal(). Модальный диалог остаётся выше.
 *
 * Значит и притворяться незачем: тост живёт выше любого z-index и любого
 * overflow: hidden, но НЕ выше модалки. Ограничение записано в справочнике
 * вместе с тем, что делать вместо.                                        */
function ensureOpen(region) {
  try {
    if (!region.matches(':popover-open')) region.showPopover();
  } catch {
    // Область без атрибута popover — приложение решило иначе, и это его
    // право: тосты будут работать, просто не в верхнем слое.
  }
}

/* Убрать с доигранным переходом.
 *
 * transitionend со страховкой по таймеру: под prefers-reduced-motion переход
 * схлопывается до 0.01ms и событие может не прийти вовсе — без страховки
 * тост остался бы на экране навсегда. */
function removeToast(el) {
  if (el.dataset.state === 'leaving') return;
  el.dataset.state = 'leaving';
  let done = false;
  const finish = () => {
    if (done) return;
    done = true;
    el.remove();
  };
  el.addEventListener('transitionend', finish, { once: true });
  setTimeout(finish, 400);
}

/**
 * Показать тост.
 *
 * Разметку строит кит, но она вся описана в справочнике: если стандартной не
 * хватает, соберите свой узел и положите в `.inst-toasts` сами — очередь и
 * таймер подхватят его так же.
 */
export function toast(options = {}) {
  const {
    title = '',
    text = '',
    tone = '',
    duration = TOAST_DURATION,
    action = null, // { label, onClick }
    doc = document,
  } = options;

  const region = toastRegion(doc);

  const el = doc.createElement('div');
  el.className = 'inst-toast';
  if (tone) el.dataset.tone = tone;
  // Ошибка обязана перебить: она сообщает, что действие НЕ выполнено.
  el.setAttribute('role', tone === 'error' ? 'alert' : 'status');

  const body = doc.createElement('div');
  body.className = 'inst-toast-body';
  if (title) {
    const t = doc.createElement('div');
    t.className = 'inst-toast-title';
    t.textContent = title;
    body.append(t);
  }
  if (text) {
    const s = doc.createElement('div');
    s.className = 'inst-toast-text';
    s.textContent = text;
    body.append(s);
  }
  el.append(body);

  if (action) {
    const wrap = doc.createElement('div');
    wrap.className = 'inst-toast-actions';
    const btn = doc.createElement('button');
    btn.className = 'inst-btn inst-btn--sm inst-btn--ghost';
    btn.type = 'button';
    btn.textContent = action.label;
    btn.addEventListener('click', () => {
      action.onClick?.();
      removeToast(el);
    });
    wrap.append(btn);
    el.append(wrap);
  }

  region.append(el);

  // Потолок очереди. Двадцать тостов подряд — не сообщение, а стена: самые
  // старые уходят, чтобы новое было видно.
  const live = [...region.children].filter((x) => x.dataset.state !== 'leaving');
  for (const old of live.slice(0, Math.max(0, live.length - TOAST_LIMIT))) {
    removeToast(old);
  }

  if (duration > 0) {
    let timer = setTimeout(() => removeToast(el), duration);
    /* Таймер замирает под курсором и на фокусе. Тост, исчезнувший ровно
       тогда, когда его начали читать, — это потерянное сообщение, и WCAG
       2.2.1 требует ровно этого: у времени должна быть пауза. */
    const hold = () => clearTimeout(timer);
    const resume = () => { timer = setTimeout(() => removeToast(el), duration); };
    el.addEventListener('pointerenter', hold);
    el.addEventListener('pointerleave', resume);
    el.addEventListener('focusin', hold);
    el.addEventListener('focusout', resume);
  }

  return el;
}

/** Расставить бегущий tabindex во всех группах внутри root. */
export function refresh(root = document) {
  for (const group of root.querySelectorAll?.(GROUP_SELECTOR) || []) {
    const spec = GROUPS[group.getAttribute('role')];
    if (spec) roving(group, spec);
  }
}

let observer = null;

/** Подключить поведение. Вызывается сам при загрузке модуля. */
export function start(root = document) {
  root.addEventListener('keydown', onKeydown);
  root.addEventListener('focusin', onFocusin);
  refresh(root);

  // Группы прибывают во время работы — это агентный интерфейс, в нём строки
  // очереди появляются по одной. Наблюдатель дешевле, чем требование звать
  // refresh() руками после каждого рендера, и не забывается.
  const target = root.body || root;
  if (target && typeof MutationObserver === 'function') {
    let queued = false;
    observer = new MutationObserver(() => {
      if (queued) return;
      queued = true;
      queueMicrotask(() => {
        queued = false;
        refresh(root);
      });
    });
    observer.observe(target, { childList: true, subtree: true });
  }
}

/** Отключить. Нужно тестам и горячей перезагрузке, а не приложению. */
export function stop(root = document) {
  root.removeEventListener('keydown', onKeydown);
  root.removeEventListener('focusin', onFocusin);
  observer?.disconnect();
  observer = null;
}

if (typeof document !== 'undefined') {
  if (document.readyState === 'loading') {
    document.addEventListener('DOMContentLoaded', () => start(), { once: true });
  } else {
    start();
  }
}
