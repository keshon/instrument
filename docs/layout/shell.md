---
title: Оболочка приложения
group: Раскладка
status: stable
source: src/layout.css
api:
  - { name: "inst-shell", kind: "класс", doc: "Грид на всю высоту окна (`100dvh`)" }
  - { name: "inst-shell-aside", kind: "класс", doc: "Боковая колонка" }
  - { name: "inst-shell-header", kind: "класс", doc: "Шапка" }
  - { name: "inst-shell-main", kind: "класс", doc: "Рабочая область. Контейнер `app`" }
  - { name: "inst-skip-link", kind: "класс", doc: "Пропуск к содержимому. Виден только по `:focus-visible`" }
  - { name: "--aside-w", kind: "токен" }
  - { name: "--pad-panel", kind: "токен" }
  - { name: "--space-3", kind: "токен" }
  - { name: "--space-4", kind: "токен" }
  - { name: "--space-5", kind: "токен" }
  - { name: "--gap-inline", kind: "токен" }
  - { name: "--hairline", kind: "токен" }
  - { name: "--border", kind: "токен" }
  - { name: "--surface-sunken", kind: "токен" }
  - { name: "--surface-page", kind: "токен" }
  - { name: "--surface-raised", kind: "токен" }
  - { name: "--radius-md", kind: "токен" }
  - { name: "--text-sm", kind: "токен" }
  - { name: "--z-modal", kind: "токен" }
title-en: "Application shell"
group-en: "Layout"
---

Каркас экрана: боковая колонка, шапка, рабочая область. Каждая часть
прокручивается сама — боковая навигация не уезжает, пока читают таблицу.

```html preview
<div class="inst-shell">
  <aside class="inst-shell-aside">
    <nav class="inst-nav" aria-label="Разделы">
      <div class="inst-nav-group">
        <span class="inst-nav-label">Работа</span>
        <a class="inst-nav-item" href="#" aria-current="page">Прогоны
          <span class="inst-badge inst-nav-count">7</span></a>
        <a class="inst-nav-item" href="#">Очередь</a>
      </div>
    </nav>
  </aside>

  <header class="inst-shell-header">
    <nav aria-label="Хлебные крошки">
      <ol class="inst-crumbs">
        <li><a href="#">Проекты</a></li>
        <li><span aria-current="page">Прогон #4127</span></li>
      </ol>
    </nav>
    <span class="inst-cluster-spacer"></span>
    <button class="inst-btn inst-btn--sm" type="button">Экспорт</button>
    <button class="inst-btn inst-btn--sm inst-btn--primary" type="button">Запустить</button>
  </header>

  <main class="inst-shell-main" id="main">
    <div class="inst-container inst-container--full inst-stack">
      <div class="inst-page-header">
        <div class="inst-page-header-main">
          <h1 class="inst-page-title">Прогоны</h1>
        </div>
      </div>
    </div>
  </main>
</div>
```

## Установка

```html
<link rel="stylesheet" href="src/kit.css">
```

## Когда использовать

| Используйте | Возьмите другое |
|---|---|
| Каркас всего приложения: одна оболочка на документ | **Внутренняя двухколоночная раскладка** — [сплит](./split.md): он переносится сам и не занимает всю высоту окна |
| Нужны независимо прокручиваемые области | **Ограничение ширины содержимого** — [контейнер](./container.md) внутри рабочей области |
| Постоянная боковая навигация — [nav](../components/navigation/nav.md) | **Область приложения внутри экрана** — [панель](../components/display/panel.md) |
| — | **Вертикальный ритм внутри области** — [примитивы потока](./flow.md) |

## Разметка

```html
<a class="inst-skip-link" href="#main">К содержанию</a>

<div class="inst-shell">
  <aside class="inst-shell-aside">…</aside>
  <header class="inst-shell-header">…</header>
  <main class="inst-shell-main" id="main">…</main>
</div>
```

Порядок в разметке — `aside`, `header`, `main`; на экране их расставляет
`grid-template-areas`, поэтому шапка стоит справа от колонки, а не над ней. На
узком экране порядок областей меняется тем же свойством — без второго набора
разметки.

Пропуск к содержимому стоит **до** оболочки: это первое, что встречает
клавиатура, и `href` обязан указывать на `id` рабочей области.

## Области

