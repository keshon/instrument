# instrument

A CSS library for interfaces that show a person the work of a machine: queues
of tasks, runs, logs, inspectors, states and uncertainty.

**[instrument.keshon.ru](https://instrument.keshon.ru)** — the documentation
with live examples: every one of them can be handled on the page, in any of 300
combinations — five themes, four accents, five scales, three densities.

The library is meant for dashboards, admin panels and agent tools: screens
where the data is the point and where a person is watching a machine work.
Visual effect is deliberately limited — no gradients, no glows, no shadows on
cards; of the weights there are only regular and medium; headings in caps are
forbidden by rule. For landing pages, shops and content sites this is not the
right thing.

## Install

There is no build step and there are no dependencies. One file — 110 KB, 17 KB
gzipped: the tokens, five themes, five scales, three densities, every
component.

```bash
npm i @keshon/instrument
```

```html
<link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/@keshon/instrument@0.7.0/dist/instrument.min.css">
```

Name the version exactly: `@latest` is cached for a week and can change at an
arbitrary moment.

Icons want the sprite inside the document itself: `<use href="#i-close">` obeys
the same-origin rule, and a reference to an external file will not draw — with
no error in the console. The contents of `assets/sprite.svg` go in as the first
element of `<body>`; the full list of symbols is on the
[icons](docs/foundations/icons.md) page.

The whole procedure, including an entire page to copy, is in
[installation](docs/start/install.md).

## An example

```html
<div role="listbox" aria-label="The queue of agents">
  <div class="inst-task" role="option" data-state="running"
       aria-selected="true" tabindex="0">
    <span class="inst-task-gutter"><span class="inst-dot"></span></span>
    <span class="inst-task-main">
      <span class="inst-task-title inst-u-truncate">Parsing the height map</span>
      <span class="inst-task-sub">terrain/heightmap.ts · 4 files</span>
    </span>
    <span class="inst-task-meta">2.1 s</span>
  </div>
</div>
```

`role`, `aria-*` and `data-state` are **part of the markup of a component**
rather than an addition to it. `data-state` colours the row, `aria-selected`
marks the choice; without `role="listbox"` and `role="option"` the first is
invalid and the second is never announced. The obligatory markup of every
component is listed in the
[design principles](docs/about/design-principles.md#the-markup-contract).

An assembled screen, from which it can be seen how the components join up, is
the [run console](docs/blocks/console.md).

## What is inside

The package ships two files: `instrument.min.css` and `instrument.js`. They are
built from the sources in `src/`:

| The file in `src/` | What |
|---|---|
| `tokens.css` | the ramps, the semantics, the roles, five themes, four accents, five scales, three densities |
| `base.css` | the reset, element defaults, focus, scrollbars, prose |
| `layout.css` | the shell, the container, flow, the page header, navigation |
| `tone.css` | the vocabulary of meaningful tone and the box of a glyph — the rest read them |
| `actions.css` | the button, the segmented control, the chip |
| `surfaces.css` | the panel, the card, the metric, the toolbar |
| `status.css` | the badge, the dot, the spinner, the caret, the meter |
| `table.css` | the table, the list of pairs, the property row |
| `rows.css` | the queue row, the inserts |
| `text.css` | the log, code, the tree |
| `feedback.css` | the note, the empty state, the banner, the accordion |
| `forms.css` | the fields, validation, the layout of a form |
| `data.css` | the avatar, the tag, the timeline, the legend, the sparkline, the ring, the calendar, the change |
| `overlay.css` | the popover, the menu, the tooltip, the modal, the sheet, the toast, the cascader |
| `agent.css` | the agent layer: the approval, the step, the diff, the output, the budget, the lane, the failure |
| `forced.css` | `forced-colors` — the forced-colours mode of Windows |
| `motion.css` | `prefers-reduced-motion` |
| `print.css` | `@media print` |
| `kit.js` | the keyboard behaviour, optional — in the package this is `instrument.js` |

There is no need to include them one by one: `kit.css` declares the order of
the layers, and it is the same one that ends up in the built
`instrument.min.css`.

```
kit.tokens → kit.base → kit.layout → kit.components → kit.overlay
  → kit.agent → kit.forced → kit.motion → kit.print
```

The application's styles lie outside every layer and always win — with no
`!important`. More on that in [the tokens](docs/foundations/tokens.md#the-order-of-layers).

The package also carries `components.json` — the registry: for every component,
what happens to it in print, under `forced-colors` and on a page break, what
ARIA it declares and which components it stands beside. It is what the build
gates check the kit against, and an application can read it too.

## What it can do

**Theme, accent and density are attributes of the markup.** They work on any
subtree rather than on the document alone.

```html
<html data-theme="dark">            <!-- five themes; with no attribute — by the OS setting -->
<html data-accent="clay">           <!-- petrol (default) · graphite · indigo · clay -->
<section data-density="compact">    <!-- compact · comfortable -->
```

The four axes do not intersect: 5 themes × 4 accents × 5 scales × 3 densities.
Scale and density are both an attribute of a container rather than a variant of
a component, and what tells them apart is the type size: scale moves it along
with the geometry, density does not touch it at all. More on that in
[colour](docs/foundations/colors.md), [scale](docs/foundations/scale.md) and
[density](docs/foundations/density.md).

**Accessibility.** The contrast of token pairs is checked by machine — 2 014
checks across five themes and four accents; tap targets are checked against
WCAG 2.2 AA in 15 combinations of scale and density.
`prefers-reduced-motion`, `forced-colors` and `@media print` are supported. A
separate [audit by pixels](docs/about/audit.md) runs in the browser on your own
application.

**The keyboard.** The arrows, `Home`/`End` and the roving `tabindex` for a
list, a tree, a menu, a radiogroup and tabs are done by a separate file:

```html
<script type="module" src="dist/instrument.js"></script>
```

It is optional and sets only the attributes declared in the markup contract.
Without it everything that does not want a keyboard still works. More on that
in [behaviour](docs/foundations/behavior.md).

## Requirements

Chrome 123+ · Safari 17.5+ · Firefox 120+

Because of `light-dark()`, `oklch()`, cascade layers, nesting, `:has()` and the
`lh` unit. There are no polyfills.

## Documentation

| | |
|---|---|
| [instrument.keshon.ru](https://instrument.keshon.ru) | the reference with live examples |
| [docs/start/install.md](docs/start/install.md) | installation |
| [docs/about/design-principles.md](docs/about/design-principles.md) | the design principles: the rules the components are derived from |
| [CONTRIBUTING.md](CONTRIBUTING.md) | developing the library itself, the checks, the release |
| [ROADMAP.md](ROADMAP.md) | what is planned |
| [CHANGELOG.md](CHANGELOG.md) | what changed and what breaks on an update |

The reference is bilingual: the English pages live at the bare address, the
Russian ones under `/ru/`. In the sources a page is `button.md` and its Russian
half is `button.ru.md`.

## Licence

[MIT](LICENSE) © Innokentiy Sokolov
