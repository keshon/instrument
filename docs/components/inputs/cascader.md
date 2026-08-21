---
title: Cascader
group: Inputs
layout: component
source: src/overlay.css
js: The arrows along a column and between the columns are done by `instrument.js`
api:
  - { name: "inst-cascader", kind: "class", doc: "The row of columns. Stands inside an `.inst-popover`" }
  - { name: "inst-cascader-col", kind: "class", doc: "One level. A `role=\"listbox\"` with scrolling of its own" }
  - { name: "inst-cascader-item", kind: "class", doc: "An item of a level. Worn together with `.inst-menu-item`" }
  - { name: "inst-cascader-path", kind: "class", doc: "The chosen path inside the trigger. The separator is a chevron" }
  - { name: "aria-expanded", kind: "attribute", value: "true · false", doc: "On an item that HAS CHILDREN. The chevron follows it. A leaf carries no attribute" }
  - { name: "aria-selected", kind: "attribute", value: "true · false", doc: "On an item that lies on the path. In the last column that is the chosen value" }
  - { name: "inst:select", kind: "event", doc: "Bubbles from an item of a column. `{ value }` — the application renders the column after it" }
  - { name: "--control-h-sm", kind: "token" }
  - { name: "--size-chevron", kind: "token" }
  - { name: "--space-1", kind: "token" }
  - { name: "--space-2", kind: "token" }
  - { name: "--hairline", kind: "token" }
  - { name: "--border-subtle", kind: "token" }
  - { name: "--surface-selected", kind: "token" }
  - { name: "--accent-text", kind: "token" }
  - { name: "--accent-mark", kind: "token" }
  - { name: "--text-muted", kind: "token" }
  - { name: "--text-faint", kind: "token" }
  - { name: "--text-primary", kind: "token" }
---

Choosing one leaf of a hierarchy, level by level. The columns turn the depth
into width: three levels of forty items are three columns that scroll
separately rather than a hundred and twenty rows.

```html preview
<span class="inst-select-wrap">
  <button class="inst-select" type="button" popovertarget="cas-region"
          aria-haspopup="listbox" aria-label="Region">
    <span class="inst-cascader-path">
      <span>Europe</span><span>Germany</span><span>Berlin</span>
    </span>
  </button>
</span>

<div class="inst-popover inst-popover--anchored" id="cas-region" popover>
  <div class="inst-cascader">
    <div class="inst-cascader-col" role="listbox" aria-label="Part of the world">
      <button class="inst-menu-item inst-cascader-item" type="button" role="option"
              aria-selected="true" aria-expanded="true" tabindex="0">Europe</button>
      <button class="inst-menu-item inst-cascader-item" type="button" role="option"
              aria-selected="false" aria-expanded="false" tabindex="-1">Asia</button>
      <button class="inst-menu-item inst-cascader-item" type="button" role="option"
              aria-selected="false" aria-expanded="false" tabindex="-1">Africa</button>
    </div>
    <div class="inst-cascader-col" role="listbox" aria-label="Country">
      <button class="inst-menu-item inst-cascader-item" type="button" role="option"
              aria-selected="true" aria-expanded="true" tabindex="0">Germany</button>
      <button class="inst-menu-item inst-cascader-item" type="button" role="option"
              aria-selected="false" aria-expanded="false" tabindex="-1">France</button>
      <button class="inst-menu-item inst-cascader-item" type="button" role="option"
              aria-selected="false" aria-expanded="false" tabindex="-1">Portugal</button>
    </div>
    <div class="inst-cascader-col" role="listbox" aria-label="City">
      <button class="inst-menu-item inst-cascader-item" type="button" role="option"
              aria-selected="true" tabindex="0">Berlin</button>
      <button class="inst-menu-item inst-cascader-item" type="button" role="option"
              aria-selected="false" tabindex="-1">Hamburg</button>
      <button class="inst-menu-item inst-cascader-item" type="button" role="option"
              aria-selected="false" tabindex="-1">Munich</button>
    </div>
  </div>
</div>
```

## Contract

**The kit draws the columns, the application fills them.** What stands in the
column after this one is the children of the chosen node, that is data. No
selector computes it, and a kit that pretended to would be holding a copy of
somebody else's tree.

Above four levels the component stops working, and that is a property of it
rather than a matter of taste: four columns of 11rem is 44rem of horizontal
room, which a narrow screen does not have. There the case is a search over the
leaves rather than a walk down the levels.

| What | Required | Why |
|---|---|---|
| A `role="listbox"` on every column with an `aria-label` | yes | Each level is a group of its own with a name of its own: "Country" rather than a second nameless list |
| A `role="option"` on every item | yes | The carrier of the keyboard contract. Without it the arrows do nothing |
| An `aria-expanded` on an item **with children** | yes | The chevron follows it. A leaf carries no attribute, and the absence of a chevron is what says the road ends here |
| An `aria-selected` on what lies on the path | yes | In the last column that is the chosen value; in the ones before it, the road to it |
| A roving `tabindex` in every column | yes, the initial one | `0` on the selected one, `-1` on the rest. One Tab stop per column |
| A `type="button"` | yes | Otherwise inside a form the item will submit it |
| The trigger is a `<button class="inst-select">` | no, but usually | `.inst-select` is worn by a button when the native list is not enough. A class of its own would be a second field look |

