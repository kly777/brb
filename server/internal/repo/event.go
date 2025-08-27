package repo

import (
	"database/sql"
	"fmt"

	"brb/internal/entity"
)

type eventRepo struct {
	base *BaseRepo[entity.Event]
}

func NewEventRepo(db *sql.DB) (*eventRepo, error) {
	// 初始化数据库表
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS events (
			id TEXT PRIMARY KEY,
			title TEXT NOT NULL,
			description TEXT,
			recurrence TEXT
		)
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to create events table: %w", err)
	}

	baseRepo := NewBaseRepo[entity.Event](db, "events")
	return &eventRepo{base: baseRepo}, nil
}

// Create 创建新的event记录
func (r *eventRepo) Create(event *entity.Event) error {
	fields := map[string]interface{}{
		"id":          event.ID,
		"title":       event.Title,
		"description": event.Description,
		"recurrence":  event.Recurrence,
	}

	_, err := r.base.Create(fields)
	return err
}

// GetByID 根据ID获取event
func (r *eventRepo) GetByID(id string) (*entity.Event, error) {
	event, err := r.base.FindByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("event not found")
		}
		return nil, err
	}
	return event, nil
}

// Update 更新event记录
func (r *eventRepo) Update(event *entity.Event) error {
	fields := map[string]interface{}{
		"title":       event.Title,
		"description": event.Description,
		"recurrence":  event.Recurrence,
	}

	return r.base.Update(event.ID, fields)
}

// Delete 删除event记录
func (r *eventRepo) Delete(id string) error {
	return r.base.Delete(id)
}