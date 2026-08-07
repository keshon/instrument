package content

import (
	"bytes"
	"fmt"
	"path"
	"regexp"
	"strings"

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
	cr := &codeRenderer{page: p}
	md := goldmark.New(
		goldmark.WithExtensions(extension.GFM),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
			parser.WithASTTransformers(util.Prioritized(&linkRewriter{dir: p.Dir}, 100)),
		),
		goldmark.WithRendererOptions(
			html.WithUnsafe(), // примеры — это HTML, в том и смысл
			renderer.WithNodeRenderers(util.Prioritized(cr, 1)),
		),
	)

	doc := md.Parser().Parse(text.NewReader(body))
	p.TOC = collectTOC(doc, body)

	var buf bytes.Buffer
	if err := md.Renderer().Render(&buf, body, doc); err != nil {
		return err
	}
	p.HTML = buf.String()
	return nil
}

// ── Ссылки ────────────────────────────────────────────────────────────────
//
// Страницы ссылаются друг на друга относительными путями к .md — так они
// читаются и в репозитории на GitHub. На сайте у страницы другой адрес,
// поэтому ссылка переписывается в маршрут. Точка отсчёта — каталог ФАЙЛА,
// а не текущий URL: относительный путь в markdown отсчитывается от файла.

type linkRewriter struct{ dir string }

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
		link.Destination = []byte("/" + strings.Trim(target, "/") + "/" + hash)
		return ast.WalkContinue, nil
	})
}

// ── Блоки кода ────────────────────────────────────────────────────────────

type codeRenderer struct {
	page *Page
	n    int
}

func (r *codeRenderer) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	reg.Register(ast.KindFencedCodeBlock, r.render)
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
		id := strings.Trim(strings.ReplaceAll(r.page.Route, "/", "-"), "-")
		if id == "" {
			id = "index"
		}
		id = fmt.Sprintf("%s-%d", id, r.n)
		ctx := contextRe.MatchString(info)
		r.page.Demos = append(r.page.Demos, Demo{ID: id, Markup: raw, Context: ctx})

		label := "Пример"
		if ctx {
			label = "В контексте"
		}

		// Заголовок сцены — это её адрес, а не украшение: пример живёт
		// отдельным документом, и его можно открыть в новой вкладке.
		fmt.Fprintf(w, `<figure class="demo" data-demo>`+
			`<figcaption class="demo-bar">`+
			`<span class="demo-chrome" aria-hidden="true"></span>`+
			`<span class="demo-label">`+label+`</span>`+
			`<span class="demo-tools">`+
			`<span class="inst-select-wrap demo-theme"><select class="inst-select inst-select--sm" aria-label="Тема примера" data-demo-theme>`+
			`<option value="">как у сайта</option>`+
			`<option value="light-neutral">светлая нейтральная</option><option value="light">светлая тёплая</option><option value="light-cool">светлая холодная</option>`+
			`<option value="dark-light">тёмная светло-серая</option><option value="dark-soft">тёмная серая</option><option value="dark">тёмная чёрная</option></select></span>`+
			`<a class="demo-open" href="/demo/%s.html" target="_blank" rel="noopener" title="Открыть пример отдельно">↗</a>`+
			`</span></figcaption>`+
			`<iframe class="demo-frame" src="/demo/%s.html" title="Живой пример" loading="lazy"></iframe>`,
			id, id)
		writeCode(w, raw, lang, true)
		w.WriteString(`</figure>`)
		return ast.WalkSkipChildren, nil
	}

	writeCode(w, raw, lang, false)
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
// платформы, и без JS блок остаётся доступен.
func writeCode(w util.BufWriter, raw, lang string, inDemo bool) {
	if inDemo {
		w.WriteString(`<details class="demo-code"><summary class="demo-code-head">Разметка</summary>`)
	}
	cls := "code-block"
	if inDemo {
		cls += " code-block--demo"
	}
	fmt.Fprintf(w, `<div class="%s"><button class="code-copy inst-btn inst-btn--sm" type="button" data-copy>копировать</button><pre><code class="lang-%s">%s</code></pre></div>`,
		cls, escape(lang), highlight(raw, lang))
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
	w.WriteString(`<div class="api"><table class="api-table"><thead><tr>` +
		`<th>Имя</th><th>Значение</th><th>Что делает</th>` +
		`</tr></thead><tbody>`)
	kind := ""
	for _, r := range p.API {
		if r.Kind != kind {
			kind = r.Kind
			fmt.Fprintf(w, `<tr class="api-group"><th colspan="3" scope="colgroup">%s</th></tr>`,
				escape(kind))
		}
		val := escape(r.Value)
		if val == "" {
			// Прочерк, а не пустота: пустая ячейка читается как «забыли».
			val = `<span class="api-none">—</span>`
		} else {
			val = `<code>` + val + `</code>`
		}
		fmt.Fprintf(w,
			`<tr data-kind="%s"><td><code class="api-name">%s</code></td><td>%s</td><td>%s</td></tr>`,
			escape(r.Kind), escape(r.Name), val, inlineCode(r.Doc))
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

func collectTOC(doc ast.Node, src []byte) []Heading {
	var out []Heading
	ast.Walk(doc, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}
		h, ok := n.(*ast.Heading)
		if !ok || h.Level < 2 || h.Level > 3 {
			return ast.WalkContinue, nil
		}
		id, _ := h.AttributeString("id")
		b, _ := id.([]byte)
		out = append(out, Heading{Level: h.Level, ID: string(b), Text: nodeText(h, src)})
		return ast.WalkContinue, nil
	})
	return out
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
