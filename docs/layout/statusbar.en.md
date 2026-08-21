---
title: Status bar
group: Layout
layout: component
source: src/layout.css
api:
  - { name: "inst-statusbar", kind: "class", doc: "The strip at the bottom edge. One line in height" }
  - { name: "inst-statusbar-item", kind: "class", doc: "A reading. A `<span>` if it is not pressed, a `<button>` or an `<a>` if it is" }
  - { name: "inst-statusbar-end", kind: "class", doc: "The right group: everything after it goes to the far edge" }
  - { name: "data-tone", kind: "attribute", value: "ok · warn · error · running · neutral", doc: "The tone of a reading. It takes `--tone-ink`" }
  - { name: "--control-h-xs", kind: "token", doc: "The area of a reading. The shape does not thicken by it" }
  - { name: "--text-2xs", kind: "token", value: "11px", doc: "The type size of the strip" }
  - { name: "--icon-size", kind: "variable", value: "--size-icon-sm", doc: "Set on a reading and inherited by the glyph" }
  - { name: "--size-icon-sm", kind: "token", value: "14px", doc: "The glyph of the strip. The usual 18 against a type size of 11 outweighs the reading" }
  - { name: "--surface-sunken", kind: "token" }
  - { name: "--border-subtle", kind: "token" }
---

A strip at the bottom edge with groups at the sides: the branch and the errors
at the left, the mode and the encoding at the right. What stands here are
**readings** rather than controls — they are read, though almost every one of
them can be pressed.

```html preview
<div class="inst-statusbar">
  <button class="inst-statusbar-item" type="button">master</button>
  <span class="inst-statusbar-item" data-tone="error">
    <svg class="inst-icon" aria-hidden="true"><use href="#i-close"/></svg>2 errors</span>
  <span class="inst-statusbar-item" data-tone="warn">7 remarks</span>
  <span class="inst-statusbar-end">
    <span class="inst-statusbar-item">Ln 1, Col 1</span>
    <span class="inst-statusbar-item">UTF-8</span>
    <button class="inst-statusbar-item" type="button">Go</button>
  </span>
</div>
```

## Contract

A reading that cannot be pressed stays a `<span>`. A pressable one is a real
button or a link: a response to hover appears on those alone, and that is
honesty rather than styling.

| What | Required | Why |
|---|---|---|
| A `<span>` for what is not pressed, a `<button>` for what is | yes | Hover is declared on a button and a link alone. A strip of identical `<div>`s lies about what will respond |
| An `aria-hidden` on a glyph beside a number | yes | Otherwise the "2" is read together with the contents of the SVG |
| A word or a glyph along with the tone | yes | Colour as the only carrier is forbidden: a red "2" by itself does not report that these are errors |
| One line | yes | A strip that went to two lines eats the bottom of the screen of the very tool the screen was opened for |

### Accessibility

| | |
|---|---|
| The role | An ordinary container. A `role="status"` is set **only** on the reading that changes by itself and has to be announced — the whole strip gets no such role, otherwise every movement of the cursor is read out |
| The tap target | `--control-h-xs` plus the gap between neighbours. Checked by `cmd/targets` |
| The tone | A `data-tone` gives the colour, but the carrier stays the sign or the word beside it |
| Contrast | `--text-secondary` on `--surface-sunken`, above 4.5:1 in five themes. The type size is `--text-2xs` — the floor of the scale, and there is no size below it |
| The order of traversal | It matches the order in the markup: the left group, then the right one |

## API

```api
```

## Related

```related
```
