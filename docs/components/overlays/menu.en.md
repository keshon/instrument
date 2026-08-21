---
title: Menu
group: Overlays
layout: component
source: src/overlay.css
js: The arrows and the roving tabindex are done by `instrument.js`. Opening and closing are taken on by the Popover API
api:
  - { name: "inst-menu", kind: "class", doc: "The container. A column of items with a gap of `--space-1`" }
  - { name: "inst-menu-item", kind: "class", doc: "An item: a `<button>` or an `<a>`" }
  - { name: "inst-menu-item-main", kind: "class", doc: "The wrapper of the name and the explanation. Its presence makes the item two lines — there is no modifier for that" }
  - { name: "inst-menu-item-sub", kind: "class", doc: "The explanation under the name. For a menu of choice, where a name alone says nothing" }
  - { name: "inst-menu-label", kind: "class", doc: "The label of a group. Not interactive" }
  - { name: "inst-menu-sep", kind: "class", doc: "A divider a hairline thick" }
  - { name: "inst-menu-shortcut", kind: "class", doc: "A hotkey at the far edge of an item" }
  - { name: "data-tone", kind: "attribute", doc: "`neutral` `running` `ok` `warn` `error`. Closed. On an item only `error` makes sense" }
  - { name: "aria-checked", kind: "attribute", doc: "`true` — the item is checked" }
  - { name: "aria-current", kind: "attribute", doc: "Any value but `false`" }
  - { name: "aria-disabled", kind: "attribute", doc: "`true` — the item is unavailable" }
  - { name: "--space-1", kind: "token" }
  - { name: "--space-2", kind: "token" }
  - { name: "--space-3", kind: "token" }
  - { name: "--space-6", kind: "token" }
  - { name: "--gap-inline", kind: "token" }
  - { name: "--control-h-sm", kind: "token" }
  - { name: "--radius-sm", kind: "token" }
  - { name: "--text-2xs", kind: "token" }
  - { name: "--text-sm", kind: "token" }
  - { name: "--weight-medium", kind: "token" }
  - { name: "--hairline", kind: "token" }
  - { name: "--border-subtle", kind: "token" }
  - { name: "--surface-hover", kind: "token" }
  - { name: "--accent-text", kind: "token" }
  - { name: "--err-text", kind: "token" }
  - { name: "--err-bg", kind: "token" }
  - { name: "--text-muted", kind: "token" }
  - { name: "--leading-ui", kind: "token" }
---

A list of actions on an object. A menu is the **content** of
[a popover](./popover.md) rather than an overlay in its own right: the top
layer, the closing and the focus come from there.

```html preview
<button class="inst-btn" type="button" popovertarget="menu-run">Actions</button>
<div class="inst-popover inst-popover--anchored" id="menu-run" popover>
  <div class="inst-menu" role="menu">
    <span class="inst-menu-label">Run #4127</span>
    <button class="inst-menu-item" type="button" role="menuitem">
      <svg class="inst-icon" aria-hidden="true"><use href="#i-refresh"/></svg>Restart
      <span class="inst-menu-shortcut"><kbd>R</kbd></span></button>
    <button class="inst-menu-item" type="button" role="menuitem">
      <svg class="inst-icon" aria-hidden="true"><use href="#i-copy"/></svg>Copy the id</button>
    <button class="inst-menu-item" type="button" role="menuitem" aria-checked="true">
      <svg class="inst-icon" aria-hidden="true"><use href="#i-list"/></svg>Show the log</button>
    <span class="inst-menu-sep"></span>
    <button class="inst-menu-item" type="button" role="menuitem" data-tone="error">Delete the run</button>
    <button class="inst-menu-item" type="button" role="menuitem" aria-disabled="true">Archive</button>
  </div>
</div>
```

## Contract

| What | Required | Why |
|---|---|---|
| A `role="menu"` on the container | yes | Without it a `role="menuitem"` is invalid |
| A `role="menuitem"` on the items | yes | The library draws the state but does not invent the role |
| A `type="button"` on items that are buttons | yes | Otherwise inside a form the menu will submit it |
| An `aria-hidden="true"` on the icon | yes | The name of an item is given by the label rather than by the icon |
| The wrapper of [a popover](./popover.md) | yes | A menu does not pop up by itself: the top layer and the closing come from the `popover` |

:::warn
Moving by arrows and the roving `tabindex` are done by
[`instrument.js`](../../foundations/behavior.md). Opening, closing on `Escape`
and on a click outside, and returning the focus to the button call for no
script at all — they are taken on by the Popover API.
:::

### Accessibility

