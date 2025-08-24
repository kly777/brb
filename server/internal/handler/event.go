package handler

import (
	"encoding/json"
	"net/http"

	"brb/internal/entity"
)

// eventHandler 处理event相关的HTTP请求
type eventHandler struct {
	eventService EventService
}

type EventService interface {
	CreateEvent(event *entity.Event) error
	GetEventByID(id string) (*entity.Event, error)
	UpdateEvent(event *entity.Event) error
	DeleteEvent(id string) error
}

// NewEventHandler 创建新的EventHandler
func NewEventHandler(eventService EventService) *eventHandler {
	return &eventHandler{eventService: eventService}
}

// CreateEvent 创建新event
func (h *eventHandler) CreateEvent(w http.ResponseWriter, r *http.Request) {
	var event entity.Event
	if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.eventService.CreateEvent(&event); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(event)
}

// GetEvent 获取单个event
func (h *eventHandler) GetEvent(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	event, err := h.eventService.GetEventByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(event)
}

// UpdateEvent 更新event
func (h *eventHandler) UpdateEvent(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var event entity.Event
	if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	event.ID = id

	if err := h.eventService.UpdateEvent(&event); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// DeleteEvent 删除event
func (h *eventHandler) DeleteEvent(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	if err := h.eventService.DeleteEvent(id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// RegisterRoutes 注册event相关路由
func (h *eventHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/events", h.CreateEvent)
	mux.HandleFunc("GET /api/events/{id}", h.GetEvent)
	mux.HandleFunc("PUT /api/events/{id}", h.UpdateEvent)
	mux.HandleFunc("DELETE /api/events/{id}", h.DeleteEvent)
}
