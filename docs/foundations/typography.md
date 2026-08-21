---
title: Typography
group: Foundations
layout: foundation
source: src/tokens.css
api:
  - { name: "inst-prose", kind: "class", doc: "Paragraph mode: 16/1.65, proportional figures, 68ch" }
  - { name: "inst-u-mono", kind: "class", doc: "A monospaced piece in the flow" }
  - { name: "inst-u-prose-nums", kind: "class", doc: "Proportional figures, pointwise" }
  - { name: "inst-u-dim", kind: "class", doc: "Muting to `--text-secondary`" }
  - { name: "inst-u-truncate", kind: "class", doc: "Truncating one line with an ellipsis" }
  - { name: "--font-sans", kind: "token" }
  - { name: "--font-mono", kind: "token" }
  - { name: "--text-2xs/xs/sm/md/lg/xl/2xl", kind: "token" }
  - { name: "--weight-normal", kind: "token" }
  - { name: "--weight-medium", kind: "token" }
  - { name: "--leading-tight/ui/prose", kind: "token" }
  - { name: "--tracking-tight", kind: "token" }
  - { name: "--tracking-normal", kind: "token" }
---

Seven steps and two weights. The base of the document is the **interface, not
prose**: 14/1.4, because labels, cells, buttons and task rows are the frequent
case.

```html preview
<div class="inst-stack inst-stack--tight">
  <span style="font-size:var(--text-2xl)">2xl · 28px · a hero number</span>
  <span style="font-size:var(--text-xl)">xl · 22px · a section heading</span>
  <span style="font-size:var(--text-lg)">lg · 18px · a block heading</span>
  <span style="font-size:var(--text-md)">md · 16px · prose, a panel name</span>
  <span style="font-size:var(--text-sm)">sm · 14px · the interface base</span>
  <span style="font-size:var(--text-xs)">xs · 12.5px · times, metadata</span>
  <span style="font-size:var(--text-2xs)">2xs · 11px · a badge, an axis label</span>
</div>
```

A step exists only if it is seen. The ratios of neighbouring sizes are
1.222 · 1.273 in the display part and 1.136 · 1.12 · 1.143 · 1.125 in the
working part: below that the eye stops telling sizes apart, and three steps
begin to work as one.

Here is the scale at work — interface, prose and a mark on one screen:

```html preview
<div class="inst-stack">
  <h2>A section heading</h2>
  <p>The interface by default: labels, buttons, cells.</p>
  <div class="inst-prose">
    <p>The body of an agent's reply is set in <code>.inst-prose</code>: 16/1.65,
    proportional figures and a line bounded at 68 characters.</p>
  </div>
  <span class="inst-badge" data-tone="ok"><span class="inst-dot"></span>done</span>
</div>
```

## Contract

The type size comes from the scale, the weight from three. A heading **has no
right** to shout louder than the data it labels.

```css
.my-label {
  font-size: var(--text-xs);
  color: var(--text-muted);
  line-height: var(--leading-ui);
}
```

| What | Required | Why |
|---|---|---|
| A size from `--text-*` | yes | Seven steps. An eighth, invented in place, breaks the ladder of headings |
| Leading by the work | yes | `--leading-ui` for the interface, `--leading-prose` for text read in paragraphs |
| Three weights | yes | `normal`, `medium`, `mono`. Weight 700 is forbidden in the library |
| No tracked capitals | yes | They shout louder than the number they label. `text-transform` does not occur in the library once |


### Accessibility

| | |
|---|---|
| Zoom | The sizes are in `rem`, so the browser's "default font size" setting works. The control heights are `min-block-size`, and the text does not hit a ceiling |
| The floor of the type size | 11px. There is no smaller size in the library and none is started |
| The threshold of reading | Any readable text stands on a token with a 4.5:1 threshold. `--text-faint` (3:1) is not applied to text — see [colour](./colors.md) |
| The measure of a line | Prose is bounded at 68 characters: a long line loses the start of the next one |
| Headings | The levels `h1`–`h6` are real. Type size does not replace structure: a screen reader reads the level, not the size |
| Truncation | `.inst-u-truncate` hides text visually but not from a screen reader. The full text is wanted in a `title` — that is what [the inspector](../components/display/kv.md) does |

