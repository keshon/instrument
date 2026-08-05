#!/usr/bin/env node
/* ==========================================================================
   agentic-ui — проверка контраста

   Читает НАСТОЯЩИЙ src/tokens.css и резолвит семантику так же, как браузер:
   var() рекурсивно, light-dark() по режиму, color-mix(... transparent) в
   альфу, oklch() в sRGB. Поэтому тест не может разойтись с китом — он не
   дублирует значения, он их вычисляет.

   Правило приёмки:
     текст  < 18px   — 4.5:1   (WCAG 1.4.3)
     текст ≥ 18px    — 3.0:1
     несущая граница — 3.0:1   (WCAG 1.4.11): чекбокс, поле, дорожка свитча,
                                индикатор состояния. Декоративная рамка, у
                                которой рядом есть перепад поверхностей, сюда
                                не входит и намеренно тише.

   Запуск:  node tools/contrast.mjs
   ========================================================================== */

import { readFileSync } from 'node:fs';
import { fileURLToPath } from 'node:url';
import { dirname, join } from 'node:path';

const ROOT = join(dirname(fileURLToPath(import.meta.url)), '..');

/* ── Разбор CSS ─────────────────────────────────────────────────────────── */

const css = readFileSync(join(ROOT, 'src/tokens.css'), 'utf8')
  .replace(/\/\*[\s\S]*?\*\//g, '');

const DECLS = new Map();
for (const m of css.matchAll(/(--[\w-]+)\s*:\s*([^;]+);/g)) {
  DECLS.set(m[1], m[2].trim().replace(/\s+/g, ' '));
}

/** Разбить список аргументов по запятым верхнего уровня. */
function splitArgs(s) {
  const out = [];
  let depth = 0, cur = '';
  for (const ch of s) {
    if (ch === '(') depth++;
    else if (ch === ')') depth--;
    if (ch === ',' && depth === 0) { out.push(cur.trim()); cur = ''; }
    else cur += ch;
  }
  if (cur.trim()) out.push(cur.trim());
  return out;
}

/** Содержимое функции: "fn(...)" -> "..." */
function inner(s) {
  return s.slice(s.indexOf('(') + 1, s.lastIndexOf(')'));
}

/* ── Цвет ───────────────────────────────────────────────────────────────── */

function oklchToSrgb(L, C, H, alpha = 1) {
  const h = (H * Math.PI) / 180;
  const a = C * Math.cos(h);
  const b = C * Math.sin(h);
  const l = (L + 0.3963377774 * a + 0.2158037573 * b) ** 3;
  const m = (L - 0.1055613458 * a - 0.0638541728 * b) ** 3;
  const s = (L - 0.0894841775 * a - 1.2914855480 * b) ** 3;
  const enc = (v) => {
    v = v <= 0.0031308 ? 12.92 * v : 1.055 * Math.pow(Math.max(v, 0), 1 / 2.4) - 0.055;
    return Math.min(1, Math.max(0, v));
  };
  return {
    r: enc(4.0767416621 * l - 3.3077115913 * m + 0.2309699292 * s),
    g: enc(-1.2684380046 * l + 2.6097574011 * m - 0.3413193965 * s),
    b: enc(-0.0041960863 * l - 0.7034186147 * m + 1.7076147010 * s),
    a: alpha,
  };
}

/**
 * Разрешить значение токена в цвет sRGB для заданного режима.
 * mode — 'light' | 'dark', им управляется light-dark().
 */
function resolve(value, mode, seen = new Set()) {
  const v = value.trim();

  if (v === 'transparent') return { r: 0, g: 0, b: 0, a: 0 };

  if (v.startsWith('var(')) {
    const name = splitArgs(inner(v))[0];
    if (seen.has(name)) throw new Error(`цикл в ${name}`);
    if (!DECLS.has(name)) throw new Error(`токен ${name} не объявлен`);
    return resolve(DECLS.get(name), mode, new Set([...seen, name]));
  }

  if (v.startsWith('light-dark(')) {
    const [l, d] = splitArgs(inner(v));
    return resolve(mode === 'dark' ? d : l, mode, seen);
  }

  if (v.startsWith('color-mix(')) {
    const args = splitArgs(inner(v));
    if (!/^in oklab$/i.test(args[0])) {
      throw new Error(`поддержан только "in oklab": ${v}`);
    }
    const [c1, p1] = args[1].match(/^(.*?)\s+([\d.]+)%$/).slice(1);
    const pct = parseFloat(p1) / 100;
    const base = resolve(c1, mode, seen);
    const other = resolve(args[2], mode, seen);
    // Смешивание премультиплицированное: с transparent это ровно
    // «тот же цвет с альфой pct».
    const a = base.a * pct + other.a * (1 - pct);
    if (a === 0) return { r: 0, g: 0, b: 0, a: 0 };
    return {
      r: (base.r * base.a * pct + other.r * other.a * (1 - pct)) / a,
      g: (base.g * base.a * pct + other.g * other.a * (1 - pct)) / a,
      b: (base.b * base.a * pct + other.b * other.a * (1 - pct)) / a,
      a,
    };
  }

  const ok = v.match(/^oklch\(\s*([\d.]+)\s+([\d.]+)\s+([\d.]+)\s*(?:\/\s*([\d.]+)\s*)?\)$/);
  if (ok) {
    return oklchToSrgb(+ok[1], +ok[2], +ok[3], ok[4] === undefined ? 1 : +ok[4]);
  }

  throw new Error(`не разобрать цвет: ${v}`);
}

const token = (name, mode) => {
  if (!DECLS.has(name)) throw new Error(`токен ${name} не объявлен`);
  return resolve(DECLS.get(name), mode);
};

/** Наложить fg на bg (оба sRGB, простое смешивание по альфе). */
function composite(fg, bg) {
  const a = fg.a + bg.a * (1 - fg.a);
  return {
    r: fg.r * fg.a + bg.r * (1 - fg.a),
    g: fg.g * fg.a + bg.g * (1 - fg.a),
    b: fg.b * fg.a + bg.b * (1 - fg.a),
    a,
  };
}

/** Схлопнуть стопку токенов в непрозрачный цвет. Первый — непрозрачная база. */
function flatten(stack, mode) {
  return stack.map((n) => token(n, mode)).reduce((bg, fg) => composite(fg, bg));
}

function luminance({ r, g, b }) {
  const f = (v) => (v <= 0.04045 ? v / 12.92 : Math.pow((v + 0.055) / 1.055, 2.4));
  return 0.2126 * f(r) + 0.7152 * f(g) + 0.0722 * f(b);
}

function ratio(fg, bg) {
  const solid = composite(fg, bg);
  const [hi, lo] = [luminance(solid), luminance(bg)].sort((x, y) => y - x);
  return (hi + 0.05) / (lo + 0.05);
}

/* ── Что именно проверяем ───────────────────────────────────────────────── */

const TEXT = 4.5;   // < 18px
const LARGE = 3.0;  // ≥ 18px, а также несущие границы и индикаторы

/** [ описание, токен-переднего-плана, стопка фона, порог ] */
const CASES = [
  // Текст на трёх поверхностях. Каждая ступень проверяется везде, где может
  // оказаться, — именно это правило и не выполнялось раньше.
  ['текст: primary на панели',        '--text-primary',   ['--surface-raised'], TEXT],
  ['текст: primary на странице',      '--text-primary',   ['--surface-page'],   TEXT],
  ['текст: primary во врезе',         '--text-primary',   ['--surface-sunken'], TEXT],
  ['текст: secondary на панели',      '--text-secondary', ['--surface-raised'], TEXT],
  ['текст: secondary во врезе',       '--text-secondary', ['--surface-sunken'], TEXT],
  ['текст: muted на панели',          '--text-muted',     ['--surface-raised'], TEXT],
  ['текст: muted на странице',        '--text-muted',     ['--surface-page'],   TEXT],
  ['текст: muted во врезе (лог)',     '--text-muted',     ['--surface-sunken'], TEXT],
  ['текст: primary в поле ввода',     '--text-primary',   ['--surface-field'],  TEXT],

  // faint — порог декорации, а не чтения. Читаемому тексту он запрещён.
  ['декор: faint на панели',          '--text-faint',     ['--surface-raised'], LARGE],
  ['декор: faint во врезе',           '--text-faint',     ['--surface-sunken'], LARGE],

  // Бейджи: 11px, значит полный текстовый порог.
  ['бейдж: accent на своём фоне',     '--accent-text', ['--surface-raised', '--accent-bg'], TEXT],
  ['бейдж: ok на своём фоне',         '--ok-text',     ['--surface-raised', '--ok-bg'],     TEXT],
  ['бейдж: warn на своём фоне',       '--warn-text',   ['--surface-raised', '--warn-bg'],   TEXT],
  ['бейдж: err на своём фоне',        '--err-text',    ['--surface-raised', '--err-bg'],    TEXT],
  ['бейдж: нейтральный во врезе',     '--text-secondary', ['--surface-sunken'], TEXT],

  // Статусный текст живёт и вне бейджа — .metric-delta, .field-error, .note.
  ['статус: ok-text на панели',       '--ok-text',   ['--surface-raised'], TEXT],
  ['статус: warn-text на панели',     '--warn-text', ['--surface-raised'], TEXT],
  ['статус: err-text на панели',      '--err-text',  ['--surface-raised'], TEXT],
  ['статус: ok-text во врезе',        '--ok-text',   ['--surface-sunken'], TEXT],
  ['статус: warn-text во врезе',      '--warn-text', ['--surface-sunken'], TEXT],
  ['статус: err-text во врезе',       '--err-text',  ['--surface-sunken'], TEXT],

  // Сплошная кнопка. Ховер обязан УВЕЛИЧИВАТЬ контраст подписи, а не ронять.
  ['кнопка: подпись на accent-solid', '--accent-on', ['--accent-solid'], TEXT],
  ['кнопка: подпись на accent-hover', '--accent-on', ['--accent-hover'], TEXT],
  ['ссылка: accent-text на странице', '--accent-text', ['--surface-page'], TEXT],

  // Индикаторы состояния — нетекстовые, но несущие: 3:1.
  ['точка: ok на панели',             '--ok-text',   ['--surface-raised'], LARGE],
  ['точка: ok во врезе',              '--ok-text',   ['--surface-sunken'], LARGE],
  ['точка: warn на панели',           '--warn-text', ['--surface-raised'], LARGE],
  ['точка: warn во врезе',            '--warn-text', ['--surface-sunken'], LARGE],
  ['точка: err на панели',            '--err-text',  ['--surface-raised'], LARGE],
  ['точка: err во врезе',             '--err-text',  ['--surface-sunken'], LARGE],
  ['точка: running на панели',        '--accent-mark', ['--surface-raised'], LARGE],
  ['точка: running во врезе',         '--accent-mark', ['--surface-sunken'], LARGE],
  ['каретка на панели',               '--accent-mark', ['--surface-raised'], LARGE],
  ['бегунок на дорожке',              '--accent-mark', ['--surface-raised', '--track'], LARGE],

  // Несущие границы: граница И ЕСТЬ контрол.
  ['граница контрола на панели',      '--border-control', ['--surface-raised'], LARGE],
  ['граница контрола на странице',    '--border-control', ['--surface-page'],   LARGE],
  ['граница контрола во врезе',       '--border-control', ['--surface-field'],  LARGE],

  // Заполнение меры относительно собственной дорожки, и дорожки — на всех
  // поверхностях, где мера может стоять.
  ['мера: заливка на дорожке (панель)', '--accent-mark', ['--surface-raised', '--track'], LARGE],
  ['мера: заливка на дорожке (врез)',   '--accent-mark', ['--surface-sunken', '--track'], LARGE],
  ['мера: ok на дорожке',               '--ok-text',   ['--surface-raised', '--track'], LARGE],
  ['мера: warn на дорожке',             '--warn-text', ['--surface-raised', '--track'], LARGE],
  ['мера: err на дорожке',              '--err-text',  ['--surface-raised', '--track'], LARGE],

  // Категориальная палитра: каждый ряд обязан отделяться от поверхности.
  ['график: ряд 1 на панели',        '--chart-1', ['--surface-raised'], LARGE],
  ['график: ряд 1 на странице',      '--chart-1', ['--surface-page'],   LARGE],
  ['график: ряд 2 на панели',        '--chart-2', ['--surface-raised'], LARGE],
  ['график: ряд 2 на странице',      '--chart-2', ['--surface-page'],   LARGE],
  ['график: ряд 3 на панели',        '--chart-3', ['--surface-raised'], LARGE],
  ['график: ряд 3 на странице',      '--chart-3', ['--surface-page'],   LARGE],
  ['график: ряд 4 на панели',        '--chart-4', ['--surface-raised'], LARGE],
  ['график: ряд 4 на странице',      '--chart-4', ['--surface-page'],   LARGE],
  ['график: ряд 5 на панели',        '--chart-5', ['--surface-raised'], LARGE],
  ['график: ряд 5 на странице',      '--chart-5', ['--surface-page'],   LARGE],
  ['график: ряд 6 на панели',        '--chart-6', ['--surface-raised'], LARGE],
  ['график: ряд 6 на странице',      '--chart-6', ['--surface-page'],   LARGE],

  // Оверлеи: всё, что лежит на --surface-overlay.
  ['поповер: текст',                  '--text-primary', ['--surface-overlay'], TEXT],
  ['меню: горячая клавиша',           '--text-muted',   ['--surface-overlay'], TEXT],
  ['меню: опасный пункт',             '--err-text',     ['--surface-overlay'], TEXT],
  ['меню: отмеченный пункт',          '--accent-text',  ['--surface-overlay'], TEXT],
  ['тултип: текст',                   '--text-primary', ['--surface-overlay'], TEXT],

  // Баннер: текст поверх тонированной заливки.
  ['баннер ok: заголовок',            '--text-primary',   ['--surface-page', '--ok-bg'],   TEXT],
  ['баннер warn: заголовок',          '--text-primary',   ['--surface-page', '--warn-bg'], TEXT],
  ['баннер error: заголовок',         '--text-primary',   ['--surface-page', '--err-bg'],  TEXT],
  ['баннер warn: пояснение',          '--text-secondary', ['--surface-page', '--warn-bg'], TEXT],
  ['баннер warn: значок',             '--warn-text',      ['--surface-page', '--warn-bg'], LARGE],

  // Формы.
  ['карточка выбора: заголовок',      '--text-primary',   ['--surface-raised', '--accent-bg'], TEXT],
  ['карточка выбора: описание',       '--text-secondary', ['--surface-raised', '--accent-bg'], TEXT],
  ['множественный выбор: выбранный пункт', '--accent-text', ['--surface-field', '--surface-selected'], TEXT],
  ['приставка поля',                  '--text-muted',     ['--surface-sunken'], TEXT],
  ['readonly: текст на врезе',        '--text-primary',   ['--surface-sunken'], TEXT],
  ['пунктир зоны файла',              '--border-control', ['--surface-field'],  LARGE],
  ['обязательность',                  '--err-text',       ['--surface-raised'], TEXT],

  // Раскладка и навигация.
  ['текст на боковой колонке',        '--text-secondary', ['--surface-sunken'], TEXT],
  ['навигация: текущий пункт',        '--accent-text', ['--surface-sunken', '--surface-selected'], TEXT],
  ['навигация: метка у края',         '--accent-solid', ['--surface-sunken', '--surface-selected'], LARGE],
  ['вкладка: подчёркивание',          '--accent-solid', ['--surface-page'], LARGE],
  ['крошки: разделитель',             '--text-faint', ['--surface-page'], LARGE],
  ['пагинация: текущая страница',     '--accent-text', ['--surface-page', '--surface-selected'], TEXT],
  ['шаги: полоса пройденного',        '--accent-mark', ['--surface-page', '--track'], LARGE],

  // Кольцо фокуса — против того, что под ним.
  ['фокус: кольцо на странице',       '--focus-ring', ['--surface-page'],   LARGE],
  ['фокус: кольцо на панели',         '--focus-ring', ['--surface-raised'], LARGE],
  ['фокус: кольцо во врезе',          '--focus-ring', ['--surface-sunken'], LARGE],
];

/* ── Прогон ─────────────────────────────────────────────────────────────── */

let failed = 0;
const rows = [];

for (const mode of ['light', 'dark']) {
  for (const [label, fg, bgStack, min] of CASES) {
    let r, err = null;
    try {
      r = ratio(token(fg, mode), flatten(bgStack, mode));
    } catch (e) {
      err = e.message;
    }
    const ok = err === null && r >= min;
    if (!ok) failed++;
    rows.push({ mode, label, r, min, ok, err });
  }
}

const W = Math.max(...rows.map((x) => x.label.length));
let lastMode = null;
for (const x of rows) {
  if (x.mode !== lastMode) {
    console.log(`\n${x.mode === 'light' ? 'СВЕТЛАЯ' : 'ТЁМНАЯ'} тема`);
    console.log('─'.repeat(W + 24));
    lastMode = x.mode;
  }
  if (x.err) {
    console.log(`  ✗ ${x.label.padEnd(W)}  ОШИБКА: ${x.err}`);
    continue;
  }
  console.log(
    `  ${x.ok ? '·' : '✗'} ${x.label.padEnd(W)}  ${x.r.toFixed(2).padStart(6)}  (нужно ${x.min.toFixed(1)})`
  );
}

console.log();
if (failed) {
  console.log(`✗ провалов: ${failed} из ${rows.length}`);
  process.exit(1);
}
console.log(`· все ${rows.length} проверок пройдены`);
