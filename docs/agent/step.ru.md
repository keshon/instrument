---
title: Шаг — вызов инструмента
group: Агентный слой
layout: component
source: src/agent.css
api:
  - { name: "inst-step", kind: "класс", doc: "Шаг, на `<details>`" }
  - { name: "inst-step-head", kind: "класс", doc: "Шапка, на `<summary>`" }
  - { name: "inst-step-twist", kind: "класс", doc: "Треугольник раскрытия. Поворачивается по `[open]`" }
  - { name: "inst-step-name", kind: "класс", doc: "Имя инструмента, моноширинное" }
  - { name: "inst-step-sub", kind: "класс", doc: "Аргумент вызова" }
  - { name: "inst-step-meta", kind: "класс", doc: "Объём и время, прижаты к концу" }
  - { name: "inst-step-body", kind: "класс", doc: "Тело: аргументы и вывод" }
  - { name: "inst-output", kind: "класс", doc: "Обёртка сворачиваемого вывода" }
  - { name: "inst-output-body", kind: "класс", doc: "Сам вывод" }
  - { name: "inst-output-more", kind: "класс", doc: "Кнопка разворота с числом" }
  - { name: "data-state", kind: "атрибут", value: "running · ok · failed", doc: "на `inst-step`" }
  - { name: "open", kind: "атрибут", value: "нативный", doc: "на `<details>`" }
  - { name: "data-details-all", kind: "атрибут", value: "селектор области · пусто", doc: "Кнопка «раскрыть все». Пусто — ближайшая панель" }
  - { name: "data-label-expand", kind: "атрибут", doc: "Подпись кнопки в свёрнутом состоянии" }
  - { name: "data-label-collapse", kind: "атрибут", doc: "Подпись в раскрытом" }
  - { name: "data-truncated", kind: "атрибут", value: "true · false", doc: "на `inst-output`" }
  - { name: "--space-1", kind: "токен" }
  - { name: "--space-3", kind: "токен" }
  - { name: "--gap-inline", kind: "токен" }
  - { name: "--pad-cell-x", kind: "токен" }
  - { name: "--text-xs", kind: "токен" }
  - { name: "--font-mono", kind: "токен" }
  - { name: "--size-dot", kind: "токен" }
  - { name: "--text-muted", kind: "токен" }
  - { name: "--tone-ink", kind: "токен" }
  - { name: "--dur-1", kind: "токен" }
---

Один вызов инструмента агентом: что вызвали, с чем и что вернулось. У шага
есть **тело**, и раскрывается оно платформой — это `<details>`, поэтому
клавиатура и состояние достаются бесплатно.

```html preview
<details class="inst-step" data-state="ok">
  <summary class="inst-step-head">
    <span class="inst-step-twist"></span>
    <span class="inst-dot"></span>
    <span class="inst-step-name">read_file</span>
    <span class="inst-step-sub inst-u-truncate">terrain/heightmap.ts</span>
    <span class="inst-step-meta">240 строк · 0,4 с</span>
  </summary>
  <div class="inst-step-body">
    <div class="inst-code">{ "path": "terrain/heightmap.ts", "range": [1, 240] }</div>
  </div>
</details>
```

## Контракт

| Что | Обязательно | Почему |
|---|---|---|
| `<details>` и `<summary>` | да | Раскрытие, клавиатура, состояние `open` и поиск по странице — от платформы |
| `data-state` на `inst-step` | да | Фаза вызова. Точка внутри наследует тон отсюда |
| `inst-u-truncate` на аргументе | да | Путь длиннее строки иначе распирает шапку и выталкивает время за край |
| `data-truncated` на `inst-output` | да, если вывод обрезан | Читатель обязан знать, что показано не всё |

### Доступность

| | |
|---|---|
| Раскрытие | Нативное. `<summary>` фокусируется, `Enter` и `Space` переключают, `aria-expanded` браузер объявляет сам |
| Имя инструмента | Моноширинное и машинное (`read_file`). Скринридер прочитает его по частям — если это мешает, добавьте человеческую подпись рядом |
| Свёрнутый вывод | Число в кнопке — часть доступного имени: «Показать все 240 строк» озвучивается целиком |
| Состояние словом | Тон имени не озвучивается. Успех и отказ должны быть в тексте — в `inst-step-meta` или бейджем |
| Уменьшенное движение | Пульсация точки `running` замедляется, а не гаснет |
| Треугольник | Декоративен, отдельного имени не имеет: состояние раскрытия несёт сам `<details>` |

