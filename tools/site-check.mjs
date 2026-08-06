/* Проверка собранного сайта. Запускается после `npm --prefix site run build`.

   docs-check сверяет ИСХОДНИКИ: что классы существуют, ссылки ведут в
   существующие файлы, токены не выдуманы. Но между исходником и страницей
   стоит сборка, и она ломает своё:

     · относительные ссылки .md превращаются в маршруты — и могут не
       превратиться, тогда вся перелинковка ведёт в 404;
     · живой пример разворачивается из ограды ```html preview — и может не
       развернуться, тогда страница показывает код без демонстрации;
     · спрайт иконок вклеивается в тело — и может не вклеиться, тогда все
       иконки в примерах пустые.

   Ни одну из трёх поломок не видно на глаз при беглом просмотре: страница
   выглядит целой. Поэтому проверка идёт по файлам, а не по впечатлению. */

import { readdirSync, readFileSync, statSync, existsSync } from 'node:fs';
import { join, relative } from 'node:path';

const P = (u) => new URL(u, import.meta.url).pathname.replace(/^\/([A-Z]:)/, '$1');
const DIST = P('../site/dist/');

if (!existsSync(DIST)) {
  console.error('сайт не собран: нет site/dist. Запустите npm --prefix site run build');
  process.exit(1);
}

const walk = (d) => readdirSync(d).flatMap((e) => {
  const p = join(d, e);
  return statSync(p).isDirectory() ? walk(p) : [p];
});

const pages = walk(DIST).filter((p) => p.endsWith('.html'));
const routes = new Set(
  pages.map((p) => '/' + relative(DIST, p).replace(/\\/g, '/').replace(/index\.html$/, '')),
);

let stages = 0, withSprite = 0;
const problems = [];

for (const p of pages) {
  const html = readFileSync(p, 'utf8');
  const route = '/' + relative(DIST, p).replace(/\\/g, '/').replace(/index\.html$/, '');

  /* Живые примеры развернулись. */
  const st = [...html.matchAll(/<div class="preview-stage"[^>]*>([\s\S]*?)<\/div>\s*<\/div>/g)];
  stages += st.length;
  for (const [, body] of st) {
    if (!body.trim()) problems.push(`${route}  пустая сцена примера`);
  }

  /* Спрайт на месте — иначе <use href="#i-..."> в примерах пуст. */
  if (/<symbol id="i-/.test(html)) withSprite++;
  else if (st.length) problems.push(`${route}  нет спрайта, а примеры есть`);

  /* Ссылки .md не переписались в маршруты. */
  for (const m of html.matchAll(/href="([^"]*\.md(?:#[^"]*)?)"/g)) {
    problems.push(`${route}  ссылка осталась файловой: ${m[1]}`);
  }

  /* Внутренние ссылки ведут в существующий маршрут. */
  for (const m of html.matchAll(/href="(\/[^"#?]*)"/g)) {
    const t = m[1].endsWith('/') ? m[1] : m[1] + '/';
    if (routes.has(t) || routes.has(m[1])) continue;
    if (/\.(css|js|svg|png|ico|xml|txt|json|woff2?)$/.test(m[1])) continue;
    if (m[1].startsWith('/_astro/') || m[1].startsWith('/pagefind/')) continue;
    problems.push(`${route}  ссылка в никуда: ${m[1]}`);
  }
}

console.log(`страниц: ${pages.length}  ·  живых примеров: ${stages}  ·  со спрайтом: ${withSprite}`);
console.log(`поиск: ${existsSync(join(DIST, 'pagefind')) ? 'индекс собран' : 'ИНДЕКСА НЕТ'}\n`);

if (problems.length) {
  const uniq = [...new Set(problems)];
  console.log(`── проблемы (${uniq.length}) ──`);
  uniq.slice(0, 40).forEach((x) => console.log('  ' + x));
  if (uniq.length > 40) console.log(`  … и ещё ${uniq.length - 40}`);
  process.exit(1);
}

console.log('· сайт собран целиком: ссылки ведут в маршруты, примеры развёрнуты, спрайт на месте');
