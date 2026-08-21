---
title: Motion
group: Foundations
layout: foundation
source: src/tokens.css
api:
  - { name: "--dur-1", kind: "token" }
  - { name: "--dur-2", kind: "token" }
  - { name: "--dur-3", kind: "token" }
  - { name: "--ease-out", kind: "token" }
  - { name: "--ease-in-out", kind: "token" }
---

Three durations and two easing functions. Fast enough for motion to read as
**a cause rather than an animation**: everything the user set off has to
resolve within `--dur-2`.

```html preview
<div class="inst-cluster">
  <button class="inst-btn" type="button">Hover and press</button>
  <button class="inst-btn inst-btn--primary" type="button" aria-busy="true">Run</button>
  <span class="inst-badge" data-tone="running"><span class="inst-dot"></span>running</span>
</div>
```

## Contract

Motion in the library answers one question: **what changed**. The duration
comes from the scale, the easing from the two.

```css
.my-thing {
  transition: background var(--dur-1) var(--ease-out);
}
```

| What | Required | Why |
|---|---|---|
| A duration from `--dur-*` | yes | Four steps. A number in place parts from the rest of the library on the first edit |
| `--ease-out` for an appearance | yes | The entry is fast at the start and soft at the end; `--ease-in-out` is only for what travels back and forth |
| Reduced motion accounted for | yes | Infinite indicators **slow down rather than go out**: a stopped indicator says "nothing is happening" while something is |
| An animation on data | no | A chart is a state rather than a process. What moves is what changes on an action |


### Accessibility

| | |
|---|---|
| Reduced motion | Transitions collapse to 0.01ms; indicators of activity **slow down rather than stop** |
| Busy | Spoken by `aria-busy` rather than by an animation. Motion is a second carrier, not the only one |
| Progress | An indeterminate meter has no `aria-valuenow`, and that is a message in itself |
| Focus | Not animated: the ring appears at once |
| Flashing | Nothing in the library flashes faster than 3 Hz |
| Layer | `kit.motion` is declared last, so it beats the components and does not beat the application — without a single `!important` |

## Scale

| Token | Value | Where |
|---|---|---|
| `--dur-1` | 80ms | Hover, press, focus, the turn of a chevron, the underline of a link. This is most of the library's transitions |
| `--dur-2` | 140ms | A switch, the appearance and disappearance of a popover. The ceiling for everything the user set off |
| `--dur-3` | 220ms | The filling of a meter — the only thing that changes geometry from data rather than from a touch |

Everything that does not fit within `--dur-2` **calls for a state rather than a
transition**: if a process is longer, what should show it is an indicator of
busyness, not a stretched animation.

### Easing functions

| Token | Value | Where |
|---|---|---|
| `--ease-out` | `cubic-bezier(0.22, 0.61, 0.36, 1)` | Every transition. A sharp start and a soft stop read as a response to an action |
| `--ease-in-out` | `cubic-bezier(0.65, 0, 0.35, 1)` | Infinite cycles: the pulse of a dot, the breathing of a skeleton, an indeterminate meter |

Linear stays where a rotation has to be even. The ring of a spinner and the
ring of a busy button turn `linear`: a slowdown within a cycle reads as a jerk.

## Behavior

| What | Duration | How |
|---|---|---|
| Hover and press of a button | `--dur-1` | `background`, `border-color`, `transform` |
| A press | `--dur-1` | `scale(0.985)` — a press is to be felt rather than examined |
| Focus | — | The ring appears instantly: a delay on focus reads as lag |
| The chevron of a select and of an opened step | `--dur-1` | `rotate` |
| A switch | `--dur-2` | The position of the thumb and the fill of the track |
| A popover | `--dur-2` | `opacity` + `translate`, plus `display` and `overlay` through `allow-discrete` |
| The filling of a meter | `--dur-3` | `inline-size` |
| A `running` dot | 1.6s, infinite | A pulse of opacity |
| A spinner and the ring of a busy button | 0.7s, infinite | Rotation |
| A skeleton | 1.6s, infinite | Breathing by opacity, **not a sheen**: gradients are forbidden |
| An indeterminate meter | 1.3s, infinite | A segment travelling along the track |
| The caret of streaming | 1s, infinite | Blinking in steps |

