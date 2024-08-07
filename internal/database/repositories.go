package database

import "github.com/GlebKirsan/go-final-project/internal/models"

type TaskRepo interface {
	Create(*models.Task) (id int64, err error)
	Get(int64) (*models.Task, error)
	Update(*models.Task) error
	Delete(int64) error
	GetAll(string) ([]models.Task, error)
	Close() error
}
