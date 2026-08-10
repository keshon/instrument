---
title: Поиск
group: Ввод
layout: component
source: src/forms.css
js-optional: Живая фильтрация. Поле и очистка работают без скрипта
api:
  - { name: "inst-search", kind: "класс", doc: "Обёртка. Рисует лупу и отводит под неё место" }
  - { name: "inst-input", kind: "класс", doc: "Само поле. Все размеры и состояния — оттуда" }
  - { name: "--control-pad-sm", kind: "токен" }
  - { name: "--size-icon", kind: "токен" }
  - { name: "--size-chevron", kind: "токен" }
  - { name: "--space-3", kind: "токен" }
  - { name: "--text-muted", kind: "токен" }
  - { name: "--text-primary", kind: "токен" }
title-en: "Search"
group-en: "Inputs"
---

Текстовое поле с лупой и нативной кнопкой очистки. Обёртка вокруг
`input[type=search]`, а не отдельный контрол: своя кнопка потребовала бы JS.

```html preview
<div class="inst-field">
  <label class="inst-label" for="q">Поиск по прогонам</label>
  <span class="inst-search">
    <input class="inst-input" id="q" type="search" placeholder="worldgen">
  </span>
</div>
```

## Использование

Три части, и порядок между ними обязателен.

```html
<span class="inst-search">
  <input class="inst-input" type="search" placeholder="worldgen" aria-label="Поиск по прогонам">
</span>
```

| Что | Обязательно | Почему |
|---|---|---|
| `inst-search` на обёртке | да | Лупа рисуется псевдоэлементом обёртки: на самом `<input>` псевдоэлементов нет |
| `inst-input` на поле | да | Оформление поля целиком приходит оттуда. `inst-search` добавляет только отступ под лупу и очистку. Отступ задан селектором `& > .inst-input` — поле без класса встанет под значок |
| `type="search"` | да | Крестик очистки — **нативный**. Без этого типа его не существует, и это обычное текстовое поле с лупой |
| `<label for>` или `aria-label` | да | `placeholder` подписью не является: он исчезает при вводе |

Лупа несёт `pointer-events: none`: клик по значку попадает в поле, а не в
пустоту перед ним.

## Когда использовать

| Используйте | Возьмите другое |
|---|---|
| Фильтрация или поиск по списку, таблице, логу | **Произвольный текст, который сохраняют** — [текстовое поле](./input.md): у него нет крестика, и очистка ему не нужна |
| Значение, которое пользователь часто стирает целиком | **Выбор из известного набора** — [селект](./select.md): поиск по трём вариантам медленнее выбора |
| Поле в панели инструментов над данными | **Переключение между заранее известными видами** — [сегментированный контрол](../actions/segmented.md) |
| Одно поле — одна строка запроса | **Несколько условий сразу** — это [форма](./form.md) фильтров, а не одно поле |

## Размеры

Все — от [текстового поля](./input.md), потому что `inst-search` не заводит
своих.

```html preview
<span class="inst-search">
  <input class="inst-input inst-input--sm" type="search" value="sm" aria-label="Поиск sm">
</span>
<span class="inst-search">
  <input class="inst-input" type="search" value="md" aria-label="Поиск md">
</span>
<span class="inst-search">
  <input class="inst-input inst-input--lg" type="search" value="lg" aria-label="Поиск lg">
</span>
```

Отступ под лупу считается от `--control-pad-sm` и ширины значка, поэтому у
всех трёх размеров текст начинается на одинаковом расстоянии от значка, а не
наезжает на него в `sm`.

## Состояния

```html preview
<span class="inst-search">
  <input class="inst-input" type="search" value="Обычное" aria-label="Поиск обычный">
</span>
<span class="inst-search">
  <input class="inst-input" type="search" value="Недоступно" disabled aria-label="Поиск недоступен">
</span>
```

Наведение, фокус, ошибка, `readonly` и `disabled` — те же, что у [текстового
поля](./input.md).

