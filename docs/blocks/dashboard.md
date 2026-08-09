---
title: Дашборд
group: Сборки
source: src/kit.css
title-en: "Dashboard"
group-en: "Blocks"
---

Сводка за сутки: ряд метрик, тренд, таблица прогонов. Ни одного класса, кроме
китовых, и ни одной строки скрипта.

```html preview context
<div class="inst-container inst-container--full inst-stack">

  <div class="inst-page-header">
    <div class="inst-page-header-main">
      <h3 class="inst-page-title">Прогоны</h3>
      <p class="inst-page-desc">Сводка за последние сутки по всем проектам.</p>
    </div>
    <div class="inst-page-actions">
      <div class="inst-segmented" role="radiogroup" aria-label="Период">
        <button type="button" role="radio" aria-checked="false">Сутки</button>
        <button type="button" role="radio" aria-checked="true">Неделя</button>
        <button type="button" role="radio" aria-checked="false">Месяц</button>
      </div>
      <button class="inst-btn inst-btn--sm inst-btn--primary" type="button">Запустить</button>
    </div>
  </div>

  <div class="inst-metric-row">
    <div class="inst-metric">
      <div class="inst-metric-label">Прогонов</div>
      <div class="inst-metric-value">128</div>
      <div class="inst-metric-delta" data-dir="up" data-tone="ok">12% к прошлой неделе</div>
    </div>
    <div class="inst-metric">
      <div class="inst-metric-label">Среднее время</div>
      <div class="inst-metric-value">42<span class="inst-metric-unit">&nbsp;с</span></div>
      <div class="inst-metric-delta" data-dir="down" data-tone="ok">18% к прошлой</div>
    </div>
    <div class="inst-metric">
      <div class="inst-metric-label">Токенов</div>
      <div class="inst-metric-value">184<span class="inst-metric-unit">&nbsp;тыс.</span></div>
      <div class="inst-metric-delta">без изменений</div>
    </div>
    <div class="inst-metric">
      <div class="inst-metric-label">Упало</div>
      <div class="inst-metric-value">7</div>
      <div class="inst-metric-delta" data-dir="up" data-tone="error">5 новых</div>
    </div>
  </div>

  <div class="inst-split">
    <div class="inst-split-main inst-panel">
      <div class="inst-panel-header">
        <span class="inst-panel-title">Последние прогоны</span>
        <span class="inst-panel-actions">
          <span class="inst-search">
            <input class="inst-input inst-input--sm" type="search" placeholder="Фильтр" aria-label="Фильтр по прогонам">
          </span>
        </span>
      </div>
      <div class="inst-panel-body inst-panel-body--flush">
        <table class="inst-table">
          <caption class="inst-u-visually-hidden">Последние прогоны</caption>
          <thead>
            <tr>
              <th aria-sort="descending"><button class="inst-sort" type="button">Прогон</button></th>
              <th>Состояние</th>
              <th class="inst-num"><button class="inst-sort" type="button">Время, с</button></th>
              <th class="inst-num">Токенов</th>
            </tr>
          </thead>
          <tbody>
            <tr>
              <td>worldgen-01</td>
              <td><span class="inst-badge" data-tone="running"><span class="inst-dot"></span>идёт</span></td>
              <td class="inst-num">18,2</td><td class="inst-num">61 204</td>
            </tr>
            <tr>
              <td>biomes-04</td>
              <td><span class="inst-badge" data-tone="ok"><span class="inst-dot"></span>готово</span></td>
              <td class="inst-num">11,5</td><td class="inst-num">38 910</td>
            </tr>
            <tr>
              <td>audio-pass</td>
              <td><span class="inst-badge" data-tone="error"><span class="inst-dot"></span>упало</span></td>
              <td class="inst-num">4,2</td><td class="inst-num">1 180</td>
            </tr>
            <tr>
              <td>nightly-4127</td>
              <td><span class="inst-badge" data-tone="warn"><span class="inst-dot"></span>с замечаниями</span></td>
              <td class="inst-num">33,9</td><td class="inst-num">82 640</td>
            </tr>
          </tbody>
          <tfoot>
            <tr><td>Итого</td><td></td><td class="inst-num">67,8</td><td class="inst-num">183 934</td></tr>
          </tfoot>
        </table>
      </div>
    </div>

    <div class="inst-split-side inst-stack inst-stack--tight">
      <div class="inst-panel">
        <div class="inst-panel-header"><span class="inst-panel-title">Токены в сутки</span></div>
        <div class="inst-panel-body inst-stack inst-stack--tight">
          <svg class="inst-sparkline" viewBox="0 0 120 32" preserveAspectRatio="none" aria-hidden="true">
            <path class="inst-sparkline-area" d="M0,26 12,22 24,24 36,16 48,18 60,10 72,13 84,7 96,9 108,4 120,6 120,32 0,32Z"/>
            <path class="inst-sparkline-line" d="M0,26 12,22 24,24 36,16 48,18 60,10 72,13 84,7 96,9 108,4 120,6"/>
          </svg>
          <ul class="inst-legend">
            <li class="inst-legend-item"><span class="inst-legend-swatch"></span>Расход <span class="inst-legend-value">184 тыс.</span></li>
          </ul>
        </div>
      </div>

      <div class="inst-budget">
        <div class="inst-budget-row">
          <span>Бюджет</span>
          <span class="inst-budget-value">184&nbsp;320 <span class="inst-budget-of">из 250&nbsp;000</span></span>
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
```

## Из чего собрано

Ни одного нового класса. Экран целиком — существующие компоненты, поставленные
рядом.

| Часть | Чем собрана |
|---|---|
| Название и действия | [Шапка экрана](../layout/page-header.md), [сегментированный контрол](../components/actions/segmented.md), [кнопка](../components/actions/button.md) |
| Ряд чисел | [Метрика](../components/display/metric.md) без рамок: четыре числа — одна группа |
| Две колонки | [Сплит](../layout/split.md). Переносится сам, без медиазапроса |
| Таблица | [Таблица](../components/display/table.md) в [панели](../components/display/panel.md) с телом `--flush` |
| Состояние строки | [Бейдж](../components/display/badge.md) с точкой: тон **и** слово |
| Тренд | [Спарклайн](../components/charts/sparkline.md) и [легенда](../components/charts/legend.md) |
| Остаток | [Бюджет прогона](../agent/budget.md) с [мерой](../components/charts/meter.md) |
| Ритм | [Стопка и сетка](../layout/flow.md), ширина — [контейнер](../layout/container.md) |

## Что здесь видно

**Плотность работает целиком.** Переключите её в «Виде»: высоты контролов,
отступы панелей и шаг таблицы меняются разом, потому что все они приходят из
ролей, а не из чисел в разметке.

**Тон — один словарь.** `data-tone` на бейдже строки, на дельте метрики и на
заливке меры означает одно и то же и красится из одного места.

**Числа выровнены.** Табличные цифры включены во всём ките, поэтому разряды в
колонке стоят друг под другом без единого правила в разметке.

**Ни одной тени.** Глубина передана поверхностями: `--surface-page` под
экраном, `--surface-raised` у панелей, врез у поля поиска.

## Связанное

[Инспектор](./inspector.md) [Экран настроек](./settings-screen.md)
[Оболочка приложения](../layout/shell.md) [Прогон](../agent/run.md)
