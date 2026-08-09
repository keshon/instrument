package content

import (
	"bytes"
	"fmt"
	"path"
	"regexp"
	"strings"

	"instrument/site/internal/i18n"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/renderer/html"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
)

func renderMarkdown(p *Page, body []byte) error {
	body, err := expandDirectives(body, p.Lang)
	if err != nil {
		return err
	}

	cr := &codeRenderer{page: p}
	md := goldmark.New(
		goldmark.WithExtensions(extension.GFM),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
			parser.WithASTTransformers(util.Prioritized(&linkRewriter{dir: p.Dir, prefix: p.Lang.Prefix()}, 100)),
		),
		goldmark.WithRendererOptions(
			html.WithUnsafe(), // примеры — это HTML, в том и смысл
			renderer.WithNodeRenderers(util.Prioritized(cr, 1)),
		),
	)

	doc := md.Parser().Parse(text.NewReader(body))

	var buf bytes.Buffer
	// Шапка страницы — всё до первого раздела: определение в одну-две фразы и
	// живой пример, показывающий характер компонента. Она открывается здесь, а
	// закрывается первым же H2: иначе пришлось бы искать её конец обходом
	// дерева ради того, что и так известно.
	buf.WriteString(`<div class="doc-hero">`)
	cr.heroOpen = true

	if err := md.Renderer().Render(&buf, body, doc); err != nil {
		return err
	}
	buf.WriteString(cr.tail())

	p.HTML = wrapTables(buf.String())
	return nil
}

// wrapTables даёт каждой таблице свой горизонтальный скролл.
//
// Жалоба «колонки не вмещают текст» — про это. Таблица из четырёх колонок, в
// каждой из которых предложение, на колонке в 68 знаков не помещается никак;
// без обёртки она распирала документ и заводила горизонтальную прокрутку
// ВСЕЙ странице — то есть уезжала не только таблица, а шапка, оглавление и
// навигация вместе с ней.
//
// Подстановкой по строке, а не своим рендерером узла: таблиц внутри таблиц в
// документации нет и не будет, а рендерер GFM-таблицы пришлось бы переписать
// целиком ради двух тегов вокруг.
//
// Обёртка ставится ТОЛЬКО таблице прозы. Справочник печатает свою
// (`<table class="api-table">`) и уже лежит в `.api` со своим скроллом.
// Отсюда сканер вместо замены по двум строкам: закрывающий тег у обеих таблиц
// одинаков, и подстановка добавила бы справочнику лишний `</div>`. Он закрыл
// бы раздел раньше времени, следующий `</section>` закрыл бы `<article>`, и
// оглавление ушло бы из своей колонки сетки в конец документа.
func wrapTables(s string) string {
	var b strings.Builder
	b.Grow(len(s) + 256)

	// wrapped хранит по одному флагу на каждую открытую таблицу: обёрнута она
	// или нет. Стек, а не булево, — на случай таблицы в ячейке; сегодня такой
	// нет, но молчаливо сломаться на ней эта функция уже не сможет.
	var wrapped []bool
	for {
		i := strings.Index(s, "<table")
		j := strings.Index(s, "</table>")
		switch {
		case i < 0 && j < 0:
			b.WriteString(s)
			return b.String()
		case i >= 0 && (j < 0 || i < j):
			b.WriteString(s[:i])
			bare := strings.HasPrefix(s[i:], "<table>")
			if bare {
				b.WriteString(`<div class="table-wrap">`)
			}
			wrapped = append(wrapped, bare)
			b.WriteString("<table")
			s = s[i+len("<table"):]
		default:
			b.WriteString(s[:j])
			b.WriteString("</table>")
			s = s[j+len("</table>"):]
			if n := len(wrapped); n > 0 {
				if wrapped[n-1] {
					b.WriteString("</div>")
				}
				wrapped = wrapped[:n-1]
			}
		}
	}
}

