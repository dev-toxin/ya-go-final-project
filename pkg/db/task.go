package db

import (
	"fmt"
	"time"
)

// Task описывает задачу планировщика.
type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

// AddTask добавляет задачу в таблицу scheduler и возвращает её идентификатор.
func AddTask(task *Task) (int64, error) {
	if DB == nil {
		return 0, fmt.Errorf("база данных не инициализирована")
	}

	result, err := DB.Exec(
		`INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)`,
		task.Date,
		task.Title,
		task.Comment,
		task.Repeat,
	)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

// GetTask возвращает задачу по её идентификатору.
func GetTask(id string) (*Task, error) {
	if DB == nil {
		return nil, fmt.Errorf("база данных не инициализирована")
	}
	if id == "" {
		return nil, fmt.Errorf("не указан идентификатор")
	}

	task := new(Task)
	err := DB.QueryRow(
		`SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?`, id,
	).Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err != nil {
		return nil, err
	}
	return task, nil
}

// UpdateTask сохраняет изменённые параметры существующей задачи.
func UpdateTask(task *Task) error {
	if DB == nil {
		return fmt.Errorf("база данных не инициализирована")
	}
	if task.ID == "" {
		return fmt.Errorf("не указан идентификатор")
	}

	result, err := DB.Exec(
		`UPDATE scheduler SET date = ?, title = ?, comment = ?, repeat = ? WHERE id = ?`,
		task.Date, task.Title, task.Comment, task.Repeat, task.ID,
	)
	if err != nil {
		return err
	}
	count, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf("задача не найдена")
	}
	return nil
}

// DeleteTask удаляет задачу по идентификатору.
func DeleteTask(id string) error {
	if DB == nil {
		return fmt.Errorf("база данных не инициализирована")
	}
	if id == "" {
		return fmt.Errorf("не указан идентификатор")
	}

	result, err := DB.Exec(`DELETE FROM scheduler WHERE id = ?`, id)
	if err != nil {
		return err
	}
	return checkAffectedRows(result)
}

// UpdateDate изменяет дату выполнения существующей задачи.
func UpdateDate(date, id string) error {
	if DB == nil {
		return fmt.Errorf("база данных не инициализирована")
	}
	if id == "" {
		return fmt.Errorf("не указан идентификатор")
	}

	result, err := DB.Exec(`UPDATE scheduler SET date = ? WHERE id = ?`, date, id)
	if err != nil {
		return err
	}
	return checkAffectedRows(result)
}

func checkAffectedRows(result interface{ RowsAffected() (int64, error) }) error {
	count, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf("задача не найдена")
	}
	return nil
}

// Tasks возвращает ближайшие задачи по возрастанию даты. Поиск ограничивает
// выборку точной датой или подстрокой в заголовке и комментарии.
func Tasks(limit int, search string) ([]*Task, error) {
	if DB == nil {
		return nil, fmt.Errorf("база данных не инициализирована")
	}
	if limit < 1 {
		return []*Task{}, nil
	}

	query := `SELECT id, date, title, comment, repeat FROM scheduler`
	args := make([]any, 0, 3)
	if search != "" {
		if _, err := time.Parse("20060102", search); err == nil {
			query += ` WHERE date = ?`
			args = append(args, search)
		} else {
			query += ` WHERE title LIKE ? OR comment LIKE ?`
			pattern := "%" + search + "%"
			args = append(args, pattern, pattern)
		}
	}
	query += ` ORDER BY date LIMIT ?`
	args = append(args, limit)

	rows, err := DB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tasks := make([]*Task, 0)
	for rows.Next() {
		task := new(Task)
		if err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat); err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return tasks, nil
}
