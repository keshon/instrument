---
title: Page header
group: Layout
layout: component
source: src/layout.css
api:
  - { name: "inst-page-header", kind: "class", doc: "The container. A row that wraps, aligned to the top" }
  - { name: "inst-page-header-main", kind: "class", doc: "The \"name + explanation\" column. The wrapping threshold is 20rem" }
  - { name: "inst-page-title", kind: "class", doc: "The name of the screen. The size is `--text-xl`" }
  - { name: "inst-page-desc", kind: "class", doc: "The explanation. `--text-sm`, `--text-secondary`, a measure of 68ch" }
  - { name: "inst-page-actions", kind: "class", doc: "The actions. Pushed to the far edge, and they wrap" }
  - { name: "--gap-inline", kind: "token" }
  - { name: "--pad-panel", kind: "token" }
  - { name: "--space-2", kind: "token" }
  - { name: "--space-3", kind: "token" }
  - { name: "--text-xl", kind: "token" }
  - { name: "--text-sm", kind: "token" }
  - { name: "--text-secondary", kind: "token" }
---

The name of a screen, its explanation and the principal actions. One per screen
— exactly as one `<h1>`.

```html preview
<div class="inst-page-header">
  <div class="inst-page-header-main">
    <h1 class="inst-page-title">Runs</h1>
    <p class="inst-page-desc">Every region scrolls by itself: the side navigation does
      not travel away while the table is being read.</p>
  </div>
  <div class="inst-page-actions">
    <button class="inst-btn inst-btn--sm" type="button">Filters</button>
    <button class="inst-btn inst-btn--sm inst-btn--primary" type="button">Run</button>
  </div>
</div>
```

## Contract

The `inst-page-header-main` wrapper is required: it holds the name and the
explanation as one column with a gap of `--space-2` and carries a
`flex: 1 1 20rem` — the threshold below which the actions wrap under the
heading rather than squeezing it into a column of one word.

The `min-inline-size: 0` on it is not decoration: without it a long name with
no spaces stretches the whole header and drives the actions past the edge.

| What | Why |
|---|---|
| An `<h1>` under `inst-page-title` | The class gives the size rather than the level. The level is the work of the tag, and a screen reader reads exactly that |
| An `inst-page-header-main` as the wrapper | Without it the name and the explanation become separate children of the row and line up side by side |
| A `type="button"` on the actions | Otherwise inside a form they will submit it |

### Accessibility

| | |
|---|---|
| The level of the heading | One `<h1>` per document. It is not assigned automatically: it is set by the author of the markup |
| The measure of the explanation | 68ch. A longer line is lost on the return to the start |
| Wrapping | The header wraps by itself; at a 200% increase of type size the actions go under the heading rather than past the edge of the screen |
| Contrast | The explanation is `--text-secondary` (4.5:1). `--text-faint` is not applied to text that gets read |
| Print | The name and the explanation are printed, `inst-page-actions` is not: actions on paper are not actions |
| Order | The actions stand after the heading in the markup and are traversed after it. `order` is not used |

## Anatomy

`inst-page-actions` is pushed to the far edge (`margin-inline-start: auto`) and
wraps within itself. Keep **one** `primary` here: a screen where everything is
primary has no primary at all — see [the button](../components/actions/button.md).

The fourth action and the ones after it move into
[a menu](../components/overlays/menu.md). A strip of six buttons above a screen
is read neither by eye nor from the keyboard.

## Composition

```html preview
<div class="inst-page-header">
  <div class="inst-page-header-main">
    <h1 class="inst-page-title">Run #4127
      <span class="inst-badge" data-tone="running"><span class="inst-dot"></span>running</span></h1>
    <p class="inst-page-desc">worldgen · started at 14:02</p>
  </div>
  <div class="inst-page-actions">
    <button class="inst-btn inst-btn--sm inst-btn--danger" type="button">Stop</button>
  </div>
</div>
```

## API

```api
```

## Related

```related
```