// ── Ссылки ────────────────────────────────────────────────────────────────
//
// Страницы ссылаются друг на друга относительными путями к .md — так они
// читаются и в репозитории на GitHub. На сайте у страницы другой адрес,
// поэтому ссылка переписывается в маршрут. Точка отсчёта — каталог ФАЙЛА,
// а не текущий URL: относительный путь в markdown отсчитывается от файла.

// prefix — префикс языка. Без него ссылка со страницы /en/ уводила бы
// читателя обратно в русскую версию, и перевод превращался бы в остров.
type linkRewriter struct{ dir, prefix string }

func (l *linkRewriter) Transform(doc *ast.Document, r text.Reader, pc parser.Context) {
	ast.Walk(doc, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}
		link, ok := n.(*ast.Link)
		if !ok {
			return ast.WalkContinue, nil
		}
		dest := string(link.Destination)
		if strings.Contains(dest, "://") || strings.HasPrefix(dest, "/") || strings.HasPrefix(dest, "#") {
			return ast.WalkContinue, nil
		}
		hash := ""
		if i := strings.IndexByte(dest, '#'); i >= 0 {
			dest, hash = dest[:i], dest[i:]
		}
		if !strings.HasSuffix(dest, ".md") {
			return ast.WalkContinue, nil
		}
		target := path.Clean(path.Join(l.dir, dest))
		target = strings.TrimSuffix(target, ".md")
		target = strings.TrimSuffix(target, "/index")
		link.Destination = []byte(l.prefix + "/" + strings.Trim(target, "/") + "/" + hash)
		return ast.WalkContinue, nil
	})
}

// ── Блоки кода ────────────────────────────────────────────────────────────

type codeRenderer struct {
	page *Page
	n    int

	// Состояние каркаса страницы: открыта шапка, открыт раздел, показан ли уже
	// главный пример. Один рендерер на документ, поэтому состояние здесь, а не
	// в обходе дерева.
	heroOpen   bool
	sectionOn  bool
	heroDemoOn bool
}

func (r *codeRenderer) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	reg.Register(ast.KindFencedCodeBlock, r.render)
	reg.Register(ast.KindHeading, r.heading)
}

// ── Разделы ───────────────────────────────────────────────────────────────
//
// Каждый H2 открывает <section> с устойчивым идентификатором из словаря.
// Идентификатор больше не выводится из текста заголовка, и это чинит две
// вещи разом: якорь `#variants` одинаков в русской и английской версии, а
// оглавление получает не «список заголовков», а список РАЗДЕЛОВ — то есть
// узнаваемую с первого взгляда поверхность компонента.
//
// Заголовок незнакомого вида остаётся как есть, с автоидентификатором: на
// странице вне контракта он не ошибка, а обычная проза. Ошибкой он становится
// только там, где контракт объявлен, — это проверяет check.

func (r *codeRenderer) heading(w util.BufWriter, source []byte, n ast.Node, entering bool) (ast.WalkStatus, error) {
	h := n.(*ast.Heading)

	if h.Level != 2 {
		if entering {
			id := headingID(h)
			if h.Level == 3 {
				r.page.TOC = append(r.page.TOC, Heading{Level: 3, ID: id, Text: nodeText(h, source)})
			}
			fmt.Fprintf(w, `<h%d id="%s">`, h.Level, escape(id))
		} else {
			fmt.Fprintf(w, `</h%d>`, h.Level)
		}
		return ast.WalkContinue, nil
	}

	if !entering {
		w.WriteString(`</h2>`)
		return ast.WalkContinue, nil
	}

	w.WriteString(r.closeOpen())

	raw := nodeText(h, source)
	def, known := LookupSection(raw)
	id, title, order := headingID(h), "", 0
	if known {
		id, title, order = def.ID, i18n.SectionTitle(r.page.Lang, def.ID), def.Order
	}

	if title == "" {
		title = raw
	}
	r.page.Sections = append(r.page.Sections, Section{
		ID: id, Title: title, Order: order, Known: known,
	})
	r.page.TOC = append(r.page.TOC, Heading{Level: 2, ID: id, Text: title})

	fmt.Fprintf(w, `<section class="doc-s" id="%s" data-section="%s"><h2 class="doc-s-title">`,
		escape(id), escape(id))
	r.sectionOn = true

	// У знакомого раздела заголовок берётся из словаря, а не из страницы:
	// иначе «Справочник» на одной странице и «API» на другой остались бы
	// разными разделами на вид, хотя это один и тот же раздел.
	if known {
		w.WriteString(escape(title))
		return ast.WalkSkipChildren, nil
	}
	return ast.WalkContinue, nil
}

