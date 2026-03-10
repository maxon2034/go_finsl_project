package db

import (
	"database/sql"
	"fmt"
)

type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

func AddTask(task Task) (int64, error) {
	if DB == nil {
		return 0, fmt.Errorf("database not initialized")
	}

	query := `INSERT INTO scheduler (date, title, comment, repeat) VALUES (:date, :title, :comment, :repeat)`
	res, err := DB.Exec(query, sql.Named("date", task.Date), sql.Named("title", task.Title), sql.Named("comment", task.Comment), sql.Named("repeat", task.Repeat))
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func GetTasks(limit int) ([]Task, error) {
	if DB == nil {
		return []Task{}, fmt.Errorf("database not initialized")
	}

	if limit <= 0 || limit > 50 {
		return []Task{}, fmt.Errorf("wrong amount of tasks")
	}

	tasks := []Task{}

	rows, err := DB.Query("SELECT id, date, title, comment, repeat FROM scheduler LIMIT ?", limit)
	if err != nil {
		return tasks, err
	}
	defer rows.Close()
	for rows.Next() {
		t := Task{}
		err := rows.Scan(&t.ID, &t.Date, &t.Title, &t.Comment, &t.Repeat)
		if err != nil {
			return tasks, err
		}
		tasks = append(tasks, t)
	}

	if err := rows.Err(); err != nil {
		return tasks, err
	}

	return tasks, nil
}

func GetTask(id string) (*Task, error) {
	task := Task{}

	if id == "" || id == "<nil>" {
		return nil, fmt.Errorf("Error in id")
	}
	err := DB.QueryRow("SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?", id).Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err != nil {
		return nil, err
	}

	return &task, nil
}

func UpdateTask(task *Task) error {
	query := "UPDATE scheduler SET date = ?, title = ?, comment = ?, repeat = ? WHERE id = ?"

	res, err := DB.Exec(query, task.Date, task.Title, task.Comment, task.Repeat, task.ID)
	if err != nil {
		return err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf("no rows were affected")
	}

	return nil
}

func UpdateDate(task *Task) error {
	query := "UPDATE scheduler SET date = ? WHERE id = ?"

	res, err := DB.Exec(query, task.Date, task.ID)
	if err != nil {
		return err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf("no rows were affected")
	}

	return nil
}

func DeleteTask(id string) error {
	if id == "" || id == "<nil>" {
		return fmt.Errorf("Error in id")
	}
	res, err := DB.Exec("DELETE FROM scheduler WHERE id = ?", id)
	if err != nil {
		return err
	}
	count, err := res.RowsAffected()
	if count == 0 {
		return fmt.Errorf("no rows were affected")
	}
	return err
}
