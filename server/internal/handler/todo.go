package handler

import (
	"encoding/json"
	"net/http"

	"brb/internal/dto"
	"brb/internal/entity"
)

// todoHandler 处理todo相关的HTTP请求
type todoHandler struct {
	todoService TodoService
}

type TodoService interface {
	CreateTodo(todo *entity.Todo) error
	GetAllTodo() ([]*entity.Todo, error)
	GetTodoByID(id string) (*entity.Todo, error)
	UpdateTodo(todo *entity.Todo) error
	DeleteTodo(id string) error
}

// NewTodoHandler 创建新的TodoHandler
func NewTodoHandler(todoService TodoService) *todoHandler {
	return &todoHandler{todoService: todoService}
}

// CreateTodo 创建新todo
func (h *todoHandler) CreateTodo(w http.ResponseWriter, r *http.Request) {
	var req dto.TodoCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	todo := req.ToEntity()
	if err := h.todoService.CreateTodo(todo); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := dto.FromTodoEntity(todo)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// GetAllTodo 获取所有todo
func (h *todoHandler) GetAllTodo(w http.ResponseWriter, r *http.Request) {
	todos, err := h.todoService.GetAllTodo()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	responses := dto.FromTodoEntities(todos)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(responses)
}

// GetTodo 获取单个todo
func (h *todoHandler) GetTodo(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	todo, err := h.todoService.GetTodoByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	response := dto.FromTodoEntity(todo)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// UpdateTodo 更新todo
func (h *todoHandler) UpdateTodo(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var req dto.TodoUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	todo := req.ToEntity(id)
	if err := h.todoService.UpdateTodo(todo); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// DeleteTodo 删除todo
func (h *todoHandler) DeleteTodo(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	if err := h.todoService.DeleteTodo(id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// RegisterRoutes 注册todo相关路由
func (h *todoHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/todos", h.CreateTodo)
	mux.HandleFunc("GET /api/todos", h.GetAllTodo)
	mux.HandleFunc("GET /api/todos/{id}", h.GetTodo)
	mux.HandleFunc("PUT /api/todos/{id}", h.UpdateTodo)
	mux.HandleFunc("DELETE /api/todos/{id}", h.DeleteTodo)
}
