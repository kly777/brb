package dto

import "brb/internal/entity"

// EventCreateRequest DTO for creating an event
type EventCreateRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Recurrence  string `json:"recurrence"`
}

// EventUpdateRequest DTO for updating an event
type EventUpdateRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Recurrence  string `json:"recurrence"`
}

// EventResponse DTO for event responses
type EventResponse struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Recurrence  string `json:"recurrence"`
}

// ToEntity converts EventCreateRequest to entity.Event
func (req *EventCreateRequest) ToEntity() *entity.Event {
	return &entity.Event{
		Title:       req.Title,
		Description: req.Description,
		Recurrence:  req.Recurrence,
	}
}

// ToEntity converts EventUpdateRequest to entity.Event
func (req *EventUpdateRequest) ToEntity(id string) *entity.Event {
	return &entity.Event{
		ID:          id,
		Title:       req.Title,
		Description: req.Description,
		Recurrence:  req.Recurrence,
	}
}

// FromEventEntity converts entity.Event to EventResponse
func FromEventEntity(event *entity.Event) *EventResponse {
	return &EventResponse{
		ID:          event.ID,
		Title:       event.Title,
		Description: event.Description,
		Recurrence:  event.Recurrence,
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