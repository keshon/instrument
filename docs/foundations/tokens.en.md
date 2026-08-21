---
title: Tokens
group: Foundations
layout: foundation
source: src/tokens.css
api-from: kit
---

The full reference. Four tiers, and the separation between them is the whole
point of the file: **a component may see tiers 2 and 3, and tier 1 never**.

```html preview
<div class="inst-panel">
  <div class="inst-panel-header">
    <span class="inst-panel-title">The tiers</span>
  </div>
  <div class="inst-panel-body">
    <dl class="inst-kv">
      <dt>1. Ramps</dt><dd>Raw steps. They do not change between themes</dd>
      <dt>2. Semantics</dt><dd>What a colour is for. One declaration through light-dark()</dd>
      <dt>3. Roles</dt><dd>Paddings and sizes by purpose. Density retunes these</dd>
      <dt>4. Component</dt><dd>--btn-bg and the like, 2–3 variables</dd>
    </dl>
  </div>
</div>
```

## Contract

A component sees **tiers 2 and 3** — the semantics and the roles. It never sees
tier 1.

```css
.my-thing {
  /* tier 2: what the colour is for */
  color: var(--text-primary);
  /* tier 3: a padding by purpose */
  padding: var(--pad-panel);
  /* tier 4: the point of adjustment of the component itself.
     A button has five of them: --btn-bg and its neighbours */
  background: var(--btn-bg);
}
```

| What | Required | Why |
|---|---|---|
| Tiers 2 and 3 in a component | yes | Tier 1 is raw material: a ramp does not change between themes, and reaching for it hard-codes the light theme |
| One instance of every parameter | yes | A duplicate is a future divergence, and it arrives sooner than it seems |
| 2–3 variables per component | yes | Tier 4 is points of adjustment, not a copy of every property |
| A value through `light-dark()` | yes | A second theme as a separate block does not exist |


### Accessibility

| | |
|---|---|
| The check | `go -C tools run ./cmd/contrast` reads the real `tokens.css` and resolves `light-dark()`, `color-mix()` and `var()` the way a browser does. A new pair of tokens is added there together with the token |
| Thresholds | Text 4.5:1, a mark 3:1 (against the surface and against the track), a load-bearing border 3:1 |
| Six themes | Every pair is checked in all six rather than in two |
| Sizes in `rem` | The browser's "default font size" setting works; the heights are applied as `min-block-size` |
| The floor of the type size | `--text-2xs`, 11px. There is no token below it |

## Anatomy

Raw steps. They do not change between themes. A ramp is continuous by
definition: a step the semantics has not taken is a reserve for the application
and for charts rather than rubbish.

| Group | Tokens | Hue |
|---|---|---|
| Neutral | `--n-0` … `--n-14` | `--hue-neutral` (75 warm / 250 cool) |
| Accent | `--a-1` … `--a-6` | 250° |
| Success | `--ok-1` `--ok-2` `--ok-3` `--ok-4` `--ok-5` `--ok-6` | 150° |
| Warning | `--warn-1` `--warn-3` `--warn-4` `--warn-5` `--warn-6` | 85° |
| Error | `--err-1` `--err-3` `--err-4` `--err-5` `--err-6` | 25° |

The statuses have no step 2: their roles are fixed (1 is the background in a
light theme, 3 the text in a dark one, 4 the fill in a light one, 5 the text in
a light one, 6 a reserve for charts), and a free step between 1 and 3 has not
been wanted once. The lightness values and the analysis of the roles are in
[colour](./colors.md).

### The hue knob

| Token | Value | What it does |
|---|---|---|
| `--hue-neutral` | 75 | The cast of the neutral. 75 is warm, 250 cool. A chroma of 0.002–0.006: below the threshold of conscious notice |

### Tier 2 — semantics

One declaration per token through `light-dark()`. A second block of tokens does
not exist, so there is nothing to diverge.

### Surfaces

| Token | Light | Dark |
|---|---|---|
| `--surface-sunken` | `--n-2` | `--n-14` |
| `--surface-page` | `--n-1` | `--n-13` |
| `--surface-raised` | `--n-0` | `--n-12` |
| `--surface-overlay` | `--n-0` | `--n-11` |
| `--surface-field` | `--surface-recessed`, in both themes | |
| `--surface-hover` | black 3.5% | white 4.5% |
| `--surface-active` | black 6.5% | white 8% |
| `--surface-selected` | `--a-4` 14%, in both themes | |
| `--surface-recessed` | black 6% | white 5% |
| `--surface-recessed-hover` | black 2.5% | white 1.5% |
| `--surface-recessed-active` | black 7.5% | white 6.2% |

