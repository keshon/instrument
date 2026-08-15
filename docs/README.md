# Документация

Исходники страниц сайта. Одна страница на компонент. Сборка —
`go -C site run ./cmd/site`; остальное про разработку — в
[CONTRIBUTING.md](../CONTRIBUTING.md).

## Структура

```
docs/
  start/         install
  foundations/   colors · density · elevation · icons · motion ·
                 spacing · tokens · typography · utilities · behavior
  layout/        container · flow · page-header · rail · section ·
                 shell · split · statusbar
  components/
    actions/     button · button-group · segmented
    charts/      legend · meter · palette · ring · sparkline
    display/     avatar · badge · calendar · card · code · kv ·
                 metric · panel · table · tag · timeline
    feedback/    accordion · banner · empty · note · skeleton ·
                 spinner · states · toast
    inputs/      choice-card · file · form · input · inserts ·
                 num-field · search · select · slider · toggles
    navigation/  breadcrumbs · nav · pagination · steps · tabs · toolbar
    overlays/    dialog · menu · popover · sheet · tooltip
  agent/         approval · budget · diff · failure · history · lane ·
                 log · run · step · task · tree
  blocks/        console
  about/         design-principles · audit
  internal/      процессные документы, не документация
```

Категории названы работой компонента: `.inst-metric` и `.inst-sparkline` лежат
в разных `.css`, но стоят на соседних страницах.

Перевод — соседний файл: `button.md` → `button.en.md`. Страница без перевода
показывает русский текст с пометкой.

## Шаблон страницы

Порядок разделов задан словарём в `site/internal/content/sections.go` и
проверяется сборкой. Полное описание — в
[internal/PAGE-CONTRACT.md](internal/PAGE-CONTRACT.md).

```
frontmatter        title · group · status · source
одна строка        какую работу делает компонент
ЖИВОЙ ПРИМЕР       первым, всегда
Установка          одна строка (опускается на страницах основ)
Когда использовать таблица в две колонки: используйте · возьмите другое
Варианты · Размеры · Состояния · С иконкой · Композиции
Справочник         КЛАССЫ · АТРИБУТЫ · CSS-ПЕРЕМЕННЫЕ · ТОКЕНЫ — таблицами
Доступность        клавиатура · контраст · движение · цели нажатия
Связанное          + исходник + ссылка на принципы дизайна
```

Эталоны: [кнопка](components/actions/button.md) — сложный API,
[бейдж](components/display/badge.md) — простой.

## Правила

**Справочник — таблицами, не прозой.** Таблицу можно просмотреть по диагонали
и проверить на полноту.

**Раздела «почему так» на странице нет.** Правила системы описаны в
[принципах дизайна](about/design-principles.md), и страница ссылается на них
один раз, строкой в подвале. Обоснование, которое меняет то, что читатель
напечатает, идёт не в прозу, а в таблицу CSS-переменных или токенов.

**Разметка берётся из работающей, а не сочиняется.** Каждый пример стоит на
странице живьём или проверен в браузере.

**Обязательные роли входят в пример, а не в примечание.** `role="listbox"` в
примере строки очереди обязателен, иначе `aria-selected` невалиден.

**«Когда не использовать» — со ссылкой на замену.** Строка «не используйте для
X» без ответа «возьмите Y» бесполезна.

**Три блока `preview`, а не десять.** Раздел, где каждый вариант получил свой
живой пример, читается дольше, чем сетка из одного.

## Разметка примеров

| Ограда | Что делает |
|---|---|
| ` ```html preview ` | Рендерится живьём, под ним — код с кнопкой копирования |
| ` ```html ` | Только код: композиции, которые незачем показывать целиком |
| ` ```css ` | Переопределение на стороне приложения |

## Проверка

```bash
go -C tools run ./cmd/docscheck
```

Сверяет страницы с библиотекой в обе стороны: выдуманный класс — ошибка,
класс без страницы — метрика. Заодно проверяются значения `data-*` против
словарей из CSS и имена токенов. Подробности — в
[CONTRIBUTING.md](../CONTRIBUTING.md#что-сверяет-docscheck).

Проверка не отличает разметку в примере от разметки, названной в прозе:
опечатка в тексте страницы тоже будет ошибкой.
