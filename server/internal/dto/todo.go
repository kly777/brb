package dto

import (
	"brb/internal/entity"
	"time"
)

// TodoCreateRequest DTO for creating a todo
type TodoCreateRequest struct {
	EventID       *uint  `json:"eventId" form:"eventId"`
	TaskID        uint   `json:"taskId" form:"taskId"`
	Status        string `json:"status" form:"status"`
	PlannedStart  string `json:"plannedStart" form:"plannedStart"`
	PlannedEnd    string `json:"plannedEnd" form:"plannedEnd"`
	ActualStart   string `json:"actualStart" form:"actualStart"`
	ActualEnd     string `json:"actualEnd" form:"actualEnd"`
}

// TodoUpdateRequest DTO for updating a todo
type TodoUpdateRequest struct {
	EventID       *uint  `json:"eventId"`
	TaskID        uint   `json:"taskId"`
	Status        string `json:"status"`
	PlannedStart  string `json:"plannedStart"`
	PlannedEnd    string `json:"plannedEnd"`
	ActualStart   string `json:"actualStart"`
	ActualEnd     string `json:"actualEnd"`
}

// TodoResponse DTO for todo responses
type TodoResponse struct {
	ID           uint    `json:"id"`
	EventID      *uint   `json:"eventId"`
	TaskID       uint    `json:"taskId"`
	Status       string  `json:"status"`
	PlannedTime  TimeSpan `json:"plannedTime"`
	ActualTime   TimeSpan `json:"actualTime"`
	CompletedTime *string `json:"completedTime"`
}

// ToEntity converts TodoCreateRequest to entity.Todo
func (req *TodoCreateRequest) ToEntity() *entity.Todo {
	todo := &entity.Todo{
		EventID: req.EventID,
		TaskID:  req.TaskID,
		Status:  entity.Status(req.Status),
	}

	// Parse planned time
	if req.PlannedStart != "" && req.PlannedEnd != "" {
		todo.PlannedTime = parseTodoTimeSpan(req.PlannedStart, req.PlannedEnd)
	}

	// Parse actual time
	if req.ActualStart != "" && req.ActualEnd != "" {
		todo.ActualTime = parseTodoTimeSpan(req.ActualStart, req.ActualEnd)
	}

	return todo
}

// ToEntity converts TodoUpdateRequest to entity.Todo
func (req *TodoUpdateRequest) ToEntity(id uint) *entity.Todo {
	todo := &entity.Todo{
		ID:      id,
		EventID: req.EventID,
		TaskID:  req.TaskID,
		Status:  entity.Status(req.Status),
	}

	// Parse planned time
	if req.PlannedStart != "" && req.PlannedEnd != "" {
		todo.PlannedTime = parseTodoTimeSpan(req.PlannedStart, req.PlannedEnd)
	}

	// Parse actual time
	if req.ActualStart != "" && req.ActualEnd != "" {
		todo.ActualTime = parseTodoTimeSpan(req.ActualStart, req.ActualEnd)
	}

	return todo
}

// FromTodoEntity converts entity.Todo to TodoResponse
func FromTodoEntity(todo *entity.Todo) *TodoResponse {
	response := &TodoResponse{
		ID:      todo.ID,
		EventID: todo.EventID,
		TaskID:  todo.TaskID,
		Status:  string(todo.Status),
	}

	// Convert planned time
	if todo.PlannedTime.Start != nil {
		startStr := todo.PlannedTime.Start.Format("2006-01-02T15:04")
		response.PlannedTime.Start = &startStr
	}
	if todo.PlannedTime.End != nil {
		endStr := todo.PlannedTime.End.Format("2006-01-02T15:04")
		response.PlannedTime.End = &endStr
	}

	// Convert actual time
	if todo.ActualTime.Start != nil {
		startStr := todo.ActualTime.Start.Format("2006-01-02T15:04")
		response.ActualTime.Start = &startStr
	}
	if todo.ActualTime.End != nil {
		endStr := todo.ActualTime.End.Format("2006-01-02T15:04")
		response.ActualTime.End = &endStr
	}

	// Convert completed time
	if todo.CompletedTime != nil {
		completedStr := todo.CompletedTime.Format("2006-01-02T15:04")
		response.CompletedTime = &completedStr
	}

	return response
}

// FromTodoEntities converts a slice of entity.Todo to a slice of TodoResponse
func FromTodoEntities(todos []*entity.Todo) []*TodoResponse {
	responses := make([]*TodoResponse, len(todos))
	for i, todo := range todos {
		responses[i] = FromTodoEntity(todo)
	}
	return responses
}

// parseTodoTimeSpan parses start and end strings into a TimeSpan for todos
func parseTodoTimeSpan(startStr, endStr string) entity.TimeSpan {
	var startTime, endTime *time.Time
	
	if startStr != "" {
		if t, err := time.Parse("2006-01-02T15:04", startStr); err == nil {
			startTime = &t
		}
	}
	
	if endStr != "" {
		if t, err := time.Parse("2006-01-02T15:04", endStr); err == nil {
			endTime = &t
		}
	}
	
	return entity.TimeSpan{
		Start: startTime,
		End:   endTime,
	}
}