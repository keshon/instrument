/**
 * audit-run — прогон tools/audit.js по ВСЕМ страницам справочника, headless.
 *
 * Зачем. Сам audit.js вставляется в консоль руками, и это правильно для своего
 * экрана: он проверяет ВАШУ композицию, а не примеры. Но у кита есть ещё одна
 * обязанность — не разъезжаться на собственных восьмидесяти страницах, и
 * проверять их руками по одной невозможно, а значит не будут.
 *
 * Отсюда этот бегунок. Он не заменяет ручной запуск, а закрывает другую дыру:
 * регресс в ките, который виден только на отрисованном и только на какой-то
 * одной из страниц.
 *
 * Зависимостей нет намеренно. Chrome умеет отдавать протокол отладки по
 * сокету, а в Node встроен WebSocket — этого достаточно. Ставить браузерный
 * фреймворк ради одного вызова Runtime.evaluate значило бы завести в проекте
 * дерево зависимостей больше самого кита.
 *
 *   node tools/audit-run.mjs                     весь справочник
 *   node tools/audit-run.mjs /components/         только раздел
 *   node tools/audit-run.mjs --base http://…:4399
 */

import { spawn } from 'node:child_process';
import { readFile } from 'node:fs/promises';
import { existsSync } from 'node:fs';

const args = process.argv.slice(2);
const baseIdx = args.indexOf('--base');
const BASE = baseIdx >= 0 ? args[baseIdx + 1] : 'http://localhost:4399';
const FILTER = args.find((a) => a.startsWith('/')) || '';
const PORT = 9222 + (process.pid % 500);

const CHROME = [
  'C:/Program Files/Google/Chrome/Application/chrome.exe',
  'C:/Program Files (x86)/Google/Chrome/Application/chrome.exe',
  '/usr/bin/google-chrome',
  '/usr/bin/chromium',
].find((p) => existsSync(p));

if (!CHROME) {
  console.error('Chrome не найден. Укажите путь в CHROME внутри audit-run.mjs.');
  process.exit(1);
}

const sleep = (ms) => new Promise((r) => setTimeout(r, ms));

/* Список страниц берётся из НАВИГАЦИИ САЙТА, а не из каталога docs: так
   проверяется ровно то, что читатель может открыть, и новая страница попадает
   под проверку сама, без правки этого файла. */
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
  throw new Error('Chrome не поднялся на порту ' + PORT);
}

/* Тонкий клиент протокола: открыть вкладку, выполнить выражение, вернуть
   значение. Больше от отладочного протокола здесь ничего не нужно.

   У каждой команды СВОЙ СРОК, и вкладка закрывается через finally. Без этих
   двух вещей бегунок вставал намертво: ответ, который не пришёл, оставлял
   промис висеть навсегда, а вкладка и сокет утекали при любом выходе мимо
   последних строк. Проверка, которая молча замирает на середине, хуже
   упавшей — та хотя бы называет страницу. */
const CMD_TIMEOUT = 90_000;

async function evalInPage(url, expr) {
  const tab = await (await fetch(
    `http://127.0.0.1:${PORT}/json/new?${encodeURIComponent(url)}`,
    { method: 'PUT' },
  )).json();

  let ws;
  try {
    ws = new WebSocket(tab.webSocketDebuggerUrl);
    await withTimeout(
      new Promise((res, rej) => { ws.onopen = res; ws.onerror = rej; }),
      CMD_TIMEOUT, 'сокет не открылся');

    let id = 0;
    const pending = new Map();
    ws.onmessage = (e) => {
      const msg = JSON.parse(e.data);
      if (msg.id && pending.has(msg.id)) { pending.get(msg.id)(msg); pending.delete(msg.id); }
    };
    const send = (method, params) => withTimeout(
      new Promise((res) => { const i = ++id; pending.set(i, res); ws.send(JSON.stringify({ id: i, method, params })); }),
      CMD_TIMEOUT, method + ' не ответил');

    await send('Page.enable', {});
    await send('Runtime.enable', {});
    await sleep(450);                     // дать шрифтам и раскладке устояться

    const out = await send('Runtime.evaluate', {
      expression: expr, awaitPromise: true, returnByValue: true,
    });

    if (out.result?.exceptionDetails) throw new Error(out.result.exceptionDetails.text);
    return out.result?.result?.value;
  } finally {
    try { ws?.close(); } catch {}
    // Вкладку закрываем ВСЕГДА. Утёкшая остаётся живой в headless-браузере и
    // продолжает есть память; на восьмидесяти страницах это кончается тем,
    // что тормозит уже сам Chrome, а виноватой выглядит проверка.
    try { await fetch(`http://127.0.0.1:${PORT}/json/close/${tab.id}`); } catch {}
  }
}

