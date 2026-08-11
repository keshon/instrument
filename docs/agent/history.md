---
title: История проверок
group: Агентный слой
layout: component
source: src/agent.css
api:
  - { name: "inst-history", kind: "класс", doc: "Полоса. Свежие исходы справа" }
  - { name: "inst-history-tick", kind: "класс", doc: "Один исход: одна проверка, один прогон" }
  - { name: "data-tone", kind: "атрибут", value: "ok · warn · error · running · neutral", doc: "Исход штриха. Без атрибута штрих серый: проверка была, результат неизвестен" }
  - { name: "role", kind: "атрибут", value: "img", doc: "Полоса — картинка данных, и у неё обязано быть имя" }
  - { name: "--size-tick", kind: "токен" }
  - { name: "--space-1", kind: "токен" }
  - { name: "--radius-xs", kind: "токен" }
  - { name: "--track", kind: "токен" }
  - { name: "--size-icon", kind: "токен" }
title-en: "Check history"
group-en: "Agent layer"
---

Ряд дискретных исходов: одна проверка — один штрих, свежие справа. Аптайм
монитора, прогоны ночной сборки, попытки повтора.

```html preview
<div class="inst-history" role="img" aria-label="Последние 24 проверки: 22 успешны, 2 с ошибкой">
  <span class="inst-history-tick" data-tone="ok"></span>
  <span class="inst-history-tick" data-tone="ok"></span>
  <span class="inst-history-tick" data-tone="error"></span>
  <span class="inst-history-tick" data-tone="ok"></span>
  <span class="inst-history-tick" data-tone="ok"></span>
  <span class="inst-history-tick" data-tone="ok"></span>
  <span class="inst-history-tick" data-tone="ok"></span>
  <span class="inst-history-tick" data-tone="warn"></span>
  <span class="inst-history-tick" data-tone="ok"></span>
  <span class="inst-history-tick" data-tone="ok"></span>
  <span class="inst-history-tick" data-tone="ok"></span>
  <span class="inst-history-tick" data-tone="ok"></span>
</div>
```

## Использование

```html
<div class="inst-history" role="img" aria-label="Последние 24 проверки: 22 успешны">
  <span class="inst-history-tick" data-tone="ok"></span>
  <span class="inst-history-tick" data-tone="error"></span>
</div>
```

| Что | Обязательно | Почему |
|---|---|---|
| `role="img"` и `aria-label` с **итогом словами** | да | Двести пустых `<span>` для скринридера — двести пустых `<span>`. Итог обязан быть текстом: «22 из 24 успешны» |
| Свежие **справа** | да | Полоса читается как время, а время в этой письменности идёт вправо |
| Тон из словаря | да | `ok` · `warn` · `error` · `running` · `neutral`. Шестого исхода не бывает |
| Штрих без тона | нет | Проверка была, результат неизвестен: штрих остаётся цветом дорожки |

## Когда использовать

| Используйте | Возьмите другое |
|---|---|
| Ряд равных событий с перечислимым исходом | **Величина во времени** — [спарклайн](../components/charts/sparkline.md): там важно значение, а не исход |
| Аптайм, прогоны, попытки повтора | **Отрезки разной длины на общей оси** — [дорожки](./lane.md) |
| Один показатель против предела | **Доля** — [мера](../components/charts/meter.md) или [кольцо](../components/charts/ring.md) |

## Устройство

Ширина штриха фиксирована (`--size-tick`), а не доля. Доля сделала бы полосу
из двадцати проверок неотличимой от полосы из двухсот, хотя означают они
разное: короткая история и длинная.

Полоса прижата к концу, поэтому лишнее уходит за **начало** и обрезается там:
старое исчезает первым, а последняя проверка видна всегда. Разметка при этом
обычная — от старых к свежим, как в массиве, который отдаёт сервер.

Высота — `--size-icon`, ростом с глиф. Полоса стоит в одном ряду с подписью и
бейджем, и на полную высоту строки она перестаёт быть отметкой рядом с ними:
становится блоком и перетягивает ряд на себя.

