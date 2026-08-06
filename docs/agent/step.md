---
title: Шаг — вызов инструмента
group: Агентный слой
status: stable
source: src/agent.css
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

## Когда использовать

| Используйте | Возьмите другое |
|---|---|
| Вызов инструмента, у которого есть вход и выход | **Задача в очереди без вывода** — [строка очереди](./task.md): она плоская, и их на экране сотни |
| Вывод длинный и его надо сворачивать | **Однородный поток строк** — [лог](./log.md) с колонками уровня и времени |
| Шагов десятки и состав заранее неизвестен | **Известные наперёд шаги мастера** — [шаги](../components/navigation/steps.md): там конечная последовательность |
| — | **Просто сворачиваемый блок текста** — [аккордеон](../components/feedback/accordion.md) |

## Почему `<details>`, а не свой раскрыватель

Своя реализация обошлась бы в `aria-expanded`, обработчик клавиатуры,
управление фокусом и синхронизацию состояния — и всё равно проиграла бы
нативной по мелочам вроде поиска по странице, который раскрывает свёрнутое.

```html
<details class="inst-step" open>…</details>
```

Открытость — атрибут `open`, а не класс и не `data-state`. Это второй канал,
которого здесь нет намеренно.

## Состояния

| `data-state` | Что происходит |
|---|---|
| нет атрибута | Шаг ещё не начался |
| `running` | Идёт. Точка пульсирует, имя берёт акцент |
| `ok` | Вернулся успешно |
| `failed` | Вызов не удался |

Обратите внимание: у шага словарь **свой** — `ok`, а не `done`, и в нём нет
`skipped`. Словари `data-state` перечислены у каждого компонента отдельно,
потому что у строки очереди и вызова инструмента разные фазы жизни.

## Свёрнутый вывод

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

## Собранный экран

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

## Справочник

### Классы

| Класс | Работа |
|---|---|
| `inst-step` | Шаг, на `<details>` |
| `inst-step-head` | Шапка, на `<summary>` |
| `inst-step-twist` | Треугольник раскрытия. Поворачивается по `[open]` |
| `inst-step-name` | Имя инструмента, моноширинное |
| `inst-step-sub` | Аргумент вызова |
| `inst-step-meta` | Объём и время, прижаты к концу |
| `inst-step-body` | Тело: аргументы и вывод |
| `inst-output` | Обёртка сворачиваемого вывода |
| `inst-output-body` | Сам вывод |
| `inst-output-more` | Кнопка разворота с числом |

### Атрибуты

| Атрибут | Значения | Где |
|---|---|---|
| `data-state` | `running` · `ok` · `failed` | на `inst-step` |
| `open` | нативный | на `<details>` |
| `data-truncated` | `true` · `false` | на `inst-output` |

### Токены

`--space-1` · `--space-3` · `--gap-inline` · `--pad-cell-x` · `--text-xs` ·
`--font-mono` · `--size-dot` · `--text-muted` · `--tone-ink` · `--dur-1`

## Доступность

| | |
|---|---|
| Раскрытие | Нативное. `<summary>` фокусируется, `Enter` и `Space` переключают, `aria-expanded` браузер объявляет сам |
| Имя инструмента | Моноширинное и машинное (`read_file`). Скринридер прочитает его по частям — если это мешает, добавьте человеческую подпись рядом |
| Свёрнутый вывод | Число в кнопке — часть доступного имени: «Показать все 240 строк» озвучивается целиком |
| Состояние словом | Тон имени не озвучивается. Успех и отказ должны быть в тексте — в `inst-step-meta` или бейджем |
| Уменьшенное движение | Пульсация точки `running` замедляется, а не гаснет |
| Треугольник | Декоративен, отдельного имени не имеет: состояние раскрытия несёт сам `<details>` |

## Связанное

[Строка очереди](./task.md) · [Диф](./diff.md) ·
[Блок отказа](./failure.md) · [Лог](./log.md) ·
[Код](../components/display/code.md) ·
[Аккордеон](../components/feedback/accordion.md)

Исходник: `src/agent.css` · Почему кит устроен именно так — [конституция](../about/design-principles.md)
