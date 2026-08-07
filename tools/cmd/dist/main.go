// Команда dist собирает кит в один файл.
//
// Кит подключается одной строкой и не требует сборки У ПОТРЕБИТЕЛЯ — это
// остаётся правдой. Но «нет сборки у потребителя» и «в репозитории нет
// собранного файла» — разные вещи, и вторая обходилась дорого:
//
//	одиннадцать запросов вместо одного (@import — водопад: браузер сначала
//	забирает kit.css, только потом узнаёт про остальные десять);
//	227 КБ сырых, из которых 47% — комментарии по-русски.
//
// Источником остаётся src/. dist/ — вывод, и он лежит в репозитории только
// затем, чтобы потребитель мог взять один файл.
//
//	go run ./cmd/dist            собрать
//	go run ./cmd/dist -check     упасть, если собранное разошлось с исходником
package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var (
	layerStmtRe = regexp.MustCompile(`@layer\s+([^;{]+);`)
	importRe    = regexp.MustCompile(`@import\s+url\(["']?\./([^"')]+)["']?\)\s*layer\(([^)]+)\)\s*;`)
)

func main() {
	var (
		src     = flag.String("src", "../src", "каталог кита")
		out     = flag.String("out", "../dist", "каталог вывода")
		version = flag.String("version", "", "версия; по умолчанию из VERSION")
		check   = flag.Bool("check", false, "не писать, а сверить")
	)
	flag.Parse()

	ver := *version
	if ver == "" {
		b, err := os.ReadFile(filepath.Join(filepath.Dir(*src), "VERSION"))
		if err != nil {
			fmt.Fprintln(os.Stderr, "не прочитать VERSION:", err)
			os.Exit(1)
		}
		ver = strings.TrimSpace(string(b))
	}

	entry, err := os.ReadFile(filepath.Join(*src, "kit.css"))
	if err != nil {
		fmt.Fprintln(os.Stderr, "не прочитать kit.css:", err)
		os.Exit(1)
	}

	// Порядок слоёв берётся из самого kit.css, а не переписывается сюда:
	// второй список разошёлся бы с первым, и разошёлся бы молча.
	stmt := layerStmtRe.FindSubmatch(entry)
	if stmt == nil {
		fmt.Fprintln(os.Stderr, "в kit.css нет объявления @layer")
		os.Exit(1)
	}

	var css bytes.Buffer
	fmt.Fprintf(&css, "/*! instrument %s — https://github.com/keshon/instrument\n"+
		"    Собран из src/. Правится ИСХОДНИК, а не этот файл. */\n", ver)
	fmt.Fprintf(&css, "@layer %s;\n\n", strings.Join(strings.Fields(string(stmt[1])), " "))

	imports := importRe.FindAllSubmatch(entry, -1)
	if len(imports) == 0 {
		fmt.Fprintln(os.Stderr, "в kit.css нет ни одного @import")
		os.Exit(1)
	}
	for _, m := range imports {
		name, layer := string(m[1]), string(m[2])
		b, err := os.ReadFile(filepath.Join(*src, name))
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		// Переносы строк нормализуются на входе. Иначе собранное зависит от
		// того, на какой системе его собрали: git отдаёт исходники с CRLF на
		// Windows, и -check не совпадал сам с собой.
		b = bytes.ReplaceAll(b, []byte("\r\n"), []byte("\n"))
		// Импорт в слой вкладывает ВСЁ содержимое файла — ровно это и
		// воспроизводится блоком. Иначе слои разъехались бы, а порядок слоёв
		// у кита несёт смысл: motion и print обязаны перебивать компоненты.
		fmt.Fprintf(&css, "/* ── %s → %s ─────────────────────────────── */\n", name, layer)
		fmt.Fprintf(&css, "@layer %s {\n%s\n}\n\n", layer, b)
	}

	js, err := os.ReadFile(filepath.Join(*src, "kit.js"))
	if err != nil {
		fmt.Fprintln(os.Stderr, "не прочитать kit.js:", err)
		os.Exit(1)
	}
	js = bytes.ReplaceAll(js, []byte("\r\n"), []byte("\n"))

	files := map[string][]byte{
		"instrument.css":     css.Bytes(),
		"instrument.min.css": []byte(minify(css.String())),
		"instrument.js":      js,
	}

	if *check {
		bad := 0
		for name, want := range files {
			got, err := os.ReadFile(filepath.Join(*out, name))
			if err != nil {
				fmt.Fprintf(os.Stderr, "  нет %s: %v\n", name, err)
				bad++
				continue
			}
			if !bytes.Equal(bytes.ReplaceAll(got, []byte("\r\n"), []byte("\n")), want) {
				fmt.Fprintf(os.Stderr, "  %s разошёлся с исходником\n", name)
				bad++
			}
		}
		if bad > 0 {
			fmt.Fprintln(os.Stderr, "\nСобранное отстало от src/. Пересоберите:")
			fmt.Fprintln(os.Stderr, "  go -C tools run ./cmd/dist")
			os.Exit(1)
		}
		fmt.Println("· dist совпадает с src")
		return
	}

	if err := os.MkdirAll(*out, 0o755); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	for name, b := range files {
		if err := os.WriteFile(filepath.Join(*out, name), b, 0o644); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}

	fmt.Printf("instrument %s\n", ver)
	for _, name := range []string{"instrument.css", "instrument.min.css", "instrument.js"} {
		fmt.Printf("  %-18s %6.1f КБ\n", name, float64(len(files[name]))/1024)
	}
	fmt.Printf("\nбыло: %d файлов, %.1f КБ, одиннадцать запросов\n",
		len(imports)+1, float64(len(css.Bytes()))/1024)
}

