package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"net/http"
)

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(cors.AllowAll().Handler)

	db := NewDB()
	fr := NewFileReader()
	svc := NewService(fr, db)
	h := NewHandler(svc)

	r.Get("/video/{id}", h.GetVideo)
	r.Get("/video/stream/{filename}", h.GetStream)

	http.ListenAndServe(":8080", r)
}
