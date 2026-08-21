---
title: Calendar
group: Data display
layout: component
source: src/data.css
js: Picking a date and switching months belong to the application layer
api:
  - { name: "inst-calendar", kind: "class", doc: "A grid of seven columns" }
  - { name: "inst-calendar-dow", kind: "class", doc: "The head: the abbreviation of a weekday" }
  - { name: "inst-calendar-day", kind: "class", doc: "The cell of a day" }
  - { name: "data-today", kind: "attribute", value: "true", doc: "on `inst-calendar-day`" }
  - { name: "data-outside", kind: "attribute", value: "true", doc: "in the same place" }
  - { name: "aria-selected", kind: "attribute", value: "true · false", doc: "in the same place" }
  - { name: "aria-disabled", kind: "attribute", value: "true", doc: "in the same place" }
  - { name: "--space-1", kind: "token" }
  - { name: "--space-2", kind: "token" }
  - { name: "--control-h-sm", kind: "token" }
  - { name: "--radius-sm", kind: "token" }
  - { name: "--text-xs", kind: "token" }
  - { name: "--text-2xs", kind: "token" }
  - { name: "--text-muted", kind: "token" }
  - { name: "--surface-hover", kind: "token" }
  - { name: "--surface-selected", kind: "token" }
  - { name: "--accent-text", kind: "token" }
  - { name: "--accent-solid", kind: "token" }
  - { name: "--size-marker", kind: "token" }
  - { name: "--weight-medium", kind: "token" }
---

A static grid of dates. The library answers for the grid being a grid and for
today being marked by more than colour.

```html preview
<div class="inst-calendar" role="grid" aria-label="August 2026">
  <span class="inst-calendar-dow" aria-hidden="true">Mo</span>
  <span class="inst-calendar-dow" aria-hidden="true">Tu</span>
  <span class="inst-calendar-dow" aria-hidden="true">We</span>
  <span class="inst-calendar-dow" aria-hidden="true">Th</span>
  <span class="inst-calendar-dow" aria-hidden="true">Fr</span>
  <span class="inst-calendar-dow" aria-hidden="true">Sa</span>
  <span class="inst-calendar-dow" aria-hidden="true">Su</span>
  <span class="inst-calendar-day" role="gridcell" data-outside="true">27</span>
  <span class="inst-calendar-day" role="gridcell" data-outside="true">28</span>
  <span class="inst-calendar-day" role="gridcell" data-outside="true">29</span>
  <span class="inst-calendar-day" role="gridcell" data-outside="true">30</span>
  <span class="inst-calendar-day" role="gridcell" data-outside="true">31</span>
  <span class="inst-calendar-day" role="gridcell">1</span>
  <span class="inst-calendar-day" role="gridcell">2</span>
  <span class="inst-calendar-day" role="gridcell">3</span>
  <span class="inst-calendar-day" role="gridcell">4</span>
  <span class="inst-calendar-day" role="gridcell">5</span>
  <span class="inst-calendar-day" role="gridcell">6</span>
  <span class="inst-calendar-day" role="gridcell">7</span>
  <span class="inst-calendar-day" role="gridcell">8</span>
  <span class="inst-calendar-day" role="gridcell">9</span>
  <span class="inst-calendar-day" role="gridcell">10</span>
  <span class="inst-calendar-day" role="gridcell">11</span>
  <span class="inst-calendar-day" role="gridcell" data-today="true">12</span>
  <span class="inst-calendar-day" role="gridcell">13</span>
  <span class="inst-calendar-day" role="gridcell">14</span>
  <span class="inst-calendar-day" role="gridcell">15</span>
  <span class="inst-calendar-day" role="gridcell">16</span>
  <span class="inst-calendar-day" role="gridcell">17</span>
  <span class="inst-calendar-day" role="gridcell" aria-selected="true">18</span>
  <span class="inst-calendar-day" role="gridcell">19</span>
  <span class="inst-calendar-day" role="gridcell">20</span>
  <span class="inst-calendar-day" role="gridcell">21</span>
  <span class="inst-calendar-day" role="gridcell" aria-disabled="true">22</span>
  <span class="inst-calendar-day" role="gridcell" aria-disabled="true">23</span>
  <span class="inst-calendar-day" role="gridcell">24</span>
  <span class="inst-calendar-day" role="gridcell">25</span>
  <span class="inst-calendar-day" role="gridcell">26</span>
  <span class="inst-calendar-day" role="gridcell">27</span>
  <span class="inst-calendar-day" role="gridcell">28</span>
  <span class="inst-calendar-day" role="gridcell" aria-disabled="true">29</span>
  <span class="inst-calendar-day" role="gridcell" aria-disabled="true">30</span>
  <span class="inst-calendar-day" role="gridcell">31</span>
  <span class="inst-calendar-day" role="gridcell" data-outside="true">1</span>
  <span class="inst-calendar-day" role="gridcell" data-outside="true">2</span>
  <span class="inst-calendar-day" role="gridcell" data-outside="true">3</span>
  <span class="inst-calendar-day" role="gridcell" data-outside="true">4</span>
  <span class="inst-calendar-day" role="gridcell" data-outside="true">5</span>
  <span class="inst-calendar-day" role="gridcell" data-outside="true">6</span>
</div>
```

