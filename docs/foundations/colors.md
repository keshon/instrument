---
title: Colour
group: Foundations
layout: foundation
source: src/tokens.css
api:
  - { name: "inst-theme", kind: "class", doc: "The root of a theme region. Redeclares `color` and `background` so that `data-theme` works on a subtree" }
  - { name: "data-theme", kind: "attribute", value: "light-neutral · light · light-cool · dark-soft · dark", doc: "The theme. On the root of the document or on any subtree together with `inst-theme`" }
  - { name: "--text-primary", kind: "token", doc: "The main text" }
  - { name: "--text-secondary", kind: "token", doc: "Secondary text. It holds the threshold of reading" }
  - { name: "--text-muted", kind: "token", doc: "Labels and the auxiliary" }
  - { name: "--text-faint", kind: "token", doc: "Decoration. Text is not painted with it" }
  - { name: "--surface-page", kind: "token", doc: "The background of the page" }
  - { name: "--surface-sunken", kind: "token", doc: "A recess: a sunken surface" }
  - { name: "--surface-raised", kind: "token", doc: "A raised surface: a panel, a card" }
  - { name: "--surface-overlay", kind: "token", doc: "Something popping up: a popover, a modal" }
  - { name: "--surface-field", kind: "token", doc: "The recess of an input field" }
  - { name: "--surface-hover", kind: "token", doc: "Under the cursor" }
  - { name: "--surface-active", kind: "token", doc: "On a press" }
  - { name: "--surface-selected", kind: "token", doc: "What is selected" }
  - { name: "--border", kind: "token", doc: "The main boundary" }
  - { name: "--border-subtle", kind: "token", doc: "A divider inside a block" }
  - { name: "--border-strong", kind: "token", doc: "A strengthened boundary" }
  - { name: "--border-control", kind: "token", doc: "The boundary of a control. It holds 3:1" }
  - { name: "--accent-text", kind: "token", doc: "The accent for text, 4.5:1" }
  - { name: "--accent-mark", kind: "token", doc: "The accent for a mark with no label, 3:1" }
  - { name: "--accent-solid", kind: "token", doc: "A solid fill in the accent" }
  - { name: "--accent-on", kind: "token", doc: "The foreground on a solid accent" }
  - { name: "--accent-bg", kind: "token", doc: "The tinted background of the accent" }
  - { name: "--accent-border", kind: "token", doc: "A boundary in the accent" }
  - { name: "--tone-ink", kind: "token", doc: "The foreground of a tone for TEXT, 4.5:1" }
  - { name: "--tone-mark", kind: "token", doc: "The foreground of a tone for a MARK, 3:1" }
  - { name: "--tone-bg", kind: "token", doc: "The tinted background" }
  - { name: "--ok-text", kind: "token", doc: "The tone of success" }
  - { name: "--warn-text", kind: "token", doc: "The tone of a remark" }
  - { name: "--err-text", kind: "token", doc: "The tone of a refusal" }
  - { name: "--focus-ring", kind: "token", doc: "The focus ring" }
  - { name: "--scrim", kind: "token", doc: "The ground under a modal" }
  - { name: "--track", kind: "token", doc: "The track of a meter, a slider, a ring" }
---

Two tiers. Below are the **ramps** — raw steps, and they never change. Above is
the **semantics** — what a colour is for. A component sees the semantics only:
whoever wrote `--n-3` has just hard-coded the light theme.

```html preview
<div class="inst-panel">
  <div class="inst-panel-header">
    <span class="inst-panel-title">Run #4127</span>
    <span class="inst-badge" data-tone="running"><span class="inst-dot"></span>running</span>
  </div>
  <div class="inst-panel-body inst-stack">
    <input class="inst-input" type="text" value="terrain_chunk_04">
    <div class="inst-note" data-tone="warn">Three tests failed after the last pass.</div>
  </div>
</div>
```

## Contract

A component reaches for the **semantics**, and for nothing else. A ramp is the
raw material the semantics is made of, and it has no business in markup.

```css
.my-thing {
  color: var(--text-primary);
  background: var(--surface-raised);
  border: var(--hairline) solid var(--border);
}
```

| What | Required | Why |
|---|---|---|
| The semantic tier only | yes | Whoever wrote `--n-3` has just hard-coded the light theme: a ramp does not change between themes |
| A tone through `data-tone` rather than a colour of your own | yes | Five values, one vocabulary for the whole library. A sixth will not be painted |
| Two foregrounds by the work | yes | `--tone-ink` for text (4.5:1), `--tone-mark` for a mark with no label (3:1) |
| Colour is not the only carrier | yes | A state travels with a word and a sign. The first law |


