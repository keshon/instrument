---
title: Icons
group: Foundations
layout: foundation
source: src/tone.css
api:
  - { name: "inst-sprite", kind: "class", doc: "The sprite in the body of the page. Hides it: without the rule it is a visible block of 150px" }
  - { name: "inst-icon", kind: "class", doc: "The base one. Always required" }
  - { name: "inst-icon--sm", kind: "modifier", doc: "Reduction to `--size-chevron`" }
  - { name: "--icon-size", kind: "variable", doc: "The size of this icon. Set on the CONTAINER and inherited downwards. Through a variable rather than through a descendant selector: the layout lies in an earlier layer, and a selector from there is dead in silence" }
  - { name: "--size-icon", kind: "token", value: "18px", doc: "The default: an icon beside a label at the base type size" }
  - { name: "--size-icon-sm", kind: "token", value: "14px", doc: "An icon at type size `--text-2xs`: the status bar" }
  - { name: "--size-icon-lg", kind: "token", value: "24px", doc: "An icon WITHOUT a label: the rail of sections. Measured against the tap target" }
  - { name: "--size-chevron", kind: "token", value: "10px", doc: "A utility glyph: a chevron, a cross" }
  - { name: "--stroke", kind: "token" }
---

Inline SVG on a grid of 16, a `currentColor` stroke of thickness 1.5. An icon
inherits the colour of the text, which is why the library has not one tone
variant for icons.

```html preview
<button class="inst-btn" type="button">
  <svg class="inst-icon" aria-hidden="true"><use href="#i-refresh"/></svg>
  Refresh
</button>
<button class="inst-btn inst-btn--icon" type="button" aria-label="Settings">
  <svg class="inst-icon" aria-hidden="true"><use href="#i-settings"/></svg>
</button>
```

## Contract

Five requirements. Breaking any of them shows at once — the icon either slips
off the baseline or stops changing colour along with the text.

| Requirement | Why |
|---|---|
| **Inline SVG** (or a `<use>` from the sprite) | An icon has to change colour along with the text. `<img>` and `background-image` cannot do that |
| **`stroke: currentColor`, `fill: none`** | One file works in all five themes and in every tone without a single variant |
| **A stroke thickness of 1.5** | The optics shared with the hairline: 1 reads as a scratch, 2 as a bold face beside a 400 |
| **A grid of 16 (`viewBox="0 0 16 16"`)** | It matches `--size-icon`. Another grid gives another optical thickness at the same `stroke-width` |
| **The colour does not live in a data URI** | The shape is drawn by a mask and the colour comes as a token. Otherwise the theme never catches the icon up |

The class `.inst-icon` declares that plus `flex: none` — otherwise the icon is
the first thing to shrink when the label beside it does not fit.

```css
.inst-icon {
  flex: none;
  inline-size: var(--icon-size, var(--size-icon));
  block-size: var(--icon-size, var(--size-icon));
  fill: none;
  stroke: currentColor;
  stroke-width: 1.5;
  stroke-linecap: round;
  stroke-linejoin: round;
}
```

### Markup

| Case | Markup |
|---|---|
| An icon **beside** a label | `aria-hidden="true"` on the `<svg>`. The name is given by the label; without the attribute a screen reader reads the element twice or unintelligibly |
| An icon **instead of** a label | `aria-hidden="true"` on the `<svg>` **and** `aria-label` on the button. Otherwise a button with no text has no accessible name |
| An icon inside a link | The same: the name is carried by the link |
| An icon as the only carrier of information | There is no such thing. Add a label or an `aria-label` — colour and shape have no right to be the only carrier |

### Accessibility

| | |
|---|---|
| A decorative icon | `aria-hidden="true"` always, without exception: the name is given by a label or an `aria-label` |
| A button with no text | Has to carry an `aria-label`. Otherwise it has no accessible name at all |
| Contrast | An icon inherits the colour of the text, so it passes the same thresholds. An icon on `--text-faint` is decoration and has no right to carry meaning |
| Colour is not the only carrier | An icon is itself a second carrier of tone (a note), but not the first: the tone has to be in the text as well |
| Forced-colours mode | A `currentColor` stroke survives the reset, because the colour is taken from the text rather than set to a value of its own |
| Focus | An `<svg>` never takes focus: what is focused is the button or the link around it, and `focusable="false"` on the sprite removes the old IE artefact |

## Anatomy

