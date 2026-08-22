/* The script of the documentation site. Not part of the kit.
   Exactly seven jobs: theme, accent, scale, density, search, copying and
   raising a toast in an example — a toast cannot be shown by markup, it is
   started by a call. */

(() => {
  const root = document.documentElement;

  /* ── Theme, accent, scale and density ───────────────────────────────────
     The site switches them with the same attributes as any application on the
     kit. This is not demonstration for its own sake: if the theming leaks
     anywhere, it will be seen here first.

     The four knobs are orthogonal: 5 themes × 4 accents × 5 scales × 3
     densities = 300 combinations, and the reference has to be able to show any
     of them. That is precisely why the switches are independent rather than
     gathered into one list of presets: a preset would hide that the axes do
     not intersect, and the first leaking combination would be found by a user
     rather than here.

     Scale was added as the fourth axis because three densities could not say
     "larger": the comfortable one gave air rather than size. */

  /* An attribute on the root: an empty value removes it, otherwise it is set.
     One function for scale, theme and accent, because the difference between
     them is exactly the name of the attribute; three identical blocks would
     drift apart on the first edit.

     The dropdown is a menu of the kit in a popover rather than a native
     <select>: a list drawn by the operating system was the one place in this
     panel the kit did not reach.

     The state is led by the SITE rather than by the kit, and that is a
     contract rather than a workaround: the kit has `follows: null` on the menu
     role, because a menu item is an action rather than a choice. Here an item
     is a choice after all, so aria-checked is carried over here. The arrows
     and Escape stay with the kit. */
  const menuKnob = (sel, attr, key) => {
    const menu = document.querySelector(sel);
    if (!menu) return;
    const pop = menu.closest('[popover]');
    const trigger = document.querySelector('[popovertarget="' + pop.id + '"]');
    const items = [...menu.querySelectorAll('[role="menuitem"]')];

    const apply = (v) => {
      if (v) { root.setAttribute(attr, v); localStorage.setItem(key, v); }
      else { root.removeAttribute(attr); localStorage.removeItem(key); }
      for (const x of items) {
        const on = (x.dataset.v || '') === (v || '');
        x.setAttribute('aria-checked', String(on));
        /* The label of the button is what is selected. Otherwise the panel
           has to be opened to recall where one stopped. */
        if (on && trigger) trigger.textContent = x.dataset.label || x.textContent.trim();
      }
    };

    apply(localStorage.getItem(key) || '');

    menu.addEventListener('click', (e) => {
      const item = e.target.closest('[role="menuitem"]');
      if (!item) return;
      apply(item.dataset.v || '');
      pop.hidePopover();
    });
  };

  menuKnob('[data-scale-picker]', 'data-scale', 'instrument-scale');
  menuKnob('[data-theme-picker]', 'data-theme', 'instrument-theme');
  menuKnob('[data-accent-picker]', 'data-accent', 'instrument-accent');

  /* Density stayed a segmented control: three short words, and all three are
     visible at once — hiding them in a list would trade one click for two. On
     scale the choice is labelled by type size rather than by a word, and there
     a list wins.

     Both knobs stand apart in the panel rather than gathered into one list of
     five positions: a preset would hide that the axes do not intersect, and
     the combination "large and compact" — the most useful one on a big monitor
     — would become inexpressible. */
  const picker = (sel, attr, key) => {
    const group = document.querySelector(sel);
    if (!group) return;
    const items = [...group.querySelectorAll('[role="radio"]')];
    const apply = (v) => {
      if (v === 'default') { delete root.dataset[attr]; localStorage.removeItem(key); }
      else { root.dataset[attr] = v; localStorage.setItem(key, v); }

      /* A DEMO STAGE SHOWS ITSELF AT COMFORTABLE, and the page around it does
         not. The kit's default is the denser one and stays so — but an example
         is the component being looked at rather than chrome being operated,
         and the comfortable geometry is what it wants.

         The moment a visitor picks a density, that choice has to reach the
         examples too, or the picker would move the page and leave the very
         thing it exists to demonstrate alone. So the stage carries the pinned
         value only while the choice is the default. */
      if (attr !== 'density') return;
      for (const stage of document.querySelectorAll('.demo-root')) {
        if (v === 'default') stage.dataset.density = 'comfortable';
        else delete stage.dataset.density;
      }
    };

    /* The initial state is set by the site: it alone knows what lies in
       localStorage. After that the kit leads the choice. */
    const saved = localStorage.getItem(key) || 'default';
    const start = items.find((x) => x.dataset.v === saved)
      || items.find((x) => x.dataset.v === 'default') || items[0];
    for (const x of items) {
      x.setAttribute('aria-checked', String(x === start));
      x.tabIndex = x === start ? 0 : -1;
    }
    apply(start.dataset.v);

    /* From here on an event of the kit, rather than a click handler of our
       own or a peek at the attribute after an arrow key. The kit carries
       aria-checked over on a click and on an arrow alike; what is left to the
       site is what the kit cannot know — what to do with the chosen value.

       This is a check of the behaviour on the live thing: let the kit stop
       working and the density stops switching before your eyes. */
    group.addEventListener('inst:select', (e) => apply(e.target.dataset.v));
  };

  picker('[data-density-picker]', 'density', 'instrument-density');

  /* ── Search ─────────────────────────────────────────────────────────────
     The index is one JSON for the whole site, loaded on the first request. A
     ready-made search engine would pull a binary and WASM onto every page;
     here the whole kit fits into a few hundred kilobytes. */

  const input = document.querySelector('[data-search]');
  const box = document.querySelector('[data-results]');
  if (input && box) {
    let index = null, active = -1;

    const load = async () => {
      if (index) return index;
      /* The index is its own per language: searching Russian bodies from the
         English version would mean handing the reader pages they will not
         read. WHICH file that is comes from the markup rather than from a rule
         spelled out again here: the build names the file, and a second rule
         would drift from the first — it did, the day the base language moved
         from Russian to English. */
      index = await (await fetch(input.dataset.index)).json();
      return index;
    };

    /* THE TWO RUSSIAN LETTERS ARE DATA, and they are the reason this file is
       named in the exception list of the language gate. A reader types «е»
       where a heading is set with «ё», and without the fold «Модалка» is found
       while «Всплывашка» is not. The pair leaves with the last Russian page.
       The one visible phrase below is data of the same kind: it IS the
       translation, not a line awaiting one. */
    const norm = (s) => s.toLowerCase().replace(/ё/g, 'е');

    /* The markup is assembled from nodes rather than from a string. An
       innerHTML with t/g/r pasted in is safe today because the content is our
       own, but it is a mine: the first page heading with a "<" or a quote in
       it breaks the results in silence. A project where an invented token is
       caught by a machine has no right to leave such a thing on "it has not
       happened yet". */
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
      /* The weights are set by WHAT a person typed rather than by where it
         matched. An exact class name is an address rather than a search:
         somebody who entered "inst-badge" knows what they want, and showing
         them five pages where the word is merely mentioned is not allowed.

         The slug stands level with the heading, because for a Latin query it
         IS the heading: the documentation is in Russian, the API in
         English. */
      const scored = (await load()).map((p) => {
        const t = norm(p.t), s0 = norm(p.s || ''), g = norm(p.g || '');
        const own = norm(p.o || '').split(' ');
        const names = norm(p.n || '').split(' ');
        const b = norm(p.b);
        let s = 0;
        for (const w of words) {
          /* DESCRIBES beats MENTIONS, and the gap has to be such that no
             number of matches in a body can outweigh it. `inst-btn` occurs on
             26 pages: without this gap the constitution came first, where
             everything is listed. */
          if (own.includes(w)) s += 120;
          else if (own.some((n) => n.startsWith(w))) s += 60;
          else if (names.includes(w)) s += 12;
          else if (names.some((n) => n.startsWith(w))) s += 6;

          if (s0 === w) s += 100;                       // the whole slug: dialog → the modal page
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
        /* The focus stays in the field, so what tells a screen reader "where
           am I" is aria-activedescendant. Without it the arrows moved the
           highlight in silence: the listbox role was declared and not
           fulfilled. */
        input.setAttribute('aria-activedescendant', opts[active].id);
      } else if (e.key === 'Enter' && active >= 0) {
        e.preventDefault();
        location.href = opts[active].getAttribute('href');
      }
    });

    document.addEventListener('click', (e) => {
      if (!e.target.closest('.site-search')) setOpen(false);
    });

    /* "/" puts the focus into search. Expected where people sit for hours. */
    document.addEventListener('keydown', (e) => {
      if (e.key !== '/' || e.metaKey || e.ctrlKey) return;
      const t = e.target;
      if (t && /^(INPUT|TEXTAREA|SELECT)$/.test(t.tagName)) return;
      e.preventDefault(); input.focus(); input.select();
    });
  }

  /* ── Copying ──────────────────────────────────────────────────────────── */

  /* Copying is done by the KIT: the class inst-copy and data-copy are its
     contract. There is deliberately no handler of our own here. Copying became
     a behaviour of the kit, and a second implementation in the kit's own
     reference would mean the kit documents one thing and shows another. */

  /* A toast in an example. It has no markup: the region is built by the kit
     on a call, so the one honest example is a button that calls it. The
     parameters live in data-demo-toast so that the page stays markup rather
     than script. */
  document.addEventListener('click', async (e) => {
    const btn = e.target.closest('[data-demo-toast]');
    if (!btn) return;
    const { toast } = await import('/kit/kit.js');
    toast(JSON.parse(btn.dataset.demoToast));
  });

  /* ── The burger ─────────────────────────────────────────────────────────
     The slide-out of the drawer is drawn by the kit from aria-expanded. What
     is left to the site is exactly what the kit declared the application's
     layer: toggling the attribute, Escape, a click outside and returning the
     focus. Without the last three the drawer opens and closes by nothing but
     a second press. */

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

    /* A click on the ground. The ground is a pseudo-element of the shell, so
       the event has no target: anything that is neither the drawer nor the
       button counts as a miss. */
    document.addEventListener('click', (e) => {
      if (!isOpen()) return;
      if (e.target.closest('#sidebar') || e.target.closest('[data-burger]')) return;
      setNav(false);
    });
  }

  /* ── Highlighting the current item of the table of contents ───────────── */

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
