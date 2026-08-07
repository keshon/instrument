---
title: Числовое поле инспектора
group: Ввод
status: stable
source: src/components.css
needs-js: перетаскивание подписи оси для изменения значения
api:
  - { name: "inst-num-field", kind: "класс", doc: "Обёртка: рамка, врез, фокус. Высота — `--control-h-sm`" }
  - { name: "inst-num-axis", kind: "класс", doc: "Буква оси. Ручка перетаскивания, **не подпись**" }
  - { name: "inst-vec", kind: "класс", doc: "Ряд из нескольких полей равной ширины" }
  - { name: "inst-prop", kind: "класс", doc: "Строка свойства инспектора" }
  - { name: "inst-prop-label", kind: "класс", doc: "Имя свойства. Обрезается, поэтому обязателен `title`" }
  - { name: "inst-prop-control", kind: "класс", doc: "Контрольная часть строки" }
  - { name: "--control-h-sm", kind: "токен" }
  - { name: "--radius-sm", kind: "токен" }
  - { name: "--surface-field", kind: "токен" }
  - { name: "--border-control", kind: "токен" }
  - { name: "--accent-border", kind: "токен" }
  - { name: "--text-2xs", kind: "токен" }
  - { name: "--text-xs", kind: "токен" }
  - { name: "--space-2", kind: "токен" }
  - { name: "--space-3", kind: "токен" }
---

Число с буквой оси слева — примитив инспектора: три таких поля составляют
вектор, и в панели свойств их десятки. Компактнее [текстового
поля](./input.md), потому что подпись оси занимает место степпера, а не
отдельную колонку.

```html preview
<span class="inst-num-field">
  <span class="inst-num-axis">X</span>
  <input type="number" value="128" aria-label="Позиция X">
</span>
<span class="inst-num-field">
  <span class="inst-num-axis">Y</span>
  <input type="number" value="0" aria-label="Позиция Y">
</span>
<span class="inst-num-field">
  <span class="inst-num-axis">Z</span>
  <input type="number" value="-64" aria-label="Позиция Z">
</span>
```

## Установка

```html
<link rel="stylesheet" href="src/kit.css">
```

## Когда использовать

| Используйте | Возьмите другое |
|---|---|
| Компонента вектора: позиция, поворот, масштаб | **Одиночное число в обычной форме** — [текстовое поле](./input.md) с `inst-input-group` и суффиксом: буква оси там ничего не значит |
| Плотный ряд однотипных чисел в инспекторе | **Число с единицей измерения** — приставка и суффикс [текстового поля](./input.md): `inst-num-axis` не единица, а ось |
| Точное значение, которое вводят и правят цифрами | **Приблизительное значение из диапазона** — [слайдер](./slider.md): по нему видно положение внутри шкалы |
| Значение внутри строки свойства | **Выбор из набора** — [селект](./select.md) |

## Контракт разметки

Здесь легко получить поле, которое выглядит подписанным и не является им.

| Что | Обязательно | Почему |
|---|---|---|
| `aria-label` на `<input>`, **включающий ось** | да | Скринридер иначе прочитает «поле, 128» и не скажет, чего именно. Ось — часть имени: «Позиция X», а не «X» |
| `inst-num-axis` **не** `<label>` | да | Это ручка перетаскивания и визуальный маркер, а не подпись. `<label for>` сделал бы буквой «X» полное доступное имя поля — вместо имени свойства получилось бы одно слово |
| `type="number"` | да | Клавиатурные стрелки, шаг и числовая клавиатура на сенсорных — от платформы |
| `<input>` прямым ребёнком | да | Правила написаны через `& > input` |

Отсюда следствие, которое проще нарушить, чем заметить: **буква оси не должна
быть единственным местом, где написано, какое это поле**. Имя свойства живёт в
`inst-prop-label` рядом и повторяется в `aria-label` каждого поля — три
одинаковых квадратика с буквами X, Y, Z вне контекста строки не значат ничего.

> **Перетаскивание — это JS.** Курсор `ew-resize` на подписи оси означает
> «отсюда тянут», но само изменение значения тягой пишет приложение. Без
> скрипта поле остаётся обычным числовым вводом, и он полностью работает:
> цифры, стрелки, колесо, вставка.

## Состояния

