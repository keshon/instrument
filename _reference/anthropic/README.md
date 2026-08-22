# Reference snapshots — Anthropic / Claude Code UI

Static captures of the reference the kit is modelled on, taken **2026-08-22**
from a logged-in `claude.ai` session: the `/code` landing view and a Claude Code
session running parallel subagents.

They exist so a later conversation can measure against the real thing without
needing a live browser, a login, or a screenshot to squint at.

## How they were made

A serializer walked each component's DOM, read `getComputedStyle` on every node,
and inlined the properties that carry design intent — layout, spacing, corners,
colour, type, borders, shadows. Framework classes and `data-*` were stripped.

**Sizes were deliberately NOT inlined** except on the root of each capture and
on `<svg>` glyphs. A snapshot full of `width: 377.562px` renders faithfully and
teaches nothing; the measured geometry lives in
[`tokens.md`](tokens.md) and in `docs/internal/AUDIT-2026-08-22.md`.

## What does not survive the capture

| | |
|---|---|
| **Icons** | The reference draws them with an **icon font**, not inline SVG — the glyph spans come across empty. Every empty `<span aria-hidden="true">` with a `font-weight` near 533 or 578 is an icon. |
| **The typeface** | `anthropic-sans` is a webfont and is not embedded here; the captures fall back to the system stack. Metrics shift a little. |
| **Variable weights** | The reference sets weights like `533.25`, `577.75` and `580` — axis positions in a variable font. A system stack snaps them to 400/700. |
| **State** | Hover, focus and open states are not captured; each file is one resting frame. |

## The files

| File | What it is |
|---|---|
| `composer.html` | The prompt field and its session row — field, send at the trailing edge, `Auto · add · mic · settings` on the left, `Sonnet 5 · Low · context ring` on the right. Contains the **context ring as real SVG**. |
| `background-task-card.html` | One finished subagent in the Background tasks panel. Three stacked registers: name, kind·state·duration, cost·link. |
| `panel-header.html` | The Background tasks panel header — a 32px row inset by 4px, no rule under it. |
| `segmented.html` | The `Home` / `Code` segmented pair, including the three-layer shadow on the selected segment. |
| `tokens.md` | Every number measured, in tables: the control ladder, type, colour, chrome. |

## Reading them

Open in any browser. Each is a fragment on the reference's own page ground
(`rgb(252,252,251)`) at the width it was captured at.
