---
title: Design principles
group: About
source: docs/about/design-principles.md
---

The rules the components of instrument are derived from. This is not a list of
components: where a rule and a component disagree, the rule has priority — the
component is fixed.

The discipline is utilitarian minimalism: a Swiss grid plus Tufte, rendered
flat. The field of application is interfaces that show a human the work of a
machine: task queues, runs, logs, inspectors, states and uncertainty.

---

## Three laws

**1. Colour is meaning.** If a colour encodes neither a state, nor a category,
nor a polarity, it is removed. Four different colours on four bars measuring
one quantity are a false message.

**2. A component knows nothing of themes.** It reaches only for the semantics
and the roles (`--text-primary`, `--pad-panel`). A `--n-3` or a `#333` in a
component is a hard-coded light theme, a `14px` a broken density. A dark theme
is not an inversion: it is declared in the same lines as the light one, through
`light-dark()`.

**3. Hierarchy matters more than decoration.** What matters on a screen is the
data. Labels, borders, headers and icons serve it and have to be quieter. A
heading in tracked capitals sounds louder than the number it labels.

---

## The vocabulary

Four tiers. Every parameter exists in a single instance: a duplicate means a
future divergence.

| Tier | Where | The rule |
|---|---|---|
| 1. Ramps | `tokens.css` | 15 neutral steps + the accent + 3 statuses. They do not change between themes. A ramp is continuous: a step the semantics has not taken is a reserve |
| 2. Semantics | `tokens.css` | What a colour is for. One declaration per token through `light-dark()` |
| 3. Roles | `tokens.css` | Paddings, heights, the sizes of glyphs and gutters by purpose — and proportions, when one thing has to travel with another. The one tier density and scale retune |
| 4. Component | beside the component | `--btn-bg` and the like. One or two on a component with no variants; one per axis on a component with a matrix of variants |

A component may see tiers 2 and 3. Tier 1 never.

**The count of tier-4 variables is a measure rather than a norm.** The median
across the library is one per component, and that is the right number for a
component that has no variants. There are more of them exactly where there is a
matrix: a button has four weights × a tone × four states, and seven variables
are the cheapest record of such a matrix, because a variant stays two lines
instead of a rewritten rule. The mark of a breach is not the count but a
variable that answers for no axis: one started just in case and read in one
place is not a tier but a name for a constant.

### The rules of the tiers

- **A tone that carries meaning has exactly two tokens:** `--*-text` (the mark:
  text, a dot, the fill of a meter, the bar of a note) and `--*-bg` (the tinted
  background under it). There is one exception — the accent, which has a fill
  under a white label (`--accent-solid`).
- **A load-bearing border and a decorative one are different tokens.** If the
  boundary IS the control (a checkbox, a field, the track of a switch, a
  slider), it holds 3:1 and takes `--border-control`. If there is a step
  between surfaces beside it (a panel, a card), it takes `--border`.

  The list of load-bearing ones is closed, and the button is not on it: a
  button has a label inside, and what identifies it is the label rather than
  the outline. A button and a key (`kbd` is a picture of a key rather than a
  control) take `--border-strong`, that is 1.69:1 against 3.4:1 on a
  load-bearing one. Everything that is empty inside stays on the load-bearing
  one.

- **An interface link and a link in prose are different things.** In an
  interface an `<a>` inherits the colour and is underlined only under the
  cursor: it is recognised by its place and its shape, as a button is. The
  accent with a permanent underline lives in `.inst-prose`, where a link stands
  among text and there is nothing else to recognise it by.

- **A focus ring differs from what it outlines.** Holding 3:1 against the
  surfaces is not enough for it: under it there may be a solid fill.
  `--focus-ring` takes the step of the accent that is never a fill.

- **The surface of an annotation falls out of the stack.** A tooltip has an
  inverted surface (`--surface-inverse` / `--text-on-inverse`) — nothing else
  on the screen is that colour, and a hint is recognised before it is read. An
  inversion against the page rather than a fixed black: a black plate on a dark
  page disappears.

