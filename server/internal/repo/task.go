package repo

import (
	"database/sql"
	"fmt"

	"brb/internal/entity"
)

type taskRepo struct {
	base *BaseRepo[entity.Task]
}

func NewTaskRepo(db *sql.DB) (*taskRepo, error) {
	// 初始化数据库表
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS tasks (
			id TEXT PRIMARY KEY,
			event_id TEXT NOT NULL,
			sub_task_id TEXT, -- 父任务ID，如果此任务是子任务
			description TEXT NOT NULL,
			start_time DATETIME NOT NULL,
			end_time DATETIME NOT NULL,
			estimate_time BIGINT NOT NULL
		)
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to create tasks table: %w", err)
	}

	baseRepo := NewBaseRepo[entity.Task](db, "tasks")
	return &taskRepo{base: baseRepo}, nil
}

// Create 创建新的task记录
func (r *taskRepo) Create(task *entity.Task) error {
	fields := map[string]interface{}{
		"id":            task.ID,
		"event_id":      task.EventID,
		"sub_task_id":   task.MainTaskID,
		"description":   task.Description,
		"start_time":    task.StartTime,
		"end_time":      task.EndTime,
		"estimate_time": task.EstimateTime,
	}

	_, err := r.base.Create(fields)
	return err
}

// HaveID 检查是否存在指定ID的task
func (r *taskRepo) HaveID(id string) bool {
	var exists bool
	err := r.base.db.QueryRow("SELECT EXISTS(SELECT 1 FROM tasks WHERE id = ?)", id).Scan(&exists)
	if err != nil {
		return false
	}
	return exists
}

// GetByID 根据ID获取task
func (r *taskRepo) GetByID(id string) (*entity.Task, error) {
	task, err := r.base.FindByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("task not found")
		}
		return nil, err
	}
	return task, nil
}

// Update 更新task记录
func (r *taskRepo) Update(task *entity.Task) error {
	fields := map[string]interface{}{
		"event_id":      task.EventID,
		"sub_task_id":   task.MainTaskID,
		"description":   task.Description,
		"start_time":    task.StartTime,
		"end_time":      task.EndTime,
		"estimate_time": task.EstimateTime,
	}

	return r.base.Update(task.ID, fields)
}

// Delete 删除task记录
func (r *taskRepo) Delete(id string) error {
	return r.base.Delete(id)
}

// DeleteByEventID 根据eventID删除所有相关的tasks
func (r *taskRepo) DeleteByEventID(eventID string) error {
	query := "DELETE FROM tasks WHERE event_id = ?"
	_, err := r.base.db.Exec(query, eventID)
	return err
}
