---
title: Прогон
group: Агентный слой
layout: component
shape: 2
source: src/agent.css
js: Живые числа и остановка — слой приложения. Свёртка фаз работает без скрипта
api:
  - { name: "inst-dots", kind: "класс", doc: "Счётная мера: сколько единиц из скольких. Единица — `inst-dot`" }
  - { name: "inst-card-head", kind: "класс", doc: "Шапка карточки: заголовок и значок перехода у дальнего края" }
  - { name: "inst-card--link", kind: "модификатор", doc: "Карточка целиком является ссылкой. Без него `base.css` красит её акцентом и подчёркивает" }
  - { name: "role", kind: "атрибут", value: "progressbar", doc: "На `inst-dots`. Вместе с `aria-valuenow/min/max`" }
  - { name: "data-tone", kind: "атрибут", value: "neutral · running · ok · warn · error", doc: "На каждой точке: `ok` — сделано, `running` — идёт, без атрибута — ещё не начиналось" }
  - { name: "--space-1", kind: "токен", doc: "Зазор между точками" }
  - { name: "--size-dot", kind: "токен", doc: "Сторона точки — от `inst-dot`" }
title-en: "Run"
group-en: "Agent layer"
---

Экран прогона — это не компонент, а **сборка**: шапка с именем и счётчиками,
фазы со свёрткой, таблица участников. Отдельного класса у него нет: всё
нужное уже есть, а свой класс означал бы панель под другим именем.

Новое здесь одно: **счётная мера**.

```html preview context
<div class="inst-panel">
  <div class="inst-panel-header">
    <span class="inst-panel-title">audit-worldbox-1</span>
    <span class="inst-badge" data-tone="running"><span class="inst-dot"></span>идёт</span>
    <span class="inst-panel-actions">
      <button class="inst-btn inst-btn--sm inst-btn--danger" type="button">Остановить</button>
    </span>
  </div>

  <div class="inst-panel-body inst-stack">
    <p class="inst-prose">Враждебный аудит worldbox-1: поиск самообмана по классам A–F.</p>

    <dl class="inst-kv">
      <dt>Идёт</dt><dd>21 с</dd>
      <dt>Агентов</dt><dd>7</dd>
      <dt>Токенов</dt><dd>186 000</dd>
    </dl>

    <div class="inst-section">
      <div class="inst-section-head">
        <span class="inst-section-title">Фазы</span>
      </div>

      <details class="inst-accordion-item" open>
        <summary class="inst-accordion-head">
          Разбор
          <span class="inst-dots" role="progressbar" aria-valuenow="4" aria-valuemin="0" aria-valuemax="7"
                aria-label="Разбор: агентов завершено">
            <span class="inst-dot" data-tone="ok"></span>
            <span class="inst-dot" data-tone="ok"></span>
            <span class="inst-dot" data-tone="ok"></span>
            <span class="inst-dot" data-tone="ok"></span>
            <span class="inst-dot" data-tone="running"></span>
            <span class="inst-dot"></span>
            <span class="inst-dot"></span>
          </span>
        </summary>
        <div class="inst-accordion-body inst-panel-body--flush">
          <table class="inst-table">
            <thead>
              <tr><th>Агент</th><th class="inst-num">Токенов</th><th class="inst-num">Вызовов</th><th class="inst-num">Время</th></tr>
            </thead>
            <tbody>
              <tr><td>разбор:docs-drift</td><td class="inst-num">38 200</td><td class="inst-num">4</td><td class="inst-num">18 с</td></tr>
              <tr><td>разбор:shared-and-chunk</td><td class="inst-num">39 600</td><td class="inst-num">5</td><td class="inst-num">18 с</td></tr>
              <tr><td>разбор:probes-assert</td><td class="inst-num">37 900</td><td class="inst-num">7</td><td class="inst-num">18 с</td></tr>
              <tr><td>разбор:silent-failure</td><td class="inst-num">36 800</td><td class="inst-num">4</td><td class="inst-num">18 с</td></tr>
              <tr data-state="running"><td>разбор:coverage-hole</td><td class="inst-num">—</td><td class="inst-num">—</td><td class="inst-num">—</td></tr>
              <tr><td>разбор:history</td><td class="inst-num">—</td><td class="inst-num">—</td><td class="inst-num">—</td></tr>
              <tr><td>разбор:eyes-only</td><td class="inst-num">—</td><td class="inst-num">—</td><td class="inst-num">—</td></tr>
            </tbody>
          </table>
        </div>
      </details>

      <details class="inst-accordion-item">
        <summary class="inst-accordion-head">
          Опровержение
          <span class="inst-dots" role="progressbar" aria-valuenow="0" aria-valuemin="0" aria-valuemax="3"
                aria-label="Опровержение: агентов завершено">
            <span class="inst-dot"></span>
            <span class="inst-dot"></span>
            <span class="inst-dot"></span>
          </span>
        </summary>
        <div class="inst-accordion-body">
          <div class="inst-empty">
            <span class="inst-empty-title">Ещё не начиналось</span>
            <span class="inst-empty-desc">Фаза стартует, когда разбор закроет все семь агентов.</span>
          </div>
        </div>
      </details>
    </div>
  </div>
</div>
```