// closeOpen закрывает то, что открыто: шапку или предыдущий раздел.
func (r *codeRenderer) closeOpen() string {
	switch {
	case r.heroOpen:
		r.heroOpen = false
		return `</div>`
	case r.sectionOn:
		r.sectionOn = false
		return `</section>`
	}
	return ""
}

// tail закрывает каркас после обхода дерева.
func (r *codeRenderer) tail() string { return r.closeOpen() }

func headingID(h *ast.Heading) string {
	id, _ := h.AttributeString("id")
	b, _ := id.([]byte)
	return string(b)
}

// copyButton — иконочная кнопка копирования.
//
// Иконка, а не слово: кнопка стоит в углу блока кода, и подпись «копировать»
// занимает там место, которого нет, — на узком экране она наезжает на первую
// строку. Иконочная кнопка у кита для этого и заведена, вместе с требованием
// доступного имени.
//
// Результат сообщается ДВАЖДЫ: значок меняется на галку для тех, кто смотрит,
// и слово уходит в живую область для тех, кто слушает. Смена значка для
// скринридера — событие без содержания, и это записано на странице «Код».
func copyButton(lg i18n.Lang) string {
	return fmt.Sprintf(
		`<button class="code-copy inst-btn inst-btn--sm inst-btn--icon" type="button" `+
			`data-copy aria-label="%s" title="%s">`+
			`<svg class="inst-icon" aria-hidden="true" data-copy-icon><use href="#i-copy"/></svg>`+
			`<span class="inst-u-visually-hidden" data-copy-msg aria-live="polite"></span>`+
			`</button>`,
		escape(i18n.T(lg, "copy")), escape(i18n.T(lg, "copy")))
}

// copyValue — кнопка копирования одного значения.
//
// Тот же обработчик, что у блока кода, но источник берётся из `data-value`, а
// не из соседнего `<code>`: в строке справочника копируют имя, а не ячейку.
func copyValue(v string, lg i18n.Lang) string {
	if v == "" {
		return ""
	}
	return fmt.Sprintf(
		`<button class="api-copy inst-btn inst-btn--sm inst-btn--icon inst-btn--ghost" type="button" `+
			`data-copy data-value="%s" aria-label="%s %s" title="%s">`+
			`<svg class="inst-icon inst-icon--sm" aria-hidden="true" data-copy-icon><use href="#i-copy"/></svg>`+
			`<span class="inst-u-visually-hidden" data-copy-msg aria-live="polite"></span>`+
			`</button>`,
		escape(v), escape(i18n.T(lg, "copy")), escape(v), escape(v))
}

// ctxClass отмечает сцену контекстного примера: у собранного экрана свои
// правила поля — оболочка занимает сцену целиком, а не лежит в ней с отступом.
func ctxClass(ctx bool) string {
	if ctx {
		return " demo-root--context"
	}
	return ""
}

var (
	previewRe = regexp.MustCompile(`(^|\s)preview(\s|$)`)
	contextRe = regexp.MustCompile(`(^|\s)context(\s|$)`)
)

