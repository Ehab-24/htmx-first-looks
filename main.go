package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"

	// "github.com/Ehab-24/test/scripts"
	"github.com/Ehab-24/test/scripts"
	"github.com/go-chi/chi/v5"
)

type RootLayoutData struct {
	Title       string
	Description string
	Scripts     template.HTML
}

type HomePageData struct {
	RootLayoutData
}

type AboutPageData struct {
	RootLayoutData
}

type ArticlePageData struct {
	RootLayoutData
	HTML template.HTML
}

func checkAPIError(w http.ResponseWriter, err error) bool {
	if err != nil {
		log.Println("~ api error:", err)
		http.Error(w, "unexpected error", http.StatusInternalServerError)
		return true
	}
	return false
}
func main() {

	log.SetFlags(0)

	cachedArticles := scripts.CacheArticles()
	r := chi.NewRouter()
	r.Get("/endpoint", func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.Header)
		w.Write([]byte("Clicked!"))
	})

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		t, err := template.ParseFiles("index.html", "views/home.html", "components/button.html")
		if checkAPIError(w, err) {
			return
		}

		data := HomePageData{
			RootLayoutData: RootLayoutData{
				Title:       "Ehab Sohail",
				Description: "I am an aspiring web developer from Pakistan",
				Scripts:     HTMLScripts[TailwindCSS] + HTMLScripts[HTMX],
			},
		}
		if checkAPIError(w, err) {
			return
		}
		checkAPIError(w, t.Execute(w, data))
	})

	r.Get("/about", func(w http.ResponseWriter, r *http.Request) {
		t, err := template.ParseFiles("index.html", "views/about.html")
		if checkAPIError(w, err) {
			return
		}
		data := AboutPageData{
			RootLayoutData: RootLayoutData{
				Title:       "About | Ehab Sohail",
				Description: "I am an aspiring web developer from Pakistan",
			},
		}
		checkAPIError(w, t.Execute(w, data))
	})

	r.Get("/editor", func(w http.ResponseWriter, r *http.Request) {
		t, err := template.ParseFiles("index.html", "viewsabout.html")
		if checkAPIError(w, err) {
			return
		}
		data := AboutPageData{
			RootLayoutData: RootLayoutData{
				Title:       "Write | Art",
				Description: "Write an article",
				Scripts:     HTMLScripts[SimpleMDE] + HTMLScripts[HTMX] + HTMLScripts[TailwindCSS],
			},
		}
		checkAPIError(w, t.Execute(w, data))
	})

	r.Get("/articles/{slug}", func(w http.ResponseWriter, r *http.Request) {
		slug := chi.URLParam(r, "slug")
		html := template.HTML(cachedArticles[slug])
		data := ArticlePageData{
			RootLayoutData: RootLayoutData{
				Title:       strings.ToTitle(slug),
				Description: fmt.Sprintf("some random description for %s", slug),
			},
			HTML: html,
		}

		t, err := template.ParseFiles("index.html", "views/article.html")
		if checkAPIError(w, err) {
			return
		}

		err = t.Execute(w, data)
		checkAPIError(w, err)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	http.ListenAndServe(":"+port, r)
}
