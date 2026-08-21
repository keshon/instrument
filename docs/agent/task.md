---
title: Queue row
group: Agent layer
layout: component
source: src/rows.css
js: The order and the states of the tasks belong to the application layer. The keyboard and the selection of a row are done by `instrument.js`
api:
  - { name: "inst-task", kind: "class", doc: "The row" }
  - { name: "inst-task-gutter", kind: "class", doc: "The left gutter for the dot" }
  - { name: "inst-task-main", kind: "class", doc: "The heading and the caption" }
  - { name: "inst-task-title", kind: "class", doc: "The name of the task" }
  - { name: "inst-task-sub", kind: "class", doc: "What exactly is being done" }
  - { name: "inst-task-meta", kind: "class", doc: "A time or a counter, pushed to the end" }
  - { name: "data-state", kind: "attribute", doc: "`queued` `running` `done` `warn` `failed` `skipped`" }
  - { name: "aria-selected", kind: "attribute", doc: "`true` · `false`" }
  - { name: "--size-gutter", kind: "token" }
  - { name: "--size-dot", kind: "token" }
  - { name: "--row-pad-y", kind: "token" }
  - { name: "--pad-cell-x", kind: "token" }
  - { name: "--gap-inline", kind: "token" }
  - { name: "--space-1", kind: "token" }
  - { name: "--text-sm", kind: "token" }
  - { name: "--text-xs", kind: "token" }
  - { name: "--surface-selected", kind: "token" }
  - { name: "--surface-hover", kind: "token" }
  - { name: "--tone-ink", kind: "token" }
---

One task in a queue of agents. The row is **flat**: it has no body, because
there are hundreds of them on a screen, and a border around each would turn the
queue into a grid.

```html preview
<div role="listbox" aria-label="The queue of agents">
  <div class="inst-task" role="option" data-state="done" aria-selected="false" tabindex="-1">
    <span class="inst-task-gutter"><span class="inst-dot"></span></span>
    <span class="inst-task-main">
      <span class="inst-task-title inst-u-truncate">Parsing the heightmap</span>
      <span class="inst-task-sub">terrain/heightmap.ts · 4 files</span>
    </span>
    <span class="inst-task-meta">2.1 s</span>
  </div>
  <div class="inst-task" role="option" data-state="running" aria-selected="true" tabindex="0">
    <span class="inst-task-gutter"><span class="inst-dot"></span></span>
    <span class="inst-task-main">
      <span class="inst-task-title inst-u-truncate">Generating the biomes</span>
      <span class="inst-task-sub">world/biomes.ts</span>
    </span>
    <span class="inst-task-meta">8.4 s</span>
  </div>
</div>
```

## Contract

| What | Required | Why |
|---|---|---|
| A `role="listbox"` outside, a `role="option"` on the row | yes | A queue is a list with a selection. Without the roles this is a stream of thirty `div`s |
| A `data-state` on the row | yes | The phase of the task. The dot in the gutter inherits its tone from here |
| A roving `tabindex` | yes | `0` on the selected one, `-1` on the rest: `Tab` enters the list once, and the arrows take over |
| An `inst-u-truncate` on the heading | yes | A long name otherwise pushes the row wide and drives the time past the edge |

### Accessibility

| | |
|---|---|
| The roles are required | The container `role="listbox"`, the row `role="option"`. Without them `aria-selected` is **invalid**, and the row is not spoken as selectable |
| The keyboard | A roving `tabindex`: `0` on the active one, `-1` on the rest. `Tab` enters the list once, then the arrows, `Home`/`End`. The behaviour belongs to the application |
| The state in words | A `data-state` paints the row but is not spoken. The state has to be in the text — in the caption or in [a badge](../components/display/badge.md) beside it |
| A truncated heading | `inst-u-truncate` cuts visually, and the full text stays in the DOM and is spoken. Add a `title` for the mouse |
| Reduced motion | The pulse of a `running` dot slows rather than going out: a queue of agents with no mark of work is useless |
| Selection and hover add up | The highlight of a selected row is translucent, so a hover over it stays visible |

## Anatomy

```html
<span class="inst-task-gutter"><span class="inst-dot"></span></span>
```

A left column of fixed width (`--size-gutter`). It holds the state dot on one
vertical across every row — without it headings of different lengths scatter
the markers, and the queue stops reading as a column.

## States

| `data-state` | What happens |
|---|---|
| `queued` | The base. There is no styling — the task waits |
| `running` | Going. The dot **pulses**, the heading takes the accent |
| `done` | Finished. The heading **steps back** to `--text-secondary` |
| `warn` | With remarks |
| `failed` | Failed |
| `skipped` | Skipped. An opacity of `0.7` |

**The finished ones step back rather than going out.** The history stays
readable while it stops arguing with what is going on now. Opacity is applied
to `skipped` alone — and it is soft: `0.55` dropped the caption below the
threshold of readability.

## JS

