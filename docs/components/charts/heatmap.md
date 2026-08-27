---
title: Heatmap
group: Charts
layout: component
source: src/data.css
api:
  - { name: "inst-heatmap", kind: "class", doc: "The field. Seven rows, weeks running along the inline axis, freshest at the right" }
  - { name: "inst-heatmap-cell", kind: "class", doc: "One day. A square mark, and nothing else goes inside it" }
  - { name: "data-level", kind: "attribute", value: "1 · 2 · 3 · 4", doc: "The intensity of the day. With no attribute the cell is the empty ground: the day happened and nothing was counted on it" }
  - { name: "data-empty", kind: "attribute", doc: "There is no such day — the range starts mid-week or ends mid-week. The room stays, the mark does not" }
  - { name: "role", kind: "attribute", value: "img", doc: "**Required.** The field is a picture of data and has to have a name" }
  - { name: "aria-label", kind: "attribute", doc: "**Required**, and it has to carry the quantity **in words**. Colour is the only carrier in the picture itself" }
  - { name: "--cell-size", kind: "variable", value: "--size-icon-sm", doc: "The side of a cell. Set it on the field to make the squares larger or smaller" }
  - { name: "--size-icon-sm", kind: "token", value: "14px", doc: "The default side of a cell. The reference's measures 14.7" }
  - { name: "--accent-solid", kind: "token", doc: "The top of the ramp. The steps below it are mixed down into the empty ground" }
  - { name: "--surface-recessed", kind: "token", doc: "The empty cell, and the ground every step is mixed into" }
  - { name: "--radius-mark", kind: "token" }
  - { name: "--space-1", kind: "token" }
---

How much of something happened per day, across weeks. A contribution field, a
load calendar, a run count. The one reading here whose axis is **intensity**
rather than a category or an outcome.

```html preview
<div class="inst-heatmap" role="img"
     aria-label="Thirteen weeks of runs: busiest in the last fortnight, at 40 a day; nine quiet days in all">
  <span class="inst-heatmap-cell" data-empty></span>
  <span class="inst-heatmap-cell" data-empty></span>
  <span class="inst-heatmap-cell"></span>
  <span class="inst-heatmap-cell" data-level="1"></span>
  <span class="inst-heatmap-cell" data-level="2"></span>
  <span class="inst-heatmap-cell" data-level="1"></span>
  <span class="inst-heatmap-cell" data-level="3"></span>
  <span class="inst-heatmap-cell"></span>
  <span class="inst-heatmap-cell" data-level="2"></span>
  <span class="inst-heatmap-cell" data-level="4"></span>
  <span class="inst-heatmap-cell" data-level="1"></span>
  <span class="inst-heatmap-cell" data-level="3"></span>
  <span class="inst-heatmap-cell" data-level="2"></span>
  <span class="inst-heatmap-cell"></span>
  <span class="inst-heatmap-cell" data-level="1"></span>
  <span class="inst-heatmap-cell" data-level="4"></span>
  <span class="inst-heatmap-cell" data-level="3"></span>
  <span class="inst-heatmap-cell" data-level="2"></span>
  <span class="inst-heatmap-cell" data-level="1"></span>
  <span class="inst-heatmap-cell"></span>
  <span class="inst-heatmap-cell" data-level="3"></span>
  <span class="inst-heatmap-cell" data-level="4"></span>
  <span class="inst-heatmap-cell" data-level="2"></span>
  <span class="inst-heatmap-cell" data-level="1"></span>
  <span class="inst-heatmap-cell" data-level="3"></span>
  <span class="inst-heatmap-cell"></span>
  <span class="inst-heatmap-cell" data-level="2"></span>
  <span class="inst-heatmap-cell" data-level="4"></span>
  <span class="inst-heatmap-cell" data-level="3"></span>
  <span class="inst-heatmap-cell" data-level="1"></span>
  <span class="inst-heatmap-cell" data-level="2"></span>
  <span class="inst-heatmap-cell" data-level="4"></span>
  <span class="inst-heatmap-cell" data-level="4"></span>
  <span class="inst-heatmap-cell" data-level="3"></span>
  <span class="inst-heatmap-cell" data-level="2"></span>
  <span class="inst-heatmap-cell" data-level="1"></span>
  <span class="inst-heatmap-cell" data-level="4"></span>
  <span class="inst-heatmap-cell" data-level="3"></span>
  <span class="inst-heatmap-cell" data-level="2"></span>
  <span class="inst-heatmap-cell" data-level="4"></span>
  <span class="inst-heatmap-cell" data-level="1"></span>
  <span class="inst-heatmap-cell" data-level="3"></span>
  <span class="inst-heatmap-cell" data-level="4"></span>
  <span class="inst-heatmap-cell" data-level="2"></span>
  <span class="inst-heatmap-cell" data-level="3"></span>
  <span class="inst-heatmap-cell" data-level="4"></span>
  <span class="inst-heatmap-cell" data-level="1"></span>
  <span class="inst-heatmap-cell" data-empty></span>
  <span class="inst-heatmap-cell" data-empty></span>
</div>
```

