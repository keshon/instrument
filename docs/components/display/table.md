---
title: Table
group: Data display
layout: component
source: src/table.css
js: Sorting and row actions. "Select all" is done by `instrument.js`
api:
  - { name: "inst-table", kind: "class", doc: "The base one, on a `<table>`" }
  - { name: "inst-num", kind: "class", doc: "A numeric cell: aligned to the end" }
  - { name: "inst-sort", kind: "class", doc: "The heading button of a sort" }
  - { name: "inst-col-select", kind: "class", doc: "The column of checkboxes, its width from the control" }
  - { name: "inst-col-actions", kind: "class", doc: "The column of actions, its width from the content" }
  - { name: "inst-row-actions", kind: "class", doc: "The wrapper of actions: it appears on hover and on focus" }
  - { name: "inst-table--cells", kind: "modifier", doc: "A cell as a tile: every value in a box of its own, with a gap between the boxes" }
  - { name: "inst-table--zebra", kind: "modifier", doc: "Stripes instead of rules. For a wide table: on fifteen columns a rule does not show WHERE you are within a row" }
  - { name: "inst-table--sticky", kind: "modifier", doc: "A sticky head" }
  - { name: "inst-table--sticky-col", kind: "modifier", doc: "A sticky first column" }
  - { name: "aria-sort", kind: "attribute", value: "ascending · descending", doc: "on a `<th>`" }
  - { name: "aria-selected", kind: "attribute", value: "true · false", doc: "on a `<tr>`" }
  - { name: "--table-head-bg", kind: "variable", value: "--surface-raised" }
  - { name: "--text-sm", kind: "token" }
  - { name: "--pad-cell-x", kind: "token" }
  - { name: "--row-pad-y", kind: "token" }
  - { name: "--space-3", kind: "token" }
  - { name: "--border", kind: "token" }
  - { name: "--border-subtle", kind: "token" }
  - { name: "--hairline", kind: "token" }
  - { name: "--text-muted", kind: "token" }
  - { name: "--surface-selected", kind: "token" }
  - { name: "--surface-hover", kind: "token" }
  - { name: "--control-h-md", kind: "token" }
  - { name: "--z-sticky", kind: "token" }
---

Comparable records in columns. A real `<table>` — with the semantics of
headings, a screen reader's navigation over cells, and printing.

```html preview
<table class="inst-table">
  <thead>
    <tr><th>Run</th><th>Model</th><th class="inst-num">Time</th></tr>
  </thead>
  <tbody>
    <tr><td>worldgen-01</td><td>opus</td><td class="inst-num">2.1 s</td></tr>
    <tr><td>biomes-04</td><td>sonnet</td><td class="inst-num">8.4 s</td></tr>
  </tbody>
</table>
```

## Contract

| What | Required | Why |
|---|---|---|
| A real `<table>` with `<th>` in `<thead>` | yes | A screen reader speaks the heading when moving over cells — without it a table turns into a set of numbers |
| A `<caption>` or an `aria-label` | yes | Two tables on a screen with no names are indistinguishable by ear |
| `inst-num` on numeric cells | yes, on numbers | Alignment to the end of the line: without it the digits cannot be compared by eye |
| A scrolling region reachable from the keyboard | yes, when it scrolls horizontally | Otherwise part of the columns exists for the mouse only |

