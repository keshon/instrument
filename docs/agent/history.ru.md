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
  - { name: "inst-history-group", kind: "класс", doc: "Пачка штрихов с общим именем: час, день, прогон. Якорь собственного тултипа, отвечает на наведение" }
  - { name: "--n", kind: "токен", doc: "Число штрихов в группе. Задаёт её долю от полосы, и то же число задаёт долю ячейки оси под ней" }
  - { name: "inst-history--sm", kind: "класс", doc: "Высота меры — для полосы, которая служит колонкой в таблице объектов" }
  - { name: "inst-history--lg", kind: "класс", doc: "Высота мелкого контрола — для полосы, которая на экране предмет" }
  - { name: "inst-history-axis", kind: "класс", doc: "Линейка под полосой. Повторяет геометрию групп, а не вычисляет положения" }
  - { name: "inst-history-axis-cell", kind: "класс", doc: "Одна подпись. Берёт тот же `--n`, что и группа над ней" }
  - { name: "data-minor", kind: "атрибут", doc: "На ячейке оси: эту подпись можно уронить, когда линейка уже 36rem. Место ячейка сохраняет" }
  - { name: "role", kind: "атрибут", value: "img", doc: "Полоса — картинка данных, и у неё обязано быть имя" }
  - { name: "--size-tick", kind: "токен" }
  - { name: "--space-1", kind: "токен" }
  - { name: "--space-3", kind: "токен" }
  - { name: "--radius-mark", kind: "токен" }
  - { name: "--track", kind: "токен" }
  - { name: "--size-icon", kind: "токен" }
  - { name: "--size-meter", kind: "токен" }
  - { name: "--control-h-sm", kind: "токен" }
  - { name: "--space-2", kind: "токен" }
  - { name: "--text-2xs", kind: "токен" }
  - { name: "--text-muted", kind: "токен" }
  - { name: "--hairline", kind: "токен" }
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

## Контракт

| Что | Обязательно | Почему |
|---|---|---|
| `role="img"` и `aria-label` с **итогом словами** | да | Двести пустых `<span>` для скринридера — двести пустых `<span>`. Итог обязан быть текстом: «22 из 24 успешны» |
| Свежие **справа** | да | Полоса читается как время, а время в этой письменности идёт вправо |
| Тон из словаря | да | `ok` · `warn` · `error` · `running` · `neutral`. Шестого исхода не бывает |
| Штрих без тона | нет | Проверка была, результат неизвестен: штрих остаётся цветом дорожки |
| `tabindex` на группе | нет | Он делает группу целью нажатия, а сутки по часам дают ей 15px против 24 по WCAG 2.2 AA. До скринридера показание доходит через `aria-label` самой полосы |

### Доступность

| | |
|---|---|
| Имя | `role="img"` плюс `aria-label` с итогом словами. Штрихи пусты и скринридеру не читаются — читается имя полосы |
| Цвет не единственный носитель | Полоса стоит рядом со статусом словом. Сама по себе она уточнение: «падало дважды за сутки», а не «сейчас упало» |
| Контраст | Тон берёт `--tone-mark` — порог 3:1 против поверхности, тот же, что у метки состояния |
| Движение | Анимации нет: история — это состояние, а не процесс |
| Печать | Печатается как есть: штрихи различаются тоном и на бумаге остаются серыми ступенями |

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
первым» здесь не нужна, и именно это позволяет повесить на группу
[тултип](../components/overlays/tooltip.md) — под `overflow: hidden` его
срезало бы.

| Без групп | С группами |
|---|---|
| Ширина штриха фиксирована — длина полосы говорит, сколько было проверок | Ширина полосы фиксирована — длина группы говорит, сколько проверок в этом часе |
| Лишнее обрезается началом | Не обрезается: доли считаются от того, что есть |
| Одна строка списка рядом с другой | Один объект крупно, на своём экране |

### Собранная форма

Выше примерно четырёх проверок в час штрих слишком тонок, чтобы его увидеть, не
говоря уже о том, чтобы в него попасть. Дайте группе **один** штрих и оставьте
`--n` как есть: час сохранит свою точную долю и вернётся одним кирпичом,
окрашенным по худшему своему исходу.

