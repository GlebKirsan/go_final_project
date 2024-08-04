package database

import (
	"database/sql"
	"log"

	"github.com/GlebKirsan/go-final-project/internal/config"
	"github.com/GlebKirsan/go-final-project/internal/database/repositories"
	_ "modernc.org/sqlite"
)

type DB struct {
	*sql.DB
}

type Storage struct {
	Db *sql.DB

	Task *repositories.TaskRepo
}

func ConnectDB() (*Storage, error) {
	cfg := config.Get()
	db, err := sql.Open("sqlite", cfg.DBFile)
	if err != nil {
		return nil, err
	}

	log.Println("Database connected")

	log.Println("Running migration")
	if err = runMigration(db); err != nil {
		return nil, err
	}

	return &Storage{Db: db, Task: repositories.NewTaskRepo(db)}, nil
}
