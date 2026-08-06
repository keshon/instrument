---
title: Мера
group: Графики
status: stable
source: src/components.css
---

Одна величина относительно своего потолка: расход бюджета, доля времени,
готовность шага. Один тон, потому что сравнивается **одна** величина.

```html preview
<div class="inst-meter-row"><span>Рендер</span><span class="inst-meter-value">18,2 с</span></div>
<div class="inst-meter" role="progressbar" aria-label="Рендер"
     aria-valuenow="43" aria-valuemin="0" aria-valuemax="100">
  <div class="inst-meter-fill" style="inline-size:43%"></div>
</div>
```

## Установка

```html
<link rel="stylesheet" href="src/kit.css">
```

## Когда использовать

| Используйте | Возьмите другое |
|---|---|
| Доля от известного целого: 43 из 100, 184 320 из 250 000 | **Сравнение нескольких рядов между собой** — [легенда](./legend.md) и категориальная палитра: четыре меры разного цвета врут о различии, которого нет |
| Расход исчерпаемого ресурса — [бюджет прогона](../../agent/budget.md) | **Доля в компактном месте, где нужен круг** — [кольцо](./ring.md) |
| Незавершённая работа неизвестной длительности — `data-state="indeterminate"` | **Форма изменения во времени** — [спарклайн](./sparkline.md): у меры нет истории |
| Полоса под числом в карточке — [метрика](../display/metric.md) | **Ожидание, у которого нет величины вовсе** — скелет или спиннер: [заглушка](../feedback/skeleton.md) |

## Два вида полосы

Определённая и неопределённая — **разные компоненты**, а не два состояния
одного. Определённая обещает остаток, неопределённая — только факт работы.

```html preview
<div class="inst-meter" role="progressbar" aria-label="Расход токенов"
     aria-valuenow="74" aria-valuemin="0" aria-valuemax="100">
  <div class="inst-meter-fill" style="inline-size:74%"></div>
</div>
<div class="inst-meter" data-state="indeterminate" role="progressbar" aria-label="Индексация">
  <div class="inst-meter-fill"></div>
</div>
```

| | Определённая | Неопределённая |
|---|---|---|
| Атрибут | без `data-state` | `data-state="indeterminate"` |
| Ширина заполнения | приходит из данных | задаёт кит — 35%, и она метётся |
| `aria-valuenow` | **обязателен** | **отсутствует, и это сообщение** |
| Что обещает читателю | остаток и темп | что работа идёт |

Отсутствие `aria-valuenow` — не упущение разметки, а её содержание: по нему
скринридер объявляет «занято» вместо ложного процента. Поставить `valuenow`
неопределённой полосе — то же самое, что нарисовать её застрявшей на 90%.

**Никогда не изображайте неизвестный прогресс определённой полосой.** Полоса,
которая врёт про остаток, хуже честного «длительность неизвестна» рядом с
метущейся: первую пользователь ждёт, вторую — планирует.

## Тон

`data-tone` ставится на `inst-meter`, заливка читает `--tone-mark` и больше
ничего о тоне не знает. Тон здесь — **исключение, а не приправа**: полоса без
атрибута берёт акцент, и это правильный ответ в подавляющем большинстве
случаев.

```html preview
<div class="inst-meter" role="progressbar" aria-label="Расход"
     aria-valuenow="62" aria-valuemin="0" aria-valuemax="100">
  <div class="inst-meter-fill" style="inline-size:62%"></div>
</div>
<div class="inst-meter" data-tone="warn" role="progressbar" aria-label="Расход у порога"
     aria-valuenow="88" aria-valuemin="0" aria-valuemax="100">
  <div class="inst-meter-fill" style="inline-size:88%"></div>
</div>
<div class="inst-meter" data-tone="error" role="progressbar" aria-label="Лимит превышен"
     aria-valuenow="100" aria-valuemin="0" aria-valuemax="100">
  <div class="inst-meter-fill" style="inline-size:100%"></div>
</div>
```

Тон меняется, когда меняется **смысл величины** — перешли порог, исчерпали
лимит. Тон, назначенный по номеру полосы в списке, сообщает различие, которого
нет.

## Ширина как канал данных

`style="inline-size:43%"` на заливке — единственный инлайновый стиль на этой
странице, и он допустим ровно потому, что это **не оформление, а значение**:
величина приходит с сервера и физически не может жить в статическом CSS. Всё
остальное — высота, радиус, дорожка, цвет — берётся из кита.

Заливка едет с переходом (`--dur-3`), поэтому обновление значения читается как
движение, а не как скачок.

## Композиции

### Строка с подписью и величиной

