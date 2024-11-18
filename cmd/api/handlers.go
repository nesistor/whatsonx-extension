package main

import (
	"net/http"
)

type jsonResponse struct {
	Error   bool        `json:"error"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// AddUser handler
func (app *Config) AddUser(w http.ResponseWriter, r *http.Request) {
	var requestPayload struct {
		Email string `json:"email"`
		Name  string `json:"name"`
	}

	err := app.readJSON(w, r, &requestPayload)
	if err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	response := jsonResponse{
		Error:   false,
		Message: "User added successfully",
		Data: map[string]interface{}{
			"email": requestPayload.Email,
			"name":  requestPayload.Name,
		},
	}

	err = app.writeJSON(w, http.StatusOK, response)
	if err != nil {
		app.errorJSON(w, err, http.StatusInternalServerError)
	}
}

// CheckAvailability handler
func (app *Config) CheckAvailability(w http.ResponseWriter, r *http.Request) {
	availability := struct {
		Available         bool   `json:"available"`
		NextAvailableTime string `json:"nextAvailableTime"`
	}{
		Available:         true,
		NextAvailableTime: "2024-11-20T10:00:00Z",
	}

	response := jsonResponse{
		Error:   false,
		Message: "Availability data",
		Data:    availability,
	}

	err := app.writeJSON(w, http.StatusOK, response)
	if err != nil {
		app.errorJSON(w, err, http.StatusInternalServerError)
	}
}

// ScheduleMeeting handler
func (app *Config) ScheduleMeeting(w http.ResponseWriter, r *http.Request) {
	var meeting struct {
		Title        string   `json:"title"`
		StartTime    string   `json:"startTime"`
		EndTime      string   `json:"endTime"`
		Participants []string `json:"participants"`
	}

	err := app.readJSON(w, r, &meeting)
	if err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	meetLink := "https://meet.google.com/xyz-abc-123"

	response := jsonResponse{
		Error:   false,
		Message: "Meeting successfully scheduled",
		Data: map[string]interface{}{
			"meetingId": "12345",
			"meetLink":  meetLink,
		},
	}

	err = app.writeJSON(w, http.StatusOK, response)
	if err != nil {
		app.errorJSON(w, err, http.StatusInternalServerError)
	}
}

// GenerateMeetLink handler
func (app *Config) GenerateMeetLink(w http.ResponseWriter, r *http.Request) {
	meetLink := "https://meet.google.com/xyz-abc-123"

	response := jsonResponse{
		Error:   false,
		Message: "Google Meet link generated",
		Data: map[string]interface{}{
			"meetLink": meetLink,
		},
	}

	err := app.writeJSON(w, http.StatusOK, response)
	if err != nil {
		app.errorJSON(w, err, http.StatusInternalServerError)
	}
}

// ListUsers handler
func (app *Config) ListUsers(w http.ResponseWriter, r *http.Request) {
	users := []struct {
		UserID string `json:"userId"`
		Email  string `json:"email"`
		Name   string `json:"name"`
	}{
		{"1", "user1@example.com", "User One"},
		{"2", "user2@example.com", "User Two"},
	}

	response := jsonResponse{
		Error:   false,
		Message: "List of users",
		Data:    users,
	}

	err := app.writeJSON(w, http.StatusOK, response)
	if err != nil {
		app.errorJSON(w, err, http.StatusInternalServerError)
	}
}

// ListGroups handler
func (app *Config) ListGroups(w http.ResponseWriter, r *http.Request) {
	groups := []struct {
		GroupID   string   `json:"groupId"`
		GroupName string   `json:"groupName"`
		Members   []string `json:"members"`
	}{
		{"1", "Group 1", []string{"user1@example.com", "user2@example.com"}},
		{"2", "Group 2", []string{"user3@example.com"}},
	}

	response := jsonResponse{
		Error:   false,
		Message: "List of groups",
		Data:    groups,
	}

	err := app.writeJSON(w, http.StatusOK, response)
	if err != nil {
		app.errorJSON(w, err, http.StatusInternalServerError)
	}
}
