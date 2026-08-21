---
title: Text field
group: Inputs
layout: component
source: src/forms.css
api:
  - { name: "inst-input", kind: "class", doc: "A single-line field" }
  - { name: "inst-textarea", kind: "class", doc: "A multi-line one" }
  - { name: "inst-field", kind: "class", doc: "The wrapper: label + control + hint" }
  - { name: "inst-label", kind: "class", doc: "The label" }
  - { name: "inst-field-hint", kind: "class", doc: "The hint" }
  - { name: "inst-field-error", kind: "class", doc: "The error message" }
  - { name: "inst-input-group", kind: "class", doc: "A group with affixes" }
  - { name: "inst-input-affix", kind: "class", doc: "A prefix or a suffix" }
  - { name: "inst-required", kind: "class", doc: "The mark of being required" }
  - { name: "inst-optional", kind: "class", doc: "The mark of being optional" }
  - { name: "inst-input--sm", kind: "modifier", doc: "The size" }
  - { name: "inst-input--lg", kind: "modifier", doc: "The size" }
  - { name: "inst-textarea--sm", kind: "modifier", doc: "The size" }
  - { name: "inst-textarea--lg", kind: "modifier", doc: "The size" }
  - { name: "inst-field--side", kind: "modifier", doc: "The label at the side instead of above" }
  - { name: "aria-invalid", kind: "attribute", doc: "`true` — an error from the application" }
  - { name: "data-when", kind: "attribute", doc: "`invalid` — the showing is led by the platform" }
  - { name: "--surface-field", kind: "token" }
  - { name: "--control-h-sm/md/lg", kind: "token" }
  - { name: "--control-pad-sm", kind: "token" }
  - { name: "--radius-sm", kind: "token" }
  - { name: "--text-xs/sm/md", kind: "token" }
  - { name: "--label-col", kind: "token" }
  - { name: "--border-control", kind: "token" }
  - { name: "--err-text", kind: "token" }
  - { name: "--leading-ui", kind: "token" }
---

Single-line and multi-line input. A field stands in a recess
(`--surface-field`) rather than on the colour of the panel: otherwise its only
boundary is the border, and the field has to be searched for.

```html preview
<div class="inst-field">
  <label class="inst-label" for="agent">Agent name</label>
  <input class="inst-input" id="agent" placeholder="worldgen-01">
  <span class="inst-field-hint">Latin letters, digits and a hyphen</span>
</div>
```

## Contract

The field, the label and the hint live in one `inst-field` wrapper.

| What | Required | Why |
|---|---|---|
| A `<label for>` or an `aria-label` | yes | A `placeholder` is **not** a label: it disappears on input and is not read by some screen readers |
| `inst-input` on the field | yes | The recess, the border, the height and every state come from there |
| `inst-field` on the wrapper | no | Wanted where there is a label or a hint. A bare field in a toolbar does without it |
| An `aria-describedby` on the field | no, but usually yes | It ties the hint to the field: without it the hint is seen by the eye and not read by a screen reader |

A multi-line field is the same contract with the class `inst-textarea`.

A `<label for>` is required. A `placeholder` is **not** a label: it disappears
on input and is not read by some screen readers. The hint is hung on through
`aria-describedby` on the `inst-field-hint`. An error message is text rather
than a red border alone: colour has no right to be the only carrier. The same
with being required — `required` plus a visible mark, and an asterisk with no
explanation in words at the start of the form is half the contract.

The border of a field takes `--border-control` and holds 3:1: here the boundary
**is** the control.

The order of the states is written into the selector: disabled > error > focus
> hover. `:hover:not(:disabled)` has a specificity of (0,3,0) and
`:user-invalid` only (0,2,0) — a cursor over an invalid field hid the red
border until hover began to yield explicitly.

## Sizes

They are symmetrical with the button's, because a field and a button stand in
one row of a form.

```html preview
<input class="inst-input inst-input--sm" value="sm" aria-label="Field sm">
<input class="inst-input" value="md" aria-label="Field md">
<input class="inst-input inst-input--lg" value="lg" aria-label="Field lg">
<textarea class="inst-textarea inst-textarea--sm" aria-label="Multi-line sm">sm</textarea>
<textarea class="inst-textarea" aria-label="Multi-line md">md</textarea>
<textarea class="inst-textarea inst-textarea--lg" aria-label="Multi-line lg">lg</textarea>
```

| Size | Height | Type size |
|---|---|---|
| `inst-input--sm` | `--control-h-sm` | `--text-xs` |
| default | `--control-h-md` | `--text-sm` |
| `inst-input--lg` | `--control-h-lg` | `--text-md` |

The horizontal padding does not change: in a field the text is pushed to the
edge rather than centred, so there is no reason for the padding to grow along
with the height.

`inst-textarea` has the same three sizes, with a height of two and a half
controls.

## States

```html preview
<input class="inst-input" value="Ordinary" aria-label="Ordinary">
<input class="inst-input" value="Read-only" readonly aria-label="Read-only">
<input class="inst-input" value="Unavailable" disabled aria-label="Unavailable">
<input class="inst-input" value="Error" aria-invalid="true" aria-label="Error">
```

| State | How it is set | What happens |
|---|---|---|
| hover | `:hover` | The border darkens. It yields to an error explicitly, see below |
| focus | `:focus-visible` | A ring plus a border in the accent |
| error | `:user-invalid` **or** `aria-invalid="true"` | A red border. The first is set by the platform, the second by the application |
| read-only | `readonly` | The recess is removed and the colour of the text is **full**: the value can be read, selected and copied |
| unavailable | `disabled` | An opacity of `0.5` — as on a button |

### `readonly` is not `disabled`

These are different things, and while they looked alike a user could not tell
whether to wait for an unlock or the value is unchangeable in principle.

| | `readonly` | `disabled` |
|---|---|---|
| Can be read and copied | yes | no |
| Takes focus | yes | no |
| Goes to the server with the form | yes | no |
| How it looks | the recess removed, the text full | an opacity of 0.5 |

## Behavior

### Validation

The error is shown by the **platform**.

```html preview
<div class="inst-field">
  <label class="inst-label" for="mail">Email<span class="inst-required">*</span></label>
  <input class="inst-input" id="mail" type="email" required
         aria-describedby="mail-msg">
  <span class="inst-field-hint" id="mail-msg">We will send a report of the run</span>
  <span class="inst-field-error" data-when="invalid">Check the address</span>
</div>
```

`:user-invalid` rather than `:invalid`. The difference is fundamental:
`:invalid` burns red from the first second, before anything has been typed into
the field — the form meets the user with a scattering of errors they have not
made. `:user-invalid` fires only after an interaction.

**The hint and the error take one slot.** The error replaces the hint rather
than shifting the form: the markup is not relaid at the moment of filling it
in.

| Who shows it | How |
|---|---|
| The platform | `data-when="invalid"` on the `inst-field-error` — it appears on `:user-invalid` |
| The application | The same block **with no attribute** — from a server's reply, for instance |

## Composition

### A prefix and a suffix

```html preview
<div class="inst-input-group">
  <span class="inst-input-affix">$</span>
  <input class="inst-input" aria-label="Budget">
  <span class="inst-input-affix">per hour</span>
</div>
```

### A label at the side

```html preview
<div class="inst-field inst-field--side">
  <label class="inst-label" for="t">Timeout</label>
  <input class="inst-input" id="t" value="30">
</div>
```

The label column is the same `--label-col` as the inspector's: two ways of
showing a name–value pair have no right to part in their metrics.

## API

```api
```

## Related

```related
```
