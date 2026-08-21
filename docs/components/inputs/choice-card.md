---
title: Choice card
group: Inputs
layout: component
source: src/forms.css
api:
  - { name: "inst-choice-card", kind: "class", doc: "The `<label>` container. The grid, the border, the states" }
  - { name: "inst-choice-card-title", kind: "class", doc: "The heading of an option" }
  - { name: "inst-choice-card-desc", kind: "class", doc: "The explanation of the consequences" }
  - { name: "--pad-card", kind: "token" }
  - { name: "--radius-md", kind: "token" }
  - { name: "--surface-raised", kind: "token" }
  - { name: "--border-control", kind: "token" }
  - { name: "--accent-border", kind: "token" }
  - { name: "--accent-bg", kind: "token" }
  - { name: "--accent-solid", kind: "token" }
  - { name: "--size-marker", kind: "token" }
  - { name: "--gap-inline", kind: "token" }
  - { name: "--text-sm", kind: "token" }
  - { name: "--text-xs", kind: "token" }
  - { name: "--leading-ui", kind: "token" }
  - { name: "--dur-1", kind: "token" }
---

A radio button or a checkbox that has a heading and an explanation of the
consequences. A bare one-line label is not enough when the options cost
differently.

```html preview
<label class="inst-choice-card">
  <input type="radio" name="mode" checked>
  <span><span class="inst-choice-card-title">In parallel</span>
    <span class="inst-choice-card-desc">Up to 12 agents at once. Faster, but the spending of tokens is higher and the order of the output is not guaranteed.</span></span>
</label>
<label class="inst-choice-card">
  <input type="radio" name="mode">
  <span><span class="inst-choice-card-title">Sequentially</span>
    <span class="inst-choice-card-desc">One agent at a time. Slower, but the log reads top to bottom with no interleaving.</span></span>
</label>
```

## Contract

| What | Required | Why |
|---|---|---|
| A `<label>` as the container | yes | A press anywhere on the card lands in the control. A `<div>` wrapper would turn the card into dead area around a small circle |
| An `<input>` as a direct child | yes | The states are written through `:has(> input…)`: a control in a sub-wrapper will not switch them on |
| One `<span>` for the text | yes | The grid of a card is two columns: the control and **one** block of text. Two sibling nodes will part the columns |
| `inst-choice-card-title` | yes | With no heading this is just a checkbox in a frame |
| A shared `name` on a radio group | yes | Without it the options do not exclude one another |
| A `fieldset` with a `legend` around the group | yes | Otherwise the question the cards answer is named nowhere |
| `inst-choice-card-desc` | no, but usually yes | With no explanation a card does not justify its area — take [the toggles](./toggles.md) |

Inside a card is a native `<input>` rather than a `div` with a role: the
keyboard, the announcement of state and the part in a form come from the
platform. `Tab` enters the group, the arrows switch the radio buttons, `Space`
presses a checkbox — nothing is intercepted.

The focus ring is drawn around the **whole card**: on the control inside it
would be around a circle the size of a dot. The heading and the description
inside the `<label>` become the name of the option entire, so a separate
`aria-label` is not needed and does harm — it would override the description.
The group is named by a `fieldset` with a `legend`, otherwise the options are
announced with no question.

The choice carries a mark at the edge and the native check of the control
rather than a background and a border alone. The tap target is the whole card,
the largest in the whole library. A `disabled` on the `<input>` dims the card
entire, and it shows **which** option is unavailable.

## Variants

Inside is a **real native control**, and the whole semantics of the choice
depends on its type. The library tells the two cases apart by nothing but what
the browser draws: the rule `:has(> input:checked)` works alike in both.

```html preview
<label class="inst-choice-card">
  <input type="checkbox" checked>
  <span><span class="inst-choice-card-title">Collect a trace</span>
    <span class="inst-choice-card-desc">A full log of every step. The run becomes about 15% slower.</span></span>
</label>
<label class="inst-choice-card">
  <input type="checkbox">
  <span><span class="inst-choice-card-title">Notify on completion</span>
    <span class="inst-choice-card-desc">An email to the address of the account when the queue empties.</span></span>
</label>
```

| Type | When | Required |
|---|---|---|
| `type="radio"` | One of several mutually exclusive modes | A shared `name` across the group, otherwise there is no mutual exclusion |
| `type="checkbox"` | Independent options, any number of them | A `name` of its own on each |

A group of radio cards is a group of fields, and it has to have a shared name:
wrap it in a `fieldset` with a `legend`, otherwise a screen reader announces
four unrelated options with no question for them to answer.

## States

```html preview
<label class="inst-choice-card">
  <input type="radio" name="st" checked>
  <span><span class="inst-choice-card-title">Chosen</span>
    <span class="inst-choice-card-desc">A border in the accent, a tinted background and a load-bearing mark at the edge.</span></span>
</label>
<label class="inst-choice-card">
  <input type="radio" name="st">
  <span><span class="inst-choice-card-title">Not chosen</span>
    <span class="inst-choice-card-desc">The ordinary border of a control.</span></span>
</label>
<label class="inst-choice-card">
  <input type="radio" name="st" disabled>
  <span><span class="inst-choice-card-title">Unavailable</span>
    <span class="inst-choice-card-desc">Unavailable on the current plan.</span></span>
</label>
```

| State | How it is set | What happens |
|---|---|---|
| hover | `:hover` | The border darkens to `--text-muted` |
| chosen | `:has(> input:checked)` | The border `--accent-border`, the background `--accent-bg` **and a load-bearing mark at the inner edge** |
| focus | `:has(> input:focus-visible)` | A ring around the whole card rather than around the circle |
| unavailable | `:has(> input:disabled)` | An opacity of `0.5`, the cursor `not-allowed` |

**The choice carries more than colour.** At the edge of a chosen card stands
the same load-bearing mark of thickness `--size-marker` as on the current
navigation item and on an active tab: one language for the whole library, and
it survives both colour blindness and print.

## Composition

### With a badge in the heading

```html preview
<label class="inst-choice-card">
  <input type="radio" name="tier">
  <span><span class="inst-choice-card-title">Distributed
      <span class="inst-badge" data-tone="warn">beta</span></span>
    <span class="inst-choice-card-desc">The agents go off to a pool of workers. No guarantees on timing yet.</span></span>
</label>
```

## Patterns

### A group with a shared question

```html preview context
<fieldset class="inst-fieldset">
  <legend>Run mode</legend>
  <div class="inst-grid">
    <label class="inst-choice-card">
      <input type="radio" name="run-mode" checked>
      <span><span class="inst-choice-card-title">In parallel</span>
        <span class="inst-choice-card-desc">Up to 12 agents at once.</span></span>
    </label>
    <label class="inst-choice-card">
      <input type="radio" name="run-mode">
      <span><span class="inst-choice-card-title">Sequentially</span>
        <span class="inst-choice-card-desc">One agent at a time.</span></span>
    </label>
  </div>
</fieldset>
```

A group of fields goes **with no border** by default: the cards are outlined
already, and a border around them would give a frame inside a frame. See
[the form](./form.md).

## API

```api
```

## Related

```related
```
