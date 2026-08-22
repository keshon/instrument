# Audit — the design language against Anthropic's

Third pass. The first (`AUDIT-2026-08-07.md`) was about architecture and is
closed. The second (`AUDIT-2026-08-20.md`) was about defects and the move to
English. This one is about **the design language itself**: where `instrument`
agrees with the system it is modelled on, where it diverges, and — for every
divergence — which of the two is right.

---

## What this project is

`instrument` is a design system in plain HTML + CSS + JS. No build, no
dependencies, no framework. An application links one file and writes the markup
by hand; the kit supplies tokens, classes and a minimum of script for roles and
states.

| | |
|---|---|
| Sources | `src/` — 19 CSS files + `kit.js`, ~9 400 lines |
| Delivery | `dist/instrument.min.css`, `instrument.js` |
| Documentation | `docs/` — 101 pages, 242 live examples, EN + RU |
| Registry | `components.json` — 73 components, four axes + an ARIA contract |
| Tools | `tools/` — 8 Go commands with no dependencies |
| Reference site | `site/` — Go + goldmark |
| Axes | 5 themes × 4 accents × 5 scales × 3 densities = 300 combinations |
| Browser floor | Chrome 111, Safari 17.5, Firefox 128 |

Run everything: `check.bat` (Windows), or the steps in
`.github/workflows/ci.yml`. The pixel audit and the behaviour checks need a
server: `check serve`, then `check pixels` / `check behavior`.

---

## The premise of this audit

The kit is **a close reading of Anthropic's design language, primarily Claude
Code**. Not a clone: the kit has obligations Claude Code does not — five themes,
four accents, five scales, three densities, a contrast gate, a tap-target gate,
forced-colors, print. Some of what the reference does would fail the kit's own
gates, and in those cases the kit is right and the reference is not the target.

So this is **not** "make it look like the screenshot". For every divergence the
question is which of three answers holds:

1. **The kit is behind.** The reference solved something the kit has not, and
   nothing in the kit's obligations prevents adopting it. → adopt.
2. **The kit is deliberately different**, and the reason is written down
   somewhere. → confirm the reason still holds, and if it does, say so and stop.
3. **The kit is deliberately different and the reason is NOT written down.**
   → the worst case. It will be "fixed" by the next person who sees it. Write it
   down, in the CSS comment, where the decision lives.

A divergence with no stated reason is a finding even when the kit's behaviour
is correct.

---

## What is already settled — do not raise these again

- **The technology.** Tailwind, CSS-in-JS, utility-first, a monorepo, a bundler,
  a framework, codegen. All considered and declined with reasons in
  `docs/internal/DECISIONS.md`. If you think a reason has expired, name it and
  refute it specifically; do not propose the option afresh.
- **Four tiers of custom properties.** Deliberate architecture, not sprawl.
- **The nine gates.** `contrast · targets · proportion · docscheck · registry ·
  mutate · dist · lang · site`, plus the pixel audit and the behaviour checks.
  Do not propose a gate that exists. Do question whether an existing gate
  measures what it claims — that has been wrong twice.
- **Light control borders.** The reference's are lighter than the kit's and
  would fail `cmd/contrast` at 3:1 for non-text (WCAG 1.4.11). Already checked.
  The kit's border is `--border-control` at 3.46:1.

---

## Focus 1 — the divergences, ranked

Work from **rendered screens**, not from reading CSS. The method that has
worked in this repository is written up in `.claude/skills/kit-change/SKILL.md`:
build the screen out of the parts that exist, let it fail, and let the failure
be the argument. A gap you can only describe is a guess; a gap that breaks a
draft is a finding.

Assemble at least these, using only kit classes:

1. **A chat session** — `docs/blocks/chat.md` is the current attempt. Compare it
   against the reference in composition, not in colour.
2. **A tool-call trace** — steps, nested output, an approval request.
3. **A dense readout** — a status bar, a queue, a run with a history strip.
4. **A settings surface** — menus, switches, selects, a cascader.

For each, name what the reference does that the kit cannot express, and what
the kit expresses more clearly than the reference. Both directions.

Then go after the axes the reference does not have to care about and the kit
does. A divergence that only shows up in one cell of the matrix is still a
divergence:

