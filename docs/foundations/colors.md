---
title: Цвет
group: Основания
status: stable
source: src/tokens.css
---

Два яруса. Внизу **рампы** — сырые шаги, они не меняются никогда. Сверху
**семантика** — для чего цвет нужен. Компонент видит только семантику: тот, кто
написал `--n-3`, только что захардкодил светлую тему.

```html preview
<div class="inst-panel">
  <div class="inst-panel-header">
    <span class="inst-panel-title">Прогон #4127</span>
    <span class="inst-badge" data-tone="running"><span class="inst-dot"></span>идёт</span>
  </div>
  <div class="inst-panel-body inst-stack">
    <input class="inst-input" type="text" value="terrain_chunk_04">
    <div class="inst-note" data-tone="warn">Три теста упали после последнего прохода.</div>
  </div>
</div>
```

## Когда использовать

| Используйте | Возьмите другое |
|---|---|
| Семантический токен: `--text-primary`, `--surface-raised`, `--ok-text` | **Шаг рампы напрямую** (`--n-3`, `--a-4`) — это светлая тема, вписанная в компонент. Список семантики целиком — [токены](./tokens.md) |
| Состояние объекта — атрибут `data-tone`, один словарь на весь кит | **Свой цвет под своё состояние** — заведите тон и покажите его [бейджем](../components/display/badge.md): цвет не имеет права быть единственным носителем |
| Ряд графика — `--chart-1` … `--chart-6` по порядку | **Статусный тон как «ещё одна серия»** — тона `ok`/`warn`/`error` зарезервированы. Ряды подписывает [легенда](../components/charts/legend.md) |
| Глубина — перепад поверхностей | **Тень ради глубины** — тень означает «временное и сверху», см. [высоту](./elevation.md) |
| Информационное сообщение — тон `neutral` | **Синий как «info»** — синий занят акцентом и состоянием «идёт», см. [сноску в карточке](../components/display/card.md) |

## Рампы

### Нейтраль — 15 шагов

```html preview
<div class="ramp">
  <span class="ramp-step" style="--c: var(--n-0)"></span>
  <span class="ramp-step" style="--c: var(--n-1)"></span>
  <span class="ramp-step" style="--c: var(--n-2)"></span>
  <span class="ramp-step" style="--c: var(--n-3)"></span>
  <span class="ramp-step" style="--c: var(--n-4)"></span>
  <span class="ramp-step" style="--c: var(--n-5)"></span>
  <span class="ramp-step" style="--c: var(--n-6)"></span>
  <span class="ramp-step" style="--c: var(--n-7)"></span>
  <span class="ramp-step" style="--c: var(--n-8)"></span>
  <span class="ramp-step" style="--c: var(--n-9)"></span>
  <span class="ramp-step" style="--c: var(--n-10)"></span>
  <span class="ramp-step" style="--c: var(--n-11)"></span>
  <span class="ramp-step" style="--c: var(--n-12)"></span>
  <span class="ramp-step" style="--c: var(--n-13)"></span>
  <span class="ramp-step" style="--c: var(--n-14)"></span>
</div>
<div class="ramp ramp-scale">
  <span>0</span><span>1</span><span>2</span><span>3</span><span>4</span>
  <span>5</span><span>6</span><span>7</span><span>8</span><span>9</span>
  <span>10</span><span>11</span><span>12</span><span>13</span><span>14</span>
</div>
```

Переключите тему стола: рампа **не изменится**. Она сырая и не знает о темах —
меняется только то, какие её шаги берёт семантика.

Одна ручка задаёт направление уклона: `--hue-neutral: 75` — тёплый, `250` —
холодный. Цветность 0.002–0.006: ниже порога осознанного замечания, и в этом её
работа. Пятнадцатый шаг существует потому, что тёмному концу нужно четыре
различимые поверхности подряд, а на четырнадцати их помещалось три.

| Шаг | Светлота | Где занят |
|---|---|---|
| `--n-0` | 0.994 | `--surface-raised`, `--surface-overlay` (светлая), `--accent-on` |
| `--n-1` | 0.978 | `--surface-page`, `--surface-field` (светлая), `--text-primary` (тёмная) |
| `--n-2` | 0.958 | `--surface-sunken` (светлая) |
| `--n-3` | 0.928 | Резерв |
| `--n-4` | 0.884 | Резерв |
| `--n-5` | 0.806 | `--text-secondary` (тёмная) |
| `--n-6` | 0.706 | `--text-muted` (тёмная) |
| `--n-7` | 0.606 | `--text-faint` (обе темы) |
| `--n-8` | 0.508 | `--text-muted` (светлая) |
| `--n-9` | 0.416 | `--text-secondary` (светлая) |
| `--n-10` | 0.322 | `--surface-overlay` (`dark-soft`) |
| `--n-11` | 0.242 | `--surface-overlay` (тёмная), `--surface-raised` (`dark-soft`) |
| `--n-12` | 0.196 | `--text-primary` (светлая), `--surface-raised` (тёмная), `--surface-page` (`dark-soft`) |
| `--n-13` | 0.155 | `--surface-page` (тёмная), `--surface-sunken` и `--surface-field` (`dark-soft`) |
| `--n-14` | 0.120 | `--surface-sunken`, `--surface-field` (тёмная) |

