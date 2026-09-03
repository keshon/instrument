---
title: Rank
group: Foundations
layout: foundation
source: src/tokens.css
api:
  - { name: "data-rank", kind: "attribute", doc: "`lead` `default` `support`. Set on a region; the absence of the attribute resolves to `default`" }
  - { name: "--region-title-lead", kind: "token" }
  - { name: "--region-title-default", kind: "token" }
  - { name: "--region-title-support", kind: "token" }
  - { name: "--region-ground-support", kind: "token" }
  - { name: "--region-edge-support", kind: "token" }
  - { name: "--region-title-size", kind: "token" }
  - { name: "--region-title-ink", kind: "token" }
  - { name: "--region-ground", kind: "token" }
  - { name: "--region-edge", kind: "token" }
---

Rank says **how much attention a place asks for**. It is set on a region and
answers a question about the SCREEN rather than about the component: the same
panel is the point of one screen and a footnote on another.

```html preview
<div class="inst-stack">
  <div class="inst-panel" data-rank="lead">
    <div class="inst-panel-header"><span class="inst-panel-title">Current run</span></div>
    <div class="inst-panel-body">Three agents working, one waiting on approval.</div>
  </div>
  <div class="inst-panel" data-rank="default">
    <div class="inst-panel-header"><span class="inst-panel-title">Recent runs</span></div>
    <div class="inst-panel-body">Four finished in the last hour.</div>
  </div>
  <div class="inst-panel" data-rank="support">
    <div class="inst-panel-header"><span class="inst-panel-title">Retention</span></div>
    <div class="inst-panel-body">Transcripts are kept for thirty days.</div>
  </div>
</div>
```

The library had this axis for a button and for nothing else. A button has four
weights, each one answer to *how loudly does this ask to be pressed*; a panel had
one appearance and no question, so two panels of wholly different importance came
out with identical computed styles and no markup could tell them apart.

