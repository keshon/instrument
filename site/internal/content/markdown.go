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
			html.WithUnsafe(), // the examples ARE HTML, that is the whole point
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
		// THE FRAGMENT IS SLUGGED BY THE SAME FUNCTION THAT MAKES THE ANCHOR.
		// A heading becomes an id through `slug`, which transliterates — so a
		// Russian heading answers to `poryadok-sloev` while the author of the
		// link writes what is printed on the page. Both existing links of this
		// kind in the corpus were broken from the day they were written, and
		// nothing said so: a fragment that names nothing lands the reader at
		// the top of the page.
		hash := ""
		if i := strings.IndexByte(dest, '#'); i >= 0 {
			hash, dest = "#"+slug(dest[i+1:]), dest[:i]
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
	section    string // id of the open section: Related links need a class of their own
	ids        map[string]int

	// Fences already eaten by the previous example as its extra targets. The
	// walk goes one node at a time while the targets are gathered forward along
	// the neighbours — without a mark the same fence would be drawn a second
	// time, on its own.
	eaten map[ast.Node]bool
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

// A Related link is a BUTTON OF THE KIT rather than a plate of the site's own.
//
// The section shows a row of transitions, and drawing them here would mean the
// site owning a height, a padding, a border and a radius that .inst-btn on an
// <a> already knows. A reference that draws its own around somebody else's
// components shows something other than what it hands out. In the other
// sections a link stays a link: there it sits inside a sentence.
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

// THE TABLE IS DATA, not text: it turns the heading of a Russian page into an
// anchor a link can point at. It leaves with the last Russian page, and the
// kind vocabulary below it leaves the same way — both spellings stand there
// side by side for as long as both kinds of page exist.
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

// The glyphs of the button, and NEITHER OF THEM NAMES A SIZE: `.inst-copy` in
// the kit decides that, and it decides it by where the button stands — beside
// a value or alone in the corner of a listing. This used to take an argument,
// and the two call sites below passed different values.
func copyIcons() string {
	return `<svg class="inst-icon copy-i" aria-hidden="true"><use href="#i-copy"/></svg>` +
		`<svg class="inst-icon copy-i-done" aria-hidden="true"><use href="#i-check"/></svg>`
}

func copyButton(text string, lg i18n.Lang) string {
	return fmt.Sprintf(
		`<button class="code-copy inst-copy inst-btn inst-btn--sm inst-btn--icon inst-btn--ghost" type="button" `+
			`data-copy="%s" data-copied-label="%s" data-failed-label="%s" `+
			`aria-label="%s" title="%s">%s</button>`,
		escape(text), escape(i18n.T(lg, "copy.done")), escape(i18n.T(lg, "copy.fail")),
		escape(i18n.T(lg, "copy")), escape(i18n.T(lg, "copy")), copyIcons())
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
		escape(i18n.T(lg, "copy")), escape(v), escape(v), copyIcons())
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

	// An extra target of an example: a fence right after the live example,
	// marked with the word target. The label comes from target=Name, otherwise
	// from the language of the fence.
	targetRe = regexp.MustCompile(`(^|\s)target(=(\S+))?(\s|$)`)
)

// codeTarget is one target of the code panel: the same component in another
// language.
type codeTarget struct {
	lang  string
	label string
	raw   string
	id    string
}

// The label of a target by the language of its fence. The list is open on
// purpose: a new target is a line here and a fence on the page, not an edit to
// the markup of the example.
//
// React and Svelte stand on the list BEFORE any adapter exists, and that is no
// promise to the reader: a tab is drawn only when the page carries a fence in
// that language. There is no such thing as an empty tab, and the place for one
// is ready.
var targetLabels = map[string]string{
	"html":   "HTML",
	"js":     "JS",
	"jsx":    "React",
	"tsx":    "React",
	"svelte": "Svelte",
	"vue":    "Vue",
}