## Поведение

### Очистка

Крестик очистки — собственная кнопка браузера
(`::-webkit-search-cancel-button`), приведённая к тону кита: тот же глиф и тот
же приглушённый цвет, что у снятия [тега](../display/tag.md). **Её не прячут и
не подменяют** — поведение «очистить и сообщить об этом форме» уже написано
платформой, и своя кнопка означала бы JS ради того, что уже работает.

:::warn
Плата названа честно: в Firefox нативного крестика нет, и поле там просто без
очистки. Крестик — улучшение, а не единственный способ стереть значение.
:::

## JS

Поле работает без единой строки: `type="search"` даёт нативную кнопку очистки,
`Escape` очищает, форма отправляется по `Enter`. Скрипт нужен, только если
результаты обновляются на лету.

### Фильтрация на лету

```js
input.addEventListener('input', () => {
  const q = input.value.trim().toLowerCase();
  for (const row of rows) {
    row.hidden = q !== '' && !row.dataset.name.includes(q);
  }
  count.textContent = rows.filter((r) => !r.hidden).length;
});
```

Читается `input`, а не `keyup`: очистка крестиком, вставка мышью и
автозаполнение `keyup` не порождают, и список молча оставался бы старым.

### Что сказать вслух

Число найденного обязано попадать в живую область — иначе для скринридера
список меняется беззвучно.

```html
<p class="inst-u-visually-hidden" aria-live="polite">
  <span id="count">12</span> совпадений
</p>
```

## Композиции

### В панели инструментов

```html preview context
<div class="inst-toolbar">
  <span class="inst-search">
    <input class="inst-input inst-input--sm" type="search" placeholder="Фильтр" aria-label="Фильтр по задачам">
  </span>
  <button class="inst-btn inst-btn--sm" type="button">Сбросить</button>
</div>
```

### С подсказкой о синтаксисе

```html preview
<div class="inst-field">
  <label class="inst-label" for="q2">Поиск по логу</label>
  <span class="inst-search">
    <input class="inst-input" id="q2" type="search" aria-describedby="q2-hint">
  </span>
  <span class="inst-field-hint" id="q2-hint">Поддерживаются кавычки и минус для исключения</span>
</div>
```

## Правила

:::do Нативный type="search"
Крестик, роль «поле поиска» и `Esc` для очистки приходят от платформы.
:::

:::dont Своя кнопка очистки
Она означает JS ради того, что браузер уже написал, и теряет клавиатурный
эквивалент.
:::

:::do Подпись даже в тулбаре
Там, где видимая подпись съела бы строку, — `aria-label`. Поле без имени не
называет, по чему оно ищет.
:::

:::dont Поиск по трём вариантам
Выбор из известного короткого набора быстрее делается
[селектом](./select.md) или [сегментированным](../actions/segmented.md).
:::

## Доступность

| | |
|---|---|
| Подпись | Видимая `<label for>` предпочтительна. `aria-label` — когда поле стоит в панели инструментов и подпись съела бы строку |
| Лупа | Псевдоэлемент, для скринридера её нет вовсе — и правильно: это украшение, а работу называет подпись |
| Очистка | Нативная кнопка браузера. Клавиатурный эквивалент — `Esc` в поле поиска, тоже нативный |
| Тип поля | `type="search"` объявляется как «поле поиска»: роль приходит от платформы, `role="search"` руками ставить не нужно |
| Контраст | Значок берёт `--text-muted` — это украшение, а не несущий знак. Крестик темнеет на наведении до `--text-primary` |
| Цель нажатия | Крестик — `--size-chevron`, 10px. Он дублирует выделение и `Backspace`, поэтому не является единственным способом стереть значение |

## API

```api
```

## Связанное

[Текстовое поле](./input.md) [Селект](./select.md) [Форма](./form.md)
[Тег](../display/tag.md) [Таблица](../display/table.md)
[Пустое состояние](../feedback/empty.md)