/*
minify — намеренно осторожный.

Он снимает комментарии и лишние пробелы и НЕ ТРОГАЕТ больше ничего. Соблазн
дожать велик, но каждая следующая оптимизация опасна на этом конкретном CSS:

  - пробел перед «:» убирать нельзя. «.a :hover» и «.a:hover» — разные
    селекторы, и отличить их без разбора селектора невозможно;
  - содержимое строк неприкосновенно. Кит рисует формы масками из data-URI, а
    внутри них живут и пробелы, и «;», и «}» — схлопнуть их значит стереть
    половину значков;
  - «url(...)» без кавычек — то же самое.

Комментарии дают 47% файла, поэтому осторожности хватает с запасом.
*/
func minify(s string) string {
	var b strings.Builder
	b.Grow(len(s))

	var quote byte     // текущая кавычка, 0 — вне строки
	inURL := false     // внутри url(...) без кавычек
	pendingWS := false // видели пробел, ещё не решили, писать ли
	keepFirst := true  // первый комментарий /*! сохраняется — это шапка

	writeByte := func(c byte) { b.WriteByte(c) }

	for i := 0; i < len(s); i++ {
		c := s[i]

		if quote != 0 {
			writeByte(c)
			if c == '\\' && i+1 < len(s) {
				i++
				writeByte(s[i])
			} else if c == quote {
				quote = 0
			}
			continue
		}

		if inURL {
			writeByte(c)
			if c == ')' {
				inURL = false
			}
			continue
		}

		switch {
		case c == '"' || c == '\'':
			flushWS(&b, &pendingWS)
			quote = c
			writeByte(c)
			continue

		case c == '/' && i+1 < len(s) && s[i+1] == '*':
			end := strings.Index(s[i+2:], "*/")
			if end < 0 {
				return b.String()
			}
			if keepFirst && i+2 < len(s) && s[i+2] == '!' {
				b.WriteString(s[i : i+2+end+2])
				b.WriteByte('\n')
				keepFirst = false
			} else {
				// Комментарий равносилен пробелу: «a/**/b» — два токена.
				pendingWS = true
			}
			i += 2 + end + 1
			continue

		case c == ' ' || c == '\t' || c == '\n' || c == '\r' || c == '\f':
			pendingWS = true
			continue
		}

		// url( без кавычек
		if c == '(' && strings.HasSuffix(strings.ToLower(b.String()), "url") {
			flushWS(&b, &pendingWS)
			writeByte(c)
			// Пропускаем пробелы сразу после «url(»; если дальше кавычка,
			// обычная ветка строки её подхватит.
			j := i + 1
			for j < len(s) && (s[j] == ' ' || s[j] == '\n' || s[j] == '\t') {
				j++
			}
			if j < len(s) && s[j] != '"' && s[j] != '\'' {
				inURL = true
			}
			i = j - 1
			continue
		}

		// Пробел не нужен рядом со скобками и разделителями объявлений.
		if pendingWS {
			prev := lastByte(&b)
			if !isDrop(prev) && !isDrop(c) {
				writeByte(' ')
			}
			pendingWS = false
		}
		writeByte(c)
	}
	return strings.TrimSpace(b.String()) + "\n"
}

func flushWS(b *strings.Builder, pending *bool) {
	if *pending {
		if p := lastByte(b); p != 0 && !isDrop(p) {
			b.WriteByte(' ')
		}
		*pending = false
	}
}

func lastByte(b *strings.Builder) byte {
	s := b.String()
	if len(s) == 0 {
		return 0
	}
	return s[len(s)-1]
}

// isDrop — символы, рядом с которыми пробел не значит ничего.
//
// «:» и «,» сюда НЕ входят намеренно: у первого пробел меняет селектор, у
// второго безопасно, но выигрыш копеечный, а правило проще держать одно.
func isDrop(c byte) bool {
	switch c {
	case '{', '}', ';', 0:
		return true
	}
	return false
}
