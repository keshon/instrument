---
title: Inserts
group: Inputs
layout: component
source: src/rows.css
js: Inserting into a field on a press is done by `instrument.js`
api:
  - { name: "inst-inserts", kind: "class", doc: "A row of inserts" }
  - { name: "inst-insert", kind: "class", doc: "One insert. A real `<button>`" }
  - { name: "data-insert-into", kind: "attribute", doc: "The selector of the field to put the value into. Without it the field of the same `inst-field` is taken" }
  - { name: "data-insert", kind: "attribute", doc: "What to insert, when it is not the text of the button itself" }
  - { name: "inst:insert", kind: "event", doc: "An insert was pressed. `detail` is `{ text, field }`. Cancelling leaves the field untouched" }
  - { name: "--control-h-sm", kind: "token" }
  - { name: "--surface-sunken", kind: "token" }
  - { name: "--border-subtle", kind: "token" }
  - { name: "--font-mono", kind: "token" }
  - { name: "--text-xs", kind: "token" }
  - { name: "--space-2", kind: "token" }
  - { name: "--space-3", kind: "token" }
  - { name: "--radius-sm", kind: "token" }
---

A row of inserts for the field beside them: the variables of a template, the
keys of an event, the fields of a record. A press puts the value **where the
caret is**.

```html preview
<div class="inst-field" style="max-inline-size:30rem">
  <label class="inst-label" for="ins-tpl">Notification template</label>
  <textarea class="inst-textarea inst-u-mono" id="ins-tpl" rows="3">{{name}} has failed
</textarea>
  <div class="inst-inserts">
    <button class="inst-insert" type="button" data-insert-into="#ins-tpl" title="The name of the monitor">{{name}}</button>
    <button class="inst-insert" type="button" data-insert-into="#ins-tpl" title="The address checked">{{target}}</button>
    <button class="inst-insert" type="button" data-insert-into="#ins-tpl" title="The state">{{status}}</button>
    <button class="inst-insert" type="button" data-insert-into="#ins-tpl" title="The text of the error">{{error}}</button>
    <button class="inst-insert" type="button" data-insert-into="#ins-tpl" title="The time of the event">{{time}}</button>
  </div>
</div>
```

The example is live: put the cursor in the field and press an insert — it lands
where the cursor stood rather than at the end.

## Contract

| What | Required | Why |
|---|---|---|
| A real `<button type="button">` | yes | This is an action. A `<span>` takes no focus, does not fire on `Enter` and is not called a button |
| `data-insert-into` | no | Without it the field of the same `inst-field` is taken. It is wanted when the field lies elsewhere in the markup |
| Every insert visible at once | yes | A couple of rows is taken in at a glance. What did not fit and has to be searched for is a reference of the format already, and its place is in documentation rather than under a field |

The text of the button is the value — a separate label is not needed. An
explanation goes into the `title`, but it **adds to** the name rather than
replacing it.

A row of ten buttons is traversed by tab, and that is right: each does its own
thing. After an insert the focus returns to the field — left on the button, it
means the field will have to be found by hand.

The label holds 4.5:1 on a recessed surface in all five themes. The monospaced
setting is not decoration here: an insert is an identifier rather than a word,
and `{{name}}` differs from `{{node}}` character by character.

## Variants

There are none. An insert is quieter than an action: there can be a dozen and
more of them to one field, and at full voice the palette would shout louder
than what the page was opened for. If an insert matters more than the text
under it, it is not an insert.

## JS

Include the module once per page.

```html
<script type="module" src="instrument.js"></script>
```

### What `instrument.js` does

It puts the value into the field **where the caret is** rather than at the end:
the cursor is placed where the insert is wanted, and that is where it is
expected. A selection is replaced. After an insert the focus returns to the
field — otherwise there is nowhere to type the next letter.

The value is taken from `data-insert`, and without one from the text of the
button.

### Events

`inst:insert` bubbles from the button, and its `detail` is `{ text, field }`.
Cancelling leaves the field untouched: an application may want a logic of
insertion of its own, with a check of the position, for instance.

```js
form.addEventListener('inst:insert', (e) => {
  if (!allowed.has(e.detail.text)) e.preventDefault();
});
```

The field itself sends a native `input` after an insert — frameworks see the
change with not one line of glue.

### Options

| Attribute | What it does |
|---|---|
| `data-insert-into` | The selector of the field. Without it the first `textarea`/`input` of the same `.inst-field` is taken |
| `data-insert` | What to insert, when it is not the text of the button itself: a short label, a long value |

## API

```api
```

## Related

```related
```