- **`--text-faint` is not applied to text that gets read.** It is the threshold
  of decoration (3:1) rather than of reading (4.5:1). A log timestamp is data,
  and `--text-muted` is due to it.

- **Text and a mark are different foregrounds.** A tone has `--tone-ink` (text,
  4.5:1) and `--tone-mark` (a dot, a caret, a thumb, the fill of a meter —
  3:1, but against the track as well). On the statuses they coincide at step 5.
  On the accent they part: `--accent-text` takes step 5, otherwise it will not
  reach 4.5:1 on white, and `--accent-mark` takes step 4, the same as the fill
  of a button. The mark and the fill have to coincide, otherwise two different
  blues stand side by side on one screen. In a dark theme they part: the fill
  cannot go up because of the white label, and the mark cannot go down because
  of the track.

- **`overscroll-behavior: contain` is set only where scrolling is
  guaranteed.** A container with `overflow: auto` that has nothing to scroll
  eats the mouse wheel under `contain`: the page under it does not move. That
  is, together with virtualisation rather than in advance.

- **A pseudo-element in a grid is placed explicitly.** A `::before` takes part
  in auto-flow on a par with the elements and comes first; one explicit
  `grid-column` on it breaks the auto-placement of the other cells. Either
  every cell is explicit or none is.

- **`overflow-wrap: anywhere` lives only on the cell that wants it.** It
  reduces the minimum content width to one character and on a container
  collapses the columns of a grid.

- **Defaults are declared through `:where()`.** A default with the same
  specificity as an override is not a default.

- **A gap is set between neighbours rather than around each.** A `margin-block`
  on every element adds up with the container's inner padding at the top and
  the bottom but not at the sides. The right form is
  `X + X { margin-block-start }`.

- **The order of states is written into the selector rather than into the order
  of the lines:** disabled > error > focus > hover. `:hover:not(:disabled)` has
  a specificity of (0,3,0) and `:user-invalid` (0,2,0), so hover yields
  explicitly, through `:not(:user-invalid)`.

- **Validation is shown by the platform: `:user-invalid` rather than
  `:invalid`.** `:invalid` burns red before any input, and the form meets the
  user with errors they have not made.

- **`readonly` and `disabled` are different things.** One cannot be changed but
  can be read and copied (a recess, the full colour of the text); the other is
  unavailable entire (an opacity of 0.5).

- **The hairline is switched on by pixel density.** A browser does not draw a
  border thinner than a physical pixel: on a 1x display `0.5px` and `1px` are
  the same thing.

### Forbidden

- Weight 700. `base.css` closes the one door the platform brings it through:
  `strong` and `b` are reset to `--weight-medium`.

  `--weight-medium` equals 600. Segoe UI has no true 500, and the system stack
  snaps it to Semibold, while SF Pro on macOS gives a true Medium — at a value
  of 500 the drawing would differ between platforms along the busiest axis of
  the hierarchy. There is no font of our own: it would cost the promise of "no
  build and no dependencies".
- HEADINGS IN CAPITALS and Title Case. An ordinary sentence only.
- Gradients, glows, blurs. A skeleton pulses by opacity rather than by a sheen.
  The boundary runs along the BLUR: what is forbidden is what has a radius. A
  hard hairline inwards — the edge of a solid button — is not a glow but an
  edge: it has neither spread nor halo, and it works by the same means a border
  does.
- A shadow on a card or a panel. A card gets a surface rather than a shadow.
- A rounding on a one-sided border. A `border-left` + a `border-radius` is
  always a bug.
- A second accent tone at the same time. Changing the accent entire is allowed:
  `data-accent` has four positions (petrol · graphite · indigo · clay), and all
  four pass the contrast gate on a par with the five themes.
- A font size below 11px.
- Colour as the only carrier of a state. A status travels with a dot and a
  word, a note carries an icon, a diff row carries a sign.
