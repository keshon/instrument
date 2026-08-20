---
title: Календарь
group: Отображение данных
layout: component
shape: 2
source: src/data.css
js: Выбор даты и переключение месяцев — слой приложения
api:
  - { name: "inst-calendar", kind: "класс", doc: "Сетка на семь колонок" }
  - { name: "inst-calendar-dow", kind: "класс", doc: "Шапка: сокращение дня недели" }
  - { name: "inst-calendar-day", kind: "класс", doc: "Ячейка дня" }
  - { name: "data-today", kind: "атрибут", value: "true", doc: "на `inst-calendar-day`" }
  - { name: "data-outside", kind: "атрибут", value: "true", doc: "там же" }
  - { name: "aria-selected", kind: "атрибут", value: "true · false", doc: "там же" }
  - { name: "aria-disabled", kind: "атрибут", value: "true", doc: "там же" }
  - { name: "--space-1", kind: "токен" }
  - { name: "--space-2", kind: "токен" }
  - { name: "--control-h-sm", kind: "токен" }
  - { name: "--radius-sm", kind: "токен" }
  - { name: "--text-xs", kind: "токен" }
  - { name: "--text-2xs", kind: "токен" }
  - { name: "--text-muted", kind: "токен" }
  - { name: "--surface-hover", kind: "токен" }
  - { name: "--surface-selected", kind: "токен" }
  - { name: "--accent-text", kind: "токен" }
  - { name: "--accent-solid", kind: "токен" }
  - { name: "--size-marker", kind: "токен" }
  - { name: "--weight-medium", kind: "токен" }
title-en: "Calendar"
group-en: "Data display"
---

Статическая сетка дат. Библиотека отвечает за то, чтобы сетка была сеткой, а
сегодняшний день был помечен не только цветом.

```html preview
<div class="inst-calendar" role="grid" aria-label="Август 2026">
  <span class="inst-calendar-dow" aria-hidden="true">пн</span>
  <span class="inst-calendar-dow" aria-hidden="true">вт</span>
  <span class="inst-calendar-dow" aria-hidden="true">ср</span>
  <span class="inst-calendar-dow" aria-hidden="true">чт</span>
  <span class="inst-calendar-dow" aria-hidden="true">пт</span>
  <span class="inst-calendar-dow" aria-hidden="true">сб</span>
  <span class="inst-calendar-dow" aria-hidden="true">вс</span>
  <span class="inst-calendar-day" role="gridcell" data-outside="true">27</span>
  <span class="inst-calendar-day" role="gridcell" data-outside="true">28</span>
  <span class="inst-calendar-day" role="gridcell" data-outside="true">29</span>
  <span class="inst-calendar-day" role="gridcell" data-outside="true">30</span>
  <span class="inst-calendar-day" role="gridcell" data-outside="true">31</span>
  <span class="inst-calendar-day" role="gridcell">1</span>
  <span class="inst-calendar-day" role="gridcell">2</span>
  <span class="inst-calendar-day" role="gridcell">3</span>
  <span class="inst-calendar-day" role="gridcell">4</span>
  <span class="inst-calendar-day" role="gridcell">5</span>
  <span class="inst-calendar-day" role="gridcell">6</span>
  <span class="inst-calendar-day" role="gridcell">7</span>
  <span class="inst-calendar-day" role="gridcell">8</span>
  <span class="inst-calendar-day" role="gridcell">9</span>
  <span class="inst-calendar-day" role="gridcell">10</span>
  <span class="inst-calendar-day" role="gridcell">11</span>
  <span class="inst-calendar-day" role="gridcell" data-today="true">12</span>
  <span class="inst-calendar-day" role="gridcell">13</span>
  <span class="inst-calendar-day" role="gridcell">14</span>
  <span class="inst-calendar-day" role="gridcell">15</span>
  <span class="inst-calendar-day" role="gridcell">16</span>
  <span class="inst-calendar-day" role="gridcell">17</span>
  <span class="inst-calendar-day" role="gridcell" aria-selected="true">18</span>
  <span class="inst-calendar-day" role="gridcell">19</span>
  <span class="inst-calendar-day" role="gridcell">20</span>
  <span class="inst-calendar-day" role="gridcell">21</span>
  <span class="inst-calendar-day" role="gridcell" aria-disabled="true">22</span>
  <span class="inst-calendar-day" role="gridcell" aria-disabled="true">23</span>
  <span class="inst-calendar-day" role="gridcell">24</span>
  <span class="inst-calendar-day" role="gridcell">25</span>
  <span class="inst-calendar-day" role="gridcell">26</span>
  <span class="inst-calendar-day" role="gridcell">27</span>
  <span class="inst-calendar-day" role="gridcell">28</span>
  <span class="inst-calendar-day" role="gridcell" aria-disabled="true">29</span>
  <span class="inst-calendar-day" role="gridcell" aria-disabled="true">30</span>
  <span class="inst-calendar-day" role="gridcell">31</span>
  <span class="inst-calendar-day" role="gridcell" data-outside="true">1</span>
  <span class="inst-calendar-day" role="gridcell" data-outside="true">2</span>
  <span class="inst-calendar-day" role="gridcell" data-outside="true">3</span>
  <span class="inst-calendar-day" role="gridcell" data-outside="true">4</span>
  <span class="inst-calendar-day" role="gridcell" data-outside="true">5</span>
  <span class="inst-calendar-day" role="gridcell" data-outside="true">6</span>
</div>
```

