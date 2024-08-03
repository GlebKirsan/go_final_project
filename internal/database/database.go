package database

import (
	"database/sql"
	"fmt"
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

func FindByDate(tasks *[]models.Task, date string) *gorm.DB {
	return DB.Db.Limit(50).Order("date").Where("date = @date", sql.Named("date", date)).Find(&tasks)
}

func FindByTitle(tasks *[]models.Task, title string) *gorm.DB {
	return DB.Db.Limit(50).Order("date").Where("title LIKE @title", sql.Named("title", fmt.Sprintf("%%%s%%", title))).Find(&tasks)
}

func FindAll(tasks *[]models.Task) *gorm.DB {
	return DB.Db.Limit(50).Order("date").Find(&tasks)
}

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
