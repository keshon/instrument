---
title: Tabs
group: Navigation
layout: component
source: src/layout.css
js: Switching the panel. The arrows and the roving tabindex are done by `instrument.js`
api:
  - { name: "inst-tabs", kind: "class", doc: "The strip of tabs. A line below and horizontal scrolling with no visible scrollbar" }
  - { name: "inst-tab", kind: "class", doc: "A tab. A link or a button with `role=\"tab\"`" }
  - { name: "inst-tabs--chips", kind: "modifier", doc: "Tabs as chips: the chosen one by a fill rather than by a bar. For OPEN DOCUMENTS, of which there are any number and which close" }
  - { name: "inst-tab-close", kind: "class", doc: "The close cross. A real `<button>`, and an `aria-label` is required. A closable tab is therefore a `<div role=\"tab\">` rather than a `<button>`. The area grows to `--tap-min` while the shape stays small" }
  - { name: "data-state", kind: "attribute", value: "modified", doc: "On `.inst-tab`: a dot in the place of the cross while the file is unsaved. The cross comes back under the cursor" }
  - { name: "--space-3", kind: "token" }
  - { name: "--space-4", kind: "token" }
  - { name: "--space-6", kind: "token" }
  - { name: "--text-sm", kind: "token" }
  - { name: "--weight-medium", kind: "token" }
  - { name: "--size-marker", kind: "token" }
  - { name: "--hairline", kind: "token" }
  - { name: "--border", kind: "token" }
  - { name: "--surface-inverse-soft", kind: "token", doc: "The surface of a chosen chip: dark in a light theme, light in a dark one. Soft — a permanent selection does not take the full inversion" }
  - { name: "--control-h-sm", kind: "token", doc: "The height of a chip. Setting it by padding gave a row taller than the button beside it" }
  - { name: "--control-h-xs", kind: "token", doc: "The box of the cross. A square with a field around the glyph" }
  - { name: "--text-on-inverse", kind: "token", doc: "The label of a chosen chip" }
  - { name: "--accent-solid", kind: "token" }
  - { name: "--text-primary", kind: "token" }
  - { name: "--text-secondary", kind: "token" }
---

A horizontal switch of views inside one section. By default on links:
switching a screen is routing, and it needs no JS.

```html preview
<nav class="inst-tabs" aria-label="Sections of the runs">
  <a class="inst-tab" href="#" aria-current="page">Active</a>
  <a class="inst-tab" href="#">History</a>
  <a class="inst-tab" href="#">Schedule</a>
</nav>
```

## Contract

| What | Required | Why |
|---|---|---|
| A `<nav>` + an `aria-label` on the routing variant | yes | There are several navigations on a screen, and with no name they are indistinguishable |
| An `aria-current="page"` | yes | The active view when switching by address |
| A `role="tablist"` + `role="tab"` + `aria-selected` + `aria-controls` | yes, on **real** tabs only | A role with no arrow navigation and no panel is a promise nobody will keep |
| A `type="button"` on a tab that is a button | yes | Otherwise inside a form it will submit it |
| The active tab by an underline as well as by colour | yes | An `aria-current` or an `aria-selected` tells a screen reader the same |

### Accessibility

| | |
|---|---|
| The keyboard, routing | Native links: `Tab` traverses, `Enter` goes |
| The keyboard, real tabs | `Tab` enters the strip once and the arrows switch — that is the **application's JS**. Without it a `role="tablist"` may not be set |
| Colour is not the only carrier | The active tab carries an underline and an `aria-current` / `aria-selected` |
| Scrolling | The strip scrolls horizontally, the scrollbar is hidden visually, but scrolling by wheel and from the keyboard stays |
| Contrast | An inactive label `--text-secondary` holds 4.5:1 in five themes |
| Focus | The `:focus-visible` ring from the base layer; the strip of tabs does not clip it |

## Variants

Links or real tabs are two different cases, and their markup differs. The one
chosen is the one that is truthful:

| The case | The markup | The mark |
|---|---|---|
| Switching a screen | `<nav class="inst-tabs">` + `<a class="inst-tab" aria-current="page">` | Every view has an address of its own. A middle click opens it in a new browser tab |
| Panels inside one page | `role="tablist"` + `role="tab"` + `aria-selected` + `aria-controls` | The address does not change, the content is swapped in place |

The styling catches both: the library's selector reacts to
`[aria-current]:not([aria-current="false"])` and to `[aria-selected="true"]`
alike. So the question "which class to put" does not arise — the question that
arises is "which markup is truthful".

### Real tabs

A `role="tablist"` is set **only** if the tabs are real. The role promises a
screen reader arrow navigation and a tie to a panel; an unbacked promise is
worse than none.

```html preview
<div class="inst-tabs" role="tablist" aria-label="Properties of the object">
  <button class="inst-tab" type="button" role="tab" id="t-geom"
          aria-selected="true" aria-controls="p-geom" tabindex="0">Geometry</button>
  <button class="inst-tab" type="button" role="tab" id="t-mat"
          aria-selected="false" aria-controls="p-mat" tabindex="-1">Material</button>
  <button class="inst-tab" type="button" role="tab" id="t-phys"
          aria-selected="false" aria-controls="p-phys" tabindex="-1">Physics</button>
</div>
<div id="p-geom" role="tabpanel" aria-labelledby="t-geom" tabindex="0">12 480 vertices, 6 240 polygons.</div>
<div id="p-mat" role="tabpanel" aria-labelledby="t-mat" tabindex="0" hidden>Standard PBR, two textures.</div>
<div id="p-phys" role="tabpanel" aria-labelledby="t-phys" tabindex="0" hidden>A convex hull, a mass of 4.2 kg.</div>
```

