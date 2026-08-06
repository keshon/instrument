---
title: Навигация
group: Навигация
status: stable
source: src/layout.css
---

Список разделов приложения в боковой колонке. Текущий пункт помечен **не только
цветом**: у его края стоит несущая полоса.

```html preview
<nav class="inst-nav" aria-label="Разделы">
  <div class="inst-nav-group">
    <span class="inst-nav-label">Работа</span>
    <a class="inst-nav-item" href="#" aria-current="page">
      <svg class="inst-icon" aria-hidden="true"><use href="#i-grid"/></svg>Прогоны
      <span class="inst-badge inst-nav-count">7</span></a>
    <a class="inst-nav-item" href="#">
      <svg class="inst-icon" aria-hidden="true"><use href="#i-list"/></svg>Очередь</a>
    <a class="inst-nav-item" href="#">
      <svg class="inst-icon" aria-hidden="true"><use href="#i-chart"/></svg>Отчёты</a>
  </div>
  <div class="inst-nav-group">
    <span class="inst-nav-label">Настройка</span>
    <a class="inst-nav-item" href="#">
      <svg class="inst-icon" aria-hidden="true"><use href="#i-users"/></svg>Доступы</a>
    <a class="inst-nav-item" href="#">
      <svg class="inst-icon" aria-hidden="true"><use href="#i-settings"/></svg>Параметры</a>
  </div>
</nav>
```

## Установка

```html
<link rel="stylesheet" href="src/kit.css">
```

## Когда использовать

| Используйте | Возьмите другое |
|---|---|
| Постоянный список разделов приложения — обычно внутри [оболочки](../../layout/shell.md) | **Переключение видов внутри одного раздела** — [вкладки](./tabs.md): у них другая работа и другое место на экране |
| Плоский список ссылок, сгруппированный подписями | **Раскрывающаяся иерархия объектов** — [дерево](../../agent/tree.md): у пункта навигации нет уровня и раскрытия |
| Пункт со счётчиком незавершённого — `inst-nav-count` | **Список действий по клику** — [меню](../overlays/menu.md): навигация ведёт на адрес, меню запускает действие |
| Путь до текущего экрана — рядом, в шапке | **Только путь, без списка разделов** — [крошки](./breadcrumbs.md) |

## Группы и подписи

`inst-nav-group` собирает пункты в смысловой блок, `inst-nav-label` его
подписывает. Зазор между группами ставится **между соседями**
(`.inst-nav-group + .inst-nav-group`), а не вокруг каждой, — иначе первая
группа отступила бы от верхнего края колонки вдвое дальше, чем от боковых.

Подпись — обычное предложение кеглем `--text-2xs` цветом `--text-muted`. Капса
здесь нет: подпись не имеет права звучать громче пунктов, которые она
подписывает.

## Текущий пункт

```html preview
<nav class="inst-nav" aria-label="Пример текущего пункта">
  <div class="inst-nav-group">
    <a class="inst-nav-item" href="#" aria-current="page">Прогоны</a>
    <a class="inst-nav-item" href="#">Очередь</a>
  </div>
</nav>
```

Метка ставится атрибутом `aria-current="page"`, а не классом. Оформление ловит
селектор `[aria-current]:not([aria-current="false"])`, поэтому строчка
`aria-current="false"` у остальных пунктов допустима и ничего не подсвечивает.

| Признак текущего | Что даёт |
|---|---|
| Полоса у внутреннего края | `box-shadow: inset` шириной `--size-marker`. Носитель, не зависящий от цвета |
| Поверхность | `--surface-selected` |
| Цвет подписи | `--accent-text` |
| Начертание | `--weight-medium`. Веса 600 и 700 в ките нет |

## Счётчик

```html
<a class="inst-nav-item" href="#">Очередь
  <span class="inst-badge inst-nav-count">7</span></a>
```

`inst-nav-count` — модификатор положения: `margin-inline-start: auto` прижимает
[бейдж](../display/badge.md) к дальнему краю пункта. Сам бейдж остаётся
бейджем и читает те же токены.

## В узкой оболочке

Внутри `.inst-shell-aside` ниже 60rem колонка разворачивается в горизонтальную
полосу: `.inst-nav` и `.inst-nav-group` встают в ряд, `.inst-nav-label`
скрывается, подписи пунктов перестают переноситься. Это единственное место, где
навигация меняет форму, и делает это [оболочка](../../layout/shell.md), а не
сама навигация.

## Справочник

### Классы

| Класс | Работа |
|---|---|
| `inst-nav` | Контейнер списка. Колонка с зазором `--space-1` |
| `inst-nav-group` | Смысловой блок пунктов |
| `inst-nav-label` | Подпись группы |
| `inst-nav-item` | Пункт. Ссылка, а не кнопка |
| `inst-nav-count` | Прижать счётчик к дальнему краю пункта |

### Обязательная разметка

| Что | Почему |
|---|---|
| `<nav>` + `aria-label` | На экране больше одной навигации (боковая, крошки, вкладки, страницы) — без имени они неразличимы в списке ориентиров |
| `<a href>` у пункта | Навигация — это адрес. Кнопка ломает средний клик, «открыть в новой вкладке» и копирование ссылки |
| `aria-current="page"` | Единственный источник текущего состояния. Класса для него нет |
| `aria-hidden="true"` у иконки | Имя пункта даёт подпись |

### Токены

`--control-h-md` · `--radius-sm` · `--text-sm` · `--text-2xs` ·
`--weight-medium` · `--space-1` · `--space-3` · `--gap-inline` · `--pad-panel` ·
`--size-marker` · `--surface-hover` · `--surface-selected` · `--accent-solid` ·
`--accent-text` · `--text-secondary` · `--text-muted` · `--dur-1` · `--ease-out`

## Доступность

| | |
|---|---|
| Клавиатура | Всё от нативных ссылок: `Tab` — обход, `Enter` — переход. Кит ничего не перехватывает, бегущий `tabindex` навигации не нужен |
| Ориентир | `<nav aria-label>` — пункт назначения для «перейти к ориентиру» в скринридере |
| Цвет не единственный носитель | Текущий пункт несёт полосу у края и `aria-current`, а не один только оттенок фона |
| Цель нажатия | Высота пункта — `--control-h-md`, 32px. В плотности `compact` она уменьшается: применяйте `compact` к табличным видам, а не к основной навигации |
| Контраст | Подпись `--text-secondary` держит 4.5:1 на приглушённой поверхности колонки в четырёх темах |
| Печать | Боковая колонка не печатается: это хром приложения, а не данные |

## Связанное

[Оболочка](../../layout/shell.md) · [Крошки](./breadcrumbs.md) ·
[Вкладки](./tabs.md) · [Бейдж](../display/badge.md) ·
[Дерево](../../agent/tree.md) · [Иконка](../../foundations/icons.md)

Исходник: `src/layout.css` · Почему кит устроен именно так — [конституция](../../about/design-principles.md)
