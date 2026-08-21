---
title: Behaviour
group: Foundations
layout: foundation
source: src/kit.js
api:
  - { name: "role", kind: "attribute", value: "listbox · tree · menu · radiogroup · tablist", doc: "On the container. Only these five have a keyboard contract" }
  - { name: "aria-orientation", kind: "attribute", value: "horizontal · vertical", doc: "The axis of the arrows. By default vertical on a list, a tree and a menu, horizontal on a radio group and tabs" }
  - { name: "aria-level", kind: "attribute", doc: "On a `treeitem`. The back arrow finds the parent by it — the nesting in the markup may be flat" }
  - { name: "aria-expanded", kind: "attribute", value: "true · false", doc: "On an expandable `treeitem`. Without it the node counts as a leaf" }
  - { name: "aria-disabled", kind: "attribute", value: "true", doc: "The item stays in the arrow traversal but does not act" }
  - { name: "data-value", kind: "attribute", doc: "On an item of a group and on a tag. What arrives in the event's `detail` instead of the text" }
  - { name: "data-copy", kind: "attribute", doc: "On `.inst-copy`. Copy this value rather than the text of the block" }
  - { name: "data-copied-label", kind: "attribute", doc: "On `.inst-copy`. A phrase of your own for a screen reader instead of «Copied»" }
  - { name: "data-failed-label", kind: "attribute", doc: "On `.inst-copy`. The same for a failure" }
  - { name: "data-copied", kind: "attribute", value: "true · false", doc: "Set and removed by `instrument.js`. The response of the copy button, it lives 1.4 s" }
  - { name: "data-instrument", kind: "attribute", value: "manual", doc: "On `<html>`. The module does not start itself — the application calls `start()` when it is ready" }
  - { name: "data-roving", kind: "attribute", value: "manual", doc: "On a group. The roving `tabindex` inside it is led by the application; the arrows and the selection work as before" }
  - { name: "inst:select", kind: "event", doc: "An item of a group was selected. `detail` is `{ value }`. Cancelling leaves the markup untouched" }
  - { name: "inst:copy", kind: "event", doc: "The copy button was pressed. `detail` is `{ text }`" }
  - { name: "inst:remove", kind: "event", doc: "A tag is being removed. `detail` is `{ value }`" }
  - { name: "inst:selectall", kind: "event", doc: "The checkbox in a table head was toggled. `detail` is `{ checked }`" }
---

CSS declares the roles, `instrument.js` carries out their keyboard contract. A
`role="listbox"` tells assistive technology that the arrows work — without the
script that statement is false.

The file deals with the keyboard and with attributes only; it does not touch
the styling.

The example below is live: enter it with `Tab` and press an arrow. One `Tab`
for the whole group, arrows after that; `Home` and `End` work, and the
traversal wraps around. There is not one handler in the markup.

```html preview
<div class="inst-segmented" role="radiogroup" aria-label="Density">
  <button type="button" role="radio" aria-checked="true" tabindex="0">Compact</button>
  <button type="button" role="radio" aria-checked="false" tabindex="-1">Default</button>
  <button type="button" role="radio" aria-checked="false" tabindex="-1">Comfortable</button>
</div>

<div class="inst-tabs" role="tablist" aria-label="Sections">
  <button class="inst-tab" type="button" role="tab" aria-selected="true" tabindex="0">Active</button>
  <button class="inst-tab" type="button" role="tab" aria-selected="false" tabindex="-1">History</button>
  <button class="inst-tab" type="button" role="tab" aria-selected="false" tabindex="-1">Schedule</button>
</div>
```

## Contract

```html
<link rel="stylesheet" href="instrument.min.css">
<script type="module" src="instrument.js"></script>
```

A module rather than a classic script: it has named exports. No build step is
required for that.

### Accessibility

| | |
|---|---|
| Disabled items | `aria-disabled="true"` **keeps** an item in the arrow traversal and in the roving `tabindex` but removes the action: neither a click nor `Enter` selects it or calls the application. Somebody moving by arrows has to learn that the action exists and why it is unavailable now — an item thrown out is indistinguishable from one that never was. `disabled` is a different matter: the platform removes it from the traversal itself |
| Invisible items | Thrown out too: neither the platform nor `instrument.js` puts focus on them |
| Nested groups | An arrow in a parent does not jump across somebody else's items — submenus and groups inside a tree are counted separately |
| Modifiers | `Ctrl`, `Alt` and `Cmd` with arrows are not intercepted: those are commands of the platform |
| Already handled | An event with `defaultPrevented` is skipped, so an application can intercept a key earlier |

