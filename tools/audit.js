/* Auditing the kit by RENDERED pixels.
 *
 * The kit's own gate is textual: contrast and targets read tokens.css and walk
 * the pairs LISTED INSIDE THE CHECK ITSELF. That catches the kit drifting away
 * from its documentation and runs in CI without a browser — but it knows
 * nothing about what ended up on screen after nesting, stacked translucency
 * and the application's cascade.
 *
 * Twice in the kit's history the textual check stayed green while the browser
 * was broken: a `*` inside a comment closed it early and silently carried off
 * both the neutral tone and the tap-area rule. No threshold was violated,
 * because the list of pairs held tokens rather than what actually got painted.
 *
 * Here it is the other way round: walk every node that carries text, take the
 * real colour and the real background with alpha composited up the ancestors.
 * Any colour notation — oklch, color-mix, light-dark — is resolved to sRGB by
 * the browser itself, through a canvas.
 *
 * USE. Open your own screen and paste the file into the console:
 *
 *     kitAudit.run()          the whole page, 5 themes and 3 densities
 *     kitAudit.run('#main')   a part of it
 *     kitAudit.contrast()     one pass in the current theme
 *     kitAudit.targets()      one pass in the current density
 *
 * It checks YOUR screens rather than the examples in the documentation — that
 * is the entire point.
 */
