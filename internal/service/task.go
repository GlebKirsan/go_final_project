package service

import (
	"errors"
	"log"
	"strconv"
	"time"

	"github.com/GlebKirsan/go-final-project/internal/database"
	"github.com/GlebKirsan/go-final-project/internal/models"
)

type TaskService struct {
	storage     *database.Storage
	dateService *DateService
}

func NewTaskService(storage *database.Storage, dateService *DateService) *TaskService {
	return &TaskService{
		storage:     storage,
		dateService: dateService,
	}
}

func (s *TaskService) GetTask(id int64) (*models.Task, error) {
	return s.storage.Task.GetTask(id)
}

func (s *TaskService) CreateTask(task *models.Task) (id int64, err error) {
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

func (s *TaskService) DeleteTask(id int64) error {
	return s.storage.Task.Delete(id)
}

func (s *TaskService) GetAllTasks(search string) ([]models.Task, error) {
	if search != "" {
		if date, err := time.Parse("02.01.2006", search); err == nil {
			return s.storage.Task.GetAllByDate(date.Format(YYYYMMDD))
		} else {
			return s.storage.Task.GetAllByTitle(search)
		}
	}

	return s.storage.Task.GetAll()
}

func (s *TaskService) MarkDone(id int64) error {
	task, err := s.storage.Task.GetTask(id)
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

	return s.storage.Task.UpdateTask(task)
}

func (s *TaskService) UpdateTask(task *models.Task) error {
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

	id, err := strconv.Atoi(task.ID)
	if err != nil {
		return err
	}
	storedTask, err := s.storage.Task.GetTask(int64(id))
	if err != nil {
		return err
	}
	storedTask.Repeat = task.Repeat
	storedTask.Title = task.Title
	storedTask.Date = task.Date
	storedTask.Comment = task.Comment

	return s.storage.Task.UpdateTask(storedTask)
}
