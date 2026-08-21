---
title: Slider
group: Inputs
layout: component
source: src/forms.css
js-optional: The number beside the thumb and the fill of the track while dragging. With no script the slider works entire, and the fill stands at the value from the markup
api:
  - { name: "inst-slider", kind: "class", doc: "The only class. Put on an `input[type=range]`" }
  - { name: "--fill", kind: "variable", value: "0%", doc: "The share of the track that is filled. **A channel of data**: the markup declares the initial value, and `instrument.js` keeps it while dragging" }
  - { name: "--size-track", kind: "token" }
  - { name: "--size-thumb", kind: "token" }
  - { name: "--size-icon", kind: "token" }
  - { name: "--track", kind: "token" }
  - { name: "--border", kind: "token" }
  - { name: "--hairline", kind: "token" }
  - { name: "--accent-mark", kind: "token" }
  - { name: "--surface-raised", kind: "token" }
  - { name: "--radius-full", kind: "token" }
  - { name: "--dur-1", kind: "token" }
  - { name: "--ease-out", kind: "token" }
---

A native `input[type=range]` brought to the metrics of the library. The value
is chosen by a position on a scale: a slider answers "more or less" rather than
"how much exactly".

```html preview
<input class="inst-slider" type="range" min="0" max="16" value="6" step="1" style="--fill:37.5%"
       id="agents" aria-label="Parallel agents">
```

## Contract

| What | Required | Why |
|---|---|---|
| `type="range"` | yes | The whole behaviour comes from the platform: the keyboard, the step, the announcement of the value. The class styles the native pseudo-elements and on another type does nothing |
| `min` and `max` | yes | With no bounds the scale is uninterpretable, and the platform derives `aria-valuemin/max` from them itself |
| A `<label for>` or an `aria-label` | yes | A slider has no text of its own |
| `step` | no | The default is `1`. Set it explicitly when the step is not one |
| `style="--fill: 45%"` | no | The fill of the track up to the thumb. With a script it is computed by itself; with no script declare the share in the markup, otherwise the track is drawn empty against a non-empty value. This is **a channel of data** rather than styling — as `inline-size` is on [a meter](../charts/meter.md) |
| An `<output for>` for the number | no | A semantic tie instead of a label changing at every step |

A slider stretches to **100% of the width of its container**. In a row with
other elements give it an explicit `flex` and `min-inline-size: 0`, otherwise
it pushes its neighbour onto the next line.

The keyboard, the `slider` role and `aria-valuenow` are the platform's. By hand
only `aria-valuetext` is wanted where the number alone is unclear — "6 agents"
rather than "6".

The thumb's target is 14px, **below the 24px of WCAG 2.5.8**: the press is taken
by the track across the whole height of `--size-icon`, and the position is
reachable from the keyboard. Where a slider is the only way to set a value, put
[a number field](./num-field.md) beside it — precision is not a slider's job.

## Anatomy

**The track is filled up to the thumb.** With no fill a slider says only "there
is a slider here": the position of a point on an even bar has to be judged by
eye, while a share is read as a shape. The empty part takes `--track` plus an
inner hairline, the filled one `--accent-mark`.

The thumb carries **a ring the colour of the surface** and therefore stays
distinguishable where it runs over the filled part of its own track.

| Part | The size token | A remark |
|---|---|---|
| The track | `--size-track`, 4px | `--track` plus a `--border` hairline |
| The thumb | `--size-thumb`, 14px | At density `compact` 12px, at `comfortable` 16px |
| The total height | `--size-icon`, 16px | So that a slider stands in a property row by the same metric as a glyph |

A press grows the thumb to `1.15` — a response is to be felt rather than
examined.

The fill is made of **two solid zones with a shared boundary** rather than the
proprietary `::-moz-range-progress` — that has no counterpart in Chromium, and
the track would fill in one browser of two. The stops coincide and there is no
transition between the colours, so this does not break the ban on gradients:
the same trick as the fill of [a sparkline](../charts/sparkline.md).

The share is held by `--fill`. CSS cannot read the `value` of an `input`, so it
is computed by [instrument.js](#js) — but **the initial value is declared by the
markup**, and with no script the track shows what the page arrived with rather
than zero.

## Behavior

### Outputting the value

:::note
**The number beside the thumb is kept by [instrument.js](#js).** Tying a
position to text by means of CSS is impossible, so the module does it: declare
an `<output for>` — the rest is its work.

Without the module the slider **works entire**: the value goes with the form
and a screen reader announces it at every step. Only the visible number does
not update.
:::

```html preview
<div class="inst-cluster">
  <input class="inst-slider" type="range" min="0" max="16" value="6" step="1" style="--fill:37.5%"
         id="sl" aria-label="Parallel agents">
  <output for="sl">6</output>
</div>
```

An `<output for>` ties the number to the slider semantically, so there is no
need to repeat it in an `aria-label` — and it should not be: a label changing
at every step turns the announcement into noise.

## JS

Include the module once per page — there is no need to initialise the
components one by one, `instrument.js` works by delegation and sees nodes that
arrived later.

```html
<script type="module" src="instrument.js"></script>
```

### What `instrument.js` does

Keeps the `<output for>` at the thumb's value, on input and on load. The tie is
declared by the `for` attribute already; synchronising it by hand would mean
writing one and the same line in every application.

### Events

There are none of its own — read the native `input`.

```js
sl.addEventListener('input', (e) => preview.setAgents(+e.target.value));
```

### Methods

`refresh(root)` sets the numbers again. Wanted if the markup arrived past the
observer — in a test, for instance.

## Composition

### In a property row

```html preview
<div class="inst-prop">
  <span class="inst-prop-label" title="Density of vegetation">Density</span>
  <span class="inst-prop-control">
    <input class="inst-slider" type="range" min="0" max="100" value="72" style="--fill:72%"
           aria-label="Density of vegetation">
  </span>
</div>
```

### With a label and a hint

```html preview
<div class="inst-field">
  <label class="inst-label" for="thr">Confidence threshold</label>
  <input class="inst-slider" type="range" min="0" max="100" value="80" step="5" style="--fill:80%"
         id="thr" aria-describedby="thr-hint">
  <span class="inst-field-hint" id="thr-hint">Below the threshold the agent asks for confirmation</span>
</div>
```

## API

```api
```

## Related

```related
```
