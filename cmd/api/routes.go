package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

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
	mux.Post("/schedule-meeting", app.ScheduleMeeting)
	mux.Get("/generate-meet-link", app.GenerateMeetLink)
	mux.Get("/check-availability", app.CheckAvailability)
	mux.Get("/list-users", app.ListUsers)
	mux.Get("/list-groups", app.ListGroups)

	return mux
}
