---
title: Pixel audit
group: About
source: tools/audit.js
---

`tools/audit.js` measures contrast and tap targets by what is drawn in the
browser. It is an addition to the textual checks: those parse `tokens.css`, the
pages of the documentation and the built files, but between the declaration of
a token and a pixel stand nesting, translucent layers, the application's
cascade and the chosen theme.

The check walks every node that carries text, takes its real colour and its
real background, gathering the alpha up the ancestors, and computes the
contrast. A colour in any notation — `oklch`, `color-mix`, `light-dark` — is
resolved to sRGB by the browser itself: the value is painted on a canvas one
pixel in size and read back.

It has to be run **on your own screen**: a list of pairs for a textual check
cannot be drawn up in advance, because the number of combinations depends on
how the layers are nested in a particular application.

## How to run it

Open your own screen, paste
[`tools/audit.js`](https://github.com/keshon/instrument/blob/master/tools/audit.js)
into the browser console and call:

```js
kitAudit.run()            // the whole page: every axis at once
kitAudit.run('#main')     // a part of it
kitAudit.contrast()       // one pass at the current theme and accent
kitAudit.targets()        // one pass at the current scale and density
kitAudit.proportion()     // the same pass: does the text fit its box
```

`run()` prints a summary as a table and unfolds the violations by cell. It
returns the same as an object, if the result has to be parsed by a program.

The axes are not multiplied all by all, and that is deliberate. Contrast goes
by **theme × accent**: an accent redefines the `--a-*` ramp, and a label on an
accent fill is different in every theme. Targets and proportions go by **scale
× density**: both axes move geometry and both round to whole pixels
independently of each other. Multiplying one by the other would stretch the run
twentyfold to repeat the same numbers.

## What it measures

| | |
|---|---|
| The contrast of text | 4.5:1, and 3:1 for large text (24px, or 18.66px in semibold). The background is gathered up the ancestors, and translucent layers are composited in the order of drawing |
| Tap targets | 24×24 under WCAG 2.2 AA (2.5.8) with two exceptions of the criterion: by spacing — a target smaller than the norm counts when `S + G ≥ 24`; inline — a link inside a line of text is exempt entirely |
| Proportions | Whether a label fits its column and whether it has wrapped onto a second line; the ink of a glyph against the cap height of the label beside it. This is the one thing that cannot be derived from the tokens: the width of a line depends on the family, the weight and the language |
| Themes and accents | 5 × 4 for contrast. A theme set on a subtree stays its own: the root one does not override it |
| Scales and densities | 5 × 3 for targets and proportions |

## What it does not see

| | |
|---|---|
| The hidden | A closed popover, a folded panel, a second tab — they are not in the layout. Open them and run it again |
| A background as an image | Only colour is counted. Text over an image is skipped by the check |
| Order and meaning | It measures pixels. The keyboard, the names of roles and the order of reading stay with a human |
| States | Hover, focus and press live at the moment of an interaction. What is measured is what is visible now |

## When the result is not to be trusted

The result depends on the moment of the measurement. Five cases give false
violations.

**A measurement during an animation.** A change of theme is animated, and a
colour read in the same frame is an interpolation. The mark: the browser
returns the colour as `oklab(…)` rather than `oklch(…)`. `run()` suppresses the
transitions itself, `contrast()` does not.

**Fractional widths.** A button of a segmented control in the compact mode is
22px of height and 2px of gap, exactly 24. At a fractional width the browser
gives back a gap of 1.9999, and the norm formally fails. A tolerance of 0.1px
suppresses that; there is no margin there, and an application that squeezes the
gap slightly will break WCAG in silence.

**A tap area grown by a pseudo-element.** Small controls stay small (the glyph
of a checkbox is 13px), and what grows them to the norm is an invisible
`::before`. The check takes the largest of the element and its
absolutely-positioned pseudo-elements; an inline pseudo-element decorates
rather than widening the target.

**A sticky header over a scroll.** The spacing is measured in the coordinates
of the window, and an element travelling over another reads as standing flush
against it. Intersecting boxes are layers rather than neighbours, and the check
does not compare them.

**Visibly hidden controls.** A "skip to content" link and labels for screen
readers live as a box one pixel in size with a `clip-path`, and become a real
target on focus alone. The mark is the `clip-path: inset(50%)` clipping rather
than the size: an element with `border-box` cannot become thinner than its
borders, so a declared 1px is drawn as 2px.

## Related

[Colour](../foundations/colors.md)
[Density](../foundations/density.md)
[Design principles](./design-principles.md)