Шаг без семантики — **резерв под приложение и графики, а не мусор**: рампа
непрерывна по определению.

### Акцент — один тон

```html preview
<div class="ramp">
  <span class="ramp-step" style="--c: var(--a-1)"></span>
  <span class="ramp-step" style="--c: var(--a-2)"></span>
  <span class="ramp-step" style="--c: var(--a-3)"></span>
  <span class="ramp-step" style="--c: var(--a-4)"></span>
  <span class="ramp-step" style="--c: var(--a-5)"></span>
  <span class="ramp-step" style="--c: var(--a-6)"></span>
</div>
<div class="ramp ramp-scale">
  <span>1</span><span>2</span><span>3</span><span>4</span><span>5</span><span>6</span>
</div>
```

Второй акцентный тон не заводится. Это не запрос фичи, а сигнал, что смысл
должно нести что-то другое.

| Шаг | Светлота | Где занят |
|---|---|---|
| `--a-1` | 0.960 | `--accent-bg` (светлая) |
| `--a-2` | 0.890 | Резерв |
| `--a-3` | 0.700 | `--accent-text`, `--accent-mark`, `--accent-border`, `--focus-ring` (тёмная) |
| `--a-4` | 0.560 | `--accent-solid`, `--accent-mark`, `--accent-border`, `--focus-ring` (светлая) |
| `--a-5` | 0.450 | `--accent-text` (светлая), `--accent-hover` |
| `--a-6` | 0.330 | Резерв |

Тон — 250°. Светлота `--accent-solid` ограничена **сверху**: светлее `0.56`
белая подпись перестаёт держать 4.5:1.

### Статусы — три зарезервированных тона

```html preview
<div class="ramp">
  <span class="ramp-step" style="--c: var(--ok-1)"></span>
  <span class="ramp-step" style="--c: var(--ok-3)"></span>
  <span class="ramp-step" style="--c: var(--ok-4)"></span>
  <span class="ramp-step" style="--c: var(--ok-5)"></span>
  <span class="ramp-step" style="--c: var(--ok-6)"></span>
</div>
<div class="ramp">
  <span class="ramp-step" style="--c: var(--warn-1)"></span>
  <span class="ramp-step" style="--c: var(--warn-3)"></span>
  <span class="ramp-step" style="--c: var(--warn-4)"></span>
  <span class="ramp-step" style="--c: var(--warn-5)"></span>
  <span class="ramp-step" style="--c: var(--warn-6)"></span>
</div>
<div class="ramp">
  <span class="ramp-step" style="--c: var(--err-1)"></span>
  <span class="ramp-step" style="--c: var(--err-3)"></span>
  <span class="ramp-step" style="--c: var(--err-4)"></span>
  <span class="ramp-step" style="--c: var(--err-5)"></span>
  <span class="ramp-step" style="--c: var(--err-6)"></span>
</div>
<div class="ramp ramp-scale">
  <span>1</span><span>3</span><span>4</span><span>5</span><span>6</span>
</div>
```

Шага 2 в статусных рядах нет: он не понадобился ни одной роли, а рампа
заводится под роли, а не под ровный счёт.

`ok` — 150°, `warn` — 85°, `err` — 25°. Никогда не используются как украшение
или как «четвёртая серия на графике».

| Шаг | Роль | Почему именно так |
|---|---|---|
| 1 | Тонированный фон в **светлой** теме | `--ok-1` `--warn-1` `--err-1` |
| 3 | Текст в **тёмных** темах | Светлее соседей по ряду: `dark-soft` поднимает поверхности на шаг, и шаг 3 обязан остаться читаемым |
| 4 | Заливка и индикатор в **светлой** | `--ok-4` `--warn-4` `--err-4` |
| 5 | Текст в **светлой**, поверх шага 1 | Шаг 4 не вытягивает 4.5:1 на собственном тонированном фоне |
| 6 | Резерв под графики | |

