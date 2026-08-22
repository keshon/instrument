---
title: Uptime monitor
group: Blocks
source: src/kit.css
template: splash
---

One object under observation, and the two questions asked of it at once: **when**
something happened and **how much** of it there was. Not a single class except
the library's own, and not a single line of script.

```html preview context
<div class="inst-shell">
  <aside class="inst-shell-aside">
    <a class="inst-shell-brand" href="#"><span class="inst-dot" data-tone="ok"></span>beacon</a>
    <nav class="inst-nav" aria-label="Sections">
      <div class="inst-nav-group">
        <span class="inst-nav-label">Watch</span>
        <a class="inst-nav-item" href="#" aria-current="page">Monitors
          <span class="inst-badge inst-nav-count">12</span></a>
        <a class="inst-nav-item" href="#">Incidents
          <span class="inst-badge inst-nav-count" data-tone="error">1</span></a>
        <a class="inst-nav-item" href="#">Summary</a>
      </div>
      <div class="inst-nav-group">
        <span class="inst-nav-label">Node</span>
        <a class="inst-nav-item" href="#">Peers</a>
        <a class="inst-nav-item" href="#">Settings</a>
      </div>
    </nav>
  </aside>

  <header class="inst-shell-header">
    <nav aria-label="Breadcrumbs">
      <ol class="inst-crumbs">
        <li><a href="#">Monitors</a></li>
        <li><span aria-current="page">absalyamov.ru</span></li>
      </ol>
    </nav>
    <span class="inst-cluster-spacer"></span>
    <span class="inst-cluster inst-cluster--tight"><span class="inst-dot" data-tone="ok"></span>checked 12 s ago</span>
  </header>

  <main class="inst-shell-main">
    <div class="inst-container inst-container--full">
      <div class="inst-stack">

      <div class="inst-page-header">
        <div class="inst-page-header-main">
          <h1 class="inst-page-title">absalyamov.ru
            <span class="inst-badge" data-tone="ok"><span class="inst-dot"></span>up</span></h1>
          <p class="inst-page-desc">https://absalyamov.ru/ · http · every 8 minutes</p>
        </div>
        <div class="inst-page-actions">
          <button class="inst-btn inst-btn--sm" type="button">Check now</button>
          <button class="inst-btn inst-btn--sm" type="button">Mute for an hour</button>
        </div>
      </div>

      <div class="inst-panel">
        <div class="inst-panel-header">
          <span class="inst-panel-title">The last 24 hours</span>
          <span class="inst-cluster-spacer"></span>
          <span class="inst-badge" data-tone="warn">3 incidents this week</span>
        </div>
        <div class="inst-panel-body inst-stack">

          <div class="inst-share" role="img" aria-label="The last 24 hours: 95.8% up, 3.1% down, 1.1% not checked">
            <span class="inst-share-part" data-tone="ok" style="--share:0.958"></span>
            <span class="inst-share-part" data-tone="error" style="--share:0.031"></span>
            <span class="inst-share-part" data-tone="neutral" style="--share:0.011"></span>
          </div>

          <div>
        <div class="inst-history inst-history--lg" role="img" aria-label="The last 24 hours: 97 checks, 4 with an error, one hour with none">
          <span class="inst-history-group" style="--n:4"><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-tooltip-text">00:00 — 4 checks</span></span>
          <span class="inst-history-group" style="--n:4"><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-tooltip-text">01:00 — 4 checks</span></span>
          <span class="inst-history-group" style="--n:4"><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-tooltip-text">02:00 — 4 checks</span></span>
          <span class="inst-history-group" style="--n:4"><span class="inst-history-tick" data-tone="error"></span><span class="inst-history-tick" data-tone="error"></span><span class="inst-history-tick" data-tone="error"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-tooltip-text">03:00 — 4 checks</span></span>
          <span class="inst-history-group" style="--n:4"><span class="inst-history-tick" data-tone="error"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-tooltip-text">04:00 — 4 checks</span></span>
          <span class="inst-history-group" style="--n:4"><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-tooltip-text">05:00 — 4 checks</span></span>
          <span class="inst-history-group" style="--n:1"><span class="inst-history-tick"></span><span class="inst-tooltip-text">06:00 — no checks</span></span>
          <span class="inst-history-group" style="--n:4"><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-tooltip-text">07:00 — 4 checks</span></span>
          <span class="inst-history-group" style="--n:4"><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-tooltip-text">08:00 — 4 checks</span></span>
          <span class="inst-history-group" style="--n:4"><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-tooltip-text">09:00 — 4 checks</span></span>
          <span class="inst-history-group" style="--n:4"><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-tooltip-text">10:00 — 4 checks</span></span>
          <span class="inst-history-group" style="--n:1"><span class="inst-history-tick" data-tone="ok"></span><span class="inst-tooltip-text">11:00 — 1 check</span></span>
          <span class="inst-history-group" style="--n:4"><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-tooltip-text">12:00 — 4 checks</span></span>
          <span class="inst-history-group" style="--n:4"><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-tooltip-text">13:00 — 4 checks</span></span>
          <span class="inst-history-group" style="--n:4"><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-tooltip-text">14:00 — 4 checks</span></span>
          <span class="inst-history-group" style="--n:4"><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-tooltip-text">15:00 — 4 checks</span></span>
          <span class="inst-history-group" style="--n:4"><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-tooltip-text">16:00 — 4 checks</span></span>
          <span class="inst-history-group" style="--n:9"><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-tooltip-text">17:00 — 9 checks</span></span>
          <span class="inst-history-group" style="--n:4"><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-tooltip-text">18:00 — 4 checks</span></span>
          <span class="inst-history-group" style="--n:4"><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-tooltip-text">19:00 — 4 checks</span></span>
          <span class="inst-history-group" style="--n:4"><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="warn"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-tooltip-text">20:00 — 4 checks</span></span>
          <span class="inst-history-group" style="--n:4"><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-tooltip-text">21:00 — 4 checks</span></span>
          <span class="inst-history-group" style="--n:4"><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-tooltip-text">22:00 — 4 checks</span></span>
          <span class="inst-history-group" style="--n:4"><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-tooltip-text">23:00 — 4 checks</span></span>
        </div>
        <div class="inst-history-axis" aria-hidden="true">
          <span class="inst-history-axis-cell" style="--n:4">00</span>
          <span class="inst-history-axis-cell" style="--n:4" data-minor>01</span>
          <span class="inst-history-axis-cell" style="--n:4" data-minor>02</span>
          <span class="inst-history-axis-cell" style="--n:4">03</span>
          <span class="inst-history-axis-cell" style="--n:4" data-minor>04</span>
          <span class="inst-history-axis-cell" style="--n:4" data-minor>05</span>
          <span class="inst-history-axis-cell" style="--n:1">06</span>
          <span class="inst-history-axis-cell" style="--n:4" data-minor>07</span>
          <span class="inst-history-axis-cell" style="--n:4" data-minor>08</span>
          <span class="inst-history-axis-cell" style="--n:4">09</span>
          <span class="inst-history-axis-cell" style="--n:4" data-minor>10</span>
          <span class="inst-history-axis-cell" style="--n:1" data-minor>11</span>
          <span class="inst-history-axis-cell" style="--n:4">12</span>
          <span class="inst-history-axis-cell" style="--n:4" data-minor>13</span>
          <span class="inst-history-axis-cell" style="--n:4" data-minor>14</span>
          <span class="inst-history-axis-cell" style="--n:4">15</span>
          <span class="inst-history-axis-cell" style="--n:4" data-minor>16</span>
          <span class="inst-history-axis-cell" style="--n:9" data-minor>17</span>
          <span class="inst-history-axis-cell" style="--n:4">18</span>
          <span class="inst-history-axis-cell" style="--n:4" data-minor>19</span>
          <span class="inst-history-axis-cell" style="--n:4" data-minor>20</span>
          <span class="inst-history-axis-cell" style="--n:4">21</span>
          <span class="inst-history-axis-cell" style="--n:4" data-minor>22</span>
          <span class="inst-history-axis-cell" style="--n:4" data-minor>23</span>
        </div>
          </div>

          <ul class="inst-legend">
            <li class="inst-legend-item" data-tone="ok"><span class="inst-legend-swatch"></span>up <span class="inst-legend-value">95.8%</span></li>
            <li class="inst-legend-item" data-tone="error"><span class="inst-legend-swatch"></span>down <span class="inst-legend-value">3.1%</span></li>
            <li class="inst-legend-item"><span class="inst-legend-swatch" data-empty></span>not checked <span class="inst-legend-value">1.1%</span></li>
          </ul>

        </div>
      </div>

      <div class="inst-split">
        <div class="inst-split-main">
          <div class="inst-panel">
            <div class="inst-panel-header"><span class="inst-panel-title">Incidents</span></div>
            <table class="inst-table">
              <thead>
                <tr><th>Started</th><th>Lasted</th><th>Reason</th></tr>
              </thead>
              <tbody>
                <tr>
                  <td>03:12</td>
                  <td>14 min</td>
                  <td><span class="inst-badge" data-tone="error"><span class="inst-dot"></span>timeout</span></td>
                </tr>
                <tr>
                  <td>04:01</td>
                  <td>4 min</td>
                  <td><span class="inst-badge" data-tone="error"><span class="inst-dot"></span>502</span></td>
                </tr>
                <tr>
                  <td>20:35</td>
                  <td>3 min</td>
                  <td><span class="inst-badge" data-tone="warn"><span class="inst-dot"></span>slow</span></td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>

        <div class="inst-split-side">
          <div class="inst-panel">
            <div class="inst-panel-header"><span class="inst-panel-title">Availability</span></div>
            <div class="inst-panel-body inst-stack inst-stack--tight">
              <div>
                <div class="inst-meter-row"><span>24 hours</span><span class="inst-meter-value">95.8%</span></div>
                <div class="inst-meter" data-tone="warn" role="progressbar" aria-label="24 hours"
                     aria-valuenow="96" aria-valuemin="0" aria-valuemax="100">
                  <div class="inst-meter-fill" style="inline-size:95.8%"></div>
                </div>
              </div>
              <div>
                <div class="inst-meter-row"><span>7 days</span><span class="inst-meter-value">99.2%</span></div>
                <div class="inst-meter" role="progressbar" aria-label="7 days"
                     aria-valuenow="99" aria-valuemin="0" aria-valuemax="100">
                  <div class="inst-meter-fill" style="inline-size:99.2%"></div>
                </div>
              </div>
              <dl class="inst-kv">
                <dt>Last check</dt><dd>12 s ago</dd>
                <dt>Latency</dt><dd>38 ms</dd>
                <dt>Certificate</dt><dd>54 days left</dd>
              </dl>
            </div>
          </div>
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
| The framework | [The application shell](../layout/shell.md) and [the split](../layout/split.md) |
| The name and the actions | [The page header](../layout/page-header.md), [a badge](../components/display/badge.md), [a button](../components/actions/button.md) |
| How much of the day was up | [The share bar](../components/charts/share.md), one part per outcome |
| When it happened | [Check history](../agent/history.md) with groups, a ruler and a hint per hour |
| The key to the colours | [The legend](../components/charts/legend.md), states rather than series |
| The window figures | [A meter](../components/charts/meter.md) and [a list of properties](../components/display/kv.md) |
| The incidents | [A table](../components/display/table.md) with a badge in the reason |

## What is seen here

**The two bars answer different questions.** The strip is ordered and the share
bar is totalled: the same 3.1% is one long outage or ninety scattered ones, and
only the strip can tell them apart. Neither replaces the other, so both stand,
sharing one legend.

**The hour carries its own weight.** A group's width is its share of the
checks — 11:00 had one and is narrow, 17:00 had nine and is wide. Without it a
monitor backed by one check an hour looks exactly like one backed by a hundred.

**The ruler is not a second layout.** An axis cell takes the same `--n` as the
group above it, so the two rows cannot drift apart.

**The legend keys states, not series.** The swatch falls through to
`--tone-mark`, and the hour nobody checked takes `data-empty` — drawn as the
track, the same value the strip draws it with.

## Related

[Check history](../agent/history.md) [The share bar](../components/charts/share.md)
[The legend](../components/charts/legend.md) [The meter](../components/charts/meter.md)
[The application shell](../layout/shell.md) [The run console](./console.md)