```html preview
<span class="inst-num-field">
  <span class="inst-num-axis">X</span>
  <input type="number" value="128" aria-label="Позиция X">
</span>
<span class="inst-num-field">
  <span class="inst-num-axis">Y</span>
  <input type="number" value="0" aria-label="Позиция Y" readonly>
</span>
<span class="inst-num-field">
  <span class="inst-num-axis">Z</span>
  <input type="number" value="-64" aria-label="Позиция Z" disabled>
</span>
```

| Состояние | Как ставится | Что происходит |
|---|---|---|
| наведение | `:hover` на обёртке | Рамка темнеет до `--text-muted` |
| фокус | `:focus-within` на обёртке | Рамка акцентом. Кольцо рисует **обёртка**, а не поле: у поля рамки нет вовсе |
| только чтение | `readonly` на `<input>` | Значение читается и копируется. Не `disabled` — разницу см. в [текстовом поле](./input.md) |
| недоступно | `disabled` на `<input>` | Нативная недоступность |

Нативные стрелки-степперы сняты (`appearance: textfield`): в плотном ряду
инспектора они съедали треть ширины поля и попадали под курсор чаще, чем сам
ввод. Клавиатурные стрелки при этом продолжают работать — снято оформление, а
не поведение.

## Композиции

### Вектор

```html preview
<span class="inst-vec">
  <span class="inst-num-field"><span class="inst-num-axis">X</span><input type="number" value="128" aria-label="Позиция X"></span>
  <span class="inst-num-field"><span class="inst-num-axis">Y</span><input type="number" value="0" aria-label="Позиция Y"></span>
  <span class="inst-num-field"><span class="inst-num-axis">Z</span><input type="number" value="-64" aria-label="Позиция Z"></span>
</span>
```

`inst-vec` раздаёт трём полям равные доли ширины, поэтому вектор не
перекашивается, когда в одну компоненту вводят шестизначное число.

### Строка свойства

```html
<div class="inst-prop">
  <span class="inst-prop-label" title="Позиция">Позиция</span>
  <span class="inst-prop-control inst-vec">
    <span class="inst-num-field"><span class="inst-num-axis">X</span><input type="number" value="128" aria-label="Позиция X"></span>
    <span class="inst-num-field"><span class="inst-num-axis">Y</span><input type="number" value="0" aria-label="Позиция Y"></span>
    <span class="inst-num-field"><span class="inst-num-axis">Z</span><input type="number" value="-64" aria-label="Позиция Z"></span>
  </span>
</div>
```

### Одиночное значение с осью-именем

```html
<span class="inst-num-field">
  <span class="inst-num-axis">R</span>
  <input type="number" value="45" step="5" aria-label="Радиус влияния">
</span>
```

## Справочник

```api
```

### Обязательная разметка

| Что | Почему |
|---|---|
| `aria-label` с осью на `<input>` | Единственный носитель имени поля |
| `type="number"` | Стрелки, шаг, числовая клавиатура |
| `title` на `inst-prop-label` | Обрезанное имя свойства иначе не прочитать |

## Доступность

| | |
|---|---|
| Имя | `aria-label`, включающий ось **и** свойство. «X» именем не является |
| Ось | `inst-num-axis` не размечена как подпись намеренно: иначе она перебила бы `aria-label` и стала бы полным именем поля |
| Клавиатура | Стрелки вверх и вниз меняют значение на `step`, `PageUp` и `PageDown` — на десять шагов. Всё нативное: степперы сняты оформлением, поведение не тронуто |
| Фокус | Кольцо на обёртке через `:focus-within` — фокус виден целиком, а не только вокруг цифр |
| Контраст | Рамка берёт `--border-control` и держит 3:1: здесь граница и есть контрол. Буква оси — `--text-muted`, вспомогательная |
| Цель нажатия | Высота `--control-h-sm` — 26px, в плотности `compact` 22px, **ниже 24px по WCAG 2.5.8**: это осознанная цена панели инспектора, где на экране десятки полей |
| Режим принудительных цветов | Обёртка получает рамку `CanvasText`: без неё поле схлопывается с поверхностью панели |

## Связанное

[Текстовое поле](./input.md) · [Слайдер](./slider.md) · [Селект](./select.md) ·
[Форма](./form.md) · [Панель](../display/panel.md) ·
[Мера](../charts/meter.md)