- A double Y axis on a chart. Two scales are two charts.
- A number in a component. A new value gets a token in the tier of roles.
- Colour inside a data URI. The shape is drawn by a mask and the colour comes
  as a token.
- `!important`. One exception, named: the `[hidden]` in `base.css` — that is
  correctness rather than styling.

### The categorical palette

An exception to the first law, agreed in advance: a series of a chart is a
category, so the exception refines the law rather than breaking it. The palette
has rules of its own:

- **The order is part of the API.** Series no. 1 is always `--chart-1`.
  Otherwise one and the same measure changes colour from screen to screen.
- **One series is not a category.** A single series takes `--accent-mark`
  rather than `--chart-1`: categorical colour exists to tell series apart.
- **The hues stand at least 25° from every status hue (25, 85, 150),** so that
  a categorical colour is not read as "it failed".

  Against the accent there is no such guarantee: it is a knob with four
  positions (215 · 250 · 275 · 45). In place of distance, a division of
  territory: the accent does not step inside a chart, a categorical colour does
  not step outside. They meet only in the legend of a multi-series chart, where
  lightness tells them apart.
- **The lightness within the palette differs,** so that the series stay
  distinguishable under colour blindness and in black-and-white print.
- **A categorical colour encodes a category, and only that.** It is allowed
  where the set of values is enumerable and closed: a series of a chart, the
  type of a syntax token, the type of a file. It is forbidden where a value
  means a state, a polarity or an interaction: a status, the fill of a button,
  the highlight of a row.

  Two conditions are not up for discussion there. Colour is never the only
  carrier: a token has a shape, a file an extension, a series a legend. And the
  threshold is taken by the work: the palette is calculated for a mark (3:1),
  and on text the text threshold (4.5:1) is wanted.
- **Six is the ceiling.** A seventh series means this is not a chart but a
  table.
- **A legend is obligatory at two series and more.**

### There is no informational tone

Blue is taken by the accent — by interaction and by the "running" state. A
third role would make it unreadable, and a second blue is forbidden. An
informational message takes the `neutral` tone: a muted fill and an ⓘ glyph.

### One overlap

The accent carries both interaction and the "running" state. Formally that
breaks the first law, and it is allowed because the carrier of "running" is the
pulse rather than the tone: nothing else in the library pulses. If something
else starts pulsing, the overlap will have to be parted.

---

## Rules of styling

**Tabular figures by default.** In a tooling interface almost any number either
stands in a column or updates in place, and proportional figures tug at both
cases. Prose switches that off through `.inst-prose`.

**The base of the document is the interface rather than prose:** 14/1.4.
Otherwise every component has to redeclare the size.

**Nested radii decrease.** An inner element with the same radius as its
container looks as if it bulges out.

**Radii are even.** A radius is multiplied by the pixel density, and at a
density of 1.5 an odd one gives half a device pixel (7 × 1.5 = 10.5): the arc
meets the straight edge off the grid, and the smoothing lands its own way on
every corner. The even 4 · 6 · 8 · 12 land whole at 1.0, 1.5 and 2.0.

**A control's padding is derived from its height.** A height and a padding set
as two independent ladders drift apart, and a button changes shape along with
density. The proportion is declared once (`--control-ratio-*`), and the padding
is computed from the height.

**The vertical of a row and the vertical of a block are different roles.**
`--row-pad-y` for rows, `--pad-block-y` for headers, toolbars, banners and
blocks of text. It is tighter vertically than horizontally: part of the
vertical air is brought by the leading.

**A press is a `scale(0.985)`.**

**Focus with a 1px gap and with no radius of its own.** A ring flush against a
control reads as a change of border. An outline follows the radius of the
element anyway, and a radius set on it changes the shape of the element itself.

**On hover the fill of a solid button moves away from the colour of the
label.** The label is white, so in both themes that is downwards: a mirrored
"lighter in the dark theme" drops the label's contrast below the threshold.

**A field stands in a recess rather than on the colour of the panel.**
Otherwise its only boundary is the border. The recess lies on the same step as
the muted surface in all five themes; what tells a field from a depression is
the load-bearing border rather than the lightness.

