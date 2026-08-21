---
title: File change
group: Data display
layout: component
source: src/data.css
api:
  - { name: "inst-change", kind: "class", doc: "The glyph of a change. An outline with a sign inside it" }
  - { name: "data-change", kind: "attribute", value: "added · modified · deleted · renamed · conflict", doc: "What happened to the file. An axis of its own rather than the tonal vocabulary" }
  - { name: "aria-label", kind: "attribute", doc: "**Required.** A screen reader reads «+» as \"plus\" rather than as \"file added\"" }
  - { name: "--size-icon", kind: "token", doc: "The box of the glyph" }
  - { name: "--text-2xs", kind: "token", doc: "The size of the sign" }
  - { name: "--chart-3", kind: "token" }
  - { name: "--chart-4", kind: "token" }
  - { name: "--chart-5", kind: "token" }
  - { name: "--chart-6", kind: "token" }
  - { name: "--change-ink", kind: "variable", doc: "The colour of the outline. The threshold of a mark, 3:1: an outline is not text" }
  - { name: "--text-secondary", kind: "token", doc: "The colour of the sign. The threshold of text, 4.5:1: a sign gets read" }
  - { name: "--warn-text", kind: "token", doc: "Conflict only: the one value of the axis that is also a judgement. It paints both the outline and the sign" }
---

What happened to a file in a list of changes: added, modified, deleted,
renamed, conflict.

```html preview
<div class="inst-stack inst-stack--tight">
  <span class="inst-change" data-change="added" aria-label="Added">+</span>
  <span class="inst-change" data-change="modified" aria-label="Modified">M</span>
  <span class="inst-change" data-change="deleted" aria-label="Deleted">−</span>
  <span class="inst-change" data-change="renamed" aria-label="Renamed">R</span>
  <span class="inst-change" data-change="conflict" aria-label="Conflict">!</span>
</div>
```

## Contract

The sign stands as the content and the name in the `aria-label`. They are
parted because the sign is data and the name is speech: «−» aloud has to sound
"deleted". The `aria-label` overrides the content, so the sign is not read
twice.

| What | Required | Why |
|---|---|---|
| An `aria-label` in words | yes | The sign «+» will be read as "plus". What exactly happened does not follow from that |
| The sign as content | yes | Colour has no right to be the only carrier. In black-and-white print and under colour blindness only the sign is left |
| A value from the vocabulary | yes | Five values, and the list is closed. A sixth will silently do nothing |

### Accessibility

| | |
|---|---|
| The accessible name | From the `aria-label`: it overrides the content, and "plus" is not said aloud |
| Colour is not the only one | The sign inside tells an addition from a deletion where colour has gone: print, colour blindness, high contrast |
| Contrast | Different carriers have different thresholds: the outline is a mark at 3:1, the sign is text at 4.5:1. Checked by a sweep over the pixels in five themes |
| The role | Neither `img` nor `status`: this is a mark beside a file name, and it is spoken together with it |

## API

```api
```

## Related

```related
```
