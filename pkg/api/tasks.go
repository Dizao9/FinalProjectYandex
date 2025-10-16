package api

import (
	"net/http"

	"main.go/pkg/database"
)

type TasksResp struct {
	Tasks []*database.Task `json:"tasks"`
}

func GetTasksHandler(w http.ResponseWriter, r *http.Request) {
	tasks, err := database.GetTasks(50)
	if err != nil {
		writeJSONResp(w, http.StatusInternalServerError, map[string]string{"error": "ошибка при получении ближайших задач из бд: " + err.Error()})
		return
	}

	response := TasksResp{Tasks: tasks}
	writeJSONResp(w, http.StatusOK, response)
}
