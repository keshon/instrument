---
title: Spinner
group: Feedback
layout: component
source: src/status.css
api:
  - { name: "inst-spinner", kind: "class", doc: "The base one. Put on an `<svg>`" }
  - { name: "inst-spinner-track", kind: "class", doc: "The track. A direct child of `inst-spinner`" }
  - { name: "inst-spinner-arc", kind: "class", doc: "The arc. A direct child of `inst-spinner`" }
  - { name: "--size-spinner", kind: "token" }
  - { name: "--border-control", kind: "token" }
  - { name: "--accent-mark", kind: "token" }
---

A ring showing that the machine is busy when there is nothing to name the share
done with. A real SVG of two circles rather than a border on a `div`.

```html preview
<svg class="inst-spinner" viewBox="0 0 16 16" role="status" aria-label="Loading">
  <circle class="inst-spinner-track" cx="8" cy="8" r="6.5"/>
  <circle class="inst-spinner-arc" cx="8" cy="8" r="6.5"/>
</svg>
```

## Contract

The markup is fixed: `viewBox="0 0 16 16"`, two circles of radius `6.5` centred
at `8 8`. It is copied entire.

| What | Why it cannot be changed |
|---|---|
| `viewBox="0 0 16 16"` | The thickness of the stroke and the `stroke-dasharray` are computed for this coordinate system |
| `r="6.5"` `cx="8"` `cy="8"` | The circumference at `r=6.5` is about 40.8, and `11 30` gives exactly a quarter to the arc |
| `role="status"` | Otherwise the appearance of the spinner says nothing |
| `aria-label` | A `role="status"` has to have a name: "Loading", "The run is going" |

Both circles are direct children. Nesting them deeper is not allowed: the
stroke will not apply.

### Accessibility

| | |
|---|---|
| The role | `role="status"` + `aria-label`. A spinner with no role is just a picture, and for a screen reader the waiting does not begin |
| Not the only carrier | A long wait comes with words: "The run is going" rather than one turning ring |
| Reduced motion | The turn **slows to 2.4s rather than stopping**. Collapsing an animation to `0.01ms` is not a squeezed duration but a stop: the spinner would freeze at a random phase, and the indicator would stop showing that the machine is busy |
| Contrast | The arc `--accent-mark` holds 3:1 against the track rather than against the background alone |
| Print | It is not printed on paper: print is static |

## Anatomy

| Element | The work |
|---|---|
| `inst-spinner` | The `<svg>`. The size, the thickness of the stroke, the rotation |
| `inst-spinner-track` | The first circle — the track, `--border-control` |
| `inst-spinner-arc` | The second circle — the arc, `--accent-mark`. A quarter of the length through `stroke-dasharray: 11 30` |

Both circles are identical; what sets them apart is only the stroke. The track
is wanted so that the ring is seen entire: one arc without it reads as a
fragment.

| Property | Value |
|---|---|
| size | `--size-spinner` — 14px; 13px in the compact mode, 16px in the comfortable one |
| stroke thickness | `2` |
| ends | `stroke-linecap: round` |
| turn | 0.7s, linear, infinite |

`--accent-mark` rather than `--accent-text`: a tone has two foregrounds, and a
mark with no label takes the one checked at 3:1 against the track.

### Why SVG

A circle made of a `border-radius` on a box of 13px at a fractional DPR
rasterises with its centre on half a pixel, and on a rotation that shows as
**wobble**. In an SVG the centre is set in the coordinate system of the
`viewBox`, the rotation is applied to the box entire and the rasterisation
happens after — there is no wobble at any DPR.

## API

```api
```

## Related

```related
```
