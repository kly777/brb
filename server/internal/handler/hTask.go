package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"brb/internal/dto"
	"brb/internal/entity"
	"brb/internal/router"
)

// taskHandler 处理task相关的HTTP请求
type taskHandler struct {
	taskService TaskService
}

type TaskService interface {
	CreateTask(task *entity.Task) error
	GetAllTasks() ([]*entity.Task, error)
	GetTaskByID(id uint) (*entity.Task, error)
	UpdateTask(task *entity.Task) error
	DeleteTask(id uint) error
}

// NewTaskHandler 创建新的TaskHandler
func NewTaskHandler(taskService TaskService) *taskHandler {
	return &taskHandler{taskService: taskService}
}

// CreateTask 创建新task
func (h *taskHandler) CreateTask(w http.ResponseWriter, r *http.Request) {
	var req dto.TaskCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	task := req.ToEntity()
	if err := h.taskService.CreateTask(task); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := dto.FromTaskEntity(task)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// GetAllTasks 获取所有task
func (h *taskHandler) GetAllTasks(w http.ResponseWriter, r *http.Request) {
	tasks, err := h.taskService.GetAllTasks()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	response := dto.FromTaskEntities(tasks)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetTask 获取单个task
func (h *taskHandler) GetTask(w http.ResponseWriter, r *http.Request) {
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

	task, err := h.taskService.GetTaskByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	response := dto.FromTaskEntity(task)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// UpdateTask 更新task
func (h *taskHandler) UpdateTask(w http.ResponseWriter, r *http.Request) {
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

	var req dto.TaskUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	task := req.ToEntity(id)
	if err := h.taskService.UpdateTask(task); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// DeleteTask 删除task
func (h *taskHandler) DeleteTask(w http.ResponseWriter, r *http.Request) {
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

	if err := h.taskService.DeleteTask(id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// RegisterRoutes 注册task相关路由（新接口）
func (h *taskHandler) RegisterRoutes(r router.Router) {
    // 为所有task路由添加统一中间件
    api := r.Group("/api/tasks")
    
    api.POST("", h.CreateTask)
    api.GET("", h.GetAllTasks)
    api.GET("/{id}", h.GetTask)
    api.PUT("/{id}", h.UpdateTask)
    api.DELETE("/{id}", h.DeleteTask)
}
