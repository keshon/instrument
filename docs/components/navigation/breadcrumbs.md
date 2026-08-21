---
title: Breadcrumbs
group: Navigation
layout: component
source: src/layout.css
api:
  - { name: "inst-crumbs", kind: "class", doc: "The list of the path. Put on an `<ol>`" }
  - { name: "--space-2", kind: "token" }
  - { name: "--text-xs", kind: "token" }
  - { name: "--size-chevron", kind: "token" }
  - { name: "--text-muted", kind: "token" }
  - { name: "--text-faint", kind: "token" }
  - { name: "--text-primary", kind: "token" }
---

The path from the root to the current screen. It answers "where am I and how do
I get one level up from here" rather than "what other sections are there".

```html preview
<nav aria-label="Breadcrumbs">
  <ol class="inst-crumbs">
    <li><a href="#">Projects</a></li>
    <li><a href="#">worldgen</a></li>
    <li><span aria-current="page">Run #4127</span></li>
  </ol>
</nav>
```

## Contract

Breadcrumbs are an **ordered list inside a `<nav>`**, and both are required:
the `<ol>` reports the order and the length of the path, and the
`<nav aria-label>` gives a landmark that is distinguishable from the other
navigations on the screen.

| What | Required | Why |
|---|---|---|
| A `<nav>` + an `aria-label` | yes | There are several navigations on a screen; with no name the breadcrumbs are indistinguishable from the other landmarks |
| An `<ol>` with `<li>` | yes | The path is ordered, and the order is part of the meaning. A `<div>` loses it |
| An `aria-current="page"` on the last one | yes | It marks the current screen and gets `--text-primary` instead of the muted colour |
| The last element is **not** an `<a>` | yes | A link to the page you are already on is a false action |

The item and the link have no class of their own: `& > li` and `& a` inside the
list are styled. The markup is the API.

### Accessibility

| | |
|---|---|
| The keyboard | Native links. `Tab` traverses, `Enter` goes |
| The landmark | A `<nav aria-label="Breadcrumbs">` — a point of its own in the list of landmarks |
| The separator is not read | The chevron is a pseudo-element; it is in neither the text nor the accessibility tree |
| Contrast | The path is `--text-muted` (4.5:1) rather than `--text-faint`: `--text-faint` is the threshold of decoration, and text that gets read is not painted with it. The chevron is the other way round — decoration |
| Hover | The link gets an underline **and** the full colour: an underline works where colour is not told apart |
| Type size | `--text-xs`. The path is quieter than the screen title it precedes |

## Anatomy

The chevron between the items is drawn by `li + li::before` — a pseudo-element
rather than a character in the markup. Two consequences it was made for:

- the separator **does not reach the text**: a screen reader and a copy get
  "Projects worldgen Run #4127" rather than a line with angle brackets;
- the shape comes as a mask and the colour as the `--text-faint` token. There
  is no colour inside a data URI in the library, so the chevron survives a
  change of theme and forced-colours mode with no second copy of the picture.

The strip wraps (`flex-wrap`): a long path in a narrow header goes onto a
second line rather than pushing the document past the viewport.

## API

```api
```

## Related

```related
```
