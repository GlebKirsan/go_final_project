package app

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

type App struct {
	router *chi.Mux
}

func New() (*App, error) {
	a := &App{}

	a.router = chi.NewRouter()

	return a, nil
}

func (a *App) Run() error {
	a.router.Handle("/*", http.FileServer(http.Dir("./web")))
	err := http.ListenAndServe(":7540", a.router)
	if err != nil {
		return err
	}

	return nil
}
