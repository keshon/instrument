---
title: Timeline
group: Data display
layout: component
source: src/data.css
api:
  - { name: "inst-timeline", kind: "class", doc: "The container" }
  - { name: "inst-timeline-item", kind: "class", doc: "One event" }
  - { name: "inst-timeline-dot", kind: "class", doc: "The dot on the line" }
  - { name: "inst-timeline-body", kind: "class", doc: "The description of the event" }
  - { name: "inst-timeline-sub", kind: "class", doc: "A detail under the description" }
  - { name: "inst-timeline-time", kind: "class", doc: "The timestamp, pushed to the end" }
  - { name: "data-tone", kind: "attribute", doc: "`neutral` `running` `ok` `warn` `error`" }
  - { name: "--size-gutter", kind: "token" }
  - { name: "--gap-inline", kind: "token" }
  - { name: "--row-pad-y", kind: "token" }
  - { name: "--size-dot", kind: "token" }
  - { name: "--size-marker", kind: "token" }
  - { name: "--border", kind: "token" }
  - { name: "--text-sm", kind: "token" }
  - { name: "--text-xs", kind: "token" }
  - { name: "--text-muted", kind: "token" }
  - { name: "--tone-mark", kind: "token" }
---

A chronology: what happened and when. The dots are joined by a line, so the
timeline reads as a sequence rather than as a list.

```html preview
<div class="inst-timeline">
  <div class="inst-timeline-item">
    <span class="inst-timeline-dot" data-tone="ok"></span>
    <div class="inst-timeline-body">
      Run finished
      <div class="inst-timeline-sub">4 files changed</div>
    </div>
    <span class="inst-timeline-time">14:36</span>
  </div>
  <div class="inst-timeline-item">
    <span class="inst-timeline-dot" data-tone="running"></span>
    <div class="inst-timeline-body">
      Indexing
      <div class="inst-timeline-sub">terrain/</div>
    </div>
    <span class="inst-timeline-time">14:32</span>
  </div>
</div>
```

## Contract

| What | Required | Why |
|---|---|---|
| `<ol>` and `<li>` | yes | The order carries meaning here, and a screen reader has to announce "a list of N" |
| `<time datetime>` | yes | "14:36" with no date is ambiguous outside the context of a day |
| The tone is repeated in words in the description | yes | A dot by itself says nothing to somebody who does not tell tones apart |
| `data-tone` | no | Without the attribute the dot is neutral |

The line is drawn by a pseudo-element and does not reach the speaking — it is
decorative.

### Accessibility

| | |
|---|---|
| This is a list | Wrap it in an `<ol>` and an event in an `<li>`: the order carries meaning here, and a screen reader has to announce "a list of N" |
| Time | A `<time datetime="…">` with a machine value. "14:36" with no date is ambiguous outside the context of a day |
| Not colour alone | The tone of a dot is repeated in words in the description of the event. A dot by itself says nothing to somebody who does not tell tones apart |
| The line | Drawn by a pseudo-element and does not reach the speaking — which is right, it is decorative |
| Reduced motion | The pulse of a `running` dot slows down rather than going out |

## Variants

A dot takes the same `data-tone` as the whole library.

```html preview
<div class="inst-timeline">
  <div class="inst-timeline-item">
    <span class="inst-timeline-dot"></span>
    <div class="inst-timeline-body">An ordinary event</div>
  </div>
  <div class="inst-timeline-item">
    <span class="inst-timeline-dot" data-tone="running"></span>
    <div class="inst-timeline-body">Happening now</div>
  </div>
  <div class="inst-timeline-item">
    <span class="inst-timeline-dot" data-tone="ok"></span>
    <div class="inst-timeline-body">Finished successfully</div>
  </div>
  <div class="inst-timeline-item">
    <span class="inst-timeline-dot" data-tone="warn"></span>
    <div class="inst-timeline-body">With remarks</div>
  </div>
  <div class="inst-timeline-item">
    <span class="inst-timeline-dot" data-tone="error"></span>
    <div class="inst-timeline-body">A refusal</div>
  </div>
</div>
```

| Tone | When |
|---|---|
| no attribute | An ordinary event |
| `data-tone="running"` | Happening now. The dot pulses |
| `data-tone="ok"` | A successful finish |
| `data-tone="warn"` | With remarks |
| `data-tone="error"` | A refusal |

## API

```api
```

## Related

```related
```
