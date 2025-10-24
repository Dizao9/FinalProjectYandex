package database

import (
	"database/sql"
	"fmt"
	"strconv"
	"time"
)

const DateFormat = "20060102"

type Task struct {
	ID      string `json:"id,omitempty"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

type GetTaskOptions struct {
	SearchString string
	Limit        int
}

func AddTask(task Task) (int64, error) {
	var id int64
	query := `INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)`

	res, err := DB.Exec(query,
		task.Date, task.Title, task.Comment, task.Repeat)
	if err != nil {
		return 0, fmt.Errorf("could not execute insert query: %w", err)
	}

	id, err = res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("could not get last insert id: %w", err)
	}
	return id, nil
}

func GetTasks(opts GetTaskOptions) ([]*Task, error) {
	var rows *sql.Rows
	var err error
	sliceArg := []interface{}{}
	var queryString string
	limitStr := strconv.Itoa(opts.Limit)
	if opts.SearchString != "" {
		searchDate, err := time.Parse("02.01.2006", opts.SearchString)
		if err != nil {
			queryString = "SELECT id, date, title, comment, repeat FROM scheduler WHERE title LIKE ? OR comment LIKE ? ORDER BY date LIMIT ?"
			searchPattern := "%" + opts.SearchString + "%"
			sliceArg = append(sliceArg, searchPattern, searchPattern, limitStr)
		} else {
			searchDateFormat := searchDate.Format(DateFormat)
			queryString = "SELECT id, date, title, comment, repeat FROM scheduler WHERE date = ? LIMIT ?"
			sliceArg = append(sliceArg, searchDateFormat, limitStr)
		}
	} else {
		queryString = "SELECT id, date, title, comment, repeat FROM scheduler ORDER BY date ASC LIMIT ?"
		sliceArg = append(sliceArg, limitStr)
	}
	rows, err = DB.Query(queryString, sliceArg...)
	defer rows.Close()

	tasks := make([]*Task, 0)
	for rows.Next() {
		var task Task

		if err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat); err != nil {
			return nil, err
		}

		tasks = append(tasks, &task)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("ошибка при итерации по задачам: %v", err)
	}

	return tasks, nil
}

func GetTaskById(id string) (*Task, error) {
	var task Task
	err := DB.QueryRow("SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?",
		id).Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("task with id %s not found", id)
		}
		return nil, fmt.Errorf("failed to get task from db: %w", err)
	}

	return &task, nil
}

func UpdateTask(task Task) error {
	res, err := DB.Exec("UPDATE scheduler SET date = ?, title = ?, comment = ?, repeat = ? WHERE id = ?",
		task.Date, task.Title, task.Comment, task.Repeat, task.ID)
	if err != nil {
		return fmt.Errorf("failed to update task: %w", err)
	}

	rowsAff, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAff == 0 {
		return fmt.Errorf("task with id %s not found for update", task.ID)
	}

	return nil
}

func DeleteTask(id string) error {
	res, err := DB.Exec("DELETE FROM scheduler WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("failed to delete task: %w", err)
	}

	rowsAff, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAff == 0 {
		return fmt.Errorf("task with id %s not found for deletion", id)
	}

	return nil
}

func UpdateTaskDate(id, newDate string) error {
	res, err := DB.Exec("UPDATE scheduler SET date=? WHERE id = ?",
		newDate, id)

	if err != nil {
		return fmt.Errorf("failed to update task date: %w", err)
	}

	rowsAff, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAff == 0 {
		return fmt.Errorf("task with id %s not found for date update", id)
	}

	return nil
}
