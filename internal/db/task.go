package db

import (
	"database/sql"
	"fmt"
	"strconv"
)

type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

const query = `
		INSERT INTO scheduler (date, title, comment, repeat)
		VALUES (?, ?, ?, ?)
	`

func AddTask(task *Task) (string, error) {
	res, err := Get().Exec(query, task.Date, task.Title, task.Comment, task.Repeat)
	if err != nil {
		return "", err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return "", err
	}
	return strconv.FormatInt(id, 10), nil
}

func GetTask(id string) (*Task, error) {
	var task Task
	var iid int64
	err := Get().QueryRow(
		`SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?`, id,
	).Scan(&iid, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err == sql.ErrNoRows {
		return nil, sql.ErrNoRows
	}
	if err != nil {
		return nil, err
	}
	task.ID = strconv.FormatInt(iid, 10)
	return &task, nil
}

func Tasks(limit int) ([]*Task, error) {
	rows, err := Get().Query(
		`SELECT id, date, title, comment, repeat FROM scheduler ORDER BY date ASC LIMIT ?`,
		limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tasks := []*Task{}
	for rows.Next() {
		var task Task
		var iid int64
		if err := rows.Scan(&iid, &task.Date, &task.Title, &task.Comment, &task.Repeat); err != nil {
			return nil, err
		}
		task.ID = strconv.FormatInt(iid, 10)
		tasks = append(tasks, &task)
	}
	return tasks, rows.Err()
}

func SearchTasksByText(search string, limit int) ([]*Task, error) {
	like := "%" + search + "%"
	rows, err := Get().Query(
		`SELECT id, date, title, comment, repeat FROM scheduler WHERE title LIKE ? OR comment LIKE ? ORDER BY date ASC LIMIT ?`,
		like, like, limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tasks := []*Task{}
	for rows.Next() {
		var task Task
		var iid int64
		if err := rows.Scan(&iid, &task.Date, &task.Title, &task.Comment, &task.Repeat); err != nil {
			return nil, err
		}
		task.ID = strconv.FormatInt(iid, 10)
		tasks = append(tasks, &task)
	}
	return tasks, rows.Err()
}

func SearchTasksByDate(date string, limit int) ([]*Task, error) {
	rows, err := Get().Query(
		`SELECT id, date, title, comment, repeat FROM scheduler WHERE date = ? LIMIT ?`,
		date, limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tasks := []*Task{}
	for rows.Next() {
		var task Task
		var iid int64
		if err := rows.Scan(&iid, &task.Date, &task.Title, &task.Comment, &task.Repeat); err != nil {
			return nil, err
		}
		task.ID = strconv.FormatInt(iid, 10)
		tasks = append(tasks, &task)
	}
	return tasks, rows.Err()
}

func UpdateTask(task *Task) error {
	res, err := Get().Exec(
		`UPDATE scheduler SET date=?, title=?, comment=?, repeat=? WHERE id=?`,
		task.Date, task.Title, task.Comment, task.Repeat, task.ID,
	)
	if err != nil {
		return err
	}
	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return fmt.Errorf("task not found")
	}
	return nil
}

func UpdateDate(nextDate string, id string) error {
	res, err := Get().Exec(
		`UPDATE scheduler SET date=? WHERE id=?`,
		nextDate, id,
	)
	if err != nil {
		return err
	}
	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return fmt.Errorf("task not found")
	}
	return nil
}

func DeleteTask(id string) error {
	res, err := Get().Exec(`DELETE FROM scheduler WHERE id=?`, id)
	if err != nil {
		return err
	}
	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return fmt.Errorf("task not found")
	}
	return nil
}
