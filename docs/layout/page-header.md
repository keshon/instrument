---
title: Шапка экрана
group: Раскладка
status: stable
source: src/layout.css
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

## Установка

```html
<link rel="stylesheet" href="src/kit.css">
```

## Когда использовать

| Используйте | Возьмите другое |
|---|---|
| Название экрана и его главные действия | **Заголовок блока внутри экрана** — [секция](./section.md): она тише и не претендует на `<h1>` |
| Пояснение, что это за экран, — одно-два предложения | **Заголовок области приложения** — `inst-panel-title`, [панель](../components/display/panel.md) |
| Место для 1–3 действий над экраном целиком | **Инструменты над содержимым панели** — [тулбар](../components/navigation/toolbar.md) |
| — | **Путь до экрана** — [крошки](../components/navigation/breadcrumbs.md) в шапке [оболочки](./shell.md), а не здесь |

## Разметка

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

## Действия

`inst-page-actions` прижимается к дальнему краю (`margin-inline-start: auto`) и
переносится внутри себя. Держите здесь **одно** `primary`: экран, где всё
primary, не имеет primary вовсе — см. [кнопку](../components/actions/button.md).

Четвёртое и следующие действия уходят в [меню](../components/overlays/menu.md).
Полоса из шести кнопок над экраном не читается ни глазом, ни с клавиатуры.

## Со статусом

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

## Справочник

### Классы

| Класс | Работа |
|---|---|
| `inst-page-header` | Контейнер. Ряд с переносом, выравнивание по верху |
| `inst-page-header-main` | Колонка «название + пояснение». Порог переноса 20rem |
| `inst-page-title` | Название экрана. Кегль `--text-xl` |
| `inst-page-desc` | Пояснение. `--text-sm`, `--text-secondary`, мера 68ch |
| `inst-page-actions` | Действия. Прижаты к дальнему краю, переносятся |

### Обязательная разметка

| Что | Почему |
|---|---|
| `<h1>` под `inst-page-title` | Класс даёт кегль, а не уровень. Уровень — работа тега, и скринридер читает именно его |
| `inst-page-header-main` как обёртка | Без неё название и пояснение становятся отдельными детьми ряда и встают в строку |
| `type="button"` у действий | Иначе внутри формы они её отправят |

### Токены

`--gap-inline` · `--pad-panel` · `--space-2` · `--space-3` · `--text-xl` ·
`--text-sm` · `--text-secondary`

## Доступность

| | |
|---|---|
| Уровень заголовка | `<h1>` один на документ. Кит его не назначает — ставит автор разметки |
| Мера пояснения | 68ch. Строка длиннее теряется на возврате к началу |
| Перенос | Шапка переносится сама; при увеличении кегля до 200% действия уходят под заголовок, а не за край экрана |
| Контраст | Пояснение — `--text-secondary` (4.5:1). `--text-faint` для текста, который читают, не применяется |
| Печать | Название и пояснение печатаются, `inst-page-actions` — нет: действия на бумаге не действия |
| Порядок | Действия стоят в разметке после заголовка и обходятся после него. `order` не применяется |

## Связанное

[Оболочка](./shell.md) · [Секция](./section.md) ·
[Контейнер](./container.md) · [Кнопка](../components/actions/button.md) ·
[Крошки](../components/navigation/breadcrumbs.md) ·
[Вкладки](../components/navigation/tabs.md)

Исходник: `src/layout.css` · Почему кит устроен именно так — [конституция](../../DESIGN.md)
