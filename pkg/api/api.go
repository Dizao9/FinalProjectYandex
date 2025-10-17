package api

import (
	"github.com/go-chi/chi/v4"
)

func Init(r chi.Router) {
	r.Get("/api/nextdate", NextDateHandler)

	r.Get("/api/tasks", GetTasksHandler)

	r.Get("/api/task", GetTaskByIdHandler)
	r.Post("/api/task", AddTaskHandler)
	r.Put("/api/task", PutTaskHandler)
	r.Delete("/api/task", DeleteTaskHandler)

	r.Post("/api/task/done", TaskIsDoneHandler)
}
