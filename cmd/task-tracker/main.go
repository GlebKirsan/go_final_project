package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/pkg/errors"

	"github.com/unrolled/render"

	"github.com/GlebKirsan/go-final-project/internal/config"
	"github.com/GlebKirsan/go-final-project/internal/database"
	"github.com/GlebKirsan/go-final-project/internal/handlers"
	"github.com/GlebKirsan/go-final-project/internal/logger"
	"github.com/GlebKirsan/go-final-project/internal/middleware"
	"github.com/GlebKirsan/go-final-project/internal/service"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	cfg := config.Get()

	logger := logger.Get()

	storage, err := database.ConnectDB()
	if err != nil {
		return errors.Wrap(err, "error when creating storage")
	}
	defer storage.Db.Close()

	manager, err := service.NewManager(storage)
	if err != nil {
		return errors.Wrap(err, "error when creating manager")
	}

	render := render.New()
	dateHandler := handlers.NewDateHandler(manager, render, logger)
	taskHandler := handlers.NewTaskHandler(manager, render, logger)
	authHandler := handlers.NewAuthHandler(manager, logger, render)

	r := chi.NewRouter()
	r.Handle("/*", http.FileServer(http.Dir("./web")))
	r.Route("/api", func(r chi.Router) {
		r.Post("/signin", authHandler.Signin)
		r.Get("/nextdate", dateHandler.GetNextDate)
		r.Route("/", func(r chi.Router) {
			r.Use(middleware.Auth)
			r.Route("/task", func(r chi.Router) {
				r.Post("/", taskHandler.PostTask)
				r.Get("/", taskHandler.GetTask)
				r.Put("/", taskHandler.UpdateTask)
				r.Delete("/", taskHandler.DeleteTask)
				r.Post("/done", taskHandler.MarkTaskDone)
			})
			r.Get("/tasks", taskHandler.GetTasks)
		})
	})
	err = http.ListenAndServe("localhost:"+cfg.Port, r)
	if err != nil {
		return err
	}
	return nil
}
