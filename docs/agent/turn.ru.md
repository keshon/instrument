---
title: Реплика
group: Агентный слой
layout: component
source: src/agent.css
api:
  - { name: "inst-turn", kind: "класс", doc: "Один шаг обмена" }
  - { name: "data-from", kind: "атрибут", value: "user", doc: "Реплика человека, и только она: она утоплена, реплика машины — умолчание, без поверхности" }
  - { name: "inst-turn-head", kind: "класс", doc: "Кто и когда, кеглем метаданных" }
  - { name: "inst-turn-who", kind: "класс", doc: "Имя в этой строке" }
  - { name: "inst-turn-body", kind: "класс", doc: "Что сказано. Может держать любой компонент этого слоя" }
  - { name: "inst-turn-actions", kind: "класс", doc: "Копировать, повторить, ветвить. Видны всегда, а не по наведению" }
  - { name: "--surface-recessed", kind: "токен" }
  - { name: "--radius-lg", kind: "токен" }
  - { name: "--pad-panel", kind: "токен" }
  - { name: "--space-3", kind: "токен" }
  - { name: "--text-xs", kind: "токен" }
  - { name: "--text-muted", kind: "токен" }
  - { name: "--text-secondary", kind: "токен" }
  - { name: "--weight-medium", kind: "токен" }
  - { name: "--leading-ui", kind: "токен" }
---

Один шаг обмена. Всё остальное в этом слое показывает работу; реплика — то,
**внутри** чего работа происходит: ответ это не абзац, а контейнер для шагов,
диффа, запроса разрешения и бюджета.

```html preview
<div class="inst-stack inst-stack--loose">
  <div class="inst-turn" data-from="user">
    <div class="inst-turn-head"><span class="inst-turn-who">You</span><span>14:02</span></div>
    <div class="inst-turn-body">Добавь полосу состава над историей проверок.</div>
  </div>
  <div class="inst-turn">
    <div class="inst-turn-head"><span class="inst-turn-who">Агент</span><span>14:02</span></div>
    <div class="inst-turn-body">Поставлю её над полосой, легенда будет общая.</div>
    <div class="inst-turn-actions">
      <button class="inst-btn inst-btn--sm inst-btn--ghost" type="button">Копировать</button>
      <button class="inst-btn inst-btn--sm inst-btn--ghost" type="button">Повторить</button>
    </div>
  </div>
</div>
```

## Контракт

| Что | Обязательно | Почему |
|---|---|---|
| `data-from="user"` у реплики человека | да | Реплика машины — умолчание. Помечать обе значило бы завести значение атрибута, которому нечего объявлять |
| Имя в `inst-turn-who` | да | Цвет и отступ не носители: утопление — одна и та же плёнка в любой теме, а скринридер не получает ни того ни другого |
| Время | нет | Полезно в длинной ветке, шум в короткой |
| Видимые `inst-turn-actions` | да, если они есть | Наведение может усилить элемент, но не может вызвать его к жизни: ряд, появляющийся под курсором, на сенсорном экране не существует |
| Ветка как список | нет | Реплики — соседи в `inst-stack--loose`; второго имени для колонки кит не заводит |

Реплика человека **утоплена**, у реплики машины поверхности нет вовсе: именно
она держит шаги, диффы и запросы разрешения, и каждый из них приносит
собственную рамку. Поверхность под ними положила бы плашку внутрь плашки.

### Доступность

| | |
|---|---|
| Кто говорит | Назван текстом, а не поверхностью. Утопление — оформление, и до вспомогательных технологий оно не доходит |
| Поток | Ещё идущий ответ несёт каретку; объявлением служит сам текст |
| Действия | Настоящие кнопки с подписями. Они стоят в потоке, поэтому `Tab` доходит до них в порядке чтения |

## Композиции

### Ответ, который сделал работу

```html preview context
<div class="inst-turn">
  <div class="inst-turn-head"><span class="inst-turn-who">Агент</span><span>14:02</span></div>
  <div class="inst-turn-body">Добавил, легенда теперь расшифровывает тона.</div>
  <details class="inst-step" open>
    <summary class="inst-step-head">
      <span class="inst-step-twist"></span>
      <span class="inst-dot" data-tone="ok"></span>
      <span class="inst-step-name">edit</span>
      <span class="inst-step-sub">src/data.css</span>
      <span class="inst-step-meta">0.4 с</span>
    </summary>
    <div class="inst-step-body">
      <div class="inst-code">.inst-share { display: flex; }</div>
    </div>
  </details>
  <div class="inst-turn-actions">
    <button class="inst-btn inst-btn--sm inst-btn--ghost" type="button">Копировать</button>
    <button class="inst-btn inst-btn--sm inst-btn--ghost" type="button">Повторить</button>
  </div>
</div>
```

## API

```api
```

## Связанное

```related
```
