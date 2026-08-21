---
title: Categorical palette
group: Charts
layout: component
source: src/tokens.css
api:
  - { name: "--chart-1", kind: "token", doc: "The hue of series no. 1. The order is part of the contract: series no. 1 is always this hue" }
  - { name: "--chart-2", kind: "token", doc: "The hue of series no. 2. The order is part of the contract: series no. 2 is always this hue" }
  - { name: "--chart-3", kind: "token", doc: "The hue of series no. 3. The order is part of the contract: series no. 3 is always this hue" }
  - { name: "--chart-4", kind: "token", doc: "The hue of series no. 4. The order is part of the contract: series no. 4 is always this hue" }
  - { name: "--chart-5", kind: "token", doc: "The hue of series no. 5. The order is part of the contract: series no. 5 is always this hue" }
  - { name: "--chart-6", kind: "token", doc: "The hue of series no. 6. The order is part of the contract: series no. 6 is always this hue" }
  - { name: "--series", kind: "variable", value: "--accent-mark", doc: "The receiver of the hue. A component of a series reads it rather than the hue directly" }
---

Six hues telling **the series of a chart** apart. The one place in the library
where colour does not encode a state — and therefore the one place with rules
of its own.

```html preview
<ul class="inst-legend">
  <li class="inst-legend-item"><span class="inst-legend-swatch" style="--series:var(--chart-1)"></span>Series 1</li>
  <li class="inst-legend-item"><span class="inst-legend-swatch" style="--series:var(--chart-2)"></span>Series 2</li>
  <li class="inst-legend-item"><span class="inst-legend-swatch" style="--series:var(--chart-3)"></span>Series 3</li>
  <li class="inst-legend-item"><span class="inst-legend-swatch" style="--series:var(--chart-4)"></span>Series 4</li>
  <li class="inst-legend-item"><span class="inst-legend-swatch" style="--series:var(--chart-5)"></span>Series 5</li>
  <li class="inst-legend-item"><span class="inst-legend-swatch" style="--series:var(--chart-6)"></span>Series 6</li>
</ul>
```

The first law asks that colour encode a state, a category or a polarity. A
series of a chart is a **category**, so the palette does not break the law but
refines it. The exception is named in
[the design principles](../../about/design-principles.md), and along with it
are named five rules that exist nowhere else in the library.

## Contract

A series is assigned through `--series`, inline, **as a channel of data**: the
number of a series comes with the data and cannot be known to static CSS.

```html preview
<svg class="inst-sparkline" viewBox="0 0 120 32" preserveAspectRatio="none" aria-hidden="true" style="--series:var(--chart-3)">
  <path class="inst-sparkline-area" d="M0,20 24,14 48,18 72,9 96,12 120,5 120,32 0,32Z"/>
  <path class="inst-sparkline-line" d="M0,20 24,14 48,18 72,9 96,12 120,5"/>
</svg>
```

| What | Required | Why |
|---|---|---|
| A `--series` on the root of a series | yes | Once on the `<svg>` or on the mark of a legend — every node of the series inherits it |
| The hue through `--series` rather than directly | yes | A component reads the receiver rather than the hue of a series: otherwise reassigning a series takes editing the markup of every node |
| A legend at two series and more | yes | With no key colour reports nothing |
| A `--series` at one series | no | The default is `--accent-mark`: one series is not a category, and there is nothing to tell apart |

The components that read `--series`:

| Component | The nodes |
|---|---|
| [The legend](./legend.md) | `inst-legend-swatch` |
| [The sparkline](./sparkline.md) | `inst-sparkline-line` · `inst-sparkline-area` · `inst-sparkline-last` |
| [The avatar](../display/avatar.md) | The hue of a group: an avatar tells **entities** apart rather than reporting a state — the same work a series does |

[The ring](./ring.md) and [the meter](./meter.md) **do not read** `--series`:
they show one quantity and take `--tone-mark`.

```css
/* A named series of the application — the number tied to a measure in one place */
.chart-render  { --series: var(--chart-1); }
.chart-physics { --series: var(--chart-2); }
```

Such a tie is the cheapest way to keep the first rule: a measure gets its
number once rather than in every template.

### Accessibility

| | |
|---|---|
| Colour is not the only carrier | A series is always named in words in the legend. The palette tells series apart but reports nothing about them — the label does |
| Colour blindness | The differing lightness within the palette was chosen for this: when hue is lost the series stay distinguishable by brightness |
| Black-and-white print | The same mechanism. Six hues of one lightness would merge into one grey |
| Contrast | Every hue holds 3:1 against the surface in five themes: the line and the mark of a series are load-bearing signs |
| Labels | A categorical hue is **not used for text**: 3:1 is enough for a sign but not for a label. Text takes `--text-primary` or `--tone-ink` |
| Forced-colours mode | The categorical hues are reset by the system, and that is right: in such a mode it is the legend with its words that has to tell the series apart |

## Variants

| Token | Hue | Lightness (light theme) | Lightness (dark) |
|---|---|---|---|
| `--chart-1` | 292 | 0.520 | 0.700 |
| `--chart-2` | 320 | 0.600 | 0.760 |
| `--chart-3` | 355 | 0.550 | 0.720 |
| `--chart-4` | 55 | 0.650 | 0.800 |
| `--chart-5` | 115 | 0.580 | 0.740 |
| `--chart-6` | 190 | 0.620 | 0.780 |

The values are written in OKLCH through `light-dark()`: the hue and the
saturation hold while the lightness turns over with the theme. The palette
stands at least 25° from the status hues: the corridors `25` (error), `85`
(warning) and `150` (success) are reserved for good, and a series landing in
them reads as "it failed" or "done".

Series no. 3 stands at `355`, that is 30° from the error hue at `25`: that is
the closest approach to a status hue in the whole palette, and widening it
further is impossible not only for the sake of telling them apart but because
there are no free corridors left.

**Against the accent there is no such guarantee and there can be none.** The
accent is [a knob with four positions](../../foundations/colors.md): petrol
`215`, graphite `250`, indigo `275`, clay `45`. Four positions cover too much
of the circle: clay stands 10° from series no. 4 (`55`). In place of an angular
guarantee what works here is a division of TERRITORY — the accent does not step
inside a chart, a categorical colour does not step outside. The one place they
meet on one screen is the legend of a multi-series chart beside the interface,
and what tells them apart there is lightness.

## API

```api
```

## Related

```related
```
