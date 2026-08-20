/* Скрипт сайта документации. В кит не входит.
   Ровно семь работ: тема, акцент, масштаб, плотность, поиск, копирование и
   вызов тоста в примере — его нельзя показать разметкой, он заводится
   вызовом. */

(() => {
  const root = document.documentElement;

  /* ── Тема, акцент, масштаб и плотность ──────────────────────────────────
     Сайт переключает их теми же атрибутами, что и любое приложение на ките.
     Это не демонстрация ради демонстрации: если темизация где-то протекает,
     видно будет здесь первым.

     Четыре ручки ортогональны: 5 тем × 4 акцента × 5 масштабов × 3 плотности =
     300 сочетаний, и справочник обязан показывать любое из них. Ровно поэтому
     переключатели независимы, а не сведены в один список пресетов: пресет
     скрыл бы, что оси не пересекаются, и первое же протёкшее сочетание
     нашлось бы у пользователя, а не здесь.

     Масштаб добавлен четвёртой осью потому, что тремя плотностями «крупнее»
     было не выразить: просторная давала воздух, а не размер. */

  /* Атрибут на корне: пустое значение — снять, иначе поставить. Одна функция
     на масштаб, тему и акцент, потому что разница между ними ровно в имени
     атрибута; три одинаковых блока разошлись бы при первой правке.

     Выпадашка — меню кита в поповере, а не нативный <select>: список,
     нарисованный операционной системой, был единственным местом в этой панели,
     до которого кит не дотягивался.

     Состояние ведёт САЙТ, а не кит, и это не обход, а контракт: у кита
     `follows: null` на роли menu, потому что пункт меню — действие, а не
     выбор. Здесь пункт всё-таки выбор, поэтому aria-checked переносим сами.
     Стрелки и Escape при этом остаются за китом. */
  const menuKnob = (sel, attr, key) => {
    const menu = document.querySelector(sel);
    if (!menu) return;
    const pop = menu.closest('[popover]');
    const trigger = document.querySelector('[popovertarget="' + pop.id + '"]');
    const items = [...menu.querySelectorAll('[role="menuitem"]')];

    const apply = (v) => {
      if (v) { root.setAttribute(attr, v); localStorage.setItem(key, v); }
      else { root.removeAttribute(attr); localStorage.removeItem(key); }
      for (const x of items) {
        const on = (x.dataset.v || '') === (v || '');
        x.setAttribute('aria-checked', String(on));
        /* Подпись кнопки — то, что выбрано. Иначе панель приходится
           открывать, чтобы вспомнить, на чём остановился. */
        if (on && trigger) trigger.textContent = x.dataset.label || x.textContent.trim();
      }
    };

    apply(localStorage.getItem(key) || '');

    menu.addEventListener('click', (e) => {
      const item = e.target.closest('[role="menuitem"]');
      if (!item) return;
      apply(item.dataset.v || '');
      pop.hidePopover();
    });
  };

  menuKnob('[data-scale-picker]', 'data-scale', 'instrument-scale');
  menuKnob('[data-theme-picker]', 'data-theme', 'instrument-theme');
  menuKnob('[data-accent-picker]', 'data-accent', 'instrument-accent');

  /* Плотность осталась сегментированным контролом: три коротких слова, и все
     три видны сразу — прятать их в список значило бы менять один щелчок на
     два. У масштаба выбор подписан кеглем, а не словом, и там список выигрывает.

     Обе ручки стоят в панели раздельно, а не сведены в один список из пяти
     положений: пресет скрыл бы, что оси не пересекаются, и сочетание «крупно
     и плотно» — самое полезное на большом мониторе — оказалось бы
     невыразимым. */
  const picker = (sel, attr, key) => {
    const group = document.querySelector(sel);
    if (!group) return;
    const items = [...group.querySelectorAll('[role="radio"]')];
    const apply = (v) => {
      if (v === 'default') { delete root.dataset[attr]; localStorage.removeItem(key); }
      else { root.dataset[attr] = v; localStorage.setItem(key, v); }
    };

    /* Начальное состояние ставит сайт: он один знает, что лежит в
       localStorage. Дальше выбор ведёт кит. */
    const saved = localStorage.getItem(key) || 'default';
    const start = items.find((x) => x.dataset.v === saved)
      || items.find((x) => x.dataset.v === 'default') || items[0];
    for (const x of items) {
      x.setAttribute('aria-checked', String(x === start));
      x.tabIndex = x === start ? 0 : -1;
    }
    apply(start.dataset.v);

    /* Дальше — событие кита, а не свой обработчик щелчка и не подглядывание
       за атрибутом после стрелки. Кит переносит aria-checked и по щелчку, и
       по стрелке; сайту остаётся то, чего кит знать не может, — что делать с
       выбранным значением.

       Это и есть проверка поведения на живом: перестань кит работать —
       плотность перестанет переключаться на глазах. */
    group.addEventListener('inst:select', (e) => apply(e.target.dataset.v));
  };

  picker('[data-density-picker]', 'density', 'instrument-density');

  /* ── Поиск ──────────────────────────────────────────────────────────────
     Индекс — один JSON на весь сайт, грузится по первому обращению.
     Готовый поисковик тянул бы бинарник и WASM на каждую страницу; здесь
     весь кит помещается в несколько сотен килобайт. */

  const input = document.querySelector('[data-search]');
  const box = document.querySelector('[data-results]');
  if (input && box) {
    let index = null, active = -1;

    const load = async () => {
      if (index) return index;
      /* Индекс свой у каждого языка: искать по русским телам из английской
         версии значило бы отдавать читателю страницы, которых он не
         прочтёт. Язык берётся из документа, а не из адреса — он там уже
         объявлен, и второй источник разошёлся бы. */
      const lang = document.documentElement.lang || 'ru';
      const file = lang === 'ru' ? '/search.json' : '/' + lang + '-search.json';
      index = await (await fetch(file)).json();
      return index;
    };

    const norm = (s) => s.toLowerCase().replace(/ё/g, 'е');

    /* Разметка собирается узлами, а не строкой. innerHTML со вставкой t/g/r
       сегодня безопасен, потому что содержимое своё, но это мина: первый же
       заголовок страницы с «<» или кавычкой ломает результаты молча. Проект,
       у которого выдуманный токен ловится машиной, не имеет права оставлять
       такое на «пока не случалось». */
    const setOpen = (on) => {
      box.hidden = !on;
      input.setAttribute('aria-expanded', String(on));
      if (!on) {
        active = -1;
        input.removeAttribute('aria-activedescendant');
      }
    };

    const render = (items) => {
      active = -1;
      input.removeAttribute('aria-activedescendant');
      box.replaceChildren();
      if (!items.length) {
        const empty = document.createElement('div');
        empty.className = 'site-result-empty';
        empty.textContent = RU ? 'Ничего не найдено' : 'Nothing found';
        box.append(empty);
      } else {
        items.forEach((p, i) => {
          const a = document.createElement('a');
          a.className = 'site-result';
          a.id = 'site-result-' + i;
          a.setAttribute('role', 'option');
          a.setAttribute('aria-selected', 'false');
          a.href = p.r;
          const t = document.createElement('div');
          t.className = 'site-result-title';
          t.textContent = p.t;
          a.append(t);
          if (p.g) {
            const g = document.createElement('div');
            g.className = 'site-result-group';
            g.textContent = p.g;
            a.append(g);
          }
          box.append(a);
        });
      }
      setOpen(true);
    };

    const search = async () => {
      const q = norm(input.value.trim());
      if (q.length < 2) { setOpen(false); return; }
      const words = q.split(/\s+/);
      /* Веса расставлены по тому, ЧТО человек набрал, а не по тому, где
         совпало. Точное имя класса — это адрес, а не поиск: тот, кто ввёл
         «inst-badge», знает, чего хочет, и ему нельзя показывать пять
         страниц, где это слово просто упомянуто.

         Слаг стоит наравне с заголовком, потому что для латинского запроса
         он И ЕСТЬ заголовок: документация по-русски, API по-английски. */
      const scored = (await load()).map((p) => {
        const t = norm(p.t), s0 = norm(p.s || ''), g = norm(p.g || '');
        const own = norm(p.o || '').split(' ');
        const names = norm(p.n || '').split(' ');
        const b = norm(p.b);
        let s = 0;
        for (const w of words) {
          /* ОПИСЫВАЕТ важнее, чем УПОМИНАЕТ, и разрыв должен быть таким,
             чтобы никакое число совпадений в теле его не перекрыло.
             `inst-btn` встречается на 26 страницах: без этого разрыва
             первой выходила конституция, где перечислено всё. */
          if (own.includes(w)) s += 120;
          else if (own.some((n) => n.startsWith(w))) s += 60;
          else if (names.includes(w)) s += 12;
          else if (names.some((n) => n.startsWith(w))) s += 6;

          if (s0 === w) s += 100;                       // слаг целиком: dialog → Модалка
          else if (s0.startsWith(w)) s += 40;

          if (t.startsWith(w)) s += 12;
          else if (t.includes(w)) s += 6;
          if (g.includes(w)) s += 2;
          if (b.includes(w)) s += 1;
        }
        return { p, s };
      }).filter((x) => x.s > 0).sort((a, b2) => b2.s - a.s).slice(0, 8);
      render(scored.map((x) => x.p));
    };

    input.addEventListener('input', search);
    input.addEventListener('focus', () => { if (input.value.trim().length >= 2) search(); });

    input.addEventListener('keydown', (e) => {
      const opts = [...box.querySelectorAll('.site-result')];
      if (e.key === 'Escape') { setOpen(false); input.blur(); return; }
      if (!opts.length) return;
      if (e.key === 'ArrowDown' || e.key === 'ArrowUp') {
        e.preventDefault();
        active = (active + (e.key === 'ArrowDown' ? 1 : -1) + opts.length) % opts.length;
        opts.forEach((o, i) => o.setAttribute('aria-selected', String(i === active)));
        opts[active].scrollIntoView({ block: 'nearest' });
        /* Фокус остаётся в поле, поэтому «где я» скринридеру сообщает
           aria-activedescendant. Без неё стрелки двигали подсветку молча:
           роль listbox была объявлена и не выполнена. */
        input.setAttribute('aria-activedescendant', opts[active].id);
      } else if (e.key === 'Enter' && active >= 0) {
        e.preventDefault();
        location.href = opts[active].getAttribute('href');
      }
    });

    document.addEventListener('click', (e) => {
      if (!e.target.closest('.site-search')) setOpen(false);
    });

    /* «/» — фокус в поиск. Ожидаемо там, где сидят часами. */
    document.addEventListener('keydown', (e) => {
      if (e.key !== '/' || e.metaKey || e.ctrlKey) return;
      const t = e.target;
      if (t && /^(INPUT|TEXTAREA|SELECT)$/.test(t.tagName)) return;
      e.preventDefault(); input.focus(); input.select();
    });
  }

  /* ── Копирование ───────────────────────────────────────────────────────── */

  /* Копирование делает КИТ: класс inst-copy и data-copy — его контракт.
     Своего обработчика здесь нет намеренно. Копирование стало поведением
     кита, и вторая реализация в его же справочнике означала бы, что кит
     документирует одно, а показывает другое. */

  /* Тост в примере. Разметки у него нет: область строит кит по вызову,
     поэтому единственный честный пример — кнопка, которая его зовёт.
     Параметры лежат в data-demo-toast, чтобы страница оставалась
     разметкой, а не скриптом. */
  document.addEventListener('click', async (e) => {
    const btn = e.target.closest('[data-demo-toast]');
    if (!btn) return;
    const { toast } = await import('/kit/kit.js');
    toast(JSON.parse(btn.dataset.demoToast));
  });

  /* ── Бургер ──────────────────────────────────────────────────────────────
     Выезд ящика рисует кит по aria-expanded. Сайту остаётся ровно то, что
     кит объявил слоем приложения: переключение атрибута, Escape, клик мимо
     и возврат фокуса. Без последних трёх ящик открывается и не закрывается
     ничем, кроме повторного нажатия. */

  const burger = document.querySelector('[data-burger]');
  const drawer = document.getElementById('sidebar');
  if (burger && drawer) {
    const setNav = (open) => {
      burger.setAttribute('aria-expanded', String(open));
      if (open) {
        const first = drawer.querySelector('a, button');
        if (first) first.focus();
      } else if (drawer.contains(document.activeElement)) {
        burger.focus();
      }
    };
    const isOpen = () => burger.getAttribute('aria-expanded') === 'true';

    burger.addEventListener('click', () => setNav(!isOpen()));

    document.addEventListener('keydown', (e) => {
      if (e.key === 'Escape' && isOpen()) setNav(false);
    });

    /* Клик по подложке. Подложка — псевдоэлемент оболочки, поэтому цели у
       события нет: считаем промахом всё, что не ящик и не кнопка. */
    document.addEventListener('click', (e) => {
      if (!isOpen()) return;
      if (e.target.closest('#sidebar') || e.target.closest('[data-burger]')) return;
      setNav(false);
    });
  }

  /* ── Подсветка текущего пункта оглавления ──────────────────────────────── */

  const toc = [...document.querySelectorAll('.site-toc-item')];
  if (toc.length) {
    const byId = new Map(toc.map((a) => [decodeURIComponent(a.getAttribute('href').slice(1)), a]));
    const spy = new IntersectionObserver((entries) => {
      for (const en of entries) {
        if (!en.isIntersecting) continue;
        toc.forEach((a) => a.removeAttribute('aria-current'));
        const a = byId.get(en.target.id);
        if (a) a.setAttribute('aria-current', 'true');
      }
    }, { root: document.querySelector('.inst-shell-main'), rootMargin: '0px 0px -70% 0px' });
    byId.forEach((_, id) => {
      const el = document.getElementById(id);
      if (el) spy.observe(el);
    });
  }
})();
