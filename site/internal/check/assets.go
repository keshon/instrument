package check

import (
	"fmt"
	"regexp"
	"sort"
	"strings"

	"instrument/site/internal/content"
)

// Assets проверяет CSS САМОГО САЙТА против кита.
//
// Ловится ошибка одного вида: объявление, которое браузер молча отбрасывает.
// `letter-spacing: var(--tracking-wide)` при токенах `--tracking-tight` и
// `--tracking-normal` не делает ничего, а страница остаётся внешне целой.
//
// Остальные проверки сюда не смотрят: docs-check сверяет документацию с китом,
// check.Verify — собранные страницы со ссылками, contrast — контраст. CSS
// сайта одевает всё, что они проверяют, и не проверяется ничем.
//
// Проверяется два правила, и оба взяты не из вкуса, а из конституции.

// ownVars — переменные, которые принадлежат САЙТУ, а не киту.
//
// Их две, обе в demo.css, и обе — канал данных, а не оформления: цвет образца
// и величина шкалы приходят инлайном, потому что показываемое значение и есть
// содержание примера. Список закрыт намеренно: третья такая переменная должна
// потребовать правки здесь и объяснения, зачем она.
var ownVars = map[string]bool{"--c": true, "--v": true}

var (
	cssVarUse = regexp.MustCompile(`var\((--[a-z0-9-]+)`)

	// text-transform во всём ките не встречается ни разу. Это не совпадение:
	// третий принцип конституции называет капс с трекингом дословно —
	// «кричит громче числа, которое он подписывает, — значит, он неправ».
	// Сайт не имеет права заводить типографический приём, которого у кита нет.
	bannedProps = regexp.MustCompile(`(?m)^\s*(text-transform)\s*:`)

	// Цвет в обход семантики — второй принцип: тот, кто написал #333, только
	// что захардкодил светлую тему.
	hexColor = regexp.MustCompile(`#[0-9a-fA-F]{3,8}\b`)

	// Комментарий объясняет, ПОЧЕМУ так сейчас, а не что было раньше.
	//
	// Граница проходит между двумя видами отрицания. «Не :root, потому что
	// тогда тему нельзя дать поддереву» — суть: контрфактическое утверждение о
	// коде, которое читателю нужно. Рассказ о том, какое правило в этом месте
	// когда-то было и как оно ломалось, — биография файла: читателю она не
	// нужна, а с каждой правкой такие фразы копятся, и комментарий
	// превращается в журнал изменений.
	//
	// «Раньше» ловится ТОЛЬКО в начале предложения.
	//
	// В середине это слово означает порядок, а не время: «kit.layout объявлен
	// РАНЬШЕ kit.components», «закрыл бы раздел раньше времени». С начала
	// предложения оно всегда начинает рассказ.
	//
	// Без \b: в Go граница слова определена по ASCII, и перед кириллической
	// буквой её нет — образец с \b не совпадает никогда, а проверка молчит.
	pastTense = regexp.MustCompile(`(?i)(^|[.!?;] )раньше|ни разу не|выяснилось|` +
		`здесь сто(ял|яла|яли|яло)|было объявлено`)
)

// Assets сверяет таблицы стилей сайта с китом.
//
// tokens — то, что объявлено в src/*.css. Источник тот же, что у справочника:
// значения токенов и здесь берутся из кита, а не из списка в голове.
func Assets(files map[string]string, tokens map[string]content.Token) []string {
	var problems []string

	for name, css := range files {
		// Комментарии вырезаются: в них законно упоминаются и #333 как пример
		// нарушения, и имена токенов, которых нет.
		code := stripComments(css)

		seen := map[string]bool{}
		for _, m := range cssVarUse.FindAllStringSubmatch(code, -1) {
			v := m[1]
			if seen[v] || ownVars[v] {
				continue
			}
			seen[v] = true
			if _, ok := tokens[v]; !ok {
				problems = append(problems, fmt.Sprintf(
					"%s  токена %s в ките нет: объявление молча отбрасывается", name, v))
			}
		}

		for _, m := range bannedProps.FindAllStringSubmatch(code, -1) {
			problems = append(problems, fmt.Sprintf(
				"%s  свойство %s не встречается в ките ни разу — сайту его заводить нечем оправдать", name, m[1]))
		}

		for _, m := range hexColor.FindAllString(code, -1) {
			// data-URI масок — это SVG, и там цвет служебный: маска красится
			// currentColor, а %23000 внутри неё браузер обязан видеть.
			problems = append(problems, fmt.Sprintf(
				"%s  цвет %s в обход семантики: захардкожена одна тема", name, m))
		}

	}

	sort.Strings(problems)
	return problems
}

// Comments проверяет комментарии — и кита, и сайта.
//
// Кит здесь важнее сайта: его исходники открывают, чтобы понять решение, а не
// чтобы узнать, какие правила из него удаляли. Проверка одна на оба, потому
// что правило одно.
func Comments(files map[string]string) []string {
	var problems []string
	for name, src := range files {
		// Переносы строк схлопываются: «…поверхность.\n   Раньше…» — это одно
		// предложение, разбитое переносом, и без нормализации начало
		// предложения не опознаётся.
		text := strings.Join(strings.Fields(comments(src)), " ")
		seen := map[string]bool{}
		for _, m := range pastTense.FindAllString(text, -1) {
			key := strings.ToLower(m)
			if seen[key] {
				continue
			}
			seen[key] = true
			problems = append(problems, fmt.Sprintf(
				"%s  комментарий рассказывает историю («%s»): нужно, почему так сейчас", name, m))
		}
	}
	sort.Strings(problems)
	return problems
}

// comments — только содержимое комментариев: обратная сторона stripComments.
//
// Понимает обе формы: /* … */ у CSS и Go, и // … у Go.
func comments(src string) string {
	var b strings.Builder
	for _, line := range strings.Split(src, "\n") {
		if i := strings.Index(line, "//"); i >= 0 {
			b.WriteString(line[i+2:])
			b.WriteByte('\n')
		}
	}
	for {
		i := strings.Index(src, "/*")
		if i < 0 {
			return b.String()
		}
		src = src[i+2:]
		j := strings.Index(src, "*/")
		if j < 0 {
			b.WriteString(src)
			return b.String()
		}
		b.WriteString(src[:j])
		b.WriteByte('\n')
		src = src[j+2:]
	}
}

// stripComments убирает /* … */ — и заодно data-URI, внутри которых стоят
// служебные цвета SVG-масок. Маска красится currentColor, поэтому %23000 в ней
// не является цветом интерфейса и под второй принцип не попадает.
func stripComments(css string) string {
	var b strings.Builder
	b.Grow(len(css))
	for {
		i := strings.Index(css, "/*")
		if i < 0 {
			break
		}
		j := strings.Index(css[i:], "*/")
		if j < 0 {
			css = css[:i]
			break
		}
		b.WriteString(css[:i])
		css = css[i+j+2:]
	}
	b.WriteString(css)

	// url("data:…") — служебное содержимое, а не оформление.
	return dataURI.ReplaceAllString(b.String(), `url("")`)
}

var dataURI = regexp.MustCompile(`url\("data:[^"]*"\)`)
