package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

// routes sets up the API endpoints for the application
// @Summary Set up API routes for the application
// @Description Defines all the routes and the associated handlers for the application.
// @Tags Routes
// @Accept  json
// @Produce  json
// @Success 200 {string} string "Routes set successfully"
// @Router /routes [get]
func (app *Config) routes() http.Handler {
	mux := chi.NewRouter()

	mux.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"POST", "PUT", "GET", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Content-Type", "Authorization", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	mux.Use(middleware.Heartbeat("/ping"))

	mux.Post("/add-user", app.AddUser)
	mux.Get("/oauth2callback", app.OAuthCallback)
	mux.Get("/check-availability", app.CheckAvailability)
	mux.Post("/add-user-to-group", app.AddUserToGroup)
	mux.Get("/list-users", app.ListUsers)
	mux.Get("/list-groups", app.ListGroups)

	return mux
}