## Scale

Seven steps, no more. Every extra one is a decision somebody will take wrongly.

| Token | Value | At a root of 16px | Ratio | What for |
|---|---|---|---|---|
| `--text-2xs` | `0.6875rem` | 11px | — | Badges, axis labels. **The minimum size**, we go no lower |
| `--text-xs` | `0.78125rem` | 12.5px | 1.136 | Metadata, times, hints |
| `--text-sm` | `0.875rem` | 14px | 1.12 | The interface base: labels, buttons, cells |
| `--text-md` | `1rem` | 16px | 1.143 | Prose, message bodies, a panel name |
| `--text-lg` | `1.125rem` | 18px | 1.125 | A block heading |
| `--text-xl` | `1.375rem` | 22px | 1.222 | A section heading |
| `--text-2xl` | `1.75rem` | 28px | 1.273 | Hero numbers only |


**A step exists only if it is seen.** In the working part of the scale the
ratios stay in the band 1.12–1.143. That is a threshold: below it the eye stops
telling sizes apart, and neighbouring steps begin to work as one — the whole
hierarchy has to be laid on colour and weight, and they have few registers. The
display steps part more widely, 1.222 and 1.273, and that too is deliberate: a
heading differs from a heading by a larger share than a caption from a caption.

The scale is counted from the bottom up, from the floor of 11px, rather than
from the root of the document. That `--text-md` coincides with 16px is a
convenient coincidence rather than a tie.

There are still two weights, and 700 is forbidden. So the large-scale hierarchy
is carried by size **alone** — all the more reason for it to be
distinguishable.

**Type size is the only tier in `rem`.** An interface has to respect an
increased default font size in the browser: people sit in an admin panel all
day, and that is the first setting somebody with weakened sight changes. Page
zoom works with `px` too, but the "default font size" setting is ignored by
`px` silently.

Everything else stays in `px`, because it is **geometry rather than text**: the
hairline, the state dot, the gutter, the radius. A consequence that cannot be
ignored: if the text grows while the control heights are set in `px`, the text
will one day hit the ceiling — so the heights are applied as `min-block-size`
rather than `block-size`.

### Headings

Defaults of the elements; they need no class. All carry `--weight-medium`,
`--leading-tight`, `--tracking-tight` and `text-wrap: balance`.

| Element | Size |
|---|---|
| `<h1>` | `--text-2xl` |
| `<h2>` | `--text-xl` |
| `<h3>` | `--text-lg` |
| `<h4>` `<h5>` `<h6>` | `--text-md` |

```html preview
<div class="inst-prose">
  <h2>A section heading</h2>
  <p>Inside <code>.inst-prose</code> a heading has a rhythm of its own: twice as
  much air above as below. A heading belongs to what is under it, and dividing
  the text evenly is not allowed to it.</p>
  <h3>A subheading</h3>
  <p>Without that rule a heading hangs midway between its own paragraph and
  somebody else's — and parts neither from the other.</p>
</div>
```

HEADINGS IN CAPITALS and Title Case are forbidden: an ordinary sentence only.

## Variants

Two. 700 exists in the font and is forbidden all the same: on a layout of
hairlines and air it reads as a shout.

**Medium is 600 rather than 500, and that is an admission of fact.** Segoe UI
has no true 500, and the system stack snaps it to Semibold: a measurement over
the pixels showed that 500 and 600 match line for line. On macOS SF Pro gives a
true Medium — that is, the drawing differed on two platforms, and what chose it
was the font rather than the author.

| Token | Value | Where |
|---|---|---|
| `--weight-normal` | 400 | Everything by default |
| `--weight-medium` | 600 | Headings, button labels, a table head, a badge |

```html preview
<div class="inst-stack inst-stack--tight">
  <span>400 — ordinary text, cells, field labels</span>
  <span style="font-weight:var(--weight-medium)">600 — a heading, a button, a table head</span>
  <span style="font-weight:700">700 — forbidden in the library: on a flat layout it reads as a shout</span>
</div>
```

### Leading

