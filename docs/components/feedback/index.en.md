---
title: Feedback
group: Feedback
layout: index
---

Eight components, and they all answer one question: **what is happening now and
what to do about it.** Two things tell them apart: how long it holds and what
it is about.

| | holds | is about |
|---|---|---|
| [The toast](./toast.md) | seconds, it leaves by itself | what the person has just done |
| [The banner](./banner.md) | as long as the condition holds | the whole screen |
| [The note](./note.md) | always | the paragraph, field or block beside it |
| [The spinner](./spinner.md) | while the work goes on | the place it stands in |
| [The skeleton](./skeleton.md) | while it loads | the shape of the content to come |
| [The empty state](./empty.md) | while it is empty | the whole list |
| [State dots](./states.md) | as long as the phase holds | one object |
| [The accordion](./accordion.md) | — | it is a layout rather than a message |

The comparison stands here rather than on the pages: "take instead" is a
statement about NEIGHBOURS, and repeated eight times it parts from them by one
line per edit.

## Waiting: what to show

| What is known | Take | Why not the neighbour |
|---|---|---|
| The shape of the content to come: a heading and three lines | [The skeleton](./skeleton.md) | A skeleton of the wrong shape lies twice |
| Nothing, and the wait is shorter than a second | [The spinner](./spinner.md) | A skeleton would promise a shape that will not be |
| The share done | [a meter](../charts/meter.md) | Neither a spinner nor a skeleton says how much is left |
| A button is busy | `aria-busy="true"` on [the button](../actions/button.md) itself | It has a ring of its own in the place of the label |
| There is no data and there will be none | [The empty state](./empty.md) | A skeleton at that moment promises what will not happen |
| Empty because of an error | [a failure block](../../agent/failure.md) | There a reason and what has already been tried are wanted |

## A message: where to show it

| What is wanted | Take | Why not the neighbour |
|---|---|---|
| The result of an action a person started themselves | [The toast](./toast.md) | A banner holds as long as the condition holds and does not leave by itself |
| A condition acting on the whole screen | [The banner](./banner.md) | A toast will leave before the condition does |
| An explanation of the block beside it | [The note](./note.md) | A banner lays claim to the whole screen |
| A question without an answer to which nothing goes on | [a modal](../overlays/dialog.md) | A toast and a banner have no question |
| A refusal that needs a way out | [a failure block](../../agent/failure.md) | It has a reason and a retry, and it stays in the flow |
| An error in one field | a message of [the form](../inputs/form.md) under the field | A screen-level banner about one field is noise |

## When none of them fits

| What is wanted | Take | Why not from here |
|---|---|---|
| Equal views of one object | [tabs](../navigation/tabs.md) | All their names are visible at once and there is no jump in height |
| Everything that has to be seen at once | [panels](../display/panel.md) in a column | An accordion hides what is being compared |
| A tool call with a body of output | [an agent step](../../agent/step.md) | It has a phase and a duration |
| A status with no word | [a badge](../display/badge.md) entire | One dot does not name a state |

## Common to all eight

**An error interrupts, the rest does not.** An error message has `role="alert"`,
a confirmation has `role="status"`: the first interrupts reading, the second
waits for a pause. Confusing them means either interrupting over trifles or
staying silent about a failure.

**A live region has to be in the layout.** The toast region is opened by a
popover once and is not closed: a closed `popover` is `display: none`, and a
live region in `display: none` is not spoken at all.

**An infinite indicator is not stopped at `prefers-reduced-motion`.** An
infinite animation squeezed to 0.01ms is a stop, and a stopped spinner says
"nothing is happening" at exactly the moment something is. Such things are
slowed to a frequency that reads as "running".
