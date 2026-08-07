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