```html preview
<div class="inst-history" style="inline-size:26rem" role="img" aria-label="Двенадцать часов по 7 проверок: два часа с ошибками">
  <span class="inst-history-group" style="--n:7"><span class="inst-history-tick" data-tone="ok"></span></span>
  <span class="inst-history-group" style="--n:7"><span class="inst-history-tick" data-tone="ok"></span></span>
  <span class="inst-history-group" style="--n:7"><span class="inst-history-tick" data-tone="error"></span></span>
  <span class="inst-history-group" style="--n:7"><span class="inst-history-tick" data-tone="warn"></span></span>
  <span class="inst-history-group" style="--n:7"><span class="inst-history-tick" data-tone="ok"></span></span>
  <span class="inst-history-group" style="--n:1"><span class="inst-history-tick" data-tone="ok"></span></span>
  <span class="inst-history-group" style="--n:7"><span class="inst-history-tick" data-tone="ok"></span></span>
  <span class="inst-history-group" style="--n:7"><span class="inst-history-tick" data-tone="ok"></span></span>
  <span class="inst-history-group" style="--n:12"><span class="inst-history-tick" data-tone="ok"></span></span>
  <span class="inst-history-group" style="--n:7"><span class="inst-history-tick" data-tone="ok"></span></span>
  <span class="inst-history-group" style="--n:7"><span class="inst-history-tick" data-tone="ok"></span></span>
  <span class="inst-history-group" style="--n:7"><span class="inst-history-tick" data-tone="ok"></span></span>
</div>
```

Вес часа держит `--n`: в шестом была одна проверка, и он узкий, в девятом
двенадцать, и он широкий.

## Размеры

```html preview
<div class="inst-stack" style="max-inline-size:22rem">
  <div class="inst-history inst-history--sm" role="img" aria-label="Мелкая: восемь проверок, одна с ошибкой">
    <span class="inst-history-tick" data-tone="ok"></span>
    <span class="inst-history-tick" data-tone="ok"></span>
    <span class="inst-history-tick" data-tone="error"></span>
    <span class="inst-history-tick" data-tone="ok"></span>
    <span class="inst-history-tick" data-tone="ok"></span>
    <span class="inst-history-tick" data-tone="ok"></span>
    <span class="inst-history-tick" data-tone="ok"></span>
    <span class="inst-history-tick" data-tone="ok"></span>
  </div>
  <div class="inst-history" role="img" aria-label="Базовая: восемь проверок, одна с ошибкой">
    <span class="inst-history-tick" data-tone="ok"></span>
    <span class="inst-history-tick" data-tone="ok"></span>
    <span class="inst-history-tick" data-tone="error"></span>
    <span class="inst-history-tick" data-tone="ok"></span>
    <span class="inst-history-tick" data-tone="ok"></span>
    <span class="inst-history-tick" data-tone="ok"></span>
    <span class="inst-history-tick" data-tone="ok"></span>
    <span class="inst-history-tick" data-tone="ok"></span>
  </div>
  <div class="inst-history inst-history--lg" role="img" aria-label="Крупная: восемь проверок, одна с ошибкой">
    <span class="inst-history-tick" data-tone="ok"></span>
    <span class="inst-history-tick" data-tone="ok"></span>
    <span class="inst-history-tick" data-tone="error"></span>
    <span class="inst-history-tick" data-tone="ok"></span>
    <span class="inst-history-tick" data-tone="ok"></span>
    <span class="inst-history-tick" data-tone="ok"></span>
    <span class="inst-history-tick" data-tone="ok"></span>
    <span class="inst-history-tick" data-tone="ok"></span>
  </div>
</div>
```

| | Высота | Где |
|---|---|---|
| `inst-history--sm` | `--size-meter` | Колонка в таблице объектов. На высоте глифа полоса перекрикивает имена, которым принадлежит |
| базовая | `--size-icon` | В строке рядом с именем и бейджем, где полоса — пометка среди пометок |
| `inst-history--lg` | `--control-h-sm` | На экране **одного** объекта, где полоса предмет, а не примечание |

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

### С осью

Полоса говорит, что случилось, ось — когда. Ячейка оси берёт **тот же `--n`**,
что и группа над ней, поэтому обе строки раскладывает одно правило и разъехаться
им нечем.

