---
title: Approval request
group: Agent layer
layout: component
source: src/agent.css
js: Sending the decision and moving the block into a decided state belong to the application layer
api:
  - { name: "inst-approval", kind: "class", doc: "The block" }
  - { name: "inst-approval-head", kind: "class", doc: "The heading with an icon" }
  - { name: "inst-approval-what", kind: "class", doc: "What will be done" }
  - { name: "inst-approval-effects", kind: "class", doc: "The list of consequences, on a `<ul>`" }
  - { name: "inst-approval-actions", kind: "class", doc: "The decision" }
  - { name: "inst-approval-verdict", kind: "class", doc: "What was decided and when. Shown after the answer" }
  - { name: "data-state", kind: "attribute", value: "pending · approved · denied", doc: "on `inst-approval`" }
  - { name: "data-tone", kind: "attribute", value: "warn · error", doc: "on an `<li>` of a consequence" }
  - { name: "--pad-card", kind: "token" }
  - { name: "--radius-lg", kind: "token" }
  - { name: "--space-3", kind: "token" }
  - { name: "--gap-inline", kind: "token" }
  - { name: "--text-sm", kind: "token" }
  - { name: "--text-xs", kind: "token" }
  - { name: "--warn-bg", kind: "token" }
  - { name: "--warn-text", kind: "token" }
  - { name: "--err-text", kind: "token" }
  - { name: "--text-muted", kind: "token" }
---

The defining interaction of agent systems: the agent stops and asks
permission. Everything else on the screen can be scrolled past — **this cannot
be**, so the block is the only one in the library entitled to stop the eye.

```html preview
<div class="inst-approval" data-state="pending" role="group" aria-labelledby="ap1">
  <div class="inst-approval-head" id="ap1">Confirmation is required</div>
  <div class="inst-approval-what">Write 4 files into <code>terrain/</code></div>
  <ul class="inst-approval-effects">
    <li>heightmap.ts — overwrite</li>
    <li data-tone="error">chunks.bin — deletion</li>
  </ul>
  <div class="inst-approval-actions">
    <button class="inst-btn inst-btn--sm inst-btn--primary" type="button">Allow</button>
    <button class="inst-btn inst-btn--sm" type="button">Deny</button>
  </div>
</div>
```

## Contract

Without any one of them the component does not do its work and turns into a
decorated "are you sure?".

| What | The class | Why it is required |
|---|---|---|
| **What exactly will be done** | `inst-approval-what` | "Allow the action?" is not a question: the user does not know what they are answering |
| **What it will touch** | `inst-approval-effects` | The list of consequences. The irreversible ones are marked with the `error` tone |
| **A decision in one press** | `inst-approval-actions` | Allow and deny side by side. The absence of "deny" turns a request into a notification |

### Accessibility

| | |
|---|---|
| The group | A `role="group"` + an `aria-labelledby` pointing at the heading. Otherwise the block falls apart into unrelated paragraphs |
| Appearing | A request arises on the machine's initiative, so it has to reach a live region: an `aria-live="assertive"` or a move of the focus to the block. A request that appeared quietly will wait for ever |
| The order of the actions | "Allow" first. The row is aligned to the start, and the principal action stands at the outer edge of the row — the one signal of which of the two is the principal one |
| Denial is always available | The "Deny" button cannot be hidden in a menu. A refusal has to cost as many presses as consent |
| Not colour alone | An irreversible consequence is marked by a tone **and** a word ("deletion") rather than by redness alone |
| The focus after an answer | The actions disappear along with the answer — the application has to move the focus, otherwise it falls into the `<body>` |

## Anatomy

```html
<ul class="inst-approval-effects">
  <li>heightmap.ts — overwrite</li>
  <li data-tone="warn">config.json — a change of settings</li>
  <li data-tone="error">chunks.bin — deletion</li>
</ul>
```

A tone here means **irreversibility**. A deletion is red because it cannot be
undone — rather than because it is "bad".

## States

```html preview
<div class="inst-approval" data-state="pending" role="group" aria-label="Waiting">
  <div class="inst-approval-head">Confirmation is required</div>
  <div class="inst-approval-what">Write 4 files</div>
  <div class="inst-approval-actions">
    <button class="inst-btn inst-btn--sm inst-btn--primary" type="button">Allow</button>
    <button class="inst-btn inst-btn--sm" type="button">Deny</button>
  </div>
</div>
<div class="inst-approval" data-state="approved" role="group" aria-label="Allowed">
  <div class="inst-approval-head">Confirmation is required</div>
  <div class="inst-approval-what">Write 4 files</div>
  <div class="inst-approval-verdict">Allowed at 14:32</div>
</div>
```

| `data-state` | What happens |
|---|---|
| `pending` | Waiting for an answer. **The only one** where the actions are shown |
| `approved` | Allowed. The block steps back, the actions are hidden |
| `denied` | Denied. The description of the action is struck through |

**After an answer the block neither disappears nor goes out — it steps back.**
The history of decisions stays readable: a user has to be able to come back and
see what exactly they allowed an hour ago. A request that vanished cannot be
checked.

## JS

Include the module once per page — there is no need to initialise the
components one by one, `instrument.js` works by delegation and sees nodes that
arrived later.

```html
<script type="module" src="instrument.js"></script>
```

### What `instrument.js` does

Nothing. Here it honestly has no part: consent to an agent's action is a
decision rather than the behaviour of a widget, and it cannot be faked by
styling.

### What the application has to do

1. Send the decision to where it is carried out.
2. Move the block into a decided state — otherwise the buttons stay live and a
   person presses a second time.
3. Leave visible **what** was decided: a block of consent is a trace in the
   history rather than a dialog that closed.

```js
block.addEventListener('click', async (e) => {
  const btn = e.target.closest('[data-decision]');
  if (!btn) return;

  for (const b of block.querySelectorAll('button')) b.disabled = true;
  await api.decide(block.dataset.id, btn.dataset.decision);

  block.dataset.decided = btn.dataset.decision;
});
```

The buttons are disabled **before** the request. The network takes seconds, and
in that time "Allow" gets pressed twice.

## Patterns

A request from a real run — with the command, the list of consequences and the
mark of the irreversible.

Note the **order**: what will be done, what it will touch, and only then the
decision. A request beginning with the buttons demands an answer before it
states the question.

```html preview context
<div class="inst-approval" data-state="pending" role="group" aria-labelledby="ap2">
  <div class="inst-approval-head" id="ap2">The agent asks permission</div>
  <div class="inst-approval-what">rm -rf build/ &amp;&amp; npm run build</div>
  <ul class="inst-approval-effects">
    <li>It will delete the <code>build/</code> directory entire — 1 284 files</li>
    <li data-tone="warn">Irreversible: the contents do not go to a recycle bin</li>
    <li>The rebuild will take about 40 s</li>
  </ul>
  <div class="inst-approval-actions">
    <button class="inst-btn inst-btn--sm inst-btn--primary" type="button">Allow</button>
    <button class="inst-btn inst-btn--sm" type="button">Deny</button>
    <button class="inst-btn inst-btn--ghost inst-btn--sm" type="button">Always allow</button>
  </div>
</div>
```

## API

```api
```

## Related

```related
```