### Accessibility

| | |
|---|---|
| The check | `go -C tools run ./cmd/contrast` reads the real `tokens.css` and resolves `light-dark()`, `color-mix()` and `var()` the way a browser does |
| The threshold of text | 4.5:1 in all five themes. `--text-faint` is excluded from it and therefore carries no text |
| The threshold of a mark | 3:1, but against **two** backgrounds at once: the surface and the track. Hence `--tone-mark` apart from `--tone-ink` |
| Colour is not the only carrier | A status carries a dot **and** a word, a note carries an icon, a diff row carries a sign |
| The dark theme | Not an inversion: the contrast between neighbouring surfaces widens, the chroma comes down, the borders move from a dark alpha to a light one |
| Forced-colours mode | The carriers of meaning survive the reset through `forced-color-adjust: none` |

## Anatomy

### The neutral — 15 steps

```html preview
<div class="ramp">
  <span class="ramp-step" style="--c: var(--n-0)"></span>
  <span class="ramp-step" style="--c: var(--n-1)"></span>
  <span class="ramp-step" style="--c: var(--n-2)"></span>
  <span class="ramp-step" style="--c: var(--n-3)"></span>
  <span class="ramp-step" style="--c: var(--n-4)"></span>
  <span class="ramp-step" style="--c: var(--n-5)"></span>
  <span class="ramp-step" style="--c: var(--n-6)"></span>
  <span class="ramp-step" style="--c: var(--n-7)"></span>
  <span class="ramp-step" style="--c: var(--n-8)"></span>
  <span class="ramp-step" style="--c: var(--n-9)"></span>
  <span class="ramp-step" style="--c: var(--n-10)"></span>
  <span class="ramp-step" style="--c: var(--n-11)"></span>
  <span class="ramp-step" style="--c: var(--n-12)"></span>
  <span class="ramp-step" style="--c: var(--n-13)"></span>
  <span class="ramp-step" style="--c: var(--n-14)"></span>
</div>
<div class="ramp ramp-scale">
  <span>0</span><span>1</span><span>2</span><span>3</span><span>4</span>
  <span>5</span><span>6</span><span>7</span><span>8</span><span>9</span>
  <span>10</span><span>11</span><span>12</span><span>13</span><span>14</span>
</div>
```

Switch the theme of the table: the ramp **will not change**. It is raw and
knows nothing of themes — what changes is only which of its steps the semantics
takes.

One knob sets the direction of the cast: `--hue-neutral: 75` is warm, `250`
cool. A chroma of 0.002–0.006: below the threshold of conscious notice, and
that is its job. The fifteenth step exists because the dark end needs four
discernible surfaces in a row, and on fourteen only three fitted.

| Step | Lightness | Where it is taken |
|---|---|---|
| `--n-0` | 1.000 | `--surface-raised`, `--surface-overlay` (light), `--accent-on` |
| `--n-1` | 0.976 | `--surface-page` (light), `--text-primary` (dark) |
| `--n-2` | 0.952 | `--surface-sunken` (light) |
| `--n-3` | 0.922 | Reserve |
| `--n-4` | 0.878 | Reserve |
| `--n-5` | 0.806 | `--text-secondary` (dark) |
| `--n-6` | 0.706 | `--text-muted` (dark) |
| `--n-7` | 0.605 | `--text-faint` (both themes) |
| `--n-8` | 0.500 | `--text-muted` (light) |
| `--n-9` | 0.400 | `--text-secondary` (light) |
| `--n-10` | 0.322 | `--surface-overlay` (`dark-soft`) |
| `--n-11` | 0.242 | `--surface-overlay` (dark), `--surface-raised` (`dark-soft`) |
| `--n-12` | 0.196 | `--text-primary` (light), `--surface-raised` (dark), `--surface-page` (`dark-soft`) |
| `--n-13` | 0.155 | `--surface-page` (dark), `--surface-sunken` (`dark-soft`) |
| `--n-14` | 0.120 | `--surface-sunken` (dark) |

The steps with no semantics are a reserve for the application and for charts.

### The accent — one hue

```html preview
<div class="ramp">
  <span class="ramp-step" style="--c: var(--a-1)"></span>
  <span class="ramp-step" style="--c: var(--a-2)"></span>
  <span class="ramp-step" style="--c: var(--a-3)"></span>
  <span class="ramp-step" style="--c: var(--a-4)"></span>
  <span class="ramp-step" style="--c: var(--a-5)"></span>
  <span class="ramp-step" style="--c: var(--a-6)"></span>
</div>
<div class="ramp ramp-scale">
  <span>1</span><span>2</span><span>3</span><span>4</span><span>5</span><span>6</span>
</div>
```

