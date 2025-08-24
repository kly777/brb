package handler

import (
	"encoding/json"
	"net/http"

	"brb/internal/dto"
	"brb/internal/entity"
)

// eventHandler 处理event相关的HTTP请求
type eventHandler struct {
	eventService eventService
}

type eventService interface {
	CreateEvent(event *entity.Event) error
	GetEventByID(id string) (*entity.Event, error)
	UpdateEvent(event *entity.Event) error
	DeleteEvent(id string) error
}

// NewEventHandler 创建新的EventHandler
func NewEventHandler(eventService eventService) *eventHandler {
	return &eventHandler{eventService: eventService}
}

// CreateEvent 创建新event
func (h *eventHandler) CreateEvent(w http.ResponseWriter, r *http.Request) {
	var req dto.EventCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	event := req.ToEntity()
	if err := h.eventService.CreateEvent(event); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := dto.FromEventEntity(event)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
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

	response := dto.FromEventEntity(event)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// UpdateEvent 更新event
func (h *eventHandler) UpdateEvent(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var req dto.EventUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	event := req.ToEntity(id)
	if err := h.eventService.UpdateEvent(event); err != nil {
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
