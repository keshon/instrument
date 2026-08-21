---
title: Container
group: Layout
layout: component
source: src/layout.css
api:
  - { name: "inst-container", kind: "class", doc: "A width of 100%, a ceiling of `--container-max`, centring, side fields of `--pad-panel`" }
  - { name: "inst-container--narrow", kind: "modifier", doc: "A ceiling of 44rem" }
  - { name: "inst-container--prose", kind: "modifier", doc: "A ceiling of 68ch" }
  - { name: "inst-container--full", kind: "modifier", doc: "The ceiling is taken off" }
  - { name: "--container-max", kind: "variable", value: "76rem" }
  - { name: "--pad-panel", kind: "token" }
---

The ceiling of the width and the horizontal air of a work area. What it bounds
is not "the design" but **the measure of readability**: a line longer than ~90
characters gets lost on the return to the start.

```html preview
<div class="inst-container inst-container--narrow inst-stack">
  <div class="inst-card">
    <div class="inst-card-title">A narrow container</div>
    <div class="inst-card-sub">44rem. A form, a dialog of settings, the screen of one object.</div>
  </div>
  <div class="inst-card">
    <div class="inst-card-title">The width arrives from outside</div>
    <div class="inst-card-sub">The container is declared a @container region, so the
      components inside answer to its width rather than to the width of the window.</div>
  </div>
</div>
```

The region is named `content`, and that is its public name: a rule of the
application addresses it by that.

```css
/* A rule of your own under a narrow container */
@container content (inline-size < 40rem) {
  .my-thing { grid-template-columns: 1fr; }
}
```

The name may also be left out — a nameless `@container (inline-size < 40rem)`
will find the nearest region-ancestor. The name is wanted where several regions
are nested: a container inside a work area inside a shell is three nested
regions, and without the name the rule will catch the wrong one.

## Contract

| What | Required | Why |
|---|---|---|
| The width arrives from outside | yes | `container-type` switches off the intrinsic along the inline axis: an element that was getting its width from its contents will collapse |
| The ceiling in `rem` rather than in pixels | yes | At an increased system type size the container grows along with the text, and the measure of reading stays the same |
| A ceiling of your own — through `--container-max` | yes | A modifier overrides the variable rather than rewriting the rule |

```css
/* A ceiling of your own — one line */
.my-screen { --container-max: 60rem; }
```

### Accessibility

| | |
|---|---|
| The measure of reading | A ceiling of the width is a requirement of readability rather than of taste: a line longer than ~90 characters gets lost on the return to the start |
| Zoom | The ceiling is in `rem`, the fields are in tokens: at an increase of the type size to 200% the content is not cut off and demands no horizontal scrolling |
| The side fields | `--pad-panel` is retuned by density, so the text does not stick to the edge in any of the three |
| Landmarks | A container is a `<div>` and carries no role. The landmark is given by the `<main>` outside |

## Variants

```html preview
<div class="inst-stack">
  <div class="inst-container"><div class="inst-card"><div class="inst-card-title">by default — 76rem</div></div></div>
  <div class="inst-container inst-container--narrow"><div class="inst-card"><div class="inst-card-title">--narrow — 44rem</div></div></div>
  <div class="inst-container inst-container--prose"><div class="inst-card"><div class="inst-card-title">--prose — 68ch</div></div></div>
</div>
```

| The class | The ceiling | When |
|---|---|---|
| with no modifier | `--container-max`, 76rem | The default: a dashboard, a table, a list of tasks |
| `inst-container--narrow` | 44rem | A form, settings, the screen of one object |
| `inst-container--prose` | 68ch | Readable text. The measure is in `ch` — it is about characters rather than about pixels |
| `inst-container--full` | `none` | The whole width of the work area: maps, timelines, wide tables |

A modifier is **one line**: it overrides `--container-max` rather than
rewriting the rule.

## Behavior

### The container of queries

`.inst-container` is declared `container: content / inline-size`. The components
inside answer to **its** width — the same card behaves identically in the narrow
column of an inspector and across the whole width of a dashboard.

Hence a rule that is easier to break than to notice: only that whose width
**arrives from outside** gets declared a container. `container-type` switches
off the intrinsic along the inline axis, and an element that was getting its
width from its contents will collapse. There are three such regions in the
library: `.inst-container`, `.inst-shell-main` and `.inst-panel`.

## Composition

### Inside the shell

```html
<main class="inst-shell-main" id="main">
  <div class="inst-container inst-stack">
    <div class="inst-page-header">…</div>
    <div class="inst-section">…</div>
  </div>
</main>
```

The vertical air is given by the work area, the horizontal one by the
container. The division is not cosmetic: the side fields are obliged to
disappear along with the bounding of the width, and the top indent is not.

### A container and a stack on one element

```html
<div class="inst-container inst-container--narrow inst-stack inst-stack--loose">…</div>
```

The classes are independent: one answers for the width, the other for the gap
between the children. No extra wrapper is wanted for that.

## API

```api
```

## Related

```related
```
