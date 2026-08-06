---
title: Таблица
group: Отображение данных
status: stable
source: src/components.css
needs-js: Сама сортировка, «выбрать всё» и действия строки
---

Сравнимые записи в колонках. Настоящая `<table>` — с семантикой заголовков,
навигацией скринридера по ячейкам и печатью.

```html preview
<table class="inst-table">
  <thead>
    <tr><th>Прогон</th><th>Модель</th><th class="inst-num">Время</th></tr>
  </thead>
  <tbody>
    <tr><td>worldgen-01</td><td>opus</td><td class="inst-num">2,1 с</td></tr>
    <tr><td>biomes-04</td><td>sonnet</td><td class="inst-num">8,4 с</td></tr>
  </tbody>
</table>
```

> **Нужен JS.** Кит даёт стили, разметку и роли. Сама сортировка, «выбрать
> всё» и действия строки — на слое приложения.

## Когда использовать

| Используйте | Возьмите другое |
|---|---|
| Записи, которые **сравнивают** по одним и тем же полям | **Разнородные объекты** — [карточки](./card.md): в таблице половина ячеек будет пустой |
| Больше трёх полей на запись | **Пара «имя — значение»** — [список свойств](./kv.md) |
| Нужны сортировка, выбор строк, итоги | **Строки с состоянием и прогрессом** — [очередь](../../agent/task.md): она плоская и не имеет колонок |
| Данные уйдут в печать | **Десятки тысяч строк** — нужна виртуализация, её в ките нет |

## Числовые колонки

```html
<th class="inst-num">Время</th>
<td class="inst-num">2,1 с</td>
```

`inst-num` выравнивает по концу строки. Табличные цифры включены по умолчанию
для всего кита, поэтому разряды встают друг под друга без дополнительных
усилий.

## Сортировка

```html
<th aria-sort="ascending">
  <button class="inst-sort" type="button">Время</button>
</th>
```

| Значение `aria-sort` | Что показывает |
|---|---|
| нет атрибута | Колонка не сортирована |
| `ascending` | По возрастанию, значок вверх |
| `descending` | По убыванию, значок вниз |

Направление несёт **значок**, а не только позиция строк: `aria-sort`
озвучивает то же самое. Заголовок-сортировка — настоящая `<button>` внутри
`<th>`, иначе до неё не добраться с клавиатуры.

У числовой колонки значок идёт перед подписью, иначе он отрывается от неё на
всю ширину ячейки.

## Выбор строк

```html
<tr aria-selected="true">
  <td class="inst-col-select"><label class="inst-checkbox"><input type="checkbox" checked></label></td>
  <td>worldgen-01</td>
</tr>
```

Состояние живёт на `aria-selected` у `<tr>`. Подсветка полупрозрачна, поэтому
состояния складываются: наведение на выбранную строку остаётся видно.

## Действия строки

```html
<td class="inst-col-actions">
  <div class="inst-row-actions">
    <button class="inst-btn inst-btn--sm inst-btn--icon" type="button" aria-label="Обновить">
      <svg class="inst-icon" aria-hidden="true"><use href="#i-refresh"/></svg>
    </button>
  </div>
</td>
```

Действия появляются на наведении — **и всегда доступны с клавиатуры**: они
проявляются также по `:focus-within` и на выбранной строке. Действие, видимое
только под курсором, для клавиатуры не существует.

## Липкие шапка и колонка

```html
<table class="inst-table inst-table--sticky inst-table--sticky-col">
```

| Класс | Что липнет |
|---|---|
| `inst-table--sticky` | Шапка при вертикальной прокрутке |
| `inst-table--sticky-col` | Первая колонка при горизонтальной |

Липкость работает относительно **прокручиваемого предка** — обычно
`inst-panel-body`. Таблица в неприкреплённом потоке ничего не приклеит.

## Итоги

```html
<tfoot>
  <tr><td>Всего</td><td></td><td class="inst-num">10,5 с</td></tr>
</tfoot>
```

`<tfoot>` идёт в разметке **после** `<tbody>` и отбивается сверху. Строка
итогов — не последняя строка данных, и путать их нельзя.

## Собранный экран

Таблица показывает себя только целиком: сортировка, выбор строк, действия и
итоги по отдельности выглядят как четыре независимых класса, а вместе — как
один инструмент.

Что проверить руками:

- **наведите на строку** — появятся действия. Теперь пройдите по строкам
  `Tab`: действия появятся снова. Это разные события (`:hover` и
  `:focus-within`), и второе существует ровно затем, чтобы действие,
  видимое только под курсором, не пропадало для клавиатуры;
- **выделите строку** и наведите на неё — подсветки складываются, потому
  что выделение полупрозрачно;
- **прокрутите тело панели** — шапка останется на месте.

