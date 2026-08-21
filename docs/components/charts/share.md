---
title: Share
group: Charts
layout: component
source: src/data.css
api:
  - { name: "inst-share", kind: "class", doc: "The bar. The carrier of the role and the name" }
  - { name: "inst-share-part", kind: "class", doc: "One outcome. Its width is its share of the whole" }
  - { name: "inst-share--lg", kind: "class", doc: "The height of a glyph, for a bar that is the subject of the screen rather than an annotation to a number" }
  - { name: "--share", kind: "token", doc: "The share of a part, 0 to 1. It is **the value rather than styling** and comes from the data" }
  - { name: "data-tone", kind: "attribute", value: "neutral · running · ok · warn · error", doc: "The outcome of a part. Closed. With no attribute, the accent" }
  - { name: "role", kind: "attribute", value: "img", doc: "The bar is a picture of data, and it has to have a name" }
  - { name: "--size-meter", kind: "token" }
  - { name: "--size-icon", kind: "token" }
  - { name: "--radius-mark", kind: "token" }
  - { name: "--track", kind: "token" }
  - { name: "--hairline", kind: "token" }
  - { name: "--border", kind: "token" }
  - { name: "--tone-mark", kind: "token" }
  - { name: "--accent-mark", kind: "token" }
---

What a period was made **of**: twenty-two hours up, an hour and a half down,
half an hour nobody looked. One bar, divided by the share of each outcome in it.

```html preview
<div class="inst-share" role="img"
     aria-label="The last 24 hours: 92% up, 6% down, 2% not checked">
  <span class="inst-share-part" data-tone="ok" style="--share:0.92"></span>
  <span class="inst-share-part" data-tone="error" style="--share:0.06"></span>
  <span class="inst-share-part" data-tone="neutral" style="--share:0.02"></span>
</div>
```

## Contract

| What | Required | Why |
|---|---|---|
| A `role="img"` and an `aria-label` with the shares **in words** | yes | Three empty `<span>`s are three empty `<span>`s to a screen reader. The reading has to be text: "92% up, 6% down" |
| A `--share` on every part | yes | With no share a part gets a grow factor of 0 and is not drawn. The number is the data and comes from the server |
| A tone from the vocabulary | yes | `neutral` · `running` · `ok` · `warn` · `error`. There is no sixth outcome |
| Shares that sum to **1** | no | What is not accounted for stays track, and a bar a third full is the honest drawing of eight hours of data out of twenty-four |
| A minimum width for a small share | **no, and it is refused** | A floor would draw 0.02% as wide as 2%, and the bar would stop reporting proportion in the units of a bar that claims to report proportion |

The `style="--share:0.92"` is **the value rather than styling**, by the same
rule as the width of a meter fill: the quantity comes from the server and
physically cannot live in static CSS. Everything else — the height, the corner,
the track, the colour — comes from the library.

The extent is drawn by the **inner hairline**, as on a meter: the track stays
light enough for a tone on it to hold 3:1.

The parts are **square inside a clipped strip**. Rounding each part would put an
arc at every internal boundary — four corners where the data has one edge — and
at `--size-meter` those arcs meet and leave a notch of track showing between two
neighbouring outcomes.

### Accessibility

| | |
|---|---|
| The name | A `role="img"` plus an `aria-label` with the shares in words. The parts are empty and are not read; what is read is the name of the bar |
| Colour is not the only carrier | The bar stands beside the numbers, in a legend or in a line of text. On its own it says "mostly one thing"; which thing is a matter of the label |
| Contrast | A part takes `--tone-mark` — a 3:1 threshold against the surface, the same as a state mark. The track is not load-bearing and is not asked for 3:1 |
| Motion | There is none. A share is a state rather than a process, and the bar does not animate as the data updates |
| Forced-colours mode | The strip gets a `CanvasText` border and the parts survive the reset, because they carry a value. The tones collapse into one fill there, so what stays readable is the **division**: each part but the last draws a line on its trailing edge |
| Print | It is printed as it is: the parts are told apart by tone and stay grey steps on paper |

## Sizes

