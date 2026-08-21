---
title: Run console
group: Blocks
source: src/kit.css
template: splash
---

The principal scenario of the library: on the left a list of what the machine
is doing, on the right the particulars of the chosen one. Not a single class
except the library's own, and not a single line of script.

```html preview context
<div class="inst-shell">
  <aside class="inst-shell-aside">
    <a class="inst-shell-brand" href="#"><span class="inst-dot" data-tone="running"></span>worldgen</a>
    <nav class="inst-nav" aria-label="Sections">
      <div class="inst-nav-group">
        <span class="inst-nav-label">Work</span>
        <a class="inst-nav-item" href="#" aria-current="page">Runs
          <span class="inst-badge inst-nav-count">4</span></a>
        <a class="inst-nav-item" href="#">Queue</a>
        <a class="inst-nav-item" href="#">Artefacts</a>
      </div>
      <div class="inst-nav-group">
        <span class="inst-nav-label">Project</span>
        <a class="inst-nav-item" href="#">Agents</a>
        <a class="inst-nav-item" href="#">Keys</a>
      </div>
    </nav>
  </aside>

  <header class="inst-shell-header">
    <nav aria-label="Breadcrumbs">
      <ol class="inst-crumbs">
        <li><a href="#">worldgen</a></li>
        <li><span aria-current="page">Runs</span></li>
      </ol>
    </nav>
    <span class="inst-cluster-spacer"></span>
    <span class="inst-search">
      <input class="inst-input inst-input--sm" type="search" placeholder="Search" aria-label="Search across runs">
    </span>
    <button class="inst-btn inst-btn--sm inst-btn--primary" type="button">Run</button>
  </header>

  <main class="inst-shell-main" id="main">
    <div class="inst-container inst-container--full">
      <div class="inst-split">

        <div class="inst-split-side inst-panel">
          <div class="inst-toolbar">
            <div class="inst-segmented" role="radiogroup" aria-label="Filter by state">
              <button type="button" role="radio" aria-checked="true"  tabindex="0">All</button>
              <button type="button" role="radio" aria-checked="false" tabindex="-1">Running</button>
              <button type="button" role="radio" aria-checked="false" tabindex="-1">Fell</button>
            </div>
          </div>
          <div class="inst-panel-body inst-panel-body--list">
            <div role="listbox" aria-label="Runs">
              <div class="inst-task" role="option" data-state="running" aria-selected="true" tabindex="0">
                <span class="inst-task-gutter"><span class="inst-dot"></span></span>
                <span class="inst-task-main">
                  <span class="inst-task-title inst-u-truncate">nightly-4127</span>
                  <span class="inst-task-sub">7 agents · main</span>
                </span>
                <span class="inst-task-meta">21 s</span>
              </div>
              <div class="inst-task" role="option" data-state="failed" aria-selected="false" tabindex="-1">
                <span class="inst-task-gutter"><span class="inst-dot"></span></span>
                <span class="inst-task-main">
                  <span class="inst-task-title inst-u-truncate">audio-pass</span>
                  <span class="inst-task-sub">1 agent · fix/mixdown</span>
                </span>
                <span class="inst-task-meta">4.2 s</span>
              </div>
              <div class="inst-task" role="option" data-state="warn" aria-selected="false" tabindex="-1">
                <span class="inst-task-gutter"><span class="inst-dot"></span></span>
                <span class="inst-task-main">
                  <span class="inst-task-title inst-u-truncate">biomes-04</span>
                  <span class="inst-task-sub">3 agents · main</span>
                </span>
                <span class="inst-task-meta">11.5 s</span>
              </div>
              <div class="inst-task" role="option" data-state="done" aria-selected="false" tabindex="-1">
                <span class="inst-task-gutter"><span class="inst-dot"></span></span>
                <span class="inst-task-main">
                  <span class="inst-task-title inst-u-truncate">terrain-88</span>
                  <span class="inst-task-sub">4 agents · main</span>
                </span>
                <span class="inst-task-meta">18.2 s</span>
              </div>
            </div>
          </div>
        </div>

        <div class="inst-split-main inst-stack">

          <div class="inst-page-header">
            <div class="inst-page-header-main">
              <h3 class="inst-page-title">nightly-4127</h3>
              <p class="inst-page-desc">An adversarial audit: the search for self-deception by classes A–F.</p>
            </div>
            <div class="inst-page-actions">
              <span class="inst-badge" data-tone="running"><span class="inst-dot"></span>running</span>
              <button class="inst-btn inst-btn--sm inst-btn--danger" type="button">Stop</button>
            </div>
          </div>

          <div class="inst-props">
            <div class="inst-prop">
              <span class="inst-prop-label" title="Model">Model</span>
              <span class="inst-prop-control">opus</span>
            </div>
            <div class="inst-prop">
              <span class="inst-prop-label" title="Agents">Agents</span>
              <span class="inst-prop-control">7 of 12</span>
            </div>
            <div class="inst-prop">
              <span class="inst-prop-label" title="Tokens spent">Tokens</span>
              <span class="inst-prop-control">186,000</span>
            </div>
          </div>

          <div class="inst-panel">
            <div class="inst-panel-header">
              <span class="inst-panel-title">Phases</span>
              <span class="inst-panel-actions">
                <button class="inst-btn inst-btn--sm inst-btn--ghost" type="button" aria-expanded="false"
                        data-details-all data-label-expand="Expand all"
                        data-label-collapse="Collapse all">Expand all</button>
              </span>
            </div>
            <div class="inst-panel-body inst-panel-body--list">
              <details class="inst-step" data-state="ok">
                <summary class="inst-step-head">
                  <span class="inst-step-twist"></span>
                  <span class="inst-dot"></span>
                  <span class="inst-step-name">parse</span>
                  <span class="inst-step-sub inst-u-truncate">7 agents, 4 classes</span>
                  <span class="inst-step-meta">38,200 · 18 s</span>
                </summary>
                <div class="inst-step-body">
                  <div class="inst-code">{ "classes": ["docs-drift", "shared-and-chunk"], "depth": 2 }</div>
                </div>
              </details>
              <details class="inst-step" data-state="running" open>
                <summary class="inst-step-head">
                  <span class="inst-step-twist"></span>
                  <span class="inst-dot"></span>
                  <span class="inst-step-name">refutation</span>
                  <span class="inst-step-sub inst-u-truncate">3 of 7 checked</span>
                  <span class="inst-step-meta">21,400 · running</span>
                </summary>
                <div class="inst-step-body">
                  <div class="inst-log" role="log" aria-label="The log of the phase">
                    <div class="inst-log-line"><span class="inst-log-time">14:32:07</span><span class="inst-log-level">info</span><span>probes-assert: hypothesis accepted</span></div>
                    <div class="inst-log-line" data-tone="warn"><span class="inst-log-time">14:32:11</span><span class="inst-log-level">warn</span><span>silent-failure: no reproduction, lowered</span></div>
                    <div class="inst-log-line"><span class="inst-log-time">14:32:16</span><span class="inst-log-level">info</span><span>coverage-hole: check started</span></div>
                  </div>
                </div>
              </details>
              <details class="inst-step">
                <summary class="inst-step-head">
                  <span class="inst-step-twist"></span>
                  <span class="inst-dot"></span>
                  <span class="inst-step-name">summary</span>
                  <span class="inst-step-sub inst-u-truncate">waiting for the refutation</span>
                  <span class="inst-step-meta">—</span>
                </summary>
                <div class="inst-step-body">
                  <div class="inst-code">waiting</div>
                </div>
              </details>
            </div>
          </div>

          <div class="inst-budget">
            <div class="inst-budget-row">
              <span>Budget</span>
              <span class="inst-budget-value">186&nbsp;000 <span class="inst-budget-of">of 250&nbsp;000</span></span>
            </div>
            <div class="inst-meter" role="progressbar" aria-label="The spending of the budget"
                 aria-valuenow="74" aria-valuemin="0" aria-valuemax="100">
              <div class="inst-meter-fill" style="inline-size:74%"></div>
            </div>
            <div class="inst-budget-note">At the current pace it will last for about 3 more runs.</div>
          </div>

        </div>
      </div>
    </div>
  </main>
</div>
```

