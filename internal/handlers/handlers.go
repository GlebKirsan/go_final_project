package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/GlebKirsan/go-final-project/internal/database"
	"github.com/GlebKirsan/go-final-project/internal/date"
	"github.com/GlebKirsan/go-final-project/internal/models"
	"github.com/sirupsen/logrus"
)

var log = logrus.New()

type ErrorResp struct {
	Error string `json:"error"`
}

func JSONError(w http.ResponseWriter, err string, code int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	resp, _ := json.Marshal(ErrorResp{
		Error: err,
	})
	w.Write(resp)
}

func GetNextDate(w http.ResponseWriter, r *http.Request) {
	n, err := time.Parse(date.YYYYMMDD, r.URL.Query().Get("now"))
	if err != nil {
		log.Error(err.Error())
		JSONError(w, err.Error(), http.StatusBadRequest)
		return
	}

	d, err := time.Parse(date.YYYYMMDD, r.URL.Query().Get("date"))
	if err != nil {
		log.Error(err)
		JSONError(w, err.Error(), http.StatusBadRequest)
		return
	}

	repeat := r.URL.Query().Get("repeat")

	next, err := date.NextDate(n, d, repeat)
	if err != nil {
		log.Error(err)
		JSONError(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(next))
}

func PostTask(w http.ResponseWriter, r *http.Request) {
	task := new(models.Task)
	var buf bytes.Buffer

	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		log.Error(err)
		JSONError(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err = json.Unmarshal(buf.Bytes(), &task); err != nil {
		log.Error(err)
		JSONError(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Info(task)

	if task.Title == "" {
		log.Error("title is empty")
		JSONError(w, "title is empty", http.StatusBadRequest)
		return
	}

	now := time.Now()
	if task.Date != "" {
		if parsed, err := time.Parse(date.YYYYMMDD, task.Date); err != nil {
			log.Error("date has to be YYYYMMDD")
			JSONError(w, "date has to be YYYYMMDD", http.StatusBadRequest)
			return
		} else {
			if date.Before(parsed, now) {
				if task.Repeat == "" {
					task.Date = now.Format(date.YYYYMMDD)
				} else {
					task.Date, err = date.NextDate(now, parsed, task.Repeat)
					if err != nil {
						log.Error(err)
						JSONError(w, err.Error(), http.StatusBadRequest)
						return
					}
				}
			}
		}
	} else {
		task.Date = now.Format(date.YYYYMMDD)
	}

	resultDB := database.DB.Db.Create(&task)
	if resultDB.Error != nil {
		log.Error(resultDB.Error.Error())
		JSONError(w, resultDB.Error.Error(), http.StatusBadRequest)
		return
	}

	resp, err := json.Marshal(map[string]any{
		"id": task.ID,
	})
	if err != nil {
		log.Error(err)
		JSONError(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}

type GetTasksResp struct {
	Tasks []models.Task `json:"tasks"`
}

func GetTasks(w http.ResponseWriter, r *http.Request) {
	tasks := []models.Task{}

	search := r.URL.Query().Get("search")
	if search != "" {
		if parsed, err := time.Parse("02.01.2006", search); err == nil {
			resultDB := database.FindByDate(&tasks, parsed.Format(date.YYYYMMDD))
			if resultDB.Error != nil {
				log.Error(resultDB.Error)
				JSONError(w, resultDB.Error.Error(), http.StatusInternalServerError)
				return
			}
		} else {
			resultDB := database.FindByTitle(&tasks, search)
			if resultDB.Error != nil {
				log.Error(resultDB.Error)
				JSONError(w, resultDB.Error.Error(), http.StatusInternalServerError)
				return
			}
		}
	} else {
		resultDB := database.FindAll(&tasks)
		if resultDB.Error != nil {
			log.Error(resultDB.Error)
			JSONError(w, resultDB.Error.Error(), http.StatusInternalServerError)
			return
		}
	}

	resp, err := json.Marshal(GetTasksResp{
		Tasks: tasks,
	})
	if err != nil {
		log.Error(err)
		JSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}

func GetTask(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	if id == "" {
		log.Error("id is empty")
		JSONError(w, "id is empty", http.StatusBadRequest)
		return
	}

	task := models.Task{}
	resultDB := database.FindById(&task, id)
	if resultDB.Error != nil {
		log.Error(resultDB.Error)
		JSONError(w, resultDB.Error.Error(), http.StatusInternalServerError)
		return
	}

	resp, err := json.Marshal(task)
	if err != nil {
		log.Error(err)
		JSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}

func UpdateTask(w http.ResponseWriter, r *http.Request) {
	task := new(models.Task)
	var buf bytes.Buffer

	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		log.Error(err)
		JSONError(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err = json.Unmarshal(buf.Bytes(), &task); err != nil {
		log.Error(err)
		JSONError(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Info("Update: ", task)

	if task.Title == "" {
		log.Error("title is empty")
		JSONError(w, "title is empty", http.StatusBadRequest)
		return
	}

	now := time.Now()
	if task.Date != "" {
		if parsed, err := time.Parse(date.YYYYMMDD, task.Date); err != nil {
			log.Error("date has to be YYYYMMDD")
			JSONError(w, "date has to be YYYYMMDD", http.StatusBadRequest)
			return
		} else {
			task.Date, err = date.NextDate(now, parsed, task.Repeat)
			if err != nil {
				log.Error(err)
				JSONError(w, err.Error(), http.StatusBadRequest)
				return
			}
			if date.Before(parsed, now) {
				if task.Repeat == "" {
					task.Date = now.Format(date.YYYYMMDD)
				}
			}
		}
	} else {
		task.Date = now.Format(date.YYYYMMDD)
	}

	storedTask := models.Task{}
	resultDB := database.FindById(&storedTask, strconv.FormatUint(uint64(task.ID), 10))
	if resultDB.Error != nil {
		log.Error(resultDB.Error.Error())
		JSONError(w, resultDB.Error.Error(), http.StatusBadRequest)
		return
	}
	storedTask.Repeat = task.Repeat
	storedTask.Title = task.Title
	storedTask.Date = task.Date
	storedTask.Comment = task.Comment

	resultDB = database.DB.Db.Save(&storedTask)
	if resultDB.Error != nil {
		log.Error(resultDB.Error.Error())
		JSONError(w, resultDB.Error.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("{}"))
}

func DeleteTask(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	if id == "" {
		log.Error("id is empty")
		JSONError(w, "id is empty", http.StatusBadRequest)
		return
	}

	resultDB := database.DeleteById(id)
	if resultDB.Error != nil {
		log.Error(resultDB.Error)
		JSONError(w, resultDB.Error.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("{}"))
}

func MarkTaskDone(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	if id == "" {
		log.Error("id is empty")
		JSONError(w, "id is empty", http.StatusBadRequest)
		return
	}

	task := models.Task{}
	resultDB := database.FindById(&task, id)
	if resultDB.Error != nil {
		log.Error(resultDB.Error)
		JSONError(w, resultDB.Error.Error(), http.StatusInternalServerError)
		return
	}

	if task.Repeat == "" {
		resultDB = database.DB.Db.Delete(&task)
		if resultDB.Error != nil {
			log.Error(resultDB.Error)
			JSONError(w, resultDB.Error.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("{}"))
		return
	}
	parsed, err := date.Parse(task.Date)
	if err != nil {
		log.Error("stored task is damaged")
		JSONError(w, "stored task is damaged", http.StatusInternalServerError)
		return
	}
	task.Date, err = date.NextDate(time.Now(), parsed, task.Repeat)
	if err != nil {
		log.Error("next date calculation gone wrong")
		JSONError(w, "next date calculation gone wrong", http.StatusInternalServerError)
		return
	}

	resultDB = database.DB.Db.Save(&task)
	if resultDB.Error != nil {
		log.Error(resultDB.Error.Error())
		JSONError(w, resultDB.Error.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("{}"))
}
