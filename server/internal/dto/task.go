package dto

import (
	"brb/internal/entity"
	"time"
)

// TaskCreateRequest DTO for creating a task
type TaskCreateRequest struct {
	EventID        uint    `json:"eventId" form:"eventId"`
	ParentTaskID   *uint   `json:"parentTaskId" form:"parentTaskId"`
	PreTaskIDs     []uint  `json:"preTaskIds" form:"preTaskIds"`
	Description    string  `json:"description" form:"description"`
	AllowedStart   *string `json:"allowedStart" form:"allowedStart"`
	AllowedEnd     *string `json:"allowedEnd" form:"allowedEnd"`
	PlannedStart   *string `json:"plannedStart" form:"plannedStart"`
	PlannedEnd     *string `json:"plannedEnd" form:"plannedEnd"`
	Status         string  `json:"status" form:"status"`
}

// TaskUpdateRequest DTO for updating a task
type TaskUpdateRequest struct {
	EventID        uint    `json:"eventId"`
	ParentTaskID   *uint   `json:"parentTaskId"`
	PreTaskIDs     []uint  `json:"preTaskIds"`
	Description    string  `json:"description"`
	AllowedStart   *string `json:"allowedStart"`
	AllowedEnd     *string `json:"allowedEnd"`
	PlannedStart   *string `json:"plannedStart"`
	PlannedEnd     *string `json:"plannedEnd"`
	Status         string  `json:"status"`
}

// TaskResponse DTO for task responses
type TaskResponse struct {
	ID             uint     `json:"id"`
	EventID        uint     `json:"eventId"`
	ParentTaskID   *uint    `json:"parentTaskId"`
	PreTaskIDs     []uint   `json:"preTaskIds"`
	Description    string   `json:"description"`
	AllowedTime    TimeSpan `json:"allowedTime"`
	PlannedTime    TimeSpan `json:"plannedTime"`
	Status         string   `json:"status"`
	CreatedAt      string   `json:"createdAt"`
}

// TimeSpan represents a time range with start and end
type TimeSpan struct {
	Start *string `json:"start"`
	End   *string `json:"end"`
}

// ToEntity converts TaskCreateRequest to entity.Task
func (req *TaskCreateRequest) ToEntity() *entity.Task {
	task := &entity.Task{
		EventID:      req.EventID,
		ParentTaskID: req.ParentTaskID,
		PreTaskIDs:   req.PreTaskIDs,
		Description:  req.Description,
		Status:       entity.Status(req.Status),
	}

	// Parse time strings into TimeSpan
	if req.AllowedStart != nil && req.AllowedEnd != nil {
		task.AllowedTime = parseTimeSpan(*req.AllowedStart, *req.AllowedEnd)
	}

	if req.PlannedStart != nil && req.PlannedEnd != nil {
		task.PlannedDuration = parseTimeSpan(*req.PlannedStart, *req.PlannedEnd)
	}

	return task
}

// ToEntity converts TaskUpdateRequest to entity.Task
func (req *TaskUpdateRequest) ToEntity(id uint) *entity.Task {
	task := &entity.Task{
		ID:           id,
		EventID:      req.EventID,
		ParentTaskID: req.ParentTaskID,
		PreTaskIDs:   req.PreTaskIDs,
		Description:  req.Description,
		Status:       entity.Status(req.Status),
	}

	// Parse time strings into TimeSpan
	if req.AllowedStart != nil && req.AllowedEnd != nil {
		task.AllowedTime = parseTimeSpan(*req.AllowedStart, *req.AllowedEnd)
	}

	if req.PlannedStart != nil && req.PlannedEnd != nil {
		task.PlannedDuration = parseTimeSpan(*req.PlannedStart, *req.PlannedEnd)
	}

	return task
}

// FromTaskEntity converts entity.Task to TaskResponse
func FromTaskEntity(task *entity.Task) *TaskResponse {
	response := &TaskResponse{
		ID:           task.ID,
		EventID:      task.EventID,
		ParentTaskID: task.ParentTaskID,
		PreTaskIDs:   task.PreTaskIDs,
		Description:  task.Description,
		Status:       string(task.Status),
		CreatedAt:    task.CreatedAt.Format(time.RFC3339),
	}

	// Convert TimeSpan to string representations
	if task.AllowedTime.Start != nil {
		startStr := task.AllowedTime.Start.Format(time.RFC3339)
		response.AllowedTime.Start = &startStr
	}
	if task.AllowedTime.End != nil {
		endStr := task.AllowedTime.End.Format(time.RFC3339)
		response.AllowedTime.End = &endStr
	}

	if task.PlannedDuration.Start != nil {
		startStr := task.PlannedDuration.Start.Format(time.RFC3339)
		response.PlannedTime.Start = &startStr
	}
	if task.PlannedDuration.End != nil {
		endStr := task.PlannedDuration.End.Format(time.RFC3339)
		response.PlannedTime.End = &endStr
	}

	return response
}

// FromTaskEntities converts a slice of entity.Task to a slice of TaskResponse
func FromTaskEntities(tasks []*entity.Task) []*TaskResponse {
	responses := make([]*TaskResponse, len(tasks))
	for i, task := range tasks {
		responses[i] = FromTaskEntity(task)
	}
	return responses
}

// parseTimeSpan parses start and end strings into a TimeSpan
func parseTimeSpan(startStr, endStr string) entity.TimeSpan {
	var startTime, endTime *time.Time
	
	if startStr != "" {
		if t, err := time.Parse(time.RFC3339, startStr); err == nil {
			startTime = &t
		}
	}
	
	if endStr != "" {
		if t, err := time.Parse(time.RFC3339, endStr); err == nil {
			endTime = &t
		}
	}
	
	return entity.TimeSpan{
		Start: startTime,
		End:   endTime,
	}
}