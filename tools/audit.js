/* Проверка кита по ОТРИСОВАННЫМ пикселям.
 *
 * Гейт кита текстовый: contrast и targets читают tokens.css и считают пары,
 * ПЕРЕЧИСЛЕННЫЕ В САМОЙ ПРОВЕРКЕ. Это ловит расхождение кита с документацией
 * и работает в CI без браузера — но не знает, что сложилось на экране после
 * вложения, наложения полупрозрачностей и каскада приложения.
 *
 * Дважды за историю кита текстовая проверка была зелёной, пока в браузере
 * было сломано: `*` внутри комментария закрыл его раньше времени и молча унёс
 * и нейтральный тон, и правило области нажатия. Порогов это не нарушало,
 * потому что в списке пар стояли токены, а не то, что реально нарисовалось.
 *
 * Здесь наоборот: обход каждого узла с текстом, настоящий цвет и настоящий
 * фон со сборкой альфы по предкам. Любую запись цвета — oklch, color-mix,
 * light-dark — приводит к sRGB сам браузер через канву.
 *
 * ПРИМЕНЕНИЕ. Откройте свой экран, вставьте файл в консоль:
 *
 *     kitAudit.run()          вся страница, 5 тем и 3 плотности
 *     kitAudit.run('#main')   только часть
 *     kitAudit.contrast()     один прогон в текущей теме
 *     kitAudit.targets()      один прогон в текущей плотности
 *
 * Проверяет ВАШИ экраны, а не примеры из документации, — в этом весь смысл.
 */