func targetLabel(lang, info string) string {
	if m := targetRe.FindStringSubmatch(info); m != nil && m[3] != "" {
		return m[3]
	}
	if l, ok := targetLabels[lang]; ok {
		return l
	}
	return lang
}

func (r *codeRenderer) render(w util.BufWriter, source []byte, n ast.Node, entering bool) (ast.WalkStatus, error) {
	if !entering {
		return ast.WalkContinue, nil
	}
	if r.eaten[n] {
		return ast.WalkSkipChildren, nil
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
		writeIcons(w, r.page.Lang)
		return ast.WalkSkipChildren, nil
	}

	if lang == "api" {
		writeAPI(w, r.page)
		return ast.WalkSkipChildren, nil
	}

	if lang == "related" {
		writeRelated(w, r.page)
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

		// An example has no theme switch of its own. The theme of the
		// reference is one knob in the header, and an example inherits it:
		// .inst-theme takes colour and background from the tokens, and a theme
		// of the kit works on any subtree. A second switch on each of two
		// hundred examples was a longer road to the same result.
		// The frame of an example is a BORDER rather than a panel.
		//
		// A panel of the kit in this role follows from the reasoning that a
		// reference drawing its own frame around somebody else's component
		// shows something other than what it hands out. The reasoning is right
		// about the drawing and wrong about the role: a panel is a PRODUCT,
		// the kit hands it out, and a demonstration dressed in one looks like
		// a shipped thing and argues for attention with what is inside it. On
		// the button page there are thirteen such panels.
		//
		// A dashed border is not handed out by the kit — and so an example can
		// be outlined with one without passing it off as a component. The
		// tokens are still the kit's: the reference starts no language of its
		// own, it merely stops dressing a demonstration as merchandise.
		//
		// The class .demo-stage stays: a couple of dozen rules with :has()
		// hang on it — the hatching of the ground, the air, the assembled
		// screen. Removing it would mean rewriting them all along with the
		// frame.
		// The targets are gathered FORWARD along the neighbours: a fence
		// marked target standing right after a live example is the same
		// example in another language rather than a separate block of code.
		// The goldmark walk is streaming and hands over one node at a time,
		// which is why what has been eaten is marked.
		targets := []codeTarget{{lang: lang, label: targetLabel(lang, info), raw: raw, id: id + "-0"}}
		for sib, k := n.NextSibling(), 1; sib != nil; sib, k = sib.NextSibling(), k+1 {
			fence, ok := sib.(*ast.FencedCodeBlock)
			if !ok {
				break
			}
			sinfo := ""
			if fence.Info != nil {
				sinfo = string(fence.Info.Segment.Value(source))
			}
			if !targetRe.MatchString(sinfo) {
				break
			}
			var b bytes.Buffer
			fl := fence.Lines()
			for i := 0; i < fl.Len(); i++ {
				seg := fl.At(i)
				b.Write(seg.Value(source))
			}
			slang := string(fence.Language(source))
			targets = append(targets, codeTarget{
				lang: slang, label: targetLabel(slang, sinfo), raw: b.String(),
				id: fmt.Sprintf("%s-%d", id, k),
			})
			if r.eaten == nil {
				r.eaten = map[ast.Node]bool{}
			}
			r.eaten[sib] = true
			if slang == "js" || slang == "javascript" {
				r.page.HasJS = true
			}
		}

		fmt.Fprintf(w, `<figure class="demo%s" data-demo>`+
			`<figcaption class="demo-label">%s</figcaption>`+
			`<div class="demo-frame demo-stage inst-theme">`+
			`<div class="demo-root%s" data-density="comfortable">%s</div></div>`,
			hero, label, ctxClass(ctx), raw)
		writeTargets(w, targets, lg)
		w.WriteString(`</figure>`)
		return ast.WalkSkipChildren, nil
	}

	writeCode(w, raw, lang, false, r.page.Lang)
	return ast.WalkSkipChildren, nil
}

