---
title: Card
group: Data display
layout: component
source: src/surfaces.css
api:
  - { name: "inst-card", kind: "class", doc: "The base one" }
  - { name: "inst-card-title", kind: "class", doc: "The heading of an object" }
  - { name: "inst-card-sub", kind: "class", doc: "The caption under the heading" }
  - { name: "inst-card--interactive", kind: "modifier", doc: "A response to hover and focus" }
  - { name: "--surface-raised", kind: "token" }
  - { name: "--border", kind: "token" }
  - { name: "--hairline", kind: "token" }
  - { name: "--radius-lg", kind: "token" }
  - { name: "--pad-card", kind: "token" }
---

An object one could have dragged. There are as many cards as came from the
data — unlike [a panel](./panel.md), which is a place on the screen.

```html preview
<div class="inst-card">
  <div class="inst-card-title">worldgen-01</div>
  <div class="inst-card-sub">terrain/heightmap.ts · 4 files</div>
</div>
```

A card is a column with a gap between its children. The header does not press
against the content, and the heading does not read as part of the first row of
data.

## Contract

| What | Required | Why |
|---|---|---|
| `inst-card-title` | no | The heading of an object. A card without one is just a frame |
| An `<a>` or a `<button>` for a clickable one | yes | The modifier gives the look, the tag gives the role. A `div` with a handler is unreachable from the keyboard |
| A real heading (`<h3>`) in a list of cards | no, but usually yes | `inst-card-title` styles but does not declare: navigation by headings does not work without one |

### Accessibility

| | |
|---|---|
| Clickability | Through an `<a>` or a `<button>` only. The modifier gives the look, the tag gives the role |
| One target rather than two | If the whole card is clickable, there must be no second link inside it: the keyboard gets two stops for one object |
| The heading | `inst-card-title` is visual. For a list of cards a real heading (`<h3>`) is in order, so that navigation by headings works |
| Focus | A ring with a 1px gap along the outer edge. The radius is inherited from the card |
| No shadow | Depth is carried by the surface. A shadow in the library means "temporary and above" and belongs to the popover and the modal alone |

## Variants

### Interactive

```html preview
<a class="inst-card inst-card--interactive" href="#">
  <div class="inst-card-title">Open the run</div>
  <div class="inst-card-sub">The whole card is the tap target</div>
</a>
```

The modifier adds a response to hover and focus but **does not make the element
clickable**: that has to be an `<a>` or a `<button>`. A card that is a `div`
with a click handler is unreachable from the keyboard and has no role.

## API

```api
```

## Related

```related
```
