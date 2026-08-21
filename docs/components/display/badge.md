---
title: Badge
group: Data display
layout: component
source: src/status.css
api:
  - { name: "inst-badge", kind: "class", doc: "The base one" }
  - { name: "inst-badge--outline", kind: "modifier", doc: "The tone is carried by the border, with no fill. For a row where a badge stands beside icons and a plate is the heaviest thing in it" }
  - { name: "inst-badge--count", kind: "modifier", doc: "A round counter laid over an icon: no horizontal padding, a solid fill in the accent" }
  - { name: "inst-dot", kind: "class", doc: "The state dot inside. It inherits the tone from the badge" }
  - { name: "inst-nav-count", kind: "class", doc: "A modifier of position for a counter in navigation" }
  - { name: "data-tone", kind: "attribute", doc: "`neutral` `running` `ok` `warn` `error`. Closed" }
  - { name: "--control-h-xs", kind: "token" }
  - { name: "--radius-xs", kind: "token" }
  - { name: "--text-2xs", kind: "token" }
  - { name: "--weight-medium", kind: "token" }
  - { name: "--space-2", kind: "token" }
  - { name: "--space-3", kind: "token" }
  - { name: "--size-dot", kind: "token" }
  - { name: "--tone-ink", kind: "token" }
  - { name: "--tone-bg", kind: "token" }
---

The state of an object, named by a word from the library's fixed vocabulary.
Five values, and a sixth cannot be added by accident.

```html preview
<span class="inst-badge"><span class="inst-dot"></span>queued</span>
<span class="inst-badge" data-tone="running"><span class="inst-dot"></span>running</span>
<span class="inst-badge" data-tone="ok"><span class="inst-dot"></span>done</span>
<span class="inst-badge" data-tone="warn"><span class="inst-dot"></span>with remarks</span>
<span class="inst-badge" data-tone="error"><span class="inst-dot"></span>failed</span>

<span class="inst-badge inst-badge--outline" data-tone="ok">New</span>
<span class="inst-badge inst-badge--outline" data-tone="warn">beta</span>
```

## Contract

| What | Required | Why |
|---|---|---|
| A word inside the badge | yes | Colour has no right to be the only carrier of a state |
| A `data-tone` from the closed vocabulary | yes, on a status | Five values. An invented sixth will not be painted at all |
| An `inst-dot` inside | yes, on a status | A second carrier of the mark beside the word |
| A `data-tone` on the dot | no | The dot inherits the tone from the badge: the attribute is set **once per group** |

### Accessibility

| | |
|---|---|
| Colour is not the only carrier | A status carries a dot **and** a word. A reader who does not tell tones apart reads the label |
| Contrast | The label of every tone is checked at 4.5:1 against its own tinted background in five themes. That is what the tone has a separate step of the ramp for text for: the step of the fill does not reach 4.5:1 on itself |
| Reduced motion | The pulse of a `running` dot **slows down rather than going out**: an indicator of busyness has to stay visible |
| Forced-colours mode | The dot survives the reset (`forced-color-adjust: none`), because it carries meaning |
| Type size | `--text-2xs` — 11px, the lower limit of type size. We go no lower |

## Anatomy

The dot is the carrier of a second mark: **colour has no right to be the only
carrier of a state**, so a status travels with a dot *and* a word. A badge with
no dot is used where it is not a status but a label of an object.

```html preview
<span class="inst-badge" data-tone="ok"><span class="inst-dot"></span>done</span>
<span class="inst-badge" data-tone="ok">done</span>
<span class="inst-badge">terrain_chunk_04</span>
```

A badge reads `--tone-ink` and `--tone-bg` and knows nothing else about the
tone.

## Variants

A tone is `data-tone`, **one vocabulary for the whole library**. The same
attribute with the same values works on a note, a banner, a dot and a log row.

| Tone | Means |
|---|---|
| no attribute | Neutral, not yet begun |
| `data-tone="neutral"` | Explicitly neutral. It also plays the part of info — blue in the library is taken by the accent |
| `data-tone="running"` | Happening now. The dot inside **pulses** |
| `data-tone="ok"` | Finished successfully |
| `data-tone="warn"` | Finished with remarks |
| `data-tone="error"` | Failed |

## Composition

### A counter in navigation

```html preview
<a class="inst-nav-item" href="#">Queue
  <span class="inst-badge inst-nav-count">7</span></a>
```

### In a panel header

```html preview context
<div class="inst-panel">
  <div class="inst-panel-header">
    <span class="inst-panel-title">Run</span>
    <span class="inst-badge" data-tone="running"><span class="inst-dot"></span>running</span>
  </div>
  <div class="inst-panel-body">worldgen-01 · 4 files</div>
</div>
```

## Customization

If you are building a component of your own on a tone: a tone has **two**
foregrounds, and the one to take is the one that matches the work.

| Variable | What for | Threshold |
|---|---|---|
| `--tone-ink` | Text | 4.5:1 |
| `--tone-mark` | A mark with no label: a dot, a caret, the fill of a meter | 3:1, and against the track as well |

On the statuses they coincide, on the accent they part.

## API

```api
```

## Related

```related
```
