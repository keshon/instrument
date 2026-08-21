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

Nothing runs past the edge, and that is a property of the construction rather
than of the numbers. The gap between ticks is capped at a share of the group
instead of a constant, so the separators are what give way when an hour gets
crowded. With a constant they could not: a gap does not shrink and a flex item
does, so at seven checks an hour the six gaps took their twelve pixels first
and the strip came out wider than the box it was given.

| With no groups | With groups |
|---|---|
| The width of a tick is fixed — the length of the strip says how many checks there were | The width of the strip is fixed — the length of a group says how many checks are in that hour |
| What is extra is clipped at the start | It is not clipped: the shares are counted from what there is |
| One row of a list beside another | One object at large, on a screen of its own |

### The aggregate form

Above roughly four checks an hour a tick stops being something a person can
see, let alone point at: on a 520px day an hour's share is 17.8px, and seven
ticks with their gaps divide it into slivers under two pixels. Drawing them
individually at that point is a promise the picture cannot keep.

The answer is not a thinner tick. Give the group **one** tick and leave `--n`
alone: the hour keeps its exact share of the strip and comes back as a single
brick, tinted by the outcome that matters — the worst one.

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

The strip still reports how much stands behind each hour — the sixth hour had
one check and is narrow, the ninth had twelve and is wide — and it has stopped
claiming that the individual checks are distinguishable. Which of the two
drawings to use is the caller's call, because only the caller knows how many
checks there are.

Nothing breaks if that call is not made: the gap between ticks is capped at a
share of the group rather than a constant, so at any density the separators
thin out and the marks keep their width. What is lost is only the ability to
count them.

## Sizes

```html preview
<div class="inst-stack" style="max-inline-size:22rem">
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
<div class="inst-stack" style="max-inline-size:22rem">
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

The strip says what happened; the axis says when. It is not a second layout but
the same one: an axis cell takes the **same `--n`** as the group above it, and
the two rows are laid out by one rule. An application that builds the ruler
itself has to reproduce the gap exactly, and one pixel of disagreement compounds
— over twenty-four hours the last label ends up two dozen pixels from the column
it names.

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

Which labels survive a narrow strip is the caller's decision: `data-minor` marks
the ones that may go, and under 36rem of ruler they do. The cell **keeps its
place** — `visibility`, not `display` — because a dropped cell that also gave up
its share would let the surviving labels slide, and an hour that moves when the
window is resized is worse than an hour with no label.

The hint on a group is an ordinary [tooltip](../components/overlays/tooltip.md):
an `inst-tooltip-text` placed straight inside the group, with no wrapper.

Pointing is answered by **the hour itself rising** — the batch grows `--space-1`
above and below and stands proud of the line. A frame around the marks and a
plate behind them were both built first and neither could be seen: everything
drawn around a group shows only in the 2px gaps between its ticks and the 2px
margin around them, and at one check an hour those ticks are 17.8px wide. The
failure was area, not strength, so the carrier had to be the one surface that is
there at any density.

It costs nothing horizontally, which is what makes it allowed: a group is a flex
item with a basis of 0, so a pixel of padding or inline margin comes out of the
share it grows into and every hour shifts as the cursor travels. Block margin is
the other axis and takes no part in inline sizing.

**A `tabindex` on a group is the caller's decision, and a day of hours is the
case where the answer is no.** The hint also shows on keyboard focus, so making
a group focusable is one attribute — but a focusable element is a target, and
WCAG 2.2 AA asks a target for 24 pixels or for spacing an hourly strip cannot
give by construction: twenty-four groups in a column are about 15px wide with
4px between them. The pixel gate measures exactly that and goes red.

The number is therefore not reachable by keyboard through the group, and it does
not need to be: the strip's own `aria-label` carries the reading in words, which
is what a screen reader gets in any case — the ticks are empty and are never
announced. Twenty-four tab stops leading to what one label already says is a
worse answer than no tab stops at all.

Where a group really is a target — six runs across a panel rather than a day of
hours, so each is wider than `--tap-min` — a `tabindex` is right, and the hint
comes with it.

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
stretches: a tick is `--size-tick` wide and never a fraction, so the LENGTH of a
plain strip is how many checks there were. Ninety-six is a day at four an hour,
which over-fills the card at every scale — so the oldest checks run past the
leading edge and are clipped there, and the freshest one is always against the
trailing edge. That is the rule from [Anatomy](#anatomy) with something to look
at.

Fewer checks would sit in the trailing corner and leave the card empty, which is
the honest drawing of a monitor that has only run ten times. If the strip has to
fill a width the caller cannot predict, that is the [grouped
form](#groups) — there the shares are counted from the width there is.

## API

```api
```

## Related

```related
```
