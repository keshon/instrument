---
title: Диф
group: Агентный слой
status: stable
source: src/agent.css
api:
  - { name: "inst-diff", kind: "класс", doc: "Контейнер" }
  - { name: "inst-diff-head", kind: "класс", doc: "Шапка: путь и статистика" }
  - { name: "inst-diff-path", kind: "класс", doc: "Путь к файлу. Обрезается с конца" }
  - { name: "inst-diff-stat", kind: "класс", doc: "`+N` / `−M`" }
  - { name: "inst-diff-body", kind: "класс", doc: "Тело, прокручивается" }
  - { name: "inst-diff-line", kind: "класс", doc: "Строка" }
  - { name: "inst-diff-num", kind: "класс", doc: "Номер строки. Не выделяется мышью" }
  - { name: "inst-diff-code", kind: "класс", doc: "Сам код" }
  - { name: "inst-diff-fold", kind: "класс", doc: "Свёрнутый участок" }
  - { name: "data-kind", kind: "атрибут", value: "add · del", doc: "на `inst-diff-line` и на числах в `inst-diff-stat`" }
  - { name: "--font-mono", kind: "токен" }
  - { name: "--text-xs", kind: "токен" }
  - { name: "--ok-bg", kind: "токен" }
  - { name: "--ok-text", kind: "токен" }
  - { name: "--err-bg", kind: "токен" }
  - { name: "--err-text", kind: "токен" }
  - { name: "--border-subtle", kind: "токен" }
  - { name: "--hairline", kind: "токен" }
  - { name: "--radius-lg", kind: "токен" }
  - { name: "--pad-cell-x", kind: "токен" }
  - { name: "--text-muted", kind: "токен" }
---

Основной артефакт работы агента: что именно он изменил в файле.

```html preview
<div class="inst-diff">
  <div class="inst-diff-head">
    <span class="inst-diff-path">terrain/heightmap.ts</span>
    <span class="inst-diff-stat"><span data-kind="add">+3</span> <span data-kind="del">−2</span></span>
  </div>
  <div class="inst-diff-body">
    <div class="inst-diff-line"><span class="inst-diff-num">41</span><span class="inst-diff-code">const size = 256;</span></div>
    <div class="inst-diff-line" data-kind="del"><span class="inst-diff-num">42</span><span class="inst-diff-code">let seed = 0;</span></div>
    <div class="inst-diff-line" data-kind="add"><span class="inst-diff-num">42</span><span class="inst-diff-code">let seed = Date.now();</span></div>
    <div class="inst-diff-fold">— пропущено 18 строк —</div>
  </div>
</div>
```

## Когда использовать

| Используйте | Возьмите другое |
|---|---|
| Показать, что агент изменил в файле | **Показать файл целиком** — [блок кода](../components/display/code.md): диф без изменений это просто код |
| Изменения нужно принять или отклонить | **Перечислить затронутые файлы без содержания** — список в [запросе подтверждения](./approval.md) |
| Изменений в файле десятки строк | **Изменений сотни файлов** — [таблица](../components/display/table.md) со счётчиками, диф каждого по требованию |
| — | **Поток строк от машины** — [лог](./log.md) |

## Знак в первой колонке — не украшение

Строка добавления несёт `+`, строка удаления `−`. Это выполнение закона кита:
**цвет не имеет права быть единственным носителем состояния**. Диф, где
изменения различаются только зелёным и красным фоном, для дальтоника — и на
чёрно-белой печати — превращается в сплошной текст.

Знак рисуется псевдоэлементом из `data-kind`, поэтому он не попадает в
выделение и не копируется вместе с кодом.

## `data-kind` — категория, а не тон

```html
<div class="inst-diff-line" data-kind="add">…</div>
```

| Значение | Что значит |
|---|---|
| нет атрибута | Строка контекста, не изменялась |
| `add` | Добавлена |
| `del` | Удалена |

Добавленная строка **не находится в состоянии «ok»** — она относится к виду
«добавление». Поэтому здесь свой атрибут, а не `data-tone`: тон соврал бы про
смысл. Тот же `data-kind` несут и числа в шапке дифа — одна ось, одна запись.

## Свёрнутые участки

```html
<div class="inst-diff-fold">— пропущено 18 строк —</div>
```

Как и свёрнутый вывод шага, участок **называет своё число словами**. Молча
пропустить строки — то же враньё про объём.

## Собранный экран

Диф с шапкой, статистикой и свёрнутым участком.

Что проверить руками: **выделите несколько строк и скопируйте**. В буфер
уйдёт чистый код — без номеров строк и без знаков `+` и `−`. Номера несут
`user-select: none`, знак рисуется псевдоэлементом. Диф, из которого нельзя
скопировать код, бесполезен ровно в тот момент, когда он нужен.

```html preview context
<div class="inst-diff">
  <div class="inst-diff-head">
    <span class="inst-diff-path">src/world/biomes.ts</span>
    <span class="inst-diff-stat"><span data-kind="add">+3</span> <span data-kind="del">−2</span></span>
  </div>
  <div class="inst-diff-body">
    <div class="inst-diff-fold">свёрнуто 118 строк</div>
    <div class="inst-diff-line"><span class="inst-diff-num">119</span><span class="inst-diff-num">119</span><span class="inst-diff-code">  const seed = opts.seed ?? 0;</span></div>
    <div class="inst-diff-line" data-kind="del"><span class="inst-diff-num">120</span><span class="inst-diff-num"></span><span class="inst-diff-code">  const noise = simplex2(seed);</span></div>
    <div class="inst-diff-line" data-kind="del"><span class="inst-diff-num">121</span><span class="inst-diff-num"></span><span class="inst-diff-code">  if (noise &gt; 0.5) return TUNDRA;</span></div>
    <div class="inst-diff-line" data-kind="add"><span class="inst-diff-num"></span><span class="inst-diff-num">120</span><span class="inst-diff-code">  const noise = simplex2(seed, { octaves: 3 });</span></div>
    <div class="inst-diff-line" data-kind="add"><span class="inst-diff-num"></span><span class="inst-diff-num">121</span><span class="inst-diff-code">  // тундра и степь перекрывались в 12 клетках — разводим по высоте</span></div>
    <div class="inst-diff-line" data-kind="add"><span class="inst-diff-num"></span><span class="inst-diff-num">122</span><span class="inst-diff-code">  if (noise &gt; 0.5 &amp;&amp; height &gt; 0.72) return TUNDRA;</span></div>
    <div class="inst-diff-line"><span class="inst-diff-num">122</span><span class="inst-diff-num">123</span><span class="inst-diff-code">  return STEPPE;</span></div>
  </div>
</div>
```

## Справочник

```api
```

## Доступность

| | |
|---|---|
| Не только цвет | Знак `+`/`−` в первой колонке. Это главное требование компонента |
| Номера строк | `user-select: none` — при копировании фрагмента номера не попадают в буфер вместе с кодом |
| Копирование | Знак — псевдоэлемент, поэтому копируется чистый код, а не `+const size` |
| Путь к файлу | Обрезается с конца многоточием. Полный путь обязан оставаться в `title` — иначе диф безымянен |
| Прокрутка | Тело прокручивается по горизонтали и должно быть достижимо с клавиатуры |
| Печать | Тонированные фоны в `@media print` заменяются на печатаемые, знак остаётся: на бумаге тон исчезает вовсе |

## Связанное

[Шаг](./step.md) · [Запрос подтверждения](./approval.md) ·
[Код](../components/display/code.md) · [Лог](./log.md)
