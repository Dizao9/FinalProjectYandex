package api

import (
	"net/http"

	"main.go/pkg/database"
)

type TasksResp struct {
	Tasks []*database.Task `json:"tasks"`
}

func GetTasksHandler(w http.ResponseWriter, r *http.Request) {
	opts := database.GetTaskOptions{
		SearchString: r.URL.Query().Get("search"),
		Limit:        50,
	}
	tasks, err := database.GetTasks(opts)
	if err != nil {
		writeJSONResp(w, http.StatusInternalServerError, map[string]string{"error": "ошибка при получении ближайших задач из бд: " + err.Error()})
		return
	}

	response := TasksResp{Tasks: tasks}
	writeJSONResp(w, http.StatusOK, response)
}
