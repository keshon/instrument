---
title: Строка очереди
group: Агентный слой
layout: component
source: src/components.css
needs-js: Порядок и состояния задач — слой приложения. Клавиатуру и выбор строки делает кит
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

> **Нужен JS.** Кит даёт стили и роли. Бегущий `tabindex`, стрелки и выбор
> строки — на слое приложения. Без них до очереди **не добраться с
> клавиатуры**: `role="listbox"` без `tabindex` не получает фокус.

## Использование

```html
<div role="listbox" aria-label="Очередь агентов">
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

| Что | Обязательно | Почему |
|---|---|---|
| `role="listbox"` снаружи, `role="option"` на строке | да | Очередь — список с выбором. Без ролей это поток из тридцати `div` |
| `data-state` на строке | да | Фаза задачи. Точка в жёлобе наследует тон отсюда |
| Бегущий `tabindex` | да | `0` у выбранной, `-1` у остальных: `Tab` входит в список один раз, дальше стрелки |
| `inst-u-truncate` на заголовке | да | Длинное название иначе распирает строку и выталкивает время за край |


## Когда использовать

Оба показывают работу машины, и путать их нельзя.

| | Строка очереди | [Шаг](./step.md) |
|---|---|---|
| Форма | Плоская, одна строка | Есть тело, раскрывается |
| Что показывает | Задачу в списке | Вызов инструмента с выводом |
| Сколько на экране | Сотни | Десятки |
| Раскрытие | нет | `<details>`, бесплатно от платформы |
| Выбор | да, `aria-selected` | нет |

| Возьмите другое | Когда |
|---|---|
| [Таблица](../components/display/table.md) | Задачи сравнивают по нескольким полям |
| [Лента событий](../components/display/timeline.md) | Важен момент времени, а не состояние |
| [Шаги мастера](../components/navigation/steps.md) | Шаги известны заранее и конечны |

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
отдельности не нужно, кит работает делегированием и видит узлы, пришедшие
позже.

```html
<script type="module" src="src/kit.js"></script>
```

### Что делает кит

Очередь объявлена как `role="listbox"`, и кит выполняет контракт этой роли:
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
| `aria-disabled="true"` | Строка выпадает из обхода стрелками, но остаётся видимой |

### Что остаётся приложению

Порядок задач, их состояния и отмена — данные. Кит рисует состояние, но не
знает, когда оно меняется.

## Сценарии

Одна строка очереди не показывает почти ничего: состояния читаются только
рядом друг с другом.

Что видно только здесь:

- **идущая пульсирует**, а завершённая — нет. Пульсация и есть носитель
  «идёт», тон вторичен;
- **завершённая отступает**, а не гаснет: заголовок уходит в
  `--text-secondary`, но остаётся читаемым. История не спорит с тем, что
  происходит сейчас;
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

## Правила

:::do Очередь — настоящий listbox
`role="listbox"` снаружи, `role="option"` на строке, бегущий `tabindex`. Иначе
это поток из тридцати `div`.
:::

:::dont Строка очереди вместо шага
У задачи есть фаза и время, у [шага](./step.md) — вывод, который раскрывают.
Это разные компоненты.
:::

:::do Фаза точкой и словом
Точка в жёлобе наследует тон от `data-state`, а подпись называет фазу.
:::

:::dont Обрезка заголовка без truncate
Длинное название распирает строку и выталкивает время за край.
:::

## Доступность

| | |
|---|---|
| Роли обязательны | Контейнер `role="listbox"`, строка `role="option"`. Без них `aria-selected` **невалиден**, и строка не озвучится как выбираемая |
| Клавиатура | Бегущий `tabindex`: `0` у активной, `-1` у остальных. `Tab` входит в список один раз, дальше стрелки, `Home`/`End`. Поведение — на приложении |
| Состояние словом | `data-state` красит строку, но не озвучивается. Состояние обязано быть в тексте — в подписи или в [бейдже](../components/display/badge.md) рядом |
| Обрезанный заголовок | `inst-u-truncate` режет визуально, полный текст остаётся в DOM и озвучивается. Добавьте `title` для мыши |
| Уменьшенное движение | Пульсация точки `running` замедляется, а не гаснет: очередь агентов без признака работы бесполезна |
| Выделение и наведение складываются | Подсветка выбранной полупрозрачна, поэтому наведение на неё остаётся видно |

## API

```api
```

Длинный заголовок обрезается через `inst-u-truncate` — иначе он выталкивает
метаданные за край.

## Связанное

[Шаг](./step.md)
[Лог](./log.md)
[Дорожки прогонов](./lane.md)
[Бейдж](../components/display/badge.md)
[Таблица](../components/display/table.md)
