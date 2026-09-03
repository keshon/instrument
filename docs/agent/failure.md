---
title: Failure block
group: Agent layer
layout: component
source: src/agent.css
js: Sending the retry and moving the block into a resolved state belong to the application layer
api:
  - { name: "inst-failure", kind: "class", doc: "The block" }
  - { name: "inst-failure-head", kind: "class", doc: "What did not work. With an icon" }
  - { name: "inst-failure-reason", kind: "class", doc: "The reason: a machine code and its explanation" }
  - { name: "inst-failure-tried", kind: "class", doc: "What has already been tried, on a `<ul>`" }
  - { name: "inst-failure-actions", kind: "class", doc: "The way out. **Required**" }
  - { name: "inst-failure-verdict", kind: "class", doc: "What happened and when. Shown after the block steps back" }
  - { name: "data-state", kind: "attribute", value: "pending · resolved · given-up", doc: "on `inst-failure`" }
  - { name: "--pad-card", kind: "token" }
  - { name: "--radius-lg", kind: "token" }
  - { name: "--space-3", kind: "token" }
  - { name: "--err-bg", kind: "token" }
  - { name: "--err-text", kind: "token" }
  - { name: "--hairline", kind: "token" }
  - { name: "--text-sm", kind: "token" }
  - { name: "--text-xs", kind: "token" }
  - { name: "--text-muted", kind: "token" }
---

The agent could not do it. The block answers three questions in a row: **what
did not work**, **what has already been tried** and **what to do next**.

```html preview
<div class="inst-failure" role="alert" data-state="pending">
  <div class="inst-failure-head">Could not read terrain/chunks.bin</div>
  <div class="inst-failure-reason">The file is locked by another process (EBUSY)</div>
  <ul class="inst-failure-tried">
    <li>A retry after 1 s — the same code</li>
    <li>A retry after 4 s — the same code</li>
  </ul>
  <div class="inst-failure-actions">
    <button class="inst-btn inst-btn--sm inst-btn--primary" type="button">Retry</button>
    <button class="inst-btn inst-btn--sm" type="button">Skip the file</button>
  </div>
</div>
```

## Contract

| What | Required | Why |
|---|---|---|
| A `role="alert"` | yes | A failure has to interrupt: it reports that an action was **not** carried out |
| An `inst-failure-actions` with at least one way out | yes | A failure block with no action is [a note](../components/feedback/note.md) in red |
| A machine code in the reason | yes | "Something went wrong" can neither be found in a log nor sent to support |
| An `inst-failure-tried` on a `<ul>` | no, but usually yes | The attempts are enumerable, and their number is announced before the reading |

### Accessibility

| | |
|---|---|
| `role="alert"` | Required: a failure arises on the machine's initiative and has to be announced at once. A block that appeared quietly stays unnoticed |
| Not colour alone | The heading carries an icon **and** a word. A red background by itself reports nothing |
| The reason in text | A machine code (`EBUSY`) comes with an explanation. A code with no explanation is not a reason but its identifier |
| Focus | When the block appears it makes sense to move the focus to the first action: the user has just lost what they were doing |
| Contrast | The heading and the text are checked on `--err-bg` in five themes. The border is a `color-mix` of `--err-text` and is decorative: there is a step between surfaces beside it |
| Repeating failures | Several `role="alert"`s in a row interrupt one another. For a series of failures one block with a counter is better |

## Anatomy

**A failure block with no way out is not a failure block** — it is a red
message there is nothing to answer. A user of an agent system is stuck: the
work has stopped and no decision has been offered.

At least one action. If going on is impossible in principle, the action is
"Cancel the run".

| Take instead | When |
|---|---|
| [A note or a banner](../components/feedback/banner.md) with the `error` tone | A message that does not stop the work |
| [An empty state](../components/feedback/empty.md) | It is empty **not** because of an error |
| [A request for approval](./approval.md) | The agent has not failed but is asking permission |

### "What has already been tried"

```html
<ul class="inst-failure-tried">
  <li>A retry after 1 s — the same code</li>
  <li>A retry after 4 s — the same code</li>
</ul>
```

This is the one thing that sets an agent's report apart from an ordinary error.
Without the list the first thing a user does is press "Retry" — that is, do
what the machine has already done twice.

## States

```html preview
<div class="inst-failure" role="alert" data-state="pending">
  <div class="inst-failure-head">Could not read terrain/chunks.bin</div>
  <div class="inst-failure-reason">The file is locked by another process (EBUSY)</div>
  <div class="inst-failure-actions">
    <button class="inst-btn inst-btn--sm inst-btn--primary" type="button">Retry</button>
    <button class="inst-btn inst-btn--sm" type="button">Skip the file</button>
  </div>
</div>
<div class="inst-failure" role="alert" data-state="resolved">
  <div class="inst-failure-head">Could not read terrain/chunks.bin</div>
  <div class="inst-failure-reason">The file is locked by another process (EBUSY)</div>
  <div class="inst-failure-verdict">Resolved at 14:33, on the third retry</div>
</div>
```

| `data-state` | What happens |
|---|---|
| `pending` | Waiting for a way out. **The only one** where the actions are shown |
| `resolved` | Retried and succeeded. The alert surface comes off, the actions are hidden |
| `given-up` | No further attempt. The alert surface comes off and the reason is struck through |

**Once resolved or given up, the block steps back the same way an approval does
after an answer** — the same law, applied to the component that had not been
given the means to honour it: the alert ground and the accent on `role="alert"`
come off, and a quiet record remains. A failure that vanished when fixed cannot
be checked later against the same file failing again.

## JS

Include the module once per page — there is no need to initialise the
components one by one, `instrument.js` works by delegation and sees nodes that
arrived later.

```html
<script type="module" src="instrument.js"></script>
```

### What `instrument.js` does

Nothing, for the same reason as an approval: retrying is a decision about the
machine's work rather than the behaviour of a widget.

### What the application has to do

1. Send the retry (or the skip) to where it is carried out.
2. Move the block into a resolved state — otherwise the buttons stay live and a
   person presses "Retry" a second time while the first is still in flight.
3. Leave visible **what** happened: a failure block is a trace in the history,
   the same way an answered approval is.

```js
block.addEventListener('click', async (e) => {
  const btn = e.target.closest('[data-decision]');
  if (!btn) return;

  for (const b of block.querySelectorAll('button')) b.disabled = true;
  const ok = await api.retry(block.dataset.id, btn.dataset.decision);

  block.dataset.state = ok ? 'resolved' : 'given-up';
});
```

The buttons are disabled **before** the request, for the same reason as an
approval's: the network takes seconds, and in that time "Retry" gets pressed
twice.

## API

```api
```

## Related

```related
```
