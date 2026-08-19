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
	"encoding/json"
	"errors"
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

	// Версия объявлена ДВАЖДЫ: в VERSION и в package.json. Разойтись им нечем
	// — кроме забывчивости, а она случается ровно в момент выпуска.
	//
	// Цена расхождения не косметическая: шапка dist/ печатает VERSION, а npm
	// публикует package.json, и в реестр уезжает пакет, внутри которого
	// написана другая версия. Дальше это не отследить: файл на CDN лежит с
	// одним номером, а комментарий внутри него называет другой.
	if err := checkPkgVersion(filepath.Dir(*src), ver); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
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
		if err := checkBraces(name, b); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		if err := checkRoleTier(name, b); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		if err := checkBans(name, b); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
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

// checkBraces сверяет баланс фигурных скобок в файле кита.
//
// Импорт в слой вкладывает файл в блок `@layer X { ... }`, поэтому лишняя `}`
// закрывает не правило, а сам слой: остаток файла — и всё, что импортировано в
// тот же слой следом, — оказывается ВНЕ слоёв. Неслойное правило выигрывает у
// любого слоя, то есть у таких компонентов перестают работать переопределения
// приложения, prefers-reduced-motion, @media print и forced-colors.
//
// В src/ через @import баг не виден: лишняя скобка на верхнем уровне —
// синтаксическая ошибка, браузер её отбрасывает, а слой назначает импорт.
// Сравнение dist/ с src/ его тоже не видит: скобка одинаково стоит в обоих.
// Поэтому проверка идёт по исходнику и до сборки.
//
// Строки, комментарии и url(...) без кавычек пропускаются: внутри масок из
// data-URI скобки — часть рисунка, а не структуры.
func checkBraces(name string, css []byte) error {
	s := string(css)
	line := 1
	depth := 0
	openLine := 0 // строка последней незакрытой «{»

	var quote byte
	inURL := false

	for i := 0; i < len(s); i++ {
		c := s[i]
		if c == '\n' {
			line++
		}

		switch {
		case quote != 0:
			if c == '\\' && i+1 < len(s) {
				i++
			} else if c == quote {
				quote = 0
			}
		case inURL:
			if c == ')' {
				inURL = false
			}
		case c == '"' || c == '\'':
			quote = c
		case c == '/' && i+1 < len(s) && s[i+1] == '*':
			end := strings.Index(s[i+2:], "*/")
			if end < 0 {
				return fmt.Errorf("%s:%d: комментарий не закрыт", name, line)
			}
			line += strings.Count(s[i:i+2+end+2], "\n")
			i += 2 + end + 1
		case c == '(' && i >= 3 && strings.EqualFold(s[i-3:i], "url"):
			j := i + 1
			for j < len(s) && (s[j] == ' ' || s[j] == '\t' || s[j] == '\n') {
				j++
			}
			if j < len(s) && s[j] != '"' && s[j] != '\'' {
				inURL = true
			}
		case c == '{':
			if depth == 0 {
				openLine = line
			}
			depth++
		case c == '}':
			depth--
			if depth < 0 {
				return fmt.Errorf(
					"%s:%d: лишняя «}».\n"+
						"Файл вкладывается в блок @layer, поэтому она закроет слой, а не правило:\n"+
						"остаток файла окажется вне слоёв и начнёт выигрывать у стилей приложения,\n"+
						"prefers-reduced-motion и @media print.",
					name, line)
			}
		}
	}

	if depth > 0 {
		return fmt.Errorf("%s: не закрыто блоков: %d, первый открыт на строке %d",
			name, depth, openLine)
	}
	if quote != 0 {
		return fmt.Errorf("%s: строка не закрыта", name)
	}
	return nil
}