```html preview
<div style="inline-size:30rem">
  <div class="inst-history inst-history--lg" role="img"
       aria-label="Двенадцать часов, 60 проверок: три с ошибкой в 03:00">
    <span class="inst-history-group" style="--n:5"><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-tooltip-text">00:00 &mdash; 5 проверок</span></span>
    <span class="inst-history-group" style="--n:5"><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-tooltip-text">01:00 &mdash; 5 проверок</span></span>
    <span class="inst-history-group" style="--n:5"><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-tooltip-text">02:00 &mdash; 5 проверок</span></span>
    <span class="inst-history-group" style="--n:5"><span class="inst-history-tick" data-tone="error"></span><span class="inst-history-tick" data-tone="error"></span><span class="inst-history-tick" data-tone="error"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-tooltip-text">03:00 &mdash; 5 проверок</span></span>
    <span class="inst-history-group" style="--n:2"><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-tooltip-text">04:00 &mdash; 2 проверок</span></span>
    <span class="inst-history-group" style="--n:5"><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-tooltip-text">05:00 &mdash; 5 проверок</span></span>
    <span class="inst-history-group" style="--n:5"><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-tooltip-text">06:00 &mdash; 5 проверок</span></span>
    <span class="inst-history-group" style="--n:5"><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-tooltip-text">07:00 &mdash; 5 проверок</span></span>
    <span class="inst-history-group" style="--n:8"><span class="inst-history-tick" data-tone="ok"></span><span class="inst-tooltip-text">08:00 &mdash; 8 проверок, показаны одним</span></span>
    <span class="inst-history-group" style="--n:5"><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-tooltip-text">09:00 &mdash; 5 проверок</span></span>
    <span class="inst-history-group" style="--n:5"><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-tooltip-text">10:00 &mdash; 5 проверок</span></span>
    <span class="inst-history-group" style="--n:5"><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-tooltip-text">11:00 &mdash; 5 проверок</span></span>
  </div>
  <div class="inst-history-axis" aria-hidden="true">
    <span class="inst-history-axis-cell" style="--n:5">00</span>
    <span class="inst-history-axis-cell" style="--n:5" data-minor>01</span>
    <span class="inst-history-axis-cell" style="--n:5" data-minor>02</span>
    <span class="inst-history-axis-cell" style="--n:5">03</span>
    <span class="inst-history-axis-cell" style="--n:2" data-minor>04</span>
    <span class="inst-history-axis-cell" style="--n:5" data-minor>05</span>
    <span class="inst-history-axis-cell" style="--n:5">06</span>
    <span class="inst-history-axis-cell" style="--n:5" data-minor>07</span>
    <span class="inst-history-axis-cell" style="--n:8" data-minor>08</span>
    <span class="inst-history-axis-cell" style="--n:5">09</span>
    <span class="inst-history-axis-cell" style="--n:5" data-minor>10</span>
    <span class="inst-history-axis-cell" style="--n:5" data-minor>11</span>
  </div>
</div>
```

`data-minor` помечает подписи, которые могут уйти, когда линейка уже 36rem.
Ячейка остаётся на месте — `visibility`, а не `display`, — иначе уцелевшие
подписи сдвинулись бы.

Подсказка — обычный [тултип](../components/overlays/tooltip.md), положенный
прямо внутрь группы. На наведение отвечает сам час, который приподнимается, и по
горизонтали это не стоит ничего: блочный отступ в инлайновом размере не
участвует, поэтому часы сдвинуться не могут.

### В карточке монитора

```html preview context
<div class="inst-card" style="max-inline-size:24rem">
  <div class="inst-card-head">
    <span class="inst-card-title">absalyamov.ru</span>
    <span class="inst-badge" data-tone="ok"><span class="inst-dot"></span>работает</span>
  </div>
  <div class="inst-card-sub">https://absalyamov.ru/ · http</div>
  <div class="inst-history" role="img" aria-label="Сутки проверок, по четыре в час: три с ошибкой в 12:45, два часа без проверок">
    <span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span>
    <span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-empty></span><span class="inst-history-tick" data-empty></span>
    <span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span>
    <span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span>
    <span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="error"></span><span class="inst-history-tick" data-tone="error"></span><span class="inst-history-tick" data-tone="error"></span><span class="inst-history-tick" data-tone="warn"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span>
    <span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span>
    <span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span>
    <span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span>
  </div>
</div>
```

Полоса занимает карточку потому, что ей дали достаточно проверок, а не потому,
что она тянется: девяносто шесть — это сутки по четыре в час, и самые старые
обрезаются передним краем. Чтобы занять ширину, не известную заранее, берут
[форму с группами](#groups).

## API

```api
```

## Связанное

```related
```
