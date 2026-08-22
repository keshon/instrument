---
title: AI chat
group: Blocks
source: src/kit.css
template: splash
---

A conversation with a machine that does work while it answers: reasoning, tool
calls and a request for permission all stand inside one turn. Not a single
class except the library's own, and not a single line of script.

```html preview context
<div class="inst-shell">
  <aside class="inst-shell-aside">
    <a class="inst-shell-brand" href="#"><span class="inst-dot" data-tone="running"></span>instrument</a>
    <nav class="inst-nav" aria-label="Sections">
      <div class="inst-nav-group">
        <span class="inst-nav-label">Today</span>
        <a class="inst-nav-item" href="#" aria-current="page">The share bar</a>
        <a class="inst-nav-item" href="#">The axis floor</a>
        <a class="inst-nav-item" href="#">Release 0.8.0</a>
      </div>
    </nav>
  </aside>

  <header class="inst-shell-header">
    <nav aria-label="Breadcrumbs">
      <ol class="inst-crumbs">
        <li><a href="#">instrument</a></li>
        <li><span aria-current="page">The share bar</span></li>
      </ol>
    </nav>
    <span class="inst-cluster-spacer"></span>
    <span class="inst-statusbar-item"><span class="inst-dot" data-tone="running"></span>working</span>
  </header>

  <main class="inst-shell-main">
    <div class="inst-container">
      <div class="inst-stack">

        <div class="inst-stack inst-stack--loose">

          <div class="inst-turn" data-from="user">
            <div class="inst-turn-head"><span class="inst-turn-who">You</span><span>14:02</span></div>
            <div class="inst-turn-body">Add a share bar above the history strip, and let one legend key both.</div>
          </div>

          <div class="inst-turn">
            <div class="inst-turn-head"><span class="inst-turn-who">Agent</span><span>14:02</span></div>

            <details class="inst-thinking">
              <summary class="inst-thinking-head"><span class="inst-step-twist"></span>Thought for 4 s</summary>
              <div class="inst-thinking-body">The strip answers when, the share answers how much. The legend swatch already falls through to the tone, so it can key both without a second vocabulary.</div>
            </details>

            <div class="inst-turn-body">I'll put the bar above the strip and let the swatch read the tone.</div>

            <details class="inst-step" open>
              <summary class="inst-step-head">
                <span class="inst-step-twist"></span>
                <span class="inst-dot" data-tone="ok"></span>
                <span class="inst-step-name">read</span>
                <span class="inst-step-sub">src/data.css</span>
                <span class="inst-step-meta">0.2 s</span>
              </summary>
              <div class="inst-step-body">
                <div class="inst-code">.inst-legend-swatch { background: var(--series, var(--tone-mark)); }</div>
              </div>
            </details>

            <details class="inst-step">
              <summary class="inst-step-head">
                <span class="inst-step-twist"></span>
                <span class="inst-dot" data-tone="ok"></span>
                <span class="inst-step-name">edit</span>
                <span class="inst-step-sub">src/data.css</span>
                <span class="inst-step-meta">0.4 s</span>
              </summary>
              <div class="inst-step-body">Added <code>.inst-share</code> with one part per outcome.</div>
            </details>

            <div class="inst-turn-actions">
              <button class="inst-btn inst-btn--sm inst-btn--icon inst-btn--ghost" type="button" aria-label="Copy"><svg class="inst-icon" aria-hidden="true"><use href="#i-copy"/></svg></button>
              <button class="inst-btn inst-btn--sm inst-btn--icon inst-btn--ghost" type="button" aria-label="Retry"><svg class="inst-icon" aria-hidden="true"><use href="#i-refresh"/></svg></button>
              <button class="inst-btn inst-btn--sm inst-btn--icon inst-btn--ghost" type="button" aria-label="More"><svg class="inst-icon" aria-hidden="true"><use href="#i-more"/></svg></button>
            </div>
          </div>

          <div class="inst-turn" data-from="user">
            <div class="inst-turn-head"><span class="inst-turn-who">You</span><span>14:05</span></div>
            <div class="inst-turn-body">Now drop the app rules this replaces.</div>
          </div>

          <div class="inst-turn">
            <div class="inst-turn-head"><span class="inst-turn-who">Agent</span><span>14:05</span></div>
            <div class="inst-turn-body">Sixty lines go. This rewrites the file in place, so I need permission.<span class="inst-caret"></span></div>

            <div class="inst-approval" data-state="pending">
              <div class="inst-approval-head">Delete 60 lines from static/beacon.css</div>
              <div class="inst-approval-what">static/beacon.css</div>
              <ul class="inst-approval-effects">
                <li>The axis moves to the kit</li>
                <li data-tone="warn">Irreversible: the file is rewritten in place</li>
              </ul>
              <div class="inst-approval-actions">
                <button class="inst-btn inst-btn--sm inst-btn--primary" type="button">Allow</button>
                <button class="inst-btn inst-btn--sm" type="button">Deny</button>
              </div>
            </div>
          </div>

        </div>

        <div class="inst-stack inst-stack--tight">
          <div class="inst-composer">
            <textarea class="inst-textarea" rows="2" placeholder="Describe a task" aria-label="Describe a task"></textarea>
            <div class="inst-composer-bar">
              <span class="inst-cluster-spacer"></span>
              <button class="inst-btn inst-btn--sm inst-btn--icon inst-btn--ghost" type="button" aria-label="Send"><svg class="inst-icon" aria-hidden="true"><use href="#i-chevron"/></svg></button>
            </div>
          </div>

          <div class="inst-cluster inst-cluster--tight">
            <button class="inst-btn inst-btn--sm inst-btn--ghost" type="button" aria-haspopup="menu">Auto</button>
            <button class="inst-btn inst-btn--sm inst-btn--icon inst-btn--ghost" type="button" aria-haspopup="menu" aria-label="Add"><svg class="inst-icon" aria-hidden="true"><use href="#i-more"/></svg></button>
            <span class="inst-cluster-spacer"></span>
            <button class="inst-btn inst-btn--sm inst-btn--ghost" type="button" aria-haspopup="menu">Opus 5</button>
            <button class="inst-btn inst-btn--sm inst-btn--ghost" type="button" aria-haspopup="menu">High</button>
            <span class="inst-badge" data-tone="warn">context 90%</span>
          </div>
        </div>

      </div>
    </div>
  </main>
</div>
```

