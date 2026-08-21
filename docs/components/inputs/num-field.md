---
title: Inspector number field
group: Inputs
layout: component
source: src/table.css
js: Dragging the label of an axis is done by `instrument.js`
api:
  - { name: "inst-num-field", kind: "class", doc: "The wrapper: the border, the recess, the focus. The height is `--control-h-sm`" }
  - { name: "inst-num-axis", kind: "class", doc: "The letter of the axis. A drag handle, **not a label**" }
  - { name: "inst-vec", kind: "class", doc: "A row of several fields of equal width" }
  - { name: "inst-prop", kind: "class", doc: "A property row of the inspector" }
  - { name: "inst-prop-label", kind: "class", doc: "The name of a property. It is truncated, so a `title` is required" }
  - { name: "inst-prop-control", kind: "class", doc: "The control part of the row" }
  - { name: "--control-h-sm", kind: "token" }
  - { name: "--radius-sm", kind: "token" }
  - { name: "--surface-field", kind: "token" }
  - { name: "--border-control", kind: "token" }
  - { name: "--accent-border", kind: "token" }
  - { name: "--text-2xs", kind: "token" }
  - { name: "--text-xs", kind: "token" }
  - { name: "--space-2", kind: "token" }
  - { name: "--space-3", kind: "token" }
---

A number with the letter of an axis at the left — a primitive of the
inspector: three such fields make a vector, and a property panel holds dozens
of them. More compact than [the text field](./input.md), because the label of
the axis takes the place of the stepper rather than a column of its own.

```html preview
<span class="inst-num-field">
  <span class="inst-num-axis">X</span>
  <input type="number" value="128" aria-label="Position X">
</span>
<span class="inst-num-field">
  <span class="inst-num-axis">Y</span>
  <input type="number" value="0" aria-label="Position Y">
</span>
<span class="inst-num-field">
  <span class="inst-num-axis">Z</span>
  <input type="number" value="-64" aria-label="Position Z">
</span>
```

## Contract

It is easy to get a field here that looks labelled and is not.

| What | Required | Why |
|---|---|---|
| An `aria-label` on the `<input>` **including the axis** | yes | Otherwise a screen reader reads "field, 128" and does not say of what. The axis is part of the name: "Position X" rather than "X" |
| `inst-num-axis` is **not** a `<label>` | yes | It is a drag handle and a visual marker rather than a label. A `<label for>` would make the letter "X" the full accessible name of the field — instead of the name of the property there would be one word |
| `type="number"` | yes | The keyboard arrows, the step and the numeric keypad on touch come from the platform |
| An `<input>` as a direct child | yes | The rules are written through `& > input` |
| A `title` on `inst-prop-label` | yes, in a property row | Otherwise a truncated property name cannot be read |

Hence a consequence that is easier to break than to notice: **the letter of the
axis must not be the only place saying which field this is**. The name of the
property lives in the `inst-prop-label` beside it and is repeated in the
`aria-label` of every field — three identical squares with the letters X, Y, Z
mean nothing outside the context of the row.

