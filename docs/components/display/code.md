---
title: Code and copyable value
group: Data display
layout: component
source: src/agent.css
js: Copying to the clipboard is done by `instrument.js`
api:
  - { name: "inst-code", kind: "class", doc: "A block of code. It scrolls horizontally" }
  - { name: "inst-code--numbered", kind: "modifier", doc: "A listing with a column of line numbers" }
  - { name: "inst-code-line", kind: "class", doc: "A line of a listing. Inside `--numbered` only" }
  - { name: "inst-code-num", kind: "class", doc: "The number. `--text-muted`, because it gets read. Not selected by the mouse" }
  - { name: "inst-code-src", kind: "class", doc: "The code itself. The one place `overflow-wrap: anywhere` lives in" }
  - { name: "inst-copyable", kind: "class", doc: "A single-line value with a button" }
  - { name: "inst-copy", kind: "class", doc: "The copy button. The work is done by [instrument.js](../../foundations/behavior.md). Inside `.inst-code` it takes the far top corner by itself" }
  - { name: "data-copy", kind: "attribute", doc: "Copy this value rather than the text of the block" }
  - { name: "data-copied-label", kind: "attribute", doc: "A phrase of your own for a screen reader instead of «Copied»" }
  - { name: "data-failed-label", kind: "attribute", doc: "The same for a failure. «Could not copy» by default" }
  - { name: "data-copied", kind: "attribute", value: "true · false", doc: "The response of the button. Set and removed by `instrument.js`" }
  - { name: "--surface-sunken", kind: "token" }
  - { name: "--border-subtle", kind: "token" }
  - { name: "--hairline", kind: "token" }
  - { name: "--radius-md", kind: "token" }
  - { name: "--space-4", kind: "token" }
  - { name: "--pad-panel", kind: "token" }
  - { name: "--text-xs", kind: "token" }
  - { name: "--font-mono", kind: "token" }
---

A block of code and a single-line value one wants to copy: a path, a hash, an
identifier.

The page is about both: the first is shown, the second is taken.

```html preview
<div class="inst-code">go -C tools run ./cmd/contrast
  <button class="inst-copy" type="button" aria-label="Copy the command">
    <svg class="inst-icon" aria-hidden="true"><use href="#i-copy"/></svg>
  </button>
</div>
<span class="inst-copyable">a4f7c2e91b0d5537
  <button class="inst-copy" type="button" aria-label="Copy a4f7c2e91b0d5537">
    <svg class="inst-icon" aria-hidden="true"><use href="#i-copy"/></svg>
  </button>
</span>
```

## Contract

| What | Required | Why |
|---|---|---|
| `tabindex="0"` and a named `role="region"` | yes, if there is nothing focusable inside | The block scrolls horizontally, and the scrolling has to be reachable from the keyboard |
| An `aria-label` **with the value** on the copy button | yes | A row of "Copy" buttons is indistinguishable by ear: "Copy a4f7c2e" |
| Nothing for speaking the copy | — | `instrument.js` starts a shared live region `[data-inst-live]` itself and speaks into it in words. There is nothing for the application to do here |

There is no line wrapping in a block of code: a wrap in a command changes its
meaning. Hence the horizontal scrolling, and `scrollbar-gutter` is reserved so
that the appearance of the bar does not shift the code by a pixel on every
update.

### Accessibility

| | |
|---|---|
| Scrolling | The block scrolls horizontally, so it has to be reachable from the keyboard: `tabindex="0"` and a named `role="region"` if there is nothing focusable inside |
| `scrollbar-gutter` | Reserved so that the appearance of the bar does not shift the code by a pixel on every update |
| The name of the copy button | An `aria-label` **with the value**: "Copy a4f7c2e". A row of "Copy" buttons is indistinguishable by ear |
| The result of a copy | `instrument.js` changes the look of the button **and** says it in words: a change of glyph is an event with no content for a screen reader. The kit has one live region per document, started at the first message and reused |
| Line wrapping | There is none in a block of code: a wrap in a command changes its meaning. Hence the horizontal scrolling |
| Type size | `--text-xs` — 12.5px. The size does not go below `--text-2xs` |

## Variants

### A copyable value

A single-line value with a button — a hash, an identifier, a path.

```html preview
<span class="inst-copyable">a4f7c2e
  <button class="inst-copy" type="button" aria-label="Copy a4f7c2e">
    <svg class="inst-icon" aria-hidden="true"><use href="#i-copy"/></svg>
  </button>
</span>
```

The button stands **in the flow**, with a gap, and is always visible. An
absolute in the corner appearing on hover is a trick of a block of code: there
the corner is free and there is plenty of text. A single-line value has no
corner, and the icon lay over the last characters — and on touch it was not
there at all, because on touch there is no hover.

The tap area of the button grows to the norm by an invisible `::before` — the
same trick as on a checkbox: the area grows rather than the glyph, and the line
stays the same height.

### A copyable value in a row

