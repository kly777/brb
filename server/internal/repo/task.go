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
		"event_id":       task.EventID,
		"parent_task_id": task.ParentTaskID,
		"description":    task.Description,
		"status":         string(task.Status),
		"created_at":     task.CreatedAt,
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

// GetAll 获取所有task
func (r *taskRepo) GetAll() ([]*entity.Task, error) {
	query := "SELECT * FROM tasks"
	rows, err := r.base.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []*entity.Task
	for rows.Next() {
		task, err := r.scanTask(rows)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return tasks, nil
}

// GetByID 根据ID获取task
func (r *taskRepo) GetByID(id uint) (*entity.Task, error) {
	query := "SELECT * FROM tasks WHERE id = ?"
	row := r.base.db.QueryRow(query, id)
	return r.scanTask(row)
}

// scanTask 从数据库行扫描Task实体
func (r *taskRepo) scanTask(row interface{}) (*entity.Task, error) {
	var task entity.Task
	var allowedStart, allowedEnd, plannedStart, plannedEnd sql.NullTime
	var parentTaskID sql.NullInt64
	var preTaskIDs sql.NullString

	var err error
	switch row := row.(type) {
	case *sql.Row:
		err = row.Scan(
			&task.ID,
			&task.EventID,
			&parentTaskID,
			&task.Description,
			&allowedStart,
			&allowedEnd,
			&plannedStart,
			&plannedEnd,
			&task.Status,
			&task.CreatedAt,
			&preTaskIDs,
		)
	case *sql.Rows:
		err = row.Scan(
			&task.ID,
			&task.EventID,
			&parentTaskID,
			&task.Description,
			&allowedStart,
			&allowedEnd,
			&plannedStart,
			&plannedEnd,
			&task.Status,
			&task.CreatedAt,
			&preTaskIDs,
		)
	default:
		return nil, fmt.Errorf("unsupported row type")
	}

	if err != nil {
		return nil, err
	}

	// 处理父任务ID
	if parentTaskID.Valid {
		parentID := uint(parentTaskID.Int64)
		task.ParentTaskID = &parentID
	}

	// 处理时间字段
	if allowedStart.Valid {
		task.AllowedTime.Start = &allowedStart.Time
	}
	if allowedEnd.Valid {
		task.AllowedTime.End = &allowedEnd.Time
	}
	if plannedStart.Valid {
		task.PlannedDuration.Start = &plannedStart.Time
	}
	if plannedEnd.Valid {
		task.PlannedDuration.End = &plannedEnd.Time
	}

	// 处理前置任务ID（这里简化处理，实际可能需要JSON解析）
	// 目前preTaskIDs字段在数据库中存储为JSON字符串，但这里我们先忽略具体解析
	// 可以根据实际需要添加JSON解析逻辑
	task.PreTaskIDs = []uint{} // 初始化为空切片

	return &task, nil
}

// Update 更新task记录
func (r *taskRepo) Update(task *entity.Task) error {
	fields := map[string]any{
		"event_id":       task.EventID,
		"parent_task_id": task.ParentTaskID,
		"description":    task.Description,
		"status":         string(task.Status),
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
