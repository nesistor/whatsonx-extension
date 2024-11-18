package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
)

var oauthConfig = &oauth2.Config{
	ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
	ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
	RedirectURL:  "http://localhost:80/oauth2callback",
	Scopes:       []string{calendar.CalendarReadonlyScope},
	Endpoint:     google.Endpoint,
}

func (app *Config) AddUser(w http.ResponseWriter, r *http.Request) {
	url := oauthConfig.AuthCodeURL("state", oauth2.AccessTypeOffline)
	response := jsonResponse{
		Error:   false,
		Message: "Click the link to authorize the app",
		Data:    url,
	}

	err := app.writeJSON(w, http.StatusOK, response)
	if err != nil {
		app.errorJSON(w, err, http.StatusInternalServerError)
	}
}

func (app *Config) OAuthCallback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	if code == "" {
		app.errorJSON(w, fmt.Errorf("no code in request"), http.StatusBadRequest)
		return
	}

	token, err := oauthConfig.Exchange(context.Background(), code)
	if err != nil {
		app.errorJSON(w, fmt.Errorf("failed to exchange token: %w", err), http.StatusInternalServerError)
		return
	}

	err = app.Models.SaveUserToken(token)
	if err != nil {
		app.errorJSON(w, fmt.Errorf("failed to save token: %w", err), http.StatusInternalServerError)
		return
	}

	response := jsonResponse{
		Error:   false,
		Message: "Authorization successful",
	}
	err = app.writeJSON(w, http.StatusOK, response)
	if err != nil {
		app.errorJSON(w, err, http.StatusInternalServerError)
	}
}

func (app *Config) CheckAvailability(w http.ResponseWriter, r *http.Request) {
	token, err := app.Models.GetUserToken("user_email@example.com")
	if err != nil {
		app.errorJSON(w, fmt.Errorf("failed to get user token: %w", err), http.StatusInternalServerError)
		return
	}

	freeSlots, err := app.Models.GetFreeSlots(r.Context(), token)
	if err != nil {
		app.errorJSON(w, fmt.Errorf("failed to get free slots: %w", err), http.StatusInternalServerError)
		return
	}

	response := jsonResponse{
		Error:   false,
		Message: "Availability data",
		Data:    freeSlots,
	}

	err = app.writeJSON(w, http.StatusOK, response)
	if err != nil {
		app.errorJSON(w, err, http.StatusInternalServerError)
	}
}

func (app *Config) AddUserToGroup(w http.ResponseWriter, r *http.Request) {
	type payload struct {
		UserEmail string `json:"user_email"`
		GroupName string `json:"group_name"`
	}

	var req payload
	err := app.readJSON(w, r, &req)
	if err != nil {
		app.errorJSON(w, fmt.Errorf("invalid request body: %w", err), http.StatusBadRequest)
		return
	}

	err = app.Models.AddUserToGroup(req.UserEmail, req.GroupName)
	if err != nil {
		app.errorJSON(w, fmt.Errorf("failed to add user to group: %w", err), http.StatusInternalServerError)
		return
	}

	response := jsonResponse{
		Error:   false,
		Message: fmt.Sprintf("User %s added to group %s", req.UserEmail, req.GroupName),
	}

	err = app.writeJSON(w, http.StatusOK, response)
	if err != nil {
		app.errorJSON(w, err, http.StatusInternalServerError)
	}
}

func (app *Config) ListUsers(w http.ResponseWriter, r *http.Request) {
	users, err := app.Models.ListUsers()
	if err != nil {
		app.errorJSON(w, fmt.Errorf("failed to list users: %w", err), http.StatusInternalServerError)
		return
	}

	response := jsonResponse{
		Error:   false,
		Message: "List of users",
		Data:    users,
	}

	err = app.writeJSON(w, http.StatusOK, response)
	if err != nil {
		app.errorJSON(w, err, http.StatusInternalServerError)
	}
}

func (app *Config) ListGroups(w http.ResponseWriter, r *http.Request) {
	groups, err := app.Models.ListGroups()
	if err != nil {
		app.errorJSON(w, fmt.Errorf("failed to list groups: %w", err), http.StatusInternalServerError)
		return
	}

	response := jsonResponse{
		Error:   false,
		Message: "List of groups",
		Data:    groups,
	}

	err = app.writeJSON(w, http.StatusOK, response)
	if err != nil {
		app.errorJSON(w, err, http.StatusInternalServerError)
	}
}
