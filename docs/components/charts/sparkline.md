---
title: Sparkline
group: Charts
layout: component
source: src/data.css
api:
  - { name: "inst-sparkline", kind: "class", doc: "The `<svg>`. Required: the other classes work only inside it" }
  - { name: "inst-sparkline-line", kind: "class", doc: "The line of the series, a stroke of 1.5, non-scaling" }
  - { name: "inst-sparkline-area", kind: "class", doc: "The fill under the line, alpha `0.12`" }
  - { name: "inst-sparkline-last", kind: "class", doc: "The dot of the last value with a ring the colour of the surface" }
  - { name: "--chart-1", kind: "token", doc: "The hue of series no. 1. The order is part of the contract" }
  - { name: "--chart-2", kind: "token", doc: "The hue of series no. 2. The order is part of the contract" }
  - { name: "--chart-3", kind: "token", doc: "The hue of series no. 3. The order is part of the contract" }
  - { name: "--chart-4", kind: "token", doc: "The hue of series no. 4. The order is part of the contract" }
  - { name: "--chart-5", kind: "token", doc: "The hue of series no. 5. The order is part of the contract" }
  - { name: "--chart-6", kind: "token", doc: "The hue of series no. 6. The order is part of the contract" }
  - { name: "--series", kind: "variable", value: "--accent-mark", doc: "The receiver of the hue: the line, the fill and the last dot read it" }
  - { name: "--control-h-md", kind: "token", doc: "The height is counted from it: ×1.25" }
  - { name: "--surface-raised", kind: "token", doc: "The ring of the last dot is the colour of the surface under it" }
  - { name: "--size-marker", kind: "token", doc: "The thickness of the ring of the last dot" }
---

The shape of a change beside a number. Deliberately **with no axes, grid or
labels**: this is not a chart but the second line of a metric. SVG, with no JS.

```html preview
<svg class="inst-sparkline" viewBox="0 0 120 32" preserveAspectRatio="none" aria-hidden="true">
  <path class="inst-sparkline-area" d="M0,27 5,24 10,26 15,21 20,23 25,18 30,20 35,14 40,17 45,11 50,15 55,9 60,13 65,6 70,10 75,5 80,12 85,8 90,11 95,4 100,7 105,3 110,6 115,4 120,2 120,32 0,32Z"/>
  <path class="inst-sparkline-line" d="M0,27 5,24 10,26 15,21 20,23 25,18 30,20 35,14 40,17 45,11 50,15 55,9 60,13 65,6 70,10 75,5 80,12 85,8 90,11 95,4 100,7 105,3 110,6 115,4 120,2"/>
</svg>
```

## Contract

The path comes from the server or from a template. The library answers for the
hues and the thicknesses and **computes nothing**: the `d` is data, exactly as
the width of a meter's fill is.

| What | Required | Why |
|---|---|---|
| A `viewBox` | yes | Without it the SVG does not scale, and the path travels outside the bounds |
| An `aria-hidden="true"` **or** a `role="img"` with an `aria-label` | yes | There is no third way: either the picture is decorative and the number beside it carries the value, or it has a name |
| Direct children of the `<svg>` | yes | The library's rules are written through `& >`: a node wrapped in a `<g>` loses its hue and its thickness |
| `inst-sparkline-area` and `-last` | no | The fill and the dot of the last value are optional |

`preserveAspectRatio="none"`: a sparkline stretches across the width of its
container and keeps its height fixed. That is precisely why the stroke is
declared non-scaling — otherwise a stretched line would grow thicker from left
to right.

### Accessibility

| | |
|---|---|
| The name | A decorative sparkline is `aria-hidden="true"`, and then the value has to be in the text beside it. A self-standing one is `role="img"` plus an `aria-label` with the outcome in words: "A rise from 92 to 184 thousand over 11 days" |
| Colour is not the only carrier | There is one sparkline — colour encodes nothing. Two series in one colour cannot be told apart, which is why a sparkline has no two series |
| Contrast | The line takes a categorical hue and holds 3:1 against the surface in five themes. The fill at alpha `0.12` is decorative and is not checked against a threshold: the value is carried by the line |
| Thickness | `vector-effect: non-scaling-stroke` — the stroke stays 1.5 at any stretching, so it does not thin below visibility |
| Motion | There is no animation at all: a chart is a state rather than a process |

