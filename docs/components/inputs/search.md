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

## Контракт

Три части, и порядок между ними обязателен.

| Что | Обязательно | Почему |
|---|---|---|
| `inst-search` на обёртке | да | Лупа рисуется псевдоэлементом обёртки: на самом `<input>` псевдоэлементов нет |
| `inst-input` на поле | да | Оформление поля целиком приходит оттуда. `inst-search` добавляет только отступ под лупу и очистку. Отступ задан селектором `& > .inst-input` — поле без класса встанет под значок |
| `type="search"` | да | Крестик очистки — **нативный**. Без этого типа его не существует, и это обычное текстовое поле с лупой |
| `<label for>` или `aria-label` | да | `placeholder` подписью не является: он исчезает при вводе |

Лупа несёт `pointer-events: none`: клик по значку попадает в поле, а не в
пустоту перед ним.

Видимая `<label for>` предпочтительна; `aria-label` — когда поле стоит в панели
инструментов и подпись съела бы строку. Тип `search` объявляет поле как поле
поиска сам, и `role="search"` руками ставить не нужно.

Лупа — псевдоэлемент, и для скринридера её нет вовсе. Это правильно: значок
украшение, а работу называет подпись. Значок берёт `--text-muted` по той же
причине, крестик темнеет на наведении до `--text-primary`.

Очистка — нативная кнопка браузера, клавиатурный эквивалент `Esc` тоже
нативный. Крестик размером `--size-chevron`, 10px, но единственным способом
стереть значение он не является: есть выделение и `Backspace`.

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
(`::-webkit-search-cancel-button`), приведённая к тону библиотеки: тот же глиф и тот
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

## API

```api
```

## Связанное

```related
```
