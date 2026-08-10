---
title: Сегментированный контрол
group: Действия
layout: component
source: src/components.css
js: Выбор, стрелки и бегущий tabindex — делает кит
api:
  - { name: "inst-segmented", kind: "класс", doc: "Контейнер. Внутри — голые `button`, свой класс им не нужен" }
  - { name: "aria-checked", kind: "атрибут", doc: "`true` · `false`. Ровно один `true` в группе" }
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
  <button type="button" role="radio" aria-checked="true">Список</button>
  <button type="button" role="radio" aria-checked="false">Сетка</button>
  <button type="button" role="radio" aria-checked="false">Таблица</button>
</div>
```

## Использование

Кит рисует состояние, но не выдумывает его. Без ролей `aria-checked`
**невалиден**, а с клавиатуры до контрола не добраться.

```html
<div class="inst-segmented" role="radiogroup" aria-label="Плотность">
  <button type="button" role="radio" aria-checked="true"  tabindex="0">Плотно</button>
  <button type="button" role="radio" aria-checked="false" tabindex="-1">Обычно</button>
</div>
```

| Что | Почему обязательно |
|---|---|
| `role="radiogroup"` на контейнере | Группа — один контрол с одним значением |
| `aria-label` на контейнере | Иначе у контрола нет имени: подписи есть только у вариантов |
| `role="radio"` на каждой кнопке | `role="tab"` без панели — невалидная роль, и раньше стили висели именно на ней |
| `aria-checked` на каждой | Носитель состояния. Стиль выбранного висит на нём, а не на классе |
| `type="button"` | Иначе внутри формы кнопка её отправит |
| Бегущий `tabindex` | `0` у выбранной, `-1` у остальных: `Tab` входит в группу один раз, дальше стрелки |

:::warn
**Нужен JS.** Кит даёт стили и роли. Переключение `aria-checked`, стрелки и
бегущий `tabindex` — на слое приложения. Это часть контракта, а не умолчание.
:::

## Когда использовать

| Используйте | Возьмите другое |
|---|---|
| 2–4 равноправных варианта одного параметра: плотность, режим показа | **Больше пяти вариантов** — [селект](../inputs/select.md): сегментированный расползётся по ширине |
| Выбор виден целиком и переключается часто | **Переключение вкл/выкл** — свитч: у одного варианта нет пары |
| Варианты называются одним-двумя словами | **Переход между разделами** — [вкладки](../navigation/tabs.md) на ссылках: это роутинг, а не выбор значения |
| — | **Разные действия вплотную** — [группа кнопок](./button-group.md) |

## Состояния

```html preview
<div class="inst-segmented" role="radiogroup" aria-label="Выбрано первое">
  <button type="button" role="radio" aria-checked="true">Выбрано</button>
  <button type="button" role="radio" aria-checked="false">Обычное</button>
  <button type="button" role="radio" aria-checked="false" disabled>Недоступно</button>
</div>
```

| Состояние | Как ставится | Что происходит |
|---|---|---|
| выбрано | `aria-checked="true"` | Поднятая поверхность, рамка волоском, начертание `medium` |
| обычное | `aria-checked="false"` | Подпись `--text-secondary`, фона нет |
| наведение | `:hover` | Меняется только цвет подписи: поверхность занята выбранным |
| фокус | `:focus-visible` | Кольцо внутри контейнера, без собственного радиуса |
| недоступно | `disabled` | Прозрачность. Вариант остаётся видимым — важно, **какой именно** недоступен |

Выбранный несёт **три** признака: поверхность, начертание и рамку. Один тон
различать не требуется.

## JS

Подключите модуль один раз на страницу — инициализировать компоненты по
отдельности не нужно, кит работает делегированием и видит узлы, пришедшие
позже.

```html
<script type="module" src="src/kit.js"></script>
```

### Что делает кит

Примеры на этой странице живые: войдите в контрол `Tab`, дальше `←` и `→`. Мышью — тоже.

Контрол объявлен как `role="radiogroup"`, и кит выполняет обещание: `←` и `→`
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
      <button type="button" role="radio" aria-checked="false">День</button>
      <button type="button" role="radio" aria-checked="true">Неделя</button>
      <button type="button" role="radio" aria-checked="false">Месяц</button>
    </div>
  </div>
</div>
```

## Правила

:::do Два-четыре варианта
Выбор виден целиком, и переключение стоит одного нажатия. Это единственное
преимущество перед селектом.
:::

:::dont Шесть вариантов в ряд
Контрол расползается по ширине и перестаёт помещаться в шапку. Шесть — это
[селект](../inputs/select.md).
:::

:::do Ровно один aria-checked="true"
Носитель состояния — атрибут, а не класс. Ноль выбранных или два — невалидная
группа радио.
:::

:::dont Вкладки под видом сегментированного
Переход между разделами — это роутинг: [вкладки](../navigation/tabs.md) на
ссылках, у которых работает средний клик.
:::

## Доступность

| | |
|---|---|
| Клавиатура | `Tab` — вход в группу и выход из неё. `←` `→` `↑` `↓` — перебор вариантов, выбор следует за фокусом. Поведение — на приложении |
| Фокус | Бегущий `tabindex` обязателен: без него `Tab` пройдёт по всем вариантам подряд и группа перестанет быть одним контролом |
| Не только цвет | Выбранный несёт **три** признака: поверхность, начертание и кольцо рамки. Один тон различать не требуется |
| Контраст | Подпись невыбранного — `--text-secondary`, проверена на 4.5:1 в шести темах |
| Цель нажатия | Высота на `--space-2` меньше `--control-h-md`. В плотности `compact` уходит ниже 24px по WCAG 2.5.8 — как и всё размера `sm` |

## API

```api
```

## Связанное

[Группа кнопок](./button-group.md) [Кнопка](./button.md)
[Вкладки](../navigation/tabs.md) [Селект](../inputs/select.md)
[Плотность](../../foundations/density.md)
