---
title: instrument
template: splash
---

A CSS library for interfaces that show a human what a machine is doing: task
queues, runs, logs, inspectors, states and uncertainty.

No build step, no dependencies, one link tag.

```html
<link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/@keshon/instrument@0.6.0/dist/instrument.min.css">
```

Or from npm: `npm i @keshon/instrument`.

## Who it is for

Dashboards, admin panels and agent tooling — screens where **data** is the
point, where people sit for hours, and where the machine works while the human
watches.

It covers what general-purpose kits usually leave out: a task queue, a step
with a body, an approval request, a failure block, a diff, parallel run lanes,
a log with columns, a budget and a tree.

**What it is not for.** Landing pages, shops, content sites. Visual effects are
limited: no gradients, no glows, no shadows on cards, no bold weights. All-caps
headings are forbidden by a rule.

## Features

**Six themes, four accents, three densities** — three independent axes, all set
through markup attributes. Density is declared on a container, not picked as a
component variant: one attribute on a panel retunes everything inside it.

**Contrast is checked by machine.** 2208 token pairs against WCAG thresholds
across every theme and accent. The tool reads the real `tokens.css` and
resolves `var()`, `light-dark()` and `color-mix()` the way a browser does. What
those tokens add up to on screen — after nesting and translucent layers — is
measured by a separate pixel audit you run on your own application.

**Tap targets are checked too.** 54 checks across all three densities,
including the WCAG 2.2 spacing exception.

**Reduced motion does not switch indicators off.** Transitions collapse, but
infinite activity indicators slow down instead of stopping: an animation
squeezed to 0.01ms is a stopped animation, and a stopped indicator reports that
work has halted.

## Install

```html
<link rel="stylesheet" href="dist/instrument.min.css">
```

97 KB, one request, 16 KB gzipped. Keyboard behaviour is a separate file and
optional:

```html
<script type="module" src="dist/instrument.js"></script>
```

Icons need the sprite inlined in the document: `<use href="#i-close">` is
subject to the same-origin rule, so a link to an external sprite file renders
nothing and reports no error. Inline `assets/sprite.svg` (99 symbols) as the
first element of `<body>`.

The kit lives entirely in cascade layers and your styles do not, so any rule of
yours wins without `!important`.

## Requirements

Chrome 123+ · Safari 17.5+ · Firefox 120+ — for `light-dark()`, `oklch()`,
cascade layers, nesting, `:has()` and the `lh` unit. There are no polyfills.

## Translation

The documentation is written in Russian. This page is currently the only
translated one; a page without its own translation shows the Russian text with
a notice rather than a 404.

**The API reference is the same in both languages**: class, token and attribute
names were English from the start.