```html preview
<div class="inst-share" role="img" aria-label="Today: 88% up, 9% degraded, 3% down">
  <span class="inst-share-part" data-tone="ok" style="--share:0.88"></span>
  <span class="inst-share-part" data-tone="warn" style="--share:0.09"></span>
  <span class="inst-share-part" data-tone="error" style="--share:0.03"></span>
</div>
<div class="inst-share inst-share--lg" role="img" aria-label="Today: 88% up, 9% degraded, 3% down">
  <span class="inst-share-part" data-tone="ok" style="--share:0.88"></span>
  <span class="inst-share-part" data-tone="warn" style="--share:0.09"></span>
  <span class="inst-share-part" data-tone="error" style="--share:0.03"></span>
</div>
```

Two, for the two places the bar stands. The base is the meter's height, because
in a row of data this is the meter's job with more than one answer. The large
one is a glyph's height: on the screen **of** a single object the bar is the
subject rather than an annotation to a number.

## Composition

### With a legend

```html preview
<div class="inst-stack">
  <div class="inst-share" role="img"
       aria-label="The last 24 hours: 92% up, 6% down, 2% not checked">
    <span class="inst-share-part" data-tone="ok" style="--share:0.92"></span>
    <span class="inst-share-part" data-tone="error" style="--share:0.06"></span>
    <span class="inst-share-part" data-tone="neutral" style="--share:0.02"></span>
  </div>
  <ul class="inst-legend">
    <li class="inst-legend-item" data-tone="ok"><span class="inst-legend-swatch"></span>up <span class="inst-legend-value">92%</span></li>
    <li class="inst-legend-item" data-tone="error"><span class="inst-legend-swatch"></span>down <span class="inst-legend-value">6%</span></li>
    <li class="inst-legend-item"><span class="inst-legend-swatch" data-empty></span>not checked <span class="inst-legend-value">2%</span></li>
  </ul>
</div>
```

The swatch reads `--tone-mark` when no `--series` is set, so **one legend keys
both** a chart of series and a bar of states. An hour nobody checked is neither:
it takes `data-empty`, the same word `inst-history-tick` uses for it.

The same statement, drawn differently on purpose. On a strip an absence is drawn
**by** absence — no mark at all — and that works because a tick has neighbours: a
hole between two coloured bricks reads as a hole. A swatch has none, and one that
paints nothing is a missing swatch beside a stray word rather than an entry about
absence. So in the legend it takes the colour the strip's ground has at that
hour: the track.

### Above a history strip

```html preview context
<div class="inst-stack">
  <div class="inst-share" role="img" aria-label="The last 12 hours: 92% up, 8% down">
    <span class="inst-share-part" data-tone="ok" style="--share:0.92"></span>
    <span class="inst-share-part" data-tone="error" style="--share:0.08"></span>
  </div>
  <div class="inst-history" role="img" aria-label="12 hours, 4 checks each: 44 succeeded, 4 with an error">
    <span class="inst-history-group" style="--n:4"><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span></span>
    <span class="inst-history-group" style="--n:4"><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span></span>
    <span class="inst-history-group" style="--n:4"><span class="inst-history-tick" data-tone="error"></span><span class="inst-history-tick" data-tone="error"></span><span class="inst-history-tick" data-tone="error"></span><span class="inst-history-tick" data-tone="error"></span></span>
    <span class="inst-history-group" style="--n:4"><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span></span>
    <span class="inst-history-group" style="--n:4"><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span></span>
    <span class="inst-history-group" style="--n:4"><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span></span>
    <span class="inst-history-group" style="--n:4"><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span></span>
    <span class="inst-history-group" style="--n:4"><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span></span>
    <span class="inst-history-group" style="--n:4"><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span></span>
    <span class="inst-history-group" style="--n:4"><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span></span>
    <span class="inst-history-group" style="--n:4"><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span></span>
    <span class="inst-history-group" style="--n:4"><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span></span>
  </div>
</div>
```

The two answer different questions and an uptime screen wants both. A strip
answers **when** — the failure was at 03:00, there were four of them, they were
next to each other. A share answers **how much** — the day was 8% down. A strip
cannot total, and a share cannot order: the same 8% is one long outage or ninety
scattered ones, and this bar draws them identically.

## API

```api
```

## Related

```related
```
