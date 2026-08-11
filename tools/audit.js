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
 *     kitAudit.run()          вся страница, 6 тем и 3 плотности
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
  function inlineInText(el) {
    var cs = getComputedStyle(el);
    if (cs.display !== 'inline') return false;
    var parent = el.parentElement;
    if (!parent) return false;
    var own = parent.textContent.replace(el.textContent, '').trim();
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
      // The "visually hidden" trick (a skip link, a screen-reader label) leaves
      // a 1×1 box with clip-path: it only becomes a real target on focus.
      // Without this filter it both failed on its own and zeroed a neighbour's
      // gap — a 22px burger failed because an invisible pixel sat against it.
      if (r.width < 2 || r.height < 2) return false;
      return !inlineInText(e);
    });
    var boxes = els.map(function (e) { return { e: e, r: e.getBoundingClientRect() }; });
    var bad = [];
    for (var i = 0; i < boxes.length; i++) {
      var e = boxes[i].e, r = boxes[i].r;
      if (r.width >= 24 && r.height >= 24) continue;
      var gap = Infinity;
      for (var j = 0; j < boxes.length; j++) {
        var o = boxes[j];
        if (o.e === e || e.contains(o.e) || o.e.contains(e)) continue;
        var dx = Math.max(0, Math.max(r.left - o.r.right, o.r.left - r.right));
        var dy = Math.max(0, Math.max(r.top - o.r.bottom, o.r.top - r.bottom));
        // Boxes that overlap are not neighbours, they are layers: a sticky
        // header scrolling over a card reads as distance zero and fails
        // everything under it. Spacing is a question about the plane.
        if (dx === 0 && dy === 0) continue;
        gap = Math.min(gap, Math.sqrt(dx * dx + dy * dy));
      }
      var s = Math.min(r.width, r.height);
      // Допуск на субпиксель: у дробных ширин зазор выходит 1.9999 вместо 2,
      // и ровно проходящая норма 22+2 ложно падала бы каждый раз.
      if (gap !== Infinity && s + gap < 23.9) {
        bad.push({
          где: e.className || e.tagName,
          текст: (e.textContent || e.getAttribute('aria-label') || '').trim().slice(0, 32),
          размер: Math.round(r.width) + '×' + Math.round(r.height),
          зазор: Math.round(gap * 10) / 10,
          сумма: Math.round((s + gap) * 10) / 10
        });
      }
    }
    return { проверено: boxes.length, нарушений: bad.length, список: bad };
  }

  var THEMES = ['', 'light-neutral', 'light-cool', 'dark', 'dark-soft', 'dark-light'];
  var DENSITIES = ['', 'compact', 'comfortable'];

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

    var res = { контраст: {}, цели: {}, всего: 0 };
    THEMES.forEach(function (t) {
      t ? html.setAttribute('data-theme', t) : html.removeAttribute('data-theme');
      flush();
      var v = contrast(sel);
      res.контраст[t || 'умолчание'] = v;
      res.всего += v.проверено;
    });
    theme0 ? html.setAttribute('data-theme', theme0) : html.removeAttribute('data-theme');

    DENSITIES.forEach(function (d) {
      d ? html.setAttribute('data-density', d) : html.removeAttribute('data-density');
      flush();
      res.цели[d || 'обычная'] = targets(sel);
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
    console.log('%cinstrument · проверка по пикселям', 'font-weight:bold');
    console.table(свод);
    if (!падений) {
      console.log('%c· ' + res.всего + ' замеров контраста в 6 темах и цели в 3 плотностях — чисто',
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
  }

  return { run: run, contrast: contrast, targets: targets, rgba: rgba, ratio: ratio };
})();
