---
title: Tooltip
group: Overlays
layout: component
source: src/overlay.css
api:
  - { name: "inst-tooltip", kind: "class", doc: "The wrapper of the trigger. `position: relative`, the position is counted from it" }
  - { name: "inst-tooltip-text", kind: "class", doc: "The label itself. Above the trigger, centred" }
  - { name: "inst-tooltip-text--wide", kind: "modifier", doc: "Line wrapping and a width ceiling of `22rem`" }
  - { name: "--z-popover", kind: "token" }
  - { name: "--surface-overlay", kind: "token" }
  - { name: "--shadow-popover", kind: "token" }
  - { name: "--border", kind: "token" }
  - { name: "--hairline", kind: "token" }
  - { name: "--radius-sm", kind: "token" }
  - { name: "--space-2", kind: "token" }
  - { name: "--space-3", kind: "token" }
  - { name: "--text-xs", kind: "token" }
  - { name: "--leading-ui", kind: "token" }
  - { name: "--dur-1", kind: "token" }
  - { name: "--ease-out", kind: "token" }
---

A short label appearing on hover and on focus. The one overlay of the library
made **without** the Popover API: by ordinary positioning, so that it works the
same everywhere.

```html preview
<span class="inst-tooltip">
  <button class="inst-btn inst-btn--icon inst-btn--ghost" type="button" aria-label="Refresh" aria-describedby="tt-refresh">
    <svg class="inst-icon" aria-hidden="true"><use href="#i-refresh"/></svg>
  </button>
  <span class="inst-tooltip-text" role="tooltip" id="tt-refresh">Refresh</span>
</span>
<span class="inst-tooltip">
  <button class="inst-btn inst-btn--ghost" type="button" aria-describedby="tt-pass">What is a pass?</button>
  <span class="inst-tooltip-text" role="tooltip" id="tt-pass">One full sweep of the map by an agent</span>
</span>
```

The label is centred on the trigger and **does not flip**: ordinary
positioning has no spare places. At the edge of a narrow region the nearest
`overflow: hidden` clips it — there [a popover](./popover.md) is wanted.

## Contract

| What | Required | Why |
|---|---|---|
| A `role="tooltip"` on the label | yes | Otherwise it is just hidden text |
| An `id` on the label | yes | The trigger refers to it |
| An `aria-describedby` on the trigger | yes | It points at the `id` of the label and is the only thing tying it to the control. Without it a tooltip does not exist for a screen reader |
| The trigger is **focusable** | yes | The showing goes through `:focus-within`. A non-focusable element wants a `tabindex="0"` |

### Accessibility

| | |
|---|---|
| The keyboard | A tooltip appears on `:focus-within`, that is, on tabbing to the trigger. Without that it would be available to the mouse alone |
| A touchscreen | There is no hover there, and the focus arrives together with the press. **Do not put into a tooltip anything that cannot be done without** |
| It does not replace a name | An `aria-describedby` is a description. The accessible name of an icon button is given by an `aria-label` rather than by a tooltip |
| Overlapping | A tooltip stands `--space-3` above the trigger: it does not cover what it explains |
| Reduced motion | The transition of opacity collapses along with all the others |
| Print | It is not printed on paper |

## Variants

```html preview
<span class="inst-tooltip">
  <button class="inst-btn inst-btn--sm" type="button" aria-describedby="tt-short">Short</button>
  <span class="inst-tooltip-text" role="tooltip" id="tt-short">19:38:04</span>
</span>
<span class="inst-tooltip" style="margin-inline-start:9rem">
  <button class="inst-btn inst-btn--sm" type="button" aria-describedby="tt-wide">In detail</button>
  <span class="inst-tooltip-text inst-tooltip-text--wide" role="tooltip" id="tt-wide">Opening a modal is the one line of script in the whole overlays section; closing, the ground and blocking the scroll do without it.</span>
</span>
```

A base tooltip is one line (`white-space: nowrap`). That is right for a label
of two words and breaks on a sentence: the line stretches past the edge of the
screen. So a long text wants an explicit `inst-tooltip-text--wide` — wrapping
and a width ceiling of `22rem`.

## Behavior

| The event | The selector |
|---|---|
| hover by mouse | `.inst-tooltip:hover` |
| focus inside the wrapper | `.inst-tooltip:focus-within` |

The label is not removed from the flow but dimmed by opacity, and the mouse is
removed from it (`pointer-events: none`): a cursor moving towards the tooltip
is not intercepted by it and does not drop the hover.

:::warn
**A tooltip is clipped inside a container with `overflow: hidden`.** It is
positioned `position: absolute` against its own wrapper and does not live in
the top layer. That is the known price of a tooltip that works the same
everywhere and does not depend on support for anchor positioning.

Where clipping is possible — inside [a panel](../display/panel.md), a cell of
[a table](../display/table.md), [an accordion](../feedback/accordion.md) —
take [a popover](./popover.md): it is in the top layer, and nothing will cut it
off.
:::

## API

```api
```

## Related

```related
```
