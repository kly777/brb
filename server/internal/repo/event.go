package repo

import (
	"database/sql"
	"fmt"

	"brb/internal/entity"
)

type eventRepo struct {
	db *sql.DB
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

	return &eventRepo{db: db}, nil
}

// Create 创建新的event记录
func (r *eventRepo) Create(event *entity.Event) error {
	_, err := r.db.Exec(
		"INSERT INTO events (id, title, description, recurrence) VALUES (?, ?, ?, ?)",
		event.ID, event.Title, event.Description, event.Recurrence,
	)
	return err
}

// GetByID 根据ID获取event
func (r *eventRepo) GetByID(id string) (*entity.Event, error) {
	event := &entity.Event{}
	err := r.db.QueryRow(
		"SELECT id, title, description, recurrence FROM events WHERE id = ?",
		id,
	).Scan(&event.ID, &event.Title, &event.Description, &event.Recurrence)

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
	_, err := r.db.Exec(
		"UPDATE events SET title = ?, description = ?, recurrence = ? WHERE id = ?",
		event.Title, event.Description, event.Recurrence, event.ID,
	)
	return err
}

// Delete 删除event记录
func (r *eventRepo) Delete(id string) error {
	_, err := r.db.Exec("DELETE FROM events WHERE id = ?", id)
	return err
}