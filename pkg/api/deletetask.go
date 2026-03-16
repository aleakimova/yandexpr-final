package api

import (
	"encoding/json"
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
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Не указан идентификатор"})
		return
	}

	if _, err := strconv.ParseInt(id, 10, 64); err != nil {
		slog.Warn("deleteTask: invalid id format", "id", id)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid id"})
		return
	}

	if err := db.DeleteTask(id); err != nil {
		slog.Error("deleteTask: DB error", "id", id, "error", err)
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	slog.Info("task deleted", "id", id)
	json.NewEncoder(w).Encode(map[string]string{})
}
