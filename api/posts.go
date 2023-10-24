package api

import (
	"context"
	"log"
	"net/http"
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
	r.With(middleware.VerifyUser).Post("/", CreatePost)

	return r
}

func CreatePost(w http.ResponseWriter, r *http.Request) {
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
