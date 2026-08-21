---
title: Scale
group: Foundations
layout: foundation
source: src/tokens.css
api:
  - { name: "data-scale", kind: "attribute", doc: "`15` `16` `17` `18` — the base type size of the interface. The absence of the attribute is a base of 14px" }
  - { name: "--text-2xs/xs/sm/md/lg/xl/2xl", kind: "token" }
  - { name: "--radius-2xs/xs/sm/md/lg", kind: "token" }
  - { name: "--control-h-xs/sm/md/lg", kind: "token" }
  - { name: "--radius-control-sm/md/lg", kind: "token" }
  - { name: "--size-dot", kind: "token" }
  - { name: "--size-icon", kind: "token" }
  - { name: "--size-icon-lg", kind: "token" }
  - { name: "--size-icon-sm", kind: "token" }
  - { name: "--size-chevron", kind: "token" }
  - { name: "--size-track", kind: "token" }
  - { name: "--size-meter", kind: "token" }
  - { name: "--size-tick", kind: "token" }
  - { name: "--size-marker", kind: "token" }
  - { name: "--label-col", kind: "token" }
  - { name: "--aside-w", kind: "token" }
  - { name: "--col-min", kind: "token" }
  - { name: "--container-max", kind: "token" }
---

Scale is an **attribute of a container**, as [density](./density.md) is. There
is one difference between them, and it is all in the type size.

**Density moves geometry and does not touch type size.** The text stays 14px in
all three modes, and what changes is the air around it: how many rows fit on a
screen.

**Scale moves type size together with geometry.** The proportions are the same,
the drawing holds, and the screen as a whole becomes larger.

A step is named by its base type size: `data-scale="16"` means "the interface
base is 16px". The name of the step is the answer to the question a person asks
when choosing a scale.

```html preview
<div class="inst-cluster" data-scale="16">
  <button class="inst-btn inst-btn--sm" type="button">sm</button>
  <button class="inst-btn" type="button">md</button>
  <span class="inst-badge" data-tone="ok"><span class="inst-dot"></span>done</span>
  <input class="inst-input inst-input--sm" type="text" value="16">
</div>
```

## Contract

```html
<html data-scale="16">               <!-- a person's setting, on the root -->
  <section data-density="compact">   <!-- a dense region inside -->
```

The usual arrangement is exactly this: the scale is chosen by a person and put
on the root, the density is chosen by the author of the screen and put on a
region. Both axes add up — the combination is described by a cell of tokens of
its own.

| Value | Base | What for |
|---|---|---|
| no attribute | 14px | The default. Working distance: a laptop or a monitor at the system scale |
| `data-scale="15"` | 15px | A step of margin: a long session, a middling monitor |
| `data-scale="16"` | 16px | A large monitor at 100% system scaling, working while standing |
| `data-scale="17"` | 17px | Showing on a projector, reading from two metres |
| `data-scale="18"` | 18px | A wall panel, a duty screen, weakened sight |

### Scale comes first

On a single match the token goes to scale: the `[data-scale]` blocks stand
after the density blocks in the tokens, and their specificity is the same.
Density takes its own back in the two-dimensional cell
`[data-scale][data-density]`.

The cell has two selectors, and the second is obligatory: the paired one
catches only both attributes on one element, while a compact panel inside a
large root is precisely the descendant case. Without it the panel would take
the geometry of the base while keeping the large type.

### Accessibility

Scale is an answer to an increased system font, but not a replacement for it.
The control heights are applied as `min-block-size`, so text grown by a browser
setting pushes a control apart rather than being cut off inside it.

### Tap targets

All fifteen combinations of scale and density pass WCAG 2.2 AA 2.5.8 —
`cmd/targets` checks it over a 5 × 3 grid.

## Scale

### Type size

| Token | 14 | 15 | 16 | 17 | 18 |
|---|---|---|---|---|---|
| `--text-2xs` | 11 | 11.5 | 12.5 | 13 | 14 |
| `--text-xs` | 12.5 | 13 | 14 | 15 | 16 |
| `--text-sm` | 14 | 15 | 16 | 17 | 18 |
| `--text-md` | 16 | 17 | 18.5 | 19.5 | 21 |
| `--text-lg` | 18 | 19.5 | 21 | 22.5 | 24 |
| `--text-xl` | 22 | 23.5 | 26 | 28 | 30 |
| `--text-2xl` | 28 | 30 | 33 | 35.5 | 38 |

The values are declared in `rem` and are exact at a root of 16px. The scale
axis **multiplies** the browser's "default font size" setting rather than
replacing it: somebody who has raised the root to 20px gets one more step on
top from the step here.

