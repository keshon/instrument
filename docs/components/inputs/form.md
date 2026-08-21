---
title: Форма
group: Ввод
layout: component
source: src/forms.css
api:
  - { name: "inst-form", kind: "класс", doc: "Форма: вертикальный поток с шагом `--pad-panel`" }
  - { name: "inst-field", kind: "класс", doc: "Одно поле: подпись + контрол + подсказка" }
  - { name: "inst-form-actions", kind: "класс", doc: "Полоса действий" }
  - { name: "inst-form-actions-note", kind: "класс", doc: "Пояснение в полосе" }
  - { name: "inst-fieldset", kind: "класс", doc: "Группа полей" }
  - { name: "inst-required", kind: "класс", doc: "Пометки обязательности" }
  - { name: "inst-optional", kind: "класс", doc: "Пометки обязательности" }
  - { name: "inst-form--side", kind: "модификатор", doc: "Все подписи сбоку" }
  - { name: "inst-field--side", kind: "модификатор", doc: "Подпись сбоку у одного поля" }
  - { name: "inst-form-actions--end", kind: "модификатор", doc: "Действия к концу строки" }
  - { name: "inst-fieldset--framed", kind: "модификатор", doc: "С рамкой" }
  - { name: "--pad-panel", kind: "токен" }
  - { name: "--label-col", kind: "токен" }
  - { name: "--gap-inline", kind: "токен" }
  - { name: "--space-2", kind: "токен" }
  - { name: "--text-xs", kind: "токен" }
  - { name: "--err-text", kind: "токен" }
  - { name: "--text-muted", kind: "токен" }
  - { name: "--z-sticky", kind: "токен" }
title-en: "Form"
group-en: "Inputs"
---

Раскладка полей и полоса действий. Два варианта раскладки; выбор между ними
определяется характером значений.

```html preview
<form class="inst-form">
  <div class="inst-field">
    <label class="inst-label" for="f1">Имя прогона</label>
    <input class="inst-input" id="f1" placeholder="worldgen-01">
  </div>
  <div class="inst-field">
    <label class="inst-label" for="f2">Модель<span class="inst-optional">необязательно</span></label>
    <span class="inst-select-wrap">
      <select class="inst-select" id="f2"><option>opus</option><option>sonnet</option></select>
    </span>
  </div>
  <div class="inst-form-actions">
    <button class="inst-btn inst-btn--primary" type="submit">Запустить</button>
    <button class="inst-btn" type="button">Отмена</button>
  </div>
</form>
```

## Контракт

Настоящий `<form>`, поля в `inst-field`, действия в `inst-form-actions`.

| Что | Обязательно | Почему |
|---|---|---|
| Настоящий `<form>` | да | `Enter` в поле отправляет форму. Своя кнопка на `div` это ломает |
| `<label for>` у каждого контрола | да | Форма без единого `for` проходит визуальную проверку и проваливает любую другую |
| `type="submit"` у главного действия | да | Иначе отправки не происходит, а `Enter` не работает |
| `<fieldset>` с `<legend>` вокруг группы | нет, но обычно да | Для радиокнопок — единственный способ озвучить вопрос, а не только варианты |

### Обязательность

```html preview
<div class="inst-field">
  <label class="inst-label" for="r1">Имя<span class="inst-required">*</span></label>
  <input class="inst-input" id="r1" required>
</div>
<div class="inst-field">
  <label class="inst-label" for="r2">Комментарий<span class="inst-optional">необязательно</span></label>
  <input class="inst-input" id="r2">
</div>
```

Обязательность несёт **и знак, и слово**: звёздочка в одиночку работает только
для тех, кто уже знает соглашение. В длинных формах помечать выгоднее
необязательные — их обычно меньше.

Настоящий `<form>` с `type="submit"`: `Enter` в поле отправляет форму, и это
ожидаемое поведение, которое своя кнопка ломает. Каждый контрол связан с
`<label for>` — форма без единого `for` проходит визуальную проверку и
проваливает любую другую. Группе имя дают `<fieldset>` и `<legend>`; для
радиокнопок это единственный способ озвучить вопрос, а не только варианты.

Обязательность требует двух вещей сразу: `required` на контроле и видимой
пометки. Атрибут без пометки виден только скринридеру, пометка без атрибута —
только глазом.

Порядок обхода совпадает с порядком в разметке. Визуальная перестановка колонок
без перестановки узлов ломает `Tab`.

## Устройство

### Полоса действий

```html preview
<div class="inst-form-actions inst-form-actions--end">
  <span class="inst-form-actions-note">Изменения применятся сразу</span>
  <button class="inst-btn" type="button">Отмена</button>
  <button class="inst-btn inst-btn--primary" type="submit">Сохранить</button>
</div>
```

В рамке — `inst-form-actions--end`; умолчание с действиями у начала строки
стоит в примере формы выше.

| Класс | Работа |
|---|---|
| `inst-form-actions` | Полоса действий: черта, воздух над ней, действия в ряд |
| `inst-form-actions--end` | Действия прижаты к концу строки |
| `inst-form-actions-note` | Пояснение рядом с действиями |

