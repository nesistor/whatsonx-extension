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
