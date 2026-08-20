---
title: Метрика
group: Отображение данных
layout: component
source: src/surfaces.css
api:
  - { name: "inst-metric", kind: "класс", doc: "Одна метрика" }
  - { name: "inst-metric-row", kind: "класс", doc: "Ряд метрик" }
  - { name: "inst-metric-label", kind: "класс", doc: "Подпись над числом" }
  - { name: "inst-metric-value", kind: "класс", doc: "Само число" }
  - { name: "inst-metric-unit", kind: "класс", doc: "Единица внутри числа: мельче и тише" }
  - { name: "inst-metric-delta", kind: "класс", doc: "Изменение к прошлому периоду" }
  - { name: "data-dir", kind: "атрибут", value: "up · down", doc: "на `inst-metric-delta`" }
  - { name: "data-tone", kind: "атрибут", value: "neutral · running · ok · warn · error", doc: "там же" }
  - { name: "--surface-sunken", kind: "токен" }
  - { name: "--radius-md", kind: "токен" }
  - { name: "--pad-panel", kind: "токен" }
  - { name: "--text-2xl", kind: "токен" }
  - { name: "--text-xs", kind: "токен" }
  - { name: "--tracking-tight", kind: "токен" }
  - { name: "--tone-ink", kind: "токен" }
  - { name: "--size-chevron", kind: "токен" }
  - { name: "--space-2", kind: "токен" }
title-en: "Metric"
group-en: "Data display"
---

Одно число, за которым следят. Единица набирается мельче и тише самого числа:
«42 с» одним кеглем читается как одно слово, а разделив размер и цвет, глаз
хватает число, и «с» уходит в фон.

```html preview
<div class="inst-metric-row">
  <div class="inst-metric">
    <div class="inst-metric-label">Прогонов за сутки</div>
    <div class="inst-metric-value">128</div>
    <div class="inst-metric-delta" data-dir="down" data-tone="ok">18% к прошлому</div>
  </div>
  <div class="inst-metric">
    <div class="inst-metric-label">Среднее время</div>
    <div class="inst-metric-value">4,2<span class="inst-metric-unit">с</span></div>
  </div>
  <div class="inst-metric">
    <div class="inst-metric-label">Предупреждений</div>
    <div class="inst-metric-value">5</div>
    <div class="inst-metric-delta" data-dir="up" data-tone="error">5 новых</div>
  </div>
</div>
```

## Контракт

| Что | Обязательно | Почему |
|---|---|---|
| `inst-metric-label` | да | «128» без ответа на «чего» не является метрикой |
| Знак изменения словом в дельте | да | Стрелка — второй носитель, а не единственный: «↓ 18%» читается и без тона |
| Доступное имя из подписи и числа | да, если метрика кликабельна | Иначе прозвучит «128» без ответа на «чего» |
| `data-dir` и `data-tone` | нет | Без первого дельта идёт без стрелки, без второго — нейтральным цветом |

### Доступность

| | |
|---|---|
| Подпись и число | Связаны визуально порядком. Если метрика кликабельна, доступное имя должно включать оба — иначе прозвучит «128» без ответа на «чего» |
| Не только цвет | Дельта несёт стрелку **и** знак изменения в тексте. Читатель, не различающий тон, читает «↓ 18%» |
| Контраст | Дельта берёт `--tone-ink` — порог текста 4.5:1, а не метки 3:1 |
| Кегль числа | `--text-2xl` предназначен **только** для числа-героя. Заголовок этого размера кричит громче данных |
| Табличные цифры | Включены по умолчанию: число, обновляемое на месте, не дёргает соседей |

## Устройство

### Направление и оценка — разные атрибуты

Это главная ловушка компонента. Стрелка вниз у времени прогона — хорошо,
стрелка вверх у предупреждений — плохо. Атрибут, названный стрелкой и
означающий оценку, — гарантированная ошибка применения.

| Атрибут | Что кодирует | Значения |
|---|---|---|
| `data-dir` | **Куда** изменилось | `up` · `down` |
| `data-tone` | **Хорошо это или плохо** | `ok` · `warn` · `error` · `neutral` · `running` |

