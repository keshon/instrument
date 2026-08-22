---
title: Composer
group: Agent layer
layout: component
source: src/agent.css
api:
  - { name: "inst-composer", kind: "class", doc: "The frame. It carries the border, the background and the focus ring" }
  - { name: "inst-composer-bar", kind: "class", doc: "The controls. Inside the frame they belong to the message; standing under it they belong to the session. Anything in it that is not a control reads as metadata, and its glyphs take the small icon" }
  - { name: "--border-control", kind: "token" }
  - { name: "--accent-border", kind: "token" }
  - { name: "--surface-field", kind: "token" }
  - { name: "--focus-ring", kind: "token" }
  - { name: "--radius-lg", kind: "token" }
  - { name: "--control-pad-sm", kind: "token" }
  - { name: "--space-3", kind: "token" }
  - { name: "--text-xs", kind: "token" }
  - { name: "--text-muted", kind: "token" }
---

The field a turn is written in: one frame holding the text and the controls
that send it.

```html preview
<div class="inst-composer">
  <textarea class="inst-textarea" rows="3" placeholder="Describe a task" aria-label="Describe a task"></textarea>
  <div class="inst-composer-bar">
    <button class="inst-btn inst-btn--xs inst-btn--ghost" type="button">Attach</button>
    <span class="inst-badge">opus</span>
    <span class="inst-cluster-spacer"></span>
    12 480 / 200 000
    <button class="inst-btn inst-btn--sm inst-btn--primary" type="button">Send</button>
  </div>
</div>
```

## Contract

| What | Required | Why |
|---|---|---|
| An `aria-label` or a `<label>` on the field | yes | A placeholder is not a name: it goes as soon as a character is typed |
| The bar **inside** the frame | yes | The controls belong to the field. Under it they are a second object, and the focus ring then lights up half the control |
| A `<button type="button">` for sending | yes | Inside a form an unnamed type submits it |
| Sending on `Enter` | no | The kit styles; whether `Enter` sends or breaks the line is the application's decision, and both are defensible |

The border lives on the **composer** and the field hands its own outwards —
the arrangement [an input group](../components/inputs/input.md) uses. A
textarea beside a toolbar gives two borders with a gap; an input group
stretches the send button to the full height of the field; a panel around both
draws a frame inside a frame.

### Accessibility

| | |
|---|---|
| The name | On the field, not on the frame: the frame is a `<div>` and takes no focus |
| Focus | The ring goes round the whole frame, so what is highlighted is the object the keyboard is in |
| Contrast | The border is `--border-control` at 3.46:1 — a load-bearing line, because here the border **is** the control |

## API

```api
```

## Related

```related
```
