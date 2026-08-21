---
title: Skeleton
group: Feedback
layout: component
source: src/agent.css
api:
  - { name: "inst-skeleton", kind: "class", doc: "The base one. Always required" }
  - { name: "inst-skeleton--title", kind: "modifier", doc: "A width of `40%`" }
  - { name: "inst-skeleton--line", kind: "modifier", doc: "A width of `100%`" }
  - { name: "inst-skeleton--short", kind: "modifier", doc: "A width of `62%`" }
  - { name: "--surface-sunken", kind: "token" }
  - { name: "--radius-xs", kind: "token" }
  - { name: "--ease-in-out", kind: "token" }
---

The place where content has not arrived yet. A library that is about waiting
needs a primitive for waiting for content.

```html preview
<div class="inst-stack inst-stack--tight">
  <div class="inst-skeleton inst-skeleton--title"></div>
  <div class="inst-skeleton inst-skeleton--line"></div>
  <div class="inst-skeleton inst-skeleton--short"></div>
</div>
```

## Contract

| What | Required | Why |
|---|---|---|
| `aria-busy="true"` on the container | yes | A skeleton is styling rather than a role. Without the attribute the page is simply empty to a screen reader |
| `inst-skeleton` on every bar | yes | The base: the surface, the radius, the pulse, a height of `1lh` |
| The shape repeats the content to come | yes | A skeleton of the wrong shape lies twice |

The height `1lh` is the real height of a line at the current type size. So a
skeleton takes exactly as much room as the text that arrives will, and nothing
jumps at the substitution.

### Accessibility

| | |
|---|---|
| The region has to declare it is busy | A skeleton is styling rather than a role. A container that is loading carries `aria-busy="true"`, otherwise for a screen reader the page is simply empty |
| It is not read | The bars are empty and hold no text: a screen reader has nothing to read in them and no need |
| Reduced motion | The pulse **slows to 3s rather than going out**: an indicator of busyness has to stay visible |
| Amplitude | The `opacity` does not go below `0.45`: a flicker of greater amplitude irritates peripheral vision |
| Print | It is not printed on paper: print is static, and waiting on paper is meaningless |

## Variants

```html preview
<div class="inst-stack inst-stack--tight">
  <div class="inst-skeleton inst-skeleton--title"></div>
  <div class="inst-skeleton inst-skeleton--line"></div>
  <div class="inst-skeleton inst-skeleton--line"></div>
  <div class="inst-skeleton inst-skeleton--short"></div>
</div>
```

| Class | Width | The work |
|---|---|---|
| `inst-skeleton` | inherited | The base: the surface, the radius, the pulse. The height is `1lh` |
| `inst-skeleton--title` | `40%` | The heading of a paragraph |
| `inst-skeleton--line` | `100%` | A full line of text |
| `inst-skeleton--short` | `62%` | The last line of a paragraph |

`--short` at the end: a real paragraph rarely ends exactly at the right edge,
and a bar of the same length reads as a table rather than as text.

The width can be set by a rule of your own too — the modifiers are only three
frequent values.

## Behavior

### The pulse

A skeleton pulses by **opacity**: from `1` to `0.45` and back over 1.6s. The
keyframes are its own, running infinitely.

A travelling sheen is a gradient, and gradients are forbidden in the library.
But the ban is not the whole of it: a sheen depicts motion from left to right,
that is, **filling**, and a skeleton fills nothing. A pulse is calmer and more
honest — the block does not pretend to be filling up.

## Composition

### In the body of a panel

```html preview context
<div class="inst-panel">
  <div class="inst-panel-header">
    <span class="inst-panel-title">Output</span>
  </div>
  <div class="inst-panel-body">
    <div class="inst-stack inst-stack--tight" aria-busy="true">
      <div class="inst-skeleton inst-skeleton--title"></div>
      <div class="inst-skeleton inst-skeleton--line"></div>
      <div class="inst-skeleton inst-skeleton--short"></div>
    </div>
  </div>
</div>
```

A skeleton repeats **the shape of what will come**. Three bars in the place of
three lines are a hint; twenty bars in the place of a card with a button are
noise.

## API

```api
```

## Related

```related
```