The last three are the recess of a control: the default button, a tag, the
fallback fill of a badge. Black laid over rather than a step of the ramp of its
own, and that is no trifle. A step sets lightness absolutely and therefore
knows one depth, while a control lies now on a page, now on a panel, now in a
modal — and it has to step away from each of them by the same amount. On a
raised panel `--surface-sunken` gave a drop of 0.087 instead of 0.033: the
panel is already a step higher, and the steps added up.

**The direction depends on where we are on the ramp**, as with
`--surface-hover`. In the light themes a control is recessed — a black film; in
the dark ones it is raised — a white one.

That is not a refusal of the drop but the only form of it available. A film
subtracts lightness, and at the bottom of the ramp there is nothing to subtract
from: the page of a dark theme is `rgb(12,12,12)`, and even a solid black gives
1.073 against it against 1.142 in the light ones. The drop that is wanted
**does not exist** as a darkening. And the default button has its fill as its
only mark: there is no border and no shadow, the drop removed both.

The invariant is therefore stated without the word "depth": the surface of a
control has to differ from its ground by a discernible step.

**The sign of the film changes, the order of loudness does not.** The middle
row serves two jobs at once: it is the hover of the default and the resting
fill of the soft weight `inst-btn--soft`. So it has to be quieter than the
default in both themes, otherwise the third step of the ladder sounds louder
than the second. Loudness here is the modulus of the step against the ground
rather than the alpha: in the light themes quieter means less black, in the
dark ones less white.

Three `cmd/contrast` pairs guard this: "a button against a panel", "a button
against a page" and "the ladder: soft against default".

### Text

| Token | Light | Dark | Threshold |
|---|---|---|---|
| `--text-primary` | `--n-12` | `--n-1` | 4.5:1 |
| `--text-secondary` | `--n-9` | `--n-5` | 4.5:1 |
| `--text-muted` | `--n-8` | `--n-6` | 4.5:1 |
| `--text-faint` | `--n-7` | `--n-7` | 3:1 — **not for text that gets read** |

### Borders

| Token | Light | Dark | Requirement |
|---|---|---|---|
| `--border-subtle` | black 7% | white 6% | Decorative |
| `--border` | black 12% | white 11% | Decorative |
| `--border-strong` | black 22% | white 20% | Decorative |
| `--border-control` | black 46% | white 38% | Load-bearing, 3:1 |

### Accent

| Token | Light | Dark | Work |
|---|---|---|---|
| `--accent-text` | `--a-5` | `--a-3` | Text, 4.5:1 |
| `--accent-mark` | `--a-4` | `--a-3` | A mark with no label, 3:1 against the surface **and** the track |
| `--accent-bg` | `--a-1` | `--a-4` 15% | A tinted background |
| `--accent-solid` | `--a-4`, in both themes | | The fill of a button, a lightness ceiling that depends on hue |
| `--accent-on` | `--n-0`, in both themes | | The label on a fill |
| `--accent-hover` | `--a-5`, in both themes | | It moves **away** from the colour of the label rather than towards it |
| `--accent-border` | `--a-4` | `--a-3` | The border of an accented control |
| `--focus-ring` | `--a-6` | `--a-2` | The focus ring. The one step of the ramp that is never a fill, so it outlines one without merging into it |

### Statuses

| Token | Light | Dark |
|---|---|---|
| `--ok-text` | `--ok-5` | `--ok-3` |
| `--ok-bg` | `--ok-1` | `--ok-4` 16% |
| `--warn-text` | `--warn-5` | `--warn-3` |
| `--warn-bg` | `--warn-1` | `--warn-4` 16% |
| `--err-text` | `--err-5` | `--err-3` |
| `--err-bg` | `--err-1` | `--err-4` 16% |

There is no third token — a "fill" between the text and the background. There
is one exception, the accent: it has a case the statuses do not, a fill under a
white label.

### Tone through an attribute

`data-tone` sets three variables, and a component reads them and knows nothing
else about the tone.

