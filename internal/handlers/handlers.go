package handlers

import (
	"net/http"
	"time"

	"github.com/GlebKirsan/go-final-project/internal/date"
)

func GetNextDate(w http.ResponseWriter, r *http.Request) {
	n, err := time.Parse(date.YYYYMMDD, r.URL.Query().Get("now"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	d, err := time.Parse(date.YYYYMMDD, r.URL.Query().Get("date"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	repeat := r.URL.Query().Get("repeat")

	next, err := date.NextDate(n, d, repeat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(next))
}
