---
title: Empty state
group: Feedback
layout: component
source: src/feedback.css
api:
  - { name: "inst-empty", kind: "class", doc: "The container; it centres the content" }
  - { name: "inst-empty-title", kind: "class", doc: "The heading: why it is empty" }
  - { name: "inst-empty-desc", kind: "class", doc: "The explanation" }
  - { name: "--space-9", kind: "token" }
  - { name: "--space-6", kind: "token" }
  - { name: "--gap-inline", kind: "token" }
  - { name: "--text-sm", kind: "token" }
  - { name: "--text-muted", kind: "token" }
---

A region with nothing to show yet. An empty state is not the absence of an
interface but a screen in its own right: it explains **why** it is empty and
gives a way out.

```html preview
<div class="inst-panel">
  <div class="inst-empty">
    <div class="inst-empty-title">No runs yet</div>
    <div class="inst-empty-desc">A launch creates a queue of agents and a log, and they will appear here.</div>
    <button class="inst-btn inst-btn--sm" type="button">Start the first one</button>
  </div>
</div>
```

## Contract

| What | Required | Why |
|---|---|---|
| A heading naming the **reason** | yes | "Nothing found" with no reason reads as a breakage |
| An `aria-live="polite"` on the container | yes, if the list emptied because of a filter | Otherwise for a screen reader simply nothing happened |
| A real heading of the right level | yes, if it takes a whole region | `inst-empty-title` styles but does not declare |
| An action | no | An ordinary [button](../actions/button.md); it has no class of its own |

### Accessibility

| | |
|---|---|
| Appearing after an action | If the list emptied because of a filter, the message has to reach an `aria-live="polite"` — otherwise for a screen reader simply nothing happened |
| The heading | A real heading of the right level, if the empty state takes a whole region |
| Not an illustration alone | There is no picture in the component: the meaning is carried by the text. An empty state made of one icon says nothing |
| Contrast | The text takes `--text-muted` (4.5:1) rather than `--text-faint`: this is data that gets read rather than decoration |

## Variants

Three different kinds of "empty". They are confused, and each calls for a text
of its own.

| Why it is empty | What to write | The action |
|---|---|---|
| **Not started yet** | What will appear here and where from | A button creating the first object |
| **A filter found nothing** | Which filter is applied | "Reset the filter" |
| **Everything is done** | That the list is empty for a good reason | Usually not needed |

```html preview
<div class="inst-empty">
  <div class="inst-empty-title">No runs yet</div>
  <div class="inst-empty-desc">A launch creates a queue of agents and a log, and they will appear here.</div>
  <button class="inst-btn inst-btn--sm" type="button">Start the first one</button>
</div>
<div class="inst-empty">
  <div class="inst-empty-title">Nothing matched the "failed" filter</div>
  <div class="inst-empty-desc">Over the last day every run finished successfully.</div>
  <button class="inst-btn inst-btn--sm" type="button">Reset the filter</button>
</div>
<div class="inst-empty">
  <div class="inst-empty-title">The queue is empty</div>
  <div class="inst-empty-desc">Every task has been taken.</div>
</div>
```

The text "Nothing found" with no reason given suits the second case only and in
the first reads as a breakage.

## API

```api
```

## Related

```related
```
