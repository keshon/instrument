---
title: Запрос подтверждения
group: Агентный слой
layout: component
source: src/agent.css
needs-js: Отправка решения и перевод блока в решённое состояние — слой приложения
api:
  - { name: "inst-approval", kind: "класс", doc: "Блок" }
  - { name: "inst-approval-head", kind: "класс", doc: "Заголовок со значком" }
  - { name: "inst-approval-what", kind: "класс", doc: "Что будет сделано" }
  - { name: "inst-approval-effects", kind: "класс", doc: "Список последствий, на `<ul>`" }
  - { name: "inst-approval-actions", kind: "класс", doc: "Решение" }
  - { name: "inst-approval-verdict", kind: "класс", doc: "Что решили и когда. Показывается после ответа" }
  - { name: "data-state", kind: "атрибут", value: "pending · approved · denied", doc: "на `inst-approval`" }
  - { name: "data-tone", kind: "атрибут", value: "warn · error", doc: "на `<li>` последствия" }
  - { name: "--pad-card", kind: "токен" }
  - { name: "--radius-lg", kind: "токен" }
  - { name: "--space-3", kind: "токен" }
  - { name: "--gap-inline", kind: "токен" }
  - { name: "--text-sm", kind: "токен" }
  - { name: "--text-xs", kind: "токен" }
  - { name: "--warn-bg", kind: "токен" }
  - { name: "--warn-text", kind: "токен" }
  - { name: "--err-text", kind: "токен" }
  - { name: "--text-muted", kind: "токен" }
title-en: "Approval request"
group-en: "Agent layer"
---

Определяющее взаимодействие агентных систем: агент останавливается и
спрашивает разрешения. Всё остальное на экране можно проскроллить — **это
нельзя**, поэтому блок единственный в ките, кто имеет право останавливать
глаз.

```html preview
<div class="inst-approval" data-state="pending" role="group" aria-labelledby="ap1">
  <div class="inst-approval-head" id="ap1">Требуется подтверждение</div>
  <div class="inst-approval-what">Записать 4 файла в <code>terrain/</code></div>
  <ul class="inst-approval-effects">
    <li>heightmap.ts — перезапись</li>
    <li data-tone="error">chunks.bin — удаление</li>
  </ul>
  <div class="inst-approval-actions">
    <button class="inst-btn inst-btn--primary" type="button">Разрешить</button>
    <button class="inst-btn" type="button">Отклонить</button>
  </div>
</div>
```

## Использование

Без любой из них компонент не выполняет свою работу и превращается в
украшенное «вы уверены?».

| Что | Класс | Почему обязательно |
|---|---|---|
| **Что именно будет сделано** | `inst-approval-what` | «Разрешить действие?» не является вопросом: пользователь не знает, на что отвечает |
| **Что это затронет** | `inst-approval-effects` | Список последствий. Необратимые помечаются тоном `error` |
| **Решение одним нажатием** | `inst-approval-actions` | Разрешить и отклонить рядом. Отсутствие «отклонить» превращает запрос в уведомление |

## Когда использовать

| Используйте | Возьмите другое |
|---|---|
| Агент остановился и ждёт разрешения продолжить | **Агент упал** — [блок отказа](./failure.md): там нужна причина и то, что уже пробовали |
| Решение принимается в потоке и остаётся в истории | **Решение обязано быть принято немедленно** — [модалка](../components/overlays/dialog.md): она блокирует всё остальное |
| Последствия перечислимы и их можно показать | **Сообщение без выбора** — [баннер](../components/feedback/banner.md) с тоном `warn` |
| — | **Обычное подтверждение удаления в форме** — модалка. Запрос подтверждения агента живёт в ленте прогона, а не поверх неё |

## Устройство

```html
<ul class="inst-approval-effects">
  <li>heightmap.ts — перезапись</li>
  <li data-tone="warn">config.json — изменение настроек</li>
  <li data-tone="error">chunks.bin — удаление</li>
</ul>
```

Тон здесь означает **необратимость**. Удаление красное потому, что его нельзя
отменить, — а не потому, что оно «плохое».

## Состояния

```html preview
<div class="inst-approval" data-state="pending" role="group" aria-label="Ожидает">
  <div class="inst-approval-head">Требуется подтверждение</div>
  <div class="inst-approval-what">Записать 4 файла</div>
  <div class="inst-approval-actions">
    <button class="inst-btn inst-btn--primary" type="button">Разрешить</button>
    <button class="inst-btn" type="button">Отклонить</button>
  </div>
</div>
<div class="inst-approval" data-state="approved" role="group" aria-label="Разрешено">
  <div class="inst-approval-head">Требуется подтверждение</div>
  <div class="inst-approval-what">Записать 4 файла</div>
  <div class="inst-approval-verdict">Разрешено в 14:32</div>
</div>
```

