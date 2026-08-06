/* Скрипт стола примера. Одна работа: сообщить родителю свою высоту.

   Иначе iframe остаётся прямоугольником фиксированной высоты — либо с
   пустотой под коротким примером, либо с внутренней полосой прокрутки у
   длинного. И то и другое читается как поломка. */

(() => {
  const send = () => {
    const h = Math.ceil(document.documentElement.getBoundingClientRect().height);
    parent.postMessage({ demoHeight: h, id: location.pathname }, location.origin);
  };

  new ResizeObserver(send).observe(document.documentElement);
  addEventListener('load', send);
  send();

  /* Родитель может сменить тему уже открытого примера — перезагружать
     кадр ради атрибута незачем. */
  addEventListener('message', (e) => {
    if (e.origin !== location.origin || !e.data || e.data.demoTheme === undefined) return;
    const r = document.documentElement;
    if (e.data.demoTheme) r.dataset.theme = e.data.demoTheme;
    else delete r.dataset.theme;
    if (e.data.demoDensity) r.dataset.density = e.data.demoDensity;
    else delete r.dataset.density;
    send();
  });
})();
