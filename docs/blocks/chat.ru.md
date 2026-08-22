---
title: Чат с ИИ
group: Сборки
source: src/kit.css
template: splash
---

Разговор с машиной, которая по ходу ответа делает работу: рассуждение, вызовы
инструментов и запрос разрешения стоят внутри одной реплики. Ни одного класса,
кроме библиотечных, и ни одной строки скрипта.

```html preview context
<div class="inst-shell">
  <aside class="inst-shell-aside">
    <a class="inst-shell-brand" href="#"><span class="inst-dot" data-tone="running"></span>instrument</a>
    <nav class="inst-nav" aria-label="Разделы">
      <div class="inst-nav-group">
        <span class="inst-nav-label">Сегодня</span>
        <a class="inst-nav-item" href="#" aria-current="page">Полоса состава</a>
        <a class="inst-nav-item" href="#">Пол у ячейки оси</a>
        <a class="inst-nav-item" href="#">Выпуск 0.8.0</a>
      </div>
    </nav>
  </aside>

  <header class="inst-shell-header">
    <nav aria-label="Хлебные крошки">
      <ol class="inst-crumbs">
        <li><a href="#">instrument</a></li>
        <li><span aria-current="page">Полоса состава</span></li>
      </ol>
    </nav>
    <span class="inst-cluster-spacer"></span>
    <span class="inst-statusbar-item"><span class="inst-dot" data-tone="running"></span>работает</span>
  </header>

  <main class="inst-shell-main">
    <div class="inst-container">
      <div class="inst-stack">

        <div class="inst-stack inst-stack--loose">

          <div class="inst-turn" data-from="user">
            <div class="inst-turn-head"><span class="inst-turn-who">Вы</span><span>14:02</span></div>
            <div class="inst-turn-body">Добавь полосу состава над историей проверок, и пусть легенда будет общая.</div>
          </div>

          <div class="inst-turn">
            <div class="inst-turn-head"><span class="inst-turn-who">Агент</span><span>14:02</span></div>

            <details class="inst-thinking">
              <summary class="inst-thinking-head"><span class="inst-step-twist"></span>Думал 4 с</summary>
              <div class="inst-thinking-body">История отвечает на «когда», состав — на «сколько». Образец легенды уже падает в тон, значит расшифрует оба без второго словаря.</div>
            </details>

            <div class="inst-turn-body">Поставлю полосу над историей, а образец пусть читает тон.</div>

            <details class="inst-step" open>
              <summary class="inst-step-head">
                <span class="inst-step-twist"></span>
                <span class="inst-dot" data-tone="ok"></span>
                <span class="inst-step-name">read</span>
                <span class="inst-step-sub">src/data.css</span>
                <span class="inst-step-meta">0.2 s</span>
              </summary>
              <div class="inst-step-body">
                <div class="inst-code">.inst-legend-swatch { background: var(--series, var(--tone-mark)); }</div>
              </div>
            </details>

            <details class="inst-step">
              <summary class="inst-step-head">
                <span class="inst-step-twist"></span>
                <span class="inst-dot" data-tone="ok"></span>
                <span class="inst-step-name">edit</span>
                <span class="inst-step-sub">src/data.css</span>
                <span class="inst-step-meta">0.4 s</span>
              </summary>
              <div class="inst-step-body">Добавил <code>.inst-share</code>, по части на исход.</div>
            </details>

            <div class="inst-turn-actions">
              <button class="inst-btn inst-btn--sm inst-btn--ghost" type="button">Копировать</button>
              <button class="inst-btn inst-btn--sm inst-btn--ghost" type="button">Повторить</button>
            </div>
          </div>

          <div class="inst-turn" data-from="user">
            <div class="inst-turn-head"><span class="inst-turn-who">Вы</span><span>14:05</span></div>
            <div class="inst-turn-body">Теперь убери правила приложения, которые это заменило.</div>
          </div>

          <div class="inst-turn">
            <div class="inst-turn-head"><span class="inst-turn-who">Агент</span><span>14:05</span></div>
            <div class="inst-turn-body">Уходит шестьдесят строк. Файл переписывается на месте, нужно разрешение.<span class="inst-caret"></span></div>

            <div class="inst-approval" data-state="pending">
              <div class="inst-approval-head">Удалить 60 строк из static/beacon.css</div>
              <div class="inst-approval-what">static/beacon.css</div>
              <ul class="inst-approval-effects">
                <li>Ось переезжает в кит</li>
                <li data-tone="warn">Необратимо: файл переписывается на месте</li>
              </ul>
              <div class="inst-approval-actions">
                <button class="inst-btn inst-btn--sm inst-btn--primary" type="button">Разрешить</button>
                <button class="inst-btn inst-btn--sm" type="button">Отклонить</button>
              </div>
            </div>
          </div>

        </div>

        <div class="inst-composer">
          <textarea class="inst-textarea" rows="3" placeholder="Опишите задачу" aria-label="Опишите задачу"></textarea>
          <div class="inst-composer-bar">
            <button class="inst-btn inst-btn--sm inst-btn--ghost" type="button">Приложить</button>
            <span class="inst-badge">opus</span>
            <span class="inst-cluster-spacer"></span>
            12 480 / 200 000
            <button class="inst-btn inst-btn--sm inst-btn--primary" type="button">Отправить</button>
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
| Каркас | [Оболочка приложения](../layout/shell.md) |
| Один шаг обмена | [Реплика](../agent/turn.md) — утоплена у человека, без поверхности у машины |
| Черновик перед ответом | [Размышление](../agent/thinking.md), закрытое |
| Вызов инструмента | [Шаг](../agent/step.md) с аргументами и результатом |
| Остановка за разрешением | [Запрос разрешения](../agent/approval.md) |
| Ещё идущий текст | Каретка, из [бейджа и точки](../components/display/badge.md) |
| Написание следующей реплики | [Поле ввода](../agent/composer.md) |

## Что здесь видно

**Ответ — контейнер, а не абзац.** Всё, что машина сделала, лежит внутри
отвечающей реплики в том порядке, в каком происходило. Ничего не висит в
боковой панели, поэтому не возникает вопроса, к какому прогону относился шаг.

**Человек утоплен, машина без поверхности.** Ответ держит вещи в рамках — шаги,
запрос разрешения, — и поверхность под ними была бы плашкой внутри плашки.
Вопросу рамка не нужна, чтобы его нашли: он короткий.

**Размышление — не шаг.** Шаг это вызов с результатом, который может
провалиться; у размышления статуса нет, поэтому нет ни точки, ни тона, и оно
закрыто, пока его не спросят.

**Поле ввода — одна рамка.** Поле и его контролы это один объект, поэтому фокус
подсвечивает объект, в котором находится клавиатура, а не его половину.

## Связанное

[Реплика](../agent/turn.md) [Размышление](../agent/thinking.md)
[Поле ввода](../agent/composer.md) [Шаг](../agent/step.md)
[Запрос разрешения](../agent/approval.md) [Консоль прогонов](./console.md)
