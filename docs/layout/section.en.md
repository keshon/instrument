---
title: Section
group: Layout
layout: component
source: src/layout.css
api:
  - { name: "inst-section", kind: "class", doc: "A column with a gap of `--gap-row`" }
  - { name: "inst-section-head", kind: "class", doc: "The row of the heading. Aligned to the baseline" }
  - { name: "inst-section-title", kind: "class", doc: "The label of a section" }
  - { name: "inst-section-actions", kind: "class", doc: "The tail of the heading row, pushed to the far edge" }
  - { name: "--gap-row", kind: "token" }
  - { name: "--gap-inline", kind: "token" }
  - { name: "--text-sm", kind: "token" }
  - { name: "--weight-medium", kind: "token" }
  - { name: "--text-secondary", kind: "token" }
---

A named block inside a screen: a label, optional actions and contents. The
label is **quieter than the data** it labels.

```html preview
<div class="inst-section">
  <div class="inst-section-head">
    <h2 class="inst-section-title">The summary for the day</h2>
    <span class="inst-section-actions">
      <button class="inst-btn inst-btn--sm inst-btn--ghost" type="button">Period</button></span>
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
| A real heading of the wanted level | yes | `inst-section-title` styles but does not declare: navigation by headings does not work without it |
| A `<section>` or a `<div>` | no | The tag changes no role: a section draws nothing except rhythm |
| An `inst-section-head` | no | Wanted where there is a heading or there are actions |
| An `inst-section-actions` — for this block alone | yes | Actions at the level of a screen live in [the page header](./page-header.md) |

A section draws nothing except rhythm: it is a stack with a gap of `--gap-row`
and a row of heading. A border around four metrics would turn one group into
four objects.

| What | Why |
|---|---|
| The tag of the heading by level (`<h2>`, `<h3>`) | The class gives the size rather than the level. A screen reader builds the table of contents from the tags |
| An `inst-section-head` as the wrapper | Without it the heading and the actions become children of the column and line up one under another |
| A `type="button"` on the actions | Otherwise inside a form they will submit it |

### Accessibility

| | |
|---|---|
| The level of the heading | Chosen by the structure of the document rather than by the look. It is not assigned automatically |
| The landmark | If the block really is a self-standing division, put a `<section>` with an `aria-labelledby` on the heading; a `<div>` gives no landmark, and that is fine for grouping |
| Contrast | The label is `--text-secondary`, 4.5:1 in five themes. It is quieter than the data by colour, but not below the floor of reading |
| Type size | `--text-sm`, the same as the base of the document: the label is not larger than the contents |
| Order | The actions stand after the heading in the markup and are traversed after it |

## Anatomy

```html preview
<div class="inst-section">
  <div class="inst-section-head">
    <h2 class="inst-section-title">Active runs</h2>
    <span class="inst-badge inst-nav-count">7</span>
    <span class="inst-section-actions">
      <button class="inst-btn inst-btn--sm inst-btn--ghost" type="button">All</button></span>
  </div>
  <div class="inst-stack inst-stack--tight">
    <div class="inst-card"><div class="inst-card-title">worldgen · #4127</div></div>
    <div class="inst-card"><div class="inst-card-title">worldgen · #4126</div></div>
  </div>
</div>
```

`inst-section-head` aligns its children **to the baseline**: the heading, the
counter and the button stand on one line of text rather than by the centres of
their boxes. So anything at all can be put into the heading row without fitting
heights to one another.

`inst-section-actions` is pushed to the far edge by `margin-inline-start: auto`
— by the same device as `inst-cluster-spacer` and `inst-page-actions`. One axis
— one mechanism.

### The heading

`inst-section-title` is `--text-sm`, `--weight-medium`, `--text-secondary`.
That is, **quieter than the contents**: the label serves the data. There are no
caps and no tracking here for the same reason as everywhere in the library, and
a weight of 700 does not exist in the library at all.

The class sets the size but not the level. The level is chosen by the place in
the structure of the document: under the `<h1>` of the page header it is
usually an `<h2>`.

## Composition

### Several sections on a screen

```html
<div class="inst-container inst-stack inst-stack--loose">
  <div class="inst-page-header">…</div>
  <div class="inst-section">…</div>
  <div class="inst-section">…</div>
</div>
```

The gap between sections is set by the outer stack rather than by the sections
themselves: a gap is between neighbours rather than around each one.

### A section with no actions

```html
<div class="inst-section">
  <div class="inst-section-head">
    <h2 class="inst-section-title">The parameters of a run</h2>
  </div>
  <div class="inst-kv">…</div>
</div>
```

## API

```api
```

## Related

```related
```