| Область | Что делает |
|---|---|
| `inst-shell-aside` | Боковая колонка шириной `--aside-w`. Прокручивается сама, лежит на `--surface-sunken` |
| `inst-shell-header` | Шапка. Ряд с обязательным переносом: крошки и контролы на узком экране уходят на вторую строку, а не выталкивают документ за вьюпорт |
| `inst-shell-main` | Рабочая область. Прокручивается сама, даёт вертикальный воздух и объявлена **контейнером** `app` |

Горизонтальный воздух даёт не рабочая область, а [контейнер](./container.md)
внутри неё: у них разная работа — один ограничивает меру чтения, другая держит
прокрутку.

## Узкий экран

```css
@media (width < 60rem) { … }
```

**Единственный порог во всём ките.** Оболочка и есть то единственное, что
действительно зависит от размера окна: всё остальное перестраивается интринсиком
или отвечает на ширину своей области через `@container` — см.
[конституцию](../about/design-principles.md).

Ниже 60rem есть **два поведения**, и выбор между ними не вкусовой, а по
**длине навигации**.

### Умолчание: горизонтальная полоса

Колонка становится лентой **над** содержимым:

| Что меняется | Как |
|---|---|
| Порядок областей | `header`, `aside`, `main` — колонкой |
| Колонка | Ряд с горизонтальной прокруткой, нижняя линия вместо боковой |
| `inst-nav`, `inst-nav-group` | Разворачиваются в ряд |
| `inst-nav-label` | Скрывается: подписи групп в полосе не помещаются |
| `inst-nav-item` | Перестаёт переносить подпись |

Работает **без единой строки JS** и потому не может «не открыться»: у
приложения, которое ничего не подключило, навигация на телефоне всё равно
есть. Годится примерно **до восьми пунктов**.

### Вариант: выезжающий ящик

```html
<aside class="inst-shell-aside inst-shell-aside--drawer">…</aside>

<button class="inst-btn inst-btn--icon inst-shell-nav-toggle" type="button"
        aria-label="Разделы" aria-expanded="false" aria-controls="sidebar">…</button>
```

Как только пунктов больше восьми, лента перестаёт работать: семьдесят один
пункт — это лента шириной почти **восемь тысяч пикселей** без заголовков
групп. Это не «неудобно», это отсутствие навигации, и найдено оно было на
собственном справочнике этого кита.

Ящик сохраняет вертикаль и подписи групп, выезжает из-за начального края и
гасит фон той же подложкой `--scrim`, что и модалка.

| Класс | Работа |
|---|---|
| `inst-shell-aside--drawer` | Колонка уходит в фиксированный слой и выезжает по требованию |
| `inst-shell-nav-toggle` | Кнопка вызова. На широком экране скрыта: колонка и так на месте, и кнопка была бы обманом |

**Что даёт приложение.** Одну строку — переключить `aria-expanded` на кнопке.
Ровно тот же контракт, что у [модалки](../components/overlays/dialog.md) с
`showModal()`. Ловушка фокуса, `Escape` и клик по подложке — тоже слой
приложения: кит рисует состояние, но не выдумывает его.

> Внутренние правила ленты намеренно исключают ящик
> (`:not(.inst-shell-aside--drawer)`). Пока они применялись ко всякой узкой
> колонке, приложение, показавшее её иначе, получало вертикальную панель с
> навигацией, выложенной **в строку**, и без подписей групп: два механизма
> на одну работу расходятся всегда.

## Собранный экран

Оболочка сама по себе — три пустые области. Работу видно только когда в них
что-то стоит: боковая навигация с активным пунктом, крошки и действия в
шапке, а в рабочей области — шапка экрана, вкладки и секции.

Главное в этом примере не видно на картинке, но проверяется руками:
**каждая область прокручивается сама**. Прокрутите таблицу — боковая
колонка и шапка останутся на месте.

