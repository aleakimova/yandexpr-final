package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/aleakimova/yandexpr-final/internal/db"
)

func doneTaskHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	id := r.FormValue("id")
	if id == "" {
		slog.Warn("doneTask: missing id")
		data, err := json.Marshal(map[string]string{"error": "Не указан идентификатор"})
		if err != nil {
			slog.Error("doneTask: failed to marshal response", "error", err)
			http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err.Error()), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		if _, err = w.Write(data); err != nil {
			slog.Error("doneTask: failed to write response", "error", err)
		}
		return
	}

	task, err := db.GetTask(id)
	if err == sql.ErrNoRows {
		slog.Warn("doneTask: task not found", "id", id)
		data, err := json.Marshal(map[string]string{"error": "Задача не найдена"})
		if err != nil {
			slog.Error("doneTask: failed to marshal response", "error", err)
			http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err.Error()), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNotFound)
		if _, err = w.Write(data); err != nil {
			slog.Error("doneTask: failed to write response", "error", err)
		}
		return
	}
	if err != nil {
		slog.Error("doneTask: DB error on GetTask", "id", id, "error", err)
		data, err := json.Marshal(map[string]string{"error": err.Error()})
		if err != nil {
			slog.Error("doneTask: failed to marshal response", "error", err)
			http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err.Error()), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		if _, err = w.Write(data); err != nil {
			slog.Error("doneTask: failed to write response", "error", err)
		}
		return
	}

	if task.Repeat == "" {
		slog.Debug("doneTask: no repeat rule, deleting task", "id", id)
		if err := db.DeleteTask(id); err != nil {
			slog.Error("doneTask: failed to delete task", "id", id, "error", err)
			data, err := json.Marshal(map[string]string{"error": err.Error()})
			if err != nil {
				slog.Error("doneTask: failed to marshal response", "error", err)
				http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err.Error()), http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusInternalServerError)
			if _, err = w.Write(data); err != nil {
				slog.Error("doneTask: failed to write response", "error", err)
			}
			return
		}
		slog.Info("task done and deleted", "id", id, "title", task.Title)
	} else {
		nextDate, err := NextDate(time.Now(), task.Date, task.Repeat)
		if err != nil {
			slog.Warn("doneTask: NextDate failed", "id", id, "repeat", task.Repeat, "error", err)
			data, err := json.Marshal(map[string]string{"error": err.Error()})
			if err != nil {
				slog.Error("doneTask: failed to marshal response", "error", err)
				http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err.Error()), http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusBadRequest)
			if _, err = w.Write(data); err != nil {
				slog.Error("doneTask: failed to write response", "error", err)
			}
			return
		}
		slog.Debug("doneTask: rescheduling task", "id", id, "old_date", task.Date, "next_date", nextDate, "repeat", task.Repeat)
		if err := db.UpdateDate(nextDate, task.ID); err != nil {
			slog.Error("doneTask: failed to update date", "id", id, "error", err)
			data, err := json.Marshal(map[string]string{"error": err.Error()})
			if err != nil {
				slog.Error("doneTask: failed to marshal response", "error", err)
				http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err.Error()), http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusInternalServerError)
			if _, err = w.Write(data); err != nil {
				slog.Error("doneTask: failed to write response", "error", err)
			}
			return
		}
		slog.Info("task done and rescheduled", "id", id, "title", task.Title, "next_date", nextDate)
	}

	data, err := json.Marshal(map[string]string{})
	if err != nil {
		slog.Error("doneTask: failed to marshal response", "error", err)
		http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err.Error()), http.StatusInternalServerError)
		return
	}
	if _, err = w.Write(data); err != nil {
		slog.Error("doneTask: failed to write response", "error", err)
	}
}
