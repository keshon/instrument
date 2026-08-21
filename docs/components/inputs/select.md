---
title: Select
group: Inputs
layout: component
source: src/forms.css
js-optional: A list of your own instead of the native one. The native one works with no script
api:
  - { name: "inst-select-wrap", kind: "class", doc: "The wrapper. **Required** — the chevron lives in it" }
  - { name: "inst-select", kind: "class", doc: "The `<select>` itself" }
  - { name: "inst-select--sm", kind: "modifier", doc: "The size" }
  - { name: "inst-select--lg", kind: "modifier", doc: "The size" }
  - { name: "--surface-field", kind: "token" }
  - { name: "--control-h-sm/md/lg", kind: "token" }
  - { name: "--control-pad-sm", kind: "token" }
  - { name: "--size-chevron", kind: "token" }
  - { name: "--space-7", kind: "token" }
  - { name: "--radius-sm", kind: "token" }
  - { name: "--border-control", kind: "token" }
  - { name: "--text-muted", kind: "token" }
---

Choosing one value from a list. A native `<select>` — with the whole keyboard,
first-letter search and the system dropdown, which on a phone becomes a wheel.

```html preview
<span class="inst-select-wrap">
  <select class="inst-select" aria-label="Model">
    <option>opus</option>
    <option>sonnet</option>
    <option>haiku</option>
  </select>
</span>
```

## Contract

The wrapper is required — the chevron lives in it: the browser does not hand
over the pseudo-elements of a `<select>` itself.

| What | Required | Why |
|---|---|---|
| `inst-select-wrap` on the wrapper | yes | The chevron is drawn by a pseudo-element of the wrapper. Without it the field is left with no arrow and is indistinguishable from a text one |
| `inst-select` on the field | yes | The recess, the border, the height and the states come from there |
| A `<label for>` or an `aria-label` | yes | The first `<option>` is not a label |

The keyboard is entirely native: arrows, `Home` and `End`, first-letter search,
`Enter`. The label is given by a `<label for>` or an `aria-label` — the first
`<option>` is not a label.

The chevron is drawn by a mask and gets `pointer-events: none`, so that a click
on the arrow opens the list rather than falling past it; in `forced-colors` the
mask loses its background colour and is therefore moved to `CanvasText`
explicitly.

The native list stays native deliberately: on a phone it turns into a system
wheel, and a home-made dropdown loses that.

## Sizes

```html preview
<span class="inst-select-wrap">
  <select class="inst-select inst-select--sm" aria-label="sm"><option>sm</option></select>
</span>
<span class="inst-select-wrap">
  <select class="inst-select" aria-label="md"><option>md</option></select>
</span>
<span class="inst-select-wrap">
  <select class="inst-select inst-select--lg" aria-label="lg"><option>lg</option></select>
</span>
```

They are symmetrical with those of the button and the field: a select, a
[field](./input.md) and a [button](../actions/button.md) of the same size stand
in one row of a form and cannot part.

## States

```html preview
<span class="inst-select-wrap">
  <select class="inst-select" aria-label="Ordinary"><option>Ordinary</option></select>
</span>
<span class="inst-select-wrap">
  <select class="inst-select" aria-label="Unavailable" disabled><option>Unavailable</option></select>
</span>
<span class="inst-select-wrap">
  <select class="inst-select" aria-label="Error" aria-invalid="true"><option>Error</option></select>
</span>
```

The recess, the border, focus, error, `disabled` and `readonly` are the same as
on [the text field](./input.md): all three controls are declared by one rule,
so they cannot part.

## JS

A native `<select>` calls for not one line: the keyboard, typing by first
letters, a full-screen list on a phone — all of it from the platform. A script
is wanted only where the native one is not enough: search within the list, two
lines in an item, an icon on a row.

### When it is worth leaving the native one

| | |
|---|---|
| More than two or three dozen options | Search is wanted, and a `<select>` has none |
| An item holds more than text | The native list draws one line with no markup |
| Choosing several | A `<select multiple>` looks like a list box and is awkward on a phone |

In every other case the native one wins, and replacing it is work with no gain.

### A list of your own

The surface and the keyboard are given by the library. What is left to the
application is the data: opening, filtering, writing the choice into a hidden
field.

```html preview
<div class="inst-field" style="max-inline-size:16rem">
  <label class="inst-label" for="pick">Run</label>
  <span class="inst-select-wrap">
    <button class="inst-select" type="button" id="pick"
            popovertarget="pick-list" aria-haspopup="listbox">worldgen-01</button>
  </span>
  <div class="inst-popover inst-popover--fill inst-popover--anchored" id="pick-list" popover>
    <input class="inst-input inst-input--sm" type="search" aria-label="Search the runs" placeholder="Search">
    <div role="listbox" aria-label="Runs">
      <div class="inst-menu-item" role="option" aria-selected="true" tabindex="0">worldgen-01</div>
      <div class="inst-menu-item" role="option" aria-selected="false" tabindex="-1">biomes-04</div>
      <div class="inst-menu-item" role="option" aria-selected="false" tabindex="-1">rivers-12</div>
    </div>
  </div>
</div>
```

The arrows, `Home`, `End` and carrying the selection inside the list are done
by [`instrument.js`](../../foundations/behavior.md) — this is a
`role="listbox"`. Opening and closing are taken by the Popover API.

### Events

```js
list.addEventListener('inst:select', (e) => {
  trigger.textContent = e.detail.value;
  hidden.value = e.detail.value;      // the value goes with the form
  list.hidePopover();
});
```

The hidden field is obligatory: a list of your own is markup, and a form sends
fields. Without it the choice is visible to a person and invisible to the
server.

## API

```api
```

## Related

```related
```
