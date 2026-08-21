---
title: Секция
group: Раскладка
layout: component
source: src/layout.css
api:
  - { name: "inst-section", kind: "класс", doc: "Колонка с зазором `--gap-row`" }
  - { name: "inst-section-head", kind: "класс", doc: "Строка заголовка. Выравнивание по базовой линии" }
  - { name: "inst-section-title", kind: "класс", doc: "Подпись секции" }
  - { name: "inst-section-actions", kind: "класс", doc: "Хвост строки заголовка, прижатый к дальнему краю" }
  - { name: "--gap-row", kind: "токен" }
  - { name: "--gap-inline", kind: "токен" }
  - { name: "--text-sm", kind: "токен" }
  - { name: "--weight-medium", kind: "токен" }
  - { name: "--text-secondary", kind: "токен" }
title-en: "Section"
group-en: "Layout"
---

Именованный блок внутри экрана: подпись, необязательные действия и содержимое.
Подпись **тише данных**, которые она подписывает.

```html preview
<div class="inst-section">
  <div class="inst-section-head">
    <h2 class="inst-section-title">Сводка за сутки</h2>
    <span class="inst-section-actions">
      <button class="inst-btn inst-btn--sm inst-btn--ghost" type="button">Период</button></span>
  </div>
  <div class="inst-grid inst-grid--tight">
    <div class="inst-metric"><div class="inst-metric-label">В работе</div><div class="inst-metric-value">7</div></div>
    <div class="inst-metric"><div class="inst-metric-label">В очереди</div><div class="inst-metric-value">5</div></div>
    <div class="inst-metric"><div class="inst-metric-label">Упало</div><div class="inst-metric-value">1</div></div>
  </div>
</div>
```

## Контракт

| Что | Обязательно | Почему |
|---|---|---|
| Настоящий заголовок нужного уровня | да | `inst-section-title` оформляет, но не объявляет: навигация по заголовкам без него не работает |
| `<section>` или `<div>` | нет | Тег роли не меняет: секция ничего не рисует, кроме ритма |
| `inst-section-head` | нет | Нужна там, где есть заголовок или действия |
| Настоящий `<h2>` или `<h3>` внутри секции | да | `inst-section-title` оформляет, но не объявляет: без заголовка навигация по заголовкам секцию не увидит |
| `inst-section-actions` — только для этого блока | да | Действия уровня экрана живут в [шапке экрана](./page-header.md) |

Секция ничего не рисует, кроме ритма: это стопка с зазором `--gap-row` и
строкой заголовка. Рамка вокруг четырёх метрик превратила бы одну группу в
четыре объекта.

| Что | Почему |
|---|---|
| Тег заголовка по уровню (`<h2>`, `<h3>`) | Класс даёт кегль, а не уровень. Скринридер строит оглавление по тегам |
| `inst-section-head` как обёртка | Без неё заголовок и действия становятся детьми колонки и встают друг под друга |
| `type="button"` у действий | Иначе внутри формы они её отправят |

### Доступность

| | |
|---|---|
| Уровень заголовка | Выбирается по структуре документа, а не по внешнему виду. Автоматически он не назначается |
| Ориентир | Если блок действительно самостоятельный раздел, ставьте `<section>` с `aria-labelledby` на заголовок; `<div>` ориентира не даёт, и это нормально для группировки |
| Контраст | Подпись — `--text-secondary`, 4.5:1 в пяти темах. Она тише данных по цвету, но не ниже порога чтения |
| Кегль | `--text-sm`, тот же, что у базы документа: подпись не крупнее содержимого |
| Порядок | Действия идут в разметке после заголовка и обходятся после него |

## Устройство

```html preview
<div class="inst-section">
  <div class="inst-section-head">
    <h2 class="inst-section-title">Активные прогоны</h2>
    <span class="inst-badge inst-nav-count">7</span>
    <span class="inst-section-actions">
      <button class="inst-btn inst-btn--sm inst-btn--ghost" type="button">Все</button></span>
  </div>
  <div class="inst-stack inst-stack--tight">
    <div class="inst-card"><div class="inst-card-title">worldgen · #4127</div></div>
    <div class="inst-card"><div class="inst-card-title">worldgen · #4126</div></div>
  </div>
</div>
```

`inst-section-head` выравнивает детей **по базовой линии**: заголовок, счётчик
и кнопка стоят на одной строке текста, а не по центру своих коробок. Поэтому в
строку заголовка можно класть что угодно, не подгоняя высоты.

`inst-section-actions` прижимается к дальнему краю `margin-inline-start: auto` —
тем же приёмом, что `inst-cluster-spacer` и `inst-page-actions`. Одна ось —
одна механика.

### Заголовок

`inst-section-title` — `--text-sm`, `--weight-medium`, `--text-secondary`. То
есть **тише содержимого**: подпись обслуживает данные. Капса и
трекинга здесь нет по той же причине, что и везде в библиотеке, а веса 700 в
библиотеке не существует вовсе.

Класс задаёт кегль, но не уровень. Уровень выбирается по месту в структуре
документа: под `<h1>` шапки экрана обычно `<h2>`.

## Композиции

### Несколько секций на экране

```html
<div class="inst-container inst-stack inst-stack--loose">
  <div class="inst-page-header">…</div>
  <div class="inst-section">…</div>
  <div class="inst-section">…</div>
</div>
```

Зазор между секциями ставит внешняя стопка, а не сами секции: зазор — между
соседями, а не вокруг каждого.

### Секция без действий

```html
<div class="inst-section">
  <div class="inst-section-head">
    <h2 class="inst-section-title">Параметры прогона</h2>
  </div>
  <div class="inst-kv">…</div>
</div>
```

## API

```api
```

## Связанное

```related
```
