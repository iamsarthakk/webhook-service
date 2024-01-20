package main

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func setupMiddleware(r chi.Router) {
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
}
