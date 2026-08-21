---
title: Popover
group: Overlays
layout: component
source: src/overlay.css
api:
  - { name: "inst-popover", kind: "class", doc: "The floating surface: the border, the radius, the background, the shadow, the inner padding" }
  - { name: "inst-popover--anchored", kind: "modifier", doc: "Anchoring to the button that called it. Under an `@supports`" }
  - { name: "inst-popover--fill", kind: "modifier", doc: "The full width of the parent instead of `max-content`. For a dropdown list under a field" }
  - { name: "popover", kind: "attribute", value: "auto · manual", doc: "It adds the top layer, closing on `Escape` and the transition of appearing to the surface" }
  - { name: "popovertarget", kind: "attribute", doc: "On the button. The tie and the implicit anchor at once" }
  - { name: "--surface-overlay", kind: "token" }
  - { name: "--shadow-popover", kind: "token" }
  - { name: "--border", kind: "token" }
  - { name: "--hairline", kind: "token" }
  - { name: "--radius-md", kind: "token" }
  - { name: "--space-1", kind: "token" }
  - { name: "--space-2", kind: "token" }
  - { name: "--space-8", kind: "token" }
  - { name: "--dur-2", kind: "token" }
  - { name: "--ease-out", kind: "token" }
---

A small block popping up over the interface on a press of a button and
disappearing on a click outside. Entirely on the platform: `popover` +
`popovertarget`, with not one line of script.

```html preview
<button class="inst-btn" type="button" popovertarget="pop-actions">Actions</button>
<div class="inst-popover inst-popover--anchored" id="pop-actions" popover>
  <div class="inst-menu" role="menu">
    <span class="inst-menu-label">Run #4127</span>
    <button class="inst-menu-item" type="button" role="menuitem">Restart</button>
    <button class="inst-menu-item" type="button" role="menuitem">Copy the id</button>
  </div>
</div>
```

## Contract

| What | Required | Why |
|---|---|---|
| A `popover` on the block | yes, on something popping up | Without it this is a floating surface in the flow: no top layer and no closing on `Escape` |
| An `id` on the block | yes | The button refers to it by name |
| A `popovertarget="id"` on the button | yes | The tie and the implicit anchor at once |
| A `type="button"` on the button | yes | Otherwise inside a form it will submit it |
| An `inst-popover--fill` for a plate under a field | yes | A list narrower or wider than the field reads as somebody else's element beside it rather than as its continuation |

The content of a popover has no class of its own: most often it is
[a menu](./menu.md), but it can be any markup.

:::note
**The look is parted from the mechanism.** `.inst-popover` is a surface; the
behaviour of the top layer is added by the `popover` attribute. So the same
look can be taken for something that floats without being a popover: the
results list of a combobox, an autocomplete hint. With no attribute the block
is simply visible, and the showing is led by the application.
:::

### Accessibility

| | |
|---|---|
| The keyboard | `Enter`/`Space` on the button opens, `Escape` closes, and the focus returns to the button. All of it native, and `instrument.js` intercepts nothing |
| There is no ground | The `::backdrop` is transparent: a popover does not block the page and does not pretend to be modal |
| The role of the content | Roles are not invented. If there is a menu inside, `role="menu"` and `role="menuitem"` are put in the markup |
| A surface with no attribute | An `.inst-popover` with no `popover` has neither a top layer nor closing on `Escape`: the showing, the hiding and the role belong to the application |
| Reduced motion | The transition collapses to `0.01ms` rather than being switched off: the state machines listening for `transitionend` keep working |
| Print | A popover is not printed on paper: it is temporary by definition |

## Variants

```html preview
<button class="inst-btn" type="button" popovertarget="pop-anchored">Under the button</button>
<div class="inst-popover inst-popover--anchored" id="pop-anchored" popover>
  It stands under the button and flips if there is no room below.
</div>
<button class="inst-btn" type="button" popovertarget="pop-centered">At the centre</button>
<div class="inst-popover" id="pop-centered" popover>
  With no modifier a popover stands at the centre of the screen.
</div>
```

| Modifier | Where it stands | How the width is measured |
|---|---|---|
| no modifier | At the centre of the screen | `max-content` within 11–22rem |
| `inst-popover--anchored` | Under the button that called it | The same |
| `inst-popover--fill` | Where it was put | The full width of the parent |

The anchor of `--anchored` is **implicit**: a button with a `popovertarget` is
already the anchor for its popover, so there is no need to start names of
anchors.

The rule is hidden under an `@supports (position-area: block-end)`. Where
`position-area` is not supported a popover stands at the centre of the screen:
it **degrades rather than breaking**. The layout is set in logical values, so
in RTL it mirrors itself. If there is no room below or at the side,
`position-try-fallbacks` flips the popover along the axis wanted rather than
letting it out past the edge of the screen.

`inst-popover--fill` is wanted where a plate has to match the width of what it
stands under: a dropdown list under a field. A list narrower or wider than the
field reads as somebody else's element that happened to end up beside it.

## Sizes

| Property | Value |
|---|---|
| minimum width | `11rem` |
| maximum width | `min(22rem, 100vw - var(--space-8))` — on a narrow screen a popover does not go past the edge |
| width by content | `max-content` within those bounds |
| with `--fill` | `100%` of the parent, with no ceiling |

The shadow is at work here: `--shadow-popover` means "it floats above and is
about to disappear" — that is, it describes a popover entire. The second and
last carrier of a shadow in the library is [the modal](./dialog.md).

## Behavior

A popover has not one line of behaviour in the library, because all of it is in
the browser already.

| What | Who does it |
|---|---|
| The top layer — a popover is cut off by neither an `overflow` nor a neighbour's `z-index` | The platform |
| Closing on `Escape` | The platform |
| Closing on a click outside (light dismiss) | The platform |
| Returning the focus to the button that called it | The platform |
| The tie between the button and the block | `popovertarget="id"` |
| Appearing and disappearing | The library: a transition of `opacity` and `translate` |

The transition plays out on closing too, because the rule declares `display`
and `overlay` with `allow-discrete`. Without them a popover would disappear
instantly, without waiting for its own animation.

## API

```api
```

## Related

```related
```
