---
title: Тепловая карта
group: Диаграммы
layout: component
source: src/data.css
api:
  - { name: "inst-heatmap", kind: "класс", doc: "Поле. Семь строк, недели идут по горизонтали, свежие справа" }
  - { name: "inst-heatmap-cell", kind: "класс", doc: "Один день. Квадратная метка, внутрь неё ничего не кладут" }
  - { name: "data-level", kind: "атрибут", value: "1 · 2 · 3 · 4", doc: "Интенсивность дня. Без атрибута — пустая земля: день был, и на нём ничего не насчитали" }
  - { name: "data-empty", kind: "атрибут", doc: "Такого дня нет — диапазон начался или кончился посреди недели. Место остаётся, метки нет" }
  - { name: "role", kind: "атрибут", value: "img", doc: "**Обязателен.** Поле — картинка данных, и у неё должно быть имя" }
  - { name: "aria-label", kind: "атрибут", doc: "**Обязателен**, и в нём должно стоять количество **словами**. В самой картинке цвет — единственный носитель" }
  - { name: "--cell-size", kind: "переменная", value: "--size-icon-sm", doc: "Сторона ячейки. Задают на поле, чтобы сделать квадраты крупнее или мельче" }
  - { name: "--size-icon-sm", kind: "токен", value: "14px", doc: "Сторона ячейки по умолчанию. У образца она 14,7" }
  - { name: "--accent-solid", kind: "токен", doc: "Верх шкалы. Ступени ниже подмешиваются к пустой земле" }
  - { name: "--surface-recessed", kind: "токен", doc: "Пустая ячейка и та земля, в которую подмешана каждая ступень" }
  - { name: "--radius-mark", kind: "токен" }
  - { name: "--space-1", kind: "токен" }
---

Сколько чего-то случилось за день, по неделям. Поле вкладов, календарь
нагрузки, число прогонов. Единственное показание здесь, у которого ось —
**интенсивность**, а не категория и не исход.

```html preview
<div class="inst-heatmap" role="img"
     aria-label="Thirteen weeks of runs: busiest in the last fortnight, at 40 a day; nine quiet days in all">
  <span class="inst-heatmap-cell" data-empty></span>
  <span class="inst-heatmap-cell" data-empty></span>
  <span class="inst-heatmap-cell"></span>
  <span class="inst-heatmap-cell" data-level="1"></span>
  <span class="inst-heatmap-cell" data-level="2"></span>
  <span class="inst-heatmap-cell" data-level="1"></span>
  <span class="inst-heatmap-cell" data-level="3"></span>
  <span class="inst-heatmap-cell"></span>
  <span class="inst-heatmap-cell" data-level="2"></span>
  <span class="inst-heatmap-cell" data-level="4"></span>
  <span class="inst-heatmap-cell" data-level="1"></span>
  <span class="inst-heatmap-cell" data-level="3"></span>
  <span class="inst-heatmap-cell" data-level="2"></span>
  <span class="inst-heatmap-cell"></span>
  <span class="inst-heatmap-cell" data-level="1"></span>
  <span class="inst-heatmap-cell" data-level="4"></span>
  <span class="inst-heatmap-cell" data-level="3"></span>
  <span class="inst-heatmap-cell" data-level="2"></span>
  <span class="inst-heatmap-cell" data-level="1"></span>
  <span class="inst-heatmap-cell"></span>
  <span class="inst-heatmap-cell" data-level="3"></span>
  <span class="inst-heatmap-cell" data-level="4"></span>
  <span class="inst-heatmap-cell" data-level="2"></span>
  <span class="inst-heatmap-cell" data-level="1"></span>
  <span class="inst-heatmap-cell" data-level="3"></span>
  <span class="inst-heatmap-cell"></span>
  <span class="inst-heatmap-cell" data-level="2"></span>
  <span class="inst-heatmap-cell" data-level="4"></span>
  <span class="inst-heatmap-cell" data-level="3"></span>
  <span class="inst-heatmap-cell" data-level="1"></span>
  <span class="inst-heatmap-cell" data-level="2"></span>
  <span class="inst-heatmap-cell" data-level="4"></span>
  <span class="inst-heatmap-cell" data-level="4"></span>
  <span class="inst-heatmap-cell" data-level="3"></span>
  <span class="inst-heatmap-cell" data-level="2"></span>
  <span class="inst-heatmap-cell" data-level="1"></span>
  <span class="inst-heatmap-cell" data-level="4"></span>
  <span class="inst-heatmap-cell" data-level="3"></span>
  <span class="inst-heatmap-cell" data-level="2"></span>
  <span class="inst-heatmap-cell" data-level="4"></span>
  <span class="inst-heatmap-cell" data-level="1"></span>
  <span class="inst-heatmap-cell" data-level="3"></span>
  <span class="inst-heatmap-cell" data-level="4"></span>
  <span class="inst-heatmap-cell" data-level="2"></span>
  <span class="inst-heatmap-cell" data-level="3"></span>
  <span class="inst-heatmap-cell" data-level="4"></span>
  <span class="inst-heatmap-cell" data-level="1"></span>
  <span class="inst-heatmap-cell" data-empty></span>
  <span class="inst-heatmap-cell" data-empty></span>
</div>
```

