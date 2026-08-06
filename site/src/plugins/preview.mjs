/* Живой пример из ограды ```html preview.

   Правило страниц: примеры не сочиняются, а берутся из работающей разметки.
   Живой рендер — вторая половина того же правила: пример, который никто не
   запускал, рано или поздно расходится с китом молча.

   Плагин remark, а не rehype, потому что мета ограды (`preview`) доступна
   именно в mdast: в hast она уже потеряна.

   Перед блоком кода вставляется узел html с той же разметкой. Она попадает
   в документ как есть — то есть рисуется настоящим китом, а не картинкой
   и не копией. Разойтись с кодом под ней она физически не может: это одна
   и та же строка. */

import { visit } from 'unist-util-visit';

export function remarkPreview() {
  return (tree) => {
    visit(tree, 'code', (node, index, parent) => {
      if (!parent || node.lang !== 'html') return;
      if (!/(^|\s)preview(\s|$)/.test(node.meta || '')) return;

      parent.children.splice(index, 0, {
        type: 'html',
        value:
          '<div class="preview">' +
          '<div class="preview-stage" data-preview>' + node.value + '</div>' +
          '</div>',
      });
      return index + 2;   /* пропускаем вставленный узел и сам блок кода */
    });
  };
}