### Reduced motion

`prefers-reduced-motion: reduce` is handled in a cascade layer of its own,
declared last. The layer is needed so that the rule beats the components
(`.inst-btn` has a higher specificity than `*`) while **not beating the
application**: without the layer that would call for `!important`, and an
`!important` inside a layer beats the application's ordinary declarations —
that is, the library would break its own promise.

### Transitions collapse

```css
*, *::before, *::after {
  transition-duration: 0.01ms;
  animation-delay: 0ms;
  scroll-behavior: auto;
}
```

`0.01ms`, not `none`. The `transitionend` event keeps arriving, and the
application's state machines that listen for it keep working. A switched-off
transition would stop them in silence.

### Indicators do not go out

Here the library parts from the customary "switch everything off".

An infinite animation **may not** be collapsed. An animation of one iteration
squeezed to `0.01ms` is not a squeezed duration, it is a stop: the spinner
freezes at a random phase, the dot stops pulsing, the caret sets. And the
indicator then says the work has halted. On a busy button a stop turns it into
an empty rectangle: the label is dimmed and the ring stands still.

So the indicators move to a frequency that reads as "running" without tugging
at peripheral vision. What makes people sick is **amplitude and frequency**,
not motion in itself.

| Indicator | Ordinarily | At `reduce` |
|---|---|---|
| `.inst-spinner` | 0.7s | 2.4s |
| `.inst-btn[aria-busy="true"]` | 0.7s | 2.4s |
| `.inst-skeleton` | 1.6s | 3s |
| `.inst-meter` in the indeterminate state | 1.3s | 3s |
| `.inst-dot` in the `running` tone | 1.6s | 3s |
| `.inst-dot` inside `.inst-task` and `.inst-step` in the `running` state | 1.6s | 3s |
| `.inst-caret` | 1s | The blinking is removed, the caret stays visible |

On a caret the blinking says only that it is there, nothing about speed. So the
blinking is suppressed and the caret itself stays.

:::note
The list of dots that get slowed has to match what starts the pulse in the
components. Let it part, and somewhere an indicator is left freezing at a
random phase.
:::

## Rules

### What to take and what not to

| Use | Take instead |
|---|---|
| A transition of surface and border on a user's action | **Motion of the layout** — a hover does not move the layout. State is shown by [a badge](../components/display/badge.md) |
| An infinite indicator for a machine being busy | **A determinate bar when the percentage is unknown** — a bar stuck at 90% lies. An indeterminate [meter](../components/charts/meter.md) has no `aria-valuenow` |
| `--dur-3` for the filling of a meter | **The same duration for a hover** — 220ms on a hover reads as lag in the interface, see [the button](../components/actions/button.md) |
| A transition that may not play out | **An animation as the sole carrier of meaning** — the state has to be in the markup too: `aria-busy`, `data-state`. See [the step](../agent/step.md) |
| The appearance of an overlay through `@starting-style` | **Opening by hand in JS** — the top layer and Escape are given by the platform, see [the popover](../components/overlays/popover.md) |

## API

```api
```

### What is not animated

| Not animated | Why |
|---|---|
| The layout on hover | A hover does not move the layout, does not change the type size and does not add a shadow |
| Content appearing on scroll | A tooling interface is read rather than watched |
| `scroll-behavior: smooth` globally | It fights keyboard navigation and makes long logs viscous. It is switched on pointwise, for a particular interaction |
| The colour of text on a change of theme | Switching the theme is not a transition but another state of the document |

## Related

[Elevation and surfaces](./elevation.md)
[Tokens](./tokens.md)
[Button](../components/actions/button.md)
[Meter](../components/charts/meter.md)
[Step](../agent/step.md)
[Queue row](../agent/task.md)
