---
title: Step — a tool call
group: Agent layer
layout: component
source: src/agent.css
api:
  - { name: "inst-step", kind: "class", doc: "A step, on a `<details>`" }
  - { name: "inst-step-head", kind: "class", doc: "The header, on a `<summary>`" }
  - { name: "inst-step-twist", kind: "class", doc: "The triangle of expansion. It turns by `[open]`" }
  - { name: "inst-step-name", kind: "class", doc: "The name of the tool, monospaced" }
  - { name: "inst-step-sub", kind: "class", doc: "The argument of the call" }
  - { name: "inst-step-meta", kind: "class", doc: "The volume and the time, pushed to the end" }
  - { name: "inst-step-body", kind: "class", doc: "The body: the arguments and the output" }
  - { name: "inst-output", kind: "class", doc: "The wrapper of a folding output" }
  - { name: "inst-output-body", kind: "class", doc: "The output itself" }
  - { name: "inst-output-more", kind: "class", doc: "The unfold button with a number" }
  - { name: "data-state", kind: "attribute", value: "running · ok · failed", doc: "on `inst-step`" }
  - { name: "open", kind: "attribute", value: "native", doc: "on the `<details>`" }
  - { name: "data-details-all", kind: "attribute", value: "a selector of the region · empty", doc: "The \"expand all\" button. Empty means the nearest panel" }
  - { name: "data-label-expand", kind: "attribute", doc: "The label of the button in the folded state" }
  - { name: "data-label-collapse", kind: "attribute", doc: "The label in the unfolded one" }
  - { name: "data-truncated", kind: "attribute", value: "true · false", doc: "on `inst-output`" }
  - { name: "--space-1", kind: "token" }
  - { name: "--space-3", kind: "token" }
  - { name: "--gap-inline", kind: "token" }
  - { name: "--pad-cell-x", kind: "token" }
  - { name: "--text-xs", kind: "token" }
  - { name: "--font-mono", kind: "token" }
  - { name: "--size-dot", kind: "token" }
  - { name: "--text-muted", kind: "token" }
  - { name: "--tone-ink", kind: "token" }
  - { name: "--dur-1", kind: "token" }
---

One call of a tool by an agent: what was called, with what, and what came back.
A step has a **body**, and it is expanded by the platform — this is a
`<details>`, so the keyboard and the state come free.

```html preview
<details class="inst-step" data-state="ok">
  <summary class="inst-step-head">
    <span class="inst-step-twist"></span>
    <span class="inst-dot"></span>
    <span class="inst-step-name">read_file</span>
    <span class="inst-step-sub inst-u-truncate">terrain/heightmap.ts</span>
    <span class="inst-step-meta">240 lines · 0.4 s</span>
  </summary>
  <div class="inst-step-body">
    <div class="inst-code">{ "path": "terrain/heightmap.ts", "range": [1, 240] }</div>
  </div>
</details>
```

## Contract

| What | Required | Why |
|---|---|---|
| A `<details>` and a `<summary>` | yes | The expansion, the keyboard, the `open` state and search on the page come from the platform |
| A `data-state` on `inst-step` | yes | The phase of the call. The dot inside inherits its tone from here |
| An `inst-u-truncate` on the argument | yes | A path longer than the line otherwise pushes the header wide and drives the time past the edge |
| A `data-truncated` on `inst-output` | yes, if the output is truncated | The reader has to know that not everything is shown |

### Accessibility

| | |
|---|---|
| Expansion | Native. The `<summary>` takes focus, `Enter` and `Space` toggle it, and the browser announces `aria-expanded` itself |
| The name of the tool | Monospaced and machine-like (`read_file`). A screen reader reads it in pieces — if that gets in the way, add a human label beside it |
| The folded output | The number in the button is part of the accessible name: "Show all 240 lines" is spoken entire |
| The state in words | The tone of the name is not spoken. Success and refusal have to be in the text — in the `inst-step-meta` or as a badge |
| Reduced motion | The pulse of a `running` dot slows rather than going out |
| The triangle | Decorative, with no name of its own: the state of expansion is carried by the `<details>` itself |

## States

| `data-state` | What happens |
|---|---|
| no attribute | The step has not begun |
| `running` | Going. The dot pulses, the name takes the accent |
| `ok` | It came back successfully |
| `failed` | The call did not succeed |

A step has a vocabulary **of its own**: `ok` instead of `done`, and there is no
`skipped` in it. The vocabularies of `data-state` are listed separately on
every component — a queue row and a tool call have different phases of life.

### The folded output

The log of a tool is almost always longer than the screen. Truncating it in
silence is not allowed — that is a lie about the volume.

```html
<div class="inst-output" data-truncated="true">
  <pre class="inst-code inst-output-body">…</pre>
  <button class="inst-output-more" type="button">Show all 240 lines</button>
</div>
```

**A folded output names its number in words.** "Show all 240 lines" rather than
"…" and rather than "Show more": a user has to know how much they are not
seeing in order to decide whether it is worth unfolding.

## Behavior

An implementation of your own would cost an `aria-expanded`, a keyboard
handler, the management of focus and the synchronisation of state — and would
still lose to the native one on small things like search on the page, which
expands what is folded.

```html
<details class="inst-step" open>…</details>
```

Being open is the `open` attribute rather than a class or a `data-state`: there
is no second channel for one state here.

### Expand all

Six steps mean six presses to read a run entire, and six more to put it back.
The button goes into the header of the region:

