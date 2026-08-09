/* Скрипт сайта документации. В кит не входит.
   Ровно четыре работы: тема, плотность, поиск, копирование. */

(() => {
  const root = document.documentElement;

  /* ── Тема и плотность ───────────────────────────────────────────────────
     Сайт переключает их теми же атрибутами, что и любое приложение на ките.
     Это не демонстрация ради демонстрации: если темизация где-то протекает,
     видно будет здесь первым. */

  const theme = document.querySelector('[data-theme-picker]');
  if (theme) {
    theme.value = localStorage.getItem('instrument-theme') || '';
    theme.addEventListener('change', () => {
      if (theme.value) {
        root.dataset.theme = theme.value;
        localStorage.setItem('instrument-theme', theme.value);
      } else {
        delete root.dataset.theme;
        localStorage.removeItem('instrument-theme');
      }
    });
  }

  const density = document.querySelector('[data-density-picker]');
  if (density) {
    const items = [...density.querySelectorAll('[role="radio"]')];
    const saved = localStorage.getItem('instrument-density') || 'default';
    let current = saved;
    const select = (btn) => {
      current = btn.dataset.v;
      items.forEach((x) => {
        const on = x === btn;
        x.setAttribute('aria-checked', String(on));
        x.tabIndex = on ? 0 : -1;
      });
      const v = btn.dataset.v;
      if (v === 'default') { delete root.dataset.density; localStorage.removeItem('instrument-density'); }
      else { root.dataset.density = v; localStorage.setItem('instrument-density', v); }
    };
    select(items.find((x) => x.dataset.v === saved) || items[1]);
    density.addEventListener('click', (e) => {
      const b = e.target.closest('[role="radio"]');
      if (b) select(b);
    });
    /* Стрелок здесь больше нет: их ставит сам кит по role="radiogroup".
       Сайт слушает только результат — смену aria-checked. Это и есть
       проверка поведения: если бы оно не работало, плотность перестала бы
       переключаться с клавиатуры на глазах.

       Контракт кита — атрибуты, а не события: он ставит aria-checked, а что
       по этому поводу делать, решает приложение. */
    density.addEventListener('keyup', (e) => {
      if (!/^Arrow|^Home$|^End$/.test(e.key)) return;
      const checked = items.find((x) => x.getAttribute('aria-checked') === 'true');
      if (checked && checked.dataset.v !== current) select(checked);
    });
  }


  /* ── Сцены примеров ────────────────────────────────────────────────────

     Пример живёт в потоке страницы, и своя тема ему достаётся атрибутом на
     сцене: тема кита объявлена как [data-theme] и работает на любом
     поддереве. Высоту подгонять не нужно — сцена и есть содержимое.

     Выбор запоминается на весь справочник, а не на один пример: человек,
     сравнивающий компонент в тёмной теме, идёт по страницам и не хочет
     переключать её на каждом примере заново. */

  const DEMO_THEME = 'instrument-demo-theme';
  let demoTheme = '';
  try { demoTheme = localStorage.getItem(DEMO_THEME) || ''; } catch (e) {}

  /* Пустое значение — «как у сайта»: сцена снимает свой атрибут и наследует
     тему справочника. Плотность у сцены собственной не бывает — она глобальна,
     и пример обязан показывать ту, которую выбрал читатель. */
  const paintDemos = (theme) => {
    document.querySelectorAll('[data-demo-stage]').forEach((stage) => {
      if (theme) stage.dataset.theme = theme;
      else delete stage.dataset.theme;
    });
    document.querySelectorAll('[data-demo-theme]').forEach((sel) => { sel.value = theme; });
  };

  document.querySelectorAll('[data-demo-theme]').forEach((sel) => {
    sel.addEventListener('change', () => {
      demoTheme = sel.value;
      try { localStorage.setItem(DEMO_THEME, demoTheme); } catch (e) {}
      paintDemos(demoTheme);
    });
  });

  paintDemos(demoTheme);

  /* ── Поиск ──────────────────────────────────────────────────────────────
     Индекс — один JSON на весь сайт, грузится по первому обращению.
     Готовый поисковик тянул бы бинарник и WASM на каждую страницу; здесь
     весь кит помещается в несколько сотен килобайт. */

  const input = document.querySelector('[data-search]');
  const box = document.querySelector('[data-results]');
  if (input && box) {
    let index = null, active = -1;

    const load = async () => {
      if (index) return index;
      /* Индекс свой у каждого языка: искать по русским телам из английской
         версии значило бы отдавать читателю страницы, которых он не
         прочтёт. Язык берётся из документа, а не из адреса — он там уже
         объявлен, и второй источник разошёлся бы. */
      const lang = document.documentElement.lang || 'ru';
      const file = lang === 'ru' ? '/search.json' : '/' + lang + '-search.json';
      index = await (await fetch(file)).json();
      return index;
    };

    const norm = (s) => s.toLowerCase().replace(/ё/g, 'е');

    /* Разметка собирается узлами, а не строкой. innerHTML со вставкой t/g/r
       сегодня безопасен, потому что содержимое своё, но это мина: первый же
       заголовок страницы с «<» или кавычкой ломает результаты молча. Проект,
       у которого выдуманный токен ловится машиной, не имеет права оставлять
       такое на «пока не случалось». */
    const setOpen = (on) => {
      box.hidden = !on;
      input.setAttribute('aria-expanded', String(on));
      if (!on) {
        active = -1;
        input.removeAttribute('aria-activedescendant');
      }
    };

    const render = (items) => {
      active = -1;
      input.removeAttribute('aria-activedescendant');
      box.replaceChildren();
      if (!items.length) {
        const empty = document.createElement('div');
        empty.className = 'site-result-empty';
        empty.textContent = RU ? 'Ничего не найдено' : 'Nothing found';
        box.append(empty);
      } else {
        items.forEach((p, i) => {
          const a = document.createElement('a');
          a.className = 'site-result';
          a.id = 'site-result-' + i;
          a.setAttribute('role', 'option');
          a.setAttribute('aria-selected', 'false');
          a.href = p.r;
          const t = document.createElement('div');
          t.className = 'site-result-title';
          t.textContent = p.t;
          a.append(t);
          if (p.g) {
            const g = document.createElement('div');
            g.className = 'site-result-group';
            g.textContent = p.g;
            a.append(g);
          }
          box.append(a);
        });
      }
      setOpen(true);
    };

    const search = async () => {
      const q = norm(input.value.trim());
      if (q.length < 2) { setOpen(false); return; }
      const words = q.split(/\s+/);
      /* Веса расставлены по тому, ЧТО человек набрал, а не по тому, где
         совпало. Точное имя класса — это адрес, а не поиск: тот, кто ввёл
         «inst-badge», знает, чего хочет, и ему нельзя показывать пять
         страниц, где это слово просто упомянуто.

         Слаг стоит наравне с заголовком, потому что для латинского запроса
         он И ЕСТЬ заголовок: документация по-русски, API по-английски. */
      const scored = (await load()).map((p) => {
        const t = norm(p.t), s0 = norm(p.s || ''), g = norm(p.g || '');
        const own = norm(p.o || '').split(' ');
        const names = norm(p.n || '').split(' ');
        const b = norm(p.b);
        let s = 0;
        for (const w of words) {
          /* ОПИСЫВАЕТ важнее, чем УПОМИНАЕТ, и разрыв должен быть таким,
             чтобы никакое число совпадений в теле его не перекрыло.
             `inst-btn` встречается на 26 страницах: без этого разрыва
             первой выходила конституция, где перечислено всё. */
          if (own.includes(w)) s += 120;
          else if (own.some((n) => n.startsWith(w))) s += 60;
          else if (names.includes(w)) s += 12;
          else if (names.some((n) => n.startsWith(w))) s += 6;

          if (s0 === w) s += 100;                       // слаг целиком: dialog → Модалка
          else if (s0.startsWith(w)) s += 40;

          if (t.startsWith(w)) s += 12;
          else if (t.includes(w)) s += 6;
          if (g.includes(w)) s += 2;
          if (b.includes(w)) s += 1;
        }
        return { p, s };
      }).filter((x) => x.s > 0).sort((a, b2) => b2.s - a.s).slice(0, 8);
      render(scored.map((x) => x.p));
    };

    input.addEventListener('input', search);
    input.addEventListener('focus', () => { if (input.value.trim().length >= 2) search(); });

    input.addEventListener('keydown', (e) => {
      const opts = [...box.querySelectorAll('.site-result')];
      if (e.key === 'Escape') { setOpen(false); input.blur(); return; }
      if (!opts.length) return;
      if (e.key === 'ArrowDown' || e.key === 'ArrowUp') {
        e.preventDefault();
        active = (active + (e.key === 'ArrowDown' ? 1 : -1) + opts.length) % opts.length;
        opts.forEach((o, i) => o.setAttribute('aria-selected', String(i === active)));
        opts[active].scrollIntoView({ block: 'nearest' });
        /* Фокус остаётся в поле, поэтому «где я» скринридеру сообщает
           aria-activedescendant. Без неё стрелки двигали подсветку молча:
           роль listbox была объявлена и не выполнена. */
        input.setAttribute('aria-activedescendant', opts[active].id);
      } else if (e.key === 'Enter' && active >= 0) {
        e.preventDefault();
        location.href = opts[active].getAttribute('href');
      }
    });

    document.addEventListener('click', (e) => {
      if (!e.target.closest('.site-search')) setOpen(false);
    });

    /* «/» — фокус в поиск. Ожидаемо там, где сидят часами. */
    document.addEventListener('keydown', (e) => {
      if (e.key !== '/' || e.metaKey || e.ctrlKey) return;
      const t = e.target;
      if (t && /^(INPUT|TEXTAREA|SELECT)$/.test(t.tagName)) return;
      e.preventDefault(); input.focus(); input.select();
    });
  }

  /* ── Копирование ───────────────────────────────────────────────────────── */

  /* Подпись возвращается из dataset, а НЕ из текущего textContent: второе
     нажатие в пределах 1400мс запоминало «скопировано» как исходную подпись,
     и кнопка оставалась такой навсегда. Таймер тоже свой на кнопку и
     сбрасывается — иначе два таймера доигрывали вразнобой.

     aria-live объявлен в разметке заранее: регион, созданный одновременно с
     изменением текста, не озвучивается. */
  const RU = (document.documentElement.lang || 'ru') === 'ru';
  const COPY_DONE = RU ? 'скопировано' : 'copied';
  const COPY_FAIL = RU ? 'не вышло' : 'failed';

  /* Копируется значение из data-value, если оно есть, иначе текст соседнего
     <code>. Первое нужно строкам справочника: там копируют одно имя, а не
     всю ячейку. */
  /* textContent, а не innerText: второй зависит от раскладки и у скрытого
     содержимого возвращает пустую строку — блок кода свёрнут в <details>, и
     копирование молча копировало бы ничто. Для <pre> textContent к тому же
     точнее: он отдаёт исходный текст, а не то, как браузер его отрисовал. */
  const copySource = (btn) => {
    if (btn.dataset.value) return btn.dataset.value;
    const code = btn.parentElement.querySelector('code');
    return code ? code.textContent : '';
  };

  document.addEventListener('click', async (e) => {
    const btn = e.target.closest('[data-copy]');
    if (!btn) return;
    const text = copySource(btn);
    if (!text) return;

    clearTimeout(+btn.dataset.timer || 0);
    const icon = btn.querySelector('[data-copy-icon] use');
    const msg = btn.querySelector('[data-copy-msg]');

    let ok = true;
    try {
      await navigator.clipboard.writeText(text);
    } catch (err) {
      /* Незащищённый origin (LAN по http) отклоняет запись молча. Молчать
         в ответ — значит соврать: кнопка выглядела бы сработавшей. */
      ok = false;
    }

    /* Результат сообщается дважды: значком тем, кто смотрит, и словом в
       живой области тем, кто слушает. Смена значка сама по себе — событие
       без содержания. */
    btn.dataset.state = ok ? 'done' : 'failed';
    if (icon && ok) icon.setAttribute('href', '#i-check');
    if (msg) msg.textContent = ok ? COPY_DONE : COPY_FAIL;

    btn.dataset.timer = setTimeout(() => {
      delete btn.dataset.state;
      if (icon) icon.setAttribute('href', '#i-copy');
      if (msg) msg.textContent = '';
    }, 1400);
  });

  /* ── Бургер ──────────────────────────────────────────────────────────────
     Выезд ящика рисует кит по aria-expanded. Сайту остаётся ровно то, что
     кит объявил слоем приложения: переключение атрибута, Escape, клик мимо
     и возврат фокуса. Без последних трёх ящик открывается и не закрывается
     ничем, кроме повторного нажатия. */

  const burger = document.querySelector('[data-burger]');
  const drawer = document.getElementById('sidebar');
  if (burger && drawer) {
    const setNav = (open) => {
      burger.setAttribute('aria-expanded', String(open));
      if (open) {
        const first = drawer.querySelector('a, button');
        if (first) first.focus();
      } else if (drawer.contains(document.activeElement)) {
        burger.focus();
      }
    };
    const isOpen = () => burger.getAttribute('aria-expanded') === 'true';

    burger.addEventListener('click', () => setNav(!isOpen()));

    document.addEventListener('keydown', (e) => {
      if (e.key === 'Escape' && isOpen()) setNav(false);
    });

    /* Клик по подложке. Подложка — псевдоэлемент оболочки, поэтому цели у
       события нет: считаем промахом всё, что не ящик и не кнопка. */
    document.addEventListener('click', (e) => {
      if (!isOpen()) return;
      if (e.target.closest('#sidebar') || e.target.closest('[data-burger]')) return;
      setNav(false);
    });
  }

  /* ── Подсветка текущего пункта оглавления ──────────────────────────────── */

  const toc = [...document.querySelectorAll('.site-toc-item')];
  if (toc.length) {
    const byId = new Map(toc.map((a) => [decodeURIComponent(a.getAttribute('href').slice(1)), a]));
    const spy = new IntersectionObserver((entries) => {
      for (const en of entries) {
        if (!en.isIntersecting) continue;
        toc.forEach((a) => a.removeAttribute('aria-current'));
        const a = byId.get(en.target.id);
        if (a) a.setAttribute('aria-current', 'true');
      }
    }, { root: document.querySelector('.inst-shell-main'), rootMargin: '0px 0px -70% 0px' });
    byId.forEach((_, id) => {
      const el = document.getElementById(id);
      if (el) spy.observe(el);
    });
  }
})();
