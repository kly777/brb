package dto

import (
	"brb/internal/entity"
	"time"
)

// TodoCreateRequest DTO for creating a todo
type TodoCreateRequest struct {
	TaskID        string     `json:"taskId"`
	Status        string     `json:"status"`
	Priority      int        `json:"priority"`
	StartTime     *time.Time `json:"startTime"`
	EndTime       *time.Time `json:"endTime"`
	CompletedTime *time.Time `json:"completedTime"`
}

// TodoUpdateRequest DTO for updating a todo
type TodoUpdateRequest struct {
	TaskID        string     `json:"taskId"`
	Status        string     `json:"status"`
	Priority      int        `json:"priority"`
	StartTime     *time.Time `json:"startTime"`
	EndTime       *time.Time `json:"endTime"`
	CompletedTime *time.Time `json:"completedTime"`
}

// TodoResponse DTO for todo responses
type TodoResponse struct {
	ID            string     `json:"id"`
	TaskID        string     `json:"taskId"`
	Status        string     `json:"status"`
	Priority      int        `json:"priority"`
	StartTime     *time.Time `json:"startTime"`
	EndTime       *time.Time `json:"endTime"`
	CompletedTime *time.Time `json:"completedTime"`
}

// ToEntity converts TodoCreateRequest to entity.Todo
func (req *TodoCreateRequest) ToEntity() *entity.Todo {
	return &entity.Todo{
		TaskID:        req.TaskID,
		Status:        entity.Status(req.Status),
		Priority:      req.Priority,
		StartTime:     req.StartTime,
		EndTime:       req.EndTime,
		CompletedTime: req.CompletedTime,
	}
}

// ToEntity converts TodoUpdateRequest to entity.Todo
func (req *TodoUpdateRequest) ToEntity(id string) *entity.Todo {
	return &entity.Todo{
		ID:            id,
		TaskID:        req.TaskID,
		Status:        entity.Status(req.Status),
		Priority:      req.Priority,
		StartTime:     req.StartTime,
		EndTime:       req.EndTime,
		CompletedTime: req.CompletedTime,
	}
}

// FromTodoEntity converts entity.Todo to TodoResponse
func FromTodoEntity(todo *entity.Todo) *TodoResponse {
	return &TodoResponse{
		ID:            todo.ID,
		TaskID:        todo.TaskID,
		Status:        string(todo.Status),
		Priority:      todo.Priority,
		StartTime:     todo.StartTime,
		EndTime:       todo.EndTime,
		CompletedTime: todo.CompletedTime,
	}
}

// FromTodoEntities converts a slice of entity.Todo to a slice of TodoResponse
func FromTodoEntities(todos []*entity.Todo) []*TodoResponse {
	responses := make([]*TodoResponse, len(todos))
	for i, todo := range todos {
		responses[i] = FromTodoEntity(todo)
	}
	return responses
}