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

func (m *Models) GetFreeSlots(ctx context.Context, token *oauth2.Token) ([]string, error) {
	client := oauthConfig.Client(ctx, token)
	srv, err := calendar.New(client)
	if err != nil {
		return nil, fmt.Errorf("unable to create calendar service: %w", err)
	}

	events, err := srv.Events.List("primary").
		TimeMin("2024-11-18T00:00:00Z").
		TimeMax("2024-11-20T00:00:00Z").
		SingleEvents(true).
		OrderBy("startTime").Do()
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve calendar events: %w", err)
	}

	var freeSlots []string
	for _, event := range events.Items {
		fmt.Printf("Event: %s (%s - %s)\n", event.Summary, event.Start.DateTime, event.End.DateTime)
		// W tym miejscu można zaimplementować logikę szukania wolnych slotów
	}

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
