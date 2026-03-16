package api

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const dateFormat = "20060102"

func nextDayHandler(rw http.ResponseWriter, req *http.Request) {
	nowStr := req.FormValue("now")
	dateStr := req.FormValue("date")
	repeat := req.FormValue("repeat")

	slog.Debug("nextdate: computing next date", "now", nowStr, "date", dateStr, "repeat", repeat)

	var now time.Time
	if nowStr == "" {
		now = time.Now()
	} else {
		var err error
		now, err = time.Parse(dateFormat, nowStr)
		if err != nil {
			slog.Warn("nextdate: invalid 'now' parameter", "now", nowStr, "error", err)
			http.Error(rw, "invalid 'now' date format", http.StatusBadRequest)
			return
		}
	}
	next, err := NextDate(now, dateStr, repeat)
	if err != nil {
		slog.Warn("nextdate: NextDate failed", "date", dateStr, "repeat", repeat, "error", err)
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	slog.Debug("nextdate: result", "date", dateStr, "repeat", repeat, "next", next)
	rw.WriteHeader(http.StatusOK)
	_, _ = rw.Write([]byte(next))
}

func NextDate(now time.Time, dstart string, repeat string) (string, error) {
	if repeat == "" {
		return "", errors.New("no repeat rule specified")
	}

	startDate, err := time.Parse(dateFormat, dstart)
	if err != nil {
		return "", fmt.Errorf("invalid dstart date: %w", err)
	}

	parts := strings.Fields(repeat)
	rule := parts[0]
	date := startDate

	switch rule {
	case "y":
		for !(date.After(now) && date.After(startDate)) {
			date = date.AddDate(1, 0, 0)
		}
		return date.Format(dateFormat), nil

	case "d":
		if len(parts) != 2 {
			return "", errors.New("invalid d rule format")
		}
		interval, err := strconv.Atoi(parts[1])
		if err != nil || interval < 1 || interval > 400 {
			return "", errors.New("invalid day interval")
		}
		for !(date.After(now) && date.After(startDate)) {
			date = date.AddDate(0, 0, interval)
		}
		return date.Format(dateFormat), nil

	case "w":
		// if !date.After(now) {
		// 	date = now.Truncate(24 * time.Hour)
		// }
		if len(parts) != 2 {
			return "", errors.New("invalid w rule format")
		}
		daysStr := strings.Split(parts[1], ",")
		weekdays := make([]bool, 7)
		for _, d := range daysStr {
			dayNum, err := strconv.Atoi(d)
			if err != nil || dayNum < 1 || dayNum > 7 {
				return "", errors.New("invalid weekday value")
			}
			// Sunday is 0 in time.Time
			if dayNum == 7 {
				dayNum = int(time.Sunday)
			}
			weekdays[dayNum] = true
		}
		for !(date.After(now) && date.After(startDate) && weekdays[int(date.Weekday())]) {
			date = date.AddDate(0, 0, 1)
		}
		return date.Format(dateFormat), nil

	case "m":
		if len(parts) < 2 {
			return "", errors.New("invalid m rule format")
		}
		daysStr := strings.Split(parts[1], ",")
		var positiveDays [32]bool
		var negativeDays [32]bool
		for _, d := range daysStr {
			dayNum, err := strconv.Atoi(d)
			if err != nil || dayNum == 0 || dayNum < -2 || dayNum > 31 {
				return "", errors.New("invalid month day")
			}
			if dayNum > 0 {
				positiveDays[dayNum] = true
			} else {
				negativeDays[-dayNum] = true
			}
		}
		matchDay := func(d time.Time) bool {
			day := d.Day()
			if positiveDays[day] {
				return true
			}
			lastDay := time.Date(d.Year(), d.Month()+1, 0, 0, 0, 0, 0, d.Location()).Day()
			for neg := 1; neg <= 31; neg++ {
				if negativeDays[neg] && day == lastDay-neg+1 {
					return true
				}
			}
			return false
		}
		var months [13]bool
		if len(parts) == 3 {
			monthsStr := strings.Split(parts[2], ",")
			for _, m := range monthsStr {
				monthNum, err := strconv.Atoi(m)
				if err != nil || monthNum < 1 || monthNum > 12 {
					return "", errors.New("invalid month number")
				}
				months[monthNum] = true
			}
		} else {
			for i := 1; i <= 12; i++ {
				months[i] = true
			}
		}

		for !(date.After(now) && date.After(startDate) && months[int(date.Month())] && matchDay(date)) {
			date = date.AddDate(0, 0, 1)
		}
		return date.Format(dateFormat), nil

	default:
		return "", errors.New("unsupported repeat rule")
	}
}
