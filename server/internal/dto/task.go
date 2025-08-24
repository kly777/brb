package dto

import (
	"brb/internal/entity"
	"time"
)

// TaskCreateRequest DTO for creating a task
type TaskCreateRequest struct {
	EventID      string         `json:"eventId"`
	MainTaskID   *string        `json:"mainTaskId"`
	Description  string         `json:"description"`
	StartTime    *time.Time     `json:"startTime"`
	EndTime      *time.Time     `json:"endTime"`
	EstimateTime *time.Duration `json:"estimateTime"`
}

// TaskUpdateRequest DTO for updating a task
type TaskUpdateRequest struct {
	EventID      string         `json:"eventId"`
	MainTaskID   *string        `json:"mainTaskId"`
	Description  string         `json:"description"`
	StartTime    *time.Time     `json:"startTime"`
	EndTime      *time.Time     `json:"endTime"`
	EstimateTime *time.Duration `json:"estimateTime"`
}

// TaskResponse DTO for task responses
type TaskResponse struct {
	ID           string         `json:"id"`
	EventID      string         `json:"eventId"`
	MainTaskID   *string        `json:"mainTaskId"`
	Description  string         `json:"description"`
	StartTime    *time.Time     `json:"startTime"`
	EndTime      *time.Time     `json:"endTime"`
	EstimateTime *time.Duration `json:"estimateTime"`
}

// ToEntity converts TaskCreateRequest to entity.Task
func (req *TaskCreateRequest) ToEntity() *entity.Task {
	return &entity.Task{
		EventID:      req.EventID,
		MainTaskID:   req.MainTaskID,
		Description:  req.Description,
		StartTime:    req.StartTime,
		EndTime:      req.EndTime,
		EstimateTime: req.EstimateTime,
	}
}

// ToEntity converts TaskUpdateRequest to entity.Task
func (req *TaskUpdateRequest) ToEntity(id string) *entity.Task {
	return &entity.Task{
		ID:           id,
		EventID:      req.EventID,
		MainTaskID:   req.MainTaskID,
		Description:  req.Description,
		StartTime:    req.StartTime,
		EndTime:      req.EndTime,
		EstimateTime: req.EstimateTime,
	}
}

// FromTaskEntity converts entity.Task to TaskResponse
func FromTaskEntity(task *entity.Task) *TaskResponse {
	return &TaskResponse{
		ID:           task.ID,
		EventID:      task.EventID,
		MainTaskID:   task.MainTaskID,
		Description:  task.Description,
		StartTime:    task.StartTime,
		EndTime:      task.EndTime,
		EstimateTime: task.EstimateTime,
	}
}

// FromTaskEntities converts a slice of entity.Task to a slice of TaskResponse
func FromTaskEntities(tasks []*entity.Task) []*TaskResponse {
	responses := make([]*TaskResponse, len(tasks))
	for i, task := range tasks {
		responses[i] = FromTaskEntity(task)
	}
	return responses
}