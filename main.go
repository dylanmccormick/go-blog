package main

import (
	"bufio"
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"time"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/parser"
	"go.abhg.dev/goldmark/frontmatter"
)

const PREVIEW_CHARACTER_LIMIT = 300

var (
	postsStore = make(map[string]MarkdownPost)
	postsList  = make([]MarkdownPost, 0)
	templates  *template.Template
)

func init() {
	templates = template.Must(template.ParseGlob("./templates/*.html"))
}

func loadPosts() {
	postsDir := os.Getenv("POSTS_DIR")
	if postsDir == "" {
		log.Println("Didn't get POSTS_DIR environment variable. Using default")
		postsDir = "./posts"
	}
	files, err := filepath.Glob(postsDir)
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
		markdown.Slug = slug
		markdown.FormattedDate = markdown.OriginalDate.Format("January 2, 2006")
		markdown.Preview = generatePreview(file)

		postsStore[slug] = markdown
		postsList = append(postsList, markdown)
	}
	slices.SortFunc(postsList, func(a, b MarkdownPost) int {
		return a.OriginalDate.Compare(b.OriginalDate) * -1
	})
}

func generatePreview(file []byte) string {
	var builder strings.Builder
	frontmatter := false
	reader := bytes.NewReader(file)
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.EqualFold(line, "---") {
			frontmatter = !frontmatter
			continue
		}

		if frontmatter {
			continue
		}

		if strings.HasPrefix(line, "# ") {
			continue
		}
		words := strings.SplitSeq(line, " ")
		for word := range words {
			if builder.Len() < PREVIEW_CHARACTER_LIMIT {
				builder.WriteString(" ")
				builder.WriteString(word)
			} else {
				break
			}
		}
		if builder.Len() >= PREVIEW_CHARACTER_LIMIT {
			break
		}

	}

	return builder.String() + "…"
}

func main() {
	loadPosts()
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/blog/{slug}", blogHandler)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	fmt.Println("Listening on :3000")
	http.ListenAndServe(":3000", nil)
}

type MarkdownPost struct {
	Title         string    `yaml:"title"`
	OriginalDate  time.Time `yaml:"original_date"`
	LastUpdated   string    `yaml:"last_updated"`
	Body          template.HTML
	Slug          string
	FormattedDate string
	Preview       string
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	err := templates.ExecuteTemplate(w, "home.html", postsList)
	if err != nil {
		fmt.Println("Uh oh couldn't write post")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func blogHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("slug: %s\n", r.PathValue("slug"))
	slug := r.PathValue("slug")
	mdPost, ok := postsStore[slug]
	if !ok {
		fmt.Println("post not found")
		http.NotFound(w, r)
		return
	}
	log.Printf("%#v\n", mdPost.OriginalDate)

	err := templates.ExecuteTemplate(w, "blog.html", mdPost)
	if err != nil {
		fmt.Println("Uh oh couldn't write post")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
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
	meta.Body = template.HTML(buf.String())

	return meta
}
