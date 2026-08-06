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

    const render = (items) => {
      active = -1;
      if (!items.length) {
        box.innerHTML = '<div class="site-result-empty">Ничего не найдено</div>';
      } else {
        box.innerHTML = items.map((p) =>
          `<a class="site-result" role="option" href="${p.r}">` +
          `<div class="site-result-title">${p.t}</div>` +
          (p.g ? `<div class="site-result-group">${p.g}</div>` : '') +
          `</a>`).join('');
      }
      box.hidden = false;
    };

    const search = async () => {
      const q = norm(input.value.trim());
      if (q.length < 2) { box.hidden = true; return; }
      const words = q.split(/\s+/);
      const scored = (await load()).map((p) => {
        const t = norm(p.t), b = norm(p.b), g = norm(p.g || '');
        let s = 0;
        for (const w of words) {
          if (t.startsWith(w)) s += 10;
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
      if (e.key === 'Escape') { box.hidden = true; input.blur(); return; }
      if (!opts.length) return;
      if (e.key === 'ArrowDown' || e.key === 'ArrowUp') {
        e.preventDefault();
        active = (active + (e.key === 'ArrowDown' ? 1 : -1) + opts.length) % opts.length;
        opts.forEach((o, i) => o.setAttribute('aria-selected', String(i === active)));
        opts[active].scrollIntoView({ block: 'nearest' });
      } else if (e.key === 'Enter' && active >= 0) {
        e.preventDefault();
        location.href = opts[active].getAttribute('href');
      }
    });

    document.addEventListener('click', (e) => {
      if (!e.target.closest('.site-search')) box.hidden = true;
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

  document.addEventListener('click', async (e) => {
    const btn = e.target.closest('[data-copy]');
    if (!btn) return;
    const code = btn.parentElement.querySelector('code');
    await navigator.clipboard.writeText(code.innerText);
    const was = btn.textContent;
    btn.textContent = 'скопировано';
    /* Смена подписи для скринридера — событие без содержания, поэтому
       результат объявляется отдельно и вежливо. */
    btn.setAttribute('aria-live', 'polite');
    setTimeout(() => { btn.textContent = was; }, 1400);
  });

  /* ── Бургер ────────────────────────────────────────────────────────────── */

  const burger = document.querySelector('[data-burger]');
  if (burger) {
    burger.addEventListener('click', () => {
      const open = burger.getAttribute('aria-expanded') === 'true';
      burger.setAttribute('aria-expanded', String(!open));
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
