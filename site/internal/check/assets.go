package check

import (
	"fmt"
	"regexp"
	"sort"
	"strings"

	"instrument/site/internal/content"
)

// Собственные переменные сайта: те, которым в ките места нет и не будет.
// --c и --v — канал данных у образца цвета и у строки API.
// --code-* — подсветка синтаксиса: кит кода не подсвечивает, это забота
// справочника о своём содержимом.
var ownVars = map[string]bool{
	"--c": true, "--v": true,
	"--code-tag": true, "--code-attr": true, "--code-val": true,
}

var (
	cssVarUse = regexp.MustCompile(`var\((--[a-z0-9-]+)`)

	bannedProps = regexp.MustCompile(`(?m)^\s*(text-transform)\s*:`)

	hexColor = regexp.MustCompile(`#[0-9a-fA-F]{3,8}\b`)

	pastTense = regexp.MustCompile(`(?i)(^|[.!?;] )раньше|ни разу не|выяснилось|` +
		`здесь сто(ял|яла|яли|яло)|было объявлено`)
)

func Assets(files map[string]string, tokens map[string]content.Token) []string {
	var problems []string

	for name, css := range files {
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
			problems = append(problems, fmt.Sprintf(
				"%s  цвет %s в обход семантики: захардкожена одна тема", name, m))
		}

	}

	sort.Strings(problems)
	return problems
}

// StrayCommentEnd ищет `*/`, который закрыл комментарий раньше времени.
//
// Комментарии в CSS не вкладываются. Строка `--pad-*/--gap-*` внутри
// пояснения закрывает его на месте: остаток текста становится мусором, а
// парсер, восстанавливаясь, доедает СЛЕДУЮЩЕЕ правило целиком. Ошибка
// молчаливая — файл грузится, стили частично работают, и найти её можно
// только по отсутствующему поведению.
//
// Так из кита пропали тон `neutral` (сноска без значка и без заливки) и
// область нажатия у чекбокса, радио и свитча — то есть требование WCAG
// 2.2 AA, ради которого правило и написано.
func StrayCommentEnd(files map[string]string) []string {
	var problems []string
	for name, src := range files {
		if !strings.HasSuffix(name, ".css") {
			continue
		}
		var code strings.Builder
		rest := src
		for {
			i := strings.Index(rest, "/*")
			if i < 0 {
				code.WriteString(rest)
				break
			}
			code.WriteString(rest[:i])
			j := strings.Index(rest[i+2:], "*/")
			if j < 0 {
				break
			}
			rest = rest[i+2+j+2:]
		}
		if k := strings.Index(code.String(), "*/"); k >= 0 {
			line := strings.Count(code.String()[:k], "\n") + 1
			problems = append(problems, fmt.Sprintf(
				"%s:%d  комментарий закрылся раньше времени: `*/` внутри пояснения съедает следующее правило",
				name, line))
		}
	}
	sort.Strings(problems)
	return problems
}

func Comments(files map[string]string) []string {
	var problems []string
	for name, src := range files {
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

// comments вынимает из исходника то, что человек написал ЧЕЛОВЕКУ.
//
// Строковые литералы при этом пропускаются, и это не педантизм. Разбор по
// первому `//` в строке считает комментарием всё, что стоит за кавычками, —
// а мутационный стенд хранит запрещённые фразы ДАННЫМИ: у него это текст
// мутации, а не пояснение к соседней функции. Без такого различения стенд
// роняет сборку сайта собственным содержимым, и роняет заслуженно: с точки
// зрения разбора он неотличим от нарушителя.
//
// Разбор построчный: многострочный сырой литерал Go он не отследит. Цена
// ошибки при этом односторонняя — пропущенный комментарий, а не выдуманный,
// — и это правильная сторона: гейт, который ругается на то, чего человек не
// писал, отучают читать первым.
func comments(src string) string {
	var b strings.Builder
	for _, line := range strings.Split(src, "\n") {
		if i := lineComment(line); i >= 0 {
			b.WriteString(line[i+2:])
			b.WriteByte('\n')
		}
	}
	for {
		i := blockOpen(src)
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

// lineComment — позиция `//`, начинающего комментарий, или -1.
func lineComment(line string) int {
	return outsideQuotes(line, "//", false)
}

// blockOpen — позиция `/*`, начинающего комментарий, или -1.
func blockOpen(src string) int {
	return outsideQuotes(src, "/*", true)
}

// outsideQuotes ищет первое вхождение want вне кавычек. Кавычка, не закрытая
// до конца строки, кавычкой не считается: это апостроф в прозе, а не литерал.
func outsideQuotes(src, want string, multiline bool) int {
	var quote byte
	for i := 0; i < len(src); i++ {
		c := src[i]
		if c == '\n' {
			if !multiline {
				break
			}
			quote = 0
			continue
		}
		if quote != 0 {
			if c == '\\' {
				i++
				continue
			}
			if c == quote {
				quote = 0
			}
			continue
		}
		if c == '"' || c == '\'' || c == '`' {
			// Одиночная кавычка без пары до конца строки — апостроф.
			if !closes(src[i+1:], c) {
				continue
			}
			quote = c
			continue
		}
		if strings.HasPrefix(src[i:], want) {
			return i
		}
	}
	return -1
}

// closes отвечает, встречается ли кавычка ещё раз до конца строки.
func closes(rest string, q byte) bool {
	for i := 0; i < len(rest); i++ {
		if rest[i] == '\n' {
			return false
		}
		if rest[i] == '\\' {
			i++
			continue
		}
		if rest[i] == q {
			return true
		}
	}
	return false
}

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

	return dataURI.ReplaceAllString(b.String(), `url("")`)
}

var dataURI = regexp.MustCompile(`url\("data:[^"]*"\)`)
