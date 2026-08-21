---
title: File field
group: Inputs
layout: component
source: src/forms.css
js: Dragging files in and the list of what was chosen belong to the application layer. Choosing through the dialog works with no script
api:
  - { name: "inst-file", kind: "class", doc: "The `<label>` zone: the dashes, the glyph, the states" }
  - { name: "inst-file-hint", kind: "class", doc: "The limits of format and size" }
  - { name: "--space-7", kind: "token" }
  - { name: "--pad-card", kind: "token" }
  - { name: "--stroke", kind: "token" }
  - { name: "--radius-md", kind: "token" }
  - { name: "--surface-field", kind: "token" }
  - { name: "--border-control", kind: "token" }
  - { name: "--accent-border", kind: "token" }
  - { name: "--accent-bg", kind: "token" }
  - { name: "--focus-ring", kind: "token" }
  - { name: "--size-icon", kind: "token" }
  - { name: "--text-sm", kind: "token" }
  - { name: "--text-xs", kind: "token" }
  - { name: "--space-3", kind: "token" }
---

Choosing a file and the zone it is dropped into are one component rather than
two. A dashed border means "something can be put here"; the native choose
button is hidden and replaced by the whole area of the zone.

```html preview
<div class="inst-field">
  <span class="inst-label">Import a list of recipients</span>
  <label class="inst-file">
    <input type="file" multiple>
    Drag files in or choose them
    <span class="inst-file-hint">Up to 20 MB, the formats .json and .csv</span>
  </label>
</div>
```

## Contract

| What | Required | Why |
|---|---|---|
| A `<label>` as the container | yes | A press anywhere in the zone opens the dialog. A `<div>` wrapper would leave only the hidden button clickable |
| An `<input type="file">` as a direct child | yes | The states of focus and unavailability are written through `:has(> input…)` |
| The call to action as text in the `<label>` | yes | It is the accessible name of the field. Without it a screen reader reads "file, button" |
| No `display: none` on the field | yes | The field is taken away by a clip and stays in the order of traversal |
| `inst-file-hint` | no | The limits of format and size. Write them **before** the choice rather than in an error message after |
| An `accept` on the `<input>` | no, but usually yes | The system dialog will filter by type. It is a hint rather than validation: check on the server anyway |

The field itself is not hidden by `display: none` but **taken away by a clip**
(`clip-path`) at a size of 1×1: a field hidden through `display` or
`visibility` falls out of the order of traversal, and the zone stops being
reachable from the keyboard.

The dashes are the **one place in the library where a line is not solid**, and
that is deliberate: there is no other sign of "something can be put here" as
understandable as this.

Inside the zone is a native `input[type=file]`: the system dialog, the filter
by `accept` and the part in a form come from the platform. The field lies in a
clip rather than being hidden — `display: none` would throw it out of the
traversal. `Tab` puts the focus on the zone, `Enter` and `Space` open the
dialog, and the ring is drawn around the whole zone rather than around the
invisible 1×1 field.

The name is given by the call to action inside the `<label>`; the hint about
the format is read next, because it too is inside.

Dragging is not the only way, and that is the point: by mouse, by keyboard and
through the system dialog the zone is equally reachable. Uploading **only** by
dragging is unavailable to almost everybody but a mouse. The zone itself is
recognised by the dashes and the glyph rather than by the colour of the border,
and a message about a wrong format is given as text beside it ([a note or a
banner](../feedback/banner.md)).

## States

```html preview
<label class="inst-file">
  <input type="file">
  Drag a file in or choose one
  <span class="inst-file-hint">Up to 20 MB</span>
</label>
<label class="inst-file">
  <input type="file" disabled>
  Import is switched off by the administrator
  <span class="inst-file-hint">Speak to the owner of the space</span>
</label>
```

| State | How it is set | What happens |
|---|---|---|
| hover | `:hover` | The border `--accent-border`, the background `--accent-bg` — the zone responds before the press |
| focus | `:has(> input:focus-visible)` | A ring around the whole zone |
| unavailable | `:has(> input:disabled)` | An opacity of `0.5`, the cursor `not-allowed` |

There is no "a file is over the zone" state in the library: it arises only from
a drag event, that is, from JS. The application hangs a class of its own and
paints with it the same two variables `:hover` uses.

## Behavior

:::warn
**Dragging is JS.** The library draws the **static** part: the dashed zone, the
glyph, the highlight on hover and on focus. Handling `dragover` and `drop`, the
"a file is over the zone" highlight, the list of what was chosen, the progress
of the upload and the errors are the application's code.

With not one line of script the component **stays working**: it is a `<label>`
around an `input[type=file]`, a press opens the system dialog, and the choice
goes with the form. A drop zone is an improvement over something that works
rather than a condition.
:::

For the progress of an upload take [a meter](../charts/meter.md), for a report
of an error [a banner or a note](../feedback/banner.md).

## JS

Include the module once per page — there is no need to initialise the
components one by one, `instrument.js` works by delegation and sees nodes that
arrived later.

```html
<script type="module" src="instrument.js"></script>
```

### What `instrument.js` does

Nothing. The zone draws the state of dragging but does not know what counts as
an acceptable file — and without that, accepting files would be a promise the
library cannot keep.

Without a script the zone **works**: inside it is a real `<input type="file">`,
and choosing through the dialog is available by mouse and by keyboard alike.

### What the application has to do

```js
for (const type of ['dragenter', 'dragover']) {
  zone.addEventListener(type, (e) => {
    e.preventDefault();          // without this the browser opens the file itself
    zone.dataset.over = '';
  });
}
for (const type of ['dragleave', 'drop']) {
  zone.addEventListener(type, () => delete zone.dataset.over);
}

zone.addEventListener('drop', (e) => {
  e.preventDefault();
  input.files = e.dataTransfer.files; // the same path as the dialog's
  input.dispatchEvent(new Event('change', { bubbles: true }));
});
```

What was dragged in is put **into the same `<input>`**. A separate list would
give the form two sources of files, and they would part on the very first
submission.

## Composition

### In a form with a label

```html preview
<div class="inst-field">
  <span class="inst-label">Import a list of recipients</span>
  <label class="inst-file">
    <input type="file" multiple accept=".json,.csv">
    Drag files in or choose them
    <span class="inst-file-hint">Up to 20 MB, the formats .json and .csv</span>
  </label>
</div>
```

The label here is a `<span class="inst-label">` rather than a `<label for>`:
the wrapper of the zone **is** a `<label>` itself, and a second label on the
same field would give two competing names.

## Patterns

### With a report of the upload

```html preview context
<div class="inst-field">
  <span class="inst-label">Dataset</span>
  <label class="inst-file">
    <input type="file" accept=".csv">
    Drag a file in or choose one
    <span class="inst-file-hint">CSV, up to 100 MB</span>
  </label>
  <div class="inst-meter-row"><span>terrain.csv</span><span class="inst-meter-value">62%</span></div>
  <div class="inst-meter" role="progressbar" aria-label="Uploading terrain.csv"
       aria-valuenow="62" aria-valuemin="0" aria-valuemax="100">
    <div class="inst-meter-fill" style="inline-size:62%"></div>
  </div>
</div>
```

## API

```api
```

## Related

```related
```
