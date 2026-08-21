---
title: Log
group: Agent layer
layout: component
source: src/text.css
js: Copying is done by `instrument.js`. Scrolling with the tail and virtualisation belong to the application layer
api:
  - { name: "inst-log", kind: "class", doc: "The container; it scrolls" }
  - { name: "inst-log-line", kind: "class", doc: "A line: three columns" }
  - { name: "inst-log-time", kind: "class", doc: "The timestamp" }
  - { name: "inst-log-level", kind: "class", doc: "The level, a width of `5ch`" }
  - { name: "data-tone", kind: "attribute", doc: "`warn` · `error`" }
  - { name: "--level-ink", kind: "variable", value: "--text-muted" }
  - { name: "--font-mono", kind: "token" }
  - { name: "--text-xs", kind: "token" }
  - { name: "--surface-sunken", kind: "token" }
  - { name: "--radius-md", kind: "token" }
  - { name: "--space-4", kind: "token" }
  - { name: "--pad-panel", kind: "token" }
  - { name: "--gap-inline", kind: "token" }
  - { name: "--warn-text", kind: "token" }
  - { name: "--err-text", kind: "token" }
---

A stream of lines from a machine: a time, a level, a message. Monospaced
setting and columns of fixed width, so a log is read in columns rather than as
continuous text.

```html preview
<div class="inst-log" role="log" aria-label="The log of the run">
  <div class="inst-log-line"><span class="inst-log-time">14:32:07</span><span class="inst-log-level">info</span><span>Starting worldgen-01</span></div>
  <div class="inst-log-line"><span class="inst-log-time">14:32:09</span><span class="inst-log-level">info</span><span>4 files read</span></div>
  <div class="inst-log-line" data-tone="warn"><span class="inst-log-time">14:32:11</span><span class="inst-log-level">warn</span><span>chunks.bin is busy, retrying in 1 s</span></div>
  <div class="inst-log-line" data-tone="error"><span class="inst-log-time">14:32:16</span><span class="inst-log-level">error</span><span>EBUSY: could not read chunks.bin</span></div>
</div>
```

## Contract

| What | Required | Why |
|---|---|---|
| A `role="log"` and an `aria-label` | yes | Lines arriving while the work goes on are otherwise not announced at all |
| Three nodes in a line: time, level, message | yes | The columns are set by a width in `ch`; a node out of place breaks the alignment of every line |
| A `data-tone` on `warn` and `error` alone | yes | An ordinary line carries no tone: in a stream of a thousand lines what is highlighted is what calls for attention |
| Monospaced setting | yes | It comes from `inst-log`. A proportional font scatters the column of times |

### Accessibility

| | |
|---|---|
| `role="log"` | A region records are added to. A screen reader announces the new lines without rereading the whole log |
| Liveness | A `role="log"` implies an `aria-live="polite"`. For a stream of hundreds of lines a second that has to be **switched off** explicitly: continuous speaking makes the interface unusable |
| The name of the region | An `aria-label` is required: two logs on a screen with no names are indistinguishable by ear |
| Time is data | `inst-log-time` takes `--text-muted` (4.5:1) rather than `--text-faint`: a timestamp is read rather than looked at |
| Not colour alone | The level is written as a word (`warn`, `error`) beside the colour |
| Scrolling | The region is reachable from the keyboard. Auto-scrolling to the bottom has to stop when the user has scrolled up |
| Wrapping | `white-space: pre-wrap` and `overflow-wrap: anywhere` on the message cell alone: on the container that would collapse the columns |

## Anatomy

The grid is declared **on the line**, that is, every line has one of its own.
If the level column is set by its content, its width changes from word to word:
`info` is narrower than `error`, and the messages go off in a staircase. So
`inst-log-level` gets a fixed width of `5ch`.

## Variants

```html
<div class="inst-log-line" data-tone="error">…</div>
```

| Tone | What is painted |
|---|---|
| no attribute | The level alone, muted |
| `warn` | The level |
| `error` | The level **and the whole message** |

An error is painted entire: it is looked for in a thousand lines by scanning
rather than by reading.

## JS

Include the module once per page — there is no need to initialise the
components one by one, `instrument.js` works by delegation and sees nodes that
arrived later.

```html
<script type="module" src="instrument.js"></script>
```

### What `instrument.js` does

The example in the header is live: hover over the log and press the copy
button.

Copying: a `.inst-copy` button inside an `.inst-copyable` puts the text in the
clipboard and answers with colour. The details are on the page of
[the block of code](../components/display/code.md#js).

```html
<div class="inst-log inst-copyable">
  <button class="inst-copy" type="button" aria-label="Copy the log"></button>
  …
</div>
```

### What is left to the application

| What | Why not `instrument.js` |
|---|---|
| Scrolling with the tail | "Hold the bottom until a person scrolls up" is a policy rather than styling |
| Virtualisation | An agent's log runs to tens of thousands of lines; the way depends on the source |
| Filtering by level | Data |

Scrolling to the tail is done by the application, and it is one line:

```js
if (atBottom) log.scrollTop = log.scrollHeight;
```

## API

```api
```

## Related

```related
```
