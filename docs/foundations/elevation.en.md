---
title: Elevation and surfaces
group: Foundations
layout: foundation
source: src/tokens.css
api:
  - { name: "--surface-sunken", kind: "token" }
  - { name: "--surface-page", kind: "token" }
  - { name: "--surface-raised", kind: "token" }
  - { name: "--surface-overlay", kind: "token" }
  - { name: "--surface-field", kind: "token" }
  - { name: "--surface-hover", kind: "token" }
  - { name: "--surface-active", kind: "token" }
  - { name: "--surface-selected", kind: "token" }
  - { name: "--shadow-popover", kind: "token" }
  - { name: "--shadow-modal", kind: "token" }
  - { name: "--shadow-color-near", kind: "token" }
  - { name: "--shadow-color-far", kind: "token" }
  - { name: "--scrim", kind: "token" }
  - { name: "--z-sticky", kind: "token" }
  - { name: "--z-popover", kind: "token" }
  - { name: "--z-modal", kind: "token" }
  - { name: "--border", kind: "token" }
  - { name: "--border-subtle", kind: "token" }
  - { name: "--hairline", kind: "token" }
---

Depth is carried by the **order of lightness**. A shadow means one thing here —
"it floats above and is about to disappear" — and it has two carriers: the
popover and the modal.

```html preview context
<div class="inst-panel">
  <div class="inst-panel-header">
    <span class="inst-panel-title">Run #4127</span>
  </div>
  <div class="inst-panel-body inst-stack">
    <div class="inst-card">
      <div class="inst-card-title">worldgen-01</div>
      <div class="inst-card-sub">A card on a panel: a step of surface, not one shadow</div>
    </div>
    <input class="inst-input" type="text" value="terrain_chunk_04">
  </div>
</div>
```

## Contract

Depth is carried by the **surface**. A shadow means "it floats above and is
about to disappear" and belongs to three components.

```css
.my-panel {
  background: var(--surface-raised);
  border: var(--hairline) solid var(--border);
}
```

| What | Required | Why |
|---|---|---|
| A surface from the semantics | yes | Four steps named by work rather than by number: the page, the recess, the raised one, the overlay |
| A shadow only on something popping up | yes | A popover, a modal, a notification. There is no fourth carrier of a shadow in the library |
| A `CanvasText` border in `forced-colors` | yes | The surfaces collapse in that mode, and without a border a block merges with the background |
| A `z-index` from the vocabulary | yes | Four tokens. A number invented in place works until the first larger number of somebody else's |


### Accessibility

| | |
|---|---|
| Depth not by colour alone | Every step of surface comes with a border or a change of rhythm: telling the stack apart does not require telling lightnesses apart |
| The contrast of the steps | Neighbouring surfaces are checked as a pair in all five themes — otherwise a field on a panel disappears |
| A shadow does not carry meaning alone | An overlay is recognised by its role and behaviour (`popover`, `<dialog>`) rather than by a shadow. In forced-colours mode the shadow disappears, and nothing breaks |
| Scrolling under a modal | It stops: `showModal()` makes the background inert to clicks but not to the wheel, and the user lost the place they were coming back to |
| Reduced motion | The appearance of an overlay collapses to 0.01ms but is not switched off — the state machines listening for `transitionend` keep working. See [motion](./motion.md) |

## Anatomy

Declared in one place so that it does not drift across components.

| Token | Value | Who |
|---|---|---|
| `--z-sticky` | 10 | The sticky head of a table, the action bar of a form |
| `--z-popover` | 100 | A popover, a menu, a tooltip |
| `--z-modal` | 1000 | A modal, a drawer |

A popover and a modal built on the Popover API and `<dialog>` live in the
browser's top layer, where `z-index` no longer decides. The tokens stay for
what does not get the top layer — and for an application that needs to fit into
the same order.

## Scale

They read as a stack from the bottom up. A step is chosen by the **role of the
region** rather than by a wanted brightness.

| Token | Step | Who carries it |
|---|---|---|
| `--surface-sunken` | The bottom | A track, the head of a table, a block of code, the body of a step, the background of a side column |
| `--surface-page` | The page | `body`, the background of the working area |
| `--surface-raised` | Raised | A panel, a card, the header of the shell, the action bar of a form, `<kbd>` |
| `--surface-overlay` | The overlay | A popover, a menu, a tooltip, a modal, a drawer |