**One layout of labels per group rather than per field.** A label at the side
is in order where the values are of one kind and are scanned rather than read,
but it is applied to a form or a `fieldset`: otherwise one form gets two
different left edges. The layout is declared by the container
(`inst-form--side`, `inst-fieldset--side`) and the fields inherit it.

**The hint and the error in one slot.** The error replaces the hint rather than
shifting the form.

**The track of a slider is filled up to the thumb.** The position of a point on
an even bar is judged by eye, and a share is read as a shape. The share is
carried by `--fill` — a channel of data, as `inline-size` is on a meter.

**A shadow has two meanings, and they do not mix.** "It floats above and is
about to disappear" — a large blur and a noticeable offset, due to a popover
and a modal. "This is a physical surface" — a blur of one pixel, almost no
offset, due to a raised button. There is one test: does the shadow tear the
layer off or thicken it? Tears it off — an overlay, and there are exactly two.
Thickens it — a control. A card and a panel get neither.

**A button has four weights, and each is taken by one answer** to the question
of how loudly the button asks to be pressed:

| the weight | the look | where it stands |
|---|---|---|
| solid | a fill in the accent | **one** action per screen |
| raised (the default) | a fill, a hairline border, a shadow of a pixel | an action standing on its own |
| soft | a muted fill, with no border | a control inside a group: a toolbar, a row of icons, a chip |
| ghost | nothing until hover | an action people look for with their eyes |

There is no outlined weight: in a row of five the outlines add up into a grid,
and the group reads as a table. A fill with no outline gathers it into one
object, and the chosen element inside differs by the density of the fill.

**Danger is a tone rather than a fifth weight.** `--danger` goes over any of the
four and paints the label without adding a border. There is no solid red: the
fill under a white label belongs to the accent, and a second solid colour would
mean a second principal action on the screen.

**An icon in a button is quieter than the label.** There are two exceptions: on
a solid fill there is nothing to dim (the label is white), and on an icon
button the glyph IS the label.

**An icon button is square.**

**A unit of measure is smaller and quieter than the number.** "42 s" at one
size reads as one word; once size and colour are parted the eye grabs the
number and the "s" recedes.

**Direction and polarity are two different attributes.** An arrow down on the
time of a run is good, an arrow up on warnings is bad. A class named after the
arrow and meaning the judgement leads to mistakes of application.

**A row of metrics with no borders.** Four numbers are one group, and outlines
around each add up into a grid. A fill is allowed: it has no line, and air
parts the neighbouring tiles. A shared fill across the whole row, on the
contrary, makes it a seamless band in which the numbers stop reading as
separate quantities.

**Finished tasks step back rather than going out.** The history stays readable.
An answered approval request does the same.

**Indeterminate progress is a separate kind of bar.** An indeterminate one has
no `aria-valuenow`, and that is a message in itself.

**A folded output names its number in words.**

**The widths of a log's columns are set in `ch`.** The grid is declared on the
line, so every line has one of its own, and the level column would otherwise
change width from word to word.

**A disabled state is transparency rather than grey.** It shows which action is
unavailable.

**Reduced motion squeezes transitions and slows indicators of activity rather
than putting them out.** An animation squeezed to 0.01ms is a stop, and a
stopped indicator reports that the work has halted.

---

## The markup contract

The styles are tied to ARIA: the state lives in one place. A component with no
markup of its own counts as not applied.

