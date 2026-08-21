---
title: Монитор доступности
group: Сборки
source: src/kit.css
template: splash
---

Один объект под наблюдением и два вопроса к нему сразу: **когда** что-то
случилось и **сколько** этого было. Ни одного класса, кроме библиотечных, и ни
одной строки скрипта.

```html preview context
<div class="inst-shell">
  <aside class="inst-shell-aside">
    <a class="inst-shell-brand" href="#"><span class="inst-dot" data-tone="ok"></span>beacon</a>
    <nav class="inst-nav" aria-label="Разделы">
      <div class="inst-nav-group">
        <span class="inst-nav-label">Наблюдение</span>
        <a class="inst-nav-item" href="#" aria-current="page">Мониторы
          <span class="inst-badge inst-nav-count">12</span></a>
        <a class="inst-nav-item" href="#">Инциденты
          <span class="inst-badge inst-nav-count" data-tone="error">1</span></a>
        <a class="inst-nav-item" href="#">Сводка</a>
      </div>
      <div class="inst-nav-group">
        <span class="inst-nav-label">Узел</span>
        <a class="inst-nav-item" href="#">Соседи</a>
        <a class="inst-nav-item" href="#">Настройки</a>
      </div>
    </nav>
  </aside>

  <header class="inst-shell-header">
    <nav aria-label="Хлебные крошки">
      <ol class="inst-crumbs">
        <li><a href="#">Мониторы</a></li>
        <li><span aria-current="page">absalyamov.ru</span></li>
      </ol>
    </nav>
    <span class="inst-cluster-spacer"></span>
    <span class="inst-statusbar-item"><span class="inst-dot" data-tone="ok"></span>проверено 12 с назад</span>
  </header>

  <main class="inst-shell-main">
    <div class="inst-container inst-container--full">
      <div class="inst-stack">

      <div class="inst-page-header">
        <div class="inst-page-header-main">
          <h1 class="inst-page-title">absalyamov.ru
            <span class="inst-badge" data-tone="ok"><span class="inst-dot"></span>работает</span></h1>
          <p class="inst-page-desc">https://absalyamov.ru/ · http · каждые 8 минут</p>
        </div>
        <div class="inst-page-actions">
          <button class="inst-btn inst-btn--sm" type="button">Проверить сейчас</button>
          <button class="inst-btn inst-btn--sm" type="button">Заглушить на час</button>
        </div>
      </div>

      <div class="inst-panel">
        <div class="inst-panel-header">
          <span class="inst-panel-title">Последние 24 часа</span>
          <span class="inst-cluster-spacer"></span>
          <span class="inst-badge" data-tone="warn">3 инцидента за неделю</span>
        </div>
        <div class="inst-panel-body inst-stack">

          <div class="inst-share" role="img" aria-label="Последние 24 часа: 95,8% работает, 3,1% не работает, 1,1% не проверяли">
            <span class="inst-share-part" data-tone="ok" style="--share:0.958"></span>
            <span class="inst-share-part" data-tone="error" style="--share:0.031"></span>
            <span class="inst-share-part" data-tone="neutral" style="--share:0.011"></span>
          </div>

          <div>
        <div class="inst-history inst-history--lg" role="img" aria-label="Последние 24 часа: 97 проверок, 4 с ошибкой, один час без проверок">
          <span class="inst-history-group" style="--n:4"><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-tooltip-text">00:00 — 4 проверки</span></span>
          <span class="inst-history-group" style="--n:4"><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-tooltip-text">01:00 — 4 проверки</span></span>
          <span class="inst-history-group" style="--n:4"><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-tooltip-text">02:00 — 4 проверки</span></span>
          <span class="inst-history-group" style="--n:4"><span class="inst-history-tick" data-tone="error"></span><span class="inst-history-tick" data-tone="error"></span><span class="inst-history-tick" data-tone="error"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-tooltip-text">03:00 — 4 проверки</span></span>
          <span class="inst-history-group" style="--n:4"><span class="inst-history-tick" data-tone="error"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-tooltip-text">04:00 — 4 проверки</span></span>
          <span class="inst-history-group" style="--n:4"><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-tooltip-text">05:00 — 4 проверки</span></span>
          <span class="inst-history-group" style="--n:1"><span class="inst-history-tick"></span><span class="inst-tooltip-text">06:00 — проверок не было</span></span>
          <span class="inst-history-group" style="--n:4"><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-tooltip-text">07:00 — 4 проверки</span></span>
          <span class="inst-history-group" style="--n:4"><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-tooltip-text">08:00 — 4 проверки</span></span>
          <span class="inst-history-group" style="--n:4"><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-tooltip-text">09:00 — 4 проверки</span></span>
          <span class="inst-history-group" style="--n:4"><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-tooltip-text">10:00 — 4 проверки</span></span>
          <span class="inst-history-group" style="--n:1"><span class="inst-history-tick" data-tone="ok"></span><span class="inst-tooltip-text">11:00 — 1 проверка</span></span>
          <span class="inst-history-group" style="--n:4"><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-tooltip-text">12:00 — 4 проверки</span></span>
          <span class="inst-history-group" style="--n:4"><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-tooltip-text">13:00 — 4 проверки</span></span>
          <span class="inst-history-group" style="--n:4"><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-tooltip-text">14:00 — 4 проверки</span></span>
          <span class="inst-history-group" style="--n:4"><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-tooltip-text">15:00 — 4 проверки</span></span>
          <span class="inst-history-group" style="--n:4"><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-tooltip-text">16:00 — 4 проверки</span></span>
          <span class="inst-history-group" style="--n:9"><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-tooltip-text">17:00 — 9 проверки</span></span>
          <span class="inst-history-group" style="--n:4"><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-tooltip-text">18:00 — 4 проверки</span></span>
          <span class="inst-history-group" style="--n:4"><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-tooltip-text">19:00 — 4 проверки</span></span>
          <span class="inst-history-group" style="--n:4"><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="warn"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-tooltip-text">20:00 — 4 проверки</span></span>
          <span class="inst-history-group" style="--n:4"><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-tooltip-text">21:00 — 4 проверки</span></span>
          <span class="inst-history-group" style="--n:4"><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-tooltip-text">22:00 — 4 проверки</span></span>
          <span class="inst-history-group" style="--n:4"><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-history-tick" data-tone="ok"></span><span class="inst-tooltip-text">23:00 — 4 проверки</span></span>
        </div>
        <div class="inst-history-axis" aria-hidden="true">
          <span class="inst-history-axis-cell" style="--n:4">00</span>
          <span class="inst-history-axis-cell" style="--n:4" data-minor>01</span>
          <span class="inst-history-axis-cell" style="--n:4" data-minor>02</span>
          <span class="inst-history-axis-cell" style="--n:4">03</span>
          <span class="inst-history-axis-cell" style="--n:4" data-minor>04</span>
          <span class="inst-history-axis-cell" style="--n:4" data-minor>05</span>
          <span class="inst-history-axis-cell" style="--n:1">06</span>
          <span class="inst-history-axis-cell" style="--n:4" data-minor>07</span>
          <span class="inst-history-axis-cell" style="--n:4" data-minor>08</span>
          <span class="inst-history-axis-cell" style="--n:4">09</span>
          <span class="inst-history-axis-cell" style="--n:4" data-minor>10</span>
          <span class="inst-history-axis-cell" style="--n:1" data-minor>11</span>
          <span class="inst-history-axis-cell" style="--n:4">12</span>
          <span class="inst-history-axis-cell" style="--n:4" data-minor>13</span>
          <span class="inst-history-axis-cell" style="--n:4" data-minor>14</span>
          <span class="inst-history-axis-cell" style="--n:4">15</span>
          <span class="inst-history-axis-cell" style="--n:4" data-minor>16</span>
          <span class="inst-history-axis-cell" style="--n:9" data-minor>17</span>
          <span class="inst-history-axis-cell" style="--n:4">18</span>
          <span class="inst-history-axis-cell" style="--n:4" data-minor>19</span>
          <span class="inst-history-axis-cell" style="--n:4" data-minor>20</span>
          <span class="inst-history-axis-cell" style="--n:4">21</span>
          <span class="inst-history-axis-cell" style="--n:4" data-minor>22</span>
          <span class="inst-history-axis-cell" style="--n:4" data-minor>23</span>
        </div>
          </div>

          <ul class="inst-legend">
            <li class="inst-legend-item" data-tone="ok"><span class="inst-legend-swatch"></span>работает <span class="inst-legend-value">95.8%</span></li>
            <li class="inst-legend-item" data-tone="error"><span class="inst-legend-swatch"></span>не работает <span class="inst-legend-value">3.1%</span></li>
            <li class="inst-legend-item"><span class="inst-legend-swatch" data-empty></span>не проверяли <span class="inst-legend-value">1.1%</span></li>
          </ul>

        </div>
      </div>

      <div class="inst-split">
        <div class="inst-split-main">
          <div class="inst-panel">
            <div class="inst-panel-header"><span class="inst-panel-title">Инциденты</span></div>
            <table class="inst-table">
              <thead>
                <tr><th>Начался</th><th>Длился</th><th>Причина</th></tr>
              </thead>
              <tbody>
                <tr>
                  <td>03:12</td>
                  <td>14 мин</td>
                  <td><span class="inst-badge" data-tone="error"><span class="inst-dot"></span>таймаут</span></td>
                </tr>
                <tr>
                  <td>04:01</td>
                  <td>4 мин</td>
                  <td><span class="inst-badge" data-tone="error"><span class="inst-dot"></span>502</span></td>
                </tr>
                <tr>
                  <td>20:35</td>
                  <td>3 мин</td>
                  <td><span class="inst-badge" data-tone="warn"><span class="inst-dot"></span>медленно</span></td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>

        <div class="inst-split-side">
          <div class="inst-panel">
            <div class="inst-panel-header"><span class="inst-panel-title">Доступность</span></div>
            <div class="inst-panel-body inst-stack inst-stack--tight">
              <div>
                <div class="inst-meter-row"><span>24 часа</span><span class="inst-meter-value">95.8%</span></div>
                <div class="inst-meter" data-tone="warn" role="progressbar" aria-label="24 часа"
                     aria-valuenow="96" aria-valuemin="0" aria-valuemax="100">
                  <div class="inst-meter-fill" style="inline-size:95.8%"></div>
                </div>
              </div>
              <div>
                <div class="inst-meter-row"><span>7 дней</span><span class="inst-meter-value">99.2%</span></div>
                <div class="inst-meter" role="progressbar" aria-label="7 дней"
                     aria-valuenow="99" aria-valuemin="0" aria-valuemax="100">
                  <div class="inst-meter-fill" style="inline-size:99.2%"></div>
                </div>
              </div>
              <dl class="inst-kv">
                <dt>Проверено</dt><dd>12 с назад</dd>
                <dt>Задержка</dt><dd>38 ms</dd>
                <dt>Сертификат</dt><dd>осталось 54 дня</dd>
              </dl>
            </div>
          </div>
        </div>
      </div>

      </div>
    </div>
  </main>
</div>
```

