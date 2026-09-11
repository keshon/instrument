---
title: Composition
group: Foundations
layout: foundation
source: src/surfaces.css
api:
  - { name: "data-rank", kind: "attribute", doc: "How much attention a region asks for: `lead` `default` `support`" }
  - { name: "data-tone", kind: "attribute", doc: "The meaning register of a region. On a region it moves the ground" }
  - { name: "data-scale", kind: "attribute", doc: "The base type size. Legal on any element, not only on the root" }
  - { name: "inst-panel--plain", kind: "modifier", doc: "No edge: the step of surface carries the boundary" }
  - { name: "inst-metric-row--joined", kind: "modifier", doc: "One surface parted by rules" }
  - { name: "inst-metric-row--headline", kind: "modifier", doc: "The band under a page header" }
---

Every other page here answers "what does this thing look like". This one answers
the question a screen asks: **what is the main thing on it.** Hierarchy is a law
of the kit — the data is what matters, and labels, frames and headers serve it —
but a law is not an instrument. These are the instruments.

Below is a screen of three tiers: the headline band, the lead region and one that withdrew.

```html preview
<div class="inst-stack">
  <div class="inst-metric-row inst-metric-row--joined inst-metric-row--headline">
    <div class="inst-metric">
      <div class="inst-metric-label">Success rate</div>
      <div class="inst-metric-value">94.2<span class="inst-metric-unit">%</span></div>
    </div>
    <div class="inst-metric">
      <div class="inst-metric-label">In the queue</div>
      <div class="inst-metric-value">37</div>
    </div>
  </div>

  <div class="inst-panel inst-panel--plain" data-rank="lead">
    <div class="inst-panel-header"><span class="inst-panel-title">The last 24 hours</span></div>
    <div class="inst-panel-body">The region the screen belongs to. No edge: the step
      of surface is enough across an area this size.</div>
  </div>

  <div class="inst-panel" data-rank="support">
    <div class="inst-panel-header"><span class="inst-panel-title">Budget</span></div>
    <div class="inst-panel-body">Withdrawn: a smaller name, quieter ink, a sunken ground
      and no edge. It did not get quieter — everything else stayed where it was.</div>
  </div>
</div>
```

## Contract

| | |
|---|---|
| One lead region per screen | Two leads are no lead. Checked on the rendered tree by `kitAudit.composition()` |
| One headline band per screen | The same reasoning, one tier up |
| Rank does not touch the data | It moves the name, the frame and the ground. Not the body text, not the numbers, not the rows |
| Rank does not inherit | A region inside a lead region is `default` unless it says otherwise. A lead region leads the screen, not its own contents |
| Depth is set by nesting, not by rank | A region inside a region sinks, and never asks |

## Anatomy

| Instrument | Kind | What it moves | Legal on |
|---|---|---|---|
| `data-rank="lead"` | attribute | The name of the region goes up one rung | panel · card · section |
| `data-rank="support"` | attribute | Name down, ink down, ground sinks one film, edge goes | panel · card · section |
| `data-tone` | attribute | The ground of the region takes the tone's wash | panel · card |
| `data-scale` | attribute | The whole type scale of the subtree | any element |
| `--plain` | modifier | The edge goes, the ground stays | panel · card |
| `--joined` | modifier | A row becomes one surface parted by rules | metric row |
| `--headline` | modifier | The value takes the ceiling of the scale | metric row |

## Rules

**Build the ladder by withdrawing, not by shouting.** Rank never raises — the
raised surface is the ceiling in every light theme, so a `lead` that tried to
rise would work in two themes of five and do nothing in three. What a screen
gets instead is `support` on everything that is not the answer, and the primary
stands out because its neighbours stepped back.

The one exception is the band under the page header. There the breadcrumb has
already named the screen, so the band is allowed to be the loudest type present
— it answers the question the reader arrived with. Everywhere else a number
stays under the name of the screen it stands on.

**Separation is one system at a time.** A region that has its own ground does
not also owe an edge, and the choice is the author's because it depends on size:
a step of 1.07 holds a boundary across a large region and is lost across a small
one. A row of cells is either separate plates (depth) or one surface parted by
rules (`--joined`) — never both.

## API

The full list is in the page header.

## Related

[Rank](./rank.md) · [Elevation](./elevation.md) · [Scale](./scale.md) ·
[Typography](./typography.md) · [Panel](../components/display/panel.md) ·
[Metric](../components/display/metric.md)