// checkPkgVersion сверяет версию в package.json с VERSION.
//
// Отсутствие package.json — не ошибка: кит работает файлом и без реестра, и
// репозиторий, из которого его берут ссылкой, package.json не обязан иметь.
// А вот РАСХОЖДЕНИЕ — ошибка, и молчаливая: она видна только после публикации.
func checkPkgVersion(root, ver string) error {
	b, err := os.ReadFile(filepath.Join(root, "package.json"))
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("не прочитать package.json: %w", err)
	}

	var pkg struct {
		Version string `json:"version"`
	}
	if err := json.Unmarshal(b, &pkg); err != nil {
		return fmt.Errorf("package.json не разбирается: %w", err)
	}
	if pkg.Version != ver {
		return fmt.Errorf(
			"версии разошлись: VERSION = %s, package.json = %s.\n"+
				"В шапке dist/ печатается первая, а в реестр уезжает вторая —\n"+
				"на CDN окажется файл, который называет сам себя не так, как пакет вокруг него.",
			ver, pkg.Version)
	}
	return nil
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

// rawGap ловит обращение к шкале пространства напрямую там, где на просвет уже
// есть роль.
//
// Правило кита: «компоненты берут отсюда только через ярус ролей». Оно держалось
// на честном слове и однажды сломалось: .inst-stack--loose нёс
// gap: var(--space-7), потому что роли на просвет между разделами не было, а
// 24px были нужны. Роль завелась (--gap-section), обход убран — но ничто не
// мешало завести его снова.
//
// Ловятся только КРУПНЫЕ ступени и только на просвете. Мелкие (1–6) уходят на
// микрорасстояния, где роли нет и заводить её не на что: просвет между значком
// и подписью — не то же, что просвет между разделами. Набивка не ловится вовсе:
// у неё свои роли, а арифметика вьюпорта вроде 100vw - var(--space-8) к ярусу
// отношения не имеет.
var rawGap = regexp.MustCompile(`(?:^|[;{[:space:]])(?:row-|column-)?gap:[[:space:]]*var\(--space-(7|8|9|10)\)`)

// checkRoleTier запрещает просвет мимо яруса ролей. tokens.css исключён: там
// роли и ОБЪЯВЛЯЮТСЯ через ступени шкалы, в этом их работа.
func checkRoleTier(name string, css []byte) error {
	if name == "tokens.css" {
		return nil
	}
	for i, line := range strings.Split(string(css), "\n") {
		if m := rawGap.FindStringSubmatch(line); m != nil {
			return fmt.Errorf("%s:%d: просвет мимо яруса ролей — var(--space-%s). "+
				"Между разделами берите --gap-section, внутри — --gap-row или --gap-inline",
				name, i+1, m[1])
		}
	}
	return nil
}

// blankComments заменяет тела комментариев пробелами, СОХРАНЯЯ переносы строк.
//
// Без этого проверки ниже ловили бы сами объяснения: слово «!important»
// встречается в шапках kit.css, motion.css, print.css и forms.css, а «#333» —
// в описании второго закона прямо в tokens.css. Вырезать комментарии целиком
// нельзя: номер строки в сообщении обязан совпадать с тем, что откроют.
func blankComments(css []byte) []byte {
	out := make([]byte, len(css))
	copy(out, css)
	for i := 0; i+1 < len(out); {
		if out[i] != '/' || out[i+1] != '*' {
			i++
			continue
		}
		j := i
		for ; j+1 < len(out); j++ {
			if out[j] == '*' && out[j+1] == '/' {
				j += 2
				break
			}
		}
		if j+1 >= len(out) {
			j = len(out)
		}
		for k := i; k < j; k++ {
			if out[k] != '\n' {
				out[k] = ' '
			}
		}
		i = j
	}
	return out
}

// Раздел «Запрещено» принципов дизайна, переведённый в регулярные выражения.
//
// Все пять правил СЕГОДНЯ ЗЕЛЁНЫЕ — это защёлки, а не работа. Смысл гейта не в
// том, чтобы что-то найти сейчас, а в том, чтобы запрет перестал держаться на
// памяти: каждое из пяти нарушается одной строкой, которая выглядит безобидно
// и не даёт ошибки в браузере.
var (
	banImportant = regexp.MustCompile(`!\s*important`)
	banBold      = regexp.MustCompile(`font-weight:\s*(700|800|900|bold)\b`)
	banColor     = regexp.MustCompile(`#[0-9a-fA-F]{3,8}\b|\brgba?\(|\bhsla?\(|\boklch\(|\bcolor\(`)
	banUtility   = regexp.MustCompile(`\.[mp][trblxy]?-[0-9]`)

	maskDecl = regexp.MustCompile(`\bmask(-image)?\s*:`)
	maskInk  = regexp.MustCompile(`^#(000|000000)$`)
	uriHex   = regexp.MustCompile(`%23[0-9a-fA-F]{3,8}`)
	urlDecl  = regexp.MustCompile(`([\w-]+)\s*:\s*[^;]*url\(`)
)

