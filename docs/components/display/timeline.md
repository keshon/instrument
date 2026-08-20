---
title: Лента событий
group: Отображение данных
layout: component
shape: 2
source: src/data.css
api:
  - { name: "inst-timeline", kind: "класс", doc: "Контейнер" }
  - { name: "inst-timeline-item", kind: "класс", doc: "Одно событие" }
  - { name: "inst-timeline-dot", kind: "класс", doc: "Точка на линии" }
  - { name: "inst-timeline-body", kind: "класс", doc: "Описание события" }
  - { name: "inst-timeline-sub", kind: "класс", doc: "Подробность под описанием" }
  - { name: "inst-timeline-time", kind: "класс", doc: "Отметка времени, прижата к концу" }
  - { name: "data-tone", kind: "атрибут", doc: "`neutral` `running` `ok` `warn` `error`" }
  - { name: "--size-gutter", kind: "токен" }
  - { name: "--gap-inline", kind: "токен" }
  - { name: "--row-pad-y", kind: "токен" }
  - { name: "--size-dot", kind: "токен" }
  - { name: "--size-marker", kind: "токен" }
  - { name: "--border", kind: "токен" }
  - { name: "--text-sm", kind: "токен" }
  - { name: "--text-xs", kind: "токен" }
  - { name: "--text-muted", kind: "токен" }
  - { name: "--tone-mark", kind: "токен" }
title-en: "Timeline"
group-en: "Data display"
---

Хронология: что произошло и когда. Точки соединены линией, поэтому лента
читается как последовательность, а не как список.

```html preview
<div class="inst-timeline">
  <div class="inst-timeline-item">
    <span class="inst-timeline-dot" data-tone="ok"></span>
    <div class="inst-timeline-body">
      Прогон завершён
      <div class="inst-timeline-sub">4 файла изменено</div>
    </div>
    <span class="inst-timeline-time">14:36</span>
  </div>
  <div class="inst-timeline-item">
    <span class="inst-timeline-dot" data-tone="running"></span>
    <div class="inst-timeline-body">
      Индексация
      <div class="inst-timeline-sub">terrain/</div>
    </div>
    <span class="inst-timeline-time">14:32</span>
  </div>
</div>
```

## Контракт

| Что | Обязательно | Почему |
|---|---|---|
| `<ol>` и `<li>` | да | Порядок здесь несёт смысл, и скринридер должен объявить «список из N» |
| `<time datetime>` | да | «14:36» без даты неоднозначно вне контекста дня |
| Тон дублируется словом в описании | да | Точка сама по себе ничего не сообщает тому, кто не различает тон |
| `data-tone` | нет | Без атрибута точка нейтральна |

Линия рисуется псевдоэлементом и в озвучивание не попадает — она декоративна.

### Доступность

| | |
|---|---|
| Это список | Оберните в `<ol>`, а событие — в `<li>`: порядок здесь несёт смысл, и скринридер должен объявить «список из N» |
| Время | `<time datetime="…">` с машинным значением. «14:36» без даты неоднозначно вне контекста дня |
| Не только цвет | Тон точки дублируется словом в описании события. Точка сама по себе ничего не сообщает тому, кто не различает тон |
| Линия | Рисуется псевдоэлементом и в озвучивание не попадает — это правильно, она декоративна |
| Уменьшенное движение | Пульсация точки `running` замедляется, а не гаснет |

## Варианты

Точка берёт тот же `data-tone`, что и вся библиотека.

```html preview
<div class="inst-timeline">
  <div class="inst-timeline-item">
    <span class="inst-timeline-dot"></span>
    <div class="inst-timeline-body">Рядовое событие</div>
  </div>
  <div class="inst-timeline-item">
    <span class="inst-timeline-dot" data-tone="running"></span>
    <div class="inst-timeline-body">Происходит сейчас</div>
  </div>
  <div class="inst-timeline-item">
    <span class="inst-timeline-dot" data-tone="ok"></span>
    <div class="inst-timeline-body">Успешно завершено</div>
  </div>
  <div class="inst-timeline-item">
    <span class="inst-timeline-dot" data-tone="warn"></span>
    <div class="inst-timeline-body">С замечаниями</div>
  </div>
  <div class="inst-timeline-item">
    <span class="inst-timeline-dot" data-tone="error"></span>
    <div class="inst-timeline-body">Отказ</div>
  </div>
</div>
```

| Тон | Когда |
|---|---|
| без атрибута | Рядовое событие |
| `data-tone="running"` | Происходит сейчас. Точка пульсирует |
| `data-tone="ok"` | Успешное завершение |
| `data-tone="warn"` | С замечаниями |
| `data-tone="error"` | Отказ |

## Правила

:::do Лента на ol и li
Порядок несёт смысл. Скринридер объявляет «список из N» и номер позиции.
:::

:::dont Точка как единственный носитель тона
Отказ обязан быть назван словом в описании события, а не только красной
точкой.
:::

:::do time с машинным значением
«14:36» без даты неоднозначно вне контекста дня.
:::

:::dont Лента для однородного потока
Строки одного формата от процесса — [лог](../../agent/log.md): у него колонки
и моноширинный набор.
:::

## API

```api
```

## Связанное

```related
```
