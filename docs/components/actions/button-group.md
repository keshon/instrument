---
title: Группа кнопок
group: Действия
layout: component
source: src/components.css
api:
  - { name: "inst-btn-group", kind: "класс", doc: "Контейнер. Внутри — только `inst-btn`" }
  - { name: "--radius-md", kind: "токен" }
  - { name: "--hairline", kind: "токен" }
title-en: "Button group"
group-en: "Actions"
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

## Использование

Контейнер и обычные кнопки внутри. Группа не добавляет ролей: каждая кнопка
объявляет себя сама — в этом и отличие от [сегментированного
контрола](./segmented.md), где роль обязательна.

```html
<div class="inst-btn-group">
  <button class="inst-btn" type="button">…</button>
  <button class="inst-btn" type="button">…</button>
</div>
```

:::warn
Размер должен быть **один на всю группу**. Соседние кнопки разной высоты дают
ступеньку на шве, и её видно на любом фоне.
:::

## Когда использовать

| Используйте | Возьмите другое |
|---|---|
| Разные действия над одним объектом: копировать · дублировать · удалить | **Выбор одного из равных вариантов** — [сегментированный контрол](./segmented.md). У группы нет выбранного элемента, и быть не может |
| Действие плюс его вариант: «Сохранить» и «Сохранить как» | **Список действий, который не влезает** — [меню](../overlays/menu.md) в поповере |
| Панель инструментов, где кнопки сгруппированы по смыслу | **Разрозненные действия** — [кластер](../../layout/flow.md) с зазором: вплотную они соврут о родстве |

## Устройство

Шов группы — **снятая рамка у соседа**, а не отрицательное поле в `-0.5px`:
браузер никогда не рисует рамку тоньше физического пикселя, поэтому сдвиг на
полпикселя оставлял видимое утолщение на каждом стыке.

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

## Правила

:::do Один размер на группу
Все кнопки внутри одного размера. Разная высота даёт ступеньку на шве.
:::

:::dont Группа как выбор
У группы нет выбранного элемента. Если один из вариантов должен оставаться
нажатым — это [сегментированный контрол](./segmented.md).
:::

:::do Tab по каждой кнопке
Это разные действия, и до каждого надо доходить с клавиатуры отдельно.
:::

:::dont Бегущий tabindex
`tabindex="-1"` у части кнопок спрятал бы действия от клавиатуры. Бегущий
`tabindex` уместен там, где контрол один, — в сегментированном.
:::

## Доступность

| | |
|---|---|
| Клавиатура | `Tab` проходит по каждой кнопке отдельно — это разные действия, а не один выбор. Бегущий `tabindex` здесь был бы **ошибкой**: он спрятал бы часть действий от клавиатуры |
| Фокус | Кнопка под фокусом поднимается по `z-index`, иначе кольцо срезал бы сосед. Элементы флекса подчиняются `z-index` без `position` |
| Имена | Каждая кнопка несёт своё имя. Группа не нуждается в `aria-label`, потому что не является одним контролом |

## API

```api
```

## Связанное

[Кнопка](./button.md) [Сегментированный контрол](./segmented.md)
[Меню](../overlays/menu.md) [Тулбар](../navigation/toolbar.md)
