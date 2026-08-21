---
title: Check history
group: Agent layer
layout: component
source: src/agent.css
api:
  - { name: "inst-history", kind: "class", doc: "The strip. The fresh outcomes are at the right" }
  - { name: "inst-history-tick", kind: "class", doc: "One outcome: one check, one run" }
  - { name: "data-tone", kind: "attribute", value: "ok · warn · error · running · neutral", doc: "The outcome of a tick. With no attribute the tick is grey: the check happened, the result is unknown" }
  - { name: "data-empty", kind: "attribute", doc: "There was no check. The room stays, the mark does not" }
  - { name: "inst-history-group", kind: "class", doc: "A batch of ticks with a shared name: an hour, a day, a run. It anchors its own tooltip and answers to hover" }
  - { name: "--n", kind: "token", doc: "The number of ticks in a group. It sets its share of the strip, and the same number sets the share of the axis cell under it" }
  - { name: "inst-history--sm", kind: "class", doc: "The height of a meter, for a strip that is a column in a table of objects" }
  - { name: "inst-history--lg", kind: "class", doc: "The height of a small control, for a strip that is the subject of the screen" }
  - { name: "inst-history-axis", kind: "class", doc: "The ruler under the strip. It repeats the geometry of the groups rather than computing positions" }
  - { name: "inst-history-axis-cell", kind: "class", doc: "One label. It takes the same `--n` as the group above it" }
  - { name: "data-minor", kind: "attribute", doc: "On an axis cell: this label may be dropped when the ruler is narrower than 36rem. The cell keeps its place" }
  - { name: "role", kind: "attribute", value: "img", doc: "The strip is a picture of data, and it has to have a name" }
  - { name: "--size-tick", kind: "token" }
  - { name: "--space-1", kind: "token" }
  - { name: "--space-3", kind: "token" }
  - { name: "--radius-mark", kind: "token" }
  - { name: "--track", kind: "token" }
  - { name: "--size-icon", kind: "token" }
  - { name: "--size-meter", kind: "token" }
  - { name: "--control-h-sm", kind: "token" }
  - { name: "--space-2", kind: "token" }
  - { name: "--text-2xs", kind: "token" }
  - { name: "--text-muted", kind: "token" }
  - { name: "--hairline", kind: "token" }
---

A row of discrete outcomes: one check, one tick, the fresh ones at the right.
The uptime of a monitor, the runs of a nightly build, the attempts of a retry.

```html preview
<div class="inst-history" role="img" aria-label="The last 24 checks: 22 succeeded, 2 with an error">
  <span class="inst-history-tick" data-tone="ok"></span>
  <span class="inst-history-tick" data-tone="ok"></span>
  <span class="inst-history-tick" data-tone="error"></span>
  <span class="inst-history-tick" data-tone="ok"></span>
  <span class="inst-history-tick" data-tone="ok"></span>
  <span class="inst-history-tick" data-tone="ok"></span>
  <span class="inst-history-tick" data-tone="ok"></span>
  <span class="inst-history-tick" data-tone="warn"></span>
  <span class="inst-history-tick" data-tone="ok"></span>
  <span class="inst-history-tick" data-tone="ok"></span>
  <span class="inst-history-tick" data-tone="ok"></span>
  <span class="inst-history-tick" data-tone="ok"></span>
</div>
```

## Contract

| What | Required | Why |
|---|---|---|
| A `role="img"` and an `aria-label` with the **outcome in words** | yes | Two hundred empty `<span>`s are two hundred empty `<span>`s to a screen reader. The outcome has to be text: "22 of 24 succeeded" |
| The fresh ones at the **right** | yes | The strip is read as time, and time in this script runs to the right |
| A tone from the vocabulary | yes | `ok` · `warn` · `error` · `running` · `neutral`. There is no sixth outcome |
| A tick with no tone | no | The check happened and the result is unknown: the tick stays the colour of the track |
| A `tabindex` on a group | no | It makes the group a target, and twenty-four hours in a column give it 15px against the 24 WCAG 2.2 AA asks. The reading reaches a screen reader through the strip's own `aria-label` |