Include the module once per page — there is no need to initialise the
components one by one, `instrument.js` works by delegation and sees nodes that
arrived later.

```html
<script type="module" src="instrument.js"></script>
```

### What `instrument.js` does

The example in the header is live: enter the queue with `Tab` and go through
the rows with the arrows — the selection follows the focus.

The queue is declared as a `role="listbox"`, and `instrument.js` carries out
the contract of that role: one `Tab` for the whole list, the arrows inside,
`Home` and `End`, a traversal that wraps. The selection follows the focus and
is carried over on a click. Rows that arrived while the work goes on are
traversed on a par with the rest — the elements are found at the moment of the
keypress.

### Events

`inst:select` bubbles from the chosen row, with `{ value }` in the `detail`.

```js
queue.addEventListener('inst:select', (e) => {
  details.show(e.detail.value);
});
```

Cancelling leaves the `aria-selected` untouched: the application leads the
selection itself.

### Options

| Attribute | What it does |
|---|---|
| `data-value` | The identifier of the task in the `detail` instead of the text of the row |
| `aria-disabled="true"` | The row stays in the arrow traversal and stays visible but is not chosen — see [behaviour](../foundations/behavior.md) |

### What is left to the application

The order of the tasks, their states and cancellation are data. The library
draws a state but does not know when it changes.

## Patterns

One queue row shows almost nothing: the states read only beside one another.

What is visible only here:

- **the one that is going pulses** and the finished one does not. The pulse is
  the carrier of "running", and the tone is secondary;
- **the finished one steps back**: the heading goes to `--text-secondary` while
  staying readable. The history does not argue with what is happening now;
- **the skipped one** is the only one opacity is applied to, and it is soft:
  `0.55` dropped the caption below the threshold of readability;
- **the dots stand on one vertical** whatever the length of the headings. That
  is the work of the gutter of fixed width.

```html preview context
<div class="inst-panel">
  <div class="inst-panel-header">
    <span class="inst-panel-title">The queue of agents</span>
    <span class="inst-panel-actions"><span class="inst-badge">7 of 12</span></span>
  </div>
  <div class="inst-panel-body inst-panel-body--list">
    <div class="tree-none" role="listbox" aria-label="The queue of agents" id="queue">
      <div class="inst-task" role="option" data-state="done" aria-selected="false" tabindex="-1">
        <span class="inst-task-gutter"><span class="inst-dot"></span></span>
        <span class="inst-task-main"><span class="inst-task-title inst-u-truncate">Parsing the heightmap</span>
          <span class="inst-task-sub">terrain/heightmap.ts · 4 files</span></span>
        <span class="inst-task-meta">2.1 s</span>
      </div>
      <div class="inst-task" role="option" data-state="done" aria-selected="false" tabindex="-1">
        <span class="inst-task-gutter"><span class="inst-dot"></span></span>
        <span class="inst-task-main"><span class="inst-task-title inst-u-truncate">Generating the biomes</span>
          <span class="inst-task-sub">world/biomes.ts · 11 files</span></span>
        <span class="inst-task-meta">8.4 s</span>
      </div>
      <div class="inst-task" role="option" data-state="running" aria-selected="true" tabindex="0">
        <span class="inst-task-gutter"><span class="inst-dot"></span></span>
        <span class="inst-task-main"><span class="inst-task-title inst-u-truncate">Placing the resources</span>
          <span class="inst-task-sub">the third pass is going<span class="inst-caret"></span></span></span>
        <span class="inst-task-meta">14.0 s</span>
      </div>
      <div class="inst-task" role="option" data-state="failed" aria-selected="false" tabindex="-1">
        <span class="inst-task-gutter"><span class="inst-dot"></span></span>
        <span class="inst-task-main"><span class="inst-task-title inst-u-truncate">Validating the paths</span>
          <span class="inst-task-sub">3 nodes unreachable</span></span>
        <span class="inst-task-meta">1.2 s</span>
      </div>
      <div class="inst-task" role="option" data-state="queued" aria-selected="false" tabindex="-1">
        <span class="inst-task-gutter"><span class="inst-dot"></span></span>
        <span class="inst-task-main"><span class="inst-task-title inst-u-truncate">Baking the navmesh</span>
          <span class="inst-task-sub">waiting for the validation</span></span>
        <span class="inst-task-meta">—</span>
      </div>
      <div class="inst-task" role="option" data-state="skipped" aria-selected="false" tabindex="-1">
        <span class="inst-task-gutter"><span class="inst-dot"></span></span>
        <span class="inst-task-main"><span class="inst-task-title inst-u-truncate">Exporting the preview</span>
          <span class="inst-task-sub">skipped by a flag</span></span>
        <span class="inst-task-meta">—</span>
      </div>
    </div>
  </div>
</div>
```

## API

```api
```

A long heading is truncated through `inst-u-truncate` — otherwise it drives the
metadata past the edge.

## Related

```related
```
