---
title: Segmented control
group: Actions
layout: component
source: src/actions.css
js: The choice, the arrows and the roving tabindex are done by `instrument.js`
api:
  - { name: "inst-segmented", kind: "class", doc: "The container. Bare `button` or `a` inside it; they need no class of their own" }
  - { name: "aria-checked", kind: "attribute", doc: "`true` · `false` on the buttons. Exactly one `true` per group" }
  - { name: "aria-current", kind: "attribute", value: "page", doc: "The carrier of the choice when the options are links" }
  - { name: "--control-h-md", kind: "token" }
  - { name: "--control-pad-md", kind: "token" }
  - { name: "--radius-md", kind: "token" }
  - { name: "--radius-sm", kind: "token" }
  - { name: "--space-1", kind: "token" }
  - { name: "--text-sm", kind: "token" }
  - { name: "--surface-sunken", kind: "token" }
  - { name: "--surface-raised", kind: "token" }
  - { name: "--border", kind: "token" }
  - { name: "--hairline", kind: "token" }
---

One choice among equal options. It always has exactly one chosen value — and
that is the whole difference from [the button group](./button-group.md), where
there is no chosen one and cannot be.

```html preview
<div class="inst-segmented" role="radiogroup" aria-label="View mode">
  <button type="button" role="radio" aria-checked="true"  tabindex="0">List</button>
  <button type="button" role="radio" aria-checked="false" tabindex="-1">Grid</button>
  <button type="button" role="radio" aria-checked="false" tabindex="-1">Table</button>
</div>
```

The [density](../../foundations/density.md) switch in the header of the
reference is this one.

## Contract

The state comes from the markup: the library draws it but does not compute it.
Without the roles `aria-checked` is **invalid**, and the control cannot be
reached from the keyboard.

| What | Required | Why |
|---|---|---|
| `role="radiogroup"` on the container | yes | A group is one control with one value |
| `aria-label` on the container | yes | Otherwise the control has no name: only the options have labels |
| `role="radio"` on every button | yes | A `role="tab"` with no panel is an invalid role |
| `aria-checked` on each | yes | The carrier of the state. The style of the chosen one hangs on it rather than on a class |
| `type="button"` | yes | Otherwise inside a form the button will submit it |
| A roving `tabindex` | yes, the initial one | `0` on the chosen one, `-1` on the rest. Without it `Tab` goes through every option and the group stops being one control |

The keyboard: `Tab` enters and leaves the group, `←` `→` `↑` `↓` traverse it,
and the choice follows the focus.

The chosen one carries **two** marks, and neither of them is colour: a raised
surface and a ring of border. The label of an unchosen one is
`--text-secondary`, checked at 4.5:1 in five themes. The height is `--space-2`
less than `--control-h-md`, and in the compact mode the target goes below 24px
under WCAG 2.5.8 — as everything of size `sm` does.

## Variants

### With an address of its own

A value that has an address does not stop being a value. A history window
survives a reload, goes into a bookmark and is handed to another person entire
— as a link; it has not stopped being a choice among equals for that. The
carrier of the state is then `aria-current="page"` rather than `aria-checked`.

```html preview
<nav class="inst-segmented" aria-label="History window">
  <a href="?window=day" aria-current="page">24 hours</a>
  <a href="?window=week">7 days</a>
  <a href="?window=month">30 days</a>
</nav>
```

## States

```html preview
<div class="inst-segmented" role="radiogroup" aria-label="The first one is chosen">
  <button type="button" role="radio" aria-checked="true"  tabindex="0">Chosen</button>
  <button type="button" role="radio" aria-checked="false" tabindex="-1">Ordinary</button>
  <button type="button" role="radio" aria-checked="false" tabindex="-1" disabled>Unavailable</button>
</div>
```

| State | How it is set | What happens |
|---|---|---|
| chosen | `aria-checked="true"` — or `aria-current="page"` on links | A raised surface and a hairline border. The weight does not change |
| ordinary | `aria-checked="false"` | The label is `--text-secondary`, with no background |
| hover | `:hover` | Only the colour of the label changes: the surface is taken by the chosen one |
| focus | `:focus-visible` | A ring inside the container, with no radius of its own |
| unavailable | `disabled` | Transparency. The option stays visible — **which one** is unavailable matters |

The chosen one carries **two** marks, and neither of them is colour: a raised
surface and a ring of border. Telling shades apart is not required.

The weight carries no state. Measured: with `medium` the width of the whole
track depended on which segment was chosen — 217.16, 216.94 and 216.64 pixels —
that is, the choice shifted the neighbours in the toolbar, and the rounded
right corner landed half a pixel off and re-rasterised. At a density of 2.0
that is a whole device pixel. Without `medium` all three positions give one
width.

## JS

The module is included [once per page](../../foundations/behavior.md) — there
is no need to initialise the components one by one.

### What `instrument.js` does

The examples on this page are live: enter the control with `Tab`, then `←` and
`→`. With a mouse too.

The control is declared as `role="radiogroup"`, and `instrument.js` carries out
its contract: `←` and `→` between the options, `Home` and `End`, one `Tab` for
the whole group. The `aria-checked` mark follows the focus and is carried over
on a click — across the whole group at once, because two checked options are a
state the markup will not get out of again.

### Events

`inst:select` bubbles from the chosen option.

```js
group.addEventListener('inst:select', (e) => {
  document.documentElement.dataset.density = e.detail.value;
});
```

### Options

| Attribute | What it does |
|---|---|
| `data-value` | The value in `detail` instead of the option's label |
| `aria-orientation="vertical"` | Arrows on the vertical |

```html
<div class="inst-segmented" role="radiogroup" aria-label="Density">
  <button type="button" role="radio" aria-checked="true"  data-value="sm" tabindex="0">Compact</button>
  <button type="button" role="radio" aria-checked="false" data-value="md" tabindex="-1">Default</button>
</div>
```

## Composition

### In a page header

```html preview context
<div class="inst-page-header">
  <div class="inst-page-header-main">
    <h2 class="inst-page-title">Runs</h2>
  </div>
  <div class="inst-page-actions">
    <div class="inst-segmented" role="radiogroup" aria-label="Period">
      <button type="button" role="radio" aria-checked="false" tabindex="-1">Day</button>
      <button type="button" role="radio" aria-checked="true"  tabindex="0">Week</button>
      <button type="button" role="radio" aria-checked="false" tabindex="-1">Month</button>
    </div>
  </div>
</div>
```

## API

```api
```

## Related

```related
```