### Accessibility

| | |
|---|---|
| The name | A `role="img"` plus an `aria-label` with the outcome in words. The ticks are empty and are not read by a screen reader — what is read is the name of the strip |
| Colour is not the only carrier | The strip stands beside a status in words. By itself it is a refinement: "it failed twice in a day" rather than "it is down now" |
| Contrast | The tone takes `--tone-mark` — a 3:1 threshold against the surface, the same as a state mark |
| Motion | There is no animation: a history is a state rather than a process |
| Print | It is printed as it is: the ticks are told apart by tone and stay grey steps on paper |

## Anatomy

The width of a tick is fixed (`--size-tick`) rather than a share. A share would
make a strip of twenty checks indistinguishable from a strip of two hundred,
though they mean different things: a short history and a long one.

The strip is pushed to the end, so what is extra goes past the **start** and is
clipped there: the old disappears first, and the last check is always visible.
The markup meanwhile is ordinary — from old to fresh, as in the array the
server hands over.

The height is `--size-icon`, as tall as a glyph. The strip stands in one row
with a label and a badge, and at the full height of a line it stops being a
mark beside them: it becomes a block and pulls the row onto itself.

## Variants

### Groups

An hour, a day, a run: a batch of checks with a shared name.

A strip with no groups answers "what happened". With groups it also answers
"how much stands behind that": an hour with one check is one brick, an hour
with seven is seven.

```html preview
<div class="inst-history" style="inline-size:26rem" role="img" aria-label="Six hours: 19 checks, two with an error">
  <span class="inst-history-group" style="--n:1"><span class="inst-history-tick" data-tone="ok"></span></span>
  <span class="inst-history-group" style="--n:7"
    ><span class="inst-history-tick" data-tone="ok"></span
    ><span class="inst-history-tick" data-tone="ok"></span
    ><span class="inst-history-tick" data-tone="error"></span
    ><span class="inst-history-tick" data-tone="error"></span
    ><span class="inst-history-tick" data-tone="ok"></span
    ><span class="inst-history-tick" data-tone="ok"></span
    ><span class="inst-history-tick" data-tone="ok"></span
  ></span>
  <span class="inst-history-group" style="--n:1"><span class="inst-history-tick" data-empty></span></span>
  <span class="inst-history-group" style="--n:4"
    ><span class="inst-history-tick" data-tone="ok"></span
    ><span class="inst-history-tick" data-tone="ok"></span
    ><span class="inst-history-tick" data-tone="ok"></span
    ><span class="inst-history-tick" data-tone="ok"></span
  ></span>
  <span class="inst-history-group" style="--n:6"
    ><span class="inst-history-tick" data-tone="ok"></span
    ><span class="inst-history-tick" data-tone="ok"></span
    ><span class="inst-history-tick" data-tone="ok"></span
    ><span class="inst-history-tick" data-tone="ok"></span
    ><span class="inst-history-tick" data-tone="ok"></span
    ><span class="inst-history-tick" data-tone="ok"></span
  ></span>
</div>
```

The width of a group is **proportional** to the number of ticks in it — the
`--n` in the markup. The width of the strip is given by the container: the
groups divide what they were given.

Equal slots will not do, and that is measurable: a strip of 520px split over 24
hours gives 14px to an hour, and three ticks fit into 14px. A proportion, on
the contrary, always takes exactly the width of the strip however many checks
there are — a dense monitor looks dense and a sparse one sparse.

A strip with groups **does not clip** its content: the clipping for "the old
goes first" is not wanted here, and it is what makes it possible to hang
[a tooltip](../components/overlays/tooltip.md) on a group — under an
`overflow: hidden` it would be cut off.