## Contract

| What | Required | Why |
|---|---|---|
| A `role="img"` and an `aria-label` with the quantity **in words** | yes | Ninety empty `<span>`s are ninety empty `<span>`s to a screen reader. Here the name is the whole of the accessible content: unlike a badge or a change mark, a cell has no sign inside it to fall back on |
| Seven rows | yes | The grid is declared with seven rows and the cells flow down them by column. A different number is a different picture and wants a different component |
| The freshest at the **right** | yes | The field is read as time, and time in this script runs to the right — the same rule as [the check history](../../agent/history.md) |
| A level from the closed set | yes | `1` · `2` · `3` · `4`, with no attribute for the empty ground. A fifth value silently does nothing |
| A `data-empty` for days outside the range | no | Only when the range starts or ends mid-week |

### Accessibility

| | |
|---|---|
| Colour is the only carrier, and here that is honest | Intensity has no shape to fall back on: a cell twice as dark is not a cell of a different form. Everything else in the kit that colours a mark also gives it a sign or a word; this cannot, so the burden moves entirely to the accessible name |
| The name | It has to state the quantity, not the shape. "Thirteen weeks of runs, busiest in the last fortnight at 40 a day" is a name; "activity heatmap" is not |
| Forced colours | There is no rule, deliberately. One fill is left in that mode, and a grid of identical squares would go on claiming to mean something it no longer does — the same decision as on [the check history](../../agent/history.md) |
| The role | `img`, not `table`. There are no headers and nothing is read cell by cell; it is one picture with one name |
| Contrast | The steps are marks rather than text and are not held to 4.5:1. What matters is that neighbouring steps are told apart, which the mix guarantees by construction: each is a quarter of the distance from the ground to the accent |

## Variants

The four steps are **mixed**, not declared:

```css
background: color-mix(in oklab, var(--accent-solid) 25%, var(--surface-recessed));
```

`--chart-1..6` are the wrong palette for this and could not be used. They are
categorical — chosen so that series are told **apart** — and a ramp needs the
opposite: steps that are plainly the same colour at different strengths.

Declaring five new tokens instead would mean five values across fifteen scale
and density cells and four accents, re-tuned every time the accent moved. A
mix from `--accent-solid` down into the empty ground follows all four accents
and all five themes for free, and the empty cell takes `--surface-recessed` —
the same fill [a metric tile](../display/metric.md) sits on.

| Level | Fill |
|---|---|
| no attribute | `--surface-recessed` |
| `data-level="1"` | 25% accent |
| `data-level="2"` | 50% accent |
| `data-level="3"` | 75% accent |
| `data-level="4"` | `--accent-solid` |

**Four steps and no more.** Five was tried on screen and the middle two stop
being told apart at a cell of 14px; the eye reads a heatmap by area, and area
needs coarse steps.

## Sizes

The cell defaults to `--size-icon-sm` — 14px, against the reference's measured
14.7 — and rides the type ladder, so the field grows with the interface rather
than staying a fixed square while everything round it moves.

```html
<div class="inst-heatmap" style="--cell-size: 10px" role="img" aria-label="…">
```

A year is fifty-two columns and will not fit a narrow panel. The field pans
with the scrollbar hidden — the device [the tab strip](../navigation/tabs.md)
and [the status bar](../../layout/statusbar.md) both use. It does not wrap:
a wrapped heatmap is two pictures.

## API

```api
```

## Related

```related
```
