/**
 * behavior-run — tools/behavior.js over the whole reference, headless.
 *
 * Same shape and same reasons as audit-run.mjs: the checks are worth pasting
 * into a console by hand for one screen, and worthless as a gate until
 * something runs them over every page without being asked.
 *
 * No dependencies, on purpose. Chrome speaks the debugging protocol over a
 * socket and Node has a WebSocket; a browser framework for one call to
 * Runtime.evaluate would out-weigh the kit.
 *
 *   node tools/behavior-run.mjs                    the whole reference
 *   node tools/behavior-run.mjs /components/       one section
 *   node tools/behavior-run.mjs --jobs 8           tabs at once (4 by default)
 *   node tools/behavior-run.mjs --mutate           the checks, checked
 *   node tools/behavior-run.mjs --base http://…
 *
 * --mutate is the same question cmd/mutate asks of the Go gates: a check
 * nobody has seen red is decoration. It cannot live in cmd/mutate, which
 * copies a tree and runs a binary — this gate needs a browser and a server, so
 * it carries its own stand. The kit is swapped by intercepting the request for
 * the module, which means the page under test is the real page: same markup,
 * same styles, same server, one file different.
 */

import { spawn } from 'node:child_process';
import { readFile } from 'node:fs/promises';
import { existsSync } from 'node:fs';

const args = process.argv.slice(2);
const baseIdx = args.indexOf('--base');
const BASE = baseIdx >= 0 ? args[baseIdx + 1] : 'http://localhost:4322';
const MUTATE = args.includes('--mutate');
const VERBOSE = args.includes('-v');
const FILTER = args.find((a) => a.startsWith('/')) || '';
/* Pages are walked several tabs at a time.
 *
 * A tab spends most of its life waiting — navigation, fonts, layout — and the
 * measurement is independent per page: each tab presses keys on its OWN
 * document, and there is nothing to share.
 *
 * Four is measured, not "more is better": the pixel gate over the same
 * eighty-five pages takes 4 min on one tab, 1:08 on four and 1:32 on eight.
 * Past four the wait is gone and the layout thread is the wall. */
const jobsIdx = args.indexOf('--jobs');
const JOBS = Math.max(1, jobsIdx >= 0 ? Number(args[jobsIdx + 1]) || 4 : 4);
const PORT = 9222 + (process.pid % 500);

const CHROME = [
  'C:/Program Files/Google/Chrome/Application/chrome.exe',
  'C:/Program Files (x86)/Google/Chrome/Application/chrome.exe',
  '/usr/bin/google-chrome',
  '/usr/bin/chromium',
].find((p) => existsSync(p));

if (!CHROME) {
  console.error('Chrome not found. Add its path to CHROME in behavior-run.mjs.');
  process.exit(1);
}

const sleep = (ms) => new Promise((r) => setTimeout(r, ms));
const CMD_TIMEOUT = 90_000;

/* A queue of n concurrent jobs. Results come back IN INPUT ORDER, not in
   completion order: pages finish interleaved, and a report whose order moves
   between runs cannot be compared with the run before. */
async function pool(items, n, fn) {
  const out = new Array(items.length);
  let next = 0;
  await Promise.all(Array.from({ length: Math.min(n, items.length) }, async () => {
    for (;;) {
      const i = next++;
      if (i >= items.length) return;
      out[i] = await fn(items[i], i);
    }
  }));
  return out;
}

function withTimeout(p, ms, what) {
  let t;
  return Promise.race([
    p.finally(() => clearTimeout(t)),
    new Promise((_, rej) => { t = setTimeout(() => rej(new Error(`${what} after ${ms / 1000}s`)), ms); }),
  ]);
}

/* The page list comes from the site's NAVIGATION, not from the docs folder:
   that way a new page is covered without editing this file, and a page nobody
   can open is not covered at all. */
