package data

import (
	"database/sql"
	"fmt"
)

// User struct to represent a user
type User struct {
	ID    int    `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
}

// Group struct to represent a group
type Group struct {
	ID      int      `json:"id"`
	Name    string   `json:"name"`
	Members []string `json:"members"`
}

// Meeting struct to represent a scheduled meeting
type Meeting struct {
	ID           int      `json:"id"`
	Title        string   `json:"title"`
	StartTime    string   `json:"startTime"`
	EndTime      string   `json:"endTime"`
	Participants []string `json:"participants"`
	MeetLink     string   `json:"meetLink"`
}

// Models struct to manage database operations
type Models struct {
	DB *sql.DB
}

// NewModels creates a new Models instance with a database connection
func NewModels(db *sql.DB) Models {
	return Models{DB: db}
}

// AddUser adds a user to the database
func (m *Models) AddUser(user User) (int, error) {
	query := `INSERT INTO users (email, name) VALUES ($1, $2) RETURNING id`
	var userID int
	err := m.DB.QueryRow(query, user.Email, user.Name).Scan(&userID)
	if err != nil {
		return 0, fmt.Errorf("unable to insert user: %w", err)
	}
	return userID, nil
}

// ListUsers retrieves all users from the database
func (m *Models) ListUsers() ([]User, error) {
	rows, err := m.DB.Query("SELECT id, email, name FROM users")
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve users: %w", err)
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.ID, &user.Email, &user.Name); err != nil {
			return nil, fmt.Errorf("unable to scan user: %w", err)
		}
		users = append(users, user)
	}
	return users, nil
}

// ListGroups retrieves all groups from the database
func (m *Models) ListGroups() ([]Group, error) {
	rows, err := m.DB.Query("SELECT id, name FROM groups")
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve groups: %w", err)
	}
	defer rows.Close()

	var groups []Group
	for rows.Next() {
		var group Group
		if err := rows.Scan(&group.ID, &group.Name); err != nil {
			return nil, fmt.Errorf("unable to scan group: %w", err)
		}
		groups = append(groups, group)
	}
	return groups, nil
}

// AddGroup adds a group to the database
func (m *Models) AddGroup(group Group) (int, error) {
	query := `INSERT INTO groups (name) VALUES ($1) RETURNING id`
	var groupID int
	err := m.DB.QueryRow(query, group.Name).Scan(&groupID)
	if err != nil {
		return 0, fmt.Errorf("unable to insert group: %w", err)
	}
	return groupID, nil
}

// CheckAvailability checks the availability of a user for a meeting
func (m *Models) CheckAvailability(userID int) (bool, string, error) {
	// Mock implementation (integrate with Google Calendar API or other sources as needed)
	// For now, just return availability and mock time.
	return true, "2024-11-20T10:00:00Z", nil
}

// ScheduleMeeting schedules a new meeting in the database
func (m *Models) ScheduleMeeting(meeting Meeting) (int, error) {
	query := `INSERT INTO meetings (title, start_time, end_time, participants, meet_link) 
	          VALUES ($1, $2, $3, $4, $5) RETURNING id`
	var meetingID int
	err := m.DB.QueryRow(query, meeting.Title, meeting.StartTime, meeting.EndTime,
		meeting.Participants, meeting.MeetLink).Scan(&meetingID)
	if err != nil {
		return 0, fmt.Errorf("unable to schedule meeting: %w", err)
	}
	return meetingID, nil
}

// GenerateMeetLink generates a Google Meet link for a meeting
func (m *Models) GenerateMeetLink(meetingID int) (string, error) {
	// Mock implementation (this would call the Google Calendar API or similar service in a real app)
	// For now, just return a mock link.
	return fmt.Sprintf("https://meet.google.com/xyz-abc-%d", meetingID), nil
}
