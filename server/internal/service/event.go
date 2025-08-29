package service

import (
	"brb/internal/entity"
	"fmt"
)

// eventService 实现handler.eventService接口
type eventService struct {
	eventRepo eventRepository
	taskRepo  taskRepository
}

type eventRepository interface {
	Create(event *entity.Event) error
	GetByID(id uint) (*entity.Event, error)
	Update(event *entity.Event) error
	Delete(id uint) error
}

// NewEventService 创建新的EventService实例
func NewEventService(eventRepo eventRepository, taskRepo taskRepository) *eventService {
	return &eventService{
		eventRepo: eventRepo,
		taskRepo:  taskRepo,
	}
}

// CreateEvent 创建新的event
func (s *eventService) CreateEvent(event *entity.Event) error {
	return s.eventRepo.Create(event)
}

// GetEventByID 根据ID获取event
func (s *eventService) GetEventByID(id uint) (*entity.Event, error) {
	return s.eventRepo.GetByID(id)
}

// UpdateEvent 更新event
func (s *eventService) UpdateEvent(event *entity.Event) error {
	return s.eventRepo.Update(event)
}

// DeleteEvent 删除event，并级联删除相关的tasks
func (s *eventService) DeleteEvent(id uint) error {
	// 先删除所有相关的tasks
	err := s.taskRepo.DeleteByEventID(id)
	if err != nil {
		return fmt.Errorf("failed to delete related tasks: %w", err)
	}

	// 然后删除event
	return s.eventRepo.Delete(id)
}