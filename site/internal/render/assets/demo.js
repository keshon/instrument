/* Скрипт стола примера. Одна работа: сообщить родителю свою высоту.

   Иначе iframe остаётся прямоугольником фиксированной высоты — либо с
   пустотой под коротким примером, либо с внутренней полосой прокрутки у
   длинного. И то и другое читается как поломка. */

(() => {
  /* Высота документа не покрывает три случая, и каждый вылезал на живых
     страницах:

       · поповер, меню и модалка рисуются в ВЕРХНЕМ СЛОЕ и в поток не
         входят вовсе;
       · тултип позиционирован абсолютно и тоже ничего не добавляет;
       · и то и другое появляется ПОСЛЕ замера — по нажатию или наведению.

     Поэтому меряется нижняя кромка всего, что видно, и меряется заново
     на каждое взаимодействие. */
  const bottom = () => {
    let h = document.documentElement.getBoundingClientRect().height;
    for (const el of document.querySelectorAll('body *')) {
      const r = el.getBoundingClientRect();
      if (r.height > 0 && r.bottom > h) h = r.bottom;
    }
    return Math.ceil(h) + 8;
  };

  let last = 0;
  const send = () => {
    const h = bottom();
    if (h === last) return;
    last = h;
    parent.postMessage({ demoHeight: h, id: location.pathname }, location.origin);
  };

  new ResizeObserver(send).observe(document.documentElement);
  addEventListener('load', send);

  /* Тултип раскрывается по наведению и фокусу, поповер и модалка — по
     нажатию. Ни одно из этих событий не меняет высоту документа, поэтому
     ResizeObserver про них не узнает.

     Замер делается СЕРИЕЙ, а не однократно, и это не перестраховка. К
     моменту, когда обработчик клика отработал, поповер ещё не открыт:
     открытие — это действие по умолчанию, оно случается после всплытия.
     А открывшись, он ещё и доезжает переходом, и первый кадр застаёт его
     прозрачным и без размера. Три замера покрывают обе задержки. */
  const remeasure = () => [0, 160, 420].forEach((t) => setTimeout(send, t));
  for (const ev of ['pointerover', 'focusin', 'click', 'keyup', 'toggle', 'transitionend']) {
    addEventListener(ev, remeasure, true);
  }

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
    last = 0;
    send();
  });
})();
