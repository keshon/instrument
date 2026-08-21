---
title: Run
group: Agent layer
layout: component
source: src/agent.css
js: Live numbers and stopping belong to the application layer. The folding of phases works with no script
api:
  - { name: "inst-dots", kind: "class", doc: "A counting meter: how many units of how many. The unit is an `inst-dot`" }
  - { name: "inst-card-head", kind: "class", doc: "The header of a card: the heading and the glyph of a transition at the far edge" }
  - { name: "inst-card--link", kind: "modifier", doc: "The card entire is a link. Without it `base.css` paints it in the accent and underlines it" }
  - { name: "role", kind: "attribute", value: "progressbar", doc: "On `inst-dots`. Together with `aria-valuenow/min/max`" }
  - { name: "data-tone", kind: "attribute", value: "neutral · running · ok · warn · error", doc: "On every dot: `ok` — done, `running` — going, with no attribute — not begun yet" }
  - { name: "--space-1", kind: "token", doc: "The gap between the dots" }
  - { name: "--size-dot", kind: "token", doc: "The side of a dot — from `inst-dot`" }
---

The screen of a run is not a component but an **assembly**: a header with a
name and counters, phases that fold, a table of participants. It has no class
of its own: everything wanted is there already, and a class of its own would
mean a panel under another name.

There is one new thing here: the **counting meter**.

```html preview context
<div class="inst-panel">
  <div class="inst-panel-header">
    <span class="inst-panel-title">audit-worldbox-1</span>
    <span class="inst-badge" data-tone="running"><span class="inst-dot"></span>running</span>
    <span class="inst-panel-actions">
      <button class="inst-btn inst-btn--sm inst-btn--danger" type="button">Stop</button>
    </span>
  </div>

  <div class="inst-panel-body inst-stack">
    <p class="inst-prose">An adversarial audit of worldbox-1: a search for self-deception by classes A–F.</p>

    <dl class="inst-kv">
      <dt>Running</dt><dd>21 s</dd>
      <dt>Agents</dt><dd>7</dd>
      <dt>Tokens</dt><dd>186 000</dd>
    </dl>

    <div class="inst-section">
      <div class="inst-section-head">
        <span class="inst-section-title">Phases</span>
      </div>

      <details class="inst-accordion-item" open>
        <summary class="inst-accordion-head">
          Analysis
          <span class="inst-dots" role="progressbar" aria-valuenow="4" aria-valuemin="0" aria-valuemax="7"
                aria-label="Analysis: agents finished">
            <span class="inst-dot" data-tone="ok"></span>
            <span class="inst-dot" data-tone="ok"></span>
            <span class="inst-dot" data-tone="ok"></span>
            <span class="inst-dot" data-tone="ok"></span>
            <span class="inst-dot" data-tone="running"></span>
            <span class="inst-dot"></span>
            <span class="inst-dot"></span>
          </span>
        </summary>
        <div class="inst-accordion-body inst-panel-body--flush">
          <table class="inst-table">
            <thead>
              <tr><th>Agent</th><th class="inst-num">Tokens</th><th class="inst-num">Calls</th><th class="inst-num">Time</th></tr>
            </thead>
            <tbody>
              <tr><td>analysis:docs-drift</td><td class="inst-num">38 200</td><td class="inst-num">4</td><td class="inst-num">18 s</td></tr>
              <tr><td>analysis:shared-and-chunk</td><td class="inst-num">39 600</td><td class="inst-num">5</td><td class="inst-num">18 s</td></tr>
              <tr><td>analysis:probes-assert</td><td class="inst-num">37 900</td><td class="inst-num">7</td><td class="inst-num">18 s</td></tr>
              <tr><td>analysis:silent-failure</td><td class="inst-num">36 800</td><td class="inst-num">4</td><td class="inst-num">18 s</td></tr>
              <tr data-state="running"><td>analysis:coverage-hole</td><td class="inst-num">—</td><td class="inst-num">—</td><td class="inst-num">—</td></tr>
              <tr><td>analysis:history</td><td class="inst-num">—</td><td class="inst-num">—</td><td class="inst-num">—</td></tr>
              <tr><td>analysis:eyes-only</td><td class="inst-num">—</td><td class="inst-num">—</td><td class="inst-num">—</td></tr>
            </tbody>
          </table>
        </div>
      </details>

      <details class="inst-accordion-item">
        <summary class="inst-accordion-head">
          Refutation
          <span class="inst-dots" role="progressbar" aria-valuenow="0" aria-valuemin="0" aria-valuemax="3"
                aria-label="Refutation: agents finished">
            <span class="inst-dot"></span>
            <span class="inst-dot"></span>
            <span class="inst-dot"></span>
          </span>
        </summary>
        <div class="inst-accordion-body">
          <div class="inst-empty">
            <span class="inst-empty-title">Not begun yet</span>
            <span class="inst-empty-desc">The phase starts when the analysis closes all seven agents.</span>
          </div>
        </div>
      </details>
    </div>
  </div>
</div>
```

