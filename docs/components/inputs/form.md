---
title: Form
group: Inputs
layout: component
source: src/forms.css
api:
  - { name: "inst-form", kind: "class", doc: "The form: a vertical flow with a step of `--pad-panel`" }
  - { name: "inst-field", kind: "class", doc: "One field: label + control + hint" }
  - { name: "inst-form-actions", kind: "class", doc: "The action bar" }
  - { name: "inst-form-actions-note", kind: "class", doc: "An explanation in the bar" }
  - { name: "inst-fieldset", kind: "class", doc: "A group of fields" }
  - { name: "inst-required", kind: "class", doc: "The marks of being required" }
  - { name: "inst-optional", kind: "class", doc: "The marks of being optional" }
  - { name: "inst-form--side", kind: "modifier", doc: "Every label at the side" }
  - { name: "inst-field--side", kind: "modifier", doc: "The label at the side on one field" }
  - { name: "inst-form-actions--end", kind: "modifier", doc: "The actions to the end of the line" }
  - { name: "inst-fieldset--framed", kind: "modifier", doc: "With a border" }
  - { name: "--pad-panel", kind: "token" }
  - { name: "--label-col", kind: "token" }
  - { name: "--gap-inline", kind: "token" }
  - { name: "--space-2", kind: "token" }
  - { name: "--text-xs", kind: "token" }
  - { name: "--err-text", kind: "token" }
  - { name: "--text-muted", kind: "token" }
  - { name: "--z-sticky", kind: "token" }
---

The layout of fields and the action bar. Two layouts; the choice between them
is settled by the nature of the values.

```html preview
<form class="inst-form">
  <div class="inst-field">
    <label class="inst-label" for="f1">Run name</label>
    <input class="inst-input" id="f1" placeholder="worldgen-01">
  </div>
  <div class="inst-field">
    <label class="inst-label" for="f2">Model<span class="inst-optional">optional</span></label>
    <span class="inst-select-wrap">
      <select class="inst-select" id="f2"><option>opus</option><option>sonnet</option></select>
    </span>
  </div>
  <div class="inst-form-actions">
    <button class="inst-btn inst-btn--primary" type="submit">Run</button>
    <button class="inst-btn" type="button">Cancel</button>
  </div>
</form>
```

## Contract

A real `<form>`, the fields in `inst-field`, the actions in
`inst-form-actions`.

| What | Required | Why |
|---|---|---|
| A real `<form>` | yes | `Enter` in a field submits the form. A button of your own on a `div` breaks that |
| A `<label for>` on every control | yes | A form with not one `for` passes a visual check and fails any other |
| `type="submit"` on the principal action | yes | Otherwise no submission happens, and `Enter` does not work |
| A `<fieldset>` with a `<legend>` around a group | no, but usually yes | For radio buttons it is the only way to speak the question rather than the options alone |

### Being required

```html preview
<div class="inst-field">
  <label class="inst-label" for="r1">Name<span class="inst-required">*</span></label>
  <input class="inst-input" id="r1" required>
</div>
<div class="inst-field">
  <label class="inst-label" for="r2">Comment<span class="inst-optional">optional</span></label>
  <input class="inst-input" id="r2">
</div>
```

Being required is carried by **a sign and a word**: an asterisk alone works
only for those who already know the convention. In long forms it pays to mark
the optional ones — there are usually fewer of them.

A real `<form>` with a `type="submit"`: `Enter` in a field submits the form,
and that is expected behaviour a button of your own breaks. Every control is
tied to a `<label for>` — a form with not one `for` passes a visual check and
fails any other. The group is named by a `<fieldset>` and a `<legend>`; for
radio buttons that is the only way to speak the question rather than the
options alone.

Being required calls for two things at once: `required` on the control and a
visible mark. The attribute with no mark is seen only by a screen reader, the
mark with no attribute only by the eye.

The order of traversal matches the order in the markup. Rearranging columns
visually without rearranging the nodes breaks `Tab`.

## Anatomy

### The action bar

```html preview
<div class="inst-form-actions inst-form-actions--end">
  <span class="inst-form-actions-note">The changes will apply at once</span>
  <button class="inst-btn" type="button">Cancel</button>
  <button class="inst-btn inst-btn--primary" type="submit">Save</button>
</div>
```

Here `inst-form-actions--end`; the default with the actions at the start of the
line stands in the example of a form above.

| Class | The work |
|---|---|
| `inst-form-actions` | The action bar: a rule, air above it, the actions in a row |
| `inst-form-actions--end` | The actions pushed to the end of the line |
| `inst-form-actions-note` | An explanation beside the actions |

