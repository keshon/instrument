---
title: Календарь
group: Отображение данных
status: stable
source: src/data.css
needs-js: Выбор даты и переключение месяцев
---

Статическая сетка дат. Кит отвечает за то, чтобы сетка была сеткой, а
сегодняшний день был помечен **не только цветом**.

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

> **Нужен JS.** Выбор даты, переключение месяцев и клавиатура — на слое
> приложения. Кит даёт сетку и состояния дня.

## Когда использовать

| Используйте | Возьмите другое |
|---|---|
| Обзор месяца, где важна раскладка по неделям | **Ввод одной даты в форме** — `<input type="date">`, он стилизован китом и даёт нативный пикер с клавиатурой |
| Плотность событий по дням видна с одного взгляда | **Диапазон дат в фильтре** — два нативных поля даты |
| Выбор дня внутри уже открытого месяца | **События во времени без сетки** — [лента](./timeline.md) |
| — | **Длительности параллельных процессов** — [дорожки](../../agent/lane.md) |

## Состояния дня

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

## Справочник

### Классы

| Класс | Работа |
|---|---|
| `inst-calendar` | Сетка на семь колонок |
| `inst-calendar-dow` | Шапка: сокращение дня недели |
| `inst-calendar-day` | Ячейка дня |

### Атрибуты

| Атрибут | Значения | Где |
|---|---|---|
| `data-today` | `true` | на `inst-calendar-day` |
| `data-outside` | `true` | там же |
| `aria-selected` | `true` · `false` | там же |
| `aria-disabled` | `true` | там же |

### Токены

`--space-1` · `--space-2` · `--control-h-sm` · `--radius-sm` ·
`--text-xs` · `--text-2xs` · `--text-muted` · `--text-faint` ·
`--surface-hover` · `--surface-selected` · `--accent-text` ·
`--accent-solid` · `--size-marker` · `--weight-medium`

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

## Связанное

[Лента событий](./timeline.md) · [Текстовое поле](../inputs/input.md) ·
[Дорожки прогонов](../../agent/lane.md) · [Таблица](./table.md)

Исходник: `src/data.css` · Почему кит устроен именно так — [конституция](../../../DESIGN.md)
