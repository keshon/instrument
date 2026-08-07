---
title: Текстовое поле
group: Ввод
status: stable
source: src/forms.css
api:
  - { name: "inst-input", kind: "класс", doc: "Однострочное поле" }
  - { name: "inst-textarea", kind: "класс", doc: "Многострочное" }
  - { name: "inst-field", kind: "класс", doc: "Обёртка: подпись + контрол + подсказка" }
  - { name: "inst-label", kind: "класс", doc: "Подпись" }
  - { name: "inst-field-hint", kind: "класс", doc: "Подсказка" }
  - { name: "inst-field-error", kind: "класс", doc: "Сообщение об ошибке" }
  - { name: "inst-input-group", kind: "класс", doc: "Группа с приставками" }
  - { name: "inst-input-affix", kind: "класс", doc: "Приставка или суффикс" }
  - { name: "inst-required", kind: "класс", doc: "Пометка обязательности" }
  - { name: "inst-optional", kind: "класс", doc: "Пометка обязательности" }
  - { name: "inst-input--sm", kind: "модификатор", doc: "Размер" }
  - { name: "inst-input--lg", kind: "модификатор", doc: "Размер" }
  - { name: "inst-textarea--sm", kind: "модификатор", doc: "Размер" }
  - { name: "inst-textarea--lg", kind: "модификатор", doc: "Размер" }
  - { name: "inst-field--side", kind: "модификатор", doc: "Подпись сбоку вместо «сверху»" }
  - { name: "aria-invalid", kind: "атрибут", doc: "`true` — ошибка от приложения" }
  - { name: "data-when", kind: "атрибут", doc: "`invalid` — показом управляет платформа" }
  - { name: "--surface-field", kind: "токен" }
  - { name: "--control-h-sm/md/lg", kind: "токен" }
  - { name: "--control-pad-sm", kind: "токен" }
  - { name: "--radius-sm", kind: "токен" }
  - { name: "--text-xs/sm/md", kind: "токен" }
  - { name: "--label-col", kind: "токен" }
  - { name: "--border-control", kind: "токен" }
  - { name: "--err-text", kind: "токен" }
  - { name: "--leading-ui", kind: "токен" }
title-en: "Text field"
group-en: "Inputs"
---

Однострочный и многострочный ввод. Поле стоит во **врезе**
(`--surface-field`), а не на цвете панели: иначе его единственной границей
остаётся рамка, и поле приходится искать.

```html preview
<div class="inst-field">
  <label class="inst-label" for="agent">Имя агента</label>
  <input class="inst-input" id="agent" placeholder="worldgen-01">
  <span class="inst-field-hint">Латиница, цифры и дефис</span>
</div>
```

## Когда использовать

| Используйте | Возьмите другое |
|---|---|
| Произвольный текст: имя, путь, запрос | **Выбор из набора** — [селект](./select.md) или [сегментированный контрол](../actions/segmented.md) |
| Число с единицей и осью — `inst-num-field` | **Число, которое тянут мышью** — [числовое поле инспектора](./num-field.md) |
| Многострочный текст — `inst-textarea` | **Поиск с очисткой** — [поиск](./search.md): `input[type=search]` даёт нативный крестик |
| Значение с приставкой: валюта, единица | **Значение только для чтения** — `readonly`, а не `disabled`: см. ниже |

## Размеры

Симметричны кнопочным, потому что поле и кнопка стоят в одной строке формы.

```html preview
<input class="inst-input inst-input--sm" value="sm" aria-label="Поле sm">
<input class="inst-input" value="md" aria-label="Поле md">
<input class="inst-input inst-input--lg" value="lg" aria-label="Поле lg">
```

| Размер | Высота | Кегль |
|---|---|---|
| `inst-input--sm` | `--control-h-sm` | `--text-xs` |
| по умолчанию | `--control-h-md` | `--text-sm` |
| `inst-input--lg` | `--control-h-lg` | `--text-md` |

Горизонтальный отступ не меняется: у поля текст прижат к краю, а не
отцентрован, поэтому отступу незачем расти вместе с высотой.

У `inst-textarea` те же три размера, высота — два с половиной контрола.

## Состояния

