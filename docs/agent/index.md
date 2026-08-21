---
title: Agent layer
group: Agent layer
layout: index
---

Twelve components, and all of them are about one thing: **showing a human the
work of a machine that is going on right now.** Hence the differences — not in
the look but in what exactly is shown: the unit of work, its output, its time
or its price.

| | the unit | the form |
|---|---|---|
| [The run](./run.md) | a phase, an attempt, a shard | countable dots |
| [The queue row](./task.md) | a task | a flat row, and there are hundreds |
| [The step](./step.md) | a tool call | it expands, with the output inside |
| [The log](./log.md) | a line from the machine | a stream with a level and a time |
| [The diff](./diff.md) | an edit in a file | two sides |
| [The lanes](./lane.md) | a worker | segments on a shared axis |
| [The history](./history.md) | an attempt with an outcome | a row of equal events |
| [The budget](./budget.md) | a resource with a limit | spending against a ceiling |
| [The tree](./tree.md) | a node | a hierarchy of unknown depth |
| [The approval](./approval.md) | a decision by a human | it stays in the flow |
| [The failure](./failure.md) | a breakdown | a reason and a retry |

The comparison stands here rather than on the pages: "take instead" is a
statement about NEIGHBOURS, and repeated twelve times it parts from them by one
line per edit.

## What to show

| What is wanted | Take | Why not the neighbour |
|---|---|---|
| Tasks as a list, hundreds of them, each with a state | [The queue](./task.md) | The row is flat: a step has a body, and at a hundred rows it does not read |
| A tool call with its input and its output | [The step](./step.md) | A queue has no output, and there is nothing to collapse |
| A uniform stream of lines from a machine | [The log](./log.md) | Steps have a body, a timeline has a varying make-up |
| What the agent changed in a file | [The diff](./diff.md) | A diff with no changes is just [code](../components/display/code.md) |
| Parallel workers on a shared axis | [The lanes](./lane.md) | A queue has no axis of time |
| A row of equal attempts with an outcome | [The history](./history.md) | On a sparkline what matters is the value rather than the outcome |
| Spending against an explicit limit | [The budget](./budget.md) | A metric has no ceiling, a meter has no forecast of exhaustion |
| Countable units of work: phases, shards | [The run](./run.md) | Fifty dots are not counted by eye — there a meter with a number is wanted |
| A hierarchy of unknown depth | [The tree](./tree.md) | A tree of one level is a list with superfluous roles |

## What to do when the work has stopped

| What happened | Take | Why not the neighbour |
|---|---|---|
| The agent is waiting for permission to go on | [The approval](./approval.md) | A modal blocks everything else, and the decision here is part of the feed |
| The agent has failed | [The failure](./failure.md) | An approval has no reason and nothing about what has already been tried |
| The decision has to be taken immediately | [a modal](../components/overlays/dialog.md) | An approval stays in the flow and waits |
| A message with no choice | [a banner](../components/feedback/banner.md) with the `warn` tone | An approval has options, and they are enumerable |

## When none of them fits

| What is wanted | Take | Why not from here |
|---|---|---|
| Named steps known in advance | [wizard steps](../components/navigation/steps.md) | There the sequence is finite, here the make-up comes from the work |
| Events of a varying make-up | [a timeline](../components/display/timeline.md) | A log has uniform lines and columns |
| A continuous quantity | [a meter](../components/charts/meter.md) | A run counts units rather than percentages |
| Comparing records by fields | [a table](../components/display/table.md) | A tree is about structure rather than about comparison |
| The properties of one object | [a property list](../components/display/kv.md) | The nodes of a tree expand, and properties have nothing to expand |

## Common to all twelve

**The state lives in `data-state` rather than in a class.** The vocabulary is
closed: a value outside it gives an error nowhere — it simply does nothing, and
that is why it is guarded by [the registry](../about/audit.md).

**Infinite motion means "running".** A pulsing dot, a turning ring, a
travelling caret — all of it slows at `prefers-reduced-motion` without
stopping: a stopped indicator says "nothing is happening" at exactly the moment
something is.

**The nodes arrive while the work goes on.** Queue rows, steps and tree nodes
appear one at a time, so `instrument.js` works by delegation and watches for
groups appearing with an observer — [more in
behaviour](../foundations/behavior.md).