func (r *codeRenderer) render(w util.BufWriter, source []byte, n ast.Node, entering bool) (ast.WalkStatus, error) {
	if !entering {
		return ast.WalkContinue, nil
	}
	node := n.(*ast.FencedCodeBlock)

	var code bytes.Buffer
	lines := node.Lines()
	for i := 0; i < lines.Len(); i++ {
		seg := lines.At(i)
		code.Write(seg.Value(source))
	}
	raw := code.String()

	lang := string(node.Language(source))
	info := ""
	if node.Info != nil {
		info = string(node.Info.Segment.Value(source))
	}

	// Ограда ```api разворачивается в справочник из frontmatter. Место в
	// потоке задаёт автор — оно разное на разных страницах, — а колонки и
	// порядок задаёт шаблон, и разойтись они не могут.
	if lang == "api" {
		writeAPI(w, r.page)
		return ast.WalkSkipChildren, nil
	}

	if lang == "html" && previewRe.MatchString(info) {
		r.n++
		// ID берётся из ПУТИ ФАЙЛА, а не из маршрута: маршрут несёт префикс
		// языка, и столы размножились бы по числу языков — 288 документов
		// вместо 144, причём с одинаковым содержимым.
		id := strings.TrimSuffix(r.page.Rel, ".md")
		id = strings.Trim(strings.ReplaceAll(id, "/", "-"), "-")
		if id == "" || id == "index" {
			id = "index"
		}
		id = fmt.Sprintf("%s-%d", id, r.n)
		ctx := contextRe.MatchString(info)
		r.page.Demos = append(r.page.Demos, Demo{ID: id, Markup: raw, Context: ctx})

		lg := r.page.Lang
		label := i18n.T(lg, "demo")
		if ctx {
			label = i18n.T(lg, "demo.context")
		}

		// Главный пример — тот, что стоит в шапке, до первого раздела. Он не
		// документирует, он показывает характер компонента: человек, пришедший
		// со страницы daisyUI, должен за одну секунду понять, на что смотрит,
		// и только потом начать читать. Отличается он не подписью, а размером
		// сцены и тем, что стоит вплотную к определению.
		hero := ""
		if r.heroOpen && !r.heroDemoOn {
			r.heroDemoOn = true
			r.page.Hero = true
			hero = " demo--hero"
			label = i18n.T(lg, "demo.hero")
		}

		// Пример рисуется В ПОТОКЕ СТРАНИЦЫ, а не в отдельном документе внутри
		// iframe.
		//
		// Кадр появился ради одного: у примера своя тема, и переключение темы
		// справочника не должно перекрашивать демонстрируемый кит. Оказалось,
		// что кадр для этого не нужен — тема в ките объявлена атрибутом
		// (`[data-theme]`), а не селектором корня, и работает на любом
		// поддереве. Не хватало только `.inst-theme`, восстанавливающей базовый
		// цвет и фон; она заведена в самом ките, потому что тёмная панель в
		// светлом приложении нужна не только справочнику.
		//
		// Кадр стоил дорого: 184 отдельных документа и три мегабайта, девять
		// загрузок на странице кнопки, синхронизация высоты через postMessage,
		// разметка примеров вне поиска и вне печати.
		//
		// data-density на сцене намеренно НЕ ставится по умолчанию: плотность
		// глобальна, и пример обязан показывать ту, которую выбрал читатель.
		fmt.Fprintf(w, `<figure class="demo%s" data-demo>`+
			`<figcaption class="demo-bar">`+
			`<span class="demo-chrome" aria-hidden="true"></span>`+
			`<span class="demo-label">%s</span>`+
			`<span class="demo-tools">`+
			`<span class="inst-select-wrap demo-theme"><select class="inst-select inst-select--sm" aria-label="%s" data-demo-theme>`+
			`<option value="">%s</option>`+
			`<option value="light-neutral">%s</option><option value="light">%s</option><option value="light-cool">%s</option>`+
			`<option value="dark-light">%s</option><option value="dark-soft">%s</option><option value="dark">%s</option></select></span>`+
			`</span></figcaption>`+
			`<div class="demo-stage inst-theme" data-demo-stage>`+
			`<div class="demo-root%s">%s</div></div>`,
			hero, label, i18n.T(lg, "demo.theme"), i18n.T(lg, "demo.same"),
			i18n.T(lg, "theme.ln"), i18n.T(lg, "theme.l"), i18n.T(lg, "theme.lc"),
			i18n.T(lg, "theme.dl"), i18n.T(lg, "theme.ds"), i18n.T(lg, "theme.d"),
			ctxClass(ctx), raw)
		writeCode(w, raw, lang, true, lg)
		w.WriteString(`</figure>`)
		return ast.WalkSkipChildren, nil
	}

	writeCode(w, raw, lang, false, r.page.Lang)
	return ast.WalkSkipChildren, nil
}

