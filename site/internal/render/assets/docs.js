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
      syncDemos();
    });
  }

  const density = document.querySelector('[data-density-picker]');
  if (density) {
    const items = [...density.querySelectorAll('[role="radio"]')];
    const saved = localStorage.getItem('instrument-density') || 'default';
    const select = (btn) => {
      items.forEach((x) => {
        const on = x === btn;
        x.setAttribute('aria-checked', String(on));
        x.tabIndex = on ? 0 : -1;
      });
      const v = btn.dataset.v;
      if (v === 'default') { delete root.dataset.density; localStorage.removeItem('instrument-density'); }
      else { root.dataset.density = v; localStorage.setItem('instrument-density', v); }
      syncDemos();
    };
    select(items.find((x) => x.dataset.v === saved) || items[1]);
    density.addEventListener('click', (e) => {
      const b = e.target.closest('[role="radio"]');
      if (b) select(b);
    });
    /* Стрелки — часть контракта radiogroup: без них с клавиатуры до
       вариантов не добраться, а роль уже обещала, что можно. */
    density.addEventListener('keydown', (e) => {
      const i = items.indexOf(document.activeElement);
      if (i < 0) return;
      const d = { ArrowRight: 1, ArrowDown: 1, ArrowLeft: -1, ArrowUp: -1 }[e.key];
      if (!d) return;
      e.preventDefault();
      const next = items[(i + d + items.length) % items.length];
      next.focus(); select(next);
    });
  }


  /* ── Столы примеров ────────────────────────────────────────────────────
     Каждый пример — отдельный документ в iframe. Родитель делает две вещи:
     подгоняет высоту под содержимое (иначе кадр либо пустует, либо
     заводит собственную полосу прокрутки — и то и другое читается как
     поломка) и передаёт выбранную тему стола. */

  const frames = new Map();
  document.querySelectorAll('[data-demo]').forEach((fig) => {
    const frame = fig.querySelector('.demo-frame');
    if (frame) frames.set(new URL(frame.src, location).pathname, frame);
  });

  addEventListener('message', (e) => {
    if (e.origin !== location.origin || !e.data || !e.data.demoHeight) return;
    const f = frames.get(e.data.id);
    if (f) f.style.blockSize = e.data.demoHeight + 'px';
  });

  document.querySelectorAll('[data-demo-theme]').forEach((sel) => {
    sel.addEventListener('change', () => {
      const frame = sel.closest('[data-demo]').querySelector('.demo-frame');
      /* Пустое значение — «как у сайта»: стол берёт тему справочника, но
         остаётся в своём документе. */
      const t = sel.value || root.dataset.theme || '';
      frame.contentWindow.postMessage(
        { demoTheme: t, demoDensity: root.dataset.density || '' }, location.origin);
    });
  });

  /* Смена темы или плотности справочника догоняет те столы, которые сами
     ничего не выбрали.

     Объявление функцией, а не const-стрелкой, намеренно: syncDemos
     вызывается из обработчика плотности, который отрабатывает при
     инициализации — то есть ВЫШЕ по файлу. Стрелка в const попадала бы во
     временную мёртвую зону и роняла весь скрипт целиком, вместе с темой,
     поиском и копированием. */
  function syncDemos() {
    document.querySelectorAll('[data-demo]').forEach((fig) => {
      const sel = fig.querySelector('[data-demo-theme]');
      if (sel && sel.value) return;
      const w = fig.querySelector('.demo-frame').contentWindow;
      if (w) w.postMessage(
        { demoTheme: root.dataset.theme || '', demoDensity: root.dataset.density || '' },
        location.origin);
    });
  }

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
      index = await (await fetch('/search.json')).json();
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
        empty.textContent = 'Ничего не найдено';
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
  document.addEventListener('click', async (e) => {
    const btn = e.target.closest('[data-copy]');
    if (!btn) return;
    const code = btn.parentElement.querySelector('code');
    if (!code) return;

    if (!btn.dataset.label) {
      btn.dataset.label = btn.textContent;
      btn.setAttribute('aria-live', 'polite');
    }
    clearTimeout(+btn.dataset.timer || 0);

    try {
      await navigator.clipboard.writeText(code.innerText);
      btn.textContent = 'скопировано';
    } catch (err) {
      /* Незащищённый origin (LAN по http) отклоняет запись молча. Молчать
         в ответ — значит соврать: кнопка выглядела бы сработавшей. */
      btn.textContent = 'не вышло';
    }
    btn.dataset.timer = setTimeout(() => {
      btn.textContent = btn.dataset.label;
    }, 1400);
  });

  /* ── Бургер ──────────────────────────────────────────────────────────────
     Выезд ящика рисует кит по aria-expanded. Сайту остаётся ровно то, что
     кит объявил слоем приложения: переключение атрибута, Escape, клик мимо
     и возврат фокуса. Раньше здесь было только первое — ящик открывался и
     не закрывался ничем, кроме повторного нажатия. */

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
