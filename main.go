package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(cors.AllowAll().Handler)

	db := NewDB()
	fr := NewFile()
	svc := NewService(fr, fr, db)
	h := NewHandler(svc)

	r.Get("/video/{id}", h.GetVideo)
	r.Get("/video/stream/{filename}", h.GetStream)
	r.Post("/video/upload", h.UploadVideo)

	http.ListenAndServe(":8080", r)
}
