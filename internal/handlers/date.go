package handlers

import (
	"net/http"
	"time"

	"github.com/unrolled/render"

	"github.com/GlebKirsan/go-final-project/internal/logger"
	"github.com/GlebKirsan/go-final-project/internal/service"
)

type DateHandler struct {
	services *service.Manager
	render   *render.Render
	logger   *logger.Logger
}

func NewDateHandler(manager *service.Manager, render *render.Render, logger *logger.Logger) *DateHandler {
	return &DateHandler{
		services: manager,
		render:   render,
		logger:   logger,
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
	next, err := handler.services.Date.NextDate(n, date, repeat)
	if err != nil {
		handler.handleError(w, err, http.StatusBadRequest)
		return
	}
	handler.render.Text(w, http.StatusOK, next)
}
