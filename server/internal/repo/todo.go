package repo

import (
	"database/sql"
	"fmt"

	"brb/internal/entity"
)

type todoRepo struct {
	db *sql.DB
}

func NewTodoRepo(db *sql.DB) (*todoRepo, error) {
	// 初始化数据库表
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS todos (
			id TEXT PRIMARY KEY,
			task_id TEXT NOT NULL,
			status TEXT NOT NULL,
			priority INTEGER NOT NULL,
			completed_time DATETIME,
			start_time DATETIME NOT NULL,
			end_time DATETIME NOT NULL
		)
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to create todos table: %w", err)
	}

	return &todoRepo{db: db}, nil
}

// Create 创建新的todo记录
func (r *todoRepo) Create(todo *entity.Todo) error {
	_, err := r.db.Exec(
		"INSERT INTO todos (id, task_id, status, priority, completed_time, start_time, end_time) VALUES (?, ?, ?, ?, ?, ?, ?)",
		todo.ID, todo.TaskID, todo.Status, todo.Priority, todo.CompletedTime, todo.StartTime, todo.EndTime,
	)
	return err
}

// GetAll 获取所有todo记录
func (r *todoRepo) GetAll() ([]*entity.Todo, error) {
	rows, err := r.db.Query("SELECT id, task_id, status, priority, completed_time, start_time, end_time FROM todos")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var todos []*entity.Todo
	for rows.Next() {
		todo := &entity.Todo{}
		if err := rows.Scan(&todo.ID, &todo.TaskID, &todo.Status, &todo.Priority, &todo.CompletedTime, &todo.StartTime, &todo.EndTime); err != nil {
			return nil, err
		}
		todos = append(todos, todo)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return todos, nil
}

// GetByID 根据ID获取todo
func (r *todoRepo) GetByID(id string) (*entity.Todo, error) {
	todo := &entity.Todo{}
	err := r.db.QueryRow(
		"SELECT id, task_id, status, priority, completed_time, start_time, end_time FROM todos WHERE id = ?",
		id,
	).Scan(&todo.ID, &todo.TaskID, &todo.Status, &todo.Priority, &todo.CompletedTime, &todo.StartTime, &todo.EndTime)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("todo not found")
		}
		return nil, err
	}
	return todo, nil
}

// Update 更新todo记录
func (r *todoRepo) Update(todo *entity.Todo) error {
	_, err := r.db.Exec(
		"UPDATE todos SET task_id = ?, status = ?, priority = ?, completed_time = ?, start_time = ?, end_time = ? WHERE id = ?",
		todo.TaskID, todo.Status, todo.Priority, todo.CompletedTime, todo.StartTime, todo.EndTime, todo.ID,
	)
	return err
}

// Delete 删除todo记录
func (r *todoRepo) Delete(id string) error {
	_, err := r.db.Exec("DELETE FROM todos WHERE id = ?", id)
	return err
}

// DeleteByTaskID 根据taskID删除所有相关的todos
func (r *todoRepo) DeleteByTaskID(taskID string) error {
	_, err := r.db.Exec("DELETE FROM todos WHERE task_id = ?", taskID)
	return err
}
