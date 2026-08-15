---
title: Таблица
group: Отображение данных
layout: component
source: src/components.css
js: Сортировка и действия строки. «Выбрать всё» делает `instrument.js`
api:
  - { name: "inst-table", kind: "класс", doc: "Базовый, на `<table>`" }
  - { name: "inst-num", kind: "класс", doc: "Числовая ячейка: выравнивание по концу" }
  - { name: "inst-sort", kind: "класс", doc: "Кнопка-заголовок сортировки" }
  - { name: "inst-col-select", kind: "класс", doc: "Колонка чекбоксов, ширина по контролу" }
  - { name: "inst-col-actions", kind: "класс", doc: "Колонка действий, ширина по содержимому" }
  - { name: "inst-row-actions", kind: "класс", doc: "Обёртка действий: проявляется на наведении и фокусе" }
  - { name: "inst-table--cells", kind: "модификатор", doc: "Ячейка-плитка: каждое значение в своей коробке, между коробками просвет" }
  - { name: "inst-table--zebra", kind: "модификатор", doc: "Полосы вместо линеек. Для широкой таблицы: на пятнадцати колонках линейка не показывает, ГДЕ ты внутри строки" }
  - { name: "inst-table--sticky", kind: "модификатор", doc: "Липкая шапка" }
  - { name: "inst-table--sticky-col", kind: "модификатор", doc: "Липкая первая колонка" }
  - { name: "aria-sort", kind: "атрибут", value: "ascending · descending", doc: "на `<th>`" }
  - { name: "aria-selected", kind: "атрибут", value: "true · false", doc: "на `<tr>`" }
  - { name: "--table-head-bg", kind: "переменная", value: "--surface-raised" }
  - { name: "--text-sm", kind: "токен" }
  - { name: "--pad-cell-x", kind: "токен" }
  - { name: "--row-pad-y", kind: "токен" }
  - { name: "--space-3", kind: "токен" }
  - { name: "--border", kind: "токен" }
  - { name: "--border-subtle", kind: "токен" }
  - { name: "--hairline", kind: "токен" }
  - { name: "--text-muted", kind: "токен" }
  - { name: "--surface-selected", kind: "токен" }
  - { name: "--surface-hover", kind: "токен" }
  - { name: "--control-h-md", kind: "токен" }
  - { name: "--z-sticky", kind: "токен" }
title-en: "Table"
group-en: "Data display"
---

Сравнимые записи в колонках. Настоящая `<table>` — с семантикой заголовков,
навигацией скринридера по ячейкам и печатью.

```html preview
<table class="inst-table">
  <thead>
    <tr><th>Прогон</th><th>Модель</th><th class="inst-num">Время</th></tr>
  </thead>
  <tbody>
    <tr><td>worldgen-01</td><td>opus</td><td class="inst-num">2,1 с</td></tr>
    <tr><td>biomes-04</td><td>sonnet</td><td class="inst-num">8,4 с</td></tr>
  </tbody>
</table>
```

## Использование

```html
<table class="inst-table">
  <caption>Модули прогона</caption>
  <thead>
    <tr><th>Прогон</th><th class="inst-num">Время</th></tr>
  </thead>
  <tbody>
    <tr><td>worldgen-01</td><td class="inst-num">2,1 с</td></tr>
  </tbody>
</table>
```

| Что | Обязательно | Почему |
|---|---|---|
| Настоящая `<table>` с `<th>` в `<thead>` | да | Скринридер озвучивает заголовок при переходе по ячейкам — без него таблица превращается в набор чисел |
| `<caption>` или `aria-label` | да | Две таблицы на экране без имён неразличимы на слух |
| `inst-num` на числовых ячейках | да, у чисел | Выравнивание по концу строки: без него разряды не сравнить глазом |
| Достижимая с клавиатуры область прокрутки | да, при горизонтальной прокрутке | Иначе часть колонок существует только для мыши |