The commonest place for one, and the one that looks least like a block of code:
a column of identifiers in a table, a machine value in a list of properties. The
markup is the same `inst-copyable` — a table changes nothing about it.

```html preview
<table class="inst-table">
  <thead><tr><th>Run</th><th>Identifier</th><th>State</th></tr></thead>
  <tbody>
    <tr><td>nightly-4127</td>
      <td><span class="inst-copyable"><code>2f8a1c04-9b7e-4d31-a5f0-c6e2</code>
        <button class="inst-copy" type="button" aria-label="Copy 2f8a1c04-9b7e-4d31-a5f0-c6e2">
          <svg class="inst-icon" aria-hidden="true"><use href="#i-copy"/></svg>
        </button></span></td>
      <td><span class="inst-badge" data-tone="running"><span class="inst-dot"></span>running</span></td></tr>
    <tr><td>audio-pass</td>
      <td><span class="inst-copyable"><code>7c31b9de-2a44-4e08-8f1b-30d7</code>
        <button class="inst-copy" type="button" aria-label="Copy 7c31b9de-2a44-4e08-8f1b-30d7">
          <svg class="inst-icon" aria-hidden="true"><use href="#i-copy"/></svg>
        </button></span></td>
      <td><span class="inst-badge" data-tone="error"><span class="inst-dot"></span>fell</span></td></tr>
  </tbody>
</table>
```

Three things about a row that are easy to get wrong, and the reference got two
of them wrong itself before this section was written:

| What | Why |
|---|---|
| The button is **always visible** rather than shown on hover over the row | A row of identifiers on a touchscreen has no hover, so a button revealed by one does not exist there. The reference used to hide it that way in its own API tables |
| The colour is `--text-muted` rather than `--text-faint` | Faint is drawn at a threshold of 3:1 and is meant for lines and dividers. A control is read, so it is due the 4.5:1 of muted |
| The `aria-label` carries **the value** | A column of ten "Copy" buttons is one word repeated ten times by ear. "Copy 2f8a1c04" says which row it is |

The size of the glyph is not written here and must not be: `.inst-copy` is
`--size-icon-sm` wherever it stands, because a glyph beside a value serves the
value. Call sites that chose for themselves chose differently — that is exactly
how one kit came to have an 18px copy button and a 10px one on the same page.

### A listing with line numbers

```html preview
<div class="inst-code inst-code--numbered">
  <div class="inst-code-line"><span class="inst-code-num">41</span><span class="inst-code-src">const size = 256;</span></div>
  <div class="inst-code-line"><span class="inst-code-num">42</span><span class="inst-code-src">let seed = Date.now();</span></div>
  <div class="inst-code-line"><span class="inst-code-num">43</span><span class="inst-code-src">return generate(size, seed);</span></div>
</div>
```

A number is not selected by the mouse — `user-select: none`: what gets copied
is the code rather than the page. The column of numbers is exactly `5ch`, and
that is not a margin by eye: a five-digit line number occurs in real files, and
a column growing from its content would shift the code on every load.

| | `inst-code` | `inst-copyable` |
|---|---|---|
| How many lines | Any number | One |
| Scrolling | Horizontal | None |
| The button | No | Yes |
| An example | A command, a config | A hash, a path, an identifier |

## JS

Include the module once per page — there is no need to initialise the
components one by one, `instrument.js` works by delegation and sees nodes that
arrived later.

```html
<script type="module" src="instrument.js"></script>
```

### What `instrument.js` does

The examples on this page are live: hover over the block of code and press the
copy button — the glyph turns into a tick.

Copying entire: the text into the clipboard, a coloured response on the button
and the result announced aloud. The response is obligatory — a clipboard shows
nothing, and without it a person presses the button a second time.

The source is the text of the nearest `.inst-code` or `.inst-copyable`
**without the button itself**. `textContent` is read rather than `innerText`:
the second returns an empty string on a collapsed block, and the copy would go
nowhere in silence.

### Options

| Attribute | What it does |
|---|---|
| `data-copy` | Copy this value rather than the text of the block. Wanted when an abbreviation is shown and the full thing has to be copied |
| `data-copied-label` | A phrase of your own for a screen reader instead of «Copied» |
| `data-failed-label` | The same for a failure. «Could not copy» by default |

```html
<span class="inst-copyable">a4f7c2e
  <button class="inst-copy" type="button" data-copy="a4f7c2e91b0d5537"
          aria-label="Copy the full hash"></button>
</span>
```

### Events

`inst:copy` bubbles with `{ text }` and is cancellable: `preventDefault()`
means "I copy it myself", and `instrument.js` does not touch the clipboard.

```js
document.addEventListener('inst:copy', (e) => {
  analytics.track('copy', { length: e.detail.text.length });
});
```

### What `instrument.js` does not do

Syntax highlighting. It calls for parsing a language, and that is the work of
the application.

## API

```api
```

## Related

```related
```
