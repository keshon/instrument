---
title: Thinking
group: Agent layer
layout: component
source: src/agent.css
api:
  - { name: "inst-thinking", kind: "class", doc: "The disclosure. A `<details>`, collapsed by default" }
  - { name: "inst-thinking-head", kind: "class", doc: "The summary line: how long it took" }
  - { name: "inst-thinking-body", kind: "class", doc: "The reasoning, set as an aside" }
  - { name: "inst-step-twist", kind: "class", doc: "The chevron. The kit has one disclosure glyph and it already has a name" }
  - { name: "--text-xs", kind: "token" }
  - { name: "--text-muted", kind: "token" }
  - { name: "--text-secondary", kind: "token" }
  - { name: "--leading-ui", kind: "token" }
  - { name: "--hairline", kind: "token" }
  - { name: "--border", kind: "token" }
  - { name: "--pad-panel", kind: "token" }
  - { name: "--space-3", kind: "token" }
---

What the machine worked out before it answered. Collapsed by default: the
answer is what was asked for.

```html preview
<details class="inst-thinking">
  <summary class="inst-thinking-head"><span class="inst-step-twist"></span>Thought for 4 s</summary>
  <div class="inst-thinking-body">The strip answers when, the share answers how much. One legend can key both.</div>
</details>
```

## Contract

| What | Required | Why |
|---|---|---|
| A `<details>` with a `<summary>` | yes | Disclosure, the keyboard and the role all come from the platform |
| Collapsed by default | yes | Open, it puts the working before the answer on every turn |
| A duration in the summary | no | It is the one fact about thinking that is worth a line when it is shut |
| A tone or a dot | **no** | Nothing here can fail |

**Not a step.** A step is a call — a name, arguments, a result, an outcome, and
a dot reporting how it ended. Thinking has none of those, and dressed as a step
it would promise a status it does not have.

## API

```api
```

## Related

```related
```
