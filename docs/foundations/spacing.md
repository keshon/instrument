---
title: Rhythm and shape
group: Foundations
layout: foundation
source: src/tokens.css
api:
  - { name: "inst-stack", kind: "class", doc: "A vertical flow, gap `--gap-inline`" }
  - { name: "inst-cluster", kind: "class", doc: "A horizontal flow that wraps" }
  - { name: "inst-cluster-spacer", kind: "class", doc: "A break in a row: everything after it goes to the right" }
  - { name: "inst-grid", kind: "class", doc: "A responsive `auto-fit` grid by `--col-min`" }
  - { name: "inst-stack--tight", kind: "modifier", doc: "A step of the gap" }
  - { name: "inst-stack--loose", kind: "modifier", doc: "A step of the gap" }
  - { name: "inst-cluster--tight", kind: "modifier", doc: "A step of the gap" }
  - { name: "inst-cluster--loose", kind: "modifier", doc: "A step of the gap" }
  - { name: "inst-grid--tight", kind: "modifier", doc: "A different minimum column" }
  - { name: "inst-grid--wide", kind: "modifier", doc: "A different minimum column" }
  - { name: "--space-1", kind: "token", doc: "A step of the scale" }
  - { name: "--space-2", kind: "token", doc: "A step of the scale" }
  - { name: "--space-3", kind: "token", doc: "A step of the scale" }
  - { name: "--space-4", kind: "token", doc: "A step of the scale" }
  - { name: "--space-5", kind: "token", doc: "A step of the scale" }
  - { name: "--space-6", kind: "token", doc: "A step of the scale" }
  - { name: "--space-7", kind: "token", doc: "A step of the scale" }
  - { name: "--space-8", kind: "token", doc: "A step of the scale" }
  - { name: "--space-9", kind: "token", doc: "A step of the scale" }
  - { name: "--space-10", kind: "token", doc: "A step of the scale" }
  - { name: "--pad-panel", kind: "token", doc: "The inner padding of a panel" }
  - { name: "--pad-card", kind: "token", doc: "The inner padding of a card" }
  - { name: "--pad-cell-x", kind: "token", doc: "The horizontal padding of a cell" }
  - { name: "--row-pad-y", kind: "token", doc: "The vertical padding of a row" }
  - { name: "--gap-row", kind: "token", doc: "The gap between rows" }
  - { name: "--gap-inline", kind: "token", doc: "The gap in a row" }
  - { name: "--control-pad-sm", kind: "token", doc: "The padding of a small control" }
  - { name: "--control-pad-md", kind: "token", doc: "The padding of a base control" }
  - { name: "--control-pad-lg", kind: "token", doc: "The padding of a large control" }
  - { name: "--radius-2xs", kind: "token", value: "2px", doc: "The bottom step of the ladder. Marks take it through --radius-mark" }
  - { name: "--radius-mark", kind: "token", doc: "The rounding of a mark: anything the kit draws below 8px — the dot, a tick, the meter and the lane" }
  - { name: "--radius-xs", kind: "token", value: "4px", doc: "Inline: badges, chips, a key" }
  - { name: "--radius-control-sm", kind: "token", doc: "The rounding of a small control. It travels with density along with the height" }
  - { name: "--radius-control-md", kind: "token", doc: "The rounding of a control" }
  - { name: "--radius-control-lg", kind: "token", doc: "The rounding of a large control" }
  - { name: "--radius-sm", kind: "token", doc: "Controls inside a card" }
  - { name: "--radius-md", kind: "token", doc: "Free-standing controls" }
  - { name: "--radius-lg", kind: "token", doc: "Cards and panels" }
  - { name: "--radius-full", kind: "token", doc: "A circle and a pill" }
  - { name: "--hairline", kind: "token", doc: "A hairline. At 1x it equals a pixel, above that a half" }
  - { name: "--stroke", kind: "token", doc: "A line that has to stay visible always" }
  - { name: "--aside-w", kind: "token", doc: "The side column of the shell" }
  - { name: "--container-max", kind: "token", doc: "The ceiling on the width of the working area" }
  - { name: "--split-side", kind: "token", doc: "The wanted width of a split's side part" }
  - { name: "--split-main", kind: "token", doc: "The threshold at which the main part wraps" }
  - { name: "--col-min", kind: "token", doc: "The minimum column of a responsive grid" }
---

A step of 4px, and the scale is deliberately **sparse at the top** — so that "a
little bigger" is simply not in the set. A component takes its padding from the
tier of roles rather than from the scale: that is what gives density one point
through which it retunes everything at once.