## Контракт

Прогон — не отдельный класс, а **сборка**: панель с шапкой, список свойств,
счётная мера и секции с шагами. Своего у него ровно две вещи — счётная мера и
карточка-ссылка.

| Что | Обязательно | Почему |
|---|---|---|
| `role="progressbar"` с тремя значениями на `inst-dots` | да | Иначе «сколько из скольких» существует только как ряд кружков |
| `data-tone` на каждой точке | да | `ok` — сделано, `running` — идёт, без атрибута — ещё не начиналось |
| `inst-card--link` у карточки-ссылки | да | Без него базовый слой красит её акцентом и подчёркивает как обычную ссылку |
| Число рядом со счётной мерой | да | Семь точек глазом не считаются: «3 из 4» словами обязательно |

### Доступность

| | |
|---|---|
| Роль | `role="progressbar"` на `inst-dots` плюс `aria-valuenow/min/max`. Без них «сколько сделано» существует только в числе закрашенных кружков |
| Имя | `aria-label` с тем, что именно считается. «4 из 7» без предмета не является сообщением |
| Точки | Декоративны для скринридера: значение несёт контейнер, а не они. Отдельных подписей им не нужно |
| Не только цвет | Идущая единица **пульсирует**, а не просто окрашена; завершённость озвучивается `aria-valuenow` |
| Уменьшенное движение | Пульсация замедляется, а не гаснет: индикатор занятости должен оставаться видимым |
| Печать | Точки печатаются: это данные, а не индикатор активности |

## Устройство

```html preview
<span class="inst-dots" role="progressbar" aria-valuenow="4" aria-valuemin="0" aria-valuemax="7"
      aria-label="Агентов завершено">
  <span class="inst-dot" data-tone="ok"></span>
  <span class="inst-dot" data-tone="ok"></span>
  <span class="inst-dot" data-tone="ok"></span>
  <span class="inst-dot" data-tone="ok"></span>
  <span class="inst-dot" data-tone="running"></span>
  <span class="inst-dot"></span>
  <span class="inst-dot"></span>
</span>
```

Единиц семь, и каждая либо сделана, либо нет. **Доля здесь — выдуманная
точность:** 4 из 7, нарисованные полосой на 57%, выглядят правдоподобно и
сообщают то, чего система не знает. Это та же ложь, что определённая полоса,
застрявшая на 90%, только тише.

Своей точки компонент не заводит: единица — обычный [`inst-dot`](../components/feedback/states.md).
Он уже умеет тон, пульсацию «идёт» и режим принудительных цветов, и второй
кружок в библиотеке был бы вторым именем для того же.

## Состояния

| Тон | Что значит |
|---|---|
| без атрибута | Ещё не начиналось |
| `data-tone="running"` | Идёт. Точка пульсирует |
| `data-tone="ok"` | Сделано |
| `data-tone="error"` | Упало |
| `data-tone="warn"` | Сделано с замечанием |

Пульсирует **только** идущая, и это та же пульсация, что у строки очереди и у
шага: в библиотеке она означает ровно одно, и заводить ей второй смысл нельзя.

## JS

Подключите модуль один раз на страницу — инициализировать компоненты по
отдельности не нужно, `instrument.js` работает делегированием и видит узлы, пришедшие
позже.

```html
<script type="module" src="instrument.js"></script>
```

### Что делает `instrument.js`

Свёртку фаз не делает никто: она на `<details>` и работает без скрипта.

### Что остаётся приложению

Живые числа и остановка — данные и команда, а не оформление.

