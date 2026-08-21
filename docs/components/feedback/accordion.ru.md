---
title: Аккордеон
group: Обратная связь
layout: component
source: src/feedback.css
api:
  - { name: "inst-accordion", kind: "класс", doc: "Контейнер. Колонка секций с зазором `--space-1`" }
  - { name: "inst-accordion-item", kind: "класс", doc: "Секция. Рамка, радиус, `overflow: hidden`" }
  - { name: "inst-accordion-head", kind: "класс", doc: "Заголовок. Высота от `--control-h-lg`, шеврон у дальнего края" }
  - { name: "inst-accordion-body", kind: "класс", doc: "Содержимое" }
  - { name: "inst-accordion-item--flush", kind: "модификатор", doc: "Секция ВНУТРИ чужой рамки: панели, карточки, поля примера. Своей рамки и радиуса нет, остаётся шов сверху" }
  - { name: "--space-1", kind: "токен" }
  - { name: "--space-3", kind: "токен" }
  - { name: "--pad-cell-x", kind: "токен" }
  - { name: "--pad-panel", kind: "токен" }
  - { name: "--gap-inline", kind: "токен" }
  - { name: "--control-h-lg", kind: "токен" }
  - { name: "--radius-md", kind: "токен" }
  - { name: "--hairline", kind: "токен" }
  - { name: "--border", kind: "токен" }
  - { name: "--border-subtle", kind: "токен" }
  - { name: "--surface-raised", kind: "токен" }
  - { name: "--surface-hover", kind: "токен" }
  - { name: "--size-chevron", kind: "токен" }
  - { name: "--text-muted", kind: "токен" }
  - { name: "--text-sm", kind: "токен" }
  - { name: "--weight-medium", kind: "токен" }
  - { name: "--dur-1", kind: "токен" }
  - { name: "--ease-out", kind: "токен" }
---

Группа сворачиваемых разделов, из которых открыт ровно один. На нативных
`<details>` с общим `name`, поэтому взаимное исключение, раскрытие и
клавиатура достаются от платформы — без единой строки JS.

```html preview
<div class="inst-accordion">
  <details class="inst-accordion-item" name="faq" open>
    <summary class="inst-accordion-head">Почему прогон занял 42 секунды?</summary>
    <div class="inst-accordion-body">Три прохода расстановки ресурсов плюс запекание навмеша.
      Основное время ушло на валидацию путей, которая перезапускалась дважды.</div>
  </details>
  <details class="inst-accordion-item" name="faq">
    <summary class="inst-accordion-head">Можно ли ограничить расход токенов?</summary>
    <div class="inst-accordion-body">Да, в параметрах проекта. Лимит проверяется перед стартом
      каждого агента, а не по факту — прогон не уйдёт в минус.</div>
  </details>
  <details class="inst-accordion-item" name="faq">
    <summary class="inst-accordion-head">Что происходит с артефактами после удаления?</summary>
    <div class="inst-accordion-body">Удаляются вместе с прогоном. Резервная копия хранится семь
      суток, восстановление — через поддержку.</div>
  </details>
</div>
```

## Контракт

| Что | Обязательно | Почему |
|---|---|---|
| `<details>` под `inst-accordion-item` | да | Раскрытие, клавиатура и состояние `open` — от платформы |
| `<summary>` под `inst-accordion-head` | да | Только `<summary>` фокусируем и переключает `<details>` |
| Общий `name` | да, для аккордеона | Делает группу взаимоисключающей. Без него секции независимы |
| Прямое вложение | да | Правила открытой секции записаны через `>`: заголовок обязан быть прямым потомком `<details>` |
| Настоящий заголовок внутри `<summary>` | нет, но обычно да | `<summary>` сам по себе в оглавление страницы не попадает |

### Доступность

| | |
|---|---|
| Клавиатура | `Tab` до заголовка, `Enter` или `Space` — переключение. Всё нативное, `instrument.js` ничего не перехватывает |
| Роль и состояние | `<summary>` сам объявляет себя переключателем и сообщает, раскрыт ли раздел. Своего `aria-expanded` ставить не надо и не нужно |
| Поиск по странице | Содержимое закрытой секции находится браузерным поиском, и секция раскрывается сама — это поведение `<details>`, ради него нативный элемент и взят |
| Заголовок | Если секции — разделы страницы, вложите в `<summary>` настоящий заголовок нужного уровня: `<summary>` сам по себе в оглавление не попадает |
| Цель нажатия | Заголовок нажимается целиком, от края до края: `--control-h-lg` в высоту |
| Уменьшенное движение | Поворот шеврона схлопывается вместе со всеми переходами |
| Печать | Закрытые секции **раскрываются на листе**: закрытый `<details>` на бумаге — невосстановимо потерянное содержимое |

## Устройство

| Элемент | Класс | Тег |
|---|---|---|
| Контейнер группы | `inst-accordion` | `<div>` |
| Секция | `inst-accordion-item` | `<details>` |
| Заголовок-переключатель | `inst-accordion-head` | `<summary>` |
| Содержимое | `inst-accordion-body` | `<div>` |

Шеврон рисуется псевдоэлементом заголовка и поворачивается на 180° при
открытии — своего класса у него нет и разметки он не требует. Нативный
треугольник `<summary>` снят.

Волосок под заголовком появляется **только у открытой секции**
(`.inst-accordion-item[open] > .inst-accordion-head`): у закрытой отделять
нечего.

### Внутри чужой рамки

Раскрывашка часто стоит не сама по себе, а внизу панели или карточки: «ещё
поля», «разметка», «журнал». Своя рамка была бы там второй в пикселе от первой,
а своё скругление — скруглением внутри скругления, ни с чем не концентричным.

```html preview context
<div class="inst-panel" style="inline-size:26rem">
  <div class="inst-panel-header"><span class="inst-panel-title">Прогон 4127</span></div>
  <div class="inst-panel-body">Собран за 4,2 с. Тестов 318, упавших нет.</div>
  <details class="inst-accordion-item inst-accordion-item--flush">
    <summary class="inst-accordion-head">Параметры сборки</summary>
    <div class="inst-accordion-body">Ветка master, коммит d7c9cd3, кэш включён.</div>
  </details>
</div>
```

Остаётся один шов сверху: секция **продолжает** то, под чем стоит, а не
объявляет себя отдельным блоком. Контейнер `.inst-accordion` при этом не нужен
— секция здесь одна и живёт в чужой стопке.

## Поведение

### Взаимное исключение

Общий атрибут `name` на элементах `<details>` делает группу
взаимоисключающей: открытие второй секции закрывает первую.

| Что нужно | Как |
|---|---|
| Открыта ровно одна секция | Общий `name` на всех `<details>` |
| Открыто сколько угодно | `name` не ставить |
| Какая открыта при загрузке | `open` на нужном `<details>` |

Уберите `name` — и секции станут независимыми, открывать можно сколько
угодно. Это единственная разница между «аккордеоном» и «списком
раскрывашек», и она в один атрибут.

```html preview
<div class="inst-accordion">
  <details class="inst-accordion-item" open>
    <summary class="inst-accordion-head">Без name — открываются независимо</summary>
    <div class="inst-accordion-body">Эту можно оставить открытой и открыть следующую.</div>
  </details>
  <details class="inst-accordion-item">
    <summary class="inst-accordion-head">Вторая секция</summary>
    <div class="inst-accordion-body">Первая при этом не закроется.</div>
  </details>
</div>
```

## API

```api
```

## Связанное

```related
```
