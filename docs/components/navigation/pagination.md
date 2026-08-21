---
title: Pagination
group: Navigation
layout: component
source: src/layout.css
api:
  - { name: "inst-pager", kind: "class", doc: "The list of pages. Put on an `<ol>`; it wraps across the width" }
  - { name: "inst-pager-item", kind: "class", doc: "A page or an arrow. A square no smaller than `--control-h-sm`" }
  - { name: "inst-pager-gap", kind: "class", doc: "The ellipsis in the place of a skipped range. Not a link" }
  - { name: "--control-h-sm", kind: "token" }
  - { name: "--radius-sm", kind: "token" }
  - { name: "--space-2", kind: "token" }
  - { name: "--space-3", kind: "token" }
  - { name: "--text-xs", kind: "token" }
  - { name: "--weight-medium", kind: "token" }
  - { name: "--hairline", kind: "token" }
  - { name: "--surface-hover", kind: "token" }
  - { name: "--surface-selected", kind: "token" }
  - { name: "--accent-text", kind: "token" }
  - { name: "--accent-border", kind: "token" }
  - { name: "--text-secondary", kind: "token" }
  - { name: "--text-faint", kind: "token" }
---

Going through the pages of a long list. Every page is an address of its own, so
these are links rather than buttons.

```html preview
<nav aria-label="Pages">
  <ol class="inst-pager">
    <li><a class="inst-pager-item" href="#" aria-disabled="true">Back</a></li>
    <li><a class="inst-pager-item" href="#" aria-current="page">1</a></li>
    <li><a class="inst-pager-item" href="#">2</a></li>
    <li><a class="inst-pager-item" href="#">3</a></li>
    <li><span class="inst-pager-gap">…</span></li>
    <li><a class="inst-pager-item" href="#">24</a></li>
    <li><a class="inst-pager-item" href="#">Forward</a></li>
  </ol>
</nav>
```

## Contract

A pager is an `<ol>` inside a `<nav aria-label>`: the pages are ordered, and
the order is part of the meaning. The library removes the markers and the
indent of the list itself, so no wrappers are needed.

| What | Required | Why |
|---|---|---|
| A `<nav>` + an `aria-label` | yes | A landmark of its own: the screen already has a side navigation and breadcrumbs |
| An `<ol>` with `<li>` | yes | The pages are ordered |
| An `aria-current="page"` | yes | The only source of the current page. There is no class for it |
| An `aria-disabled="true"` | yes, at the edge of the range | A link has no `disabled` — the ARIA attribute is the only truthful one here |

`inst-pager-gap` is a `<span>` rather than a link: an ellipsis marks a gap
rather than a page. It stays inside an `<li>` so that the list does not lose
its structure.

### Accessibility

| | |
|---|---|
| The keyboard | Native links: `Tab` traverses, `Enter` goes. Nothing is intercepted |
| The tap target | At least `--control-h-sm` on both axes — a square rather than a narrow rectangle. At density `compact` it drops to 22px, below the 24px threshold of WCAG 2.5.8: keep a pager out of `compact` |
| The current page | It carries an `aria-current`, a fill **and** a border: colour is not the only carrier |
| Unavailability | An `aria-disabled` is spoken and keeps the focus, unlike removing the link |
| Contrast | The label is `--text-secondary`, the ellipsis `--text-faint` — it is decoration rather than data |
| Print | A pager is not printed: on paper there are no pages left |

## States

```html preview
<nav aria-label="An example of the states">
  <ol class="inst-pager">
    <li><a class="inst-pager-item" href="#" aria-disabled="true">Back</a></li>
    <li><a class="inst-pager-item" href="#">7</a></li>
    <li><a class="inst-pager-item" href="#" aria-current="page">8</a></li>
    <li><a class="inst-pager-item" href="#">9</a></li>
  </ol>
</nav>
```

| State | How it is set | What happens |
|---|---|---|
| current | `aria-current="page"` | `--surface-selected`, `--accent-text`, a border of `--accent-border` by an inner shadow, `--weight-medium` |
| unavailable | `aria-disabled="true"` | An opacity of `0.5` and the mouse removed. A link has no `disabled` — the ARIA attribute is the only truthful one here |
| hover | `:hover` | `--surface-hover` and the full colour of the label |

An `aria-disabled` **does not remove the link from the order of traversal** —
it stays reachable from the keyboard and is spoken as unavailable.
Intercepting the press belongs to the application layer: only the mouse is
removed.

## API

```api
```

## Related

```related
```