The sprite is declared once per document, before the first `<use>`. It is
itself invisible and hidden from a screen reader.

```html
<svg class="inst-sprite" aria-hidden="true" focusable="false" xmlns="http://www.w3.org/2000/svg">
  <symbol id="i-refresh" viewBox="0 0 16 16">…</symbol>
</svg>
```

Use:

```html
<svg class="inst-icon" aria-hidden="true"><use href="#i-refresh"/></svg>
```

The `<svg>` at the place of use has **no `viewBox` of its own**: it comes from
the `<symbol>`. A `viewBox` on the outer element beats the symbol's and cuts
part of the glyphs off.

### Available ids

The set is deliberately small: an icon is started for a task rather than just
in case. The list is printed **from the sprite itself** — typed by hand it has
already parted from the library by four symbols once.

```icons
```

## Sizes

| Class or variable | Size | Token | Where |
|---|---|---|---|
| `inst-icon` | 18px | `--size-icon` | The default: a button, a menu item, a note, a header |
| `inst-icon--sm` | 10px | `--size-chevron` | The chevron of a select, the cross of a tag — a utility glyph |
| `--icon-size: var(--size-icon-sm)` | 14px | `--size-icon-sm` | A reading in [the status bar](../layout/statusbar.md): the type size there is `--text-2xs` |
| `--icon-size: var(--size-icon-lg)` | 24px | `--size-icon-lg` | An icon in [the rail](../layout/rail.md), where there is no label at all |

All four are one measure with different anchors. An icon beside a label is
measured against ITS cap height, so at the eleven-pixel size the icon has one
of its own; an icon without a label is measured against the tap target, because
there is nothing else to compare it with. `cmd/proportion` checks this, and the
bands are there too.

An icon larger than these four is either an illustration (its place is in
[the empty state](../components/feedback/empty.md)) or a mistake of hierarchy:
an icon serves the data and has to be quieter than it.

The sizes deliberately do not react to density — see [density](./density.md):
an icon is a carrier of meaning, and at 12px a stroke of 1.5 turns into a blot.

## Composition

```html
<div class="inst-menu" role="menu">
  <button class="inst-menu-item" type="button" role="menuitem">
    <svg class="inst-icon" aria-hidden="true"><use href="#i-refresh"/></svg>Restart
  </button>
  <button class="inst-menu-item" type="button" role="menuitem">
    <svg class="inst-icon" aria-hidden="true"><use href="#i-copy"/></svg>Copy the id
  </button>
</div>
```

An icon button is **square**: the width is taken from the same control height,
so the square holds at all three sizes and at all densities.

## Rules

### What to take and what not to

| Use | Take instead |
|---|---|
| An icon **beside** a label — a second carrier of meaning | **An icon instead of a label in a list of actions** — in a menu the label is obligatory, see [the menu in a popover](../components/overlays/popover.md) |
| An icon instead of a label where the action is common knowledge and there is an `aria-label` | **An icon instead of a status** — a status carries a dot and a word, see [the badge](../components/display/badge.md) |
| An icon as a second mark of tone on a note | **Colour as the only carrier** — the general rule. Tone is described in [colour](./colors.md) |
| A `<use href="#i-*">` from the sprite | **A text glyph ⟳ or ▶** — different metrics in different fonts and the risk of an emoji substitution |
| `.inst-icon--sm` in a chevron and in a small control | **A width of your own inline** — the size comes from [the roles](./spacing.md) |
| `--icon-size` on a container when a component needs a size of its own | **A descendant selector `.my-block > .inst-icon`** — if your layer comes before `kit.components`, the rule will not fire and will not complain |

## Customization

1. Draw on a 16×16 grid, in strokes only.
2. Remove `fill`, `stroke` and `stroke-width` from the paths themselves — they
   come from `.inst-icon`.
3. Round the caps and joins by the class rather than in the file:
   `stroke-linecap` and `stroke-linejoin` are already there.
4. Put it in the sprite as `<symbol id="i-…" viewBox="0 0 16 16">`.
5. Check it on a dark theme: a stroke of 1.5 is optically thinner in a dark
   theme, and too fine a detail disappears.

## API

```api
```

## Related

[Rhythm and shape](./spacing.md)
[Colour](./colors.md)
[Button](../components/actions/button.md)
[Popover and menu](../components/overlays/popover.md)
[Empty state](../components/feedback/empty.md)