| Value | `--tone-ink` | `--tone-mark` | `--tone-bg` |
|---|---|---|---|
| `neutral` | `--text-secondary` | `--text-secondary` | `--surface-sunken` |
| `running` | `--accent-text` | `--accent-mark` | `--accent-bg` |
| `ok` | `--ok-text` | `--ok-text` | `--ok-bg` |
| `warn` | `--warn-text` | `--warn-text` | `--warn-bg` |
| `error` | `--err-text` | `--err-text` | `--err-bg` |

### Charts

| Token | Hue | Light | Dark |
|---|---|---|---|
| `--chart-1` | 292° | L 0.520 | L 0.700 |
| `--chart-2` | 320° | L 0.600 | L 0.760 |
| `--chart-3` | 355° | L 0.550 | L 0.720 |
| `--chart-4` | 55° | L 0.650 | L 0.800 |
| `--chart-5` | 115° | L 0.580 | L 0.740 |
| `--chart-6` | 190° | L 0.620 | L 0.780 |

### The rest

| Token | Light | Dark |
|---|---|---|
| `--track` | black 10% | white 16% |
| `--scrim` | black 32% | black 58% |
| `--shadow-color-near` | black 6% | black 30% |
| `--shadow-color-far` | black 10% | black 45% |
| `--shadow-popover` | two layers, short | |
| `--shadow-modal` | two layers, long | |

### The order of layers

| Token | Value |
|---|---|
| `--z-sticky` | 10 |
| `--z-popover` | 100 |
| `--z-modal` | 1000 |

### Tier 3 — roles

The one tier density retunes, and the one reason there are no constants in the
components. Declared through `:where()` — zero specificity, so an application
overrides without a fight. The values below are for the default density, and
the other two are in [density](./density.md).

### Control heights and paddings

| Token | Value |
|---|---|
| `--control-h-xs` | 20px |
| `--control-h-sm` | 26px |
| `--control-h-md` | 32px |
| `--control-h-lg` | 38px |
| `--control-pad-sm` | `--space-4` |
| `--control-pad-md` | `--space-5` |
| `--control-pad-lg` | `--space-6` |

### Paddings by purpose

| Token | Value |
|---|---|
| `--pad-panel` | `--space-5` |
| `--pad-card` | `--space-6` |
| `--pad-cell-x` | `--space-4` |
| `--row-pad-y` | `--space-3` |
| `--gap-row` | `--space-3` |
| `--gap-inline` | `--space-4` |
| `--gap-section` | `--space-7` |

`--gap-section` is the fourth step of the rhythm: the gap **between sections**,
where `--gap-row` works inside one. The value is not appointed but derived: it
is `--gap-row` raised by four steps of the scale. The scale is built so that +4
steps give exactly ×4 at every point — 4→16, 6→24, 8→32, 12→48, 16→64 — so the
ratio holds across all fifteen combinations of scale and density with no
separate fitting.

The role appeared after assembled screens were measured. Without it the tier
stopped at 16px: the loudest thing available to part two blocks was of the same
order as the gap inside a block — and the screen read as one continuous flow.
The steps 24, 32, 48 and 64 meanwhile lay in the scale without a single role.

### Glyphs and gutters

| Token | Value |
|---|---|
| `--size-dot` | 6px |
| `--size-check` | 15px |
| `--size-switch-w` | 30px |
| `--size-switch-h` | 18px |
| `--size-thumb` | 14px |
| `--size-track` | 4px |
| `--size-meter` | 6px |
| `--size-spinner` | 14px |
| `--size-icon` | 18px |
| `--size-chevron` | 10px |
| `--size-marker` | 2px |
| `--size-gutter` | 14px |
| `--size-indent` | 14px |
| `--label-col` | 104px |

### Layout

| Token | Value |
|---|---|
| `--aside-w` | 232px |
| `--container-max` | 76rem |
| `--split-side` | 18rem |
| `--split-main` | 30rem |
| `--col-min` | 260px |

## Scale

Not tied to a theme and not retuned by density. The details are in
[typography](./typography.md) and [rhythm and shape](./spacing.md).

### Typography

| Group | Tokens |
|---|---|
| Families | `--font-sans` `--font-mono` |
| Sizes | `--text-2xs` `--text-xs` `--text-sm` `--text-md` `--text-lg` `--text-xl` `--text-2xl` |
| Weights | `--weight-normal` `--weight-medium` |
| Leading | `--leading-tight` `--leading-ui` `--leading-code` `--leading-prose` |
| Tracking | `--tracking-tight` `--tracking-normal` |
| Measure | `--measure` |

