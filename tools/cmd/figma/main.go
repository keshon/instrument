// Команда figma читает дамп «Copy as CSS» из Figma и отвечает на вопросы о нём.
//
// Дамп — это плоский список узлов макета: имя в комментарии, следом объявления.
// Вложенность в него не попадает, порядок — обход дерева слоёв. Файл на четверть
// миллиона знаков, и читать его целиком незачем: нужны выборки.
//
// Пятая проверка в tools/ он НЕ является и в CI не входит. Это читалка чужого
// макета, а не сверка кита с самим собой: она ничего не утверждает про
// instrument и потому не может упасть.
package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"regexp"
	"sort"
	"strings"
)

// Пометки, которые Figma вставляет тем же синтаксисом, что и имена узлов.
// Без этого списка «Auto layout» становится компонентом, встречающимся 402 раза.
var notNames = map[string]bool{
	"Auto layout":        true,
	"Inside auto layout": true,
}

var notNamePrefix = []string{
	"identical to box height",
	"or ",
}

// node — один узел макета: имя, объявления в порядке дампа, номер по порядку.
type node struct {
	idx   int
	name  string
	props []prop
}

type prop struct{ key, val string }

func (n node) get(k string) string {
	for _, p := range n.props {
		if p.key == k {
			return p.val
		}
	}
	return ""
}

// signature — строка, по которой узел опознаётся в списке, не разворачивая его.
func (n node) signature() string {
	var b []string
	if w, h := n.get("width"), n.get("height"); w != "" || h != "" {
		b = append(b, strings.TrimSuffix(w, "px")+"×"+strings.TrimSuffix(h, "px"))
	}
	for _, k := range []string{"background", "border", "border-radius", "font-size", "color"} {
		if v := n.get(k); v != "" {
			b = append(b, k+"="+v)
		}
	}
	return strings.Join(b, "  ")
}

var reComment = regexp.MustCompile(`^/\*\s*(.*?)\s*\*/$`)
var reDecl = regexp.MustCompile(`^([a-zA-Z-]+):\s*(.*?);?$`)

func isName(s string) bool {
	if notNames[s] {
		return false
	}
	for _, p := range notNamePrefix {
		if strings.HasPrefix(s, p) {
			return false
		}
	}
	return s != ""
}

func parse(path string) ([]node, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var nodes []node
	var cur *node
	sc := bufio.NewScanner(f)
	sc.Buffer(make([]byte, 0, 64*1024), 1024*1024)
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" {
			continue
		}
		if m := reComment.FindStringSubmatch(line); m != nil {
			if !isName(m[1]) {
				continue
			}
			if cur != nil {
				nodes = append(nodes, *cur)
			}
			cur = &node{idx: len(nodes), name: m[1]}
			continue
		}
		if cur == nil {
			continue
		}
		if m := reDecl.FindStringSubmatch(line); m != nil {
			cur.props = append(cur.props, prop{m[1], strings.TrimSuffix(m[2], ";")})
		}
	}
	if cur != nil {
		nodes = append(nodes, *cur)
	}
	return nodes, sc.Err()
}

func printNode(n node, full bool) {
	fmt.Printf("#%d  %s\n", n.idx, n.name)
	if !full {
		if s := n.signature(); s != "" {
			fmt.Printf("      %s\n", s)
		}
		return
	}
	for _, p := range n.props {
		fmt.Printf("      %s: %s\n", p.key, p.val)
	}
}

func main() {
	var (
		path  = flag.String("f", "../temp/components-dump.css", "файл дампа")
		find  = flag.String("find", "", "показать узлы, чьё имя совпало с регулярным выражением")
		at    = flag.Int("at", -1, "показать узел по номеру")
		ctx   = flag.Int("ctx", 0, "сколько соседей показать вокруг -at")
		props = flag.String("props", "", "таблица перечисленных свойств по всем узлам")
		list  = flag.Bool("list", false, "имена всех узлов с числом повторов")
		sum   = flag.Bool("sum", false, "опись: какие значения вообще встречаются")
		full  = flag.Bool("full", false, "разворачивать узлы целиком, а не подписью")
	)
	flag.Parse()

	nodes, err := parse(*path)
	if err != nil {
		fmt.Fprintln(os.Stderr, "не прочитан дамп:", err)
		os.Exit(1)
	}
	fmt.Fprintf(os.Stderr, "· %s: %d узлов\n\n", *path, len(nodes))

	switch {
	case *list:
		count := map[string]int{}
		for _, n := range nodes {
			count[n.name]++
		}
		names := make([]string, 0, len(count))
		for k := range count {
			names = append(names, k)
		}
		sort.Slice(names, func(i, j int) bool {
			if count[names[i]] != count[names[j]] {
				return count[names[i]] > count[names[j]]
			}
			return names[i] < names[j]
		})
		for _, n := range names {
			fmt.Printf("%4d  %s\n", count[n], n)
		}

	case *find != "":
		re, err := regexp.Compile("(?i)" + *find)
		if err != nil {
			fmt.Fprintln(os.Stderr, "плохое выражение:", err)
			os.Exit(1)
		}
		hits := 0
		for _, n := range nodes {
			if re.MatchString(n.name) {
				printNode(n, *full)
				hits++
			}
		}
		fmt.Fprintf(os.Stderr, "\n· совпало узлов: %d\n", hits)

	case *at >= 0:
		lo, hi := *at-*ctx, *at+*ctx
		if lo < 0 {
			lo = 0
		}
		if hi >= len(nodes) {
			hi = len(nodes) - 1
		}
		for i := lo; i <= hi; i++ {
			printNode(nodes[i], *full || i == *at)
		}

	case *props != "":
		keys := strings.Split(*props, ",")
		for i := range keys {
			keys[i] = strings.TrimSpace(keys[i])
		}
		for _, n := range nodes {
			var cells []string
			for _, k := range keys {
				if v := n.get(k); v != "" {
					cells = append(cells, k+"="+v)
				}
			}
			if len(cells) > 0 {
				fmt.Printf("#%-5d %-38s %s\n", n.idx, trunc(n.name, 38), strings.Join(cells, "  "))
			}
		}

	case *sum:
		for _, k := range []string{"font-family", "font-size", "font-weight", "line-height",
			"letter-spacing", "border-radius", "border", "padding", "gap", "background", "color", "box-shadow"} {
			vals := map[string]int{}
			for _, n := range nodes {
				if v := n.get(k); v != "" {
					vals[v]++
				}
			}
			if len(vals) == 0 {
				continue
			}
			type kv struct {
				v string
				c int
			}
			out := make([]kv, 0, len(vals))
			for v, c := range vals {
				out = append(out, kv{v, c})
			}
			sort.Slice(out, func(i, j int) bool {
				if out[i].c != out[j].c {
					return out[i].c > out[j].c
				}
				return out[i].v < out[j].v
			})
			fmt.Printf("── %s (%d различных)\n", k, len(out))
			for i, e := range out {
				if i == 24 {
					fmt.Printf("   … и ещё %d\n", len(out)-24)
					break
				}
				fmt.Printf("   %4d  %s\n", e.c, e.v)
			}
			fmt.Println()
		}

	default:
		flag.Usage()
	}
}

func trunc(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n-1] + "…"
}