The fifth role stands **outside** the stack, because it is not about depth but
about input:

| Token | Work |
|---|---|
| `--surface-field` | The recess under a field. A role of its own rather than a reuse of `raised`: while a field was painted the colour of a panel, its only boundary was a border at 1.31:1 — that is, a field on a panel was not visible at all |

In a dark theme the distance between neighbouring steps is **wider** than in a
light one: the dark end of the ramp is harder to tell apart, and a mirrored
light arrangement would give four identical black rectangles. In the `dark-soft`
theme the whole stack rises by one step of the ramp — four surfaces are
overridden, and the rest of the semantics does not know the theme has changed.
There is no room above that: the next step of the ramp no longer holds text.
The light end behaves in exactly the same way, where `raised` and `overlay`
both sit on `--n-0`.

### States over a surface

Hover and selection are **alpha**. An opaque fill would keep them from adding
up: a hover over a selected row would disappear.

| Token | Work |
|---|---|
| `--surface-hover` | Hover |
| `--surface-active` | Press |
| `--surface-selected` | A selected row |

## Variants

Both are functional. Both are declared entire in the tokens — a component
substitutes a ready value and does not assemble a shadow itself.

| Token | Value | Who carries it |
|---|---|---|
| `--shadow-popover` | `0 4px 12px -2px` far + `0 2px 4px -2px` near | A popover, a menu, a tooltip |
| `--shadow-modal` | `0 16px 40px -8px` far + `0 4px 10px -4px` near | A modal, a drawer |

Every shadow has two layers: the far one gives the distance to the background,
the near one the contact of the edge. In one layer it looks either glued down
or floating with no support.

The colour of a shadow is kept separate, because `light-dark()` accepts only a
`<color>` and a whole shadow value does not pass through it:

| Token | Light | Dark |
|---|---|---|
| `--shadow-color-near` | `oklch(0 0 0 / 0.06)` | `oklch(0 0 0 / 0.30)` |
| `--shadow-color-far` | `oklch(0 0 0 / 0.10)` | `oklch(0 0 0 / 0.45)` |

In a dark theme a shadow is **denser**. A weak shadow on a dark background is
invisible, and a light shadow is a glow already, and glows are forbidden.

### The rule

**A shadow is given only to what appeared on an action and will disappear on
Escape.** Everything that stands on the screen permanently is parted by a
surface and a border. A card gets no shadow — it gets a surface.

| Element | Shadow |
|---|---|
| A panel, a card, a table, an action bar | None |
| A popover, a menu, a tooltip | `--shadow-popover` |
| A modal, a drawer | `--shadow-modal` |
| Everything else | None, and no new recipient is started |

### The ground

| Token | Light | Dark |
|---|---|---|
| `--scrim` | `oklch(0 0 0 / 0.32)` | `oklch(0 0 0 / 0.58)` |

In a dark theme the ground is denser: on a light page a dimming creates
contrast by itself, on a dark one it has to be added. This is a **dimming, not
a blur** — gradients, glows and blurs are forbidden entirely.

The ground is put on the `::backdrop` of a native `<dialog>`, so the top layer,
the focus trap and closing on Escape come from the platform rather than from
the library.

## Rules

### What to take and what not to

| Use | Take instead |
|---|---|
| A step of surface for permanent structure | **A shadow on a card or a panel** — a shadow means "temporary". Structure is built by [the card](../components/display/card.md) and [the panel](../components/display/panel.md) |
| `--shadow-popover` for what will disappear on Escape | **A shadow "for volume"** — it has no recipient other than [the popover](../components/overlays/popover.md) |
| `--shadow-modal` + `--scrim` for what blocks the screen | **A popover instead of a modal** when a decision is obligatory — take [the modal](../components/overlays/dialog.md) |
| `--surface-field` under input | **`--surface-raised` under input** — the field will merge with the panel, see [the field](../components/inputs/input.md) |
| Blurring the background under an overlay | **Blurs are not used** — the dimming is done by [`--scrim`](#the-ground) |

## API

```api
```

## Related

[Colour](./colors.md)
[Motion](./motion.md)
[Panel](../components/display/panel.md)
[Card](../components/display/card.md)
[Popover](../components/overlays/popover.md)
[Modal](../components/overlays/dialog.md)
