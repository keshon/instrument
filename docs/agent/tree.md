---
title: Дерево
group: Агентный слой
layout: component
source: src/text.css
js: Ленивая подгрузка детей — слой приложения. Стрелки, раскрытие и бегущий tabindex делает `instrument.js`
api:
  - { name: "inst-tree", kind: "класс", doc: "Контейнер" }
  - { name: "inst-tree--guides", kind: "модификатор", doc: "Вертикальные направляющие вложенности. Нужны от четырёх уровней; на двух были бы шумом" }
  - { name: "inst-tree-item", kind: "класс", doc: "Узел" }
  - { name: "inst-tree-twist", kind: "класс", doc: "Треугольник раскрытия. Поворачивается по `aria-expanded`; щелчок по нему раскрывает узел и не переносит выбор" }
  - { name: "aria-expanded", kind: "атрибут", doc: "`true` · `false`. Только у узлов с детьми" }
  - { name: "aria-selected", kind: "атрибут", doc: "`true` · `false`" }
  - { name: "aria-level", kind: "атрибут", doc: "целое, с 1" }
  - { name: "--depth", kind: "переменная", value: "0" }
  - { name: "--size-indent", kind: "токен" }
  - { name: "--control-h-sm", kind: "токен" }
  - { name: "--space-3", kind: "токен" }
  - { name: "--radius-sm", kind: "токен" }
  - { name: "--text-sm", kind: "токен" }
  - { name: "--surface-hover", kind: "токен" }
  - { name: "--surface-selected", kind: "токен" }
  - { name: "--size-chevron", kind: "токен" }
title-en: "Tree"
group-en: "Agent layer"
---

Иерархия: файлы проекта, структура сцены, вложенность вызовов.

В примере четыре уровня и `inst-tree--guides` — вертикальные направляющие
вложенности. От четырёх уровней глаз теряет, к какой ветке относится строка;
на двух те же линии были бы шумом, и модификатор для того и отдельный.

```html preview
<div class="inst-tree inst-tree--guides" role="tree" aria-label="Файлы проекта">
  <div class="inst-tree-item" role="treeitem" aria-level="1" aria-expanded="true"
       tabindex="0" style="--depth:0">
    <span class="inst-tree-twist"></span>src
  </div>
  <div class="inst-tree-item" role="treeitem" aria-level="2" aria-expanded="true"
       tabindex="-1" style="--depth:1">
    <span class="inst-tree-twist"></span>world
  </div>
  <div class="inst-tree-item" role="treeitem" aria-level="3" aria-expanded="true"
       tabindex="-1" style="--depth:2">
    <span class="inst-tree-twist"></span>terrain
  </div>
  <div class="inst-tree-item" role="treeitem" aria-level="4"
       tabindex="-1" style="--depth:3">heightmap.ts</div>
  <div class="inst-tree-item" role="treeitem" aria-level="4" aria-selected="true"
       tabindex="-1" style="--depth:3">chunks.bin</div>
</div>
```

## Контракт

Состояние приходит из разметки: библиотека его рисует, но не вычисляет.

| Что | Почему |
|---|---|
| `role="tree"` на контейнере | Иначе это набор `<div>` |
| `role="treeitem"` на узле | Носитель роли |
| `aria-level` | Единственный способ озвучить глубину: визуальный отступ скринридеру не виден |
| `aria-expanded` | **Только у узлов с детьми.** На листе он соврёт, что там что-то есть |
| `aria-selected` | Если выбор поддерживается |
| Бегущий `tabindex` | `0` у активного, `-1` у остальных |

### Доступность

| | |
|---|---|
| Клавиатура | `↑` `↓` — по видимым узлам, `→` — раскрыть или войти, `←` — свернуть или к родителю, `Home`/`End`. Поведение целиком на приложении |
| Глубина | Только `aria-level`. Отступ — визуальный признак, для скринридера его нет |
| `aria-expanded` на листе | Частая ошибка: атрибут со значением `false` на узле без детей объявляет несуществующее свёрнутое поддерево |
| Треугольник | Декоративен. Состояние раскрытия несёт `aria-expanded`, а не поворот значка |
| Цель нажатия | Высота — `--control-h-sm`, то есть 26px. В плотности `compact` — 22px, **ниже 24px по WCAG 2.5.8** |
| Виртуализация | При виртуализации обязателен `aria-setsize` и `aria-posinset`: иначе скринридер объявит «3 из 3», когда узлов три тысячи |

## Устройство

```html
<div class="inst-tree-item" style="--depth:2">…</div>
```

Отступ считается из `--depth` умножением на `--size-indent`, поэтому он
реагирует на плотность контейнера. Глубина задаётся инлайном как **значение** —
тот же разрешённый канал, что у ширины заполнения меры и границ отрезка
дорожки.

Плоский список узлов с `--depth`, а не вложенные контейнеры: так строку можно
виртуализировать, а бегущий `tabindex` не должен обходить дерево рекурсивно.

## JS

Подключите модуль один раз на страницу — инициализировать компоненты по
отдельности не нужно, `instrument.js` работает делегированием и видит узлы, пришедшие
позже.

```html
<script type="module" src="instrument.js"></script>
```

### Что делает `instrument.js`

Пример в шапке живой: `↓` и `↑` по узлам, `→` раскрывает, `←` сворачивает и уводит к родителю.

Контракт роли `treeitem` целиком: стрелки вверх и вниз по узлам, `→`
раскрывает свёрнутый и уходит внутрь раскрытого, `←` сворачивает и уводит к
родителю. Родитель ищется по `aria-level`, поэтому вложенность в разметке
может быть плоской — для дерева на десять тысяч узлов это единственный рабочий
вариант.

Раскрытие меняет `aria-expanded`; как оно выглядит, решает CSS.

### События

`inst:select` — на выбранном узле, `{ value }` в `detail`.

```js
tree.addEventListener('inst:select', (e) => open(e.detail.value));
```

Раскрытие: `←` и `→` с клавиатуры, щелчок по завитку мышью. Щелчок по самой
строке выбирает узел и раскрытия не трогает — завиток и строка две разные
цели, и то, что они соседи, не повод склеивать их смысл.

```js
tree.addEventListener('inst:expand', (e) => {
  if (!e.detail.open) return;
  e.preventDefault();          // строки снимет и покажет приложение
  loadChildren(e.target);
});
```

### Что остаётся приложению

Ленивая подгрузка детей. `instrument.js` переключил `aria-expanded` — приложение решает,
надо ли идти за данными. Событие `inst:expand` отменяемо: приложение, которое
рисует дерево из данных, снимет и покажет строки само.

## API

```api
```

## Связанное

```related
```