## Состояния

| `data-state` | Что происходит |
|---|---|
| нет атрибута | Шаг ещё не начался |
| `running` | Идёт. Точка пульсирует, имя берёт акцент |
| `ok` | Вернулся успешно |
| `failed` | Вызов не удался |

У шага **свой** словарь: `ok` вместо `done`, и `skipped` в нём нет. Словари
`data-state` перечислены у каждого компонента отдельно — у строки очереди и у
вызова инструмента разные фазы жизни.

### Свёрнутый вывод

Лог инструмента почти всегда длиннее экрана. Молча обрезать нельзя — это враньё
про объём.

```html
<div class="inst-output" data-truncated="true">
  <pre class="inst-code inst-output-body">…</pre>
  <button class="inst-output-more" type="button">Показать все 240 строк</button>
</div>
```

**Свёрнутый вывод называет своё число словами.** «Показать все 240 строк», а
не «…» и не «Показать ещё»: пользователь должен знать, сколько он не видит,
чтобы решить, стоит ли разворачивать.

## Поведение

Своя реализация обошлась бы в `aria-expanded`, обработчик клавиатуры,
управление фокусом и синхронизацию состояния — и всё равно проиграла бы
нативной по мелочам вроде поиска по странице, который раскрывает свёрнутое.

```html
<details class="inst-step" open>…</details>
```

Открытость — атрибут `open`, а не класс и не `data-state`: второго канала для
одного состояния здесь нет.

### Раскрыть все

Шесть шагов — это шесть нажатий, чтобы прочитать прогон целиком, и ещё шесть,
чтобы вернуть как было. Кнопка ставится в шапку области:

```html preview
<div class="inst-panel">
  <div class="inst-panel-header">
    <span class="inst-panel-title">Последняя проверка</span>
    <span class="inst-panel-actions">
      <button class="inst-btn inst-btn--sm" type="button" aria-expanded="false"
              data-details-all data-label-expand="Раскрыть все"
              data-label-collapse="Свернуть все">Раскрыть все</button>
    </span>
  </div>
  <div class="inst-panel-body inst-panel-body--list">
    <details class="inst-step" data-state="ok">
      <summary class="inst-step-head">
        <span class="inst-step-twist"></span>
        <span class="inst-step-name">TCP</span>
        <span class="inst-step-sub inst-u-truncate">сокет открыт</span>
        <span class="inst-step-meta">38 мс</span>
      </summary>
      <div class="inst-step-body"><div class="inst-code">93.184.216.34:443</div></div>
    </details>
    <details class="inst-step" data-state="ok">
      <summary class="inst-step-head">
        <span class="inst-step-twist"></span>
        <span class="inst-step-name">TLS</span>
        <span class="inst-step-sub inst-u-truncate">рукопожатие завершено</span>
        <span class="inst-step-meta">64 мс</span>
      </summary>
      <div class="inst-step-body"><div class="inst-code">TLS 1.3</div></div>
    </details>
  </div>
</div>
```

Состояние кнопка определяет по содержимому, а не помнит: пока хоть один шаг
закрыт, нажатие раскрывает всё. Кнопка со своей памятью рано или поздно
расходится с тем, что человек открыл руками.

| Что | Зачем |
|---|---|
| `data-details-all` пустой | Область — ближайшая `inst-panel`. Со значением — селектор нужного узла |
| `aria-expanded` | `instrument.js` переставляет его сам; в разметке нужно начальное значение |
| `data-label-expand` / `data-label-collapse` | Подписи. Без них кнопка молча меняет смысл, не меняя надписи |

Событие `inst:details-all` всплывает с кнопки, `detail` — `{ expand, count }`.

## Сценарии

Транскрипт прогона — то, ради чего шаг существует. Порознь шаги выглядят
списком, вместе — рассказом о том, что делала машина.

