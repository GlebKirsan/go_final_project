package app

import (
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
)

type App struct {
	router *chi.Mux
}

func GetEnvOrDefault(key string, def string) string {
	if variable := os.Getenv(key); len(variable) != 0 {
		return variable
	}
	return def
}

func New() (*App, error) {
	a := &App{}

	a.router = chi.NewRouter()

	return a, nil
}

func (a *App) Run() error {
	port := GetEnvOrDefault("TODO_PORT", "7540")
	a.router.Handle("/*", http.FileServer(http.Dir("./web")))
	err := http.ListenAndServe(":"+port, a.router)
	if err != nil {
		return err
	}

	return nil
}
