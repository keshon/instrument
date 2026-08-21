---
title: Rail
group: Layout
layout: component
source: src/layout.css
api:
  - { name: "inst-rail", kind: "class", doc: "A column of glyphs. The width is intrinsic: the padding plus the button" }
  - { name: "inst-rail-item", kind: "class", doc: "A section. A link or a button. An `aria-label` is **required**" }
  - { name: "inst-rail-spacer", kind: "class", doc: "A break: everything after it is pushed to the bottom" }
  - { name: "inst-rail-count", kind: "class", doc: "The room for a counter in the corner of a glyph. Put on `.inst-badge--count`" }
  - { name: "aria-current", kind: "attribute", value: "page", doc: "The current section: a fill plus a mark at the edge" }
  - { name: "aria-label", kind: "attribute", doc: "**Required.** A glyph with no name has no accessible name at all" }
  - { name: "--icon-size", kind: "variable", value: "--size-icon-lg", doc: "Set on the item and inherited by the glyph. Through a variable, because the layout lies in an earlier layer than the glyph" }
  - { name: "--size-icon-lg", kind: "token", value: "24px", doc: "A glyph with no label. Larger than usual deliberately" }
  - { name: "--control-h-lg", kind: "token", doc: "The button. Larger than `--tap-min`, otherwise the glyph fills it entire" }
  - { name: "--size-marker", kind: "token", doc: "The mark of the current one at the edge" }
  - { name: "--surface-sunken", kind: "token" }
  - { name: "--surface-selected", kind: "token" }
  - { name: "--accent-mark", kind: "token" }
---

Vertical navigation in glyphs alone, pushed to the edge of the window. The
sections are recognised **by the shape of the glyph** rather than by a word, so
there is no label in a rail — it moves into the tooltip and into the accessible
name.

```html preview
<nav class="inst-rail" aria-label="Sections">
  <a class="inst-rail-item" href="#rail" aria-current="page" aria-label="Files">
    <svg class="inst-icon" aria-hidden="true"><use href="#i-list"/></svg>
  </a>
  <a class="inst-rail-item" href="#rail" aria-label="Projects">
    <svg class="inst-icon" aria-hidden="true"><use href="#i-grid"/></svg>
  </a>
  <a class="inst-rail-item" href="#rail" aria-label="Runs">
    <svg class="inst-icon" aria-hidden="true"><use href="#i-refresh"/></svg>
    <span class="inst-badge inst-badge--count inst-rail-count">11</span>
  </a>
  <a class="inst-rail-item inst-rail-spacer" href="#rail" aria-label="Parameters">
    <svg class="inst-icon" aria-hidden="true"><use href="#i-settings"/></svg>
  </a>
</nav>
```

## Contract

A rail is a `<nav>` with a name of its own and a row of links inside. The glyph
has to carry an `aria-hidden`, and the name of the section has to live in the
`aria-label` of the link itself.

| What | Required | Why |
|---|---|---|
| An `aria-label` on every item | yes | A button with no text has no accessible name. Without it a rail reads as a row of empty buttons |
| An `aria-hidden` on the glyph | yes | Otherwise a screen reader reads out the contents of the SVG over the name |
| An `aria-current="page"` on the current one | yes | The fill reports the choice to the eye, the attribute to everything else |
| A tooltip on hover | no, but wanted | The shape of a glyph is not recognised by everybody and not at once |
| An `<a>` for going to an address | yes | A section has an address, and the middle click has to work |

### Accessibility

| | |
|---|---|
| The accessible name | From the `aria-label` alone. A glyph is a picture and there is no text in it |
| The tap target | `--control-h-lg`, that is 38px at the ordinary density and 32 at the compact one — above the norm of 24 in both |
| The current section | An `aria-current="page"` rather than colour alone: colour as the only carrier is forbidden |
| The contrast of the glyph | `--text-secondary` at rest, `--accent-text` on the current one — both above 4.5:1 in five themes |
| The landmark | A `<nav>` with an `aria-label`, otherwise the page has two nameless navigation landmarks |

## Anatomy

| The part | What it does |
|---|---|
| `.inst-rail` | The column. The width is not set as a number: it is made of the padding and the button, so it travels with density |
| `.inst-rail-item` | A button of `--control-h-lg`. **Larger than `--tap-min` deliberately**: a 24px glyph in a 24px target would fill it entire |
| `.inst-rail-spacer` | It pushes the rest to the bottom. The account and the parameters stand at the bottom rather than being set off by a margin |
| `.inst-rail-count` | The corner of the glyph, on the outside. From the inside a counter would cover the shape the section is recognised by |

The current section carries **two carriers**: the fill and the mark at the
edge. The fill alone is not enough — on the muted surface of a rail it barely
differs from a hover; the mark is seen by peripheral vision and does not argue
with the glyph.

## API

```api
```

## Related

```related
```
