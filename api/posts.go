package api

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/Ehab-24/test/api/middleware"
	"github.com/Ehab-24/test/db"
	"github.com/Ehab-24/test/lib"
	"github.com/Ehab-24/test/types"
	"github.com/go-chi/chi/v5"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func getPostsRouter() http.Handler {
	r := chi.NewRouter()
	r.With(middleware.VerifyUser).Post("/", createPost)
	r.Get("/", getPosts)

	return r
}

func createPost(w http.ResponseWriter, r *http.Request) {
	c := r.FormValue("content")
	html := lib.MdToHTML([]byte(c))

	authorID, err := primitive.ObjectIDFromHex(r.Context().Value("user_id").(string))
	if CheckAPIError(w, err) {
		return
	}
	article := types.Article{
		ID:          primitive.NewObjectID(),
		Title:       r.FormValue("title"),
		Slug:        lib.KebabCase(r.FormValue("title")),
		Description: r.FormValue("description"),
		Content:     string(html),
		Author:      authorID,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	if _, err := db.Db.Collection("posts").InsertOne(context.Background(), article); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write([]byte("success"))
}

func getPosts(w http.ResponseWriter, r *http.Request) {

	filter := &types.ArticleFilter{
		Title:       r.URL.Query().Get("title"),
		Description: r.URL.Query().Get("desctription"),
		Content:     r.URL.Query().Get("content"),
		Author:      r.URL.Query().Get("author"),
	}

	l, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if CheckAPIErrorWithStatus(w, err, "invalid query parameter 'limit'", http.StatusBadRequest) {
		return
	}

	s, err := strconv.Atoi(r.URL.Query().Get("skip"))
	if CheckAPIErrorWithStatus(w, err, "invalid query parameter 'skip'", http.StatusBadRequest) {
		return
	}

	opts := &types.ArticleOptions{
		SortBy:        db.CreatedAt,
		SortDirection: db.ASC,
		Limit:         l,
		Skip:          s,
	}

	articles, err := db.GetArticles(filter, opts)
	if CheckAPIError(w, err) {
		return
	}
	json.NewEncoder(w).Encode(articles)
}
