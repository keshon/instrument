---
title: Run lanes
group: Agent layer
layout: component
source: src/agent.css
api:
  - { name: "inst-lanes", kind: "class", doc: "The set of lanes" }
  - { name: "inst-lane", kind: "class", doc: "One lane: a label and an axis" }
  - { name: "inst-lane-label", kind: "class", doc: "The name of the worker; it is truncated" }
  - { name: "inst-lane-track", kind: "class", doc: "The axis of time" }
  - { name: "inst-lane-span", kind: "class", doc: "A segment of work" }
  - { name: "data-tone", kind: "attribute", doc: "`neutral` `running` `ok` `warn` `error`" }
  - { name: "--label-col", kind: "token" }
  - { name: "--gap-inline", kind: "token" }
  - { name: "--size-meter", kind: "token" }
  - { name: "--radius-mark", kind: "token" }
  - { name: "--track", kind: "token" }
  - { name: "--border", kind: "token" }
  - { name: "--hairline", kind: "token" }
  - { name: "--text-xs", kind: "token" }
  - { name: "--tone-mark", kind: "token" }
---

Several agents working in parallel on a shared axis of time. It answers the
question neither the queue nor the log answers: **what ran at the same time and
what waited for what**.

```html preview
<div class="inst-lanes">
  <div class="inst-lane">
    <span class="inst-lane-label">worldgen-01</span>
    <span class="inst-lane-track">
      <span class="inst-lane-span" data-tone="ok" style="inset-inline-start:0%;inline-size:34%"></span>
    </span>
  </div>
  <div class="inst-lane">
    <span class="inst-lane-label">biomes-04</span>
    <span class="inst-lane-track">
      <span class="inst-lane-span" data-tone="running" style="inset-inline-start:30%;inline-size:52%"></span>
    </span>
  </div>
</div>
```

## Contract

| What | Required | Why |
|---|---|---|
| `inset-inline-start` and `inline-size` inline | yes | The bounds of a segment are **data** rather than styling: the start and the duration come from the server |
| Logical properties rather than `left` and `width` | yes | In RTL the axis of time mirrors itself |
| A label or a `title` on a segment | yes | A rectangle with no name and no time says nothing |
| `data-tone` | no | With no attribute a segment is neutral |

### Accessibility

| | |
|---|---|
| A segment is not text | By itself it is spoken by nothing. A lane has to carry an `aria-label` with the duration: "worldgen-01, from 0 to 34 per cent, 12 seconds" |
| The label is truncated | The full name goes in a `title`, otherwise the worker has no name |
| Not colour alone | The tone of a segment is repeated in the accessible name in words |
| Forced-colours mode | A segment carries meaning, so it survives the reset: `forced-color-adjust: none` and `Highlight` |
| The alternative | On a narrow screen lanes are unreadable. The fallback view is [a table](../components/display/table.md) with the start, the end and the duration |
| Reduced motion | A `running` segment is not animated — the duration changes in jumps of data anyway |

## Anatomy

`inset-inline-start` and `inline-size` are set **inline**, and that is the one
case of an inline style allowed in the library: the position of a segment on
the axis is a value rather than styling. The same channel carries the width of
a meter's fill and the depth of a tree node.

```html
<span class="inst-lane-span" style="inset-inline-start:30%;inline-size:52%"></span>
```

The percentages are counted from the total duration of the run, so every lane
shares one scale — otherwise comparing them is meaningless.

## API

```api
```

## Related

```related
```
