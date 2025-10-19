package api

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const DateFormat = "20060102"

func afterNow(date, now time.Time) bool {
	date = date.Truncate(24 * time.Hour)
	now = now.Truncate(24 * time.Hour)
	return date.After(now)
}

func NextDate(now time.Time, dstart string, repeat string) (string, error) {
	if repeat == "" {
		return "", fmt.Errorf("передана пустая строка в правила повторения")
	}

	date, err := time.Parse(DateFormat, dstart)
	if err != nil {
		return "", fmt.Errorf("неправильно указана дата: %v", err)
	}

	params := strings.Split(repeat, " ")
	if len(params) < 1 {
		return "", fmt.Errorf("неверный формат правила повторения")
	}

	switch params[0] {
	case "d":
		if len(params) < 2 {
			return "", fmt.Errorf("для правила d необходимо указать количество дней")
		}
		interval, err := strconv.Atoi(params[1])
		if err != nil {
			return "", fmt.Errorf("неверный параметр для d: %v", err)
		}
		if interval < 1 || interval > 400 {
			return "", fmt.Errorf("неверное количество дней для d, число должно быть в диапазоне от 1 до 400")
		}
		if date.Before(now) {
			diffDays := int(now.Sub(date).Hours() / 24)
			if diffDays > 0 {
				intervalsToSkip := diffDays / interval
				if intervalsToSkip > 0 {
					date = date.AddDate(0, 0, intervalsToSkip*interval)
				}
			}
		}

		for {
			date = date.AddDate(0, 0, interval)
			if afterNow(date, now) {
				break
			}
		}
		return date.Format(DateFormat), nil

	case "y":
		if len(params) > 1 {
			return "", fmt.Errorf("для правила y не требуется дополнительных параметров")
		}

		if date.Year() < now.Year() {
			yrsToAdd := now.Year() - date.Year()
			date = date.AddDate(yrsToAdd-1, 0, 0)
		}

		for {
			isLeapDay := date.Month() == time.February && date.Day() == 29
			date = date.AddDate(1, 0, 0)
			if isLeapDay && date.Month() == time.February && date.Day() == 28 {
				date = date.AddDate(0, 0, 1)
			}

			if afterNow(date, now) {
				break
			}
		}
		return date.Format(DateFormat), nil

	default:
		return "", fmt.Errorf("неподдерживаемый формат правила: %s", params[0])
	}
}

func NextDateHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")

	nowParam := r.FormValue("now")
	dateParam := r.FormValue("date")
	repeatParam := r.FormValue("repeat")

	if dateParam == "" {
		http.Error(w, "Отсутствует параметр date", http.StatusBadRequest)
		return
	}

	if repeatParam == "" {
		http.Error(w, "Отсутствует правило повтора", http.StatusBadRequest)
		return
	}
	var now time.Time
	var err error
	if nowParam == "" {
		now = time.Now()
	} else {
		now, err = time.Parse(DateFormat, nowParam)
		if err != nil {
			http.Error(w, "Неверный формат параметра now", http.StatusBadRequest)
			return
		}
	}
	nextDate, err := NextDate(now, dateParam, repeatParam)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(nextDate))
}
