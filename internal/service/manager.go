package service

import (
	"errors"

	"github.com/GlebKirsan/go-final-project/internal/database"
)

type Manager struct {
	Task *TaskService
	Date *DateService
}

func NewManager(storage *database.Storage) (*Manager, error) {
	if storage == nil {
		return nil, errors.New("storage not provided")
	}

	dateService := NewDateService()
	return &Manager{
		Task: NewTaskService(storage, dateService),
		Date: dateService,
	}, nil
}
