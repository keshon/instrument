---
title: Key–value list
group: Data display
layout: component
source: src/table.css
api:
  - { name: "inst-kv", kind: "class", doc: "A read-only list, on a `<dl>`" }
  - { name: "inst-kv--tight", kind: "class", doc: "The name column by its content. For narrow containers" }
  - { name: "inst-props", kind: "class", doc: "The inspector: a container of editable properties" }
  - { name: "inst-prop", kind: "class", doc: "One row of the inspector" }
  - { name: "inst-prop-label", kind: "class", doc: "The name of a property. It is truncated — a `title` is wanted" }
  - { name: "inst-prop-control", kind: "class", doc: "The cell of a control" }
  - { name: "--label-col", kind: "token" }
  - { name: "--pad-panel", kind: "token" }
  - { name: "--pad-cell-x", kind: "token" }
  - { name: "--gap-inline", kind: "token" }
  - { name: "--space-2", kind: "token" }
  - { name: "--space-3", kind: "token" }
  - { name: "--text-sm", kind: "token" }
  - { name: "--text-muted", kind: "token" }
  - { name: "--surface-hover", kind: "token" }
  - { name: "--control-h-sm", kind: "token" }
---

Name–value pairs of one object. A real `<dl>`: the tie between a name and its
value lives in the markup rather than in the layout alone.

```html preview
<dl class="inst-kv">
  <dt>Model</dt><dd>opus</dd>
  <dt>Started</dt><dd>14:32:07</dd>
  <dt>Files</dt><dd>4</dd>
</dl>
```

## Contract

| What | Required | Why |
|---|---|---|
| A real `<dl>` with `<dt>`/`<dd>` | yes | The tie between a name and its value is given by the markup. Two `<div>`s in two columns look the same and are tied by nothing |
| A `title` on `inst-prop-label` | yes, in the inspector | The name is truncated by the width of the column — without the full text the value is left with no name |
| A `<label for>` rather than a `<div>` in the inspector | yes | Otherwise the field has no name |

### Accessibility

| | |
|---|---|
| The tie between name and value | Given by `<dl>`/`<dt>`/`<dd>`. Two `<div>`s in two columns look the same and are tied by nothing |
| A truncated name | `inst-prop-label` is truncated by the width of the column, so it has to carry a `title` with the full text — otherwise the value is left with no name |
| The label of a control | In the inspector `inst-prop-label` is a real `<label for>` rather than a `<div>`: otherwise the field has no name |
| The order of reading | It matches the markup. A grid layout does not change the order of speaking, and it should not be changed visually either |

## Variants

### The inspector

`inst-kv` shows values, `inst-props` edits them.

```html preview
<div class="inst-props">
  <div class="inst-prop">
    <label class="inst-prop-label" title="Token limit" for="mx">Token limit</label>
    <div class="inst-prop-control">
      <input class="inst-input inst-input--sm" id="mx" value="500000">
    </div>
  </div>
  <div class="inst-prop">
    <label class="inst-prop-label" title="Parallel agents" for="pa">Parallel agents</label>
    <div class="inst-prop-control">
      <input class="inst-input inst-input--sm" id="pa" value="6">
    </div>
  </div>
</div>
```

The name column of both is one and the same `--label-col`: two ways of showing
a name–value pair have no right to part in their metrics.

### A narrow container

```html preview
<div style="max-width: 235px">
  <div class="inst-card">
    <div class="inst-card-head"><span class="inst-card-title">email</span></div>
    <dl class="inst-kv inst-kv--tight">
      <dt>Sent</dt><dd>1</dd>
      <dt>Suppressed</dt><dd>0</dd>
    </dl>
  </div>
</div>
```

A shared column is a promise that the rows are aligned **between** containers,
and it is worth its width. In a card of 235px a fixed `--label-col` ate four
tenths of it, and the value "1" ended up in the middle of emptiness. Inside a
card there is no such promise: what is compared is the rows within it rather
than with the one next door.

The other side is honest: two cards side by side will **not** align with each
other. If that is wanted, take the ordinary `inst-kv` and live with the width.

| | `inst-kv` | `inst-props` |
|---|---|---|
| Values | Read-only | Editable |
| The carrier | `<dl>` / `<dt>` / `<dd>` | `<div>` + `<label for>` |
| The name is tied | By the markup of the list | By the control's label |

## API

```api
```

## Related

```related
```