### Radii and heights

| Token | 14 | 15 | 16 | 17 | 18 |
|---|---|---|---|---|---|
| `--radius-2xs` | 2 | 2 | 2 | 2 | 2 |
| `--radius-xs` | 4 | 4 | 6 | 6 | 6 |
| `--radius-sm` | 6 | 6 | 8 | 8 | 8 |
| `--radius-md` | 8 | 10 | 10 | 10 | 10 |
| `--radius-lg` | 12 | 12 | 14 | 14 | 16 |
| `--control-h-xs` | 20 | 22 | 24 | 25 | 26 |
| `--control-h-sm` | 26 | 28 | 30 | 32 | 34 |
| `--control-h-md` | 32 | 35 | 38 | 40 | 42 |
| `--control-h-lg` | 38 | 41 | 44 | 47 | 50 |

The ladder of radii grows in steps rather than by a multiplier: a radius has to
stay even, otherwise at a pixel density of 1.5 the arc lands off the device
grid. Because of that neighbouring steps here and there share one value — there
is no even number between 8 and 10.

The paddings move along the `--space-*` scale, and that one is sparse on
purpose, so their step is coarser than the type size's: neighbouring steps of
scale often take one and the same step of padding.

### What is NOT scaled

| What | Why |
|---|---|
| `--hairline` | Not on principle but physically: the browser snaps a border to a device pixel, and ×1.07 of one unit comes back to the same unit |
| `--tap-min` | This is a norm rather than a size: 24px is the floor of WCAG 2.2 AA 2.5.8. The ladder goes only upwards, so a fixed floor squeezes nothing |
| `--radius-2xs` | The one nailed-down number: a dot grows to 8px, and a radius of 4 would give exactly half of it, that is, a circle. The circle is taken by the mark of a reversible item |
| `--space-*` | The scale is the foundation. Scale moves the roles, not the scale under them: the paddings move to another step of the same ladder |
| Colour and contrast | The thresholds do not depend on type size. That is why `cmd/contrast` does not check scale |

## Rules

**The ladder starts at 14 and cannot start lower.** That is not a choice but
the arithmetic of two rules of the kit multiplied by one another: the floor of
the type size is 11px, the minimum step of the ladder is 1.12, and under the
base lie two steps, `2xs` and `xs`. Hence

```
the minimum base = 11 × 1.12 × 1.12 = 13.8px
```

and 14 is the first whole number that fits. At a base of 13 the second step
down gets 10.36px, at a base of 12 it gets 9.57px; both are below the floor.
The ladder can be lowered only by moving the floor, the step or the number of
steps under the base — each of the three a separate decision. "Smaller" stays
the work of density.

**The values of the table are written by hand rather than computed.** A
proportional recount fails the rule of the step regularly: at a base of 16 the
honest share gives `xs` = 14.5, that is a step of 1.103 against a norm of 1.12.
Multiplication breaks the evenness of a radius, the growth of the ladder and
the band of the control radius against its height alike.

A component that breaks at step 18 holds a hard-coded type size or a hard-coded
height — both of them ought to be roles.

### What to take and what not to

Three densities could not say "larger". On a 27″ 4K at 100% system scaling the
kit was drawn small, and the one knob that looked like an answer — the
comfortable density — gave air rather than size: the rows parted and the
letters stayed as they were.

The test that tells an axis from a renaming is the same as for the attribute
vocabularies: if two knobs give one result, it is one knob with two names. Here
they part, and it shows in the numbers.

| | type size | `--control-h-md` | `--pad-card` |
|---|---|---|---|
| `18` + `compact` | 18px | 34px | 24px |
| base + `comfortable` | 14px | 36px | 24px |

The height and the padding coincide almost exactly, the type size does not.
These are two different screens rather than one seen from two sides.

| Use | Take instead |
|---|---|
| The steps, to make a screen larger as a whole | **Scale, to thin out a list** — air is added by [density](./density.md), and it does not touch type size |
| The attribute on the root as a person's setting | **The attribute on a single component** — scale is not a variant of a control. The size is chosen by `--sm`/`--lg` on [the button](../components/actions/button.md) |
| A combination with density: "large and compact" is the working mode of a big monitor | **A step of your own as a number** in an application — a new value gets a role, and the list is in [tokens](./tokens.md) |

## API

```api
```

## Related

[Density](./density.md)
[Typography](./typography.md)
[Rhythm and shape](./spacing.md)
[Tokens](./tokens.md)
