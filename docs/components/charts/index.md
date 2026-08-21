---
title: Charts
group: Charts
layout: index
---

Six components, and not one of them is a chart in the ordinary sense: they
have no axes. These are **readings** built into the interface — a share, a
composition, a trend, a key to the colours. As soon as a picture is asked for
values it needs axes, and with axes it deserves a place of its own on the screen
rather than a line under a number.

The comparison stands here rather than on the pages: "take instead" is a
statement about NEIGHBOURS, and repeated six times it parts from them by one
line per edit.

## What to choose

| What is wanted | Take | Why not the neighbour |
|---|---|---|
| A share of a known whole across the width of a block | [The meter](./meter.md) | A bar is longer than a ring, so it is more precise |
| What a whole was made **of** — several outcomes at once | [The share bar](./share.md) | A meter compares one quantity against a ceiling; here the parts are compared against each other |
| The same share in a tight place: beside a card title, in a row | [The ring](./ring.md) | A full-width meter does not fit there |
| The direction of a change under a number | [The sparkline](./sparkline.md) | A meter and a ring have no history |
| A key to series told apart by colour | [The legend](./legend.md) | With two or more series it is **an obligation rather than an option** |
| The colour of a series by its number in the data | [The palette](./palette.md) | The status hues are reserved and take no part in the palette |

## When none of them fits

| What is wanted | Take | Why not from here |
|---|---|---|
| Values rather than a shape | a full chart with axes | A sparkline answers "rising or falling", not "how much" |
| Exact values read row by row | [a table](../display/table.md) | Figures are read as figures rather than as a shape |
| A seventh series | [a table](../display/table.md) | Five or six series is the ceiling of telling them apart by colour |
| The state of an object from the vocabulary | [a badge](../display/badge.md) with a dot | A status is not a category, and its colour is fixed |
| **When** something happened rather than how much of it | [a check history](../../agent/history.md) | A share bar cannot order: the same 8% is one long outage or ninety scattered ones |
| Waiting with no quantity at all | [a skeleton](../feedback/skeleton.md) or [a spinner](../feedback/spinner.md) | A meter has to have a ceiling |
| A number by itself | [a metric](../display/metric.md) | A ring with no quantity at the ceiling adds nothing |
| A list of name–value pairs with no chart | [a property list](../display/kv.md) | A series square would report a colour that is not there |
| The spending of an exhaustible resource in a run | [a run budget](../../agent/budget.md) | It knows about the limit and about what happens when it is reached |

## Common to all six

**The quantity lives in an attribute rather than in the width.**
`aria-valuenow` is obligatory: without it the value exists in pixels only and
is available neither to a screen reader nor to a test.

**A categorical colour belongs to a chart only.** The fill of a button, the
highlight of a row, the tone of a badge are taken from the semantics; the
palette of series outside a chart is forbidden, otherwise colour stops meaning
anything.

**Five or six series is the ceiling.** It is not about the palette but about
the eye: a seventh colour is no longer matched to the legend, and the picture
stops answering the question it was drawn for.