```js
source.addEventListener('message', (e) => {
  const { phase, done } = JSON.parse(e.data);
  const dots = run.querySelectorAll('#' + phase + ' .inst-dot');

  dots.forEach((dot, i) => {
    if (i < done) dot.dataset.tone = 'ok';
    else if (i === done) dot.dataset.tone = 'running';
    else delete dot.dataset.tone;
  });
  dots[0].closest('.inst-dots').setAttribute('aria-valuenow', done);
});
```

`aria-valuenow` обновляется вместе с точками, а не «когда-нибудь потом».
Точки — это картинка; для скринридера прогресс существует только в атрибуте,
и забытый атрибут означает фазу, которая для него навсегда осталась в нуле.

## Композиции

Тот же прогон одной строкой — когда их несколько и нужно выбрать. Здесь
счётная мера работает лучше всего: пятнадцать точек читаются одним взглядом,
и видно не только «сколько», но и **что одна упала с замечанием**. Полоса на
20% этого не сообщила бы.

```html preview context
<div class="inst-stack inst-stack--tight">
  <a class="inst-card inst-card--interactive inst-card--link" href="#">
    <span class="inst-card-head">
      <span class="inst-card-title">audit-worldbox-1</span>
      <svg class="inst-icon" aria-hidden="true"><use href="#i-chevron"/></svg>
    </span>
    <span class="inst-card-sub">Workflow · 15 агентов · 7 мин 58 с</span>
    <span class="inst-dots" role="progressbar" aria-valuenow="4" aria-valuemin="0" aria-valuemax="15"
          aria-label="Агентов завершено">
      <span class="inst-dot" data-tone="ok"></span>
      <span class="inst-dot" data-tone="ok"></span>
      <span class="inst-dot" data-tone="warn"></span>
      <span class="inst-dot" data-tone="ok"></span>
      <span class="inst-dot" data-tone="running"></span>
      <span class="inst-dot"></span><span class="inst-dot"></span>
      <span class="inst-dot"></span><span class="inst-dot"></span>
      <span class="inst-dot"></span><span class="inst-dot"></span>
      <span class="inst-dot"></span><span class="inst-dot"></span>
      <span class="inst-dot"></span><span class="inst-dot"></span>
    </span>
  </a>

  <a class="inst-card inst-card--interactive inst-card--link" href="#">
    <span class="inst-card-head">
      <span class="inst-card-title">review-terrain-08</span>
      <svg class="inst-icon" aria-hidden="true"><use href="#i-chevron"/></svg>
    </span>
    <span class="inst-card-sub">Workflow · 3 агента · 41 с</span>
    <span class="inst-dots" role="progressbar" aria-valuenow="3" aria-valuemin="0" aria-valuemax="3"
          aria-label="Агентов завершено">
      <span class="inst-dot" data-tone="ok"></span>
      <span class="inst-dot" data-tone="ok"></span>
      <span class="inst-dot" data-tone="ok"></span>
    </span>
  </a>
</div>
```

Карточка — **ссылка целиком**. Цель нажатия получается во всю строку, и
растить её отдельно не приходится. Шеврон декоративен, имя ссылке даёт
заголовок.

## Сценарии

Ни одного нового класса, кроме `inst-dots`:

| Часть экрана | Чем собрана |
|---|---|
| Рамка и шапка | [Панель](../components/display/panel.md) |
| Состояние прогона | [Бейдж](../components/display/badge.md) с точкой |
| Остановка | [Кнопка](../components/actions/button.md), вариант `--danger` |
| Счётчики | [Список пар](../components/display/kv.md) |
| Фаза со свёрткой | [Аккордеон](../components/feedback/accordion.md) на `<details>` |
| Участники | [Таблица](../components/display/table.md) с `inst-num` |
| Пустая фаза | [Пустое состояние](../components/feedback/empty.md) |

**Это и есть проверка покрытия.** Экран, который не собирается без
инлайнового оформления, означает дыру в библиотеке: правило про живую спецификацию из
[принципов дизайна](../about/design-principles.md) распространяется и сюда.

## Правила

:::do Счётная мера с числом рядом
Семь точек глазом не считаются. «3 из 4» словами обязательно.
:::

:::dont Точки без роли
Ряд кружков без `role="progressbar"` и трёх значений существует только для
зрячих.
:::

:::do Прогон собирается из готовых компонентов
Панель, список свойств, секции, шаги. Своего у прогона ровно две вещи — и это
правильная пропорция.
:::

:::dont Новый класс под каждый экран прогона
Экран, который не собирается из библиотеки, означает дыру в библиотеке, а не потребность в
ещё одном классе.
:::

## API

```api
```

## Связанное

```related
```