window.kitAudit = (function () {
  'use strict';

  var cv = document.createElement('canvas');
  cv.width = cv.height = 1;
  var ctx = cv.getContext('2d', { willReadFrequently: true });

  /* Разбор любой записи цвета: рисуем на канве и читаем пиксель. Так oklch,
     color-mix и light-dark разрешает браузер, а не мы. */
  function rgba(css) {
    ctx.clearRect(0, 0, 1, 1);
    ctx.fillStyle = '#000';
    ctx.fillStyle = css;                 // невалидное значение оставит чёрный
    ctx.fillRect(0, 0, 1, 1);
    var d = ctx.getImageData(0, 0, 1, 1).data;
    // Канва уже вмешала альфу в чёрный фон, поэтому берём её из нормализованной
    // записи отдельно: цвет нам нужен ЧИСТЫЙ, композицию делаем сами.
    var s = ctx.fillStyle;
    var m = /rgba?\([^)]*?([\d.]+)\s*\)$/.exec(s) || /\/\s*([\d.]+%?)\s*\)/.exec(s);
    var a = 1;
    if (m) { a = parseFloat(m[1]); if (String(m[1]).indexOf('%') >= 0) a /= 100; }
    return [d[0], d[1], d[2], a];
  }

  function over(fg, bg) {
    return [0, 1, 2].map(function (i) { return fg[i] * fg[3] + bg[i] * (1 - fg[3]); });
  }

  function lum(c) {
    function f(v) { v /= 255; return v <= 0.03928 ? v / 12.92 : Math.pow((v + 0.055) / 1.055, 2.4); }
    return 0.2126 * f(c[0]) + 0.7152 * f(c[1]) + 0.0722 * f(c[2]);
  }

  function ratio(a, b) {
    var l = [lum(a), lum(b)].sort(function (x, y) { return y - x; });
    return (l[0] + 0.05) / (l[1] + 0.05);
  }

  /* Настоящий фон под элементом: вверх по предкам, пока не наберётся
     непрозрачность. Полупрозрачные слои складываются в порядке рисования. */
  function bgOf(el) {
    var stack = [];
    for (var n = el; n && n.nodeType === 1; n = n.parentElement) {
      var c = rgba(getComputedStyle(n).backgroundColor);
      if (c[3] > 0) { stack.push(c); if (c[3] === 1) break; }
    }
    var out = [255, 255, 255];
    for (var i = stack.length - 1; i >= 0; i--) out = over(stack[i], out);
    return out;
  }

  /* Свой видимый текст, а не текст потомков: иначе один и тот же абзац
     проверялся бы столько раз, сколько над ним обёрток. */
  function ownText(el) {
    for (var i = 0; i < el.childNodes.length; i++) {
      var n = el.childNodes[i];
      if (n.nodeType === 3 && n.textContent.trim().length) return true;
    }
    return false;
  }

  function root(sel) {
    return (typeof sel === 'string' ? document.querySelector(sel) : sel) || document.body;
  }

  function contrast(sel) {
    var bad = [], checked = 0;
    var all = root(sel).querySelectorAll('*');
    for (var i = 0; i < all.length; i++) {
      var el = all[i];
      if (!ownText(el)) continue;
      var cs = getComputedStyle(el);
      if (cs.visibility === 'hidden' || cs.display === 'none' || cs.opacity === '0') continue;
      var r = el.getBoundingClientRect();
      if (!r.width || !r.height) continue;

      // ПРОЗРАЧНАЯ подпись — не нарушение, а приём, и он описан в конституции.
      // Занятая кнопка гасит текст цветом, а не убирает его из потока: подпись
      // держит ширину, чтобы полоса действий не дёргалась под курсором, и
      // остаётся скринридеру. Меряя её, проверка честно получала 1.00 и валила
      // кнопку за то, что та сделана правильно.
      //
      // Отличать по alpha, а не по «color: transparent»: браузер возвращает
      // вычисленное значение, и rgba(0,0,0,0) от прозрачного акцента не
      // отличается — оба означают «этих пикселей на экране нет».
      if (rgba(cs.color)[3] === 0) continue;

      var size = parseFloat(cs.fontSize);
      var weight = parseInt(cs.fontWeight, 10) || 400;
      // Порог 3:1 — только для КРУПНОГО текста по определению WCAG.
      var large = size >= 24 || (size >= 18.66 && weight >= 700);
      var need = large ? 3 : 4.5;

      var bg = bgOf(el);
      var got = ratio(over(rgba(cs.color), bg), bg);
      checked++;
      if (got < need - 0.01) {
        bad.push({
          где: el.className || el.tagName,
          текст: el.textContent.trim().slice(0, 48),
          кегль: Math.round(size * 10) / 10,
          нужно: need,
          получилось: Math.round(got * 100) / 100
        });
      }
    }
    return { проверено: checked, нарушений: bad.length, список: bad };
  }

  var SEL = 'button, a[href], input:not([type="hidden"]), select, textarea, summary,' +
            '[role="option"], [role="radio"], [role="menuitem"], [role="tab"],' +
            '[tabindex]:not([tabindex="-1"])';

  /* The inline exception, and the only one of the four that applies to a kit:
     a target inside a line of text is exempt. A link in a sentence cannot be
     given a tap area without pushing the line apart, which is why the criterion
     lets it go. Detected structurally — inline-level box whose parent carries
     text of its own — not by guessing from the class name. */
  /* The area a target really offers to a finger.
     A kit keeps small controls SMALL on purpose — a checkbox glyph is 13px —
     and grows an invisible ::before to the tap minimum instead. Measuring the
     element's own box reports the glyph and fails a control that is perfectly
     compliant. Only absolutely positioned pseudo-elements count: an inline one
     decorates, it does not extend the hit area. */
  function hitBox(el) {
    var r = el.getBoundingClientRect();
    var w = r.width, h = r.height;
    ['::before', '::after'].forEach(function (which) {
      var cs = getComputedStyle(el, which);
      if (!cs || cs.content === 'none') return;
      if (cs.position !== 'absolute' && cs.position !== 'fixed') return;
      var pw = parseFloat(cs.width), ph = parseFloat(cs.height);
      if (pw > w) w = pw;
      if (ph > h) h = ph;
    });
    return { w: w, h: h, rect: r };
  }

  /* ПЛОСКОСТЬ, в которой лежит элемент: сам открытый поповер или модалка,
     а для всего остального — null, то есть страница.

     Возвращается именно узел, а не «да/нет». Двух причин достаточно, и обе
     измерены. Открытый поповер лежит в верхнем слое и плавает НАД страницей:
     палец, целящийся в пункт меню, ни во что под ним не попадёт, и мерить
     расстояние между слоями значит мерить то, чего на экране нет. А два
     РАЗНЫХ открытых поповера — это две разные плавающие панели, и соседями
     они друг другу тоже не приходятся.

     Второе важнее, чем кажется: бегунок раскрывает ВСЕ поповеры страницы
     разом, чего в жизни не бывает, и без разделения по плоскостям три демо-
     меню на одной странице начинают мерить расстояние друг до друга. */
  function planeOf(el) {
    return el.closest('[popover]:popover-open, dialog[open]');
  }

  function inlineInText(el) {
    if (getComputedStyle(el).display !== 'inline') return false;

    /* Текст ищется у ближайшего БЛОЧНОГО предка, а не у непосредственного.
       Инлайновые обёртки для этого правила прозрачны: <strong><a>…</a></strong>
       внутри абзаца — это ссылка в строке текста, хотя у самой обёртки, кроме
       ссылки, внутри ничего нет. Пока смотрели на непосредственного родителя,
       такая ссылка теряла освобождение и падала как самостоятельная цель —
       ложное срабатывание на титульной странице справочника. */
    var host = el.parentElement;
    while (host && getComputedStyle(host).display === 'inline') host = host.parentElement;
    if (!host) return false;
    var own = host.textContent.replace(el.textContent, '').trim();
    return own.length > 0;
  }

  /* Цели нажатия, WCAG 2.2 AA (2.5.8), норма 24×24.
     У критерия есть исключение по РАССТОЯНИЮ: цель меньше нормы засчитывается,
     если круг диаметром 24 в её центре не пересекает круг соседа, то есть при
     S + G >= 24. Без него плотный режим был бы нарушением по построению. */
  function targets(sel) {
    var els = [].slice.call(root(sel).querySelectorAll(SEL)).filter(function (e) {
      var cs = getComputedStyle(e);
      if (cs.display === 'none' || cs.visibility === 'hidden') return false;
      var r = e.getBoundingClientRect();
      if (r.width === 0 || r.height === 0) return false;
      // The "visually hidden" trick (a skip link, a screen-reader label) leaves
      // a tiny box that is clipped away entirely: it only becomes a real target
      // on focus. Without this filter it both failed on its own and zeroed a
      // neighbour's gap — a 22px burger failed because an invisible pixel sat
      // against it.
      //
      // The test is the CLIP, not the size. A size threshold was the first
      // attempt and it was wrong twice over: `< 2px` let the kit's own skip
      // link through (declared 1px, rendered 2px, because a border-box element
      // cannot be thinner than its borders), and any threshold big enough to
      // catch it would start excusing genuinely too-small controls.
      if (cs.clipPath.indexOf('inset(50%') === 0) return false;

      // Целей у ЗАКРЫТОГО поповера нет, даже пока он ещё виден.
      //
      // Поповеры типа auto образуют стек: открытие соседнего закрывает
      // предыдущий. Бегунок раскрывает их все подряд, значит открытым
      // остаётся последний, а прочие закрываются — и с
      // transition-behavior: allow-discrete продолжают рисоваться, пока
      // гаснут. Коробка у них при этом настоящая, позиция — умолчание
      // поповера, то есть все они складываются в одну стопку по центру.
      //
      // Так и набирались 244 «нарушения» на восьмидесяти пяти страницах:
      // мерилось расстояние между пунктами меню, которых на экране в этот
      // момент уже нет, а те, что остались, стояли друг на друге.
      var pop = e.closest('[popover]');
      if (pop && !pop.matches(':popover-open')) return false;

      return !inlineInText(e);
    });
    var boxes = els.map(function (e) {
      var hb = hitBox(e);
      return { e: e, r: hb.rect, w: hb.w, h: hb.h };
    });
    var bad = [];
    for (var i = 0; i < boxes.length; i++) {
      var e = boxes[i].e, r = boxes[i].r;
      if (boxes[i].w >= 24 && boxes[i].h >= 24) continue;
      var gap = Infinity;
      for (var j = 0; j < boxes.length; j++) {
        var o = boxes[j];
        if (o.e === e || e.contains(o.e) || o.e.contains(e)) continue;
        // Соседство — вопрос ПЛОСКОСТИ, и у каждой всплывающей панели она своя.
        //
        // Проверка ниже ловила только ПЕРЕКРЫТИЕ, а этого мало: край поповера
        // может оказаться в шести сотых пикселя от ссылки в навигации под ним,
        // не накрыв её. Так и вышло — 244 ложных нарушения на восьмидесяти
        // пяти страницах, все про переключатель плотности в открытом меню.
        if (planeOf(e) !== planeOf(o.e)) continue;
        var dx = Math.max(0, Math.max(r.left - o.r.right, o.r.left - r.right));
        var dy = Math.max(0, Math.max(r.top - o.r.bottom, o.r.top - r.bottom));
        // Перекрывшиеся коробки — тоже слои, а не соседи: закреплённая шапка,
        // проехавшая над карточкой, читается как расстояние ноль и валит всё
        // под собой.
        if (dx === 0 && dy === 0) continue;
        gap = Math.min(gap, Math.sqrt(dx * dx + dy * dy));
      }
      var s = Math.min(boxes[i].w, boxes[i].h);
      // Допуск на субпиксель: у дробных ширин зазор выходит 1.9999 вместо 2,
      // и ровно проходящая норма 22+2 ложно падала бы каждый раз.
      if (gap !== Infinity && s + gap < 23.9) {
        bad.push({
          где: e.className || e.tagName,
          текст: (e.textContent || e.getAttribute('aria-label') || '').trim().slice(0, 32),
          размер: Math.round(boxes[i].w) + '×' + Math.round(boxes[i].h),
          зазор: Math.round(gap * 10) / 10,
          сумма: Math.round((s + gap) * 10) / 10
        });
      }
    }
    return { проверено: boxes.length, нарушений: bad.length, список: bad };
  }

  /* Все пять тем названы ПОИМЁННО, включая нейтральную.
   *
   * Соблазн начать список пустой строкой — «умолчание» — обходится в целую
   * тему: у корня --tint: 0, то есть снятый атрибут вычисляется ровно в
   * light-neutral, и пара «умолчание + light-neutral» меряет одно и то же
   * дважды. Светлая тёплая при этом не меряется вовсе, хотя именно у неё
   * уклон нейтрали ненулевой (--tint: 2) и именно она может увести
   * приглушённый текст под порог.
   *
   * Замерено: у корня и у light-neutral фон страницы oklch(0.976 0 75) —
   * побайтово один; у light он oklch(0.976 0.006 75). */
  var THEMES = ['light-neutral', 'light', 'light-cool', 'dark', 'dark-soft'];
  var DENSITIES = ['', 'compact', 'comfortable'];


  /* ── Пропорции: влезает ли текст в свою коробку ─────────────────────────
   *
   * Границу между этой проверкой и текстовым гейтом стоит назвать прямо.
   * cmd/proportion сторожит ОТНОШЕНИЯ токенов — лестницу кеглей, форму
   * контрола, чётность радиусов; всё это считается из tokens.css и не зависит
   * от шрифта. А влезет ли «Лимит токенов» в свои 92 пикселя, зависит от
   * гарнитуры, начертания и языка — то есть измеримо только на отрисованном.
   *
   * Ровно этого не увидел ни один из четырёх гейтов, когда база кегля выросла
   * с 13 на 14: колонка подписей переполнилась, подпись легла в две строки и
   * порвала выравнивание формы, а проверки остались зелёными.
   */

  var LABELS = '.inst-label,.inst-prop-label,.inst-kv > dt,.inst-metric-label,' +
               '.inst-section-title,.inst-panel-title,.inst-task-meta,.inst-nav-label';

  function proportion(sel) {
    var bad = [], checked = 0;
    var root0 = root(sel);

    /* Подпись, которой не хватило места.
       Мерить канвой нельзя: measureText и раскладка расходятся на пиксель-полтора
       из-за хинтинга и субпиксельных долей, и на коротких строках это давало
       ложные срабатывания там, где колонка меряется содержимым и переполниться
       не может в принципе. Спрашиваем вердикт у самого браузера.

       Случая два, и они разные.
       ОБРЕЗКА: подписи запрещено переноситься (nowrap или многоточие), и она
       вылезла за свою коробку — scrollWidth больше clientWidth.
       ПЕРЕНОС: подпись стоит в колонке ФИКСИРОВАННОЙ ширины, переноситься ей
       не запрещено, и она легла в две строки, сдвинув соседей по сетке. */
    var labels = root0.querySelectorAll(LABELS);
    for (var i = 0; i < labels.length; i++) {
      var el = labels[i];
      var txt = (el.textContent || '').trim();
      if (!txt) continue;
      var cs = getComputedStyle(el);
      if (cs.display === 'none' || cs.visibility === 'hidden') continue;
      var r = el.getBoundingClientRect();
      if (!r.width || !r.height) continue;
      checked++;

      var single = cs.whiteSpace === 'nowrap' || cs.textOverflow === 'ellipsis';
      if (single) {
        var cut = el.scrollWidth > el.clientWidth + 1;
        /* Строке свойства обрезка РАЗРЕШЕНА: имя свойства бывает длиннее любой
           разумной колонки, и конституция это принимает — но требует за это
           `title`, иначе полное имя прочитать нечем. Значит, здесь проверяется
           не сама обрезка, а уплаченная за неё цена. */
        if (el.matches('.inst-prop-label')) {
          if (cut && !el.getAttribute('title')) {
            bad.push({ вид: 'обрезано без title', где: el.className || el.tagName,
                       текст: txt.slice(0, 40) });
          }
          continue;
        }
        if (cut) {
          bad.push({ вид: 'обрезано', где: el.className || el.tagName,
                     текст: txt.slice(0, 40),
                     нужно: el.scrollWidth, есть: el.clientWidth });
        }
        continue;
      }

      /* Фиксированной считаем ту колонку, чью ширину задаёт сетка, а не
         содержимое: у строки свойства и у списка пар это --label-col. */
      var fixed = el.matches('.inst-label, .inst-prop-label') ||
                  (el.tagName === 'DT' && el.parentElement &&
                   el.parentElement.classList.contains('inst-kv') &&
                   !el.parentElement.classList.contains('inst-kv--tight'));
      if (!fixed) continue;

      var lh = parseFloat(cs.lineHeight) || parseFloat(cs.fontSize) * 1.4;
      var lines = Math.round(r.height / lh);
      if (lines > 1) {
        bad.push({ вид: 'перенос', где: el.className || el.tagName,
                   текст: txt.slice(0, 40), строк: lines,
                   колонка: Math.round(el.clientWidth) });
      }
    }

    /* Значок против прописной подписи, рядом с которой он стоит. Значок
       меряется не коробкой, а чернилами: по спрайту кита они занимают 9.02
       из 16, то есть 0.564 коробки. Прописная — около 0.71 кегля. */
    var icons = root0.querySelectorAll('.inst-icon');
    for (var j = 0; j < icons.length; j++) {
      var ic = icons[j];
      /* Служебный глиф в эту меру не попадает. .inst-icon--sm — это размер
         ШЕВРОНА (--size-chevron), и он обязан быть мельче прописной: он не
         содержимое строки, а указатель на раскрытие. Своя полоса у него уже
         есть в текстовом гейте (cmd/proportion, «шеврон к базе кегля»), и
         мерить его здесь второй, чужой мерой значит выдавать верную работу за
         нарушение. */
      if (ic.classList.contains('inst-icon--sm')) continue;
      var host = ic.parentElement;
      if (!host) continue;
      var hostText = '';
      for (var n = host.firstChild; n; n = n.nextSibling) {
        if (n.nodeType === 3) hostText += n.textContent.trim();
      }
      if (!hostText) continue;           /* значок один — сравнивать не с чем */
      var ics = getComputedStyle(ic), hcs = getComputedStyle(host);
      var box = parseFloat(ics.width);
      if (!box) continue;
      ctx.font = hcs.fontWeight + ' ' + hcs.fontSize + ' ' + hcs.fontFamily;
      var m = ctx.measureText('H');
      var cap = m.actualBoundingBoxAscent;
      if (!cap) continue;
      checked++;
      var ink = box * 0.564;
      var d = ink / cap;
      if (d < 0.88 || d > 1.16) {
        bad.push({
          где: host.className || host.tagName,
          текст: hostText.slice(0, 24),
          коробка: box,
          чернила: Math.round(ink * 10) / 10,
          прописная: Math.round(cap * 10) / 10,
          отношение: Math.round(d * 100) / 100
        });
      }
    }

    return { проверено: checked, нарушений: bad.length, список: bad };
  }

  function run(sel) {
    var html = document.documentElement;
    var theme0 = html.getAttribute('data-theme');
    var dens0 = html.getAttribute('data-density');

    /* Смена темы в ките АНИМИРОВАНА. Прочитанный в тот же кадр цвет — это
       интерполяция, а не итог: первый прогон этой проверки дал 16 ложных
       нарушений в тёмных темах именно так. Глушим переходы на время замера. */
    var kill = document.createElement('style');
    kill.textContent = '*,*::before,*::after{transition:none!important;animation:none!important}';
    document.head.appendChild(kill);
    function flush() { return document.body.offsetHeight; }

    var res = { контраст: {}, цели: {}, пропорции: {}, всего: 0 };
    THEMES.forEach(function (t) {
      html.setAttribute('data-theme', t);
      flush();
      var v = contrast(sel);
      res.контраст[t] = v;
      res.всего += v.проверено;
    });
    theme0 ? html.setAttribute('data-theme', theme0) : html.removeAttribute('data-theme');

    DENSITIES.forEach(function (d) {
      d ? html.setAttribute('data-density', d) : html.removeAttribute('data-density');
      flush();
      var v = targets(sel);
      res.цели[d || 'обычная'] = v;
      // Цели входят в общий счёт наравне с контрастом: итоговая строка
      // обещает «и цели в 3 плотностях», и без этой строки она обещала
      // больше, чем считала.
      res.всего += v.проверено;
    });
    DENSITIES.forEach(function (d) {
      d ? html.setAttribute('data-density', d) : html.removeAttribute('data-density');
      flush();
      var p = proportion(sel);
      res.пропорции[d || 'обычная'] = p;
      res.всего += p.проверено;
    });
    dens0 ? html.setAttribute('data-density', dens0) : html.removeAttribute('data-density');
    flush();
    kill.remove();

    report(res);
    return res;
  }

  function report(res) {
    var падений = 0, свод = {};
    Object.keys(res.контраст).forEach(function (k) {
      var v = res.контраст[k];
      падений += v.нарушений;
      свод['тема ' + k] = v.нарушений + ' из ' + v.проверено;
    });
    Object.keys(res.цели).forEach(function (k) {
      var v = res.цели[k];
      падений += v.нарушений;
      свод['цели ' + k] = v.нарушений + ' из ' + v.проверено;
    });
    Object.keys(res.пропорции || {}).forEach(function (k) {
      var v = res.пропорции[k];
      падений += v.нарушений;
      свод['пропорции ' + k] = v.нарушений + ' из ' + v.проверено;
    });
    console.log('%cinstrument · проверка по пикселям', 'font-weight:bold');
    console.table(свод);
    if (!падений) {
      console.log('%c· ' + res.всего + ' замеров: контраст в ' + THEMES.length +
                  ' темах, цели и пропорции в ' + DENSITIES.length + ' плотностях — чисто',
                  'color:green');
      return;
    }
    Object.keys(res.контраст).forEach(function (k) {
      if (res.контраст[k].нарушений) {
        console.group('контраст, тема ' + k);
        console.table(res.контраст[k].список);
        console.groupEnd();
      }
    });
    Object.keys(res.цели).forEach(function (k) {
      if (res.цели[k].нарушений) {
        console.group('цели, плотность ' + k);
        console.table(res.цели[k].список);
        console.groupEnd();
      }
    });
    Object.keys(res.пропорции || {}).forEach(function (k) {
      if (res.пропорции[k].нарушений) {
        console.group('пропорции, плотность ' + k);
        console.table(res.пропорции[k].список);
        console.groupEnd();
      }
    });
  }

  return { run: run, contrast: contrast, targets: targets, proportion: proportion, rgba: rgba, ratio: ratio };
})();
