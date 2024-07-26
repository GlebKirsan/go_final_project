package database

import (
	"log"
	"os"

	"github.com/GlebKirsan/go-final-project/internal/env"
	"github.com/GlebKirsan/go-final-project/internal/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Dbinstance struct {
	Db *gorm.DB
}

var DB Dbinstance

func ConnectDB() {
	dbFile := env.GetEnvOrDefault("TODO_DBFILE", "scheduler.db")
	db, err := gorm.Open(sqlite.Open(dbFile), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		log.Fatal("Failed to connect to database.\n", err)
		os.Exit(1)
	}

	log.Println("Database connected")
	db.Logger = logger.Default.LogMode(logger.Info)

	log.Println("Running migration")
	db.AutoMigrate(&models.Task{})

	DB = Dbinstance{
		Db: db,
	}
}