### Accessibility

| | |
|---|---|
| The name of a column | An `aria-label` on the column: "Part of the world", "Country", "City". Three lists with no names are read as three lists |
| Traversal | `Tab` enters a column and goes to the next one. `↑` `↓` along a column, `←` `→` between the columns, `Home` and `End` to the ends |
| Selection | Follows the focus, as in any single-select listbox. Walking a column with the arrows chooses as it goes — that is the promise of the role rather than a decision of this component |
| The mark on the path | A surface plus a colour rather than colour alone: `--surface-selected` with `--accent-text`, the same carrier as a selected row of a queue |
| A branch and a leaf | Told apart by the chevron, and the chevron by `aria-expanded`. By ear the same difference is "expanded" against nothing |
| The tap target | `--control-h-sm` plus the gap between neighbours. Checked by `cmd/targets` |
| Scrolling | Per column. One scrollbar over all three would move six regions along with two hundred cities |

## Anatomy

| The part | What it does |
|---|---|
| `.inst-cascader` | The row of columns. A ceiling of ten rows in height, counted in `--control-h-sm` rather than in pixels, so it rides density |
| `.inst-cascader-col` | One level. 11rem wide — about twenty characters at the base scale; below that the column becomes a strip of ellipses, and choosing from ellipses is not choosing |
| `.inst-cascader-item` | Worn **together with** `.inst-menu-item`: the height, the radius and the hover are already right there, and a second set of them would drift |
| `.inst-cascader-path` | The path in the trigger. The separator is the chevron of [breadcrumbs](../navigation/breadcrumbs.md) — both say "this lies inside that", and one shape per meaning is the kit's icon rule |

The seam between columns is a hairline, the same one the shell puts between its
regions. That is not the dividing line the kit refuses inside a menu: that one
separates **items of one list**, this one separates lists — and without it
three columns read as one list that has wrapped.

The scrollbar gutter is reserved rather than left to appear: a column that
gains one on the fortieth item shifts the two beside it at the moment of a
click.

## States

```html preview
<div class="inst-cascader" style="max-block-size:none">
  <div class="inst-cascader-col" role="listbox" aria-label="States">
    <button class="inst-menu-item inst-cascader-item" type="button" role="option"
            aria-selected="true" aria-expanded="true" tabindex="0">On the path</button>
    <button class="inst-menu-item inst-cascader-item" type="button" role="option"
            aria-selected="false" aria-expanded="false" tabindex="-1">A branch</button>
    <button class="inst-menu-item inst-cascader-item" type="button" role="option"
            aria-selected="false" tabindex="-1">A leaf</button>
    <button class="inst-menu-item inst-cascader-item" type="button" role="option"
            aria-selected="false" aria-disabled="true" tabindex="-1">Unavailable</button>
  </div>
</div>
```

| State | How it is set | What happens |
|---|---|---|
| on the path | `aria-selected="true"` | `--surface-selected` and `--accent-text`. In the last column the same mark is the chosen value |
| a branch | `aria-expanded` | A chevron at the far edge. `true` and `false` are drawn alike: whether the column beside is open is said by the column beside |
| a leaf | no `aria-expanded` | No chevron. The absence is what says the road ends |
| unavailable | `aria-disabled="true"` | Transparency. The arrow still walks to it — **which** item is unavailable matters — and it cannot be acted on |

## JS

The module is included [once per page](../../foundations/behavior.md) — there is
no need to initialise the components one by one.

### What `instrument.js` does

A column is an ordinary `listbox`, and everything a listbox promises it gets
for free: one Tab stop, `↑` `↓` along it, `Home` and `End`, selection following
focus.

What a listbox has no notion of is the column **beside** it, and there the
cross-axis arrows are free: in a vertical group `←` and `→` mean nothing at
all. So `←` and `→` move between the columns — onto whatever the next column
already has selected, or onto its first item.

That behaviour is scoped by **class rather than by role**, and the distinction
is the point. Teaching every listbox in the kit to walk sideways would be
changing a role's promise for the sake of one component. A listbox that happens
to stand inside `.inst-cascader` gains a behaviour of the cascader instead, and
a listbox anywhere else answers the arrows exactly as it did.

### What the application does

```js
cascader.addEventListener('inst:select', (e) => {
  const col = e.target.closest('.inst-cascader-col');
  const level = [...cascader.children].indexOf(col);
  render(level + 1, childrenOf(e.detail.value));   // the column after it
});
```

One handler. The kit sets `aria-selected`, moves the roving `tabindex` and
sends the event; what the next column contains is known to the application
alone.

### What `instrument.js` does not do

It does not open the popover — that is the platform's `popovertarget` — and it
does not write the path into the trigger. The path is the value of the field,
and the value belongs to the application.

## API

```api
```

## Related

```related
```
