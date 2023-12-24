package main

import (
	"log"
	"net/http"
	"os"

	"github.com/Ehab-24/test/api"
	"github.com/Ehab-24/test/cache"
	"github.com/Ehab-24/test/db"
	"github.com/Ehab-24/test/views"
	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
)

func main() {
	log.SetFlags(0)

	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
	db.ConnectToDatabase()
	cache.Init()

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	http.ListenAndServe(":"+port, getRouter())
}

func getRouter() *chi.Mux {
	r := chi.NewRouter()

	r.Mount("/", views.GetLocalRouter())
	r.Mount("/api", api.GetLocalRouter())

	return r
}