```html preview context
<div class="inst-panel">
  <div class="inst-panel-header">
    <span class="inst-panel-title">Очередь агентов</span>
    <span class="inst-panel-actions"><span class="inst-badge">7 из 12</span></span>
  </div>
  <div class="inst-panel-body inst-panel-body--list">
    <div class="tree-none" role="listbox" aria-label="Очередь агентов" id="queue">
      <div class="inst-task" role="option" data-state="done" aria-selected="false" tabindex="-1">
        <span class="inst-task-gutter"><span class="inst-dot"></span></span>
        <span class="inst-task-main"><span class="inst-task-title inst-u-truncate">Разбор карты высот</span>
          <span class="inst-task-sub">terrain/heightmap.ts · 4 файла</span></span>
        <span class="inst-task-meta">2,1 с</span>
      </div>
      <div class="inst-task" role="option" data-state="done" aria-selected="false" tabindex="-1">
        <span class="inst-task-gutter"><span class="inst-dot"></span></span>
        <span class="inst-task-main"><span class="inst-task-title inst-u-truncate">Генерация биомов</span>
          <span class="inst-task-sub">world/biomes.ts · 11 файлов</span></span>
        <span class="inst-task-meta">8,4 с</span>
      </div>
      <div class="inst-task" role="option" data-state="running" aria-selected="true" tabindex="0">
        <span class="inst-task-gutter"><span class="inst-dot"></span></span>
        <span class="inst-task-main"><span class="inst-task-title inst-u-truncate">Расстановка ресурсов</span>
          <span class="inst-task-sub">идёт третий проход<span class="inst-caret"></span></span></span>
        <span class="inst-task-meta">14,0 с</span>
      </div>
      <div class="inst-task" role="option" data-state="failed" aria-selected="false" tabindex="-1">
        <span class="inst-task-gutter"><span class="inst-dot"></span></span>
        <span class="inst-task-main"><span class="inst-task-title inst-u-truncate">Валидация путей</span>
          <span class="inst-task-sub">3 узла недостижимы</span></span>
        <span class="inst-task-meta">1,2 с</span>
      </div>
      <div class="inst-task" role="option" data-state="queued" aria-selected="false" tabindex="-1">
        <span class="inst-task-gutter"><span class="inst-dot"></span></span>
        <span class="inst-task-main"><span class="inst-task-title inst-u-truncate">Запекание навмеша</span>
          <span class="inst-task-sub">ждёт валидацию</span></span>
        <span class="inst-task-meta">—</span>
      </div>
      <div class="inst-task" role="option" data-state="skipped" aria-selected="false" tabindex="-1">
        <span class="inst-task-gutter"><span class="inst-dot"></span></span>
        <span class="inst-task-main"><span class="inst-task-title inst-u-truncate">Экспорт превью</span>
          <span class="inst-task-sub">пропущено флагом</span></span>
        <span class="inst-task-meta">—</span>
      </div>
    </div>
  </div>
</div>
```

## Справочник

### Классы

| Класс | Работа |
|---|---|
| `inst-table` | Базовый, на `<table>` |
| `inst-table--sticky` | Липкая шапка |
| `inst-table--sticky-col` | Липкая первая колонка |
| `inst-num` | Числовая ячейка: выравнивание по концу |
| `inst-sort` | Кнопка-заголовок сортировки |
| `inst-col-select` | Колонка чекбоксов, ширина по контролу |
| `inst-col-actions` | Колонка действий, ширина по содержимому |
| `inst-row-actions` | Обёртка действий: проявляется на наведении и фокусе |

### Атрибуты

| Атрибут | Значения | Где |
|---|---|---|
| `aria-sort` | `ascending` · `descending` | на `<th>` |
| `aria-selected` | `true` · `false` | на `<tr>` |

### CSS-переменные

| Переменная | По умолчанию |
|---|---|
| `--table-head-bg` | `--surface-raised` — фон липкой шапки. Прозрачный фон пропустил бы под неё строки |

### Токены

`--text-sm` · `--pad-cell-x` · `--row-pad-y` · `--space-3` · `--border` ·
`--border-subtle` · `--hairline` · `--text-muted` · `--surface-selected` ·
`--surface-hover` · `--control-h-md` · `--z-sticky`

## Доступность

| | |
|---|---|
| Заголовки | `<th>` в `<thead>`. Скринридер озвучивает заголовок при переходе по ячейкам — без него таблица превращается в набор чисел |
| Сортировка | `aria-sort` на `<th>` плюс `<button>` внутри. Атрибут без кнопки недостижим, кнопка без атрибута не озвучивает направление |
| Выбор | `aria-selected` на `<tr>`. Чекбокс несёт своё имя, иначе колонка озвучится как «отметка, отметка, отметка» |
| Действия строки | Проявляются по `:focus-within`, а не только по `:hover`. Иначе они существуют лишь для мыши |
| Подпись таблицы | `<caption>` или `aria-label`. Две таблицы на экране без имён неразличимы на слух |
| Не только цвет | Выбранная строка несёт `aria-selected`, а не одну лишь подсветку |
| Горизонтальная прокрутка | Область прокрутки должна быть достижима с клавиатуры |

## Связанное

[Карточка](./card.md) · [Список свойств](./kv.md) · [Панель](./panel.md) ·
[Очередь](../../agent/task.md) · [Пагинация](../navigation/pagination.md) ·
[Пустое состояние](../feedback/empty.md)

Исходник: `src/components.css` · Почему кит устроен именно так — [конституция](../../about/design-principles.md)
