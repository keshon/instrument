---
title: Выбираемая карточка
group: Ввод
layout: component
source: src/forms.css
api:
  - { name: "inst-choice-card", kind: "класс", doc: "`<label>`-контейнер. Сетка, рамка, состояния" }
  - { name: "inst-choice-card-title", kind: "класс", doc: "Заголовок варианта" }
  - { name: "inst-choice-card-desc", kind: "класс", doc: "Объяснение последствий" }
  - { name: "--pad-card", kind: "токен" }
  - { name: "--radius-md", kind: "токен" }
  - { name: "--surface-raised", kind: "токен" }
  - { name: "--border-control", kind: "токен" }
  - { name: "--accent-border", kind: "токен" }
  - { name: "--accent-bg", kind: "токен" }
  - { name: "--accent-solid", kind: "токен" }
  - { name: "--size-marker", kind: "токен" }
  - { name: "--gap-inline", kind: "токен" }
  - { name: "--text-sm", kind: "токен" }
  - { name: "--text-xs", kind: "токен" }
  - { name: "--leading-ui", kind: "токен" }
  - { name: "--dur-1", kind: "токен" }
title-en: "Choice card"
group-en: "Inputs"
---

Радиокнопка или чекбокс, у которых есть заголовок и объяснение последствий.
Голой подписи в одну строку не хватает, когда у вариантов разная цена.

```html preview
<label class="inst-choice-card">
  <input type="radio" name="mode" checked>
  <span><span class="inst-choice-card-title">Параллельно</span>
    <span class="inst-choice-card-desc">До 12 агентов сразу. Быстрее, но расход токенов выше и порядок вывода не гарантирован.</span></span>
</label>
<label class="inst-choice-card">
  <input type="radio" name="mode">
  <span><span class="inst-choice-card-title">Последовательно</span>
    <span class="inst-choice-card-desc">По одному агенту. Медленнее, зато лог читается сверху вниз без чересполосицы.</span></span>
</label>
```

## Использование

```html
<label class="inst-choice-card">
  <input type="radio" name="plan">
  <span><span class="inst-choice-card-title">Заголовок</span>
    <span class="inst-choice-card-desc">Объяснение в одну-две строки.</span></span>
</label>
```

| Что | Обязательно | Почему |
|---|---|---|
| `<label>` как контейнер | да | Нажатие в любое место карточки попадает в контрол. Обёртка `<div>` превратила бы карточку в мёртвую площадь вокруг маленького кружка |
| `<input>` прямым ребёнком | да | Состояния написаны через `:has(> input…)`: контрол в подобёртке их не включит |
| Один `<span>` под текст | да | Сетка карточки двухколоночная: контрол и **один** текстовый блок. Два соседних узла разъедут колонки |
| `inst-choice-card-title` | да | Без заголовка это просто чекбокс в рамке |
| Общий `name` у радиогруппы | да | Без него варианты не исключают друг друга |
| `fieldset` с `legend` вокруг группы | да | Иначе вопрос, на который отвечают карточки, нигде не назван |
| `inst-choice-card-desc` | нет, но обычно да | Без объяснения карточка не оправдывает своей площади — берите [переключатели](./toggles.md) |

## Когда использовать

| Используйте | Возьмите другое |
|---|---|
| Выбор режима, у которого есть последствия и их надо объяснить | **Вариант, объяснимый одним словом** — [переключатели](./toggles.md): карточка вокруг слова «Да» это лишняя рамка |
| Два-четыре варианта, показанные все сразу | **Много вариантов или длинный список** — [селект](./select.md): двенадцать карточек это не выбор, а страница |
| Взаимоисключающий выбор одного из режимов | **Переключение вида одних и тех же данных** — [сегментированный контрол](../actions/segmented.md): он живёт в панели инструментов и не занимает высоту |
| Набор независимых опций с описаниями — чекбоксы внутри | **Включение и выключение одной настройки** — свитч в [переключателях](./toggles.md) |
| Карточка выбирается | **Карточка открывается или ведёт куда-то** — [карточка](../display/card.md) со ссылкой: у выбираемой нет перехода |

## Варианты

Внутри — **настоящий нативный контрол**, и от его типа зависит вся семантика
выбора. Кит не различает эти два случая ничем, кроме того, что рисует
браузер: правило `:has(> input:checked)` одинаково работает и там, и там.

```html preview
<label class="inst-choice-card">
  <input type="checkbox" checked>
  <span><span class="inst-choice-card-title">Собирать трассировку</span>
    <span class="inst-choice-card-desc">Полный лог каждого шага. Прогон становится примерно на 15% медленнее.</span></span>
</label>
<label class="inst-choice-card">
  <input type="checkbox">
  <span><span class="inst-choice-card-title">Оповестить по завершении</span>
    <span class="inst-choice-card-desc">Письмо на адрес учётной записи, когда очередь опустеет.</span></span>
</label>
```

