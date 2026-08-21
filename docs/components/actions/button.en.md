---
title: Button
group: Actions
layout: component
source: src/actions.css
api:
  - { name: "inst-btn",          kind: "class",     doc: "The base one. On a `<button>` or an `<a>`. Always required" }
  - { name: "inst-btn-group",    kind: "class",     doc: "The container of a group: it removes the neighbour's border and rounds only the ends" }

  - { name: "inst-btn--primary", kind: "modifier", doc: "A solid fill in the accent. **One** principal action per screen" }
  - { name: "inst-btn--soft",    kind: "modifier", doc: "A muted fill, with no border. A control **inside a group**: a toolbar, a row of icons, a selection chip" }
  - { name: "inst-btn--stacked", kind: "modifier", doc: "Two lines: a quiet label above, the value below. For a control whose point is the value (repository, branch, environment)" }
  - { name: "inst-btn--field",   kind: "modifier", doc: "The look of a recess, the behaviour of a button. It opens a palette or a list; substituting a real `<input>` is not allowed — that is a lie to a screen reader" }
  - { name: "inst-btn-stack",    kind: "class",     doc: "The column of text inside `--stacked`. **Required**: what becomes a column is the text rather than the button, otherwise the icon travels above the label" }
  - { name: "inst-btn-label",    kind: "class",     doc: "The label inside `.inst-btn-stack`. Quieter and smaller than the value" }
  - { name: "inst-btn-value",    kind: "class",     doc: "The value inside `.inst-btn-stack`. A long one is truncated: this is an identifier, and wrapping it is not allowed" }
  - { name: "inst-btn--ghost",   kind: "modifier", doc: "No border and no fill. For secondary actions in a dense row" }
  - { name: "inst-btn--danger",  kind: "modifier", doc: "A destructive action: the tone of an error in the label. This is a TONE over a weight rather than a fifth weight — it combines with any of them" }
  - { name: "inst-btn--sm",      kind: "modifier", doc: "The height from `--control-h-sm`, the size `--text-xs`" }
  - { name: "inst-btn--lg",      kind: "modifier", doc: "The height from `--control-h-lg`, the size `--text-md`" }
  - { name: "inst-btn--icon",    kind: "modifier", doc: "Square, an icon only. It calls for an `aria-label`" }

  - { name: "type",       kind: "attribute", value: "button · submit", doc: "Required. Without it a button inside a form will submit it" }
  - { name: "aria-busy",  kind: "attribute", value: "true",            doc: "The action has started: the label dims and a ring turns in its place. The width does not change" }
  - { name: "aria-label", kind: "attribute",                           doc: "Required on `inst-btn--icon`: a button with no text has no accessible name" }
  - { name: "disabled",   kind: "attribute",                           doc: "Unavailable. **Not applied** for busyness — it throws the button out of the keyboard traversal" }

  - { name: "--btn-fg",       kind: "variable", value: "--text-primary",   doc: "The colour of the label" }
  - { name: "--btn-bg",       kind: "variable", value: "--surface-raised", doc: "The fill. Hover changes **only** this" }
  - { name: "--btn-bg-hover", kind: "variable", value: "--surface-hover",  doc: "The fill under the cursor" }
  - { name: "--btn-bg-active",kind: "variable", value: "--surface-active", doc: "The fill on a press" }
  - { name: "--btn-border",   kind: "variable", value: "--border",         doc: "The border. A hairline: what identifies a button is the label rather than the outline" }
  - { name: "inst-btn--outline", kind: "class", doc: "An outline: a quiet border, with no fill. The step between ghost and soft" }
  - { name: "data-tone", kind: "attribute", value: "running · ok · warn · error · neutral", doc: "The tone of the label from the kit's shared vocabulary. Set on the button itself" }
  - { name: "--btn-shadow",   kind: "variable", value: "none", doc: "The default has no shadow — the depth is carried by the recess. It is set on `--primary`: there, together with the edge, it gives thickness to a solid fill" }

  - { name: "--control-h-md",   kind: "token", value: "32px", doc: "The height. Retuned by density" }
  - { name: "--control-pad-md", kind: "token", value: "12px", doc: "The horizontal padding. Derived from the height through `--control-ratio-md`" }
  - { name: "--radius-md",      kind: "token", value: "8px",  doc: "The rounding of a free-standing control" }
  - { name: "--shadow-control", kind: "token", doc: "A shadow of one pixel. Not \"it floats above\" but \"it has thickness\"" }
  - { name: "--gap-inline",     kind: "token", value: "8px",  doc: "The gap between the icon and the label" }
  - { name: "--text-muted",     kind: "token", doc: "The colour of the icon: it is quieter than the label it accompanies" }
  - { name: "--text-sm",        kind: "token", value: "14px", doc: "The size of the label" }
  - { name: "--weight-medium",  kind: "token", value: "600",  doc: "The weight of the label" }
  - { name: "--size-spinner",   kind: "token", value: "14px", doc: "The ring of busyness" }
  - { name: "--dur-1",          kind: "token", value: "80ms", doc: "The duration of the fill transition" }
