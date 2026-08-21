---
title: Toolbar
group: Navigation
layout: component
source: src/surfaces.css
api:
  - { name: "inst-toolbar", kind: "class", doc: "The strip. A row with a gap of `--space-3` and a line below" }
  - { name: "inst-toolbar-sep", kind: "class", doc: "A vertical divider between groups" }
  - { name: "inst-toolbar-spacer", kind: "class", doc: "Pushes the tail of the strip to the far edge" }
  - { name: "--space-2", kind: "token" }
  - { name: "--space-3", kind: "token" }
  - { name: "--pad-cell-x", kind: "token" }
  - { name: "--hairline", kind: "token" }
  - { name: "--border-subtle", kind: "token" }
  - { name: "--border", kind: "token" }
---

A strip of tools inside [a panel](../display/panel.md): modes, filters, small
actions on its content. Not a page header and not a menu.

```html preview
<div class="inst-panel">
  <div class="inst-panel-header"><span class="inst-panel-title">Inspector</span></div>
  <div class="inst-toolbar">
    <button class="inst-btn inst-btn--sm inst-btn--ghost inst-btn--icon" type="button" aria-label="Move">
      <svg class="inst-icon" aria-hidden="true"><use href="#i-move"/></svg></button>
    <button class="inst-btn inst-btn--sm inst-btn--ghost inst-btn--icon" type="button" aria-label="Rotate">
      <svg class="inst-icon" aria-hidden="true"><use href="#i-rotate"/></svg></button>
    <button class="inst-btn inst-btn--sm inst-btn--ghost inst-btn--icon" type="button" aria-label="Scale">
      <svg class="inst-icon" aria-hidden="true"><use href="#i-scale"/></svg></button>
    <span class="inst-toolbar-sep"></span>
    <span class="inst-toolbar-spacer"></span>
    <span class="inst-badge">terrain_chunk_04</span>
  </div>
  <div class="inst-panel-body">…</div>
</div>
```

## Contract

| What | Required | Why |
|---|---|---|
| A `type="button"` on every button | yes | Otherwise inside a form it will submit it |
| An `aria-label` on an icon button | yes | A button with no text has no accessible name: "Move" rather than "Tool 1" |
| An `aria-hidden="true"` on the icon | yes | Otherwise a screen reader reads it twice or does not read it at all |
| An `aria-pressed` on a mode button | yes, if the tool stays switched on | The state has to be in the markup. Set by the application |
| One `inst-toolbar-spacer` per strip | yes | The spacer sets off the tail — a status, a counter, a close. There is no such thing as two points of break in one strip |

The `toolbar` role is not assigned automatically: it promises arrow navigation
and a roving `tabindex`, and that behaviour belongs to the application layer. A
strip of ordinary buttons with no role is accessible without the promises.

### Accessibility

| | |
|---|---|
| The keyboard | Every button is traversed by `Tab`. Arrow navigation within the strip comes only along with a `role="toolbar"` and a roving `tabindex`, that is, with the application's JS |
| The names | Icon buttons have to carry an `aria-label`: "Move" rather than "Tool 1" |
| The tap target | `inst-btn--sm` is 26px. At density `compact` it is 22px, below the 24px threshold of WCAG 2.5.8 |
| The divider | An empty element with no text and no role — it is not in the accessibility tree, and that is right: it says nothing beyond the order |
| Contrast | The lines `--border-subtle` and `--border` are decoration (a 3:1 threshold) rather than text |
| Print | A toolbar is not printed: it is the chrome of the application rather than data |

## Anatomy

Two service elements that are easy to confuse:

| Element | What it does | When |
|---|---|---|
| `inst-toolbar-sep` | A vertical line the full height of the strip | Between **meaningful groups** of tools: the transforms apart, the display apart |
| `inst-toolbar-spacer` | `margin-inline-start: auto` — everything after it travels to the far edge | Once per strip, before the tail: a status, a counter, a close button |

The divider is an empty `<span>`, that is, pure graphics with no text and no
role. It shows the grouping to the eye; for a screen reader the grouping is
carried by the order and by the labels of the buttons.

## Composition

### Its place in a panel

A toolbar stands **between the header and the body** of a panel and brings the
line below itself. The header then removes its own line
(`:has(+ .inst-toolbar)`) — otherwise there would be two rules in a row within
thirty pixels.

```html preview context
<div class="inst-panel">
  <div class="inst-panel-header"><span class="inst-panel-title">Modules</span></div>
  <div class="inst-toolbar">
    <div class="inst-segmented" role="radiogroup" aria-label="View">
      <button type="button" role="radio" aria-checked="true"  tabindex="0">List</button>
      <button type="button" role="radio" aria-checked="false" tabindex="-1">Grid</button>
    </div>
    <span class="inst-toolbar-sep"></span>
    <span class="inst-search">
      <input class="inst-input inst-input--sm" type="search" placeholder="Filter" aria-label="Filter the modules">
    </span>
    <span class="inst-toolbar-spacer"></span>
    <span class="inst-badge">3 of 12</span>
  </div>
  <div class="inst-panel-body">The content of the region.</div>
</div>
```

A toolbar with no header is legitimate too: a panel that has no name but has
modes.

## API

```api
```

## Related

```related
```
