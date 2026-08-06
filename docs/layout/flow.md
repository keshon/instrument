---
title: Примитивы потока
group: Раскладка
status: stable
source: src/layout.css
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

## Установка

```html
<link rel="stylesheet" href="src/kit.css">
```

## Почему здесь нет `mt-3` и `p-2`

Утилит отступов в ките **нет и не будет**, и это не эстетика.

Шкала кита нарочно разрежена сверху: между соседними шагами разрыв растёт, чтобы
«чуть побольше» не было доступным решением. Набор из двухсот утилит возвращает
его первым классом — и вместе с ним переносит решение о вертикальном ритме из
кита в разметку каждого экрана. Через полгода одинаковые на вид экраны
отличаются на два пикселя в тридцати местах, и починить это можно только
поштучно.

Вместо утилит — три примитива с тремя шагами зазора, названными намерением
(обычный, `--tight`, `--loose`). Плотность перенастраивает все три разом, потому
что зазор приходит из ролей (`--pad-panel`, `--gap-row`, `--gap-inline`), а не
из числа в разметке.

Практическое следствие: **зазор ставит контейнер, а не элемент**. Если между
двумя блоками нужен воздух, они кладутся в стопку — им не дописывается поле.
Отступ, поставленный элементу, складывается с внутренним отступом контейнера
сверху и снизу, но не по бокам, и блок отъезжает от рамки по вертикали вдвое
дальше, чем по горизонтали.

Полное рассуждение — в [конституции](../about/design-principles.md).

## Когда использовать

| Используйте | Возьмите другое |
|---|---|
| Вертикальный ритм между блоками — `inst-stack` | **Ограничение ширины и боковые поля** — [контейнер](./container.md) |
| Горизонтальный ряд, который переносится, — `inst-cluster` | **Две колонки разной важности** — [сплит](./split.md): у него есть порог переноса |
| Карточки в адаптивную сетку — `inst-grid` | **Табличные данные** — [таблица](../components/display/table.md): у сетки нет строк и заголовков колонок |
| Полоса действий у заголовка | **Действия экрана** — `inst-page-actions`, [шапка экрана](./page-header.md); **действия секции** — [секция](./section.md) |

## Стопка

```html preview
<div class="inst-stack inst-stack--tight">
  <div class="inst-card"><div class="inst-card-title">Первый</div></div>
  <div class="inst-card"><div class="inst-card-title">Второй</div></div>
</div>
```

| Класс | Зазор | Когда |
|---|---|---|
| `inst-stack` | `--pad-panel` | Умолчание: блоки экрана, секции, панели |
| `inst-stack--tight` | `--gap-row` | Строки внутри блока: пары ключ-значение, список задач |
| `inst-stack--loose` | `--space-7` | Крупные смысловые разделы длинного экрана |

## Кластер

```html preview
<div class="inst-cluster inst-cluster--loose">
  <span class="inst-badge" data-tone="ok"><span class="inst-dot"></span>готово</span>
  <span class="inst-badge" data-tone="running"><span class="inst-dot"></span>идёт</span>
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
`margin-inline-start: auto`: всё после него уезжает к дальнему краю. Один на
ряд; второй ничего не даст, потому что первый уже забрал остаток.

## Сетка

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
| `inst-grid` | `--col-min`, 260px | Умолчание: карточки, панели |
| `inst-grid--tight` | 180px | Мелкие ячейки: метрики, плитки состояний |
| `inst-grid--wide` | 380px | Крупные блоки с текстом внутри |

Число колонок не задаётся: `repeat(auto-fit, minmax(min(var(--col-min), 100%), 1fr))`
считает его сам. `min(…, 100%)` обязателен — без него колонка шириной 380px не
помещается в контейнер шириной 320px и вылезает наружу вместе с горизонтальной
прокруткой.

Это первый уровень адаптивности, **интринсик**: работает всегда, в том числе
там, где контейнера-предка нет.

## Справочник

### Классы

| Класс | Работа |
|---|---|
| `inst-stack` | Колонка. Зазор `--pad-panel` |
| `inst-stack--tight` `inst-stack--loose` | Шаг зазора |
| `inst-cluster` | Ряд с переносом, выравнивание по центру. Зазор `--gap-inline` |
| `inst-cluster--tight` `inst-cluster--loose` | Шаг зазора |
| `inst-cluster-spacer` | Прижать хвост ряда к дальнему краю |
| `inst-grid` | Адаптивная сетка. Зазор `--pad-panel` |
| `inst-grid--tight` `inst-grid--wide` | Минимальная ширина колонки |

### CSS-переменные

| Переменная | По умолчанию |
|---|---|
| `--col-min` | `260px` |

```css
/* Своя плотность сетки — одна строка */
.my-board { --col-min: 320px; }
```

### Токены

`--pad-panel` · `--gap-row` · `--gap-inline` · `--space-2` · `--space-7` ·
`--col-min`

## Доступность

| | |
|---|---|
| Порядок | Все три примитива сохраняют порядок разметки: визуальный и клавиатурный совпадают. `order` и `row-reverse` в ките не применяются |
| Перенос | Кластер и сетка переносятся сами — при увеличении кегля до 200% содержимое не обрезается и не даёт горизонтальной прокрутки |
| Распорка и разделитель | `inst-cluster-spacer` пуст и в дереве доступности отсутствует: он не сообщает ничего сверх порядка |
| Плотность | Все три зазора приходят из ролей и перенастраиваются `data-density` разом. Компонент, у которого зашито число, ломается здесь первым |
| Семантика | Примитивы — `<div>`. Если группа элементов — список, ставьте `<ul>` и класс на него |

## Связанное

[Контейнер](./container.md) · [Сплит](./split.md) · [Секция](./section.md) ·
[Шапка экрана](./page-header.md) · [Плотность](../foundations/density.md) ·
[Шкала отступов](../foundations/spacing.md)

Исходник: `src/layout.css` · Почему кит устроен именно так — [конституция](../about/design-principles.md)