## Contract

A run is not a class of its own but an **assembly**: a panel with a header, a
property list, a counting meter and sections with steps. It has exactly two
things of its own — the counting meter and the card as a link.

| What | Required | Why |
|---|---|---|
| A `role="progressbar"` with three values on `inst-dots` | yes | Otherwise "how many of how many" exists only as a row of circles |
| A `data-tone` on every dot | yes | `ok` — done, `running` — going, with no attribute — not begun yet |
| An `inst-card--link` on a card that is a link | yes | Without it the base layer paints it in the accent and underlines it as an ordinary link |
| A number beside the counting meter | yes | Seven dots are not counted by eye: "3 of 4" in words is obligatory |

### Accessibility

| | |
|---|---|
| The role | A `role="progressbar"` on `inst-dots` plus `aria-valuenow/min/max`. Without them "how much is done" exists only in the number of filled circles |
| The name | An `aria-label` with what exactly is being counted. "4 of 7" with no subject is not a message |
| The dots | Decorative to a screen reader: the value is carried by the container rather than by them. They need no labels of their own |
| Not colour alone | The unit that is going **pulses** rather than merely being coloured; completion is spoken by the `aria-valuenow` |
| Reduced motion | The pulse slows rather than going out: an indicator of busyness has to stay visible |
| Print | The dots are printed: they are data rather than an indicator of activity |

## Anatomy

```html preview
<span class="inst-dots" role="progressbar" aria-valuenow="4" aria-valuemin="0" aria-valuemax="7"
      aria-label="Agents finished">
  <span class="inst-dot" data-tone="ok"></span>
  <span class="inst-dot" data-tone="ok"></span>
  <span class="inst-dot" data-tone="ok"></span>
  <span class="inst-dot" data-tone="ok"></span>
  <span class="inst-dot" data-tone="running"></span>
  <span class="inst-dot"></span>
  <span class="inst-dot"></span>
</span>
```

There are seven units, and each is either done or not. **A share here is
invented precision:** 4 of 7 drawn as a bar at 57% looks plausible and reports
what the system does not know. That is the same lie as a determinate bar stuck
at 90%, only quieter.

The component starts no dot of its own: the unit is an ordinary
[`inst-dot`](../components/feedback/states.md). It already knows the tone, the
"running" pulse and forced-colours mode, and a second circle in the library
would be a second name for the same thing.

## States

| Tone | What it means |
|---|---|
| no attribute | Not begun yet |
| `data-tone="running"` | Going. The dot pulses |
| `data-tone="ok"` | Done |
| `data-tone="error"` | Failed |
| `data-tone="warn"` | Done with a remark |

**Only** the one that is going pulses, and it is the same pulse as on a queue
row and on a step: in the library it means exactly one thing, and starting a
second meaning for it is not allowed.

## JS

Include the module once per page — there is no need to initialise the
components one by one, `instrument.js` works by delegation and sees nodes that
arrived later.

```html
<script type="module" src="instrument.js"></script>
```

### What `instrument.js` does

