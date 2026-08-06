---
title: Тулбар
group: Навигация
status: stable
source: src/components.css
---

Полоса инструментов внутри [панели](../display/panel.md): режимы, фильтры,
мелкие действия над её содержимым. Не шапка экрана и не меню.

```html preview
<div class="inst-panel">
  <div class="inst-panel-header"><span class="inst-panel-title">Инспектор</span></div>
  <div class="inst-toolbar">
    <button class="inst-btn inst-btn--sm inst-btn--ghost inst-btn--icon" type="button" aria-label="Переместить">
      <svg class="inst-icon" aria-hidden="true"><use href="#i-move"/></svg></button>
    <button class="inst-btn inst-btn--sm inst-btn--ghost inst-btn--icon" type="button" aria-label="Повернуть">
      <svg class="inst-icon" aria-hidden="true"><use href="#i-rotate"/></svg></button>
    <button class="inst-btn inst-btn--sm inst-btn--ghost inst-btn--icon" type="button" aria-label="Масштаб">
      <svg class="inst-icon" aria-hidden="true"><use href="#i-scale"/></svg></button>
    <span class="inst-toolbar-sep"></span>
    <span class="inst-toolbar-spacer"></span>
    <span class="inst-badge">terrain_chunk_04</span>
  </div>
  <div class="inst-panel-body">…</div>
</div>
```

## Установка

```html
<link rel="stylesheet" href="src/kit.css">
```

## Когда использовать

| Используйте | Возьмите другое |
|---|---|
| Инструменты, действующие на содержимое **этой** панели | **Действия над экраном целиком** — [шапка экрана](../../layout/page-header.md) и её `inst-page-actions` |
| Постоянно видимая полоса: режим переключают часто | **Редкие действия по клику** — [меню](../overlays/menu.md): полоса из десяти иконок дороже одной кнопки «ещё» |
| Иконочные кнопки и мелкие контролы в один ряд | **Просто ряд кнопок вне панели** — [кластер](../../layout/flow.md): у тулбара есть рамка и место в поверхности |
| Пара действий у заголовка секции | **Заголовок секции** — `inst-section-actions`, [секция](../../layout/section.md) |

## Место в панели

Тулбар стоит **между шапкой и телом** панели и приносит нижнюю линию сам. Шапка
при этом свою линию снимает (`:has(+ .inst-toolbar)`) — иначе на тридцати
пикселях оказалось бы два правила подряд.

```html
<div class="inst-panel">
  <div class="inst-panel-header"><span class="inst-panel-title">Инспектор</span></div>
  <div class="inst-toolbar">…</div>
  <div class="inst-panel-body inst-panel-body--flush">…</div>
</div>
```

Тулбар без шапки тоже законен: панель, у которой нет названия, но есть режимы.

## Разделитель и распорка

Два служебных элемента, которые легко перепутать:

| Элемент | Что делает | Когда |
|---|---|---|
| `inst-toolbar-sep` | Вертикальная линия на всю высоту полосы | Между **смысловыми группами** инструментов: трансформации отдельно, отображение отдельно |
| `inst-toolbar-spacer` | `margin-inline-start: auto` — всё после него уезжает к дальнему краю | Один раз на полосу, перед хвостом: статусом, счётчиком, кнопкой закрытия |

Разделитель — пустой `<span>`, то есть чистая графика без текста и без роли.
Группировку он показывает глазу; для скринридера её несёт порядок и подписи
кнопок.

## Справочник

### Классы

| Класс | Работа |
|---|---|
| `inst-toolbar` | Полоса. Ряд с зазором `--space-3` и нижней линией |
| `inst-toolbar-sep` | Вертикальный разделитель групп |
| `inst-toolbar-spacer` | Прижать хвост полосы к дальнему краю |

### Обязательная разметка

| Что | Почему |
|---|---|
| `type="button"` у каждой кнопки | Иначе внутри формы она её отправит |
| `aria-label` у иконочной кнопки | Кнопка без текста не имеет доступного имени |
| `aria-hidden="true"` у иконки | Иначе скринридер прочитает её дважды или не прочитает вовсе |
| `aria-pressed` у кнопки-режима | Если инструмент остаётся включённым, состояние обязано быть в разметке. Ставит приложение |

Роль `toolbar` кит не назначает: она обещает стрелочную навигацию и бегущий
`tabindex`, а это поведение — слой приложения. Полоса из обычных кнопок без
роли доступна и без обещаний.

### Токены

`--space-2` · `--space-3` · `--pad-cell-x` · `--hairline` ·
`--border-subtle` · `--border`

## Доступность

| | |
|---|---|
| Клавиатура | Каждая кнопка обходится `Tab`. Стрелочная навигация внутри полосы — только вместе с `role="toolbar"` и бегущим `tabindex`, то есть с JS приложения |
| Имена | Иконочные кнопки обязаны нести `aria-label`: «Переместить», а не «Инструмент 1» |
| Цель нажатия | `inst-btn--sm` — 26px. В плотности `compact` — 22px, ниже порога 24px по WCAG 2.5.8 |
| Разделитель | Пустой элемент без текста и роли — в дереве доступности его нет, и это правильно: он не сообщает ничего сверх порядка |
| Контраст | Линии `--border-subtle` и `--border` — декорация (порог 3:1), а не текст |
| Печать | Тулбар не печатается: это хром приложения, а не данные |

## Связанное

[Панель](../display/panel.md) · [Кнопка](../actions/button.md) ·
[Группа кнопок](../actions/button-group.md) ·
[Сегментированный контрол](../actions/segmented.md) · [Меню](../overlays/menu.md) ·
[Секция](../../layout/section.md)

Исходник: `src/components.css` · Почему кит устроен именно так — [конституция](../../about/design-principles.md)
