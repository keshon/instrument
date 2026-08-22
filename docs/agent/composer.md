---
title: Composer
group: Agent layer
layout: component
source: src/agent.css
api:
  - { name: "inst-composer", kind: "class", doc: "The frame. It carries the border, the background and the focus ring" }
  - { name: "inst-composer-bar", kind: "class", doc: "The controls. Inside the frame they belong to the message; standing under it they belong to the session. Anything in it that is not a control reads as metadata, and its glyphs take the small icon" }
  - { name: "inst-composer-bar--edge", kind: "modifier", doc: "The bar in the trailing slot beside the field rather than on its own line. Sunk to the bottom, so an empty composer is one line tall" }
  - { name: "--control-h-md", kind: "token", value: "32px", doc: "The cap on growth is eight of these" }
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
<div class="inst-stack inst-stack--tight">
  <div class="inst-composer">
    <textarea class="inst-textarea" rows="1" placeholder="Describe a task or ask a question" aria-label="Describe a task or ask a question"></textarea>
    <div class="inst-composer-bar inst-composer-bar--edge">
      <button class="inst-btn inst-btn--sm inst-btn--icon inst-btn--ghost" type="button" aria-label="Send"><svg class="inst-icon" aria-hidden="true"><use href="#i-enter"/></svg></button>
    </div>
  </div>
  <div class="inst-composer-bar">
    <button class="inst-btn inst-btn--xs inst-btn--ghost" type="button" aria-haspopup="menu">Auto</button>
    <button class="inst-btn inst-btn--xs inst-btn--icon inst-btn--ghost" type="button" aria-label="Add context"><svg class="inst-icon" aria-hidden="true"><use href="#i-plus"/></svg></button>
    <span class="inst-cluster-spacer"></span>
    <button class="inst-btn inst-btn--xs inst-btn--ghost" type="button" aria-haspopup="menu">Opus 5</button>
    <span class="inst-badge" data-tone="warn">context 90%</span>
  </div>
</div>
```

## Contract

| What | Required | Why |
|---|---|---|
| An `aria-label` or a `<label>` on the field | yes | A placeholder is not a name: it goes as soon as a character is typed |
| The **message's** controls inside the frame | yes | Attaching and sending belong to the field. Outside they are a second object, and the focus ring then lights up half the control. What outlives the message — the mode, the model, the context — stands **under** the frame in a bar of its own, and being a second object is the point |
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

## Variants

The bar can take the **whole width inside the frame** instead of the trailing
slot. That is the shape for a composer whose controls are more than sending —
an attachment, a model, a count of what is left — and it costs the empty field
its single line, which is the trade.

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

## Sizes

An empty composer is **one line tall**. It was two and a half — the height a
[textarea](../components/inputs/input.md) has by definition, which is right
for a form, where the field says how much answer is expected. A composer says
something else: three rows of nothing promise a paragraph nobody has typed, and
the shape of the object then lies about its content.

So the floor goes and the content sets the height, through `field-sizing:
content`. The cap is eight controls, after which the field scrolls instead of
eating the thread above it. Where `field-sizing` is not implemented the field
takes the height of its `rows` and scrolls — a smaller field, not a broken one,
and the alternative is a script the kit does not have.

`inst-composer-bar--edge` puts the bar in the **trailing slot** rather than on
its own line, which is what makes one line possible: sending is one control and
does not deserve a row of its own. It sinks to the bottom of the field, so on
one line it reads as centred and on ten it stays with the last one.

```html preview
<div class="inst-composer">
  <textarea class="inst-textarea" rows="1" aria-label="A grown composer">The field grows with what is typed into it, and the send affordance stays with the last line rather than floating in the middle of an empty box.</textarea>
  <div class="inst-composer-bar inst-composer-bar--edge">
    <button class="inst-btn inst-btn--sm inst-btn--icon inst-btn--ghost" type="button" aria-label="Send"><svg class="inst-icon" aria-hidden="true"><use href="#i-enter"/></svg></button>
  </div>
</div>
```

## API

```api
```

## Related

```related
```