## The session, not the message

The row under the composer holds what outlives one turn — the permission mode,
the model, the effort, how much context is left. Those are **not** the field's
controls, so they stand outside its frame in an ordinary
[cluster](../layout/flow.md). Sending belongs to the message, so it stays
inside the frame at the field's trailing edge.

```html preview
<div class="inst-cluster inst-cluster--tight">
    <button class="inst-btn inst-btn--sm inst-btn--ghost" type="button" aria-haspopup="menu">Auto</button>
    <span class="inst-cluster-spacer"></span>
    <button class="inst-btn inst-btn--sm inst-btn--ghost" type="button" aria-haspopup="menu">Opus 5</button>
    <button class="inst-btn inst-btn--sm inst-btn--ghost" type="button" aria-haspopup="menu">High</button>
    <span class="inst-badge" data-tone="warn">context 90%</span>
  </div>
```

## Choosing among named things

```html preview
<div class="inst-popover" style="inline-size:20rem">
  <div class="inst-menu" role="menu" aria-label="Models">
    <span class="inst-menu-label">Models</span>
    <button class="inst-menu-item" role="menuitemradio" aria-checked="false" type="button">Fable 5<span class="inst-menu-shortcut">1</span></button>
    <button class="inst-menu-item" role="menuitemradio" aria-checked="true" type="button">Opus 5<span class="inst-menu-shortcut">2</span></button>
    <button class="inst-menu-item" role="menuitemradio" aria-checked="false" type="button">Sonnet 5<span class="inst-menu-shortcut">3</span></button>
    <button class="inst-menu-item" type="button" role="menuitem" aria-haspopup="menu">More models<span class="inst-menu-shortcut">&rsaquo;</span></button>
    <hr class="inst-menu-sep">
    <span class="inst-menu-label">Fast mode</span>
    <label class="inst-switch"><input type="checkbox">Enable fast mode</label>
  </div>
</div>
```