---

An action the user starts themselves. By default a button is quiet: a screen
where everything is `primary` has no `primary` at all.

```html preview
<button class="inst-btn inst-btn--primary" type="button">Run</button>
<button class="inst-btn" type="button">Cancel</button>
<button class="inst-btn inst-btn--soft" type="button">Filter</button>
<button class="inst-btn inst-btn--ghost" type="button">Details</button>
<button class="inst-btn inst-btn--danger" type="button">Delete</button>
```

## Contract

| What | Required | Why |
|---|---|---|
| `.inst-btn` | yes | The base class. On a `<button>` or on an `<a>` |
| `type` on a `<button>` | yes | Without it a button inside a form has `type="submit"` by default and submits the form — including one that was merely showing a hint |
| An `<a>` for going to an address | yes | A button leading to another page breaks the middle click, "open in a new tab" and copying the address. The class on a link gives exactly the same button |
| `aria-busy="true"` while it works | yes, if the button waits | It is spoken, and the label stays in place |
| `disabled` | the native one | Transparency does not drop the label below the threshold: it is read together with the background |
| A roving `tabindex` | **no** | `Tab`, `Enter` and `Space` come entirely from a native `<button>`; `instrument.js` intercepts nothing here |
| `aria-label` on `inst-btn--icon` | yes | A button with no accessible name is a button there is nothing to call |
| A `<button>` rather than an `<input>` for the look | yes | Substituting a field for a button is a lie to a screen reader about being able to type here, and a broken keyboard: a field does not open on space |

`:focus-visible` — the ring appears from the keyboard and does not appear from
the mouse. In a group the focused button rises by `z-index` so that its
neighbour does not cut the ring off.

The label of every variant is checked at 4.5:1 in five themes, and the
`primary` fill is bounded from above in lightness: the bound depends on hue and
is 0.545 for the default petrol. The tap target is 32px at `md` and 26px at
`sm`; at density `compact` `sm` gives 22px, which is **below the 24px of WCAG
2.5.8**, and so `compact` is applied to dense tabular views rather than to the
main navigation. At `prefers-reduced-motion` the ring of busyness slows down
but does not stop: a stopped indicator says "nothing is happening" while
something is.

## Variants

A variant is two lines: it changes only the inner variables rather than
rewriting the rule.

```html preview
<button class="inst-btn" type="button">Default</button>
<button class="inst-btn inst-btn--primary" type="button">Primary</button>
<button class="inst-btn inst-btn--soft" type="button">Soft</button>
<button class="inst-btn inst-btn--outline" type="button">Outline</button>
<button class="inst-btn inst-btn--ghost" type="button">Ghost</button>
<button class="inst-btn inst-btn--danger" type="button">Danger</button>
```

| Variant | When |
|---|---|
| no modifier | The default is **recessed**: a drop in lightness, with neither border nor shadow. The depth is counted from the ground, so the step is the same on a page, in a panel and in a modal |
| `inst-btn--primary` | The principal action of the screen. **One per screen** |
| `inst-btn--soft` | A control inside a group: a muted fill, with no border |
| `inst-btn--outline` | An outline: a quiet border, with no fill. An action that has to look pressable without laying claim to attention |
| `inst-btn--ghost` | Third-rank: neither border nor fill — text alone |
| `inst-btn--danger` | An irreversible or destructive action. **A tone rather than a weight**: it goes over any of the four and paints the label without adding a border |
| `data-tone` | The tone of the label from the shared vocabulary: `running` for an accented action with no chrome, `ok`, `warn`, `error`, `neutral`. It goes over any weight except `--primary` |

### The ladder of weights

Weight answers the question "how loudly does the button ask to be pressed", and
the answer is measurable. The measure is **presence**: area weighted by the
step in lightness against the ground. A fill works with the whole box, a border
only with the perimeter, which on a 100×32 button is 8% of the area — so a
border, at a higher contrast of line, comes out six times quieter than a fill
in total.

| Weight | light | dark |
|---|---|---|
| `--ghost` | 0 | 0 |
| `--outline` | 0.007 | 0.007 |
| `--soft` | 0.017 | 0.022 |
| default | 0.046 | 0.048 |
| `--primary` | 0.47 | 0.17 |

