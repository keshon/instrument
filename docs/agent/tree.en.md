---
title: Tree
group: Agent layer
layout: component
source: src/text.css
js: Lazy loading of children belongs to the application layer. The arrows, the expansion and the roving tabindex are done by `instrument.js`
api:
  - { name: "inst-tree", kind: "class", doc: "The container" }
  - { name: "inst-tree--guides", kind: "modifier", doc: "Vertical guides of nesting. Wanted from four levels; at two they would be noise" }
  - { name: "inst-tree-item", kind: "class", doc: "A node" }
  - { name: "inst-tree-twist", kind: "class", doc: "The triangle of expansion. It turns by `aria-expanded`; a click on it expands the node and does not carry the selection" }
  - { name: "aria-expanded", kind: "attribute", doc: "`true` · `false`. On nodes with children alone" }
  - { name: "aria-selected", kind: "attribute", doc: "`true` · `false`" }
  - { name: "aria-level", kind: "attribute", doc: "an integer, from 1" }
  - { name: "--depth", kind: "variable", value: "0" }
  - { name: "--size-indent", kind: "token" }
  - { name: "--control-h-sm", kind: "token" }
  - { name: "--space-3", kind: "token" }
  - { name: "--radius-sm", kind: "token" }
  - { name: "--text-sm", kind: "token" }
  - { name: "--surface-hover", kind: "token" }
  - { name: "--surface-selected", kind: "token" }
  - { name: "--size-chevron", kind: "token" }
---

A hierarchy: the files of a project, the structure of a scene, the nesting of
calls.

The example has four levels and an `inst-tree--guides` — vertical guides of
nesting. From four levels the eye loses which branch a row belongs to; at two
the same lines would be noise, and that is what the modifier is separate for.

```html preview
<div class="inst-tree inst-tree--guides" role="tree" aria-label="The files of the project">
  <div class="inst-tree-item" role="treeitem" aria-level="1" aria-expanded="true"
       tabindex="0" style="--depth:0">
    <span class="inst-tree-twist"></span>src
  </div>
  <div class="inst-tree-item" role="treeitem" aria-level="2" aria-expanded="true"
       tabindex="-1" style="--depth:1">
    <span class="inst-tree-twist"></span>world
  </div>
  <div class="inst-tree-item" role="treeitem" aria-level="3" aria-expanded="true"
       tabindex="-1" style="--depth:2">
    <span class="inst-tree-twist"></span>terrain
  </div>
  <div class="inst-tree-item" role="treeitem" aria-level="4"
       tabindex="-1" style="--depth:3">heightmap.ts</div>
  <div class="inst-tree-item" role="treeitem" aria-level="4" aria-selected="true"
       tabindex="-1" style="--depth:3">chunks.bin</div>
</div>
```

## Contract

The state comes from the markup: the library draws it but does not compute it.

| What | Why |
|---|---|
| A `role="tree"` on the container | Otherwise this is a set of `<div>`s |
| A `role="treeitem"` on a node | The carrier of the role |
| An `aria-level` | The only way to speak the depth: a visual indent is invisible to a screen reader |
| An `aria-expanded` | **On nodes with children alone.** On a leaf it lies that there is something there |
| An `aria-selected` | If selection is supported |
| A roving `tabindex` | `0` on the active one, `-1` on the rest |

### Accessibility

| | |
|---|---|
| The keyboard | `↑` `↓` over the visible nodes, `→` expands or enters, `←` collapses or goes to the parent, `Home`/`End`. The behaviour belongs entirely to the application |
| Depth | The `aria-level` alone. The indent is a visual mark; for a screen reader it is not there |
| An `aria-expanded` on a leaf | A frequent mistake: the attribute with a value of `false` on a node with no children announces a collapsed subtree that does not exist |
| The triangle | Decorative. The state of expansion is carried by the `aria-expanded` rather than by the turn of the glyph |
| The tap target | The height is `--control-h-sm`, that is 26px. At density `compact` it is 22px, **below the 24px of WCAG 2.5.8** |
| Virtualisation | With virtualisation an `aria-setsize` and an `aria-posinset` are required: otherwise a screen reader announces "3 of 3" when there are three thousand nodes |

## Anatomy

```html
<div class="inst-tree-item" style="--depth:2">…</div>
```

The indent is computed from `--depth` by multiplication by `--size-indent`, so
it reacts to the density of the container. The depth is set inline as a
**value** — the same allowed channel as the width of a meter's fill and the
bounds of a lane segment.

A flat list of nodes with a `--depth` rather than nested containers: this way a
row can be virtualised, and a roving `tabindex` need not walk the tree
recursively.

## JS

Include the module once per page — there is no need to initialise the
components one by one, `instrument.js` works by delegation and sees nodes that
arrived later.

```html
<script type="module" src="instrument.js"></script>
```

### What `instrument.js` does

The example in the header is live: `↓` and `↑` over the nodes, `→` expands,
`←` collapses and goes to the parent.

The contract of the `treeitem` role entire: the up and down arrows over the
nodes, `→` expands a collapsed one and goes inside an expanded one, `←`
collapses and goes to the parent. The parent is found by `aria-level`, so the
nesting in the markup may be flat — for a tree of ten thousand nodes that is
the only workable option.

Expanding changes `aria-expanded`; how it looks is decided by CSS.

### Events

`inst:select` — on the chosen node, with `{ value }` in the `detail`.

```js
tree.addEventListener('inst:select', (e) => open(e.detail.value));
```

Expanding: `←` and `→` from the keyboard, a click on the twist with the mouse.
A click on the row itself selects the node and does not touch the expansion —
the twist and the row are two different targets, and their being neighbours is
no reason to glue their meanings together.

```js
tree.addEventListener('inst:expand', (e) => {
  if (!e.detail.open) return;
  e.preventDefault();          // the application will remove and show the rows
  loadChildren(e.target);
});
```

### What is left to the application

Lazy loading of children. `instrument.js` toggled the `aria-expanded` — the
application decides whether to go for the data. The `inst:expand` event is
cancellable: an application that draws a tree from data will remove and show
the rows itself.

## API

```api
```

## Related

```related
```
