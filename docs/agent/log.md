---
title: Лог
group: Агентный слой
status: stable
source: src/components.css
needs-js: Копирование в буфер, виртуализация
api:
  - { name: "inst-log", kind: "класс", doc: "Контейнер, прокручивается" }
  - { name: "inst-log-line", kind: "класс", doc: "Строка: три колонки" }
  - { name: "inst-log-time", kind: "класс", doc: "Отметка времени" }
  - { name: "inst-log-level", kind: "класс", doc: "Уровень, ширина `5ch`" }
  - { name: "data-tone", kind: "атрибут", doc: "`warn` · `error`" }
  - { name: "--level-ink", kind: "переменная", value: "--text-muted" }
  - { name: "--font-mono", kind: "токен" }
  - { name: "--text-xs", kind: "токен" }
  - { name: "--surface-sunken", kind: "токен" }
  - { name: "--radius-md", kind: "токен" }
  - { name: "--space-4", kind: "токен" }
  - { name: "--pad-panel", kind: "токен" }
  - { name: "--gap-inline", kind: "токен" }
  - { name: "--warn-text", kind: "токен" }
  - { name: "--err-text", kind: "токен" }
title-en: "Log"
group-en: "Agent layer"
---

Поток строк от машины: время, уровень, сообщение. Моноширинный набор и
колонки фиксированной ширины, поэтому лог читается столбцами, а не сплошным
текстом.

```html preview
<div class="inst-log" role="log" aria-label="Лог прогона">
  <div class="inst-log-line"><span class="inst-log-time">14:32:07</span><span class="inst-log-level">info</span><span>Запуск worldgen-01</span></div>
  <div class="inst-log-line"><span class="inst-log-time">14:32:09</span><span class="inst-log-level">info</span><span>Прочитано 4 файла</span></div>
  <div class="inst-log-line" data-tone="warn"><span class="inst-log-time">14:32:11</span><span class="inst-log-level">warn</span><span>chunks.bin занят, повтор через 1 с</span></div>
  <div class="inst-log-line" data-tone="error"><span class="inst-log-time">14:32:16</span><span class="inst-log-level">error</span><span>EBUSY: не удалось прочитать chunks.bin</span></div>
</div>
```

## Ширины колонок в `ch`, а не по содержимому

Сетка объявлена **на строке**, то есть у каждой строки она своя. Если задать
колонку уровня по содержимому, её ширина будет меняться от слова к слову:
`info` уже, чем `error`, и сообщения поедут лесенкой. Поэтому
`inst-log-level` получает фиксированную ширину в `5ch`.

## Тон строки

```html
<div class="inst-log-line" data-tone="error">…</div>
```

| Тон | Что красится |
|---|---|
| нет атрибута | Только уровень, приглушённо |
| `warn` | Уровень |
| `error` | Уровень **и всё сообщение** |

Ошибка красится целиком намеренно: её надо найти в тысяче строк, пролистывая
глазом, а не читая.

## Когда использовать

| Используйте | Возьмите другое |
|---|---|
| Поток однородных строк от машины | **События с разным составом** — [лента](../components/display/timeline.md) |
| Строк сотни и тысячи | **Вызовы инструментов с выводом** — [шаги](./step.md): у них есть тело |
| Порядок строго хронологический | **Фрагмент кода или команда** — [блок кода](../components/display/code.md) |

## Справочник

```api
```

## Доступность

| | |
|---|---|
| `role="log"` | Область, в которую добавляются записи. Скринридер объявит новые строки, не перечитывая весь лог |
| Живость | `role="log"` подразумевает `aria-live="polite"`. Для потока в сотни строк в секунду это надо **отключить** явно: непрерывное озвучивание делает интерфейс неработоспособным |
| Имя области | `aria-label` обязателен: двa лога на экране без имён неразличимы на слух |
| Время — данные | `inst-log-time` берёт `--text-muted` (4.5:1), а не `--text-faint`: таймстемп читают, а не разглядывают |
| Не только цвет | Уровень написан словом (`warn`, `error`) рядом с цветом |
| Прокрутка | Область достижима с клавиатуры. Автопрокрутка вниз обязана останавливаться, когда пользователь прокрутил вверх |
| Перенос | `white-space: pre-wrap` и `overflow-wrap: anywhere` только на ячейке сообщения: на контейнере это схлопнуло бы колонки |

## Связанное

[Шаг](./step.md) · [Лента событий](../components/display/timeline.md) ·
[Код](../components/display/code.md) · [Блок отказа](./failure.md)
