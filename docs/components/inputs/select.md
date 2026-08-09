---
title: Селект
group: Ввод
layout: component
source: src/forms.css
api:
  - { name: "inst-select-wrap", kind: "класс", doc: "Обёртка. **Обязательна** — в ней шеврон" }
  - { name: "inst-select", kind: "класс", doc: "Сам `<select>`" }
  - { name: "inst-select--sm", kind: "модификатор", doc: "Размер" }
  - { name: "inst-select--lg", kind: "модификатор", doc: "Размер" }
  - { name: "--surface-field", kind: "токен" }
  - { name: "--control-h-sm/md/lg", kind: "токен" }
  - { name: "--control-pad-sm", kind: "токен" }
  - { name: "--size-chevron", kind: "токен" }
  - { name: "--space-7", kind: "токен" }
  - { name: "--radius-sm", kind: "токен" }
  - { name: "--border-control", kind: "токен" }
  - { name: "--text-muted", kind: "токен" }
title-en: "Select"
group-en: "Inputs"
---

Выбор одного значения из списка. Нативный `<select>` — со всей клавиатурой,
поиском по первым буквам и системным выпадающим списком, который на телефоне
станет колесом.

```html preview
<span class="inst-select-wrap">
  <select class="inst-select" aria-label="Модель">
    <option>opus</option>
    <option>sonnet</option>
    <option>haiku</option>
  </select>
</span>
```

## Использование

Обёртка обязательна — в ней живёт шеврон: псевдоэлемент самого `<select>`
браузер не отдаёт.

```html
<span class="inst-select-wrap">
  <select class="inst-select" aria-label="Модель">…</select>
</span>
```

| Что | Обязательно | Почему |
|---|---|---|
| `inst-select-wrap` на обёртке | да | Шеврон рисуется псевдоэлементом обёртки. Без неё поле останется без стрелки и будет неотличимо от текстового |
| `inst-select` на поле | да | Врез, рамка, высота и состояния приходят оттуда |
| `<label for>` или `aria-label` | да | Первый `<option>` подписью не является |

## Когда использовать

| Используйте | Возьмите другое |
|---|---|
| Больше пяти взаимоисключающих вариантов | **2–4 варианта, все на виду** — [сегментированный контрол](../actions/segmented.md) |
| Список известен заранее и не меняется | **Список, который нужно искать или дополнять** — комбобокс. В ките его пока нет, см. роадмап |
| Значение выбирается редко | **Независимые флаги** — [чекбоксы](./toggles.md): в селекте выбран всегда ровно один |
| Множественный выбор — `multiple` | **Действия, а не значения** — [меню](../overlays/menu.md) в поповере |

## Размеры

```html preview
<span class="inst-select-wrap">
  <select class="inst-select inst-select--sm" aria-label="sm"><option>sm</option></select>
</span>
<span class="inst-select-wrap">
  <select class="inst-select" aria-label="md"><option>md</option></select>
</span>
<span class="inst-select-wrap">
  <select class="inst-select inst-select--lg" aria-label="lg"><option>lg</option></select>
</span>
```

Симметричны кнопочным и полевым: селект, [поле](./input.md) и
[кнопка](../actions/button.md) одного размера стоят в одной строке формы и
разойтись не могут.

## Состояния

```html preview
<span class="inst-select-wrap">
  <select class="inst-select" aria-label="Обычный"><option>Обычный</option></select>
</span>
<span class="inst-select-wrap">
  <select class="inst-select" aria-label="Недоступен" disabled><option>Недоступен</option></select>
</span>
<span class="inst-select-wrap">
  <select class="inst-select" aria-label="Ошибка" aria-invalid="true"><option>Ошибка</option></select>
</span>
```

Врез, рамка, фокус, ошибка, `disabled` и `readonly` — те же, что у
[текстового поля](./input.md): все три контрола объявлены одним правилом,
поэтому разойтись не могут.

## Правила

:::do Селект от шести вариантов
Ниже шести выбор быстрее делается [сегментированным
контролом](../actions/segmented.md): все варианты видны, и переключение стоит
одного нажатия.
:::

:::dont Селект вместо действий
Список действий над объектом — [меню](../overlays/menu.md) в поповере. Селект
хранит значение, а не запускает работу.
:::

:::do Нативный select
Клавиатура, поиск по первым буквам и системное колесо на телефоне приходят от
платформы и не пишутся заново.
:::

:::dont Самодельный выпадающий список
Он теряет всё перечисленное и заводит JS ради того, что уже работает.
:::

## Доступность

| | |
|---|---|
| Клавиатура | Стрелки, `Home`/`End`, поиск по первым буквам, `Enter`. Всё нативное |
| Подпись | `<label for>` или `aria-label`. Первый `<option>` подписью не является |
| Шеврон | `pointer-events: none` — клик по стрелке открывает список, а не проваливается мимо |
| Режим принудительных цветов | Шеврон — маска, в `forced-colors` теряющая цвет фона, поэтому переводится на `CanvasText` явно |
| Мобильные | Нативный список превращается в системное колесо. Самодельный выпадающий список это теряет |

## API

```api
```

## Связанное

[Текстовое поле](./input.md) [Переключатели](./toggles.md)
[Сегментированный контрол](../actions/segmented.md) [Форма](./form.md)
