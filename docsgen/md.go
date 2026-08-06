package main

import (
	"bytes"
	"fmt"
	"os"
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
	"gopkg.in/yaml.v3"
)

type frontmatter struct {
	Title    string `yaml:"title"`
	Group    string `yaml:"group"`
	Source   string `yaml:"source"`
	Status   string `yaml:"status"`
	NeedsJS  string `yaml:"needs-js"`
	Template string `yaml:"template"`
}

var fmRe = regexp.MustCompile(`(?s)\A---\r?\n(.*?)\r?\n---\r?\n`)

func parsePage(fsPath, rel string) (*Page, error) {
	raw, err := os.ReadFile(fsPath)
	if err != nil {
		return nil, err
	}

	var meta frontmatter
	body := raw
	if m := fmRe.FindSubmatch(raw); m != nil {
		if err := yaml.Unmarshal(m[1], &meta); err != nil {
			return nil, fmt.Errorf("frontmatter: %w", err)
		}
		body = raw[len(m[0]):]
	}

	slug := strings.TrimSuffix(path.Base(rel), ".md")
	dir := path.Dir(rel)
	if dir == "." {
		dir = ""
	}
	route := "/"
	if slug != "index" {
		route = "/" + path.Join(dir, slug) + "/"
	} else if dir != "" {
		route = "/" + dir + "/"
	}

	p := &Page{
		Route: route, Rel: rel, Dir: dir, Slug: slug,
		Title: meta.Title, Group: meta.Group, Source: meta.Source,
		NeedsJS: meta.NeedsJS, Status: meta.Status,
		Splash: meta.Template == "splash",
	}
	if p.Title == "" {
		p.Title = slug
	}

	md := newMarkdown(dir)
	doc := md.Parser().Parse(text.NewReader(body))
	p.TOC = collectTOC(doc, body)

	var buf bytes.Buffer
	if err := md.Renderer().Render(&buf, body, doc); err != nil {
		return nil, err
	}
	p.HTML = buf.String()
	p.Text = plainText(body)
	return p, nil
}

func newMarkdown(dir string) goldmark.Markdown {
	return goldmark.New(
		goldmark.WithExtensions(extension.GFM),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
			parser.WithASTTransformers(util.Prioritized(&linkRewriter{dir: dir}, 100)),
		),
		goldmark.WithRendererOptions(
			html.WithUnsafe(), // примеры — это HTML, в том и смысл
			renderer.WithNodeRenderers(util.Prioritized(&codeRenderer{}, 1)),
		),
	)
}

// ── Ссылки ────────────────────────────────────────────────────────────────
//
// Страницы ссылаются друг на друга относительными путями к .md — так они
// читаются и в репозитории на GitHub. На сайте у страницы другой адрес,
// поэтому ссылка переписывается в маршрут. Точка отсчёта — каталог ФАЙЛА,
// а не текущий URL: относительный путь в markdown отсчитывается от файла.

type linkRewriter struct{ dir string }

func (l *linkRewriter) Transform(doc *ast.Document, reader text.Reader, pc parser.Context) {
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
//
// Ограда ```html preview разворачивается в живой пример: та же разметка
// сначала рисуется настоящим китом, потом показывается кодом. Разойтись
// они не могут — это одна и та же строка.
//
// Сцена намеренно БЕЗ рамки. Панель и карточка приносят свою, и рамка
// вокруг рамки — ровно то, что кит запрещает у ряда метрик.

type codeRenderer struct{}

func (r *codeRenderer) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	reg.Register(ast.KindFencedCodeBlock, r.renderFenced)
}

func (r *codeRenderer) renderFenced(w util.BufWriter, source []byte, n ast.Node, entering bool) (ast.WalkStatus, error) {
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
	preview := lang == "html" && regexp.MustCompile(`(^|\s)preview(\s|$)`).MatchString(info)

	if preview {
		fmt.Fprintf(w, `<div class="demo"><div class="demo-stage">%s</div>`, raw)
		writeCode(w, raw, lang, true)
		w.WriteString(`</div>`)
		return ast.WalkSkipChildren, nil
	}
	writeCode(w, raw, lang, false)
	return ast.WalkSkipChildren, nil
}

func writeCode(w util.BufWriter, raw, lang string, inDemo bool) {
	cls := "code-block"
	if inDemo {
		cls += " code-block--demo"
	}
	fmt.Fprintf(w, `<div class="%s"><button class="code-copy inst-btn inst-btn--sm" type="button" data-copy>копировать</button><pre><code class="lang-%s">%s</code></pre></div>`,
		cls, htmlEscape(lang), highlight(raw, lang))
}

// ── Подсветка ─────────────────────────────────────────────────────────────
//
// Своя, на три цвета из палитры кита, а не библиотека: подсветка в 60 цветов
// на макете из хайрлайнов и воздуха кричит громче кода, который подсвечивает.
// Различать надо ровно три вещи: имя тега, имя атрибута и значение.

var (
	// Захватываются ОТДЕЛЬНО открывающая часть и имя тега. Иначе поиск
	// первой буквы попадает внутрь мнемоники: в «&lt;button» первая буква —
	// это «l» из «lt», и подсветка разрывала саму мнемонику на «&» и
	// «lt;button», после чего в коде проступали сущности.
	tagRe     = regexp.MustCompile(`(&lt;/?)([a-zA-Z][\w-]*)`)
	attrRe    = regexp.MustCompile(`([a-zA-Z-]+)=(&#34;[^&]*&#34;)`)
	commentRe = regexp.MustCompile(`(?s)&lt;!--.*?--&gt;|/\*.*?\*/`)
	cssPropRe = regexp.MustCompile(`(?m)^(\s*)([a-z-]+):`)
	cssVarRe  = regexp.MustCompile(`(--[a-z][\w-]*)`)
)

func highlight(raw, lang string) string {
	s := htmlEscape(raw)
	switch lang {
	case "html":
		s = tagRe.ReplaceAllString(s, `$1<b class="t-tag">$2</b>`)
		s = attrRe.ReplaceAllString(s, `<b class="t-attr">$1</b>=<b class="t-val">$2</b>`)
	case "css":
		s = cssPropRe.ReplaceAllString(s, `$1<b class="t-attr">$2</b>:`)
		s = cssVarRe.ReplaceAllString(s, `<b class="t-val">$1</b>`)
	}
	s = commentRe.ReplaceAllString(s, `<i class="t-com">$0</i>`)
	return s
}

func htmlEscape(s string) string {
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

// ── Оглавление и текст для поиска ─────────────────────────────────────────

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
		idStr, _ := id.([]byte)
		out = append(out, Heading{Level: h.Level, ID: string(idStr), Text: nodeText(h, src)})
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

var stripRe = regexp.MustCompile("(?s)```.*?```|`[^`]*`|[|*_#>\\[\\]()]")

func plainText(body []byte) string {
	s := stripRe.ReplaceAllString(string(body), " ")
	return strings.Join(strings.Fields(s), " ")
}
