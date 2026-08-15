---
title: Меню
group: Оверлеи
layout: component
source: src/overlay.css
js: Стрелки и бегущий tabindex делает `instrument.js`. Открытие и закрытие берёт на себя Popover API
api:
  - { name: "inst-menu", kind: "класс", doc: "Контейнер. Колонка пунктов с зазором `--space-1`" }
  - { name: "inst-menu-item", kind: "класс", doc: "Пункт: `<button>` или `<a>`" }
  - { name: "inst-menu-item-main", kind: "класс", doc: "Обёртка названия и пояснения. Её наличие делает пункт двухстрочным — модификатора для этого нет" }
  - { name: "inst-menu-item-sub", kind: "класс", doc: "Пояснение под названием. Для меню выбора, где одно название ни о чём не говорит" }
  - { name: "inst-menu-label", kind: "класс", doc: "Подпись группы. Не интерактивна" }
  - { name: "inst-menu-sep", kind: "класс", doc: "Разделитель толщиной в волосок" }
  - { name: "inst-menu-shortcut", kind: "класс", doc: "Горячая клавиша у дальнего края пункта" }
  - { name: "data-tone", kind: "атрибут", doc: "`neutral` `running` `ok` `warn` `error`. Закрыт. На пункте осмыслен только `error`" }
  - { name: "aria-checked", kind: "атрибут", doc: "`true` — пункт отмечен" }
  - { name: "aria-current", kind: "атрибут", doc: "Любое значение, кроме `false`" }
  - { name: "aria-disabled", kind: "атрибут", doc: "`true` — пункт недоступен" }
  - { name: "--space-1", kind: "токен" }
  - { name: "--space-2", kind: "токен" }
  - { name: "--space-3", kind: "токен" }
  - { name: "--space-6", kind: "токен" }
  - { name: "--gap-inline", kind: "токен" }
  - { name: "--control-h-sm", kind: "токен" }
  - { name: "--radius-sm", kind: "токен" }
  - { name: "--text-2xs", kind: "токен" }
  - { name: "--text-sm", kind: "токен" }
  - { name: "--weight-medium", kind: "токен" }
  - { name: "--hairline", kind: "токен" }
  - { name: "--border-subtle", kind: "токен" }
  - { name: "--surface-hover", kind: "токен" }
  - { name: "--accent-text", kind: "токен" }
  - { name: "--err-text", kind: "токен" }
  - { name: "--err-bg", kind: "токен" }
  - { name: "--text-muted", kind: "токен" }
  - { name: "--leading-ui", kind: "токен" }
title-en: "Menu"
group-en: "Overlays"
---

Список действий над объектом. Меню — это **содержимое**
[поповера](./popover.md), а не самостоятельный оверлей: верхний слой,
закрытие и фокус приходят оттуда.

```html preview
<button class="inst-btn" type="button" popovertarget="menu-run">Действия</button>
<div class="inst-popover inst-popover--anchored" id="menu-run" popover>
  <div class="inst-menu" role="menu">
    <span class="inst-menu-label">Прогон #4127</span>
    <button class="inst-menu-item" type="button" role="menuitem">
      <svg class="inst-icon" aria-hidden="true"><use href="#i-refresh"/></svg>Перезапустить
      <span class="inst-menu-shortcut"><kbd>R</kbd></span></button>
    <button class="inst-menu-item" type="button" role="menuitem">
      <svg class="inst-icon" aria-hidden="true"><use href="#i-copy"/></svg>Скопировать id</button>
    <button class="inst-menu-item" type="button" role="menuitem" aria-checked="true">
      <svg class="inst-icon" aria-hidden="true"><use href="#i-list"/></svg>Показывать лог</button>
    <span class="inst-menu-sep"></span>
    <button class="inst-menu-item" type="button" role="menuitem" data-tone="error">Удалить прогон</button>
    <button class="inst-menu-item" type="button" role="menuitem" aria-disabled="true">Архивировать</button>
  </div>
</div>
```

## Использование

```html
<div class="inst-popover inst-popover--anchored" id="menu-1" popover>
  <div class="inst-menu" role="menu">
    <button class="inst-menu-item" type="button" role="menuitem">Перезапустить</button>
  </div>
</div>
```

