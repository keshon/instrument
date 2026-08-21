---
title: Modal
group: Overlays
layout: component
source: src/overlay.css
js: "Opening is one line: dlg.showModal(). Closing, the ground, blocking the scroll of the background and returning the focus call for no script"
api:
  - { name: "inst-dialog", kind: "class", doc: "The base one. Put on a `<dialog>`" }
  - { name: "inst-dialog-head", kind: "class", doc: "The header" }
  - { name: "inst-dialog-title", kind: "class", doc: "The heading" }
  - { name: "inst-dialog-close", kind: "class", doc: "The cross at the far edge of the header" }
  - { name: "inst-dialog-body", kind: "class", doc: "The scrolling body" }
  - { name: "inst-dialog-foot", kind: "class", doc: "The footer with the actions" }
  - { name: "inst-dialog-foot-note", kind: "class", doc: "An explanation pushing the buttons to the far edge" }
  - { name: "inst-sheet", kind: "class", doc: "Another layout of the same `<dialog>` — [the sheet](./sheet.md)" }
  - { name: "inst-dialog-sub", kind: "class", doc: "A line under the heading: what the dialog wants, if that does not follow from the heading entire" }
  - { name: "inst-dialog-foot--start", kind: "modifier", doc: "The actions at the near edge. The default is the far one" }
  - { name: "inst-dialog--scroll", kind: "modifier", doc: "The body scrolls: it brings back the rules under the header and above the footer" }
  - { name: "--surface-overlay", kind: "token" }
  - { name: "--shadow-modal", kind: "token" }
  - { name: "--scrim", kind: "token" }
  - { name: "--border", kind: "token" }
  - { name: "--border-subtle", kind: "token" }
  - { name: "--hairline", kind: "token" }
  - { name: "--radius-lg", kind: "token" }
  - { name: "--pad-card", kind: "token" }
  - { name: "--space-3", kind: "token" }
  - { name: "--space-5", kind: "token" }
  - { name: "--space-8", kind: "token" }
  - { name: "--gap-inline", kind: "token" }
  - { name: "--text-md", kind: "token" }
  - { name: "--text-xs", kind: "token" }
  - { name: "--weight-medium", kind: "token" }
  - { name: "--text-muted", kind: "token" }
---

A window that takes the focus and dims the background, because there is no
going on without an answer. A native `<dialog>` is styled, so the modality is
real rather than drawn.

```html preview
<button class="inst-btn inst-btn--danger" type="button"
        onclick="document.getElementById('dlg-confirm').showModal()">Delete the run</button>

<dialog class="inst-dialog" id="dlg-confirm">
  <form method="dialog">
    <div class="inst-dialog-head">
      <span class="inst-dialog-title">Delete run #4127?</span>
      <span class="inst-dialog-sub">The artefacts go with it</span>
      <button class="inst-btn inst-btn--sm inst-btn--ghost inst-btn--icon inst-dialog-close" type="submit" aria-label="Close">
        <svg class="inst-icon" aria-hidden="true"><use href="#i-close"/></svg>
      </button>
    </div>
    <div class="inst-dialog-body inst-stack">
      <p>The log, the artefacts and the diff will be deleted. Links to the run in reports will stop opening.</p>
      <div class="inst-note" data-tone="warn">Irreversible: restoring from a backup takes up to a day.</div>
    </div>
    <div class="inst-dialog-foot">
      <span class="inst-dialog-foot-note">Escape closes it too</span>
      <button class="inst-btn" type="submit">Cancel</button>
      <button class="inst-btn inst-btn--danger" type="submit">Delete</button>
    </div>
  </form>
</dialog>
```

## Contract

The whole content is wrapped in a `<form method="dialog">`. Opening is one line
of script; closing, the ground and returning the focus call for no script.

:::note
The `onclick` in the example above is a shorthand for the sake of one file. In
an application the handler is attached in code; how it is attached is no
concern of the library's. The declarative `command` and `commandfor` are
progressive for now and are not taken into the library's contract.
:::

### Accessibility

