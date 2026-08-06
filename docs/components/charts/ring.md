---
title: Кольцо
group: Графики
status: stable
source: src/data.css
---

Та же доля, что у [меры](./meter.md), но в габарите контрола: кольцо встаёт
рядом с текстом карточки, где полосе на всю ширину места нет. SVG, без JS.

```html preview
<span class="inst-ring-wrap">
  <svg class="inst-ring" viewBox="0 0 20 20" role="img" aria-label="Израсходовано 74% бюджета">
    <circle class="inst-ring-track" cx="10" cy="10" r="8"/>
    <circle class="inst-ring-fill" cx="10" cy="10" r="8" style="--value:0.74"/>
  </svg>
  <span class="inst-ring-label">74</span>
</span>
```

## Установка

```html
<link rel="stylesheet" href="src/kit.css">
```

## Когда использовать

| Используйте | Возьмите другое |
|---|---|
| Доля от целого в тесном месте: рядом с заголовком карточки, в строке | **Доля во всю ширину блока** — [мера](./meter.md): полоса длиннее, значит точнее |
| Одна величина, у которой есть потолок | **Несколько долей, складывающихся в целое** — это не кольцо, а [легенда](./legend.md) с рядами: кольцо кита однодуговое |
| Число внутри кольца как главный носитель значения | **Число само по себе** — [метрика](../display/metric.md): кольцо без величины у потолка ничего не добавляет |
| Статичное значение, приходящее из данных | **Неизвестная длительность** — [мера](./meter.md) с `data-state="indeterminate"`: у кольца нет неопределённого вида |

## Разметка

Три обязательные части, и ни одну нельзя опустить.

| Узел | Работа |
|---|---|
| `inst-ring-wrap` | Позиционирующая обёртка. Только она умеет держать подпись по центру |
| `inst-ring` | Сам `<svg>`. Повёрнут на `-90°`, поэтому дуга начинается сверху, а не справа |
| `inst-ring-track` | Дорожка — полная окружность |
| `inst-ring-fill` | Дуга значения |
| `inst-ring-label` | Число по центру. Табличные цифры, чтобы не дёргалось при обновлении |

`viewBox="0 0 20 20"` и `r="8"` **не произвольны**: длина окружности при
`r=8` равна 50.27, и пунктир дуги считается прямо из неё. Другой радиус —
другая длина, и дуга перестанет соответствовать значению. Меняйте размер
кольца снаружи (`inline-size`/`block-size`), а не радиусом внутри.

## Значение

```html preview
<span class="inst-ring-wrap">
  <svg class="inst-ring" viewBox="0 0 20 20" role="img" aria-label="Готово 25%">
    <circle class="inst-ring-track" cx="10" cy="10" r="8"/>
    <circle class="inst-ring-fill" cx="10" cy="10" r="8" style="--value:0.25"/>
  </svg>
  <span class="inst-ring-label">25</span>
</span>
<span class="inst-ring-wrap">
  <svg class="inst-ring" viewBox="0 0 20 20" role="img" aria-label="Готово 60%">
    <circle class="inst-ring-track" cx="10" cy="10" r="8"/>
    <circle class="inst-ring-fill" cx="10" cy="10" r="8" style="--value:0.6"/>
  </svg>
  <span class="inst-ring-label">60</span>
</span>
<span class="inst-ring-wrap">
  <svg class="inst-ring" viewBox="0 0 20 20" role="img" aria-label="Готово 96%">
    <circle class="inst-ring-track" cx="10" cy="10" r="8"/>
    <circle class="inst-ring-fill" cx="10" cy="10" r="8" style="--value:0.96"/>
  </svg>
  <span class="inst-ring-label">96</span>
</span>
```

Значение — **доля 0…1** в `--value`, а не проценты и не градусы: доля
единственная величина, из которой пунктир считается без пересчёта.

`style="--value:0.74"` — инлайновый стиль как **канал данных**, а не
оформление. Значение приходит с сервера и не может жить в статическом CSS;
всё остальное — толщина, скругление, цвет, поворот — берётся из кита.

## Тон

```html preview
<span class="inst-ring-wrap">
  <svg class="inst-ring" viewBox="0 0 20 20" role="img" aria-label="Расход 88%, у порога">
    <circle class="inst-ring-track" cx="10" cy="10" r="8"/>
    <circle class="inst-ring-fill" cx="10" cy="10" r="8" style="--value:0.88" data-tone="warn"/>
  </svg>
  <span class="inst-ring-label">88</span>
</span>
```

