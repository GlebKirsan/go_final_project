package service

import (
	"time"

	"github.com/GlebKirsan/go-final-project/internal/models"
)

type TaskService interface {
	Create(*models.Task) (int64, error)
	GetAll(string) ([]models.Task, error)
	Get(int64) (*models.Task, error)
	Update(*models.Task) error
	Delete(int64) error
	MarkDone(int64) error
}

type DateService interface {
	NextDate(time.Time, string, string) (string, error)
	Before(d1 string, d2 string) bool
}

type AuthService interface {
	Authorize(*models.AuthRequest) (string, error)
}
