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

// AddUser handles user authorization process
// @Summary Initiates user authorization
// @Description Redirects the user to Google OAuth2 authorization page to allow app access.
// @Tags User
// @Accept  json
// @Produce  json
// @Success 200 {string} string "User authorization link"
// @Failure 500 {string} string "Error initiating authorization"
// @Router /add-user [post]
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

// OAuthCallback handles the callback from Google after user authorization
// @Summary Handles OAuth2 callback
// @Description Handles the Google OAuth2 callback and stores the access token.
// @Tags User
// @Accept  json
// @Produce  json
// @Success 200 {string} string "Authorization successful"
// @Failure 500 {string} string "Error during OAuth2 callback"
// @Router /oauth2callback [get]
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

// CheckAvailability checks user's calendar availability
// @Summary Check user calendar availability
// @Description Retrieves free slots from the user's Google Calendar within a given time range.
// @Tags Calendar
// @Accept  json
// @Produce  json
// @Success 200 {array} string "List of free slots"
// @Failure 500 {string} string "Error retrieving availability"
// @Router /check-availability [get]
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

// AddUserToGroup adds a user to a group
// @Summary Add a user to a group
// @Description Adds a specified user to a specified group.
// @Tags Group
// @Accept  json
// @Produce  json
// @Param user_data body AddUserToGroupRequest true "User and Group Data"
// @Success 200 {string} string "User added to group"
// @Failure 500 {string} string "Error adding user to group"
// @Router /add-user-to-group [post]

type AddUserToGroupRequest struct {
	UserEmail string `json:"user_email"`
	GroupName string `json:"group_name"`
}

func (app *Config) AddUserToGroup(w http.ResponseWriter, r *http.Request) {

	var req AddUserToGroupRequest
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

// ListUsers lists all users
// @Summary List all users
// @Description Retrieves the list of all users from the database.
// @Tags User
// @Accept  json
// @Produce  json
// @Success 200 {array} string "List of users"
// @Failure 500 {string} string "Error listing users"
// @Router /list-users [get]
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

// ListGroups lists all groups
// @Summary List all groups
// @Description Retrieves the list of all groups from the database.
// @Tags Group
// @Accept  json
// @Produce  json
// @Success 200 {array} string "List of groups"
// @Failure 500 {string} string "Error listing groups"
// @Router /list-groups [get]
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
