/**
 * behavior — what instrument.js PROMISES the keyboard, checked on a live page.
 *
 * Every other gate reads text. `registry` reads the markup printed in the
 * reference and asks whether the author wrote one roving tabindex; `docscheck`
 * reads tables. None of them runs the script, so none of them can see the one
 * failure mode that matters here: markup that is correct on the page and wrong
 * a keystroke later.
 *
 * That gap was not theoretical. Collapsing a tree node moved the only Tab stop
 * onto a hidden descendant and dropped the whole tree out of the tab order —
 * the defect was published in the reference itself, and seven green gates had
 * nothing to say about it.
 *
 * So this file asks the running page, not the source:
 *
 *   roving    every group the kit owns has exactly one Tab stop
 *   collapse  collapsing a tree node keeps the tree reachable
 *   popover   opening a popover gives its group a roving tabindex
 *   roots     two roots keep two live observers
 *   disabled  an aria-disabled item is walked to and cannot be acted on
 *   expand    clicking a tree twist opens the node it is drawn on
 *
 * Paste it into the console on any page of the reference and call
 * `kitBehavior.run()`. tools/behavior-run.mjs does the same over every page.
 *
 * Sections report {checked, violations, list} — the same shape audit.js uses,
 * so a runner can treat both alike.
 */
