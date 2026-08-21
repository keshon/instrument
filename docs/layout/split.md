---
title: Split
group: Layout
layout: component
source: src/layout.css
api:
  - { name: "inst-split", kind: "class", doc: "The container. A row that wraps, with a gap of `--pad-panel`" }
  - { name: "inst-split-main", kind: "class", doc: "The main part. It takes the remainder and goes no narrower than `--split-main`" }
  - { name: "inst-split-side", kind: "class", doc: "The side part. Its wanted width is `--split-side`" }
  - { name: "--split-side", kind: "variable", value: "18rem", doc: "The wanted width of the side part" }
  - { name: "--split-main", kind: "variable", value: "30rem", doc: "Below this the main part wraps underneath" }
  - { name: "--pad-panel", kind: "token" }
---

Two columns of differing importance: the main one and the side one. It wraps by
itself when the main part stops having its threshold — **with not a single
media query**.

```html preview
<div class="inst-split">
  <div class="inst-split-side inst-card">
    <div class="inst-card-title">split-side</div>
    <div class="inst-card-sub">The wanted width is 18rem. Filters, properties, a summary.</div>
  </div>
  <div class="inst-split-main inst-card">
    <div class="inst-card-title">split-main</div>
    <div class="inst-card-sub">It takes the remainder. Narrow the region — the columns stand
      one under another, and not one @media was needed for that.</div>
  </div>
</div>
```

## Contract

| What | Required | Why |
|---|---|---|
| An `inst-split-main` and an `inst-split-side` | yes | Both parts are wanted: a split is about a layout of **differing** importance, and without the second it degenerates into a block |
| The order in the markup = the order on wrapping | yes | `order` parts the visual order from the keyboard one, and `Tab` starts jumping across the screen |
| The thresholds through `--split-side` and `--split-main` | yes | Roles of tier 3, overridden on the container rather than inside the component |
| A media query | no | The wrapping is done by the intrinsic size: it works inside a narrow panel and in a modal alike |

### Accessibility

| | |
|---|---|
| Order | The visual order equals the order of the markup. `order` is not used: it parts `Tab` from the picture |
| Zoom | The wrapping fires on an increase of type size too: the thresholds are set in `rem`, so the columns fold along with the growth of the text rather than after it |
| Landmarks | The parts are `<div>`s. If the side part is complementary content, put an `<aside>` |
| Scrolling | A split scrolls nothing itself. Independent scrolling is the work of [the shell](./shell.md) and [the panel](../components/display/panel.md) |

## Anatomy

The wrapping is achieved not by a query but by two numbers:

| The part | The rule | Which means |
|---|---|---|
| `inst-split-side` | `flex: 1 1 var(--split-side)` | A wanted width of 18rem. Ready to squeeze and to stretch |
| `inst-split-main` | `flex: 999 1 0` + `min-inline-size: min(var(--split-main), 100%)` | It takes almost the whole remainder but does not squeeze below 30rem |

When the main part stops having its minimum, `flex-wrap` moves the side one
below, and both take the full width. That is the first level of responsiveness
— the **intrinsic** one: it works where there is no container ancestor, inside
a narrow panel and in a modal alike.

The `min()` in the minimum is required: without it a main part of 30rem does
not fit into a region of 24rem and sticks out along with a horizontal
scrollbar.

### The order

The order in the markup is the order on wrapping. If you want the main part on
top on a narrow screen, put it first:

```html
<div class="inst-split">
  <div class="inst-split-main">…</div>
  <div class="inst-split-side">…</div>
</div>
```

`order` is not used for that: it parts the visual order from the keyboard one,
and `Tab` starts jumping across the screen.

## Composition

### A list and an inspector

```html
<div class="inst-split">
  <div class="inst-split-main inst-panel">
    <div class="inst-panel-header"><span class="inst-panel-title">Queue</span></div>
    <div class="inst-panel-body inst-panel-body--list">…</div>
  </div>
  <div class="inst-split-side inst-panel">
    <div class="inst-panel-header"><span class="inst-panel-title">Properties</span></div>
    <div class="inst-panel-body">…</div>
  </div>
</div>
```

### With a stack inside

```html
<div class="inst-split">
  <div class="inst-split-main inst-stack">…</div>
  <div class="inst-split-side inst-stack inst-stack--tight">…</div>
</div>
```

## Customization

```css
/* A narrow column of properties, an early wrap */
.inspector-split {
  --split-side: 14rem;
  --split-main: 24rem;
}
```

Both variables are roles of tier 3, so they are overridden on the container
rather than inside the component.

## API

```api
```

## Related

```related
```