## Anatomy

| Node | Required | The work |
|---|---|---|
| `inst-sparkline` | yes | The `<svg>` with a `viewBox`. Width 100%, height a quarter above a control |
| `inst-sparkline-line` | yes | The line. `vector-effect: non-scaling-stroke`, so stretching across the width does not thicken the stroke |
| `inst-sparkline-area` | no | The fill under the line. The same path, closed down to the bottom edge |
| `inst-sparkline-last` | no | The dot of the last value |

```html preview
<svg class="inst-sparkline" viewBox="0 0 120 32" preserveAspectRatio="none" aria-hidden="true">
  <path class="inst-sparkline-area" d="M0,27 5,24 10,26 15,21 20,23 25,18 30,20 35,14 40,17 45,11 50,15 55,9 60,13 65,6 70,10 75,5 80,12 85,8 90,11 95,4 100,7 105,3 110,6 115,4 120,2 120,32 0,32Z"/>
  <path class="inst-sparkline-line" d="M0,27 5,24 10,26 15,21 20,23 25,18 30,20 35,14 40,17 45,11 50,15 55,9 60,13 65,6 70,10 75,5 80,12 85,8 90,11 95,4 100,7 105,3 110,6 115,4 120,2"/>
  <circle class="inst-sparkline-last" cx="120" cy="6" r="2.5"/>
</svg>
```

The order of the nodes inside the SVG is the order of drawing: the fill first,
the line over it, the dot last. Swap them and the line goes under its own fill.

The fill under the line is the **one gradient allowed in the library, and it is
not a gradient**: a solid hue at alpha `0.12`. A real gradient would report
that the quantity changes on the vertical, which it does not.

The dot of the last value carries a ring the colour of the surface and so stays
distinguishable when it runs over its own line. Its coordinates are the same as
the last point of the path: they part in silence.

## Variants

```html preview
<svg class="inst-sparkline" viewBox="0 0 120 32" preserveAspectRatio="none" aria-hidden="true" style="--series:var(--chart-2)">
  <path class="inst-sparkline-area" d="M0,10 12,14 24,11 36,18 48,15 60,21 72,19 84,24 96,22 108,27 120,25 120,32 0,32Z"/>
  <path class="inst-sparkline-line" d="M0,10 12,14 24,11 36,18 48,15 60,21 72,19 84,24 96,22 108,27 120,25"/>
</svg>
```

All three nodes read `--series` and without it take `--chart-1`. The variable
is set **once on the `<svg>`**. The line, the fill and the dot of one series
have no right to part in colour.

The number of a series is part of the contract: the first series is always
`--chart-1`. The rules of categorical colour entire are on the page of
[the palette](./palette.md).

```css
/* The series is set from outside rather than by rewriting stroke */
.my-chart-2 { --series: var(--chart-2); }
```

## Composition

### Under a metric

```html preview context
<div class="inst-card inst-stack inst-stack--tight">
  <div class="inst-metric-label">Tokens per day</div>
  <div class="inst-metric-value">184<span class="inst-metric-unit">&nbsp;k</span></div>
  <svg class="inst-sparkline" viewBox="0 0 120 32" preserveAspectRatio="none" aria-hidden="true">
    <path class="inst-sparkline-area" d="M0,27 5,24 10,26 15,21 20,23 25,18 30,20 35,14 40,17 45,11 50,15 55,9 60,13 65,6 70,10 75,5 80,12 85,8 90,11 95,4 100,7 105,3 110,6 115,4 120,2 120,32 0,32Z"/>
    <path class="inst-sparkline-line" d="M0,27 5,24 10,26 15,21 20,23 25,18 30,20 35,14 40,17 45,11 50,15 55,9 60,13 65,6 70,10 75,5 80,12 85,8 90,11 95,4 100,7 105,3 110,6 115,4 120,2"/>
  </svg>
  <div class="inst-metric-label">Over 11 days</div>
</div>
```

### In a table cell

```html
<td>
  <svg class="inst-sparkline" viewBox="0 0 120 32" preserveAspectRatio="none" aria-hidden="true">
    <path class="inst-sparkline-line" d="M0,26 24,24 48,18 72,13 96,9 120,6"/>
  </svg>
</td>
```

## API

```api
```

## Related

```related
```
