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

// The `do` / `dont` pair USED to live here, and its removal is the mark of a
// step being finished rather than a loss of an affordance.
//
// It rendered a two-column row of banners: this way, not that way. Two hundred
// and sixty-four of them stood on sixty-seven pages, five thousand eight
// hundred words, and a measurement took them apart: two thirds were pure
// advice naming neither attribute nor class, a quarter compared the component
// with a NEIGHBOUR, and the rest stated a requirement of the markup. The
// comparison had already moved to the section index when the page shape
// changed — checked by name, forty-five links of forty-seven were there. The
// requirements moved into the contract table, and nineteen of twenty-eight
// turned out to be standing in it already.
//
// What was left was teaching how to use the thing, and a catalogue does not
// teach. The directive goes with the last of its blocks: a renderer that can
// still draw something nobody writes is dead code somebody will one day take
// for live. Nothing guards this — nothing has to. An unknown ::: block is a
// build error, so the day someone writes :::do again the site will say so.

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
			case noteTones[kind] != "":
				// A note has NO heading, and text on the opening line used to be
				// dropped without a word. It was harmless while `do` and `dont`
				// took a heading and the habit of writing one existed; now that
				// they are gone, the habit outlives them, and silently eating what
				// the author typed is the one thing this build refuses to do.
				if head != "" {
					return nil, fmt.Errorf("строка %d: у :::%s нет заголовка, а в строке открытия стоит %q — этот текст пропал бы молча",
						i+1, kind, head)
				}
				closeRow()
				out = append(out, "",
					fmt.Sprintf(`<div class="inst-note" data-tone="%s"><div class="note-body">`, noteTones[kind]),
					"")
				closes = "</div></div>"
			default:
				return nil, fmt.Errorf("строка %d: неизвестный блок :::%s (есть note · warn · danger · ok)", i+1, kind)
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
			closeRow()
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
