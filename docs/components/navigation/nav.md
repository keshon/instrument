---
title: Навигация
group: Навигация
layout: component
source: src/layout.css
api:
  - { name: "inst-nav", kind: "класс", doc: "Контейнер списка. Колонка с зазором `--space-1`" }
  - { name: "inst-nav-group", kind: "класс", doc: "Смысловой блок пунктов" }
  - { name: "inst-nav-label", kind: "класс", doc: "Подпись группы" }
  - { name: "inst-nav-item", kind: "класс", doc: "Пункт. Ссылка, а не кнопка" }
  - { name: "inst-nav-count", kind: "класс", doc: "Прижать счётчик к дальнему краю пункта" }
  - { name: "--control-h-md", kind: "токен" }
  - { name: "--radius-sm", kind: "токен" }
  - { name: "--text-sm", kind: "токен" }
  - { name: "--text-2xs", kind: "токен" }
  - { name: "--weight-medium", kind: "токен" }
  - { name: "--space-1", kind: "токен" }
  - { name: "--space-3", kind: "токен" }
  - { name: "--gap-inline", kind: "токен" }
  - { name: "--pad-panel", kind: "токен" }
  - { name: "--size-marker", kind: "токен" }
  - { name: "--surface-hover", kind: "токен" }
  - { name: "--surface-selected", kind: "токен" }
  - { name: "--accent-solid", kind: "токен" }
  - { name: "--accent-text", kind: "токен" }
  - { name: "--text-secondary", kind: "токен" }
  - { name: "--text-muted", kind: "токен" }
  - { name: "--dur-1", kind: "токен" }
  - { name: "--ease-out", kind: "токен" }
title-en: "Navigation"
group-en: "Navigation"
---

Список разделов приложения в боковой колонке. Текущий пункт помечен не только
цветом: у его края стоит несущая полоса.

```html preview
<nav class="inst-nav" aria-label="Разделы">
  <div class="inst-nav-group">
    <span class="inst-nav-label">Работа</span>
    <a class="inst-nav-item" href="#" aria-current="page">
      <svg class="inst-icon" aria-hidden="true"><use href="#i-grid"/></svg>Прогоны
      <span class="inst-badge inst-nav-count">7</span></a>
    <a class="inst-nav-item" href="#">
      <svg class="inst-icon" aria-hidden="true"><use href="#i-list"/></svg>Очередь</a>
    <a class="inst-nav-item" href="#">
      <svg class="inst-icon" aria-hidden="true"><use href="#i-chart"/></svg>Отчёты</a>
  </div>
  <div class="inst-nav-group">
    <span class="inst-nav-label">Настройка</span>
    <a class="inst-nav-item" href="#">
      <svg class="inst-icon" aria-hidden="true"><use href="#i-users"/></svg>Доступы</a>
    <a class="inst-nav-item" href="#">
      <svg class="inst-icon" aria-hidden="true"><use href="#i-settings"/></svg>Параметры</a>
  </div>
</nav>
```

## Контракт

| Что | Обязательно | Почему |
|---|---|---|
| `<nav>` + `aria-label` | да | На экране больше одной навигации (боковая, крошки, вкладки, страницы) — без имени они неразличимы в списке ориентиров |
| `<a href>` у пункта | да | Навигация — это адрес. Кнопка ломает средний клик, «открыть в новой вкладке» и копирование ссылки |
| `aria-current="page"` | да | Единственный источник текущего состояния. Класса для него нет |
| `aria-hidden="true"` у иконки | да | Имя пункта даёт подпись |

### Доступность

| | |
|---|---|
| Клавиатура | Всё от нативных ссылок: `Tab` — обход, `Enter` — переход. `instrument.js` ничего не перехватывает, бегущий `tabindex` навигации не нужен |
| Ориентир | `<nav aria-label>` — пункт назначения для «перейти к ориентиру» в скринридере |
| Цвет не единственный носитель | Текущий пункт несёт полосу у края и `aria-current`, а не один только оттенок фона |
| Цель нажатия | Высота пункта — `--control-h-md`, 32px. В плотности `compact` она уменьшается: применяйте `compact` к табличным видам, а не к основной навигации |
| Контраст | Подпись `--text-secondary` держит 4.5:1 на приглушённой поверхности колонки в пяти темах |
| Печать | Боковая колонка не печатается: это хром приложения, а не данные |

## Устройство

`inst-nav-group` собирает пункты в смысловой блок, `inst-nav-label` его
подписывает. Зазор между группами ставится **между соседями**
(`.inst-nav-group + .inst-nav-group`), а не вокруг каждой, — иначе первая
группа отступила бы от верхнего края колонки вдвое дальше, чем от боковых.

Подпись — обычное предложение кеглем `--text-2xs` цветом `--text-muted`. Капса
здесь нет: подпись не имеет права звучать громче пунктов, которые она
подписывает.

## Состояния

```html preview
<nav class="inst-nav" aria-label="Пример текущего пункта">
  <div class="inst-nav-group">
    <a class="inst-nav-item" href="#" aria-current="page">Прогоны</a>
    <a class="inst-nav-item" href="#">Очередь</a>
  </div>
</nav>
```

Метка ставится атрибутом `aria-current="page"`, а не классом. Оформление ловит
селектор `[aria-current]:not([aria-current="false"])`, поэтому строчка
`aria-current="false"` у остальных пунктов допустима и ничего не подсвечивает.

| Признак текущего | Что даёт |
|---|---|
| Полоса у внутреннего края | `box-shadow: inset` шириной `--size-marker`. Носитель, не зависящий от цвета |
| Поверхность | `--surface-selected` |
| Цвет подписи | `--accent-text` |
| Начертание | `--weight-medium`. Веса 700 в библиотеке нет |

## Поведение

### В узкой оболочке

Внутри `.inst-shell-aside` ниже 60rem колонка разворачивается в горизонтальную
полосу: `.inst-nav` и `.inst-nav-group` встают в ряд, `.inst-nav-label`
скрывается, подписи пунктов перестают переноситься. Это единственное место, где
навигация меняет форму, и делает это [оболочка](../../layout/shell.md), а не
сама навигация.

## Композиции

### Счётчик в пункте

```html preview
<nav class="inst-nav" aria-label="Пример со счётчиком">
  <div class="inst-nav-group">
    <a class="inst-nav-item" href="#">Очередь
      <span class="inst-badge inst-nav-count">7</span></a>
    <a class="inst-nav-item" href="#">История</a>
  </div>
</nav>
```

`inst-nav-count` — модификатор положения: `margin-inline-start: auto` прижимает
[бейдж](../display/badge.md) к дальнему краю пункта. Сам бейдж остаётся
бейджем и читает те же токены.

## Правила

:::do Пункт — ссылка
Навигация ведёт на адрес. Средний клик, «открыть в новой вкладке» и копирование
работают бесплатно.
:::

:::dont Пункт — кнопка с обработчиком
Она ломает всё перечисленное и требует роли, которой у неё нет.
:::

:::do Текущий пункт полосой у края
Носитель, не зависящий от цвета. `aria-current` сообщает то же самое
скринридеру.
:::

:::dont Капс в подписи группы
Подпись не имеет права звучать громче пунктов, которые она подписывает.
Различают её кегль и цвет.
:::

## API

```api
```

## Связанное

```related
```
