/* Сверка документации с китом. Обе стороны разом.

   Документация о ките расходится с китом в двух направлениях, и оба
   одинаково дороги:

     вперёд   в примере стоит класс, которого в ките нет — читатель копирует
              разметку, она молча не работает, и виноват оказывается кит;
     назад    класс в ките есть, а страницы у него нет — компонент существует
              и не находится, что для потребителя равно его отсутствию.

   Проверка читает настоящие исходники кита и настоящие страницы, поэтому
   разойтись с ними не может — тот же принцип, что у contrast.mjs.

   Заодно проверяются значения data-атрибутов: словари закрыты (DESIGN.md),
   и опечатка вроде data-state="quued" выглядит ровно как базовое значение,
   то есть молча ничего не делает. Найти её глазами нельзя.               */

import { readFileSync, readdirSync, existsSync, statSync } from 'node:fs';
import { join, relative } from 'node:path';

const P = (u) => new URL(u, import.meta.url).pathname.replace(/^\/([A-Z]:)/, '$1');
const SRC = P('../src/');
const DOCS = P('../docs/');
const ROOT = P('../');

const strip = (css) => css
  .replace(/\/\*[\s\S]*?\*\//g, ' ')
  .replace(/url\([^)]*\)/g, ' ')
  .replace(/"(?:[^"\\]|\\.)*"/g, ' ')
  .replace(/'(?:[^'\\]|\\.)*'/g, ' ');

/* ── Что есть в ките ─────────────────────────────────────────────────────── */
const kit = new Set();
const rawCss = [];
for (const f of readdirSync(SRC).filter(f => f.endsWith('.css'))) {
  const css = readFileSync(join(SRC, f), 'utf8');
  rawCss.push(css);
  for (const line of strip(css).split('\n')) {
    if (!line.includes('{')) continue;
    for (const m of line.slice(0, line.indexOf('{')).matchAll(/\.(-?[A-Za-z_][\w-]*)/g)) kit.add(m[1]);
  }
}
const allCss = rawCss.join('\n');

/* Словари data-атрибутов — из селекторов, а не из головы. */
const vocab = new Map();
for (const m of allCss.matchAll(/\[data-([a-z-]+)="([^"]+)"\]/g)) {
  if (!vocab.has(m[1])) vocab.set(m[1], new Set());
  vocab.get(m[1]).add(m[2]);
}
/* Базовые значения оформления не имеют, поэтому в селекторах их нет.
   Они объявлены в конституции и перечисляются здесь явно — иначе проверка
   ругалась бы ровно на правильную разметку. */
const BASE_VALUES = { state: ['queued', 'todo', 'approved'] };
for (const [attr, vals] of Object.entries(BASE_VALUES)) {
  if (!vocab.has(attr)) vocab.set(attr, new Set());
  vals.forEach(v => vocab.get(attr).add(v));
}

/* ── Что есть в документации ─────────────────────────────────────────────── */
function walk(dir) {
  if (!existsSync(dir)) return [];
  return readdirSync(dir).flatMap(e => {
    const p = join(dir, e);
    return statSync(p).isDirectory() ? walk(p) : (p.endsWith('.md') ? [p] : []);
  });
}

const pages = walk(DOCS).filter(p => !p.includes('internal'));
const documented = new Set();
const problems = [];
const pending = new Map();   /* ссылка → откуда. Ещё не написанные страницы */

for (const p of pages) {
  const md = readFileSync(p, 'utf8');
  const rel = relative(ROOT, p).replace(/\\/g, '/');

  md.split('\n').forEach((line, i) => {
    const at = `${rel}:${i + 1}`;

    for (const m of line.matchAll(/\binst-[a-z0-9-]+/g)) {
      documented.add(m[0]);
      if (!kit.has(m[0])) problems.push(`${at}  класса нет в ките: .${m[0]}`);
    }

    /* Ссылки между страницами. Не «ошибка» и не «метрика», а третий вид:
       ссылка на ещё не написанную страницу — нормальный след недоделанной
       работы, но забытая — дыра, которую никто не заметит, пока читатель в
       неё не провалится. Поэтому они собираются отдельным списком. */
    for (const m of line.matchAll(/\]\((\.[^)#]+\.md)(#[^)]*)?\)/g)) {
      const target = new URL(m[1], 'file:///' + p.replace(/\\/g, '/')).pathname.replace(/^\/([A-Z]:)/, '$1');
      /* Наличие на диске, а не членство в списке страниц: цель может лежать
         и вне docs/ — например, конституция в корне. */
      if (!existsSync(target)) {
        const rt = relative(ROOT, target).replace(/\\/g, '/');
        if (!pending.has(rt)) pending.set(rt, []);
        pending.get(rt).push(at);
      }
    }

    for (const m of line.matchAll(/data-([a-z-]+)="([^"]*)"/g)) {
      const [, attr, val] = m;
      if (['theme', 'density', 'dir'].includes(attr)) continue;
      const known = vocab.get(attr);
      if (!known) { problems.push(`${at}  неизвестный атрибут: data-${attr}`); continue; }
      if (val && !known.has(val)) {
        problems.push(`${at}  data-${attr}="${val}" — нет в словаре (${[...known].sort().join(' ')})`);
      }
    }
  });
}

/* ── Отчёт ───────────────────────────────────────────────────────────────── */
const undocumented = [...kit].filter(c => !documented.has(c)).sort();
const pct = kit.size ? Math.round((kit.size - undocumented.length) / kit.size * 100) : 0;

console.log(`страниц: ${pages.length}  ·  классов в ките: ${kit.size}  ·  покрыто: ${kit.size - undocumented.length} (${pct}%)\n`);

if (problems.length) {
  console.log(`── документация ссылается на несуществующее (${problems.length}) ──`);
  problems.forEach(p => console.log('  ' + p));
  console.log();
}

if (pending.size) {
  console.log(`── ссылки на ещё не написанные страницы (${pending.size}) ──`);
  [...pending].sort().forEach(([t, from]) => console.log(`  ${t}  ← ${from.length}`));
  console.log();
}

if (undocumented.length && process.argv.includes('-v')) {
  console.log(`── без страницы (${undocumented.length}) ──`);
  console.log('  ' + undocumented.map(c => '.' + c).join(' ') + '\n');
} else if (undocumented.length) {
  console.log(`── без страницы: ${undocumented.length} классов (список — с ключом -v) ──\n`);
}

if (!problems.length && !undocumented.length && !pending.size) {
  console.log('· документация и кит сходятся полностью');
}

process.exit(problems.length ? 1 : 0);