| Component | The required markup |
|---|---|
| `.inst-segmented` | `role="radiogroup"` + `role="radio"` + `aria-checked` + a roving `tabindex` |
| `.inst-tabs` | `role="tablist"` + `role="tab"` + `aria-selected` + a roving `tabindex`; a tab with a panel takes `aria-controls` |
| `.inst-task` | the container `role="listbox"`, the row `role="option"` + `aria-selected` + a roving `tabindex` |
| `.inst-tree-item` | the container `role="tree"`, the node `role="treeitem"` + `aria-level` + `aria-expanded` + a roving `tabindex` |
| `.inst-meter` | `role="progressbar"` + `aria-valuenow/min/max`; an indeterminate one has no `aria-valuenow` |
| `.inst-step` | `<details>` / `<summary>` — the expansion and the keyboard come from the platform |
| `.inst-select` | the wrapper `.inst-select-wrap` (the chevron lives in it) |
| `.inst-num-field` | an `aria-label` with the axis on the `input`; `.inst-num-axis` is not a label |
| `.inst-prop-label` | a `title` with the full text, because the label is truncated |
| `.inst-failure` | a `role="alert"` and at least one action: with no way out this is not a failure block |
| `.inst-popover` | a `popover` on the block + a `popovertarget` on the button. The button becomes the implicit anchor — there is no need to start names of anchors |
| `.inst-menu` | `role="menu"` + `role="menuitem"`. The arrows and the roving `tabindex` are done by `instrument.js` |
| `.inst-tooltip-text` | `role="tooltip"` + an `id`, and the trigger carries an `aria-describedby`. Inside a panel with `overflow: hidden` it is clipped — there a `.inst-popover` is wanted |
| `.inst-dialog` | a native `<dialog>`; opening by `showModal()`, closing by buttons inside a `<form method="dialog">` |
| `.inst-accordion-item` | a `<details>`; a shared `name` makes the group mutually exclusive with no JS |
| any button | a `type="button"`, otherwise inside a form it will submit it |
| a busy button | an `aria-busy="true"`. The label stays in the markup: it holds the width and is spoken. `disabled` is not set — it throws the button out of the order of traversal. The mouse is removed through `pointer-events`, and protection from a second press is set by the handler |

The arrows, `Home` and `End` are done by
[`instrument.js`](../foundations/behavior.md). A role promises a keyboard, so
the keyboard has to work; a separate file is what keeps the CSS usable with no
script.

**The initial roving `tabindex` is written by the author** — everywhere the
group is visible at once: a segmented control, a queue, a tree, tabs. The
script moves it afterwards, but until it arrives a group with no `tabindex` is
not "without a keyboard for now" but a traversal that goes through every item
in a row. There is one exception, named in the table: the items of a menu lie
in a closed popover and do not reach the traversal until it is opened — there
the initial `tabindex` is set by the script.

Who owns what is written in
[`components.json`](../../components.json) in the `roving` field, and
`cmd/registry` checks it against the live examples of the pages.

### State and variant

- **A variant** is what the author of the markup chooses: a `--` modifier
  (`.inst-btn--primary`, `.inst-panel-body--flush`).
- **A state** is what changes over a lifetime: an attribute.

There are four attribute vocabularies, and all of them are closed. A class of
the form `.is-*` for a state is forbidden. A new value is started together with
a line in the table below, otherwise it silently does nothing.

The test that tells an axis from a duplicate: if two attributes give one and
the same colour, it is one attribute with two names.

**1. `data-tone` — the register of meaning. One vocabulary for the whole
library.**

`neutral` · `running` · `ok` · `warn` · `error`

It means the same thing everywhere: on a badge, a note, a banner, a dot, a log
row, the delta of a metric. It sets `--tone-ink`, `--tone-mark`, `--tone-bg`;
the component reads them and knows nothing else about the tone.

**2. `data-state` — the phase of a life. The vocabulary is its own on every
component.**

A queue row and a wizard step have different phases. Every set is listed here
entire, the base value included — the one that has no styling of its own.

| Component | The vocabulary | The base |
|---|---|---|
| `.inst-task` | `queued` `running` `done` `warn` `failed` `skipped` | `queued` |
| `.inst-step` | `running` `ok` `failed` | no attribute |
| `.inst-approval` | `pending` `approved` `denied` | `pending` — the only one with actions |
| `.inst-stepper-item` | `todo` `current` `done` | `todo` |
| `.inst-meter` | `indeterminate` | no attribute = determinate |