| | |
|---|---|
| Focus | `showModal()` takes the focus inside, holds it there and returns it to the button that called it on closing. None of that has to be written |
| The keyboard | `Escape` closes. `Tab` does not leave the modal |
| The background is inert | The content under the ground is available neither to the mouse nor to a screen reader — the platform does that |
| Scrolling the background | It stops by the rule `html:has(dialog:modal)`. Otherwise the user loses their place |
| The heading | `inst-dialog-title` styles but does not declare. If the modal wants an accessible name, tie it to the heading by `aria-labelledby` |
| The ground | `--scrim` — 0.32 in a light theme and 0.58 in a dark one: the background reads as switched off without disappearing |
| Print | A modal is not printed on paper |
| The required markup | A native `<dialog>`, `showModal()` rather than `show()`, a `<form method="dialog">` inside, a `type="submit"` on the closing buttons, an `aria-label` on the cross |

## Anatomy

| The part | The work |
|---|---|
| `inst-dialog-head` | The header: the heading and the cross. Parted by a hairline |
| `inst-dialog-title` | The heading, `--text-md` |
| `inst-dialog-close` | The cross. Pushed to the far edge of the header |
| `inst-dialog-body` | The body. The only scrolling part |
| `inst-dialog-foot` | The footer with the actions |
| `inst-dialog-foot--start` | The actions at the near edge. The default is the far one |
| `inst-dialog-foot-note` | An explanation at the near edge, pushing the buttons to the far one |

The header and the footer do not shrink (`flex: none`), and only the body
scrolls. Because of that the heading and the buttons stay in place at any
length of content, and the modal does not grow taller than `80dvh`.

| Property | Value |
|---|---|
| width | `min(34rem, 100vw - var(--space-8))` |
| maximum height | `min(80dvh, 100dvh - var(--space-8))` |

The layout is set for the open one alone: `.inst-dialog[open] { display: flex }`.
An unconditional `display: flex` would show the modal always — the closed state
of a `<dialog>` is `display: none`.

A modal is the **second and last carrier of a shadow in the library**. The
first is [the popover](./popover.md).

## Behavior

### Closing with no script

Any `type="submit"` button inside a `<form method="dialog">` closes the modal
and passes its value to `dialog.returnValue` — **not one of them needs a click
handler**, the cross and "Cancel" included.

| The way of closing | What is wanted |
|---|---|
| A button inside a `<form method="dialog">` | A `type="submit"` |
| `Escape` | Nothing. The platform |
| A click on the ground | A `closedby="any"` on the `<dialog>` |

### Blocking the scroll of the background

`showModal()` makes the background inert to clicks but **does not stop the
scrolling**: the wheel keeps turning the page, and the user loses the place
they were coming back to. The library closes that with one rule:

```css
html:has(dialog:modal) { overflow: hidden; }
```

The selector is `:modal` rather than `[open]`, and the difference here is
fundamental. The `open` attribute is set by a non-modal `show()` too, whose
entire job is **not** to block the page. While the rule looked at the
attribute, a non-modal panel froze the scroll of the whole document — and did
so with any `<dialog>` on the page, even one with not a single class of the
library.

What stops is the user's scrolling — the wheel, the touchpad, the keys. A
programmatic `scrollTo` keeps working, and that is right: an application may
want to bring the background to the place wanted while a modal is open.

The price is known: if what scrolls is the **document** itself, the scrollbar
disappears at the moment of opening and the content shifts sideways by its
width. `scrollbar-gutter` does not save it here — the reservation of room works
at `overflow: auto/scroll`, and here it is `hidden`.

```css
/* If what scrolls in your case is the document */
html { scrollbar-gutter: stable; }
```

In an application built on [the shell](../../layout/shell.md) the document does
not scroll at all, so there is no shift and nothing to pay. The library does
not set `scrollbar-gutter` itself: on a screen that does not scroll it would be
an empty 10px strip at the right edge.

## JS

Opening is one line. Everything else is taken on by the platform.

```js
document.getElementById('dlg').showModal();
```

| What | Who does it |
|---|---|
| Opening | The application: `showModal()` |
| Closing | A `<form method="dialog">` inside — with no script |
| The ground, the top layer, `Escape` | The platform |
| Returning the focus to the button that opened it | The platform |
| Blocking the scroll of the background | The library, by the rule `html:has(dialog:modal)` |

:::note
The declarative `command` and `commandfor` would spare even that line, but they
are not supported everywhere yet and are not taken into the library's contract.
:::

## API

```api
```

## Related

```related
```
