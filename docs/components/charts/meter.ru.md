---
title: Мера
group: Графики
layout: component
source: src/status.css
api:
  - { name: "inst-meter", kind: "класс", doc: "Дорожка. Носитель роли и значений" }
  - { name: "inst-meter-fill", kind: "класс", doc: "Заливка. Обязательный единственный ребёнок" }
  - { name: "inst-meter-row", kind: "класс", doc: "Строка «подпись — величина» над полосой" }
  - { name: "inst-meter-value", kind: "класс", doc: "Величина в этой строке, приглушённым цветом" }
  - { name: "data-state", kind: "атрибут", doc: "`indeterminate` — единственное значение у меры. Без атрибута полоса определённая" }
  - { name: "data-tone", kind: "атрибут", doc: "`neutral` `running` `ok` `warn` `error`. Закрыт. Без атрибута — акцент" }
  - { name: "--size-meter", kind: "токен" }
  - { name: "--radius-mark", kind: "токен" }
  - { name: "--track", kind: "токен" }
  - { name: "--hairline", kind: "токен" }
  - { name: "--border", kind: "токен" }
  - { name: "--tone-mark", kind: "токен" }
  - { name: "--accent-mark", kind: "токен" }
  - { name: "--text-sm", kind: "токен" }
  - { name: "--space-3", kind: "токен" }
  - { name: "--gap-inline", kind: "токен" }
  - { name: "--dur-3", kind: "токен" }
  - { name: "--ease-out", kind: "токен" }
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

## Контракт

| Что | Обязательно | Почему |
|---|---|---|
| `role="progressbar"` на `inst-meter` | да | Без роли значение существует только в пикселях ширины и недоступно ни скринридеру, ни тесту |
| `aria-valuenow` `aria-valuemin` `aria-valuemax` | да, у определённой | Определённая полоса обязана сообщать число, а не картинку |
| `aria-label` или `aria-labelledby` | да | У полосы нет своего текста. Подпись из `inst-meter-row` связывается через `aria-labelledby`, если она видима |
| **отсутствие** `aria-valuenow` у `indeterminate` | да | Ровно так объявляется «занято, длительность неизвестна» |
| `inst-meter-fill` внутри | да | Дорожка сама ничего не заливает |

`style="inline-size:43%"` на заливке — **не оформление, а значение**: величина
приходит с сервера и физически не может жить в статическом CSS. Всё остальное —
высота, радиус, дорожка, цвет — берётся из библиотеки. Заливка едет с переходом
(`--dur-3`), поэтому обновление значения читается как движение, а не как
скачок.

Протяжённость дорожки показывает **внутренний хайрлайн**. Сама дорожка при
этом остаётся светлой, и заливка держит на ней 3:1.

### Доступность

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

## Варианты

### Определённая и неопределённая

Это **два разных компонента**. Определённая мера обещает остаток,
неопределённая — только факт работы.

```html preview
<div class="inst-meter" role="progressbar" aria-label="Расход токенов"
     aria-valuenow="74" aria-valuemin="0" aria-valuemax="100">
  <div class="inst-meter-fill" style="inline-size:74%"></div>
</div>
<div class="inst-meter" data-state="indeterminate" role="progressbar" aria-label="Индексация"
     aria-valuemin="0" aria-valuemax="100">
  <div class="inst-meter-fill"></div>
</div>
```

| | Определённая | Неопределённая |
|---|---|---|
| Атрибут | без `data-state` | `data-state="indeterminate"` |
| Ширина заполнения | приходит из данных | задаёт CSS — 35%, и она метётся |
| `aria-valuenow` | **обязателен** | **отсутствует, и это сообщение** |
| Что обещает читателю | остаток и темп | что работа идёт |

Отсутствие `aria-valuenow` — не упущение разметки, а её содержание: по нему
скринридер объявляет «занято» вместо ложного процента. Поставить `valuenow`
неопределённой полосе — то же самое, что нарисовать её застрявшей на 90%.

### Тон

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

`data-tone` ставится на `inst-meter`, заливка читает `--tone-mark` и больше
ничего о тоне не знает. Тон здесь — **исключение, а не приправа**: полоса без
атрибута берёт акцент, и это правильный ответ в подавляющем большинстве
случаев. Тон меняется, когда меняется **смысл величины** — перешли порог,
исчерпали лимит.

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
  <div class="inst-meter" data-state="indeterminate" role="progressbar" aria-label="Индексация"
     aria-valuemin="0" aria-valuemax="100">
    <div class="inst-meter-fill"></div>
  </div>
</div>
```

`inst-meter-row` разносит подпись и величину по краям и стоит **над** полосой:
число слева от полосы съедало бы её длину, а длина здесь и есть данные.

### В бюджете прогона

```html preview context
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

## API

```api
```

## Связанное

```related
```