A second accent hue is not started. That is not a feature request but a signal
that the meaning ought to be carried by something else.

| Step | Lightness | Where it is taken |
|---|---|---|
| `--a-1` | 0.960 | `--accent-bg` (light) |
| `--a-2` | 0.800 | `--focus-ring` (dark) |
| `--a-3` | 0.700 | `--accent-text`, `--accent-mark`, `--accent-border` (dark) |
| `--a-4` | 0.545 | `--accent-solid`, `--accent-mark`, `--accent-border` (light) |
| `--a-5` | 0.450 | `--accent-text` (light), `--accent-hover` |
| `--a-6` | 0.330 | `--focus-ring` (light) |

The hue is 215°. The lightness of `--accent-solid` is bounded **from above** by
the white label, and the bound DEPENDS ON HUE: 0.545 at 215°, 0.568 at 250°,
0.576 at 275°, 0.580 at 45°. There can be no single number here — "lighter than
0.56 and the label loses 4.5:1" is the measurement of one blue. Each set holds
its own bound, and the gate checks it.

The focus ring stands on steps 6 and 2 rather than 4 and 3, and that is not a
third blue by oversight: the ring outlines a fill, and in the light themes step
4 is EXACTLY `--accent-solid`. The ring around a solid button was the colour of
the button, contrast 1.00.

### Statuses — three reserved hues

```html preview
<div class="ramp">
  <span class="ramp-step" style="--c: var(--ok-1)"></span>
  <span class="ramp-step" style="--c: var(--ok-2)"></span>
  <span class="ramp-step" style="--c: var(--ok-3)"></span>
  <span class="ramp-step" style="--c: var(--ok-4)"></span>
  <span class="ramp-step" style="--c: var(--ok-5)"></span>
  <span class="ramp-step" style="--c: var(--ok-6)"></span>
</div>
<div class="ramp">
  <span class="ramp-step" style="--c: var(--warn-1)"></span>
  <span class="ramp-step" style="--c: var(--warn-3)"></span>
  <span class="ramp-step" style="--c: var(--warn-4)"></span>
  <span class="ramp-step" style="--c: var(--warn-5)"></span>
  <span class="ramp-step" style="--c: var(--warn-6)"></span>
</div>
<div class="ramp">
  <span class="ramp-step" style="--c: var(--err-1)"></span>
  <span class="ramp-step" style="--c: var(--err-3)"></span>
  <span class="ramp-step" style="--c: var(--err-4)"></span>
  <span class="ramp-step" style="--c: var(--err-5)"></span>
  <span class="ramp-step" style="--c: var(--err-6)"></span>
</div>
<div class="ramp ramp-scale">
  <span>1</span><span>3</span><span>4</span><span>5</span><span>6</span>
</div>
```

There is no step 2 in the status rows — no role has wanted it.

`ok` is 150°, `warn` 85°, `err` 25°. They are never used as decoration or as "a
fourth series on a chart".

| Step | Role | Why exactly so |
|---|---|---|
| 1 | The tinted background in a **light** theme | `--ok-1` `--warn-1` `--err-1` |
| 2 | Text on **raised** dark surfaces | It exists for the sake of `dark-soft`: a raised panel shows through the translucent tone fill there, and step 3 drops the label to 4.15 |
| 3 | Text in the **dark** themes at the bottom of the ramp | Lighter than its neighbours in the row |
| 4 | The fill and the indicator in a **light** theme | `--ok-4` `--warn-4` `--err-4` |
| 5 | Text in a **light** theme, over step 1 | Step 4 does not reach 4.5:1 on its own tinted background |
| 6 | A reserve for charts | |

Step 5 exists precisely because step 4 does not pass the threshold on itself.
While it was not there, a yellow badge lived on a contrast of 3.02.

## Scale

One declaration per token through `light-dark()`. A second theme as a separate
block does not exist — so there is nothing to diverge.

### Surfaces

```html preview
<div class="swatches">
  <div class="sw"><div class="sw-chip" style="--c: var(--surface-sunken)"></div><div class="sw-cap">--surface-sunken</div></div>
  <div class="sw"><div class="sw-chip" style="--c: var(--surface-page)"></div><div class="sw-cap">--surface-page</div></div>
  <div class="sw"><div class="sw-chip" style="--c: var(--surface-raised)"></div><div class="sw-cap">--surface-raised</div></div>
  <div class="sw"><div class="sw-chip" style="--c: var(--surface-overlay)"></div><div class="sw-cap">--surface-overlay</div></div>
  <div class="sw"><div class="sw-chip" style="--c: var(--surface-field)"></div><div class="sw-cap">--surface-field</div></div>
</div>
```

