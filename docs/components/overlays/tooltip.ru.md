---
title: Тултип
group: Оверлеи
layout: component
source: src/overlay.css
api:
  - { name: "inst-tooltip", kind: "класс", doc: "Обёртка триггера. `position: relative`, отсюда считается положение" }
  - { name: "inst-tooltip-text", kind: "класс", doc: "Сама подпись. Над триггером, по центру" }
  - { name: "inst-tooltip-text--wide", kind: "модификатор", doc: "Перенос строк и потолок ширины `22rem`" }
  - { name: "--z-popover", kind: "токен" }
  - { name: "--surface-overlay", kind: "токен" }
  - { name: "--shadow-popover", kind: "токен" }
  - { name: "--border", kind: "токен" }
  - { name: "--hairline", kind: "токен" }
  - { name: "--radius-sm", kind: "токен" }
  - { name: "--space-2", kind: "токен" }
  - { name: "--space-3", kind: "токен" }
  - { name: "--text-xs", kind: "токен" }
  - { name: "--leading-ui", kind: "токен" }
  - { name: "--dur-1", kind: "токен" }
  - { name: "--ease-out", kind: "токен" }
---

Короткая подпись, появляющаяся при наведении и при фокусе. Единственный оверлей
библиотеки, сделанный **без** Popover API: обычным позиционированием, чтобы работать
везде одинаково.

```html preview
<span class="inst-tooltip">
  <button class="inst-btn inst-btn--icon inst-btn--ghost" type="button" aria-label="Обновить" aria-describedby="tt-refresh">
    <svg class="inst-icon" aria-hidden="true"><use href="#i-refresh"/></svg>
  </button>
  <span class="inst-tooltip-text" role="tooltip" id="tt-refresh">Обновить</span>
</span>
<span class="inst-tooltip">
  <button class="inst-btn inst-btn--ghost" type="button" aria-describedby="tt-pass">Что такое проход?</button>
  <span class="inst-tooltip-text" role="tooltip" id="tt-pass">Один полный обход карты агентом</span>
</span>
```

Подпись центрируется по триггеру и **не переворачивается**: у обычного
позиционирования нет запасных мест. У края узкой области её обрежет ближайший
`overflow: hidden` — там нужен [поповер](./popover.md).

## Контракт

| Что | Обязательно | Почему |
|---|---|---|
| `role="tooltip"` на подписи | да | Иначе это просто спрятанный текст |
| `id` на подписи | да | На него ссылается триггер |
| `aria-describedby` на триггере | да | Указывает на `id` подписи и есть единственное, что связывает её с контролом. Без него тултип для скринридера не существует |
| Триггер **фокусируем** | да | Показ идёт через `:focus-within`. Нефокусируемому элементу нужен `tabindex="0"` |

### Доступность

| | |
|---|---|
| Клавиатура | Тултип появляется по `:focus-within`, то есть при табе на триггер. Без этого он был бы доступен только мыши |
| Сенсорный экран | Наведения там нет, а фокус приходит вместе с нажатием. **Не кладите в тултип то, без чего нельзя обойтись** |
| Не заменяет имя | `aria-describedby` — это описание. Доступное имя иконочной кнопке даёт `aria-label`, а не тултип |
| Перекрытие | Тултип стоит над триггером на `--space-3`: он не закрывает то, что поясняет |
| Уменьшенное движение | Переход прозрачности схлопывается вместе со всеми остальными |
| Печать | На листе не печатается |

## Варианты

```html preview
<span class="inst-tooltip">
  <button class="inst-btn inst-btn--sm" type="button" aria-describedby="tt-short">Коротко</button>
  <span class="inst-tooltip-text" role="tooltip" id="tt-short">19:38:04</span>
</span>
<span class="inst-tooltip" style="margin-inline-start:9rem">
  <button class="inst-btn inst-btn--sm" type="button" aria-describedby="tt-wide">Подробно</button>
  <span class="inst-tooltip-text inst-tooltip-text--wide" role="tooltip" id="tt-wide">Открытие модалки — единственная строка скрипта во всём разделе оверлеев; закрытие, подложка и блокировка прокрутки обходятся без него.</span>
</span>
```

Базовый тултип — в одну строку (`white-space: nowrap`). Это верно для подписи
из двух слов и ломается на предложении: строка растягивается за край экрана.
Поэтому длинному тексту нужен явный `inst-tooltip-text--wide` — перенос и
потолок ширины в `22rem`.

## Поведение

| Событие | Селектор |
|---|---|
| наведение мышью | `.inst-tooltip:hover` |
| фокус внутри обёртки | `.inst-tooltip:focus-within` |

Подпись не убирается из потока, а гасится прозрачностью, и на ней снята мышь
(`pointer-events: none`): курсор, идущий к тултипу, не перехватывается им и не
роняет наведение.

:::warn
**Тултип обрежется внутри контейнера с `overflow: hidden`.** Он позиционируется
`position: absolute` относительно своей обёртки и не живёт в верхнем слое. Это
известная плата за то, что тултип работает одинаково всюду и не зависит от
поддержки якорного позиционирования.

Там, где обрезка возможна — внутри [панели](../display/panel.md), ячейки
[таблицы](../display/table.md), [аккордеона](../feedback/accordion.md), —
берите [поповер](./popover.md): он в верхнем слое, и его не срежет ничто.
:::

## API

```api
```

## Связанное

```related
```
