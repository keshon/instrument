---
title: Failure block
group: Agent layer
layout: component
source: src/agent.css
api:
  - { name: "inst-failure", kind: "class", doc: "The block" }
  - { name: "inst-failure-head", kind: "class", doc: "What did not work. With an icon" }
  - { name: "inst-failure-reason", kind: "class", doc: "The reason: a machine code and its explanation" }
  - { name: "inst-failure-tried", kind: "class", doc: "What has already been tried, on a `<ul>`" }
  - { name: "inst-failure-actions", kind: "class", doc: "The way out. **Required**" }
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
<div class="inst-failure" role="alert">
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

## API

```api
```

## Related

```related
```
