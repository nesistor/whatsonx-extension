package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"calendar-extension/data"

	_ "github.com/jackc/pgx/v4/stdlib"
)

const (
	webPort = "80"
)

var counts int64

type Config struct {
	DB     *sql.DB
	Models data.Models
}

// @title Google Meeting Scheduler API
// @version 1.0
// @description API do synchronizacji kalendarza Google i planowania spotkań

// @BasePath /

// main function starts the application, connects to the database, and listens for HTTP requests
// @Summary Starts the Google Meeting Scheduler Service
// @Description Initializes the service and listens for incoming HTTP requests.
// @Tags Startup
// @Accept  json
// @Produce  json
// @Success 200 {string} string "Service Started"
// @Router / [get]
func main() {
	log.Println("Starting calendar meeting scheduler service")

	conn := connectToDB()
	if conn == nil {
		log.Panic("Can't connect to Postgres!")
	}

	app := Config{
		DB:     conn,
		Models: data.NewModels(conn),
	}

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	go func() {
		err := srv.ListenAndServe()
		if err != nil {
			log.Panic(err)
		}
	}()

	select {}
}

// connectToDB establishes a connection to the PostgreSQL database
// @Summary Connect to the PostgreSQL database
// @Description Tries to connect to the database and retries on failure.
// @Tags Database
// @Accept  json
// @Produce  json
// @Success 200 {string} string "Connected to database"
// @Failure 500 {string} string "Error connecting to the database"
// @Router /connect-to-db [get]
func connectToDB() *sql.DB {
	dsn := os.Getenv("DSN")
	log.Println("DSN:", dsn)

	for {
		connection, err := openDB(dsn)
		if err != nil {
			fmt.Errorf("Error opening database: %w", err)
			counts++
		} else {
			log.Printf("Connected to Postgres!")
			return connection
		}

		if counts > 10 {
			log.Println(err)
			return nil
		}

		log.Println("Backing off for two seconds....")
		time.Sleep(2 * time.Second)
		continue
	}
}

// openDB tries to open a database connection and pings it
// @Summary Open a database connection
// @Description Attempts to open a connection and check if the database is reachable.
// @Tags Database
// @Accept  json
// @Produce  json
// @Success 200 {string} string "Database connection successful"
// @Failure 500 {string} string "Error opening database"
// @Router /open-db [get]
func openDB(dsn string) (*sql.DB, error) {
	log.Println("Trying to connect to the database...")
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, fmt.Errorf("error opening database: %w", err)
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("error pinging database: %w", err)
	}

	log.Println("Successfully connected to the database")
	return db, nil
}
