---
title: Notification
group: Feedback
layout: component
source: src/overlay.css
js: "The queue, the timer and the pause under the cursor are on instrument.js. The look, the appearing and the leaving are entirely CSS"
api:
  - { name: "inst-toasts", kind: "class", doc: "The region. One per document, declared `popover=\"manual\"`" }
  - { name: "inst-toast", kind: "class", doc: "One notification" }
  - { name: "inst-toast-body", kind: "class", doc: "The text part" }
  - { name: "inst-toast-title", kind: "class", doc: "What happened. One line" }
  - { name: "inst-toast-text", kind: "class", doc: "The detail. Optional" }
  - { name: "inst-toast-actions", kind: "class", doc: "The action: \"Undo\", \"Show\"" }
  - { name: "inst-toasts--start", kind: "modifier", doc: "The region at the top edge. The default is the bottom" }
  - { name: "data-tone", kind: "attribute", value: "neutral · running · ok · warn · error", doc: "The icon and the colour. With no attribute there is no icon" }
  - { name: "data-state", kind: "attribute", value: "leaving", doc: "Set by `instrument.js` before a removal, so that the transition plays out" }
  - { name: "role", kind: "attribute", value: "status · alert", doc: "`alert` on an error alone: it interrupts, the rest reports politely" }
  - { name: "popover", kind: "attribute", value: "manual", doc: "On the region. The top layer, with no light dismiss" }
  - { name: "--surface-overlay", kind: "token", doc: "The surface of a notification" }
  - { name: "--shadow-popover", kind: "token", doc: "The third and last carrier of a shadow in the library" }
---

The result of an action a person started, when the result has no place on the
screen they are looking at. "The run is queued", "the settings are saved",
"could not send".

```html preview
<div class="inst-toasts" style="position:static;inset:auto;pointer-events:auto">
  <div class="inst-toast" data-tone="ok" role="status">
    <div class="inst-toast-body">
      <div class="inst-toast-title">The run is queued</div>
      <div class="inst-toast-text">worldbox-1 · seventh in the queue</div>
    </div>
    <div class="inst-toast-actions">
      <button class="inst-btn inst-btn--sm inst-btn--ghost" type="button">Show</button>
    </div>
  </div>

  <div class="inst-toast" data-tone="error" role="alert">
    <div class="inst-toast-body">
      <div class="inst-toast-title">Could not send</div>
      <div class="inst-toast-text">The network is unavailable. Attempt 3 of 5.</div>
    </div>
  </div>

  <div class="inst-toast">
    <div class="inst-toast-body">
      <div class="inst-toast-title">The settings are saved</div>
    </div>
  </div>
</div>
```

:::note
The inline style in the example is the only one in the whole reference, and it
is forced here: the real region lives in the **top layer** and is pushed into
the corner of the screen, while it has to be shown inside a frame. In an
application there is no need to write it this way.
:::

## Contract

The markup is built by `instrument.js`: a toast is started by a call rather
than by a node written by hand. If the standard one is not enough, assemble
your own and put it into `.inst-toasts` yourself.

| What | Required | Why |
|---|---|---|
| `role="status"`, and `role="alert"` on an error | yes | An error interrupts, the rest reports politely |
| `duration: 0` on an error | yes | A message that an action was **not** carried out has no right to leave unnoticed |
| One action rather than two | yes | A toast leaves by itself; a choice between two options takes time it does not have |
| One region per document | yes | `popover="manual"`, open from initialisation |

### Accessibility

| | |
|---|---|
| The role | `status` on an ordinary one, `alert` on an error. An error interrupts, the rest reports politely |
| The live region | The region exists and is **open** before the first toast appears. A region created at the same time as its content is not spoken |
| The pause | The timer freezes under the cursor and on focus — WCAG 2.2.1 |
| An error | `duration: 0`: a message about an action not carried out does not leave by itself |
| Not colour alone | The tone carries an **icon** rather than a fill alone — the same set of masks as [the note](./note.md) and [the banner](./banner.md) |
| Reduced motion | The transition collapses; the removal is backed by a timer, otherwise the toast would stay for good |
| The keyboard | The action inside is an ordinary button and is reachable by `Tab`. The region does not trap the focus: it reports rather than asks |
| Print | Notifications are not printed on paper: they are temporary by definition |
| A modal | A toast is not visible under one. Report the result inside the dialog or after it closes — see above |

