// @title Calendar Meeting Scheduler API
// @version 1.0
// @description API for scheduling meetings and managing users/groups using Google Calendar.
// @termsOfService http://example.com/terms/

// @contact.name API Support
// @contact.url
// @contact.email karolmalicki.001@gmail.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:80
// @BasePath /
// @schemes http
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
