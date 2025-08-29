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
	GetByID(id uint) (*entity.Todo, error)
	Update(todo *entity.Todo) error
	Delete(id uint) error
	DeleteByTaskID(taskID uint) error
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
	if !s.taskRepo.HaveID(todo.TaskID) {
		return fmt.Errorf("关联的Task不存在")
	}

	// 验证Todo时间范围是否在Task的时间范围内
	task, err := s.taskRepo.GetByID(todo.TaskID)
	if err != nil {
		return fmt.Errorf("failed to get task %d: %w", todo.TaskID, err)
	}

	// 检查计划时间是否在任务时间范围内
	if todo.PlannedTime.Start != nil && task.PlannedDuration.Start != nil && task.PlannedDuration.End != nil {
		if todo.PlannedTime.Start.Before(*task.PlannedDuration.Start) || todo.PlannedTime.Start.After(*task.PlannedDuration.End) {
			return fmt.Errorf("todo planned start time must be within task time range")
		}
	}

	if todo.PlannedTime.End != nil && task.PlannedDuration.Start != nil && task.PlannedDuration.End != nil {
		if todo.PlannedTime.End.Before(*task.PlannedDuration.Start) || todo.PlannedTime.End.After(*task.PlannedDuration.End) {
			return fmt.Errorf("todo planned end time must be within task time range")
		}
	}

	if todo.PlannedTime.Start != nil && todo.PlannedTime.End != nil {
		if todo.PlannedTime.End.Before(*todo.PlannedTime.Start) {
			return fmt.Errorf("todo planned end time cannot be before start time")
		}
	}

	// 检查实际时间是否在任务时间范围内
	if todo.ActualTime.Start != nil && task.PlannedDuration.Start != nil && task.PlannedDuration.End != nil {
		if todo.ActualTime.Start.Before(*task.PlannedDuration.Start) || todo.ActualTime.Start.After(*task.PlannedDuration.End) {
			return fmt.Errorf("todo actual start time must be within task time range")
		}
	}

	if todo.ActualTime.End != nil && task.PlannedDuration.Start != nil && task.PlannedDuration.End != nil {
		if todo.ActualTime.End.Before(*task.PlannedDuration.Start) || todo.ActualTime.End.After(*task.PlannedDuration.End) {
			return fmt.Errorf("todo actual end time must be within task time range")
		}
	}

	if todo.ActualTime.Start != nil && todo.ActualTime.End != nil {
		if todo.ActualTime.End.Before(*todo.ActualTime.Start) {
			return fmt.Errorf("todo actual end time cannot be before start time")
		}
	}

	return s.todoRepo.Create(todo)
}

// GetAllTodo 获取所有todo
func (s *todoService) GetAllTodo() ([]*entity.Todo, error) {
	return s.todoRepo.GetAll()
}

// GetTodoByID 根据ID获取todo
func (s *todoService) GetTodoByID(id uint) (*entity.Todo, error) {
	return s.todoRepo.GetByID(id)
}

// UpdateTodo 更新todo
func (s *todoService) UpdateTodo(todo *entity.Todo) error {
	// 验证Todo时间范围是否在Task的时间范围内
	task, err := s.taskRepo.GetByID(todo.TaskID)
	if err != nil {
		return fmt.Errorf("failed to get task %d: %w", todo.TaskID, err)
	}

	// 检查计划时间是否在任务时间范围内
	if todo.PlannedTime.Start != nil && task.PlannedDuration.Start != nil && task.PlannedDuration.End != nil {
		if todo.PlannedTime.Start.Before(*task.PlannedDuration.Start) || todo.PlannedTime.Start.After(*task.PlannedDuration.End) {
			return fmt.Errorf("todo planned start time must be within task time range")
		}
	}

	if todo.PlannedTime.End != nil && task.PlannedDuration.Start != nil && task.PlannedDuration.End != nil {
		if todo.PlannedTime.End.Before(*task.PlannedDuration.Start) || todo.PlannedTime.End.After(*task.PlannedDuration.End) {
			return fmt.Errorf("todo planned end time must be within task time range")
		}
	}

	if todo.PlannedTime.Start != nil && todo.PlannedTime.End != nil {
		if todo.PlannedTime.End.Before(*todo.PlannedTime.Start) {
			return fmt.Errorf("todo planned end time cannot be before start time")
		}
	}

	// 检查实际时间是否在任务时间范围内
	if todo.ActualTime.Start != nil && task.PlannedDuration.Start != nil && task.PlannedDuration.End != nil {
		if todo.ActualTime.Start.Before(*task.PlannedDuration.Start) || todo.ActualTime.Start.After(*task.PlannedDuration.End) {
			return fmt.Errorf("todo actual start time must be within task time range")
		}
	}

	if todo.ActualTime.End != nil && task.PlannedDuration.Start != nil && task.PlannedDuration.End != nil {
		if todo.ActualTime.End.Before(*task.PlannedDuration.Start) || todo.ActualTime.End.After(*task.PlannedDuration.End) {
			return fmt.Errorf("todo actual end time must be within task time range")
		}
	}

	if todo.ActualTime.Start != nil && todo.ActualTime.End != nil {
		if todo.ActualTime.End.Before(*todo.ActualTime.Start) {
			return fmt.Errorf("todo actual end time cannot be before start time")
		}
	}

	return s.todoRepo.Update(todo)
}

// DeleteTodo 删除todo
func (s *todoService) DeleteTodo(id uint) error {
	return s.todoRepo.Delete(id)
}
