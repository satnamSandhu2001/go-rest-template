package routes

import (
	"go-rest-template/internal/db"
	"go-rest-template/internal/handlers"
	"go-rest-template/internal/middlewares"
	"go-rest-template/internal/repository"

	"github.com/go-chi/chi/v5"
)

func RegisterUserRoutes(r chi.Router, q *db.Queries) {
	repo := repository.NewUserRepository(q)
	h := handlers.NewUserHandler(repo)

	r.Route("/auth", func(r chi.Router) {
		r.Post("/signup", h.Signup)
		r.Post("/login", h.Login)
		r.Get("/logout", h.Logout)
	})

	r.Route("/user", func(r chi.Router) {
		r.With(middlewares.IsAuthenticated(q)).Get("/profile", h.GetMyDetails)
	})
}
