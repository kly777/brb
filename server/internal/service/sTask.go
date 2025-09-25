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
	HaveID(id uint) bool
	GetAll() ([]*entity.Task, error)
	GetByID(id uint) (*entity.Task, error)
	Update(task *entity.Task) error
	Delete(id uint) error
	DeleteByEventID(eventID uint) error
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

// GetAllTasks 获取所有task
func (s *taskService) GetAllTasks() ([]*entity.Task, error) {
	return s.taskRepo.GetAll()
}

// GetTaskByID 根据ID获取task
func (s *taskService) GetTaskByID(id uint) (*entity.Task, error) {
	return s.taskRepo.GetByID(id)
}

// UpdateTask 更新task
func (s *taskService) UpdateTask(task *entity.Task) error {
	return s.taskRepo.Update(task)
}

// DeleteTask 删除task，并级联删除相关的todos
func (s *taskService) DeleteTask(id uint) error {
	// 先删除所有相关的todos
	err := s.todoRepo.DeleteByTaskID(id)
	if err != nil {
		return fmt.Errorf("failed to delete related todos: %w", err)
	}

	// 然后删除task
	return s.taskRepo.Delete(id)
}
