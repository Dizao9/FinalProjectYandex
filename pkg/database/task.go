package database

import (
	"database/sql"
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

func GetTaskById(id string) (*Task, error) {
	var task Task
	err := DB.QueryRow("SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?",
		id).Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("задача не найдена")
		}
		return nil, fmt.Errorf("Возникла ошибка в процессе поиска задачи в бд: %v", err)
	}

	return &task, nil
}

func UpdateTask(task Task) error {
	res, err := DB.Exec("UPDATE scheduler SET date = ?, title = ?, comment = ?, repeat = ? WHERE id = ?",
		task.Date, task.Title, task.Comment, task.Repeat, task.ID)
	if err != nil {
		return fmt.Errorf("Возникла ошибка при попытке редактировать задачу: %v", err)
	}

	rowsAff, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAff == 0 {
		return fmt.Errorf("Не найдена задача с id = %s для обновления", task.ID)
	}

	return nil
}

func DeleteTask(id string) error {
	res, err := DB.Exec("DELETE FROM scheduler WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("ошибка при удалении задачи: %v", err)
	}

	rowsAff, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAff == 0 {
		return fmt.Errorf("Задача с id %s не найдена для удаления", id)
	}

	return nil
}

func UpdateTaskDate(id, newDate string) error {
	res, err := DB.Exec("UPDATE scheduler SET date=? WHERE id = ?",
		newDate, id)

	if err != nil {
		return fmt.Errorf("Ошибка при попытке обновить поле date: %v", err)
	}

	rowsAff, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAff == 0 {
		return fmt.Errorf("Задача с id %s не найдена для обновления даты", id)
	}

	return nil
}
