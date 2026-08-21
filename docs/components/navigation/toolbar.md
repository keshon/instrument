---
title: Тулбар
group: Навигация
layout: component
source: src/surfaces.css
api:
  - { name: "inst-toolbar", kind: "класс", doc: "Полоса. Ряд с зазором `--space-3` и нижней линией" }
  - { name: "inst-toolbar-sep", kind: "класс", doc: "Вертикальный разделитель групп" }
  - { name: "inst-toolbar-spacer", kind: "класс", doc: "Прижать хвост полосы к дальнему краю" }
  - { name: "--space-2", kind: "токен" }
  - { name: "--space-3", kind: "токен" }
  - { name: "--pad-cell-x", kind: "токен" }
  - { name: "--hairline", kind: "токен" }
  - { name: "--border-subtle", kind: "токен" }
  - { name: "--border", kind: "токен" }
title-en: "Toolbar"
group-en: "Navigation"
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

## Контракт

| Что | Обязательно | Почему |
|---|---|---|
| `type="button"` у каждой кнопки | да | Иначе внутри формы она её отправит |
| `aria-label` у иконочной кнопки | да | Кнопка без текста не имеет доступного имени: «Переместить», а не «Инструмент 1» |
| `aria-hidden="true"` у иконки | да | Иначе скринридер прочитает её дважды или не прочитает вовсе |
| `aria-pressed` у кнопки-режима | да, если инструмент остаётся включённым | Состояние обязано быть в разметке. Ставит приложение |
| Одна `inst-toolbar-spacer` на полосу | да | Распорка отбивает хвост — статус, счётчик, закрытие. Двух точек разрыва в одной полосе не бывает |

Роль `toolbar` не назначается автоматически: она обещает стрелочную навигацию и бегущий
`tabindex`, а это поведение — слой приложения. Полоса из обычных кнопок без
роли доступна и без обещаний.

### Доступность

| | |
|---|---|
| Клавиатура | Каждая кнопка обходится `Tab`. Стрелочная навигация внутри полосы — только вместе с `role="toolbar"` и бегущим `tabindex`, то есть с JS приложения |
| Имена | Иконочные кнопки обязаны нести `aria-label`: «Переместить», а не «Инструмент 1» |
| Цель нажатия | `inst-btn--sm` — 26px. В плотности `compact` — 22px, ниже порога 24px по WCAG 2.5.8 |
| Разделитель | Пустой элемент без текста и роли — в дереве доступности его нет, и это правильно: он не сообщает ничего сверх порядка |
| Контраст | Линии `--border-subtle` и `--border` — декорация (порог 3:1), а не текст |
| Печать | Тулбар не печатается: это хром приложения, а не данные |

## Устройство

Два служебных элемента, которые легко перепутать:

| Элемент | Что делает | Когда |
|---|---|---|
| `inst-toolbar-sep` | Вертикальная линия на всю высоту полосы | Между **смысловыми группами** инструментов: трансформации отдельно, отображение отдельно |
| `inst-toolbar-spacer` | `margin-inline-start: auto` — всё после него уезжает к дальнему краю | Один раз на полосу, перед хвостом: статусом, счётчиком, кнопкой закрытия |

Разделитель — пустой `<span>`, то есть чистая графика без текста и без роли.
Группировку он показывает глазу; для скринридера её несёт порядок и подписи
кнопок.

## Композиции

### Место в панели

Тулбар стоит **между шапкой и телом** панели и приносит нижнюю линию сам. Шапка
при этом свою линию снимает (`:has(+ .inst-toolbar)`) — иначе на тридцати
пикселях оказалось бы два правила подряд.

```html preview context
<div class="inst-panel">
  <div class="inst-panel-header"><span class="inst-panel-title">Модули</span></div>
  <div class="inst-toolbar">
    <div class="inst-segmented" role="radiogroup" aria-label="Вид">
      <button type="button" role="radio" aria-checked="true"  tabindex="0">Список</button>
      <button type="button" role="radio" aria-checked="false" tabindex="-1">Сетка</button>
    </div>
    <span class="inst-toolbar-sep"></span>
    <span class="inst-search">
      <input class="inst-input inst-input--sm" type="search" placeholder="Фильтр" aria-label="Фильтр по модулям">
    </span>
    <span class="inst-toolbar-spacer"></span>
    <span class="inst-badge">3 из 12</span>
  </div>
  <div class="inst-panel-body">Содержимое области.</div>
</div>
```

Тулбар без шапки тоже законен: панель, у которой нет названия, но есть режимы.

## API

```api
```

## Связанное

```related
```