(() => {
  /* The kit's own contract, restated. Keeping the list here rather than
     importing it is deliberate: a check that reads its expectations out of the
     thing it checks agrees with any answer that thing gives. */
  const GROUPS = {
    listbox: '[role="option"]',
    tree: '[role="treeitem"]',
    menu: '[role^="menuitem"]',
    radiogroup: '[role="radio"]',
    tablist: '[role="tab"]',
  };
  const GROUP_SELECTOR = Object.keys(GROUPS).map((r) => `[role="${r}"]`).join(',');

  const sleep = (ms) => new Promise((r) => setTimeout(r, ms));

  /* Waiting in frames, not milliseconds, wherever layout has to settle. A
     popover fades out with transition-behavior: allow-discrete, so for a few
     frames after it closes it still has a box and its items still answer
     offsetParent — and a check that waited a fixed 30ms counted a different
     number of items depending on which frame the read landed in. */
  const frames = (n) => new Promise((res) => {
    const step = () => (--n <= 0 ? res() : requestAnimationFrame(step));
    requestAnimationFrame(step);
  });

  /* Items of a group, as the kit counts them.
   *
   * aria-disabled is NOT a reason to leave an item out: a disabled item stays
   * discoverable by arrow, so it takes part in the roving tabindex like any
   * other. `disabled` is, because the platform removes it from the tab order
   * itself. */
  function itemsOf(group, { withHidden = false } = {}) {
    const sel = GROUPS[group.getAttribute('role')];
    return [...group.querySelectorAll(sel)].filter(
      (el) =>
        el.closest(GROUP_SELECTOR) === group &&
        !el.disabled &&
        (withHidden || el.offsetParent !== null),
    );
  }

  const isDisabled = (el) => el.getAttribute('aria-disabled') === 'true';

  const owned = (group) => group.dataset.roving !== 'manual';

  /* Rendered, asked of the element itself.
   *
   * offsetParent is the wrong question here and was asked anyway: it answers
   * null for anything positioned fixed, which is every popover and every
   * pinned toolbar — so a whole class of live groups looked absent. It also
   * cannot tell "hidden at this width" from "box has not arrived yet", and
   * that difference decides whether a skip is nothing to do or a hole. */
  const rendered = (el) => el.checkVisibility();

  /** How many of a group's items a Tab press can land on. */
  const tabStops = (items) => items.filter((el) => el.tabIndex === 0).length;

  function where(el) {
    const id = el.id ? '#' + el.id : '';
    const label = el.getAttribute('aria-label') || '';
    return `${el.getAttribute('role')}${id}${label ? ` "${label}"` : ''}`;
  }

  /* A section always reports how much it looked at. A page with no tree makes
     `collapse` report zero checks, and zero checks is a fact worth printing:
     without it a section that silently found nothing to do is indistinguishable
     from a section that passed.

     `skipped` counts the other thing: something that WAS applicable and could
     not be measured. That number has to surface, because a run that quietly
     measured nine things fewer than the one before it is green for the wrong
     reason — and this section drifted by nine before the count existed. */
  const section = () => ({ checked: 0, violations: 0, skipped: 0, list: [], skips: [] });

  /* A skip is recorded by NAME, not just counted. A number that went from 361
     to 352 says something was dropped; it does not say what, and "something"
     cannot be argued with. */
  function skip(s, what, why) {
    s.skipped++;
    if (s.skips.length < 12) s.skips.push({ what, why });
  }
  function fail(s, at, expected, got) {
    s.violations++;
    if (s.list.length < 12) s.list.push({ at, expected, got });
  }

  /* ── roving ──────────────────────────────────────────────────────────────
     One Tab stop per group. A group where every item has one costs a press
     per item to walk past; a group where none has one cannot be reached at
     all. Both are failures of the same promise. */
  function checkRoving(root = document) {
    const s = section();
    for (const group of root.querySelectorAll(GROUP_SELECTOR)) {
      if (!owned(group) || !rendered(group)) continue;
      const items = itemsOf(group);
      if (!items.length) continue;
      s.checked++;
      const stops = tabStops(items);
      if (stops !== 1) {
        fail(s, where(group), '1 tab stop', `${stops} of ${items.length}`);
      }
    }
    return s;
  }

  /* ── collapse ────────────────────────────────────────────────────────────
     State in this kit lives in attributes, and a collapsed node hides its
     descendants with `hidden` — an attribute change. A tab stop that was
     sitting on a descendant goes with it, and the tree becomes unreachable
     while still looking correct.

     The key is dispatched, not simulated by calling the handler: the point is
     that the kit's own listener runs, on the kit's own element, in the kit's
     own order. */
  async function checkCollapse() {
    const s = section();
    for (const tree of document.querySelectorAll('[role="tree"]')) {
      if (!owned(tree) || !rendered(tree)) continue;
      const node = [...tree.querySelectorAll('[role="treeitem"][aria-expanded="true"]')]
        .find((el) => el.closest(GROUP_SELECTOR) === tree && el.offsetParent !== null);
      if (!node) continue;

      s.checked++;
      node.focus();
      node.dispatchEvent(new KeyboardEvent('keydown', { key: 'ArrowLeft', bubbles: true }));
      await sleep(0);

      const stops = tabStops(itemsOf(tree));
      if (stops !== 1) {
        fail(s, where(tree) + ' after collapse', '1 tab stop', String(stops));
      }

      // Put the page back. A check that leaves the tree collapsed changes what
      // every later check and every later page-load screenshot sees.
      node.dispatchEvent(new KeyboardEvent('keydown', { key: 'ArrowRight', bubbles: true }));
      await sleep(0);
    }
    return s;
  }

  /* ── popover ─────────────────────────────────────────────────────────────
     A closed popover is out of layout, so its items are invisible and the kit
     skips them. Opening one is not a childList edit, so nothing wakes the
     observer, and the items keep the native tabIndex 0 they were born with
     until focus happens to land inside.

     Focus is what hides this: the first Tab into a menu repairs it, so by the
     time anyone looks, it looks right. The check therefore opens the popover
     and reads the group WITHOUT focusing anything. */
  /* Open popovers are handed to `disabled` on the way past.
   *
   * Every aria-disabled item in the reference lives in a menu, every menu
   * lives in a popover, and a closed popover is not rendered — so a `disabled`
   * section that only ever read the document reported 0/0 on the one page
   * built to exercise it. Opening popovers is work this section already does;
   * doing it twice would double the slowest part of the run. */
  async function checkPopover(disabled) {
    const s = section();
    for (const pop of document.querySelectorAll('[popover]')) {
      const groups = [...pop.querySelectorAll(GROUP_SELECTOR)].filter(owned);
      if (!groups.length) continue;

      try {
        pop.showPopover();
      } catch {
        continue; // already open, or not connected — not this check's business
      }
      await frames(2);

      // An auto popover closes whatever was open before it, and a closing one
      // keeps its box while it fades. :popover-open is the only answer that
      // does not depend on where in the fade the read landed — and without it
      // the section counted four items on one run and five on the next.
      if (!pop.matches(':popover-open')) {
        try { pop.hidePopover(); } catch {}
        continue;
      }


      for (const group of groups) {
        // Open is not the same as rendered, and the difference is not a
        // timing one. #site-view carries copies of the header menus for the
        // compact layout; at this width they have no box, the kit leaves them
        // alone, and counting them would charge the kit four violations it
        // does not owe. Nothing to measure is not a hole.
        if (!rendered(group)) continue;

        // A rendered group is owed rendered items. An anchored popover can
        // need more than a frame to place them, and two frames was sometimes
        // not enough — that is where nine checks went missing between runs.
        let items = itemsOf(group);
        for (let i = 0; !items.length && i < 10; i++) {
          await frames(1);
          items = itemsOf(group);
        }
        if (!items.length) {
          // Rendered group, no rendered items, ten frames apart. Recorded by
          // name rather than counted: a silent skip is how a check goes
          // quietly empty, and "something was dropped" cannot be argued with.
          skip(s, where(group) + ' in ' + (pop.id || 'an unnamed popover'),
            'the group is rendered, its items never got a box');
          continue;
        }
        s.checked++;
        const stops = tabStops(items);
        if (stops !== 1) {
          fail(s, where(group) + ' in an open popover', '1 tab stop', `${stops} of ${items.length}`);
        }
      }
      checkDisabled(pop, disabled);
      try { pop.hidePopover(); } catch {}
      await frames(1);
    }
    return s;
  }

  /* ── roots ───────────────────────────────────────────────────────────────
     start(root) is documented as taking an arbitrary root, which means two of
     them can be live at once: a dialog that mounts its own subtree, a preview
     pane, a test harness. Each is owed its own observer.

     The failure this catches is quiet and total: the second start() replaces
     the first one's observer, the first root stops updating, and stop() cannot
     reach the orphan because nothing holds it any more. Nothing throws where
     the reader is looking.

     Built from scratch rather than borrowed from the page: the check adds a
     row and waits for the kit to notice, and doing that to a real component
     would leave the page changed.

     The document's own observer is stopped first, and that is not tidiness.
     It watches document.body with subtree:true, so it wakes on every edit
     inside the probes and repairs them — a replaced per-root observer stays
     invisible behind it, and the check reports a pass it did not earn. */
  async function checkRoots(moduleURL) {
    const s = section();
    let kit;
    try {
      kit = await import(moduleURL);
    } catch (e) {
      fail(s, moduleURL, 'module loads', String(e));
      return s;
    }
    if (typeof kit.start !== 'function' || typeof kit.stop !== 'function') {
      skip(s, moduleURL, 'the module has no start/stop to test');
      return s;
    }

    const make = (name) => {
      const host = document.createElement('div');
      host.dataset.behaviorProbe = name;
      host.innerHTML =
        '<div role="tablist" aria-label="' + name + '">' +
        '<button type="button" role="tab" aria-selected="true">one</button>' +
        '<button type="button" role="tab" aria-selected="false">two</button>' +
        '</div>';
      document.body.append(host);
      return host;
    };
    const add = async (host) => {
      const list = host.querySelector('[role="tablist"]');
      const btn = document.createElement('button');
      btn.type = 'button';
      btn.setAttribute('role', 'tab');
      btn.setAttribute('aria-selected', 'false');
      btn.textContent = 'more';
      list.append(btn);
      // Two turns: the observer answers in a microtask, and the kit's own
      // writes inside that turn must settle before the attribute is read.
      await sleep(20);
      return btn.tabIndex;
    };

    // Stopping the document is itself part of what is under test: start and
    // stop are a pair, and a stop that cannot reach its own observer is the
    // same defect seen from the other side.
    kit.stop(document);
    const a = make('probe-a');
    const b = make('probe-b');
    try {
      kit.start(a);
      kit.start(b);

      s.checked++;
      if ((await add(b)) !== -1) fail(s, 'root B, first edit', 'observer alive', 'no roving tabindex');
      s.checked++;
      if ((await add(a)) !== -1) fail(s, 'root A, first edit', 'observer alive', 'no roving tabindex');
      s.checked++;
      if ((await add(b)) !== -1) {
        fail(s, 'root B, after root A started', 'observer alive', 'observer replaced by the second start()');
      }
    } finally {
      try { kit.stop(a); } catch {}
      try { kit.stop(b); } catch {}
      a.remove();
      b.remove();
      try { kit.start(document); } catch {}
    }
    return s;
  }

  /* ── disabled ────────────────────────────────────────────────────────────
     An aria-disabled item is a promise with two halves, and the halves pull
     against each other: it must stay reachable, so the arrows still walk to
     it and the roving tabindex still writes to it; and it must do nothing, so
     neither a click nor Enter may move the group's selection.

     Leaving it out of the group satisfies neither half. The item then never
     receives -1, the native button keeps 0, and the result is the exact
     inverse of what was promised: arrows skip it, Tab stops on it. */
  function checkDisabled(root, s) {
    for (const group of root.querySelectorAll(GROUP_SELECTOR)) {
      if (!owned(group) || !rendered(group)) continue;
      const items = itemsOf(group);
      const off = items.filter(isDisabled);
      if (!off.length) continue;

      for (const el of off) {
        // Half one: written to, therefore counted. An item the kit never
        // touches keeps the tabIndex the platform gave it.
        s.checked++;
        if (!el.hasAttribute('tabindex')) {
          fail(s, where(group) + ' disabled item', 'a roving tabindex', 'no tabindex attribute at all');
          continue;
        }

        /* Half two: acted on, therefore inert.
         *
         * Watched through inst:select rather than through the state
         * attribute. A menu has no state attribute at all — its items are
         * actions, not choices — and that is exactly where aria-disabled is
         * used most, so a check that read attributes would have skipped its
         * whole population and reported the halves it did measure as if they
         * were the promise. */
        s.checked++;
        let fired = false;
        const spy = () => { fired = true; };
        group.addEventListener('inst:select', spy);
        el.click();
        group.removeEventListener('inst:select', spy);
        if (fired) {
          fail(s, where(group) + ' disabled item', 'a click does nothing', 'inst:select was emitted');
        }

        /* Half three: Enter must not be answered with a synthesised click.
         *
         * For a MENU this is the whole promise and the only half there is.
         * Its items are actions rather than choices, so no state moves and
         * inst:select is emitted for nobody — which means half two can never
         * fail there, and a check that stopped at half two would leave the
         * component where aria-disabled is used most entirely unguarded while
         * reporting a number. */
        s.checked++;
        let clicked = false;
        const watch = () => { clicked = true; };
        el.addEventListener('click', watch);
        el.focus();
        el.dispatchEvent(new KeyboardEvent('keydown', { key: 'Enter', bubbles: true }));
        el.removeEventListener('click', watch);
        if (clicked) {
          fail(s, where(group) + ' disabled item', 'Enter does nothing', 'it synthesised a click');
        }
      }
    }
  }

  /* ── expand ──────────────────────────────────────────────────────────────
     The twist carries cursor: pointer and turns ninety degrees on
     aria-expanded. That is a control, drawn. A drawn control that does
     nothing under the mouse is the kit lying about itself in the one language
     a reader cannot help but read. */
  async function checkExpand() {
    const s = section();
    for (const tree of document.querySelectorAll('[role="tree"]')) {
      if (!owned(tree) || !rendered(tree)) continue;
      for (const node of tree.querySelectorAll('[role="treeitem"][aria-expanded]')) {
        const twist = node.querySelector('.inst-tree-twist');
        if (!twist || !rendered(twist)) continue;

        const before = node.getAttribute('aria-expanded');
        s.checked++;
        twist.click();
        await sleep(0);
        const after = node.getAttribute('aria-expanded');
        if (after === before) {
          fail(s, where(tree) + ' twist', `aria-expanded flips from ${before}`, 'unchanged');
        }

        // Put it back, whichever way it went.
        if (after !== before) {
          twist.click();
          await sleep(0);
        }
      }
    }
    return s;
  }

  /* The page names its own module, so the same file works against a served
     site, a local build and an intercepted one. Guessing the path here would
     make the check pass against a kit that is not the one under test. */
  function moduleURL() {
    const tag = document.querySelector('script[type="module"][src]');
    return tag ? tag.src : null;
  }

  window.kitBehavior = {
    async run() {
      // Layout has to be settled before anything is counted. A group that is
      // still mid-layout answers offsetParent with null and is skipped, and a
      // skipped group is a silent hole: the run stays green and quietly
      // checks five things fewer than the run before it. Seen once on a cold
      // server — 356 checks where three later runs said 361.
      try { await document.fonts?.ready; } catch {}
      await frames(2);

      const url = moduleURL();
      const disabled = section();
      checkDisabled(document, disabled);
      const out = {
        roving: checkRoving(),
        collapse: await checkCollapse(),
        popover: await checkPopover(disabled),
        disabled,
        expand: await checkExpand(),
        // Last on purpose: it calls start() on roots of its own, and until the
        // observer is per-root that call takes the document's observer with it.
        //
        // No module tag is not "nothing to do": this section builds its own
        // probes, so on a real page it always has three checks to report.
        // Silence here means the page is not the page that was asked for.
        roots: url ? await checkRoots(url) : (() => {
          const s = section();
          skip(s, 'the document', 'no module script to import — is this the page you asked for?');
          return s;
        })(),
      };
      out.total = Object.values(out).reduce((n, v) => n + (v.checked || 0), 0);
      return out;
    },
  };
})();
