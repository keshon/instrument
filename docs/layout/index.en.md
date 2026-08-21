---
title: Layout
group: Layout
layout: index
---

Eight primitives, and not one of them draws — they only place. The difference
between them is one of **scale**: the shell holds the window, the split and the
container hold the screen, the stack and the cluster hold the blocks inside it.

| | what it holds | how many per document |
|---|---|---|
| [The shell](./shell.md) | the window entire, regions with independent scrolling | one |
| [The rail](./rail.md) | a narrow column of sections as glyphs | one, inside the shell |
| [The status bar](./statusbar.md) | the readings at the bottom of the window | one |
| [The container](./container.md) | the width and the side fields of the content | any number |
| [The split](./split.md) | two columns of differing importance | any number |
| [Flow](./flow.md) | the rhythm: a stack, a cluster, a grid | everywhere |
| [The page header](./page-header.md) | the name of a screen and its actions | one per screen |
| [The section](./section.md) | the heading of a block inside a screen | any number |

The comparison stands here rather than on the pages: "take instead" is a
statement about NEIGHBOURS, and repeated eight times it parts from them by one
line per edit.

## What to choose

| What is wanted | Take | Why not the neighbour |
|---|---|---|
| The framework of a whole application with independently scrolling regions | [The shell](./shell.md) | A split wraps by itself and does not take the whole height of the window |
| Two parts of differing importance: a list and its inspector | [The split](./split.md) | In a grid every column is the same |
| To bound the width of the content and give it side fields | [The container](./container.md) | It sets the width alone; the rhythm is the work of flow |
| A vertical rhythm between blocks | `inst-stack` — [flow](./flow.md) | A container is about width rather than about gaps |
| A row that wraps across the width | `inst-cluster` — [flow](./flow.md) | A split has a threshold of wrapping, a cluster has none |
| Cards in a responsive grid | `inst-grid` — [flow](./flow.md) | A table has rows and column headings |
| The name of a screen and 1–3 actions on it | [The page header](./page-header.md) | A section is quieter and lays no claim to an `<h1>` |
| The heading of a block inside a screen | [The section](./section.md) | A panel has a border, a background and scrolling of its own |
| Sections as glyphs, when room is worth more than a label | [The rail](./rail.md) | Out of ten similar glyphs a rail does not read — there labels are wanted |
| Permanent readings of a window: the branch, the encoding, the count of errors | [The status bar](./statusbar.md) | A toolbar holds controls and is therefore twice as tall |

## When none of them fits

| What is wanted | Take | Why not from here |
|---|---|---|
| A region of the application with a border, a header and scrolling | [a panel](../components/display/panel.md) | The layout primitives do not draw |
| Tools over the content of a panel | [a toolbar](../components/navigation/toolbar.md) | It is part of a panel rather than of the layout of a screen |
| Tabular data | [a table](../components/display/table.md) | A grid has no rows and no column headings |
| A message at the level of a screen | [a banner](../components/feedback/banner.md) | A status bar is read out of the corner of an eye |
| The total of a table | a `tfoot` | It stands under its own data rather than at the bottom of the window |
| An action rather than a transition | [an icon button](../components/actions/button.md) | An item of a rail leads to an address |

## Common to all eight

**A primitive does not draw.** No border, no background, no shadow: as soon as
a layout starts painting, it cannot be nested inside something that paints
otherwise. What draws are the panel and the card, and that is their work.

**A gap is named by intent rather than by a number.** `--gap-section` between
sections, `--gap-row` inside one, `--gap-inline` in a row. The `--space-*`
scale is taken directly only where there is no role and there is nothing to
start one for — and [the build gate](../about/audit.md) guards that.

**A component answers to the width of ITS OWN region rather than of the
window.** So the kit has two named `@container`s — `app` and `content` — and
exactly one media-query point in the whole library, on the shell.
