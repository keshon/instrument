---
title: Density
group: Foundations
layout: foundation
source: src/tokens.css
api:
  - { name: "data-density", kind: "attribute", doc: "`compact` `comfortable`. The absence of the attribute is the third, base value" }
  - { name: "--control-h-xs/sm/md/lg", kind: "token" }
  - { name: "--control-pad-sm/md/lg", kind: "token" }
  - { name: "--pad-panel", kind: "token" }
  - { name: "--pad-card", kind: "token" }
  - { name: "--pad-cell-x", kind: "token" }
  - { name: "--row-pad-y", kind: "token" }
  - { name: "--gap-row", kind: "token" }
  - { name: "--gap-inline", kind: "token" }
  - { name: "--size-check", kind: "token" }
  - { name: "--size-switch-w", kind: "token" }
  - { name: "--size-switch-h", kind: "token" }
  - { name: "--size-thumb", kind: "token" }
  - { name: "--size-spinner", kind: "token" }
  - { name: "--size-gutter", kind: "token" }
  - { name: "--size-indent", kind: "token" }
  - { name: "--label-col", kind: "token" }
---

Density is an **attribute of a container**. One attribute on a panel retunes
everything inside it: the tier of roles changes, and the components hold no
numbers of their own.

Density moves geometry and **does not touch type size**: in all three modes the
text stays 14px, and what changes is the air around it. Size is the business of
the second axis — [scale](./scale.md) — and the axes add up.

```html preview
<div class="inst-cluster" data-density="compact">
  <button class="inst-btn inst-btn--sm" type="button">sm</button>
  <button class="inst-btn" type="button">md</button>
  <span class="inst-badge" data-tone="ok"><span class="inst-dot"></span>done</span>
  <input class="inst-input inst-input--sm" type="text" value="compact">
</div>
```

## Contract

```html
<section data-density="compact">
  <!-- everything inside is retuned -->
</section>
```

The attribute is inherited through the cascade of custom properties, so it
works at any depth and on any container. Nesting is allowed: an inner region
with a different value overrides the outer one.

| Value | What for |
|---|---|
| no attribute | The default. A panel tool, a dashboard, most screens |
| `data-density="compact"` | A dense tabular view: an agent queue, a log, a tree, a property inspector |
| `data-density="comfortable"` | Forms, settings, wizards — where there are few decisions and each wants attention |

### Accessibility

### Tap targets

:::warn
WCAG 2.5.8 (AA) asks for a tap target no smaller than 24×24 CSS pixels. At
density `compact` a control of size `sm` gives **22px** — below the threshold.
That is a deliberate trade: a dense tabular view exists for the sake of
density.
:::

| Density | `--control-h-sm` | `--control-h-md` | The 24px threshold |
|---|---|---|---|
| `compact` | 22px | 26px | `sm` **does not pass** |
| default | 26px | 32px | Passes |
| `comfortable` | 30px | 36px | Passes |

The rule of application: **`compact` is for dense tabular views, not for the
main navigation and not for the principal actions of a screen.** If a user hits
this control with a mouse every day, it must not be 22px. The WCAG exception
for targets that have an equivalent of a normal size applies only when such an
equivalent really exists.

### Everything else

| | |
|---|---|
| Type size | It does not shrink at any density: the floor of 11px stays the floor |
| Contrast | It does not depend on density — the same pairs are checked once |
| Text zoom | The heights are applied as `min-block-size`, so an increased browser font stretches a control rather than cutting it off. True at all three densities |
| Keyboard | Density changes neither the order of traversal nor the targets it hides: it is pure geometry |

## Scale

Exactly the tier of roles, and that is enough.

### Control heights and paddings

| Token | `compact` | default | `comfortable` |
|---|---|---|---|
| `--control-h-xs` | 18px | 20px | 22px |
| `--control-h-sm` | 22px | 26px | 30px |
| `--control-h-md` | 26px | 32px | 36px |
| `--control-h-lg` | 32px | 38px | 44px |
| `--control-pad-sm` | `--space-3` | `--space-4` | `--space-5` |
| `--control-pad-md` | `--space-4` | `--space-5` | `--space-6` |
| `--control-pad-lg` | `--space-5` | `--space-6` | `--space-7` |

### Region paddings

| Token | `compact` | default | `comfortable` |
|---|---|---|---|
| `--pad-panel` | `--space-4` | `--space-5` | `--space-6` |
| `--pad-card` | `--space-5` | `--space-6` | `--space-7` |
| `--pad-cell-x` | `--space-3` | `--space-4` | `--space-5` |
| `--row-pad-y` | `--space-2` | `--space-3` | `--space-4` |
| `--gap-row` | `--space-2` | `--space-3` | `--space-4` |
| `--gap-inline` | `--space-3` | `--space-4` | `--space-5` |

### Glyphs and gutters

| Token | `compact` | default | `comfortable` |
|---|---|---|---|
| `--size-check` | 13px | 15px | 17px |
| `--size-switch-w` | 26px | 30px | 34px |
| `--size-switch-h` | 15px | 18px | 20px |
| `--size-thumb` | 12px | 14px | 16px |
| `--size-spinner` | 13px | 14px | 16px |
| `--size-gutter` | 12px | 14px | 16px |
| `--size-indent` | 12px | 14px | 16px |

### What is NOT retuned

| What | Why |
|---|---|
| Type size | Density is about geometry. Shrinking the text along with the paddings means hitting the 11px floor within one step. Type size is moved by [scale](./scale.md) — the second axis, which adds up with this one |
| Colour and contrast | The thresholds do not depend on how many rows fit on a screen |
| Radii | The shape of an object does not change with how tightly it stands |
| `--size-dot`, `--size-marker`, `--size-icon` | Carriers of meaning. A dot at 4px does not read, a mark at 1px is not seen |
| `--space-*` | The scale is the foundation. Density moves the roles, not the scale under them |

## Behavior

The tier of roles is declared at zero specificity: a default with the same
specificity as an override stops being a default.

So an application overrides any role without a fight — including those density
does not retune:

```css
.my-inspector { --label-col: 120px; }
```

## Rules

A component that breaks at `compact` holds a hard-coded constant that ought to
be a role. Three densities are an obligatory run when a component is added, on
a par with the contrast check.

### What to take and what not to

| Use | Take instead |
|---|---|
| `compact` for a dense tabular view: a queue, a log, an inspector | **`compact` for the main navigation** — the tap targets go below the threshold, see [the warning](#tap-targets) and [the shell](../layout/shell.md) |
| `comfortable` for a form and for settings filled in thoughtfully | **`comfortable` for a long list** — half as many rows fit on a screen as [a queue row](../agent/task.md) needs |
| The attribute on a container — a panel, a section, a dialog | **The attribute on every component** — density is not a variant. The size of a control is chosen by `--sm`/`--lg` on [the button](../components/actions/button.md) |
| The default — writing nothing | **A density of your own as a number** in an application — a new value gets a role, and the list of roles is in [tokens](./tokens.md) |
| Changing the density of a whole region | **Mixing densities in one table** — the columns will part, see [the table](../components/display/table.md) |

## API

```api
```

## Related

[Rhythm and shape](./spacing.md)
[Tokens](./tokens.md)
[Table](../components/display/table.md)
[Queue row](../agent/task.md)
[Form](../components/inputs/form.md)
[Shell](../layout/shell.md)
