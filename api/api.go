package api

import (
	"log"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func randomInt(w http.ResponseWriter, r *http.Request) {
	log.Println("~ hit")
	w.Write([]byte(strconv.Itoa(rand.Intn(100))))
}

func click(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("clicked"))
}

func GetLocalRouter() http.Handler {
	r := chi.NewRouter()

	r.Get("/", randomInt)
	r.Get("/click", click)

	r.Mount("/posts", getPostsRouter())
	r.Mount("/auth", getAuthRouter())

	return r
}
