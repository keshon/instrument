---
title: Banner
group: Feedback
layout: component
source: src/feedback.css
js: Closing a banner belongs to the application layer. A banner has no cross
api:
  - { name: "inst-banner", kind: "class", doc: "The container. The fill of the tone, the border, the icon" }
  - { name: "inst-banner-body", kind: "class", doc: "The text part. `min-inline-size: 0`, so that a long word does not push the banner wide" }
  - { name: "inst-banner-title", kind: "class", doc: "The heading, at medium weight" }
  - { name: "inst-banner-text", kind: "class", doc: "The detail, `--text-secondary`" }
  - { name: "inst-banner-actions", kind: "class", doc: "The actions at the far edge. They do not shrink and wrap when there is not enough room" }
  - { name: "data-tone", kind: "attribute", doc: "`neutral` `running` `ok` `warn` `error`. Closed" }
  - { name: "--tone-bg", kind: "token" }
  - { name: "--tone-ink", kind: "token" }
  - { name: "--border", kind: "token" }
  - { name: "--hairline", kind: "token" }
  - { name: "--radius-md", kind: "token" }
  - { name: "--pad-panel", kind: "token" }
  - { name: "--space-3", kind: "token" }
  - { name: "--space-4", kind: "token" }
  - { name: "--gap-inline", kind: "token" }
  - { name: "--size-icon", kind: "token" }
  - { name: "--text-sm", kind: "token" }
  - { name: "--text-secondary", kind: "token" }
  - { name: "--weight-medium", kind: "token" }
  - { name: "--leading-ui", kind: "token" }
---

A message at the level of the page: it is about the whole screen rather than
about what it stands beside. So a banner has a fill, a border and room for an
action.

```html preview
<div class="inst-banner" data-tone="warn">
  <div class="inst-banner-body">
    <div class="inst-banner-title">Access expires in 3 days</div>
    <div class="inst-banner-text">After 8 August starting runs will stop.</div>
  </div>
  <div class="inst-banner-actions">
    <button class="inst-btn inst-btn--sm" type="button">Renew</button>
  </div>
</div>
```

## Contract

| What | Required | Why |
|---|---|---|
| An `inst-banner-body` around the text | yes | `min-inline-size: 0` — otherwise a long word pushes the banner wide |
| An `aria-live="polite"` or a `role="alert"` | yes, if the banner appeared in answer to an action | The role is not set automatically: the library does not know whether the banner appeared just now or was there from the start |
| An `aria-label` on the close button | yes, if the banner closes | A cross with no name does not say what it closes |
| `inst-banner-text` | no | The detail is not always wanted: the heading may stay the only thing |

### Accessibility

| | |
|---|---|
| Appearing | A banner that arose in answer to an action has to reach an `aria-live="polite"`, and a message about a refusal a `role="alert"`. The role is not set automatically: the library does not know whether the banner appeared just now or was on the page from the start |
| Colour is not the only carrier | The tone is carried by the icon **and** the text. The icon appeared precisely because `warn` and `error` used to differ by colour alone |
| Contrast | The heading takes `--text-primary` and the detail `--text-secondary` — both at 4.5:1 against `--tone-bg` in five themes |
| Closing | It is no business of the library's. If a banner closes, the close button is an ordinary icon [button](../actions/button.md) with an `aria-label` |
| Print | A banner is printed and does not break between sheets: `break-inside: avoid`. A message at the level of the page has to reach the report |

## Anatomy

| Class | The work |
|---|---|
| `inst-banner` | The container: the fill of the tone, the border, the icon |
| `inst-banner-body` | The text part. `min-inline-size: 0` against being pushed wide |
| `inst-banner-title` | The heading. It may be the only content |
| `inst-banner-text` | The detail in one phrase |
| `inst-banner-actions` | The actions at the far edge. They do not shrink and wrap when there is not enough room |

The actions are pushed to the far edge and do not shrink, so a long text wraps
and the button stays whole.

The icon is centred on the **first line** of the text:
`calc((1lh - var(--size-icon)) / 2)`. The formula adjusts to the type size and
the density by itself, so a banner of one line and a banner of three look
equally aligned.

## Variants

```html preview
<div class="inst-banner" data-tone="warn">
  <div class="inst-banner-body">
    <div class="inst-banner-title">Access expires in 3 days</div>
    <div class="inst-banner-text">After 8 August starting runs will stop.</div>
  </div>
  <div class="inst-banner-actions">
    <button class="inst-btn inst-btn--sm" type="button">Renew</button>
  </div>
</div>
<div class="inst-banner" data-tone="error">
  <div class="inst-banner-body">
    <div class="inst-banner-title">Storage maintenance</div>
    <div class="inst-banner-text">The artefacts are unavailable until 21:00.</div>
  </div>
</div>
<div class="inst-banner" data-tone="ok">
  <div class="inst-banner-body"><div class="inst-banner-title">The changes are saved</div></div>
</div>
```

The tone is the same `data-tone` as on [the badge](../display/badge.md) and
everywhere else. It carries an **icon**. The glyph comes from a shared set of
masks — one for the banner and [the note](./note.md).

| Tone | The icon | When |
|---|---|---|
| no attribute | none | A neutral message with no glyph |
| `data-tone="neutral"` | an "i" in a circle | Information. The part of info is played by `neutral` — blue in the library is taken by the accent |
| `data-tone="running"` | an "i" in a circle | Happening right now: maintenance, a migration |
| `data-tone="ok"` | a tick in a circle | It has happened successfully |
| `data-tone="warn"` | a triangle | A condition that will soon become a problem |
| `data-tone="error"` | a cross in a circle | Something is not working right now |

With no `data-tone` the icon is not drawn at all (`content: none`), and the
banner stays on a neutral muted surface.

## JS

Include the module once per page — there is no need to initialise the
components one by one, `instrument.js` works by delegation and sees nodes that
arrived later.

```html
<script type="module" src="instrument.js"></script>
```

### What `instrument.js` does

Nothing. A banner has no cross: a banner reports the state of the system, and a
state does not disappear from being closed. Closing makes sense only where the
application knows what to do with something closed.

### If closing is wanted after all

The button is an ordinary one and has to be added yourself. Along with it comes
the decision of what "closed" means: until a reload, until the end of the
session, or for good.

```js
document.addEventListener('click', (e) => {
  const id = e.target.closest('[data-close]')?.dataset.close;
  if (!id) return;
  document.getElementById(id).hidden = true;
  sessionStorage.setItem('hidden:' + id, '1'); // until the session ends, not for good
});
```

`hidden` rather than removing the node: the banner will come back when the
state repeats, and there is no point in recreating the markup.

## API

```api
```

## Related

```related
```
