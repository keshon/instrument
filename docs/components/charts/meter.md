---
title: Meter
group: Charts
layout: component
source: src/status.css
api:
  - { name: "inst-meter", kind: "class", doc: "The track. The carrier of the role and the values" }
  - { name: "inst-meter-fill", kind: "class", doc: "The fill. A required only child" }
  - { name: "inst-meter-row", kind: "class", doc: "The \"label — quantity\" line above the bar" }
  - { name: "inst-meter-value", kind: "class", doc: "The quantity in that line, in a muted colour" }
  - { name: "data-state", kind: "attribute", doc: "`indeterminate` — the only value a meter has. With no attribute the bar is determinate" }
  - { name: "data-tone", kind: "attribute", doc: "`neutral` `running` `ok` `warn` `error`. Closed. With no attribute, the accent" }
  - { name: "--size-meter", kind: "token" }
  - { name: "--radius-mark", kind: "token" }
  - { name: "--track", kind: "token" }
  - { name: "--hairline", kind: "token" }
  - { name: "--border", kind: "token" }
  - { name: "--tone-mark", kind: "token" }
  - { name: "--accent-mark", kind: "token" }
  - { name: "--text-sm", kind: "token" }
  - { name: "--space-3", kind: "token" }
  - { name: "--gap-inline", kind: "token" }
  - { name: "--dur-3", kind: "token" }
  - { name: "--ease-out", kind: "token" }
---

One quantity against its own ceiling: the spending of a budget, a share of the
time, the readiness of a step. One tone, because **one** quantity is being
compared.

```html preview
<div class="inst-meter-row"><span>Render</span><span class="inst-meter-value">18.2 s</span></div>
<div class="inst-meter" role="progressbar" aria-label="Render"
     aria-valuenow="43" aria-valuemin="0" aria-valuemax="100">
  <div class="inst-meter-fill" style="inline-size:43%"></div>
</div>
```

## Contract

| What | Required | Why |
|---|---|---|
| A `role="progressbar"` on `inst-meter` | yes | With no role the value exists in pixels of width alone and is available neither to a screen reader nor to a test |
| `aria-valuenow` `aria-valuemin` `aria-valuemax` | yes, on a determinate one | A determinate bar has to report a number rather than a picture |
| An `aria-label` or an `aria-labelledby` | yes | A bar has no text of its own. The label from `inst-meter-row` is tied by `aria-labelledby` if it is visible |
| The **absence** of an `aria-valuenow` on an `indeterminate` one | yes | That is exactly how "busy, the duration is unknown" is declared |
| An `inst-meter-fill` inside | yes | The track fills nothing itself |

The `style="inline-size:43%"` on the fill is **the value rather than styling**:
the quantity comes from the server and physically cannot live in static CSS.
Everything else — the height, the radius, the track, the colour — comes from the
library. The fill travels with a transition (`--dur-3`), so an update of the
value reads as motion rather than as a jump.

The extent of the track is shown by the **inner hairline**. The track itself
stays light, and the fill holds 3:1 on it.

### Accessibility

| | |
|---|---|
| The role | `progressbar` — in the example rather than in a note: an `aria-valuenow` with no role is invalid |
| The name | An `aria-label` is always required: a bar is a fill with no text |
| Indeterminacy | The absence of an `aria-valuenow` is announced as "busy". The words "the duration is unknown" beside it repeat that for those who look |
| Colour is not the only carrier | The quantity is repeated as a number in the `inst-meter-value`. A bar whose value is only in its colour and its length is unavailable to somebody colour-blind and unreproducible in print |
| Contrast | The fill holds 3:1 against the track. The track is not load-bearing, so 3:1 against the surface is not asked of it |
| Reduced motion | The sweep of an indeterminate bar **slows to 3 s rather than going out**: a stopped indicator says "nothing is happening" while something is |
| Forced-colours mode | The track gets a `CanvasText` border, and the fill survives the reset (`forced-color-adjust: none`), because it carries a value |
| Print | A sweep on paper is meaningless: the fill of an indeterminate bar is not printed, and an empty track and the label are left |

