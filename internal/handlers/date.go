package handlers

import (
	"net/http"
	"time"

	"github.com/GlebKirsan/go-final-project/internal/logger"
	"github.com/GlebKirsan/go-final-project/internal/service"
	"github.com/unrolled/render"
)

type DateHandler struct {
	dateService *service.DateService
	render      *render.Render
	logger      *logger.Logger
}

func NewDateHandler(manager *service.Manager, render *render.Render, logger *logger.Logger) *DateHandler {
	return &DateHandler{
		dateService: manager.Date,
		render:      render,
		logger:      logger,
	}
}

func (handler *DateHandler) handleError(w http.ResponseWriter, err error, code int) {
	handler.logger.Error().Msg(err.Error())
	handler.render.JSON(w, code, map[string]any{"error": err.Error()})
}

func (handler *DateHandler) GetNextDate(w http.ResponseWriter, r *http.Request) {
	n, err := time.Parse(service.YYYYMMDD, r.URL.Query().Get("now"))
	if err != nil {
		handler.handleError(w, err, http.StatusBadRequest)
		return
	}

	date := r.URL.Query().Get("date")
	repeat := r.URL.Query().Get("repeat")
	next, err := handler.dateService.NextDate(n, date, repeat)
	if err != nil {
		handler.handleError(w, err, http.StatusBadRequest)
		return
	}
	handler.render.Text(w, http.StatusOK, next)
}
