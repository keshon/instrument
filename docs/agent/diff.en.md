---
title: Diff
group: Agent layer
layout: component
source: src/agent.css
api:
  - { name: "inst-diff", kind: "class", doc: "The container" }
  - { name: "inst-diff--flush", kind: "modifier", doc: "With no border and no radius: a diff takes the panel entire, and a border of its own would lie as a second line over its boundary" }
  - { name: "inst-diff-head", kind: "class", doc: "The header: the path and the statistics" }
  - { name: "inst-diff-path", kind: "class", doc: "The path to the file. Truncated from the end" }
  - { name: "inst-diff-stat", kind: "class", doc: "`+N` / `−M`" }
  - { name: "inst-diff-body", kind: "class", doc: "The body; it scrolls" }
  - { name: "inst-diff-line", kind: "class", doc: "A line" }
  - { name: "inst-diff-num", kind: "class", doc: "The line number. Not selected by the mouse" }
  - { name: "inst-diff-code", kind: "class", doc: "The code itself" }
  - { name: "inst-diff-fold", kind: "class", doc: "A folded stretch" }
  - { name: "data-kind", kind: "attribute", value: "add · del", doc: "on `inst-diff-line` and on the numbers in `inst-diff-stat`" }
  - { name: "--font-mono", kind: "token" }
  - { name: "--text-xs", kind: "token" }
  - { name: "--ok-bg", kind: "token" }
  - { name: "--ok-text", kind: "token" }
  - { name: "--err-bg", kind: "token" }
  - { name: "--err-text", kind: "token" }
  - { name: "--border-subtle", kind: "token" }
  - { name: "--hairline", kind: "token" }
  - { name: "--radius-lg", kind: "token" }
  - { name: "--pad-cell-x", kind: "token" }
  - { name: "--text-muted", kind: "token" }
---

The main artefact of an agent's work: what exactly it changed in a file.

```html preview
<div class="inst-diff">
  <div class="inst-diff-head">
    <span class="inst-diff-path">terrain/heightmap.ts</span>
    <span class="inst-diff-stat"><span data-kind="add">+3</span> <span data-kind="del">−2</span></span>
  </div>
  <div class="inst-diff-body">
    <div class="inst-diff-line"><span class="inst-diff-num">41</span><span class="inst-diff-code">const size = 256;</span></div>
    <div class="inst-diff-line" data-kind="del"><span class="inst-diff-num">42</span><span class="inst-diff-code">let seed = 0;</span></div>
    <div class="inst-diff-line" data-kind="add"><span class="inst-diff-num">42</span><span class="inst-diff-code">let seed = Date.now();</span></div>
    <div class="inst-diff-fold">— 18 lines skipped —</div>
  </div>
</div>
```

## Contract

| What | Required | Why |
|---|---|---|
| A `data-kind` on a line | yes | The carrier of the category. A line without one is context, and that is the right default |
| An `inst-diff-num` as a separate node | yes | The number is not selected by the mouse: otherwise a copied diff arrives mixed up with line numbers |
| An `inst-diff-code` around the code | yes | Only it holds the monospaced setting and does not wrap the line |
| An `inst-diff-fold` in the place of a gap | no, but usually yes | A folded stretch is data about not everything being shown |

### Accessibility

| | |
|---|---|
| Not colour alone | A `+`/`−` sign in the first column. That is the main requirement of the component |
| Line numbers | `user-select: none` — when a fragment is copied the numbers do not reach the clipboard along with the code |
| Copying | The sign is a pseudo-element, so what is copied is clean code rather than `+const size` |
| The path to the file | Truncated from the end with an ellipsis. The full path has to stay in a `title` — otherwise the diff has no name |
| Scrolling | The body scrolls horizontally and has to be reachable from the keyboard |
| Print | The tinted backgrounds are replaced with printable ones in `@media print`, and the sign stays: on paper the tone disappears altogether |

## Anatomy

A line of addition carries a `+`, a line of deletion a `−`. That is the
library's law being carried out: **colour has no right to be the only carrier
of a state**. A diff where the changes are told apart by green and red
backgrounds alone turns into a continuous text for somebody colour-blind — and
in black-and-white print.

The sign is drawn by a pseudo-element from `data-kind`, so it does not reach a
selection and is not copied along with the code.

## Variants

```html
<div class="inst-diff-line" data-kind="add">…</div>
```

| Value | What it means |
|---|---|
| no attribute | A line of context, unchanged |
| `add` | Added |
| `del` | Deleted |

An added line **is not in the "ok" state** — it belongs to the kind
"addition". So it has an attribute of its own rather than a `data-tone`: a tone
would lie about the meaning. The numbers in the header of a diff carry the same
`data-kind` — one axis, one record.

## Behavior

```html
<div class="inst-diff-fold">— 18 lines skipped —</div>
```

As with the folded output of a step, a stretch **names its number in words**.
Skipping lines in silence is the same lie about the volume.

## Patterns

A diff with a header, statistics and a folded stretch.

What to check by hand: **select several lines and copy them**. What reaches the
clipboard is clean code — with no line numbers and no `+` and `−` signs. The
numbers carry `user-select: none` and the sign is drawn by a pseudo-element. A
diff whose code cannot be copied is useless at the moment it is wanted.

```html preview context
<div class="inst-diff inst-diff--flush">
  <div class="inst-diff-head">
    <span class="inst-diff-path">src/world/biomes.ts</span>
    <span class="inst-diff-stat"><span data-kind="add">+3</span> <span data-kind="del">−2</span></span>
  </div>
  <div class="inst-diff-body">
    <div class="inst-diff-fold">118 lines folded</div>
    <div class="inst-diff-line"><span class="inst-diff-num">119</span><span class="inst-diff-num">119</span><span class="inst-diff-code">  const seed = opts.seed ?? 0;</span></div>
    <div class="inst-diff-line" data-kind="del"><span class="inst-diff-num">120</span><span class="inst-diff-num"></span><span class="inst-diff-code">  const noise = simplex2(seed);</span></div>
    <div class="inst-diff-line" data-kind="del"><span class="inst-diff-num">121</span><span class="inst-diff-num"></span><span class="inst-diff-code">  if (noise &gt; 0.5) return TUNDRA;</span></div>
    <div class="inst-diff-line" data-kind="add"><span class="inst-diff-num"></span><span class="inst-diff-num">120</span><span class="inst-diff-code">  const noise = simplex2(seed, { octaves: 3 });</span></div>
    <div class="inst-diff-line" data-kind="add"><span class="inst-diff-num"></span><span class="inst-diff-num">121</span><span class="inst-diff-code">  // tundra and steppe overlapped in 12 cells — parting them by height</span></div>
    <div class="inst-diff-line" data-kind="add"><span class="inst-diff-num"></span><span class="inst-diff-num">122</span><span class="inst-diff-code">  if (noise &gt; 0.5 &amp;&amp; height &gt; 0.72) return TUNDRA;</span></div>
    <div class="inst-diff-line"><span class="inst-diff-num">122</span><span class="inst-diff-num">123</span><span class="inst-diff-code">  return STEPPE;</span></div>
  </div>
</div>
```

## API

```api
```

## Related

```related
```