A [menu](../components/overlays/menu.md) already carries every part this needs:
a group label, a name with an explanation under it, a shortcut at the far edge,
`aria-checked` for the one in force, and a rule between sections. A switch put
in a menu reads as a row of it — name first, state at the trailing edge.

## What is left of the budget

```html preview
<div class="inst-panel" style="max-inline-size:24rem">
  <div class="inst-panel-body inst-stack inst-stack--tight">
    <div>
      <div class="inst-meter-row"><span>Context window</span><span class="inst-meter-value">904.2k / 1M</span></div>
      <div class="inst-share" role="img" aria-label="Context window: 90% used">
        <span class="inst-share-part" data-tone="running" style="--share:0.62"></span>
        <span class="inst-share-part" data-tone="warn" style="--share:0.14"></span>
        <span class="inst-share-part" data-tone="error" style="--share:0.06"></span>
        <span class="inst-share-part" data-tone="neutral" style="--share:0.08"></span>
      </div>
    </div>
    <div>
      <div class="inst-meter-row"><span>5-hour limit</span><span class="inst-meter-value">resets in 4 hr · 10%</span></div>
      <div class="inst-meter" role="progressbar" aria-label="5-hour limit" aria-valuenow="10" aria-valuemin="0" aria-valuemax="100">
        <div class="inst-meter-fill" style="inline-size:10%"></div>
      </div>
    </div>
    <div>
      <div class="inst-meter-row"><span>Weekly</span><span class="inst-meter-value">resets Fri · 8%</span></div>
      <div class="inst-meter" role="progressbar" aria-label="Weekly" aria-valuenow="8" aria-valuemin="0" aria-valuemax="100">
        <div class="inst-meter-fill" style="inline-size:8%"></div>
      </div>
    </div>
  </div>
</div>
```

The context window is a [share bar](../components/charts/share.md): one part per
kind of occupant, and what is unaccounted for stays track. Each limit under it
is a [meter](../components/charts/meter.md) — one quantity against its own
ceiling, which is what a meter is for and what a share bar is not.

## What it is assembled from

| The part | What it is assembled with |
|---|---|
| The framework | [The application shell](../layout/shell.md) |
| One step of the exchange | [The turn](../agent/turn.md) — recessed for the person, bare for the machine |
| The working before the answer | [Thinking](../agent/thinking.md), collapsed |
| A tool call | [The step](../agent/step.md) with its arguments and result |
| A stop for permission | [The approval](../agent/approval.md) |
| Text still arriving | The caret, from [the badge and the dot](../components/display/badge.md) |
| Writing the next turn | [The composer](../agent/composer.md) |

## What is seen here

**An answer is a container, not a paragraph.** Everything the machine did is
inside the turn that answers, in the order it happened. Nothing floats in a
side panel, so there is no question of which run a step belonged to.

**The person is recessed and the machine is bare.** The answer holds framed
things — steps, an approval — and a surface under them would be a plate inside
a plate. The question needs no frame to be found: it is the short one.

**Thinking is not a step.** A step is a call with a result that can fail;
thinking has no status, so it carries no dot and no tone, and it is shut until
asked for.

**The composer is one frame.** The field and its controls are a single object,
so focus lights up the object the keyboard is in rather than half of it.

## Related

[The turn](../agent/turn.md) [Thinking](../agent/thinking.md)
[The composer](../agent/composer.md) [The step](../agent/step.md)
[The approval](../agent/approval.md) [The run console](./console.md)
