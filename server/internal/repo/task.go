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
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			event_id INTEGER NOT NULL,
			parent_task_id INTEGER,
			description TEXT NOT NULL,
			allowed_start DATETIME,
			allowed_end DATETIME,
			planned_start DATETIME,
			planned_end DATETIME,
			status TEXT NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			
			-- 存储pre_task_ids作为JSON数组
			pre_task_ids TEXT
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
		"event_id":      task.EventID,
		"parent_task_id": task.ParentTaskID,
		"description":   task.Description,
		"status":        string(task.Status),
		"created_at":    task.CreatedAt,
	}

	// 处理时间字段
	if task.AllowedTime.Start != nil {
		fields["allowed_start"] = task.AllowedTime.Start
	}
	if task.AllowedTime.End != nil {
		fields["allowed_end"] = task.AllowedTime.End
	}
	if task.PlannedDuration.Start != nil {
		fields["planned_start"] = task.PlannedDuration.Start
	}
	if task.PlannedDuration.End != nil {
		fields["planned_end"] = task.PlannedDuration.End
	}

	// 处理pre_task_ids作为JSON数组
	if len(task.PreTaskIDs) > 0 {
		// 这里需要将uint切片转换为JSON字符串
		// 在实际实现中可能需要使用json.Marshal
		fields["pre_task_ids"] = "[]" // 临时占位符，需要实际实现
	}
	
	// If ID is set (for updates), include it, otherwise it will be auto-generated
	if task.ID != 0 {
		fields["id"] = task.ID
	}

	_, err := r.base.Create(fields)
	return err
}

// HaveID 检查是否存在指定ID的task
func (r *taskRepo) HaveID(id uint) bool {
	var exists bool
	err := r.base.db.QueryRow("SELECT EXISTS(SELECT 1 FROM tasks WHERE id = ?)", id).Scan(&exists)
	if err != nil {
		return false
	}
	return exists
}

// GetByID 根据ID获取task
func (r *taskRepo) GetByID(id uint) (*entity.Task, error) {
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
		"parent_task_id": task.ParentTaskID,
		"description":   task.Description,
		"status":        string(task.Status),
	}

	// 处理时间字段
	if task.AllowedTime.Start != nil {
		fields["allowed_start"] = task.AllowedTime.Start
	}
	if task.AllowedTime.End != nil {
		fields["allowed_end"] = task.AllowedTime.End
	}
	if task.PlannedDuration.Start != nil {
		fields["planned_start"] = task.PlannedDuration.Start
	}
	if task.PlannedDuration.End != nil {
		fields["planned_end"] = task.PlannedDuration.End
	}

	// 处理pre_task_ids作为JSON数组
	if len(task.PreTaskIDs) > 0 {
		fields["pre_task_ids"] = "[]" // 临时占位符，需要实际实现
	}

	return r.base.Update(task.ID, fields)
}

// Delete 删除task记录
func (r *taskRepo) Delete(id uint) error {
	return r.base.Delete(id)
}

// DeleteByEventID 根据eventID删除所有相关的tasks
func (r *taskRepo) DeleteByEventID(eventID uint) error {
	query := "DELETE FROM tasks WHERE event_id = ?"
	_, err := r.base.db.Exec(query, eventID)
	return err
}
