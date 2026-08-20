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
 *   node tools/audit-run.mjs --jobs 8             вкладок разом (по умолчанию 4)
 *   node tools/audit-run.mjs --mutate            сама проверка, проверенная
 *   node tools/audit-run.mjs --base http://…:4399
 *
 * --mutate задаёт этому гейту тот же вопрос, что cmd/mutate задаёт гейтам на
 * Go: проверка, которую никто не видел красной, — украшение. В cmd/mutate
 * он не помещается — там копируется дерево и запускается двоичный файл, а
 * здесь нужны браузер и живой сервер, — поэтому стенд свой. Кит подменяется
 * перехватом запроса за токенами: страница под проверкой настоящая, отличается
 * один файл.
 */

import { spawn } from 'node:child_process';
import { readFile } from 'node:fs/promises';
import { existsSync } from 'node:fs';

const args = process.argv.slice(2);
const baseIdx = args.indexOf('--base');
const BASE = baseIdx >= 0 ? args[baseIdx + 1] : 'http://localhost:4399';
const FILTER = args.find((a) => a.startsWith('/')) || '';
/* Страницы обходятся НЕСКОЛЬКИМИ вкладками разом.
 *
 * Вкладка почти всё время ждёт: навигация, шрифты, раскладка — около трети
 * прогона уходит на ожидание, в которое процессор не занят ничем. Замер при
 * этом независим постранично: каждая вкладка крутит атрибуты на СВОЁМ
 * документе, и делить между ними нечего.
 *
 * Четыре — замеренное, а не «побольше». Полный прогон: одна вкладка 4 мин,
 * четыре — 1:08, восемь — 1:32. Дальше упор не в ожидание, а в раскладку:
 * каждая вкладка на смене темы или масштаба заставляет пересчитать документ
 * целиком, и восьми на это уже не хватает процессора. */
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
    /* Фоновая вкладка в Chrome тормозится: таймеры разрежаются, а
       requestAnimationFrame в невидимой вкладке может не сработать вовсе.
       При одной вкладке это незаметно, при четырёх прогон встаёт намертво —
       ожидание кадра, который не придёт, упирается в срок команды.

       Флаги снимают именно торможение, а не что-то ещё: вкладки остаются
       невидимыми, но живыми. */
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

