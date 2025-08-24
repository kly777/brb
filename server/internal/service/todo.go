package service

import (
	"brb/internal/entity"
	"fmt"
)

// todoService 实现handler.todoService接口
type todoService struct {
	todoRepo todoRepository
	taskRepo taskRepository
}

type todoRepository interface {
	Create(todo *entity.Todo) error
	GetAll() ([]*entity.Todo, error)
	GetByID(id string) (*entity.Todo, error)
	Update(todo *entity.Todo) error
	Delete(id string) error
	DeleteByTaskID(taskID string) error
}

// NewTodoService 创建新的TodoService实例
func NewTodoService(todoRepo todoRepository, taskRepo taskRepository) *todoService {
	return &todoService{
		todoRepo: todoRepo,
		taskRepo: taskRepo,
	}
}

// CreateTodo 创建新的todo
func (s *todoService) CreateTodo(todo *entity.Todo) error {
	// 检查关联的Task是否存在
	if !s.taskRepo.HaveID(todo.TaskID)|| todo.TaskID=="" {
		return fmt.Errorf("关联的Task不存在")
	}

	// 验证Todo时间范围是否在Task的时间范围内
	task, err := s.taskRepo.GetByID(todo.TaskID)
	if err != nil {
		return fmt.Errorf("failed to get task %s: %w", todo.TaskID, err)
	}

	// 检查时间是否为nil，如果为nil则跳过时间验证
	if todo.StartTime != nil && task.StartTime != nil && task.EndTime != nil {
		if todo.StartTime.Before(*task.StartTime) || todo.StartTime.After(*task.EndTime) {
			return fmt.Errorf("todo start time must be within task time range")
		}
	}

	if todo.EndTime != nil && task.StartTime != nil && task.EndTime != nil {
		if todo.EndTime.Before(*task.StartTime) || todo.EndTime.After(*task.EndTime) {
			return fmt.Errorf("todo end time must be within task time range")
		}
	}

	if todo.StartTime != nil && todo.EndTime != nil {
		if todo.EndTime.Before(*todo.StartTime) {
			return fmt.Errorf("todo end time cannot be before start time")
		}
	}

	return s.todoRepo.Create(todo)
}

// GetAllTodo 获取所有todo
func (s *todoService) GetAllTodo() ([]*entity.Todo, error) {
	return s.todoRepo.GetAll()
}

// GetTodoByID 根据ID获取todo
func (s *todoService) GetTodoByID(id string) (*entity.Todo, error) {
	return s.todoRepo.GetByID(id)
}

// UpdateTodo 更新todo
func (s *todoService) UpdateTodo(todo *entity.Todo) error {
	// 验证Todo时间范围是否在Task的时间范围内
	task, err := s.taskRepo.GetByID(todo.TaskID)
	if err != nil {
		return fmt.Errorf("failed to get task: %w", err)
	}

	// 检查时间是否为nil，如果为nil则跳过时间验证
	if todo.StartTime != nil && task.StartTime != nil && task.EndTime != nil {
		if todo.StartTime.Before(*task.StartTime) || todo.StartTime.After(*task.EndTime) {
			return fmt.Errorf("todo start time must be within task time range")
		}
	}

	if todo.EndTime != nil && task.StartTime != nil && task.EndTime != nil {
		if todo.EndTime.Before(*task.StartTime) || todo.EndTime.After(*task.EndTime) {
			return fmt.Errorf("todo end time must be within task time range")
		}
	}

	if todo.StartTime != nil && todo.EndTime != nil {
		if todo.EndTime.Before(*todo.StartTime) {
			return fmt.Errorf("todo end time cannot be before start time")
		}
	}

	return s.todoRepo.Update(todo)
}

// DeleteTodo 删除todo
func (s *todoService) DeleteTodo(id string) error {
	return s.todoRepo.Delete(id)
}
