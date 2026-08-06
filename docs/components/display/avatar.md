---
title: Аватар
group: Отображение данных
status: stable
source: src/data.css
---

Опознавательный знак человека или агента: изображение либо инициалы.

```html preview
<span class="inst-avatar">ИС</span>
<span class="inst-avatar inst-avatar--sm">W1</span>
<span class="inst-avatar inst-avatar--lg">B4</span>
```

## Группа с наложением

```html preview
<span class="inst-avatars">
  <span class="inst-avatar">ИС</span>
  <span class="inst-avatar">АП</span>
  <span class="inst-avatar">МК</span>
  <span class="inst-avatars-more">+4</span>
</span>
```

Кольцо вокруг каждого аватара в группе — не обводка, а **цвет поверхности
под ним**: оно отделяет соседей друг от друга. На другой поверхности группа
потребует переопределить это кольцо, иначе появится ореол.

## Справочник

### Классы

| Класс | Работа |
|---|---|
| `inst-avatar` | Один аватар. Размер `--control-h-md` |
| `inst-avatar--sm` `inst-avatar--lg` | Размер |
| `inst-avatars` | Группа с наложением |
| `inst-avatars-more` | Счётчик «ещё N» в конце группы |

### Токены

`--control-h-sm/md/lg` · `--radius-full` · `--surface-sunken` ·
`--surface-raised` · `--text-xs` · `--weight-medium` · `--size-marker`

## Доступность

| | |
|---|---|
| Аватар — не подпись | Если имя есть рядом текстом, аватар декоративен: `aria-hidden="true"`. Иначе он озвучится как «И С» рядом с уже прочитанным именем |
| Аватар вместо подписи | Тогда он несёт имя: `<img alt="Иннокентий Соколов">` или `aria-label` на контейнере |
| Инициалы | Не являются доступным именем: «ИС» скринридер прочитает по буквам |
| Группа | Дайте контейнеру `aria-label` с полным составом или числом: «7 участников». Иначе прозвучит череда инициалов |
| Счётчик | `inst-avatars-more` содержит «+4» текстом — это данные, а не декорация |
| Контраст | Инициалы берут `--text-secondary` на `--surface-sunken` и держат 4.5:1 в четырёх темах |

## Связанное

[Тег](./tag.md) · [Бейдж](./badge.md) · [Таблица](./table.md) ·
[Лента событий](./timeline.md)

Исходник: `src/data.css` · Почему кит устроен именно так — [конституция](../../../DESIGN.md)