Шаг 5 существует именно потому, что шаг 4 не проходит порог на самом себе.
Пока его не было, жёлтый бейдж жил на контрасте 3.02.

## Семантика

Одно объявление на токен через `light-dark()`. Второй темы как отдельного блока
не существует — значит, и расходиться нечему.

### Поверхности

```html preview
<div class="swatches">
  <div class="sw"><div class="sw-chip" style="--c: var(--surface-sunken)"></div><div class="sw-cap">--surface-sunken</div></div>
  <div class="sw"><div class="sw-chip" style="--c: var(--surface-page)"></div><div class="sw-cap">--surface-page</div></div>
  <div class="sw"><div class="sw-chip" style="--c: var(--surface-raised)"></div><div class="sw-cap">--surface-raised</div></div>
  <div class="sw"><div class="sw-chip" style="--c: var(--surface-overlay)"></div><div class="sw-cap">--surface-overlay</div></div>
  <div class="sw"><div class="sw-chip" style="--c: var(--surface-field)"></div><div class="sw-cap">--surface-field</div></div>
</div>
```

Переключите тему стола — вот здесь всё и меняется. Рампа осталась прежней,
а семантика взяла другие её шаги: в тёмной теме перепад между соседними
поверхностями шире, потому что тёмная тема **не инверсия**, а авторская.

Читаются как стопка. Глубина передаётся **порядком светлоты и никогда тенью**.

| Токен | Работа |
|---|---|
| `--surface-sunken` | Дно: врез, дорожка, шапка таблицы |
| `--surface-page` | Фон документа |
| `--surface-raised` | Панель, карточка, шапка оболочки |
| `--surface-overlay` | Поповер, меню, тултип, модалка |
| `--surface-field` | Врез под ввод. Отдельная роль, а не переиспользование `raised`: пока поле красилось цветом панели, его единственной границей была рамка на 1.31:1 |
| `--surface-hover` `--surface-active` | Альфа поверх чего угодно |
| `--surface-selected` | Выделение. Тоже альфа — чтобы наведение по выделенной строке оставалось видно |

### Текст

| Токен | Порог | Для чего |
|---|---|---|
| `--text-primary` | 4.5:1 | Данные, заголовки, подписи контролов |
| `--text-secondary` | 4.5:1 | Второстепенное: описания, единицы |
| `--text-muted` | 4.5:1 | Метаданные и таймстемпы — это **данные**, им положен читаемый порог |
| `--text-faint` | 3:1 | **Не применяется к тексту, который читают.** Порог декорации |

### Рамки

Несущая и декоративная — разные токены, потому что к ним разные требования.

| Токен | Когда |
|---|---|
| `--border-subtle` | Разделитель внутри уже ограниченной области |
| `--border` | Панель, карточка, ячейка: рядом есть перепад поверхностей |
| `--border-strong` | Скроллбар, акцентированный шов |
| `--border-control` | Граница, которая **и есть контрол**: чекбокс, поле, дорожка свитча. Обязана держать 3:1 |

### Тоны

У смыслового тона **ровно два токена**, и это результат расчёта, а не вкуса.

| Тон | Метка | Фон |
|---|---|---|
| акцент | `--accent-text` (4.5:1) и `--accent-mark` (3:1) | `--accent-bg` |
| ok | `--ok-text` | `--ok-bg` |
| warn | `--warn-text` | `--warn-bg` |
| error | `--err-text` | `--err-bg` |

Третьего токена — «заливки» между текстом и фоном — нет намеренно: он жил на
шаге 4 и не проходил нигде, где под ним оказывалась дорожка (жёлтая заливка
меры давала 2.49 при норме 3.0). Исключение одно — акцент: у него есть случай,
которого нет у статусов, заливка **под белой подписью** (`--accent-solid` +
`--accent-on`).

### Тон как атрибут

`data-tone` ставится один раз на группу; вложенные элементы читают
`--tone-ink`, `--tone-mark` и `--tone-bg` и больше ничего о тоне не знают.

| Значение | Значит |
|---|---|
| `neutral` | Явно нейтральное. Оно же исполняет роль info |
| `running` | Идёт сейчас. Носитель — пульсация, тон вторичен |
| `ok` | Успешно завершено |
| `warn` | Завершено с замечаниями |
| `error` | Упало |

Словарь закрыт. Шестое значение не заводится.

## Четыре темы

Тема — не набор токенов, а комбинация **двух независимых ручек**: тон нейтрали
и глубина тёмных поверхностей. Поэтому пятая тема — это одна строка, а не копия
семантики.