function withTimeout(p, ms, what) {
  let t;
  return Promise.race([
    p.finally(() => clearTimeout(t)),
    new Promise((_, rej) => { t = setTimeout(() => rej(new Error(what + ` за ${ms / 1000}с`)), ms); }),
  ]);
}

const EXPR = (auditSrc) => `(async () => {
  ${auditSrc}
  // Поповеры и меню закрыты по умолчанию, а закрытого нет в раскладке —
  // проверка их просто не увидит. Раскрываем всё, что умеет раскрываться.
  //
  // Открытым останется ОДИН: поповеры типа auto образуют стек, и открытие
  // соседнего закрывает предыдущий. Это ограничение платформы, а не недосмотр,
  // и обойти его можно было бы только прогоном на каждый поповер по
  // отдельности — то есть умножив полный обход тем и плотностей на число
  // поповеров страницы.
  //
  // Важно другое: закрытые при этом ПРОДОЛЖАЮТ рисоваться, пока гаснут
  // (transition-behavior: allow-discrete), и коробка у них настоящая. Целей у
  // них нет — audit.js их отсеивает по :popover-open, иначе меряется
  // расстояние между пунктами меню, которых на экране уже нет.
  document.querySelectorAll('[popover]').forEach(p => { try { p.showPopover(); } catch {} });
  document.querySelectorAll('details:not([open])').forEach(d => d.open = true);
  await new Promise(r => setTimeout(r, 120));
  const r = await window.kitAudit.run();
  const pack = o => Object.fromEntries(Object.entries(o).map(([k, v]) =>
    [k, { нарушений: v.нарушений, проверено: v.проверено, список: v.список.slice(0, 6) }]));
  return { контраст: pack(r.контраст), цели: pack(r.цели), пропорции: pack(r.пропорции || {}), всего: r.всего };
})()`;

const proc = await launch();
try {
  const src = await readFile(new URL('./audit.js', import.meta.url), 'utf8');
  const list = await pages();
  console.log(`страниц: ${list.length}  ·  ${BASE}\n`);

  let checked = 0, failed = 0;
  const problems = [];

  for (const p of list) {
    let r;
    try {
      r = await evalInPage(BASE + p, EXPR(src));
    } catch (e) {
      problems.push({ страница: p, раздел: '—', что: 'не выполнилось: ' + e.message });
      failed++;
      continue;
    }
    checked += r.всего;
    for (const [раздел, набор] of Object.entries({ контраст: r.контраст, цели: r.цели, пропорции: r.пропорции })) {
      for (const [ключ, v] of Object.entries(набор)) {
        if (!v.нарушений) continue;
        failed += v.нарушений;
        for (const item of v.список) problems.push({ страница: p, раздел: `${раздел}/${ключ}`, что: JSON.stringify(item) });
      }
    }
    process.stdout.write(problems.length ? '×' : '·');
    if (process.env.AUDIT_VERBOSE) console.log(' ' + p);
  }

  console.log('\n');
  if (problems.length) {
    console.log(`── нарушений: ${failed} ──`);
    for (const pr of problems.slice(0, 60)) console.log(`  ${pr.страница}\n    ${pr.раздел}  ${pr.что}`);
    if (problems.length > 60) console.log(`  … и ещё ${problems.length - 60}`);
    process.exitCode = 1;
  } else {
    console.log(`· ${checked} замеров на ${list.length} страницах — чисто`);
  }
} finally {
  proc.kill();
}