## Anatomy

| Role | Arrows | And also |
|---|---|---|
| `listbox` | ↑ ↓ over rows | The selection follows the focus |
| `tree` | ↑ ↓ over nodes, → ← expansion | ← on a collapsed node goes to the parent by `aria-level` |
| `menu` | ↑ ↓ over items | There is no selection: an item is an action rather than a choice |
| `radiogroup` | ← → over options | The check follows the focus |
| `tablist` | ← → over tabs | The selection follows the focus |

`Home` and `End` are in all five. The traversal **wraps around**: a task queue
has no "end" past which there is nowhere to go, and a dead end on the last row
gives nothing.

## Behavior

A roving `tabindex`: exactly one element of a group is reachable by `Tab` — one
`Tab` in, arrows inside, one `Tab` out.

The anchor is the selected element, and if there is no selection, the first
one.

### Rows arriving while the page works

The listeners hang on the document, and the elements are found **at the moment
of the keypress**. Queue rows, steps and tree nodes appear one at a time, and
an `init(el)` at start-up would miss everything that came later.

The one thing that cannot be done lazily is the initial `tabindex`: a group
without it is not reachable by `Tab`. What watches for groups appearing is a
`MutationObserver` rather than a human. Every root has an observer of its own:
`start()` on a second root does not take the observer away from the first, and
`stop()` removes exactly the one it started.

A group appearing is not always an edit to the markup. A popover and a
`<details>` open by an attribute, and an observer does not wake for that, so
the module also listens for `toggle`: menu items get their roving `tabindex` at
the moment of opening rather than at the moment focus first lands in them.

```js
import { refresh } from '@keshon/instrument/js';

// Needed only when the observer is off — in a test, for instance.
refresh(container);
```

### Who owns the markup

Everything the module writes it writes through a cancellable event:
`preventDefault` leaves the markup untouched, and the focus goes where it was
going anyway — the application refused to lead the state, it did not forbid a
person to move.

There is one exception, the roving `tabindex`. An event will not do for it: it
is laid out from `refresh()`, that is on every batch of changes to the
document, and the answer there is always the same. So the ownership is declared
rather than decided each time.

Three knobs, and all three exist for one thing: so that an attribute never has
two owners.

| Knob | Where | What it does |
|---|---|---|
| `data-instrument="manual"` | on `<html>` | The module does not start itself. Read before the first subscription: by the time of the import it is too late to cancel |
| `start(root, { observe: false })` | in code | Do not raise a `MutationObserver`. A framework knows about its own commit anyway, and a second observer means a full walk of the document on every render |
| `data-roving="manual"` | on a group | The `tabindex` inside it is led by the application. The arrows, `Home`, `End` and the selection keep working |

Without these three words everything stays as it was: a `<script type="module">`
and a side-effect import start the module themselves.

```html
<html data-instrument="manual">
```

```js
import { start, refresh } from '@keshon/instrument/js';

start(document, { observe: false });
// and refresh() after every render
```

### The lines a person hears

The module speaks four phrases — the name of the notification region and three
responses to actions. **The defaults are English**: the kit travels in a
package, and the language of a screen is chosen by the application rather than
by the library. All four are gathered in one place and change together:

```js
import { strings } from '@keshon/instrument/js';

strings.copied = 'Copié';
```

The attributes `data-copied-label` and `data-failed-label` still win: they have
another job — to say something of their own for one button rather than to
change the language of the screen.

### Drawn controls

The copy button, the removal cross, the `ew-resize` cursor on the label of an
axis are controls declared by the styling. Their work is done by the same
module.

