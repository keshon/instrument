---
title: Search
group: Inputs
layout: component
source: src/forms.css
js-optional: Live filtering. The field and the clearing work with no script
api:
  - { name: "inst-search", kind: "class", doc: "The wrapper. It draws the magnifier and leaves room for it" }
  - { name: "inst-input", kind: "class", doc: "The field itself. Every size and state comes from there" }
  - { name: "--control-pad-sm", kind: "token" }
  - { name: "--size-icon", kind: "token" }
  - { name: "--size-chevron", kind: "token" }
  - { name: "--space-3", kind: "token" }
  - { name: "--text-muted", kind: "token" }
  - { name: "--text-primary", kind: "token" }
---

A text field with a magnifier and a native clear button. A wrapper around an
`input[type=search]` rather than a control of its own: a button of our own
would call for JS.

```html preview
<div class="inst-field">
  <label class="inst-label" for="q">Search the runs</label>
  <span class="inst-search">
    <input class="inst-input" id="q" type="search" placeholder="worldgen">
  </span>
</div>
```

## Contract

Three parts, and the order among them is obligatory.

| What | Required | Why |
|---|---|---|
| `inst-search` on the wrapper | yes | The magnifier is drawn by a pseudo-element of the wrapper: an `<input>` itself has no pseudo-elements |
| `inst-input` on the field | yes | The whole styling of the field comes from there. `inst-search` adds only the padding for the magnifier and the clearing. The padding is set by the selector `& > .inst-input` — a field with no class will stand under the glyph |
| `type="search"` | yes | The clear cross is **native**. Without that type it does not exist, and this is an ordinary text field with a magnifier |
| A `<label for>` or an `aria-label` | yes | A `placeholder` is not a label: it disappears on input |

The magnifier carries `pointer-events: none`: a click on the glyph lands in the
field rather than in the emptiness in front of it.

A visible `<label for>` is preferable; an `aria-label` is for when the field
stands in a toolbar and a label would eat a line. The type `search` declares
the field as a search field by itself, and `role="search"` need not be set by
hand.

The magnifier is a pseudo-element, and for a screen reader it is not there at
all. That is right: the glyph is decoration, and the work is named by the
label. The glyph takes `--text-muted` for the same reason, and the cross
darkens to `--text-primary` on hover.

The clearing is the browser's own button, and the keyboard equivalent `Esc` is
native too. The cross is `--size-chevron`, 10px, but it is not the only way to
erase a value: there is selecting and `Backspace`.

## Sizes

All of them come from [the text field](./input.md), because `inst-search`
starts none of its own.

```html preview
<span class="inst-search">
  <input class="inst-input inst-input--sm" type="search" value="sm" aria-label="Search sm">
</span>
<span class="inst-search">
  <input class="inst-input" type="search" value="md" aria-label="Search md">
</span>
<span class="inst-search">
  <input class="inst-input inst-input--lg" type="search" value="lg" aria-label="Search lg">
</span>
```

The padding for the magnifier is counted from `--control-pad-sm` and the width
of the glyph, so at all three sizes the text starts at the same distance from
the glyph rather than running into it at `sm`.

## States

```html preview
<span class="inst-search">
  <input class="inst-input" type="search" value="Ordinary" aria-label="Search ordinary">
</span>
<span class="inst-search">
  <input class="inst-input" type="search" value="Unavailable" disabled aria-label="Search unavailable">
</span>
```

Hover, focus, error, `readonly` and `disabled` are the same as on [the text
field](./input.md).

## Behavior

### Clearing

The clear cross is the browser's own button
(`::-webkit-search-cancel-button`), brought to the tone of the library: the
same glyph and the same muted colour as the removal of [a tag](../display/tag.md).
**It is neither hidden nor replaced** — the behaviour "clear and tell the form
about it" is already written by the platform, and a button of our own would
mean JS for the sake of what already works.

:::warn
The cost is named honestly: Firefox has no native cross, and the field there is
simply without clearing. The cross is an improvement rather than the only way
to erase a value.
:::

## JS

The field works with not one line: `type="search"` gives the native clear
button, `Escape` clears it, the form is submitted on `Enter`. A script is
wanted only if the results update on the fly.

### Filtering on the fly

```js
input.addEventListener('input', () => {
  const q = input.value.trim().toLowerCase();
  for (const row of rows) {
    row.hidden = q !== '' && !row.dataset.name.includes(q);
  }
  count.textContent = rows.filter((r) => !r.hidden).length;
});
```

`input` is read rather than `keyup`: clearing by the cross, pasting with the
mouse and autofill produce no `keyup`, and the list would silently stay as it
was.

### What to say aloud

The number found has to reach a live region — otherwise for a screen reader the
list changes in silence.

```html
<p class="inst-u-visually-hidden" aria-live="polite">
  <span id="count">12</span> matches
</p>
```

## Composition

### In a toolbar

```html preview context
<div class="inst-toolbar">
  <span class="inst-search">
    <input class="inst-input inst-input--sm" type="search" placeholder="Filter" aria-label="Filter the tasks">
  </span>
  <button class="inst-btn inst-btn--sm" type="button">Reset</button>
</div>
```

### With a hint about the syntax

```html preview
<div class="inst-field">
  <label class="inst-label" for="q2">Search the log</label>
  <span class="inst-search">
    <input class="inst-input" id="q2" type="search" aria-describedby="q2-hint">
  </span>
  <span class="inst-field-hint" id="q2-hint">Quotes and a minus for exclusion are supported</span>
</div>
```

## API

```api
```

## Related

```related
```
