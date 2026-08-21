---
title: Legend
group: Charts
layout: component
source: src/data.css
api:
  - { name: "inst-legend", kind: "class", doc: "The container list" }
  - { name: "inst-legend-item", kind: "class", doc: "An item: a mark plus a label" }
  - { name: "inst-legend-swatch", kind: "class", doc: "The square mark of a series" }
  - { name: "inst-legend-value", kind: "class", doc: "The quantity of a series" }
  - { name: "--chart-1", kind: "token", doc: "The hue of series no. 1. The order is part of the contract" }
  - { name: "--chart-2", kind: "token", doc: "The hue of series no. 2. The order is part of the contract" }
  - { name: "--chart-3", kind: "token", doc: "The hue of series no. 3. The order is part of the contract" }
  - { name: "--chart-4", kind: "token", doc: "The hue of series no. 4. The order is part of the contract" }
  - { name: "--chart-5", kind: "token", doc: "The hue of series no. 5. The order is part of the contract" }
  - { name: "--chart-6", kind: "token", doc: "The hue of series no. 6. The order is part of the contract" }
  - { name: "--series", kind: "variable", value: "--accent-mark", doc: "The receiver of the hue: the swatch of a legend reads it" }
  - { name: "--size-dot", kind: "token", doc: "The side of the swatch" }
  - { name: "--radius-xs", kind: "token", doc: "The rounding of the swatch" }
  - { name: "--text-xs", kind: "token", doc: "The size of the label" }
  - { name: "--space-3", kind: "token", doc: "The gap between the swatch and the label" }
  - { name: "--gap-inline", kind: "token", doc: "The gap between the items" }
  - { name: "--weight-medium", kind: "token", doc: "The weight of the quantity" }
---

The key to the colours of a chart. **Obligatory at two series and more**:
without it colour reports nothing, and a chart with meaningless colour is a
coloured-in picture.

```html preview
<ul class="inst-legend">
  <li class="inst-legend-item"><span class="inst-legend-swatch" style="--series:var(--chart-1)"></span>Render <span class="inst-legend-value">18.2 s</span></li>
  <li class="inst-legend-item"><span class="inst-legend-swatch" style="--series:var(--chart-2)"></span>Physics <span class="inst-legend-value">11.5 s</span></li>
  <li class="inst-legend-item"><span class="inst-legend-swatch" style="--series:var(--chart-3)"></span>Assets <span class="inst-legend-value">8.1 s</span></li>
  <li class="inst-legend-item"><span class="inst-legend-swatch" style="--series:var(--chart-4)"></span>Audio <span class="inst-legend-value">4.2 s</span></li>
</ul>
```

## Contract

A legend is a **list**. Series are enumerable, so this is a `<ul>`: a screen
reader announces "a list of four items" before it begins to read them.

| What | Required | Why |
|---|---|---|
| A `<ul>` and `<li>` | yes | Series are enumerable. A list announces their number before the reading |
| A `--series` on every swatch | yes | Otherwise every swatch is one colour — a legend with no legend |
| An order matching the chart | yes | The number of a series is a contract rather than an order of sorting |
| A quantity or a distinguishable label | yes | Colour has no right to be the only mark of a series |

### Accessibility

| | |
|---|---|
| The structure | A `<ul>` — the number of series is announced before the enumeration |
| Colour is not the only carrier | A series has a label in words and, as a rule, a quantity. A legend where the series are told apart by the tone of the swatch alone works neither under colour blindness nor in black-and-white print |
| Contrast | The swatch holds 3:1 against the surface: it is a load-bearing sign rather than decoration. The labels hold 4.5:1 |
| Telling the series apart | The lightness within the categorical palette **differs**: the series stay distinguishable when the hue disappears altogether |
| Type size | `--text-xs` — 12.5px. The swatch is `--size-dot`, that is 6px: making it smaller is not allowed, it would stop reading as a colour |
| The order of reading | The swatch comes before the label: a screen reader reads the label rather than an empty node, because the swatch has no text and is not announced |

## Anatomy

| Class | The work |
|---|---|
| `inst-legend` | The `<ul>`. The markers and indents of a list are removed, and the items wrap |
| `inst-legend-item` | The `<li>`. The swatch and the label in one line |
| `inst-legend-swatch` | The mark of a series. It reads `--series` |
| `inst-legend-value` | The quantity of a series. The full colour of the text and a medium weight |

The mark of a series is a **square rather than a circle**: the circle in the
library is taken by [the state dot](../display/badge.md), and a round mark of a
series would read as a status.

## Variants

The colour of a series is assigned through `--series` inline — that is a
**channel of data**: the number of a series comes with the data and cannot be
known to static CSS. The styling of the swatch — the size, the radius, the
position — is entirely in the library.

```html preview
<ul class="inst-legend">
  <li class="inst-legend-item"><span class="inst-legend-swatch" style="--series:var(--chart-1)"></span>Series 1</li>
  <li class="inst-legend-item"><span class="inst-legend-swatch" style="--series:var(--chart-2)"></span>Series 2</li>
  <li class="inst-legend-item"><span class="inst-legend-swatch" style="--series:var(--chart-3)"></span>Series 3</li>
  <li class="inst-legend-item"><span class="inst-legend-swatch" style="--series:var(--chart-4)"></span>Series 4</li>
</ul>
```

The order is part of the contract: the first series is always `--chart-1`, and
**the order in the legend has to match the order of the series on the chart**.
A legend sorted by quantity and a chart sorted by time give two different
numberings of one and the same thing.

**One series is not a category.** If there is one series, `--series` is not set
at all, and the default is `--accent-mark` rather than `--chart-1`. Categorical
colour exists to **tell series apart**; at one series there is nothing to tell
apart, and the first hue of the palette there would mean a category that does
not exist.

`--chart-1` comes into force when series no. 2 appears — and then `--series` is
set explicitly on **all** the series, the first included.

## Composition

### With quantities

```html preview
<ul class="inst-legend">
  <li class="inst-legend-item"><span class="inst-legend-swatch" style="--series:var(--chart-1)"></span>Render <span class="inst-legend-value">18.2 s</span></li>
  <li class="inst-legend-item"><span class="inst-legend-swatch" style="--series:var(--chart-2)"></span>Physics <span class="inst-legend-value">11.5 s</span></li>
</ul>
```

The quantity in a legend is a **second carrier**: it gives the exact number to
somebody who does not tell the tones apart and makes the legend readable apart
from the chart. That is precisely why `inst-legend-value` takes the full colour
of the text and the label the muted one: what gets scanned is the numbers.

### Under a chart in a card

```html preview context
<div class="inst-card inst-stack inst-stack--tight">
  <div class="inst-card-title">Time by module</div>
  <svg class="inst-sparkline" viewBox="0 0 120 32" preserveAspectRatio="none" aria-hidden="true">
    <path class="inst-sparkline-line" d="M0,26 24,24 48,18 72,13 96,9 120,6"/>
  </svg>
  <ul class="inst-legend">
    <li class="inst-legend-item"><span class="inst-legend-swatch" style="--series:var(--chart-1)"></span>Render <span class="inst-legend-value">18.2 s</span></li>
    <li class="inst-legend-item"><span class="inst-legend-swatch" style="--series:var(--chart-2)"></span>Physics <span class="inst-legend-value">11.5 s</span></li>
  </ul>
</div>
```

The legend stands **under** the chart: above it, it pushes the data below the
fold for the sake of what is read second.

## API

```api
```

## Related

```related
```