| What is drawn | What `instrument.js` does |
|---|---|
| `.inst-copy` inside `.inst-copyable` or `.inst-code` | Puts the text in the clipboard, answers with colour and announces the result |
| `.inst-tag-remove` | Removes [the tag](../components/display/tag.md) and moves the focus to its neighbour |
| `.inst-num-axis` with an `ew-resize` cursor | Drags the value of [a number field](../components/inputs/num-field.md) |
| An `<output for>` beside [a slider](../components/inputs/slider.md) | Keeps the current value in it |
| The checkbox in the head of a selection column of [a table](../components/display/table.md) | Selects every row, shows a partial selection |
| A selected tab, option, list row | Carries `aria-selected` and `aria-checked` over — on an arrow and on a click alike |

Sorting, uploading files, picking a date and closing a banner are **not** here:
those actions work with data, that is, they stay on the application's side.

## JS

The module exports nothing as classes and needs no initialisation: it is
included once and works by delegation. Everything listed below is the surface
for when that is not enough.

### Methods

| | |
|---|---|
| `start(root = document)` | Attach. Called by itself when the module loads. There may be several roots — each with an observer of its own |
| `stop(root = document)` | Detach. Wanted by tests and hot reloading rather than by an application |
| `refresh(root = document)` | Walk a subtree again: the roving `tabindex`, the `<output>` of sliders, "select all" in tables |
| `toast(options)` | Show [a notification](../components/feedback/toast.md) |

### Events

All of them are `CustomEvent`, they bubble and they are **cancellable**.
`preventDefault()` means the application takes the work on itself:
`instrument.js` stops and does not touch the markup.

| Event | Where | `detail` |
|---|---|---|
| `inst:select` | on an item of a group with a role | `{ value }` — the item's `data-value`, otherwise its text |
| `inst:copy` | on `.inst-copy` | `{ text }` — what will go to the clipboard |
| `inst:remove` | on `.inst-tag` | `{ value }` — the tag's `data-value`, otherwise its text |
| `inst:selectall` | on `.inst-table` | `{ checked }` — the new state of every row |
| `inst:expand` | on a node of [a tree](../agent/tree.md) | `{ open }` — the new state of `aria-expanded` |

```js
document.addEventListener('inst:remove', (e) => {
  e.preventDefault();            // the tag will stay in the markup
  store.dropTag(e.detail.value); // a re-render from the data will do it
});
```

A cancelled `inst:select` leaves `aria-selected` and `aria-checked` as they
were, but the focus still moves: the application refused to lead the state
itself, it did not forbid a person to move.

Dragging an axis and the slider start no events of their own: they change the
`<input>` and send the native `input` and `change`. A framework sees them with
no glue.

### Options

The settings are attributes of the markup rather than a configuration object:
the library has no instances to attach one to.

| Attribute | Where | What it does |
|---|---|---|
| `data-copy` | on `.inst-copy` | Copy this rather than the text of the block |
| `data-copied-label` | on `.inst-copy` | A phrase of your own for a screen reader instead of «Copied» |
| `data-failed-label` | on `.inst-copy` | The same for a failure |
| `data-value` | on `.inst-tag` and on an item of a group | What arrives in `detail` instead of the text |
| `step`, `min`, `max` | on an `input[type=number]` | The step and the bounds of dragging an axis |

## Rules

**It does not change the styling.** Not one assignment to `style`, not one
class of styling: the module sets attributes from the markup contract —
`tabindex`, `aria-selected`, `aria-checked`, `aria-expanded` — and how they
look is decided by CSS. So an application can override the look without
touching the behaviour.

**It is not obligatory.** Without it everything that does not need a keyboard
works: expansion on `<details>`, [the popover](../components/overlays/popover.md)
on the Popover API, [the modal](../components/overlays/dialog.md) on
`<dialog>`, validation on `:user-invalid`.

## API

```api
```

### What is not there yet

- **typeahead**, selection by first letters. In a list of two hundred rows it is
  wanted, and APG recommends it;
- **multiple selection** with `Shift`+arrows in a `listbox`;
- **closing a tooltip on `Escape`**: as it stands it does not meet criterion
  1.4.13, because it can be neither hovered nor dismissed.

## Related

[Queue](../agent/task.md)
[Tree](../agent/tree.md)
[Menu](../components/overlays/menu.md)
[Segmented control](../components/actions/segmented.md)
[Tabs](../components/navigation/tabs.md)
[Design principles](../about/design-principles.md)
