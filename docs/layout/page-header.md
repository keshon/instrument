---
title: Шапка экрана
group: Раскладка
layout: component
source: src/layout.css
api:
  - { name: "inst-page-header", kind: "класс", doc: "Контейнер. Ряд с переносом, выравнивание по верху" }
  - { name: "inst-page-header-main", kind: "класс", doc: "Колонка «название + пояснение». Порог переноса 20rem" }
  - { name: "inst-page-title", kind: "класс", doc: "Название экрана. Кегль `--text-xl`" }
  - { name: "inst-page-desc", kind: "класс", doc: "Пояснение. `--text-sm`, `--text-secondary`, мера 68ch" }
  - { name: "inst-page-actions", kind: "класс", doc: "Действия. Прижаты к дальнему краю, переносятся" }
  - { name: "--gap-inline", kind: "токен" }
  - { name: "--pad-panel", kind: "токен" }
  - { name: "--space-2", kind: "токен" }
  - { name: "--space-3", kind: "токен" }
  - { name: "--text-xl", kind: "токен" }
  - { name: "--text-sm", kind: "токен" }
  - { name: "--text-secondary", kind: "токен" }
title-en: "Page header"
group-en: "Layout"
---

Название экрана, его пояснение и главные действия. Одна на экран — ровно как
одно `<h1>`.

```html preview
<div class="inst-page-header">
  <div class="inst-page-header-main">
    <h1 class="inst-page-title">Прогоны</h1>
    <p class="inst-page-desc">Каждая область прокручивается сама: боковая навигация
      не уезжает, пока читают таблицу.</p>
  </div>
  <div class="inst-page-actions">
    <button class="inst-btn inst-btn--sm" type="button">Фильтры</button>
    <button class="inst-btn inst-btn--sm inst-btn--primary" type="button">Запустить</button>
  </div>
</div>
```

## Использование

```html
<div class="inst-page-header">
  <div class="inst-page-header-main">
    <h1 class="inst-page-title">Прогон #4127</h1>
    <p class="inst-page-desc">…</p>
  </div>
  <div class="inst-page-actions">…</div>
</div>
```

Обёртка `inst-page-header-main` обязательна: она держит название и пояснение
одной колонкой с зазором `--space-2` и несёт `flex: 1 1 20rem` — порог, ниже
которого действия переносятся под заголовок, а не сжимают его в столбик из
одного слова.

`min-inline-size: 0` на ней — не украшение: без него длинное название без
пробелов растягивает всю шапку и выталкивает действия за край.

| Что | Почему |
|---|---|
| `<h1>` под `inst-page-title` | Класс даёт кегль, а не уровень. Уровень — работа тега, и скринридер читает именно его |
| `inst-page-header-main` как обёртка | Без неё название и пояснение становятся отдельными детьми ряда и встают в строку |
| `type="button"` у действий | Иначе внутри формы они её отправят |

## Когда использовать

| Используйте | Возьмите другое |
|---|---|
| Название экрана и его главные действия | **Заголовок блока внутри экрана** — [секция](./section.md): она тише и не претендует на `<h1>` |
| Пояснение, что это за экран, — одно-два предложения | **Заголовок области приложения** — `inst-panel-title`, [панель](../components/display/panel.md) |
| Место для 1–3 действий над экраном целиком | **Инструменты над содержимым панели** — [тулбар](../components/navigation/toolbar.md) |
| — | **Путь до экрана** — [крошки](../components/navigation/breadcrumbs.md) в шапке [оболочки](./shell.md), а не здесь |

## Устройство

`inst-page-actions` прижимается к дальнему краю (`margin-inline-start: auto`) и
переносится внутри себя. Держите здесь **одно** `primary`: экран, где всё
primary, не имеет primary вовсе — см. [кнопку](../components/actions/button.md).

Четвёртое и следующие действия уходят в [меню](../components/overlays/menu.md).
Полоса из шести кнопок над экраном не читается ни глазом, ни с клавиатуры.

## Композиции

```html preview
<div class="inst-page-header">
  <div class="inst-page-header-main">
    <h1 class="inst-page-title">Прогон #4127
      <span class="inst-badge" data-tone="running"><span class="inst-dot"></span>идёт</span></h1>
    <p class="inst-page-desc">worldgen · начат 14:02</p>
  </div>
  <div class="inst-page-actions">
    <button class="inst-btn inst-btn--sm inst-btn--danger" type="button">Остановить</button>
  </div>
</div>
```

## Правила

:::do Одно название экрана
`<h1>` на экран. Шапка отвечает на «где я», а не перечисляет разделы.
:::

:::dont Шапка экрана внутри блока
Заголовок блока тише и не претендует на `<h1>` — это [секция](./section.md).
:::

:::do Один-три действия
Главное действие экрана и его ближайшие соседи. Больше — значит, часть из них
принадлежит содержимому, а не экрану.
:::

:::dont Инструменты содержимого в шапке экрана
Режимы и фильтры над данными панели живут в
[тулбаре](../components/navigation/toolbar.md).
:::

## Доступность

| | |
|---|---|
| Уровень заголовка | `<h1>` один на документ. Автоматически он не назначается: ставит автор разметки |
| Мера пояснения | 68ch. Строка длиннее теряется на возврате к началу |
| Перенос | Шапка переносится сама; при увеличении кегля до 200% действия уходят под заголовок, а не за край экрана |
| Контраст | Пояснение — `--text-secondary` (4.5:1). `--text-faint` для текста, который читают, не применяется |
| Печать | Название и пояснение печатаются, `inst-page-actions` — нет: действия на бумаге не действия |
| Порядок | Действия стоят в разметке после заголовка и обходятся после него. `order` не применяется |

## API

```api
```

## Связанное

[Оболочка](./shell.md)
[Секция](./section.md)
[Контейнер](./container.md)
[Кнопка](../components/actions/button.md)
[Крошки](../components/navigation/breadcrumbs.md)
[Вкладки](../components/navigation/tabs.md)
