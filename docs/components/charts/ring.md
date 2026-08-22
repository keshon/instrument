---
title: Ring
group: Charts
layout: component
source: src/data.css
api:
  - { name: "inst-ring-wrap", kind: "class", doc: "The wrapper. It positions the label" }
  - { name: "inst-ring", kind: "class", doc: "The `<svg>`. The size is `--control-h-lg`, turned `-90°`" }
  - { name: "inst-ring--sm", kind: "modifier", doc: "The box from `--control-h-sm`, the stroke 2.5. In a row of controls" }
  - { name: "inst-ring--xs", kind: "modifier", doc: "The box from `--icon-size` — where a glyph stands rather than where a control does. The stroke 2, the optical weight of an icon" }
  - { name: "--ring-size", kind: "variable", value: "--control-h-lg", doc: "The box. Set by the size modifiers; a caller may set it directly" }
  - { name: "--ring-stroke", kind: "variable", value: "3", doc: "The stroke, in viewBox units. It scales with the box, so it comes down separately" }
  - { name: "inst-ring-track", kind: "class", doc: "The full circle of the track" }
  - { name: "inst-ring-fill", kind: "class", doc: "The arc of the value. It reads `--value`" }
  - { name: "inst-ring-label", kind: "class", doc: "The label at the centre" }
  - { name: "--value", kind: "variable", value: "0", doc: "A share of 0…1. Set inline as data" }
  - { name: "--tone-mark", kind: "variable", value: "--accent-mark", doc: "The colour of the arc. It comes from `data-tone`" }
  - { name: "--track", kind: "variable", value: "—", doc: "The colour of the track" }
  - { name: "--control-h-lg", kind: "token" }
  - { name: "--accent-mark", kind: "token" }
  - { name: "--text-2xs", kind: "token" }
  - { name: "--radius-full", kind: "token" }
---

The same share as [a meter](./meter.md) but in the size of a control: a ring
stands beside the text of a card, where a full-width bar has no room. SVG, with
no JS.

```html preview
<span class="inst-ring-wrap">
  <svg class="inst-ring" viewBox="0 0 20 20" role="img" aria-label="74% of the budget spent">
    <circle class="inst-ring-track" cx="10" cy="10" r="8"/>
    <circle class="inst-ring-fill" cx="10" cy="10" r="8" style="--value:0.74"/>
  </svg>
  <span class="inst-ring-label">74</span>
</span>
```

## Contract

| What | Required | Why |
|---|---|---|
| A `role="img"` on the `<svg>` | standing alone | Otherwise the graphic has no role, and an `aria-label` on it is not guaranteed. **Inside a labelled control** the ring is `aria-hidden="true"` and the control carries the name — two names on one object are read twice |
| An `aria-label` with the value in words | yes | A "74" inside a ring is a figure with no unit and no context |
| `viewBox="0 0 20 20"`, `cx="10" cy="10" r="8"` | yes | The dash of the arc is computed from the circumference at `r=8` |
| A `--value` as a share of 0…1 | yes | Not percentages and not degrees |

The `viewBox="0 0 20 20"` and the `r="8"` are **not arbitrary**: the
circumference at `r=8` equals 50.27, and the dash of the arc is computed
straight from it. Another radius means another length, and the arc stops
matching the value.

```css
/* A larger ring — from outside rather than by the radius inside */
.my-ring { inline-size: 48px; block-size: 48px; }
```

### Accessibility

| | |
|---|---|
| The role | `role="img"` — in the example rather than in a note |
| The name | An `aria-label` as a sentence: "74% of the budget spent" rather than "74" |
| Colour is not the only carrier | The value is repeated as a number in the `inst-ring-label` or in the label beside it. An arc with no number is read only by the sighted and only by colour |
| Contrast | The arc holds 3:1 against the track. The track is not load-bearing, so 3:1 against the surface is not asked of it |
| The size of the label | `--text-2xs` — 11px, the lower limit of type size. A ring smaller than `--control-h-lg` does not fit a label inside: take it outside |
| Figures | `font-variant-numeric: tabular-nums` — an update of the value does not tug at the centring |

## Anatomy

| Node | The work |
|---|---|
| `inst-ring-wrap` | The positioning wrapper. It alone can hold the label at the centre |
| `inst-ring` | The `<svg>` itself. Turned `-90°`, so the arc begins at the top rather than at the right |
| `inst-ring-track` | The track — a full circle |
| `inst-ring-fill` | The arc of the value |
| `inst-ring-label` | The number at the centre. Tabular figures, so that it does not twitch on an update |

The value is a **share of 0…1** in `--value`. Neither percentages nor degrees:
from a share the dash is computed with no conversion. The
`style="--value:0.74"` is an inline style as a **channel of data** rather than
styling: the value comes from the server and cannot live in static CSS.

```html preview
<span class="inst-ring-wrap">
  <svg class="inst-ring" viewBox="0 0 20 20" role="img" aria-label="25% done">
    <circle class="inst-ring-track" cx="10" cy="10" r="8"/>
    <circle class="inst-ring-fill" cx="10" cy="10" r="8" style="--value:0.25"/>
  </svg>
  <span class="inst-ring-label">25</span>
</span>
<span class="inst-ring-wrap">
  <svg class="inst-ring" viewBox="0 0 20 20" role="img" aria-label="60% done">
    <circle class="inst-ring-track" cx="10" cy="10" r="8"/>
    <circle class="inst-ring-fill" cx="10" cy="10" r="8" style="--value:0.6"/>
  </svg>
  <span class="inst-ring-label">60</span>
</span>
<span class="inst-ring-wrap">
  <svg class="inst-ring" viewBox="0 0 20 20" role="img" aria-label="96% done">
    <circle class="inst-ring-track" cx="10" cy="10" r="8"/>
    <circle class="inst-ring-fill" cx="10" cy="10" r="8" style="--value:0.96"/>
  </svg>
  <span class="inst-ring-label">96</span>
</span>
```

