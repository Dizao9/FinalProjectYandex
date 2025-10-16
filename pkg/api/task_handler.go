package api

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"main.go/pkg/database"
)

func writeJSONResp(w http.ResponseWriter, statCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(statCode)
	json.NewEncoder(w).Encode(data)
}

func AddTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task database.Task

	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		writeJSONResp(w, http.StatusBadRequest, map[string]string{"error": "ошибка дессериализации JSON: " + err.Error()})
		return
	}

	if strings.TrimSpace(task.Title) == "" {
		writeJSONResp(w, http.StatusBadRequest, map[string]string{"error": "не указан заголовок для задачи"})
		return
	}

	now := time.Now()
	if task.Date == "" {
		task.Date = now.Format("20060102")
	}

	t, err := time.Parse("20060102", task.Date)
	if err != nil {
		writeJSONResp(w, http.StatusBadRequest, map[string]string{"error": "ошибка при парсинге даты: " + err.Error()})
		return
	}

	if task.Repeat != "" {
		nextDay, err := NextDate(now, task.Date, task.Repeat)
		if err != nil {
			writeJSONResp(w, http.StatusBadRequest, map[string]string{"error": "ошибка при вычислении следующей даты для задачи: " +
				err.Error()})
			return
		}

		if afterNow(now, t) {
			task.Date = nextDay
		}
	} else {
		if afterNow(now, t) {
			task.Date = now.Format("20060102")
		}
	}

	id, err := database.AddTask(task)
	if err != nil {
		writeJSONResp(w, http.StatusInternalServerError, map[string]string{"error": "ошибка при добавлении задачи в БД: " + err.Error()})
		return
	}
	writeJSONResp(w, http.StatusCreated, map[string]int64{"id": id})
}
