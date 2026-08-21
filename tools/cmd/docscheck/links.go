package main

import (
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"unicode"
)

// checkLinks holds every Markdown file of the repository to the files and the
// headings it points at.
//
// THE SITE ALREADY CHECKS THIS AND CANNOT CHECK IT HERE. Its gate walks the
// RENDERED pages, and the four files a reader meets first — the README in both
// languages, CONTRIBUTING, ROADMAP — are never rendered: they are read on
// GitHub, out of the repository. So they wanted the one gate every doc page
// has and had none, and it showed the moment the base language flipped: three
// of their fragments named Russian headings on pages that had been English for
// four commits, and one link named a file that had been renamed. Nothing was
// broken visibly — a dead fragment lands the reader at the top of the page,
// and a dead file link is a 404 nobody clicks in their own repository.
//
// The anchor is slugged BY GITHUB'S RULES rather than by the site's, and the
// difference is not academic: the site TRANSLITERATES a Russian heading into
// Latin letters, GitHub keeps the Cyrillic. These files are read where GitHub
// renders them, so GitHub's rules are the ones that decide.
func checkLinks(root string) []string {
	files := map[string]bool{}
	heads := map[string]map[string]bool{}

	filepath.WalkDir(root, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		if d.IsDir() {
			switch d.Name() {
			case "node_modules", ".git", "dist":
				return filepath.SkipDir
			}
			return nil
		}
		rel, _ := filepath.Rel(root, p)
		rel = filepath.ToSlash(rel)
		files[rel] = true
		if !strings.HasSuffix(rel, ".md") {
			return nil
		}
		b, err := os.ReadFile(p)
		if err != nil {
			return nil
		}
		ids := map[string]bool{}
		for _, line := range strings.Split(string(b), "\n") {
			line = strings.TrimRight(line, "\r")
			if !strings.HasPrefix(line, "#") {
				continue
			}
			t := strings.TrimLeft(line, "#")
			if t == line || (t != "" && !strings.HasPrefix(t, " ")) {
				continue // "#!/…" and the like are not headings
			}
			ids[ghSlug(strings.TrimSpace(t))] = true
		}
		heads[rel] = ids
		return nil
	})

	var out []string
	for rel := range files {
		if !strings.HasSuffix(rel, ".md") {
			continue
		}
		b, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(rel)))
		if err != nil {
			continue
		}
		// AN ANCHOR IS ONLY CHECKED WHERE GITHUB IS THE READER. On a page of
		// the reference the same fragment is right by a different rule: the
		// site transliterates a Russian heading and rewrites the link to match,
		// so a Cyrillic fragment written in the markdown is correct there and
		// would be a false finding here. Those pages have a gate of their own,
		// over the rendered HTML, where the site's own rule applies.
		//
		// A link to a FILE is checked everywhere: it is the same fact under
		// both rules, and a renamed file is how `docs/index.en.md` stayed in
		// the roadmap after the page had moved.
		rendered := rendersOnSite(rel)
		for _, m := range mdLinkRe.FindAllStringSubmatch(string(b), -1) {
			dest, frag := m[1], ""
			if i := strings.IndexByte(dest, '#'); i >= 0 {
				dest, frag = dest[:i], dest[i+1:]
			}
			if dest == "" {
				if !rendered && !heads[rel][frag] {
					out = append(out, rel+"  a heading of its own that is not there: #"+frag)
				}
				continue
			}
			target := path.Clean(path.Join(path.Dir(rel), dest))
			if !files[target] {
				out = append(out, rel+"  a link to a file that is not there: "+dest)
				continue
			}
			if !rendered && frag != "" && heads[target] != nil && !heads[target][frag] {
				out = append(out, rel+"  a link into a heading that is not there: "+dest+"#"+frag)
			}
		}
	}
	sort.Strings(out)
	return out
}

// rendersOnSite reports whether a file becomes a page of the reference. Those
// pages are checked by the site build, with the site's own rule for an anchor;
// everything else is read on GitHub and is checked here.
func rendersOnSite(rel string) bool {
	return strings.HasPrefix(rel, "docs/") &&
		rel != "docs/README.md" &&
		!strings.HasPrefix(rel, "docs/internal/")
}

// A relative link to a file of the repository. Addresses with a scheme and
// absolute ones are somebody else's business; an image is skipped by the `[^!]`
// the caller cannot express, so the leading bracket is matched exactly.
var mdLinkRe = regexp.MustCompile(`\]\(([^)\s]+\.md(?:#[^)\s]*)?|#[^)\s]+)\)`)

// ghSlug is GitHub's rule for turning a heading into an anchor: lower-case,
// spaces become dashes, everything that is neither a letter nor a digit nor a
// dash is dropped. Letters of ANY alphabet survive — which is the whole
// difference from the site's slug, and the reason this function exists rather
// than reusing that one.
func ghSlug(s string) string {
	var b strings.Builder
	for _, r := range strings.ToLower(s) {
		switch {
		case r == ' ':
			b.WriteByte('-')
		case r == '-' || r == '_' || unicode.IsLetter(r) || unicode.IsDigit(r):
			b.WriteRune(r)
		}
	}
	return b.String()
}

func reportLinks(root string) bool {
	bad := checkLinks(root)
	if len(bad) == 0 {
		return false
	}
	fmt.Printf("── links that lead nowhere (%d) ──\n", len(bad))
	for _, s := range bad {
		fmt.Println("  " + s)
	}
	fmt.Println()
	return true
}