```html preview
<div class="inst-panel">
  <div class="inst-panel-header">
    <span class="inst-panel-title">The last check</span>
    <span class="inst-panel-actions">
      <button class="inst-btn inst-btn--sm" type="button" aria-expanded="false"
              data-details-all data-label-expand="Expand all"
              data-label-collapse="Collapse all">Expand all</button>
    </span>
  </div>
  <div class="inst-panel-body inst-panel-body--list">
    <details class="inst-step" data-state="ok">
      <summary class="inst-step-head">
        <span class="inst-step-twist"></span>
        <span class="inst-step-name">TCP</span>
        <span class="inst-step-sub inst-u-truncate">the socket is open</span>
        <span class="inst-step-meta">38 ms</span>
      </summary>
      <div class="inst-step-body"><div class="inst-code">93.184.216.34:443</div></div>
    </details>
    <details class="inst-step" data-state="ok">
      <summary class="inst-step-head">
        <span class="inst-step-twist"></span>
        <span class="inst-step-name">TLS</span>
        <span class="inst-step-sub inst-u-truncate">the handshake is finished</span>
        <span class="inst-step-meta">64 ms</span>
      </summary>
      <div class="inst-step-body"><div class="inst-code">TLS 1.3</div></div>
    </details>
  </div>
</div>
```

The button works out its state from the content rather than remembering it:
while at least one step is closed, a press expands everything. A button with a
memory of its own sooner or later parts from what a person opened by hand.

| What | What for |
|---|---|
| An empty `data-details-all` | The region is the nearest `inst-panel`. With a value, a selector of the node wanted |
| `aria-expanded` | `instrument.js` moves it itself; the markup wants an initial value |
| `data-label-expand` / `data-label-collapse` | The labels. Without them the button changes its meaning in silence, without changing its text |

The `inst:details-all` event bubbles from the button, with `{ expand, count }`
in the `detail`.

## Patterns

The transcript of a run — what a step exists for. Apart, the steps look like a
list; together, like an account of what the machine did.

What to check by hand:

- **expand a step** — it is a `<details>`, and the expansion cost nothing:
  neither a line of script nor an `aria-expanded` by hand;
- **find the text with search on the page** (`Ctrl+F`) — the browser expands a
  folded step itself. An implementation of your own loses that;
- **the folded output names its number in words.** "Show all 240 lines" rather
  than "…": a user has to know how much they are not seeing in order to decide
  whether to unfold.

```html preview context
<div class="inst-panel">
      <div class="inst-panel-body inst-panel-body--list">
        <details class="inst-step" data-state="ok">
          <summary class="inst-step-head">
            <span class="inst-step-twist"></span><span class="inst-dot"></span>
            <span class="inst-step-name">read_file</span>
            <span class="inst-step-sub inst-u-truncate">terrain/heightmap.ts</span>
            <span class="inst-step-meta">0.3 s</span>
          </summary>
          <div class="inst-step-body">
            <div class="inst-code">{ "path": "terrain/heightmap.ts", "range": [1, 240] }</div>
            <div class="inst-output" data-truncated="true">
              <pre class="inst-code inst-output-body">export function decodeHeightmap(buf: ArrayBuffer): Float32Array {
  const view = new DataView(buf);
  const size = view.getUint32(0, true);
  const out = new Float32Array(size * size);
  for (let i = 0; i &lt; out.length; i++) {
    out[i] = view.getUint16(4 + i * 2, true) / 65535;
  }
  return out;
}

export function sampleBilinear(map: Float32Array, size: number, x: number, y: number) {
  const x0 = Math.floor(x), y0 = Math.floor(y);
  const x1 = Math.min(x0 + 1, size - 1), y1 = Math.min(y0 + 1, size - 1);
  const fx = x - x0, fy = y - y0;
  return  map[y0 * size + x0] * (1 - fx) * (1 - fy)
        + map[y0 * size + x1] * fx * (1 - fy)
        + map[y1 * size + x0] * (1 - fx) * fy
        + map[y1 * size + x1] * fx * fy;
}</pre>
              <button class="inst-output-more" type="button">Show all 240 lines</button>
            </div>
          </div>
        </details>

        <details class="inst-step" data-state="running" open>
          <summary class="inst-step-head">
            <span class="inst-step-twist"></span><span class="inst-dot"></span>
            <span class="inst-step-name">place_resources</span>
            <span class="inst-step-sub inst-u-truncate">pass 3 of 3</span>
            <span class="inst-step-meta">14.0 s</span>
          </summary>
          <div class="inst-step-body">
            <div class="inst-stack inst-stack--tight">
              <div class="inst-skeleton inst-skeleton--title"></div>
              <div class="inst-skeleton inst-skeleton--line"></div>
              <div class="inst-skeleton inst-skeleton--short"></div>
            </div>
          </div>
        </details>

        <details class="inst-step" data-state="failed">
          <summary class="inst-step-head">
            <span class="inst-step-twist"></span><span class="inst-dot"></span>
            <span class="inst-step-name">validate_paths</span>
            <span class="inst-step-sub inst-u-truncate">3 nodes unreachable</span>
            <span class="inst-step-meta">1.2 s</span>
          </summary>
          <div class="inst-step-body">
            <div class="inst-failure" role="alert">
              <div class="inst-failure-head">The validation of the paths did not pass</div>
              <div class="inst-failure-reason">PathError: unreachable nodes (18,204) (19,204) (20,205) — no navmesh link from region 7</div>
              <div class="inst-failure-tried">Attempts: 3, the last at 19:38:16</div>
              <div class="inst-failure-actions">
                <button class="inst-btn inst-btn--sm" type="button">Retry</button>
                <button class="inst-btn inst-btn--sm inst-btn--ghost" type="button">Skip the step</button>
                <button class="inst-btn inst-btn--sm inst-btn--ghost" type="button">Open region 7</button>
              </div>
            </div>
          </div>
        </details>
      </div>
    </div>
```

## API

```api
```

## Related

```related
```
