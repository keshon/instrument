---
title: Budget
group: Agent layer
layout: component
source: src/agent.css
api:
  - { name: "inst-budget", kind: "class", doc: "The container" }
  - { name: "inst-budget-row", kind: "class", doc: "A line: the name of the resource and the value" }
  - { name: "inst-budget-value", kind: "class", doc: "What is spent" }
  - { name: "inst-budget-of", kind: "class", doc: "\"of N\" — smaller and quieter" }
  - { name: "inst-budget-note", kind: "class", doc: "The forecast" }
  - { name: "--space-3", kind: "token" }
  - { name: "--gap-inline", kind: "token" }
  - { name: "--text-sm", kind: "token" }
  - { name: "--text-xs", kind: "token" }
  - { name: "--weight-medium", kind: "token" }
  - { name: "--weight-normal", kind: "token" }
  - { name: "--text-muted", kind: "token" }
---

Spending against a limit: tokens, time, money. A component of agent systems,
because there the spending is **continuous and not evident** — a user learns of
an overrun after the fact unless it is shown in advance.

```html preview
<div class="inst-budget">
  <div class="inst-budget-row">
    <span>Tokens</span>
    <span class="inst-budget-value">184 200<span class="inst-budget-of">of 250 000</span></span>
  </div>
  <div class="inst-meter" role="progressbar" aria-label="Tokens"
       aria-valuenow="74" aria-valuemin="0" aria-valuemax="100">
    <div class="inst-meter-fill" style="inline-size:74%"></div>
  </div>
  <div class="inst-budget-note">At the current pace the limit will be reached in about 12 minutes</div>
</div>
```

## Contract

| What | Required | Why |
|---|---|---|
| An `inst-budget-of` **inside** the value | yes | "of 250 000" is part of the number rather than its neighbour. As a separate node beside it, it reads as a second measure |
| [A meter](../components/charts/meter.md) with a role and values | yes | The share exists in pixels of width alone while the bar has no `aria-valuenow` |
| An `inst-budget-note` with a forecast | yes | What is left, with no pace, does not answer the question people come here for: will it be enough |
| A non-breaking space in the groups of digits | yes | `184 200` otherwise breaks at the space and reads as two numbers |

### Accessibility

| | |
|---|---|
| The value | The bar carries a `role="progressbar"` with `aria-valuenow/min/max`. Without them the value exists only as a width in pixels and is available neither to a screen reader nor to a test |
| The name of the bar | An `aria-label` with the name of the resource. Three nameless bars in a row are indistinguishable |
| The forecast | It updates as the run goes on, so it has to live in an `aria-live="polite"` — but update **no more than once in 30 seconds**, otherwise the speaking never falls silent |
| Approaching the limit | It changes the tone of the bar, and that has to be repeated in words in the forecast: colour is not the only carrier |
| Tabular figures | On by default: the number grows in place and does not tug at its neighbours |

## Anatomy

`inst-budget-of` is set smaller and quieter, exactly as the unit of measure on
[a metric](../components/display/metric.md). At one size "184 200 of 250 000"
reads as one long number, and the eye cannot pick out what matters.

### The forecast matters more than what is left

```html
<div class="inst-budget-note">At the current pace the limit will be reached in about 12 minutes</div>
```

`inst-budget-note` is not a footnote but the main value of the component. "65
800 left" does not say whether the run will make it; "in about 12 minutes"
does. If there is no forecast, it is better not to show a budget at all — a
figure with no interpretation only alarms.

## API

```api
```

The bar is given by [the meter](../components/charts/meter.md); the budget has
no primitive of its own.

## Related

```related
```
