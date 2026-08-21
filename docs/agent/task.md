---
title: Строка очереди
group: Агентный слой
layout: component
source: src/rows.css
js: Порядок и состояния задач — слой приложения. Клавиатуру и выбор строки делает `instrument.js`
api:
  - { name: "inst-task", kind: "класс", doc: "Строка" }
  - { name: "inst-task-gutter", kind: "класс", doc: "Левый жёлоб под точку" }
  - { name: "inst-task-main", kind: "класс", doc: "Заголовок и подпись" }
  - { name: "inst-task-title", kind: "класс", doc: "Название задачи" }
  - { name: "inst-task-sub", kind: "класс", doc: "Что именно делается" }
  - { name: "inst-task-meta", kind: "класс", doc: "Время или счётчик, прижат к концу" }
  - { name: "data-state", kind: "атрибут", doc: "`queued` `running` `done` `warn` `failed` `skipped`" }
  - { name: "aria-selected", kind: "атрибут", doc: "`true` · `false`" }
  - { name: "--size-gutter", kind: "токен" }
  - { name: "--size-dot", kind: "токен" }
  - { name: "--row-pad-y", kind: "токен" }
  - { name: "--pad-cell-x", kind: "токен" }
  - { name: "--gap-inline", kind: "токен" }
  - { name: "--space-1", kind: "токен" }
  - { name: "--text-sm", kind: "токен" }
  - { name: "--text-xs", kind: "токен" }
  - { name: "--surface-selected", kind: "токен" }
  - { name: "--surface-hover", kind: "токен" }
  - { name: "--tone-ink", kind: "токен" }
title-en: "Queue row"
group-en: "Agent layer"
---

Одна задача в очереди агентов. Строка **плоская**: у неё нет тела, потому что
их на экране сотни, и рамка вокруг каждой превратила бы очередь в сетку.

```html preview
<div role="listbox" aria-label="Очередь агентов">
  <div class="inst-task" role="option" data-state="done" aria-selected="false" tabindex="-1">
    <span class="inst-task-gutter"><span class="inst-dot"></span></span>
    <span class="inst-task-main">
      <span class="inst-task-title inst-u-truncate">Разбор карты высот</span>
      <span class="inst-task-sub">terrain/heightmap.ts · 4 файла</span>
    </span>
    <span class="inst-task-meta">2,1 с</span>
  </div>
  <div class="inst-task" role="option" data-state="running" aria-selected="true" tabindex="0">
    <span class="inst-task-gutter"><span class="inst-dot"></span></span>
    <span class="inst-task-main">
      <span class="inst-task-title inst-u-truncate">Генерация биомов</span>
      <span class="inst-task-sub">world/biomes.ts</span>
    </span>
    <span class="inst-task-meta">8,4 с</span>
  </div>
</div>
```

## Контракт

| Что | Обязательно | Почему |
|---|---|---|
| `role="listbox"` снаружи, `role="option"` на строке | да | Очередь — список с выбором. Без ролей это поток из тридцати `div` |
| `data-state` на строке | да | Фаза задачи. Точка в жёлобе наследует тон отсюда |
| Бегущий `tabindex` | да | `0` у выбранной, `-1` у остальных: `Tab` входит в список один раз, дальше стрелки |
| `inst-u-truncate` на заголовке | да | Длинное название иначе распирает строку и выталкивает время за край |

### Доступность

| | |
|---|---|
| Роли обязательны | Контейнер `role="listbox"`, строка `role="option"`. Без них `aria-selected` **невалиден**, и строка не озвучится как выбираемая |
| Клавиатура | Бегущий `tabindex`: `0` у активной, `-1` у остальных. `Tab` входит в список один раз, дальше стрелки, `Home`/`End`. Поведение — на приложении |
| Состояние словом | `data-state` красит строку, но не озвучивается. Состояние обязано быть в тексте — в подписи или в [бейдже](../components/display/badge.md) рядом |
| Обрезанный заголовок | `inst-u-truncate` режет визуально, полный текст остаётся в DOM и озвучивается. Добавьте `title` для мыши |
| Уменьшенное движение | Пульсация точки `running` замедляется, а не гаснет: очередь агентов без признака работы бесполезна |
| Выделение и наведение складываются | Подсветка выбранной полупрозрачна, поэтому наведение на неё остаётся видно |

## Устройство

```html
<span class="inst-task-gutter"><span class="inst-dot"></span></span>
```

Левая колонка фиксированной ширины (`--size-gutter`). Она держит точку
состояния на одной вертикали во всех строках — без неё заголовки разной длины
разъезжают маркеры, и очередь перестаёт читаться столбцом.

## Состояния

| `data-state` | Что происходит |
|---|---|
| `queued` | Базовое. Оформления нет — задача ждёт |
| `running` | Идёт. Точка **пульсирует**, заголовок берёт акцент |
| `done` | Завершено. Заголовок **отступает** в `--text-secondary` |
| `warn` | С замечаниями |
| `failed` | Упало |
| `skipped` | Пропущено. Прозрачность `0.7` |

**Завершённые отступают, а не гаснут.** История остаётся читаемой, но
перестаёт спорить с тем, что идёт сейчас. Прозрачность применяется только к
`skipped` — и она мягкая: `0.55` роняло подпись под порог читаемости.

## JS

Подключите модуль один раз на страницу — инициализировать компоненты по
отдельности не нужно, `instrument.js` работает делегированием и видит узлы, пришедшие
позже.

```html
<script type="module" src="instrument.js"></script>
```

### Что делает `instrument.js`

Пример в шапке живой: войдите в очередь `Tab` и пройдите строки стрелками — выделение идёт за фокусом.

Очередь объявлена как `role="listbox"`, и `instrument.js` выполняет контракт этой роли:
один `Tab` на весь список, стрелки внутри, `Home` и `End`, перебор по кругу.
Выделение следует за фокусом и переносится щелчком. Строки, пришедшие во время
работы, обходятся наравне с остальными — элементы ищутся в момент нажатия.

### События

`inst:select` всплывает с выбранной строки, `detail` — `{ value }`.

```js
queue.addEventListener('inst:select', (e) => {
  details.show(e.detail.value);
});
```

Отмена оставляет `aria-selected` нетронутым: приложение ведёт выделение само.

### Опции

| Атрибут | Что делает |
|---|---|
| `data-value` | Идентификатор задачи в `detail` вместо текста строки |
| `aria-disabled="true"` | Строка остаётся в обходе стрелками и видимой, но не выбирается — см. [поведение](../foundations/behavior.md) |

### Что остаётся приложению

Порядок задач, их состояния и отмена — данные. Библиотека рисует состояние, но не
знает, когда оно меняется.

## Сценарии

Одна строка очереди не показывает почти ничего: состояния читаются только
рядом друг с другом.

Что видно только здесь:

- **идущая пульсирует**, а завершённая — нет. Пульсация и есть носитель
  «идёт», тон вторичен;
- **завершённая отступает**: заголовок уходит в `--text-secondary`, но
  остаётся читаемым. История не спорит с тем, что происходит сейчас;
- **пропущенная** — единственная, к кому применена прозрачность, и она
  мягкая: `0.55` роняло подпись под порог читаемости;
- **точки стоят на одной вертикали** независимо от длины заголовков. За это
  отвечает жёлоб фиксированной ширины.

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

## API

```api
```

Длинный заголовок обрезается через `inst-u-truncate` — иначе он выталкивает
метаданные за край.

## Связанное

```related
```
