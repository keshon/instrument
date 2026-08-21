---
title: Примитивы потока
group: Раскладка
layout: component
source: src/layout.css
api:
  - { name: "inst-stack", kind: "класс", doc: "Колонка. Зазор `--pad-panel`" }
  - { name: "inst-cluster", kind: "класс", doc: "Ряд с переносом, выравнивание по центру. Зазор `--gap-inline`" }
  - { name: "inst-cluster-spacer", kind: "класс", doc: "Прижать хвост ряда к дальнему краю" }
  - { name: "inst-grid", kind: "класс", doc: "Адаптивная сетка. Зазор `--pad-panel`" }
  - { name: "inst-stack--tight", kind: "модификатор", doc: "Шаг зазора" }
  - { name: "inst-stack--loose", kind: "модификатор", doc: "Шаг зазора" }
  - { name: "inst-cluster--tight", kind: "модификатор", doc: "Шаг зазора" }
  - { name: "inst-cluster--loose", kind: "модификатор", doc: "Шаг зазора" }
  - { name: "inst-grid--tight", kind: "модификатор", doc: "Минимальная ширина колонки" }
  - { name: "inst-grid--wide", kind: "модификатор", doc: "Минимальная ширина колонки" }
  - { name: "--col-min", kind: "переменная", value: "260px" }
  - { name: "--flow-self", kind: "переменная", doc: "Канал колонки: чем она говорит прямым детям, что те меряются содержимым. Незаследуемый, поэтому не протекает вглубь" }
  - { name: "--pad-panel", kind: "токен" }
  - { name: "--gap-row", kind: "токен" }
  - { name: "--gap-inline", kind: "токен" }
  - { name: "--space-2", kind: "токен" }
  - { name: "--gap-section", kind: "токен" }
title-en: "Flow primitives"
group-en: "Layout"
---

Три способа расставить элементы: стопкой, рядом, сеткой. У каждого три шага
зазора, названные **намерением, а не числом**.

```html preview
<div class="inst-stack">
  <div class="inst-cluster">
    <button class="inst-btn inst-btn--sm" type="button">Фильтры</button>
    <button class="inst-btn inst-btn--sm" type="button">Период</button>
    <span class="inst-cluster-spacer"></span>
    <button class="inst-btn inst-btn--sm inst-btn--primary" type="button">Запустить</button>
  </div>
  <div class="inst-grid inst-grid--tight">
    <div class="inst-metric"><div class="inst-metric-label">В работе</div><div class="inst-metric-value">7</div></div>
    <div class="inst-metric"><div class="inst-metric-label">В очереди</div><div class="inst-metric-value">5</div></div>
    <div class="inst-metric"><div class="inst-metric-label">Упало</div><div class="inst-metric-value">1</div></div>
  </div>
</div>
```

## Контракт

| Что | Обязательно | Почему |
|---|---|---|
| Зазор ставит **контейнер**, а не элемент | да | Отступ у элемента складывается с внутренним отступом контейнера сверху и снизу, но не по бокам: блок отъезжает от рамки по вертикали вдвое дальше, чем по горизонтали |
| Шаг зазора — модификатором | да | `--tight` и `--loose` названы намерением. Число в разметке ломается на первой же смене плотности |
| Один `inst-cluster-spacer` на ряд | да | Второй ничего не даст: первый уже забрал остаток |
| `<ul>` вместо `<div>`, если это список | нет, но обычно да | Примитивы семантики не несут |

### Доступность

| | |
|---|---|
| Порядок | Все три примитива сохраняют порядок разметки: визуальный и клавиатурный совпадают. `order` и `row-reverse` в библиотеке не применяются |
| Перенос | Кластер и сетка переносятся сами — при увеличении кегля до 200% содержимое не обрезается и не даёт горизонтальной прокрутки |
| Распорка | `inst-cluster-spacer` пуст и в дереве доступности отсутствует: он не сообщает ничего сверх порядка |
| Плотность | Все три зазора приходят из ролей и перенастраиваются `data-density` разом. Компонент, у которого зашито число, ломается здесь первым |
| Семантика | Примитивы — `<div>`. Если группа элементов — список, ставьте `<ul>` и класс на него |

## Варианты

### Стопка

