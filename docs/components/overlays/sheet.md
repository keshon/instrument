---
title: Sheet
group: Overlays
layout: component
source: src/overlay.css
js: "Opening is one line: dlg.showModal(). Closing, the ground, blocking the scroll of the background and returning the focus call for no script"
api:
  - { name: "inst-sheet", kind: "class", doc: "A full-height layout at the end edge. Only together with `inst-dialog`" }
  - { name: "inst-sheet--start", kind: "modifier", doc: "The same block at the start edge" }
  - { name: "--surface-overlay", kind: "token" }
  - { name: "--shadow-modal", kind: "token" }
  - { name: "--scrim", kind: "token" }
  - { name: "--border", kind: "token" }
  - { name: "--hairline", kind: "token" }
  - { name: "--pad-card", kind: "token" }
  - { name: "--space-5", kind: "token" }
  - { name: "--gap-inline", kind: "token" }
  - { name: "--text-md", kind: "token" }
---

A full-height panel at the edge of the screen. The same native `<dialog>` as
[the modal](./dialog.md) — a different layout rather than a different
component.

```html preview
<button class="inst-btn" type="button"
        onclick="document.getElementById('dlg-params').showModal()">From the end edge</button>
<button class="inst-btn" type="button"
        onclick="document.getElementById('dlg-nav').showModal()">From the start edge</button>

<dialog class="inst-dialog inst-sheet inst-sheet--start" id="dlg-nav">
  <form method="dialog">
    <div class="inst-dialog-head">
      <span class="inst-dialog-title">Sections</span>
      <button class="inst-btn inst-btn--sm inst-btn--ghost inst-btn--icon inst-dialog-close" type="submit" aria-label="Close">
        <svg class="inst-icon" aria-hidden="true"><use href="#i-close"/></svg>
      </button>
    </div>
    <div class="inst-dialog-body">
      <nav class="inst-nav" aria-label="Sections">
        <a class="inst-nav-item" href="#" aria-current="page">Runs</a>
        <a class="inst-nav-item" href="#">Artefacts</a>
        <a class="inst-nav-item" href="#">Settings</a>
      </nav>
    </div>
  </form>
</dialog>

<dialog class="inst-dialog inst-sheet" id="dlg-params">
  <form method="dialog">
    <div class="inst-dialog-head">
      <span class="inst-dialog-title">Parameters of the run</span>
      <button class="inst-btn inst-btn--sm inst-btn--ghost inst-btn--icon inst-dialog-close" type="submit" aria-label="Close">
        <svg class="inst-icon" aria-hidden="true"><use href="#i-close"/></svg>
      </button>
    </div>
    <div class="inst-dialog-body">
      <dl class="inst-kv">
        <dt>Model</dt><dd>opus</dd>
        <dt>Agents</dt><dd>7 of 12</dd>
        <dt>Tokens</dt><dd>184 320</dd>
        <dt>Started</dt><dd>19:38:04</dd>
      </dl>
    </div>
    <div class="inst-dialog-foot">
      <button class="inst-btn" type="submit">Close</button>
    </div>
  </form>
</dialog>
```

## Contract

`inst-sheet` is a modifier over `inst-dialog` rather than a replacement: the
class of the base stays, and the parts (`-head`, `-body`, `-foot`) are the
same.

| What | Required | Why |
|---|---|---|
| An `inst-dialog` together with `inst-sheet` | yes | The surface, the shadow, the ground and the layout of the parts come from the base |
| A native `<dialog>` | yes | Modality and `Escape` come from the platform |
| `showModal()` | yes | `show()` opens it non-modally: no ground and no inert background |
| A `<form method="dialog">` inside | yes | Closing by a button with no script |
| An `aria-label` on the cross | yes | Otherwise a button with no text has no accessible name |

### Accessibility

| | |
|---|---|
| Focus | The focus trap and the return to the button that called it come from `showModal()` |
| The keyboard | `Escape` closes. Buttons inside a `<form method="dialog">` close on a press |
| The background is inert | The ground switches the background off for the mouse and for a screen reader. If access to the background is wanted, this is not a sheet but a column of [the shell](../../layout/shell.md) |
| Scrolling the background | It stops by the same `html:has(dialog:modal)` rule as on a modal |
| Height | `100dvh` rather than `100vh`: on mobile the browser's panel does not cut off the footer with its buttons |
| Print | A sheet is not printed on paper |

## Anatomy

Every part comes from [the modal](./dialog.md): `inst-dialog-head`,
`inst-dialog-title`, `inst-dialog-close`, `inst-dialog-body`,
`inst-dialog-foot`, `inst-dialog-foot--start`, `inst-dialog-foot-note`.

The modifier changes the layout:

| Property | The modal | The sheet |
|---|---|---|
| width | `min(34rem, 100vw - var(--space-8))` | `min(28rem, 100vw)` |
| height | `min(80dvh, 100dvh - var(--space-8))` | `100dvh` |
| margins | `margin: auto` — centred | Pushed to the edge |
| radius | `--radius-lg` | `0` — the block meets the edge of the screen |
| border | Around the perimeter | On the side of the content only |

On a narrow screen `min(28rem, 100vw)` gives a full-width sheet — there is no
need to start a separate mobile variant.

## Variants

```html
<!-- At the end edge (the default) -->
<dialog class="inst-dialog inst-sheet">…</dialog>

<!-- At the start edge -->
<dialog class="inst-dialog inst-sheet inst-sheet--start">…</dialog>
```

| Class | The edge | The boundary |
|---|---|---|
| `inst-sheet` | The end one (the right in LTR) | A hairline along the start edge |
| `inst-sheet--start` | The start one (the left in LTR) | A hairline along the end edge |

The sides are logical rather than physical: in RTL a sheet mirrors itself, and
nothing has to be edited.

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
