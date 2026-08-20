/**
 * audit-run — running tools/audit.js over EVERY page of the reference,
 * headless.
 *
 * Why. audit.js itself is pasted into the console by hand, and that is the
 * right thing for your own screen: it checks YOUR composition rather than the
 * examples. But the kit has one more duty — not to fall apart on its own
 * eighty-odd pages — and checking those by hand one at a time is impossible,
 * which means it will not happen.
 *
 * Hence this runner. It does not replace the manual pass; it closes a
 * different hole: a regression in the kit that shows up only on what was
 * rendered, and only on one particular page.
 *
 * There are no dependencies, deliberately. Chrome can serve its debugging
 * protocol over a socket, and Node ships a WebSocket — that is enough. Pulling
 * in a browser framework for a single Runtime.evaluate call would mean a
 * dependency tree larger than the kit itself.
 *
 *   node tools/audit-run.mjs                     the whole reference
 *   node tools/audit-run.mjs /components/        one section only
 *   node tools/audit-run.mjs --jobs 8            tabs at once (4 by default)
 *   node tools/audit-run.mjs --mutate            the check, checked
 *   node tools/audit-run.mjs --base http://…:4399
 *
 * --mutate asks this gate the same question cmd/mutate asks the Go gates: a
 * check nobody has seen red is decoration. It does not fit inside cmd/mutate —
 * there a tree is copied and a binary is run, while here a browser and a live
 * server are needed — so the harness is its own. The kit is swapped by
 * intercepting the request for the tokens: the page under test is the real
 * one, and a single file differs.
 */

import { spawn } from 'node:child_process';
import { readFile } from 'node:fs/promises';
import { existsSync } from 'node:fs';

const args = process.argv.slice(2);
const baseIdx = args.indexOf('--base');
const BASE = baseIdx >= 0 ? args[baseIdx + 1] : 'http://localhost:4399';
const FILTER = args.find((a) => a.startsWith('/')) || '';
/* Pages are walked by SEVERAL tabs at once.
 *
 * A tab spends most of its time waiting: navigation, fonts, layout — about a
 * third of the run goes into waiting with the processor doing nothing. The
 * measurement is independent per page meanwhile: every tab turns attributes on
 * ITS OWN document, and there is nothing to share.
 *
 * Four is measured rather than "a bit more". A full pass: one tab 4 min, four
 * 1:08, eight 1:32. Beyond that the limit is not waiting but layout: on every
 * theme or scale switch each tab forces a full recalculation of the document,
 * and there is no longer enough processor for eight of them. */
const MUTATE = args.includes('--mutate');
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
  console.error('Chrome not found. Set the path in CHROME inside audit-run.mjs.');
  process.exit(1);
}

const sleep = (ms) => new Promise((r) => setTimeout(r, ms));

/* The page list comes from the SITE NAVIGATION rather than from the docs
   directory: that way exactly what a reader can open is what gets checked, and
   a new page falls under the check by itself, without editing this file. */
