---
title: Overlays
group: Overlays
layout: index
---

Five components, and they are confused constantly. What tells them apart are
two questions: **does it take the focus** and **do people hold it open**.

| | takes the focus | held open | width |
|---|---|---|---|
| [The modal](./dialog.md) | yes, the background is dimmed | no, people answer and close it | by its content |
| [The drawer](./sheet.md) | yes, the background is dimmed | yes, while they look | a column at the side |
| [The popover](./popover.md) | no | no | by its content |
| [The menu](./menu.md) | no | no | by its items |
| [The tooltip](./tooltip.md) | no, and the cursor cannot be moved into it | no | one phrase |

The comparison stands here rather than on the pages: "take instead" is a
statement about NEIGHBOURS, and repeated five times it parts from them by one
line per edit.

## What to choose

| What is wanted | Take | Why not the neighbour |
|---|---|---|
| An answer without which nothing goes on: confirming a destructive action | [The modal](./dialog.md) | A popover does not take the focus, and the question can be missed |
| Details of an object looked at without losing sight of the list | [The drawer](./sheet.md) | A modal centres itself and tears one away from the context |
| A list of actions on an object, called by a button | [The menu](./menu.md) in [a popover](./popover.md) | A select has a value and is sent with the form, a menu has neither |
| Details or a form at the side, without blocking the work | [The popover](./popover.md) | A modal dims the background where standing beside would have done |
| An explanation of an icon button or an abbreviation | [The tooltip](./tooltip.md) | Everything else calls for a press |

## When none of them fits

| What is wanted | Take | Why not from here |
|---|---|---|
| A message with no question | [a banner](../feedback/banner.md) at page level | A modal demands a press where reading would have done |
| A permanent panel nobody closes | a column of [the shell](../../layout/shell.md) | A drawer has a ground under it, and the background is switched off |
| Two or three frequent actions | [a button group](../actions/button-group.md) in plain view | A menu hides what gets pressed every time |
| Choosing a value in a form | [a select](../inputs/select.md) | A menu has no value and is not sent |
| Text that cannot be done without | a label in the flow or [a note](../feedback/note.md) | A tooltip does not open on a touchscreen |
| An approval in an agent queue | [an approval block](../../agent/approval.md) | A modal tears one out of the flow, where the decision is part of the feed |

## Common to all five

**A clipping container decides for you.** Inside [a panel](../display/panel.md)
with `overflow: hidden` a tooltip is clipped — there people take a popover: it
lives in the top layer. That is the one case where the choice is dictated by
the layout rather than by the meaning.

**A tooltip does not open to a finger.** On a touchscreen it is not there at
all, so everything written in it has to be optional.

**A modal is above the toasts.** A toast lives in the top layer, but
`showModal()` raises a dialog higher still: a notification that arrives while a
modal is open will be seen after it closes. A limit of the platform rather than
an oversight.
