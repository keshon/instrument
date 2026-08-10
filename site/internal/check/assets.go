package check

import (
	"fmt"
	"regexp"
	"sort"
	"strings"

	"instrument/site/internal/content"
)

var ownVars = map[string]bool{"--c": true, "--v": true}

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