The laws rank is held to — and the reason it is an attribute rather than a
modifier — are in
[design principles](../about/design-principles.md#state-and-variant). This page
is the material: the values, what each resolves to, and how to choose.

## Contract

| What | Required | Why |
|---|---|---|
| A value from the vocabulary | yes | `lead` · `default` · `support`. A value outside it does nothing at all, which is the failure mode of every closed vocabulary here |
| On a region | yes | A panel, a card or a section. Elsewhere it is inert — nothing else reads the channel |
| At most one `lead` per screen | yes | Two leads is no lead. Checked on the rendered tree by `kitAudit.composition()` |
| `data-rank="default"` written out | no, but wanted | The absence resolves to it. Written, a typo looks like a mistake instead of looking like the base |
| Rank on every region of a composition | yes, once any is ranked | Half a screen cannot have declared its hierarchy while the other half has not |

### Accessibility

| | |
|---|---|
| Rank is not a state | It carries no meaning and is announced by nothing. What a region MEANS is [tone](./colors.md), and a status always travels with a word |
| Contrast at every rank | Every rung is measured against every ground it can stand on, across five themes and four accents. `support` takes `--text-muted`, the quietest tier text may be read in — `--text-faint` is the decoration threshold and is forbidden for reading |
| Type floor | `support` takes `--region-title-support`, which is `--text-2xs` — the declared bottom of the scale, and the one place a badge-sized label is right, since a support region is a footnote by definition. Nothing goes below it |
| Heading level | Unchanged by rank. `<h2 class="inst-section-title">` is a real `h2` whatever its rank: the tag carries the outline, rank carries the size |

## Scale

Three values, and the list is closed.

| Value | The name | Ink | Ground | Edge | Inset |
|---|---|---|---|---|---|
| `lead` | `--text-sm` | unchanged | unchanged | unchanged | unchanged |
| `default` | `--text-xs` | unchanged | unchanged | unchanged | unchanged |
| `support` | `--text-2xs` | `--text-muted` | one recess film | none | one step tighter |

A quiet region is physically smaller, not merely labelled in smaller type. The
inset costs no token of its own: `--pad-block-y` is exactly one step under
`--pad-panel` in all fifteen scale and density cells, and `--row-pad-y` one step
under that, so `support` steps down by pointing at roles that already answer to
density. A panel's header, body and toolbar move together — they share one
inline inset, and a tightened one is still one.

**The whole ladder sits below the base,** and nothing grows to signal
importance. A lead region's name reaches `--text-sm`, the size of ordinary
interface text; an ordinary region's name sits a rung under that, and a quiet
one two rungs under.

So an ordinary label is smaller than the data it labels, which the third law
always meant and the kit used to leave to weight alone. A **lead** label is the
one exception: it matches its data's size and is told from it by weight, the way
every label once was. There is one lead per screen, and it is where the eye is
meant to stop.

Hierarchy is made by the surroundings getting quieter rather than by the
important thing getting louder. That also keeps chrome clear of content:
`--text-md` belongs to prose and to an agent's answer, and no region name
reaches it at any rank.

**Lead moves one thing** — its name, one rung. There is nowhere to lift a region
to, since `--surface-raised` is the top of the ramp in every light theme, so a
lead that lifted would find headroom in two themes of five and silently do
nothing in three.

**Support may spend more,** because withdrawing is the cheap direction: a rung
down, an ink tier down, a film and no edge at all.

### The tokens behind it

The rungs are references into the type ladder rather than numbers, so
[scale](./scale.md) carries them along with everything else it moves.
[Density](./density.md) does not touch them: it does not touch type size.

| Token | Value |
|---|---|
| `--region-title-lead` | `--text-sm` |
| `--region-title-default` | `--text-xs` |
| `--region-title-support` | `--text-2xs` |
| `--region-ground-support` | `--surface-recessed` |
| `--region-edge-support` | transparent |

A rank gets a token only where it **differs from `default`**. Lead has no
ground and no edge token of its own, and that absence is the model rather than
an omission: lead does not lift, so there is nothing to name.

## Behavior

### Rank does not descend

A region inside a lead region, carrying no rank of its own, is `default`. What a
lead region leads is the screen, not its own contents — a rank that inherited
would be a dimmer with a nicer name.

```html preview
<div class="inst-panel" data-rank="lead">
  <div class="inst-panel-header"><span class="inst-panel-title">Lead panel</span></div>
  <div class="inst-panel-body inst-stack inst-stack--tight">
    <div class="inst-card"><div class="inst-card-title">An ordinary card</div></div>
    <div class="inst-card" data-rank="support"><div class="inst-card-title">A quiet one</div></div>
  </div>
</div>
```

### Rank is not depth

Depth owns **direction**: a region inside a region recedes, always, without
consulting rank. Rank owns **presence**. The two are visible not touching in a
lead card inside a support panel — the card rises in its name and still recedes
in its ground, because those are different questions.

A region takes at most **one** film. Depth asks for it and `support` asks for it,
and a region receiving both requests still paints one: the depth of a nesting
shows through the stack of painted elements, not through the strength of one
fill. So rank's ground lever bites at depth 0 only, and a support region deeper
in shows its rank in its name and its missing edge. See
[elevation](./elevation.md).

### Rank is not tone

A lead region is the principal one, not the good one, and the two axes are
independent: rank says this is the place to look, [tone](./colors.md) says what
it means. They combine, and that combination is the point of keeping them apart
— a region that owns its screen **and** is in error:

```html preview
<div class="inst-panel" data-rank="lead" data-tone="error">
  <div class="inst-panel-header">
    <span class="inst-panel-title">Database unavailable</span>
    <span class="inst-panel-actions"><span class="inst-badge" data-tone="error"><span class="inst-dot"></span>P1</span></span>
  </div>
  <div class="inst-panel-body">Started 03:12. Three retries, none successful.</div>
</div>
```

**The tone reaches the ground and not the edge,** which is what
[a banner](../components/feedback/banner.md) already does: it keeps `--border`
and puts the tone on its background. A toned box outlined in the tone as well
states one fact twice.

**Tone outranks rank on the ground.** A support region that is also in error is
an error, drawn quietly: the rung, the ink and the missing edge still say
`support`, and the ground says `error`. Meaning wins the one property they
contend for.

**It does not descend,** for the same reason rank does not. A card inside a
toned panel shows that tone only by asking for it.

## Rules

### What to take and what not to

| Use | Take instead |
|---|---|
| Rank, to say which region is the point of a screen | **A larger heading in the markup** — the size of an interface heading is a matter of rank, not of level. See [typography](./typography.md) |
| [Tone](./colors.md), to say what a region means | **`lead`, to mean "good" or "urgent"** — rank carries no meaning and nothing announces it |
| [Density](./density.md), to make a region tighter | **`support`, to fit more in** — rank changes weight of presence, not size of geometry |
| A [button](../components/actions/button.md) weight, for how loudly an action asks | **Rank on a control** — a button already answers that question with four weights, and two axes for one question are one axis with two names |
| A [panel](../components/display/panel.md) inside a split half, ranked | **Rank on the half** — a layout primitive does not draw, and one that painted could not be nested inside something that paints otherwise |

### When not to use it at all

A screen with no hierarchy takes no rank. A conversation is the example the
library ships: in [the chat block](../blocks/chat.md) the turn is the unit and
every turn is equal, so nothing there carries a rank. An axis that applies
everywhere is a decoration; one that declines a case is doing its job.

## API

```api
```

## Related

[Elevation and surfaces](./elevation.md)
[Typography](./typography.md)
[Density](./density.md)
[Scale](./scale.md)
[Panel](../components/display/panel.md)
[Card](../components/display/card.md)
[Section](../layout/section.md)