async function pages() {
  const html = await (await fetch(BASE + '/')).text();
  const hrefs = [...html.matchAll(/href="(\/[^"#]*?)"/g)].map((m) => m[1]);
  const uniq = [...new Set(hrefs)].filter(
    (h) => h.endsWith('/') && !h.startsWith('/en/') && h.startsWith(FILTER || '/'),
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
    '--user-data-dir=' + process.env.TEMP + '/instrument-audit-' + process.pid,
    '--window-size=1440,900',
    'about:blank',
  ], { stdio: 'ignore' });

  for (let i = 0; i < 100; i++) {
    try {
      const r = await fetch(`http://127.0.0.1:${PORT}/json/version`);
      if (r.ok) return proc;
    } catch {}
    await sleep(100);
  }
  proc.kill();
  throw new Error('Chrome did not come up on port ' + PORT);
}

/* A thin protocol client: open a tab, evaluate an expression, return the
   value. Nothing else is needed from the debugging protocol here.

   Every command has ITS OWN deadline, and the tab is closed through finally.
   Without those two things the runner wedged solid: an answer that never came
   left a promise hanging forever, and the tab and the socket leaked on any
   exit past the last lines. A check that quietly freezes halfway is worse than
   one that fails — the failing one at least names the page. */
const CMD_TIMEOUT = 90_000;

async function evalInPage(url, expr, swap) {
  /* Normally the tab opens STRAIGHT at the address wanted. A blank one
     followed by a navigation is needed for one case only: an interceptor
     installed after the navigation has nothing left to intercept. */
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
      if (msg.id && pending.has(msg.id)) { pending.get(msg.id)(msg); pending.delete(msg.id); return; }
      if (msg.method && events.has(msg.method)) events.get(msg.method)(msg.params);
    };
    const send = (method, params) => withTimeout(
      new Promise((res) => { const i = ++id; pending.set(i, res); ws.send(JSON.stringify({ id: i, method, params })); }),
      CMD_TIMEOUT, method + ' did not answer');

    await send('Page.enable', {});
    await send('Runtime.enable', {});

    if (swap) {
      await send('Fetch.enable', { patterns: [{ urlPattern: '*' + swap.path + '*' }] });
      events.set('Fetch.requestPaused', async (prm) => {
        await send('Fetch.fulfillRequest', {
          requestId: prm.requestId,
          responseCode: 200,
          responseHeaders: [{ name: 'Content-Type', value: 'text/css; charset=utf-8' }],
          body: Buffer.from(swap.body, 'utf8').toString('base64'),
        });
      });
      const loaded = new Promise((res) => events.set('Page.loadEventFired', res));
      await send('Page.navigate', { url });
      await withTimeout(loaded, CMD_TIMEOUT, 'page did not load');
    }

    /* The page is asked who it is, and only then measured.
     *
     * A flat 450 ms was a bet that everything would be ready by then, and with
     * one tab the bet usually paid. With four it does not: the very first pass
     * dropped a page on `document.body` being null.
     *
     * The error is loud only by luck, at that. A slightly later moment yields
     * a document with a body and no content — zero measurements, zero
     * violations and the word "clean". A 404 reads exactly the same way: a
     * valid address, a parsed document, nineteen bytes of "404 page not
     * found".
     *
     * So what is asked for is the marker the check depends on: the module tag
     * the page raises the kit with. The reference has no page without it. */
    for (let i = 0; i < 100; i++) {
      const at = await send('Runtime.evaluate', {
        expression: `document.readyState === 'complete' && !!document.body &&
          !!document.querySelector('script[type="module"][src]') && location.href`,
        returnByValue: true,
      });
      const href = at.result?.result?.value;
      if (typeof href === 'string' && href.startsWith(url)) break;
      if (i === 99) throw new Error(`page never became ${url}: ${href}`);
      await sleep(100);
    }
    // Fonts and layout: PIXELS are what gets measured, and type set in a
    // fallback font gives a different cap height and different targets.
    await send('Runtime.evaluate', {
      expression: `(async () => {
        try { await document.fonts.ready; } catch {}
        // A timer rather than a frame: a frame in a tab nobody is showing can
        // be waited for forever — and one pass already waited ninety seconds
        // each on two pages.
        await new Promise(r => setTimeout(r, 60));
      })()`,
      awaitPromise: true,
    });

    const out = await send('Runtime.evaluate', {
      expression: expr, awaitPromise: true, returnByValue: true,
    });

    if (out.result?.exceptionDetails) throw new Error(out.result.exceptionDetails.text);
    return out.result?.result?.value;
  } finally {
    try { ws?.close(); } catch {}
    // The tab is ALWAYS closed. A leaked one stays alive in the headless
    // browser and keeps eating memory; over eighty pages that ends with Chrome
    // itself slowing down while the check looks to blame.
    try { await fetch(`http://127.0.0.1:${PORT}/json/close/${tab.id}`); } catch {}
  }
}

/* A queue of n concurrent jobs. Results come back IN INPUT ORDER rather than
   in completion order: pages are walked out of order, and the report has to
   read the same way from run to run. */
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
    new Promise((_, rej) => { t = setTimeout(() => rej(new Error(what + ` after ${ms / 1000}s`)), ms); }),
  ]);
}

const EXPR = (auditSrc) => `(async () => {
  ${auditSrc}
  // Popovers and menus are closed by default, and what is closed is not in
  // layout — the check simply will not see them. Open everything that opens.
  //
  // ONE will stay open: popovers of type auto form a stack, and opening a
  // sibling closes the previous one. That is a platform constraint rather than
  // an oversight, and the only way round it would be a pass per popover — that
  // is, multiplying the full sweep of themes and densities by the number of
  // popovers on the page.
  //
  // What matters more: the closed ones KEEP being painted while they fade
  // (transition-behavior: allow-discrete), and their box is real. They have no
  // targets — audit.js filters them out by :popover-open, otherwise the
  // distance gets measured between menu items that are no longer on screen.
  document.querySelectorAll('[popover]').forEach(p => { try { p.showPopover(); } catch {} });
  document.querySelectorAll('details:not([open])').forEach(d => d.open = true);
  await new Promise(r => setTimeout(r, 120));
  const r = await window.kitAudit.run();
  const pack = o => Object.fromEntries(Object.entries(o).map(([k, v]) =>
    [k, { failed: v.failed, checked: v.checked, list: v.list.slice(0, 6) }]));
  return { contrast: pack(r.contrast), targets: pack(r.targets), proportion: pack(r.proportion || {}), total: r.total };
})()`;

/* Every mutation breaks exactly one thing and names the section that has to
   fail. `from` has to exist in the tokens: a replacement that replaced nothing
   would hand out a free "caught" on an untouched file. */
