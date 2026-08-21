---
title: Note
group: Feedback
layout: component
source: src/feedback.css
api:
  - { name: "inst-note", kind: "class", doc: "All of it. A note has no parts: it never has a heading or actions, by definition" }
  - { name: "data-tone", kind: "attribute", doc: "`neutral` `running` `ok` `warn` `error`. Closed. With no attribute the icon is not drawn" }
  - { name: "--size-marker", kind: "token" }
  - { name: "--size-icon", kind: "token" }
  - { name: "--tone-ink", kind: "token" }
  - { name: "--border-strong", kind: "token" }
  - { name: "--pad-panel", kind: "token" }
  - { name: "--space-3", kind: "token" }
  - { name: "--gap-inline", kind: "token" }
  - { name: "--text-sm", kind: "token" }
  - { name: "--text-secondary", kind: "token" }
  - { name: "--leading-ui", kind: "token" }
---

An inset in the flow: it explains what it stands beside. A bar at the near edge
and an icon — no fill, no border, no room for a button.

```html preview
<div class="inst-note" data-tone="warn">Irreversible: restoring from a backup takes up to a day.</div>
<div class="inst-note" data-tone="neutral">The limit is checked before every agent starts rather than after the fact.</div>
```

## Contract

| What | Required | Why |
|---|---|---|
| An `aria-live="polite"` on the container | yes, if the note appeared in answer to an action | The role is not set automatically: the library does not know whether it appeared just now or was there from the start |
| Text naming the substance | yes | The tone is carried by the icon, but the meaning is carried by the words |
| `data-tone` | no | With no attribute there is no icon and the bar is neutral |

A note has no parts: it never has a heading or actions, by definition. If they
are wanted, this is [a banner](./banner.md).

### Accessibility

| | |
|---|---|
| Appearing | A note that arose in answer to an action has to reach an `aria-live="polite"`. The role is not set automatically: the library does not know whether it appeared just now or was there from the start |
| Colour is not the only carrier | The tone is carried by the icon **and** the text. That is what the icon appeared for |
| Contrast | The text takes `--text-secondary` — 4.5:1 in five themes |
| The bar is not a carrier of meaning | The meaning is carried by the icon and the words; the bar only groups |
| Print | It is printed: it is part of the content rather than styling |
| Forced-colours mode | The icon moves to the system text colour and survives the reset |

## Anatomy

The corners are square: a note has one boundary only — the left one — and a
rounding on one side looks like a clipped block.

The icon is centred on the **first line** of the text:
`calc((1lh - var(--size-icon)) / 2)`. `1lh` is the real height of a line, so
the formula adjusts to the type size and the density by itself. A manual shift
of `-1px`, which was here before, lifted the icon 2.1px too high.

## Variants

```html preview
<div class="inst-note">No tone: there is no icon and the bar is neutral.</div>
<div class="inst-note" data-tone="neutral">Information. Blue in the library is taken by the accent, so the part of info is played by neutral.</div>
<div class="inst-note" data-tone="running">Happening right now.</div>
<div class="inst-note" data-tone="ok">The check has passed.</div>
<div class="inst-note" data-tone="warn">The run will finish, but with remarks.</div>
<div class="inst-note" data-tone="error">The validation of the paths did not pass.</div>
```

The tone is carried by the **icon**, and the bar alone is not enough.
Previously `warn` and `error` differed from each other by colour alone — a
direct breach of the library's own law that colour cannot be the only carrier.
The set of masks is shared with [the banner](./banner.md): one shape of glyph
for two components.

| Tone | The icon |
|---|---|
| no attribute | There is no icon, the bar is `--border-strong` |
| `data-tone="neutral"` | An "i" in a circle |
| `data-tone="running"` | An "i" in a circle |
| `data-tone="ok"` | A tick in a circle |
| `data-tone="warn"` | A triangle |
| `data-tone="error"` | A cross in a circle |

## API

```api
```

## Related

```related
```