## Контракт

| Что | Обязательно | Почему |
|---|---|---|
| `role="grid"` на сетке, `role="gridcell"` на дне | да | Без ролей календарь озвучивается как поток из 42 чисел |
| `aria-label` с месяцем и годом на сетке | да | «2» без месяца не является датой |
| `aria-label` с полной датой на каждом дне | да | Одно число вне контекста бесполезно |
| `aria-hidden="true"` на днях недели | да | «пн», «вт» дублируют информацию из полной даты ячейки |

### Доступность

| | |
|---|---|
| Роли | `role="grid"` на сетке, `role="gridcell"` на дне, `role="row"` на неделе. Без них календарь озвучивается как поток из 42 чисел |
| Имя | `aria-label` с месяцем и годом на сетке. «2» без месяца не является датой |
| Дни недели | `aria-hidden="true"`: сокращения «пн», «вт» дублируют информацию, которую скринридер берёт из полной даты ячейки |
| Полная дата | Каждый день несёт `aria-label` с полной датой («2 августа 2026, суббота»). Одно число вне контекста бесполезно |
| Клавиатура | Стрелки по дням, `PageUp`/`PageDown` по месяцам, `Home`/`End` по неделе, бегущий `tabindex`. Поведение — на приложении |
| Сегодня | Не только цвет: подчёркивание плюс начертание. Скринридеру это надо сообщить словом в `aria-label` |
| Цель нажатия | Высота — `--control-h-sm`, то есть 26px. В плотности `compact` — 22px, **ниже 24px по WCAG 2.5.8** |
| Табличные цифры | Включены, поэтому колонки чисел стоят ровно независимо от разрядности |

## Состояния

| Атрибут | Что значит |
|---|---|
| нет атрибутов | Обычный день текущего месяца |
| `data-today="true"` | Сегодня. Подчёркивание **и** полужирное начертание |
| `data-outside="true"` | День соседнего месяца, добивающий неделю |
| `aria-selected="true"` | Выбран |
| `aria-disabled="true"` | Недоступен. Прозрачность `0.5`, нажатие снято |

**Сегодня помечено подчёркиванием, а не только цветом.** Это выполнение закона
библиотеки: цвет не имеет права быть единственным носителем состояния. Подчёркивание
рисуется внутренней тенью толщиной `--size-marker` — той же, что несёт полосу
активной вкладки и край текущего пункта навигации.

## JS

Подключите модуль один раз на страницу — инициализировать компоненты по
отдельности не нужно, `instrument.js` работает делегированием и видит узлы, пришедшие
позже.

```html
<script type="module" src="instrument.js"></script>
```

### Что делает `instrument.js`

Ничего. Календарь — это сетка дат, а даты библиотека не считает: месяц, високосный
год, первый день недели и локаль — данные, а не оформление. Библиотека рисует сетку и
состояния дня, разметку строит приложение.

### Что должно сделать приложение

| Что | Что для этого есть в библиотеке |
|---|---|
| Построить сетку месяца | Классы дня и `aria-label` с полной датой |
| Отметить выбранное | `aria-selected="true"` на дне |
| Переключить месяц | Кнопки шапки — обычные `<button>` |

```js
grid.addEventListener('click', (e) => {
  const day = e.target.closest('[data-date]');
  if (!day || day.getAttribute('aria-disabled') === 'true') return;

  for (const d of grid.querySelectorAll('[aria-selected="true"]')) {
    d.setAttribute('aria-selected', 'false');
  }
  day.setAttribute('aria-selected', 'true');
  onPick(day.dataset.date);
});
```

Дата хранится в `data-date` в формате ISO, а не читается из подписи: подпись
зависит от локали, и парсить её обратно — способ получить разъезд на первом же
языке.

## Правила

:::do Сегодня помечено начертанием
Подчёркивание и полужирное — два признака помимо цвета. Скринридеру это
сообщается словом в `aria-label`.
:::

:::dont Календарь для ввода одной даты
`<input type="date">` уже стилизован библиотекой и даёт нативный пикер с
клавиатурой и мобильным колесом.
:::

:::do Полная дата в aria-label дня
«2 августа 2026, суббота». Одно число вне контекста ничего не значит.
:::

:::dont Сетка без ролей
Без `role="grid"` и `role="gridcell"` это поток из 42 чисел.
:::

## API

```api
```

## Связанное

```related
```
