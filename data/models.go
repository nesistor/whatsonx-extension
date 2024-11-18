package data

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
)

var oauthConfig = &oauth2.Config{
	ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
	ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
	RedirectURL:  "http://localhost:80/oauth2callback",
	Scopes:       []string{calendar.CalendarReadonlyScope},
	Endpoint:     google.Endpoint,
}

type Models struct {
	DB *sql.DB
}

func NewModels(db *sql.DB) Models {
	return Models{DB: db}
}

func (m *Models) SaveUserToken(token *oauth2.Token) error {
	query := `INSERT INTO user_tokens (access_token, refresh_token, expiry) VALUES ($1, $2, $3)`
	_, err := m.DB.Exec(query, token.AccessToken, token.RefreshToken, token.Expiry)
	if err != nil {
		return fmt.Errorf("failed to save token: %w", err)
	}
	return nil
}

func (m *Models) GetUserToken(email string) (*oauth2.Token, error) {
	query := `SELECT access_token, refresh_token, expiry FROM user_tokens WHERE email = $1`
	row := m.DB.QueryRow(query, email)

	var accessToken, refreshToken string
	var expiry time.Time
	err := row.Scan(&accessToken, &refreshToken, &expiry)
	if err != nil {
		return nil, fmt.Errorf("failed to get token: %w", err)
	}

	return &oauth2.Token{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		Expiry:       expiry,
	}, nil
}

// GetFreeSlots retrieves free time slots for the next week for the authenticated user.
func (m *Models) GetFreeSlots(ctx context.Context, token *oauth2.Token) ([]string, error) {
	// Create a custom HTTP client
	client := oauthConfig.Client(ctx, token)

	// Create the Google Calendar service using the new recommended method
	srv, err := calendar.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return nil, fmt.Errorf("unable to create calendar service: %w", err)
	}

	// Get the current time and define the time range for the next week
	now := time.Now()
	startOfWeek := now.UTC().Format(time.RFC3339)                       // Start of today in UTC
	endOfWeek := now.Add(7 * 24 * time.Hour).UTC().Format(time.RFC3339) // One week from now

	// Fetch the events for the next week
	events, err := srv.Events.List("primary").
		TimeMin(startOfWeek).
		TimeMax(endOfWeek).
		SingleEvents(true).
		OrderBy("startTime").
		Do()
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve calendar events: %w", err)
	}

	// Create a slice to store the free time slots
	var freeSlots []string

	// Let's assume work hours are from 9:00 AM to 5:00 PM, adjust according to your needs.
	workStart := 9
	workEnd := 17

	// Loop through the events and find gaps between them
	var lastEndTime time.Time
	for _, event := range events.Items {
		// Parse the event start and end time
		startTime, err := time.Parse(time.RFC3339, event.Start.DateTime)
		if err != nil {
			startTime, err = time.Parse(time.RFC3339, event.Start.Date) // handle full-day events
			if err != nil {
				return nil, fmt.Errorf("error parsing event start time: %w", err)
			}
		}

		// If this is the first event, check if there's a gap before it
		if lastEndTime.IsZero() {
			if startTime.Hour() > workStart {
				// Found a free slot before the first event
				freeSlot := fmt.Sprintf("%s to %s", fmt.Sprintf("%02d:00", workStart), startTime.Format("15:04"))
				freeSlots = append(freeSlots, freeSlot)
			}
		} else {
			// Check for a gap between the last event and the current one
			if startTime.Sub(lastEndTime) > time.Minute*30 { // Allow a 30 min gap
				freeSlot := fmt.Sprintf("%s to %s", lastEndTime.Format("15:04"), startTime.Format("15:04"))
				freeSlots = append(freeSlots, freeSlot)
			}
		}

		// Update the last end time
		endTime, err := time.Parse(time.RFC3339, event.End.DateTime)
		if err != nil {
			endTime, err = time.Parse(time.RFC3339, event.End.Date) // handle full-day events
			if err != nil {
				return nil, fmt.Errorf("error parsing event end time: %w", err)
			}
		}
		lastEndTime = endTime
	}

	// Check if there are free slots after the last event for the rest of the day
	if !lastEndTime.IsZero() && lastEndTime.Hour() < workEnd {
		freeSlot := fmt.Sprintf("%s to %s", lastEndTime.Format("15:04"), fmt.Sprintf("%02d:00", workEnd))
		freeSlots = append(freeSlots, freeSlot)
	}

	// Return the list of free slots
	return freeSlots, nil
}

func (m *Models) AddUserToGroup(userEmail, groupName string) error {
	tx, err := m.DB.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Ensure group exists
	queryGroup := `INSERT INTO groups (name) VALUES ($1) ON CONFLICT (name) DO NOTHING`
	_, err = tx.Exec(queryGroup, groupName)
	if err != nil {
		return fmt.Errorf("failed to insert group: %w", err)
	}

	// Link user to group
	queryLink := `
		INSERT INTO user_groups (user_email, group_name) 
		VALUES ($1, $2) 
		ON CONFLICT DO NOTHING
	`
	_, err = tx.Exec(queryLink, userEmail, groupName)
	if err != nil {
		return fmt.Errorf("failed to link user to group: %w", err)
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (m *Models) ListUsers() ([]string, error) {
	query := `SELECT email FROM users`
	rows, err := m.DB.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query users: %w", err)
	}
	defer rows.Close()

	var users []string
	for rows.Next() {
		var email string
		err := rows.Scan(&email)
		if err != nil {
			return nil, fmt.Errorf("failed to scan user email: %w", err)
		}
		users = append(users, email)
	}

	return users, nil
}

func (m *Models) ListGroups() ([]string, error) {
	query := `SELECT name FROM groups`
	rows, err := m.DB.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query groups: %w", err)
	}
	defer rows.Close()

	var groups []string
	for rows.Next() {
		var name string
		err := rows.Scan(&name)
		if err != nil {
			return nil, fmt.Errorf("failed to scan group name: %w", err)
		}
		groups = append(groups, name)
	}

	return groups, nil
}
func (m *Models) InitializeDatabase() error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS groups (
			id SERIAL PRIMARY KEY,
			name VARCHAR(255) UNIQUE NOT NULL
		);`,
		`CREATE TABLE IF NOT EXISTS user_groups (
			id SERIAL PRIMARY KEY,
			user_email VARCHAR(255) NOT NULL,
			group_name VARCHAR(255) NOT NULL,
			FOREIGN KEY (user_email) REFERENCES users(email),
			FOREIGN KEY (group_name) REFERENCES groups(name),
			UNIQUE (user_email, group_name)
		);`,
	}

	for _, query := range queries {
		_, err := m.DB.Exec(query)
		if err != nil {
			return fmt.Errorf("failed to execute query: %w", err)
		}
	}

	return nil
}