## Variants

```html preview
<span class="inst-ring-wrap">
  <svg class="inst-ring" viewBox="0 0 20 20" role="img" aria-label="62% spent">
    <circle class="inst-ring-track" cx="10" cy="10" r="8"/>
    <circle class="inst-ring-fill" cx="10" cy="10" r="8" style="--value:0.62"/>
  </svg>
  <span class="inst-ring-label">62</span>
</span>
<span class="inst-ring-wrap">
  <svg class="inst-ring" viewBox="0 0 20 20" role="img" aria-label="88% spent, near the threshold">
    <circle class="inst-ring-track" cx="10" cy="10" r="8"/>
    <circle class="inst-ring-fill" cx="10" cy="10" r="8" style="--value:0.88" data-tone="warn"/>
  </svg>
  <span class="inst-ring-label">88</span>
</span>
<span class="inst-ring-wrap">
  <svg class="inst-ring" viewBox="0 0 20 20" role="img" aria-label="The limit is spent">
    <circle class="inst-ring-track" cx="10" cy="10" r="8"/>
    <circle class="inst-ring-fill" cx="10" cy="10" r="8" style="--value:1" data-tone="error"/>
  </svg>
  <span class="inst-ring-label">100</span>
</span>
```

The arc reads `--tone-mark` and without the attribute takes `--accent-mark`. As
on a meter, the tone follows the **meaning of the quantity**. The categorical
hues `--chart-*` are not due to an arc: a ring shows one quantity rather than a
series among equals — see [the palette](./palette.md).

## Sizes

| Size | Box | Stroke | Where |
|---|---|---|---|
| default | `--control-h-lg` | 3 | Beside the text of a card. The only one that fits a label inside |
| `inst-ring--sm` | `--control-h-sm` | 2.5 | In a row of controls, beside a button |
| `inst-ring--xs` | `--icon-size` | 2 | Where a **glyph** stands: the end of a session bar, beside the model and the effort |

The smallest rung is not a control height. There the ring is not a control
standing among controls — it is one more glyph in the row, so it reads
`--icon-size`, the same variable every icon in that row reads, and it moves
when the row moves. Given `--control-h-xs` instead it came out 22px against
14px glyphs and read as a third object between them.

The stroke follows the same argument: an icon is drawn at 1.5 on a grid of 16,
or 9.4%, and the ring's grid is 20 — so 2 is the same optical weight.

**The label does not come down with the box.** `--text-2xs` is the floor of the
type scale and does not fit inside 22px, let alone 14. Below `--control-h-lg`
the number goes beside the ring or into a [tooltip](../overlays/tooltip.md) —
an arc with no number anywhere is read by its length and its colour alone.

```html preview
<span class="inst-ring-wrap">
  <svg class="inst-ring inst-ring--sm" viewBox="0 0 20 20" role="img" aria-label="40% of the context used">
    <circle class="inst-ring-track" cx="10" cy="10" r="8"/>
    <circle class="inst-ring-fill" cx="10" cy="10" r="8" style="--value:0.4"/>
  </svg>
</span>
<span class="inst-tooltip">
  <span class="inst-ring-wrap" aria-describedby="tt-ring-xs">
    <svg class="inst-ring inst-ring--xs" viewBox="0 0 20 20" role="img" aria-label="90% of the context used">
      <circle class="inst-ring-track" cx="10" cy="10" r="8"/>
      <circle class="inst-ring-fill" cx="10" cy="10" r="8" style="--value:0.9"/>
    </svg>
  </span>
  <span class="inst-tooltip-text" role="tooltip" id="tt-ring-xs">Context used: 90%</span>
</span>
```

## Composition

### In a card

```html preview context
<div class="inst-card inst-cluster inst-cluster--loose">
  <span class="inst-ring-wrap">
    <svg class="inst-ring" viewBox="0 0 20 20" role="img" aria-label="74% of the budget spent">
      <circle class="inst-ring-track" cx="10" cy="10" r="8"/>
      <circle class="inst-ring-fill" cx="10" cy="10" r="8" style="--value:0.74"/>
    </svg>
    <span class="inst-ring-label">74</span>
  </span>
  <div>
    <div class="inst-card-title">The budget is spent</div>
    <div class="inst-card-sub">184 320 of 250 000 tokens</div>
  </div>
</div>
```

### With no label inside

```html preview
<span class="inst-ring-wrap">
  <svg class="inst-ring" viewBox="0 0 20 20" role="img" aria-label="40% done">
    <circle class="inst-ring-track" cx="10" cy="10" r="8"/>
    <circle class="inst-ring-fill" cx="10" cy="10" r="8" style="--value:0.4"/>
  </svg>
</span>
```

The wrapper stays: it gives the `inline-flex` and removes the extra line gap
under an SVG. The `aria-label` stays all the more — with no label inside it is
the only carrier of the value for a screen reader.

## API

```api
```

## Related

```related
```
