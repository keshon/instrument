---
title: Data display
group: Data display
layout: index
---

Twelve components, and the choice between them almost always comes down to one
question: **how many objects, and are they compared with one another.** One
object is a property list; many of one kind, compared, is a table; many of
different kinds is cards.

The comparison stands here rather than on the pages: "take instead" is a
statement about NEIGHBOURS, and repeated twelve times it parts from them by one
line per edit.

## How to show the data

| What is wanted | Take | Why not the neighbour |
|---|---|---|
| The properties of **one** object, 3–15 pairs | [The property list](./kv.md) | In a table the columns compare, and there is nothing to compare with |
| Records compared by the same fields | [The table](./table.md) | Half the cells of cards would be empty |
| Objects of a list with a heading and a make-up of their own | [Cards](./card.md) | In a table things of different kinds do not line up into columns |
| A hero number people watch | [The metric](./metric.md) | A table is wanted when there are many numbers and a breakdown |
| Events of different kinds in time | [The timeline](./timeline.md) | A log has columns and monospaced setting |
| An overview of a month laid out by weeks | [The calendar](./calendar.md) | Entering one date is an `<input type="date">` |
| A command, a fragment of config, markup | [Code](./code.md) | A log has a level and a time, a diff has sides |

## How to mark an object

| What is wanted | Take | Why not the neighbour |
|---|---|---|
| A state from an enumerable set: running, done, failed | [The badge](./badge.md) | A badge is the library's vocabulary, a tag is the user's words |
| A label started by the user and removable | [The tag](./tag.md) | A badge is not removed |
| The fate of a file in a list of changes | [The change mark](./change.md) | "Deleted" is not a judgement, and a badge with a tone implies one |
| Recognising a person or an agent | [The avatar](./avatar.md) | An avatar answers "who" rather than "how are things" |
| A region of the application with a header and scrolling | [The panel](./panel.md) | There are a few panels and they do not change; there are any number of cards and they come from data |

## When none of them fits

| What is wanted | Take | Why not from here |
|---|---|---|
| Choosing a filter | [a chip](../actions/chip.md) | A tag shows what is already chosen but does not toggle |
| Choosing an option | [a choice card](../inputs/choice-card.md) | A card has no "chosen" state |
| Rows with a state and progress | [a queue](../../agent/task.md) | It is flat: at a hundred rows the borders turn into a grid |
| A stream of rows from a machine | [a log](../../agent/log.md) | A level, a time and a tail people read instead of a page |
| Changes in a file | [a diff](../../agent/diff.md) | It has two sides and a count of lines |
| Progress towards a limit | [a meter](../charts/meter.md) | A metric has no maximum |
| A series of values in time | [a sparkline](../charts/sparkline.md) | A metric shows one number |
| Editable properties | the inspector (`inst-props`) | A property list is read-only |
| Tens of thousands of rows | virtualisation | The library has none, and a table will lie down under them |

## Common to all twelve

**A tone carries meaning rather than a mood.** A badge, a change mark and a
state dot take their tone from the vocabulary; "deleted" is not "bad", and
painting it as an error lies about the meaning.

**A border declares an object.** Four metrics without borders are one group;
four metrics in borders are four objects. So a row of numbers gets no borders
and a card does.

**The figures are tabular everywhere.** `font-variant-numeric: tabular-nums`
stands in the base: in a tooling interface almost any number either stands in a
column or updates in place, and proportional figures tug at both cases.