Switch the theme of the table — here is where everything changes. The ramp has
stayed as it was, and the semantics has taken other steps of it: in a dark
theme the step between neighbouring surfaces is wider, because a dark theme is
**not an inversion** but written by hand.

They read as a stack. Depth is carried by the **order of lightness and never by
a shadow**.

| Token | Work |
|---|---|
| `--surface-sunken` | The bottom: a recess, a track, the head of a table |
| `--surface-page` | The background of the document |
| `--surface-raised` | A panel, a card, the header of the shell |
| `--surface-overlay` | A popover, a menu, a tooltip, a modal |
| `--surface-field` | The recess under input. A role of its own rather than a reuse of `raised`: while a field was painted the colour of a panel, its only boundary was a border at 1.31:1 |
| `--surface-hover` `--surface-active` | Alpha over anything at all |
| `--surface-selected` | Selection. Alpha too — so that a hover over a selected row stays visible |

### Text

| Token | Threshold | What for |
|---|---|---|
| `--text-primary` | 4.5:1 | Data, headings, control labels |
| `--text-secondary` | 4.5:1 | The secondary: descriptions, units |
| `--text-muted` | 4.5:1 | Metadata and timestamps are **data**, and a readable threshold is due to them |
| `--text-faint` | 3:1 | **Not applied to text that gets read.** The threshold of decoration |

### Borders

Load-bearing and decorative are different tokens, because the requirements on
them differ.

| Token | When |
|---|---|
| `--border-subtle` | A divider inside an already bounded region |
| `--border` | A panel, a card, a cell: there is a step between surfaces beside it |
| `--border-strong` | A scrollbar, an accented seam |
| `--border-control` | A boundary that **is the control**: a checkbox, a field, the track of a switch. It has to hold 3:1 |

### Tones

A tone that carries meaning has three tokens: text, mark and background.

Text and mark parted for the same reason as with the accent. Step 5 is
calculated for reading, and while the tone stood in the label of a badge there
was no difference. Across a full-width fill — a history bar, the arc of a ring,
the ribbon of a meter — it reads darker and duller than the accent mark beside
it, and the palette falls into two registers. A mark makes do with 3:1, and on
step 4 every minimal mark of the library stands in one.

| Tone | Text (4.5:1) | Mark (3:1) | Background |
|---|---|---|---|
| accent | `--accent-text` | `--accent-mark` | `--accent-bg` |
| ok | `--ok-text` | `--ok-mark` | `--ok-bg` |
| warn | `--warn-text` | `--warn-mark` | `--warn-bg` |
| error | `--err-text` | `--err-mark` | `--err-bg` |

There is no third one — a "fill" between the text and the background. It lived
on step 4 and did not pass where a track ended up under it: a yellow meter fill
gave 2.49 against a norm of 3.0.

The accent has a case the statuses do not — a fill under a white label
(`--accent-solid` and `--accent-on`).

### Tone as an attribute

`data-tone` is set once on a group; the elements inside read `--tone-ink`,
`--tone-mark` and `--tone-bg` and know nothing else about the tone.

| Value | Means |
|---|---|
| `neutral` | Explicitly neutral. It also plays the part of info |
| `running` | Happening now. The carrier is the pulse, the tone is secondary |
| `ok` | Finished successfully |
| `warn` | Finished with remarks |
| `error` | Failed |

The vocabulary is closed. A sixth value is not started.

## Variants

A theme is **three independent knobs**: the hue of the neutral, the strength of
the cast and the depth of the dark surfaces. A seventh theme costs one line.

The light ones differ by the cast, the dark ones by depth, and that is not
carelessness. On a light background the eye catches temperature: warm paper and
cool paper read as different. On a dark one the temperature is barely visible,
while how far the panel parts from the page shows at once. The knobs differ
because different things are noticed in light and in dark.

| Attribute | `color-scheme` | `--hue-neutral` | `--tint` | And also |
|---|---|---|---|---|
| none | by the system | 75 | **0** | It follows the OS setting |
| `data-theme="light-neutral"` | `light` | 75 | **0** | The default: a pure grey |
| `data-theme="light"` | `light` | 75 | **2** | Warm |
| `data-theme="light-cool"` | `light` | 250 | **2** | Cool |
| `data-theme="dark-soft"` | `dark` | 75 | **0** | The stack rises by a step: exactly five tokens are overridden |
| `data-theme="dark"` | `dark` | 75 | **0** | The bottom of the ramp, with no overrides |