| With no groups | With groups |
|---|---|
| The width of a tick is fixed — the length of the strip says how many checks there were | The width of the strip is fixed — the length of a group says how many checks are in that hour |
| What is extra is clipped at the start | It is not clipped: the shares are counted from what there is |
| One row of a list beside another | One object at large, on a screen of its own |

### The aggregate form

Above roughly four checks an hour a tick is too thin to see, let alone point
at. Give the group **one** tick and leave `--n` alone: the hour keeps its exact
share and comes back as a single brick, tinted by its worst outcome.

```html preview
<div class="inst-history" style="inline-size:26rem" role="img" aria-label="Twelve hours, 7 checks each: two hours with errors">
  <span class="inst-history-group" style="--n:7"><span class="inst-history-tick" data-tone="ok"></span></span>
  <span class="inst-history-group" style="--n:7"><span class="inst-history-tick" data-tone="ok"></span></span>
  <span class="inst-history-group" style="--n:7"><span class="inst-history-tick" data-tone="error"></span></span>
  <span class="inst-history-group" style="--n:7"><span class="inst-history-tick" data-tone="warn"></span></span>
  <span class="inst-history-group" style="--n:7"><span class="inst-history-tick" data-tone="ok"></span></span>
  <span class="inst-history-group" style="--n:1"><span class="inst-history-tick" data-tone="ok"></span></span>
  <span class="inst-history-group" style="--n:7"><span class="inst-history-tick" data-tone="ok"></span></span>
  <span class="inst-history-group" style="--n:7"><span class="inst-history-tick" data-tone="ok"></span></span>
  <span class="inst-history-group" style="--n:12"><span class="inst-history-tick" data-tone="ok"></span></span>
  <span class="inst-history-group" style="--n:7"><span class="inst-history-tick" data-tone="ok"></span></span>
  <span class="inst-history-group" style="--n:7"><span class="inst-history-tick" data-tone="ok"></span></span>
  <span class="inst-history-group" style="--n:7"><span class="inst-history-tick" data-tone="ok"></span></span>
</div>
```

`--n` is what keeps the hour's weight: the sixth had one check and is narrow,
the ninth had twelve and is wide.

## Sizes

```html preview
<div class="inst-stack" style="align-items:start">
  <div class="inst-history inst-history--sm" role="img" aria-label="Small: eight checks, one with an error">
    <span class="inst-history-tick" data-tone="ok"></span>
    <span class="inst-history-tick" data-tone="ok"></span>
    <span class="inst-history-tick" data-tone="error"></span>
    <span class="inst-history-tick" data-tone="ok"></span>
    <span class="inst-history-tick" data-tone="ok"></span>
    <span class="inst-history-tick" data-tone="ok"></span>
    <span class="inst-history-tick" data-tone="ok"></span>
    <span class="inst-history-tick" data-tone="ok"></span>
  </div>
  <div class="inst-history" role="img" aria-label="Base: eight checks, one with an error">
    <span class="inst-history-tick" data-tone="ok"></span>
    <span class="inst-history-tick" data-tone="ok"></span>
    <span class="inst-history-tick" data-tone="error"></span>
    <span class="inst-history-tick" data-tone="ok"></span>
    <span class="inst-history-tick" data-tone="ok"></span>
    <span class="inst-history-tick" data-tone="ok"></span>
    <span class="inst-history-tick" data-tone="ok"></span>
    <span class="inst-history-tick" data-tone="ok"></span>
  </div>
  <div class="inst-history inst-history--lg" role="img" aria-label="Large: eight checks, one with an error">
    <span class="inst-history-tick" data-tone="ok"></span>
    <span class="inst-history-tick" data-tone="ok"></span>
    <span class="inst-history-tick" data-tone="error"></span>
    <span class="inst-history-tick" data-tone="ok"></span>
    <span class="inst-history-tick" data-tone="ok"></span>
    <span class="inst-history-tick" data-tone="ok"></span>
    <span class="inst-history-tick" data-tone="ok"></span>
    <span class="inst-history-tick" data-tone="ok"></span>
  </div>
</div>
```