const MUTATIONS = [
  {
    name: 'icon lagging behind the type size at scale 15',
    section: 'proportion',
    page: '/layout/statusbar/',
    from: 'steps coincide after rounding. */\n  --size-icon-sm:  16px;',
    to: 'steps coincide after rounding. */\n  --size-icon-sm:  14px;',
    why: 'exactly the defect the band was recomputed for; caught ONLY by the scale sweep',
  },
  {
    name: 'label pushed below the threshold',
    section: 'contrast',
    page: '/components/actions/button/',
    from: '--text-muted:     light-dark(var(--n-8),  var(--n-6));',
    to: '--text-muted:     light-dark(var(--n-5),  var(--n-6));',
    why: 'muted text stops taking 4.5:1 on a panel',
  },
  {
    name: 'tap target pushed below the minimum',
    section: 'targets',
    page: '/components/actions/button/',
    from: '--control-h-sm: 26px;',
    to: '--control-h-sm: 14px;',
    why: 'a button under 24 and without the gap that compensates for the size',
  },
  {
    name: 'the clay accent went white under the text',
    section: 'contrast',
    page: '/components/actions/button/',
    from: '  --a-4: oklch(0.560 0.130 45);',
    to: '  --a-4: oklch(0.880 0.130 45);',
    why: 'ONE accent of four breaks — without the accent sweep this is invisible',
  },
];

const proc = await launch();
try {
  const src = await readFile(new URL('./audit.js', import.meta.url), 'utf8');
  if (MUTATE) {
    const tokens = await readFile(new URL('../src/tokens.css', import.meta.url), 'utf8');
    console.log('checking the check\n');
    let missed = 0;
    const width = Math.max(...MUTATIONS.map((m) => [...m.name].length));
    for (const m of MUTATIONS) {
      const pad = ' '.repeat(width - [...m.name].length);
      if (!tokens.includes(m.from)) {
        console.log(`  ${m.name}${pad}  ✗ MUTATION DID NOT APPLY — the harness has drifted from the kit`);
        missed++;
        continue;
      }
      let r;
      try {
        r = await evalInPage(BASE + m.page, EXPR(src),
          { path: '/kit/tokens.css', body: tokens.replace(m.from, m.to) });
      } catch (e) {
        console.log(`  ${m.name}${pad}  ✗ DID NOT RUN — ${e.message}`);
        missed++;
        continue;
      }
      const bad = Object.values(r[m.section] || {}).reduce((n, v) => n + v.failed, 0);
      if (bad) {
        console.log(`  ${m.name}${pad}  · caught    ${m.section} (${bad})`);
      } else {
        console.log(`  ${m.name}${pad}  ✗ MISSED   ${m.section}`);
        console.log(`      ${m.why}`);
        missed++;
      }
    }
    console.log();
    if (missed) {
      console.log(`── holes in the gate: ${missed} of ${MUTATIONS.length} ──`);
      console.log('A missed mutation means the invariant is declared but not guarded.');
      process.exitCode = 1;
    } else {
      console.log(`· ${MUTATIONS.length} mutations, all caught: the matrix measures what it promises`);
    }
    proc.kill();
    process.exit(process.exitCode || 0);
  }
  const list = await pages();
  console.log(`pages: ${list.length}  ·  ${BASE}  ·  tabs at once: ${JOBS}\n`);

  let checked = 0, failed = 0;
  const problems = [];

  /* A dot is printed as each page finishes, while the list of violations is
     collected IN PAGE ORDER: the order tabs finish in is arbitrary, and a
     report that changes order from run to run cannot be compared with the last
     one. */
  const measured = await pool(list, JOBS, async (p) => {
    try {
      const r = await evalInPage(BASE + p, EXPR(src));
      process.stdout.write('·');
      return { p, r };
    } catch (e) {
      process.stdout.write('×');
      return { p, err: e.message };
    }
  });

  for (const { p, r, err } of measured) {
    if (err) {
      problems.push({ page: p, section: '—', what: 'did not run: ' + err });
      failed++;
      continue;
    }
    checked += r.total;
    for (const [section, set] of Object.entries({ contrast: r.contrast, targets: r.targets, proportion: r.proportion })) {
      for (const [key, v] of Object.entries(set)) {
        if (!v.failed) continue;
        failed += v.failed;
        for (const item of v.list) problems.push({ page: p, section: `${section}/${key}`, what: JSON.stringify(item) });
      }
    }
    if (process.env.AUDIT_VERBOSE) console.log(' ' + p);
  }

  console.log('\n');
  if (problems.length) {
    console.log(`── violations: ${failed} ──`);
    for (const pr of problems.slice(0, 60)) console.log(`  ${pr.page}\n    ${pr.section}  ${pr.what}`);
    if (problems.length > 60) console.log(`  … and ${problems.length - 60} more`);
    process.exitCode = 1;
  } else {
    console.log(`· ${checked} measurements across ${list.length} pages — clean`);
  }
} finally {
  proc.kill();
}
