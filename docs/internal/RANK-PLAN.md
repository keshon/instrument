# Rank — implementation plan

**SUPERSEDED. Implemented on 2026-09-03; kept as the record of what was decided
before the code existed, not as a description of the code.**

Where each part of it now lives:

| | |
|---|---|
| The contract — vocabulary, four laws, where rank is legal | [design principles](../about/design-principles.md#state-and-variant) |
| Containment, and why it is a film | [elevation](../foundations/elevation.md) |
| What rank does to a panel, and to a card | [panel](../components/display/panel.md) · [card](../components/display/card.md) |
| Where the heading ladder went | [typography](../foundations/typography.md) |
| The derivations behind the choices | [DECISIONS](./DECISIONS.md) |
| The invariants themselves | `cmd/contrast`, `cmd/proportion`, `cmd/registry`, `cmd/mutate`, `kitAudit.composition()` |

**Four things below turned out to be wrong, and each is marked in place.** They
are left standing rather than edited out, on the precedent of the audits in this
directory: a plan that quietly rewrites itself teaches nothing about how much of
a plan survives contact with a browser.

1. `@property { inherits: false }` cannot carry the channel — a region's title
   is its grandchild, and the value stops one element short. The mechanism is a
   `:where()` reset at every region boundary.
2. The five tokens do not all belong in the role tier. Two of them are colours
   and belong in the semantics, which is also the only tier `cmd/contrast` can
   see.
3. A region takes at most ONE film. The two rules collide on specificity, and
   films stack in the paint rather than in the value.
4. Moving the heading ladder into `.inst-prose` leaves a bare heading at the
   USER AGENT's size, not the body size. `font-size: inherit` is the line the
   plan meant.

A fifth was found by the gates rather than by measurement: the composition
mutation pointed at a page with no card inside a ranked region, so it was missed
on its first run and guarded nothing.

---

Turning the accepted conclusions of the 2026-09-03 audit into work.
Written against commit `d39cba0`, instrument 0.8.0.

This document decides. An engineer following it should not have to make an
architectural choice of their own; where a choice was open, the choice and its
reason are written here.

**Scope.** One new semantic axis (`data-rank`), one generalised existing rule
(relative containment depth), five new role tokens, four inheriting channel
variables with a `:where()` reset at every region boundary, one closed list of
three components. No new components. No new layer. No new file in `src/`.

---

## 0. The four laws this plan is held to

1. **Depth owns direction.** A region inside a region recedes. Always, and
   without consulting rank.
2. **Rank owns presence** and may modulate ground by at most one recess step,
   in the receding direction only. Rank never lifts.
3. **Rank touches a region's chrome, never its data.** Title, edge, ground.
   Not the body text, not the numbers, not the rows.
4. **Rank does not inherit.** A nested region with no rank of its own is
   `default`, whatever its container is.

Law 2's second clause is not taste. `--surface-raised` resolves to `--n-0` =
`oklch(1.000 …)` in all three light themes — the ceiling of the ramp. There is
no lift available below depth 0 in a light theme, and `surfaces.css:130` already
says so in the comment that made a card recede rather than rise. A `lead` that
lifted would work in `dark`/`dark-soft` and silently do nothing in the three
light themes: exactly the class of theme divergence `light-dark()` exists to
make impossible.

---

## 1. Existing architecture — what is touched and what is reused

### Files that change

| File | Layer | What changes |
|---|---|---|
| `src/tokens.css` | `kit.tokens` | 5 role tokens added to the tier-3 block |
| `src/surfaces.css` | `kit.components` | 2 `:where()` resets; 2 `[data-rank]` blocks; panel and card read the channel; `.inst-panel .inst-card` retired |
| `src/layout.css` | `kit.layout` | `.inst-section` / `.inst-section-title` read the channel |
| `src/base.css` | `kit.base` | four `font-size` lines move from bare `h1–h6` into `.inst-prose` |
| `tools/cmd/contrast/main.go` | — | 5 `kase` rows |
| `tools/cmd/proportion/main.go` | — | 2 `rule` rows |
| `tools/cmd/mutate/main.go` | — | 4 `mutation` rows |
| `tools/cmd/docscheck/main.go` | — | ~6 lines: generalise the base-value list |
| `tools/cmd/registry/main.go` | — | 1 axis-style list check |
| `tools/audit.js` | — | new `composition()` section |
| `tools/audit-run.mjs` | — | 1 `MUTATIONS` row |
| `components.json` | — | `"region": true` on three entries |

### Files that do **not** change, and why

- `src/agent.css`, `src/table.css`, `src/forms.css`, `src/data.css`,
  `src/actions.css`, `src/rows.css`, `src/text.css`, `src/feedback.css`,
  `src/status.css`, `src/overlay.css` — nothing in them is a region.
- `src/forced.css` — rank resolves to tokens that forced-colors already
  overrides; a `support` region loses its film and its edge there exactly as a
  panel does today. No new entry.
- `src/print.css` — same reasoning. Rank is screen presence; print already
  flattens grounds.
- `src/motion.css` — rank animates nothing.
- The density and scale cells in `tokens.css` — see §7. The new tokens are
  references into ladders those cells already move. **Zero cell edits.**

### Mechanisms reused rather than reinvented

| Need | Existing mechanism | Where |
|---|---|---|
| A channel variable a component reads, with the default at the point of consumption | `var(--tone-bg, var(--surface-sunken))`, `var(--icon-size, var(--size-icon))`, `var(--panel-body-pad, var(--pad-panel))` | `feedback.css:101`, `tone.css`, `surfaces.css:82` |
| A default that any override must beat | `:where()` at zero specificity — *"a default with the same specificity as an override is not a default"* | `design-principles.md`, and throughout `tokens.css` tier 3 |
| A closed attribute vocabulary derived from selectors | `data-tone` / `data-state` / `data-kind` / `data-when` | `docscheck/main.go:182` |
| A base value with no styling, spellable in markup | `queued` · `todo` · `approved` | `docscheck/main.go:199` |
| A cross-cutting list in CSS cross-checked against a registry | the forced / print / flow lists | `registry/main.go` |
| A relative ground counted from what is underneath | `--surface-recessed` (an alpha film, sign flipped by theme) | `tokens.css` |
| A tone reaching a **ground** and an **edge** | `.inst-note`, `.inst-banner` | `feedback.css:101,116` |

> **On `--surface-recess` / `--surface-lift`.** Both were proposed in review and
> both are refused. The recess mechanism already exists as `--surface-recessed`,
> is already an alpha film counted from the backdrop, already flips sign by
> theme, and is already guarded by three `cmd/contrast` pairs and one mutation.
> Nothing is added. `--surface-lift` is not added at all, because law 2 forbids
> the operation it would name.

---

## 2. The variable model

### 2.1 Source

One attribute, one closed vocabulary, three values:

```
data-rank = "lead" | "default" | "support"
```

### 2.2 Role tokens (tier 3, `tokens.css`, inside `:where(:root)`)

Five, and only five. A rank gets a token **only where it differs from
`default`** — a token that restates the default is the duplicate this kit spends
its budget removing.

```css
/* The name of a region, by rank. References into the type ladder rather than
   numbers: scale already moves the ladder, so these move with it. */
--region-title-lead:    var(--text-md);
--region-title-default: var(--text-sm);
--region-title-support: var(--text-xs);

/* A support region recedes by ONE film, the same film a control recesses by.
   There is no --region-ground-lead: lead does not lift, and there is nothing
   to name. */
--region-ground-support: var(--surface-recessed);

/* A support region gives up its edge. Chrome is the cheapest thing to spend
   and the first thing to withdraw. */
--region-edge-support:   transparent;
```

Ink takes no token. `--text-primary` and `--text-secondary` are already a gated
four-tier ladder; naming one of its steps a second time is a duplicate. The
`[data-rank]` block references them directly.

**Why the title rungs must be tokens rather than values inside a `[data-rank]`
block.** `cmd/proportion` resolves names out of `tokens.css` and nothing else.
A rung written as `font-size: var(--text-md)` inside `surfaces.css` is
invisible to the gate, and the ladder would be unguarded. This is what makes
the gate rows in §5 possible with **no change to the proportion tool**.

### 2.3 Channel variables — ordinary inheritance, reset at every region boundary

**MEASURED CORRECTION.** An earlier draft of this plan specified
`@property { inherits: false }` on all four channels, by analogy with
`--flow-self`. That was tested in the browser against the real kit and **it does
not work**, for a reason the analogy hides.

`--flow-self` is declared on `> *` and read *by that same element*.
`inherits: false` there stops the value going any deeper than the child it was
meant for. Rank is declared on the **region** and read by its **title**, which
is a grandchild:

```
.inst-panel[data-rank="lead"]        --region-title-size: 16px   ← declared here
  > .inst-panel-header               (empty)                     ← stopped
      > .inst-panel-title            (empty)                     ← never arrives
```

Measured chain, all three values read with `getComputedStyle`. With
`inherits: false` the title resolves the fallback in every case, so a lead
panel's name renders at the default rung and **the rank axis is inert for its
primary carrier**.

The mechanism that does work is the one the kit already states as law:

> *Defaults are declared through `:where()`. A default with the same specificity
> as an override is not a default.* — `design-principles.md`

So: the channels are **ordinary inheriting custom properties**, and every region
**resets them at its own boundary** at zero specificity.

```css
/* No @property. The reset is the mechanism. */
:where(.inst-panel, .inst-card, .inst-section) {
  --region-title-size: var(--region-title-default);
  --region-title-ink:  var(--text-primary);   /* --text-muted on the section */
  --region-ground:     var(--surface-raised); /* transparent on the section */
  --region-edge:       var(--border);         /* transparent on the section */
}
```

Zero specificity, so any `[data-rank]` block at `(0,1,0)` beats it. The value
inherits freely down to the title, and stops at the next region because that
region resets it. Law 4 holds by the reset rather than by registration.

Measured, on the real kit:

| Element | `--region-title-size` | |
|---|---|---|
| `.inst-panel[data-rank="lead"]` title | **16px** | rank reaches its carrier |
| unranked `.inst-card` title, inside that lead panel | **14px** | **the leak is stopped** |
| `.inst-card[data-rank="support"]` title, inside that lead panel | **12.5px** | rank overrides the reset |
| unranked `.inst-panel` title | **14px** | unchanged from today |

`--region-ground` and `--region-edge` are read on the region itself and would
have worked either way; they use the same mechanism so there is one rule and not
two. Nothing outside the three region components reads any of the four, so
inheritance into a badge or a paragraph is inert.

There is still **no `[data-rank="default"]` CSS block** — the `:where()` reset
already is it, and a second block at higher specificity would restate it. See
§3.2 for how `default` stays spellable.

**Consequence for the gate.** Law 4's mutation is no longer
`inherits: false → true`. It is the removal of a class from the `:where()` reset
list — see §5.7.

### 2.4 Defaults, stated once, at the point of consumption

| Channel | Fallback | Consumer |
|---|---|---|
| `--region-title-size` | `var(--region-title-default)` | `.inst-panel-title`, `.inst-card-title`, `.inst-section-title` |
| `--region-title-ink` | `var(--text-primary)` on panel/card, `var(--text-muted)` on section | the same three |
| `--region-ground` | `var(--surface-raised)` on panel/card, `transparent` on section | `.inst-panel`, `.inst-card`, `.inst-section` |
| `--region-edge` | `var(--border)` on panel/card, `transparent` on section | the same three |

The section's own defaults differ because a section **draws nothing** —
`docs/layout/index.md`: *"A primitive does not draw."* Rank on a section moves
its title's rung and ink and nothing else. This is deliberate and is what keeps
`.inst-section` a layout primitive.

### 2.5 Rank against depth — how they stay orthogonal

Ground is composed, never chosen:

```
ground(region) = ground(container)  ∘  depthFilm?  ∘  rankFilm?

  depthFilm  applied when the region sits inside another region.
             Direction: recede. Always. Set by the containment rule, not by rank.
  rankFilm   applied only when rank = support.
             Direction: recede. Always. Bounded at one step.
```

**MEASURED CORRECTION — one film per region, and the stacking is in the paint.**

An earlier draft said a support card inside a panel "gets both films". Two
things were tested and both refute it.

*The two rules collide.* The containment rule
`:is(.inst-panel, …) :is(.inst-panel, .inst-card)` is `(0,2,0)` and
`[data-rank="support"]` is `(0,1,0)`. They write the same variable, so the
containment rule wins outright and the rank film is **silently lost**. Measured:
a support card and an unranked card inside the same panel resolved to the same
ground, `rgb(239,239,239)`, in both cases. No text gate can see this — it is a
cascade fact, not a token fact.

*Stacking does not need two values anyway.* Films composite in **paint**, down
the element tree. Three nested 6% blacks on white measured
`232.2 → 218.2 → 205.2`, ratios of 1.223 and 1.141 — each level genuinely
deeper. So nesting depth already composes for free, one film per region.

The rule is therefore:

> **A region receives at most ONE recess film. Depth requests it; rank requests
> it; a region that receives both requests still paints one. Depth of nesting is
> carried by the stack of elements, not by the strength of a single film.**

Both rules assign the identical value (`--surface-recessed`), so the collision is
harmless — but the two must never be allowed to want *different* values, because
the cascade will silently pick one. That constraint is recorded here and belongs
in the comment beside the containment rule.

Consequences to hold on to:

- `lead` and `default` have **identical grounds and identical edges.** Lead is
  carried by one type rung, and by nothing else. The kit's own third law makes
  this correct rather than thin: size is the only carrier of large-scale
  hierarchy, because weight 700 is banned and 600 is already spent.
- **Rank's ground lever bites only at depth 0.** A support panel sitting on the
  page recedes; a support card already inside a panel has spent the film on
  depth, and its rank shows in its title rung, its ink and its missing edge.
  This is honest and it is the price of the one-film rule. It is also
  symmetrical with lead, which never touches ground at all.
- A `lead` region nested inside another region still recedes, because depth owns
  direction. Its rank shows in its title.

### 2.6 Comparison against the variables already in the file

| Existing | Verdict |
|---|---|
| `--surface-raised` | unchanged. Becomes the panel's *fallback* rather than a literal in the rule. |
| `--surface-recessed` | unchanged, and now read by two more callers. No new film token. |
| `--surface-sunken` | **loses one caller.** `.inst-panel .inst-card` used it as an absolute step; the depth film replaces it. Every other caller (table head, log, code, aside) keeps it — those are grounds, not nested regions. |
| `--card-surface` | **absorbed into `--region-ground`.** It is already this variable under a component-scoped name, declared at `surfaces.css:123` and assigned at `:156`. One declaration, one consumer, one comment block. |
| `--panel-body-pad` | unchanged. Padding is not a rank concern. |
| `--flow-self` | unchanged, and **not** the precedent for this channel — §2.3 measured why the analogy fails. Left alone entirely. |
| `--tone-ink` / `--tone-mark` / `--tone-bg` | unchanged in P1–P2. `--tone-bg` reaching a region ground is the P1 follow-on in §10, and it is a generalisation of `.inst-note`, not a new mechanism. |

---

## 3. Markup contract

### 3.1 Vocabulary

A fifth entry in the table in `docs/about/design-principles.md`, alongside
`data-tone`, `data-state`, `data-kind`, `data-when`:

**`data-rank` — how much attention this place asks for. One vocabulary for the
whole library, closed.**

`lead` · `default` · `support`

### 3.2 `default` is explicit

`default` is a real, spellable, documented value, on the precedent already
stated for `data-state` in `design-principles.md:426`:

> *The base value is written in the markup though there are no rules for it: a
> `data-state="queued"` is read and its absence is not, and a typo like `quued`
> looks like the base.*

It has no CSS block (§2.3). It is made legal by adding it to the base-value list
in `docscheck` — see §5.5 — exactly as `queued`, `todo` and `approved` are.

**Where it is written.** On every region participating in a rank decision on
that screen. That means: the three `blocks/` pages write `data-rank` on every
region, because they demonstrate composition. Component reference pages, which
show one component in isolation, write nothing — a panel standing alone has no
screen to have a rank within, and stamping `data-rank="default"` on it would
teach that rank is a property of the panel, which is the error the model exists
to prevent.

The rule is enforceable in the DOM, not in text: see check **C1**/**C4** in §6.

### 3.3 Where rank is legal

Closed list of three, phase one:

```
.inst-panel      a region of the application
.inst-card       a bounded object in a region
.inst-section    a named group inside a screen
```

**Illegal everywhere else**, and specifically:

| Not legal on | Because |
|---|---|
| `.inst-btn`, `.inst-badge`, `.inst-tag`, `.inst-chip`, inputs, `.inst-icon` | these have emphasis, not rank. The button already answers "how loudly does this ask to be pressed" with four weights. Two axes for one question is one axis with two names. |
| `.inst-shell`, `.inst-container`, `.inst-split`, `.inst-split-side`, `.inst-split-main`, `.inst-stack`, `.inst-cluster`, `.inst-grid` | *a primitive does not draw.* A split half carrying a ground would stop being nestable inside something that paints otherwise. The panel inside it takes the rank. |
| `.inst-task`, `.inst-tree-item`, table rows and cells | a row is data. Law 3. |
| `.inst-metric` | deliberately deferred. A tile has no title in the region sense, and `design-principles` already carries a specific law about metric prominence. Speculative today; one line to add later. |

Rank on an arbitrary descendant is not *forbidden by CSS* — it simply does
nothing, because only the three components read the channel. That is the same
guarantee every closed vocabulary in this kit already gives, and it is why the
registry check in §5.6 exists: a silent no-op is the failure mode, so the list
is guarded from both sides.

### 3.4 Interaction with the existing vocabularies

| Combination | Meaning | Legal |
|---|---|---|
| `data-rank` + `data-tone` | the principal region of the screen, in a state | **yes.** `<div class="inst-panel" data-rank="lead" data-tone="error">` is a firing incident that dominates the screen. This is the point of separating the axes. |
| `data-rank` + `data-state` | a region whose lifecycle phase is shown | yes, where the component has a state vocabulary. None of the three do today. |
| `data-rank` + `data-kind` | orthogonal, no interaction | yes |
| `data-rank` + `data-density` | orthogonal. Density is an environment and inherits; rank is a local role and does not. | yes |

`data-tone="danger"` is **not** a value. The tone vocabulary is
`neutral · running · ok · warn · error`; `--danger` exists only as
`.inst-btn--danger`. A destructive region takes `data-tone="error"`.

### 3.5 Headings stay headings

No `<div class="inst-page-title">`. `docs/layout/section.md:40` already makes
the real tag a contract requirement — *"`inst-section-title` styles but does not
declare: navigation by headings does not work without it"* — and that stands.

The form is, and remains:

```html
<h2 class="inst-section-title">Billing history</h2>
```

What changes is only that the **bare-element font sizes** move out of the global
scope into `.inst-prose`. The `h1–h6` block in `base.css:43` keeps
`font-weight`, `line-height`, `letter-spacing`, `color` and `text-wrap` —
those are harmless globally and wanted. Only the four `font-size` lines move.

---

## 4. CSS implementation

Every change, with current and desired behaviour.

### 4.1 `src/tokens.css` — role tier

| | |
|---|---|
| **Where** | inside `:where(:root) { … }`, after `--radius-mark`, before `--control-ratio-*` |
| **Current** | no rank tokens exist |
| **Desired** | the five tokens of §2.2 |
| **How** | plain declarations; values are `var()` references into the type ladder and the surface semantics |
| **Tier legality** | tier 3 may read tier 2 (`--surface-recessed`) and tier 3 (`--text-*`). It does. No component sees tier 1. |
| **Density / scale** | not repeated in any cell. They resolve through ladders the cells already move. |

### 4.2 `src/surfaces.css` — the channel and the rank blocks

| | |
|---|---|
| **Where** | a new section at the top of the file, above `.inst-panel` |
| **Current** | nothing |
| **Desired** | one `:where()` reset; two `[data-rank]` blocks. **No `@property`** — see §2.3 |

```css
/* The reset. Zero specificity, so every [data-rank] block outranks it, and a
   nested region stops its container's rank at its own boundary. */
:where(.inst-panel, .inst-card) {
  --region-title-size: var(--region-title-default);
  --region-title-ink:  var(--text-primary);
  --region-ground:     var(--surface-raised);
  --region-edge:       var(--border);
}
:where(.inst-section) {
  --region-title-size: var(--region-title-default);
  --region-title-ink:  var(--text-muted);
  --region-ground:     transparent;
  --region-edge:       transparent;
}

[data-rank="lead"] {
  --region-title-size: var(--region-title-lead);
}
[data-rank="support"] {
  --region-title-size: var(--region-title-support);
  --region-title-ink:  var(--text-secondary);
  --region-ground:     var(--region-ground-support);
  --region-edge:       var(--region-edge-support);
}
```

Two rank blocks, not three. `lead` sets one property because it differs from
default in one property (§2.5). Bare attribute selectors rather than
`.inst-panel[data-rank]`: the channel is inert on anything that does not read
it, and a per-component selector would be three copies of one rule.

The `:where(.inst-section)` reset is separate only because a section's defaults
differ (§2.4). Both resets are `(0,0,0)`.

> **Cascade note, measured.** The rank blocks and the resets both live in
> `kit.components`; within a layer, specificity decides, so `(0,1,0)` beats
> `(0,0,0)` regardless of source order. `.inst-section` is declared in
> `kit.layout`, an **earlier** layer, so a reset written here would beat any
> `.inst-section` declaration there whatever its specificity — which is why the
> section's own title rule (§4.7) must read the channel rather than fight it.
> Do not reach for a descendant selector anywhere in this section; that is the
> failure `tone.css` records for `--icon-size` and the rail.

> **Cascade note.** These blocks live in `kit.components`. `.inst-section` is in
> `kit.layout`, an **earlier** layer — but a later layer beats an earlier one at
> any specificity, so the declarations still reach it, because what crosses the
> layer boundary is an inherited-or-fallback *variable resolved on the element
> itself*, not a selector. This is the mechanism `tone.css` documents for
> `--icon-size` reaching the rail. Do not attempt a descendant selector here;
> that is the failure that comment records.

### 4.3 `src/surfaces.css:19` — `.inst-panel`

| | |
|---|---|
| **Current** | `background: var(--surface-raised); border: var(--hairline) solid var(--border);` |
| **Desired** | both read the channel with today's value as the fallback |
| **How** | `background: var(--region-ground, var(--surface-raised)); border: var(--hairline) solid var(--region-edge, var(--border));` |
| **Why legal** | fallback-at-consumption is the kit's established shape (`--tone-bg`, `--icon-size`, `--panel-body-pad`). A panel with no rank renders byte-identically to today. |

### 4.4 `src/surfaces.css:49` — `.inst-panel-title`

| | |
|---|---|
| **Current** | `font-size: var(--text-sm);` and inherited colour |
| **Desired** | `font-size: var(--region-title-size, var(--region-title-default)); color: var(--region-title-ink, var(--text-primary));` |
| **Note** | `font-weight: var(--weight-medium)` and `line-height: var(--leading-ui)` are unchanged. Weight is binary and is not a rank carrier. The long comment at `:49` explaining why a panel name gets no step up must be **amended, not deleted**: it is still true of `default`, and rank is the mechanism that finally makes the exception expressible. |

### 4.5 `src/surfaces.css:112,123` — `.inst-card`

| | |
|---|---|
| **Current** | `background: var(--card-surface, var(--surface-raised));` |
| **Desired** | `background: var(--region-ground, var(--surface-raised));` — `--card-surface` is retired, being this variable under a narrower name |
| **Also** | `border` reads `--region-edge` as in 4.3; `.inst-card-title` reads the title channel as in 4.4 |

### 4.6 `src/surfaces.css:156` — the containment rule

| | |
|---|---|
| **Current** | `.inst-panel .inst-card { --card-surface: var(--surface-sunken); }` — one hard-coded pair, an absolute step |
| **Desired** | a general rule: any region inside a region recedes by one film |
| **How** | `:is(.inst-panel, .inst-card, .inst-section) :is(.inst-panel, .inst-card) { --region-ground: var(--surface-recessed); }` |
| **Why this shape** | the *outer* list includes `.inst-section` because a section groups regions; the *inner* list excludes it because a section draws nothing and has no ground to recede. The film composes with the rank film automatically — a support card inside a panel gets both. |
| **Why it does not violate the token tiers** | it assigns a tier-2 token to a channel variable. No component gains a literal. |
| **What it fixes** | audit finding 02, measured: card-in-card and panel-in-panel both sat at exactly 0.000 lightness step. |
| **Risk** | this is the widest-blast-radius change in the plan. It ships **alone** — see §10, P2a. |

### 4.7 `src/layout.css:402,432` — `.inst-section`

| | |
|---|---|
| **Current** | `.inst-section-title { font-size: var(--text-xs); font-weight: var(--weight-normal); color: var(--text-muted); }` |
| **Desired** | `font-size: var(--region-title-size, var(--region-title-default)); color: var(--region-title-ink, var(--text-muted));` |
| **Behaviour change** | **yes, and it is the point.** A default section title moves from `--text-xs` (12.5) to `--text-sm` (14) and from `--text-muted` to… still `--text-muted`, because the fallback preserves it. Weight stays `normal`. |
| **Decision** | the size moves; the ink does not. Finding 03 of the audit is that a section heading has two settings and neither is right; one rung up at unchanged ink is the smallest correction that makes `default` a usable middle. `lead` then reaches 16px and `support` returns to 12.5. |
| **Consequence** | `docs/layout/section.md` must be corrected in the same wave — it currently claims `--text-sm` / `--weight-medium` / `--text-secondary` in three places and is wrong about all three today. After this change it becomes right about the size and still wrong about the other two. See §8. |

### 4.8 `src/base.css:51–54` — heading sizes

| | |
|---|---|
| **Current** | `h1{--text-2xl} h2{--text-xl} h3{--text-lg} h4,h5,h6{--text-md}` at global scope |
| **Desired** | the same four lines, nested inside `.inst-prose` |
| **Keep** | the `h1–h6` block at `:43` (weight, leading, tracking, colour, `text-wrap`) stays global |
| **Why** | audit finding 03: a bare `<h2>` is 22/600/primary while `.inst-section-title` on the same element is 12.5/400/muted, and nothing mediates. Confining the ladder to prose leaves exactly one document-scope answer, and it is the kit class. |
| **Risk** | any application heading without a kit class drops to the body size. That is the intended behaviour and it must be announced in `CHANGELOG.md` as breaking. |
| **Migration** | §9. |

---

## 5. Gates — first, and red before green

**The discipline.** A `kase` referencing a token that does not exist makes the
gate error, not fail a comparison. So "gates first" here means:

> **P0 writes the gate rows and the mutations and confirms they are RED against
> the current tree. P1 adds the tokens and confirms they go GREEN.**

A row that was never red proves nothing. `cmd/mutate` exists in this repository
precisely because that has happened before.

### 5.1 The measurable invariant

> **Lead is discernibly ahead of default, and support discernibly behind it, in
> every one of the five themes and every one of the five scales.**

It decomposes into two measurable halves, because rank has two carriers:

| Carrier | Quantity | Threshold | Gate |
|---|---|---|---|
| Title rung | ratio of resolved px | `stepMin` = 1.12, the existing type-ladder threshold | `cmd/proportion` |
| Support ground | OKLCH lightness delta against the container | `step` = 0.022, the existing surface-stack threshold | `cmd/contrast` |

Neither threshold is new. Both are already derived and defended in the tools.

### 5.2 `cmd/contrast` — five `kase` rows, no tool change

Append after the `SURFACE STACK STEPS` block:

```go
// ── RANK ────────────────────────────────────────────────────────────────
// Support recedes from default AT THE SAME DEPTH. Default's ground IS the
// container, so "support differs from default" is exactly "the support film
// clears a step against the container". alt is NOT the instrument here: alt
// compares two alternatives that both sit on a ground, and default has no
// film to compare.
{label: "rank: support region on page",  fg: "--region-ground-support", bg: []string{"--surface-page"},   min: step},
{label: "rank: support region on panel", fg: "--region-ground-support", bg: []string{"--surface-raised"}, min: step},

// The COMPOSED case, which no pair above asks: a support region nested one
// level deep carries two films, and the text on it still has to be read.
{label: "rank: text in a nested support region", fg: "--text-primary",
	bg: []string{"--surface-raised", "--surface-recessed", "--region-ground-support"}, min: text},
{label: "rank: support title in a nested support region", fg: "--text-secondary",
	bg: []string{"--surface-raised", "--surface-recessed", "--region-ground-support"}, min: text},

// A support title on its own ground, at depth 0.
{label: "rank: support title on page", fg: "--text-secondary",
	bg: []string{"--surface-page", "--region-ground-support"}, min: text},
```

Rows read `--region-ground-support` rather than `--surface-recessed` **on
purpose**: the two are the same value and different questions, and a mutation
of the rank token must turn only the rank rows red. This is the precedent
`--surface-field: var(--surface-raised)` already sets — *"two roles on one step
are not a duplicate."*

### 5.3 `cmd/proportion` — two `rule` rows, no tool change

```go
{label: "rank: lead title over default", a: "--region-title-lead", b: "--region-title-default",
	min: stepMin, max: 1.30,
	why: "a lead region's name stops differing from an ordinary one"},
{label: "rank: default title over support", a: "--region-title-default", b: "--region-title-support",
	min: stepMin, max: 1.30,
	why: "a support region's name stops receding"},
```

`perDens: false` — density does not move type size, by law. `baseOnly: false` —
the ladder must hold in all five scales.

**These are not duplicates of the `type: xs → sm` and `type: sm → md` rows.**
Those gate the *ladder*. These gate the *mapping from rank onto the ladder*. If
someone repoints `--region-title-lead` at `--text-sm`, the ladder rows stay
green and these go red. An implementer tempted to "simplify these away" should
read this paragraph twice.

**A margin warning that must be recorded in the comment.** Resolved across the
five scales:

| scale | sm/xs (default:support) | md/sm (lead:default) |
|---|---|---|
| base | **1.1200** | 1.1429 |
| 15 | 1.1538 | 1.1333 |
| 16 | 1.1429 | 1.1562 |
| 17 | 1.1333 | 1.1471 |
| 18 | 1.1250 | 1.1667 |

The base cell sits at **exactly 1.1200 against a threshold of 1.12**, with the
comparison being `got < min`. Zero margin. Any future nudge to `--text-xs` or
`--text-sm` breaks rank as well as the type ladder. That is a correct coupling —
the rungs *are* the ladder — but it must be written beside the rule so the next
person who touches the ladder knows what else they are touching.

### 5.4 `cmd/mutate` — four rows

```go
// ── rank ───────────────────────────────────────────────────────────────
{"rank: lead levelled with default", "proportion", "src/tokens.css",
	"--region-title-lead:    var(--text-md);",
	"--region-title-lead:    var(--text-sm);",
	"two ranks produce one look, and the ladder loses its top step"},
{"rank: support levelled with default", "proportion", "src/tokens.css",
	"--region-title-support: var(--text-xs);",
	"--region-title-support: var(--text-sm);",
	"de-emphasis stops existing and every region reads at one level"},
{"rank: the ladder inverted", "proportion", "src/tokens.css",
	"--region-title-lead:    var(--text-md);",
	"--region-title-lead:    var(--text-2xs);",
	"the principal region's name is smaller than an ordinary one"},
{"rank: support stopped receding", "contrast", "src/tokens.css",
	"--region-ground-support: var(--surface-recessed);",
	"--region-ground-support: transparent;",
	"a support region's ground levels with its container"},
```

Mutation 3 is caught because an inverted ratio (0.786) falls below `min`, not
because of a sign test — the band does the work, as it does for the control
ladder.

Mutation 4 turns **only** the five rank rows red, which is the proof that they
are wired rather than incidentally satisfied by an existing pair. It is not a
duplicate of *"button levelled with its backdrop"*: that one mutates
`--surface-recessed` and turns the button pairs red as well.

### 5.5 `cmd/docscheck` — ~6 lines

Current, at `main.go:199`:

```go
for _, v := range []string{"queued", "todo", "approved"} {
    if vocab["state"] == nil { vocab["state"] = map[string]bool{} }
    vocab["state"][v] = true
}
```

Generalise to a map so `rank` can carry a base value the same way:

```go
// Base values have no styling, so they are absent from selectors. They are
// declared by the constitution — without them the check would complain about
// valid markup.
for attr, vals := range map[string][]string{
    "state": {"queued", "todo", "approved"},
    "rank":  {"default"},
} { … }
```

This is what makes `data-rank="default"` legal in documentation while
`data-rank="lede"` fails the build. Nothing else in `docscheck` changes:
`lead` and `support` are learned from the selectors in §4.2 automatically.

### 5.6 `cmd/registry` — the region list, guarded both ways

`components.json` gains one field on three entries:

```json
"panel":   { "region": true, … }
"card":    { "region": true, … }
"section": { "region": true, … }
```

`registry` cross-checks it against the `:is(…)` lists in §4.2/§4.6, in both
directions — a component marked `region` whose class is absent from the CSS
list is an error, and a class in the list without the flag is an error too. This
is the same shape as the forced / print / flow lists the tool already guards,
and it exists for the same stated reason: *"a forgotten entry causes no browser
error."*

It is **not** a new axis in the `vocab` map. Rank is a per-instance decision;
`region` is a per-component fact, which is exactly what the registry is for.

### 5.7 `tools/audit.js` + `audit-run.mjs` — the rendered gates

Text gates cannot see inheritance. §6 defines the section; the mutation that
proves it is live:

```js
{
  name: 'a region stopped resetting its container’s rank',
  section: 'composition',
  file: 'surfaces.css',
  page: '/blocks/console/',
  from: ':where(.inst-panel, .inst-card) {',
  to:   ':where(.inst-panel) {',
  why:  'a card inside a lead panel becomes lead, which is the defect the model exists to prevent',
}
```

This is the single most important mutation in the plan. It is the only one that
can fail if law 4 is implemented wrongly, and no Go gate can express it.

It replaces an earlier `inherits: false → inherits: true` mutation, which was
written against a mechanism §2.3 measured and discarded. Removing `.inst-card`
from the reset list is the exact equivalent under the mechanism that actually
works: the card stops blocking its container's rank, and check **C3** must go
red.

---

## 6. Orthogonality tests — `kitAudit.composition()`

A fourth section in `tools/audit.js`, beside `contrast`, `targets` and
`proportion`. It measures the rendered DOM, which is the only place these
questions have answers.

| # | Check | Fails when |
|---|---|---|
| **C1** | at most one `[data-rank="lead"]` per document | a screen has two leads, which is no lead |
| **C2** | every `[data-rank]` value ∈ {lead, default, support} | a typo silently does nothing — the failure mode of every closed vocabulary |
| **C3** | **the nesting test.** For every region inside `[data-rank="lead"]` that carries no `data-rank` of its own: its title's computed `font-size` equals the resolved `--region-title-default`, not `--region-title-lead` | law 4 is broken — a region stopped resetting the channel |
| **C4** | if a document contains any `[data-rank]`, every sibling region of a ranked region also carries `[data-rank]` | half a screen has declared its composition and half has not |
| **C5** | for every region inside a region, the composited ground is strictly darker (light themes) / lighter (dark themes) than its container's, by ≥ the step | depth stopped owning direction |

The named test cases, to be built as one page under the existing preview
mechanism:

| Case | Expectation |
|---|---|
| default panel on page | ground `--surface-raised`, edge `--border`, title `--text-sm` — **byte-identical to today** |
| lead panel on page | ground and edge unchanged; title `--text-md` |
| support panel on page | ground = page + one film; edge transparent; title `--text-xs`, secondary ink |
| card in default panel | ground = raised + one depth film; title `--text-sm` |
| **card in lead panel** | ground = raised + one depth film; title **`--text-sm`** — *not* `--text-md`. **This is the critical case.** |
| support card in lead panel | ground = raised + **one** film (depth already spent it — §2.5); edge transparent; title `--text-xs`, secondary ink |
| lead card in support panel | ground = page + one film (the panel's rank) then one film (the card's depth) — two films, because they are on two elements; title `--text-md`. Rank rises where depth has fallen: the two axes visibly not touching. |
| panel in panel | ground = raised + one depth film. Today: exactly 0.000 step. |
| card in card in card | three films. Today: 0.000 at both nestings. |

`kitAudit` is documented as a tool consumers run on **their own** screens, so
this section ships as consumer-facing capability, not only as CI.

---

## 7. The theme matrix — which dimensions are load-bearing

Not all 300 cells carry the same risk, and saying which do is part of the plan.

| Axis | Count | Load-bearing for rank? | Why |
|---|---|---|---|
| **Theme** | 5 | **Yes, mandatory** | the ground film's sign flips by scheme; `--surface-raised` is the ramp ceiling in the three light themes and mid-ramp in `dark-soft`. Law 2's "no lift" rule exists *because* of this axis. `dark-soft` is the one theme where a lift would have worked, which is why the ban must be a rule rather than a per-theme discovery. `cmd/contrast` already iterates all five. |
| **Scale** | 5 | **Yes, mandatory** | the title rungs are three positions on a ladder that is hand-written per scale. The margins in §5.3 differ per scale and the base cell has none. `cmd/proportion` already iterates all five. |
| **Accent** | 4 | swept, no rank-specific risk | no rank token reads `--a-*`. Kept in the sweep because it costs nothing and excluding it would be a special case. |
| **Density** | 3 | swept, no rank-specific risk | density does not move type size, by law, so the title ladder is density-invariant. The ground film is an alpha and does not move either. Hence `perDens: false` on the proportion rows. |

The rendered sweep (`audit-run.mjs`) runs the full matrix regardless; the table
above says where to look first when it goes red.

---

## 8. Documentation

One canonical contract, referenced from everywhere else. No rule stated twice.

| File | Change | Canonical? |
|---|---|---|
| `docs/about/design-principles.md` | **the canonical contract.** A fifth entry in "State and variant": the `data-rank` vocabulary, the four laws of §0, the legality list of §3.3, and the interaction table of §3.4. Plus: bare `h1–h6` sizing now belongs to `.inst-prose`. | **yes** |
| `docs/foundations/elevation.md` | the Scale table gains the containment rule: a step is chosen by the role of the region **and by its depth**, and depth is relative. The claim *"a step is chosen by the role of the region rather than by a wanted brightness"* becomes true rather than aspirational. | for depth |
| `docs/components/display/panel.md` | a `variants` section: what rank does to a panel. Links to the contract; does not restate the laws. | no |
| `docs/components/display/card.md` | same, plus the containment behaviour. The bare `<h3>` in its preview is migrated (§9). | no |
| `docs/layout/section.md` | **corrected and extended.** Three places currently contradict `layout.css:432` — the `api` block, the a11y "Type size" row, and the anatomy sentence at `:93`. All three become `--text-sm` / `--weight-normal` / `--text-muted` after §4.7, so two of the three become right by the code change and the third must be edited. | no |
| `docs/layout/index.md` | one row in "Common to all eight": a layout primitive does not take rank, and why. | no |
| `docs/foundations/typography.md` | states where the bare-element ladder now lives. Its own demo is migrated (§9). | no |
| `docs/blocks/*.md` × 3 | rank added to the markup; a paragraph in "What is seen here" naming the lead region and why. | no |
| `CHANGELOG.md` | breaking: bare `h1–h6` sizing is now `.inst-prose`-only. Additive: `data-rank`. | — |
| `docs/internal/DECISIONS.md` | the derivations that do not belong on a public page: why lead does not lift; why the ink ladder gets no rank token; why `default` has no CSS block. | for the reasoning |

Not written: a `RANK.md` of its own. This kit already has a canonical home for
markup vocabularies and it is `design-principles.md`; a second document would be
the parallel abstraction §1 refuses. This plan file is internal and is deleted
or superseded when the work lands.

---

## 9. Migration

Searched at `d39cba0`.

### Zero conflicts

`data-rank`, `--region-*` and `inst-region` appear **nowhere** in `src/`,
`docs/`, `site/` or `components.json`. The namespace is clean.

### Required before the rule can be enabled

| # | Item | Count | Wave |
|---|---|---|---|
| 1 | `--card-surface` → `--region-ground` | 1 declaration (`surfaces.css:123`), 1 assignment (`:156`), 1 comment block (`:129–155`) | P2a |
| 2 | `.inst-panel .inst-card` retired for the general containment rule | 1 selector | P2a |
| 3 | Bare `<h1>–<h6>` inside preview markup | **42** occurrences, in 5 English pages + their `.ru` pairs: `components/display/card.md`, `foundations/typography.md`, `layout/index.md`, `layout/page-header.md`, `layout/section.md` | P3 |
| 4 | `docs/layout/section.md` three-way contradiction | 3 lines × 2 languages | P4 |
| 5 | Rank added to the three `blocks/` pages | 3 pages × 2 languages | P3 |

### Not required

- **The 274 region markup sites in the documentation** — 104 `inst-panel`, 78
  `inst-card`, 20 `inst-section`, 72 `inst-metric` — need **no edit**. Absence
  of `data-rank` resolves to `default`, and `default` is today's rendering
  exactly. This is the property the whole migration rests on.
- Any consumer application markup. Same reason. The only breaking change in the
  plan is item 3's global effect, and it is announced.

### The one real risk, item 3

`docs/foundations/typography.md` demonstrates the type ladder with bare
headings. If previews render outside `.inst-prose` — and the `.inst-theme`
comment at `tokens.css:906` says the site had exactly this leakage problem
before — that demo loses its sizes when §4.8 lands.

**Verify before P3:** render `/foundations/typography/` and read the computed
`font-size` of the demo's `<h1>`. If the preview stage is not prose, the demo
gains a `.inst-prose` wrapper, and that is the correct fix rather than a
carve-out: the ladder being demonstrated *is* the prose ladder after §4.8.

---

## 10. Implementation order

Each phase names its files, its purpose, its acceptance criterion, and what it
depends on. No phase begins before its predecessor's criterion is met.

### P0 — Gates, red

- **Modify:** `tools/cmd/contrast/main.go` (5 rows), `tools/cmd/proportion/main.go`
  (2 rows), `tools/cmd/mutate/main.go` (4 rows), `tools/cmd/docscheck/main.go`
  (~6 lines).
- **Create:** nothing.
- **Purpose:** state the invariant in machine-readable form before anything can
  satisfy it.
- **Acceptance:** `check contrast` and `check proportion` are **RED**, and red
  for the stated reason: the five and two new rows fail on an unresolvable
  token. Every pre-existing row still passes. `check docscheck` is green.
- **Depends on:** nothing.

> This phase deliberately leaves the tree failing. It is one commit with P1, or
> it is a branch. Do not push it to `master` alone.

### P1 — Token and property infrastructure

- **Modify:** `src/tokens.css` (5 role tokens), `src/surfaces.css` (2 `:where()`
  resets, 2 `[data-rank]` blocks).
- **Create:** nothing.
- **Purpose:** make the invariant satisfiable and the channel exist.
- **Acceptance:** `check` is fully green — all eight gates. `check mutate`
  reports the four new mutations **caught**. Critically: **nothing renders
  differently yet**, because no component reads the channel. Verify by
  `check pixels` against the pre-change baseline: zero diffs.
- **Depends on:** P0.

### P2a — Containment, alone

- **Modify:** `src/surfaces.css` (`.inst-panel`, `.inst-card`,
  `.inst-panel .inst-card` → the general rule; `--card-surface` retired).
- **Purpose:** depth becomes relative. Audit finding 02 closes.
- **Acceptance:** `check` green; `check pixels` green; and the C5 cases of §6
  measured by hand on a scratch page — card-in-card and panel-in-panel show a
  step ≥ 0.022 where they showed 0.000.
- **Depends on:** P1.
- **Ships alone.** This is the widest visual blast radius in the plan, and a
  pixel sweep that mixes it with rank is unreadable.

### P2b — Rank on the panel

- **Modify:** `src/surfaces.css` (`.inst-panel-title` reads the channel).
- **Purpose:** the smallest possible proof that the model works.
- **Acceptance:** the three panel cases of §6 rendered side by side read as
  *"clearly important / ordinary / quiet"* and **not** as *"big panel / normal
  panel / tiny panel"*. If they read as the latter, the rung choices are wrong
  and this is where that is discovered — before card, section and the blocks
  are built on them.
- **Depends on:** P2a.

### P2c — Rank on the card, and the orthogonality proof

- **Modify:** `src/surfaces.css` (`.inst-card`, `.inst-card-title`).
- **Create:** `tools/audit.js` gains `composition()`; `tools/audit-run.mjs`
  gains the inheritance mutation.
- **Purpose:** prove law 4 by measurement.
- **Acceptance:** C1–C5 pass on the test page; `node tools/audit-run.mjs
  --mutate` reports *"a region stopped resetting its container's rank"*
  **caught**. If the mutation is missed, the `:where()` reset is not doing what
  it claims and P3 does not begin.
- **Depends on:** P2b.

### P2d — Rank on the section, and the heading move

- **Modify:** `src/layout.css` (`.inst-section-title`), `src/base.css`
  (four `font-size` lines into `.inst-prose`).
- **Purpose:** audit finding 03 closes; one document-scope answer for headings.
- **Acceptance:** `check` green including `docscheck`; the site builds; a bare
  `<h2>` outside prose resolves to the body size; `<h2 class="inst-section-title">`
  resolves to `--text-sm`.
- **Depends on:** P2c. **Blocks:** migration item 3.

### P3 — Markup migration

- **Modify:** `docs/blocks/{console,monitor,chat}.md` + `.ru` — rank added to
  every region. The five pages of migration item 3 + `.ru` — bare headings.
- **Purpose:** the real test of the model. Not new demos: the existing three.
- **Acceptance:** each block's markup gains `data-rank` and **loses** at least
  as much as it gains — no local overrides, no wrapper divs added to force a
  hierarchy. If a block does not get simpler, the model is wrong and P4 does not
  begin. `check site` green; `check pixels /blocks/` green.
- **Depends on:** P2d.

### P4 — Documentation

- **Modify:** the nine files of §8.
- **Purpose:** move the grammar out of the CSS comments, once.
- **Acceptance:** `check docscheck` green with `data-rank` markup present in
  pages — this is the first phase where the §5.5 change is actually exercised.
  `check lang` green. `check site` green.
- **Depends on:** P3.

### P5 — Full matrix validation

- **Modify:** nothing.
- **Purpose:** confirm across all five themes and all five scales.
- **Acceptance:** `check` green; `check pixels` green across the full matrix;
  `check pixels --mutate` all caught; `check behavior` green; `check mutate` all
  caught. Then delete or supersede this file.
- **Depends on:** P4.

### Deliberately after P5, not inside it

**Tone on a region ground.** `data-tone` reaching `--region-ground` and
`--region-edge`, generalising what `.inst-note` and `.inst-banner` already do at
`feedback.css:101,116`. It closes audit finding 04 (a danger zone, a firing
incident, a failed run as a *region*). It is not in P0–P5 because three changes
in one pixel sweep is one change nobody can read, and because the composed
contrast cases it needs — a tone ground inside a ranked region inside a depth
film — are a larger gate addition than everything in §5 combined.

---

## Implementation Contract

Everything an implementer needs, with no decisions left open.

### Vocabulary

```
data-rank = "lead" | "default" | "support"
```

- Base value: `default`. Spellable in markup; no CSS block; made legal by the
  `docscheck` base-value list.
- Legal on exactly: `.inst-panel`, `.inst-card`, `.inst-section`.
- Illegal (inert) everywhere else. Guarded by `cmd/registry` against
  `components.json`'s `"region": true`.
- Orthogonal to `data-tone`, `data-state`, `data-kind`, `data-when`,
  `data-density`, `data-scale`, `data-theme`, `data-accent`.
- `data-tone="danger"` does not exist. Use `error`.

### Role tokens — `tokens.css`, `:where(:root)`

```css
--region-title-lead:     var(--text-md);
--region-title-default:  var(--text-sm);
--region-title-support:  var(--text-xs);
--region-ground-support: var(--surface-recessed);
--region-edge-support:   transparent;
```

No `--region-ground-lead`, no `--region-edge-lead`, no ink token. No entry in
any of the fifteen scale × density cells.

### Channel — `surfaces.css`, ordinary inheritance, reset at every region

**No `@property`.** Measured: `inherits: false` stops the value one element
short of the title, which is a grandchild of the ranked region. The mechanism is
a `:where()` reset on every region class, at zero specificity.

| Property | Reset value (panel, card) | Reset value (section) |
|---|---|---|
| `--region-title-size` | `var(--region-title-default)` | `var(--region-title-default)` |
| `--region-title-ink` | `var(--text-primary)` | `var(--text-muted)` |
| `--region-ground` | `var(--surface-raised)` | `transparent` |
| `--region-edge` | `var(--border)` | `transparent` |

The same values are also written as fallbacks at the point of consumption, for a
title used outside any region.

### Resolution table

| | title rung | title ink | ground | edge |
|---|---|---|---|---|
| `lead` | `--text-md` | primary | *unchanged* | *unchanged* |
| `default` | `--text-sm` | primary | *unchanged* | *unchanged* |
| `support` | `--text-xs` | secondary | requests one recess film | transparent |
| *containment, any rank* | — | — | requests one recess film | — |

**One film per region.** Depth and rank both request the same film and a region
paints at most one. Depth of nesting shows through the stack of painted
elements, measured at `232.2 → 218.2 → 205.2` on white. The two rules must never
be allowed to want different values: the containment rule at `(0,2,0)` beats
`[data-rank]` at `(0,1,0)` and the loser is dropped silently.

### Invariants

| ID | Statement | Measured by | Threshold |
|---|---|---|---|
| **I1** | lead title rung / default title rung ≥ 1.12, in all 5 scales | `cmd/proportion` | `stepMin` |
| **I2** | default title rung / support title rung ≥ 1.12, in all 5 scales | `cmd/proportion` | `stepMin` |
| **I3** | support ground differs from its container by ≥ 0.022 OKLCH lightness, in all 5 themes × 4 accents | `cmd/contrast` | `step` |
| **I4** | text on a doubly-filmed support region clears 4.5:1, in all 5 themes | `cmd/contrast` | `text` |
| **I5** | rank does not cross a region boundary: a region inside a lead region, with no rank, resolves `default` | `kitAudit.composition` C3 | equality |
| **I6** | at most one `lead` per document | `kitAudit.composition` C1 | count ≤ 1 |
| **I7** | every region inside a region recedes, never rises | `kitAudit.composition` C5 | `step`, signed |
| **I8** | a region with no `data-rank` renders byte-identically to `d39cba0` | `check pixels` baseline diff | zero |

### Mutations that must be caught

| Mutation | Gate | Proves |
|---|---|---|
| `--region-title-lead` → `--text-sm` | proportion | I1 is live |
| `--region-title-support` → `--text-sm` | proportion | I2 is live |
| `--region-title-lead` → `--text-2xs` | proportion | I1 catches inversion |
| `--region-ground-support` → `transparent` | contrast | I3 is live, and only the rank rows |
| `.inst-card` removed from the `:where()` reset list | `audit-run.mjs` composition | I5 is live |

### Test cases

The nine cases in §6, of which one is critical:

> **A `.inst-card` with no `data-rank`, inside a `.inst-panel[data-rank="lead"]`,
> must resolve its title to `--region-title-default`.**
>
> If it resolves to `--region-title-lead`, the `:where()` reset is not in effect
> and nothing downstream of P2c may proceed.

### The line that decides an ambiguous case

When an implementer is unsure whether something belongs to rank or to depth:

> **Does it change with the screen, or with the nesting?** With the screen —
> rank. With the nesting — depth. If it changes with neither, it is a component
> property and belongs in neither.
