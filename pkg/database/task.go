package database

import (
	"fmt"
)

type Task struct {
	ID      string `json:"id,omitempty"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

func AddTask(task Task) (int64, error) {
	var id int64
	query := `INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)`

	res, err := DB.Exec(query,
		task.Date, task.Title, task.Comment, task.Repeat)
	if err != nil {
		return 0, fmt.Errorf("Ошибка при выполнении запроса на размещение в бд: %v", err)
	}

	id, err = res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("Ошибка при получении последнего добавленного айди: %v", err)
	}
	return id, nil
}

func GetTasks(limit int) ([]*Task, error) {
	rows, err := DB.Query("SELECT id, date, title, comment, repeat FROM scheduler ORDER BY date ASC LIMIT ?",
		limit)
	if err != nil {
		return nil, fmt.Errorf("Ошибка при получении задач из БД: %v", err)
	}

	defer rows.Close()

	tasks := make([]*Task, 0)
	for rows.Next() {
		var task Task

		if err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat); err != nil {
			return nil, err
		}

		tasks = append(tasks, &task)
	}

	return tasks, nil
}