:::note
**The division of labour.** "Select all" is done by
[instrument.js](#js), the partial selection included. Sorting and row actions
belong to the application layer: the order of rows is data, and the library
does not know it.
:::

### Accessibility

| | |
|---|---|
| Headings | `<th>` in `<thead>`. A screen reader speaks the heading when moving over cells — without it a table turns into a set of numbers |
| Sorting | `aria-sort` on the `<th>` plus a `<button>` inside. The attribute with no button is unreachable, the button with no attribute does not speak the direction |
| Selection | `aria-selected` on the `<tr>`. A checkbox carries a name of its own, otherwise the column is spoken as "tick, tick, tick" |
| Row actions | They appear on `:focus-within` rather than on `:hover` alone. Otherwise they exist for the mouse only |
| The caption of the table | A `<caption>` or an `aria-label`. Two tables on a screen with no names are indistinguishable by ear |
| Not colour alone | A selected row carries `aria-selected` rather than a highlight alone |
| Horizontal scrolling | The scrolling region has to be reachable from the keyboard |

## Anatomy

### Numeric columns

```html
<th class="inst-num">Time</th>
<td class="inst-num">2.1 s</td>
```

`inst-num` aligns to the end of the line. Tabular figures are on by default for
the whole library, so the digits stand one under another with no extra effort.

### Totals

```html
<tfoot>
  <tr><td>Total</td><td></td><td class="inst-num">10.5 s</td></tr>
</tfoot>
```

A `<tfoot>` goes **after** the `<tbody>` in the markup and is set off above. A
row of totals is not the last row of the data, and confusing them is not
allowed.

## States

### Selecting rows

```html preview
<table class="inst-table">
  <thead><tr><th class="inst-col-select"><label class="inst-checkbox"><input type="checkbox" aria-label="Select all"></label></th><th>Run</th></tr></thead>
  <tbody>
    <tr aria-selected="true">
      <td class="inst-col-select"><label class="inst-checkbox"><input type="checkbox" checked aria-label="Select worldgen-01"></label></td>
      <td>worldgen-01</td>
    </tr>
    <tr>
      <td class="inst-col-select"><label class="inst-checkbox"><input type="checkbox" aria-label="Select biomes-04"></label></td>
      <td>biomes-04</td>
    </tr>
  </tbody>
</table>
```

The state lives on the `aria-selected` of the `<tr>`. The highlight is
translucent, so the states add up: a hover over a selected row stays visible.

## Behavior

### Sorting

```html
<th aria-sort="ascending">
  <button class="inst-sort" type="button">Time</button>
</th>
```

| The value of `aria-sort` | What it shows |
|---|---|
| no attribute | The column is not sorted |
| `ascending` | Ascending, the glyph points up |
| `descending` | Descending, the glyph points down |

The direction is carried by the **glyph**, and `aria-sort` speaks the same
thing. The position of the rows alone is not enough. A sorting heading is a
real `<button>` inside the `<th>`, otherwise it cannot be reached from the
keyboard.

In a numeric column the glyph goes before the label, otherwise it is torn away
from it across the whole width of the cell.

### Row actions

```html
<td class="inst-col-actions">
  <div class="inst-row-actions">
    <button class="inst-btn inst-btn--sm inst-btn--icon" type="button" aria-label="Refresh">
      <svg class="inst-icon" aria-hidden="true"><use href="#i-refresh"/></svg>
    </button>
  </div>
</td>
```

The actions appear on hover — **and are always available from the keyboard**:
they also appear on `:focus-within` and on a selected row. An action visible
only under the cursor does not exist for the keyboard.

### A sticky head and column

```html
<table class="inst-table inst-table--sticky inst-table--sticky-col">
```

| Class | What sticks |
|---|---|
| `inst-table--sticky` | The head on vertical scrolling |
| `inst-table--sticky-col` | The first column on horizontal scrolling |

Stickiness works against the **scrolling ancestor** — usually
`inst-panel-body`. A table in an unattached flow will stick nothing.

## JS

Include the module once per page — there is no need to initialise the
components one by one, `instrument.js` works by delegation and sees nodes that
arrived later.

```html
<script type="module" src="instrument.js"></script>
```

### What `instrument.js` does

The example in the "States" section is live: tick the checkbox in the head —
every row is selected; untick one and it becomes indeterminate.

Only the "select all" checkbox in the head of the selection column:

- ticks and unticks every row except the `disabled` ones;
- keeps `aria-selected` on the `<tr>`;
- shows a partial selection through `indeterminate` — "some of the rows are
  selected" is inexpressible either by `checked` or by its absence;
- recomputes itself when a row is ticked on its own.

### Events

`inst:selectall` bubbles from the table, and its `detail` is `{ checked }`.
Cancelling leaves the rows as they are.

```js
table.addEventListener('inst:selectall', (e) => {
  bulkBar.hidden = !e.detail.checked;
});
```

### What is left to the application

| What | Why not `instrument.js` |
|---|---|
| Sorting | The order of rows is data. The library draws `aria-sort` and the arrow, the application rearranges the rows |
| Row actions | What "Stop" does is known to the application alone |
| Pagination and virtualisation | They depend on the source of the data |

## Patterns

### A table in a panel

A table shows itself only entire: sorting, row selection, actions and totals
look like four independent classes apart, and like one instrument together.

What to check by hand:

- **hover over a row** — the actions appear. Now go through the rows with
  `Tab`: the actions appear again. These are different events (`:hover` and
  `:focus-within`), and the second exists so that an action visible only under
  the cursor does not disappear for the keyboard;
- **select a row** and hover over it — the highlights add up, because the
  selection is translucent;
- **scroll the body of the panel** — the head stays where it is.

```html preview context
<div class="inst-panel">
  <div class="inst-panel-header">
    <span class="inst-panel-title">Modules of the run</span>
    <span class="inst-panel-actions"><span class="inst-badge">1 selected</span></span>
  </div>
  <div class="inst-panel-body inst-panel-body--flush">
    <table class="inst-table inst-table--sticky">
      <thead>
        <tr>
          <th class="inst-col-select"><label class="inst-checkbox"><input type="checkbox" aria-label="Select every row"></label></th>
          <th aria-sort="descending"><button class="inst-sort" type="button">Module</button></th>
          <th>State</th>
          <th class="inst-num"><button class="inst-sort" type="button">Time, s</button></th>
          <th class="inst-num"><button class="inst-sort" type="button">Tokens</button></th>
          <th class="inst-col-actions"><span class="inst-u-visually-hidden">Actions</span></th>
        </tr>
      </thead>
      <tbody>
        <tr aria-selected="true">
          <td><label class="inst-checkbox"><input type="checkbox" checked aria-label="Select Render"></label></td>
          <td>Render</td>
          <td><span class="inst-badge" data-tone="ok"><span class="inst-dot"></span>done</span></td>
          <td class="inst-num">18.2</td><td class="inst-num">61 204</td>
          <td class="inst-col-actions"><span class="inst-row-actions">
            <button class="inst-btn inst-btn--sm inst-btn--ghost inst-btn--icon" type="button" aria-label="Restart Render">
              <svg class="inst-icon" aria-hidden="true"><use href="#i-refresh"/></svg></button>
          </span></td>
        </tr>
        <tr>
          <td><label class="inst-checkbox"><input type="checkbox" aria-label="Select Physics"></label></td>
          <td>Physics</td>
          <td><span class="inst-badge" data-tone="ok"><span class="inst-dot"></span>done</span></td>
          <td class="inst-num">11.5</td><td class="inst-num">38 910</td>
          <td class="inst-col-actions"><span class="inst-row-actions">
            <button class="inst-btn inst-btn--sm inst-btn--ghost inst-btn--icon" type="button" aria-label="Restart Physics">
              <svg class="inst-icon" aria-hidden="true"><use href="#i-refresh"/></svg></button>
          </span></td>
        </tr>
        <tr>
          <td><label class="inst-checkbox"><input type="checkbox" aria-label="Select Audio"></label></td>
          <td>Audio</td>
          <td><span class="inst-badge" data-tone="error"><span class="inst-dot"></span>failed</span></td>
          <td class="inst-num">4.2</td><td class="inst-num">1 180</td>
          <td class="inst-col-actions"><span class="inst-row-actions">
            <button class="inst-btn inst-btn--sm inst-btn--ghost inst-btn--icon" type="button" aria-label="Restart Audio">
              <svg class="inst-icon" aria-hidden="true"><use href="#i-refresh"/></svg></button>
          </span></td>
        </tr>
      </tbody>
      <tfoot>
        <tr><td></td><td>Total</td><td></td><td class="inst-num">33.9</td><td class="inst-num">101 294</td><td></td></tr>
      </tfoot>
    </table>
  </div>
</div>
```

### A cell as a tile

```html preview
<table class="inst-table inst-table--cells">
  <thead><tr><th>Run</th><th>Branch</th><th>Model</th><th class="inst-num">Tokens</th></tr></thead>
  <tbody>
    <tr><td>nightly-4127</td><td>main</td><td>opus</td><td class="inst-num">186 000</td></tr>
    <tr><td>audio-pass</td><td>fix/mixdown</td><td>sonnet</td><td class="inst-num">42 300</td></tr>
    <tr><td>biomes-04</td><td>main</td><td>opus</td><td class="inst-num">98 700</td></tr>
  </tbody>
</table>
```

A rule parts a row from a row and says nothing about a **column**: on fifteen
columns the eye going to the right loses both the row and the boundary of the
value. A tile holds both at once — every value has a box of its own.

The head does not become a tile: it is not a value but the name of a column,
and a box around a name would make it equal to the data.

### Stripes instead of rules

```html preview
<table class="inst-table inst-table--zebra">
  <thead><tr><th>Run</th><th>Branch</th><th>Model</th><th class="inst-num">Tokens</th><th class="inst-num">Time</th></tr></thead>
  <tbody>
    <tr><td>nightly-4127</td><td>main</td><td>opus</td><td class="inst-num">186 000</td><td class="inst-num">21 s</td></tr>
    <tr><td>audio-pass</td><td>fix/mixdown</td><td>sonnet</td><td class="inst-num">42 300</td><td class="inst-num">4.2 s</td></tr>
    <tr><td>biomes-04</td><td>main</td><td>opus</td><td class="inst-num">98 700</td><td class="inst-num">11.5 s</td></tr>
    <tr><td>terrain-88</td><td>main</td><td>sonnet</td><td class="inst-num">150 400</td><td class="inst-num">18.2 s</td></tr>
  </tbody>
</table>
```

A rule under every row works while there are five columns. On fifteen it parts
a row from its neighbour but does not show **where you are within the row** —
the eye travels off on the vertical while going to the right. A fill on every
other row does show it: it runs the full width and does not disappear on
horizontal scrolling.

It may not be put together with the rules: two dividers on one boundary give a
grid. So the variant **removes** the rules rather than adding to them.

| Take stripes | Take rules |
|---|---|
| More than seven columns, the table is scrolled sideways | Up to five columns, the row is visible entire |
| The rows are read entire, left to right | The rows are compared with one another by one column |
| The table is long and stands on a surface of its own | The table is short and lies inside a card |

## API

```api
```

## Related

```related
```