```html preview
<input class="inst-input" value="Обычное" aria-label="Обычное">
<input class="inst-input" value="Только чтение" readonly aria-label="Только чтение">
<input class="inst-input" value="Недоступно" disabled aria-label="Недоступно">
<input class="inst-input" value="Ошибка" aria-invalid="true" aria-label="Ошибка">
```

| Состояние | Как ставится | Что происходит |
|---|---|---|
| наведение | `:hover` | Рамка темнеет. Уступает ошибке явно, см. ниже |
| фокус | `:focus-visible` | Кольцо плюс рамка акцентом |
| ошибка | `:user-invalid` **или** `aria-invalid="true"` | Красная рамка. Первое ставит платформа, второе — приложение |
| только чтение | `readonly` | Врез снимается, цвет текста **полный**: значение можно прочитать, выделить и скопировать |
| недоступно | `disabled` | Прозрачность `0.5` — как у кнопки |

### `readonly` — не `disabled`

Это разные вещи, и пока они выглядели одинаково, пользователь не мог понять,
ждать ему разблокировки или значение неизменяемо в принципе.

| | `readonly` | `disabled` |
|---|---|---|
| Можно прочитать и скопировать | да | нет |
| Попадает в фокус | да | нет |
| Уходит с формой на сервер | да | нет |
| Как выглядит | врез снят, текст полный | прозрачность 0.5 |

## Валидация

Ошибку показывает **платформа**, а не скрипт.

```html preview
<div class="inst-field">
  <label class="inst-label" for="mail">Почта<span class="inst-required">*</span></label>
  <input class="inst-input" id="mail" type="email" required
         aria-describedby="mail-msg">
  <span class="inst-field-hint" id="mail-msg">Пришлём отчёт о прогоне</span>
  <span class="inst-field-error" data-when="invalid">Проверьте адрес</span>
</div>
```

`:user-invalid`, а не `:invalid`. Разница принципиальная: `:invalid` горит
красным с первой секунды, ещё до того, как в поле что-то ввели, — форма
встречает пользователя россыпью ошибок, которых он не совершал.
`:user-invalid` срабатывает только после взаимодействия.

**Подсказка и ошибка занимают один слот.** Ошибка заменяет подсказку, а не
сдвигает форму: разметка не перекладывается в момент заполнения.

| Кто показывает | Как |
|---|---|
| Платформа | `data-when="invalid"` на `inst-field-error` — покажется по `:user-invalid` |
| Приложение | тот же блок **без атрибута** — например, по ответу сервера |

## Композиции

### Приставка и суффикс

```html
<div class="inst-input-group">
  <span class="inst-input-affix">₽</span>
  <input class="inst-input" aria-label="Бюджет">
  <span class="inst-input-affix">в час</span>
</div>
```

### Подпись сбоку

```html
<div class="inst-field inst-field--side">
  <label class="inst-label" for="t">Таймаут</label>
  <input class="inst-input" id="t" value="30">
</div>
```

Колонка подписей — тот же `--label-col`, что у инспектора: два способа
показать пару «имя — значение» не имеют права разъезжаться по метрике.

## Справочник

```api
```

## Доступность

| | |
|---|---|
| Подпись | `<label for>` обязателен. `placeholder` подписью **не является**: он исчезает при вводе и не читается частью скринридеров |
| Подсказка | `aria-describedby` на поле, указывающий на `inst-field-hint` |
| Ошибка | Сообщение текстом, а не только красной рамкой: цвет не имеет права быть единственным носителем |
| Обязательность | `required` плюс видимая пометка. Звёздочка без словесного объяснения в начале формы — половина контракта |
| Контраст | Рамка поля берёт `--border-control` и держит 3:1: здесь граница **и есть** контрол |
| Порядок состояний | Записан селектором: отключено > ошибка > фокус > наведение. У `:hover:not(:disabled)` специфичность (0,3,0), а у `:user-invalid` всего (0,2,0) — курсор над невалидным полем прятал красную рамку, пока наведение не начало уступать явно |

## Связанное

[Селект](./select.md) · [Переключатели](./toggles.md) · [Форма](./form.md) ·
[Поиск](./search.md) · [Числовое поле](./num-field.md) ·
[Кнопка](../actions/button.md)
