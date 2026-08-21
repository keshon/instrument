---
title: Inputs
group: Inputs
layout: index
---

Ten components, and the first question to any of them is one: **is the value
typed or chosen?** Typed means a field; chosen depends on how many options
there are and whether they are all visible at once.

The comparison stands here rather than on the pages: "take instead" is a
statement about NEIGHBOURS, and repeated ten times it parts from them by one
line per edit.

## The value is typed

| What is wanted | Take | Why not the neighbour |
|---|---|---|
| Arbitrary text: a name, a path, a query | [The text field](./input.md) | A select has no free input |
| A number dragged with the mouse, in a dense inspector row | [The number field](./num-field.md) | In an ordinary form the letter of an axis means nothing — there a text field with a suffix does |
| An approximate value from a continuous range: volume, a threshold | [The slider](./slider.md) | Hitting 47 of 100 with a thumb is impossible — if an exact value is wanted, that is a number field |
| Filtering a list, often cleared entirely | [Search](./search.md) | An `input[type=search]` gives the native cross, a text field has none |
| Uploading a file or a drop zone | [The file field](./file.md) | "Export" is [a button](../actions/button.md) rather than a zone |
| Template keys under a field, pasted into the text | [Inserts](./inserts.md) | They have no chosen state: this is pasting rather than choosing |

## The value is chosen

| How many options | Take | Why not the neighbour |
|---|---|---|
| Two states: on or not | [a switch or a checkbox](./toggles.md) | A single option has no pair |
| 2–4 equals, all in view | [a segmented control](../actions/segmented.md) | A select hides the chosen one behind an expansion |
| 2–4, and each has something to explain | [The choice card](./choice-card.md) | A card around the word "Yes" is a superfluous frame |
| More than five mutually exclusive | [The select](./select.md) | Twelve cards are not a choice but a page |
| Several switched on at once and all in view | [a chip](../actions/chip.md) | In a select exactly one is chosen unless it is `multiple` |

## When none of them fits

| What is wanted | Take | Why not from here |
|---|---|---|
| A setting applied at once, with no submitting | [switches](./toggles.md) in [a property list](../display/kv.md) | A form with a button would lie there: it promises "it will apply on the press" |
| A list that has to be searched or added to | there is no combobox in the kit yet | A select takes a ready list and is not added to |
| Actions rather than values | [a menu](../overlays/menu.md) in a popover | A select holds what is chosen, a menu does not |
| Showing a quantity with no input | [a meter](../charts/meter.md) | A slider looks like a control and invites dragging |
| A read-only value | `readonly` rather than `disabled` | `disabled` throws a field out of the traversal and out of the submission |

## Common to all ten

**A form is wanted when there are more than two fields and they apply on
submission.** One field by itself calls for no form — [a form](./form.md)
around it is excessive.

**A label is obligatory on every field.** A `placeholder` is not a label: it
disappears on input, and somebody distracted halfway through loses the name of
what they are filling in.

**The error state lives on `:user-invalid` rather than on a class.** The browser
knows whether a field has been touched; a class does not and paints an empty
field red before the first touch.