// Слой токенов. Здесь сырой цвет — это РАБОТА, а не нарушение: tokens.css
// объявляет рампы и семантику, print.css возвращает светлую тему на бумаге,
// переприсваивая те же семантические имена. Оба яруса 1–2, и обоим положено
// называть цвет числом. Всем остальным — нет.
var colorLayer = map[string]bool{"tokens.css": true, "print.css": true}

// checkBans сторожит запреты, которые до сих пор держались на внимательности.
//
// Проверка построчная и потому не видит объявления, разорванного переносом.
// Для четырёх правил из пяти это безразлично — свойство и значение стоят на
// одной строке, — а для маски проверено отдельно: во всех 27 объявлениях кита
// маска однострочная.
func checkBans(name string, raw []byte) error {
	css := blankComments(raw)
	for i, line := range strings.Split(string(css), "\n") {
		at := fmt.Sprintf("%s:%d", name, i+1)

		// !important. Исключение ровно одно и названо поимённо: [hidden] в
		// base.css — это корректность, а не оформление.
		if banImportant.MatchString(line) {
			if !(name == "base.css" && strings.Contains(line, "[hidden]")) {
				return fmt.Errorf("%s: !important. Единственное разрешённое — [hidden] в base.css. "+
					"Порядок слоёв решает то же самое и не ломает обещание «стили приложения выигрывают»", at)
			}
		}

		// Вес 700. base.css закрывает единственную дверь, через которую его
		// приносит платформа (strong и b сброшены на --weight-medium); эта
		// закрывает дверь, через которую его принесёт правка.
		if m := banBold.FindStringSubmatch(line); m != nil {
			return fmt.Errorf("%s: вес %s. В ките два веса — --weight-normal и --weight-medium (600). "+
				"Настоящего 500 у Segoe UI нет, а 700 в инструментальном интерфейсе кричит громче данных", at, m[1])
		}

		// Утилиты отступов. Шкала разрежена сверху намеренно, и набор утилит
		// вернул бы «чуть побольше» первым классом.
		if m := banUtility.FindString(line); m != "" {
			return fmt.Errorf("%s: утилита отступа %q. Ритм задают примитивы потока "+
				"(.inst-stack, .inst-cluster, .inst-grid) с зазором, названным намерением", at, m)
		}

		// Сырой цвет вне слоя токенов — это захардкоженная светлая тема.
		if !colorLayer[name] {
			for _, c := range banColor.FindAllString(line, -1) {
				// Маска — исключение, и оно узкое: цвет в ней не цвет, а
				// альфа-канал, которым вырезается форма. Заливка приходит
				// токеном, поэтому чёрный здесь единственное осмысленное
				// значение, и только он и разрешён.
				if maskDecl.MatchString(line) && maskInk.MatchString(c) {
					continue
				}
				return fmt.Errorf("%s: сырой цвет %q. Компонент обращается к семантике "+
					"(--text-primary, --surface-raised), а не к числу: число — это захардкоженная светлая тема", at, c)
			}
		}

		// Цвет внутри data-URI. Форма рисуется маской, цвет приходит токеном —
		// значит в самой картинке цвета нет, есть только чернила формы.
		for _, h := range uriHex.FindAllString(line, -1) {
			if h != "%23000" {
				return fmt.Errorf("%s: цвет %s внутри data-URI. Форма рисуется маской и красится токеном; "+
					"внутри картинки допустим только %%23000 — чернила самой формы", at, h)
			}
		}

		// url() только в маске. Картинка, поставленная фоном, красится собой и
		// потому не умеет следовать теме.
		if m := urlDecl.FindStringSubmatch(line); m != nil && !strings.HasPrefix(m[1], "mask") {
			return fmt.Errorf("%s: url() в свойстве %q. Картинка допускается только маской (mask, mask-image): "+
				"фоновая красится собой и не следует ни теме, ни тону", at, m[1])
		}
	}
	return nil
}
