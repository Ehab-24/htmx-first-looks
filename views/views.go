package views

import (
	"context"
	"html/template"
	"net/http"
	"time"

	"github.com/Ehab-24/test/api"
	"github.com/Ehab-24/test/db"
	"github.com/Ehab-24/test/types"
	"github.com/go-chi/chi/v5"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type RootLayoutData struct {
	Title       string
	Description string
	HTMLScripts template.HTML
}

type HomePageData struct {
	RootLayoutData
}

type RegisterPageData struct {
	RootLayoutData
}

type AboutPageData struct {
	RootLayoutData
}

type ArticlesPageData struct {
	RootLayoutData
	Articles []types.Article
}

type ArticlePageData struct {
	RootLayoutData
	HTML template.HTML
}

func GetLocalRouter() http.Handler {
	r := chi.NewRouter()

	r.Get("/", homePage)
	r.Get("/about", aboutPage)

	r.Get("/editor", editorPage)

	r.Get("/articles/{slug}", articlePage)
	r.Get("/articles", articlesPage)

	r.Get("/register", registerPage)
	r.Get("/login", loginPage)
	return r
}

func homePage(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("index.html", "views/home.html", "components/button.html")
	if api.CheckAPIError(w, err) {
		return
	}

	data := HomePageData{
		RootLayoutData: RootLayoutData{
			Title:       "Art",
			Description: "A simple blog",
			HTMLScripts: htmlScripts[TailwindCSS] + htmlScripts[HTMX],
		},
	}
	if api.CheckAPIError(w, err) {
		return
	}
	api.CheckAPIError(w, t.Execute(w, data))
}

func registerPage(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("index.html", "views/register.html")
	if api.CheckAPIError(w, err) {
		return
	}
	data := RegisterPageData{
		RootLayoutData: RootLayoutData{
			Title:       "Register | Art",
			Description: "Register",
			HTMLScripts: htmlScripts[TailwindCSS] + htmlScripts[HTMX],
		},
	}
	api.CheckAPIError(w, t.Execute(w, data))
}

func loginPage(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("index.html", "views/login.html")
	if api.CheckAPIError(w, err) {
		return
	}
	data := RegisterPageData{
		RootLayoutData: RootLayoutData{
			Title:       "Login | Art",
			Description: "Login",
			HTMLScripts: htmlScripts[TailwindCSS] + htmlScripts[HTMX],
		},
	}
	api.CheckAPIError(w, t.Execute(w, data))
}

func editorPage(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("index.html", "views/editor.html")
	if api.CheckAPIError(w, err) {
		return
	}
	data := AboutPageData{
		RootLayoutData: RootLayoutData{
			Title:       "Editor | Art",
			Description: "Write an article",
			HTMLScripts: htmlScripts[SimpleMDE] + htmlScripts[HTMX],
		},
	}
	api.CheckAPIError(w, t.Execute(w, data))
}

func aboutPage(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("index.html", "views/about.html")
	if api.CheckAPIError(w, err) {
		return
	}
	data := AboutPageData{
		RootLayoutData: RootLayoutData{
			Title:       "About | Art",
			Description: "I am an aspiring web developer from Pakistan",
			HTMLScripts: htmlScripts[TailwindCSS],
		},
	}
	api.CheckAPIError(w, t.Execute(w, data))
}

func articlesPage(w http.ResponseWriter, r *http.Request) {
	funcMap := template.FuncMap{
		"formatDate": func(date time.Time) string {
			return date.Format("January 2, 2006")
		},
	}
	t, err := template.New("index.html").Funcs(funcMap).ParseFiles("index.html", "views/articles.html")
	if api.CheckAPIError(w, err) {
		return
	}

	opts := options.Find().SetProjection(bson.D{{Key: "title", Value: 1}, {Key: "slug", Value: 1}, {Key: "description", Value: 1}, {Key: "created_at", Value: 1}})
	filter := bson.D{}
	_articles, err := db.Db.Collection("posts").Find(context.Background(), filter, opts)
	if api.CheckAPIError(w, err) {
		return
	}
	var articles []types.Article
	if api.CheckAPIError(w, _articles.All(context.Background(), &articles)) {
		return
	}

	data := ArticlesPageData{
		RootLayoutData: RootLayoutData{
			Title:       "Articles | Art",
			Description: "Articles",
			HTMLScripts: htmlScripts[TailwindCSS],
		},
		Articles: articles,
	}
	api.CheckAPIError(w, t.Execute(w, data))
}

func articlePage(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	_article := db.Db.Collection("posts").FindOne(context.Background(), bson.D{{Key: "slug", Value: slug}})
	var article types.Article
	if api.CheckAPIError(w, _article.Decode(&article)) {
		return
	}

	data := ArticlePageData{
		RootLayoutData: RootLayoutData{
			Title:       article.Title,
			Description: article.Description,
			HTMLScripts: htmlScripts[TailwindCSS],
		},
		HTML: template.HTML(article.Content),
	}

	t, err := template.ParseFiles("index.html", "views/article.html")
	if api.CheckAPIError(w, err) {
		return
	}

	err = t.Execute(w, data)
	api.CheckAPIError(w, err)
}
