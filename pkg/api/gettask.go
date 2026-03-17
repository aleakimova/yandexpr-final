package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/aleakimova/yandexpr-final/internal/db"
)

func getTaskHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	id := r.FormValue("id")
	if id == "" {
		slog.Warn("getTask: missing id")
		data, err := json.Marshal(map[string]string{"error": "Не указан идентификатор"})
		if err != nil {
			slog.Error("getTask: failed to marshal response", "error", err)
			http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err.Error()), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		if _, err = w.Write(data); err != nil {
			slog.Error("getTask: failed to write response", "error", err)
		}
		return
	}

	task, err := db.GetTask(id)
	if err == sql.ErrNoRows {
		slog.Warn("getTask: task not found", "id", id)
		data, err := json.Marshal(map[string]string{"error": "Задача не найдена"})
		if err != nil {
			slog.Error("getTask: failed to marshal response", "error", err)
			http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err.Error()), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNotFound)
		if _, err = w.Write(data); err != nil {
			slog.Error("getTask: failed to write response", "error", err)
		}
		return
	}
	if err != nil {
		slog.Error("getTask: DB error", "id", id, "error", err)
		data, err := json.Marshal(map[string]string{"error": err.Error()})
		if err != nil {
			slog.Error("getTask: failed to marshal response", "error", err)
			http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err.Error()), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		if _, err = w.Write(data); err != nil {
			slog.Error("getTask: failed to write response", "error", err)
		}
		return
	}

	slog.Debug("getTask: task retrieved", "id", id, "title", task.Title)
	data, err := json.Marshal(task)
	if err != nil {
		slog.Error("getTask: failed to marshal response", "error", err)
		http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err.Error()), http.StatusInternalServerError)
		return
	}
	if _, err = w.Write(data); err != nil {
		slog.Error("getTask: failed to write response", "error", err)
	}
}