```html preview
<div class="inst-card inst-stack inst-stack--tight">
  <div class="inst-card-title">worldgen-01</div>
  <div class="inst-card-sub">Generating biomes, 3 passes</div>
  <div class="inst-cluster">
    <button class="inst-btn inst-btn--sm" type="button">Retry</button>
    <span class="inst-badge" data-tone="ok"><span class="inst-dot"></span>done</span>
  </div>
</div>
```

## Contract

A padding is taken **by purpose**: `--pad-panel`, not `--space-5`. A role is
retuned by density, a step of the scale is not.

```css
.my-panel {
  padding: var(--pad-panel);
  gap: var(--gap-row);
  border-radius: var(--radius-lg);
}
```

| What | Required | Why |
|---|---|---|
| A role rather than a step | yes | `--pad-panel` changes along with density, `--space-5` does not. A component with a hard-coded step is the first to break |
| The gap is set by the container | yes | A margin on an element adds up with the parent's inner padding vertically and does not add up horizontally |
| A radius from the vocabulary | yes | Five steps named by the size of a control. A number of your own parts from its neighbours at the joint |
| A hairline through `--hairline` | yes | At a fractional DPR that is 0.5px, and the line stays a line rather than a band |


### Accessibility

| | |
|---|---|
| Tap targets | The padding of a control takes part in the size of the target on a par with the height. At `compact` small controls go below 24px under WCAG 2.5.8 — see [density](./density.md) |
| Geometry in `px` | Paddings and radii do not grow along with type size: a dense tooling look would otherwise fall apart. The text grows, and the heights work as a minimum |
| Parting groups | Rhythm is a carrier of meaning too: four metrics with no borders read as one group precisely by the gap |
| The hairline | On a 1x display it stays `1px`: a line that has to be visible is visible on any screen |

## Scale

A component reaches it only through the tier of roles. Using `--space-*`
directly is allowed in an application and in the flow primitives.

| Token | Value |
|---|---|
| `--space-1` | 2px |
| `--space-2` | 4px |
| `--space-3` | 6px |
| `--space-4` | 8px |
| `--space-5` | 12px |
| `--space-6` | 16px |
| `--space-7` | 24px |
| `--space-8` | 32px |
| `--space-9` | 48px |
| `--space-10` | 64px |

The sparseness at the top is not thrift but a restriction: there is no step
between 32 and 48, and that is right. A step that does not exist cannot be
chosen by mistake.

### Paddings by purpose

The tier of roles. The values below are for the default density; at `compact`
and `comfortable` they differ, see [density](./density.md).

| Token | Default | Where |
|---|---|---|
| `--pad-panel` | `--space-5` | The inner padding of a panel |
| `--pad-card` | `--space-6` | The inner padding of a card |
| `--pad-cell-x` | `--space-4` | The horizontal padding of a table cell |
| `--row-pad-y` | `--space-3` | The vertical padding of a row: a task, a cell, a menu item |
| `--gap-row` | `--space-3` | The gap between rows in a stack |
| `--gap-inline` | `--space-4` | The gap between neighbours in a row |
| `--control-pad-sm` | `--space-4` | The horizontal padding of a small control |
| `--control-pad-md` | `--space-5` | The same for a medium one |
| `--control-pad-lg` | `--space-6` | The same for a large one |

**The gap is set between neighbours.** A `margin-block` on every element adds
up with the container's inner padding at the top and the bottom but not at the
sides — and a highlight steps away from the border twice as far vertically as
horizontally. The right form is `X + X { margin-block-start }`.

### The sizes of glyphs and gutters

The same tier of roles: every size has a name, so there are no constants in the
components.

| Token | Default | What it is |
|---|---|---|
| `--size-dot` | 6px | The state dot |
| `--size-check` | 15px | The checkbox and the radio |
| `--size-switch-w` / `--size-switch-h` | 30 / 18px | The track of a switch |
| `--size-thumb` | 14px | The thumb |
| `--size-track` | 4px | The track of a slider |
| `--size-meter` | 6px | The thickness of a meter |
| `--size-spinner` | 14px | The ring of busyness |
| `--size-icon` | 18px | An icon — see [icons](./icons.md) |
| `--size-chevron` | 10px | A chevron, a small icon |
| `--size-marker` | 2px | A load-bearing mark: a tab's bar, the edge of a note, the edge of the current navigation item. Thinner and it is not seen, thicker and it argues with the text |
| `--size-gutter` | 14px | The left gutter of a task row |
| `--size-indent` | 14px | The indent step of a tree |
| `--label-col` | 104px | The label column of the inspector |