| `data-state` | Что происходит |
|---|---|
| `pending` | Ожидает ответа. **Единственное**, где показаны действия |
| `approved` | Разрешено. Блок отступает, действия скрыты |
| `denied` | Отклонено. Описание действия зачёркнуто |

**После ответа блок не исчезает и не гаснет — он отступает.** История решений
остаётся читаемой: пользователь должен иметь возможность вернуться и увидеть,
что именно он разрешил час назад. Исчезнувший запрос невозможно проверить.

## JS

Подключите модуль один раз на страницу — инициализировать компоненты по
отдельности не нужно, кит работает делегированием и видит узлы, пришедшие
позже.

```html
<script type="module" src="src/kit.js"></script>
```

### Что делает кит

Ничего. Здесь он честно не при чём: согласие на действие агента — это решение,
а не поведение виджета, и подделать его оформлением нельзя.

### Что должно сделать приложение

1. Отправить решение туда, где оно исполняется.
2. Перевести блок в решённое состояние — иначе кнопки останутся живыми и
   человек нажмёт второй раз.
3. Оставить видимым, **что** было решено: блок согласия — это след в истории,
   а не диалог, который закрылся.

```js
block.addEventListener('click', async (e) => {
  const btn = e.target.closest('[data-decision]');
  if (!btn) return;

  for (const b of block.querySelectorAll('button')) b.disabled = true;
  await api.decide(block.dataset.id, btn.dataset.decision);

  block.dataset.decided = btn.dataset.decision;
});
```

Кнопки блокируются **до** запроса. Сеть занимает секунды, и за это время
«Разрешить» успевают нажать дважды.

## Сценарии

Запрос из настоящего прогона — с командой, списком последствий и пометкой
необратимого.

Обратите внимание на **порядок**: что будет сделано, что это затронет, и
только потом решение. Запрос, начинающийся с кнопок, требует ответа раньше,
чем сообщает вопрос.

```html preview context
<div class="inst-approval" data-state="pending" role="group" aria-labelledby="ap2">
  <div class="inst-approval-head" id="ap2">Агент просит разрешение</div>
  <div class="inst-approval-what">rm -rf build/ &amp;&amp; npm run build</div>
  <ul class="inst-approval-effects">
    <li>Удалит каталог <code>build/</code> целиком — 1 284 файла</li>
    <li data-tone="warn">Необратимо: содержимое не попадает в корзину</li>
    <li>Пересборка займёт около 40 с</li>
  </ul>
  <div class="inst-approval-actions">
    <button class="inst-btn inst-btn--primary inst-btn--sm" type="button">Разрешить</button>
    <button class="inst-btn inst-btn--sm" type="button">Отклонить</button>
    <button class="inst-btn inst-btn--ghost inst-btn--sm" type="button">Разрешать всегда</button>
  </div>
</div>
```

## Правила

:::do Последствия перечислены до кнопки
Человек соглашается с тем, что прочитал. Список последствий — содержание
согласования, а не украшение.
:::

:::dont Согласование без отказа
«Отклонить» обязателен и обязан быть равноправной кнопкой. Согласование с одним
выходом — это уведомление.
:::

:::do Разрушительное действие названо глаголом
«Удалить 4 файла», а не «Продолжить». Кнопка — последнее, что читают.
:::

:::dont Тихое согласование по таймауту
Молчание не является ответом. Истёкшее ожидание — отказ, и об этом надо
сказать словом.
:::

## Доступность

| | |
|---|---|
| Группа | `role="group"` + `aria-labelledby` на заголовок. Иначе блок разваливается на несвязанные абзацы |
| Появление | Запрос возникает по инициативе машины, поэтому обязан попасть в живую область: `aria-live="assertive"` или перевод фокуса на блок. Тихо появившийся запрос будет ждать вечно |
| Порядок действий | «Разрешить» первым: главное действие идёт первым, и это единственный сигнал, какое из двух главное |
| Отклонение доступно всегда | Кнопка «Отклонить» не может быть спрятана в меню. Отказ должен стоить столько же нажатий, сколько согласие |
| Не только цвет | Необратимое последствие помечено тоном **и** словом («удаление»), а не одной краснотой |
| Фокус после ответа | Действия исчезают вместе с ответом — приложение обязано перевести фокус, иначе он провалится в `<body>` |

## API

```api
```

## Связанное

[Блок отказа](./failure.md)
[Шаг](./step.md)
[Диф](./diff.md)
[Модалка](../components/overlays/dialog.md)
[Баннер](../components/feedback/banner.md)