## Contract

| What | Required | Why |
|---|---|---|
| `role="grid"` on the grid, `role="gridcell"` on a day | yes | Without the roles a calendar is spoken as a stream of 42 numbers |
| An `aria-label` with the month and the year on the grid | yes | "2" with no month is not a date |
| An `aria-label` with the full date on every day | yes | One number out of context is useless |
| `aria-hidden="true"` on the weekdays | yes | "Mo", "Tu" repeat what the full date of the cell already carries |
| Today marked by weight rather than by colour alone | yes | An underline and a semibold are two marks besides colour. The same is told to a screen reader in words in the `aria-label` |

### Accessibility

| | |
|---|---|
| The roles | `role="grid"` on the grid, `role="gridcell"` on a day, `role="row"` on a week. Without them a calendar is spoken as a stream of 42 numbers |
| The name | An `aria-label` with the month and the year on the grid. "2" with no month is not a date |
| The weekdays | `aria-hidden="true"`: the abbreviations "Mo", "Tu" repeat what a screen reader takes from the full date of the cell |
| The full date | Every day carries an `aria-label` with the full date ("2 August 2026, Saturday"). One number out of context is useless |
| The keyboard | Arrows over days, `PageUp`/`PageDown` over months, `Home`/`End` over a week, a roving `tabindex`. The behaviour belongs to the application |
| Today | Not colour alone: an underline plus the weight. A screen reader has to be told this in words in the `aria-label` |
| The tap target | The height is `--control-h-sm`, that is 26px. At density `compact` it is 22px, **below the 24px of WCAG 2.5.8** |
| Tabular figures | On, so the columns of numbers stand level whatever the number of digits |

## States

| Attribute | What it means |
|---|---|
| no attributes | An ordinary day of the current month |
| `data-today="true"` | Today. An underline **and** a semibold weight |
| `data-outside="true"` | A day of a neighbouring month filling out the week |
| `aria-selected="true"` | Selected |
| `aria-disabled="true"` | Unavailable. An opacity of `0.5`, with the press removed |

**Today is marked by an underline rather than by colour alone.** That is the
library's law being carried out: colour has no right to be the only carrier of
a state. The underline is drawn by an inner shadow of thickness `--size-marker`
— the same that carries the bar of an active tab and the edge of the current
navigation item.

## JS

Include the module once per page — there is no need to initialise the
components one by one, `instrument.js` works by delegation and sees nodes that
arrived later.

```html
<script type="module" src="instrument.js"></script>
```

### What `instrument.js` does

Nothing. A calendar is a grid of dates, and the library does not compute dates:
the month, the leap year, the first day of the week and the locale are data
rather than styling. The library draws the grid and the states of a day, and
the application builds the markup.

### What the application has to do

| What | What the library has for it |
|---|---|
| Build the grid of a month | The classes of a day and an `aria-label` with the full date |
| Mark what is selected | `aria-selected="true"` on a day |
| Switch the month | The buttons of the header — ordinary `<button>`s |

```js
grid.addEventListener('click', (e) => {
  const day = e.target.closest('[data-date]');
  if (!day || day.getAttribute('aria-disabled') === 'true') return;

  for (const d of grid.querySelectorAll('[aria-selected="true"]')) {
    d.setAttribute('aria-selected', 'false');
  }
  day.setAttribute('aria-selected', 'true');
  onPick(day.dataset.date);
});
```

The date is kept in `data-date` in ISO format rather than read from the label:
the label depends on the locale, and parsing it back is a way of getting a
divergence on the very first language.

## API

```api
```

## Related

```related
```
