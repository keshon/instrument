---
title: Utilities
group: Foundations
layout: foundation
source: src/base.css
api:
  - { name: "--text-secondary", kind: "token" }
  - { name: "--font-mono", kind: "token" }
---

Five utilities, and that is the **ceiling**. If a sixth is wanted, what is
usually wanted is a component.

```html preview
<span class="inst-u-dim">muted text</span>
<span class="inst-u-mono">a4f7c2e</span>
<span class="inst-u-prose-nums">1 234 567</span>
```

## Contract

Each one exists because it describes work no component does.

```html
<span class="inst-u-truncate">A very long file name that does not fit</span>
<span class="inst-u-visually-hidden">Actions</span>
```

| What | Required | Why |
|---|---|---|
| A utility names **work** | yes | "Truncate", "hide visually", "tabular figures". Not "12px of padding" |
| `inst-u-visually-hidden` rather than `display: none` | yes | What is hidden through `display` falls out of the accessibility tree along with its meaning |
| Spacing by a utility | no, and there will be none | Rhythm is the work of the [flow primitives](../layout/flow.md), not of a screen's markup |


### Accessibility

| | |
|---|---|
| `inst-u-visually-hidden` | The one right way to hide text from the eye while leaving it to a screen reader. `display: none` and `visibility: hidden` remove it from everywhere, the speaking included |
| Focus inside something hidden | If a focusable element sits inside a hidden block, the focus travels to an invisible place. The classic case is "skip to content": what is wanted there is [`inst-skip-link`](../layout/shell.md), which appears on focus |
| `inst-u-truncate` | It cuts **visually**: the full text stays in the DOM and is spoken. For a mouse add a `title`, otherwise the value is available by ear only |
| `inst-u-dim` | `--text-secondary` holds 4.5:1. Do not confuse it with `--text-faint` — that is the threshold of decoration (3:1), and text may not be painted with it |
| `inst-u-prose-nums` | Proportional figures **may not** be put in a column: the digits stop standing one under another |

## Anatomy

`inst-u-mono` is declared **in one rule** with `code`, `kbd`, `samp`, `pre`:

```css
code, kbd, samp, pre, .inst-u-mono {
  font-family: var(--font-mono);
  font-size: 0.925em;
}
```

A monospaced font is optically larger than a sans at the same size, so the
correction is obligatory. It is declared once deliberately: while the utility
had a rule of its own, a `<code class="inst-u-mono">` got `0.925 × 0.925` and
went under the declared floor of 11px.

Inside a [block of code](../components/display/code.md) the correction is
already applied by the container — a second one is not needed.

## Scale

| Class | What it does | When |
|---|---|---|
| `inst-u-truncate` | Truncates a line with an ellipsis | A title that would otherwise push its neighbour off the edge: a row of a queue, a path in a diff, the label of a lane |
| `inst-u-visually-hidden` | Hides visually, leaves it to a screen reader | A label an icon replaces; the heading of a region, superfluous on screen but wanted by ear |
| `inst-u-mono` | Monospaced setting | An identifier or a hash **inside** ordinary text |
| `inst-u-prose-nums` | Proportional figures | A number in prose, where tabular figures look loose |
| `inst-u-dim` | `--text-secondary` | A secondary part of a line, for which starting a class is excessive |

## Rules

**Spacing utilities** — `mt-3`, `p-2`, `pt-0d25`. The reason is not aesthetic:
the scale is sparse at the top on purpose, so that "a little bigger" is simply
not in the set, and a set of spacing utilities takes that decision back with
its very first class and moves the responsibility for rhythm out of the library
and into the markup of every screen.

In their place, three flow primitives with a gap named by **intent rather than
by a number**: [the stack, the cluster and the grid](../layout/flow.md) with
steps `--tight` / ordinary / `--loose`. The density of a container retunes all
three at once.

| What is wanted | Take instead |
|---|---|
| Spacing between elements | [`inst-stack`](../layout/flow.md) or [`inst-cluster`](../layout/flow.md) |
| Alignment and layout | [the layout primitives](../layout/flow.md) |
| A text colour | A semantic token in a rule of your own, not a utility |
| A font size | A component. Type size is part of an element's role rather than decoration |

## API

```api
```

The utilities deliberately have no tokens of their own: each either switches a
property of the platform or takes one existing semantic token.

## Related

[Flow primitives](../layout/flow.md)
[Typography](./typography.md)
[Spacing and shape](./spacing.md)
[Code](../components/display/code.md)
[Shell](../layout/shell.md)