| | Height | Where |
|---|---|---|
| `inst-history--sm` | `--size-meter` | A column in a table of objects. At glyph height a strip out-shouts the names it belongs to |
| base | `--size-icon` | In a row beside a name and a badge, where the strip is a mark among marks |
| `inst-history--lg` | `--control-h-sm` | On the screen **of** one object, where the strip is the subject rather than an annotation |

## States

```html preview
<div class="inst-stack" style="align-items:start">
  <div class="inst-history" role="img" aria-label="Everything succeeded">
    <span class="inst-history-tick" data-tone="ok"></span>
    <span class="inst-history-tick" data-tone="ok"></span>
    <span class="inst-history-tick" data-tone="ok"></span>
    <span class="inst-history-tick" data-tone="ok"></span>
    <span class="inst-history-tick" data-tone="ok"></span>
    <span class="inst-history-tick" data-tone="ok"></span>
  </div>
  <div class="inst-history" role="img" aria-label="It failed twice and is working now">
    <span class="inst-history-tick" data-tone="ok"></span>
    <span class="inst-history-tick" data-tone="error"></span>
    <span class="inst-history-tick" data-tone="error"></span>
    <span class="inst-history-tick" data-tone="ok"></span>
    <span class="inst-history-tick" data-tone="ok"></span>
    <span class="inst-history-tick" data-tone="ok"></span>
  </div>
  <div class="inst-history" role="img" aria-label="There have been no checks yet">
    <span class="inst-history-tick"></span>
    <span class="inst-history-tick"></span>
    <span class="inst-history-tick"></span>
    <span class="inst-history-tick"></span>
    <span class="inst-history-tick"></span>
    <span class="inst-history-tick"></span>
  </div>
  <div class="inst-history" role="img" aria-label="Two hours with no checks">
    <span class="inst-history-tick" data-tone="ok"></span>
    <span class="inst-history-tick" data-tone="ok"></span>
    <span class="inst-history-tick" data-empty></span>
    <span class="inst-history-tick" data-empty></span>
    <span class="inst-history-tick" data-tone="ok"></span>
    <span class="inst-history-tick" data-tone="error"></span>
  </div>
</div>
```

| State | How it is set | What it means |
|---|---|---|
| the outcome is known | `data-tone` | The check happened and this is how it ended |
| the outcome is unknown | no attributes | The check happened, and we have no result |
| there was no check | `data-empty` | **Nobody looked.** The room stays, the mark does not |

The difference between the last two matters more than it seems. "We do not
know the result" is about the check, "there was no check" is about us.

A gap is neither painted nor outlined, and both were tried. A border knocks the
density of the row about: the eye reads a change of **texture** as a change of
rhythm rather than as a change of meaning, and the strip "swims". Colour
promises a judgement — and on real data a gap turns out to be the norm rather
than an exception: in an application where the history came from an import, a
tone would paint 288 ticks of 336, and nothing would be an exception any more.

If a gap deserves alarm, its place is **in the line beside it**: "no checks
since 07:14". One statement, one place.

## Composition

### With an axis

The strip says what happened, the axis says when. An axis cell takes the **same
`--n`** as the group above it, so the two rows are laid out by one rule and
cannot drift apart.

