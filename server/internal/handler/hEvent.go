package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"brb/internal/dto"
	"brb/internal/entity"
	"brb/internal/router"
	"brb/pkg/logger"
)

// eventHandler 处理event相关的HTTP请求
type eventHandler struct {
	eventService eventService
}

type eventService interface {
	CreateEvent(event *entity.Event) error
	GetAllEvents() ([]*entity.Event, error)
	GetEventByID(id uint) (*entity.Event, error)
	UpdateEvent(event *entity.Event) error
	DeleteEvent(id uint) error
}

// NewEventHandler 创建新的EventHandler
func NewEventHandler(eventService eventService) *eventHandler {
	return &eventHandler{eventService: eventService}
}

// CreateEvent 创建新event
func (h *eventHandler) CreateEvent(w http.ResponseWriter, r *http.Request) {
	var req dto.EventCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Tip.Println("Invalid request body:", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	event := req.ToEntity()
	if err := h.eventService.CreateEvent(event); err != nil {
		logger.Tip.Println("Failed to create event:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := dto.FromEventEntity(event)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// GetAllEvents 获取所有event
func (h *eventHandler) GetAllEvents(w http.ResponseWriter, r *http.Request) {
	events, err := h.eventService.GetAllEvents()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	response := dto.FromEventEntities(events)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetEvent 获取单个event
func (h *eventHandler) GetEvent(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	if idStr == "" {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var id uint
	if _, err := fmt.Sscanf(idStr, "%d", &id); err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
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
	idStr := r.PathValue("id")
	if idStr == "" {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var id uint
	if _, err := fmt.Sscanf(idStr, "%d", &id); err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
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
	idStr := r.PathValue("id")
	if idStr == "" {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var id uint
	if _, err := fmt.Sscanf(idStr, "%d", &id); err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	if err := h.eventService.DeleteEvent(id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *eventHandler) RegisterRoutes(r router.Router) {
	// 为所有event路由添加统一中间件
	api := r.Group("/api/events")

	api.GET("", h.GetAllEvents)
	api.GET("/{id}", h.GetEvent)
	api.POST("", h.CreateEvent)
	api.PUT("/{id}", h.UpdateEvent)
	api.DELETE("/{id}", h.DeleteEvent)
}
