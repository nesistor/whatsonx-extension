package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

const (
	webPort = "80"
)

var counts int64

type Config struct {
	DB     *sql.DB
	Models Models
}

func main() {
	log.Println("Starting calendar meeting scheduler service")

	conn := connectToDB()
	if conn == nil {
		log.Panic("Can't connect to Postgres!")
	}

	app := Config{
		DB:     conn,
		Models: NewModels(conn),
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
