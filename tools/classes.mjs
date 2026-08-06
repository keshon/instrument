/* Инвентаризация классов. Временный инструмент этапа 0 — удаляется после
   переименования.

   Наивный grep по файлу даёт мусор: `.w3` и `.org` приезжают из
   xmlns='http://www.w3.org/2000/svg' внутри data-URI, `.tokens` и `.base` —
   из имён файлов в комментариях, `.is-up` — из исторической справки.
   Поэтому сначала вырезаются комментарии, строки и url(), и только потом
   ищутся селекторы. */

import { readFileSync, readdirSync } from 'node:fs';
import { join } from 'node:path';

const SRC = new URL('../src/', import.meta.url).pathname.replace(/^\/([A-Z]:)/, '$1');
const ROOT = new URL('../', import.meta.url).pathname.replace(/^\/([A-Z]:)/, '$1');

/* Вырезает то, что не является селектором: комментарии, кавычки, url(). */
function strip(css) {
  return css
    .replace(/\/\*[\s\S]*?\*\//g, ' ')
    .replace(/url\([^)]*\)/g, ' ')
    .replace(/"(?:[^"\\]|\\.)*"/g, ' ')
    .replace(/'(?:[^'\\]|\\.)*'/g, ' ');
}

/* Классы берутся только из ПРЕЛЮДИИ правила — до `{`, после `}` или `;`.
   Тело правила пропускается: там встречается, например, grid-area с точкой. */
function classesOf(css) {
  const out = new Set();
  const clean = strip(css);
  let depth = 0, buf = '';
  for (const ch of clean) {
    if (ch === '{') {
      if (depth === 0) for (const m of buf.matchAll(/\.(-?[A-Za-z_][\w-]*)/g)) out.add(m[1]);
      buf = ''; depth++;
    } else if (ch === '}') {
      depth = Math.max(0, depth - 1); buf = '';
    } else if (depth === 0) {
      buf += ch;
    } else if (depth >= 1) {
      /* Вложенность CSS: внутри правила тоже бывают селекторы (&:hover,
         & > .x). Собираем их из строк, начинающихся с & или содержащих `{`. */
      buf += ch;
      if (ch === ';') buf = '';
    }
  }
  return out;
}

/* Вложенные селекторы собираем отдельным проходом: строки, где встречается
   `&` или которые заканчиваются на `{`. */
function nestedClasses(css) {
  const out = new Set();
  for (const line of strip(css).split('\n')) {
    if (!line.includes('{')) continue;
    const prelude = line.slice(0, line.indexOf('{'));
    for (const m of prelude.matchAll(/\.(-?[A-Za-z_][\w-]*)/g)) out.add(m[1]);
  }
  return out;
}

const kit = new Set();
const perFile = new Map();
for (const f of readdirSync(SRC).filter(f => f.endsWith('.css'))) {
  const css = readFileSync(join(SRC, f), 'utf8');
  const set = new Set([...classesOf(css), ...nestedClasses(css)]);
  perFile.set(f, set);
  set.forEach(c => kit.add(c));
}

/* Леса страницы — из <style> в index.html. */
const html = readFileSync(join(ROOT, 'index.html'), 'utf8');
const styleBlock = html.slice(html.indexOf('<style>'), html.indexOf('</style>'));
const page = new Set([...classesOf(styleBlock), ...nestedClasses(styleBlock)]);

/* Что реально стоит в разметке. */
const used = new Set();
for (const m of html.matchAll(/class="([^"]*)"/g)) {
  for (const c of m[1].split(/\s+/)) if (c) used.add(c);
}

const overlap = [...kit].filter(c => page.has(c)).sort();
const kitUnused = [...kit].filter(c => !used.has(c)).sort();
const pageOnly = [...page].filter(c => !kit.has(c)).sort();
const orphan = [...used].filter(c => !kit.has(c) && !page.has(c)).sort();

const args = process.argv.slice(2);
if (args[0] === '--list') {
  console.log([...kit].sort().join('\n'));
} else if (args[0] === '--page') {
  console.log(pageOnly.join('\n'));
} else {
  console.log(`классов в ките:        ${kit.size}`);
  console.log(`классов в лесах:       ${page.size} (только лесов: ${pageOnly.length})`);
  console.log(`использовано в разметке: ${used.size}`);
  console.log();
  console.log(`── СТОЛКНОВЕНИЯ кит ↔ леса (${overlap.length}) ──`);
  console.log(overlap.length ? overlap.map(c => '  .' + c).join('\n') : '  нет');
  console.log();
  console.log(`── классы кита БЕЗ демо в спецификации (${kitUnused.length}) ──`);
  console.log(kitUnused.length ? '  ' + kitUnused.map(c => '.' + c).join(' ') : '  нет');
  console.log();
  console.log(`── в разметке, но нигде не объявлены (${orphan.length}) ──`);
  console.log(orphan.length ? '  ' + orphan.map(c => '.' + c).join(' ') : '  нет');
  console.log();
  console.log('── по файлам ──');
  for (const [f, s] of perFile) console.log(`  ${f.padEnd(16)} ${s.size}`);
}
