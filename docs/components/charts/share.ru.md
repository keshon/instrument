---
title: Полоса состава
group: Графики
layout: component
source: src/data.css
api:
  - { name: "inst-share", kind: "класс", doc: "Полоса. Носитель роли и имени" }
  - { name: "inst-share-part", kind: "класс", doc: "Один исход. Его ширина — его доля целого" }
  - { name: "inst-share--lg", kind: "класс", doc: "Высота глифа — для полосы, которая на экране предмет, а не примечание к числу" }
  - { name: "--share", kind: "токен", doc: "Доля части, от 0 до 1. Это **значение, а не оформление**, и приходит оно из данных" }
  - { name: "data-tone", kind: "атрибут", value: "neutral · running · ok · warn · error", doc: "Исход части. Словарь закрыт. Без атрибута — акцент" }
  - { name: "role", kind: "атрибут", value: "img", doc: "Полоса — картинка данных, и у неё обязано быть имя" }
  - { name: "--size-meter", kind: "токен" }
  - { name: "--size-icon", kind: "токен" }
  - { name: "--radius-mark", kind: "токен" }
  - { name: "--track", kind: "токен" }
  - { name: "--hairline", kind: "токен" }
  - { name: "--border", kind: "токен" }
  - { name: "--tone-mark", kind: "токен" }
  - { name: "--accent-mark", kind: "токен" }
---

Из чего период **состоял**: двадцать два часа наверху, полтора внизу, полчаса
никто не смотрел. Одна полоса, поделённая по доле каждого исхода в ней.

```html preview
<div class="inst-share" role="img"
     aria-label="Последние 24 часа: 92% наверху, 6% внизу, 2% не проверяли">
  <span class="inst-share-part" data-tone="ok" style="--share:0.92"></span>
  <span class="inst-share-part" data-tone="error" style="--share:0.06"></span>
  <span class="inst-share-part" data-tone="neutral" style="--share:0.02"></span>
</div>
```

## Контракт

| Что | Обязательно | Почему |
|---|---|---|
| `role="img"` и `aria-label` с долями **словами** | да | Три пустых `<span>` для скринридера — это три пустых `<span>`. Показание обязано быть текстом: «92% наверху, 6% внизу» |
| `--share` у каждой части | да | Без доли часть получает коэффициент роста 0 и не рисуется. Число — это данные, и приходит оно с сервера |
| Тон из словаря | да | `neutral` · `running` · `ok` · `warn` · `error`. Шестого исхода нет |
| Сумма долей, равная **1** | нет | Неучтённое остаётся дорожкой, и полоса, заполненная на треть, — честный рисунок восьми часов данных из двадцати четырёх |
| Минимальная ширина маленькой доли | **нет, и это отказ** | Порог нарисовал бы 0,02% такой же ширины, что и 2%, и полоса перестала бы сообщать пропорцию — в единицах полосы, которая заявляет, что сообщает пропорцию |

`style="--share:0.92"` — это **значение, а не оформление**, по тому же правилу,
что и ширина заливки меры: величина приходит с сервера и физически не может
жить в статическом CSS. Всё остальное — высота, угол, дорожка, цвет — из
библиотеки.

Протяжённость рисует **внутренний волосок**, как у меры: дорожка остаётся
достаточно светлой, чтобы тон на ней держал 3:1.

Части **прямоугольные внутри обрезанной полосы**. Скругление каждой части
поставило бы дугу на каждой внутренней границе — четыре угла там, где у данных
один стык, — и на высоте `--size-meter` эти дуги смыкаются, оставляя между
двумя соседними исходами зазубрину дорожки.

### Доступность

| | |
|---|---|
| Имя | `role="img"` плюс `aria-label` с долями словами. Части пустые и не читаются; читается имя полосы |
| Цвет не единственный носитель | Полоса стоит рядом с числами — в легенде или в строке текста. Сама по себе она говорит «в основном одно»; какое именно — дело подписи |
| Контраст | Часть берёт `--tone-mark` — порог 3:1 к поверхности, как у пометки состояния. Дорожка не несущая, и 3:1 с неё не спрашивают |
| Движение | Его нет. Состав — состояние, а не процесс, и при обновлении данных полоса не анимируется |
| Режим forced-colors | Полоса получает рамку `CanvasText`, а части переживают сброс, потому что несут значение. Тона там схлопываются в одну заливку, поэтому читаемым остаётся **деление**: каждая часть, кроме последней, рисует линию по своему краю |
| Печать | Печатается как есть: части различаются тоном и остаются серыми ступенями на бумаге |