The ladder is monotone in all five themes, and neighbouring steps part by no
less than 0.022 — the same threshold the kit measures its stack of surfaces by.
That is not a wish: the pair "ladder: soft against default" is checked by
`cmd/contrast`.

The border of the outline is **decorative**. The kit already had an outlined
weight and removed it precisely because it rested on a load-bearing line at
3.46:1: the outline of a button turned out to be the loudest line on the
screen. On a quiet border that does not happen.

### Tone

A tone is the kit's shared vocabulary, the same one the badge, the dot, the
meter and the banner use. A button reads `--tone-ink` from it for the label and
`--tone-bg` for hover.

```html preview
<button class="inst-btn inst-btn--ghost" type="button" data-tone="running">Add context</button>
<button class="inst-btn" type="button" data-tone="ok">Accept</button>
<button class="inst-btn" type="button" data-tone="error">Delete</button>
<button class="inst-btn inst-btn--danger" type="button">Delete</button>
```

The last two lines give one and the same thing: `inst-btn--danger` is the tone
`error` written as a class. The class stayed because a destructive action
occurs often and deserves a short spelling.

A ghost with the tone `running` is that action-as-a-link: the label in the
accent, with neither border nor fill. It is what "add context", "add a row",
"create the first object" are set in — things that invite a press without
arguing with the principal button of the screen.

| What | Why |
|---|---|
| The tone is set **on the button itself** | `--tone-ink` is inherited, and a rule that merely read the variable would paint a button inside an emergency banner by itself. A tone is declared rather than picked up |
| `--primary` **takes no** tone | The label there stands on a fill in the accent and holds 4.5:1 only in white. A red or a green on that fill gives 1.5–2.0, so the combination is disarmed explicitly |
| A tone does not replace a weight | Weight answers how loudly the button asks to be pressed; tone answers what kind of action it is. The questions differ, so the axes differ |

### A switch of a value

There is a control whose point is the **value**, with the label only saying
what the value is of: the current repository, branch, environment, account.
"master" by itself means nothing, and "Current branch" cannot be removed.

```html preview context
<button class="inst-btn inst-btn--stacked" type="button" style="inline-size:15rem">
  <svg class="inst-icon" aria-hidden="true"><use href="#i-book"/></svg>
  <span class="inst-btn-stack">
    <span class="inst-btn-label">Current repository</span>
    <span class="inst-btn-value">instrument-ui-kit</span>
  </span>
  <svg class="inst-icon inst-icon--sm" aria-hidden="true"><use href="#i-chevron"/></svg>
</button>
```

What becomes a column is the **text** rather than the button. The whole button
was the column at first, and on the very first real screen that fell apart:
such a control almost always has a leading icon and a chevron at the end, and a
column would lay the icon above the label and the chevron below the value.
Hence `.inst-btn-stack` is required.

The single line of the base button still holds: a command has a short label,
and a wrap in one is almost always a sign of a poor label. Two lines are a
property of this variant rather than permission to wrap any button.

### A button in the clothes of a field

```html preview context
<button class="inst-btn inst-btn--field" type="button" style="inline-size:18rem">
  <svg class="inst-icon" aria-hidden="true"><use href="#i-list"/></svg>
  Choose a file or a command
</button>
```

It opens a palette or a list but stays a button. The look is taken from a
recess and the behaviour is a button's: the border here is load-bearing,
because it **is** the control, with one hint inside it.

### As a link

The class on an `<a>` gives exactly the same button, with no underline.

```html preview
<a class="inst-btn" href="#button">Open the report</a>
<a class="inst-btn inst-btn--primary" href="#button">Go to the run</a>
<button class="inst-btn" type="button">Run</button>
```

## Sizes

```html preview
<button class="inst-btn inst-btn--sm" type="button">sm</button>
<button class="inst-btn" type="button">md</button>
<button class="inst-btn inst-btn--lg" type="button">lg</button>
```

The height comes from the tier of roles (`--control-h-*`), so a button and a
[field](../inputs/input.md) of the same size standing beside it physically
cannot part. The density of a container retunes all three sizes at once.

| Size | Height | Padding | Type size |
|---|---|---|---|
| `inst-btn--sm` | `--control-h-sm` | `--control-pad-sm` | `--text-xs` |
| default | `--control-h-md` | `--control-pad-md` | `--text-sm` |
| `inst-btn--lg` | `--control-h-lg` | `--control-pad-lg` | `--text-md` |

## States

```html preview
<button class="inst-btn" type="button">Ordinary</button>
<button class="inst-btn" type="button" disabled>Unavailable</button>
<button class="inst-btn inst-btn--primary" type="button" aria-busy="true">Run</button>
```