```html preview
<div>
  <div class="inst-meter-row"><span>Физика</span><span class="inst-meter-value">11,5 с</span></div>
  <div class="inst-meter" role="progressbar" aria-label="Физика"
       aria-valuenow="27" aria-valuemin="0" aria-valuemax="100">
    <div class="inst-meter-fill" style="inline-size:27%"></div>
  </div>
</div>
<div>
  <div class="inst-meter-row"><span>Индексация</span><span class="inst-meter-value">длительность неизвестна</span></div>
  <div class="inst-meter" data-state="indeterminate" role="progressbar" aria-label="Индексация">
    <div class="inst-meter-fill"></div>
  </div>
</div>
```

`inst-meter-row` разносит подпись и величину по краям и стоит **над** полосой:
число слева от полосы съедало бы её длину, а длина здесь и есть данные.

### В бюджете прогона

```html
<div class="inst-budget">
  <div class="inst-budget-row">
    <span>Токены</span>
    <span class="inst-budget-value">184 320 <span class="inst-budget-of">из 250 000</span></span>
  </div>
  <div class="inst-meter" role="progressbar" aria-label="Расход токенов"
       aria-valuenow="74" aria-valuemin="0" aria-valuemax="100">
    <div class="inst-meter-fill" style="inline-size:74%"></div>
  </div>
  <div class="inst-budget-note">При текущем темпе хватит ещё примерно на 3 прогона.</div>
</div>
```

## Справочник

### Классы

| Класс | Работа |
|---|---|
| `inst-meter` | Дорожка. Носитель роли и значений |
| `inst-meter-fill` | Заливка. Обязательный единственный ребёнок |
| `inst-meter-row` | Строка «подпись — величина» над полосой |
| `inst-meter-value` | Величина в этой строке, приглушённым цветом |

### Обязательная разметка

| Что | Почему |
|---|---|
| `role="progressbar"` на `inst-meter` | Без роли значение существует только в пикселях ширины и недоступно ни скринридеру, ни тесту |
| `aria-valuenow` `aria-valuemin` `aria-valuemax` | Определённая полоса обязана сообщать число, а не картинку |
| `aria-label` или `aria-labelledby` | У полосы нет своего текста. Подпись из `inst-meter-row` связывается через `aria-labelledby`, если она видима |
| **отсутствие** `aria-valuenow` у `indeterminate` | Ровно так объявляется «занято, длительность неизвестна» |
| `inst-meter-fill` внутри | Дорожка сама ничего не заливает |

### Атрибуты

| Атрибут | Словарь |
|---|---|
| `data-state` | `indeterminate` — единственное значение у меры. Без атрибута полоса определённая |
| `data-tone` | `neutral` `running` `ok` `warn` `error`. Закрыт. Без атрибута — акцент |

### Токены

`--size-meter` · `--radius-xs` · `--track` · `--hairline` · `--border` ·
`--tone-mark` · `--accent-mark` · `--text-sm` · `--space-3` · `--gap-inline` ·
`--dur-3` · `--ease-out`

Протяжённость дорожки показывает **внутренний хайрлайн**, а не её собственный
цвет: так дорожка остаётся достаточно светлой, чтобы заливка держала на ней
3:1.

## Доступность

| | |
|---|---|
| Роль | `progressbar` — в примере, а не в примечании: `aria-valuenow` без роли невалиден |
| Имя | `aria-label` обязателен всегда: полоса — это заливка без текста |
| Неопределённость | Отсутствие `aria-valuenow` объявляется как «занято». Слова «длительность неизвестна» рядом дублируют это для тех, кто смотрит |
| Цвет не единственный носитель | Величина продублирована числом в `inst-meter-value`. Полоса, у которой значение есть только в цвете и длине, недоступна дальтонику и невоспроизводима на печати |
| Контраст | Заливка держит 3:1 против дорожки. Дорожка не несущая, поэтому 3:1 к поверхности с неё не спрашивается |
| Уменьшенное движение | Метение неопределённой полосы **замедляется до 3 с, а не гаснет**: остановленный индикатор сообщает «ничего не происходит», когда происходит |
| Режим принудительных цветов | Дорожка получает рамку `CanvasText`, заливка переживает сброс (`forced-color-adjust: none`), потому что несёт значение |
| Печать | Метение на бумаге бессмысленно: заливка неопределённой полосы не печатается, остаётся пустая дорожка и подпись |

## Связанное

[Кольцо](./ring.md) · [Спарклайн](./sparkline.md) · [Легенда](./legend.md) ·
[Метрика](../display/metric.md) · [Бюджет прогона](../../agent/budget.md) ·
[Дорожка агента](../../agent/lane.md) · [Заглушка](../feedback/skeleton.md)

Исходник: `src/components.css` · Почему кит устроен именно так — [конституция](../../about/design-principles.md)
