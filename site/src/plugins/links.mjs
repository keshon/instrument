/* Переписывание ссылок вида ./button.md в маршруты сайта.

   Страницы ссылаются друг на друга относительными путями к .md — так они
   работают при чтении репозитория на GitHub. На сайте тот же путь ведёт в
   404: маршрут у страницы другой (/components/actions/button/), файла .md
   по сети нет.

   Терять одно ради другого не нужно: связи в исходнике остаются файловыми,
   а в сборке становятся маршрутами. Один и тот же текст читается в обоих
   местах — ровно поэтому контент и лежит в ../docs, а не внутри сайта.

   Считается от пути файла, а не от текущего URL: относительный путь в
   markdown отсчитывается от файла, и подменять эту точку отсчёта нельзя. */

import { visit } from 'unist-util-visit';
import { relative, dirname, resolve, sep } from 'node:path';

export function remarkDocLinks({ base }) {
  return (tree, file) => {
    if (!file.path) return;
    const fromDir = dirname(relative(base, file.path));

    visit(tree, 'link', (node) => {
      const m = /^([^#?]+\.md)(#.*)?$/.exec(node.url || '');
      if (!m) return;
      if (/^[a-z]+:/i.test(node.url)) return;          /* внешние — не трогаем */

      const [, target, hash = ''] = m;
      /* resolve по абсолютному пути и обратно: так «..» схлопываются
         правильно, включая выход за пределы своей папки. */
      const abs = resolve('/', fromDir, target);
      const route = abs.split(sep).join('/').replace(/\.md$/, '');

      /* index страницы группы адресуется корнем своей папки. */
      node.url = (route.endsWith('/index') ? route.slice(0, -6) || '/' : route)
        .replace(/\/?$/, '/') + hash;
    });
  };
}
