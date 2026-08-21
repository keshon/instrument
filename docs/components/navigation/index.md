---
title: Navigation
group: Navigation
layout: index
---

Six components, and the first question to any of them: **is this going to an
address or choosing a value?** Going is navigation, tabs, breadcrumbs, the
pager. Choosing is not from here: that is [a segmented
control](../actions/segmented.md) or [a select](../inputs/select.md).

The comparison stands here rather than on the pages: "take instead" is a
statement about NEIGHBOURS, and repeated six times it parts from them by one
line per edit.

## What to choose

| What is wanted | Take | Why not the neighbour |
|---|---|---|
| A permanent list of an application's sections | [Navigation](./nav.md) in the column of [the shell](../../layout/shell.md) | Tabs have another job and another place on the screen |
| Several views of the same data, each with an address of its own | [Tabs](./tabs.md) | Navigation leads between sections, tabs within one |
| The path to the current screen deeper than the second level | [Breadcrumbs](./breadcrumbs.md) | Breadcrumbs have no completeness, only a path: a list of sections cannot be replaced by them |
| A list split into pages with addresses `?page=3` | [The pager](./pagination.md) | In an endless feed a page number means nothing |
| A process with a known number of stages, led by a person | [Wizard steps](./steps.md) | Tabs have neither order nor incompleteness |
| Tools acting on the content of **this** panel | [The toolbar](./toolbar.md) | Actions on the whole screen belong in [the page header](../../layout/page-header.md) |

## When none of them fits

| What is wanted | Take | Why not from here |
|---|---|---|
| Choosing one value out of equals | [a segmented control](../actions/segmented.md) | That is input rather than going somewhere: the value is sent with the form |
| An expandable hierarchy of objects | [a tree](../../agent/tree.md) | A navigation item has neither a level nor expansion |
| A list of actions on a click | [a menu](../overlays/menu.md) | Navigation leads to an address, a menu starts an action |
| A stream with no end | loading more on scroll | A page number in an endless feed means nothing |
| A live feed of events | [a log](../../agent/log.md) | People read its tail rather than page 7 |
| The work of a machine with an expandable body | [a run step](../../agent/step.md) | It has output, a wizard step has only a label |
| The share done as a number | [a meter](../charts/meter.md) | A bar with no names of stages |
| The name of a screen | [the page header](../../layout/page-header.md) | Breadcrumbs do not replace a heading |

## Common to all six

**Going somewhere is an `<a>`.** A button leading to an address breaks the
middle click, "open in a new tab" and copying the link. The one exception is
real tabs with `role="tablist"`, which have no address of their own.

**The current position is declared by an attribute rather than by a class.**
`aria-current="page"` on a link, `aria-selected` on a real tab. A class paints
but is not spoken.

**The order of traversal matches the order in the markup.** Rearranging columns
visually without rearranging the nodes breaks `Tab` — that applies to all six
alike.