## Размеры

```html preview
<div class="inst-share" role="img" aria-label="Сегодня: 88% наверху, 9% с деградацией, 3% внизу">
  <span class="inst-share-part" data-tone="ok" style="--share:0.88"></span>
  <span class="inst-share-part" data-tone="warn" style="--share:0.09"></span>
  <span class="inst-share-part" data-tone="error" style="--share:0.03"></span>
</div>
<div class="inst-share inst-share--lg" role="img" aria-label="Сегодня: 88% наверху, 9% с деградацией, 3% внизу">
  <span class="inst-share-part" data-tone="ok" style="--share:0.88"></span>
  <span class="inst-share-part" data-tone="warn" style="--share:0.09"></span>
  <span class="inst-share-part" data-tone="error" style="--share:0.03"></span>
</div>
```

Два, по числу мест, где полоса стоит. Базовый — высота меры: в строке данных
это работа меры с более чем одним ответом. Крупный — высота глифа: на экране
**одного** объекта полоса становится предметом, а не примечанием к числу.

## Композиции

### С легендой

```html preview
<div class="inst-stack">
  <div class="inst-share" role="img"
       aria-label="Последние 24 часа: 92% наверху, 6% внизу, 2% не проверяли">
    <span class="inst-share-part" data-tone="ok" style="--share:0.92"></span>
    <span class="inst-share-part" data-tone="error" style="--share:0.06"></span>
    <span class="inst-share-part" data-tone="neutral" style="--share:0.02"></span>
  </div>
  <ul class="inst-legend">
    <li class="inst-legend-item" data-tone="ok"><span class="inst-legend-swatch"></span>наверху <span class="inst-legend-value">92%</span></li>
    <li class="inst-legend-item" data-tone="error"><span class="inst-legend-swatch"></span>внизу <span class="inst-legend-value">6%</span></li>
    <li class="inst-legend-item"><span class="inst-legend-swatch" data-empty></span>не проверяли <span class="inst-legend-value">2%</span></li>
  </ul>
</div>
```

Образец читает `--tone-mark`, когда `--series` не задан, поэтому **одна легенда
расшифровывает и то и другое**: и график рядов, и полосу состояний. Час, в
который никто не смотрел, — ни то ни другое: он берёт `data-empty`, то же
слово, каким это называет `inst-history-tick`, и рисуется дорожкой — потому что
ключ обязан быть нарисован тем же значением, что и полоса.

### Над полосой истории

```html preview context
<div class="inst-stack">
  <div class="inst-share" role="img" aria-label="Последние 12 часов: 92% наверху, 8% внизу">
    <span class="inst-share-part" data-tone="ok" style="--share:0.92"></span>
    <span class="inst-share-part" data-tone="error" style="--share:0.08"></span>
  </div>
  <div class="inst-history" role="img" aria-label="12 часов по 4 проверки: 44 успешных, 4 с ошибкой">
    <span class="inst-history-group" style="--n:4"><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span></span>
    <span class="inst-history-group" style="--n:4"><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span></span>
    <span class="inst-history-group" style="--n:4"><span class="inst-history-tick" data-tone="error"></span><span class="inst-history-tick" data-tone="error"></span><span class="inst-history-tick" data-tone="error"></span><span class="inst-history-tick" data-tone="error"></span></span>
    <span class="inst-history-group" style="--n:4"><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span></span>
    <span class="inst-history-group" style="--n:4"><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span></span>
    <span class="inst-history-group" style="--n:4"><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span></span>
    <span class="inst-history-group" style="--n:4"><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span></span>
    <span class="inst-history-group" style="--n:4"><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span></span>
    <span class="inst-history-group" style="--n:4"><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span></span>
    <span class="inst-history-group" style="--n:4"><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span></span>
    <span class="inst-history-group" style="--n:4"><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span></span>
    <span class="inst-history-group" style="--n:4"><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span></span>
  </div>
</div>
```

Эти двое отвечают на разные вопросы, и экрану аптайма нужны оба. Полоса истории
отвечает на **когда**: провал был в 03:00, их было четыре, они шли подряд.
Полоса состава отвечает на **сколько**: день был на 8% внизу. История не умеет
подводить итог, состав не умеет в порядок: те же 8% — это один долгий провал
или девяносто разбросанных, и эта полоса нарисует их одинаково.

## API

```api
```

## Связанное

```related
```
