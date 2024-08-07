package database

import (
	"database/sql"
	"log"

	_ "modernc.org/sqlite"

	"github.com/GlebKirsan/go-final-project/internal/config"
	"github.com/GlebKirsan/go-final-project/internal/database/repositories"
)

type Storage struct {
	Db *sql.DB

	Task TaskRepo
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

	taskRepo, err := repositories.NewTaskRepo(db)
	if err != nil {
		return nil, err
	}
	return &Storage{Db: db, Task: taskRepo}, nil
}

func (s *Storage) Close() {
	s.Task.Close()
	s.Db.Close()
}