| Тип | Когда | Обязательно |
|---|---|---|
| `type="radio"` | Один из взаимоисключающих режимов | Общий `name` у всей группы, иначе взаимоисключения не будет |
| `type="checkbox"` | Независимые опции, любое количество | `name` у каждого свой |

Группа радиокарточек — это группа полей, и у неё должно быть общее имя:
оберните её в `fieldset` с `legend`, иначе скринридер объявит четыре
несвязанных варианта без вопроса, на который они отвечают.

## Состояния

```html preview
<label class="inst-choice-card">
  <input type="radio" name="st" checked>
  <span><span class="inst-choice-card-title">Выбрано</span>
    <span class="inst-choice-card-desc">Рамка акцентом, тонированный фон и несущая метка у края.</span></span>
</label>
<label class="inst-choice-card">
  <input type="radio" name="st">
  <span><span class="inst-choice-card-title">Не выбрано</span>
    <span class="inst-choice-card-desc">Обычная рамка контрола.</span></span>
</label>
<label class="inst-choice-card">
  <input type="radio" name="st" disabled>
  <span><span class="inst-choice-card-title">Недоступно</span>
    <span class="inst-choice-card-desc">Недоступно на текущем тарифе.</span></span>
</label>
```

| Состояние | Как ставится | Что происходит |
|---|---|---|
| наведение | `:hover` | Рамка темнеет до `--text-muted` |
| выбрано | `:has(> input:checked)` | Рамка `--accent-border`, фон `--accent-bg` **и несущая метка у внутреннего края** |
| фокус | `:has(> input:focus-visible)` | Кольцо вокруг всей карточки, а не вокруг кружка |
| недоступно | `:has(> input:disabled)` | Прозрачность `0.5`, курсор `not-allowed` |

**Выбор несёт не только цвет.** У края выбранной карточки встаёт та же несущая
метка толщиной `--size-marker`, что у текущего пункта навигации и у активной
вкладки: один язык на весь кит, и он переживает и дальтонизм, и печать.

## Композиции

### С бейджем в заголовке

```html preview
<label class="inst-choice-card">
  <input type="radio" name="tier">
  <span><span class="inst-choice-card-title">Распределённо
      <span class="inst-badge" data-tone="warn">бета</span></span>
    <span class="inst-choice-card-desc">Агенты уезжают на пул воркеров. Пока без гарантий по времени.</span></span>
</label>
```

## Сценарии

### Группа с общим вопросом

```html preview context
<fieldset class="inst-fieldset">
  <legend>Режим прогона</legend>
  <div class="inst-grid">
    <label class="inst-choice-card">
      <input type="radio" name="run-mode" checked>
      <span><span class="inst-choice-card-title">Параллельно</span>
        <span class="inst-choice-card-desc">До 12 агентов сразу.</span></span>
    </label>
    <label class="inst-choice-card">
      <input type="radio" name="run-mode">
      <span><span class="inst-choice-card-title">Последовательно</span>
        <span class="inst-choice-card-desc">По одному агенту.</span></span>
    </label>
  </div>
</fieldset>
```

Группа полей идёт **без рамки** по умолчанию: карточки уже обведены, и рамка
вокруг них дала бы рамку внутри рамки. См. [форму](./form.md).

## Правила

:::do Карточка там, где есть последствия
Заголовок называет вариант, описание объясняет, чем за него платят.
:::

:::dont Карточка вокруг одного слова
Вариант, объяснимый словом «Да», — это [переключатель](./toggles.md). Рамка
вокруг него не добавляет смысла, только площадь.
:::

:::do fieldset с legend вокруг группы
Иначе скринридер объявит четыре несвязанных варианта без вопроса.
:::

:::dont Двенадцать карточек
Это уже не выбор, а страница. От пяти вариантов — [селект](./select.md).
:::

## Доступность

| | |
|---|---|
| Контрол | Нативный `<input>`, не `div` с `role`. Клавиатура, объявление состояния и участие в форме приходят от платформы |
| Клавиатура | `Tab` — вход в группу, стрелки — переключение между радиокнопками, `Space` — чекбокс. Ничего не перехватывается |
| Фокус | Кольцо вокруг **всей карточки**: у контрола внутри своё кольцо было бы вокруг кружка размером с точку |
| Имя | Заголовок и описание внутри `<label>` целиком становятся доступным именем варианта. Отдельный `aria-label` не нужен и вредит: он перекрыл бы описание |
| Группа | `fieldset` с `legend` — иначе варианты объявляются без вопроса |
| Цвет не единственный носитель | Выбор несёт метку у края и нативную отметку контрола, а не только фон и рамку |
| Цель нажатия | Вся карточка. Это самая большая цель во всём ките |
| Недоступность | `disabled` на `<input>` — нативная. Карточка гаснет целиком, и видно, **какой именно** вариант недоступен |

## API

```api
```

## Связанное

[Переключатели](./toggles.md) [Селект](./select.md) [Форма](./form.md)
[Сегментированный контрол](../actions/segmented.md)
[Карточка](../display/card.md) [Бейдж](../display/badge.md)
