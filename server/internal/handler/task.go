package handler

import (
	"encoding/json"
	"net/http"

	"brb/internal/dto"
	"brb/internal/entity"
)

// taskHandler 处理task相关的HTTP请求
type taskHandler struct {
	taskService TaskService
}

type TaskService interface {
	CreateTask(task *entity.Task) error
	GetTaskByID(id string) (*entity.Task, error)
	UpdateTask(task *entity.Task) error
	DeleteTask(id string) error
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

// GetTask 获取单个task
func (h *taskHandler) GetTask(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
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
	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
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
	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	if err := h.taskService.DeleteTask(id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// RegisterRoutes 注册task相关路由
func (h *taskHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/tasks", h.CreateTask)
	mux.HandleFunc("GET /api/tasks/{id}", h.GetTask)
	mux.HandleFunc("PUT /api/tasks/{id}", h.UpdateTask)
	mux.HandleFunc("DELETE /api/tasks/{id}", h.DeleteTask)
}