// writeTargets draws the code panel of an example: the fold-out label and, if
// there is more than one target, a tab strip of the kit above the blocks of
// code.
//
// A TAB IS DRAWN ONLY WHEN THERE IS CODE BEHIND IT. Three equal HTML/JS/React
// tabs would be a drawn promise: React does not exist, and `## JS` occurs on
// eighteen pages out of fifty — on the rest the tab would say "the same as in
// HTML" or stand empty. The kit forbids itself such promises, and the kit's
// reference is no exception.
//
// The strip is the kit's role="tablist" entire: instrument.js already fulfils
// its contract (arrows, Home and End, a roving tabindex) and switches the
// panels by aria-controls itself. A strip of our own would mean a third
// implementation of a switch in a repository that already has two proven
// ones.
func writeTargets(w util.BufWriter, targets []codeTarget, lg i18n.Lang) {
	head := i18n.T(lg, "demo.markup")
	if len(targets) > 1 {
		head = i18n.T(lg, "demo.code")
	}
	fmt.Fprintf(w, `<details class="inst-accordion-item demo-code">`+
		`<summary class="inst-accordion-head">%s</summary>`, head)

	if len(targets) > 1 {
		fmt.Fprintf(w, `<div class="inst-tabs demo-targets" role="tablist" aria-label="%s">`,
			i18n.T(lg, "demo.code"))
		for i, t := range targets {
			sel, tab := "false", "-1"
			if i == 0 {
				sel, tab = "true", "0"
			}
			fmt.Fprintf(w, `<button class="inst-tab" type="button" role="tab" id="%s-tab"`+
				` aria-selected="%s" aria-controls="%s" tabindex="%s">%s</button>`,
				t.id, sel, t.id, tab, escape(t.label))
		}
		w.WriteString(`</div>`)
	}

	for i, t := range targets {
		attrs := ""
		if len(targets) > 1 {
			attrs = fmt.Sprintf(` id="%s" role="tabpanel" aria-labelledby="%s-tab" tabindex="0"`, t.id, t.id)
			if i > 0 {
				attrs += " hidden"
			}
		}
		fmt.Fprintf(w, `<div class="code-block code-block--demo"%s>%s<pre><code class="lang-%s">%s</code></pre></div>`,
			attrs, copyButton(t.raw, lg), escape(t.lang), highlight(t.raw, t.lang))
	}
	w.WriteString(`</details>`)
}

func writeCode(w util.BufWriter, raw, lang string, inDemo bool, lg i18n.Lang) {
	cls := "code-block"
	if inDemo {
		cls += " code-block--demo"
	}
	fmt.Fprintf(w, `<div class="%s">%s<pre><code class="lang-%s">%s</code></pre></div>`,
		cls, copyButton(raw, lg), escape(lang), highlight(raw, lang))
}

// writeRelated prints the neighbours of a page from the registry.
//
// The links are not typed by hand, and the reason is measurable: across four
// pages of the Actions section the lists of neighbours drifted apart from one
// another — the chip named the segmented control, the segmented did not name
// the chip. A list that has to be kept in agreement across four files falls out
// of agreement on the first edit.
//
// The address is assembled from the directory of the neighbour's page rather
// than stored: the name of a component in the registry coincides with the slug
// of its page across all sixty-eight entries, and a second mapping would be a
// second source of truth.
func writeRelated(w util.BufWriter, p *Page) {
	names := p.Related()
	if len(names) == 0 {
		return
	}
	w.WriteString(`<p class="related">`)
	for _, n := range names {
		t, ok := pageBySlug[relKey{p.Lang, n}]
		if !ok {
			continue
		}
		fmt.Fprintf(w, `<a class="inst-btn inst-btn--sm" href="%s">%s</a>`,
			escape(t.Route), escape(t.Title))
	}
	w.WriteString(`</p>`)
}