async function evalInPage(url, expr, swap) {
  /* Обычно вкладка открывается СРАЗУ на нужном адресе. Пустая с последующим
     переходом нужна одному случаю: перехватчик запроса, поставленный после
     навигации, перехватывать уже нечего. */
  const tab = await (await fetch(
    `http://127.0.0.1:${PORT}/json/new?${encodeURIComponent(swap ? 'about:blank' : url)}`,
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
    const events = new Map();
    ws.onmessage = (e) => {
      const msg = JSON.parse(e.data);
      if (msg.id && pending.has(msg.id)) { pending.get(msg.id)(msg); pending.delete(msg.id); return; }
      if (msg.method && events.has(msg.method)) events.get(msg.method)(msg.params);
    };
    const send = (method, params) => withTimeout(
      new Promise((res) => { const i = ++id; pending.set(i, res); ws.send(JSON.stringify({ id: i, method, params })); }),
      CMD_TIMEOUT, method + ' не ответил');

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
      await withTimeout(loaded, CMD_TIMEOUT, 'страница не загрузилась');
    }

    /* Страницу спрашивают, кто она, и только потом меряют.
     *
     * Плоские 450 мс были ставкой на то, что за это время всё успеет, и при
     * одной вкладке ставка обычно проходила. При четырёх она не проходит:
     * первый же прогон уронил страницу на `document.body` равном null.
     *
     * Ошибка при этом громкая только по везению. Чуть более поздний момент
     * даёт документ с телом и без содержимого — ноль замеров, ноль нарушений
     * и слово «чисто». Ровно так же читается и 404: верный адрес, разобранный
     * документ, девятнадцать байт «404 page not found».
     *
     * Поэтому спрашивается признак, от которого проверка зависит: тег модуля,
     * которым страница поднимает кит. Страницы без него у справочника нет. */
    for (let i = 0; i < 100; i++) {
      const at = await send('Runtime.evaluate', {
        expression: `document.readyState === 'complete' && !!document.body &&
          !!document.querySelector('script[type="module"][src]') && location.href`,
        returnByValue: true,
      });
      const href = at.result?.result?.value;
      if (typeof href === 'string' && href.startsWith(url)) break;
      if (i === 99) throw new Error(`страница так и не стала ${url}: ${href}`);
      await sleep(100);
    }
    // Шрифты и раскладка: измеряются ПИКСЕЛИ, и кегль, набранный запасным
    // шрифтом, даст другую прописную и другие цели.
    await send('Runtime.evaluate', {
      expression: `(async () => {
        try { await document.fonts.ready; } catch {}
        // Таймер, а не кадр: кадра во вкладке, которую никто не показывает,
        // можно ждать вечно — и один прогон уже прождал по девяносто секунд
        // на двух страницах.
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
    // Вкладку закрываем ВСЕГДА. Утёкшая остаётся живой в headless-браузере и
    // продолжает есть память; на восьмидесяти страницах это кончается тем,
    // что тормозит уже сам Chrome, а виноватой выглядит проверка.
    try { await fetch(`http://127.0.0.1:${PORT}/json/close/${tab.id}`); } catch {}
  }
}

/* Очередь на n одновременных работ. Результат возвращается В ПОРЯДКЕ ВХОДА,
   а не завершения: страницы обходятся вперемешку, а отчёт обязан читаться
   одинаково от прогона к прогону. */
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

/* Каждая мутация ломает ровно одно и называет раздел, который обязан
   упасть. `from` обязан существовать в токенах: замена, ничего не заменившая,
   дала бы даровое «поймана» на нетронутом файле. */
const MUTATIONS = [
  {
    name: 'значок отстал от кегля на масштабе 15',
    section: 'пропорции',
    page: '/layout/statusbar/',
    from: 'ступени после округления совпадают. */\n  --size-icon-sm:  16px;',
    to: 'ступени после округления совпадают. */\n  --size-icon-sm:  14px;',
    why: 'ровно тот дефект, из-за которого полосу пересчитывали; ловится ТОЛЬКО свипом по масштабам',
  },
  {
    name: 'подпись ушла ниже порога',
    section: 'контраст',
    page: '/components/actions/button/',
    from: '--text-muted:     light-dark(var(--n-8),  var(--n-6));',
    to: '--text-muted:     light-dark(var(--n-5),  var(--n-6));',
    why: 'приглушённый текст перестаёт брать 4.5:1 на панели',
  },
  {
    name: 'цель нажатия ушла под норму',
    section: 'цели',
    page: '/components/actions/button/',
    from: '--control-h-sm: 26px;',
    to: '--control-h-sm: 14px;',
    why: 'кнопка меньше 24 и без зазора, компенсирующего размер',
  },
  {
    name: 'акцент «глина» побелел под текстом',
    section: 'контраст',
    page: '/components/actions/button/',
    from: '  --a-4: oklch(0.560 0.130 45);',
    to: '  --a-4: oklch(0.880 0.130 45);',
    why: 'ломается ОДИН акцент из четырёх — без свипа по акцентам это невидимо',
  },
];

const proc = await launch();
try {
  const src = await readFile(new URL('./audit.js', import.meta.url), 'utf8');
  if (MUTATE) {
    const tokens = await readFile(new URL('../src/tokens.css', import.meta.url), 'utf8');
    console.log('проверка проверки\n');
    let missed = 0;
    const width = Math.max(...MUTATIONS.map((m) => [...m.name].length));
    for (const m of MUTATIONS) {
      const pad = ' '.repeat(width - [...m.name].length);
      if (!tokens.includes(m.from)) {
        console.log(`  ${m.name}${pad}  ✗ МУТАЦИЯ НЕ ПРИМЕНИЛАСЬ — стенд разошёлся с китом`);
        missed++;
        continue;
      }
      let r;
      try {
        r = await evalInPage(BASE + m.page, EXPR(src),
          { path: '/kit/tokens.css', body: tokens.replace(m.from, m.to) });
      } catch (e) {
        console.log(`  ${m.name}${pad}  ✗ НЕ ВЫПОЛНИЛАСЬ — ${e.message}`);
        missed++;
        continue;
      }
      const bad = Object.values(r[m.section] || {}).reduce((n, v) => n + v.нарушений, 0);
      if (bad) {
        console.log(`  ${m.name}${pad}  · поймана   ${m.section} (${bad})`);
      } else {
        console.log(`  ${m.name}${pad}  ✗ ПРОПУЩЕНА  ${m.section}`);
        console.log(`      ${m.why}`);
        missed++;
      }
    }
    console.log();
    if (missed) {
      console.log(`── дыр в гейте: ${missed} из ${MUTATIONS.length} ──`);
      console.log('Пропущенная мутация означает, что инвариант объявлен, но не охраняется.');
      process.exitCode = 1;
    } else {
      console.log(`· ${MUTATIONS.length} мутаций, все пойманы: матрица меряет то, что обещает`);
    }
    proc.kill();
    process.exit(process.exitCode || 0);
  }
  const list = await pages();
  console.log(`страниц: ${list.length}  ·  ${BASE}  ·  вкладок разом: ${JOBS}\n`);

  let checked = 0, failed = 0;
  const problems = [];

  /* Точка печатается по мере готовности, а список нарушений собирается ПО
     ПОРЯДКУ СТРАНИЦ: порядок завершения вкладок случаен, и отчёт, который
     меняет порядок от прогона к прогону, невозможно сравнить с прошлым. */
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
      problems.push({ страница: p, раздел: '—', что: 'не выполнилось: ' + err });
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
