---
title: Application shell
group: Layout
layout: component
source: src/layout.css
api:
  - { name: "inst-shell", kind: "class", doc: "A grid the full height of the window (`100dvh`)" }
  - { name: "inst-shell-aside", kind: "class", doc: "The side column" }
  - { name: "inst-shell-brand", kind: "class", doc: "The name of the application above the sections. The first child of the column" }
  - { name: "inst-shell-header", kind: "class", doc: "The header" }
  - { name: "inst-shell-main", kind: "class", doc: "The work area. The `app` container" }
  - { name: "--scroll-pad-end", kind: "token", doc: "The bottom indent of a scrolling area. Declared by the shell, subtracted by a sticky action bar" }
  - { name: "inst-skip-link", kind: "class", doc: "A skip to the content. Visible by `:focus-visible` alone" }
  - { name: "--aside-w", kind: "token" }
  - { name: "--pad-panel", kind: "token" }
  - { name: "--space-3", kind: "token" }
  - { name: "--space-4", kind: "token" }
  - { name: "--space-5", kind: "token" }
  - { name: "--gap-inline", kind: "token" }
  - { name: "--hairline", kind: "token" }
  - { name: "--border", kind: "token" }
  - { name: "--surface-sunken", kind: "token" }
  - { name: "--surface-page", kind: "token" }
  - { name: "--surface-raised", kind: "token" }
  - { name: "--radius-md", kind: "token" }
  - { name: "--text-sm", kind: "token" }
  - { name: "--text-lg", kind: "token" }
  - { name: "--control-h-md", kind: "token" }
  - { name: "--tracking-tight", kind: "token" }
  - { name: "--text-primary", kind: "token" }
  - { name: "--z-modal", kind: "token" }
---

The framework of a screen: a side column, a header, a work area. Every part
scrolls by itself — the side navigation does not travel away while the table is
being read.

```html preview
<div class="inst-shell">
  <aside class="inst-shell-aside">
    <a class="inst-shell-brand" href="#"><span class="inst-dot" data-tone="ok"></span>instrument</a>
    <nav class="inst-nav" aria-label="Sections">
      <div class="inst-nav-group">
        <span class="inst-nav-label">Work</span>
        <a class="inst-nav-item" href="#" aria-current="page">Runs
          <span class="inst-badge inst-nav-count">7</span></a>
        <a class="inst-nav-item" href="#">Queue</a>
      </div>
    </nav>
  </aside>

  <header class="inst-shell-header">
    <nav aria-label="Breadcrumbs">
      <ol class="inst-crumbs">
        <li><a href="#">Projects</a></li>
        <li><span aria-current="page">Run #4127</span></li>
      </ol>
    </nav>
    <span class="inst-cluster-spacer"></span>
    <button class="inst-btn inst-btn--sm" type="button">Export</button>
    <button class="inst-btn inst-btn--sm inst-btn--primary" type="button">Run</button>
  </header>

  <main class="inst-shell-main" id="main">
    <div class="inst-container inst-container--full inst-stack">
      <div class="inst-page-header">
        <div class="inst-page-header-main">
          <h1 class="inst-page-title">Runs</h1>
        </div>
      </div>
    </div>
  </main>
</div>
```

## Contract

```html
<a class="inst-skip-link" href="#main">To the content</a>

<div class="inst-shell">
  <aside class="inst-shell-aside">
    <a class="inst-shell-brand" href="/">The application</a>
    …
  </aside>
  <header class="inst-shell-header">…</header>
  <main class="inst-shell-main" id="main">…</main>
</div>
```

The order in the markup is `aside`, `header`, `main`; on the screen they are
placed by `grid-template-areas`, so the header stands to the right of the
column rather than above it. On a narrow screen the order of the areas is
changed by the same property — with no second set of markup.

The skip to the content stands **before** the shell: it is the first thing the
keyboard meets, and the `href` is obliged to point at the `id` of the work
area.

| What | Why |
|---|---|
| An `inst-skip-link` first in the `<body>` | Otherwise the keyboard walks the whole navigation all over again every time |
| An `id` on `inst-shell-main` | The target of the skip |
| An `<aside>` / `<header>` / `<main>` | The landmarks of the page. A `<div>` gives none |
| An `aria-label` on every `<nav>` inside | There are usually three of them in a shell: the sections, the crumbs, the tabs |

The width of the column is retuned in one line:

```css
.my-app { --aside-w: 280px; }
```

### Accessibility

| | |
|---|---|
| The skip to the content | The first element in the order of traversal. Hidden until focus by `clip-path` rather than by `display: none`: what is hidden that way stays focusable |
| Landmarks | An `<aside>`, a `<header>` and a `<main>` give navigation by landmarks with not a single `role` |
| Scrolling | Three independent areas with `overflow: auto`. `overscroll-behavior` is not set: on a container with nothing to scroll, `contain` eats the mouse wheel |
| Height | `100dvh` rather than `100vh` — on mobile the panel of the browser otherwise cuts off the bottom of the work area |
| Wrapping in the header | Required: without it a narrow screen gets horizontal scrolling of the document, and `100dvh` stops fitting by height |
| Print | The side column is hidden, the shell unfolds into the ordinary flow, the scrolling areas are printed entire |

