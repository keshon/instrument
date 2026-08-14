package content

import (
	"bytes"
	"fmt"
	"path"
	"regexp"
	"sort"
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
	buf.WriteString(`<div class="doc-hero">`)
	cr.heroOpen = true

	if err := md.Renderer().Render(&buf, body, doc); err != nil {
		return err
	}
	buf.WriteString(cr.tail())

	p.HTML = wrapTables(buf.String())
	return nil
}
func wrapTables(s string) string {
	var b strings.Builder
	b.Grow(len(s) + 256)

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

type codeRenderer struct {
	page       *Page
	n          int
	heroOpen   bool
	sectionOn  bool
	heroDemoOn bool
	section    string // id открытого раздела: ссылкам «Связанного» нужен свой класс
	ids        map[string]int
}

func (r *codeRenderer) uniq(id string) string {
	if id == "" {
		id = "s"
	}
	if r.ids == nil {
		r.ids = map[string]int{}
	}
	r.ids[id]++
	if n := r.ids[id]; n > 1 {
		return fmt.Sprintf("%s-%d", id, n)
	}
	return id
}

func (r *codeRenderer) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	reg.Register(ast.KindFencedCodeBlock, r.render)
	reg.Register(ast.KindHeading, r.heading)
	reg.Register(ast.KindLink, r.link)
}

func (r *codeRenderer) heading(w util.BufWriter, source []byte, n ast.Node, entering bool) (ast.WalkStatus, error) {
	h := n.(*ast.Heading)

	if h.Level != 2 {
		if entering {
			id := r.uniq(headingID(h, source))
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
	id, title, order := "", "", 0
	if known {
		id, title, order = def.ID, i18n.SectionTitle(r.page.Lang, def.ID), def.Order
	} else {
		id = r.uniq(headingID(h, source))
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
	r.section = id
	if known {
		w.WriteString(escape(title))
		return ast.WalkSkipChildren, nil
	}
	return ast.WalkContinue, nil
}

// Ссылка «Связанного» — КИТОВАЯ КНОПКА, а не своя плашка сайта.
//
// Раздел показывает ряд переходов, и раньше сайт рисовал их сам: высота,
// отступ, рамка, радиус — всё то, что уже умеет .inst-btn на <a>. Справочник,
// который вокруг чужих компонентов рисует свои, показывает не то, что раздаёт.
// В остальных разделах ссылка остаётся ссылкой: там она внутри предложения.
func (r *codeRenderer) link(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	n := node.(*ast.Link)
	if !entering {
		w.WriteString(`</a>`)
		return ast.WalkContinue, nil
	}
	cls := ""
	if r.section == "related" {
		cls = ` class="inst-btn inst-btn--sm"`
	}
	title := ""
	if len(n.Title) > 0 {
		title = fmt.Sprintf(` title="%s"`, escape(string(n.Title)))
	}
	fmt.Fprintf(w, `<a%s href="%s"%s>`, cls, escape(string(n.Destination)), title)
	return ast.WalkContinue, nil
}

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

func (r *codeRenderer) tail() string { return r.closeOpen() }

func headingID(h *ast.Heading, source []byte) string {
	return slug(nodeText(h, source))
}

var translit = map[rune]string{
	'а': "a", 'б': "b", 'в': "v", 'г': "g", 'д': "d", 'е': "e", 'ё': "e",
	'ж': "zh", 'з': "z", 'и': "i", 'й': "y", 'к': "k", 'л': "l", 'м': "m",
	'н': "n", 'о': "o", 'п': "p", 'р': "r", 'с': "s", 'т': "t", 'у': "u",
	'ф': "f", 'х': "h", 'ц': "c", 'ч': "ch", 'ш': "sh", 'щ': "sch",
	'ъ': "", 'ы': "y", 'ь': "", 'э': "e", 'ю': "yu", 'я': "ya",
}

func slug(s string) string {
	var b strings.Builder
	dash := false
	for _, r := range strings.ToLower(s) {
		switch {
		case r >= 'a' && r <= 'z', r >= '0' && r <= '9':
			b.WriteRune(r)
			dash = false
		case translit[r] != "":
			b.WriteString(translit[r])
			dash = false
		case r == 'ъ' || r == 'ь':
		default:
			if !dash && b.Len() > 0 {
				b.WriteByte('-')
				dash = true
			}
		}
	}
	return strings.Trim(b.String(), "-")
}

func copyIcons(mod string) string {
	return fmt.Sprintf(
		`<svg class="inst-icon%s copy-i" aria-hidden="true"><use href="#i-copy"/></svg>`+
			`<svg class="inst-icon%s copy-i-done" aria-hidden="true"><use href="#i-check"/></svg>`,
		mod, mod)
}

func copyButton(text string, lg i18n.Lang) string {
	return fmt.Sprintf(
		`<button class="code-copy inst-copy inst-btn inst-btn--sm inst-btn--icon inst-btn--ghost" type="button" `+
			`data-copy="%s" data-copied-label="%s" data-failed-label="%s" `+
			`aria-label="%s" title="%s">%s</button>`,
		escape(text), escape(i18n.T(lg, "copy.done")), escape(i18n.T(lg, "copy.fail")),
		escape(i18n.T(lg, "copy")), escape(i18n.T(lg, "copy")), copyIcons(""))
}

func copyValue(v string, lg i18n.Lang) string {
	if v == "" {
		return ""
	}
	return fmt.Sprintf(
		`<button class="api-copy inst-copy inst-btn inst-btn--sm inst-btn--icon inst-btn--ghost" type="button" `+
			`data-copy="%s" data-copied-label="%s" data-failed-label="%s" `+
			`aria-label="%s %s" title="%s">%s</button>`,
		escape(v), escape(i18n.T(lg, "copy.done")), escape(i18n.T(lg, "copy.fail")),
		escape(i18n.T(lg, "copy")), escape(v), escape(v), copyIcons(" inst-icon--sm"))
}

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

	if lang == "icons" {
		writeIcons(w)
		return ast.WalkSkipChildren, nil
	}

	if lang == "api" {
		writeAPI(w, r.page)
		return ast.WalkSkipChildren, nil
	}

	if lang == "js" || lang == "javascript" {
		r.page.HasJS = true
	}

	if lang == "html" && previewRe.MatchString(info) {
		r.n++
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

		hero := ""
		if r.heroOpen && !r.heroDemoOn {
			r.heroDemoOn = true
			r.page.Hero = true
			hero = " demo--hero"
			label = i18n.T(lg, "demo.hero")
		}

		// Своего переключателя темы у примера нет. Тема справочника — одна
		// ручка в шапке, и пример её наследует: .inst-theme берёт цвет и фон
		// из токенов, а тема кита работает на любом поддереве. Второй
		// переключатель у каждого из двухсот примеров был лишней дорогой к
		// тому же результату.
		// Рама примера собрана КИТОМ, а не своими стилями сайта. Справочник,
		// который рисует свою рамку вокруг чужого компонента, показывает не
		// то, что раздаёт: панель с шапкой, телом и раскрывающимся низом у
		// кита есть, и повод писать её заново отсутствует.
		fmt.Fprintf(w, `<figure class="inst-panel demo%s" data-demo>`+
			`<figcaption class="inst-panel-header demo-bar">`+
			`<span class="inst-panel-title">%s</span>`+
			`</figcaption>`+
			`<div class="inst-panel-body demo-stage inst-theme">`+
			`<div class="demo-root%s">%s</div></div>`,
			hero, label, ctxClass(ctx), raw)
		writeCode(w, raw, lang, true, lg)
		w.WriteString(`</figure>`)
		return ast.WalkSkipChildren, nil
	}

	writeCode(w, raw, lang, false, r.page.Lang)
	return ast.WalkSkipChildren, nil
}

func writeCode(w util.BufWriter, raw, lang string, inDemo bool, lg i18n.Lang) {
	if inDemo {
		fmt.Fprintf(w, `<details class="inst-accordion-item inst-accordion-item--flush demo-code">`+
			`<summary class="inst-accordion-head">%s</summary>`,
			i18n.T(lg, "demo.markup"))
	}
	cls := "code-block"
	if inDemo {
		cls += " code-block--demo"
	}
	fmt.Fprintf(w, `<div class="%s">%s<pre><code class="lang-%s">%s</code></pre></div>`,
		cls, copyButton(raw, lg), escape(lang), highlight(raw, lang))
	if inDemo {
		w.WriteString(`</details>`)
	}
}

func writeAPI(w util.BufWriter, p *Page) {
	if len(p.API) == 0 {
		return
	}
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
			val = `<span class="api-none">—</span>`
		} else {
			val = copyValue(r.Value, p.Lang) + `<code>` + val + `</code>`
		}
		fmt.Fprintf(w,
			`<tr data-kind="%s"><td class="api-cell">%s<code class="api-name">%s</code></td>`+
				`<td class="api-cell">%s</td><td>%s</td></tr>`,
			escape(r.Kind), copyValue(r.Name, p.Lang), escape(r.Name), val, inlineCode(r.Doc))
	}
	w.WriteString(`</tbody></table></div>`)
}