window.kitAudit = (function () {
  'use strict';

  var cv = document.createElement('canvas');
  cv.width = cv.height = 1;
  var ctx = cv.getContext('2d', { willReadFrequently: true });

  /* Parsing any colour notation: paint it on a canvas and read the pixel. That
     way oklch, color-mix and light-dark are resolved by the browser, not us. */
  function rgba(css) {
    ctx.clearRect(0, 0, 1, 1);
    ctx.fillStyle = '#000';
    ctx.fillStyle = css;                 // an invalid value leaves black
    ctx.fillRect(0, 0, 1, 1);
    var d = ctx.getImageData(0, 0, 1, 1).data;
    // The canvas has already mixed alpha into the black backdrop, so take it
    // from the normalised notation separately: we need the colour PURE and do
    // the compositing ourselves.
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

  /* The real background under an element: up the ancestors until opacity adds
     up. Translucent layers are composited in paint order. */
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

  /* Its own visible text, not the text of descendants: otherwise the same
     paragraph would be checked once per wrapper above it. */
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

      // A TRANSPARENT label is not a violation but a technique, and the design
      // principles describe it. A busy button dims its text with colour rather
      // than pulling it out of flow: the label holds the width so the action
      // bar does not twitch under the cursor, and it stays available to a
      // screen reader. Measuring it, the check honestly got 1.00 and failed
      // the button for being built correctly.
      //
      // Tell it apart by alpha rather than by `color: transparent`: the
      // browser returns the computed value, and rgba(0,0,0,0) is
      // indistinguishable from a transparent accent — both mean "these pixels
      // are not on screen".
      if (rgba(cs.color)[3] === 0) continue;

      var size = parseFloat(cs.fontSize);
      var weight = parseInt(cs.fontWeight, 10) || 400;
      // The 3:1 threshold is for LARGE text only, by the WCAG definition.
      var large = size >= 24 || (size >= 18.66 && weight >= 700);
      var need = large ? 3 : 4.5;

      var bg = bgOf(el);
      var got = ratio(over(rgba(cs.color), bg), bg);
      checked++;
      if (got < need - 0.01) {
        bad.push({
          where: el.className || el.tagName,
          text: el.textContent.trim().slice(0, 48),
          size: Math.round(size * 10) / 10,
          need: need,
          got: Math.round(got * 100) / 100
        });
      }
    }
    return { checked: checked, failed: bad.length, list: bad };
  }

  var SEL = 'button, a[href], input:not([type="hidden"]), select, textarea, summary,' +
            '[role="option"], [role="radio"], [role="menuitem"], [role="tab"],' +
            '[tabindex]:not([tabindex="-1"])';

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

  /* The PLANE an element lies in: the open popover or modal itself, and null —
     that is, the page — for everything else.

     What comes back is the node rather than a yes/no. Two reasons are enough,
     and both are measured. An open popover sits in the top layer and floats
     ABOVE the page: a finger aiming at a menu item will hit nothing underneath
     it, and measuring the distance between layers means measuring something
     that is not on screen. And two DIFFERENT open popovers are two different
     floating panels, so they are not neighbours to each other either.

     The second reason matters more than it looks: the runner opens EVERY
     popover on the page at once, which never happens in life, and without
     splitting by plane three demo menus on one page start measuring the
     distance to each other. */
  function planeOf(el) {
    return el.closest('[popover]:popover-open, dialog[open]');
  }

  /* The inline exception, and the only one of the four that applies to a kit:
     a target inside a line of text is exempt. A link in a sentence cannot be
     given a tap area without pushing the line apart, which is why the criterion
     lets it go. Detected structurally — inline-level box whose parent carries
     text of its own — not by guessing from the class name. */
  function inlineInText(el) {
    if (getComputedStyle(el).display !== 'inline') return false;

    /* The text is looked for on the nearest BLOCK ancestor rather than on the
       immediate one. Inline wrappers are transparent to this rule:
       <strong><a>…</a></strong> inside a paragraph is a link in a line of
       text, even though the wrapper itself holds nothing but the link. Look at
       the immediate parent instead and such a link loses its exemption and
       fails as a standalone target — a false positive on the front page of the
       reference. */
    var host = el.parentElement;
    while (host && getComputedStyle(host).display === 'inline') host = host.parentElement;
    if (!host) return false;
    var own = host.textContent.replace(el.textContent, '').trim();
    return own.length > 0;
  }

  /* Tap targets, WCAG 2.2 AA (2.5.8), the minimum being 24×24.
     The criterion has a DISTANCE exception: a target below the minimum counts
     if a 24-diameter circle at its centre does not cross a neighbour's circle,
     that is when S + G >= 24. Without it the compact mode would be a violation
     by construction. */
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

      // A CLOSED popover has no targets, even while it is still visible.
      //
      // Popovers of type auto form a stack: opening a sibling closes the
      // previous one. The runner opens them all in turn, so the last one stays
      // open and the rest close — and with
      // transition-behavior: allow-discrete they keep being painted while they
      // fade. Their box is real meanwhile, and their position is the popover
      // default, so they all pile up in one stack at the centre.
      //
      // That is how 244 "violations" accumulated across eighty-five pages: the
      // distance was measured between menu items that are no longer on screen,
      // while the ones that remained stood on top of each other.
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
        // Neighbourhood is a question of PLANE, and every floating panel has
        // its own.
        //
        // The test below caught OVERLAP only, and that is not enough: the edge
        // of a popover can land six hundredths of a pixel from a navigation
        // link underneath without covering it. Which is exactly what happened
        // — 244 false violations across eighty-five pages, all of them about
        // the density switch inside an open menu.
        if (planeOf(e) !== planeOf(o.e)) continue;
        var dx = Math.max(0, Math.max(r.left - o.r.right, o.r.left - r.right));
        var dy = Math.max(0, Math.max(r.top - o.r.bottom, o.r.top - r.bottom));
        // Overlapping boxes are layers too, not neighbours: a sticky header
        // that has scrolled over a card reads as distance zero and fails
        // everything underneath it.
        if (dx === 0 && dy === 0) continue;
        gap = Math.min(gap, Math.sqrt(dx * dx + dy * dy));
      }
      var s = Math.min(boxes[i].w, boxes[i].h);
      // Subpixel tolerance: with fractional widths the gap comes out as 1.9999
      // instead of 2, and a 22+2 that passes exactly would fail every time.
      if (gap !== Infinity && s + gap < 23.9) {
        bad.push({
          where: e.className || e.tagName,
          text: (e.textContent || e.getAttribute('aria-label') || '').trim().slice(0, 32),
          box: Math.round(boxes[i].w) + '×' + Math.round(boxes[i].h),
          gap: Math.round(gap * 10) / 10,
          sum: Math.round((s + gap) * 10) / 10
        });
      }
    }
    return { checked: boxes.length, failed: bad.length, list: bad };
  }

  /* All five themes are named EXPLICITLY, the neutral one included.
   *
   * The temptation to start the list with an empty string — "the default" —
   * costs a whole theme: the root carries --tint: 0, so a removed attribute
   * computes to exactly light-neutral, and the pair "default + light-neutral"
   * measures the same thing twice. The warm light theme is then not measured
   * at all, even though it is the one whose neutral tint is non-zero
   * (--tint: 2) and the one that can push muted text below the threshold.
   *
   * Measured: at the root and at light-neutral the page background is
   * oklch(0.976 0 75) — byte for byte the same; at light it is
   * oklch(0.976 0.006 75). */
  var THEMES = ['light-neutral', 'light', 'light-cool', 'dark', 'dark-soft'];
  var DENSITIES = ['', 'compact', 'comfortable'];
  /* The empty string is a removed attribute, that is petrol: the kit has no
     separate [data-accent="petrol"], the default is it. */
  var ACCENTS = ['', 'graphite', 'indigo', 'clay'];
  var SCALES = ['', '15', '16', '17', '18'];


  /* ── Proportion: does the text fit its box ──────────────────────────────
   *
   * The border between this check and the textual gate is worth naming
   * outright. cmd/proportion guards the RATIOS between tokens — the type
   * ladder, the shape of a control, the parity of radii; all of it is computed
   * from tokens.css and does not depend on the font. Whether "Token limit"
   * fits into its 92 pixels depends on the family, the weight and the
   * language — that is, it is measurable only on what was rendered.
   *
   * This is precisely what none of the four gates saw when the base type size
   * grew from 13 to 14: the label column overflowed, the label wrapped onto
   * two lines and tore the alignment of the form apart, and the checks stayed
   * green.
   */

  var LABELS = '.inst-label,.inst-prop-label,.inst-kv > dt,.inst-metric-label,' +
               '.inst-section-title,.inst-panel-title,.inst-task-meta,.inst-nav-label';

  function proportion(sel) {
    var bad = [], checked = 0;
    var root0 = root(sel);

    /* A label that ran out of room.
       Measuring with a canvas will not do: measureText and layout disagree by
       a pixel and a half because of hinting and subpixel fractions, and on
       short strings that produced false positives where the column is measured
       by its content and cannot overflow in principle. Ask the browser itself
       for the verdict.

       There are two cases, and they are different.
       CLIPPING: the label is forbidden to wrap (nowrap or an ellipsis) and it
       has escaped its box — scrollWidth exceeds clientWidth.
       WRAPPING: the label sits in a column of FIXED width, is not forbidden to
       wrap, and has laid out over two lines, pushing its grid neighbours. */
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
        /* A property row is ALLOWED to clip: a property name can be longer
           than any sensible column, and the design principles accept that —
           but they charge a `title` for it, otherwise there is no way to read
           the full name. So what is checked here is not the clipping itself
           but the price paid for it. */
        if (el.matches('.inst-prop-label')) {
          if (cut && !el.getAttribute('title')) {
            bad.push({ kind: 'clipped without title', where: el.className || el.tagName,
                       text: txt.slice(0, 40) });
          }
          continue;
        }
        if (cut) {
          bad.push({ kind: 'clipped', where: el.className || el.tagName,
                     text: txt.slice(0, 40),
                     need: el.scrollWidth, have: el.clientWidth });
        }
        continue;
      }

      /* Fixed means the column whose width is set by the grid rather than by
         its content: for a property row and a key-value list that is
         --label-col. */
      var fixed = el.matches('.inst-label, .inst-prop-label') ||
                  (el.tagName === 'DT' && el.parentElement &&
                   el.parentElement.classList.contains('inst-kv') &&
                   !el.parentElement.classList.contains('inst-kv--tight'));
      if (!fixed) continue;

      var lh = parseFloat(cs.lineHeight) || parseFloat(cs.fontSize) * 1.4;
      var lines = Math.round(r.height / lh);
      if (lines > 1) {
        bad.push({ kind: 'wrapped', where: el.className || el.tagName,
                   text: txt.slice(0, 40), lines: lines,
                   column: Math.round(el.clientWidth) });
      }
    }

    /* An icon against the cap height of the label it stands next to. The icon
       is measured by its ink rather than by its box: across the kit's sprite
       the ink takes 9.02 of 16, that is 0.564 of the box.

       The band 0.95…1.07 is not taste but a measurement: 106 035 ratios over
       all 85 pages in all fifteen cells of scale and density lie within
       0.974…1.041. The earlier 0.88…1.16 admitted 1.16 against an observed
       maximum of 1.041 and sat three thousandths below the single outlier —
       that is, it was fitted to a defect rather than derived from the corpus.

       The margin on each side is about 0.025, and it is SMALLER than the step
       this quantity can move by: the box and the cap height are integers, and
       one device pixel of cap height shifts the ratio by roughly 0.1. That is
       not fragility of the band but its point: an icon that failed to step up
       with the type size is exactly the defect the check exists for, and the
       band must not stay quiet about it. */
    var icons = root0.querySelectorAll('.inst-icon');
    for (var j = 0; j < icons.length; j++) {
      var ic = icons[j];
      /* A utility glyph is out of scope for this measure. .inst-icon--sm is
         the size of the CHEVRON (--size-chevron), and it has to be smaller
         than the cap height: it is not content of the line but a pointer at
         disclosure. It already has a band of its own in the textual gate
         (cmd/proportion, "chevron against the base type size"), and measuring
         it here by a second, foreign measure means reporting correct work as a
         violation. */
      if (ic.classList.contains('inst-icon--sm')) continue;
      var host = ic.parentElement;
      if (!host) continue;
      var hostText = '';
      for (var n = host.firstChild; n; n = n.nextSibling) {
        if (n.nodeType === 3) hostText += n.textContent.trim();
      }
      if (!hostText) continue;           /* a lone icon — nothing to compare */
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
      if (d < 0.95 || d > 1.07) {
        bad.push({
          where: host.className || host.tagName,
          text: hostText.slice(0, 24),
          box: box,
          ink: Math.round(ink * 10) / 10,
          cap: Math.round(cap * 10) / 10,
          ratio: Math.round(d * 100) / 100
        });
      }
    }

    return { checked: checked, failed: bad.length, list: bad };
  }

  function run(sel) {
    var html = document.documentElement;
    var theme0 = html.getAttribute('data-theme');
    var dens0 = html.getAttribute('data-density');

    /* Switching the theme in the kit is ANIMATED. A colour read in the same
       frame is an interpolation rather than the result: the first pass of this
       check produced 16 false violations in the dark themes exactly that way.
       Kill transitions for the duration of the measurement. */
    var kill = document.createElement('style');
    kill.textContent = '*,*::before,*::after{transition:none!important;animation:none!important}';
    document.head.appendChild(kill);
    function flush() { return document.body.offsetHeight; }

    var accent0 = html.getAttribute('data-accent');
    var scale0 = html.getAttribute('data-scale');

    /* A matrix rather than a single cell.
     *
     * The check used to measure 5 themes at ONE accent and 3 densities at ONE
     * scale — roughly six cells out of three hundred. The only thing that sees
     * real pixels was looking through a keyhole.
     *
     * The axes are split by what a quantity depends on rather than merged into
     * one common matrix: colour does not depend on scale, geometry does not
     * depend on theme, and multiplying everything by everything would multiply
     * the runtime twentyfold to repeat the same numbers.
     *
     * Contrast goes by theme × accent: an accent redefines --accent-*, and a
     * label on an accent fill is different in every theme.
     * Targets and proportion go by scale × density: both axes move geometry,
     * and both round to whole pixels independently of each other. It was
     * exactly this product that turned up the icon lagging behind the type
     * size at scale 15.
     */
    var res = { contrast: {}, targets: {}, proportion: {}, total: 0 };
    THEMES.forEach(function (t) {
      html.setAttribute('data-theme', t);
      ACCENTS.forEach(function (a) {
        a ? html.setAttribute('data-accent', a) : html.removeAttribute('data-accent');
        flush();
        var v = contrast(sel);
        res.contrast[t + ' · ' + (a || 'petrol')] = v;
        res.total += v.checked;
      });
    });
    theme0 ? html.setAttribute('data-theme', theme0) : html.removeAttribute('data-theme');
    accent0 ? html.setAttribute('data-accent', accent0) : html.removeAttribute('data-accent');

    SCALES.forEach(function (sc) {
      sc ? html.setAttribute('data-scale', sc) : html.removeAttribute('data-scale');
      DENSITIES.forEach(function (d) {
        d ? html.setAttribute('data-density', d) : html.removeAttribute('data-density');
        flush();
        var key = (sc || '14') + ' · ' + (d || 'default');
        // Targets go into the overall count alongside contrast: the closing
        // line promises "and targets", and without this line it promised more
        // than it counted.
        var v = targets(sel);
        res.targets[key] = v;
        res.total += v.checked;
        var pr = proportion(sel);
        res.proportion[key] = pr;
        res.total += pr.checked;
      });
    });
    dens0 ? html.setAttribute('data-density', dens0) : html.removeAttribute('data-density');
    scale0 ? html.setAttribute('data-scale', scale0) : html.removeAttribute('data-scale');
    flush();
    kill.remove();

    report(res);
    return res;
  }

  function report(res) {
    var failures = 0, summary = {};
    Object.keys(res.contrast).forEach(function (k) {
      var v = res.contrast[k];
      failures += v.failed;
      summary['theme ' + k] = v.failed + ' of ' + v.checked;
    });
    Object.keys(res.targets).forEach(function (k) {
      var v = res.targets[k];
      failures += v.failed;
      summary['targets ' + k] = v.failed + ' of ' + v.checked;
    });
    Object.keys(res.proportion || {}).forEach(function (k) {
      var v = res.proportion[k];
      failures += v.failed;
      summary['proportion ' + k] = v.failed + ' of ' + v.checked;
    });
    console.log('%cinstrument · pixel audit', 'font-weight:bold');
    console.table(summary);
    if (!failures) {
      console.log('%c· ' + res.total + ' measurements: contrast over ' + THEMES.length +
                  '×' + ACCENTS.length + ' themes and accents, targets and proportion over ' +
                  SCALES.length + '×' + DENSITIES.length + ' scales and densities — clean',
                  'color:green');
      return;
    }
    Object.keys(res.contrast).forEach(function (k) {
      if (res.contrast[k].failed) {
        console.group('contrast, theme ' + k);
        console.table(res.contrast[k].list);
        console.groupEnd();
      }
    });
    Object.keys(res.targets).forEach(function (k) {
      if (res.targets[k].failed) {
        console.group('targets, density ' + k);
        console.table(res.targets[k].list);
        console.groupEnd();
      }
    });
    Object.keys(res.proportion || {}).forEach(function (k) {
      if (res.proportion[k].failed) {
        console.group('proportion, density ' + k);
        console.table(res.proportion[k].list);
        console.groupEnd();
      }
    });
  }

  return { run: run, contrast: contrast, targets: targets, proportion: proportion, rgba: rgba, ratio: ratio };
})();
