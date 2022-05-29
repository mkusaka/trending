package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"github.com/gorilla/feeds"
	"github.com/mmcdole/gofeed"
	_ "github.com/motemen/go-loghttp/global"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
)

const baseUrl = "https://mshibanami.github.io/GitHubTrendingRSS"

func main() {
	err := run()
	if err != nil {
		fmt.Printf("%+v", err)
		os.Exit(1)
	}
	os.Exit(0)
}

type UrlAndDirectory struct {
	Url       string
	Directory string
}

func run() error {
	periodTypes := []string{"daily", "weekly", "monthly"}
	//periodTypes := []string{"weekly"}
	languageTypes := []string{"all", "unknown", "c++", "html", "java", "javascript", "php", "python", "ruby", "1c-enterprise", "2-dimensional-array", "4d", "abap", "abap-cds", "abnf", "actionscript", "ada", "adobe-font-metrics", "agda", "ags-script", "aidl", "al", "alloy", "alpine-abuild", "altium-designer", "ampl", "angelscript", "ant-build-system", "antlr", "apacheconf", "apex", "api-blueprint", "apl", "apollo-guidance-computer", "applescript", "arc", "asciidoc", "asl", "asn.1", "classic-asp", "asp.net", "aspectj", "assembly", "astro", "asymptote", "ats", "augeas", "autohotkey", "autoit", "avro-idl", "awk", "ballerina", "basic", "batchfile", "beef", "befunge", "berry", "bibtex", "bicep", "bison", "bitbake", "blade", "blitzbasic", "blitzmax", "bluespec", "boo", "boogie", "brainfuck", "brightscript", "zeek", "browserslist", "c", "c%23", "c-objdump", "c2hs-haskell", "cabal-config", "cadence", "cairo", "cap'n-proto", "cartocss", "ceylon", "chapel", "charity", "chuck", "cil", "cirru", "clarion", "clarity", "clean", "click", "clips", "clojure", "closure-templates", "cloud-firestore-security-rules", "cmake", "cobol", "codeowners", "codeql", "coffeescript", "coldfusion", "coldfusion-cfc", "collada", "common-lisp", "common-workflow-language", "component-pascal", "conll-u", "cool", "coq", "cpp-objdump", "creole", "crystal", "cson", "csound", "csound-document", "csound-score", "css", "csv", "cuda", "cue", "curl-config", "curry", "cweb", "cycript", "cython", "d", "d-objdump", "dafny", "darcs-patch", "dart", "dataweave", "debian-package-control-file", "denizenscript", "desktop", "dhall", "diff", "digital-command-language", "dircolors", "directx-3d-file", "dm", "dns-zone", "dockerfile", "dogescript", "dtrace", "dylan", "e", "e-mail", "eagle", "earthly", "easybuild", "ebnf", "ec", "ecere-projects", "ecl", "eclipse", "editorconfig", "edje-data-collection", "edn", "eiffel", "ejs", "elixir", "elm", "emacs-lisp", "emberscript", "eq", "erlang", "euphoria", "f%23", "f*", "factor", "fancy", "fantom", "faust", "fennel", "figlet-font", "filebench-wml", "filterscript", "fish", "fluent", "flux", "formatted", "forth", "fortran", "fortran-free-form", "freebasic", "freemarker", "frege", "futhark", "g-code", "game-maker-language", "gaml", "gams", "gap", "gcc-machine-description", "gdb", "gdscript", "gedcom", "gemfile.lock", "genero", "genero-forms", "genie", "genshi", "gentoo-ebuild", "gentoo-eclass", "gerber-image", "gettext-catalog", "gherkin", "git-attributes", "git-config", "gleam", "glsl", "glyph", "glyph-bitmap-distribution-format", "gn", "gnuplot", "go", "go-checksums", "go-module", "golo", "gosu", "grace", "gradle", "grammatical-framework", "graph-modeling-language", "graphql", "graphviz-(dot)", "groovy", "groovy-server-pages", "gsc", "hack", "haml", "handlebars", "haproxy", "harbour", "haskell", "haxe", "hcl", "hiveql", "hlsl", "holyc", "hoon", "jinja", "html+ecr", "html+eex", "html+erb", "html+php", "html+razor", "http", "hxml", "hy", "hyphy", "idl", "idris", "ignore-list", "igor-pro", "imagej-macro", "inform-7", "ini", "inno-setup", "io", "ioke", "irc-log", "isabelle", "isabelle-root", "j", "janet", "jar-manifest", "jasmin", "java-properties", "java-server-pages", "javascript+erb", "jest-snapshot", "jflex", "jison", "jison-lex", "jolie", "jq", "json", "json-with-comments", "json5", "jsoniq", "jsonld", "jsonnet", "julia", "jupyter-notebook", "kaitai-struct", "kakounescript", "kicad-layout", "kicad-legacy-layout", "kicad-schematic", "kit", "kotlin", "krl", "kusto", "kvlang", "labview", "lark", "lasso", "latte", "lean", "less", "lex", "lfe", "ligolang", "lilypond", "limbo", "linker-script", "linux-kernel-module", "liquid", "literate-agda", "literate-coffeescript", "literate-haskell", "livescript", "llvm", "logos", "logtalk", "lolcode", "lookml", "loomscript", "lsl", "ltspice-symbol", "lua", "m", "m4", "m4sugar", "macaulay2", "makefile", "mako", "markdown", "marko", "mask", "mathematica", "matlab", "maven-pom", "max", "maxscript", "mcfunction", "wikitext", "mercury", "meson", "metal", "microsoft-developer-studio-project", "microsoft-visual-studio-solution", "minid", "mint", "mirah", "mirc-script", "mlir", "modelica", "modula-2", "modula-3", "module-management-system", "monkey", "monkey-c", "moocode", "moonscript", "motoko", "motorola-68k-assembly", "mql4", "mql5", "mtml", "muf", "mupad", "muse", "mustache", "myghty", "nanorc", "nasl", "ncl", "nearley", "nemerle", "neon", "nesc", "netlinx", "netlinx+erb", "netlogo", "newlisp", "nextflow", "nginx", "nim", "ninja", "nit", "nix", "nl", "npm-config", "nsis", "nu", "numpy", "nunjucks", "nwscript", "objdump", "object-data-instance-notation", "objective-c", "objective-c++", "objective-j", "objectscript", "ocaml", "odin", "omgrofl", "ooc", "opa", "opal", "open-policy-agent", "opencl", "openedge-abl", "openqasm", "openrc-runscript", "openscad", "openstep-property-list", "opentype-feature-file", "org", "ox", "oxygene", "oz", "p4", "pan", "papyrus", "parrot", "parrot-assembly", "parrot-internal-representation", "pascal", "pawn", "peg.js", "pep8", "perl", "pic", "pickle", "picolisp", "piglatin", "pike", "plantuml", "plpgsql", "plsql", "pod", "pod-6", "pogoscript", "pony", "postcss", "postscript", "pov-ray-sdl", "powerbuilder", "powershell", "prisma", "processing", "procfile", "proguard", "prolog", "promela", "propeller-spin", "protocol-buffer", "protocol-buffer-text-format", "public-key", "pug", "puppet", "pure-data", "purebasic", "purescript", "python-console", "python-traceback", "q", "q%23", "qmake", "qml", "qt-script", "quake", "r", "racket", "ragel", "raku", "raml", "rascal", "raw-token-data", "rdoc", "readline-config", "realbasic", "reason", "rebol", "record-jar", "red", "redcode", "redirect-rules", "regular-expression", "ren'py", "renderscript", "rescript", "restructuredtext", "rexx", "rich-text-format", "ring", "riot", "rmarkdown", "robotframework", "robots.txt", "roff", "roff-manpage", "rouge", "rpc", "rpgle", "rpm-spec", "runoff", "rust", "sage", "saltstack", "sas", "sass", "scala", "scaml", "scheme", "scilab", "scss", "sed", "self", "selinux-policy", "shaderlab", "shell", "shellcheck-config", "shellsession", "shen", "sieve", "singularity", "slash", "slice", "slim", "smali", "smalltalk", "smarty", "smpl", "smt", "solidity", "soong", "sourcepawn", "sparql", "spline-font-database", "sqf", "sql", "sqlpl", "squirrel", "srecode-template", "ssh-config", "stan", "standard-ml", "starlark", "stata", "ston", "stringtemplate", "stylus", "subrip-text", "sugarss", "supercollider", "svelte", "svg", "swift", "swig", "systemverilog", "talon", "tcl", "tcsh", "tea", "terra", "tex", "texinfo", "text", "textile", "textmate-properties", "thrift", "ti-program", "tla", "toml", "tsql", "tsv", "tsx", "turing", "turtle", "twig", "txl", "type-language", "typescript", "unified-parallel-c", "unity3d-asset", "unix-assembly", "uno", "unrealscript", "urweb", "v", "vala", "valve-data-format", "vba", "vbscript", "vcl", "verilog", "vhdl", "vim-help-file", "vim-script", "vim-snippet", "visual-basic-.net", "volt", "vue", "vyper", "wavefront-material", "wavefront-object", "wdl", "web-ontology-language", "webassembly", "webidl", "webvtt", "wget-config", "windows-registry-entries", "wisp", "witcher-script", "wollok", "world-of-warcraft-addon-data", "x-bitmap", "x-font-directory-index", "x-pixmap", "x10", "xbase", "xc", "xcompose", "xml", "xml-property-list", "xojo", "xonsh", "xpages", "xproc", "xquery", "xs", "xslt", "xtend", "yacc", "yaml", "yang", "yara", "yasnippet", "zap", "zenscript", "zephir", "zig", "zil", "zimpl"}
	//languageTypes := []string{"typescript"}

	var urlAndDirectories []*UrlAndDirectory
	for _, periodType := range periodTypes {
		for _, languageType := range languageTypes {
			urlAndDirectories = append(urlAndDirectories, &UrlAndDirectory{
				Url:       fmt.Sprintf("%s/%s/%s.xml", baseUrl, periodType, languageType),
				Directory: path.Join("out", languageType, periodType),
			})
		}
	}

	eg := errgroup.Group{}
	for _, urlAndDirectory := range urlAndDirectories {
		und := urlAndDirectory
		eg.Go(func() error {
			fp := gofeed.NewParser()
			parsedFeed, err := fp.ParseURL(und.Url)
			if err != nil {
				fmt.Printf("failed to parse: %s with: %+v\n", und.Url, errors.WithStack(err))
				return nil

				// some feed sometime failed..
				//return errors.WithStack(err)
			}

			feed := &feeds.Feed{
				Title:       parsedFeed.Title,
				Description: parsedFeed.Description,
				//Updated:     *parsedFeed.UpdatedParsed,
				Created: *parsedFeed.PublishedParsed,
				Link: &feeds.Link{
					Href: parsedFeed.Link,
				},
			}

			for _, item := range parsedFeed.Items {
				feed.Items = append(feed.Items, &feeds.Item{
					Title: item.Title,
					Link: &feeds.Link{
						Href: item.Link,
					},
					Description: item.Description,
					Updated:     *parsedFeed.PublishedParsed,
					Created:     *parsedFeed.PublishedParsed,
				})
			}
			f, err := feed.ToAtom()
			if err != nil {
				return errors.WithStack(err)
			}

			fmt.Printf("create directory %s start...\n", und.Directory)
			err = os.MkdirAll(und.Directory, os.ModePerm)
			fmt.Printf("create directory %s done...\n", und.Directory)
			if err != nil {
				return errors.WithStack(err)
			}

			file := path.Join(und.Directory, "index.xml")
			fmt.Printf("write file %s start...\n", file)
			err = ioutil.WriteFile(file, []byte(f), os.ModePerm)
			fmt.Printf("write file %s done...\n", file)

			if err != nil {
				return errors.WithStack(err)
			}
			return nil
		})
		if err := eg.Wait(); err != nil {
			return err
		}
	}

	return nil
}
