---
title: Navigation
group: Navigation
layout: component
source: src/layout.css
api:
  - { name: "inst-nav", kind: "class", doc: "The container of the list. A column with a gap of `--space-1`" }
  - { name: "inst-nav-group", kind: "class", doc: "A meaningful block of items" }
  - { name: "inst-nav-label", kind: "class", doc: "The label of a group" }
  - { name: "inst-nav-item", kind: "class", doc: "An item. A link rather than a button" }
  - { name: "inst-nav-count", kind: "class", doc: "Pushes a counter to the far edge of an item" }
  - { name: "--control-h-md", kind: "token" }
  - { name: "--radius-sm", kind: "token" }
  - { name: "--text-sm", kind: "token" }
  - { name: "--text-2xs", kind: "token" }
  - { name: "--weight-medium", kind: "token" }
  - { name: "--space-1", kind: "token" }
  - { name: "--space-3", kind: "token" }
  - { name: "--gap-inline", kind: "token" }
  - { name: "--pad-panel", kind: "token" }
  - { name: "--size-marker", kind: "token" }
  - { name: "--surface-hover", kind: "token" }
  - { name: "--surface-selected", kind: "token" }
  - { name: "--accent-solid", kind: "token" }
  - { name: "--accent-text", kind: "token" }
  - { name: "--text-secondary", kind: "token" }
  - { name: "--text-muted", kind: "token" }
  - { name: "--dur-1", kind: "token" }
  - { name: "--ease-out", kind: "token" }
---

The list of an application's sections in the side column. The current item is
marked by more than colour: a load-bearing bar stands at its edge.

```html preview
<nav class="inst-nav" aria-label="Sections">
  <div class="inst-nav-group">
    <span class="inst-nav-label">Work</span>
    <a class="inst-nav-item" href="#" aria-current="page">
      <svg class="inst-icon" aria-hidden="true"><use href="#i-grid"/></svg>Runs
      <span class="inst-badge inst-nav-count">7</span></a>
    <a class="inst-nav-item" href="#">
      <svg class="inst-icon" aria-hidden="true"><use href="#i-list"/></svg>Queue</a>
    <a class="inst-nav-item" href="#">
      <svg class="inst-icon" aria-hidden="true"><use href="#i-chart"/></svg>Reports</a>
  </div>
  <div class="inst-nav-group">
    <span class="inst-nav-label">Setup</span>
    <a class="inst-nav-item" href="#">
      <svg class="inst-icon" aria-hidden="true"><use href="#i-users"/></svg>Access</a>
    <a class="inst-nav-item" href="#">
      <svg class="inst-icon" aria-hidden="true"><use href="#i-settings"/></svg>Parameters</a>
  </div>
</nav>
```

## Contract

| What | Required | Why |
|---|---|---|
| A `<nav>` + an `aria-label` | yes | There is more than one navigation on a screen (the side one, breadcrumbs, tabs, pages) — with no name they are indistinguishable in the list of landmarks |
| An `<a href>` on an item | yes | Navigation is an address. A button breaks the middle click, "open in a new tab" and copying the link |
| An `aria-current="page"` | yes | The only source of the current state. There is no class for it |
| An `aria-hidden="true"` on the icon | yes | The name of the item is given by the label |
| The current item by a bar at the edge rather than by colour alone | yes | A carrier that does not depend on colour. `aria-current` tells a screen reader the same |

### Accessibility

| | |
|---|---|
| The keyboard | All of it from native links: `Tab` traverses, `Enter` goes. `instrument.js` intercepts nothing, and navigation needs no roving `tabindex` |
| The landmark | A `<nav aria-label>` is a destination for "go to a landmark" in a screen reader |
| Colour is not the only carrier | The current item carries a bar at its edge and an `aria-current` rather than a shade of background alone |
| The tap target | The height of an item is `--control-h-md`, 32px. At density `compact` it shrinks: apply `compact` to tabular views rather than to the main navigation |
| Contrast | The label `--text-secondary` holds 4.5:1 on the muted surface of the column in five themes |
| Print | The side column is not printed: it is the chrome of the application rather than data |

## Anatomy

`inst-nav-group` gathers the items into a meaningful block and
`inst-nav-label` labels it. The gap between groups is set **between neighbours**
(`.inst-nav-group + .inst-nav-group`) rather than around each — otherwise the
first group would step away from the top edge of the column twice as far as
from the sides.

The label is an ordinary sentence at size `--text-2xs` in the colour
`--text-muted`. There are no capitals here: a label has no right to sound
louder than the items it labels.

## States

```html preview
<nav class="inst-nav" aria-label="An example of the current item">
  <div class="inst-nav-group">
    <a class="inst-nav-item" href="#" aria-current="page">Runs</a>
    <a class="inst-nav-item" href="#">Queue</a>
  </div>
</nav>
```

The mark is set by the attribute `aria-current="page"` rather than by a class.
The styling is caught by the selector
`[aria-current]:not([aria-current="false"])`, so an `aria-current="false"` on
the other items is allowed and highlights nothing.

| The mark of the current one | What it gives |
|---|---|
| A bar at the inner edge | A `box-shadow: inset` of width `--size-marker`. A carrier that does not depend on colour |
| The surface | `--surface-selected` |
| The colour of the label | `--accent-text` |
| The weight | `--weight-medium`. There is no weight 700 in the library |

## Behavior

### In a narrow shell

Inside a `.inst-shell-aside` below 60rem the column unfolds into a horizontal
strip: `.inst-nav` and `.inst-nav-group` line up in a row, `.inst-nav-label` is
hidden and the labels of the items stop wrapping. That is the one place where
navigation changes shape, and it is done by [the shell](../../layout/shell.md)
rather than by the navigation itself.

## Composition

### A counter in an item

```html preview
<nav class="inst-nav" aria-label="An example with a counter">
  <div class="inst-nav-group">
    <a class="inst-nav-item" href="#">Queue
      <span class="inst-badge inst-nav-count">7</span></a>
    <a class="inst-nav-item" href="#">History</a>
  </div>
</nav>
```

`inst-nav-count` is a modifier of position: a `margin-inline-start: auto`
pushes [the badge](../display/badge.md) to the far edge of the item. The badge
itself stays a badge and reads the same tokens.

## API

```api
```

## Related

```related
```
