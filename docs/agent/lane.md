---
title: Дорожки прогонов
group: Агентный слой
status: stable
source: src/agent.css
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

## Границы отрезка — данные, а не оформление

`inset-inline-start` и `inline-size` задаются **инлайном**, и это единственный
разрешённый случай инлайнового стиля в этом ките: положение отрезка на оси —
значение, а не оформление. Тем же каналом идут ширина заполнения меры и глубина
узла дерева.

```html
<span class="inst-lane-span" style="inset-inline-start:30%;inline-size:52%"></span>
```

Проценты считаются от общей длительности прогона, поэтому все дорожки делят
одну шкалу — иначе сравнивать их бессмысленно.

## Когда использовать

| Используйте | Возьмите другое |
|---|---|
| Параллельные исполнители на общей оси | **Последовательные задачи** — [очередь](./task.md): у неё нет оси времени |
| Важны наложения и простои | **Один процесс с прогрессом** — [мера](../components/charts/meter.md) |
| Исполнителей 2–10 | **События без длительности** — [лента](../components/display/timeline.md) |

## Справочник

### Классы

| Класс | Работа |
|---|---|
| `inst-lanes` | Набор дорожек |
| `inst-lane` | Одна дорожка: подпись и ось |
| `inst-lane-label` | Имя исполнителя, обрезается |
| `inst-lane-track` | Ось времени |
| `inst-lane-span` | Отрезок работы |

### Атрибуты

| Атрибут | Значения |
|---|---|
| `data-tone` | `neutral` `running` `ok` `warn` `error` |

### Токены

`--label-col` · `--gap-inline` · `--size-meter` · `--radius-xs` ·
`--track` · `--border` · `--hairline` · `--text-xs` · `--tone-mark`

## Доступность

| | |
|---|---|
| Отрезок — не текст | Сам по себе он не озвучивается ничем. Дорожка обязана нести `aria-label` с длительностью: «worldgen-01, с 0 до 34 процентов, 12 секунд» |
| Подпись обрезается | Полное имя — в `title`, иначе исполнитель безымянен |
| Не только цвет | Тон отрезка дублируется в доступном имени словом |
| Режим принудительных цветов | Отрезок несёт значение, поэтому переживает сброс: `forced-color-adjust: none` и `Highlight` |
| Альтернатива | На узком экране дорожки нечитаемы. Резервный вид — [таблица](../components/display/table.md) с началом, концом и длительностью |
| Уменьшенное движение | Отрезок `running` не анимируется — длительность и так меняется скачками данных |

## Связанное

[Строка очереди](./task.md) · [Бюджет](./budget.md) ·
[Мера](../components/charts/meter.md) ·
[Лента событий](../components/display/timeline.md)

Исходник: `src/agent.css` · Почему кит устроен именно так — [конституция](../about/design-principles.md)
