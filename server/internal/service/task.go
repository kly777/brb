package service

import (
	"brb/internal/entity"
	"fmt"
)

// taskService 实现handler.taskService接口
type taskService struct {
	taskRepo taskRepository
	todoRepo todoRepository
}

type taskRepository interface {
	Create(task *entity.Task) error
	GetByID(id string) (*entity.Task, error)
	Update(task *entity.Task) error
	Delete(id string) error
	DeleteByEventID(eventID string) error
}

// NewTaskService 创建新的TaskService实例
func NewTaskService(taskRepo taskRepository, todoRepo todoRepository) *taskService {
	return &taskService{
		taskRepo: taskRepo,
		todoRepo: todoRepo,
	}
}

// CreateTask 创建新的task
func (s *taskService) CreateTask(task *entity.Task) error {
	return s.taskRepo.Create(task)
}

// GetTaskByID 根据ID获取task
func (s *taskService) GetTaskByID(id string) (*entity.Task, error) {
	return s.taskRepo.GetByID(id)
}

// UpdateTask 更新task
func (s *taskService) UpdateTask(task *entity.Task) error {
	return s.taskRepo.Update(task)
}

// DeleteTask 删除task，并级联删除相关的todos
func (s *taskService) DeleteTask(id string) error {
	// 先删除所有相关的todos
	err := s.todoRepo.DeleteByTaskID(id)
	if err != nil {
		return fmt.Errorf("failed to delete related todos: %w", err)
	}

	// 然后删除task
	return s.taskRepo.Delete(id)
}