---
title: Menu bar
group: Layout
layout: component
source: src/layout.css
api:
  - { name: "inst-menubar", kind: "class", doc: "The strip at the top edge of the window. One line in height. A row of triggers and nothing else" }
  - { name: "inst-menubar-item", kind: "class", doc: "One trigger. A real `<button>`: it opens a menu and is reached by keyboard" }
  - { name: "aria-expanded", kind: "attribute", value: "true · false", doc: "Marks the open trigger. Needed only when the menu is not the trigger's next sibling — otherwise the strip reads the popover itself" }
  - { name: "--control-h-sm", kind: "token", doc: "The height of a trigger — the same rung as a menu item. The gap between neighbours makes up the rest of the target" }
  - { name: "--text-xs", kind: "token", doc: "The type size of the strip" }
  - { name: "--surface-sunken", kind: "token" }
  - { name: "--surface-selected", kind: "token", doc: "The ground of the open trigger" }
  - { name: "--border-subtle", kind: "token" }
---

The strip of a desktop application: File, Edit, View, Help. It is a row of
**triggers** and nothing else — what drops from them is the kit's
[popover](../components/overlays/popover.md) and [menu](../components/overlays/menu.md),
unchanged. That is why the component is this small: what was missing was never
the menus.

A popover is anchored to the button that opened it with no anchor names to
declare, so a menu lands under its own trigger by itself.

Why a kit for dashboards has one: a screen that is a window rather than a page
is read as a window largely through this strip, and a person looks for commands
at the top left whatever the kit thinks. An application built on the kit put its
commands behind a single icon button in the top right; the owner of the
application read the screen and said it felt unfamiliar. The second attempt put
the strip back. The strip is also where commands go that have nowhere else to
live — select all, invert the selection, quit — and without it they end up as a
drawer of odds and ends, or they do not exist.

```html preview
<div class="inst-menubar" role="menubar" aria-label="Main menu">
  <button class="inst-menubar-item" type="button" role="menuitem" popovertarget="mb-file">File</button>
  <div class="inst-popover inst-popover--anchored" id="mb-file" popover>
    <div class="inst-menu" role="menu">
      <button class="inst-menu-item" type="button" role="menuitem">
        Open a folder… <span class="inst-menu-shortcut"><kbd>Ctrl</kbd><kbd>O</kbd></span></button>
      <button class="inst-menu-item" type="button" role="menuitem">
        Reread <span class="inst-menu-shortcut"><kbd>F5</kbd></span></button>
      <span class="inst-menu-sep"></span>
      <button class="inst-menu-item" type="button" role="menuitem" data-tone="error">Quit</button>
    </div>
  </div>

  <button class="inst-menubar-item" type="button" role="menuitem" popovertarget="mb-view">View</button>
  <div class="inst-popover inst-popover--anchored" id="mb-view" popover>
    <div class="inst-menu" role="menu">
      <span class="inst-menu-label">Theme</span>
      <button class="inst-menu-item" type="button" role="menuitem" aria-checked="true">Dark</button>
      <button class="inst-menu-item" type="button" role="menuitem">Light</button>
    </div>
  </div>

  <button class="inst-menubar-item" type="button" role="menuitem" popovertarget="mb-help">Help</button>
  <div class="inst-popover inst-popover--anchored" id="mb-help" popover>
    <div class="inst-menu" role="menu">
      <button class="inst-menu-item" type="button" role="menuitem">About</button>
    </div>
  </div>
</div>
```

## Contract

| What | Required | Why |
|---|---|---|
| A real `<button>` for a trigger | yes | It opens a menu and is reached by keyboard. A row of `<div>`s says neither |
| The menu stands immediately after its trigger | recommended | Then the strip marks the open one by itself. Otherwise the markup has to set `aria-expanded` |
| Names, not verbs with objects | yes | "File", not "File operations". The strip is an index, and an index that explains itself stops being scannable |
| One strip per window | yes | The same reasoning as the status bar: it belongs to the window, and there is one |
| The strip goes in the [shell](./shell.md) | recommended | The shell keeps a row for it. Outside the shell the application has to own the window height itself |

### Accessibility

| | |
|---|---|
| The role | `role="menubar"` on the strip, `role="menuitem"` on a trigger, and an `aria-label` naming the strip. Without the label a screen reader announces "menu bar" and nothing else |
| The tap target | `--control-h-sm` plus the gap between neighbours. The rung below is the status bar's, and it fails the 24px minimum at the smallest scale — readings are read, triggers are pressed. Checked by `cmd/targets` |
| The open state | Carried by the ground of the trigger. A popovertarget button holds its expanded state in the accessibility tree, so a reader is told; the ground is what tells everyone else |
| Contrast | `--text-secondary` on `--surface-sunken`, and `--text-primary` on `--surface-selected` when open |
| The order of traversal | It matches the order in the markup. The arrow keys between triggers are the application's to add: the kit does not ship behaviour here |

## API

```api
```

## Related

```related
```
