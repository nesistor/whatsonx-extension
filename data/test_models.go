package data

import (
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"golang.org/x/oauth2"
)

func TestNewModels(t *testing.T) {
	db, _, _ := sqlmock.New()
	defer db.Close()

	models := NewModels(db)
	if models.DB != db {
		t.Fatalf("expected db to be set")
	}
}

func TestSaveUserToken(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	models := NewModels(db)
	token := &oauth2.Token{
		AccessToken:  "access-token",
		RefreshToken: "refresh-token",
		Expiry:       time.Now(),
	}

	mock.ExpectExec(`INSERT INTO user_tokens`).
		WithArgs(token.AccessToken, token.RefreshToken, token.Expiry).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := models.SaveUserToken(token)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestGetUserToken(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	models := NewModels(db)
	email := "test@example.com"
	token := &oauth2.Token{
		AccessToken:  "access-token",
		RefreshToken: "refresh-token",
		Expiry:       time.Now(),
	}

	mock.ExpectQuery(`SELECT access_token, refresh_token, expiry FROM user_tokens WHERE email =`).
		WithArgs(email).
		WillReturnRows(sqlmock.NewRows([]string{"access_token", "refresh_token", "expiry"}).
			AddRow(token.AccessToken, token.RefreshToken, token.Expiry))

	result, err := models.GetUserToken(email)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.AccessToken != token.AccessToken || result.RefreshToken != token.RefreshToken || !result.Expiry.Equal(token.Expiry) {
		t.Fatalf("unexpected token result: %v", result)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestGetFreeSlots(t *testing.T) {
	// This requires integration with Google Calendar API, so it would be best to mock
	// the external API calls. A package like `gomock` or `httptest` can help.
	t.Skip("Integration test required for Google Calendar API")
}

func TestAddUserToGroup(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	models := NewModels(db)
	userEmail := "test@example.com"
	groupName := "test-group"

	mock.ExpectBegin()
	mock.ExpectExec(`INSERT INTO groups`).
		WithArgs(groupName).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec(`INSERT INTO user_groups`).
		WithArgs(userEmail, groupName).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := models.AddUserToGroup(userEmail, groupName)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestListUsers(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	models := NewModels(db)
	users := []string{"user1@example.com", "user2@example.com"}

	mock.ExpectQuery(`SELECT email FROM users`).
		WillReturnRows(sqlmock.NewRows([]string{"email"}).
			AddRow(users[0]).
			AddRow(users[1]))

	result, err := models.ListUsers()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result) != len(users) {
		t.Fatalf("expected %d users, got %d", len(users), len(result))
	}

	for i, email := range result {
		if email != users[i] {
			t.Errorf("expected %s, got %s", users[i], email)
		}
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestListGroups(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	models := NewModels(db)
	groups := []string{"group1", "group2"}

	mock.ExpectQuery(`SELECT name FROM groups`).
		WillReturnRows(sqlmock.NewRows([]string{"name"}).
			AddRow(groups[0]).
			AddRow(groups[1]))

	result, err := models.ListGroups()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result) != len(groups) {
		t.Fatalf("expected %d groups, got %d", len(groups), len(result))
	}

	for i, name := range result {
		if name != groups[i] {
			t.Errorf("expected %s, got %s", groups[i], name)
		}
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestInitializeDatabase(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	models := NewModels(db)

	mock.ExpectExec(`CREATE TABLE IF NOT EXISTS groups`).
		WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectExec(`CREATE TABLE IF NOT EXISTS user_groups`).
		WillReturnResult(sqlmock.NewResult(0, 0))

	err := models.InitializeDatabase()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}
