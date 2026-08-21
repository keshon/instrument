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
  - { name: "inst-history-group", kind: "class", doc: "A batch of ticks with a shared name: an hour, a day, a run" }
  - { name: "--n", kind: "token", doc: "The number of ticks in a group. It sets its share of the strip" }
  - { name: "role", kind: "attribute", value: "img", doc: "The strip is a picture of data, and it has to have a name" }
  - { name: "--size-tick", kind: "token" }
  - { name: "--space-1", kind: "token" }
  - { name: "--space-3", kind: "token" }
  - { name: "--radius-xs", kind: "token" }
  - { name: "--track", kind: "token" }
  - { name: "--size-icon", kind: "token" }
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
goes first" is not wanted here, everything fits by construction. That is also
what makes it possible to hang [a tooltip](../components/overlays/tooltip.md)
on a tick — under an `overflow: hidden` it would be cut off.

| With no groups | With groups |
|---|---|
| The width of a tick is fixed — the length of the strip says how many checks there were | The width of the strip is fixed — the length of a group says how many checks are in that hour |
| What is extra is clipped at the start | It is not clipped: the shares are counted from what there is |
| One row of a list beside another | One object at large, on a screen of its own |

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

### In a monitor card

```html preview context
<div class="inst-card" style="max-inline-size:24rem">
  <div class="inst-card-head">
    <span class="inst-card-title">absalyamov.ru</span>
    <span class="inst-badge" data-tone="ok"><span class="inst-dot"></span>up</span>
  </div>
  <div class="inst-card-sub">https://absalyamov.ru/ · http</div>
  <div class="inst-history" role="img" aria-label="The last 20 checks: one with an error">
    <span class="inst-history-tick" data-tone="ok"></span>
    <span class="inst-history-tick" data-tone="ok"></span>
    <span class="inst-history-tick" data-tone="ok"></span>
    <span class="inst-history-tick" data-tone="error"></span>
    <span class="inst-history-tick" data-tone="ok"></span>
    <span class="inst-history-tick" data-tone="ok"></span>
    <span class="inst-history-tick" data-tone="ok"></span>
    <span class="inst-history-tick" data-tone="ok"></span>
    <span class="inst-history-tick" data-tone="ok"></span>
    <span class="inst-history-tick" data-tone="ok"></span>
  </div>
</div>
```

## API

```api
```

## Related

```related
```