### Radii

Six steps, and what chooses among them is **nesting**. All the values are even:
at a density of 1.5 an odd radius rasterises off the device grid.

| Token | Value | What for |
|---|---|---|
| `--radius-2xs` | 2px | The bottom step. Marks reach it through `--radius-mark` rather than by name |
| `--radius-xs` | 4px | Inline: badges, chips, a key |
| `--radius-sm` | 6px | Controls **inside** a card or a panel |
| `--radius-md` | 8px | Free-standing controls: a button, a field, a popover |
| `--radius-lg` | 12px | Cards, panels, the modal |
| `--radius-full` | 999px | Tracks, thumbs, the avatar |

The rounding of a CONTROL is taken not from here directly but from
`--radius-control-sm/md/lg`: it travels with density along with the height,
otherwise a small button reads as a pill. See [density](./density.md).

The rounding of a MARK is taken from `--radius-mark` for a reason of the
opposite kind. Below 8px the ladder stops offering a choice: when two radii on
one edge add up to more than the edge is long, the browser scales every radius
on the box until they fit, silently. `--radius-xs` on a 6px dot arrives as 3 —
half the box, which is a circle; on a 6px meter it arrives as 3 as well, which
is a pill. And where the other side of the box is data — the width of a lane
span, the width of a tick — the shape then changes with the data. The mark tier
is the one step that survives every mark in every scale.

### The rule of nested radii

An inner corner needs a radius **smaller** than the outer one, otherwise the
inner element looks as if it bulges out of its container. Optically the inner
radius equals the outer minus the thickness of the gap between them, and the
scale is already chosen so that the neighbouring step down gives the right
result at a padding from the tier of roles.

```html
<!-- panel lg → a control inside sm; a free-standing control is md -->
<div class="inst-panel">
  <div class="inst-panel-body">
    <div class="inst-card"></div>
  </div>
</div>
```

**A rounding on a one-sided border is always a bug.** A `border-left` together
with a `border-radius` gives a stump of an arc: a note and the edge of the
current navigation item carry a square corner deliberately.

### The hairline

A browser never draws a border thinner than one physical pixel. On a 1x display
`0.5px` and `1px` give **the same thing**, while a negative margin of `-0.5px`
against a border drawn at `1px` leaves half a pixel of doubled line. So the
hairline is switched on only where there is room for it.

| Token | Value | When |
|---|---|---|
| `--hairline` | `1px` | A pixel density below 1.5dppx |
| `--hairline` | `0.5px` | `@media (min-resolution: 1.5dppx)` |
| `--stroke` | `1px` | A line that has to stay visible always |

The seam of a button group is made by **removing the neighbour's border** — for
the same reason.

## Rules

| Token | Value | What it is |
|---|---|---|
| `--aside-w` | 232px | The side column of the shell |
| `--container-max` | 76rem | The ceiling on the width of the working area. In `rem`, because it is a measure of readability rather than geometry |
| `--split-side` | 18rem | The wanted width of a split's side part |
| `--split-main` | 30rem | Below this the main part wraps underneath |
| `--col-min` | 260px | The minimum column of a responsive grid |

### What to take and what not to

| Use | Take instead |
|---|---|
| A role: `--pad-panel`, `--gap-row`, `--row-pad-y` | **`--space-*` in a component** — density retunes the tier of roles, not the scale. The whole list of roles is in [tokens](./tokens.md) |
| A gap between neighbours — a flow primitive `.inst-stack` / `.inst-cluster` / `.inst-grid` | **A spacing utility** (`mt-3`, `p-2`) — the library has none and will have none, see [flow](../layout/flow.md) |
| A step of the gap by name: ordinary, `--tight`, `--loose` | **A gap of your own as a number** — "a little bigger" must not be an available decision. The three steps are described in [flow](../layout/flow.md) |
| A radius by nesting: a container `lg`, a control inside `sm` | **One radius for everything** — an inner element with the container's radius bulges out. An example is [the card](../components/display/card.md) |
| `--hairline` for a line that must not shout | **`1px` as a number** — on a 2x display that is twice the wanted thickness, see [the hairline](#the-hairline) |

## API

```api
```

```api
```

## Related

[Density](./density.md)
[Tokens](./tokens.md)
[Flow](../layout/flow.md)
[Shell](../layout/shell.md)
[Panel](../components/display/panel.md)
[Card](../components/display/card.md)