Полоса **закреплена** у нижнего края области прокрутки — и закрепляется за
вычетом её нижнего отступа, `--scroll-pad-end` (см.
[оболочку](../../layout/shell.md)). Без вычета она останавливается выше
видимого низа, и под ней остаётся щель в высоту отступа, сквозь которую видно
проезжающее содержимое: полоса перестаёт быть дном и делается плашкой посреди
чужих полей.

Внутри области полоса дополнительно гасит её внутренний отступ по бокам и
снизу — доходит до краёв и берёт поверхность панели вместо поверхности
страницы. Практический вывод: **полоса живёт там же, где поля, которые она
сохраняет.**

Главное действие стоит у **внешнего края полосы**: первым, когда полоса
выровнена к началу строки, и последним, когда её прижимает к концу `--end`
или пояснение. В форме, где кнопок две, это положение — единственный сигнал,
какая из них главная, и он остаётся тем же при любом выравнивании.

### Группа полей

```html preview
<fieldset class="inst-fieldset inst-fieldset--framed inst-fieldset--side">
  <legend>Ограничения прогона</legend>
  <div class="inst-field">
    <label class="inst-label" for="lim-mem">Память</label>
    <input class="inst-input" id="lim-mem" value="4 ГиБ">
  </div>
  <div class="inst-field">
    <label class="inst-label" for="lim-time">Время</label>
    <input class="inst-input" id="lim-time" value="30 мин">
  </div>
</fieldset>
```

`inst-fieldset` **без рамки** по умолчанию, и это осознанно: работа fieldset —
сгруппировать поля и дать группе имя, а рамка приходит из браузера, а не из
смысла. Форма почти всегда живёт внутри панели, и рамка группы дала бы рамку
внутри рамки. Границу задаёт заголовок и воздух.

Рамка возвращается модификатором `inst-fieldset--framed`, когда группа стоит
сама по себе.

## Варианты

Подпись сверху или сбоку — выбор по типу значений.

| | Когда |
|---|---|
| **Сверху** (умолчание) | Поля разнородны и требуют объяснения. Подпись читается прямо перед контролом |
| **Сбоку** (`inst-form--side`) | Поля коротки и однотипны. Колонка подписей превращает форму в таблицу, которую можно просканировать, не читая каждую строку |

```html preview
<form class="inst-form inst-form--side">
  <div class="inst-field">
    <label class="inst-label" for="s1">Таймаут</label>
    <input class="inst-input" id="s1" value="30">
  </div>
  <div class="inst-field">
    <label class="inst-label" for="s2">Попыток</label>
    <input class="inst-input" id="s2" value="3">
  </div>
</form>
```

Модификатор ставится на форму целиком или на отдельное поле
(`inst-field--side`). Колонка подписей — тот же `--label-col`, что у
инспектора: два способа показать пару «имя — значение» не имеют права
разъезжаться по метрике.

## Сценарии

### Форма внутри панели

Форма редко живёт сама по себе — она стоит внутри панели, и половина её
поведения видна только там. Прокрутите тело панели: **полоса действий
прилипает к нижнему краю**, а не уезжает вместе с полями. В длинной форме
кнопка «Сохранить», до которой надо доскроллить, — это кнопка, которую не
нажмут.

```html preview context
<div class="inst-panel">
  <div class="inst-panel-header"><span class="inst-panel-title">Уведомления</span></div>
  <div class="inst-panel-body">
    <div class="inst-form">
      <div class="inst-field">
        <span class="inst-label">Импорт списка получателей</span>
        <label class="inst-file">
          <input type="file" multiple>
          Перетащите файлы или выберите
          <span class="inst-file-hint">До 20 МБ, форматы .json и .csv</span>
        </label>
      </div>

      <div class="inst-field">
        <label class="inst-label" for="f13">Кого оповещать</label>
        <span class="inst-select-wrap">
          <select class="inst-select" id="f13" multiple>
            <option selected>Дежурный инженер</option>
            <option selected>Владелец проекта</option>
            <option>Вся команда</option>
          </select>
        </span>
        <span class="inst-field-hint">Шеврон не рисуется — список ничего не раскрывает</span>
      </div>

      <fieldset class="inst-fieldset">
        <legend>Когда оповещать</legend>
        <div class="inst-cluster inst-cluster--loose">
          <label class="inst-checkbox"><input type="checkbox" checked>При падении</label>
          <label class="inst-checkbox"><input type="checkbox">При каждом прогоне</label>
          <label class="inst-switch"><input type="checkbox" checked>Дублировать в почту</label>
        </div>
      </fieldset>

    </div>

    <div class="inst-form-actions">
      <span class="inst-form-actions-note">Изменения не сохранены</span>
      <button class="inst-btn" type="button">Отмена</button>
      <button class="inst-btn inst-btn--primary" type="button">Сохранить</button>
    </div>
  </div>
</div>
```

## API

```api
```

## Связанное

```related
```