- the smallest scale at `compact`, the largest at `comfortable`;
- every theme against every accent — 20 cells;
- `forced-colors`, `prefers-reduced-motion`, print.

### The specific questions

- **Type.** The reference's scale, weights and leading against the kit's.
  `--leading-ui` 1.4 vs `--leading-prose` 1.65 — is each used where it belongs?
  An agent's answer is prose; a control's label is not.
- **Density.** The kit ships three densities and defaults to the dense one. The
  reference has one. Which of the kit's three actually matches it, and is the
  default the right one?
- **Surfaces and depth.** The kit takes fields out of the depth stack
  (`--surface-field` = `--surface-raised`) and lets the border carry them. The
  reference does something else. Which reads better on a long thread, and why?
- **Chrome weight.** Count the visible borders, fills and shadows on one
  assembled screen in each. The reference is quieter. Where does the kit's extra
  weight come from — and is each instance load-bearing?
- **The glyph vocabulary.** Sizes, stroke, and which jobs share a token. The
  chevron already turned out to answer three different jobs with one number; see
  the H3 "Which size a chevron takes" on `docs/foundations/icons.md`. Look for
  other tokens carrying more than one job.
- **Motion.** Durations, easing, what animates at all. The kit has `--dur-1..n`
  and one easing family.
- **Composition.** Not tokens — *arrangement*. Where things sit relative to each
  other, what shares a row, what gets its own line, what is allowed to be
  quiet. This is where the kit has been furthest off, and it is the hardest to
  see from CSS.

---

## Focus 2 — the doctrine, not just the pixels

The kit's real asset is that every decision carries its reasoning, in the CSS
comment, next to the rule. The audit has to check that asset the same way it
checks the pixels.

- **Find decisions with no stated reason.** Grep for rules whose comment says
  *what* rather than *why*. A number with no derivation is the failure mode.
- **Find reasons that no longer hold.** A comment citing a measurement the kit
  has since changed is worse than no comment: it is a wrong answer with
  authority. Two have already been found this way — a canvas that rounded cap
  height to a whole pixel, and a band widened to admit the error it caused.
- **Find the same decision stated twice** in two places that can drift.
- **Check the doctrine against the reference.** Where the kit states a law
  (`docs/about/design-principles.md`), does the reference obey it? If the
  reference breaks a law the kit declares, either the law is wrong or the
  reference is — say which, with evidence.

---

## Focus 3 — what the kit has that the reference does not

The kit is not only a subset. It covers dashboards, queues, monitoring, audit
trails — surfaces Claude Code has no equivalent of. Those parts have no
reference to be measured against, and that is a risk of its own: they were
designed without one.

- Which components exist with no reference and no stated derivation?
- Do they obey the kit's own laws, or did they grow their own vocabulary?
- Is any of them a solved problem elsewhere that the kit reinvented?

---

## How to answer

1. **Verdict** — three or four sentences. Is the kit a faithful reading of the
   reference, and where is the largest divergence.
2. **Divergences** — ranked by how much they cost a reader of an assembled
   screen. Format:
   **what differs → shown on which screen → which is right and why → what it
   costs to close → P0/P1/P2**.
   For each, state explicitly which of the three answers above applies.
3. **Undocumented decisions** — separately from divergences. For each: the rule,
   the number with no derivation, and what should be written where.
4. **Stale reasoning** — comments whose evidence no longer holds.
5. **The unreferenced half** — the components with no model, and how they fare
   against the kit's own laws.
6. **What NOT to change** — things that look like divergences and are correct.
   Name them, so the next audit does not spend its budget here.

---

## Constraints

- **Audit first.** Change nothing in `src/`, `docs/`, `tools/` before the
  verdict. Reading, running the gates, measuring in the browser, and building
  throwaway drafts outside the repository are all expected.
- **Every claim carries evidence.** A line number, a command's output, a
  measurement from the rendered page. A claim without one does not count.
- **Prefer a rendered measurement to a read of the CSS.** The cascade has
  fooled this project before; the pixel audit exists because of it.
- **Say plainly what is done well**, and briefly. An audit that finds problems
  everywhere cannot tell the important from the trivial.
- **If this prompt is wrong** — about the premise, about what is settled, about
  where to look — say so with reasons.
