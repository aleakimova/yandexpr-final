package api

import (
	"database/sql"
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/aleakimova/yandexpr-final/internal/db"
)

func getTaskHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	id := r.FormValue("id")
	if id == "" {
		slog.Warn("getTask: missing id")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Не указан идентификатор"})
		return
	}

	task, err := db.GetTask(id)
	if err == sql.ErrNoRows {
		slog.Warn("getTask: task not found", "id", id)
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "Задача не найдена"})
		return
	}
	if err != nil {
		slog.Error("getTask: DB error", "id", id, "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	slog.Debug("getTask: task retrieved", "id", id, "title", task.Title)
	json.NewEncoder(w).Encode(task)
}