// writeCode печатает блок кода. У примера он СВЁРНУТ.
//
// Разметка примера — это подтверждение, а не содержание: сначала смотрят, что
// компонент делает, и только потом, из чего он собран. Развёрнутым кодом
// страница получала в среднем 31 лишнюю строку вертикального шума между
// читателем и справочником, ради которого сюда чаще всего и приходят.
//
// Носитель — <details>, поэтому раскрытие, клавиатура и роль достаются от
// платформы, и без JS блок остаётся доступен. Поиск по странице находит текст
// внутри закрытого блока и раскрывает его сам — то же соображение, по которому
// на <details> построен аккордеон.
func writeCode(w util.BufWriter, raw, lang string, inDemo bool, lg i18n.Lang) {
	if inDemo {
		fmt.Fprintf(w, `<details class="demo-code"><summary class="demo-code-head">%s</summary>`,
			i18n.T(lg, "demo.markup"))
	}
	cls := "code-block"
	if inDemo {
		cls += " code-block--demo"
	}
	fmt.Fprintf(w, `<div class="%s">%s<pre><code class="lang-%s">%s</code></pre></div>`,
		cls, copyButton(lg), escape(lang), highlight(raw, lang))
	if inDemo {
		w.WriteString(`</details>`)
	}
}

// writeAPI печатает справочник одной таблицей с постоянными колонками.
//
// Колонок ровно четыре, и они не меняются от страницы к странице: имя, вид,
// значение, работа. Именно постоянство колонок делает таблицу
// просматриваемой — глаз учится читать её один раз, а не заново на каждом
// компоненте.
func writeAPI(w util.BufWriter, p *Page) {
	if len(p.API) == 0 {
		return
	}
	// Вид вынесен в строку-заголовок, а не повторяется колонкой. Колонка
	// «Вид» печатала одно и то же слово по семь раз подряд — это ровно та
	// трудночитаемость, ради устранения которой таблицу и объединяли.
	// Строка-заголовок даёт то же знание один раз и заодно делит длинный
	// список на куски, по которым можно прыгать глазом.
	fmt.Fprintf(w, `<div class="api"><table class="api-table"><thead><tr>`+
		`<th>%s</th><th>%s</th><th>%s</th>`+
		`</tr></thead><tbody>`,
		i18n.T(p.Lang, "api.name"), i18n.T(p.Lang, "api.value"), i18n.T(p.Lang, "api.doc"))
	kind := ""
	for _, r := range p.API {
		if r.Kind != kind {
			kind = r.Kind
			fmt.Fprintf(w, `<tr class="api-group"><th colspan="3" scope="colgroup">%s</th></tr>`,
				escape(i18n.Kind(p.Lang, kind)))
		}
		val := escape(r.Value)
		if val == "" {
			// Прочерк, а не пустота: пустая ячейка читается как «забыли».
			val = `<span class="api-none">—</span>`
		} else {
			val = copyValue(r.Value, p.Lang) + `<code>` + val + `</code>`
		}
		// Имя копируется по кнопке, а не выделением мышью.
		//
		// За справочником приходят за именем, и оно почти всегда уезжает в
		// разметку или в CSS. `inst-dialog-foot--end` и `--control-h-sm`
		// набирают по памяти с ошибкой, а выделяют мышью — вместе с пробелом
		// или половиной соседней ячейки.
		//
		// Значение копируется отдельной кнопкой, потому что копируют его тоже
		// отдельно: `min(34rem, 100vw - var(--space-8))` руками не наберёт
		// никто.
		fmt.Fprintf(w,
			`<tr data-kind="%s"><td class="api-cell">%s<code class="api-name">%s</code></td>`+
				`<td class="api-cell">%s</td><td>%s</td></tr>`,
			escape(r.Kind), copyValue(r.Name, p.Lang), escape(r.Name), val, inlineCode(r.Doc))
	}
	w.WriteString(`</tbody></table></div>`)
}