:::note
**Разделение труда.** «Выбрать всё» делает [instrument.js](#js), включая частичный
выбор. Сортировка и действия строки — слой приложения: порядок строк это
данные, и библиотеке они неизвестны.
:::

## Когда использовать

| Используйте | Возьмите другое |
|---|---|
| Записи, которые **сравнивают** по одним и тем же полям | **Разнородные объекты** — [карточки](./card.md): в таблице половина ячеек будет пустой |
| Больше трёх полей на запись | **Пара «имя — значение»** — [список свойств](./kv.md) |
| Нужны сортировка, выбор строк, итоги | **Строки с состоянием и прогрессом** — [очередь](../../agent/task.md): она плоская и не имеет колонок |
| Данные уйдут в печать | **Десятки тысяч строк** — нужна виртуализация, её в библиотеке нет |

## Устройство

### Числовые колонки

```html
<th class="inst-num">Время</th>
<td class="inst-num">2,1 с</td>
```

`inst-num` выравнивает по концу строки. Табличные цифры включены по умолчанию
для всей библиотеки, поэтому разряды встают друг под друга без дополнительных
усилий.

### Итоги

```html
<tfoot>
  <tr><td>Всего</td><td></td><td class="inst-num">10,5 с</td></tr>
</tfoot>
```

`<tfoot>` идёт в разметке **после** `<tbody>` и отбивается сверху. Строка
итогов — не последняя строка данных, и путать их нельзя.

## Состояния

### Выбор строк

```html preview
<table class="inst-table">
  <thead><tr><th class="inst-col-select"><label class="inst-checkbox"><input type="checkbox" aria-label="Выбрать все"></label></th><th>Прогон</th></tr></thead>
  <tbody>
    <tr aria-selected="true">
      <td class="inst-col-select"><label class="inst-checkbox"><input type="checkbox" checked aria-label="Выбрать worldgen-01"></label></td>
      <td>worldgen-01</td>
    </tr>
    <tr>
      <td class="inst-col-select"><label class="inst-checkbox"><input type="checkbox" aria-label="Выбрать biomes-04"></label></td>
      <td>biomes-04</td>
    </tr>
  </tbody>
</table>
```

Состояние живёт на `aria-selected` у `<tr>`. Подсветка полупрозрачна, поэтому
состояния складываются: наведение на выбранную строку остаётся видно.

## Поведение

### Сортировка

```html
<th aria-sort="ascending">
  <button class="inst-sort" type="button">Время</button>
</th>
```

| Значение `aria-sort` | Что показывает |
|---|---|
| нет атрибута | Колонка не сортирована |
| `ascending` | По возрастанию, значок вверх |
| `descending` | По убыванию, значок вниз |

Направление несёт **значок**, и то же самое озвучивает `aria-sort`. Одной
позиции строк мало. Заголовок-сортировка — настоящая `<button>` внутри `<th>`,
иначе до неё не добраться с клавиатуры.

У числовой колонки значок идёт перед подписью, иначе он отрывается от неё на
всю ширину ячейки.

### Действия строки

```html
<td class="inst-col-actions">
  <div class="inst-row-actions">
    <button class="inst-btn inst-btn--sm inst-btn--icon" type="button" aria-label="Обновить">
      <svg class="inst-icon" aria-hidden="true"><use href="#i-refresh"/></svg>
    </button>
  </div>
</td>
```

Действия появляются на наведении — **и всегда доступны с клавиатуры**: они
проявляются также по `:focus-within` и на выбранной строке. Действие, видимое
только под курсором, для клавиатуры не существует.

### Липкие шапка и колонка

```html
<table class="inst-table inst-table--sticky inst-table--sticky-col">
```

| Класс | Что липнет |
|---|---|
| `inst-table--sticky` | Шапка при вертикальной прокрутке |
| `inst-table--sticky-col` | Первая колонка при горизонтальной |

Липкость работает относительно **прокручиваемого предка** — обычно
`inst-panel-body`. Таблица в неприкреплённом потоке ничего не приклеит.

## JS

Подключите модуль один раз на страницу — инициализировать компоненты по
отдельности не нужно, `instrument.js` работает делегированием и видит узлы, пришедшие
позже.

```html
<script type="module" src="instrument.js"></script>
```

### Что делает `instrument.js`

Пример в разделе «Состояния» живой: отметьте чекбокс в шапке — выберутся все строки, снимите одну — он станет промежуточным.

Только чекбокс «выбрать всё» в шапке колонки выбора:

- отмечает и снимает все строки, кроме `disabled`;
- держит `aria-selected` на `<tr>`;
- показывает частичный выбор через `indeterminate` — «часть строк выбрана»
  невыразимо ни через `checked`, ни через его отсутствие;
- пересчитывает себя, когда строку отметили поодиночке.

### События

`inst:selectall` всплывает с таблицы, `detail` — `{ checked }`. Отмена
оставляет строки как есть.

```js
table.addEventListener('inst:selectall', (e) => {
  bulkBar.hidden = !e.detail.checked;
});
```

### Что остаётся приложению

| Что | Почему не `instrument.js` |
|---|---|
| Сортировка | Порядок строк — это данные. Библиотека рисует `aria-sort` и стрелку, переставляет строки приложение |
| Действия строки | Что делает «Остановить» — знает только приложение |
| Пагинация и виртуализация | Зависят от источника данных |

## Сценарии

### Таблица в панели

Таблица показывает себя только целиком: сортировка, выбор строк, действия и
итоги по отдельности выглядят как четыре независимых класса, а вместе — как
один инструмент.

Что проверить руками:

- **наведите на строку** — появятся действия. Теперь пройдите по строкам
  `Tab`: действия появятся снова. Это разные события (`:hover` и
  `:focus-within`), и второе существует затем, чтобы действие,
  видимое только под курсором, не пропадало для клавиатуры;
- **выделите строку** и наведите на неё — подсветки складываются, потому
  что выделение полупрозрачно;
- **прокрутите тело панели** — шапка останется на месте.

```html preview context
<div class="inst-panel">
  <div class="inst-panel-header">
    <span class="inst-panel-title">Модули прогона</span>
    <span class="inst-panel-actions"><span class="inst-badge">выбрано 1</span></span>
  </div>
  <div class="inst-panel-body inst-panel-body--flush">
    <table class="inst-table inst-table--sticky">
      <thead>
        <tr>
          <th class="inst-col-select"><label class="inst-checkbox"><input type="checkbox" aria-label="Выбрать все строки"></label></th>
          <th aria-sort="descending"><button class="inst-sort" type="button">Модуль</button></th>
          <th>Состояние</th>
          <th class="inst-num"><button class="inst-sort" type="button">Время, с</button></th>
          <th class="inst-num"><button class="inst-sort" type="button">Токенов</button></th>
          <th class="inst-col-actions"><span class="inst-u-visually-hidden">Действия</span></th>
        </tr>
      </thead>
      <tbody>
        <tr aria-selected="true">
          <td><label class="inst-checkbox"><input type="checkbox" checked aria-label="Выбрать «Рендер»"></label></td>
          <td>Рендер</td>
          <td><span class="inst-badge" data-tone="ok"><span class="inst-dot"></span>готово</span></td>
          <td class="inst-num">18,2</td><td class="inst-num">61 204</td>
          <td class="inst-col-actions"><span class="inst-row-actions">
            <button class="inst-btn inst-btn--sm inst-btn--ghost inst-btn--icon" type="button" aria-label="Перезапустить Рендер">
              <svg class="inst-icon" aria-hidden="true"><use href="#i-refresh"/></svg></button>
          </span></td>
        </tr>
        <tr>
          <td><label class="inst-checkbox"><input type="checkbox" aria-label="Выбрать «Физика»"></label></td>
          <td>Физика</td>
          <td><span class="inst-badge" data-tone="ok"><span class="inst-dot"></span>готово</span></td>
          <td class="inst-num">11,5</td><td class="inst-num">38 910</td>
          <td class="inst-col-actions"><span class="inst-row-actions">
            <button class="inst-btn inst-btn--sm inst-btn--ghost inst-btn--icon" type="button" aria-label="Перезапустить Физика">
              <svg class="inst-icon" aria-hidden="true"><use href="#i-refresh"/></svg></button>
          </span></td>
        </tr>
        <tr>
          <td><label class="inst-checkbox"><input type="checkbox" aria-label="Выбрать «Аудио»"></label></td>
          <td>Аудио</td>
          <td><span class="inst-badge" data-tone="error"><span class="inst-dot"></span>упало</span></td>
          <td class="inst-num">4,2</td><td class="inst-num">1 180</td>
          <td class="inst-col-actions"><span class="inst-row-actions">
            <button class="inst-btn inst-btn--sm inst-btn--ghost inst-btn--icon" type="button" aria-label="Перезапустить Аудио">
              <svg class="inst-icon" aria-hidden="true"><use href="#i-refresh"/></svg></button>
          </span></td>
        </tr>
      </tbody>
      <tfoot>
        <tr><td></td><td>Итого</td><td></td><td class="inst-num">33,9</td><td class="inst-num">101 294</td><td></td></tr>
      </tfoot>
    </table>
  </div>
</div>
```

### Ячейка-плитка

```html preview
<table class="inst-table inst-table--cells">
  <thead><tr><th>Прогон</th><th>Ветка</th><th>Модель</th><th class="inst-num">Токенов</th></tr></thead>
  <tbody>
    <tr><td>nightly-4127</td><td>main</td><td>opus</td><td class="inst-num">186 000</td></tr>
    <tr><td>audio-pass</td><td>fix/mixdown</td><td>sonnet</td><td class="inst-num">42 300</td></tr>
    <tr><td>biomes-04</td><td>main</td><td>opus</td><td class="inst-num">98 700</td></tr>
  </tbody>
</table>
```

Линейка отделяет строку от строки и ничего не говорит про **колонку**: на
пятнадцати колонках глаз, идущий вправо, теряет и строку, и границу значения.
Плитка держит и то и другое сразу — у каждого значения своя коробка.

Шапка плиткой не становится: она не значение, а имя колонки, и коробка вокруг
имени сделала бы её равной данным.

### Полосы вместо линеек

```html preview
<table class="inst-table inst-table--zebra">
  <thead><tr><th>Прогон</th><th>Ветка</th><th>Модель</th><th class="inst-num">Токенов</th><th class="inst-num">Время</th></tr></thead>
  <tbody>
    <tr><td>nightly-4127</td><td>main</td><td>opus</td><td class="inst-num">186 000</td><td class="inst-num">21 с</td></tr>
    <tr><td>audio-pass</td><td>fix/mixdown</td><td>sonnet</td><td class="inst-num">42 300</td><td class="inst-num">4,2 с</td></tr>
    <tr><td>biomes-04</td><td>main</td><td>opus</td><td class="inst-num">98 700</td><td class="inst-num">11,5 с</td></tr>
    <tr><td>terrain-88</td><td>main</td><td>sonnet</td><td class="inst-num">150 400</td><td class="inst-num">18,2 с</td></tr>
  </tbody>
</table>
```

Линейка под каждой строкой работает, пока колонок пять. На пятнадцати она
отделяет строку от соседа, но не показывает, **где вы внутри строки**, — глаз
уезжает по вертикали, пока идёт вправо. Заливка через одну показывает: она
тянется на всю ширину и не пропадает при горизонтальной прокрутке.

Вместе с линейками её ставить нельзя: два разделителя на одной границе дают
решётку. Поэтому вариант линейки **снимает**, а не дополняет.

| Берите полосы | Берите линейки |
|---|---|
| Колонок больше семи, таблицу прокручивают вбок | Колонок до пяти, строка видна целиком |
| Строки читают целиком, слева направо | Строки сравнивают между собой по одной колонке |
| Таблица длинная и стоит на своей поверхности | Таблица короткая и лежит внутри карточки |

## Правила

:::do Настоящая table с th в thead
Скринридер озвучивает заголовок при переходе по ячейкам, работает печать,
работает выделение колонки.
:::

:::dont Сетка из div с role="grid"
Она теряет навигацию по ячейкам, печать и половину поведения, которое таблица
даёт бесплатно.
:::

:::do Действия и по focus-within
Действие, видимое только под курсором, для клавиатуры не существует.
:::

:::dont Таблица для разнородных объектов
Половина ячеек окажется пустой — это [карточки](./card.md).
:::

## Доступность

| | |
|---|---|
| Заголовки | `<th>` в `<thead>`. Скринридер озвучивает заголовок при переходе по ячейкам — без него таблица превращается в набор чисел |
| Сортировка | `aria-sort` на `<th>` плюс `<button>` внутри. Атрибут без кнопки недостижим, кнопка без атрибута не озвучивает направление |
| Выбор | `aria-selected` на `<tr>`. Чекбокс несёт своё имя, иначе колонка озвучится как «отметка, отметка, отметка» |
| Действия строки | Проявляются по `:focus-within`, а не только по `:hover`. Иначе они существуют лишь для мыши |
| Подпись таблицы | `<caption>` или `aria-label`. Две таблицы на экране без имён неразличимы на слух |
| Не только цвет | Выбранная строка несёт `aria-selected`, а не одну лишь подсветку |
| Горизонтальная прокрутка | Область прокрутки должна быть достижима с клавиатуры |

## API

```api
```

## Связанное

[Карточка](./card.md) [Список свойств](./kv.md) [Панель](./panel.md)
[Очередь](../../agent/task.md) [Пагинация](../navigation/pagination.md)
[Пустое состояние](../feedback/empty.md)
