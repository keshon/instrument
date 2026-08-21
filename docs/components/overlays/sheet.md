---
title: Шторка
group: Оверлеи
layout: component
source: src/overlay.css
js: "Открытие — одна строка: dlg.showModal(). Закрытие, подложка, блокировка прокрутки фона и возврат фокуса скрипта не требуют"
api:
  - { name: "inst-sheet", kind: "класс", doc: "Раскладка во всю высоту у конечного края. Только вместе с `inst-dialog`" }
  - { name: "inst-sheet--start", kind: "модификатор", doc: "Тот же блок у начального края" }
  - { name: "--surface-overlay", kind: "токен" }
  - { name: "--shadow-modal", kind: "токен" }
  - { name: "--scrim", kind: "токен" }
  - { name: "--border", kind: "токен" }
  - { name: "--hairline", kind: "токен" }
  - { name: "--pad-card", kind: "токен" }
  - { name: "--space-5", kind: "токен" }
  - { name: "--gap-inline", kind: "токен" }
  - { name: "--text-md", kind: "токен" }
title-en: "Sheet"
group-en: "Overlays"
---

Панель во всю высоту у края экрана. Тот же нативный `<dialog>`, что и у
[модалки](./dialog.md), — другая раскладка, не другой компонент.

```html preview
<button class="inst-btn" type="button"
        onclick="document.getElementById('dlg-params').showModal()">От конечного края</button>
<button class="inst-btn" type="button"
        onclick="document.getElementById('dlg-nav').showModal()">От начального</button>

<dialog class="inst-dialog inst-sheet inst-sheet--start" id="dlg-nav">
  <form method="dialog">
    <div class="inst-dialog-head">
      <span class="inst-dialog-title">Разделы</span>
      <button class="inst-btn inst-btn--sm inst-btn--ghost inst-btn--icon inst-dialog-close" type="submit" aria-label="Закрыть">
        <svg class="inst-icon" aria-hidden="true"><use href="#i-close"/></svg>
      </button>
    </div>
    <div class="inst-dialog-body">
      <nav class="inst-nav" aria-label="Разделы">
        <a class="inst-nav-item" href="#" aria-current="page">Прогоны</a>
        <a class="inst-nav-item" href="#">Артефакты</a>
        <a class="inst-nav-item" href="#">Настройки</a>
      </nav>
    </div>
  </form>
</dialog>

<dialog class="inst-dialog inst-sheet" id="dlg-params">
  <form method="dialog">
    <div class="inst-dialog-head">
      <span class="inst-dialog-title">Параметры прогона</span>
      <button class="inst-btn inst-btn--sm inst-btn--ghost inst-btn--icon inst-dialog-close" type="submit" aria-label="Закрыть">
        <svg class="inst-icon" aria-hidden="true"><use href="#i-close"/></svg>
      </button>
    </div>
    <div class="inst-dialog-body">
      <dl class="inst-kv">
        <dt>Модель</dt><dd>opus</dd>
        <dt>Агентов</dt><dd>7 из 12</dd>
        <dt>Токенов</dt><dd>184 320</dd>
        <dt>Начат</dt><dd>19:38:04</dd>
      </dl>
    </div>
    <div class="inst-dialog-foot">
      <button class="inst-btn" type="submit">Закрыть</button>
    </div>
  </form>
</dialog>
```

## Контракт

`inst-sheet` — модификатор поверх `inst-dialog`, а не замена: класс базы
остаётся, части (`-head`, `-body`, `-foot`) те же.

| Что | Обязательно | Почему |
|---|---|---|
| `inst-dialog` вместе с `inst-sheet` | да | Поверхность, тень, подложка и раскладка частей приходят из базы |
| Нативный `<dialog>` | да | Модальность и `Escape` — от платформы |
| `showModal()` | да | `show()` открывает немодально: ни подложки, ни инертного фона |
| `<form method="dialog">` внутри | да | Закрытие кнопкой без скрипта |
| `aria-label` у крестика | да | Кнопка без текста иначе не имеет доступного имени |

### Доступность

| | |
|---|---|
| Фокус | Ловушка фокуса и возврат на вызвавшую кнопку — от `showModal()` |
| Клавиатура | `Escape` закрывает. Кнопки внутри `<form method="dialog">` закрывают нажатием |
| Фон инертен | Подложка выключает фон для мыши и для скринридера. Если нужен доступ к фону, это не шторка, а колонка [оболочки](../../layout/shell.md) |
| Прокрутка фона | Останавливается тем же правилом `html:has(dialog:modal)`, что и у модалки |
| Высота | `100dvh`, а не `100vh`: на мобильных панель браузера не срезает подвал с кнопками |
| Печать | Шторка на листе не печатается |

## Устройство

Все части — от [модалки](./dialog.md): `inst-dialog-head`,
`inst-dialog-title`, `inst-dialog-close`, `inst-dialog-body`,
`inst-dialog-foot`, `inst-dialog-foot--start`, `inst-dialog-foot-note`.

Модификатор меняет раскладку:

| Свойство | Модалка | Шторка |
|---|---|---|
| ширина | `min(34rem, 100vw - var(--space-8))` | `min(28rem, 100vw)` |
| высота | `min(80dvh, 100dvh - var(--space-8))` | `100dvh` |
| поля | `margin: auto` — по центру | Прижата к краю |
| радиус | `--radius-lg` | `0` — блок примыкает к краю экрана |
| рамка | По периметру | Только со стороны содержимого |

На узком экране `min(28rem, 100vw)` даёт шторку во всю ширину — отдельного
мобильного варианта заводить не нужно.

## Варианты

```html
<!-- У конечного края (по умолчанию) -->
<dialog class="inst-dialog inst-sheet">…</dialog>

<!-- У начального края -->
<dialog class="inst-dialog inst-sheet inst-sheet--start">…</dialog>
```

| Класс | Край | Граница |
|---|---|---|
| `inst-sheet` | Конечный (справа в LTR) | Волосок по начальному краю |
| `inst-sheet--start` | Начальный (слева в LTR) | Волосок по конечному краю |

Стороны логические, а не физические: в RTL шторка зеркалится сама, править
ничего не нужно.

## JS

Открытие — одна строка. Всё остальное берёт на себя платформа.

```js
document.getElementById('dlg').showModal();
```

| Что | Кто делает |
|---|---|
| Открытие | Приложение: `showModal()` |
| Закрытие | `<form method="dialog">` внутри — без скрипта |
| Подложка, верхний слой, `Escape` | Платформа |
| Возврат фокуса на открывшую кнопку | Платформа |
| Блокировка прокрутки фона | Библиотека, правилом `html:has(dialog:modal)` |

:::note
Декларативные `command` и `commandfor` избавили бы и от этой строки, но пока
поддержаны не везде и в контракт библиотеки не берутся.
:::

## API

```api
```

## Связанное

```related
```
