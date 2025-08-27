package repo

import (
	"database/sql"
	"fmt"

	"brb/internal/entity"
)

type signRepo struct {
	base *BaseRepo[entity.Sign]
}

func NewSignRepo(db *sql.DB) (*signRepo, error) {
	// 初始化数据库表
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS signs (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			signifier TEXT NOT NULL,
			signified TEXT NOT NULL
		)
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to create signs table: %w", err)
	}

	baseRepo := NewBaseRepo[entity.Sign](db, "signs")
	return &signRepo{base: baseRepo}, nil
}

// Create 创建新的sign记录
func (r *signRepo) Create(sign *entity.Sign) error {
	fields := map[string]interface{}{
		"signifier": sign.Signifier,
		"signified": sign.Signified,
	}

	result, err := r.base.Create(fields)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	sign.ID = id
	return nil
}

// GetByID 根据ID获取sign
func (r *signRepo) GetByID(id int64) (*entity.Sign, error) {
	sign, err := r.base.FindByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("sign not found")
		}
		return nil, err
	}
	return sign, nil
}

// Update 更新sign记录
func (r *signRepo) Update(sign *entity.Sign) error {
	fields := map[string]interface{}{
		"signifier": sign.Signifier,
		"signified": sign.Signified,
	}

	return r.base.Update(sign.ID, fields)
}

// Delete 删除sign记录
func (r *signRepo) Delete(id int64) error {
	return r.base.Delete(id)
}
