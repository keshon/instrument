---
title: Flow primitives
group: Layout
layout: component
source: src/layout.css
api:
  - { name: "inst-stack", kind: "class", doc: "A column. A gap of `--pad-panel`" }
  - { name: "inst-cluster", kind: "class", doc: "A row with wrapping, aligned to the centre. A gap of `--gap-inline`" }
  - { name: "inst-cluster-spacer", kind: "class", doc: "Push the tail of a row to the far edge" }
  - { name: "inst-grid", kind: "class", doc: "A responsive grid. A gap of `--pad-panel`" }
  - { name: "inst-stack--tight", kind: "modifier", doc: "A step of the gap" }
  - { name: "inst-stack--loose", kind: "modifier", doc: "A step of the gap" }
  - { name: "inst-cluster--tight", kind: "modifier", doc: "A step of the gap" }
  - { name: "inst-cluster--loose", kind: "modifier", doc: "A step of the gap" }
  - { name: "inst-grid--tight", kind: "modifier", doc: "The minimum width of a column" }
  - { name: "inst-grid--wide", kind: "modifier", doc: "The minimum width of a column" }
  - { name: "--col-min", kind: "variable", value: "260px" }
  - { name: "--flow-self", kind: "variable", doc: "The channel of a column: what it tells its direct children with, that they are sized by their content. Non-inheriting, so it does not leak deeper" }
  - { name: "--pad-panel", kind: "token" }
  - { name: "--gap-row", kind: "token" }
  - { name: "--gap-inline", kind: "token" }
  - { name: "--space-2", kind: "token" }
  - { name: "--gap-section", kind: "token" }
---

Three ways to place elements: as a stack, as a row, as a grid. Each has three
steps of the gap, named **by intent rather than by a number**.

```html preview
<div class="inst-stack">
  <div class="inst-cluster">
    <button class="inst-btn inst-btn--sm" type="button">Filters</button>
    <button class="inst-btn inst-btn--sm" type="button">Period</button>
    <span class="inst-cluster-spacer"></span>
    <button class="inst-btn inst-btn--sm inst-btn--primary" type="button">Run</button>
  </div>
  <div class="inst-grid inst-grid--tight">
    <div class="inst-metric"><div class="inst-metric-label">In work</div><div class="inst-metric-value">7</div></div>
    <div class="inst-metric"><div class="inst-metric-label">In the queue</div><div class="inst-metric-value">5</div></div>
    <div class="inst-metric"><div class="inst-metric-label">Fell</div><div class="inst-metric-value">1</div></div>
  </div>
</div>
```

## Contract

| What | Required | Why |
|---|---|---|
| The gap is set by the **container** rather than by the element | yes | An indent on an element adds up with the padding of the container at the top and at the bottom, but not at the sides: the block travels away from the border twice as far vertically as horizontally |
| The step of the gap through a modifier | yes | `--tight` and `--loose` are named by intent. A number in the markup breaks at the very first change of density |
| One `inst-cluster-spacer` per row | yes | A second one will give nothing: the first has already taken the remainder |
| A `<ul>` instead of a `<div>` if this is a list | no, but usually yes | The primitives carry no semantics |

### Accessibility

| | |
|---|---|
| Order | All three primitives keep the order of the markup: the visual one and the keyboard one coincide. `order` and `row-reverse` are not applied in the library |
| Wrapping | A cluster and a grid wrap by themselves — at an increase of the type size to 200% the content is not cut off and gives no horizontal scrolling |
| The spacer | `inst-cluster-spacer` is empty and absent from the accessibility tree: it reports nothing beyond the order |
| Density | All three gaps arrive from roles and are retuned by `data-density` at once. A component with a number wired into it breaks here first |
| Semantics | The primitives are `<div>`s. If a group of elements is a list, put a `<ul>` and the class on it |

## Variants

### The stack

