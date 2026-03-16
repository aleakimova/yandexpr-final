package api

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"github.com/aleakimova/yandexpr-final/internal/db"
)

func tasksHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	search := r.FormValue("search")
	var (
		tasks []*db.Task
		err   error
	)
	if search != "" {
		if t, parseErr := time.Parse("02.01.2006", search); parseErr == nil {
			slog.Debug("tasks: searching by date", "date", t.Format(dateFormat))
			tasks, err = db.SearchTasksByDate(t.Format(dateFormat), 50)
		} else {
			slog.Debug("tasks: searching by text", "query", search)
			tasks, err = db.SearchTasksByText(search, 50)
		}
	} else {
		slog.Debug("tasks: listing all tasks")
		tasks, err = db.Tasks(50)
	}
	if err != nil {
		slog.Error("tasks: DB error", "search", search, "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	slog.Info("tasks fetched", "count", len(tasks), "search", search)
	json.NewEncoder(w).Encode(map[string]any{"tasks": tasks})
}