## What it is assembled from

Not a single new class. The screen entire is existing components put side by
side.

| The part | What it is assembled with |
|---|---|
| The framework | [The application shell](../layout/shell.md): the column, the header, the work area |
| The two columns | [The split](../layout/split.md). It wraps by itself, with no media query |
| The filter | [A toolbar](../components/navigation/toolbar.md) with [a segmented control](../components/actions/segmented.md) |
| The list | [A task row](../agent/task.md) in a `role="listbox"`: the keyboard from [instrument.js](../foundations/behavior.md) |
| The name and the actions | [The page header](../layout/page-header.md), [a badge](../components/display/badge.md), [a button](../components/actions/button.md) |
| The facts about the run | [A list of properties](../components/display/kv.md) — `inst-props` |
| The course of the work | [A step](../agent/step.md) with a fold and [a log](../agent/log.md) inside |
| The remainder | [The budget of a run](../agent/budget.md) with [a meter](../components/charts/meter.md) |

## What is seen here

**The screen answers "what is happening now" rather than "how much in
total".** There is no row of large numbers at the top: the summary for the day
is a separate task and a separate screen, and here a person is looking at one
piece of work and wants to know its state. The numbers stand where the question
about them arises: the spending beside the budget, the duration in the row of
the run.

**The filter lives in a row rather than in a panel.** Three states are a
segmented control on one line rather than a dropdown or a side form with an
"Apply" button.