Что проверить руками:

- **раскройте шаг** — это `<details>`, и раскрытие ничего не стоило: ни
  строки скрипта, ни `aria-expanded` руками;
- **найдите текст поиском по странице** (`Ctrl+F`) — браузер раскроет
  свёрнутый шаг сам. Своя реализация это теряет;
- **свёрнутый вывод называет своё число словами.** «Показать все 240 строк»,
  а не «…»: пользователь должен знать, сколько он не видит, чтобы решить,
  разворачивать ли.

```html preview context
<div class="inst-panel">
      <div class="inst-panel-body inst-panel-body--list">
        <details class="inst-step" data-state="ok">
          <summary class="inst-step-head">
            <span class="inst-step-twist"></span><span class="inst-dot"></span>
            <span class="inst-step-name">read_file</span>
            <span class="inst-step-sub inst-u-truncate">terrain/heightmap.ts</span>
            <span class="inst-step-meta">0,3 с</span>
          </summary>
          <div class="inst-step-body">
            <div class="inst-code">{ "path": "terrain/heightmap.ts", "range": [1, 240] }</div>
            <div class="inst-output" data-truncated="true">
              <pre class="inst-code inst-output-body">export function decodeHeightmap(buf: ArrayBuffer): Float32Array {
  const view = new DataView(buf);
  const size = view.getUint32(0, true);
  const out = new Float32Array(size * size);
  for (let i = 0; i &lt; out.length; i++) {
    out[i] = view.getUint16(4 + i * 2, true) / 65535;
  }
  return out;
}

export function sampleBilinear(map: Float32Array, size: number, x: number, y: number) {
  const x0 = Math.floor(x), y0 = Math.floor(y);
  const x1 = Math.min(x0 + 1, size - 1), y1 = Math.min(y0 + 1, size - 1);
  const fx = x - x0, fy = y - y0;
  return  map[y0 * size + x0] * (1 - fx) * (1 - fy)
        + map[y0 * size + x1] * fx * (1 - fy)
        + map[y1 * size + x0] * (1 - fx) * fy
        + map[y1 * size + x1] * fx * fy;
}</pre>
              <button class="inst-output-more" type="button">Показать все 240 строк</button>
            </div>
          </div>
        </details>

        <details class="inst-step" data-state="running" open>
          <summary class="inst-step-head">
            <span class="inst-step-twist"></span><span class="inst-dot"></span>
            <span class="inst-step-name">place_resources</span>
            <span class="inst-step-sub inst-u-truncate">проход 3 из 3</span>
            <span class="inst-step-meta">14,0 с</span>
          </summary>
          <div class="inst-step-body">
            <div class="inst-stack inst-stack--tight">
              <div class="inst-skeleton inst-skeleton--title"></div>
              <div class="inst-skeleton inst-skeleton--line"></div>
              <div class="inst-skeleton inst-skeleton--short"></div>
            </div>
          </div>
        </details>

        <details class="inst-step" data-state="failed">
          <summary class="inst-step-head">
            <span class="inst-step-twist"></span><span class="inst-dot"></span>
            <span class="inst-step-name">validate_paths</span>
            <span class="inst-step-sub inst-u-truncate">3 узла недостижимы</span>
            <span class="inst-step-meta">1,2 с</span>
          </summary>
          <div class="inst-step-body">
            <div class="inst-failure" role="alert">
              <div class="inst-failure-head">Валидация путей не прошла</div>
              <div class="inst-failure-reason">PathError: unreachable nodes (18,204) (19,204) (20,205) — no navmesh link from region 7</div>
              <div class="inst-failure-tried">Попыток: 3, последняя в 19:38:16</div>
              <div class="inst-failure-actions">
                <button class="inst-btn inst-btn--sm" type="button">Повторить</button>
                <button class="inst-btn inst-btn--sm inst-btn--ghost" type="button">Пропустить шаг</button>
                <button class="inst-btn inst-btn--sm inst-btn--ghost" type="button">Открыть регион 7</button>
              </div>
            </div>
          </div>
        </details>
      </div>
    </div>
```

## API

```api
```

## Связанное

```related
```
