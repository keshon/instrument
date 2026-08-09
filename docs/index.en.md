---
title: instrument
template: splash
---

A CSS kit for interfaces that show a human what a machine is doing: task
queues, runs, logs, inspectors, states and uncertainty.

## Who it is for

Dashboards, admin panels and agent tooling — screens where **data** is the
point, where people sit for hours, and where the machine works while the
human watches.

It has what general-purpose kits do not: a task queue, a step with a body, an
approval request, a failure block, a diff, parallel run lanes, a log with
columns, a budget and a tree.

**What it is not for.** Landing pages, shops, content sites. The kit is
deliberately quiet: no gradients, no glows, no shadows on cards, no bold
weights. An all-caps heading is forbidden by a rule, not missing by oversight.

## Three things that are usually not done

**Contrast is checked by machine.** 492 colour pairs against WCAG thresholds
across six themes. The tool reads the real `tokens.css` and resolves `var()`,
`light-dark()` and `color-mix()` the way a browser does, so it cannot drift
from the kit.

**Tap targets are checked by machine too.** 51 targets across all three
densities, including the spacing exception — otherwise compact density would
be a violation by construction.

**Reduced motion does not switch indicators off.** Transitions collapse, but
infinite activity indicators **slow down rather than stop**: a kit whose whole
job is showing that a machine is busy must keep showing it. An infinite
animation squeezed to 0.01ms is not a shorter animation, it is a stopped one.

**Density is a container attribute, not a component variant.** One attribute
on a panel retunes everything inside it, because components carry no numbers
of their own.

## Install

```html
<link rel="stylesheet" href="dist/instrument.min.css">
```

82 KB, one request, 14 KB gzipped. Keyboard behaviour is a separate file and
optional:

```html
<script type="module" src="dist/instrument.js"></script>
```

No build step. `instrument.css` declares the layer order, and that order is
what guarantees application styles always win — without `!important`.

## Requirements

Chrome 123+ · Safari 17.5+ · Firefox 120+ — for `light-dark()`, `oklch()`,
cascade layers, nesting, `:has()` and the `lh` unit. There are no polyfills
and there will not be.

## Translation

The documentation is written in Russian and translated page by page. A page
without its own translation shows the Russian text with a notice rather than
a 404 — otherwise translating seventy-five pages would be all-or-nothing.

**The API reference is the same in both languages**: class, token and
attribute names were English from the start.
