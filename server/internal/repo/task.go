package repo

import (
	"database/sql"
	"fmt"

	"brb/internal/entity"
)

type taskRepo struct {
	db *sql.DB
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

	return &taskRepo{db: db}, nil
}

// Create 创建新的task记录
func (r *taskRepo) Create(task *entity.Task) error {
	_, err := r.db.Exec(
		"INSERT INTO tasks (id, event_id, sub_task_id, description, start_time, end_time, estimate_time) VALUES (?, ?, ?, ?, ?, ?, ?)",
		task.ID, task.EventID, task.MainTaskID, task.Description, task.StartTime, task.EndTime, task.EstimateTime,
	)
	return err
}

// GetByID 根据ID获取task
func (r *taskRepo) GetByID(id string) (*entity.Task, error) {
	task := &entity.Task{}
	err := r.db.QueryRow(
		"SELECT id, event_id, sub_task_id, description, start_time, end_time, estimate_time FROM tasks WHERE id = ?",
		id,
	).Scan(&task.ID, &task.EventID, &task.MainTaskID, &task.Description, &task.StartTime, &task.EndTime, &task.EstimateTime)

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
	_, err := r.db.Exec(
		"UPDATE tasks SET event_id = ?, sub_task_id = ?, description = ?, start_time = ?, end_time = ?, estimate_time = ? WHERE id = ?",
		task.EventID, task.MainTaskID, task.Description, task.StartTime, task.EndTime, task.EstimateTime, task.ID,
	)
	return err
}

// Delete 删除task记录
func (r *taskRepo) Delete(id string) error {
	_, err := r.db.Exec("DELETE FROM tasks WHERE id = ?", id)
	return err
}

// DeleteByEventID 根据eventID删除所有相关的tasks
func (r *taskRepo) DeleteByEventID(eventID string) error {
	_, err := r.db.Exec("DELETE FROM tasks WHERE event_id = ?", eventID)
	return err
}
