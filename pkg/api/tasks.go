package api

import (
	"encoding/json"
	"fmt"
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
			tasks, err = db.SearchTasksByDate(t.Format(dateFormat), maxTasks)
		} else {
			slog.Debug("tasks: searching by text", "query", search)
			tasks, err = db.SearchTasksByText(search, maxTasks)
		}
	} else {
		slog.Debug("tasks: listing all tasks")
		tasks, err = db.Tasks(maxTasks)
	}
	if err != nil {
		slog.Error("tasks: DB error", "search", search, "error", err)
		data, err := json.Marshal(map[string]string{"error": err.Error()})
		if err != nil {
			slog.Error("tasks: failed to marshal response", "error", err)
			http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err.Error()), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		if _, err = w.Write(data); err != nil {
			slog.Error("tasks: failed to write response", "error", err)
		}
		return
	}

	slog.Info("tasks fetched", "count", len(tasks), "search", search)
	data, err := json.Marshal(map[string]any{"tasks": tasks})
	if err != nil {
		slog.Error("tasks: failed to marshal response", "error", err)
		http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err.Error()), http.StatusInternalServerError)
		return
	}
	if _, err = w.Write(data); err != nil {
		slog.Error("tasks: failed to write response", "error", err)
	}
}
