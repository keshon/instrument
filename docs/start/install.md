---
title: Установка
group: Начало
source: dist/instrument.min.css
title-en: "Install"
group-en: "Start"
---

Кит подключается двумя тегами и не требует сборки. Ниже — целая страница,
которую можно скопировать и открыть.

## Один файл

```html
<link rel="stylesheet" href="instrument.min.css">
```

Возьмите `dist/instrument.min.css` — 82 КБ, 14 КБ в gzip. Это весь кит: токены,
шесть тем, три плотности, все компоненты.

`dist/instrument.css` — то же самое без сжатия, если хочется читать исходник в
инспекторе.

:::note
Пакета в npm нет. Скопируйте файл к себе или подключите репозиторий
сабмодулем — версия лежит в `VERSION` и в имени тега.
:::

## Иконки

Иконка рисуется ссылкой на символ: `<use href="#i-refresh">`. Значит, символы
должны быть в документе. Вставьте содержимое `assets/sprite.svg` первым
элементом внутри `<body>` — один раз на страницу.

```html
<body>
  <svg class="inst-sprite" aria-hidden="true" focusable="false" xmlns="http://www.w3.org/2000/svg">
    <!-- сюда содержимое assets/sprite.svg -->
  </svg>
```

Без спрайта всё остальное работает, просто иконки не рисуются. Полный список
из девяноста символов — на странице [иконок](../foundations/icons.md).

## Поведение с клавиатуры

```html
<script type="module" src="instrument.js"></script>
```

Отдельный файл и по желанию. Он выполняет то, что кит обещает разметкой:
стрелки в списках и вкладках, копирование, снятие тега, перетаскивание оси.
Ничего не рисует. Подробности — в [поведении](../foundations/behavior.md).

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
<section data-density="compact">      <!-- compact · comfortable -->
```

| Тема | |
|---|---|
| `light-neutral` `light-warm` `light-cool` | светлые, различаются уклоном |
| `dark-light` `dark-soft` `dark` | тёмные, различаются глубиной поверхностей |

## Ваши стили выигрывают

Кит целиком лежит в каскадных слоях, а ваши стили — вне их. Поэтому любое ваше
правило перебивает кит, и `!important` не нужен нигде.

```css
.inst-btn { border-radius: 0; }   /* сработает, специфичность неважна */
```

## Требования

Chrome 123+, Safari 17.5+, Firefox 120+ — по `light-dark()`, `oklch()`,
каскадным слоям, вложенности, `:has()` и единице `lh`. Полифилов нет.

## Дальше

[Консоль прогонов](../blocks/console.md) — собранный экран, из которого видно,
как компоненты стыкуются. Оттуда быстрее всего копировать.

[Оболочка приложения](../layout/shell.md) — шапка, боковая панель, основная
область.

[Конституция](../about/design-principles.md) — правила, по которым кит устроен.
Читать не обязательно, но если правило и компонент разойдутся, прав будет
правило.