| Token | Value | What for |
|---|---|---|
| `--leading-tight` | 1.25 | Headings, a single-line interface |
| `--leading-ui` | 1.4 | The base of the document: wrapping labels, table cells |
| `--leading-prose` | 1.65 | Prose |

### Tracking

An optical correction rather than a style: large type looks loose by default.

| Token | Value | Where |
|---|---|---|
| `--tracking-tight` | `-0.01em` | Headings |
| `--tracking-normal` | `0` | Everything else |

Small type is never loosened: letterspacing at 11px is not an accent but the
loss of a word.

### Families

| Token | Stack |
|---|---|
| `--font-sans` | `ui-sans-serif`, `system-ui`, `-apple-system`, `Segoe UI`, `Roboto`, `Helvetica Neue`, `Arial` |
| `--font-mono` | `ui-monospace`, `Cascadia Code`, `JetBrains Mono`, `SF Mono`, `Consolas`, `Liberation Mono` |

The font is the system one: a web font would add a network request and a flash
of substitution for the sake of an effect that does not read in a flat tooling
interface.

The monospaced one comes with an optical correction of `0.925em` — mono is
perceived as larger than a sans at the same size. The correction is declared
**once** on `code`, `kbd`, `samp`, `pre` and `.inst-u-mono`: while the utility
had a rule of its own, a `<code class="inst-u-mono">` got `0.925 × 0.925` and
went under the floor of 11px. Inside a block of code the correction is already
applied by the container and is not laid on a second time.

### Tabular figures

`font-variant-numeric: tabular-nums` stands **on `body`**, that is, on
everything by default. In a tooling interface almost any number either stands
in a column or updates in place, and proportional figures tug at both cases.

| What | How to switch it off |
|---|---|
| A paragraph of prose | `.inst-prose` — it switches them off already |
| A single piece of text | `.inst-u-prose-nums` |

There is no switch the other way: tabular is everywhere anyway.

## Composition

Opt-in rather than the default. The base has to be the frequent case,
otherwise every component is forced to redeclare the size, and any `div` of an
application without a class inherits the wrong one.

```html
<div class="inst-prose">
  <p>The body of an agent's reply.</p>
  <p>The second paragraph steps away from the first by itself.</p>
</div>
```

| What `.inst-prose` does | Value |
|---|---|
| Size and leading | `--text-md` / `--leading-prose` |
| Figures | Proportional |
| The width of a line | `68ch` — the measure of readability |
| Wrapping | `text-wrap: pretty` on the paragraphs |
| Rhythm | The gap is set **between** neighbours: `> * + *` gets `--space-5` |
| Lists | An indent of `--space-7`, `--space-2` between the items |

There is no separate component for streaming and none is needed: the body is
`.inst-prose`, and the mark that text is still arriving is `.inst-caret`.

### Keys

A tool lives on the keyboard, so a key is an element of the interface rather
than a nameless monospaced piece. `<kbd>` is styled in the base and needs no
class.

```html preview
<span><kbd>Ctrl</kbd> <kbd>K</kbd> — search the runs</span>
<span><kbd>Esc</kbd> — stop the run</span>
```

## Rules

### What to take and what not to

| Use | Take instead |
|---|---|
| `--text-sm` — the interface default. You write nothing, it is already so | **Paragraphs read one after another** — the class `.inst-prose`, which is also the body of an agent's reply |
| `--text-2xl` — a hero number, one per block | **A screen title** — `--text-xl` through an `<h2>`. A large size on a label shouts louder than the number it labels |
| `--text-2xs` — a badge, an axis label | **Metadata and times** — `--text-xs`: 11px is a floor rather than a working size. An example is [a log row](../agent/log.md) |
| `--weight-medium` for an accent | **700** — it is not in the library. Part things by size and colour, as in [the metric](../components/display/metric.md) |
| Monospaced for identifiers and paths | **Monospaced for prose** — ordinary text stays a sans, see [the block of code](../components/display/code.md) |

## API

```api
```

## Related

[Colour](./colors.md)
[Density](./density.md)
[Block of code](../components/display/code.md)
[Metric](../components/display/metric.md)
[Log row](../agent/log.md)