| | |
|---|---|
| The keyboard | `Tab` to the button, `Enter` to open, `Escape` to close. **The arrows inside the menu are set by the application**: without them a menu is traversed by tab, which for a long list is agony |
| The roles | A `role="menu"` + a `role="menuitem"` are required. An `aria-checked` with no `role="menuitem"` is not spoken |
| An unavailable item | `aria-disabled="true"` rather than `disabled`: the item stays in the arrow traversal and gets a roving `tabindex` but is not chosen — neither a click nor `Enter` calls the application ([behaviour](../../foundations/behavior.md)). A `disabled` removes the item from the traversal by the platform's own means, and the action is never learned of |
| Colour is not the only carrier | A destructive item is red **and** named by the word "Delete". A checked one carries an `aria-checked` rather than a colour alone |
| The tap target | An item of height `--control-h-sm` across the whole width of the popover: it is easier to hit with a mouse than the text |
| Wrapping | `white-space: nowrap` — an item does not wrap. Shorten a long label rather than hoping for a wrap |

## Anatomy

The label of a group, the divider and the hotkey are three different jobs, and
confusing them is expensive.

```html
<div class="inst-menu" role="menu">
  <span class="inst-menu-label">Run #4127</span>
  <button class="inst-menu-item" type="button" role="menuitem">Restart
    <span class="inst-menu-shortcut"><kbd>R</kbd></span></button>
  <span class="inst-menu-sep"></span>
  <button class="inst-menu-item" type="button" role="menuitem">Export</button>
</div>
```

| Class | The work |
|---|---|
| `inst-menu-label` | The label of a group. Not interactive and not in the traversal |
| `inst-menu-item` | An item. A `<button>` for an action, an `<a>` for a transition |
| `inst-menu-sep` | The divider between groups. A `<span>` rather than an `<hr>`: it is styling, and there is nothing to speak |
| `inst-menu-shortcut` | A hotkey at the far edge |

A hotkey is pushed to the far edge and is quieter than the label: it is a
**hint rather than a second name of the item**.

An item is a `<button>` or an `<a>`, by the work. An action is a button, a
transition to an address a link: a button leading to an address breaks the
middle click and "open in a new tab".

## States

```html preview
<div class="inst-popover" popover id="menu-kinds">
  <div class="inst-menu" role="menu">
    <span class="inst-menu-label">Kind</span>
    <button class="inst-menu-item" type="button" role="menuitem">Ordinary</button>
    <button class="inst-menu-item" type="button" role="menuitem" aria-checked="true">Checked</button>
    <a class="inst-menu-item" href="#menu" role="menuitem" aria-current="page">The current address</a>
    <span class="inst-menu-sep"></span>
    <button class="inst-menu-item" type="button" role="menuitem" data-tone="error">Destructive</button>
    <button class="inst-menu-item" type="button" role="menuitem" aria-disabled="true">Unavailable</button>
  </div>
</div>
<button class="inst-btn" type="button" popovertarget="menu-kinds">Kinds of item</button>
```

An item with an explanation. The `.inst-menu-item-main` wrapper makes the item
two lines by itself — there is no modifier for that, because it is not a choice
of the author of the markup but a consequence of what is inside. A menu of
choice with no explanation turns into a list of code words.

```html preview
<div class="inst-popover" popover id="menu-model">
  <div class="inst-menu" role="menu">
    <button class="inst-menu-item" type="button" role="menuitem">
      <span class="inst-menu-item-main">Fast
        <span class="inst-menu-item-sub">For short runs</span></span></button>
    <button class="inst-menu-item" type="button" role="menuitem" aria-checked="true">
      <span class="inst-menu-item-main">Full
        <span class="inst-menu-item-sub">For complex tasks</span></span></button>
    <button class="inst-menu-item" type="button" role="menuitem">
      <span class="inst-menu-item-main">Thrifty
        <span class="inst-menu-item-sub">Cheapest of all on everyday work</span></span></button>
  </div>
</div>
<button class="inst-btn" type="button" popovertarget="menu-model">Run mode</button>
```

| The state of an item | How it is set | What happens |
|---|---|---|
| ordinary | — | `--text-primary` |
| hover | `:hover` | `--surface-hover` |
| checked | `aria-checked="true"` | `--accent-text` and a medium weight |
| current | `aria-current` (any value but `false`) | The same styling as a checked one |
| destructive | `data-tone="error"` | `--err-text`, and `--err-bg` on hover |
| unavailable | `aria-disabled="true"` or `disabled` | An opacity of `0.5`, the mouse removed |

## JS

Include the module once per page — there is no need to initialise the
components one by one, `instrument.js` works by delegation and sees nodes that
arrived later.

```html
<script type="module" src="instrument.js"></script>
```

### What `instrument.js` does

The example in the header is live: open the menu and go through the items with
the arrows.

Arrows between the items, `Home`, `End`, a traversal that wraps and a roving
`tabindex`. An item has no selection and will have none: a menu item is an
action rather than a choice, and an `aria-selected` on it would lie to
assistive technology.

Opening, closing, returning the focus and closing on `Escape` are taken on by
the Popover API — they call for no script.

### Events

There are none of its own. An item is a `<button>`, and its `click` works the
same way as everywhere.

```js
menu.addEventListener('click', (e) => {
  const item = e.target.closest('[role="menuitem"]');
  if (!item) return;
  menu.hidePopover();
  run(item.dataset.action);
});
```

The menu is closed **before** the action: the action may open a modal, and the
menu would be left hanging over it.

## API

```api
```

## Related

```related
```
