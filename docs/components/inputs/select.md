---
title: Селект
group: Ввод
layout: component
source: src/forms.css
js-optional: Свой список вместо нативного. Нативный работает без скрипта
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
| Список известен заранее и не меняется | **Список, который нужно искать или дополнять** — комбобокс. В библиотеке его пока нет |
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

## JS

Нативный `<select>` не требует ни строки: клавиатура, набор по первым буквам,
список на весь экран телефона — всё от платформы. Скрипт нужен только там, где
нативного мало: поиск по списку, две строки в пункте, иконка ряда.

### Когда стоит уходить от нативного

| | |
|---|---|
| Больше двух-трёх десятков вариантов | Нужен поиск, а в `<select>` его нет |
| В пункте не только текст | Нативный список рисует одну строку без разметки |
| Выбор нескольких | `<select multiple>` выглядит как окно списка и на телефоне неудобен |

Во всех остальных случаях нативный выигрывает, и заменять его — работа без
выигрыша.

### Свой список

Поверхность и клавиатуру даёт библиотека. Приложению остаются данные: открыть,
отфильтровать, записать выбор в скрытое поле.

```html preview
<div class="inst-field" style="max-inline-size:16rem">
  <label class="inst-label" for="pick">Прогон</label>
  <span class="inst-select-wrap">
    <button class="inst-select" type="button" id="pick"
            popovertarget="pick-list" aria-haspopup="listbox">worldgen-01</button>
  </span>
  <div class="inst-popover inst-popover--fill inst-popover--anchored" id="pick-list" popover>
    <input class="inst-input inst-input--sm" type="search" aria-label="Поиск по прогонам" placeholder="Поиск">
    <div role="listbox" aria-label="Прогоны">
      <div class="inst-menu-item" role="option" aria-selected="true" tabindex="0">worldgen-01</div>
      <div class="inst-menu-item" role="option" aria-selected="false" tabindex="-1">biomes-04</div>
      <div class="inst-menu-item" role="option" aria-selected="false" tabindex="-1">rivers-12</div>
    </div>
  </div>
</div>
```

Стрелки, `Home`, `End` и перенос выбора внутри списка выполняет
[`instrument.js`](../../foundations/behavior.md) — это `role="listbox"`. Открытие и
закрытие берёт Popover API.

### События

```js
list.addEventListener('inst:select', (e) => {
  trigger.textContent = e.detail.value;
  hidden.value = e.detail.value;      // значение уходит с формой
  list.hidePopover();
});
```

Скрытое поле обязательно: свой список — это разметка, а форма отправляет
поля. Без него выбор виден человеку и не виден серверу.

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