## Behavior

### Why a popover

The region is declared a `popover`, because `z-index` does not settle the task
entire. A toast has to be visible over **everything the application drew**: a
panel with `overflow: hidden`, a sticky header, an open popover. The top layer
gives that by construction, and a number in a `z-index` only until the first
larger number of somebody else's.

Hence there is **no separate z-index for toasts** in the library: there would
be nowhere to apply it.

### A modal is above a toast

**A toast does not cover an open modal.** Verified by measurement: a region
reopened after `showModal()` — in the same task and a task later alike — stays
under the dialog. The top layer is ordered by the time of entry, but a modal
dialog holds above manipulations with `popover`.

What to do instead:

| The situation | How |
|---|---|
| The action was started **from a modal** and it is closing | Close the dialog, then show the toast |
| The modal stays open | The result goes **inside it**: [a note](./note.md) or [a failure block](../../agent/failure.md) in the body of the dialog |
| A long operation with a modal open | [A meter](../charts/meter.md) inside the dialog, and a toast after it closes |

That is more correct anyway: a message about the result of an action started in
a dialog ought to arrive where the person is looking.

### Why the region is always open

`popover="manual"` is opened once at initialisation and is not closed.

A closed `popover` is `display: none`, and a live region in `display: none` is
**not spoken**: a screen reader would report not one notification. An empty
open region is invisible and intercepts nothing — `pointer-events` come back
only to the toasts themselves.

### The pause and the ceiling of the queue

The timer freezes under the cursor and on focus. A notification that vanished
just as it began to be read is a lost message, and WCAG 2.2.1 requires this:
time has to have a pause.

More than four at once are not shown — twenty notifications in a row read as a
wall. The oldest leave so that a new one is visible.

## JS

A toast is not written as markup — it is called. So there is one example here
and it is live: press it, and a notification arrives in the bottom right corner
of the screen.

```html preview
<button class="inst-btn inst-btn--primary" type="button"
        data-demo-toast='{"tone":"ok","title":"The run is queued","text":"worldbox-1 · seventh in the queue"}'>Success</button>
<button class="inst-btn" type="button"
        data-demo-toast='{"tone":"error","title":"Could not send","duration":0}'>An error that does not leave</button>
<button class="inst-btn" type="button"
        data-demo-toast='{"tone":"running","title":"The build is going","text":"Two steps left"}'>Running</button>
```

:::note
`data-demo-toast` is an attribute of **this site** rather than of the library:
it is wanted so that the page stays markup. In an application you call
`toast()` from code.
:::

```js
import { toast } from '@keshon/instrument/js';

toast({ tone: 'ok', title: 'The run is queued',
        text: 'worldbox-1 · seventh in the queue' });
```

With no script there are no toasts at all: the queue, the timer and the pause
under the cursor are work CSS does not do. The look, the appearing and the
leaving are entirely on CSS.

### Methods

| | |
|---|---|
| `toast(options)` | Show a notification. It returns the node — it can be taken away early |

### Options

| Field | What it does |
|---|---|
| `title` | What happened. One line, with no full stop |
| `text` | The detail. Optional |
| `tone` | `ok` · `warn` · `error` · `running` · `neutral` |
| `duration` | Milliseconds. `0` — it does not leave by itself |
| `action` | `{ label, onClick }` — one action rather than two |

```js
toast({ tone: 'error', title: 'Could not send', duration: 0 });

toast({ title: 'The task is deleted',
        action: { label: 'Undo', onClick: () => restore() } });
```

### Events

There are none of its own. The action arrives as the `action.onClick` callback.

## API

```api
```

## Related

```related
```