| Что | Обязательно | Почему |
|---|---|---|
| `role="menu"` на контейнере | да | Без него `role="menuitem"` невалиден |
| `role="menuitem"` на пунктах | да | Библиотека рисует состояние, но не выдумывает роль |
| `type="button"` у пунктов-кнопок | да | Иначе внутри формы меню её отправит |
| `aria-hidden="true"` у иконки | да | Имя пункту даёт подпись, а не иконка |
| Обёртка [поповера](./popover.md) | да | Меню само по себе не всплывает: верхний слой и закрытие приходят от `popover` |

:::warn
Перемещение стрелками и бегущий `tabindex` выполняет
[`instrument.js`](../../foundations/behavior.md). Открытие, закрытие по `Escape` и по
клику мимо, возврат фокуса на кнопку скрипта не требуют вовсе — их берёт на
себя Popover API.
:::

## Когда использовать

| Используйте | Возьмите другое |
|---|---|
| Несколько действий над объектом, которые не помещаются в строку | **Два-три частых действия** — [группа кнопок](../actions/button-group.md) на виду: меню прячет то, что нажимают каждый раз |
| Редкие и разрушительные действия, убранные из основного вида | **Выбор значения формы** — [селект](../inputs/select.md): у меню нет значения и оно не отправляется |
| Переключаемая настройка вида — пункт с `aria-checked` | **Выбор одного из равных режимов на виду** — [сегментированный контрол](../actions/segmented.md) |
| Навигация по разделам во всплывающем списке | **Постоянная навигация приложения** — [боковая навигация](../navigation/nav.md) |

## Устройство

Подпись группы, разделитель и горячая клавиша — три разные работы, и путать их
дорого.

```html
<div class="inst-menu" role="menu">
  <span class="inst-menu-label">Прогон #4127</span>
  <button class="inst-menu-item" type="button" role="menuitem">Перезапустить
    <span class="inst-menu-shortcut"><kbd>R</kbd></span></button>
  <span class="inst-menu-sep"></span>
  <button class="inst-menu-item" type="button" role="menuitem">Экспорт</button>
</div>
```

| Класс | Работа |
|---|---|
| `inst-menu-label` | Подпись группы. Не интерактивна и не попадает в обход |
| `inst-menu-item` | Пункт. `<button>` для действия, `<a>` для перехода |
| `inst-menu-sep` | Разделитель групп. `<span>`, а не `<hr>`: он оформление, и озвучивать его нечем |
| `inst-menu-shortcut` | Горячая клавиша у дальнего края |

Горячая клавиша прижимается к дальнему краю и тише подписи: она **подсказка, а
не второе название пункта**.

Пункт — это `<button>` или `<a>`, смотря по работе. Действие — кнопка, переход
по адресу — ссылка: кнопка, ведущая на адрес, ломает средний клик и «открыть в
новой вкладке».

## Состояния

```html preview
<div class="inst-popover" popover id="menu-kinds">
  <div class="inst-menu" role="menu">
    <span class="inst-menu-label">Вид</span>
    <button class="inst-menu-item" type="button" role="menuitem">Обычный</button>
    <button class="inst-menu-item" type="button" role="menuitem" aria-checked="true">Отмеченный</button>
    <a class="inst-menu-item" href="#menu" role="menuitem" aria-current="page">Текущий адрес</a>
    <span class="inst-menu-sep"></span>
    <button class="inst-menu-item" type="button" role="menuitem" data-tone="error">Разрушительный</button>
    <button class="inst-menu-item" type="button" role="menuitem" aria-disabled="true">Недоступный</button>
  </div>
</div>
<button class="inst-btn" type="button" popovertarget="menu-kinds">Виды пунктов</button>
```

Пункт с пояснением. Обёртка `.inst-menu-item-main` делает пункт двухстрочным
сама — модификатора для этого нет, потому что это не выбор автора разметки, а
следствие того, что внутри. Меню выбора без пояснения превращается в список
кодовых слов.