`--measure` is the limit on the length of a line, `68ch`. The classical norm is
45–75 characters; 68 stands at the upper edge, because interface descriptions
are short and wrap rarely. The unit is `ch` rather than pixels: otherwise the
measure parts ways with the type size on each of the five steps of scale.

The role appeared after a measurement: `68ch` was written as a literal in two
places, while the body of a banner had no measure at all and reached 105
characters on a container of 1440px. What is truncated is the **body** of the
banner rather than the banner: the ground of a page-level message runs the full
width, the text on it does not.

### Space and shape

| Group | Tokens |
|---|---|
| The scale | `--space-1` … `--space-10` (2 · 4 · 6 · 8 · 12 · 16 · 24 · 32 · 48 · 64) |
| Radii | `--radius-2xs` 2 · `--radius-xs` 4 · `--radius-sm` 6 · `--radius-md` 8 · `--radius-lg` 12 · `--radius-full`. Controls take `--radius-control-sm/md/lg`, which travel with density |
| Lines | `--hairline` (0.5px from 1.5dppx, otherwise 1px) · `--stroke` 1px |

### Motion

| Group | Tokens |
|---|---|
| Durations | `--dur-1` 80ms · `--dur-2` 140ms · `--dur-3` 220ms |
| Easing | `--ease-out` · `--ease-in-out` |

### Tier 4 — the component

It lives beside the component rather than here. There is one rule: **2–3
variables, no more**, and a variant is described by them rather than by
rewriting the block.

```css
.my-btn--accent {
  --btn-fg: var(--accent-on);
  --btn-bg: var(--accent-solid);
}
```

The list of a particular component's variables is on its page: the
[button](../components/actions/button.md) has five, the
[table](../components/display/table.md) has its own.

## Variants

| Attribute | `color-scheme` | `--hue-neutral` | And also |
|---|---|---|---|
| none | by the system | 75 | |
| `data-theme="light-neutral"` | `light` | 75 | `--tint: 0` — the default, a pure grey |
| `data-theme="light"` | `light` | 75 | `--tint: 2` — warm |
| `data-theme="light-cool"` | `light` | 250 | `--tint: 2` — cool |
| `data-theme="dark-soft"` | `dark` | 75 | The stack one step up, plus the status labels |
| `data-theme="dark"` | `dark` | 75 | |

## Rules

| Rule | Consequence |
|---|---|
| A component sees tiers 2 and 3 | A `--n-3` in a component is a hard-coded light theme; a `14px` is a broken density |
| A tone that carries meaning has exactly two tokens | `--*-text` (the mark) and `--*-bg` (the background). A third did not pass the threshold on a track |
| Text and a mark are different foregrounds | `--tone-ink` holds 4.5:1, `--tone-mark` 3:1, but against two backgrounds |
| A load-bearing and a decorative border are different tokens | `--border-control` owes 3:1, `--border` does not |
| Defaults are declared through `:where()` | A default with the same specificity as an override is not a default |
| A new value is a role rather than a constant | Otherwise it will not survive `compact` |
| There is no such thing as a second accent tone | That is not a feature request but a signal that the meaning is carried by something else |

### What to take and what not to

| Use | Take instead |
|---|---|
| Tier 2 — the semantics: `--text-primary`, `--surface-raised`, `--ok-text` | **Tier 1** (`--n-3`, `--a-4`) — that is a hard-coded light theme. The ramps are analysed in [colour](./colors.md) |
| Tier 3 — the roles: `--pad-panel`, `--control-h-md`, `--size-dot` | **A number in a component** — it will break density, see [density](./density.md) |
| Tier 4 — 2–3 variables of a component for a variant | **A copy of the whole block of rules** — a variant should be two lines, as on [the button](../components/actions/button.md) |
| A new parameter — as a role | **A new constant** — a constant will outlive you and break `compact`. The scales are in [rhythm and shape](./spacing.md) |
| Overriding a role in an application | **`!important`** — the roles are declared through `:where()`, there will be no fight. An example is [the form](../components/inputs/form.md) |

## API

Listed by a generator reading `src/tokens.css`. The list is not typed by hand
and therefore cannot fall behind the sources.

```api
```
## Related

[Colour](./colors.md)
[Typography](./typography.md)
[Rhythm and shape](./spacing.md)
[Elevation and surfaces](./elevation.md)
[Motion](./motion.md)
[Density](./density.md)
[Icons](./icons.md)
