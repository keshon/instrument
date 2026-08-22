# Anthropic / Claude Code — measured values

Everything on this page was read from `getComputedStyle` on a live, logged-in
`claude.ai` session on **2026-08-22**: the `/code` landing view, a chat
transcript, and a Claude Code session running parallel subagents.

Nothing here is estimated from a screenshot.

---

## The control ladder

Buttons on one screen, at one density, cross-tabulated height × radius × type.

| Height | Radius | Ratio | Type | Padding-x | Wears it |
|---|---|---|---|---|---|
| **20** | 5 | **0.250** | 12px / 17px | 6 · 4 · 0 | `Auto`, `Sonnet 5`, `Low`, add, mic, dictation |
| **24** | 6 | **0.250** | 13px / 16px | 6 · 4 · 0 | send, filter, `Default`, file rows |
| 24 | 8 | 0.333 | 13px | 4 | sidebar section headers |
| 30 | 6 | 0.200 | 14px / 21px | 4 | the `Home` / `Code` segmented pair |
| **32** | 8 | **0.250** | 14px / 21px | 2 · 0 | collapse, search, account row |

**radius = height ÷ 4** holds at 20, 24 and 32 — three rungs, one ratio. The
two exceptions are the segmented tab (0.20) and the sidebar section header
(0.33).

Icon-only at the smallest rung: a **12px glyph in a 20px box** — 0.60.

---

## Type

| Size | Leading | Ratio | Where |
|---|---|---|---|
| 9px | 13.5 | 1.50 | a badge |
| 12px | 15 | 1.25 | task-card metadata, group headers |
| 12px | 17 | 1.42 | session-row controls |
| 13px | 16 | 1.23 | sidebar rows |
| 13px | 19 | 1.46 | panel titles, task-card names |
| 14px | 18 | 1.29 | the composer field |
| 14px | 21 | 1.50 | header controls, segmented labels |
| 15px | 20 | 1.33 | **a user's message** |
| 16px | 24 | 1.50 | **an agent's answer body** |
| 20px | 25 | 1.25 | a page heading |

The split worth noting: **the ask is tight (1.33), the answer is loose (1.50)**,
and the answer is the larger of the two.

---

## Weight

Counted over every visible element carrying its own text (554 elements):

| Weight | Count | Share |
|---|---|---|
| 400 | 94 | 83% |
| 500 | 9 | 8% |
| 580 | 10 | 9% |

Icon glyphs sit at `533.25` and `577.75`. All of these are **variable-font axis
positions in `anthropic-sans`** — they are not reachable on a system font
stack, which offers 400 and 600 with nothing between.

The transferable fact is the **share**: 17% of text is emphasised at all.
Even a background task's own name is 400; hierarchy is carried by colour.

---

## Colour

| Role | Value |
|---|---|
| Page ground | `rgb(252,252,251)` |
| Recess (segmented track) | `rgb(237,236,232)` |
| Raised (selected segment) | `rgb(249,249,247)` |
| Card film | `rgba(11,11,11,0.04)` |
| User message film | `rgba(11,11,11,0.05)` |
| Code block | `rgb(255,255,255)` at 50% |
| Text primary | `rgb(11,11,11)` |
| Text primary, softened | `rgba(11,11,11,0.8)` |
| Text secondary | `rgb(82,81,78)` |
| Text muted | `rgb(137,135,129)` |
| Link | `rgb(24,79,149)` |
| Accent (ring arc) | `rgb(42,120,214)` |
| Placeholder | `rgba(11,11,11,0.25)` |
| Hairline (as shadow spread) | `rgba(11,11,11,0.1)` |

Three text tiers plus an alpha-softened primary. Note the screens measured carry
**no tone vocabulary** — no ok/warn/error readouts — so this is not evidence
about a status palette.

---

## Chrome

| | Count | Share |
|---|---|---|
| Visible elements | 554 | — |
| With a visible border | 8 | **1.4%** |
| With a shadow | 6 | 1.1% |

Borders are rare because the reference draws its edges as **spread shadows**:
`box-shadow: 0 0 0 1px rgba(11,11,11,0.1)` on the composer,
`inset 0 0 0 0.5px` on the raised segment. No layout box, sub-pixel control.

---

## Turn rendering

| | Value |
|---|---|
| User message | 15px / 20px (1.33), film `rgba(11,11,11,0.05)`, corner **12px**, padding **12px 16px** |
| Agent answer | 16px / 24px (1.50), **no surface at all** |
| Code block | corner 8px, white at 50%, 14px |

---

## Panel

| | Value |
|---|---|
| Border | **none** |
| Corner | 10px |
| Body padding | **0** |
| Header row | 32px tall, inset `0 4px` |
| Header title | 13px / 19px, weight 400, `rgba(11,11,11,0.8)` |
| Rule under the header | **none** |
| Title inset from panel edge | **8px** total (4 row + 4 title) |

---

## Composer

| | Value |
|---|---|
| Frame | corner 10px, `rgb(255 255 255 / 0.8)`, edge as `0 0 0 1px rgba(11,11,11,0.1)`, lift `0 2px 8px rgba(11,11,11,0.04)` |
| Height, one line | **44px** |
| Field | 14px / 18px, padding `13px 4px 13px 10px` |
| Send | corner 6px (the 24 rung), `align-self: flex-end`, wrapper padded `10px 10px 10px 4px` |
| Session row | 12px / 17px, corner 5px, padding `0 6px`, gaps 6px (left) and 5px (right) |

### The context ring

| | Value |
|---|---|
| Box | **12 × 12** |
| viewBox | `0 0 12 12`, `r="5"` |
| Stroke | **2** — 16.7% of the viewBox |
| Track | `rgba(11,11,11,0.1)` |
| Arc | `rgb(42,120,214)` |
| Mechanism | `stroke-dasharray: 31.41592653589793` (2πr) + `stroke-dashoffset` |
| Cap | `round` |
| Rotation | **none** — the arc begins at 3 o'clock |

---

## Background task card

| | Value |
|---|---|
| Fill | `rgba(11,11,11,0.04)` |
| Border | **none** |
| Corner | 8px |
| Padding | 8px, 5px between registers |
| Height | 78px for three lines |
| Name | 13px / 19px, weight **400**, muted |
| Metadata | 12px / 15px, muted, `flex-wrap` with `gap: 3px 8px`, `align-items: baseline` |
| Link | 12px, `rgb(24,79,149)`, corner 3px |
