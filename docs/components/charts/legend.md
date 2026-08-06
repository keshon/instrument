---
title: Легенда
group: Графики
status: stable
source: src/data.css
---

Расшифровка цветов графика. **Обязательна при двух и более рядах**: без неё
цвет не сообщает ничего, а график с бессмысленным цветом — это раскрашенная
картинка.

```html preview
<ul class="inst-legend">
  <li class="inst-legend-item"><span class="inst-legend-swatch" style="--series:var(--chart-1)"></span>Рендер <span class="inst-legend-value">18,2 с</span></li>
  <li class="inst-legend-item"><span class="inst-legend-swatch" style="--series:var(--chart-2)"></span>Физика <span class="inst-legend-value">11,5 с</span></li>
  <li class="inst-legend-item"><span class="inst-legend-swatch" style="--series:var(--chart-3)"></span>Ассеты <span class="inst-legend-value">8,1 с</span></li>
  <li class="inst-legend-item"><span class="inst-legend-swatch" style="--series:var(--chart-4)"></span>Аудио <span class="inst-legend-value">4,2 с</span></li>
</ul>
```

## Установка

```html
<link rel="stylesheet" href="src/kit.css">
```

## Когда использовать

| Используйте | Возьмите другое |
|---|---|
| Два и более рядов на одном графике — **не опция, а обязанность** | **Один ряд** — легенда не нужна: цвет ничего не различает. Подпишите [метрику](../display/metric.md) над [спарклайном](./sparkline.md) |
| Список рядов с их итоговыми величинами | **Список пар «имя — значение» без графика** — [список пар](../display/kv.md): квадратик ряда там сообщал бы о цвете, которого нет |
| Расшифровка категорий, различаемых цветом | **Состояния из словаря кита** — [бейдж](../display/badge.md) с точкой: статус не категория, и цвет у него закреплён |
| Пять-шесть рядов максимум | **Седьмой ряд** — это уже не график, а [таблица](../display/table.md). См. [палитру](./palette.md) |

## Разметка

Легенда — **список**, а не строка из `<span>`. Ряды перечислимы, значит это
`<ul>`, и скринридер объявит «список из четырёх элементов» до того, как начнёт
их читать.

```html preview
<ul class="inst-legend">
  <li class="inst-legend-item"><span class="inst-legend-swatch" style="--series:var(--chart-1)"></span>Успешно</li>
  <li class="inst-legend-item"><span class="inst-legend-swatch" style="--series:var(--chart-2)"></span>С замечаниями</li>
  <li class="inst-legend-item"><span class="inst-legend-swatch" style="--series:var(--chart-3)"></span>Упало</li>
</ul>
```

| Узел | Работа |
|---|---|
| `inst-legend` | `<ul>`. Маркеры и отступы списка сняты, элементы переносятся |
| `inst-legend-item` | `<li>`. Метка и подпись в одну строку |
| `inst-legend-swatch` | Метка ряда. Читает `--series` |
| `inst-legend-value` | Величина ряда. Полный цвет текста и средний вес |

Метка ряда — **квадрат, а не круг**: круг в этом ките занят [точкой
состояния](../display/badge.md), и круглая метка ряда прочиталась бы как
статус.

## Цвет ряда

`style="--series:var(--chart-2)"` — инлайновый стиль как **канал данных**:
номер ряда приходит вместе с данными и не может быть известен статическому
CSS. Оформление метки — размер, радиус, положение — целиком в ките.

Порядок — часть контракта: первый ряд всегда `--chart-1`, и **порядок в
легенде обязан совпадать с порядком рядов на графике**. Легенда,
отсортированная по величине, а график — по времени, дают две разные нумерации
одного и того же.

Если рядов один, `--series` не ставится вовсе: умолчание `--chart-1`
достаточно.

## С величинами

```html preview
<ul class="inst-legend">
  <li class="inst-legend-item"><span class="inst-legend-swatch" style="--series:var(--chart-1)"></span>Рендер <span class="inst-legend-value">18,2 с</span></li>
  <li class="inst-legend-item"><span class="inst-legend-swatch" style="--series:var(--chart-2)"></span>Физика <span class="inst-legend-value">11,5 с</span></li>
</ul>
```

Величина в легенде — не украшение, а **второй носитель**: она даёт точное
число тому, кто не различает тона, и делает легенду читаемой отдельно от
графика. Именно поэтому `inst-legend-value` берёт полный цвет текста, а
подпись — приглушённый: сканируют числа.

## Композиции

### Под графиком в карточке

```html
<div class="inst-card inst-stack inst-stack--tight">
  <div class="inst-card-title">Время по модулям</div>
  <svg class="inst-sparkline" viewBox="0 0 120 32" preserveAspectRatio="none" aria-hidden="true" style="--series:var(--chart-1)">
    <path class="inst-sparkline-line" d="M0,26 24,24 48,18 72,13 96,9 120,6"/>
  </svg>
  <ul class="inst-legend">
    <li class="inst-legend-item"><span class="inst-legend-swatch" style="--series:var(--chart-1)"></span>Рендер <span class="inst-legend-value">18,2 с</span></li>
    <li class="inst-legend-item"><span class="inst-legend-swatch" style="--series:var(--chart-2)"></span>Физика <span class="inst-legend-value">11,5 с</span></li>
  </ul>
</div>
```

Легенда стоит **под** графиком: сверху она отодвигает данные ниже линии сгиба
ради того, что читают вторым.

## Справочник

### Классы

| Класс | Работа |
|---|---|
| `inst-legend` | Контейнер-список |
| `inst-legend-item` | Элемент: метка плюс подпись |
| `inst-legend-swatch` | Квадратная метка ряда |
| `inst-legend-value` | Величина ряда |

### Обязательная разметка

| Что | Почему |
|---|---|
| `<ul>` и `<li>` | Ряды перечислимы. Список объявляет их количество до чтения |
| `--series` на каждой метке | Иначе все метки одного цвета — легенда без легенды |
| Порядок, совпадающий с графиком | Номер ряда — контракт, а не порядок сортировки |
| Величина или отличимая подпись | Цвет не имеет права быть единственным признаком ряда |

### CSS-переменные

| Переменная | По умолчанию |
|---|---|
| `--series` | `--chart-1` |

### Токены

`--chart-1` … `--chart-6` · `--size-dot` · `--radius-xs` · `--text-xs` ·
`--space-3` · `--gap-inline` · `--weight-medium`

## Доступность

| | |
|---|---|
| Структура | `<ul>` — количество рядов объявляется до перечисления |
| Цвет не единственный носитель | У ряда есть подпись словом и, как правило, величина. Легенда, где ряды различимы только тоном метки, не работает ни при дальтонизме, ни на чёрно-белой печати |
| Контраст | Метка держит 3:1 против поверхности: это несущий знак, а не украшение. Подписи — 4.5:1 |
| Различимость рядов | Светлота внутри категориальной палитры **разная намеренно**: ряды остаются различимы, когда тон исчезает вовсе |
| Кегль | `--text-xs` — 12px. Метка размером `--size-dot`, то есть 6px: меньше делать нельзя, она перестанет читаться как цвет |
| Порядок чтения | Метка перед подписью: скринридер читает подпись, а не пустой узел, потому что метка не имеет текста и не объявляется |

## Связанное

[Палитра](./palette.md) · [Спарклайн](./sparkline.md) · [Кольцо](./ring.md) ·
[Мера](./meter.md) · [Бейдж](../display/badge.md) ·
[Список пар](../display/kv.md) · [Таблица](../display/table.md)

Исходник: `src/data.css` · Почему кит устроен именно так — [конституция](../../about/design-principles.md)
