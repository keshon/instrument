---
title: Легенда
group: Графики
layout: component
source: src/data.css
api:
  - { name: "inst-legend", kind: "класс", doc: "Контейнер-список" }
  - { name: "inst-legend-item", kind: "класс", doc: "Элемент: метка плюс подпись" }
  - { name: "inst-legend-swatch", kind: "класс", doc: "Квадратная метка ряда" }
  - { name: "inst-legend-value", kind: "класс", doc: "Величина ряда" }
  - { name: "--chart-1", kind: "токен", doc: "Тон ряда №1. Порядок — часть контракта" }
  - { name: "--chart-2", kind: "токен", doc: "Тон ряда №2. Порядок — часть контракта" }
  - { name: "--chart-3", kind: "токен", doc: "Тон ряда №3. Порядок — часть контракта" }
  - { name: "--chart-4", kind: "токен", doc: "Тон ряда №4. Порядок — часть контракта" }
  - { name: "--chart-5", kind: "токен", doc: "Тон ряда №5. Порядок — часть контракта" }
  - { name: "--chart-6", kind: "токен", doc: "Тон ряда №6. Порядок — часть контракта" }
  - { name: "--series", kind: "переменная", value: "--accent-mark", doc: "Приёмник тона: образец легенды читает её" }
  - { name: "--size-dot", kind: "токен", doc: "Сторона образца" }
  - { name: "--radius-xs", kind: "токен", doc: "Скругление образца" }
  - { name: "--text-xs", kind: "токен", doc: "Кегль подписи" }
  - { name: "--space-3", kind: "токен", doc: "Зазор между образцом и подписью" }
  - { name: "--gap-inline", kind: "токен", doc: "Зазор между пунктами" }
  - { name: "--weight-medium", kind: "токен", doc: "Начертание значения" }
title-en: "Legend"
group-en: "Charts"
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

## Использование

Легенда — **список**. Ряды перечислимы, значит это `<ul>`: скринридер объявит
«список из четырёх элементов» до того, как начнёт их читать.

```html
<ul class="inst-legend">
  <li class="inst-legend-item">
    <span class="inst-legend-swatch" style="--series:var(--chart-1)"></span>Рендер
    <span class="inst-legend-value">18,2 с</span>
  </li>
</ul>
```

| Что | Обязательно | Почему |
|---|---|---|
| `<ul>` и `<li>` | да | Ряды перечислимы. Список объявляет их количество до чтения |
| `--series` на каждой метке | да | Иначе все метки одного цвета — легенда без легенды |
| Порядок, совпадающий с графиком | да | Номер ряда — контракт, а не порядок сортировки |
| Величина или отличимая подпись | да | Цвет не имеет права быть единственным признаком ряда |

## Когда использовать

| Используйте | Возьмите другое |
|---|---|
| Два и более рядов на одном графике — **не опция, а обязанность** | **Один ряд** — легенда не нужна: цвет ничего не различает. Подпишите [метрику](../display/metric.md) над [спарклайном](./sparkline.md) |
| Список рядов с их итоговыми величинами | **Список пар «имя — значение» без графика** — [список свойств](../display/kv.md): квадратик ряда там сообщал бы о цвете, которого нет |
| Расшифровка категорий, различаемых цветом | **Состояния из словаря библиотеки** — [бейдж](../display/badge.md) с точкой: статус не категория, и цвет у него закреплён |
| Пять-шесть рядов максимум | **Седьмой ряд** — это уже не график, а [таблица](../display/table.md). См. [палитру](./palette.md) |

## Устройство

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

Метка ряда — **квадрат, а не круг**: круг в библиотеке занят [точкой
состояния](../display/badge.md), и круглая метка ряда прочиталась бы как
статус.

## Варианты

Цвет ряда назначается через `--series` инлайном — это **канал данных**: номер
ряда приходит вместе с данными и не может быть известен статическому CSS.
Оформление метки — размер, радиус, положение — целиком в библиотеке.

