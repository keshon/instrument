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
              <button class="inst-btn inst-btn--sm inst-btn--ghost" type="button">Copy</button>
              <button class="inst-btn inst-btn--sm inst-btn--ghost" type="button">Retry</button>
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

        <div class="inst-composer">
          <textarea class="inst-textarea" rows="3" placeholder="Describe a task" aria-label="Describe a task"></textarea>
          <div class="inst-composer-bar">
            <button class="inst-btn inst-btn--sm inst-btn--ghost" type="button">Attach</button>
            <span class="inst-badge">opus</span>
            <span class="inst-cluster-spacer"></span>
            12 480 / 200 000
            <button class="inst-btn inst-btn--sm inst-btn--primary" type="button">Send</button>
          </div>
        </div>

      </div>
    </div>
  </main>
</div>
```

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
