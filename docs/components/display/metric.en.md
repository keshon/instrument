---
title: Metric
group: Data display
layout: component
source: src/surfaces.css
api:
  - { name: "inst-metric", kind: "class", doc: "One metric" }
  - { name: "inst-metric-row", kind: "class", doc: "A row of metrics" }
  - { name: "inst-metric-label", kind: "class", doc: "The label above the number" }
  - { name: "inst-metric-value", kind: "class", doc: "The number itself" }
  - { name: "inst-metric-unit", kind: "class", doc: "The unit inside the number: smaller and quieter" }
  - { name: "inst-metric-delta", kind: "class", doc: "The change against the previous period" }
  - { name: "data-dir", kind: "attribute", value: "up · down", doc: "on `inst-metric-delta`" }
  - { name: "data-tone", kind: "attribute", value: "neutral · running · ok · warn · error", doc: "in the same place" }
  - { name: "--surface-sunken", kind: "token" }
  - { name: "--radius-md", kind: "token" }
  - { name: "--pad-panel", kind: "token" }
  - { name: "--text-2xl", kind: "token" }
  - { name: "--text-xs", kind: "token" }
  - { name: "--tracking-tight", kind: "token" }
  - { name: "--tone-ink", kind: "token" }
  - { name: "--size-chevron", kind: "token" }
  - { name: "--space-2", kind: "token" }
---

One number people watch. The unit is set smaller and quieter than the number
itself: "42 s" at one size reads as one word, and once size and colour are
parted the eye grabs the number and the "s" recedes.

```html preview
<div class="inst-metric-row">
  <div class="inst-metric">
    <div class="inst-metric-label">Runs per day</div>
    <div class="inst-metric-value">128</div>
    <div class="inst-metric-delta" data-dir="down" data-tone="ok">18% against the last</div>
  </div>
  <div class="inst-metric">
    <div class="inst-metric-label">Average time</div>
    <div class="inst-metric-value">4.2<span class="inst-metric-unit">s</span></div>
  </div>
  <div class="inst-metric">
    <div class="inst-metric-label">Warnings</div>
    <div class="inst-metric-value">5</div>
    <div class="inst-metric-delta" data-dir="up" data-tone="error">5 new</div>
  </div>
</div>
```

## Contract

| What | Required | Why |
|---|---|---|
| `inst-metric-label` | yes | "128" with no answer to "of what" is not a metric |
| The sign of the change in words in the delta | yes | The arrow is a second carrier rather than the only one: "↓ 18%" reads without a tone too |
| An accessible name from the label and the number | yes, if the metric is clickable | Otherwise "128" is spoken with no answer to "of what" |
| `data-dir` and `data-tone` | no | Without the first the delta goes with no arrow, without the second in a neutral colour |

### Accessibility

| | |
|---|---|
| The label and the number | Tied visually by their order. If the metric is clickable, the accessible name has to include both — otherwise "128" is spoken with no answer to "of what" |
| Not colour alone | The delta carries an arrow **and** the sign of the change in the text. A reader who does not tell tones apart reads "↓ 18%" |
| Contrast | The delta takes `--tone-ink` — the threshold of text at 4.5:1 rather than of a mark at 3:1 |
| The size of the number | `--text-2xl` is meant for a hero number **only**. A heading of that size shouts louder than the data |
| Tabular figures | On by default: a number updating in place does not tug at its neighbours |

## Anatomy

### Direction and judgement are different attributes

This is the main trap of the component. An arrow down on the time of a run is
good, an arrow up on warnings is bad. An attribute named after the arrow and
meaning the judgement is a guaranteed mistake of application.

| Attribute | What it encodes | Values |
|---|---|---|
| `data-dir` | **Which way** it changed | `up` · `down` |
| `data-tone` | **Whether that is good or bad** | `ok` · `warn` · `error` · `neutral` · `running` |

```html preview
<div class="inst-metric-row">
  <div class="inst-metric">
    <div class="inst-metric-label">Run time</div>
    <div class="inst-metric-value">42<span class="inst-metric-unit">&nbsp;s</span></div>
    <div class="inst-metric-delta" data-dir="down" data-tone="ok">18% — it fell, and that is good</div>
  </div>
  <div class="inst-metric">
    <div class="inst-metric-label">Warnings</div>
    <div class="inst-metric-value">12</div>
    <div class="inst-metric-delta" data-dir="up" data-tone="error">5 new — it grew, and that is bad</div>
  </div>
</div>
```

Both attributes are optional: a delta with no `data-dir` goes with no arrow,
with no `data-tone` in a neutral colour.

## Composition

### A row of metrics — with no borders

```html preview
<div class="inst-metric-row">
  <div class="inst-metric">
    <div class="inst-metric-label">Total</div>
    <div class="inst-metric-value">128</div>
  </div>
  <div class="inst-metric">
    <div class="inst-metric-label">Succeeded</div>
    <div class="inst-metric-value">121</div>
  </div>
  <div class="inst-metric">
    <div class="inst-metric-label">Failed</div>
    <div class="inst-metric-value">7</div>
  </div>
</div>
```

Four numbers are **one** group. A border around each turns them into four
objects, between which the reader starts looking for differences that are not
there. So a metric stands on a muted surface rather than in a card.

## Patterns

### The summary of a run

On its own a metric is just a large number; the meaning appears in a row.

What is visible only here:

- **none of them has a border.** Four numbers are one group; a border around
  each would turn them into four objects, between which the reader starts
  looking for differences;
- **the unit is smaller and quieter than the number.** "42 s" at one size reads
  as one word; once size and colour are parted the eye grabs the number and the
  "s" recedes;
- **direction and judgement are independent.** An arrow down on the time is
  good, an arrow up on warnings is bad.

```html preview context
<div class="inst-metric-row">
  <div class="inst-metric">
    <div class="inst-metric-label">Run time</div>
    <div class="inst-metric-value">42<span class="inst-metric-unit">&nbsp;s</span></div>
    <div class="inst-metric-delta" data-dir="down" data-tone="ok">18% against the last</div>
  </div>
  <div class="inst-metric">
    <div class="inst-metric-label">Tests passed</div>
    <div class="inst-metric-value">248<span class="inst-metric-unit">/251</span></div>
    <div class="inst-metric-delta" data-tone="warn">3 failing</div>
  </div>
  <div class="inst-metric">
    <div class="inst-metric-label">Bundle size</div>
    <div class="inst-metric-value">7.4<span class="inst-metric-unit">&nbsp;MB</span></div>
    <div class="inst-metric-delta">unchanged</div>
  </div>
  <div class="inst-metric">
    <div class="inst-metric-label">Warnings</div>
    <div class="inst-metric-value">12</div>
    <div class="inst-metric-delta" data-dir="up" data-tone="error">5 new</div>
  </div>
</div>
```

## API

```api
```

## Related

```related
```