| State | How it is set | What happens |
|---|---|---|
| hover | `:hover:not(:disabled)` | Only `--btn-bg` changes |
| press | `:active:not(:disabled)` | `--btn-bg` + `scale(0.985)` — a press is to be felt rather than examined |
| focus | `:focus-visible` | A ring with a 1px gap, with no radius of its own |
| unavailable | `disabled` | An opacity of `0.5` rather than a grey fill: the button keeps its identity, and it shows **which** action is unavailable |
| busy | `aria-busy="true"` | The label dims and a ring takes its place. See below |

### Busyness

```html preview
<button class="inst-btn inst-btn--primary" type="button" aria-busy="true">Run</button>
<button class="inst-btn" type="button" aria-busy="true">Cancel</button>
<button class="inst-btn inst-btn--icon" type="button" aria-busy="true" aria-label="Refresh">
  <svg class="inst-icon" aria-hidden="true"><use href="#i-refresh"/></svg>
</button>
```

Four decisions, each of which is easy to make the other way round and get a
defect:

- **The label is not removed from the flow but dimmed by colour.** The width of
  the button stays as it was. Otherwise the action bar twitches at the moment
  the user is looking at the cursor.
- **The label stays in the markup** — a screen reader reads it, and a fallback
  `aria-label` is not needed.
- **The button keeps its focus.** `disabled` would throw it out of the order of
  traversal right under the hands of somebody who came by keyboard.
- **Only the mouse is removed** (`pointer-events`). Protection from a second
  press belongs in the handler: that is behaviour rather than styling.

## With icon

```html preview
<button class="inst-btn" type="button">
  <svg class="inst-icon" aria-hidden="true"><use href="#i-refresh"/></svg>
  Refresh
</button>
<button class="inst-btn inst-btn--icon" type="button" aria-label="Refresh">
  <svg class="inst-icon" aria-hidden="true"><use href="#i-refresh"/></svg>
</button>
```

An icon button is **square**: a rectangle of 32×28 gives away inconsistent
sizes. The width is taken from the same control height, so the square holds at
all three sizes and at all densities.

:::warn
A button with no text has to carry an `aria-label`. The icon is then
`aria-hidden="true"` — otherwise a screen reader reads it twice or does not
read it at all.
:::

## Composition

### A group

```html preview
<div class="inst-btn-group">
  <button class="inst-btn inst-btn--sm" type="button">Left</button>
  <button class="inst-btn inst-btn--sm" type="button">Centre</button>
  <button class="inst-btn inst-btn--sm" type="button">Right</button>
</div>
```

Several **different actions** standing flush. Not to be confused with
[the segmented control](./segmented.md): there is one choice among equals
there, and it always has exactly one chosen value.

The seam of a group is the neighbour's border removed rather than a negative
margin of `-0.5px`: a browser never draws a border thinner than a physical
pixel, so a shift of half a pixel left a visible thickening at every joint.

## Patterns

### The action bar of a form

The principal action is last, cancel before it, and both are pushed to the far
edge.

```html preview context
<form class="inst-panel">
  <div class="inst-panel-body">
    <div class="inst-field">
      <label class="inst-label" for="p-name">Run name</label>
      <input class="inst-input" id="p-name" type="text" value="nightly-4127">
    </div>
    <div class="inst-form-actions">
      <button class="inst-btn" type="button">Cancel</button>
      <button class="inst-btn inst-btn--primary" type="submit">Save</button>
    </div>
  </div>
</form>
```

### The actions of a table row

A dense row takes `ghost` and `sm`: a border on every button in a table would
turn the column into a grid. A destructive action keeps its border — it is the
one thing in the row that has to be noticeable.

```html preview context
<div class="inst-cluster">
  <button class="inst-btn inst-btn--sm inst-btn--ghost" type="button">Open</button>
  <button class="inst-btn inst-btn--sm inst-btn--ghost" type="button">Retry</button>
  <button class="inst-btn inst-btn--sm inst-btn--danger" type="button">Delete</button>
</div>
```

## Customization

A variant of your own is two lines. What is overridden are the variables rather
than the rule.

```css
.my-btn--accent {
  --btn-fg: var(--accent-on);
  --btn-bg: var(--accent-solid);
}
```

Two limits that are easier to break than to notice:

- **Hover changes only `--btn-bg`.** Do not override `background` directly — the
  variable will stop being read, and the variant will lose its hover.
- **A fill under a white label cannot be lighter than the bound for its hue in
  OKLCH** — 0.545 at 215°, 0.568 at 250°, 0.576 at 275°, 0.580 at 45°. Above
  that white text stops holding 4.5:1. For the same reason `--btn-bg-hover`
  moves **away** from the colour of the label. A mirrored "lighter in the dark
  theme" drops the label's contrast to 2.61.

## API

```api
```

## Related

```related
```