// The API used to be ONE table of three columns for three kinds of row, and
// two thirds of its cells were empty by construction: a class has no value, a
// token has no prose. Filler reads as filler because that is what it was.
//
// Three blocks instead, each with the columns its own rows actually fill.
// The split is not cosmetic — it made room for the column that was missing.
func writeAPI(w util.BufWriter, p *Page) {
	if len(p.API) == 0 {
		return
	}
	var classes, attrs, tokens []APIRow
	for _, r := range p.API {
		switch r.Kind {
		case "класс", "модификатор", "class", "modifier":
			classes = append(classes, r)
		case "атрибут", "событие", "attribute", "event":
			attrs = append(attrs, r)
		default:
			tokens = append(tokens, r)
		}
	}
	w.WriteString(`<div class="api">`)
	writeAPIBlock(w, p, "api.classes", classes, 2)
	writeAPIBlock(w, p, "api.attrs", attrs, 3)
	writeAPIBlock(w, p, "api.tokens", tokens, 4)
	w.WriteString(`</div>`)
}

// writeAPIBlock renders one block. `cols` picks the shape rather than a count:
// 2 is name and prose, 3 is name, value and prose, 4 is name, value and the
// cell the value came from.
func writeAPIBlock(w util.BufWriter, p *Page, title string, rows []APIRow, cols int) {
	if len(rows) == 0 {
		return
	}
	span := 2
	if cols > 2 {
		span = 3
	}
	fmt.Fprintf(w, `<div class="api-scroll"><table class="api-table api-table--%d"><caption>%s</caption><thead><tr>`,
		cols, escape(i18n.T(p.Lang, title)))
	switch cols {
	case 2:
		fmt.Fprintf(w, `<th>%s</th><th>%s</th>`,
			i18n.T(p.Lang, "api.name"), i18n.T(p.Lang, "api.doc"))
	case 3:
		fmt.Fprintf(w, `<th>%s</th><th>%s</th><th>%s</th>`,
			i18n.T(p.Lang, "api.name"), i18n.T(p.Lang, "api.value"), i18n.T(p.Lang, "api.doc"))
	default:
		fmt.Fprintf(w, `<th>%s</th><th>%s</th><th>%s</th>`,
			i18n.T(p.Lang, "api.name"), i18n.T(p.Lang, "api.value"), i18n.T(p.Lang, "api.where"))
	}
	w.WriteString(`</tr></thead><tbody>`)

	kind := ""
	for _, r := range rows {
		if r.Kind != kind {
			kind = r.Kind
			fmt.Fprintf(w, `<tr class="api-group"><th colspan="%d" scope="colgroup">%s</th></tr>`,
				span, escape(i18n.Kind(p.Lang, kind)))
		}
		name := copyValue(r.Name, p.Lang) + `<code class="api-name">` + escape(r.Name) + `</code>`
		if cols == 4 && r.Doc != "" {
			// The prose about a token is the exception — a hundred and twenty of
			// seven hundred and fifty seven — so it goes UNDER the name rather
			// than into a column that would stand empty five rows out of six.
			name += `<small class="api-note">` + inlineCode(r.Doc) + `</small>`
		}
		fmt.Fprintf(w, `<tr data-kind="%s"><td class="api-cell">%s</td>`, escape(r.Kind), name)
		if cols > 2 {
			val := escape(r.Value)
			if val == "" {
				val = `<span class="api-none">—</span>`
			} else {
				val = copyValue(r.Value, p.Lang) + `<code>` + val + `</code>`
			}
			fmt.Fprintf(w, `<td class="api-cell">%s</td>`, val)
		}
		if cols == 4 {
			fmt.Fprintf(w, `<td class="api-where">%s</td>`, whereFrom(p, r))
		} else {
			fmt.Fprintf(w, `<td>%s</td>`, inlineCode(r.Doc))
		}
		w.WriteString(`</tr>`)
	}
	w.WriteString(`</tbody></table></div>`)
}