:::note
**The dragging is done by [instrument.js](#js).** An `ew-resize` cursor on the
label of an axis means "drag from here", and the module keeps that promise:
`Shift` speeds it up tenfold, `Alt` slows it down tenfold.

Without the module the field stays an ordinary number input, and that works
entire: digits, arrows, the wheel, pasting.
:::

The `aria-label` has to include the axis **and** the property: "X" is not a
name. `inst-num-axis` itself is not marked up as a label — otherwise it would
beat the `aria-label` and become the full name of the field.

The keyboard is native: the arrows change the value by `step`, `PageUp` and
`PageDown` by ten steps. The steppers are removed by styling, the behaviour is
untouched. The focus ring hangs on the wrapper through `:focus-within`, so the
whole control is seen rather than the digits alone.

The border takes `--border-control` and holds 3:1: here the boundary is the
control. In `forced-colors` the wrapper gets a `CanvasText` border — without it
the field collapses into the surface of the panel. The height `--control-h-sm`
gives 26px, and 22px in the compact mode, that is **below the 24px of WCAG
2.5.8**: a deliberate price of an inspector panel, where a screen holds dozens
of fields.

## States

```html preview
<span class="inst-num-field">
  <span class="inst-num-axis">X</span>
  <input type="number" value="128" aria-label="Position X">
</span>
<span class="inst-num-field">
  <span class="inst-num-axis">Y</span>
  <input type="number" value="0" aria-label="Position Y" readonly>
</span>
<span class="inst-num-field">
  <span class="inst-num-axis">Z</span>
  <input type="number" value="-64" aria-label="Position Z" disabled>
</span>
```

| State | How it is set | What happens |
|---|---|---|
| hover | `:hover` on the wrapper | The border darkens to `--text-muted` |
| focus | `:focus-within` on the wrapper | The border in the accent. The ring is drawn by the **wrapper** rather than by the field: the field has no border at all |
| read-only | `readonly` on the `<input>` | The value is read and copied. Not `disabled` — the difference is in [the text field](./input.md) |
| unavailable | `disabled` on the `<input>` | Native unavailability |

The native stepper arrows are removed (`appearance: textfield`): in a dense row
of the inspector they ate a third of the width of a field and came under the
cursor more often than the input itself. The keyboard arrows keep working — the
styling was removed rather than the behaviour.

## JS

Include the module once per page — there is no need to initialise the
components one by one, `instrument.js` works by delegation and sees nodes that
arrived later.

```html
<script type="module" src="instrument.js"></script>
```

### What `instrument.js` does

The example in the header is live: drag the letter of the axis left and right,
with `Shift` ten times faster.

Dragging the letter of an axis changes the value. The `ew-resize` cursor on it
is a promise given by a picture, and `instrument.js` keeps it.

| Modifier | The multiplier of the step |
|---|---|
| — | ×1 |
| `Shift` | ×10 |
| `Alt` | ×0.1 |

The bounds are taken from `min` and `max`, the step from `step`. A field with
`disabled` or `readonly` is not dragged.

### Events

There are no events of its own: `instrument.js` changes the `<input>` and sends
the native `input` while dragging and `change` on release. A framework sees
them with not one line of glue.

```js
input.addEventListener('change', (e) => scene.setX(+e.target.value));
```

### Options

| Attribute | What it does |
|---|---|
| `step` | The step of the dragging. 1 by default |
| `min`, `max` | The bounds. Without them the value is unbounded |

```html
<span class="inst-num-field">
  <span class="inst-num-axis">X</span>
  <input type="number" value="0.5" step="0.01" min="0" max="1"
         aria-label="Opacity X">
</span>
```

## Composition

### A vector

```html preview
<span class="inst-vec">
  <span class="inst-num-field"><span class="inst-num-axis">X</span><input type="number" value="128" aria-label="Position X"></span>
  <span class="inst-num-field"><span class="inst-num-axis">Y</span><input type="number" value="0" aria-label="Position Y"></span>
  <span class="inst-num-field"><span class="inst-num-axis">Z</span><input type="number" value="-64" aria-label="Position Z"></span>
</span>
```

`inst-vec` hands the three fields equal shares of the width, so a vector does
not skew when a six-digit number is typed into one component.

### A single value with the axis as a name

```html preview
<span class="inst-num-field">
  <span class="inst-num-axis">R</span>
  <input type="number" value="45" step="5" aria-label="Radius of influence">
</span>
```

## Patterns

### A property row of the inspector

```html preview context
<div class="inst-prop">
  <span class="inst-prop-label" title="Position">Position</span>
  <span class="inst-prop-control inst-vec">
    <span class="inst-num-field"><span class="inst-num-axis">X</span><input type="number" value="128" aria-label="Position X"></span>
    <span class="inst-num-field"><span class="inst-num-axis">Y</span><input type="number" value="0" aria-label="Position Y"></span>
    <span class="inst-num-field"><span class="inst-num-axis">Z</span><input type="number" value="-64" aria-label="Position Z"></span>
  </span>
</div>
```

## API

```api
```

## Related

```related
```