// inlineCode — минимальный markdown внутри ячейки: `код` и **важное**.
// Полный парсер здесь был бы избыточен, а совсем без него в описании нельзя
// назвать ни класс, ни значение.
var (
	tickRe = regexp.MustCompile("`([^`]+)`")
	boldRe = regexp.MustCompile(`\*\*([^*]+)\*\*`)
)

func inlineCode(s string) string {
	out := escape(s)
	out = tickRe.ReplaceAllString(out, "<code>$1</code>")
	return boldRe.ReplaceAllString(out, "<b>$1</b>")
}

// ── Подсветка ─────────────────────────────────────────────────────────────
//
// Своя, на три цвета из палитры кита, а не библиотека: подсветка в шестьдесят
// цветов на макете из хайрлайнов и воздуха кричит громче кода, который
// подсвечивает. Различать надо ровно три вещи: имя тега, имя атрибута и
// значение.

var (
	// Открывающая часть и имя тега захватываются ОТДЕЛЬНО. Иначе поиск имени
	// «по первой букве» попадает внутрь мнемоники: в «&lt;button» первая
	// буква — это «l» из «lt», и подсветка разрывала саму мнемонику.
	tagRe     = regexp.MustCompile(`(&lt;/?)([a-zA-Z][\w-]*)`)
	attrRe    = regexp.MustCompile(`([a-zA-Z-]+)=(&#34;[^&]*&#34;)`)
	commentRe = regexp.MustCompile(`(?s)&lt;!--.*?--&gt;|/\*.*?\*/`)
	cssPropRe = regexp.MustCompile(`(?m)^(\s*)([a-z-]+):`)
	cssVarRe  = regexp.MustCompile(`(--[a-z][\w-]*)`)
)

func highlight(raw, lang string) string {
	s := escape(raw)
	switch lang {
	case "html":
		s = tagRe.ReplaceAllString(s, `$1<b class="t-tag">$2</b>`)
		s = attrRe.ReplaceAllString(s, `<b class="t-attr">$1</b>=<b class="t-val">$2</b>`)
	case "css":
		s = cssPropRe.ReplaceAllString(s, `$1<b class="t-attr">$2</b>:`)
		s = cssVarRe.ReplaceAllString(s, `<b class="t-val">$1</b>`)
	}
	return commentRe.ReplaceAllString(s, `<i class="t-com">$0</i>`)
}

func escape(s string) string {
	var b strings.Builder
	for _, r := range s {
		switch r {
		case '&':
			b.WriteString("&amp;")
		case '<':
			b.WriteString("&lt;")
		case '>':
			b.WriteString("&gt;")
		case '"':
			b.WriteString("&#34;")
		default:
			b.WriteRune(r)
		}
	}
	return b.String()
}

// ── Оглавление ────────────────────────────────────────────────────────────
//
// Собирается прямо во время разметки, а не отдельным обходом: у знакомого
// раздела заголовок берётся из словаря, и второй обход всё равно не знал бы,
// что напечатано на странице.

func nodeText(n ast.Node, src []byte) string {
	var b strings.Builder
	ast.Walk(n, func(c ast.Node, entering bool) (ast.WalkStatus, error) {
		if entering {
			if t, ok := c.(*ast.Text); ok {
				b.Write(t.Segment.Value(src))
			}
		}
		return ast.WalkContinue, nil
	})
	return strings.TrimSpace(b.String())
}
