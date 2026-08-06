---
title: Строка очереди
group: Агентный слой
status: stable
source: src/components.css
needs-js: Бегущий tabindex, стрелки, выбор строки
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

## Строка очереди или шаг

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

## Жёлоб

```html
<span class="inst-task-gutter"><span class="inst-dot"></span></span>
```

Левая колонка фиксированной ширины (`--size-gutter`). Она держит точку
состояния на одной вертикали во всех строках — без неё заголовки разной длины
разъезжают маркеры, и очередь перестаёт читаться столбцом.

## Справочник

### Классы

| Класс | Работа |
|---|---|
| `inst-task` | Строка |
| `inst-task-gutter` | Левый жёлоб под точку |
| `inst-task-main` | Заголовок и подпись |
| `inst-task-title` | Название задачи |
| `inst-task-sub` | Что именно делается |
| `inst-task-meta` | Время или счётчик, прижат к концу |

Длинный заголовок обрезается через `inst-u-truncate` — иначе он выталкивает
метаданные за край.

### Атрибуты

| Атрибут | Значения |
|---|---|
| `data-state` | `queued` `running` `done` `warn` `failed` `skipped` |
| `aria-selected` | `true` · `false` |

### Токены

`--size-gutter` · `--size-dot` · `--row-pad-y` · `--pad-cell-x` ·
`--gap-inline` · `--space-1` · `--text-sm` · `--text-xs` ·
`--surface-selected` · `--surface-hover` · `--tone-ink`

## Доступность

| | |
|---|---|
| Роли обязательны | Контейнер `role="listbox"`, строка `role="option"`. Без них `aria-selected` **невалиден**, и строка не озвучится как выбираемая |
| Клавиатура | Бегущий `tabindex`: `0` у активной, `-1` у остальных. `Tab` входит в список один раз, дальше стрелки, `Home`/`End`. Поведение — на приложении |
| Состояние словом | `data-state` красит строку, но не озвучивается. Состояние обязано быть в тексте — в подписи или в [бейдже](../components/display/badge.md) рядом |
| Обрезанный заголовок | `inst-u-truncate` режет визуально, полный текст остаётся в DOM и озвучивается. Добавьте `title` для мыши |
| Уменьшенное движение | Пульсация точки `running` замедляется, а не гаснет: очередь агентов без признака работы бесполезна |
| Выделение и наведение складываются | Подсветка выбранной полупрозрачна, поэтому наведение на неё остаётся видно |

## Связанное

[Шаг](./step.md) · [Лог](./log.md) · [Дорожки прогонов](./lane.md) ·
[Бейдж](../components/display/badge.md) ·
[Таблица](../components/display/table.md)

Исходник: `src/components.css` · Почему кит устроен именно так — [конституция](../about/design-principles.md)