```html preview context
<div class="inst-shell">
  <aside class="inst-shell-aside">
    <nav class="inst-nav" aria-label="Разделы">
      <div class="inst-nav-group">
        <span class="inst-nav-label">Работа</span>
        <a class="inst-nav-item" href="#" aria-current="page">
          <svg class="inst-icon" aria-hidden="true"><use href="#i-grid"/></svg>Прогоны
          <span class="inst-badge inst-nav-count">7</span></a>
        <a class="inst-nav-item" href="#">
          <svg class="inst-icon" aria-hidden="true"><use href="#i-list"/></svg>Очередь</a>
        <a class="inst-nav-item" href="#">
          <svg class="inst-icon" aria-hidden="true"><use href="#i-chart"/></svg>Отчёты</a>
      </div>
      <div class="inst-nav-group">
        <span class="inst-nav-label">Настройка</span>
        <a class="inst-nav-item" href="#">
          <svg class="inst-icon" aria-hidden="true"><use href="#i-users"/></svg>Доступы</a>
        <a class="inst-nav-item" href="#">
          <svg class="inst-icon" aria-hidden="true"><use href="#i-settings"/></svg>Параметры</a>
      </div>
    </nav>
  </aside>

  <header class="inst-shell-header">
    <nav aria-label="Хлебные крошки">
      <ol class="inst-crumbs">
        <li><a href="#">Проекты</a></li>
        <li><a href="#">worldgen</a></li>
        <li><span aria-current="page">Прогон #4127</span></li>
      </ol>
    </nav>
    <span class="inst-cluster-spacer"></span>
    <button class="inst-btn inst-btn--sm" type="button">Экспорт</button>
    <button class="inst-btn inst-btn--sm inst-btn--primary" type="button">Запустить</button>
  </header>

  <main class="inst-shell-main">
    <div class="inst-container inst-container--full inst-stack">
      <div class="inst-page-header">
        <div class="inst-page-header-main">
          <h3 class="inst-page-title">Прогоны</h3>
          <p class="inst-page-desc">Каждая область прокручивается сама: боковая навигация
            не уезжает, пока читают таблицу.</p>
        </div>
        <div class="inst-page-actions">
          <button class="inst-btn inst-btn--sm" type="button">Фильтры</button>
        </div>
      </div>

      <nav class="inst-tabs" aria-label="Разделы прогонов">
        <a class="inst-tab" href="#" aria-current="page">Активные</a>
        <a class="inst-tab" href="#">История</a>
        <a class="inst-tab" href="#">Расписание</a>
      </nav>

      <div class="inst-section">
        <div class="inst-section-head">
          <h4 class="inst-section-title">Сводка за сутки</h4>
          <span class="inst-section-actions"><button class="inst-btn inst-btn--sm inst-btn--ghost" type="button">Период</button></span>
        </div>
        <div class="inst-grid inst-grid--tight">
          <div class="inst-metric"><div class="inst-metric-label">В работе</div><div class="inst-metric-value">7</div></div>
          <div class="inst-metric"><div class="inst-metric-label">В очереди</div><div class="inst-metric-value">5</div></div>
          <div class="inst-metric"><div class="inst-metric-label">Упало</div><div class="inst-metric-value">1</div></div>
        </div>
      </div>
    </div>
  </main>
</div>
```

## Справочник

```api
```

### Обязательная разметка

| Что | Почему |
|---|---|
| `inst-skip-link` первым в `<body>` | Иначе клавиатура каждый раз проходит всю навигацию заново |
| `id` на `inst-shell-main` | Цель пропуска |
| `<aside>` / `<header>` / `<main>` | Ориентиры страницы. `<div>` их не даёт |
| `aria-label` у каждой `<nav>` внутри | В оболочке их обычно три: разделы, крошки, вкладки |

Ширина колонки перенастраивается одной строкой:

```css
.my-app { --aside-w: 280px; }
```

## Доступность

| | |
|---|---|
| Пропуск к содержимому | Первый элемент в порядке обхода. Скрыт до фокуса `clip-path`, а не `display: none`: скрытое так остаётся фокусируемым |
| Ориентиры | `<aside>`, `<header>`, `<main>` дают навигацию по ориентирам без единого `role` |
| Прокрутка | Три независимые области с `overflow: auto`. `overscroll-behavior` не ставится: контейнеру, которому нечего прокручивать, `contain` съедает колесо мыши |
| Высота | `100dvh`, а не `100vh` — на мобильных панель браузера иначе срезает нижнюю часть рабочей области |
| Перенос в шапке | Обязателен: без него узкий экран получает горизонтальную прокрутку документа, и `100dvh` перестаёт помещаться по высоте |
| Печать | Боковая колонка скрывается, оболочка разворачивается в обычный поток, прокручиваемые области печатаются целиком |

## Связанное

[Контейнер](./container.md) · [Примитивы потока](./flow.md) ·
[Шапка экрана](./page-header.md) · [Сплит](./split.md) ·
[Навигация](../components/navigation/nav.md) ·
[Крошки](../components/navigation/breadcrumbs.md) ·
[Панель](../components/display/panel.md)
