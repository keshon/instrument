---
title: Строка меню
group: Раскладка
layout: component
source: src/layout.css
api:
  - { name: "inst-menubar", kind: "класс", doc: "Полоса у верхнего края окна. Высота в одну строку. Ряд триггеров и ничего больше" }
  - { name: "inst-menubar-item", kind: "класс", doc: "Один триггер. Настоящая `<button>`: открывает меню и достижима с клавиатуры" }
  - { name: "aria-expanded", kind: "атрибут", value: "true · false", doc: "Отмечает открытый триггер. Нужен только если меню не стоит сразу за ним — иначе полоса видит сам поповер" }
  - { name: "--control-h-sm", kind: "токен", doc: "Высота триггера — та же ступень, что у пункта меню. Остальную область цели добирает зазор между соседями" }
  - { name: "--text-xs", kind: "токен", doc: "Кегль полосы" }
  - { name: "--surface-sunken", kind: "токен" }
  - { name: "--surface-selected", kind: "токен", doc: "Подложка открытого триггера" }
  - { name: "--border-subtle", kind: "токен" }
---

Полоса настольного приложения: Файл, Правка, Вид, Справка. Это ряд
**триггеров** и ничего больше — то, что из них выпадает, это
[поповер](../components/overlays/popover.md) и [меню](../components/overlays/menu.md)
кита, без единой правки. Поэтому компонент такой маленький: не хватало никогда
не меню.

Поповер якорится к кнопке, которая его открыла, и объявлять якоря не нужно —
меню само падает под свой триггер.

Зачем она киту для дашбордов: экран, который является окном, а не страницей,
читается окном во многом благодаря этой полосе, и человек ищет команды сверху
слева независимо от того, что думает кит. В приложении на ките команды спрятали
за одну иконку в правом верхнем углу; владелец приложения посмотрел на экран и
сказал, что непривычно — со второго захода полосу вернули. Полоса ещё и место
для команд, которым больше негде жить: выделить всё, обратить выделение, выход.
Без неё они превращаются в ящик с мелочью или не существуют вовсе.

```html preview
<div class="inst-menubar" role="menubar" aria-label="Главное меню">
  <button class="inst-menubar-item" type="button" role="menuitem" popovertarget="mbru-file">Файл</button>
  <div class="inst-popover inst-popover--anchored" id="mbru-file" popover>
    <div class="inst-menu" role="menu">
      <button class="inst-menu-item" type="button" role="menuitem">
        Открыть папку… <span class="inst-menu-shortcut"><kbd>Ctrl</kbd><kbd>O</kbd></span></button>
      <button class="inst-menu-item" type="button" role="menuitem">
        Пересмотреть <span class="inst-menu-shortcut"><kbd>F5</kbd></span></button>
      <span class="inst-menu-sep"></span>
      <button class="inst-menu-item" type="button" role="menuitem" data-tone="error">Выход</button>
    </div>
  </div>

  <button class="inst-menubar-item" type="button" role="menuitem" popovertarget="mbru-view">Вид</button>
  <div class="inst-popover inst-popover--anchored" id="mbru-view" popover>
    <div class="inst-menu" role="menu">
      <span class="inst-menu-label">Тема</span>
      <button class="inst-menu-item" type="button" role="menuitem" aria-checked="true">Тёмная</button>
      <button class="inst-menu-item" type="button" role="menuitem">Светлая</button>
    </div>
  </div>

  <button class="inst-menubar-item" type="button" role="menuitem" popovertarget="mbru-help">Справка</button>
  <div class="inst-popover inst-popover--anchored" id="mbru-help" popover>
    <div class="inst-menu" role="menu">
      <button class="inst-menu-item" type="button" role="menuitem">О программе</button>
    </div>
  </div>
</div>
```

## Контракт

| Что | Обязательно | Почему |
|---|---|---|
| Настоящая `<button>` для триггера | да | Он открывает меню и достижим с клавиатуры. Ряд `<div>` не говорит ни того, ни другого |
| Меню стоит сразу за своим триггером | рекомендуется | Тогда полоса сама отмечает открытый. Иначе разметка обязана ставить `aria-expanded` |
| Имена, а не глаголы с дополнением | да | «Файл», а не «Операции с файлом». Полоса — это указатель, а указатель, объясняющий себя, перестаёт просматриваться |
| Одна полоса на окно | да | То же рассуждение, что у строки состояния: она принадлежит окну, а окно одно |
| Полоса стоит в [шелле](./shell.md) | рекомендуется | Шелл держит под неё строку. Снаружи приложение обязано само владеть высотой окна |

### Доступность

| | |
|---|---|
| Роль | `role="menubar"` на полосе, `role="menuitem"` на триггере и `aria-label` с именем полосы. Без метки скринридер объявит «строка меню» и ничего больше |
| Область нажатия | `--control-h-sm` плюс зазор между соседями. Ступенью ниже стоит строка состояния, и на мелкой шкале она не добирает 24 пикселей: показания читают, триггеры нажимают. Проверяется `cmd/targets` |
| Открытое состояние | Несёт подложка триггера. Кнопка с `popovertarget` держит развёрнутость в дереве доступности, поэтому читателю скажут; подложка говорит всем остальным |
| Контраст | `--text-secondary` на `--surface-sunken`, и `--text-primary` на `--surface-selected` в открытом состоянии |
| Порядок обхода | Совпадает с порядком в разметке. Стрелки влево-вправо между триггерами — дело приложения: поведение кит здесь не поставляет |

## API

```api
```

## Related

```related
```
