---
title: Checkbox, radio and switch
group: Inputs
layout: component
source: src/forms.css
api:
  - { name: "inst-checkbox", kind: "class", doc: "The label wrapper of a checkbox" }
  - { name: "inst-radio", kind: "class", doc: "The same for a radio" }
  - { name: "inst-switch", kind: "class", doc: "The same for a switch" }
  - { name: "--size-check", kind: "token" }
  - { name: "--size-switch-w", kind: "token" }
  - { name: "--size-switch-h", kind: "token" }
  - { name: "--size-thumb", kind: "token" }
  - { name: "--border-control", kind: "token" }
  - { name: "--accent-mark", kind: "token" }
  - { name: "--accent-on", kind: "token" }
  - { name: "--track", kind: "token" }
  - { name: "--radius-xs", kind: "token" }
  - { name: "--radius-full", kind: "token" }
  - { name: "--dur-1", kind: "token" }
---

A checkbox, a radio and a switch. All three are built on a native `input`: the
keyboard, form submission and screen readers come free.

```html preview
<label class="inst-checkbox"><input type="checkbox" checked>Cache the result</label>
<label class="inst-radio"><input type="radio" name="mode" checked>In parallel</label>
<label class="inst-radio"><input type="radio" name="mode">Sequentially</label>
<label class="inst-switch"><input type="checkbox" checked>Autostart</label>
```

## Contract

The control lies **inside** the `<label>`, so `for` and `id` are not needed,
and the whole label becomes the tap area.

| What | Required | Why |
|---|---|---|
| An `<input>` inside a `<label>` | yes | `for` and `id` are not needed, and the tap area becomes the whole label — the one way to meet WCAG 2.5.8 with a glyph of 15px |
| A shared `name` on a radio group | yes | Without it the options are not mutually exclusive, and keyboard navigation by arrows will not work |
| A `<fieldset>` with a `<legend>` around a radio group | yes | Otherwise a screen reader reads the options but not the question |

The `input` has no class of its own: it is recognised by its type and its
parent.

The keyboard is entirely native, and `instrument.js` intercepts nothing here:
`Space` toggles, in a radio group the arrows traverse the options, and `Tab`
enters and leaves the group.

The border of an unchecked one takes `--border-control` and holds 3:1: on all
three the boundary **is** the control, and for the same reason a ring is added
to the track of a switch. Colour is not the only mark: a checkbox carries a
tick, a radio a dot, a switch the position of the thumb.

A switch stays a native checkbox. `role="switch"` is added if the application
wants the state spoken as "on/off" rather than "checked".

## Variants

The choice among the three is settled by **when the value applies**.

```html preview
<label class="inst-checkbox"><input type="checkbox" checked>Checkbox</label>
<label class="inst-radio"><input type="radio" name="v" checked>Radio</label>
<label class="inst-switch"><input type="checkbox" checked>Switch</label>
```

| | When | It applies |
|---|---|---|
| **Checkbox** | An independent flag. Any number of them may be checked, zero included | On submitting the form |
| **Radio** | One of several mutually exclusive. **Exactly one** is always chosen | On submitting the form |
| **Switch** | Switching a mode on and off | **Immediately**, with no "Apply" button |

## States

```html preview
<label class="inst-checkbox"><input type="checkbox">Ordinary</label>
<label class="inst-checkbox"><input type="checkbox" checked>Checked</label>
<label class="inst-checkbox"><input type="checkbox" disabled>Unavailable</label>
<label class="inst-checkbox"><input type="checkbox" checked disabled>Checked and unavailable</label>
```

The indeterminate state of a checkbox is the `indeterminate` property; it is
set by script alone and has no attribute in the markup:

```js
el.indeterminate = true;   // "some of the nested ones are checked"
```

## API

```api
```

## Related

```related
```
