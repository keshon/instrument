---
title: Turn
group: Agent layer
layout: component
source: src/agent.css
api:
  - { name: "inst-turn", kind: "class", doc: "One step of an exchange" }
  - { name: "data-from", kind: "attribute", value: "user", doc: "The person's turn, and only that one: it is recessed, the machine's turn is the default and bare" }
  - { name: "inst-turn-head", kind: "class", doc: "Who and when, in the type of metadata" }
  - { name: "inst-turn-who", kind: "class", doc: "The name inside that line" }
  - { name: "inst-turn-body", kind: "class", doc: "What was said. It may hold any component of this layer" }
  - { name: "inst-turn-actions", kind: "class", doc: "Copy, retry, branch. Always visible, never revealed on hover" }
  - { name: "--surface-recessed", kind: "token" }
  - { name: "--radius-lg", kind: "token" }
  - { name: "--pad-panel", kind: "token" }
  - { name: "--space-3", kind: "token" }
  - { name: "--text-xs", kind: "token" }
  - { name: "--text-muted", kind: "token" }
  - { name: "--text-secondary", kind: "token" }
  - { name: "--weight-medium", kind: "token" }
  - { name: "--leading-ui", kind: "token" }
---

One step of an exchange. Everything else in this layer shows work; a turn is
what the work happens **inside** — an answer is not a paragraph but a container
for steps, a diff, an approval, a budget.

```html preview
<div class="inst-stack inst-stack--loose">
  <div class="inst-turn" data-from="user">
    <div class="inst-turn-head"><span class="inst-turn-who">You</span><span>14:02</span></div>
    <div class="inst-turn-body">Add a share bar above the history strip.</div>
  </div>
  <div class="inst-turn">
    <div class="inst-turn-head"><span class="inst-turn-who">Agent</span><span>14:02</span></div>
    <div class="inst-turn-body">I'll put it above the strip and let one legend key both.</div>
    <div class="inst-turn-actions">
      <button class="inst-btn inst-btn--sm inst-btn--ghost" type="button">Copy</button>
      <button class="inst-btn inst-btn--sm inst-btn--ghost" type="button">Retry</button>
    </div>
  </div>
</div>
```

## Contract

| What | Required | Why |
|---|---|---|
| `data-from="user"` on the person's turn | yes | The machine's turn is the default. Marking both would mean an attribute value with nothing to declare |
| A name in `inst-turn-who` | yes | Colour and indent are not carriers: the recess is the same film in every theme, and a screen reader gets neither |
| A time | no | Useful in a long thread, noise in a short one |
| `inst-turn-actions` **visible** | yes, if present | Hover may strengthen an element, but it may not bring one into being — a row revealed on hover does not exist on a touchscreen |
| A thread as a list | no | Turns are siblings in an `inst-stack--loose`; the kit does not add a second name for a column |

The person's turn is **recessed**, the machine's has no surface at all: the
answer is the one that holds steps, diffs and approvals, and each of those
brings its own frame. A surface under them would put a plate inside a plate.

### Accessibility

| | |
|---|---|
| Who spoke | Named in the text, not by the surface. The recess is decoration and reaches no assistive technology |
| Streaming | An answer still arriving carries [the caret](../components/display/badge.md); the text itself is the announcement |
| Actions | Real buttons with labels. They stand in the flow, so `Tab` reaches them in reading order |

## Composition

### An answer that did work

```html preview context
<div class="inst-turn">
  <div class="inst-turn-head"><span class="inst-turn-who">Agent</span><span>14:02</span></div>
  <div class="inst-turn-body">Added, and the legend now keys the tones.</div>
  <details class="inst-step" open>
    <summary class="inst-step-head">
      <span class="inst-step-twist"></span>
      <span class="inst-dot" data-tone="ok"></span>
      <span class="inst-step-name">edit</span>
      <span class="inst-step-sub">src/data.css</span>
      <span class="inst-step-meta">0.4 s</span>
    </summary>
    <div class="inst-step-body">
      <div class="inst-code">.inst-share { display: flex; }</div>
    </div>
  </details>
  <div class="inst-turn-actions">
    <button class="inst-btn inst-btn--sm inst-btn--ghost" type="button">Copy</button>
    <button class="inst-btn inst-btn--sm inst-btn--ghost" type="button">Retry</button>
  </div>
</div>
```

## API

```api
```

## Related

```related
```
