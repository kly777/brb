package repo

import (
	"database/sql"
	"fmt"

	"brb/internal/entity"
)

type SignRepo struct {
	db *sql.DB
}

type SignRepository interface {
	Create(sign *entity.Sign) error
	GetByID(id int64) (*entity.Sign, error)
	Update(sign *entity.Sign) error
	Delete(id int64) error
}

func NewSignRepo(db *sql.DB) (*SignRepo, error) {
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

	return &SignRepo{db: db}, nil
}

// Create 创建新的sign记录
func (r *SignRepo) Create(sign *entity.Sign) error {
	res, err := r.db.Exec(
		"INSERT INTO signs (signifier, signified) VALUES (?, ?)",
		sign.Signifier, sign.Signified,
	)
	if err != nil {
		return err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	sign.ID = id
	return nil
}

// GetByID 根据ID获取sign
func (r *SignRepo) GetByID(id int64) (*entity.Sign, error) {
	sign := &entity.Sign{}
	err := r.db.QueryRow(
		"SELECT id, signifier, signified FROM signs WHERE id = ?",
		id,
	).Scan(&sign.ID, &sign.Signifier, &sign.Signified)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("sign not found")
		}
		return nil, err
	}
	return sign, nil
}

// Update 更新sign记录
func (r *SignRepo) Update(sign *entity.Sign) error {
	_, err := r.db.Exec(
		"UPDATE signs SET signifier = ?, signified = ? WHERE id = ?",
		sign.Signifier, sign.Signified, sign.ID,
	)
	return err
}

// Delete 删除sign记录
func (r *SignRepo) Delete(id int64) error {
	_, err := r.db.Exec("DELETE FROM signs WHERE id = ?", id)
	return err
}
