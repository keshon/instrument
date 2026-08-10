package content

import (
	"fmt"
	"regexp"
	"strings"

	"instrument/site/internal/i18n"
)

var (
	directiveOpen  = regexp.MustCompile(`^:::([a-z][a-z-]*)\s*(.*)$`)
	directiveClose = regexp.MustCompile(`^:::\s*$`)
	fenceRe        = regexp.MustCompile("^(```|~~~)")
)

var ruleKinds = map[string]bool{"do": true, "dont": true}

var noteTones = map[string]string{
	"note":   "info",
	"warn":   "warn",
	"danger": "error",
	"ok":     "ok",
}

func expandDirectives(body []byte, lang i18n.Lang) ([]byte, error) {
	lines := strings.Split(string(body), "\n")
	var out []string

	inFence := false
	rowOpen := false

	type openBlock struct{ kind, closes string }
	var open []openBlock

	closeRow := func() {
		if rowOpen {
			out = append(out, "", "</div>", "")
			rowOpen = false
		}
	}

	for i := 0; i < len(lines); i++ {
		line := lines[i]

		if fenceRe.MatchString(strings.TrimSpace(line)) {
			inFence = !inFence
			out = append(out, line)
			continue
		}
		if inFence {
			out = append(out, line)
			continue
		}

		if m := directiveOpen.FindStringSubmatch(strings.TrimSpace(line)); m != nil {
			kind, head := m[1], strings.TrimSpace(m[2])
			var closes string
			switch {
			case ruleKinds[kind]:
				if !rowOpen {
					out = append(out, "", `<div class="rule-row">`, "")
					rowOpen = true
				}
				label := i18n.T(lang, "rule."+kind)
				out = append(out,
					fmt.Sprintf(`<div class="rule" data-rule="%s">`, kind),
					"",
					fmt.Sprintf(`<p class="rule-tag">%s%s</p>`, escape(label), ruleHead(head)),
					"")
				closes = "</div>"
			case noteTones[kind] != "":
				closeRow()
				out = append(out, "",
					fmt.Sprintf(`<div class="inst-note" data-tone="%s"><div class="note-body">`, noteTones[kind]),
					"")
				closes = "</div></div>"
			default:
				return nil, fmt.Errorf("строка %d: неизвестный блок :::%s (есть do · dont · note · warn · danger · ok)", i+1, kind)
			}
			open = append(open, openBlock{kind, closes})
			continue
		}

		if directiveClose.MatchString(strings.TrimSpace(line)) {
			if len(open) == 0 {
				return nil, fmt.Errorf("строка %d: закрытие ::: без открытия", i+1)
			}
			b := open[len(open)-1]
			open = open[:len(open)-1]
			out = append(out, "", b.closes, "")
			if !ruleKinds[b.kind] {
				closeRow()
			}
			continue
		}

		if rowOpen && len(open) == 0 && strings.TrimSpace(line) != "" {
			closeRow()
		}
		out = append(out, line)
	}

	if len(open) > 0 {
		return nil, fmt.Errorf("блок :::%s не закрыт", open[len(open)-1].kind)
	}
	closeRow()
	return []byte(strings.Join(out, "\n")), nil
}

func ruleHead(s string) string {
	if s == "" {
		return ""
	}
	return ` <span class="rule-head">` + escape(s) + `</span>`
}