## Контракт

| Что | Обязательно | Почему |
|---|---|---|
| `role="img"` и `aria-label` с количеством **словами** | да | Девяносто пустых `<span>` для скринридера — девяносто пустых `<span>`. Здесь имя и есть всё доступное содержимое: в отличие от бейджа или метки изменения, внутри ячейки нет знака, на который можно опереться |
| Семь строк | да | Сетка объявлена в семь строк, ячейки текут по ним колонками. Другое число — другая картинка, и ей нужен другой компонент |
| Свежие **справа** | да | Поле читается как время, а время в этом письме идёт вправо — то же правило, что у [истории проверок](../../agent/history.md) |
| Уровень из закрытого набора | да | `1` · `2` · `3` · `4`, без атрибута — пустая земля. Пятое значение молча не сделает ничего |
| `data-empty` на дни вне диапазона | нет | Только когда диапазон начинается или кончается посреди недели |

### Доступность

| | |
|---|---|
| Цвет здесь единственный носитель, и это честно | У интенсивности нет формы, на которую можно отступить: вдвое темнее — не значит другой формы. Всё прочее в ките, что красит метку, даёт ей ещё и знак или слово; здесь этого нельзя, и вся тяжесть переходит на доступное имя |
| Имя | В нём должно стоять количество, а не форма. «Тринадцать недель прогонов, плотнее всего последние две недели, до 40 в день» — имя; «тепловая карта активности» — нет |
| Принудительные цвета | Правила нет, и намеренно. В этом режиме остаётся одна заливка, и сетка одинаковых квадратов продолжала бы утверждать то, чего уже не значит, — то же решение, что у [истории проверок](../../agent/history.md) |
| Роль | `img`, не `table`. Заголовков нет, по ячейкам ничего не читают: это одна картинка с одним именем |
| Контраст | Ступени — метки, а не текст, и порога 4,5:1 на них нет. Важно, чтобы соседние ступени различались, и подмешивание это гарантирует по построению: каждая — четверть пути от земли до акцента |

## Варианты

Четыре ступени **подмешаны**, а не объявлены:

```css
background: color-mix(in oklab, var(--accent-solid) 25%, var(--surface-recessed));
```

`--chart-1..6` для этого не годятся и не были взяты. Они категориальные —
подобраны так, чтобы ряды **различались**, — а шкале нужно обратное: ступени,
про которые сразу видно, что это один цвет разной силы.

Объявить пять новых токенов означало бы пять значений на пятнадцать ячеек
масштаба и плотности и на четыре акцента, и перенастраивать их при каждом
сдвиге акцента. Подмешивание от `--accent-solid` к пустой земле следует за
всеми четырьмя акцентами и всеми пятью темами даром, а пустая ячейка берёт
`--surface-recessed` — ту же заливку, на которой стоит
[плитка метрики](../display/metric.md).

| Ступень | Заливка |
|---|---|
| без атрибута | `--surface-recessed` |
| `data-level="1"` | 25% акцента |
| `data-level="2"` | 50% акцента |
| `data-level="3"` | 75% акцента |
| `data-level="4"` | `--accent-solid` |

**Четыре ступени и не больше.** Пять пробовали на экране: средние две
перестают различаться на ячейке в 14px. Тепловую карту читают площадью, а
площади нужны грубые ступени.

## Размеры

Ячейка по умолчанию — `--size-icon-sm`, 14px против замеренных у образца 14,7,
— и едет по лестнице кегля, так что поле растёт вместе с интерфейсом, а не
остаётся неподвижным квадратом, пока всё вокруг движется.

```html
<div class="inst-heatmap" style="--cell-size: 10px" role="img" aria-label="…">
```

Год — это пятьдесят две колонки, и в узкую панель они не встанут. Поле
прокручивается со скрытой полосой — тем же приёмом, что
[полоса вкладок](../navigation/tabs.md) и
[строка состояния](../../layout/statusbar.md). Оно не переносится: перенесённая
тепловая карта — это две картинки.

## API

```api
```

## Связанное

```related
```
