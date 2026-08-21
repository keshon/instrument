---
title: Wizard steps
group: Navigation
layout: component
source: src/layout.css
api:
  - { name: "inst-stepper", kind: "class", doc: "The list of steps. Put on an `<ol>`" }
  - { name: "inst-stepper-item", kind: "class", doc: "A step. The carrier of `data-state`" }
  - { name: "inst-stepper-name", kind: "class", doc: "The name of the step" }
  - { name: "inst-stepper-note", kind: "class", doc: "The word of the state under the name" }
  - { name: "data-state", kind: "attribute", doc: "`todo` `current` `done`. Closed, with `todo` as the base" }
  - { name: "--gap-inline", kind: "token" }
  - { name: "--space-2", kind: "token" }
  - { name: "--space-4", kind: "token" }
  - { name: "--text-sm", kind: "token" }
  - { name: "--text-xs", kind: "token" }
  - { name: "--weight-medium", kind: "token" }
  - { name: "--size-marker", kind: "token" }
  - { name: "--radius-full", kind: "token" }
  - { name: "--track", kind: "token" }
  - { name: "--accent-mark", kind: "token" }
  - { name: "--text-primary", kind: "token" }
  - { name: "--text-secondary", kind: "token" }
  - { name: "--text-muted", kind: "token" }
---

A sequence of the stages of a process with a mark of where it is now. The state
is carried by a **word** under the label — the colour of the bar is not enough.

```html preview
<ol class="inst-stepper">
  <li class="inst-stepper-item" data-state="done">
    <span class="inst-stepper-name">Source</span><span class="inst-stepper-note">done</span></li>
  <li class="inst-stepper-item" data-state="done">
    <span class="inst-stepper-name">Rules</span><span class="inst-stepper-note">done</span></li>
  <li class="inst-stepper-item" data-state="current">
    <span class="inst-stepper-name">Check</span><span class="inst-stepper-note">now</span></li>
  <li class="inst-stepper-item" data-state="todo">
    <span class="inst-stepper-name">Launch</span><span class="inst-stepper-note">next</span></li>
</ol>
```

## Contract

| What | Required | Why |
|---|---|---|
| An `<ol>` with `<li>` | yes | The steps are ordered, and the order is the whole meaning of the component |
| A `data-state` on every step, `todo` included | yes | The absence of the attribute is indistinguishable from a typo in it |
| An `inst-stepper-note` with the word of the state | yes | The colour of the bar has no right to be the only carrier |

### Accessibility

| | |
|---|---|
| The order | An `<ol>` reports the number of a step and their total without labels of the form "3 of 5" |
| Colour is not the only carrier | The state travels as a word (`inst-stepper-note`) **and** as a bar |
| Contrast | The name is `--text-secondary` (4.5:1), the word of the state `--text-muted`; this is data rather than decoration, and `--text-faint` is not due to it |
| The bar | `--accent-mark` holds 3:1 against the `--track` — it is a mark rather than text |
| Type size | The word of the state is `--text-xs`; the size does not go below 11px |
| Print | The steps are printed: the word of the state reads on paper, where a bar might merge into it |

## Anatomy

### The label of the state

`inst-stepper-note` is neither decoration nor a subheading. It is the second
carrier of the state: "done", "now", "next" read on a printout and to somebody
who does not tell the shade of the bar apart. A line with a date or a counter
in that slot reports no state and does not replace it.

### The layout

The strip wraps by itself: an item has `flex: 1 1 10rem`, and on a narrow
screen the steps stand in several rows with not a single media query. The bar
of the state is a `::before` the full width of the item, so it does not part
from the label at any number of steps.

## States

`data-state` is a phase of life, and the vocabulary is the component's **own
and closed**:

```html preview
<ol class="inst-stepper">
  <li class="inst-stepper-item" data-state="done">
    <span class="inst-stepper-name">Passed</span><span class="inst-stepper-note">done</span></li>
  <li class="inst-stepper-item" data-state="current">
    <span class="inst-stepper-name">Current</span><span class="inst-stepper-note">now</span></li>
  <li class="inst-stepper-item" data-state="todo">
    <span class="inst-stepper-name">Upcoming</span><span class="inst-stepper-note">next</span></li>
</ol>
```

| Value | Means | What it draws |
|---|---|---|
| `todo` | Not begun yet. **The base** | The bar `--track`, the label `--text-secondary` |
| `current` | Happening now | The bar `--accent-mark`, the label `--text-primary` and `--weight-medium` |
| `done` | Finished | The bar `--accent-mark`, the label `--text-secondary` |

The base value `todo` **is written in the markup**, though there are no rules
for it: a `data-state="todo"` is read and its absence is not, and a typo would
look like the base.

A passed step and the current one share one colour of bar deliberately: what
tells them apart is the `inst-stepper-note` — "done" against "now" — and the
weight of the label. There is no second accent tone in the library.

## API

```api
```

## Related

```related
```
