---
name: kit-change
description: How to change instrument — prove a hole with a throwaway HTML draft before adding anything, measure in the browser instead of eyeballing, and know which gate answers which question. Use when adding or changing a component, a token or a docs page in this repo, or when an app built on the kit needs something the kit does not have.
---

# Changing the kit

The rule this repo runs on: **a claim you have not measured is a guess.** Every
gate here exists because a guess shipped once. What follows is how to work with
that, not around it.

## 1. Prove the hole with a draft before adding anything

Never add a component because a screen seems to want one. Build the screen from
what already exists, look at it, and let it fail. The failure is the argument.

```bash
# a throwaway page, gitignored, served by the docs site
.preview-tmp/probe.html          # write it here
cp .preview-tmp/probe.html site/dist/probe.html
```

Link `/kit/kit.css` and write three or four honest attempts side by side, each
labelled. Then look at it and measure it.

This is how `turn`, `thinking` and `composer` were justified: a card is a raised
plate for one sentence, a panel adds a header bar, a note carries a glyph and a
tone. A textarea beside a toolbar is two frames; an input group stretches Send
to the field's full height; a panel around both nests borders. Six attempts, six
recorded failures, and the reason each is wrong went into the CSS comment.

If an attempt *works*, you have saved a component. That is the better outcome.

**A caller having to re-dress a component is the same signal.** When an app adds
CSS to make a kit component presentable, the hole is in the kit. `beacon.css`
says so in its own header, and sixty of its lines duly moved into the kit.

## 2. Measure in the browser; do not trust the screenshot

The pane is often not compositing, and eyes are worse than numbers anyway.

```js
// via the browser tool — geometry, not vibes
el.getBoundingClientRect()          // does it overflow? by how much?
getComputedStyle(el).borderRadius   // what was DECLARED
Math.min(declared, w/2, h/2)        // what gets PAINTED — the browser clamps
main.scrollWidth - main.clientWidth // is there a scrollbar, and whose fault
```

Sweep the axes rather than checking one state: five scales × three densities,
and container widths from 260 to 720px. Set `data-scale` / `data-density` on
`documentElement` and re-measure in a loop — several real defects this session
appeared only at one scale or one column width.

Three habits that caught real bugs:

- **Compare declared against painted.** Below ~8px the browser scales every
  radius on a box until two fit an edge. Seven components declared a corner
  nobody drew.
- **Ask who owns the scrollbar.** Walk up from the offending element until
  `overflowX !== 'visible'`; the culprit is usually invisible and laid out.
- **Simulate a state you cannot hover.** Inject a `@layer kit.agent` rule with a
  probe class, measure, remove it. Remember opacity transitions report their
  *start* value — kill `transition` first or you will read 0 and misdiagnose.

## 3. Know which gate answers which question

```bash
go -C tools run ./cmd/contrast     # colour pairs against WCAG
go -C tools run ./cmd/targets      # tap targets, 15 scale/density cells
go -C tools run ./cmd/proportion   # ratios between roles
go -C tools run ./cmd/docscheck    # docs against the kit, both directions
go -C tools run ./cmd/registry     # components.json against the CSS lists
go -C tools run ./cmd/lang         # Cyrillic in translated zones
go -C tools run ./cmd/dist -check  # dist matches src
go -C tools run ./cmd/mutate       # THE GATES THEMSELVES
node tools/audit-run.mjs --base http://localhost:4322   # rendered pixels
```

`mutate` is the one people forget, and forgetting it is how a gate goes quietly
dead: it breaks each invariant and demands red. **Anything that edits an anchor
a mutation points at — a nav list, a token line — must re-run it.** That anchor
drifted twice in one session and left a gate unguarded for a commit each time.

The pixel gate needs a running server and finds what token gates cannot: a
wrapped table term at scale 18, a focusable element 15px wide. Two behaviour
timeouts on the heaviest pages are pre-existing environment flake, not you —
confirm by stashing and re-running before you chase them.

## 4. Adding a component is eight edits, not one

Miss any and a gate stops you:

1. `src/*.css` — the rule, with the reasoning in the comment
2. `src/forced.css` and `src/print.css` — if the registry entry declares a policy
3. `components.json` — four axes, all required, plus **reciprocal** `related`
4. `docs/<section>/<name>.md` **and** `.ru.md`
5. `assets/sprite.svg` — an `i-p-<name>` symbol, or the build stops
6. `site/internal/nav/nav.go` — the order of the side column
7. the section index — its count and its "what to choose" table
8. `go -C tools run ./cmd/dist` — rebuild the delivery

Document only attribute values that have a CSS rule. `data-from="agent"` was
cut for exactly this: nothing to declare, so nothing to check.

## 5. Where prose goes

`docs/internal/DOCS-SHAPE.md` is the authority: **tables grow, prose does not.**
Reference — contract, API, measured thresholds — is exhaustive. Narrative goes
to the section index ("why this and not the neighbour", once for the kit) or the
CSS comment ("why it is built this way"). A component page sits near 300 prose
words; the complex ones reach ~800.

The commonest mistake is writing the reasoning into the CSS comment *and* the
page. Before adding a paragraph, check whether the comment already says it.

## 6. Record what you rejected

The CSS comments here carry failed attempts, and they earn their space: they
stop the next person rebuilding something already measured and thrown away.

State the defect in numbers. "51px of overflow on a shell whose main region is
`overflow-x: auto`, with the pointer nowhere near the strip" outlives "the
tooltip caused a scrollbar".
