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
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			event_id INTEGER,
			task_id INTEGER NOT NULL,
			status TEXT NOT NULL,
			planned_start DATETIME,
			planned_end DATETIME,
			actual_start DATETIME,
			actual_end DATETIME,
			completed_time DATETIME
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
		"event_id":       todo.EventID,
		"task_id":        todo.TaskID,
		"status":         string(todo.Status),
		"completed_time": todo.CompletedTime,
	}

	// 处理计划时间
	if todo.PlannedTime.Start != nil {
		fields["planned_start"] = todo.PlannedTime.Start
	}
	if todo.PlannedTime.End != nil {
		fields["planned_end"] = todo.PlannedTime.End
	}

	// 处理实际时间
	if todo.ActualTime.Start != nil {
		fields["actual_start"] = todo.ActualTime.Start
	}
	if todo.ActualTime.End != nil {
		fields["actual_end"] = todo.ActualTime.End
	}

	// If ID is set (for updates), include it, otherwise it will be auto-generated
	if todo.ID != 0 {
		fields["id"] = todo.ID
	}

	_, err := r.base.Create(fields)
	return err
}

// GetAll 获取所有todo记录
func (r *todoRepo) GetAll() ([]*entity.Todo, error) {
	return r.base.FindAll()
}

// GetByID 根据ID获取todo
func (r *todoRepo) GetByID(id uint) (*entity.Todo, error) {
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
		"event_id":       todo.EventID,
		"task_id":        todo.TaskID,
		"status":         string(todo.Status),
		"completed_time": todo.CompletedTime,
	}

	// 处理计划时间
	if todo.PlannedTime.Start != nil {
		fields["planned_start"] = todo.PlannedTime.Start
	}
	if todo.PlannedTime.End != nil {
		fields["planned_end"] = todo.PlannedTime.End
	}

	// 处理实际时间
	if todo.ActualTime.Start != nil {
		fields["actual_start"] = todo.ActualTime.Start
	}
	if todo.ActualTime.End != nil {
		fields["actual_end"] = todo.ActualTime.End
	}

	return r.base.Update(todo.ID, fields)
}

// Delete 删除todo记录
func (r *todoRepo) Delete(id uint) error {
	return r.base.Delete(id)
}

// DeleteByTaskID 根据taskID删除所有相关的todos
func (r *todoRepo) DeleteByTaskID(taskID uint) error {
	query := "DELETE FROM todos WHERE task_id = ?"
	_, err := r.base.db.Exec(query, taskID)
	return err
}
