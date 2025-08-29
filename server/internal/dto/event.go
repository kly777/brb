package dto

import "brb/internal/entity"

// EventCreateRequest DTO for creating an event
type EventCreateRequest struct {
	IsTemplate  bool   `json:"isTemplate" form:"isTemplate"`
	Title       string `json:"title" form:"title"`
	Description string `json:"description" form:"description"`
	Location    string `json:"location" form:"location"`
	Priority    int    `json:"priority" form:"priority"`
	Category    string `json:"category" form:"category"`
}

// EventUpdateRequest DTO for updating an event
type EventUpdateRequest struct {
	IsTemplate  bool   `json:"isTemplate"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Location    string `json:"location"`
	Priority    int    `json:"priority"`
	Category    string `json:"category"`
}

// EventResponse DTO for event responses
type EventResponse struct {
	ID          uint   `json:"id"`
	IsTemplate  bool   `json:"isTemplate"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Location    string `json:"location"`
	Priority    int    `json:"priority"`
	Category    string `json:"category"`
}

// ToEntity converts EventCreateRequest to entity.Event
func (req *EventCreateRequest) ToEntity() *entity.Event {
	return &entity.Event{
		IsTemplate:  req.IsTemplate,
		Title:       req.Title,
		Description: req.Description,
		Location:    req.Location,
		Priority:    req.Priority,
		Category:    req.Category,
	}
}

// ToEntity converts EventUpdateRequest to entity.Event
func (req *EventUpdateRequest) ToEntity(id uint) *entity.Event {
	return &entity.Event{
		ID:          id,
		IsTemplate:  req.IsTemplate,
		Title:       req.Title,
		Description: req.Description,
		Location:    req.Location,
		Priority:    req.Priority,
		Category:    req.Category,
	}
}

// FromEventEntity converts entity.Event to EventResponse
func FromEventEntity(event *entity.Event) *EventResponse {
	return &EventResponse{
		ID:          event.ID,
		IsTemplate:  event.IsTemplate,
		Title:       event.Title,
		Description: event.Description,
		Location:    event.Location,
		Priority:    event.Priority,
		Category:    event.Category,
	}
}

// FromEventEntities converts a slice of entity.Event to a slice of EventResponse
func FromEventEntities(events []*entity.Event) []*EventResponse {
	responses := make([]*EventResponse, len(events))
	for i, event := range events {
		responses[i] = FromEventEntity(event)
	}
	return responses
}