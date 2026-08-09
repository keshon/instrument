---
title: Календарь
group: Отображение данных
layout: component
source: src/data.css
needs-js: Выбор даты и переключение месяцев
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

Статическая сетка дат. Кит отвечает за то, чтобы сетка была сеткой, а
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
  <span class="inst-calendar-day" data-outside="true">28</span>
  <span class="inst-calendar-day" data-outside="true">29</span>
  <span class="inst-calendar-day">30</span>
  <span class="inst-calendar-day">31</span>
  <span class="inst-calendar-day" data-today="true">1</span>
  <span class="inst-calendar-day" aria-selected="true">2</span>
  <span class="inst-calendar-day" aria-disabled="true">3</span>
</div>
```

## Использование

```html
<div class="inst-calendar" role="grid" aria-label="Август 2026">
  <span class="inst-calendar-dow" aria-hidden="true">пн</span>
  …
  <span class="inst-calendar-day" role="gridcell"
        aria-label="2 августа 2026, суббота">2</span>
</div>
```

| Что | Обязательно | Почему |
|---|---|---|
| `role="grid"` на сетке, `role="gridcell"` на дне | да | Без ролей календарь озвучивается как поток из 42 чисел |
| `aria-label` с месяцем и годом на сетке | да | «2» без месяца не является датой |
| `aria-label` с полной датой на каждом дне | да | Одно число вне контекста бесполезно |
| `aria-hidden="true"` на днях недели | да | «пн», «вт» дублируют информацию из полной даты ячейки |

:::warn
**Нужен JS.** Выбор даты, переключение месяцев и клавиатура — на слое
приложения. Кит даёт сетку и состояния дня.
:::

## Когда использовать

| Используйте | Возьмите другое |
|---|---|
| Обзор месяца, где важна раскладка по неделям | **Ввод одной даты в форме** — `<input type="date">`, он стилизован китом и даёт нативный пикер с клавиатурой |
| Плотность событий по дням видна с одного взгляда | **Диапазон дат в фильтре** — два нативных поля даты |
| Выбор дня внутри уже открытого месяца | **События во времени без сетки** — [лента](./timeline.md) |
| — | **Длительности параллельных процессов** — [дорожки](../../agent/lane.md) |

## Состояния

| Атрибут | Что значит |
|---|---|
| нет атрибутов | Обычный день текущего месяца |
| `data-today="true"` | Сегодня. Подчёркивание **и** полужирное начертание |
| `data-outside="true"` | День соседнего месяца, добивающий неделю |
| `aria-selected="true"` | Выбран |
| `aria-disabled="true"` | Недоступен. Прозрачность `0.5`, нажатие снято |

**Сегодня помечено подчёркиванием, а не только цветом.** Это выполнение закона
кита: цвет не имеет права быть единственным носителем состояния. Подчёркивание
рисуется внутренней тенью толщиной `--size-marker` — той же, что несёт полосу
активной вкладки и край текущего пункта навигации.

## Правила

:::do Сегодня помечено начертанием
Подчёркивание и полужирное — два признака помимо цвета. Скринридеру это
сообщается словом в `aria-label`.
:::

:::dont Календарь для ввода одной даты
`<input type="date">` уже стилизован китом и даёт нативный пикер с
клавиатурой и мобильным колесом.
:::

:::do Полная дата в aria-label дня
«2 августа 2026, суббота». Одно число вне контекста ничего не значит.
:::

:::dont Сетка без ролей
Без `role="grid"` и `role="gridcell"` это поток из 42 чисел.
:::

## Доступность

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

## API

```api
```

## Связанное

[Лента событий](./timeline.md) [Текстовое поле](../inputs/input.md)
[Дорожки прогонов](../../agent/lane.md) [Таблица](./table.md)
