package api

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/aleakimova/yandexpr-final/internal/db"
)

func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	id := r.FormValue("id")
	if id == "" {
		slog.Warn("deleteTask: missing id")
		data, err := json.Marshal(map[string]string{"error": "Не указан идентификатор"})
		if err != nil {
			slog.Error("deleteTask: failed to marshal response", "error", err)
			http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err.Error()), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		if _, err = w.Write(data); err != nil {
			slog.Error("deleteTask: failed to write response", "error", err)
		}
		return
	}

	if _, err := strconv.ParseInt(id, 10, 64); err != nil {
		slog.Warn("deleteTask: invalid id format", "id", id)
		data, err := json.Marshal(map[string]string{"error": "invalid id"})
		if err != nil {
			slog.Error("deleteTask: failed to marshal response", "error", err)
			http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err.Error()), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		if _, err = w.Write(data); err != nil {
			slog.Error("deleteTask: failed to write response", "error", err)
		}
		return
	}

	if err := db.DeleteTask(id); err != nil {
		slog.Error("deleteTask: DB error", "id", id, "error", err)
		data, err := json.Marshal(map[string]string{"error": err.Error()})
		if err != nil {
			slog.Error("deleteTask: failed to marshal response", "error", err)
			http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err.Error()), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNotFound)
		if _, err = w.Write(data); err != nil {
			slog.Error("deleteTask: failed to write response", "error", err)
		}
		return
	}

	slog.Info("task deleted", "id", id)
	data, err := json.Marshal(map[string]string{})
	if err != nil {
		slog.Error("deleteTask: failed to marshal response", "error", err)
		http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err.Error()), http.StatusInternalServerError)
		return
	}
	if _, err = w.Write(data); err != nil {
		slog.Error("deleteTask: failed to write response", "error", err)
	}
}
