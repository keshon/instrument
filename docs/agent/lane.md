---
title: Дорожки прогонов
group: Агентный слой
layout: component
source: src/agent.css
api:
  - { name: "inst-lanes", kind: "класс", doc: "Набор дорожек" }
  - { name: "inst-lane", kind: "класс", doc: "Одна дорожка: подпись и ось" }
  - { name: "inst-lane-label", kind: "класс", doc: "Имя исполнителя, обрезается" }
  - { name: "inst-lane-track", kind: "класс", doc: "Ось времени" }
  - { name: "inst-lane-span", kind: "класс", doc: "Отрезок работы" }
  - { name: "data-tone", kind: "атрибут", doc: "`neutral` `running` `ok` `warn` `error`" }
  - { name: "--label-col", kind: "токен" }
  - { name: "--gap-inline", kind: "токен" }
  - { name: "--size-meter", kind: "токен" }
  - { name: "--radius-xs", kind: "токен" }
  - { name: "--track", kind: "токен" }
  - { name: "--border", kind: "токен" }
  - { name: "--hairline", kind: "токен" }
  - { name: "--text-xs", kind: "токен" }
  - { name: "--tone-mark", kind: "токен" }
title-en: "Run lanes"
group-en: "Agent layer"
---

Несколько агентов, работающих параллельно, на общей оси времени. Отвечает на
вопрос, на который не отвечает ни очередь, ни лог: **что шло одновременно и
что кого ждало**.

```html preview
<div class="inst-lanes">
  <div class="inst-lane">
    <span class="inst-lane-label">worldgen-01</span>
    <span class="inst-lane-track">
      <span class="inst-lane-span" data-tone="ok" style="inset-inline-start:0%;inline-size:34%"></span>
    </span>
  </div>
  <div class="inst-lane">
    <span class="inst-lane-label">biomes-04</span>
    <span class="inst-lane-track">
      <span class="inst-lane-span" data-tone="running" style="inset-inline-start:30%;inline-size:52%"></span>
    </span>
  </div>
</div>
```

## Контракт

| Что | Обязательно | Почему |
|---|---|---|
| `inset-inline-start` и `inline-size` инлайном | да | Границы отрезка — **данные**, а не оформление: начало и длительность приходят с сервера |
| Логические свойства, а не `left` и `width` | да | В RTL ось времени зеркалится сама |
| Подпись или `title` у отрезка | да | Прямоугольник без имени и времени не сообщает ничего |
| `data-tone` | нет | Без атрибута отрезок нейтрален |

### Доступность

| | |
|---|---|
| Отрезок — не текст | Сам по себе он не озвучивается ничем. Дорожка обязана нести `aria-label` с длительностью: «worldgen-01, с 0 до 34 процентов, 12 секунд» |
| Подпись обрезается | Полное имя — в `title`, иначе исполнитель безымянен |
| Не только цвет | Тон отрезка дублируется в доступном имени словом |
| Режим принудительных цветов | Отрезок несёт значение, поэтому переживает сброс: `forced-color-adjust: none` и `Highlight` |
| Альтернатива | На узком экране дорожки нечитаемы. Резервный вид — [таблица](../components/display/table.md) с началом, концом и длительностью |
| Уменьшенное движение | Отрезок `running` не анимируется — длительность и так меняется скачками данных |

## Устройство

`inset-inline-start` и `inline-size` задаются **инлайном**, и это единственный
разрешённый случай инлайнового стиля в библиотеке: положение отрезка на оси —
значение, а не оформление. Тем же каналом идут ширина заполнения меры и глубина
узла дерева.

```html
<span class="inst-lane-span" style="inset-inline-start:30%;inline-size:52%"></span>
```

Проценты считаются от общей длительности прогона, поэтому все дорожки делят
одну шкалу — иначе сравнивать их бессмысленно.

## Правила

:::do Границы отрезка — данные
Начало и длительность приходят с сервера и живут в инлайновом стиле. Всё
остальное берётся из библиотеки.
:::

:::dont Дорожки вместо расписания
У дорожек нет шкалы с подписями времени. Точные моменты читаются
[таблицей](../components/display/table.md) или [лентой](../components/display/timeline.md).
:::

:::do Логические свойства оси
`inset-inline-start` вместо `left`: в RTL ось зеркалится сама.
:::

:::dont Отрезок без имени и времени
Прямоугольник на оси не сообщает ничего, пока рядом не написано, чей он и
сколько длился.
:::

## API

```api
```

## Связанное

```related
```
