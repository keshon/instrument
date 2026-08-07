---
title: Селект
group: Ввод
status: stable
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

## Обёртка обязательна

`inst-select-wrap` — не украшение и не привычка. В ней живёт шеврон:
псевдоэлемент самого `<select>` браузер не отдаёт.

```html
<span class="inst-select-wrap">
  <select class="inst-select">…</select>
</span>
```

Без обёртки поле останется без стрелки и будет неотличимо от текстового.

## Когда использовать

| Используйте | Возьмите другое |
|---|---|
| Больше пяти взаимоисключающих вариантов | **2–4 варианта, все на виду** — [сегментированный контрол](../actions/segmented.md) |
| Список известен заранее и не меняется | **Список, который нужно искать или дополнять** — комбобокс. В ките его пока нет, см. роадмап |
| Значение выбирается редко | **Независимые флаги** — [чекбоксы](./toggles.md): в селекте выбран всегда ровно один |
| Множественный выбор — `multiple` | **Действия, а не значения** — [меню](../overlays/menu.md) в поповере |

## Размеры и состояния

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

Размеры, врез, рамка, фокус, ошибка, `disabled` и `readonly` — те же, что у
[текстового поля](./input.md): все три контрола объявлены одним правилом,
поэтому разойтись не могут.

## Справочник

```api
```

## Доступность

| | |
|---|---|
| Клавиатура | Стрелки, `Home`/`End`, поиск по первым буквам, `Enter`. Всё нативное |
| Подпись | `<label for>` или `aria-label`. Первый `<option>` подписью не является |
| Шеврон | `pointer-events: none` — клик по стрелке открывает список, а не проваливается мимо |
| Режим принудительных цветов | Шеврон — маска, в `forced-colors` теряющая цвет фона, поэтому переводится на `CanvasText` явно |
| Мобильные | Нативный список превращается в системное колесо. Самодельный выпадающий список это теряет |

## Связанное

[Текстовое поле](./input.md) · [Переключатели](./toggles.md) ·
[Сегментированный контрол](../actions/segmented.md) · [Форма](./form.md)
