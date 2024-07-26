package app

import (
	"net/http"

	"github.com/GlebKirsan/go-final-project/internal/database"
	"github.com/GlebKirsan/go-final-project/internal/env"
	"github.com/go-chi/chi/v5"
)

type App struct {
	Router *chi.Mux
}

func New() (*App, error) {
	a := &App{}

	a.Router = chi.NewRouter()
	database.ConnectDB()

	return a, nil
}

func (a *App) Run() error {
	port := env.GetEnvOrDefault("TODO_PORT", "7540")
	a.Router.Handle("/*", http.FileServer(http.Dir("./web")))
	err := http.ListenAndServe(":"+port, a.Router)
	if err != nil {
		return err
	}

	return nil
}