var (
	tickRe = regexp.MustCompile("`([^`]+)`")
	boldRe = regexp.MustCompile(`\*\*([^*]+)\*\*`)
)

func inlineCode(s string) string {
	out := escape(s)
	out = tickRe.ReplaceAllString(out, "<code>$1</code>")
	return boldRe.ReplaceAllString(out, "<b>$1</b>")
}

var (
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

// Символы спрайта — регуляркой по источнику: второй список имён разошёлся бы
// с первым, а он уже разошёлся.
var symbolRe = regexp.MustCompile(`<symbol id="(i-[a-z0-9-]+)"`)

func writeIcons(w util.BufWriter) {
	var ui, pages []string
	for _, m := range symbolRe.FindAllStringSubmatch(sprite, -1) {
		if strings.HasPrefix(m[1], "i-p-") {
			pages = append(pages, m[1])
		} else {
			ui = append(ui, m[1])
		}
	}
	sort.Strings(ui)
	sort.Strings(pages)

	w.WriteString(`<div class="icon-gallery">`)
	for _, id := range ui {
		fmt.Fprintf(w,
			`<span class="icon-cell"><svg class="inst-icon" aria-hidden="true"><use href="#%s"/></svg>`+
				`<code>%s</code></span>`, escape(id), escape(id))
	}
	w.WriteString(`</div>`)

	fmt.Fprintf(w, `<p class="icon-note">Плюс %d символов страниц справочника `+
		`(<code>i-p-*</code>) — они принадлежат сайту, а не киту.</p>`, len(pages))
}