```html preview
<div class="inst-metric-row">
  <div class="inst-metric">
    <div class="inst-metric-label">Время прогона</div>
    <div class="inst-metric-value">42<span class="inst-metric-unit">&nbsp;с</span></div>
    <div class="inst-metric-delta" data-dir="down" data-tone="ok">18% — упало, и это хорошо</div>
  </div>
  <div class="inst-metric">
    <div class="inst-metric-label">Предупреждений</div>
    <div class="inst-metric-value">12</div>
    <div class="inst-metric-delta" data-dir="up" data-tone="error">5 новых — выросло, и это плохо</div>
  </div>
</div>
```

Оба атрибута необязательны: дельта без `data-dir` идёт без стрелки, без
`data-tone` — нейтральным цветом.

## Композиции

### Ряд метрик — без рамок

```html preview
<div class="inst-metric-row">
  <div class="inst-metric">
    <div class="inst-metric-label">Всего</div>
    <div class="inst-metric-value">128</div>
  </div>
  <div class="inst-metric">
    <div class="inst-metric-label">Успешно</div>
    <div class="inst-metric-value">121</div>
  </div>
  <div class="inst-metric">
    <div class="inst-metric-label">Упало</div>
    <div class="inst-metric-value">7</div>
  </div>
</div>
```

Четыре числа — это **одна** группа. Рамка вокруг каждого превращает их в
четыре объекта, между которыми читатель начинает искать различия, которых нет.
Поэтому метрика стоит на приглушённой поверхности, а не в карточке.

## Сценарии

### Сводка прогона

Порознь метрика — просто крупное число; смысл появляется в ряду.

Что видно только здесь:

- **рамок нет ни у одной.** Четыре числа — одна группа; рамка вокруг каждого
  превратила бы их в четыре объекта, между которыми читатель начнёт искать
  различия;
- **единица мельче и тише числа.** «42 с» одним кеглем читается как одно
  слово; разделив размер и цвет, глаз хватает число, а «с» уходит в фон;
- **направление и оценка независимы.** Стрелка вниз у времени — хорошо,
  стрелка вверх у предупреждений — плохо.

```html preview context
<div class="inst-metric-row">
  <div class="inst-metric">
    <div class="inst-metric-label">Время прогона</div>
    <div class="inst-metric-value">42<span class="inst-metric-unit">&nbsp;с</span></div>
    <div class="inst-metric-delta" data-dir="down" data-tone="ok">18% к прошлому</div>
  </div>
  <div class="inst-metric">
    <div class="inst-metric-label">Тесты пройдены</div>
    <div class="inst-metric-value">248<span class="inst-metric-unit">/251</span></div>
    <div class="inst-metric-delta" data-tone="warn">3 падают</div>
  </div>
  <div class="inst-metric">
    <div class="inst-metric-label">Размер бандла</div>
    <div class="inst-metric-value">7,4<span class="inst-metric-unit">&nbsp;МБ</span></div>
    <div class="inst-metric-delta">без изменений</div>
  </div>
  <div class="inst-metric">
    <div class="inst-metric-label">Предупреждений</div>
    <div class="inst-metric-value">12</div>
    <div class="inst-metric-delta" data-dir="up" data-tone="error">5 новых</div>
  </div>
</div>
```

## Правила

:::do Направление и оценка порознь
`data-dir` говорит «куда», `data-tone` — «хорошо или плохо». Падение времени и
падение выручки — одна стрелка и разные тона.
:::

:::dont Один атрибут на оба смысла
Атрибут, названный стрелкой и означающий оценку, — гарантированная ошибка
применения на первом же обратном показателе.
:::

:::do Ряд без рамок
Сравнимые числа — одна группа. Приглушённая поверхность вместо рамки у
каждого.
:::

:::dont --text-2xl у заголовка
Этот кегль предназначен только для числа-героя. Заголовок такого размера
кричит громче данных.
:::

## API

```api
```

## Связанное

```related
```
