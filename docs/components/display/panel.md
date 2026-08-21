---
title: Panel
group: Data display
layout: component
source: src/surfaces.css
api:
  - { name: "inst-panel", kind: "class", doc: "The region. A container on the inline axis" }
  - { name: "inst-panel-header", kind: "class", doc: "The header: a name and actions" }
  - { name: "inst-panel-title", kind: "class", doc: "The name of the region" }
  - { name: "inst-panel-actions", kind: "class", doc: "The actions in the header, pushed to the end" }
  - { name: "inst-panel-body", kind: "class", doc: "The body, with scrolling" }
  - { name: "inst-panel-body--flush", kind: "modifier", doc: "With no inner padding" }
  - { name: "inst-panel-body--list", kind: "modifier", doc: "Padding on the vertical only" }
  - { name: "--surface-raised", kind: "token" }
  - { name: "--border", kind: "token" }
  - { name: "--hairline", kind: "token" }
  - { name: "--radius-lg", kind: "token" }
  - { name: "--pad-panel", kind: "token" }
  - { name: "--text-sm", kind: "token" }
---

A region of an application: it stands on the screen permanently and has a name
of its own. A panel is not a card: it cannot be dragged, because it IS the
place.

```html preview
<div class="inst-panel">
  <div class="inst-panel-header">
    <span class="inst-panel-title">Runs</span>
    <div class="inst-panel-actions">
      <button class="inst-btn inst-btn--sm" type="button">Refresh</button>
    </div>
  </div>
  <div class="inst-panel-body">
    The content of the region.
  </div>
</div>
```

## Contract

| What | Required | Why |
|---|---|---|
| `role="region"` and `aria-labelledby` | yes, if the panel is a landmark | Otherwise it will not appear in a screen reader's list of regions |
| `tabindex="0"` on the body | yes, if there is nothing focusable inside | A scrolling region has to be able to take focus from the keyboard |
| The width comes from outside | yes | A panel is declared a container; `container-type` switches off the intrinsic size and will collapse an element that takes its width from its content |

### Accessibility

| | |
|---|---|
| The heading | `inst-panel-title` is visual. If the panel is a landmark, give it `role="region"` and an `aria-labelledby` pointing at the heading, otherwise it will not appear in the list of regions |
| Scrolling | The body scrolls by itself. A scrolling region has to be able to take focus from the keyboard — browsers give that to a `tabindex="0"` when there is nothing focusable inside |
| Depth | Carried by the surface rather than by a shadow: `--surface-raised` over `--surface-page`. In `forced-colors` the surfaces collapse, so a panel gets an explicit `CanvasText` border |

## Anatomy

### The body of a panel

```html
<div class="inst-panel-body inst-panel-body--list">…</div>
```

| Class | When |
|---|---|
| `inst-panel-body` | The default: inner padding, scrolling of its own |
| `inst-panel-body--flush` | The padding is removed — the content knows about the edges itself (a table, a diff) |
| `inst-panel-body--list` | For a list of rows: padding on the vertical but not on the horizontal |

## Behavior

### A panel is a container

A panel is declared a container (`container: panel / inline-size`), so the
components inside answer to **its** width. One and the same panel stands both
in the narrow column of an inspector and across the whole width of a
dashboard.

```css
/* A rule of your own for a narrow panel */
@container panel (inline-size < 20rem) {
  .my-thing { grid-template-columns: 1fr; }
}
```

The width always comes to a panel from outside, so declaring it a container is
safe. Declaring something a container when it takes its width from its content
is not allowed: `container-type` switches off the intrinsic size on the inline
axis and will collapse the element.

### Something popping up inside a panel

:::warn
**A popover inside a panel will be clipped.** A panel has `overflow: hidden`
for the sake of its radius. The library's contract: what pops up is drawn in
the top layer — `[popover]` or a portal — rather than inside the flow of the
panel. The same goes for [the tooltip](../overlays/tooltip.md): inside a panel
[a popover](../overlays/popover.md) is wanted in its place.
:::

## API

```api
```

## Related

```related
```