## Anatomy

| The area | What it does |
|---|---|
| `inst-shell-aside` | The side column, `--aside-w` wide. It scrolls by itself and lies on `--surface-sunken` |
| `inst-shell-header` | The header. A row with required wrapping: on a narrow screen the crumbs and the controls go to a second line rather than pushing the document past the viewport |
| `inst-shell-main` | The work area. It scrolls by itself, gives the vertical air and is declared the **container** `app` |

The horizontal air is given not by the work area but by [the container](./container.md)
inside it: their work differs — one bounds the measure of reading, the other
holds the scrolling.

### The heading of the column

The name of the application is the first child of the column, before the
navigation.

The class lives in the library because of one line — the space below it. In a
**column** it is required: without it the name stands in a row with the items
and becomes the first of them. In a **row** (a narrow screen) that same line
shifts the name up by half of itself, and the logo stops standing on one line
with the items. A divergence of six pixels, which is seen at once and never
found — because it is looked for in the alignment, and it lies in the spacing.

### The sticky action bar

The work area declares its own bottom indent through the variable
`--scroll-pad-end`.

`position: sticky` fixes itself relative to the **padding** box of the
scrolling area rather than to its border. An action bar fixed at
`inset-block-end: 0` stops the height of the indent above the visible bottom —
and a live slit is left under it, through which the content is seen travelling.
The bar stops being the floor.

Only the one who set that indent knows about it. Hence the variable:
[the action bar](../components/inputs/form.md) subtracts it itself. A scrolling
area of your own — declare `--scroll-pad-end` on it.

## Behavior

```css
@media (width < 60rem) { … }
```

**The only threshold in the whole library.** The shell is the one thing that
really does depend on the size of the window: everything else rebuilds by the
intrinsic or answers to the width of its own area through `@container` — see
[the design principles](../about/design-principles.md).

The work area is declared a `@container` region under the name `app`: its width
arrives from the grid of the shell rather than from the contents, so declaring
it a container is safe.

```css
/* A rule of your own under a narrow work area */
@container app (inline-size < 50rem) {
  .my-board { grid-template-columns: 1fr; }
}
```

Below 60rem there are **two behaviours**, and the choice between them is by
**the length of the navigation**.

### The default: a horizontal strip

The column becomes a ribbon **above** the content:

| What changes | How |
|---|---|
| The order of the areas | `header`, `aside`, `main` — as a column |
| The column | A row with horizontal scrolling, a bottom line instead of a side one |
| `inst-nav`, `inst-nav-group` | They unfold into a row |
| `inst-nav-label` | Hidden: the labels of groups do not fit into a strip |
| `inst-nav-item` | It stops wrapping its label |

It works **with not a single line of JS** and therefore cannot "fail to open":
an application that has connected nothing still has navigation on a phone. It
is fit for roughly **up to eight items**.

### The variant: a sliding drawer

```html
<aside class="inst-shell-aside inst-shell-aside--drawer">…</aside>

<button class="inst-btn inst-btn--icon inst-shell-nav-toggle" type="button"
        aria-label="Sections" aria-expanded="false" aria-controls="sidebar">…</button>
```

As soon as there are more than eight items, the ribbon stops working:
seventy-one items is a ribbon almost **eight thousand pixels** wide with no
headings of groups. That is no longer "inconvenient", that is an absence of
navigation.

The drawer keeps the vertical and the labels of the groups, slides out from
behind the start edge and dims the background with the same `--scrim` backing
as a modal.

| The class | The work |
|---|---|
| `inst-shell-aside--drawer` | The column goes into a fixed layer and slides out on demand |
| `inst-shell-nav-toggle` | The button that calls it. On a wide screen it is hidden: the column is in place as it is, and the button would be a deceit |

**What the application gives.** One line — to toggle the `aria-expanded` on the
button. Exactly the same contract as [a modal](../components/overlays/dialog.md)
with `showModal()`. The focus trap, the `Escape` and the click on the backing
are the layer of the application too: the library draws the state but does not
work it out.

## Patterns

A shell by itself is three empty areas. The work is seen when something stands
in them, and the main thing in it is not seen in a picture at all: **every area
scrolls by itself**. Scroll the table — the column and the header will stay in
place.

A screen assembled on this shell lies in "Assemblies":
[the console of runs](../blocks/console.md). There it is also seen how the
navigation, the page header, the split and the states — empty, loading, refusal
— are seated into it.

## API

```api
```

## Related

```related
```
