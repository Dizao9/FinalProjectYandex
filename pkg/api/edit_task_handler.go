package api

import (
	"encoding/json"
	"net/http"
	"time"

	"main.go/pkg/database"
)

func GetTaskByIdHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		writeJSONResp(w, http.StatusBadRequest, map[string]string{"error": "не указан айди задачи"})
		return
	}

	task, err := database.GetTaskById(id)
	if err != nil {
		writeJSONResp(w, http.StatusNotFound, map[string]string{"error": err.Error()})
		return
	}

	writeJSONResp(w, http.StatusOK, task)
}

func PutTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task database.Task

	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		writeJSONResp(w, http.StatusBadRequest, map[string]string{"error": "Ошибка десериализации JSON: " + err.Error()})
		return
	}

	if task.ID == "" {
		writeJSONResp(w, http.StatusBadRequest, map[string]string{"error": "передан пустой айди"})
		return
	}

	if task.Title == "" {
		writeJSONResp(w, http.StatusBadRequest, map[string]string{"error": "не указан заголовок"})
		return
	}
	now := time.Now()
	if task.Date == "" {
		task.Date = now.Format("20060102")
	}

	date, err := time.Parse("20060102", task.Date)
	if err != nil {
		writeJSONResp(w, http.StatusBadRequest, map[string]string{"error": "неверный формат времени: %v" + err.Error()})
		return
	}

	if task.Repeat != "" {
		nextDate, err := NextDate(now, task.Date, task.Repeat)
		if err != nil {
			writeJSONResp(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
			return
		}

		if afterNow(now, date) {
			task.Date = nextDate
		}

	} else if afterNow(now, date) {
		task.Date = now.Format("20060102")
	}

	if err = database.UpdateTask(task); err != nil {
		writeJSONResp(w, http.StatusNotFound, map[string]string{"error": err.Error()})
		return
	}

	writeJSONResp(w, http.StatusOK, map[string]string{})
}

func TaskIsDoneHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	if id == "" {
		writeJSONResp(w, http.StatusBadRequest, map[string]string{"error": "не указан айди для задачи"})
		return
	}

	task, err := database.GetTaskById(id)
	if err != nil {
		writeJSONResp(w, http.StatusNotFound, map[string]string{"error": err.Error()})
		return
	}

	if task.Repeat == "" {
		if err := database.DeleteTask(task.ID); err != nil {
			writeJSONResp(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}
		writeJSONResp(w, http.StatusOK, map[string]string{})
	} else {
		now := time.Now()

		nextDate, err := NextDate(now, task.Date, task.Repeat)

		if err != nil {
			writeJSONResp(w, http.StatusInternalServerError, map[string]string{"error": "Не удалось рассчитать следующую дату: " + err.Error()})
			return
		}

		if err := database.UpdateTaskDate(id, nextDate); err != nil {
			writeJSONResp(w, http.StatusInternalServerError, map[string]string{"error": "Ошибка при обновлении даты в базе данных: " + err.Error()})
			return
		}
		writeJSONResp(w, http.StatusOK, map[string]string{})
	}

}

func DeleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		writeJSONResp(w, http.StatusBadRequest, map[string]string{"error": "передан пустой айди для задачи"})
		return
	}

	if err := database.DeleteTask(id); err != nil {
		writeJSONResp(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	writeJSONResp(w, http.StatusOK, map[string]string{})
}