## Из чего собрано

| Часть | Чем собрана |
|---|---|
| Каркас | [Оболочка приложения](../layout/shell.md) и [раздвижка](../layout/split.md) |
| Имя и действия | [Шапка экрана](../layout/page-header.md), [бейдж](../components/display/badge.md), [кнопка](../components/actions/button.md) |
| Сколько суток работало | [Полоса состава](../components/charts/share.md), по части на исход |
| Когда это было | [История проверок](../agent/history.md) с группами, линейкой и подсказкой на час |
| Расшифровка цветов | [Легенда](../components/charts/legend.md), состояния, а не ряды |
| Числа по окнам | [Мера](../components/charts/meter.md) и [список свойств](../components/display/kv.md) |
| Инциденты | [Таблица](../components/display/table.md) с бейджем в причине |

## Что здесь видно

**Две полосы отвечают на разное.** История упорядочена, состав подведён: те же
3,1% — это один долгий провал или девяносто разбросанных, и различить их умеет
только история. Одно другого не заменяет, поэтому стоят обе и с общей легендой.

**Час несёт свой вес.** Ширина группы есть её доля проверок: в 11:00 была одна,
и час узкий, в 17:00 их девять, и час широкий. Без этого монитор с одной
проверкой в час выглядит ровно как монитор с сотней.

**Линейка — не вторая раскладка.** Ячейка оси берёт тот же `--n`, что и группа
над ней, поэтому разъехаться этим двум строкам нечем.

**Легенда расшифровывает состояния, а не ряды.** Образец падает в `--tone-mark`,
а час, в который никто не смотрел, берёт `data-empty` — нарисованный дорожкой,
тем же значением, каким его рисует полоса.

## Связанное

[История проверок](../agent/history.md) [Полоса состава](../components/charts/share.md)
[Легенда](../components/charts/legend.md) [Мера](../components/charts/meter.md)
[Оболочка приложения](../layout/shell.md) [Консоль прогонов](./console.md)
