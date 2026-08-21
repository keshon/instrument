---
title: Avatar
group: Data display
layout: component
source: src/data.css
api:
  - { name: "inst-avatar", kind: "class", doc: "One avatar. The size is `--control-h-md`" }
  - { name: "inst-avatars", kind: "class", doc: "A group with overlap" }
  - { name: "inst-avatars-more", kind: "class", doc: "The \"N more\" counter at the end of a group" }
  - { name: "inst-avatar--sm", kind: "modifier", doc: "The size" }
  - { name: "inst-avatar--lg", kind: "modifier", doc: "The size" }
  - { name: "--control-h-sm/md/lg", kind: "token" }
  - { name: "--radius-full", kind: "token" }
  - { name: "--surface-sunken", kind: "token" }
  - { name: "--surface-raised", kind: "token" }
  - { name: "--text-xs", kind: "token" }
  - { name: "--weight-medium", kind: "token" }
  - { name: "--size-marker", kind: "token" }
---

The identifying mark of a person or an agent: an image or initials.

```html preview
<span class="inst-avatar">IS</span>
<span class="inst-avatar">AP</span>
<span class="inst-avatar">MK</span>
```

## Contract

| What | Required | Why |
|---|---|---|
| `aria-hidden="true"` if the name is beside it as text | yes | Otherwise the avatar is spoken as "I S" next to a name already read |
| An `alt` or an `aria-label` if there is no name beside it | yes | Initials are not an accessible name: a screen reader reads them letter by letter |
| An `aria-label` on the group, by its make-up or by a number | yes | "7 participants" rather than a string of initials |

### Accessibility

| | |
|---|---|
| An avatar is not a label | If the name is beside it as text, the avatar is decorative: `aria-hidden="true"`. Otherwise it is spoken as "I S" next to a name already read |
| An avatar instead of a label | Then it carries the name: `<img alt="Innokentiy Sokolov">` or an `aria-label` on the container |
| Initials | They are not an accessible name: a screen reader reads "IS" letter by letter |
| A group | Give the container an `aria-label` with the full make-up or a number: "7 participants". Otherwise a string of initials is spoken |
| The counter | `inst-avatars-more` holds "+4" as text — that is data rather than decoration |
| Contrast | The initials take `--text-secondary` on `--surface-sunken` and hold 4.5:1 in five themes |

## Sizes

```html preview
<span class="inst-avatar inst-avatar--sm">W1</span>
<span class="inst-avatar">IS</span>
<span class="inst-avatar inst-avatar--lg">B4</span>
```

The sizes come from the tier of roles (`--control-h-*`), so an avatar and a
[button](../actions/button.md) of the same size standing beside it physically
cannot part in height.

## Composition

### A group with overlap

```html preview
<span class="inst-avatars" aria-label="7 participants">
  <span class="inst-avatar">IS</span>
  <span class="inst-avatar">AP</span>
  <span class="inst-avatar">MK</span>
  <span class="inst-avatars-more">+4</span>
</span>
```

The ring around every avatar in a group is not an outline but **the colour of
the surface under it**: it parts the neighbours from one another. On another
surface a group will call for that ring to be overridden, otherwise a halo
appears.

## API

```api
```

## Related

```related
```
