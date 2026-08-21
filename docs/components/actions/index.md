---
title: Actions
group: Actions
layout: index
---

Four components, and all four are about pressing. What sets them apart is not
the look but what happens after: a button **does**, a segmented control
**picks one**, a chip **switches several on**, a group **gathers neighbours**.

The comparison stands here rather than on the pages: "take instead" is a
statement about NEIGHBOURS, and repeated on each of the four pages it parts
from them by one line per edit.

## What to choose

| What is wanted | Take | Why not the neighbour |
|---|---|---|
| An action that changes state: run, save, delete | [The button](./button.md) | A segmented control has no action, and neither has a chip |
| Exactly one value out of 2–4 equals: density, view mode | [The segmented control](./segmented.md) | Three buttons do not show which is chosen: a button has no "pressed" state |
| Several filters switched on at once and all in view | [The chip](./chip.md) | In a segmented control the chosen one is single by construction |
| Different actions on one object standing flush: copy · duplicate · delete | [The button group](./button-group.md) | A group has no selected element, and cannot have one |

## When none of them fits

| What is wanted | Take | Why not from here |
|---|---|---|
| Going to an address | an ordinary `<a>` | A button leading to another page breaks the middle click, "open in a new tab" and copying the link |
| Toggling on and off | [a switch or a checkbox](../inputs/toggles.md) | A single option has no pair, and a button has no "pressed" state |
| More than five options | [a select](../inputs/select.md) | A segmented control sprawls across the width, and the chips wrap onto three lines |
| A list of actions that does not fit | [a menu](../overlays/menu.md) in a popover | A button group does not fold away |
| Moving between sections | [tabs](../navigation/tabs.md) | That is routing rather than choosing a value |
| Showing what is already chosen, with no toggling | [a tag](../display/tag.md) | A chip is pressed, a tag is not |
| The state of an object from the vocabulary | [a badge](../display/badge.md) | A chip names a filter rather than a state |
| Unrelated actions | [a cluster](../../layout/flow.md) with a gap | Standing flush they would lie about kinship |

## Common to all four

**Exactly one `primary` per screen.** There is no such thing as two principal
actions: the second takes weight from the first, and a person stops
understanding what is wanted of them.

**A tone carries meaning rather than a mood.** `--danger` is put where an action
is irreversible, and nowhere else. A chip has no tone at all: it names a
filter, and a filter painted the red of "it failed" would say the filter itself
is the emergency.

**The keyboard differs across the four, and that is not a muddle.** A button and
a group take one `Tab` each: they are separate actions. A segmented control and
a chip take one `Tab` for the whole strip with arrows inside: they are one
control with several values. The markup that declares this is named by each
page in its own "Contract" section.