## Variants

### Determinate and indeterminate

These are **two different components**. A determinate meter promises what is
left, an indeterminate one only the fact of work.

```html preview
<div class="inst-meter" role="progressbar" aria-label="Token spending"
     aria-valuenow="74" aria-valuemin="0" aria-valuemax="100">
  <div class="inst-meter-fill" style="inline-size:74%"></div>
</div>
<div class="inst-meter" data-state="indeterminate" role="progressbar" aria-label="Indexing"
     aria-valuemin="0" aria-valuemax="100">
  <div class="inst-meter-fill"></div>
</div>
```

| | Determinate | Indeterminate |
|---|---|---|
| The attribute | no `data-state` | `data-state="indeterminate"` |
| The width of the fill | comes from the data | set by CSS — 35%, and it sweeps |
| `aria-valuenow` | **required** | **absent, and that is the message** |
| What it promises the reader | what is left and the pace | that work is going on |

The absence of an `aria-valuenow` is not an omission of the markup but its
content: by it a screen reader announces "busy" instead of a false percentage.
Putting a `valuenow` on an indeterminate bar is the same as drawing it stuck at
90%.

### The tone

```html preview
<div class="inst-meter" role="progressbar" aria-label="Spending"
     aria-valuenow="62" aria-valuemin="0" aria-valuemax="100">
  <div class="inst-meter-fill" style="inline-size:62%"></div>
</div>
<div class="inst-meter" data-tone="warn" role="progressbar" aria-label="Spending near the threshold"
     aria-valuenow="88" aria-valuemin="0" aria-valuemax="100">
  <div class="inst-meter-fill" style="inline-size:88%"></div>
</div>
<div class="inst-meter" data-tone="error" role="progressbar" aria-label="The limit is exceeded"
     aria-valuenow="100" aria-valuemin="0" aria-valuemax="100">
  <div class="inst-meter-fill" style="inline-size:100%"></div>
</div>
```

`data-tone` is set on the `inst-meter`, the fill reads `--tone-mark` and knows
nothing else about the tone. A tone here is **an exception rather than a
seasoning**: a bar with no attribute takes the accent, and that is the right
answer in the vast majority of cases. The tone changes when the **meaning of
the quantity** changes — a threshold crossed, a limit spent.

## Composition

### A line with a label and a quantity

```html preview
<div>
  <div class="inst-meter-row"><span>Physics</span><span class="inst-meter-value">11.5 s</span></div>
  <div class="inst-meter" role="progressbar" aria-label="Physics"
       aria-valuenow="27" aria-valuemin="0" aria-valuemax="100">
    <div class="inst-meter-fill" style="inline-size:27%"></div>
  </div>
</div>
<div>
  <div class="inst-meter-row"><span>Indexing</span><span class="inst-meter-value">the duration is unknown</span></div>
  <div class="inst-meter" data-state="indeterminate" role="progressbar" aria-label="Indexing"
     aria-valuemin="0" aria-valuemax="100">
    <div class="inst-meter-fill"></div>
  </div>
</div>
```

`inst-meter-row` sets the label and the quantity at the edges and stands
**above** the bar: a number to the left of the bar would eat its length, and
the length here is the data.

### In a run budget

```html preview context
<div class="inst-budget">
  <div class="inst-budget-row">
    <span>Tokens</span>
    <span class="inst-budget-value">184 320 <span class="inst-budget-of">of 250 000</span></span>
  </div>
  <div class="inst-meter" role="progressbar" aria-label="Token spending"
       aria-valuenow="74" aria-valuemin="0" aria-valuemax="100">
    <div class="inst-meter-fill" style="inline-size:74%"></div>
  </div>
  <div class="inst-budget-note">At the current pace this is enough for about 3 more runs.</div>
</div>
```

## API

```api
```

## Related

```related
```
