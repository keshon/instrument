---
title: Install
group: Start
source: dist/instrument.min.css
---

The library is added with two tags and needs no build step. Below is a whole
page you can copy and open.

## Install

Three ways. The file is the same in all three.

**npm**

```bash
npm i @keshon/instrument
```

```js
import '@keshon/instrument/min';        // dist/instrument.min.css
import '@keshon/instrument/js';         // behaviour, optional
```

The package also exports the uncompressed CSS (`@keshon/instrument`), the
sprite (`@keshon/instrument/sprite.svg`) and the sources
(`@keshon/instrument/src/*`).

**CDN**

```html
<link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/@keshon/instrument@0.8.0/dist/instrument.min.css">
```

Name the version exactly. `@latest` is cached for a week and may update at an
arbitrary moment.

**As a file**

```html
<link rel="stylesheet" href="instrument.min.css">
```

Take `dist/instrument.min.css` — 108 KB, 17 KB gzipped. That is the whole
library: the tokens, five themes, five scales, three densities, every
component. `dist/instrument.css` is the same thing uncompressed and with the
comments.

## Icons

An icon is drawn by a reference to a symbol: `<use href="#i-refresh">`. The
symbols have to be in the document itself: a reference of the form
`href="https://cdn…/sprite.svg#i-refresh"` is subject to the same-origin rule
and will not render — with no error in the console.

Paste the contents of `assets/sprite.svg` as the first element inside `<body>`
— once per page.

```html
<body>
  <svg class="inst-sprite" aria-hidden="true" focusable="false" xmlns="http://www.w3.org/2000/svg">
    <!-- the contents of assets/sprite.svg go here -->
  </svg>
```

Without the sprite everything else works, the icons simply are not drawn. The
full list is on the [icons](../foundations/icons.md) page, and the count stands
there alone: a number repeated in a second place drifts from the sprite on the
first symbol added.

## Keyboard behaviour

```html
<script type="module" src="instrument.js"></script>
```

A separate file, and optional. It carries out what the library promises by its
markup: arrows in lists and tabs, copying, removing a tag, dragging an axis. It
does not touch the styling. The details are in
[behaviour](../foundations/behavior.md).

## A whole page

```html
<!doctype html>
<html lang="en">
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <title>My application</title>
  <link rel="stylesheet" href="instrument.min.css">
  <script type="module" src="instrument.js"></script>
</head>
<body>
  <div class="inst-shell">
    <header class="inst-shell-header">
      <strong>My application</strong>
      <button class="inst-btn inst-btn--primary" type="button">Run</button>
    </header>

    <main class="inst-shell-main">
      <div class="inst-container">
        <div class="inst-panel">
          <div class="inst-panel-header">
            <span class="inst-panel-title">Queue</span>
          </div>
          <div class="inst-panel-body">
            <div role="listbox" aria-label="Task queue">
              <div class="inst-task" role="option" data-state="running"
                   aria-selected="true" tabindex="0">
                <span class="inst-task-gutter"><span class="inst-dot"></span></span>
                <span class="inst-task-main">
                  <span class="inst-task-title">Parsing the heightmap</span>
                  <span class="inst-task-sub">terrain/heightmap.ts</span>
                </span>
                <span class="inst-task-meta">2.1 s</span>
              </div>
            </div>
          </div>
        </div>
      </div>
    </main>
  </div>
</body>
</html>
```

Here is what it looks like:

```html preview context
<div class="inst-panel">
  <div class="inst-panel-header">
    <span class="inst-panel-title">Queue</span>
  </div>
  <div class="inst-panel-body">
    <div role="listbox" aria-label="Task queue">
      <div class="inst-task" role="option" data-state="running" aria-selected="true" tabindex="0">
        <span class="inst-task-gutter"><span class="inst-dot"></span></span>
        <span class="inst-task-main">
          <span class="inst-task-title">Parsing the heightmap</span>
          <span class="inst-task-sub">terrain/heightmap.ts</span>
        </span>
        <span class="inst-task-meta">2.1 s</span>
      </div>
      <div class="inst-task" role="option" data-state="queued" aria-selected="false" tabindex="-1">
        <span class="inst-task-gutter"><span class="inst-dot"></span></span>
        <span class="inst-task-main">
          <span class="inst-task-title">Building the biomes</span>
        </span>
        <span class="inst-task-meta">—</span>
      </div>
    </div>
  </div>
</div>
```

## Theme and density

Both are attributes, and both work on any subtree rather than only on the
document.

```html
<html data-theme="dark">              <!-- five themes; without the attribute, the OS setting -->
<html data-accent="clay">             <!-- petrol (default) · graphite · indigo · clay -->
<section data-density="compact">      <!-- compact · comfortable -->
```

| Theme | |
|---|---|
| `light-neutral` `light` `light-cool` | light: neutral, warm, cool |
| `dark-soft` `dark` | dark, differing in the depth of the surfaces |

Without the attribute the theme follows the OS setting. What exactly each knob
changes is in [colour](../foundations/colors.md) and
[density](../foundations/density.md).

## Your styles win

The library lies entirely in cascade layers, and your styles lie outside them.
So any rule of yours beats it, and `!important` is needed nowhere.

```css
.inst-btn { border-radius: 0; }   /* works, specificity does not matter */
```

The order of the layers is in [tokens](../foundations/tokens.md#the-order-of-layers).

## Requirements

Chrome 123+, Safari 17.5+, Firefox 128+ — because of `light-dark()`, `oklch()`,
cascade layers, nesting, `:has()`, the `lh` unit and `@property`. There are no
polyfills.

The boundary is held by `@property`: a non-inheriting custom property is how a
column tells its children that they are measured by their content, without
naming them one by one. Firefox has had it since 128, the other two browsers
long before their thresholds above.

## Next

[The run console](../blocks/console.md) — an assembled screen that shows how
the components join up.

[The application shell](../layout/shell.md) — the header, the side panel, the
working area.

[Design principles](../about/design-principles.md) — the rules the components
are derived from.