The base value is written in the markup though there are no rules for it: a
`data-state="queued"` is read and its absence is not, and a typo like `quued`
looks like the base.

**3. `data-kind` — a category rather than a state.** A diff: `add` · `del`. An
added line is not in the "ok" state — it belongs to the kind "addition". One
and the same `data-kind` is carried by the rows of a diff and by the numbers in
its header.

**4. `data-when` — who leads the showing.**
`.inst-field-error[data-when="invalid"]` is shown by the platform on
`:user-invalid`; the same block with no attribute is shown by the application,
from a server's reply, for instance.

**Busyness is on `aria-busy` rather than in these vocabularies.** A button has
no phase of life, it has a transient flag.

---

## Layers and responsiveness

The order of the layers is declared explicitly and therefore does not depend on
the order of the imports; the application's styles lie outside the layers and
always win. The full list with the explanation is in
[tokens](../foundations/tokens.md#the-order-of-layers).

### Responsiveness

Three levels, and the order between them is obligatory.

1. **Intrinsic.** The layout rearranges itself with not a single query: an
   `auto-fit` on a grid, a `flex-wrap` with a threshold on a split. It works
   always, including where there is no container ancestor. It is tried first.
2. **`@container`.** A component answers to the width of its own region rather
   than of the window: one and the same panel stands in the narrow column of an
   inspector and across the whole width of a dashboard. Declared as containers
   are `.inst-panel`, `.inst-container` and `.inst-shell-main` — regions whose
   width comes from outside. Declaring something a container when it takes its
   width from its content is not allowed: `container-type` switches off the
   intrinsic size on the inline axis and will collapse the element.
3. **`@media`.** One threshold in the whole library — 60rem, on the application
   shell.

### Spacing is set by primitives

There are no utilities of the form `mt-3`, `p-2`, `pt-0d25`. The scale is
sparse at the top so that "a little bigger" is not an available decision; a set
of spacing utilities would bring it back with its first class and move the
decision about rhythm into the markup of every screen.

In their place, three flow primitives (`.inst-stack`, `.inst-cluster`,
`.inst-grid`) with three steps of gap named by intent rather than by a number:
ordinary, `--tight`, `--loose`. Density retunes all three at once.

### Themes, scale and density

```html
<html data-theme="dark">               <!-- forced -->
<html>                                 <!-- by the system setting -->
<html data-accent="clay">              <!-- the accent: orthogonal to the theme -->
<html data-scale="16">                 <!-- scale: the base type size of the interface -->
<section data-density="compact">       <!-- density as an attribute of a container -->
```

Four knobs, and they do not intersect: 5 themes × 4 accents × 5 scales × 3
densities — 300 combinations. The default of the neutral is `--tint: 0`, a pure
grey; the cast belongs to the two named light themes `light` and `light-cool`,
where the multiplier is raised to 2.

Switching the theme changes `color-scheme`, and `light-dark()` in the semantics
starts handing out the second branch. A second block of tokens does not exist.

**Scale and density are different axes, and what tells them apart is the type
size.** Density moves geometry and does not touch the text: in all three modes
it stays 14px. Scale moves the type size together with the geometry, keeping
the drawing. The test is the same as for the attribute vocabularies: `18` +
`compact` gives a control height of 34px at a type size of 18, and the default
+ `comfortable` gives 36px at a type size of 14. The heights almost coincide,
the screens differ.

**Scale comes first.** On a single match the token goes to it: the
`[data-scale]` blocks stand after the density blocks at equal specificity.
Density takes its own back in the two-dimensional cell.

**A step is named by its base type size, and the ladder does not go below 14.**
That is the arithmetic of two rules multiplied by one another: the floor of the
type size is 11px, the minimum step of the ladder is 1.12, and under the base
lie two steps — so the minimum base equals 11 × 1.12² = 13.8px. "Smaller" stays
the work of density.

What each knob does is in [colour](../foundations/colors.md),
[scale](../foundations/scale.md) and [density](../foundations/density.md).
