package main

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/parser"
	"go.abhg.dev/goldmark/frontmatter"
)

var postsStore = make(map[string]MarkdownPost)

func loadPosts() {
	files, err := filepath.Glob("./posts/*")
	if err != nil {
		log.Fatal("poop you couldn't load any files. Skill issue")
	}
	for _, post := range files {
		fileName := filepath.Base(post)
		slug := strings.TrimSuffix(fileName, ".md")
		file, err := os.ReadFile(post)
		if err != nil {
			log.Fatal(err)
		}
		markdown := MarkdownToHtml(string(file))

		postsStore[slug] = markdown
	}
}

func main() {
	loadPosts()
	http.HandleFunc("/hello", handler)
	http.HandleFunc("/blog/{slug}", blogHandler)
	fmt.Println("Listening on :3000")
	http.ListenAndServe(":3000", nil)
}

type MarkdownPost struct {
	Title        string `yaml:"title"`
	OriginalDate string `yaml:"original_date"`
	LastUpdated  string `yaml:"last_updated"`
	Body         template.HTML
}

func blogHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("slug: %s", r.PathValue("slug"))
	slug := r.PathValue("slug")
	mdPost, ok := postsStore[slug]
	if !ok {
		fmt.Println("post not found")
		http.NotFound(w, r)
		return
	}
	const tpl = `
	<!DOCTYPE html>
	<html>
		<title>
			{{.Title}}
		</title>
		<body>
			<div>
				{{.Body}}
			<div>
			<div>
				Last updated:{{.LastUpdated}}
			</div>
		</body>
	</html>
	`

	t, err := template.New("blogpage").Parse(tpl)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	err = t.Execute(w, mdPost)
}

func handler(w http.ResponseWriter, r *http.Request) {
	file, err := os.ReadFile("./hello_world.md")
	if err != nil {
		log.Fatal(err)
	}
	const tpl = `
	<!DOCTYPE html>
	<html>
		<title>
			{{.Title}}
		</title>
		<body>
			<div>
				{{.Body}}
			<div>
			<div>
				Last updated:{{.LastUpdated}}
			</div>
		</body>
	</html>
	`

	t, err := template.New("blogpage").Parse(tpl)
	markdown := MarkdownToHtml(string(file))

	err = t.Execute(w, markdown)
}

func MarkdownToHtml(content string) MarkdownPost {
	var buf bytes.Buffer
	meta := MarkdownPost{}
	md := goldmark.New(
		goldmark.WithExtensions(
			&frontmatter.Extender{},
		),
	)

	ctx := parser.NewContext()

	if err := md.Convert([]byte(content), &buf, parser.WithContext(ctx)); err != nil {
		log.Fatal(err)
	}
	fm := frontmatter.Get(ctx)
	err := fm.Decode(&meta)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%#v", meta)
	meta.Body = template.HTML(buf.String())

	return meta
}
