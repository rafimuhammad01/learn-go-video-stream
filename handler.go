package main

import (
	"context"
	"io"
	"net/http"
	"path/filepath"

	"github.com/go-chi/chi/v5"
)

type Service interface {
	GetMedia(ctx context.Context, id string) ([]byte, error)
	GetMediaStream(ctx context.Context, path string) ([]byte, error)
	UploadMedia(ctx context.Context, file io.Reader, fileName string) error
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

func (h Handler) UploadVideo(w http.ResponseWriter, r *http.Request) {
	mr, err := r.MultipartReader()
	if err != nil {
		HandleJSON(w, JSONResponse{Error: err.Error()}, http.StatusInternalServerError)
		return
	}

	for {
		part, err := mr.NextPart()
		if err != nil {
			if err == io.EOF {
				break
			}
			HandleJSON(w, JSONResponse{Error: err.Error()}, http.StatusInternalServerError)
			return
		}

		if part.FileName() != "" {
			err = h.service.UploadMedia(r.Context(), part, filepath.Ext(part.FileName()))
			if err != nil {
				HandleJSON(w, JSONResponse{Error: err.Error()}, http.StatusInternalServerError)
				return
			}
		}

	}

	w.Header().Set("Content-Type", "application/json")
	HandleJSON(w, JSONResponse{
		Data: "success",
	}, http.StatusOK)
}

func NewHandler(svc Service) Handler {
	return Handler{
		service: svc,
	}
}
