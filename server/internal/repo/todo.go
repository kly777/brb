package repo

import (
	"database/sql"
	"fmt"

	"brb/internal/entity"
)

type todoRepo struct {
	base *BaseRepo[entity.Todo]
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

	baseRepo := NewBaseRepo[entity.Todo](db, "todos")
	return &todoRepo{base: baseRepo}, nil
}

// Create 创建新的todo记录
func (r *todoRepo) Create(todo *entity.Todo) error {
	fields := map[string]interface{}{
		"id":             todo.ID,
		"task_id":        todo.TaskID,
		"status":         todo.Status,
		"priority":       todo.Priority,
		"completed_time": todo.CompletedTime,
		"start_time":     todo.StartTime,
		"end_time":       todo.EndTime,
	}

	_, err := r.base.Create(fields)
	return err
}

// GetAll 获取所有todo记录
func (r *todoRepo) GetAll() ([]*entity.Todo, error) {
	return r.base.FindAll()
}

// GetByID 根据ID获取todo
func (r *todoRepo) GetByID(id string) (*entity.Todo, error) {
	todo, err := r.base.FindByID(id)
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
	fields := map[string]interface{}{
		"task_id":        todo.TaskID,
		"status":         todo.Status,
		"priority":       todo.Priority,
		"completed_time": todo.CompletedTime,
		"start_time":     todo.StartTime,
		"end_time":       todo.EndTime,
	}

	return r.base.Update(todo.ID, fields)
}

// Delete 删除todo记录
func (r *todoRepo) Delete(id string) error {
	return r.base.Delete(id)
}

// DeleteByTaskID 根据taskID删除所有相关的todos
func (r *todoRepo) DeleteByTaskID(taskID string) error {
	query := "DELETE FROM todos WHERE task_id = ?"
	_, err := r.base.db.Exec(query, taskID)
	return err
}