The arrows, `Home`, `End` and the roving `tabindex` are
[`instrument.js`](../../foundations/behavior.md)'s. Showing the right panel is
the application's: the library does not know which panel is behind which tab.

## States

```html preview
<nav class="inst-tabs" aria-label="An example of an active tab">
  <a class="inst-tab" href="#" aria-current="page">Active</a>
  <a class="inst-tab" href="#">Ordinary</a>
  <a class="inst-tab" href="#">One more</a>
</nav>
```

| The mark of the active one | What it gives |
|---|---|
| An underline | A `border-block-end` of width `--size-marker` in the colour `--accent-solid`. A carrier on a par with colour |
| The colour of the label | `--text-primary` instead of `--text-secondary` |

The weight **does not change** on selection: it would change the width of the
label, and a measurement on this reference gave neighbours shifting 2–3px — a
choice by mouse would move what is being aimed at. Two carriers are left and
neither is colour alone: an active tab differs from an ordinary one by the
**presence** of the bar, and from a hovered one by tone.

The bar of a tab covers the line of the container by a negative margin of
`--hairline` rather than by a shift of half a pixel: a browser does not draw a
border thinner than a physical pixel, and the remainder would show at every
joint.

## JS

Include the module once per page.

```html
<script type="module" src="instrument.js"></script>
```

### What `instrument.js` does

Carries out the contract of `role="tablist"`: `←` and `→` between the tabs,
`Home`, `End`, a traversal that wraps, one `Tab` for the whole strip.
`aria-selected` moves on an arrow and on a click alike.

### Events

`inst:select` bubbles from the chosen tab, and its `detail` is `{ value }`: the
`data-value` if there is one, otherwise the label.

```js
tabs.addEventListener('inst:select', (e) => {
  for (const panel of panels) {
    panel.hidden = panel.id !== e.detail.value;
  }
});
```

### What is left to the application

Showing the right panel and hiding the rest.

## Composition

### With an icon and a counter

```html preview
<nav class="inst-tabs" aria-label="An example with a counter">
  <a class="inst-tab" href="#" aria-current="page">
    <svg class="inst-icon" aria-hidden="true"><use href="#i-list"/></svg>Queue
    <span class="inst-badge">7</span></a>
  <a class="inst-tab" href="#">History</a>
</nav>
```

The gap inside a tab is `--space-3` and between tabs `--space-6`: the distance
between neighbours has to be noticeably larger than the distance within one,
otherwise the icon of the next reads as part of the previous.

### Tabs as chips

A second kind of tab, and it is another job. A tab as a bar switches a
**section**: there are few sections, they are set in advance and they do not
close. A chip holds an **open document**: there are any number of them, they
appear and close as work goes on, and each has a state of its own of being
unsaved.

```html preview
<div class="inst-tabs inst-tabs--chips" role="tablist" aria-label="Open files">
  <div class="inst-tab" role="tab" aria-selected="false" tabindex="-1">tokens.css
    <button class="inst-tab-close" type="button" aria-label="Close tokens.css">
      <svg class="inst-icon" aria-hidden="true"><use href="#i-close"/></svg></button></div>
  <div class="inst-tab" role="tab" aria-selected="true" tabindex="0" data-state="modified">layout.css
    <button class="inst-tab-close" type="button" aria-label="Close layout.css">
      <svg class="inst-icon" aria-hidden="true"><use href="#i-close"/></svg></button></div>
  <div class="inst-tab" role="tab" aria-selected="false" tabindex="-1">components.css
    <button class="inst-tab-close" type="button" aria-label="Close components.css">
      <svg class="inst-icon" aria-hidden="true"><use href="#i-close"/></svg></button></div>
</div>
```

A bar under such a row stops working: it draws one line under all of them, and
at twenty tabs the chosen one has to be searched for. A fill reads at once and
at any length of row.

A chosen chip takes a **soft inverted surface** — dark in a light theme, light
in a dark one. Raised is too close: the row already lies on a raised surface. A
full inversion is too loud for something that stands there all the time. The
height comes from the control ladder (`--control-h-sm`) rather than from
padding.

**The cross is an overlay at the far edge rather than a place in the line** — a
name is read all the time and a tab is closed once. It appears from the
keyboard too: an invisible button that `Tab` has reached is a trap. Room for
the "unsaved" dot is **reserved**, because that state stands there permanently.

A chip is rounded on all four sides and lies on the strip rather than growing
into the content below.

**A closable tab is a `<div role="tab">` rather than a `<button>`.** The cross
has to be a real button: tabs are closed from the keyboard, and the cross has a
tap target of its own. And HTML does not let a button inside a button at all —
interactive content is forbidden as a descendant of a `<button>` or an `<a>`. A
tab as a bar stays a button or a link: there is nothing to close there.

| What | Required | Why |
|---|---|---|
| A `<div role="tab">` on a closable one | yes | A real cross inside a `<button>` is invalid markup |
| A roving `tabindex`: `0` on the chosen one, `-1` on the rest | yes | Thirty open files must not cost thirty presses of `Tab` |
| An `aria-label` on the cross | yes | A button with no text has no accessible name, and "close" with no file name does not say what exactly |
| A `data-state="modified"` | no | A dot in the place of the cross while the file is unsaved. Room is reserved for it, and the cross comes back under the cursor |
| The rounding from `--radius-control-md` | — | A chip is a control: in the compact mode it is rounded less, along with the button |

## API

```api
```

## Related

```related
```
