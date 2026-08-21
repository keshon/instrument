---
title: Консоль прогонов
group: Сборки
source: src/kit.css
template: splash
---

Основной сценарий библиотеки: слева список того, что машина делает,
справа — подробности выбранного. Ни одного класса, кроме библиотечных, и ни одной
строки скрипта.

```html preview context
<div class="inst-shell">
  <aside class="inst-shell-aside">
    <a class="inst-shell-brand" href="#"><span class="inst-dot" data-tone="running"></span>worldgen</a>
    <nav class="inst-nav" aria-label="Разделы">
      <div class="inst-nav-group">
        <span class="inst-nav-label">Работа</span>
        <a class="inst-nav-item" href="#" aria-current="page">Прогоны
          <span class="inst-badge inst-nav-count">4</span></a>
        <a class="inst-nav-item" href="#">Очередь</a>
        <a class="inst-nav-item" href="#">Артефакты</a>
      </div>
      <div class="inst-nav-group">
        <span class="inst-nav-label">Проект</span>
        <a class="inst-nav-item" href="#">Агенты</a>
        <a class="inst-nav-item" href="#">Ключи</a>
      </div>
    </nav>
  </aside>

  <header class="inst-shell-header">
    <nav aria-label="Хлебные крошки">
      <ol class="inst-crumbs">
        <li><a href="#">worldgen</a></li>
        <li><span aria-current="page">Прогоны</span></li>
      </ol>
    </nav>
    <span class="inst-cluster-spacer"></span>
    <span class="inst-search">
      <input class="inst-input inst-input--sm" type="search" placeholder="Поиск" aria-label="Поиск по прогонам">
    </span>
    <button class="inst-btn inst-btn--sm inst-btn--primary" type="button">Запустить</button>
  </header>

  <main class="inst-shell-main" id="main">
    <div class="inst-container inst-container--full">
      <div class="inst-split">

        <div class="inst-split-side inst-panel">
          <div class="inst-toolbar">
            <div class="inst-segmented" role="radiogroup" aria-label="Фильтр по состоянию">
              <button type="button" role="radio" aria-checked="true"  tabindex="0">Все</button>
              <button type="button" role="radio" aria-checked="false" tabindex="-1">Идут</button>
              <button type="button" role="radio" aria-checked="false" tabindex="-1">Упало</button>
            </div>
          </div>
          <div class="inst-panel-body inst-panel-body--list">
            <div role="listbox" aria-label="Прогоны">
              <div class="inst-task" role="option" data-state="running" aria-selected="true" tabindex="0">
                <span class="inst-task-gutter"><span class="inst-dot"></span></span>
                <span class="inst-task-main">
                  <span class="inst-task-title inst-u-truncate">nightly-4127</span>
                  <span class="inst-task-sub">7 агентов · main</span>
                </span>
                <span class="inst-task-meta">21 с</span>
              </div>
              <div class="inst-task" role="option" data-state="failed" aria-selected="false" tabindex="-1">
                <span class="inst-task-gutter"><span class="inst-dot"></span></span>
                <span class="inst-task-main">
                  <span class="inst-task-title inst-u-truncate">audio-pass</span>
                  <span class="inst-task-sub">1 агент · fix/mixdown</span>
                </span>
                <span class="inst-task-meta">4,2 с</span>
              </div>
              <div class="inst-task" role="option" data-state="warn" aria-selected="false" tabindex="-1">
                <span class="inst-task-gutter"><span class="inst-dot"></span></span>
                <span class="inst-task-main">
                  <span class="inst-task-title inst-u-truncate">biomes-04</span>
                  <span class="inst-task-sub">3 агента · main</span>
                </span>
                <span class="inst-task-meta">11,5 с</span>
              </div>
              <div class="inst-task" role="option" data-state="done" aria-selected="false" tabindex="-1">
                <span class="inst-task-gutter"><span class="inst-dot"></span></span>
                <span class="inst-task-main">
                  <span class="inst-task-title inst-u-truncate">terrain-88</span>
                  <span class="inst-task-sub">4 агента · main</span>
                </span>
                <span class="inst-task-meta">18,2 с</span>
              </div>
            </div>
          </div>
        </div>

        <div class="inst-split-main inst-stack">

          <div class="inst-page-header">
            <div class="inst-page-header-main">
              <h3 class="inst-page-title">nightly-4127</h3>
              <p class="inst-page-desc">Враждебный аудит: поиск самообмана по классам A–F.</p>
            </div>
            <div class="inst-page-actions">
              <span class="inst-badge" data-tone="running"><span class="inst-dot"></span>идёт</span>
              <button class="inst-btn inst-btn--sm inst-btn--danger" type="button">Остановить</button>
            </div>
          </div>

          <div class="inst-props">
            <div class="inst-prop">
              <span class="inst-prop-label" title="Модель">Модель</span>
              <span class="inst-prop-control">opus</span>
            </div>
            <div class="inst-prop">
              <span class="inst-prop-label" title="Агентов">Агентов</span>
              <span class="inst-prop-control">7 из 12</span>
            </div>
            <div class="inst-prop">
              <span class="inst-prop-label" title="Токенов израсходовано">Токенов</span>
              <span class="inst-prop-control">186 000</span>
            </div>
          </div>

          <div class="inst-panel">
            <div class="inst-panel-header">
              <span class="inst-panel-title">Фазы</span>
              <span class="inst-panel-actions">
                <button class="inst-btn inst-btn--sm inst-btn--ghost" type="button" aria-expanded="false"
                        data-details-all data-label-expand="Раскрыть все"
                        data-label-collapse="Свернуть все">Раскрыть все</button>
              </span>
            </div>
            <div class="inst-panel-body inst-panel-body--list">
              <details class="inst-step" data-state="ok">
                <summary class="inst-step-head">
                  <span class="inst-step-twist"></span>
                  <span class="inst-dot"></span>
                  <span class="inst-step-name">разбор</span>
                  <span class="inst-step-sub inst-u-truncate">7 агентов, 4 класса</span>
                  <span class="inst-step-meta">38 200 · 18 с</span>
                </summary>
                <div class="inst-step-body">
                  <div class="inst-code">{ "classes": ["docs-drift", "shared-and-chunk"], "depth": 2 }</div>
                </div>
              </details>
              <details class="inst-step" data-state="running" open>
                <summary class="inst-step-head">
                  <span class="inst-step-twist"></span>
                  <span class="inst-dot"></span>
                  <span class="inst-step-name">опровержение</span>
                  <span class="inst-step-sub inst-u-truncate">3 из 7 проверено</span>
                  <span class="inst-step-meta">21 400 · идёт</span>
                </summary>
                <div class="inst-step-body">
                  <div class="inst-log" role="log" aria-label="Лог фазы">
                    <div class="inst-log-line"><span class="inst-log-time">14:32:07</span><span class="inst-log-level">info</span><span>probes-assert: гипотеза принята</span></div>
                    <div class="inst-log-line" data-tone="warn"><span class="inst-log-time">14:32:11</span><span class="inst-log-level">warn</span><span>silent-failure: нет воспроизведения, понижено</span></div>
                    <div class="inst-log-line"><span class="inst-log-time">14:32:16</span><span class="inst-log-level">info</span><span>coverage-hole: проверка запущена</span></div>
                  </div>
                </div>
              </details>
              <details class="inst-step">
                <summary class="inst-step-head">
                  <span class="inst-step-twist"></span>
                  <span class="inst-dot"></span>
                  <span class="inst-step-name">сводка</span>
                  <span class="inst-step-sub inst-u-truncate">ждёт опровержения</span>
                  <span class="inst-step-meta">—</span>
                </summary>
                <div class="inst-step-body">
                  <div class="inst-code">ожидание</div>
                </div>
              </details>
            </div>
          </div>

          <div class="inst-budget">
            <div class="inst-budget-row">
              <span>Бюджет</span>
              <span class="inst-budget-value">186&nbsp;000 <span class="inst-budget-of">из 250&nbsp;000</span></span>
            </div>
            <div class="inst-meter" role="progressbar" aria-label="Расход бюджета"
                 aria-valuenow="74" aria-valuemin="0" aria-valuemax="100">
              <div class="inst-meter-fill" style="inline-size:74%"></div>
            </div>
            <div class="inst-budget-note">При текущем темпе хватит ещё примерно на 3 прогона.</div>
          </div>

        </div>
      </div>
    </div>
  </main>
</div>
```

