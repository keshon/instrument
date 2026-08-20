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

    /* The load event is not proof that THIS page loaded.
     *
     * The tab starts at about:blank, and about:blank fires a load of its own.
     * On a warm run that has happened before the listener exists; on a cold
     * one it arrives just after, gets taken for the navigation, and the page
     * is measured while it is still blank. The result is not an error — it is
     * four sections reporting 0/0 and a run that says "clean" while checking
     * five things fewer than the run before it.
     *
/* The page is asked who it is before a word is measured.
     *
     * Address and readyState are not enough on their own. A 404 answers both:
     * right URL, parsing finished, nineteen bytes of "404 page not found" and
     * no page at all. That is not hypothetical — a second build writing into
     * site/dist under a running server produces exactly it for a few seconds,
     * and the run before this check existed reported four sections of 0/0 and
     * the word "clean".
     *
     * So the last condition is the one this gate actually depends on: the
     * module tag it imports the kit from. A page without it is not a page this
     * check can say anything about, and saying nothing quietly is the failure
     * being guarded against.
     *
     * `roots` is the tell that found this: it builds its own probes, so it
     * cannot legitimately report 0/0 on any real page. */
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
    console.log(`pages: ${list.length}  ·  ${BASE}\n`);

    let checked = 0, violations = 0, skipped = 0;
    const problems = [];
    for (const p of list) {
      let r;
      try {
        r = await inPage(BASE + p, EXPR(src));
      } catch (e) {
        problems.push({ page: p, section: '-', what: 'did not run: ' + e.message });
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
      process.stdout.write(problems.length ? 'x' : '·');
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