```html preview
<div class="inst-stack inst-stack--loose">
  <div class="inst-stack inst-stack--tight">
    <div class="inst-card"><div class="inst-card-title">The first</div></div>
    <div class="inst-card"><div class="inst-card-title">The second</div></div>
  </div>
  <div class="inst-stack inst-stack--tight">
    <div class="inst-card"><div class="inst-card-title">The third</div></div>
    <div class="inst-card"><div class="inst-card-title">The fourth</div></div>
  </div>
</div>
```

| The class | The gap | When |
|---|---|---|
| `inst-stack` | `--pad-panel` | The default: the blocks of a screen, sections, panels |
| `inst-stack--tight` | `--gap-row` | Rows inside a block: key-value pairs, a list of tasks |
| `inst-stack--loose` | `--gap-section` | The large divisions of meaning of a long screen |

A button, a badge, a tag and a segmented control in a column **do not
stretch**: a button the full width of a card is the same mistake as an icon
button drawn as a rectangle. Fields and cards, on the contrary, do take the
width of the column.

A column does not know who is inside it: it declares an intent to its direct
children, and the one it concerns reads it. A component that is sized by its
content carries `align-self: var(--flow-self, center)` in its own rule — beside
its `display: inline-flex`. So "I am sized by my content" stays a property of
the component rather than a list of other people's names in the layout.

The fallback value `center` explains why nothing changes in a row: the intent
is not declared there, and the component centres as it should.

This works in **every column of the library** — the stack, the form, the
fieldset, the field, the section. In a column of your own application it does
not: a `display: flex; flex-direction: column` written by hand declares no
intent, and the controls stretch all over again. That is an argument for taking
the primitive rather than writing the column — or else for declaring the
channel yourself:

```css
.my-column > * { --flow-self: start; }
```

### The cluster

```html preview
<div class="inst-cluster inst-cluster--loose">
  <span class="inst-cluster inst-cluster--tight">
    <span class="inst-badge" data-tone="ok"><span class="inst-dot"></span>done</span>
    <span class="inst-badge" data-tone="running"><span class="inst-dot"></span>running</span>
  </span>
  <span class="inst-cluster-spacer"></span>
  <button class="inst-btn inst-btn--sm" type="button">More</button>
</div>
```

| The class | The gap | When |
|---|---|---|
| `inst-cluster` | `--gap-inline` | The default: buttons, badges, controls in a row |
| `inst-cluster--tight` | `--space-2` | Elements that read as one group: tags, chips |
| `inst-cluster--loose` | `--pad-panel` | Groups of differing meaning in one line |

A cluster **always wraps** (`flex-wrap: wrap`) and aligns its children to the
centre of the cross axis. `inst-cluster-spacer` is an empty element with a
`margin-inline-start: auto`: everything after it travels to the far edge.

### The grid

```html preview
<div class="inst-grid inst-grid--wide">
  <div class="inst-card"><div class="inst-card-title">A card</div>
    <div class="inst-card-sub">The columns rebuild themselves: auto-fit, with not a single query.</div></div>
  <div class="inst-card"><div class="inst-card-title">A card</div>
    <div class="inst-card-sub">min() keeps a column from spilling out of the container on a narrow screen.</div></div>
</div>
```

| The class | The minimum column | When |
|---|---|---|
| `inst-grid` | `--col-min`, 260px at the base scale | The default: cards, panels |
| `inst-grid--tight` | 180px | Small cells: metrics, tiles of states |
| `inst-grid--wide` | 380px | Large blocks with text inside |

The count of columns is not set: `repeat(auto-fit, minmax(min(var(--col-min), 100%), 1fr))`
works it out itself. The `min(…, 100%)` is required — without it a column 380px
wide does not fit into a container 320px wide and spills outside along with a
horizontal scrollbar.

That is the first level of responsiveness, the **intrinsic** one: it always
works, including where there is no container ancestor.

```css
/* A density of your own for a grid — one line */
.my-board { --col-min: 320px; }
```

## API

```api
```

## Related

```related
```
