---
title: Диф
group: Агентный слой
layout: component
source: src/agent.css
api:
  - { name: "inst-diff", kind: "класс", doc: "Контейнер" }
  - { name: "inst-diff--flush", kind: "модификатор", doc: "Без рамки и радиуса: диф занимает панель целиком, и своя рамка легла бы второй линией поверх её границы" }
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
title-en: "Diff"
group-en: "Agent layer"
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

## Контракт

| Что | Обязательно | Почему |
|---|---|---|
| `data-kind` на строке | да | Носитель категории. Строка без него — контекст, и это верное умолчание |
| `inst-diff-num` отдельным узлом | да | Номер не выделяется мышью: скопированный диф иначе приходит вперемешку с номерами строк |
| `inst-diff-code` вокруг кода | да | Только он держит моноширинный набор и не переносит строку |
| `inst-diff-fold` на месте пропуска | нет, но обычно да | Свёрнутый участок — данные о том, что показано не всё |

### Доступность

| | |
|---|---|
| Не только цвет | Знак `+`/`−` в первой колонке. Это главное требование компонента |
| Номера строк | `user-select: none` — при копировании фрагмента номера не попадают в буфер вместе с кодом |
| Копирование | Знак — псевдоэлемент, поэтому копируется чистый код, а не `+const size` |
| Путь к файлу | Обрезается с конца многоточием. Полный путь обязан оставаться в `title` — иначе диф безымянен |
| Прокрутка | Тело прокручивается по горизонтали и должно быть достижимо с клавиатуры |
| Печать | Тонированные фоны в `@media print` заменяются на печатаемые, знак остаётся: на бумаге тон исчезает вовсе |

## Устройство

Строка добавления несёт `+`, строка удаления `−`. Это выполнение закона библиотеки:
**цвет не имеет права быть единственным носителем состояния**. Диф, где
изменения различаются только зелёным и красным фоном, для дальтоника — и на
чёрно-белой печати — превращается в сплошной текст.

Знак рисуется псевдоэлементом из `data-kind`, поэтому он не попадает в
выделение и не копируется вместе с кодом.

## Варианты

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

## Поведение

```html
<div class="inst-diff-fold">— пропущено 18 строк —</div>
```

Как и свёрнутый вывод шага, участок **называет своё число словами**. Молча
пропустить строки — то же враньё про объём.

## Сценарии

Диф с шапкой, статистикой и свёрнутым участком.

Что проверить руками: **выделите несколько строк и скопируйте**. В буфер
уйдёт чистый код — без номеров строк и без знаков `+` и `−`. Номера несут
`user-select: none`, знак рисуется псевдоэлементом. Диф, из которого нельзя
скопировать код, бесполезен в тот момент, когда он нужен.

```html preview context
<div class="inst-diff inst-diff--flush">
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

## Правила

:::do Знак в первой колонке
`+` и `−` печатаются и читаются без цвета. Тон — второй признак, а не
единственный.
:::

:::dont Диф только цветом
Зелёное и красное неразличимы примерно для 8% мужчин и ни для одной
чёрно-белой распечатки.
:::

:::do Номер строки невыделяемым
Скопированный диф иначе приходит вперемешку с номерами и не вставляется никуда.
:::

:::dont Перенос строки в коде
Перенос меняет смысл строки. Тело дифа прокручивается по горизонтали.
:::

## API

```api
```

## Связанное

```related
```