## Состояния

```html preview
<div class="inst-stack" style="max-inline-size:22rem">
  <div class="inst-history" role="img" aria-label="Всё успешно">
    <span class="inst-history-tick" data-tone="ok"></span>
    <span class="inst-history-tick" data-tone="ok"></span>
    <span class="inst-history-tick" data-tone="ok"></span>
    <span class="inst-history-tick" data-tone="ok"></span>
    <span class="inst-history-tick" data-tone="ok"></span>
    <span class="inst-history-tick" data-tone="ok"></span>
  </div>
  <div class="inst-history" role="img" aria-label="Падало дважды, сейчас работает">
    <span class="inst-history-tick" data-tone="ok"></span>
    <span class="inst-history-tick" data-tone="error"></span>
    <span class="inst-history-tick" data-tone="error"></span>
    <span class="inst-history-tick" data-tone="ok"></span>
    <span class="inst-history-tick" data-tone="ok"></span>
    <span class="inst-history-tick" data-tone="ok"></span>
  </div>
  <div class="inst-history" role="img" aria-label="Проверок ещё не было">
    <span class="inst-history-tick"></span>
    <span class="inst-history-tick"></span>
    <span class="inst-history-tick"></span>
    <span class="inst-history-tick"></span>
    <span class="inst-history-tick"></span>
    <span class="inst-history-tick"></span>
  </div>
</div>
```

## Композиции

### В карточке монитора

```html preview context
<div class="inst-card" style="max-inline-size:24rem">
  <div class="inst-card-head">
    <span class="inst-card-title">absalyamov.ru</span>
    <span class="inst-badge" data-tone="ok"><span class="inst-dot"></span>работает</span>
  </div>
  <div class="inst-card-sub">https://absalyamov.ru/ · http</div>
  <div class="inst-history" role="img" aria-label="Последние 20 проверок: одна с ошибкой">
    <span class="inst-history-tick" data-tone="ok"></span>
    <span class="inst-history-tick" data-tone="ok"></span>
    <span class="inst-history-tick" data-tone="ok"></span>
    <span class="inst-history-tick" data-tone="error"></span>
    <span class="inst-history-tick" data-tone="ok"></span>
    <span class="inst-history-tick" data-tone="ok"></span>
    <span class="inst-history-tick" data-tone="ok"></span>
    <span class="inst-history-tick" data-tone="ok"></span>
    <span class="inst-history-tick" data-tone="ok"></span>
    <span class="inst-history-tick" data-tone="ok"></span>
  </div>
</div>
```

## Правила

:::do Итог словами в aria-label
«22 из 24 проверок успешны, последняя с ошибкой». Полоса без имени для
скринридера не существует.
:::

:::dont Штрих на каждую секунду
Полоса показывает последние события, а не весь журнал. Двести штрихов по три
пикселя — это шестьсот пикселей ряби, в которой ничего не видно.
:::

:::do Пустая история — серые штрихи
Проверок не было — это тоже состояние, и его видно.
:::

:::dont Цвет как единственный признак
Рядом с полосой стоит статус словом: [бейдж](../components/display/badge.md)
или подпись. Полоса уточняет, а не сообщает.
:::

## Доступность

| | |
|---|---|
| Имя | `role="img"` плюс `aria-label` с итогом словами. Штрихи пусты и скринридеру не читаются — читается имя полосы |
| Цвет не единственный носитель | Полоса стоит рядом со статусом словом. Сама по себе она уточнение: «падало дважды за сутки», а не «сейчас упало» |
| Контраст | Тон берёт `--tone-mark` — порог 3:1 против поверхности, тот же, что у метки состояния |
| Движение | Анимации нет: история — это состояние, а не процесс |
| Печать | Печатается как есть: штрихи различаются тоном и на бумаге остаются серыми ступенями |

## API

```api
```

## Связанное

[Дорожки прогонов](./lane.md)
[Спарклайн](../components/charts/sparkline.md)
[Мера](../components/charts/meter.md)
[Бейдж](../components/display/badge.md)
[Строка очереди](./task.md)
