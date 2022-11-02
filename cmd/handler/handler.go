package handler

import "github.com/go-chi/chi"

type handler struct {
	Router *chi.Mux
}

func New() *handler {
	var h handler
	r := chi.NewMux()
	h.Router = r
	return &h
}
