---
title: Chip
group: Actions
layout: component
source: src/actions.css
js: Toggling, arrows and the roving tabindex are done by `instrument.js`
api:
  - { name: "inst-chips", kind: "class", doc: "The strip. It carries `role=\"listbox\"` and wraps across the width" }
  - { name: "inst-chip", kind: "class", doc: "One filter. It carries `role=\"option\"`" }
  - { name: "aria-orientation", kind: "attribute", value: "horizontal", doc: "On `.inst-chips`. Obligatory: the strip runs sideways, and a `listbox` is vertical by default" }
  - { name: "aria-multiselectable", kind: "attribute", value: "true", doc: "On `.inst-chips`. It switches several on at once: an arrow then moves only the focus, and space switches on" }
  - { name: "aria-selected", kind: "attribute", doc: "`true` · `false` on every chip. The carrier of the choice" }
  - { name: "data-value", kind: "attribute", doc: "On `.inst-chip`. What arrives in the `detail` of `inst:select` instead of the text" }
  - { name: "aria-disabled", kind: "attribute", value: "true", doc: "The filter is unavailable: there is not one row under it, for instance" }
  - { name: "--control-h-sm", kind: "token" }
  - { name: "--radius-sm", kind: "token" }
  - { name: "--text-xs", kind: "token" }
  - { name: "--space-2", kind: "token" }
  - { name: "--space-3", kind: "token" }
  - { name: "--surface-recessed", kind: "token" }
  - { name: "--accent-text", kind: "token" }
  - { name: "--accent-border", kind: "token" }
  - { name: "--border", kind: "token" }
  - { name: "--hairline", kind: "token" }
---

A filter that is visible and is switched on by pressing. Several are on at
once — and that is the difference from [the segmented
control](./segmented.md), where exactly one value is always chosen.

```html preview
<div class="inst-chips" role="listbox" aria-multiselectable="true"
     aria-orientation="horizontal" aria-label="Filter by state">
  <span class="inst-chip" role="option" aria-selected="true" tabindex="0">running</span>
  <span class="inst-chip" role="option" aria-selected="true" tabindex="-1">failed</span>
  <span class="inst-chip" role="option" aria-selected="false" tabindex="-1">queued</span>
  <span class="inst-chip" role="option" aria-selected="false" tabindex="-1">done</span>
  <span class="inst-chip" role="option" aria-selected="false" aria-disabled="true" tabindex="-1">cancelled</span>
</div>
```
```js target
chips.addEventListener('inst:select', (e) => {
  const on = e.detail.selected ? 'on' : 'off';
  console.log(`filter ${e.detail.value} ${on}`);
});
```

## Contract

| What | Required | Why |
|---|---|---|
| `role="listbox"` on the strip | yes | Without it `role="option"` is invalid, and the arrows will not work |
| `aria-label` on the strip | yes | "Filter by state": otherwise a screen reader announces a list with no name |
| `aria-selected` on every chip | yes | The carrier of the choice. `false` is written too — the absence of the attribute means "no choice is provided here" |
| `tabindex` | yes, the initial one | Zero on one, `-1` on the rest. After that it is led by [instrument.js](../../foundations/behavior.md). Without a roving `tabindex` `Tab` goes through every chip and the strip stops being one control |
| `aria-orientation="horizontal"` | yes | The strip runs sideways, and a `listbox` is vertical by default: without the attribute the traversal would hang on `↑` `↓` while the row runs to the right |
| `aria-multiselectable="true"` | no | Set when several are on at once. Without it the strip behaves as a single-choice list: the selection follows the focus |

The keyboard: `Tab` enters and leaves the strip, `←` `→` traverse it, space
switches on and off.

A selected chip carries **two** marks rather than one colour: the colour of the
label and the fill. Hover takes a third channel — the border — and so does not
argue with the selection. The label of an unselected one is `--text-secondary`,
of a selected one `--accent-text`, and both are checked at 4.5:1 in five themes
and four accents. The tap target is the height `--control-h-sm` with the
strip's gap of `--space-2`, checked by `cmd/targets` in all fifteen
combinations of scale and density.

The `inst:select` event arrives with `detail.selected` — switched on or off. A
cancelled event leaves the markup untouched: the application has taken the
state on itself and will set `aria-selected` itself.

## States

A chip **has no tone**. It names a filter rather than a state: painting
"failed" red would say the filter itself is the emergency. Colour in a chip
means exactly one thing — whether it is on.

Three channels, and they do not overlap, so they read together:

| Channel | What it encodes |
|---|---|
| The colour of the label | selected — `--accent-text`, not — `--text-secondary` |
| The colour of the border | under the cursor — `--accent-border`, otherwise — `--border` |
| The fill | present when selected or under the cursor |

The weight deliberately does not change on selection: the strip wraps across
the width, and a label grown thicker would move its neighbours onto another
line at the moment of the press.

## JS

The module is included [once per page](../../foundations/behavior.md) — there
is no need to initialise the components one by one.

### What `instrument.js` does

The examples on this page are live: enter the strip with `Tab`, then `←` and
`→`, switching on with space. With a mouse too.

The strip is declared as `role="listbox"`, and `instrument.js` carries out its
contract: `←` `→` between the chips, `Home` and `End`, one `Tab` for the whole
strip.

### Single and multiple

The `aria-multiselectable` attribute changes not only the declaration but the
behaviour of the keyboard, and that is a requirement of
[APG](https://www.w3.org/WAI/ARIA/apg/patterns/listbox/) rather than a decision
of the library.

| | arrow | space |
|---|---|---|
| Without the attribute | moves the focus **and the selection** | nothing new |
| `aria-multiselectable="true"` | moves **the focus only** | switches on and off |

Otherwise going through the strip with arrows would mean switching on
everything passed along the way.

### Events

`inst:select` bubbles from the toggled chip. In `detail.selected` is whether it
became on or off; the code is on the JS tab of the example in the header.

### Options

| Attribute | What it does |
|---|---|
| `data-value` | The value in `detail` instead of the chip's label |
| `aria-orientation="vertical"` | Arrows on the vertical. For a strip set in a column |

## API

```api
```

## Related

```related
```
