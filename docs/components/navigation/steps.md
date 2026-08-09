---
title: Шаги мастера
group: Навигация
layout: component
source: src/layout.css
api:
  - { name: "inst-stepper", kind: "класс", doc: "Список шагов. Ставится на `<ol>`" }
  - { name: "inst-stepper-item", kind: "класс", doc: "Шаг. Носитель `data-state`" }
  - { name: "inst-stepper-name", kind: "класс", doc: "Название шага" }
  - { name: "inst-stepper-note", kind: "класс", doc: "Слово состояния под названием" }
  - { name: "data-state", kind: "атрибут", doc: "`todo` `current` `done`. Закрыт, база — `todo`" }
  - { name: "--gap-inline", kind: "токен" }
  - { name: "--space-2", kind: "токен" }
  - { name: "--space-4", kind: "токен" }
  - { name: "--text-sm", kind: "токен" }
  - { name: "--text-xs", kind: "токен" }
  - { name: "--weight-medium", kind: "токен" }
  - { name: "--size-marker", kind: "токен" }
  - { name: "--radius-full", kind: "токен" }
  - { name: "--track", kind: "токен" }
  - { name: "--accent-mark", kind: "токен" }
  - { name: "--text-primary", kind: "токен" }
  - { name: "--text-secondary", kind: "токен" }
  - { name: "--text-muted", kind: "токен" }
title-en: "Wizard steps"
group-en: "Navigation"
---

Последовательность этапов процесса с отметкой, где он сейчас. Состояние несёт
**слово** под подписью, а не только цвет полосы.

```html preview
<ol class="inst-stepper">
  <li class="inst-stepper-item" data-state="done">
    <span class="inst-stepper-name">Источник</span><span class="inst-stepper-note">готово</span></li>
  <li class="inst-stepper-item" data-state="done">
    <span class="inst-stepper-name">Правила</span><span class="inst-stepper-note">готово</span></li>
  <li class="inst-stepper-item" data-state="current">
    <span class="inst-stepper-name">Проверка</span><span class="inst-stepper-note">сейчас</span></li>
  <li class="inst-stepper-item" data-state="todo">
    <span class="inst-stepper-name">Запуск</span><span class="inst-stepper-note">дальше</span></li>
</ol>
```

## Использование

```html
<ol class="inst-stepper">
  <li class="inst-stepper-item" data-state="current">
    <span class="inst-stepper-name">Проверка</span>
    <span class="inst-stepper-note">сейчас</span>
  </li>
</ol>
```

| Что | Обязательно | Почему |
|---|---|---|
| `<ol>` с `<li>` | да | Шаги упорядочены, и порядок — весь смысл компонента |
| `data-state` на каждом шаге, включая `todo` | да | Отсутствие атрибута неотличимо от опечатки в нём |
| `inst-stepper-note` со словом состояния | да | Цвет полосы не имеет права быть единственным носителем |

## Когда использовать

| Используйте | Возьмите другое |
|---|---|
| Процесс с известным заранее числом этапов, который ведёт человек | **Работа машины с раскрывающимся телом** — [шаг прогона](../../agent/step.md): у него есть вывод, а здесь только подпись |
| Форма, разбитая на экраны | **Свободный переход между видами** — [вкладки](./tabs.md): у шагов есть порядок и незавершённость |
| Ход выполнения, у которого этапы названы | **Доля выполненного числом** — [мера](../charts/meter.md): полоса без имён этапов |
| — | **История событий с временем** — [лента](../display/timeline.md) |

## Устройство

### Подпись состояния

`inst-stepper-note` — не украшение и не подзаголовок. Это второй носитель
состояния: «готово», «сейчас», «дальше» читаются и на распечатке, и тем, кто не
различает оттенок полосы. Строка с датой или счётчиком в этом слоте состояние
не сообщает и его не заменяет.

### Раскладка

Полоса переносится сама: у элемента `flex: 1 1 10rem`, и на узком экране шаги
встают в несколько рядов без единого медиазапроса. Полоса состояния — это
`::before` шириной во весь элемент, поэтому она не расходится с подписью при
любом числе шагов.

## Состояния

`data-state` — фаза жизни, словарь у компонента **свой и закрытый**:

```html preview
<ol class="inst-stepper">
  <li class="inst-stepper-item" data-state="done">
    <span class="inst-stepper-name">Пройден</span><span class="inst-stepper-note">готово</span></li>
  <li class="inst-stepper-item" data-state="current">
    <span class="inst-stepper-name">Текущий</span><span class="inst-stepper-note">сейчас</span></li>
  <li class="inst-stepper-item" data-state="todo">
    <span class="inst-stepper-name">Предстоящий</span><span class="inst-stepper-note">дальше</span></li>
</ol>
```

| Значение | Значит | Что рисует |
|---|---|---|
| `todo` | Ещё не начат. **База** | Полоса `--track`, подпись `--text-secondary` |
| `current` | Идёт сейчас | Полоса `--accent-mark`, подпись `--text-primary` и `--weight-medium` |
| `done` | Завершён | Полоса `--accent-mark`, подпись `--text-secondary` |

Базовое значение `todo` **пишется в разметке**, хотя правил под него нет:
`data-state="todo"` читается, а его отсутствие — нет, и опечатка выглядела бы
ровно как база.

Пройденный и текущий шаг делят один цвет полосы намеренно: различает их
`inst-stepper-note` — «готово» против «сейчас» — и начертание подписи. Второго
акцентного тона в ките нет.

## Правила

:::do Слово состояния в каждом шаге
«готово», «сейчас», «дальше». Оно читается на распечатке и тем, кто не различает
оттенок полосы.
:::

:::dont Дата вместо слова состояния
`inst-stepper-note` — носитель состояния, а не слот для произвольной подписи.
Дата состояния не сообщает.
:::

:::do data-state="todo" явно
Отсутствие атрибута неотличимо от опечатки в нём. База пишется в разметке.
:::

:::dont Шаги вместо вкладок
У шагов есть порядок и незавершённость. Свободный переход между видами — это
[вкладки](./tabs.md).
:::

## Доступность

| | |
|---|---|
| Порядок | `<ol>` сообщает номер шага и их общее число без подписей вида «3 из 5» |
| Цвет не единственный носитель | Состояние ходит словом (`inst-stepper-note`) **и** полосой |
| Контраст | Название — `--text-secondary` (4.5:1), слово состояния — `--text-muted`; это данные, а не декорация, и `--text-faint` им не положен |
| Полоса | `--accent-mark` держит 3:1 против дорожки `--track` — это метка, а не текст |
| Кегль | Слово состояния — `--text-xs`, ниже 11px кит не опускается |
| Печать | Шаги печатаются: слово состояния читается на бумаге, где полоса могла бы слиться |

## API

```api
```

## Связанное

[Шаг прогона](../../agent/step.md) [Мера](../charts/meter.md)
[Вкладки](./tabs.md) [Лента событий](../display/timeline.md)
[Форма](../inputs/form.md)
