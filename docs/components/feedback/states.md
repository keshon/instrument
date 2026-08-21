---
title: Dot and caret
group: Feedback
layout: component
source: src/status.css
api:
  - { name: "inst-dot", kind: "class", doc: "The state dot. It reads `--tone-mark`" }
  - { name: "inst-caret", kind: "class", doc: "The caret of streaming text" }
  - { name: "data-tone", kind: "attribute", doc: "`neutral` `running` `ok` `warn` `error`. Closed. Set on the dot or on any ancestor of it" }
  - { name: "--size-dot", kind: "token" }
  - { name: "--radius-full", kind: "token" }
  - { name: "--tone-mark", kind: "token" }
  - { name: "--text-faint", kind: "token" }
  - { name: "--accent-mark", kind: "token" }
  - { name: "--ease-in-out", kind: "token" }
---

The two smallest marks of the library. A dot names the state of an object, a
caret shows that text is still arriving. Both are marks with no label, and so
both take `--tone-mark` rather than `--tone-ink`.

```html preview
<span class="inst-badge"><span class="inst-dot"></span>queued</span>
<span class="inst-badge" data-tone="running"><span class="inst-dot"></span>running</span>
<span class="inst-badge" data-tone="ok"><span class="inst-dot"></span>done</span>
<span class="inst-badge" data-tone="warn"><span class="inst-dot"></span>with remarks</span>
<span class="inst-badge" data-tone="error"><span class="inst-dot"></span>failed</span>
```

## Contract

| What | Required | Why |
|---|---|---|
| A word beside the dot | yes | A dot with no label is colour as the only carrier, that is, a defect |
| A `data-tone` on the dot **or** on an ancestor of it | yes | A dot does not know its tone: it reads `--tone-mark` from the nearest ancestor that set one |
| The caret flush against the last character | yes | Set off by a gap of `0.1em`. A caret moved off by a space points at the wrong place |

### Accessibility

| | |
|---|---|
| Colour is not the only carrier | A dot always stands beside a word. A dot with no label is colour as the only carrier, that is, a defect |
| Reduced motion on the dot | The pulse **slows to 3s rather than going out**: an indicator of busyness has to stay visible |
| Reduced motion on the caret | The blinking is **switched off** and the caret stays visible (`opacity: 1`). An exception to the exception: the blinking carries no information about speed, only about the fact, and the fact is carried by the caret itself |
| Forced-colours mode | The dot survives the reset (`forced-color-adjust: none`) and moves to `Highlight`, because it carries meaning |
| Print | The dot is printed, the caret is not: a blinking cursor on paper means nothing |
| Size | 6px is a mark rather than a tap target. What is pressed is the whole row, not the dot |

## Anatomy

### The dot

A circle of `--size-dot` (6px). The carrier of the **second mark** of a state:
colour has no right to be the only carrier, so a status travels with a dot
*and* a word.

### The caret

A rectangle the height of a line, blinking in steps. It is a **terminal
cursor**: it shows the point the next character will arrive at.

```html preview
<span>Streaming the answer<span class="inst-caret"></span></span>
```

Three decisions, each of which is easy to make the other way round:

- **It stands flush against the last character**, set off by a gap (`0.1em`).
  A caret moved off by a space points at the wrong place.
- **It blinks in steps** (`steps(2, start)`). A smooth fade reads as an
  artefact of the rendering.
- **The size is in `em`**: `0.45em` by `1.05em`. The caret is tied to the type
  size of the line it stands in and needs no adjusting for every size of text.

## Variants

`data-tone` is one vocabulary for the whole library, and it is closed.

| Tone | Means | The dot |
|---|---|---|
| no attribute | Not begun yet | `--text-faint` |
| `data-tone="neutral"` | Explicitly neutral. It also plays the part of info | Neutral |
| `data-tone="running"` | Happening now | **It pulses**: `opacity` 1 → 0.35 over 1.6s |
| `data-tone="ok"` | Finished successfully | The tone of success |
| `data-tone="warn"` | Finished with remarks | The tone of a remark |
| `data-tone="error"` | Failed | The tone of a refusal |

## Behavior

### Inheriting a tone

A dot does not know its tone — it reads `--tone-mark` from the nearest ancestor
that set one. So `data-tone` is set **once per group**.

```html preview
<span class="inst-badge" data-tone="ok"><span class="inst-dot"></span>done</span>
<span class="inst-dot" data-tone="error"></span>
```

| Where the `data-tone` stands | When it is so |
|---|---|
| On [the badge](../display/badge.md), [the queue row](../../agent/task.md), [the step](../../agent/step.md) | The ordinary case. The dot inherits |
| On the dot itself | The dot stands apart, with no wrapper carrying a tone |

The pulse is started not only by `data-tone="running"` but by the phase of a
component: `.inst-task[data-state="running"]` and
`.inst-step[data-state="running"]` give the dot the same animation. The list of
animated dots in `src/motion.css` has to match this — otherwise some of them
stop slowing down at `prefers-reduced-motion`.

### The animations

| What | How |
|---|---|
| The pulse of a `running` dot | `opacity` 1 → 0.35 → 1 over 1.6s, `--ease-in-out`, infinite |
| The blinking of the caret | `opacity` 0 for half the cycle, `steps(2, start)`, 1s, infinite |

## Customization

`--tone-mark` rather than `--tone-ink`: a tone has **two** foregrounds, and the
one to take is the one that matches the work.

| Variable | What for | Threshold |
|---|---|---|
| `--tone-ink` | Text | 4.5:1 |
| `--tone-mark` | A mark with no label | 3:1, and against the track as well |

On the statuses they coincide, on the accent they part.

## API

```api
```

## Related

```related
```