| Атрибут | `color-scheme` | `--hue-neutral` | Что ещё |
|---|---|---|---|
| нет | по системе | 75 | Слушает настройку ОС |
| `data-theme="light"` | `light` | 75 | Тёплая светлая |
| `data-theme="light-cool"` | `light` | 250 | Холодная светлая |
| `data-theme="dark"` | `dark` | 75 | Тёплая тёмная |
| `data-theme="dark-soft"` | `dark` | 250 | Вся стопка поверхностей поднимается на шаг: переопределены ровно пять токенов |

```html
<html data-theme="dark">   <!-- принудительно -->
<html>                     <!-- по системной настройке -->
```

Светлая глубина не регулируется намеренно: сверху рампы шаги идут плотно, и
«менее светлая» светлая тема — это просто серый фон, а не тема.

## Категориальная палитра

```html preview
<div class="ramp">
  <span class="ramp-step" style="--c: var(--chart-1)"></span>
  <span class="ramp-step" style="--c: var(--chart-2)"></span>
  <span class="ramp-step" style="--c: var(--chart-3)"></span>
  <span class="ramp-step" style="--c: var(--chart-4)"></span>
  <span class="ramp-step" style="--c: var(--chart-5)"></span>
  <span class="ramp-step" style="--c: var(--chart-6)"></span>
</div>
<div class="ramp ramp-scale">
  <span>1</span><span>2</span><span>3</span><span>4</span><span>5</span><span>6</span>
</div>
```

Светлота внутри ряда намеренно **разная**, а не выровненная: так ряды
остаются различимыми и при дальтонизме, и на чёрно-белой печати, где тон
исчезает вовсе. Прищурьтесь — порядок по светлоте всё ещё читается.

Единственное место, где цвет кодирует **ряд**, а не состояние. Поэтому у неё
свои правила.

| Токен | Тон | Ряд |
|---|---|---|
| `--chart-1` | 280° | Первый |
| `--chart-2` | 320° | Второй |
| `--chart-3` | 355° | Третий |
| `--chart-4` | 55° | Четвёртый |
| `--chart-5` | 115° | Пятый |
| `--chart-6` | 190° | Шестой |

- **Порядок — часть контракта.** Ряд №1 всегда `--chart-1`, иначе один и тот же
  показатель меняет цвет от экрана к экрану.
- **Минимум 25° от каждого статусного тона** (25, 85, 150) и от акцента (250):
  категориальный цвет не должен прочитаться как «упало».
- **Светлота внутри палитры разная намеренно** — ряды остаются различимыми при
  дальтонизме и на чёрно-белой печати, где тон исчезает вовсе.
- **Шесть — потолок.** Седьмой ряд означает, что это не график, а
  [таблица](../components/display/table.md).
- **Вне графика и его легенды запрещена.** В роли статуса, заливки кнопки или
  подсветки строки категориального цвета не бывает.

## Прочее

| Токен | Работа |
|---|---|
| `--focus-ring` | Кольцо фокуса, одно на весь кит |
| `--track` | Пустая дорожка меры и слайдера. Намеренно тихая: её работа — дать заливке отделиться, а не отделиться самой |
| `--scrim` | Подложка модалки. В тёмной теме гуще. Затемнение, а не размытие |
| `--shadow-color-near` `--shadow-color-far` | Цвета двух теней — см. [высоту](./elevation.md) |

## Доступность

| | |
|---|---|
| Проверка | `go -C tools run ./cmd/contrast` читает настоящий `tokens.css` и резолвит `light-dark()`, `color-mix()` и `var()` так же, как браузер |
| Порог текста | 4.5:1 во всех четырёх темах. `--text-faint` из этого исключён и потому не носит текст |
| Порог метки | 3:1, но против **двух** фонов сразу: поверхности и дорожки. Отсюда `--tone-mark` отдельно от `--tone-ink` |
| Цвет не единственный носитель | Статус несёт точку **и** слово, сноска несёт иконку, строка дифа несёт знак |
| Тёмная тема | Не инверсия: контраст между соседними поверхностями расширяется, цветность опускается, рамки переходят с тёмной альфы на светлую |
| Режим принудительных цветов | Носители значения переживают сброс через `forced-color-adjust: none` |

## Связанное

[Токены](./tokens.md) · [Высота и поверхности](./elevation.md) ·
[Бейдж](../components/display/badge.md) · [Легенда графика](../components/charts/legend.md) ·
[Метрика](../components/display/metric.md)

Исходник: `src/tokens.css` · Почему кит устроен именно так — [конституция](../about/design-principles.md)
