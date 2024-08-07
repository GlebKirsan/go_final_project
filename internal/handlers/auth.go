package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/unrolled/render"

	"github.com/GlebKirsan/go-final-project/internal/logger"
	"github.com/GlebKirsan/go-final-project/internal/models"
	"github.com/GlebKirsan/go-final-project/internal/service"
)

type AuthHandler struct {
	services *service.Manager
	logger   *logger.Logger
	render   *render.Render
}

func NewAuthHandler(manager *service.Manager, logger *logger.Logger, render *render.Render) *AuthHandler {
	return &AuthHandler{
		services: manager,
		logger:   logger,
		render:   render,
	}
}

func (handler *AuthHandler) handleError(w http.ResponseWriter, err error, code int) {
	handler.logger.Error().Msg(err.Error())
	handler.render.JSON(w, code, map[string]any{"error": err.Error()})
}

func (handler *AuthHandler) Signin(w http.ResponseWriter, r *http.Request) {
	var buf bytes.Buffer

	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		handler.handleError(w, err, http.StatusUnauthorized)
		return
	}

	var auth models.AuthRequest
	if err := json.Unmarshal(buf.Bytes(), &auth); err != nil {
		handler.handleError(w, err, http.StatusUnauthorized)
		return
	}

	token, err := handler.services.Auth.Authorize(&auth)
	if err != nil {
		handler.handleError(w, err, http.StatusUnauthorized)
		return
	}

	handler.render.JSON(w, http.StatusOK, &models.AuthResponse{Token: token})
}
