package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/aleakimova/yandexpr-final/internal/db"
)

func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		slog.Warn("addTask: failed to decode request body", "error", err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	if task.Title == "" {
		slog.Warn("addTask: missing title")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Не указан заголовок задачи"})
		return
	}

	now := time.Now()
	today := now.Format(dateFormat)
	todayTime, _ := time.Parse(dateFormat, today)

	if task.Date == "" {
		slog.Debug("addTask: empty date, using today", "date", today)
		task.Date = today
	}

	parsed, err := time.Parse(dateFormat, task.Date)
	if err != nil {
		slog.Warn("addTask: invalid date format", "date", task.Date, "error", err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": fmt.Sprintf("invalid date: %v", err)})
		return
	}

	if err := checkRepeatRule(task.Repeat); err != nil {
		slog.Warn("addTask: invalid repeat rule", "repeat", task.Repeat, "error", err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	if !parsed.After(now) {
		// past or today
		if task.Repeat != "" && parsed.Before(todayTime) {
			// strictly past date + has repeat → advance to next occurrence
			slog.Debug("addTask: past date with repeat, computing next occurrence", "date", task.Date, "repeat", task.Repeat)
			task.Date, err = NextDate(now, task.Date, task.Repeat)
			if err != nil {
				slog.Warn("addTask: NextDate failed", "repeat", task.Repeat, "error", err)
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
				return
			}
			slog.Debug("addTask: next occurrence computed", "next_date", task.Date)
		} else {
			// today or past without repeat → use today
			slog.Debug("addTask: past/today without repeat, using today", "date", today)
			task.Date = today
		}
	}
	// future date: keep as-is

	id, err := db.AddTask(&task)
	if err != nil {
		slog.Error("addTask: failed to insert task into DB", "title", task.Title, "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	slog.Info("task created", "id", id, "title", task.Title, "date", task.Date)
	json.NewEncoder(w).Encode(map[string]string{"id": id})
}

func checkRepeatRule(repeat string) error {
	if repeat == "" {
		return nil
	}
	parts := strings.Fields(repeat)
	rule := parts[0]
	switch rule {
	case "y":
		if len(parts) != 1 {
			return fmt.Errorf("Invalid number of params for %s rule", rule)
		}

	case "d":
		if len(parts) != 2 {
			return fmt.Errorf("Invalid number of params for %s rule", rule)
		}
		interval, err := strconv.Atoi(parts[1])
		if err != nil || interval < 1 || interval > 400 {
			return fmt.Errorf("invalid day interval")
		}

	case "w":
		if len(parts) != 2 {
			return fmt.Errorf("Invalid number of params for %s rule", rule)
		}
		daysStr := strings.Split(parts[1], ",")
		for _, d := range daysStr {
			dayNum, err := strconv.Atoi(d)
			if err != nil || dayNum < 1 || dayNum > 7 {
				return errors.New("invalid weekday value")
			}
		}

	case "m":
		if len(parts) < 2 || len(parts) > 3 {
			return fmt.Errorf("Invalid number of params for %s rule", rule)
		}
		daysStr := strings.Split(parts[1], ",")
		for _, d := range daysStr {
			dayNum, err := strconv.Atoi(d)
			if err != nil || dayNum == 0 || dayNum < -2 || dayNum > 31 {
				return errors.New("invalid month day")
			}
		}
		if len(parts) == 3 {
			monthsStr := strings.Split(parts[2], ",")
			for _, m := range monthsStr {
				monthNum, err := strconv.Atoi(m)
				if err != nil || monthNum < 1 || monthNum > 12 {
					return errors.New("invalid month number")
				}
			}
		}

	default:
		return errors.New("unsupported repeat rule")
	}
	return nil
}
