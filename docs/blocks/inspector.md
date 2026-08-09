---
title: Инспектор
group: Сборки
source: src/kit.css
title-en: "Inspector"
group-en: "Blocks"
---

Список объектов слева, свойства выбранного справа. Экран, ради которого кит и
затевался: плотность здесь не украшение, а вместимость.

```html preview context
<div class="inst-split" style="--split-side: 20rem">

  <div class="inst-split-main inst-panel">
    <div class="inst-panel-header">
      <span class="inst-panel-title">Объекты сцены</span>
      <span class="inst-badge">4</span>
      <span class="inst-panel-actions">
        <button class="inst-btn inst-btn--sm inst-btn--ghost inst-btn--icon" type="button" aria-label="Обновить">
          <svg class="inst-icon" aria-hidden="true"><use href="#i-refresh"/></svg></button>
      </span>
    </div>
    <div class="inst-toolbar">
      <div class="inst-btn-group">
        <button class="inst-btn inst-btn--sm inst-btn--icon" type="button" aria-label="Переместить">
          <svg class="inst-icon" aria-hidden="true"><use href="#i-move"/></svg></button>
        <button class="inst-btn inst-btn--sm inst-btn--icon" type="button" aria-label="Повернуть">
          <svg class="inst-icon" aria-hidden="true"><use href="#i-rotate"/></svg></button>
        <button class="inst-btn inst-btn--sm inst-btn--icon" type="button" aria-label="Масштаб">
          <svg class="inst-icon" aria-hidden="true"><use href="#i-scale"/></svg></button>
      </div>
      <span class="inst-toolbar-spacer"></span>
      <span class="inst-badge" data-tone="ok"><span class="inst-dot"></span>сохранено</span>
    </div>
    <div class="inst-panel-body inst-panel-body--list">
      <div class="inst-tree" role="tree" aria-label="Объекты сцены">
        <div class="inst-tree-item" role="treeitem" aria-level="1" aria-expanded="true"
             tabindex="-1" style="--depth:0"><span class="inst-tree-twist"></span>terrain</div>
        <div class="inst-tree-item" role="treeitem" aria-level="2" aria-selected="true"
             tabindex="0" style="--depth:1">heightmap.ts</div>
        <div class="inst-tree-item" role="treeitem" aria-level="2"
             tabindex="-1" style="--depth:1">biomes.ts</div>
        <div class="inst-tree-item" role="treeitem" aria-level="1" aria-expanded="false"
             tabindex="-1" style="--depth:0"><span class="inst-tree-twist"></span>audio</div>
      </div>
    </div>
  </div>

  <div class="inst-split-side inst-panel">
    <div class="inst-panel-header"><span class="inst-panel-title">Свойства</span></div>
    <div class="inst-panel-body">
      <div class="inst-props">
        <div class="inst-prop">
          <span class="inst-prop-label" title="Позиция">Позиция</span>
          <span class="inst-prop-control inst-vec">
            <span class="inst-num-field"><span class="inst-num-axis">X</span><input type="number" value="128" aria-label="Позиция X"></span>
            <span class="inst-num-field"><span class="inst-num-axis">Y</span><input type="number" value="0" aria-label="Позиция Y"></span>
            <span class="inst-num-field"><span class="inst-num-axis">Z</span><input type="number" value="-64" aria-label="Позиция Z"></span>
          </span>
        </div>
        <div class="inst-prop">
          <span class="inst-prop-label" title="Плотность растительности">Плотность</span>
          <span class="inst-prop-control">
            <input class="inst-slider" type="range" min="0" max="100" value="72" aria-label="Плотность растительности">
          </span>
        </div>
        <div class="inst-prop">
          <label class="inst-prop-label" title="Материал" for="ins-mat">Материал</label>
          <span class="inst-prop-control">
            <span class="inst-select-wrap">
              <select class="inst-select inst-select--sm" id="ins-mat">
                <option>rock-01</option><option>sand-02</option><option>grass-03</option>
              </select>
            </span>
          </span>
        </div>
        <div class="inst-prop">
          <span class="inst-prop-label" title="Отбрасывать тень">Тень</span>
          <span class="inst-prop-control">
            <label class="inst-switch"><input type="checkbox" checked>Отбрасывать</label>
          </span>
        </div>
      </div>

      <dl class="inst-kv">
        <dt>Полигонов</dt><dd>48 210</dd>
        <dt>Обновлён</dt><dd>14:32:07</dd>
      </dl>
    </div>
  </div>
</div>
```

## Из чего собрано

| Часть | Чем собрана |
|---|---|
| Две колонки | [Сплит](../layout/split.md). Порог боковой части задан переменной `--split-side` |
| Список объектов | [Дерево](../agent/tree.md): глубина приходит из данных, а не из вложенных обёрток |
| Полоса инструментов | [Тулбар](../components/navigation/toolbar.md) с [группой кнопок](../components/actions/button-group.md) |
| Строки свойств | [Список свойств](../components/display/kv.md) — `inst-props` и `inst-prop` |
| Вектор | [Числовое поле инспектора](../components/inputs/num-field.md) в `inst-vec` |
| Остальные контролы | [Слайдер](../components/inputs/slider.md), [селект](../components/inputs/select.md), [свитч](../components/inputs/toggles.md) |
| Итоги только для чтения | [Список свойств](../components/display/kv.md) на `<dl>` |

## Что здесь видно

**Колонка имён одна на два примитива.** `inst-props` и `inst-kv` берут один и
тот же `--label-col`, поэтому редактируемые и нередактируемые свойства стоят по
одной вертикали.

**Контрол в строке свойства — `sm`.** Высота 26px против 32px у формы: на
экране, где свойств десятки, это разница в целый экран прокрутки. Цена названа
честно на [странице числового поля](../components/inputs/num-field.md): в
плотности `compact` цель нажатия уходит ниже нормы WCAG.

**Инлайновые стили здесь — данные.** `--depth` у узла дерева и `--split-side` у
раскладки: глубина приходит из модели, а порог — решение экрана. Оформление
целиком в ките.

**Панель — контейнер запросов.** Сузьте окно: строка свойства перестанет быть
двухколоночной, и подпись встанет над контролом. Это `@container`, а не
медиазапрос: реагирует панель, а не окно.

## Связанное

[Дашборд](./dashboard.md) [Экран настроек](./settings-screen.md)
[Сплит](../layout/split.md) [Числовое поле](../components/inputs/num-field.md)
[Дерево](../agent/tree.md)