```html preview
<div class="inst-stack inst-stack--loose">
  <div class="inst-stack inst-stack--tight">
    <div class="inst-card"><div class="inst-card-title">Первый</div></div>
    <div class="inst-card"><div class="inst-card-title">Второй</div></div>
  </div>
  <div class="inst-stack inst-stack--tight">
    <div class="inst-card"><div class="inst-card-title">Третий</div></div>
    <div class="inst-card"><div class="inst-card-title">Четвёртый</div></div>
  </div>
</div>
```

| Класс | Зазор | Когда |
|---|---|---|
| `inst-stack` | `--pad-panel` | Умолчание: блоки экрана, секции, панели |
| `inst-stack--tight` | `--gap-row` | Строки внутри блока: пары ключ-значение, список задач |
| `inst-stack--loose` | `--gap-section` | Крупные смысловые разделы длинного экрана |

Кнопка, бейдж, тег и сегментированный контрол в колонке **не растягиваются**:
кнопка во всю ширину карточки — та же ошибка, что и иконочная кнопка
прямоугольником. Поля и карточки ширину колонки, наоборот, занимают.

Колонка не знает, кто внутри неё: она объявляет намерение своим прямым детям,
а читает его тот, кого оно касается. Компонент, который меряется содержимым,
несёт `align-self: var(--flow-self, center)` в собственном правиле — рядом со
своим `display: inline-flex`. Поэтому «я меряюсь содержимым» остаётся
свойством компонента, а не списком чужих имён в раскладке.

Запасное значение `center` объясняет, почему в ряду ничего не меняется: там
намерение не объявлено, и компонент центрируется, как и должен.

Это работает во **всех колонках библиотеки** — стопке, форме, филдсете, поле,
секции. В своей колонке приложения — нет: `display: flex; flex-direction:
column`, написанный руками, намерения не объявит, и контролы растянутся
заново. Это довод взять примитив, а не писать колонку — либо объявить канал
самому:

```css
.моя-колонка > * { --flow-self: start; }
```

### Кластер

```html preview
<div class="inst-cluster inst-cluster--loose">
  <span class="inst-cluster inst-cluster--tight">
    <span class="inst-badge" data-tone="ok"><span class="inst-dot"></span>готово</span>
    <span class="inst-badge" data-tone="running"><span class="inst-dot"></span>идёт</span>
  </span>
  <span class="inst-cluster-spacer"></span>
  <button class="inst-btn inst-btn--sm" type="button">Ещё</button>
</div>
```

| Класс | Зазор | Когда |
|---|---|---|
| `inst-cluster` | `--gap-inline` | Умолчание: кнопки, бейджи, контролы в ряд |
| `inst-cluster--tight` | `--space-2` | Элементы, читающиеся как одна группа: теги, чипы |
| `inst-cluster--loose` | `--pad-panel` | Разные по смыслу группы в одной строке |

Кластер **переносится всегда** (`flex-wrap: wrap`) и выравнивает детей по
центру поперечной оси. `inst-cluster-spacer` — пустой элемент с
`margin-inline-start: auto`: всё после него уезжает к дальнему краю.

### Сетка

```html preview
<div class="inst-grid inst-grid--wide">
  <div class="inst-card"><div class="inst-card-title">Карточка</div>
    <div class="inst-card-sub">Колонки перестраиваются сами: auto-fit, без единого запроса.</div></div>
  <div class="inst-card"><div class="inst-card-title">Карточка</div>
    <div class="inst-card-sub">min() не даёт колонке вылезти за контейнер на узком экране.</div></div>
</div>
```

| Класс | Минимальная колонка | Когда |
|---|---|---|
| `inst-grid` | `--col-min`, 260px на базовом масштабе | Умолчание: карточки, панели |
| `inst-grid--tight` | 180px | Мелкие ячейки: метрики, плитки состояний |
| `inst-grid--wide` | 380px | Крупные блоки с текстом внутри |

Число колонок не задаётся: `repeat(auto-fit, minmax(min(var(--col-min), 100%), 1fr))`
считает его сам. `min(…, 100%)` обязателен — без него колонка шириной 380px не
помещается в контейнер шириной 320px и вылезает наружу вместе с горизонтальной
прокруткой.

Это первый уровень адаптивности, **интринсик**: работает всегда, в том числе
там, где контейнера-предка нет.

```css
/* Своя плотность сетки — одна строка */
.my-board { --col-min: 320px; }
```

## API

```api
```

## Связанное

```related
```