async function pages() {
  /* BOTH LANGUAGES, and the second one is not a formality. What this gate
     measures is the rendered pixel: whether a label fits its column, whether
     a glyph outweighs the type beside it. An English label is a different
     length from the Russian one, so a page that passed in one spelling says
     nothing about the other. */
  const seen = new Set();
  const roots = ['/'];
  for (const root of roots) {
    const r = await fetch(BASE + root);
    if (!r.ok) continue;
    const html = await r.text();
    for (const m of html.matchAll(/href="(\/[^"#]*?)"/g)) {
      seen.add(m[1]);
      /* A two-letter root is the OTHER language, and it is found rather than
         named: which language sits at the bare root is a decision of the site,
         and it has already changed once. */
      if (/^\/[a-z]{2}\/$/.test(m[1]) && !roots.includes(m[1])) roots.push(m[1]);
    }
  }
  const uniq = [...seen].filter(
    (h) => h.endsWith('/') && h.startsWith(FILTER || '/'),
  );
  return uniq.length ? uniq : ['/'];
}

async function launch() {
  const proc = spawn(CHROME, [
    '--headless=new',
    `--remote-debugging-port=${PORT}`,
    '--disable-gpu',
    /* A background tab in Chrome is throttled: timers are thinned out, and
       requestAnimationFrame in an invisible tab may never fire at all. With one
       tab that goes unnoticed; with four the run wedges solid — waiting for a
       frame that will not come runs into the command timeout.

       The flags remove the throttling and nothing else: the tabs stay
       invisible but alive. */
    '--disable-background-timer-throttling',
    '--disable-backgrounding-occluded-windows',
    '--disable-renderer-backgrounding',
    '--no-first-run',
    '--no-default-browser-check',
    '--user-data-dir=' + (process.env.TEMP || '/tmp') + '/instrument-behavior-' + process.pid,
    '--window-size=1440,900',
    'about:blank',
  ], { stdio: 'ignore' });

  for (let i = 0; i < 100; i++) {
    try {
      if ((await fetch(`http://127.0.0.1:${PORT}/json/version`)).ok) return proc;
    } catch {}
    await sleep(100);
  }
  proc.kill();
  throw new Error('Chrome did not come up on port ' + PORT);
}

/* Two ways in, and only one of them needs explaining.
 *
 * Normally the tab is opened ALREADY POINTING at the page, the way
 * audit-run.mjs does it. The blank-then-navigate path exists for --mutate
 * alone: a request interceptor installed after the navigation has nothing
 * left to intercept, and swapping the kit is the whole point there. */
async function inPage(url, expr, swap) {
  const tab = await (await fetch(
    `http://127.0.0.1:${PORT}/json/new?${encodeURIComponent(swap ? 'about:blank' : url)}`,
    { method: 'PUT' },
  )).json();

  let ws;
  try {
    ws = new WebSocket(tab.webSocketDebuggerUrl);
    await withTimeout(
      new Promise((res, rej) => { ws.onopen = res; ws.onerror = rej; }),
      CMD_TIMEOUT, 'socket did not open');

    let id = 0;
    const pending = new Map();
    const events = new Map();
    ws.onmessage = (e) => {
      const msg = JSON.parse(e.data);
      if (msg.id && pending.has(msg.id)) {
        pending.get(msg.id)(msg);
        pending.delete(msg.id);
        return;
      }
      if (msg.method && events.has(msg.method)) events.get(msg.method)(msg.params);
    };
    const send = (method, params) => withTimeout(
      new Promise((res) => {
        const i = ++id;
        pending.set(i, res);
        ws.send(JSON.stringify({ id: i, method, params }));
      }),
      CMD_TIMEOUT, method + ' did not answer');
    const on = (method, fn) => events.set(method, fn);

    await send('Page.enable', {});
    await send('Runtime.enable', {});

    if (swap) {
      await send('Fetch.enable', { patterns: [{ urlPattern: '*' + swap.path + '*' }] });
      on('Fetch.requestPaused', async (p) => {
        await send('Fetch.fulfillRequest', {
          requestId: p.requestId,
          responseCode: 200,
          responseHeaders: [{ name: 'Content-Type', value: 'text/javascript; charset=utf-8' }],
          body: Buffer.from(swap.body, 'utf8').toString('base64'),
        });
      });
    }

    if (swap) {
      const loaded = new Promise((res) => on('Page.loadEventFired', res));
      await send('Page.navigate', { url });
      await withTimeout(loaded, CMD_TIMEOUT, 'page did not load');
    }

    /* The page is asked who it is before a word is measured.
     *
     * A load event is not proof that the page arrived, and neither are the
     * address and readyState: a 404 answers all three. Right URL, parsing
     * finished, nineteen bytes of "404 page not found" and no page at all —
     * which reads as four sections of 0/0, zero violations and the word
     * "clean". Not hypothetical: a second build writing into site/dist under
     * a running server serves exactly that for a few seconds.
     *
     * So the last condition is the one this gate actually depends on: the
     * module tag it imports the kit from. A page without it is not a page
     * this check can say anything about, and saying nothing quietly is the
     * failure being guarded against.
     *
     * `roots` is the tell that found it: that section builds its own probes,
     * so it cannot honestly report 0/0 on any real page. */
    let ready = null;
    for (let i = 0; i < 100; i++) {
      const at = await send('Runtime.evaluate', {
        expression: `(document.readyState === "complete") && location.href +
          (document.querySelector('script[type="module"][src]') ? ' +module' : '')`,
        returnByValue: true,
      });
      ready = at.result?.result?.value;
      if (typeof ready === 'string' && ready.startsWith(url) && ready.endsWith('+module')) break;
      if (i === 99) throw new Error(`never became ${url} with a module tag — last seen: ${ready}`);
      await sleep(100);
    }
    await sleep(150); // let the module run

    const out = await send('Runtime.evaluate', {
      expression: expr, awaitPromise: true, returnByValue: true,
    });
    if (out.result?.exceptionDetails) throw new Error(out.result.exceptionDetails.text);
    return out.result?.result?.value;
  } finally {
    try { ws?.close(); } catch {}
    // Always close the tab. A leaked one stays alive in the headless browser
    // and keeps eating memory; across eighty pages that ends with Chrome
    // itself being slow and the check looking guilty.
    try { await fetch(`http://127.0.0.1:${PORT}/json/close/${tab.id}`); } catch {}
  }
}

const EXPR = (src) => `(async () => {
  ${src}
  return await window.kitBehavior.run();
})()`;

/* ── the stand ──────────────────────────────────────────────────────────────
   Each mutation breaks exactly one promise and names the section that owes it.
   `from` must exist in the kit: a replacement that changed nothing would score
   a free catch on an untouched file. */
const MUTATIONS = [
  {
    name: 'every item is a tab stop',
    section: 'roving',
    page: '/components/actions/segmented/',
    from: '    for (const el of items) el.tabIndex = el === active ? 0 : -1;',
    to: '    for (const el of items) el.tabIndex = 0;',
    why: 'a list of two hundred rows costs two hundred presses to walk past',
  },
  {
    name: 'no item is a tab stop',
    section: 'roving',
    page: '/components/actions/segmented/',
    from: '    for (const el of items) el.tabIndex = el === active ? 0 : -1;',
    to: '    for (const el of items) el.tabIndex = -1;',
    why: 'a group nobody can reach is worse than one that costs presses',
  },
  {
    name: 'collapsing forgets the tab stop',
    section: 'collapse',
    page: '/agent/tree/',
    from: '\n  roving(group, spec, item);\n}',
    to: '\n}',
    why: 'the tab stop stays on a hidden descendant and the tree leaves the tab order',
  },
  {
    name: 'opening a popover wakes nothing',
    section: 'popover',
    page: '/components/overlays/menu/',
    from: "  root.addEventListener('toggle', onToggle, true);\n",
    to: '',
    why: 'menu items keep the native tabIndex 0 until focus happens to land inside',
  },
  {
    name: 'the second start() takes the first root observer',
    section: 'roots',
    page: '/components/actions/segmented/',
    from: '    observers.get(root)?.disconnect();',
    to: '    observers.get([...observers.keys()][0])?.disconnect();',
    why: 'one root stops updating silently, and stop() can no longer reach the orphan',
  },
  {
    name: 'a disabled item is dropped from the group',
    section: 'disabled',
    page: '/components/overlays/menu/',
    from: '      el.closest(GROUP_SELECTOR) === group &&\n      !el.disabled &&',
    to: "      el.closest(GROUP_SELECTOR) === group &&\n      el.getAttribute('aria-disabled') !== 'true' &&\n      !el.disabled &&",
    why: 'never written to, the native button keeps tabIndex 0 — arrows skip it and Tab stops on it',
  },
  {
    name: 'a disabled item still selects',
    section: 'disabled',
    page: '/components/actions/chip/',
    from: '  if (isDisabled(to)) return;',
    to: '  if (false) return;',
    why: 'the kit answers a press on "you cannot" by calling the application',
  },
  {
    name: 'Enter presses a disabled item anyway',
    section: 'disabled',
    page: '/components/overlays/menu/',
    from: '      if (!isDisabled(item)) item.click();',
    to: '      item.click();',
    why: 'the kit answers a press on "you cannot" by calling the application on the human behalf',
  },
  {
    name: 'the twist is drawn pressable and does nothing',
    section: 'expand',
    page: '/agent/tree/',
    from: "  const twist = e.target.closest?.('.inst-tree-twist');",
    to: "  const twist = null && e.target.closest?.('.inst-tree-twist');",
    why: 'cursor: pointer and a ninety-degree turn promise a control the mouse cannot use',
  },
];

function pack(r) {
  return Object.entries(r)
    .filter(([, v]) => v && typeof v === 'object')
    .map(([k, v]) => `${k} ${v.checked - v.violations}/${v.checked}` + (v.skipped ? `+${v.skipped}?` : ''))
    .join(' · ');
}

const proc = await launch();
try {
  const src = await readFile(new URL('./behavior.js', import.meta.url), 'utf8');
  const kitSrc = await readFile(new URL('../src/kit.js', import.meta.url), 'utf8');

  if (MUTATE) {
    console.log('checks, checked\n');
    let missed = 0;
    const width = Math.max(...MUTATIONS.map((m) => m.name.length));
    for (const m of MUTATIONS) {
      if (!kitSrc.includes(m.from)) {
        console.log(`  ${m.name.padEnd(width)}  x DID NOT APPLY — the stand and the kit have drifted apart`);
        missed++;
        continue;
      }
      const body = kitSrc.replace(m.from, m.to);
      let r;
      try {
        r = await inPage(BASE + m.page, EXPR(src), { path: '/kit/kit.js', body });
      } catch (e) {
        console.log(`  ${m.name.padEnd(width)}  x DID NOT RUN — ${e.message}`);
        missed++;
        continue;
      }
      const s = r[m.section];
      if (s && s.violations > 0) {
        console.log(`  ${m.name.padEnd(width)}  · caught    ${m.section}`);
        if (VERBOSE) for (const it of s.list) console.log(`      ${it.at}: expected ${it.expected}, got ${it.got}`);
      } else {
        console.log(`  ${m.name.padEnd(width)}  x MISSED    ${m.section}`);
        console.log(`      ${m.why}`);
        missed++;
      }
    }
    console.log();
    if (missed) {
      console.log(`-- holes: ${missed} of ${MUTATIONS.length} --`);
      console.log('A missed mutation means the promise is stated but not guarded.');
      process.exitCode = 1;
    } else {
      console.log(`· ${MUTATIONS.length} mutations, all caught: every promise on the list is guarded`);
    }
  } else {
    const list = await pages();
    console.log(`pages: ${list.length}  ·  ${BASE}  ·  tabs at once: ${JOBS}\n`);

    let checked = 0, violations = 0, skipped = 0;
    const problems = [];
    const measured = await pool(list, JOBS, async (p) => {
      try {
        const r = await inPage(BASE + p, EXPR(src));
        process.stdout.write('·');
        return { p, r };
      } catch (e) {
        process.stdout.write('x');
        return { p, err: e.message };
      }
    });

    for (const { p, r, err } of measured) {
      if (err) {
        problems.push({ page: p, section: '-', what: 'did not run: ' + err });
        violations++;
        continue;
      }
      checked += r.total;
      for (const s of Object.values(r)) {
        if (s && typeof s === 'object') skipped += s.skipped || 0;
      }
      for (const [name, s] of Object.entries(r)) {
        if (!s || typeof s !== 'object' || !s.violations) continue;
        violations += s.violations;
        for (const it of s.list) {
          problems.push({ page: p, section: name, what: `${it.at}: expected ${it.expected}, got ${it.got}` });
        }
      }
      if (VERBOSE) {
        console.log(' ' + p + '  ' + pack(r));
        for (const s of Object.values(r)) {
          if (!s || typeof s !== 'object') continue;
          for (const sk of s.skips || []) console.log(`      skipped ${sk.what} — ${sk.why}`);
        }
      }
    }

    console.log('\n');
    if (problems.length) {
      console.log(`-- violations: ${violations} --`);
      for (const pr of problems.slice(0, 60)) console.log(`  ${pr.page}\n    ${pr.section}  ${pr.what}`);
      if (problems.length > 60) console.log(`  … and ${problems.length - 60} more`);
      process.exitCode = 1;
    } else {
      // A skip is a failure of the RUN, not a footnote to it: it means the
      // gate did not do its job on something it was supposed to measure, and
      // a green line under that is the exact shape of the defect this whole
      // file exists to prevent. Printed even at zero, so the number is worth
      // comparing against the run before.
      if (skipped) {
        console.log(`-- could not measure: ${skipped} --`);
        console.log('Nothing was found wrong, and nothing was checked either. Run with -v for the list.');
        process.exitCode = 1;
      } else {
        console.log(`· ${checked} checks over ${list.length} pages — clean, ${skipped} skipped`);
      }
    }
  }
} finally {
  proc.kill();
}
