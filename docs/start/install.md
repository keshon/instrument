---
title: Установка
group: Начало
source: dist/instrument.min.css
title-en: "Install"
group-en: "Start"
---

Библиотека подключается двумя тегами и не требует сборки. Ниже — целая
страница, которую можно скопировать и открыть.

## Установка

Три способа. Файл во всех трёх один и тот же.

**npm**

```bash
npm i @keshon/instrument
```

```js
import '@keshon/instrument/min';        // dist/instrument.min.css
import '@keshon/instrument/js';         // поведение, по желанию
```

Пакет также экспортирует несжатый CSS (`@keshon/instrument`), спрайт
(`@keshon/instrument/sprite.svg`) и исходники (`@keshon/instrument/src/*`).

**CDN**

```html
<link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/@keshon/instrument@0.6.2/dist/instrument.min.css">
```

Версия в адресе указывается точно. `@latest` кэшируется на неделю и может
обновиться в произвольный момент.

**Файлом**

```html
<link rel="stylesheet" href="instrument.min.css">
```

Возьмите `dist/instrument.min.css` — 106 КБ, 17 КБ в gzip. Это вся библиотека:
токены, шесть тем, пять масштабов, три плотности, все компоненты. `dist/instrument.css` — то же
самое без сжатия и с комментариями.

## Иконки

Иконка рисуется ссылкой на символ: `<use href="#i-refresh">`. Символы должны
быть в самом документе: ссылка вида `href="https://cdn…/sprite.svg#i-refresh"`
подчиняется правилу одного источника и не отрисуется — без ошибки в консоли.

Вставьте содержимое `assets/sprite.svg` первым элементом внутри `<body>` — один
раз на страницу.

```html
<body>
  <svg class="inst-sprite" aria-hidden="true" focusable="false" xmlns="http://www.w3.org/2000/svg">
    <!-- сюда содержимое assets/sprite.svg -->
  </svg>
```

Без спрайта работает всё остальное, иконки просто не рисуются. Полный список из
99 символов — на странице [иконок](../foundations/icons.md).

## Поведение с клавиатуры

```html
<script type="module" src="instrument.js"></script>
```

Отдельный файл и по желанию. Он выполняет то, что библиотека обещает разметкой:
стрелки в списках и вкладках, копирование, снятие тега, перетаскивание оси.
Оформление он не трогает. Подробности — в
[поведении](../foundations/behavior.md).

## Целая страница

```html
<!doctype html>
<html lang="ru">
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <title>Моё приложение</title>
  <link rel="stylesheet" href="instrument.min.css">
  <script type="module" src="instrument.js"></script>
</head>
<body>
  <div class="inst-shell">
    <header class="inst-shell-header">
      <strong>Моё приложение</strong>
      <button class="inst-btn inst-btn--primary" type="button">Запустить</button>
    </header>

    <main class="inst-shell-main">
      <div class="inst-container">
        <div class="inst-panel">
          <div class="inst-panel-header">
            <span class="inst-panel-title">Очередь</span>
          </div>
          <div class="inst-panel-body">
            <div role="listbox" aria-label="Очередь задач">
              <div class="inst-task" role="option" data-state="running"
                   aria-selected="true" tabindex="0">
                <span class="inst-task-gutter"><span class="inst-dot"></span></span>
                <span class="inst-task-main">
                  <span class="inst-task-title">Разбор карты высот</span>
                  <span class="inst-task-sub">terrain/heightmap.ts</span>
                </span>
                <span class="inst-task-meta">2,1 с</span>
              </div>
            </div>
          </div>
        </div>
      </div>
    </main>
  </div>
</body>
</html>
```

Вот как это выглядит:

```html preview context
<div class="inst-panel">
  <div class="inst-panel-header">
    <span class="inst-panel-title">Очередь</span>
  </div>
  <div class="inst-panel-body">
    <div role="listbox" aria-label="Очередь задач">
      <div class="inst-task" role="option" data-state="running" aria-selected="true" tabindex="0">
        <span class="inst-task-gutter"><span class="inst-dot"></span></span>
        <span class="inst-task-main">
          <span class="inst-task-title">Разбор карты высот</span>
          <span class="inst-task-sub">terrain/heightmap.ts</span>
        </span>
        <span class="inst-task-meta">2,1 с</span>
      </div>
      <div class="inst-task" role="option" data-state="queued" aria-selected="false" tabindex="-1">
        <span class="inst-task-gutter"><span class="inst-dot"></span></span>
        <span class="inst-task-main">
          <span class="inst-task-title">Сборка биомов</span>
        </span>
        <span class="inst-task-meta">—</span>
      </div>
    </div>
  </div>
</div>
```

## Тема и плотность

Оба — атрибуты, и оба работают на любом поддереве, а не только на документе.

```html
<html data-theme="dark">              <!-- шесть тем; без атрибута — по настройке ОС -->
<html data-accent="clay">             <!-- петроль (умолчание) · graphite · indigo · clay -->
<section data-density="compact">      <!-- compact · comfortable -->
```

| Тема | |
|---|---|
| `light-neutral` `light` `light-cool` | светлые: нейтральная, тёплая, холодная |
| `dark-light` `dark-soft` `dark` | тёмные, различаются глубиной поверхностей |

Без атрибута тема выбирается по настройке ОС. Что именно меняет каждая ручка —
в [цвете](../foundations/colors.md) и [плотности](../foundations/density.md).

## Ваши стили выигрывают

Библиотека целиком лежит в каскадных слоях, а ваши стили — вне их. Поэтому любое
ваше правило перебивает её, и `!important` не нужен нигде.

```css
.inst-btn { border-radius: 0; }   /* сработает, специфичность неважна */
```

Порядок слоёв — в [токенах](../foundations/tokens.md#порядок-слоёв).

## Требования

Chrome 123+, Safari 17.5+, Firefox 120+ — по `light-dark()`, `oklch()`,
каскадным слоям, вложенности, `:has()` и единице `lh`. Полифилов нет.

## Дальше

[Консоль прогонов](../blocks/console.md) — собранный экран, из которого видно,
как компоненты стыкуются.

[Оболочка приложения](../layout/shell.md) — шапка, боковая панель, основная
область.

[Принципы дизайна](../about/design-principles.md) — правила, из которых
выводятся компоненты.
