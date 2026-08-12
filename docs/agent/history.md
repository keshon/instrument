---
title: История проверок
group: Агентный слой
layout: component
source: src/agent.css
api:
  - { name: "inst-history", kind: "класс", doc: "Полоса. Свежие исходы справа" }
  - { name: "inst-history-tick", kind: "класс", doc: "Один исход: одна проверка, один прогон" }
  - { name: "data-tone", kind: "атрибут", value: "ok · warn · error · running · neutral", doc: "Исход штриха. Без атрибута штрих серый: проверка была, результат неизвестен" }
  - { name: "data-empty", kind: "атрибут", doc: "Проверки не было. Место остаётся, метки нет" }
  - { name: "inst-history-group", kind: "класс", doc: "Пачка штрихов с общим именем: час, день, прогон" }
  - { name: "--n", kind: "токен", doc: "Число штрихов в группе. Задаёт её долю от полосы" }
  - { name: "role", kind: "атрибут", value: "img", doc: "Полоса — картинка данных, и у неё обязано быть имя" }
  - { name: "--size-tick", kind: "токен" }
  - { name: "--space-1", kind: "токен" }
  - { name: "--space-3", kind: "токен" }
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

## Варианты

### Группы

Час, день, прогон: пачка проверок, у которой есть общее имя.

Полоса без групп отвечает на «что было». С группами она отвечает ещё и на
«сколько за этим стоит»: час с одной проверкой — один кирпич, час с семью —
семь.

```html preview
<div class="inst-history" style="inline-size:26rem" role="img" aria-label="Шесть часов: 19 проверок, две с ошибкой">
  <span class="inst-history-group" style="--n:1"><span class="inst-history-tick" data-tone="ok"></span></span>
  <span class="inst-history-group" style="--n:7"
    ><span class="inst-history-tick" data-tone="ok"></span
    ><span class="inst-history-tick" data-tone="ok"></span
    ><span class="inst-history-tick" data-tone="error"></span
    ><span class="inst-history-tick" data-tone="error"></span
    ><span class="inst-history-tick" data-tone="ok"></span
    ><span class="inst-history-tick" data-tone="ok"></span
    ><span class="inst-history-tick" data-tone="ok"></span
  ></span>
  <span class="inst-history-group" style="--n:1"><span class="inst-history-tick" data-empty></span></span>
  <span class="inst-history-group" style="--n:4"
    ><span class="inst-history-tick" data-tone="ok"></span
    ><span class="inst-history-tick" data-tone="ok"></span
    ><span class="inst-history-tick" data-tone="ok"></span
    ><span class="inst-history-tick" data-tone="ok"></span
  ></span>
  <span class="inst-history-group" style="--n:6"
    ><span class="inst-history-tick" data-tone="ok"></span
    ><span class="inst-history-tick" data-tone="ok"></span
    ><span class="inst-history-tick" data-tone="ok"></span
    ><span class="inst-history-tick" data-tone="ok"></span
    ><span class="inst-history-tick" data-tone="ok"></span
    ><span class="inst-history-tick" data-tone="ok"></span
  ></span>
</div>
```

Ширина группы **пропорциональна** числу штрихов в ней — `--n` в разметке.
Ширину полосе даёт контейнер: группы делят то, что им дали.

Равные слоты не годятся, и это меряется: полоса в 520px, разбитая на 24 часа,
даёт 14px на час, а в 14px помещается три штриха. Пропорция, наоборот, всегда
занимает ровно ширину полосы, сколько бы проверок ни было — густой монитор
выглядит густым, редкий редким.

Полоса с группами **не обрезает** содержимое: обрезка ради «старое уходит
первым» здесь не нужна, всё уместилось по построению. Заодно это и позволяет
повесить на штрих [тултип](../components/overlays/tooltip.md) — под
`overflow: hidden` его срезало бы.

| Без групп | С группами |
|---|---|
| Ширина штриха фиксирована — длина полосы говорит, сколько было проверок | Ширина полосы фиксирована — длина группы говорит, сколько проверок в этом часе |
| Лишнее обрезается началом | Не обрезается: доли считаются от того, что есть |
| Одна строка списка рядом с другой | Один объект крупно, на своём экране |

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
  <div class="inst-history" role="img" aria-label="Два часа без проверок">
    <span class="inst-history-tick" data-tone="ok"></span>
    <span class="inst-history-tick" data-tone="ok"></span>
    <span class="inst-history-tick" data-empty></span>
    <span class="inst-history-tick" data-empty></span>
    <span class="inst-history-tick" data-tone="ok"></span>
    <span class="inst-history-tick" data-tone="error"></span>
  </div>
</div>
```

| Состояние | Как ставится | Что значит |
|---|---|---|
| исход известен | `data-tone` | Проверка была и вот чем кончилась |
| исход неизвестен | без атрибутов | Проверка была, результата у нас нет |
| проверки не было | `data-empty` | **Никто не смотрел.** Место остаётся, метки нет |

Разница между двумя последними важнее, чем кажется. «Не знаем результат» —
про проверку, «не было проверки» — про нас.

Пропуск не красится и не обводится, и то и другое пробовалось. Рамка сбивает
плотность ряда: глаз читает смену **текстуры** как смену ритма, а не как смену
значения, и полоса «плывёт». Цвет обещает оценку — а на реальных данных
пропуск бывает не исключением, а нормой: в приложении, где история пришла из
импорта, тоном закрасилось бы 288 штрихов из 336, и исключением перестало бы
быть всё.

Если пропуск заслуживает тревоги, её место — **в строке рядом**: «проверок нет
с 07:14». Одно утверждение — одно место.

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
