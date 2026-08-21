---
title: Button group
group: Actions
layout: component
source: src/actions.css
api:
  - { name: "inst-btn-group", kind: "class", doc: "The container. Only `inst-btn` inside it" }
  - { name: "--radius-md", kind: "token" }
  - { name: "--hairline", kind: "token" }
---

Several different actions standing flush. A group says "these actions are about
the same thing", while each of them stays on its own.

```html preview
<div class="inst-btn-group">
  <button class="inst-btn inst-btn--sm" type="button">Left</button>
  <button class="inst-btn inst-btn--sm" type="button">Centre</button>
  <button class="inst-btn inst-btn--sm" type="button">Right</button>
</div>
```

## Contract

A group adds no roles: every button declares itself — that is the difference
from [the segmented control](./segmented.md), where the role is obligatory.

| What | Required | Why |
|---|---|---|
| `.inst-btn-group` on the container | yes | It removes the neighbour's border and rounds only the ends; without it the buttons stand side by side rather than in a group |
| Ordinary `.inst-btn` inside | yes | A group is a layout rather than a new component: a button inside it stays a button |
| One size for the whole group | yes | Neighbouring buttons of different heights give a step at the seam, and it shows on any background |
| A name of its own on every button | yes | A group needs no `aria-label`: it is not one control but several actions side by side |
| A roving `tabindex` | **no, that is a mistake** | `Tab` goes through every button separately — these are different actions rather than one choice. A roving `tabindex` would hide part of the actions from the keyboard |

The button under focus rises by `z-index`, otherwise its neighbour would cut
the ring off: flex items obey `z-index` without `position`.

## Anatomy

The seam of a group is **the neighbour's border removed**. A negative margin of
`-0.5px` will not do: a browser does not draw a border thinner than a physical
pixel, and the shift left a thickening at every joint.

| What the container does | How |
|---|---|
| Removes the seam | On every button but the first the near border is taken away |
| Rounds the ends | The radius stays only on the first and the last |
| Raises the focus | The button under focus goes higher by `z-index`, otherwise its neighbour would cut the ring off |

## Sizes

The size is set on each button rather than on the container: a group is a
composition rather than a control, and it has no size of its own.

```html preview
<div class="inst-btn-group">
  <button class="inst-btn inst-btn--sm" type="button">One</button>
  <button class="inst-btn inst-btn--sm" type="button">Two</button>
</div>
<div class="inst-btn-group">
  <button class="inst-btn" type="button">One</button>
  <button class="inst-btn" type="button">Two</button>
</div>
<div class="inst-btn-group">
  <button class="inst-btn inst-btn--lg" type="button">One</button>
  <button class="inst-btn inst-btn--lg" type="button">Two</button>
</div>
```

## Composition

### With icon buttons

```html preview
<div class="inst-btn-group">
  <button class="inst-btn inst-btn--sm inst-btn--icon" type="button" aria-label="Refresh">
    <svg class="inst-icon" aria-hidden="true"><use href="#i-refresh"/></svg>
  </button>
  <button class="inst-btn inst-btn--sm inst-btn--icon" type="button" aria-label="Copy">
    <svg class="inst-icon" aria-hidden="true"><use href="#i-copy"/></svg>
  </button>
  <button class="inst-btn inst-btn--sm inst-btn--icon" type="button" aria-label="Settings">
    <svg class="inst-icon" aria-hidden="true"><use href="#i-settings"/></svg>
  </button>
</div>
```

Every icon button carries an `aria-label` of its own: a group has no shared
name, because it is not one control.

## Patterns

### In a toolbar

Groups are parted by a gap, and inside a group things stand flush. The gap is
the message about where one kinship ends and another begins.

```html preview context
<div class="inst-toolbar">
  <div class="inst-btn-group">
    <button class="inst-btn inst-btn--sm" type="button">List</button>
    <button class="inst-btn inst-btn--sm" type="button">Grid</button>
  </div>
  <span class="inst-toolbar-sep"></span>
  <div class="inst-btn-group">
    <button class="inst-btn inst-btn--sm" type="button">Copy</button>
    <button class="inst-btn inst-btn--sm" type="button">Duplicate</button>
  </div>
</div>
```

## API

```api
```

## Related

```related
```
