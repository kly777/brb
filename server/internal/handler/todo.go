package handler

import (
	"encoding/json"
	"net/http"

	"brb/internal/entity"
)

// todoHandler 处理todo相关的HTTP请求
type todoHandler struct {
	todoService TodoService
}

type TodoService interface {
	CreateTodo(todo *entity.Todo) error
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
	var todo entity.Todo
	if err := json.NewDecoder(r.Body).Decode(&todo); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.todoService.CreateTodo(&todo); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(todo)
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

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(todo)
}

// UpdateTodo 更新todo
func (h *todoHandler) UpdateTodo(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var todo entity.Todo
	if err := json.NewDecoder(r.Body).Decode(&todo); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	todo.ID = id

	if err := h.todoService.UpdateTodo(&todo); err != nil {
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
	mux.HandleFunc("GET /api/todos/{id}", h.GetTodo)
	mux.HandleFunc("PUT /api/todos/{id}", h.UpdateTodo)
	mux.HandleFunc("DELETE /api/todos/{id}", h.DeleteTodo)
}