## Из чего собрано

Ни одного нового класса. Экран целиком — существующие компоненты, поставленные
рядом.

| Часть | Чем собрана |
|---|---|
| Каркас | [Оболочка приложения](../layout/shell.md): колонка, шапка, рабочая область |
| Две колонки | [Сплит](../layout/split.md). Переносится сам, без медиазапроса |
| Фильтр | [Тулбар](../components/navigation/toolbar.md) с [сегментированным контролом](../components/actions/segmented.md) |
| Список | [Строка задачи](../agent/task.md) в `role="listbox"`: клавиатура от [instrument.js](../foundations/behavior.md) |
| Название и действия | [Шапка экрана](../layout/page-header.md), [бейдж](../components/display/badge.md), [кнопка](../components/actions/button.md) |
| Факты о прогоне | [Список свойств](../components/display/kv.md) — `inst-props` |
| Ход работы | [Шаг](../agent/step.md) со свёрткой и [лог](../agent/log.md) внутри |
| Остаток | [Бюджет прогона](../agent/budget.md) с [мерой](../components/charts/meter.md) |

## Что здесь видно

**Экран отвечает на «что сейчас», а не на «сколько всего».** Ряда крупных чисел
наверху нет: сводка за сутки — это отдельная задача и отдельный экран,
а здесь человек смотрит на одну работу и хочет знать её состояние. Числа стоят
там, где к ним возникает вопрос: расход — рядом с бюджетом, длительность — в
строке прогона.

