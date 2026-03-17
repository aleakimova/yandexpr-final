package api

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/aleakimova/yandexpr-final/internal/db"
)

func editTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		slog.Warn("editTask: failed to decode request body", "error", err)
		data, err := json.Marshal(map[string]string{"error": err.Error()})
		if err != nil {
			slog.Error("editTask: failed to marshal response", "error", err)
			http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err.Error()), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		if _, err = w.Write(data); err != nil {
			slog.Error("editTask: failed to write response", "error", err)
		}
		return
	}

	if task.ID == "" {
		slog.Warn("editTask: missing id")
		data, err := json.Marshal(map[string]string{"error": "Не указан идентификатор"})
		if err != nil {
			slog.Error("editTask: failed to marshal response", "error", err)
			http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err.Error()), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		if _, err = w.Write(data); err != nil {
			slog.Error("editTask: failed to write response", "error", err)
		}
		return
	}

	if _, err := strconv.ParseInt(task.ID, 10, 64); err != nil {
		slog.Warn("editTask: invalid id format", "id", task.ID)
		data, err := json.Marshal(map[string]string{"error": "invalid id"})
		if err != nil {
			slog.Error("editTask: failed to marshal response", "error", err)
			http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err.Error()), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		if _, err = w.Write(data); err != nil {
			slog.Error("editTask: failed to write response", "error", err)
		}
		return
	}

	if task.Title == "" {
		slog.Warn("editTask: missing title", "id", task.ID)
		data, err := json.Marshal(map[string]string{"error": "Не указан заголовок задачи"})
		if err != nil {
			slog.Error("editTask: failed to marshal response", "error", err)
			http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err.Error()), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		if _, err = w.Write(data); err != nil {
			slog.Error("editTask: failed to write response", "error", err)
		}
		return
	}

	now := time.Now()
	today := now.Format(dateFormat)
	todayTime, _ := time.Parse(dateFormat, today)

	if task.Date == "" {
		slog.Debug("editTask: empty date, using today", "id", task.ID, "date", today)
		task.Date = today
	}

	parsed, err := time.Parse(dateFormat, task.Date)
	if err != nil {
		slog.Warn("editTask: invalid date format", "id", task.ID, "date", task.Date, "error", err)
		data, err := json.Marshal(map[string]string{"error": fmt.Sprintf("invalid date: %v", err)})
		if err != nil {
			slog.Error("editTask: failed to marshal response", "error", err)
			http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err.Error()), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		if _, err = w.Write(data); err != nil {
			slog.Error("editTask: failed to write response", "error", err)
		}
		return
	}

	if err := checkRepeatRule(task.Repeat); err != nil {
		slog.Warn("editTask: invalid repeat rule", "id", task.ID, "repeat", task.Repeat, "error", err)
		data, err := json.Marshal(map[string]string{"error": err.Error()})
		if err != nil {
			slog.Error("editTask: failed to marshal response", "error", err)
			http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err.Error()), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		if _, err = w.Write(data); err != nil {
			slog.Error("editTask: failed to write response", "error", err)
		}
		return
	}

	if !parsed.After(now) {
		if task.Repeat != "" && parsed.Before(todayTime) {
			slog.Debug("editTask: past date with repeat, computing next occurrence", "id", task.ID, "date", task.Date, "repeat", task.Repeat)
			task.Date, err = NextDate(now, task.Date, task.Repeat)
			if err != nil {
				slog.Warn("editTask: NextDate failed", "id", task.ID, "repeat", task.Repeat, "error", err)
				data, err := json.Marshal(map[string]string{"error": err.Error()})
				if err != nil {
					slog.Error("editTask: failed to marshal response", "error", err)
					http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err.Error()), http.StatusInternalServerError)
					return
				}
				w.WriteHeader(http.StatusBadRequest)
				if _, err = w.Write(data); err != nil {
					slog.Error("editTask: failed to write response", "error", err)
				}
				return
			}
			slog.Debug("editTask: next occurrence computed", "id", task.ID, "next_date", task.Date)
		} else {
			slog.Debug("editTask: past/today without repeat, using today", "id", task.ID, "date", today)
			task.Date = today
		}
	}

	if err := db.UpdateTask(&task); err != nil {
		slog.Error("editTask: failed to update task in DB", "id", task.ID, "error", err)
		data, err := json.Marshal(map[string]string{"error": err.Error()})
		if err != nil {
			slog.Error("editTask: failed to marshal response", "error", err)
			http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err.Error()), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNotFound)
		if _, err = w.Write(data); err != nil {
			slog.Error("editTask: failed to write response", "error", err)
		}
		return
	}

	slog.Info("task updated", "id", task.ID, "title", task.Title, "date", task.Date)
	data, err := json.Marshal(map[string]string{})
	if err != nil {
		slog.Error("editTask: failed to marshal response", "error", err)
		http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err.Error()), http.StatusInternalServerError)
		return
	}
	if _, err = w.Write(data); err != nil {
		slog.Error("editTask: failed to write response", "error", err)
	}
}