**The choice is a state of the list rather than a transition.** A
`role="listbox"` plus an `aria-selected`: the arrows walk across the runs, the
particulars on the right change. The row is not a link, because it leads
nowhere.

**Tone is one vocabulary.** A `data-tone` on the dot of the brand, on the badge
of the run, on the line of the log and on the stroke of a phase means one and
the same thing and is coloured from one place.

**What is running pulses rather than being highlighted.** The carrier of the
"running" state is the pulsing of the dot; tone here is secondary, and
therefore the same tone can mean interaction.

## States

A real screen shows something other than data half of the time. These three
cases are part of the contract rather than a separate piece of work "for
later".

```html preview
<div class="inst-panel">
  <div class="inst-empty">
    <div class="inst-empty-title">Nothing matched the "fell" filter</div>
    <div class="inst-empty-desc">Over the last day every run finished successfully.</div>
    <button class="inst-btn inst-btn--sm" type="button">Reset the filter</button>
  </div>
</div>
```

While the list is travelling, the room is taken by
[a skeleton](../components/feedback/skeleton.md) — as many rows as are expected,
and of the same height:

```html preview
<div class="inst-panel">
  <div class="inst-panel-body inst-stack inst-stack--tight">
    <div class="inst-skeleton inst-skeleton--title"></div>
    <div class="inst-skeleton inst-skeleton--line"></div>
    <div class="inst-skeleton inst-skeleton--line inst-skeleton--short"></div>
  </div>
</div>
```

A refusal is [a block of refusal](../agent/failure.md) rather than a line of red
text: it has a reason, a list of the attempts and at least one way out.

```html preview
<div class="inst-failure" role="alert">
  <div class="inst-failure-head">audio-pass did not start</div>
  <div class="inst-failure-reason">EBUSY: could not read <code>chunks.bin</code> — the file is taken by another run.</div>
  <ul class="inst-failure-tried">
    <li>Retry after 1 s — the same code</li>
    <li>Retry after 4 s — the same code</li>
  </ul>
  <div class="inst-failure-actions">
    <button class="inst-btn inst-btn--sm" type="button">Retry</button>
    <button class="inst-btn inst-btn--sm inst-btn--ghost" type="button">Open the log</button>
  </div>
</div>
```

## Parameters as a sheet rather than as a screen

The settings of a run do not deserve a screen of their own with a "Cancel /
Save" bar: it is the same object, only from another side.
[A sheet](../components/overlays/sheet.md) opens over it, losing nothing from
sight.

```html preview
<button class="inst-btn" type="button"
        onclick="document.getElementById('con-params').showModal()">The parameters of the run</button>

<dialog class="inst-dialog inst-sheet" id="con-params">
  <form method="dialog">
    <div class="inst-dialog-head">
      <span class="inst-dialog-title">Parameters</span>
      <button class="inst-btn inst-btn--sm inst-btn--ghost inst-btn--icon inst-dialog-close" type="submit" aria-label="Close">
        <svg class="inst-icon" aria-hidden="true"><use href="#i-close"/></svg>
      </button>
    </div>
    <div class="inst-dialog-body">
      <form class="inst-form inst-form--side">
        <div class="inst-field">
          <label class="inst-label" for="con-model">Model</label>
          <span class="inst-select-wrap">
            <select class="inst-select" id="con-model"><option>opus</option><option>sonnet</option></select>
          </span>
        </div>
        <div class="inst-field">
          <label class="inst-label" for="con-limit">The limit of tokens</label>
          <div class="inst-input-group">
            <input class="inst-input" id="con-limit" value="250000">
            <span class="inst-input-affix">per run</span>
          </div>
        </div>
        <div class="inst-field">
          <label class="inst-label" for="con-par">Agents at once</label>
          <input class="inst-slider" type="range" min="1" max="16" value="12" id="con-par" style="--fill:73%">
        </div>
      </form>
      <div class="inst-note" data-tone="neutral">The changes apply to the next run and are saved at once.</div>
    </div>
  </form>
</dialog>
```

**There is no "Cancel / Save" bar here.** It reports that nothing has happened
until it is pressed, whereas the switch has already fired and the field has
already sent its request. The changes apply at once, and that is said in a line
under the form.

## Related

[The application shell](../layout/shell.md) [A run](../agent/run.md)
[A task row](../agent/task.md) [A step](../agent/step.md)
[A block of refusal](../agent/failure.md) [A sheet](../components/overlays/sheet.md)
