package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/GlebKirsan/go-final-project/internal/logger"
	"github.com/GlebKirsan/go-final-project/internal/models"
	"github.com/GlebKirsan/go-final-project/internal/service"
	"github.com/unrolled/render"
)

type TaskHandler struct {
	taskService *service.TaskService
	render      *render.Render
	logger      *logger.Logger
}

func NewTaskHandler(manager *service.Manager, render *render.Render, logger *logger.Logger) *TaskHandler {
	return &TaskHandler{
		taskService: manager.Task,
		render:      render,
		logger:      logger,
	}
}

func (handler *TaskHandler) handleError(w http.ResponseWriter, err error, code int) {
	handler.logger.Error().Msg(err.Error())
	handler.render.JSON(w, code, map[string]any{"error": err.Error()})
}

func (handler *TaskHandler) PostTask(w http.ResponseWriter, r *http.Request) {
	task := models.Task{}
	var buf bytes.Buffer

	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		handler.handleError(w, err, http.StatusBadRequest)
		return
	}

	if err = json.Unmarshal(buf.Bytes(), &task); err != nil {
		handler.handleError(w, err, http.StatusBadRequest)
		return
	}

	handler.logger.Debug().Msg(task.String())

	id, err := handler.taskService.CreateTask(&task)
	if err != nil {
		handler.handleError(w, err, http.StatusBadRequest)
		return
	}

	handler.render.JSON(w, http.StatusOK, map[string]any{
		"id": id,
	})
}

type GetTasksResp struct {
	Tasks []models.Task `json:"tasks"`
}

func (handler *TaskHandler) GetTasks(w http.ResponseWriter, r *http.Request) {
	search := r.URL.Query().Get("search")
	tasks, err := handler.taskService.GetAllTasks(search)
	if err != nil {
		handler.handleError(w, err, http.StatusBadRequest)
	}

	handler.render.JSON(w, http.StatusOK, GetTasksResp{
		Tasks: tasks,
	})
}

func (handler *TaskHandler) GetTask(w http.ResponseWriter, r *http.Request) {
	i := r.URL.Query().Get("id")

	if i == "" {
		handler.handleError(w, errors.New("id is empty"), http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(i)
	if err != nil {
		handler.handleError(w, err, http.StatusBadRequest)
		return
	}

	task, err := handler.taskService.GetTask(int64(id))
	if err != nil {
		handler.handleError(w, err, http.StatusInternalServerError)
		return
	}
	handler.render.JSON(w, http.StatusOK, task)
}

func (handler *TaskHandler) UpdateTask(w http.ResponseWriter, r *http.Request) {
	task := models.Task{}
	var buf bytes.Buffer

	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		handler.handleError(w, err, http.StatusBadRequest)
		return
	}

	if err = json.Unmarshal(buf.Bytes(), &task); err != nil {
		handler.handleError(w, err, http.StatusBadRequest)
		return
	}

	handler.logger.Debug().Msg(task.String())
	err = handler.taskService.UpdateTask(&task)
	if err != nil {
		handler.handleError(w, err, http.StatusBadRequest)
		return
	}

	handler.render.JSON(w, http.StatusOK, map[string]any{})
}

func (handler *TaskHandler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	if id == "" {
		handler.handleError(w, errors.New("id is empty"), http.StatusBadRequest)
		return
	}

	i, err := strconv.Atoi(id)
	if err != nil {
		handler.handleError(w, err, http.StatusBadRequest)
		return
	}
	err = handler.taskService.DeleteTask(int64(i))
	if err != nil {
		handler.handleError(w, err, http.StatusInternalServerError)
		return
	}
	handler.render.JSON(w, http.StatusOK, map[string]any{})
}

func (handler *TaskHandler) MarkTaskDone(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	if id == "" {
		handler.handleError(w, errors.New("id is empty"), http.StatusBadRequest)
		return
	}

	i, err := strconv.Atoi(id)
	if err != nil {
		handler.handleError(w, err, http.StatusBadRequest)
		return
	}

	err = handler.taskService.MarkDone(int64(i))
	if err != nil {
		handler.handleError(w, err, http.StatusBadRequest)
		return
	}

	handler.render.JSON(w, http.StatusOK, map[string]any{})
}
