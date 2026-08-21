---
title: Группа кнопок
group: Действия
layout: component
source: src/actions.css
api:
  - { name: "inst-btn-group", kind: "класс", doc: "Контейнер. Внутри — только `inst-btn`" }
  - { name: "--radius-md", kind: "токен" }
  - { name: "--hairline", kind: "токен" }
---

Несколько разных действий, поставленных вплотную. Группа говорит «эти действия
про одно и то же», но каждое остаётся самостоятельным.

```html preview
<div class="inst-btn-group">
  <button class="inst-btn inst-btn--sm" type="button">Слева</button>
  <button class="inst-btn inst-btn--sm" type="button">Центр</button>
  <button class="inst-btn inst-btn--sm" type="button">Справа</button>
</div>
```

## Контракт

Группа не добавляет ролей: каждая кнопка объявляет себя сама — в этом и отличие
от [сегментированного контрола](./segmented.md), где роль обязательна.

| Что | Обязательно | Почему |
|---|---|---|
| `.inst-btn-group` на контейнере | да | Снимает рамку у соседа и скругляет только торцы; без него кнопки стоят рядом, а не в группе |
| Обычные `.inst-btn` внутри | да | Группа — это раскладка, а не новый компонент: кнопка внутри неё остаётся кнопкой |
| Один размер на всю группу | да | Соседние кнопки разной высоты дают ступеньку на шве, и её видно на любом фоне |
| Своё имя у каждой кнопки | да | Группа не нуждается в `aria-label`: она не один контрол, а несколько действий рядом |
| Бегущий `tabindex` | **нет, это ошибка** | `Tab` проходит по каждой кнопке отдельно — это разные действия, а не один выбор. Бегущий `tabindex` спрятал бы часть действий от клавиатуры |

Кнопка под фокусом поднимается по `z-index`, иначе кольцо срезал бы сосед:
элементы флекса подчиняются `z-index` без `position`.

## Устройство

Шов группы — **снятая рамка у соседа**. Отрицательное поле в `-0.5px` не
годится: браузер не рисует рамку тоньше физического пикселя, и сдвиг оставлял
утолщение на каждом стыке.

| Что делает контейнер | Как |
|---|---|
| Снимает шов | У всех кнопок, кроме первой, убирается ближняя рамка |
| Скругляет торцы | Радиус остаётся только у первой и последней |
| Поднимает фокус | Кнопка под фокусом идёт выше по `z-index`, иначе кольцо срезал бы сосед |

## Размеры

Размер задаётся каждой кнопке, а не контейнеру: группа не контрол, а
композиция, и своего размера у неё нет.

```html preview
<div class="inst-btn-group">
  <button class="inst-btn inst-btn--sm" type="button">Один</button>
  <button class="inst-btn inst-btn--sm" type="button">Два</button>
</div>
<div class="inst-btn-group">
  <button class="inst-btn" type="button">Один</button>
  <button class="inst-btn" type="button">Два</button>
</div>
<div class="inst-btn-group">
  <button class="inst-btn inst-btn--lg" type="button">Один</button>
  <button class="inst-btn inst-btn--lg" type="button">Два</button>
</div>
```

## Композиции

### С иконочными кнопками

```html preview
<div class="inst-btn-group">
  <button class="inst-btn inst-btn--sm inst-btn--icon" type="button" aria-label="Обновить">
    <svg class="inst-icon" aria-hidden="true"><use href="#i-refresh"/></svg>
  </button>
  <button class="inst-btn inst-btn--sm inst-btn--icon" type="button" aria-label="Копировать">
    <svg class="inst-icon" aria-hidden="true"><use href="#i-copy"/></svg>
  </button>
  <button class="inst-btn inst-btn--sm inst-btn--icon" type="button" aria-label="Настройки">
    <svg class="inst-icon" aria-hidden="true"><use href="#i-settings"/></svg>
  </button>
</div>
```

Каждая иконочная кнопка несёт свой `aria-label`: у группы нет общего имени,
потому что она не один контрол.

## Сценарии

### В тулбаре

Группы разделены зазором, внутри группы — вплотную. Зазор и есть сообщение о
том, где кончается одно родство и начинается другое.

```html preview context
<div class="inst-toolbar">
  <div class="inst-btn-group">
    <button class="inst-btn inst-btn--sm" type="button">Список</button>
    <button class="inst-btn inst-btn--sm" type="button">Сетка</button>
  </div>
  <span class="inst-toolbar-sep"></span>
  <div class="inst-btn-group">
    <button class="inst-btn inst-btn--sm" type="button">Копировать</button>
    <button class="inst-btn inst-btn--sm" type="button">Дублировать</button>
  </div>
</div>
```

## API

```api
```

## Связанное

```related
```