```html preview
<ul class="inst-legend">
  <li class="inst-legend-item"><span class="inst-legend-swatch" style="--series:var(--chart-1)"></span>Ряд 1</li>
  <li class="inst-legend-item"><span class="inst-legend-swatch" style="--series:var(--chart-2)"></span>Ряд 2</li>
  <li class="inst-legend-item"><span class="inst-legend-swatch" style="--series:var(--chart-3)"></span>Ряд 3</li>
  <li class="inst-legend-item"><span class="inst-legend-swatch" style="--series:var(--chart-4)"></span>Ряд 4</li>
</ul>
```

Порядок — часть контракта: первый ряд всегда `--chart-1`, и **порядок в
легенде обязан совпадать с порядком рядов на графике**. Легенда,
отсортированная по величине, а график — по времени, дают две разные нумерации
одного и того же.

**Один ряд — не категория.** Если ряд один, `--series` не ставится вовсе, и
умолчание — не `--chart-1`, а `--accent-mark`. Категориальный цвет существует,
чтобы **различать** ряды; при одном ряде различать нечего, и первый тон палитры
там означал бы категорию, которой нет.

`--chart-1` вступает в силу тогда, когда появляется ряд №2 — и тогда `--series`
проставляется явно на **всех** рядах, включая первый.

## Композиции

### С величинами

```html preview
<ul class="inst-legend">
  <li class="inst-legend-item"><span class="inst-legend-swatch" style="--series:var(--chart-1)"></span>Рендер <span class="inst-legend-value">18,2 с</span></li>
  <li class="inst-legend-item"><span class="inst-legend-swatch" style="--series:var(--chart-2)"></span>Физика <span class="inst-legend-value">11,5 с</span></li>
</ul>
```

Величина в легенде — **второй носитель**: она даёт точное
число тому, кто не различает тона, и делает легенду читаемой отдельно от
графика. Именно поэтому `inst-legend-value` берёт полный цвет текста, а
подпись — приглушённый: сканируют числа.

### Под графиком в карточке

```html preview context
<div class="inst-card inst-stack inst-stack--tight">
  <div class="inst-card-title">Время по модулям</div>
  <svg class="inst-sparkline" viewBox="0 0 120 32" preserveAspectRatio="none" aria-hidden="true">
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

## Правила

:::do Легенда при двух рядах и больше
Без расшифровки цвет ничего не сообщает, и назначать его было незачем.
:::

:::dont Легенда при одном ряде
Различать нечего: цвет ничего не кодирует. Подпишите [метрику](../display/metric.md)
над графиком.
:::

:::do Порядок легенды = порядок графика
Ряд №1 всегда `--chart-1`. Легенда по величине при графике по времени даёт две
нумерации одного и того же.
:::

:::dont Ряды, различимые только тоном
Подпись словом и величина — второй носитель. Иначе легенда не работает ни при
дальтонизме, ни на чёрно-белой печати.
:::

## Доступность

| | |
|---|---|
| Структура | `<ul>` — количество рядов объявляется до перечисления |
| Цвет не единственный носитель | У ряда есть подпись словом и, как правило, величина. Легенда, где ряды различимы только тоном метки, не работает ни при дальтонизме, ни на чёрно-белой печати |
| Контраст | Метка держит 3:1 против поверхности: это несущий знак, а не декорация. Подписи — 4.5:1 |
| Различимость рядов | Светлота внутри категориальной палитры **разная**: ряды остаются различимы, когда тон исчезает вовсе |
| Кегль | `--text-xs` — 12,5px. Метка размером `--size-dot`, то есть 6px: меньше делать нельзя, она перестанет читаться как цвет |
| Порядок чтения | Метка перед подписью: скринридер читает подпись, а не пустой узел, потому что метка не имеет текста и не объявляется |

## API

```api
```

## Связанное

[Палитра](./palette.md) [Спарклайн](./sparkline.md) [Кольцо](./ring.md)
[Мера](./meter.md) [Бейдж](../display/badge.md)
[Список свойств](../display/kv.md) [Таблица](../display/table.md)
