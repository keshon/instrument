---
title: Tag
group: Data display
layout: component
source: src/data.css
js: Removing a tag is done by `instrument.js`
api:
  - { name: "inst-tag", kind: "class", doc: "The label" }
  - { name: "inst-tag-remove", kind: "class", doc: "The removal button. A real `<button>`; the removal is carried out by [instrument.js](../../foundations/behavior.md)" }
  - { name: "data-value", kind: "attribute", doc: "On `.inst-tag`. What arrives in the `detail` of `inst:remove` instead of the label's text" }
  - { name: "--control-h-sm", kind: "token" }
  - { name: "--radius-sm", kind: "token" }
  - { name: "--text-xs", kind: "token" }
  - { name: "--space-2", kind: "token" }
  - { name: "--space-3", kind: "token" }
  - { name: "--surface-recessed", kind: "token" }
  - { name: "--border", kind: "token" }
  - { name: "--hairline", kind: "token" }
  - { name: "--size-chevron", kind: "token" }
  - { name: "--text-muted", kind: "token" }
  - { name: "--err-text", kind: "token" }
---

An arbitrary label started by the user. It can be removed — and that is the
main difference from [a badge](./badge.md), which has five values from the
library's vocabulary and cannot be removed.

```html preview
<span class="inst-tag">worldgen<button class="inst-tag-remove" type="button" aria-label="Remove the worldgen label"></button></span>
<span class="inst-tag">nightly<button class="inst-tag-remove" type="button" aria-label="Remove the nightly label"></button></span>
<span class="inst-tag">priority:high</span>
```

## Contract

| What | Required | Why |
|---|---|---|
| A real `<button>` for removal | yes | A cross as text cannot be pressed from the keyboard and there is nothing to speak |
| An `aria-label` **with the word of the label** | yes | "Remove the worldgen label" rather than "Remove": a row of ten "Remove"s is indistinguishable by ear |
| Moving the focus after a removal | yes, in the application | A vanished element takes the focus nowhere |
| The removal button | no | It is added only where removal is possible |

### Accessibility

| | |
|---|---|
| The removal button | A real `<button>` rather than a cross as text: it has to be pressable from the keyboard and speakable |
| The name of the button | An `aria-label` **with the word of the label**: "Remove the worldgen label" rather than "Remove". A row of ten "Remove" buttons is indistinguishable by ear |
| The focus after a removal | It goes nowhere if the element has vanished. The application has to move it to the next tag or to the container |
| The tap target | The glyph is `--size-chevron`, that is 10px. This is **below the 24px of WCAG 2.5.8**: the tap area is widened by a pseudo-element while the glyph itself stays small |
| Tone | A tag has none, so the question of "colour as the only carrier" does not arise |

## Variants

### Without a removal button

```html preview
<span class="inst-tag">read-only</span>
<span class="inst-tag">priority:high</span>
```

The button is added only where removal is really possible. A cross that does
nothing is worse than no cross.

## JS

Include the module once per page — there is no need to initialise the
components one by one, `instrument.js` works by delegation and sees nodes that
arrived later.

```html
<script type="module" src="instrument.js"></script>
```

### What `instrument.js` does

The example in the header is live: press the cross — the tag is removed and the
focus goes to its neighbour.

A press on `.inst-tag-remove` removes the tag and **moves the focus to the
neighbouring cross** — otherwise a deleted element takes the focus nowhere, and
somebody on a keyboard ends up at the start of the document.

### Events

`inst:remove` bubbles from the tag, and its `detail` is `{ value }`: the
`data-value` if there is one, otherwise the text of the label.

The event is cancellable, and in a data-driven application cancelling it is the
norm: the markup there is recreated by a re-render rather than by
`instrument.js`.

```js
list.addEventListener('inst:remove', (e) => {
  e.preventDefault();            // the tag will stay in the markup
  store.dropTag(e.detail.value); // the re-render will take it away
});
```

### Options

| Attribute | What it does |
|---|---|
| `data-value` | What arrives in `detail` instead of the label's text — an identifier rather than a human name |

## API

```api
```

## Related

```related
```
