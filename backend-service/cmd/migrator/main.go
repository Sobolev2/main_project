package main

import (
	"context"
	"database/sql"
	"log"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
)

const (
	gooseUpCommandName   = "up"
	gooseDownCommandName = "down"
)

func main() {
	if godotenv.Load() == nil {
		log.Println("env variables loaded from .env file")
	}

	postgresDSN := os.Getenv("POSTGRES_DSN")
	if postgresDSN == "" {
		log.Fatal("POSTGRES_DSN environment variable is required")
	}

	if len(os.Args) < 2 {
		log.Fatal("You must pass required argument: migration directory")
	}
	migrationDir := os.Args[1]

	db, err := sql.Open("postgres", postgresDSN)
	if err != nil {
		log.Fatalf("can not connect to database: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Fatalf("failed to close database connection: %v", err)
		}
	}()

	if err := db.Ping(); err != nil {
		log.Fatalf("can not ping database: %v", err)
	}

	command := gooseUpCommandName
	if len(os.Args) > 2 {
		switch strings.ToLower(os.Args[2]) {
		case gooseUpCommandName:
			command = gooseUpCommandName
		case gooseDownCommandName:
			command = gooseDownCommandName
		default:
			log.Fatalf("Invalid command: %s. Use 'up' or 'down'", os.Args[2])
		}
	}

	timeout, cancelFunc := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancelFunc()

	err = goose.RunContext(timeout, command, db, migrationDir)
	if err != nil {
		log.Fatalf("goose.Run() error: %v", err)
	}

	log.Printf("Migration '%s' completed successfully", command)
}
