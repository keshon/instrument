// go run ./cmd/site                 build into ../dist
// go run ./cmd/site -serve :4321    build and raise a server
package main

import (
	"flag"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"instrument/site/internal/check"
	"instrument/site/internal/content"
	"instrument/site/internal/i18n"
	"instrument/site/internal/nav"
	"instrument/site/internal/render"
)

func main() {
	var (
		docs   = flag.String("docs", "../docs", "directory holding the page sources")
		out    = flag.String("out", "dist", "output directory")
		kit    = flag.String("kit", "../src", "directory of the kit")
		assets = flag.String("assets", "../assets", "directory of the kit's assets")
		public = flag.String("public", "public", "directory copied into the output root as it is")
		serve  = flag.String("serve", "", "raise a server after the build, for example :4321")

		registry = flag.String("registry", "../components.json", "component registry: the relations behind the Related section")

		verbose = flag.Bool("contract", false, "print contract remarks for pages not yet moved over")
	)
	flag.Parse()

	byLang, err := content.Collect(*docs)
	if err != nil {
		log.Fatalf("collecting pages: %v", err)
	}
	pages := content.Flat(byLang)
	if len(pages) == 0 {
		log.Fatalf("not one page found in %s", *docs)
	}

	sprite, err := os.ReadFile(filepath.Join(*assets, "sprite.svg"))
	if err != nil {
		log.Fatalf("cannot read the sprite: %v", err)
	}
	var missing []string
	for _, p := range pages {
		if p.Icon != "" && !strings.Contains(string(sprite), `id="`+p.Icon+`"`) {
			missing = append(missing, fmt.Sprintf("%s  no symbol %s in the sprite", p.Route, p.Icon))
		}
	}

	tokens, err := content.TokenValues(*kit)
	if err != nil {
		log.Fatalf("cannot read the kit's tokens: %v", err)
	}
	content.ResolveTokens(pages, tokens)

	content.SetSprite(string(sprite))
	rel, err := content.LoadRelations(*registry)
	if err != nil {
		log.Fatalf("cannot read the registry: %v", err)
	}
	content.SetRelations(rel)
	if err := content.Render(pages); err != nil {
		log.Fatalf("rendering: %v", err)
	}

	contractErrs, contractWarns := check.Contract(pages)

	styles, err := render.Stylesheets()
	if err != nil {
		log.Fatalf("cannot read the site's styles: %v", err)
	}
	assetErrs := check.Assets(styles, tokens)
	assetErrs = append(assetErrs, check.Base(tokens)...)

	sources := map[string]string{}
	for n, s := range styles {
		sources["site/"+n] = s
	}
	// The comment rule is one for the whole repository, so the gate's zone is
	// one as well. While `tools` and `cmd` were outside it, the rule rested on
	// attentiveness in the very place the checks themselves live: 4 384 lines
	// of Go and two commands in JS stayed out, and a chronicle line about all
	// three sizes once taking --radius-md went through in silence.
	//
	// The tools directory is looked for beside the kit rather than given by a
	// flag: it lies at the same level as `src`, and a second flag would drift
	// apart from the first.
	tools := filepath.Join(filepath.Dir(*kit), "tools")
	for _, dir := range []string{*kit, "internal", "cmd", tools} {
		// A missing directory is not an error: the site builds from a module
		// of its own, and outside its tree there may be nothing at all.
		if _, err := os.Stat(dir); err != nil {
			continue
		}
		if err := collectSources(sources, dir); err != nil {
			log.Fatalf("cannot read the sources: %v", err)
		}
	}
	assetErrs = append(assetErrs, check.Comments(sources)...)
	assetErrs = append(assetErrs, check.StrayCommentEnd(sources)...)

	problems := check.Verify(pages, string(sprite))
	problems = append(problems, missing...)
	problems = append(problems, contractErrs...)
	problems = append(problems, assetErrs...)
	problems = append(problems, nav.Check(byLang[i18n.Base])...)
	if len(problems) > 0 {
		for _, p := range problems {
			fmt.Fprintln(os.Stderr, "  "+p)
		}
		log.Fatalf("build stopped: %d problems", len(problems))
	}

	sections := map[i18n.Lang][]nav.Section{}
	for lang, ps := range byLang {
		sections[lang] = nav.Build(lang, ps)
	}
	if err := render.Site(byLang, sections, render.Options{
		Out: *out, Kit: *kit, Assets: *assets, Public: *public,
	}); err != nil {
		log.Fatalf("building: %v", err)
	}

	demos := map[string]bool{}
	for _, p := range pages {
		for _, d := range p.Demos {
			demos[d.ID] = true
		}
	}
	fmt.Printf("pages: %d  ·  live examples: %d  ·  navigation sections: %d\n",
		len(byLang[i18n.Base]), len(demos), len(sections[i18n.Base]))
	onContract := 0
	for _, p := range byLang[i18n.Base] {
		if p.Layout == "component" {
			onContract++
		}
	}
	fmt.Printf("pages under the contract: %d  ·  awaiting the move: %d remarks\n",
		onContract, len(contractWarns))
	if *verbose {
		for _, w := range contractWarns {
			fmt.Fprintln(os.Stderr, "  ~ "+w)
		}
	}

	if *serve != "" {
		fmt.Printf("server: http://localhost%s\n", *serve)
		log.Fatal(http.ListenAndServe(*serve, http.FileServer(http.Dir(*out))))
	}
}

func collectSources(out map[string]string, dir string) error {
	return filepath.WalkDir(dir, func(p string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		// `.mjs` is on the list because the runner of the pixel check is
		// written that way, and without the extension it would be the one file
		// in the repository whose comments nobody guards.
		switch filepath.Ext(p) {
		case ".css", ".js", ".mjs", ".go":
		default:
			return nil
		}
		b, err := os.ReadFile(p)
		if err != nil {
			return err
		}
		out[filepath.ToSlash(p)] = string(b)
		return nil
	})
}
