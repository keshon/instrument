---
title: Экран настроек
group: Сборки
source: src/kit.css
title-en: "Settings screen"
group-en: "Blocks"
---

Форма, разбитая на смысловые блоки, с липкой полосой действий. Экран, где
важнее всего не показать, а не соврать: что применится сразу, а что по кнопке.

```html preview context
<div class="inst-container inst-container--narrow inst-stack">

  <div class="inst-page-header">
    <div class="inst-page-header-main">
      <h3 class="inst-page-title">Параметры прогона</h3>
      <p class="inst-page-desc">Применяются к новым прогонам проекта worldgen.</p>
    </div>
  </div>

  <div class="inst-banner" data-tone="warn">
    <div class="inst-banner-body">
      <div class="inst-banner-title">Есть несохранённые изменения</div>
      <div class="inst-banner-text">Они не применятся, пока вы не нажмёте «Сохранить».</div>
    </div>
  </div>

  <div class="inst-panel">
    <div class="inst-panel-header"><span class="inst-panel-title">Настройки</span></div>
    <div class="inst-panel-body">
      <form class="inst-form">

        <fieldset class="inst-fieldset">
          <legend>Режим прогона</legend>
          <div class="inst-grid">
            <label class="inst-choice-card">
              <input type="radio" name="set-mode" checked>
              <span><span class="inst-choice-card-title">Параллельно</span>
                <span class="inst-choice-card-desc">До 12 агентов сразу. Быстрее, но расход токенов выше.</span></span>
            </label>
            <label class="inst-choice-card">
              <input type="radio" name="set-mode">
              <span><span class="inst-choice-card-title">Последовательно</span>
                <span class="inst-choice-card-desc">По одному агенту. Лог читается сверху вниз.</span></span>
            </label>
          </div>
        </fieldset>

        <div class="inst-field">
          <label class="inst-label" for="set-name">Имя прогона<span class="inst-required">*</span></label>
          <input class="inst-input" id="set-name" value="nightly" aria-describedby="set-name-hint" required>
          <span class="inst-field-hint" id="set-name-hint">Латиница, цифры и дефис</span>
        </div>

        <div class="inst-field inst-field--side">
          <label class="inst-label" for="set-model">Модель</label>
          <span class="inst-select-wrap">
            <select class="inst-select" id="set-model"><option>opus</option><option>sonnet</option></select>
          </span>
        </div>

        <div class="inst-field inst-field--side">
          <label class="inst-label" for="set-limit">Лимит токенов</label>
          <div class="inst-input-group">
            <input class="inst-input" id="set-limit" value="250000">
            <span class="inst-input-affix">на прогон</span>
          </div>
        </div>

        <fieldset class="inst-fieldset">
          <legend>Применяются сразу</legend>
          <div class="inst-stack inst-stack--tight">
            <label class="inst-switch"><input type="checkbox" checked>Показывать лог в реальном времени</label>
            <label class="inst-switch"><input type="checkbox">Останавливать прогон на первом отказе</label>
          </div>
          <div class="inst-note" data-tone="neutral">Свитч применяется в момент нажатия — кнопка «Сохранить» на него не влияет.</div>
        </fieldset>

      </form>

      <div class="inst-form-actions">
        <span class="inst-form-actions-note">Изменения не сохранены</span>
        <button class="inst-btn" type="button">Отмена</button>
        <button class="inst-btn inst-btn--primary" type="button">Сохранить</button>
      </div>
    </div>
  </div>
</div>
```

## Из чего собрано

| Часть | Чем собрана |
|---|---|
| Название экрана | [Шапка экрана](../layout/page-header.md) с `<h1>` и пояснением |
| Сообщение уровня экрана | [Баннер](../components/feedback/banner.md): относится ко всему, а не к соседнему полю |
| Выбор режима | [Выбираемая карточка](../components/inputs/choice-card.md) в [сетке](../layout/flow.md), обёрнута в `fieldset` |
| Поля | [Текстовое поле](../components/inputs/input.md), [селект](../components/inputs/select.md), приставка `inst-input-affix` |
| Немедленные настройки | [Свитч](../components/inputs/toggles.md) и [сноска](../components/feedback/note.md) рядом |
| Полоса действий | [Форма](../components/inputs/form.md): липкая к нижнему краю тела панели |
| Ширина | [Контейнер](../layout/container.md) `--narrow`, 44rem |

## Что здесь видно

**Экран честен насчёт момента применения.** Свитчи стоят в отдельной группе со
сноской: они срабатывают сразу. Всё остальное ждёт «Сохранить», и об этом
говорят баннер сверху и пояснение в полосе действий. Свитч, требующий нажать
кнопку, — самая частая ложь в настройках.

**Обязательность несёт знак и слово.** `required` на контроле плюс видимая
пометка: атрибут без пометки виден только скринридеру, пометка без атрибута —
только глазом.

**Подпись сбоку там, где поля однотипны.** Короткие однородные значения
превращаются в таблицу, которую сканируют не читая; разнородное поле с
подсказкой оставлено с подписью сверху.

**Группа полей без рамки.** `fieldset` даёт группе имя, а границу задают
заголовок и воздух: карточки внутри уже обведены, и рамка вокруг них дала бы
рамку внутри рамки.

## Связанное

[Дашборд](./dashboard.md) [Инспектор](./inspector.md)
[Форма](../components/inputs/form.md)
[Выбираемая карточка](../components/inputs/choice-card.md)
[Контейнер](../layout/container.md)