**Фильтр живёт в строке, а не в панели.** Три состояния — это сегментированный
контрол на одну строку, а не выпадающий список и не боковая форма с кнопкой
«Применить».

**Выбор — это состояние списка, а не переход.** `role="listbox"` плюс
`aria-selected`: стрелки водят по прогонам, деталь справа меняется. Строка не
ссылка, потому что никуда не ведёт.

**Тон — один словарь.** `data-tone` на точке бренда, бейдже прогона, строке лога
и штрихе фазы означает одно и то же и красится из одного места.

**Идущее пульсирует, а не подсвечено.** Носитель состояния «идёт» — пульсация
точки; тон здесь вторичен, и поэтому он же может означать интеракцию.

## Состояния

Настоящий экран половину времени показывает не данные. Эти три случая — часть
контракта, а не отдельная работа «на потом».

```html preview
<div class="inst-panel">
  <div class="inst-empty">
    <div class="inst-empty-title">Под фильтр «упало» ничего не подошло</div>
    <div class="inst-empty-desc">За последние сутки все прогоны завершились успешно.</div>
    <button class="inst-btn inst-btn--sm" type="button">Сбросить фильтр</button>
  </div>
</div>
```

Пока список едет, место занимает [скелетон](../components/feedback/skeleton.md)
— столько строк, сколько ожидается, и той же высоты:

```html preview
<div class="inst-panel">
  <div class="inst-panel-body inst-stack inst-stack--tight">
    <div class="inst-skeleton inst-skeleton--title"></div>
    <div class="inst-skeleton inst-skeleton--line"></div>
    <div class="inst-skeleton inst-skeleton--line inst-skeleton--short"></div>
  </div>
</div>
```

Отказ — это [блок отказа](../agent/failure.md), а не строка красного текста: у
него есть причина, перечень попыток и хотя бы один выход.

```html preview
<div class="inst-failure" role="alert">
  <div class="inst-failure-head">audio-pass не запустился</div>
  <div class="inst-failure-reason">EBUSY: не удалось прочитать <code>chunks.bin</code> — файл занят другим прогоном.</div>
  <ul class="inst-failure-tried">
    <li>Повтор через 1 с — тот же код</li>
    <li>Повтор через 4 с — тот же код</li>
  </ul>
  <div class="inst-failure-actions">
    <button class="inst-btn inst-btn--sm" type="button">Повторить</button>
    <button class="inst-btn inst-btn--sm inst-btn--ghost" type="button">Открыть лог</button>
  </div>
</div>
```

## Параметры — шторкой, а не экраном

Настройки прогона не заслуживают собственного экрана с полосой «Отмена /
Сохранить»: это тот же объект, только с другой стороны. [Шторка](../components/overlays/sheet.md)
открывается поверх, ничего не теряя из виду.

```html preview
<button class="inst-btn" type="button"
        onclick="document.getElementById('con-params').showModal()">Параметры прогона</button>

<dialog class="inst-dialog inst-sheet" id="con-params">
  <form method="dialog">
    <div class="inst-dialog-head">
      <span class="inst-dialog-title">Параметры</span>
      <button class="inst-btn inst-btn--sm inst-btn--ghost inst-btn--icon inst-dialog-close" type="submit" aria-label="Закрыть">
        <svg class="inst-icon" aria-hidden="true"><use href="#i-close"/></svg>
      </button>
    </div>
    <div class="inst-dialog-body">
      <form class="inst-form inst-form--side">
        <div class="inst-field">
          <label class="inst-label" for="con-model">Модель</label>
          <span class="inst-select-wrap">
            <select class="inst-select" id="con-model"><option>opus</option><option>sonnet</option></select>
          </span>
        </div>
        <div class="inst-field">
          <label class="inst-label" for="con-limit">Лимит токенов</label>
          <div class="inst-input-group">
            <input class="inst-input" id="con-limit" value="250000">
            <span class="inst-input-affix">на прогон</span>
          </div>
        </div>
        <div class="inst-field">
          <label class="inst-label" for="con-par">Агентов сразу</label>
          <input class="inst-slider" type="range" min="1" max="16" value="12" id="con-par" style="--fill:73%">
        </div>
      </form>
      <div class="inst-note" data-tone="neutral">Изменения применяются к следующему прогону и сохраняются сразу.</div>
    </div>
  </form>
</dialog>
```

**Полосы «Отмена / Сохранить» здесь нет.** Она сообщает, что до нажатия ничего
не произошло, тогда как свитч уже сработал, а поле уже отправило запрос.
Изменения применяются сразу, и об этом сказано строкой под формой.

## Связанное

[Оболочка приложения](../layout/shell.md) [Прогон](../agent/run.md)
[Строка задачи](../agent/task.md) [Шаг](../agent/step.md)
[Блок отказа](../agent/failure.md) [Шторка](../components/overlays/sheet.md)