`--tint` is a multiplier rather than an absolute quantity: the steps of the ramp
keep their own shares of the cast (weaker at the light end, stronger at the
dark one), and at zero the whole row becomes grey at once. The neutral theme is
not a third hue but the absence of one: any `--hue-neutral` at zero chroma
gives one and the same grey.

```html
<html data-theme="dark">   <!-- forced -->
<html>                     <!-- by the system setting -->
```

The light depth is not adjustable. At the top of the ramp the steps run close
together, and a "less light" light theme is just a grey background.

### The ceiling of depth

There are two dark themes, and there will be no third one upwards. Above the
"grey" one the ramp stops being dark: at `--n-10` even `--text-muted` gives
3.29 against a norm of 4.5, and a translucent tone fill, through which the
panel shows, drops the label of a badge earlier still. That is exactly why
`dark-soft` overrides the status labels along with the four surfaces — they
move to step 2 of the ramp.

### The categorical palette

```html preview
<div class="ramp">
  <span class="ramp-step" style="--c: var(--chart-1)"></span>
  <span class="ramp-step" style="--c: var(--chart-2)"></span>
  <span class="ramp-step" style="--c: var(--chart-3)"></span>
  <span class="ramp-step" style="--c: var(--chart-4)"></span>
  <span class="ramp-step" style="--c: var(--chart-5)"></span>
  <span class="ramp-step" style="--c: var(--chart-6)"></span>
</div>
<div class="ramp ramp-scale">
  <span>1</span><span>2</span><span>3</span><span>4</span><span>5</span><span>6</span>
</div>
```

The lightness within the row **differs**. That keeps the series discernible
under colour blindness and in black-and-white print, where hue disappears
altogether. Squint — the order by lightness still reads.

Here colour encodes a series rather than a state — the one such place in the
library. Hence the separate rules.

| Token | Hue | Series |
|---|---|---|
| `--chart-1` | 292° | The first |
| `--chart-2` | 320° | The second |
| `--chart-3` | 355° | The third |
| `--chart-4` | 55° | The fourth |
| `--chart-5` | 115° | The fifth |
| `--chart-6` | 190° | The sixth |

- **The order is part of the contract.** Series no. 1 is always `--chart-1`,
  otherwise one and the same measure changes colour from screen to screen.
- **At least 25° from every status hue** (25, 85, 150): a categorical colour
  must not read as "it failed". Against the accent there is no angular
  guarantee — it is a knob with four positions (215, 250, 275, 45), and clay
  stands 10° from series no. 4. In place of a guarantee, a division of
  territory: the accent does not step inside a chart, a categorical colour does
  not step outside.
- **The lightness within the palette differs deliberately** — the series stay
  discernible under colour blindness and in black-and-white print, where hue
  disappears altogether.
- **Six is the ceiling.** A seventh series means this is not a chart but
  [a table](../components/display/table.md).
- **Forbidden outside a chart and its legend.** There is no categorical colour
  in the role of a status, of a button fill or of a row highlight.

### The rest

| Token | Work |
|---|---|
| `--focus-ring` | The focus ring, one for the whole library |
| `--track` | The empty track of a meter and a slider. Deliberately quiet: its job is to let the fill part from it, not to part from anything itself |
| `--scrim` | The ground under a modal. Denser in a dark theme. A dimming rather than a blur |
| `--shadow-color-near` `--shadow-color-far` | The colours of the two shadows — see [elevation](./elevation.md) |

## Rules

### What to take and what not to

| Use | Take instead |
|---|---|
| A semantic token: `--text-primary`, `--surface-raised`, `--ok-text` | **A step of a ramp directly** (`--n-3`, `--a-4`) — that is the light theme written into a component. The whole list of the semantics is in [tokens](./tokens.md) |
| The state of an object — the `data-tone` attribute, one vocabulary for the whole library | **A colour of your own for a state of your own** — start a tone and show it with [a badge](../components/display/badge.md): colour has no right to be the only carrier |
| A series of a chart — `--chart-1` … `--chart-6` in order | **A status hue as "one more series"** — the tones `ok`/`warn`/`error` are reserved. The series are labelled by [the legend](../components/charts/legend.md) |
| Depth — a step between surfaces | **A shadow for the sake of depth** — a shadow means "temporary and above", see [elevation](./elevation.md) |
| An informational message — the `neutral` tone | **Blue as "info"** — blue is taken by the accent and by the "running" state, see [the note in a card](../components/display/card.md) |

## API

```api
```

## Related

[Tokens](./tokens.md)
[Elevation and surfaces](./elevation.md)
[Badge](../components/display/badge.md)
[Chart legend](../components/charts/legend.md)
[Metric](../components/display/metric.md)
