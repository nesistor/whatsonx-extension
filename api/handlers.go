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

	// Logic to add user (save to DB or Google Calendar)
	// Assume we successfully added the user

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
		app.errorJSON(w, err)
	}
}

// CheckAvailability handler
func (app *Config) CheckAvailability(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")
	groupID := r.URL.Query().Get("group_id")

	// Check availability logic here (use Google Calendar API or mock data)
	// Assume we found availability

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
		app.errorJSON(w, err)
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

	// Logic to schedule meeting (save to DB or Google Calendar)

	// Assume meeting scheduled successfully
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
		app.errorJSON(w, err)
	}
}

// GenerateMeetLink handler
func (app *Config) GenerateMeetLink(w http.ResponseWriter, r *http.Request) {
	// Logic to generate Google Meet link
	meetingID := r.URL.Query().Get("meetingId")
	meetLink := "https://meet.google.com/xyz-abc-123" // Mock link

	response := jsonResponse{
		Error:   false,
		Message: "Google Meet link generated",
		Data: map[string]interface{}{
			"meetLink": meetLink,
		},
	}

	err := app.writeJSON(w, http.StatusOK, response)
	if err != nil {
		app.errorJSON(w, err)
	}
}

// ListUsers handler
func (app *Config) ListUsers(w http.ResponseWriter, r *http.Request) {
	// List all users logic
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
		app.errorJSON(w, err)
	}
}

// ListGroups handler
func (app *Config) ListGroups(w http.ResponseWriter, r *http.Request) {
	// List all groups logic
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
		app.errorJSON(w, err)
	}
}
