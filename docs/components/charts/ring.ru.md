---
title: Кольцо
group: Графики
layout: component
source: src/data.css
api:
  - { name: "inst-ring-wrap", kind: "класс", doc: "Обёртка. Позиционирует подпись" }
  - { name: "inst-ring", kind: "класс", doc: "`<svg>`. Размер — `--control-h-lg`, поворот `-90°`" }
  - { name: "inst-ring-track", kind: "класс", doc: "Полная окружность дорожки" }
  - { name: "inst-ring-fill", kind: "класс", doc: "Дуга значения. Читает `--value`" }
  - { name: "inst-ring-label", kind: "класс", doc: "Подпись по центру" }
  - { name: "--value", kind: "переменная", value: "0", doc: "Доля 0…1. Ставится инлайном как данные" }
  - { name: "--tone-mark", kind: "переменная", value: "--accent-mark", doc: "Цвет дуги. Приходит из `data-tone`" }
  - { name: "--track", kind: "переменная", value: "—", doc: "Цвет дорожки" }
  - { name: "--control-h-lg", kind: "токен" }
  - { name: "--accent-mark", kind: "токен" }
  - { name: "--text-2xs", kind: "токен" }
  - { name: "--radius-full", kind: "токен" }
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

## Контракт

| Что | Обязательно | Почему |
|---|---|---|
| `role="img"` на `<svg>` | да | Иначе у графики нет роли, и `aria-label` на ней не гарантирован |
| `aria-label` со значением словами | да | «74» внутри кольца — это цифра без единицы и без контекста |
| `viewBox="0 0 20 20"`, `cx="10" cy="10" r="8"` | да | Пунктир дуги посчитан из длины окружности при `r=8` |
| `--value` в долях 0…1 | да | Не проценты и не градусы |

`viewBox="0 0 20 20"` и `r="8"` **не произвольны**: длина окружности при
`r=8` равна 50.27, и пунктир дуги считается прямо из неё. Другой радиус —
другая длина, и дуга перестанет соответствовать значению.

```css
/* Кольцо крупнее — снаружи, а не радиусом внутри */
.my-ring { inline-size: 48px; block-size: 48px; }
```

### Доступность

| | |
|---|---|
| Роль | `role="img"` — в примере, а не в примечании |
| Имя | `aria-label` предложением: «Израсходовано 74% бюджета», а не «74» |
| Цвет не единственный носитель | Значение продублировано числом в `inst-ring-label` либо в подписи рядом. Дуга без числа читается только зрячим и только на цвет |
| Контраст | Дуга держит 3:1 против дорожки. Дорожка не несущая, поэтому 3:1 к поверхности с неё не спрашивается |
| Кегль подписи | `--text-2xs` — 11px, нижний предел кегля. Кольцу меньше `--control-h-lg` подпись внутрь не помещается: выносите её наружу |
| Цифры | `font-variant-numeric: tabular-nums` — обновление значения не дёргает центровку |

## Устройство

| Узел | Работа |
|---|---|
| `inst-ring-wrap` | Позиционирующая обёртка. Только она умеет держать подпись по центру |
| `inst-ring` | Сам `<svg>`. Повёрнут на `-90°`, поэтому дуга начинается сверху, а не справа |
| `inst-ring-track` | Дорожка — полная окружность |
| `inst-ring-fill` | Дуга значения |
| `inst-ring-label` | Число по центру. Табличные цифры, чтобы не дёргалось при обновлении |

Значение — **доля 0…1** в `--value`. Ни проценты, ни градусы: из доли пунктир
считается без пересчёта.
`style="--value:0.74"` — инлайновый стиль как **канал данных**, а не
оформление: значение приходит с сервера и не может жить в статическом CSS.

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

## Варианты

```html preview
<span class="inst-ring-wrap">
  <svg class="inst-ring" viewBox="0 0 20 20" role="img" aria-label="Расход 62%">
    <circle class="inst-ring-track" cx="10" cy="10" r="8"/>
    <circle class="inst-ring-fill" cx="10" cy="10" r="8" style="--value:0.62"/>
  </svg>
  <span class="inst-ring-label">62</span>
</span>
<span class="inst-ring-wrap">
  <svg class="inst-ring" viewBox="0 0 20 20" role="img" aria-label="Расход 88%, у порога">
    <circle class="inst-ring-track" cx="10" cy="10" r="8"/>
    <circle class="inst-ring-fill" cx="10" cy="10" r="8" style="--value:0.88" data-tone="warn"/>
  </svg>
  <span class="inst-ring-label">88</span>
</span>
<span class="inst-ring-wrap">
  <svg class="inst-ring" viewBox="0 0 20 20" role="img" aria-label="Лимит исчерпан">
    <circle class="inst-ring-track" cx="10" cy="10" r="8"/>
    <circle class="inst-ring-fill" cx="10" cy="10" r="8" style="--value:1" data-tone="error"/>
  </svg>
  <span class="inst-ring-label">100</span>
</span>
```

Дуга читает `--tone-mark` и без атрибута берёт `--accent-mark`. Как и у меры,
тон меняется вслед за **смыслом величины**.
Категориальные тона `--chart-*` дуге не положены: кольцо показывает одну
величину, а не ряд среди равных — см. [палитру](./palette.md).

## Композиции

### В карточке

```html preview context
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

```html preview
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

## API

```api
```

## Связанное

```related
```
