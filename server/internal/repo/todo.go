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
	fields := map[string]any{
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
	query := "SELECT id, event_id, task_id, status, planned_start, planned_end, actual_start, actual_end, completed_time FROM todos"
	rows, err := r.base.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query todos: %w", err)
	}
	defer rows.Close()

	var todos []*entity.Todo
	for rows.Next() {
		todo, err := r.scanTodo(rows)
		if err != nil {
			return nil, fmt.Errorf("failed to scan todo: %w", err)
		}
		todos = append(todos, todo)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return todos, nil
}

// GetByID 根据ID获取todo
func (r *todoRepo) GetByID(id uint) (*entity.Todo, error) {
	query := "SELECT id, event_id, task_id, status, planned_start, planned_end, actual_start, actual_end, completed_time FROM todos WHERE id = ?"
	row := r.base.db.QueryRow(query, id)

	todo, err := r.scanTodo(row)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("todo not found")
		}
		return nil, fmt.Errorf("failed to scan todo: %w", err)
	}
	return todo, nil
}

// scanTodo 从数据库行扫描Todo实体
func (r *todoRepo) scanTodo(row any) (*entity.Todo, error) {
	var (
		id            uint
		eventID       sql.NullInt64
		taskID        uint
		status        string
		plannedStart  sql.NullTime
		plannedEnd    sql.NullTime
		actualStart   sql.NullTime
		actualEnd     sql.NullTime
		completedTime sql.NullTime
	)

	var err error
	switch row := row.(type) {
	case *sql.Row:
		err = row.Scan(&id, &eventID, &taskID, &status, &plannedStart, &plannedEnd, &actualStart, &actualEnd, &completedTime)
	case *sql.Rows:
		err = row.Scan(&id, &eventID, &taskID, &status, &plannedStart, &plannedEnd, &actualStart, &actualEnd, &completedTime)
	default:
		return nil, fmt.Errorf("unsupported row type")
	}

	if err != nil {
		return nil, fmt.Errorf("failed to scan todo row: %w", err)
	}

	todo := &entity.Todo{
		ID:     id,
		TaskID: taskID,
		Status: entity.Status(status),
	}

	// 处理可空的event_id
	if eventID.Valid {
		eventIDVal := uint(eventID.Int64)
		todo.EventID = &eventIDVal
	}

	// 处理计划时间
	if plannedStart.Valid {
		todo.PlannedTime.Start = &plannedStart.Time
	}
	if plannedEnd.Valid {
		todo.PlannedTime.End = &plannedEnd.Time
	}

	// 处理实际时间
	if actualStart.Valid {
		todo.ActualTime.Start = &actualStart.Time
	}
	if actualEnd.Valid {
		todo.ActualTime.End = &actualEnd.Time
	}

	// 处理完成时间
	if completedTime.Valid {
		todo.CompletedTime = &completedTime.Time
	}

	return todo, nil
}

// Update 更新todo记录
func (r *todoRepo) Update(todo *entity.Todo) error {
	fields := map[string]any{
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