```html preview
<div class="inst-popover" popover id="menu-model">
  <div class="inst-menu" role="menu">
    <button class="inst-menu-item" type="button" role="menuitem">
      <span class="inst-menu-item-main">Быстрый
        <span class="inst-menu-item-sub">Для коротких прогонов</span></span></button>
    <button class="inst-menu-item" type="button" role="menuitem" aria-checked="true">
      <span class="inst-menu-item-main">Полный
        <span class="inst-menu-item-sub">Для сложных задач</span></span></button>
    <button class="inst-menu-item" type="button" role="menuitem">
      <span class="inst-menu-item-main">Экономный
        <span class="inst-menu-item-sub">Дешевле всего на повседневном</span></span></button>
  </div>
</div>
<button class="inst-btn" type="button" popovertarget="menu-model">Режим прогона</button>
```

| Состояние пункта | Как ставится | Что происходит |
|---|---|---|
| обычный | — | `--text-primary` |
| наведение | `:hover` | `--surface-hover` |
| отмеченный | `aria-checked="true"` | `--accent-text` и средняя насыщенность |
| текущий | `aria-current` (любое значение, кроме `false`) | То же оформление, что у отмеченного |
| разрушительный | `data-tone="error"` | `--err-text`, а на наведении `--err-bg` |
| недоступный | `aria-disabled="true"` или `disabled` | Прозрачность `0.5`, мышь снята |

## JS

Подключите модуль один раз на страницу — инициализировать компоненты по
отдельности не нужно, `instrument.js` работает делегированием и видит узлы, пришедшие
позже.

```html
<script type="module" src="instrument.js"></script>
```

### Что делает `instrument.js`

Пример в шапке живой: откройте меню и пройдите по пунктам стрелками.

Стрелки между пунктами, `Home`, `End`, перебор по кругу и бегущий `tabindex`.
Выделения у пункта нет и не будет: пункт меню — действие, а не выбор, и
`aria-selected` на нём соврал бы вспомогательной технологии.

Открытие, закрытие, возврат фокуса и закрытие по `Escape` берёт на себя
Popover API — скрипта они не требуют.

### События

Своих нет. Пункт — это `<button>`, и его `click` работает так же, как везде.

```js
menu.addEventListener('click', (e) => {
  const item = e.target.closest('[role="menuitem"]');
  if (!item) return;
  menu.hidePopover();
  run(item.dataset.action);
});
```

Меню закрывается **до** действия: действие может открыть модалку, и меню
осталось бы висеть поверх неё.

## Правила

:::do Действие — кнопка, переход — ссылка
Пункт, ведущий на адрес, обязан быть `<a>`: иначе ломаются средний клик и
«открыть в новой вкладке».
:::

:::dont Частое действие в меню
Меню прячет то, что нажимают каждый раз. Два-три частых действия стоят на виду
— [группа кнопок](../actions/button-group.md).
:::

:::do Разрушительный пункт назван словом
«Удалить прогон» красным. Тон — второй признак, а не единственный.
:::

:::dont role="menu" без стрелок
Роль обещает перемещение стрелками. Без него длинное меню проходится табом, что
мучительно.
:::

## Доступность

| | |
|---|---|
| Клавиатура | `Tab` до кнопки, `Enter` — открыть, `Escape` — закрыть. **Стрелки внутри меню ставит приложение**: без них меню проходится табом, что для длинного списка мучительно |
| Роли | `role="menu"` + `role="menuitem"` обязательны. `aria-checked` без `role="menuitem"` не озвучивается |
| Недоступный пункт | `aria-disabled="true"`, а не `disabled`, если пункт должен оставаться в порядке обхода и объяснять, почему он недоступен |
| Цвет не единственный носитель | Разрушительный пункт красный **и** назван словом «Удалить». Отмеченный несёт `aria-checked`, а не только цвет |
| Цель нажатия | Пункт высотой `--control-h-sm` во всю ширину поповера: попасть мышью проще, чем в текст |
| Перенос | `white-space: nowrap` — пункт не переносится. Длинную подпись сокращайте, а не надейтесь на перенос |

## API

```api
```

## Связанное

[Поповер](./popover.md) [Кнопка](../actions/button.md)
[Селект](../inputs/select.md) [Тулбар](../navigation/toolbar.md)
[Иконка](../../foundations/icons.md)
