---
title: Accordion
group: Feedback
layout: component
source: src/feedback.css
api:
  - { name: "inst-accordion", kind: "class", doc: "The container. A column of sections with a gap of `--space-1`" }
  - { name: "inst-accordion-item", kind: "class", doc: "A section. A border, a radius, `overflow: hidden`" }
  - { name: "inst-accordion-head", kind: "class", doc: "The heading. The height from `--control-h-lg`, the chevron at the far edge" }
  - { name: "inst-accordion-body", kind: "class", doc: "The content" }
  - { name: "inst-accordion-item--flush", kind: "modifier", doc: "A section INSIDE somebody else's frame: a panel, a card, the field of an example. It has no border or radius of its own, and a seam above it is left" }
  - { name: "--space-1", kind: "token" }
  - { name: "--space-3", kind: "token" }
  - { name: "--pad-cell-x", kind: "token" }
  - { name: "--pad-panel", kind: "token" }
  - { name: "--gap-inline", kind: "token" }
  - { name: "--control-h-lg", kind: "token" }
  - { name: "--radius-md", kind: "token" }
  - { name: "--hairline", kind: "token" }
  - { name: "--border", kind: "token" }
  - { name: "--border-subtle", kind: "token" }
  - { name: "--surface-raised", kind: "token" }
  - { name: "--surface-hover", kind: "token" }
  - { name: "--size-chevron", kind: "token" }
  - { name: "--text-muted", kind: "token" }
  - { name: "--text-sm", kind: "token" }
  - { name: "--weight-medium", kind: "token" }
  - { name: "--dur-1", kind: "token" }
  - { name: "--ease-out", kind: "token" }
---

A group of collapsible sections with exactly one of them open. Built on native
`<details>` with a shared `name`, so mutual exclusion, expansion and the
keyboard come from the platform — with not one line of JS.

```html preview
<div class="inst-accordion">
  <details class="inst-accordion-item" name="faq" open>
    <summary class="inst-accordion-head">Why did the run take 42 seconds?</summary>
    <div class="inst-accordion-body">Three passes of laying out the resources plus baking the navmesh.
      Most of the time went on validating the paths, which restarted twice.</div>
  </details>
  <details class="inst-accordion-item" name="faq">
    <summary class="inst-accordion-head">Can the spending of tokens be limited?</summary>
    <div class="inst-accordion-body">Yes, in the parameters of the project. The limit is checked before
      every agent starts rather than after the fact — a run will not go into the red.</div>
  </details>
  <details class="inst-accordion-item" name="faq">
    <summary class="inst-accordion-head">What happens to the artefacts after a deletion?</summary>
    <div class="inst-accordion-body">They are deleted along with the run. A backup is kept for seven
      days, and restoring it goes through support.</div>
  </details>
</div>
```

## Contract

| What | Required | Why |
|---|---|---|
| A `<details>` under `inst-accordion-item` | yes | The expansion, the keyboard and the `open` state come from the platform |
| A `<summary>` under `inst-accordion-head` | yes | Only a `<summary>` takes focus and toggles a `<details>` |
| A shared `name` | yes, for an accordion | It makes the group mutually exclusive. Without it the sections are independent |
| Direct nesting | yes | The rules of an open section are written through `>`: the heading has to be a direct child of the `<details>` |
| A real heading inside the `<summary>` | no, but usually yes | A `<summary>` by itself does not reach the table of contents of a page |

### Accessibility

| | |
|---|---|
| The keyboard | `Tab` to the heading, `Enter` or `Space` to toggle. All of it native, and `instrument.js` intercepts nothing |
| Role and state | A `<summary>` declares itself a toggle and reports whether the section is expanded. An `aria-expanded` of your own need not and should not be set |
| Search on the page | The content of a closed section is found by the browser's search, and the section expands by itself — that is the behaviour of `<details>`, and it is what the native element was taken for |
| The heading | If the sections are sections of the page, nest a real heading of the right level inside the `<summary>`: a `<summary>` by itself does not reach the table of contents |
| The tap target | The heading is pressed entire, from edge to edge: `--control-h-lg` in height |
| Reduced motion | The turn of the chevron collapses along with every transition |
| Print | Closed sections **expand on paper**: a closed `<details>` on paper is content irrecoverably lost |

## Anatomy

| Element | Class | Tag |
|---|---|---|
| The container of the group | `inst-accordion` | `<div>` |
| A section | `inst-accordion-item` | `<details>` |
| The heading toggle | `inst-accordion-head` | `<summary>` |
| The content | `inst-accordion-body` | `<div>` |

The chevron is drawn by a pseudo-element of the heading and turns 180° on
opening — it has no class of its own and calls for no markup. The native
triangle of `<summary>` is removed.

The hairline under the heading appears **only on an open section**
(`.inst-accordion-item[open] > .inst-accordion-head`): on a closed one there is
nothing to part.

### Inside somebody else's frame

A fold-out often stands not on its own but at the bottom of a panel or a card:
"more fields", "the markup", "the journal". A border of its own would be a
second one a pixel from the first there, and a rounding of its own would be a
rounding inside a rounding, concentric with nothing.

```html preview context
<div class="inst-panel" style="inline-size:26rem">
  <div class="inst-panel-header"><span class="inst-panel-title">Run 4127</span></div>
  <div class="inst-panel-body">Built in 4.2 s. 318 tests, none failing.</div>
  <details class="inst-accordion-item inst-accordion-item--flush">
    <summary class="inst-accordion-head">Build parameters</summary>
    <div class="inst-accordion-body">Branch master, commit d7c9cd3, the cache on.</div>
  </details>
</div>
```

One seam above is left: the section **continues** what it stands under rather
than declaring itself a separate block. The `.inst-accordion` container is not
needed here — there is one section and it lives in somebody else's stack.

## Behavior

### Mutual exclusion

A shared `name` attribute on the `<details>` elements makes the group mutually
exclusive: opening the second section closes the first.

| What is wanted | How |
|---|---|
| Exactly one section open | A shared `name` on every `<details>` |
| Any number open | Do not set a `name` |
| Which is open on load | An `open` on the `<details>` wanted |

Remove the `name` and the sections become independent, and any number can be
opened. That is the only difference between an "accordion" and a "list of
fold-outs", and it is one attribute.

```html preview
<div class="inst-accordion">
  <details class="inst-accordion-item" open>
    <summary class="inst-accordion-head">With no name — they open independently</summary>
    <div class="inst-accordion-body">This one can be left open while the next is opened.</div>
  </details>
  <details class="inst-accordion-item">
    <summary class="inst-accordion-head">The second section</summary>
    <div class="inst-accordion-body">The first will not close.</div>
  </details>
</div>
```

## API

```api
```

## Related

```related
```