// whereFrom says which cell the printed value was taken from.
//
// This column is the point of the whole split. --control-h-sm printed as 26px
// is TRUE and means less than it looks: the kit declares it fifteen times, and
// at compact density it is 22. Fifty-six of the hundred and forty-three tokens
// the pages name are like that. A number without its cell is the same defect as
// "coverage 100%" counted over the wrong set.
func whereFrom(p *Page, r APIRow) string {
	if r.Kind == "переменная" || r.Kind == "variable" {
		return escape(i18n.T(p.Lang, "api.set"))
	}
	if len(r.Cells) <= 1 {
		return escape(i18n.T(p.Lang, "api.one"))
	}
	seen := map[string]bool{}
	var axes []string
	for _, c := range r.Cells[1:] {
		for _, ax := range axesOf(c) {
			if !seen[ax] {
				seen[ax] = true
				axes = append(axes, escape(i18n.T(p.Lang, "api.axis."+ax)))
			}
		}
	}
	return fmt.Sprintf(`%s <span class="api-note">+%d: %s</span>`,
		escape(i18n.T(p.Lang, "api.base")), len(r.Cells)-1, strings.Join(axes, " · "))
}

// axesOf names the axes a selector moves along. A cell can name two at once —
// the ten scale-and-density combinations do — so this returns a list.
func axesOf(sel string) []string {
	var out []string
	for _, a := range []struct{ mark, name string }{
		{"data-scale", "scale"},
		{"data-density", "density"},
		{"data-theme", "theme"},
		{"data-accent", "accent"},
		{"@media", "screen"},
	} {
		if strings.Contains(sel, a.mark) {
			out = append(out, a.name)
		}
	}
	return out
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

	// A string and a keyword are caught in ONE pass rather than two.
	//
	// Two passes cannot work: whichever is painted first, the second pass
	// steps inside what is already painted and paints a `const` inside a
	// string or a class name inside a label. An alternation in one regular
	// expression parses the line left to right, and the inside of a string
	// simply does not exist for the second branch — it has been eaten by the
	// first.
	jsRe = regexp.MustCompile("(&#34;[^&\n]*&#34;|'[^'\n]*'|`[^`\n]*`)|" +
		`\b(await|async|case|catch|class|const|continue|default|delete|else|export|extends|` +
		`false|finally|for|from|function|if|import|in|instanceof|let|new|null|of|return|switch|` +
		`this|throw|true|try|typeof|undefined|var|void|while|yield)\b`)

	// A whole comment line only. A `//` in the middle of a line is most often
	// an address: in `https://` the very first pass would eat half the
	// example.
	jsLineRe = regexp.MustCompile(`(?m)^(\s*)(//.*)$`)
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
	case "js", "javascript", "json":
		s = jsRe.ReplaceAllStringFunc(s, func(m string) string {
			if strings.HasPrefix(m, "&#34;") || strings.HasPrefix(m, "'") || strings.HasPrefix(m, "`") {
				return `<b class="t-val">` + m + `</b>`
			}
			return `<b class="t-tag">` + m + `</b>`
		})
		s = jsLineRe.ReplaceAllString(s, `$1<i class="t-com">$2</i>`)
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

// The symbols of the sprite come by a regexp over the source: a second list of
// names would drift from the first, and one already did.
var symbolRe = regexp.MustCompile(`<symbol id="(i-[a-z0-9-]+)"`)

func writeIcons(w util.BufWriter, lg i18n.Lang) {
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

	// The symbols of the reference's own pages do not go into the catalogue:
	// they belong to the site, they are not part of what the kit ships, and
	// their place is not here. They are named by a number so that the reader
	// sees the sprite file is larger than what the kit hands out.
	fmt.Fprintf(w, i18n.T(lg, "icons.note"), len(pages))
}
