package service

import (
	"errors"
	"log"
	"time"

	"github.com/GlebKirsan/go-final-project/internal/database"
	"github.com/GlebKirsan/go-final-project/internal/models"
)

type taskService struct {
	storage     *database.Storage
	dateService DateService
}

func NewTaskService(storage *database.Storage, dateService DateService) *taskService {
	return &taskService{
		storage:     storage,
		dateService: dateService,
	}
}

func (s *taskService) Get(id int64) (*models.Task, error) {
	return s.storage.Task.Get(id)
}

func (s *taskService) Create(task *models.Task) (id int64, err error) {
	if task.Title == "" {
		return 0, errors.New("title is empty")
	}

	if _, err := time.Parse(YYYYMMDD, task.Date); task.Date != "" && err != nil {
		return 0, err
	}

	now := time.Now()
	if task.Date != "" {
		if s.dateService.Before(task.Date, now.Format(YYYYMMDD)) {
			if task.Repeat == "" {
				task.Date = now.Format(YYYYMMDD)
			} else {
				task.Date, err = s.dateService.NextDate(now, task.Date, task.Repeat)
				if err != nil {
					return 0, err
				}
			}
		}
	} else {
		task.Date = now.Format(YYYYMMDD)
	}

	log.Println(task)
	id, err = s.storage.Task.Create(task)
	return
}

func (s *taskService) Delete(id int64) error {
	return s.storage.Task.Delete(id)
}

func (s *taskService) GetAll(search string) ([]models.Task, error) {
	return s.storage.Task.GetAll(search)
}

func (s *taskService) MarkDone(id int64) error {
	task, err := s.storage.Task.Get(id)
	if err != nil {
		return err
	}

	if task.Repeat == "" {
		return s.storage.Task.Delete(id)
	}

	task.Date, err = s.dateService.NextDate(time.Now(), task.Date, task.Repeat)
	if err != nil {
		return err
	}

	return s.storage.Task.Update(task)
}

func (s *taskService) Update(task *models.Task) error {
	if task.Title == "" {
		return errors.New("title is empty")
	}

	if _, err := time.Parse(YYYYMMDD, task.Date); task.Date != "" && err != nil {
		return err
	}

	now := time.Now()
	if task.Date != "" {
		if s.dateService.Before(task.Date, now.Format(YYYYMMDD)) {
			if task.Repeat == "" {
				task.Date = now.Format(YYYYMMDD)
			} else {
				var err error
				task.Date, err = s.dateService.NextDate(now, task.Date, task.Repeat)
				if err != nil {
					return err
				}
			}
		}
	} else {
		task.Date = now.Format(YYYYMMDD)
	}

	storedTask, err := s.storage.Task.Get(task.ID)
	if err != nil {
		return err
	}
	storedTask.Repeat = task.Repeat
	storedTask.Title = task.Title
	storedTask.Date = task.Date
	storedTask.Comment = task.Comment

	return s.storage.Task.Update(storedTask)
}