Дуга читает `--tone-mark` и без атрибута берёт `--accent-mark`. Как и у меры,
тон меняется, когда меняется **смысл величины**, а не для разнообразия.
Категориальные тона `--chart-*` дуге не положены: кольцо показывает одну
величину, а не ряд среди равных — см. [палитру](./palette.md).

## Композиции

### В карточке

```html
<div class="inst-card inst-cluster inst-cluster--loose">
  <span class="inst-ring-wrap">
    <svg class="inst-ring" viewBox="0 0 20 20" role="img" aria-label="Израсходовано 74% бюджета">
      <circle class="inst-ring-track" cx="10" cy="10" r="8"/>
      <circle class="inst-ring-fill" cx="10" cy="10" r="8" style="--value:0.74"/>
    </svg>
    <span class="inst-ring-label">74</span>
  </span>
  <div>
    <div class="inst-card-title">Бюджет израсходован</div>
    <div class="inst-card-sub">184 320 из 250 000 токенов</div>
  </div>
</div>
```

### Без подписи внутри

```html
<span class="inst-ring-wrap">
  <svg class="inst-ring" viewBox="0 0 20 20" role="img" aria-label="Готово 40%">
    <circle class="inst-ring-track" cx="10" cy="10" r="8"/>
    <circle class="inst-ring-fill" cx="10" cy="10" r="8" style="--value:0.4"/>
  </svg>
</span>
```

Обёртка остаётся: она даёт `inline-flex` и снимает лишний зазор строки под
SVG. `aria-label` остаётся тем более — без подписи внутри он единственный
носитель значения для скринридера.

## Справочник

### Классы

| Класс | Работа |
|---|---|
| `inst-ring-wrap` | Обёртка. Позиционирует подпись |
| `inst-ring` | `<svg>`. Размер — `--control-h-lg`, поворот `-90°` |
| `inst-ring-track` | Полная окружность дорожки |
| `inst-ring-fill` | Дуга значения. Читает `--value` |
| `inst-ring-label` | Подпись по центру |

### Обязательная разметка

| Что | Почему |
|---|---|
| `role="img"` на `<svg>` | Иначе у графики нет роли, и `aria-label` на ней не гарантирован |
| `aria-label` со значением словами | «74» внутри кольца — это цифра без единицы и без контекста |
| `viewBox="0 0 20 20"`, `cx="10" cy="10" r="8"` | Пунктир дуги посчитан из длины окружности при `r=8` |
| `--value` в долях 0…1 | Не проценты и не градусы |

### CSS-переменные

| Переменная | По умолчанию | Что делает |
|---|---|---|
| `--value` | `0` | Доля 0…1. Ставится инлайном как данные |
| `--tone-mark` | `--accent-mark` | Цвет дуги. Приходит из `data-tone` |
| `--track` | — | Цвет дорожки |

```css
/* Кольцо крупнее — снаружи, а не радиусом внутри */
.my-ring { inline-size: 48px; block-size: 48px; }
```

### Токены

`--control-h-lg` · `--track` · `--tone-mark` · `--accent-mark` ·
`--text-2xs` · `--radius-full`

## Доступность

| | |
|---|---|
| Роль | `role="img"` — в примере, а не в примечании |
| Имя | `aria-label` предложением: «Израсходовано 74% бюджета», а не «74» |
| Цвет не единственный носитель | Значение продублировано числом в `inst-ring-label` либо в подписи рядом. Дуга без числа читается только зрячим и только на цвет |
| Контраст | Дуга держит 3:1 против дорожки. Дорожка не несущая, поэтому 3:1 к поверхности с неё не спрашивается |
| Кегль подписи | `--text-2xs` — 11px, пол кита. Кольцу меньше `--control-h-lg` подпись внутрь не помещается: выносите её наружу |
| Цифры | `font-variant-numeric: tabular-nums` — обновление значения не дёргает центровку |

## Связанное

[Мера](./meter.md) · [Спарклайн](./sparkline.md) · [Легенда](./legend.md) ·
[Палитра](./palette.md) · [Метрика](../display/metric.md) ·
[Карточка](../display/card.md) · [Бюджет прогона](../../agent/budget.md)

Исходник: `src/data.css` · Почему кит устроен именно так — [конституция](../../about/design-principles.md)
