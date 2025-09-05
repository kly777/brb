package repo

import (
	"database/sql"
	"fmt"
	"reflect"
	"strings"
)

// BaseRepo 提供基础的ORM风格CRUD操作，使用泛型提高类型安全性
type BaseRepo[T any] struct {
	db        *sql.DB
	tableName string
}

// NewBaseRepo 创建一个新的BaseRepo实例
func NewBaseRepo[T any](db *sql.DB, tableName string) *BaseRepo[T] {
	return &BaseRepo[T]{
		db:        db,
		tableName: tableName,
	}
}

// Create 插入新记录
func (r *BaseRepo[T]) Create(fields map[string]interface{}) (sql.Result, error) {
	if len(fields) == 0 {
		return nil, fmt.Errorf("no fields to insert")
	}

	columns := make([]string, 0, len(fields))
	placeholders := make([]string, 0, len(fields))
	values := make([]any, 0, len(fields))

	for column, value := range fields {
		columns = append(columns, column)
		placeholders = append(placeholders, "?")
		values = append(values, value)
	}

	query := fmt.Sprintf(
		"INSERT INTO %s (%s) VALUES (%s)",
		r.tableName,
		strings.Join(columns, ", "),
		strings.Join(placeholders, ", "),
	)

	return r.db.Exec(query, values...)
}

// Update 更新记录
func (r *BaseRepo[T]) Update(id any, fields map[string]interface{}) error {
	if len(fields) == 0 {
		return fmt.Errorf("no fields to update")
	}

	setClauses := make([]string, 0, len(fields))
	values := make([]interface{}, 0, len(fields)+1)

	for column, value := range fields {
		setClauses = append(setClauses, fmt.Sprintf("%s = ?", column))
		values = append(values, value)
	}
	values = append(values, id)

	query := fmt.Sprintf(
		"UPDATE %s SET %s WHERE id = ?",
		r.tableName,
		strings.Join(setClauses, ", "),
	)

	_, err := r.db.Exec(query, values...)
	return err
}

// Delete 删除记录
func (r *BaseRepo[T]) Delete(id interface{}) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id = ?", r.tableName)
	_, err := r.db.Exec(query, id)
	return err
}

// FindByID 根据ID查询记录
func (r *BaseRepo[T]) FindByID(id interface{}) (*T, error) {
	query := fmt.Sprintf("SELECT * FROM %s WHERE id = ?", r.tableName)
	row := r.db.QueryRow(query, id)

	var entity T
	err := scanRow(row, &entity)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("record not found")
		}
		return nil, err
	}
	return &entity, nil
}

// FindAll 查询所有记录
func (r *BaseRepo[T]) FindAll() ([]*T, error) {
	return r.FindAllWithFields(nil)
}

// FindAllWithFields 查询所有记录，可指定字段
func (r *BaseRepo[T]) FindAllWithFields(fields []string) ([]*T, error) {
	query := fmt.Sprintf("SELECT * FROM %s", r.tableName)
	if len(fields) > 0 {
		query = fmt.Sprintf("SELECT %s FROM %s", strings.Join(fields, ", "), r.tableName)
	}

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entities []*T
	for rows.Next() {
		var entity T
		err := scanRows(rows, &entity)
		if err != nil {
			return nil, err
		}
		entities = append(entities, &entity)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return entities, nil
}

// FindWhere 条件查询
func (r *BaseRepo[T]) FindWhere(conditions map[string]interface{}) ([]*T, error) {
	return r.FindWhereWithFields(conditions, nil)
}

// FindWhereWithFields 条件查询，可指定字段
func (r *BaseRepo[T]) FindWhereWithFields(conditions map[string]any, fields []string) ([]*T, error) {
	if len(conditions) == 0 {
		return r.FindAllWithFields(fields)
	}

	whereClauses := make([]string, 0, len(conditions))
	values := make([]any, 0, len(conditions))

	for column, value := range conditions {
		whereClauses = append(whereClauses, fmt.Sprintf("%s = ?", column))
		values = append(values, value)
	}

	selectClause := "*"
	if len(fields) > 0 {
		selectClause = strings.Join(fields, ", ")
	}

	query := fmt.Sprintf(
		"SELECT %s FROM %s WHERE %s",
		selectClause,
		r.tableName,
		strings.Join(whereClauses, " AND "),
	)

	rows, err := r.db.Query(query, values...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entities []*T
	for rows.Next() {
		var entity T
		err := scanRows(rows, &entity)
		if err != nil {
			return nil, err
		}
		entities = append(entities, &entity)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return entities, nil
}

// scanRow 辅助函数，用于扫描单行数据到结构体
func scanRow(row *sql.Row, dest any) error {
	return row.Scan(getFieldPointers(dest)...)
}

// scanRows 辅助函数，用于扫描多行数据到结构体
func scanRows(rows *sql.Rows, dest any) error {
	return rows.Scan(getFieldPointers(dest)...)
}

// getFieldPointers 使用反射获取结构体字段的指针
func getFieldPointers(dest any) []any {
	// 使用反射获取字段指针
	// 这里简化实现，实际项目中可能需要更复杂的反射处理
	// 或者使用第三方库如 github.com/jmoiron/sqlx 来处理扫描
	val := reflect.ValueOf(dest).Elem()
	numField := val.NumField()
	fields := make([]any, numField)

	for i := 0; i < numField; i++ {
		fields[i] = val.Field(i).Addr().Interface()
	}

	return fields
}
