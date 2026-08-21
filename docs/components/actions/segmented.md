---
title: Сегментированный контрол
group: Действия
layout: component
source: src/actions.css
js: Выбор, стрелки и бегущий tabindex — делает `instrument.js`
api:
  - { name: "inst-segmented", kind: "класс", doc: "Контейнер. Внутри — голые `button` или `a`, свой класс им не нужен" }
  - { name: "aria-checked", kind: "атрибут", doc: "`true` · `false` на кнопках. Ровно один `true` в группе" }
  - { name: "aria-current", kind: "атрибут", value: "page", doc: "Носитель выбора, если варианты — ссылки" }
  - { name: "--control-h-md", kind: "токен" }
  - { name: "--control-pad-md", kind: "токен" }
  - { name: "--radius-md", kind: "токен" }
  - { name: "--radius-sm", kind: "токен" }
  - { name: "--space-1", kind: "токен" }
  - { name: "--text-sm", kind: "токен" }
  - { name: "--surface-sunken", kind: "токен" }
  - { name: "--surface-raised", kind: "токен" }
  - { name: "--border", kind: "токен" }
  - { name: "--hairline", kind: "токен" }
title-en: "Segmented control"
group-en: "Actions"
---

Один выбор среди равных вариантов. У него всегда ровно одно выбранное значение
— и в этом всё отличие от [группы кнопок](./button-group.md), где выбранного
нет и быть не может.

```html preview
<div class="inst-segmented" role="radiogroup" aria-label="Режим отображения">
  <button type="button" role="radio" aria-checked="true"  tabindex="0">Список</button>
  <button type="button" role="radio" aria-checked="false" tabindex="-1">Сетка</button>
  <button type="button" role="radio" aria-checked="false" tabindex="-1">Таблица</button>
</div>
```

Переключатель [плотности](../../foundations/density.md) в шапке справочника —
это он.

## Контракт

Состояние приходит из разметки: библиотека его рисует, но не вычисляет. Без
ролей `aria-checked` **невалиден**, а с клавиатуры до контрола не добраться.

| Что | Обязательно | Почему |
|---|---|---|
| `role="radiogroup"` на контейнере | да | Группа — один контрол с одним значением |
| `aria-label` на контейнере | да | Иначе у контрола нет имени: подписи есть только у вариантов |
| `role="radio"` на каждой кнопке | да | `role="tab"` без панели — невалидная роль |
| `aria-checked` на каждой | да | Носитель состояния. Стиль выбранного висит на нём, а не на классе |
| `type="button"` | да | Иначе внутри формы кнопка её отправит |
| Бегущий `tabindex` | да, начальный | `0` у выбранной, `-1` у остальных. Без него `Tab` пройдёт по всем вариантам и группа перестанет быть одним контролом |

Клавиатура: `Tab` — вход в группу и выход из неё, `←` `→` `↑` `↓` — перебор,
выбор следует за фокусом.

Выбранный несёт **два** признака, и ни один из них не цвет: поднятую
поверхность и кольцо рамки. Подпись невыбранного — `--text-secondary`,
проверена на 4.5:1 в пяти темах. Высота на `--space-2` меньше
`--control-h-md`, и в плотном режиме цель уходит ниже 24px по WCAG 2.5.8 —
как и всё размера `sm`.

## Варианты

### Со своим адресом

Значение, у которого есть адрес, не перестаёт быть значением. Окно истории
переживает перезагрузку, ложится в закладку и передаётся другому человеку
целиком — ссылкой; выбором одного из равных оно от этого быть не перестало.
Носитель состояния тогда `aria-current="page"`, а не `aria-checked`.

```html preview
<nav class="inst-segmented" aria-label="Окно истории">
  <a href="?window=day" aria-current="page">24 часа</a>
  <a href="?window=week">7 дней</a>
  <a href="?window=month">30 дней</a>
</nav>
```

## Состояния

```html preview
<div class="inst-segmented" role="radiogroup" aria-label="Выбрано первое">
  <button type="button" role="radio" aria-checked="true"  tabindex="0">Выбрано</button>
  <button type="button" role="radio" aria-checked="false" tabindex="-1">Обычное</button>
  <button type="button" role="radio" aria-checked="false" tabindex="-1" disabled>Недоступно</button>
</div>
```

| Состояние | Как ставится | Что происходит |
|---|---|---|
| выбрано | `aria-checked="true"` — или `aria-current="page"` у ссылок | Поднятая поверхность и рамка волоском. Начертание не меняется |
| обычное | `aria-checked="false"` | Подпись `--text-secondary`, фона нет |
| наведение | `:hover` | Меняется только цвет подписи: поверхность занята выбранным |
| фокус | `:focus-visible` | Кольцо внутри контейнера, без собственного радиуса |
| недоступно | `disabled` | Прозрачность. Вариант остаётся видимым — важно, **какой именно** недоступен |

Выбранный несёт **два** признака, и ни один из них не цвет: поднятую
поверхность и кольцо рамки. Различать оттенки не требуется.

Начертание состояние не несёт. Замерено: с `medium` ширина всей дорожки зависела
от того, какой сегмент выбран — 217.16, 216.94 и 216.64 пикселя, — то есть выбор
сдвигал соседей по тулбару, а правый скруглённый угол садился на полпиксель и
перерастрировался. При плотности 2.0 это целый устройственный пиксель. Без
`medium` все три положения дают одну ширину.

## JS

Модуль подключается [один раз на страницу](../../foundations/behavior.md) —
инициализировать компоненты по отдельности не нужно.

### Что делает `instrument.js`

Примеры на этой странице живые: войдите в контрол `Tab`, дальше `←` и `→`. Мышью — тоже.

Контрол объявлен как `role="radiogroup"`, и `instrument.js` выполняет её контракт: `←` и `→`
между вариантами, `Home` и `End`, один `Tab` на всю группу. Отметка
`aria-checked` следует за фокусом и переносится щелчком — сразу на всю группу,
потому что два отмеченных варианта это состояние, из которого разметка уже не
выйдет.

### События

`inst:select` всплывает с выбранного варианта.

```js
group.addEventListener('inst:select', (e) => {
  document.documentElement.dataset.density = e.detail.value;
});
```

### Опции

| Атрибут | Что делает |
|---|---|
| `data-value` | Значение в `detail` вместо подписи варианта |
| `aria-orientation="vertical"` | Стрелки по вертикали |

```html
<div class="inst-segmented" role="radiogroup" aria-label="Плотность">
  <button type="button" role="radio" aria-checked="true"  data-value="sm" tabindex="0">Плотно</button>
  <button type="button" role="radio" aria-checked="false" data-value="md" tabindex="-1">Обычно</button>
</div>
```

## Композиции

### В шапке экрана

```html preview context
<div class="inst-page-header">
  <div class="inst-page-header-main">
    <h2 class="inst-page-title">Прогоны</h2>
  </div>
  <div class="inst-page-actions">
    <div class="inst-segmented" role="radiogroup" aria-label="Период">
      <button type="button" role="radio" aria-checked="false" tabindex="-1">День</button>
      <button type="button" role="radio" aria-checked="true"  tabindex="0">Неделя</button>
      <button type="button" role="radio" aria-checked="false" tabindex="-1">Месяц</button>
    </div>
  </div>
</div>
```

## API

```api
```

## Связанное

```related
```