```html preview
<div style="inline-size:30rem">
  <div class="inst-history inst-history--lg" role="img"
       aria-label="Twelve hours, 60 checks: three with an error at 03:00">
    <span class="inst-history-group" style="--n:5"><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-tooltip-text">00:00 &mdash; 5 checks</span></span>
    <span class="inst-history-group" style="--n:5"><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-tooltip-text">01:00 &mdash; 5 checks</span></span>
    <span class="inst-history-group" style="--n:5"><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-tooltip-text">02:00 &mdash; 5 checks</span></span>
    <span class="inst-history-group" style="--n:5"><span class="inst-history-tick" data-tone="error"></span><span class="inst-history-tick" data-tone="error"></span><span class="inst-history-tick" data-tone="error"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-tooltip-text">03:00 &mdash; 5 checks</span></span>
    <span class="inst-history-group" style="--n:2"><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-tooltip-text">04:00 &mdash; 2 checks</span></span>
    <span class="inst-history-group" style="--n:5"><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-tooltip-text">05:00 &mdash; 5 checks</span></span>
    <span class="inst-history-group" style="--n:5"><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-tooltip-text">06:00 &mdash; 5 checks</span></span>
    <span class="inst-history-group" style="--n:5"><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-tooltip-text">07:00 &mdash; 5 checks</span></span>
    <span class="inst-history-group" style="--n:8"><span class="inst-history-tick" data-tone="ok"></span><span class="inst-tooltip-text">08:00 &mdash; 8 checks, shown as one</span></span>
    <span class="inst-history-group" style="--n:5"><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-tooltip-text">09:00 &mdash; 5 checks</span></span>
    <span class="inst-history-group" style="--n:5"><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-tooltip-text">10:00 &mdash; 5 checks</span></span>
    <span class="inst-history-group" style="--n:5"><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-tooltip-text">11:00 &mdash; 5 checks</span></span>
  </div>
  <div class="inst-history-axis" aria-hidden="true">
    <span class="inst-history-axis-cell" style="--n:5">00</span>
    <span class="inst-history-axis-cell" style="--n:5" data-minor>01</span>
    <span class="inst-history-axis-cell" style="--n:5" data-minor>02</span>
    <span class="inst-history-axis-cell" style="--n:5">03</span>
    <span class="inst-history-axis-cell" style="--n:2" data-minor>04</span>
    <span class="inst-history-axis-cell" style="--n:5" data-minor>05</span>
    <span class="inst-history-axis-cell" style="--n:5">06</span>
    <span class="inst-history-axis-cell" style="--n:5" data-minor>07</span>
    <span class="inst-history-axis-cell" style="--n:8" data-minor>08</span>
    <span class="inst-history-axis-cell" style="--n:5">09</span>
    <span class="inst-history-axis-cell" style="--n:5" data-minor>10</span>
    <span class="inst-history-axis-cell" style="--n:5" data-minor>11</span>
  </div>
</div>
```

`data-minor` marks the labels that may go when the ruler is under 36rem. The
cell keeps its place — `visibility`, not `display` — or the surviving labels
would slide.

The hint is an ordinary [tooltip](../components/overlays/tooltip.md) placed
straight inside the group. Pointing is answered by the hour itself rising, which
costs nothing horizontally: block margin takes no part in inline sizing, so the
hours cannot shift as the cursor travels.

### In a monitor card

```html preview context
<div class="inst-card" style="max-inline-size:24rem">
  <div class="inst-card-head">
    <span class="inst-card-title">absalyamov.ru</span>
    <span class="inst-badge" data-tone="ok"><span class="inst-dot"></span>up</span>
  </div>
  <div class="inst-card-sub">https://absalyamov.ru/ · http</div>
  <div class="inst-history" role="img" aria-label="A day of checks, four an hour: three with an error at 12:45, two hours with no checks">
    <span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span>
    <span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-empty></span><span class="inst-history-tick" data-empty></span>
    <span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span>
    <span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span>
    <span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="error"></span><span class="inst-history-tick" data-tone="error"></span><span class="inst-history-tick" data-tone="error"></span><span class="inst-history-tick" data-tone="warn"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span>
    <span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span>
    <span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span>
    <span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span>
  </div>
</div>
```

The strip spans the card because it was given enough checks to, not because it
stretches: ninety-six is a day at four an hour, and the oldest are clipped at the
leading edge. To fill a width the caller cannot predict, use the
[grouped form](#groups).

## API

```api
```

## Related

```related
```
