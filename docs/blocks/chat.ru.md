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
              <button class="inst-btn inst-btn--xs inst-btn--icon inst-btn--ghost" type="button" aria-label="Копировать"><svg class="inst-icon" aria-hidden="true"><use href="#i-copy"/></svg></button>
              <button class="inst-btn inst-btn--xs inst-btn--icon inst-btn--ghost" type="button" aria-label="Повторить"><svg class="inst-icon" aria-hidden="true"><use href="#i-refresh"/></svg></button>
              <button class="inst-btn inst-btn--xs inst-btn--icon inst-btn--ghost" type="button" aria-label="Ещё"><svg class="inst-icon" aria-hidden="true"><use href="#i-more"/></svg></button>
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

        <div class="inst-stack inst-stack--tight">
          <div class="inst-composer">
            <textarea class="inst-textarea" rows="1" placeholder="Опишите задачу или задайте вопрос" aria-label="Опишите задачу или задайте вопрос"></textarea>
            <div class="inst-composer-bar inst-composer-bar--edge">
              <button class="inst-btn inst-btn--sm inst-btn--icon inst-btn--ghost" type="button" aria-label="Отправить"><svg class="inst-icon" aria-hidden="true"><use href="#i-enter"/></svg></button>
            </div>
          </div>

          <div class="inst-composer-bar">
            <button class="inst-btn inst-btn--xs inst-btn--ghost" type="button" aria-haspopup="menu">Авто</button>
            <button class="inst-btn inst-btn--xs inst-btn--icon inst-btn--ghost" type="button" aria-label="Добавить контекст"><svg class="inst-icon" aria-hidden="true"><use href="#i-plus"/></svg></button>
            <button class="inst-btn inst-btn--xs inst-btn--icon inst-btn--ghost" type="button" aria-label="Диктовать"><svg class="inst-icon" aria-hidden="true"><use href="#i-mic"/></svg></button>
            <button class="inst-btn inst-btn--xs inst-btn--icon inst-btn--ghost" type="button" aria-haspopup="menu" aria-label="Ещё"><svg class="inst-icon inst-icon--down" aria-hidden="true"><use href="#i-chevron"/></svg></button>
            <span class="inst-cluster-spacer"></span>
            <button class="inst-btn inst-btn--xs inst-btn--ghost" type="button" aria-haspopup="menu">Opus 5</button>
            <button class="inst-btn inst-btn--xs inst-btn--ghost" type="button" aria-haspopup="menu">Высокое</button>
            <span class="inst-tooltip">
              <button class="inst-btn inst-btn--xs inst-btn--icon inst-btn--ghost" type="button" aria-describedby="tt-ctx-ru" aria-label="Контекст израсходован на 90%">
                <svg class="inst-ring inst-ring--xs" viewBox="0 0 20 20" aria-hidden="true">
                  <circle class="inst-ring-track" cx="10" cy="10" r="8"/>
                  <circle class="inst-ring-fill" cx="10" cy="10" r="8" style="--value:0.9"/>
                </svg>
              </button>
              <span class="inst-tooltip-text inst-tooltip-text--end" role="tooltip" id="tt-ctx-ru">Контекст израсходован на 90%</span>
            </span>
          </div>
        </div>

      </div>
    </div>
  </main>
</div>
```

## Сессия, а не сообщение

Строка под полем держит то, что переживает одну реплику: режим разрешений,
модель, усилие, сколько осталось контекста. Это **не** контролы поля, поэтому
они стоят снаружи его рамки — тем же
[`inst-composer-bar`](../agent/composer.md), который может стоять внутри рамки
и под ней. Отправка принадлежит сообщению и остаётся внутри, у заднего края
поля.

```html preview
<div class="inst-composer-bar">
    <button class="inst-btn inst-btn--xs inst-btn--ghost" type="button" aria-haspopup="menu">Авто</button>
    <span class="inst-cluster-spacer"></span>
    <button class="inst-btn inst-btn--xs inst-btn--ghost" type="button" aria-haspopup="menu">Opus 5</button>
    <button class="inst-btn inst-btn--xs inst-btn--ghost" type="button" aria-haspopup="menu">Высокое</button>
    <span class="inst-badge" data-tone="warn">контекст 90%</span>
  </div>
```

## Выбор среди именованного

```html preview
<div class="inst-popover">
  <div class="inst-menu" role="menu" aria-label="Модели">
    <span class="inst-menu-label">Модели</span>
    <button class="inst-menu-item" role="menuitemradio" aria-checked="false" type="button">Fable 5<span class="inst-menu-shortcut">1</span></button>
    <button class="inst-menu-item" role="menuitemradio" aria-checked="true" type="button">Opus 5<span class="inst-menu-shortcut">2</span></button>
    <button class="inst-menu-item" role="menuitemradio" aria-checked="false" type="button">Sonnet 5<span class="inst-menu-shortcut">3</span></button>
    <button class="inst-menu-item" type="button" role="menuitem" aria-haspopup="menu">Ещё модели</button>
    <hr class="inst-menu-sep">
    <span class="inst-menu-label">Быстрый режим</span>
    <label class="inst-switch"><input type="checkbox">Включить быстрый режим</label>
  </div>
</div>
```

У [меню](../components/overlays/menu.md) уже есть все нужные части: подпись
группы, имя с пояснением под ним, хоткей у дальнего края, `aria-checked` у
действующего и линейка между разделами. Переключатель, положенный в меню,
читается его строкой: имя впереди, состояние у заднего края.

## Что осталось от бюджета

```html preview
<div class="inst-panel" style="max-inline-size:24rem">
  <div class="inst-panel-body inst-stack inst-stack--tight">
    <div>
      <div class="inst-meter-row"><span>Окно контекста</span><span class="inst-meter-value">904.2k / 1M</span></div>
      <div class="inst-share" role="img" aria-label="Окно контекста: занято 90%">
        <span class="inst-share-part" data-tone="running" style="--share:0.62"></span>
        <span class="inst-share-part" data-tone="warn" style="--share:0.14"></span>
        <span class="inst-share-part" data-tone="error" style="--share:0.06"></span>
        <span class="inst-share-part" data-tone="neutral" style="--share:0.08"></span>
      </div>
    </div>
    <div>
      <div class="inst-meter-row"><span>Лимит на 5 часов</span><span class="inst-meter-value">сброс через 4 ч · 10%</span></div>
      <div class="inst-meter" role="progressbar" aria-label="Лимит на 5 часов" aria-valuenow="10" aria-valuemin="0" aria-valuemax="100">
        <div class="inst-meter-fill" style="inline-size:10%"></div>
      </div>
    </div>
    <div>
      <div class="inst-meter-row"><span>Недельный</span><span class="inst-meter-value">сброс в пятницу · 8%</span></div>
      <div class="inst-meter" role="progressbar" aria-label="Недельный" aria-valuenow="8" aria-valuemin="0" aria-valuemax="100">
        <div class="inst-meter-fill" style="inline-size:8%"></div>
      </div>
    </div>
  </div>
</div>
```

Окно контекста — [полоса состава](../components/charts/share.md): по части на
род занимающего, а неучтённое остаётся дорожкой. Каждый лимит под ней —
[мера](../components/charts/meter.md): одна величина против собственного
потолка, для чего мера и есть, а полоса состава — нет.

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