Nobody does the folding of phases: it is on `<details>` and works with no
script.

### What is left to the application

Live numbers and stopping — data and a command rather than styling.

```js
source.addEventListener('message', (e) => {
  const { phase, done } = JSON.parse(e.data);
  const dots = run.querySelectorAll('#' + phase + ' .inst-dot');

  dots.forEach((dot, i) => {
    if (i < done) dot.dataset.tone = 'ok';
    else if (i === done) dot.dataset.tone = 'running';
    else delete dot.dataset.tone;
  });
  dots[0].closest('.inst-dots').setAttribute('aria-valuenow', done);
});
```

The `aria-valuenow` is updated along with the dots rather than "some time
later". The dots are a picture; for a screen reader the progress exists in the
attribute alone, and a forgotten attribute means a phase that stayed at zero
for it for good.

## Composition

The same run as one row — for when there are several and one has to be chosen.
Here the counting meter works best of all: fifteen dots are taken in at a
glance, and what is visible is not only "how many" but **that one failed with a
remark**. A bar at 20% would not report that.

```html preview context
<div class="inst-stack inst-stack--tight">
  <a class="inst-card inst-card--interactive inst-card--link" href="#">
    <span class="inst-card-head">
      <span class="inst-card-title">audit-worldbox-1</span>
      <svg class="inst-icon" aria-hidden="true"><use href="#i-chevron"/></svg>
    </span>
    <span class="inst-card-sub">Workflow · 15 agents · 7 min 58 s</span>
    <span class="inst-dots" role="progressbar" aria-valuenow="4" aria-valuemin="0" aria-valuemax="15"
          aria-label="Agents finished">
      <span class="inst-dot" data-tone="ok"></span>
      <span class="inst-dot" data-tone="ok"></span>
      <span class="inst-dot" data-tone="warn"></span>
      <span class="inst-dot" data-tone="ok"></span>
      <span class="inst-dot" data-tone="running"></span>
      <span class="inst-dot"></span><span class="inst-dot"></span>
      <span class="inst-dot"></span><span class="inst-dot"></span>
      <span class="inst-dot"></span><span class="inst-dot"></span>
      <span class="inst-dot"></span><span class="inst-dot"></span>
      <span class="inst-dot"></span><span class="inst-dot"></span>
    </span>
  </a>

  <a class="inst-card inst-card--interactive inst-card--link" href="#">
    <span class="inst-card-head">
      <span class="inst-card-title">review-terrain-08</span>
      <svg class="inst-icon" aria-hidden="true"><use href="#i-chevron"/></svg>
    </span>
    <span class="inst-card-sub">Workflow · 3 agents · 41 s</span>
    <span class="inst-dots" role="progressbar" aria-valuenow="3" aria-valuemin="0" aria-valuemax="3"
          aria-label="Agents finished">
      <span class="inst-dot" data-tone="ok"></span>
      <span class="inst-dot" data-tone="ok"></span>
      <span class="inst-dot" data-tone="ok"></span>
    </span>
  </a>
</div>
```

The card is **a link entire**. The tap target comes out the width of the row,
and there is no need to grow it separately. The chevron is decorative, and the
name of the link is given by the heading.

## Patterns

Not one new class besides `inst-dots`:

| The part of the screen | What it is assembled from |
|---|---|
| The frame and the header | [The panel](../components/display/panel.md) |
| The state of the run | [The badge](../components/display/badge.md) with a dot |
| Stopping | [The button](../components/actions/button.md), the `--danger` variant |
| The counters | [The key–value list](../components/display/kv.md) |
| A phase that folds | [The accordion](../components/feedback/accordion.md) on `<details>` |
| The participants | [The table](../components/display/table.md) with `inst-num` |
| An empty phase | [The empty state](../components/feedback/empty.md) |

**That is the test of coverage.** A screen that does not assemble without
inline styling means a hole in the library: the rule about a live specification
from [the design principles](../about/design-principles.md) reaches here too.

## API

```api
```

## Related

```related
```