The bar is **pinned** to the bottom edge of the scrolling region — and pinned
minus that region's bottom padding, `--scroll-pad-end` (see
[the shell](../../layout/shell.md)). Without the subtraction it stops above the
visible bottom and leaves a slit under it, through which the passing content
shows: the bar stops being a floor and becomes a plate in the middle of
somebody else's fields.

Inside a region the bar additionally cancels its inner padding at the sides and
below — it reaches the edges and takes the surface of the panel instead of the
surface of the page. The practical conclusion: **the bar lives where the fields
it saves live.**

The principal action stands at the **outer edge of the bar**: first when the
bar is aligned to the start, last when it is pushed to the end by
`--end` or by a note. In a form with two buttons that position is the only
signal of which of them is the principal one, and it stays the same signal
whichever way the bar is aligned.

### A group of fields

```html preview
<fieldset class="inst-fieldset inst-fieldset--framed inst-fieldset--side">
  <legend>Limits of the run</legend>
  <div class="inst-field">
    <label class="inst-label" for="lim-mem">Memory</label>
    <input class="inst-input" id="lim-mem" value="4 GiB">
  </div>
  <div class="inst-field">
    <label class="inst-label" for="lim-time">Time</label>
    <input class="inst-input" id="lim-time" value="30 min">
  </div>
</fieldset>
```

`inst-fieldset` comes **with no border** by default, and that is deliberate:
the work of a fieldset is to group fields and give the group a name, and the
border comes from the browser rather than from the meaning. A form almost
always lives inside a panel, and a border on the group would give a frame
inside a frame. The boundary is set by the heading and the air.

The border comes back with the modifier `inst-fieldset--framed`, when a group
stands on its own.

## Variants

A label above or at the side — the choice is by the type of the values.

| | When |
|---|---|
| **Above** (the default) | The fields differ and call for explanation. The label is read right before the control |
| **At the side** (`inst-form--side`) | The fields are short and of one kind. A column of labels turns a form into a table that can be scanned without reading every row |

```html preview
<form class="inst-form inst-form--side">
  <div class="inst-field">
    <label class="inst-label" for="s1">Timeout</label>
    <input class="inst-input" id="s1" value="30">
  </div>
  <div class="inst-field">
    <label class="inst-label" for="s2">Attempts</label>
    <input class="inst-input" id="s2" value="3">
  </div>
</form>
```

The modifier is put on the whole form or on a single field
(`inst-field--side`). The label column is the same `--label-col` as the
inspector's: two ways of showing a name–value pair have no right to part in
their metrics.

## Patterns

### A form inside a panel

A form rarely lives on its own — it stands inside a panel, and half its
behaviour is visible only there. Scroll the body of the panel: **the action bar
sticks to the bottom edge** rather than travelling away with the fields. In a
long form a "Save" button that has to be scrolled to is a button nobody
presses.

```html preview context
<div class="inst-panel">
  <div class="inst-panel-header"><span class="inst-panel-title">Notifications</span></div>
  <div class="inst-panel-body">
    <div class="inst-form">
      <div class="inst-field">
        <span class="inst-label">Import a list of recipients</span>
        <label class="inst-file">
          <input type="file" multiple>
          Drag files in or choose them
          <span class="inst-file-hint">Up to 20 MB, the formats .json and .csv</span>
        </label>
      </div>

      <div class="inst-field">
        <label class="inst-label" for="f13">Who to notify</label>
        <span class="inst-select-wrap">
          <select class="inst-select" id="f13" multiple>
            <option selected>The engineer on duty</option>
            <option selected>The owner of the project</option>
            <option>The whole team</option>
          </select>
        </span>
        <span class="inst-field-hint">The chevron is not drawn — a list box expands nothing</span>
      </div>

      <fieldset class="inst-fieldset">
        <legend>When to notify</legend>
        <div class="inst-cluster inst-cluster--loose">
          <label class="inst-checkbox"><input type="checkbox" checked>On a failure</label>
          <label class="inst-checkbox"><input type="checkbox">On every run</label>
          <label class="inst-switch"><input type="checkbox" checked>Copy to email</label>
        </div>
      </fieldset>

    </div>

    <div class="inst-form-actions">
      <span class="inst-form-actions-note">The changes are not saved</span>
      <button class="inst-btn" type="button">Cancel</button>
      <button class="inst-btn inst-btn--primary" type="button">Save</button>
    </div>
  </div>
</div>
```

## API

```api
```

## Related

```related
```
