package main

import (
	"context"
	"github.com/go-chi/chi/v5"
	"net/http"
)

type Service interface {
	GetMedia(ctx context.Context, id string) ([]byte, error)
	GetMediaStream(ctx context.Context, path string) ([]byte, error)
}

type Handler struct {
	service Service
}

func (h Handler) GetVideo(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	b, err := h.service.GetMedia(r.Context(), id)
	if err != nil {
		HandleJSON(w, JSONResponse{Error: err.Error()}, http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/x-mpegURL")
	w.Write(b)
}

func (h Handler) GetStream(w http.ResponseWriter, r *http.Request) {
	path := chi.URLParam(r, "filename")

	b, err := h.service.GetMediaStream(r.Context(), path)
	if err != nil {
		HandleJSON(w, JSONResponse{Error: err.Error()}, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "video/MP2T")
	w.Write(b)
}

func NewHandler(svc Service) Handler {
	return Handler{
		service: svc,
	}
}
