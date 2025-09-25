package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"brb/internal/dto"
	"brb/internal/entity"
	"brb/internal/router"
)

// signHandler 处理sign相关的HTTP请求
type signHandler struct {
	signService signService
}

type signService interface {
	CreateSign(sign *entity.Sign) error
	GetSignByID(id int64) (*entity.Sign, error)
	UpdateSign(sign *entity.Sign) error
	DeleteSign(id int64) error
}

// NewSignHandler 创建新的SignHandler
func NewSignHandler(signService signService) *signHandler {
	return &signHandler{signService: signService}
}

// CreateSign 创建新sign
func (h *signHandler) CreateSign(w http.ResponseWriter, r *http.Request) {
	var req dto.SignCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	sign := req.ToEntity()
	if err := h.signService.CreateSign(sign); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := dto.FromSignEntity(sign)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// GetSign 获取单个sign
func (h *signHandler) GetSign(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	if idStr == "" {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	sign, err := h.signService.GetSignByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	response := dto.FromSignEntity(sign)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// UpdateSign 更新sign
func (h *signHandler) UpdateSign(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	if idStr == "" {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var req dto.SignUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	sign := req.ToEntity(id)
	if err := h.signService.UpdateSign(sign); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// DeleteSign 删除sign
func (h *signHandler) DeleteSign(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	if idStr == "" {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	if err := h.signService.DeleteSign(id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// RegisterRoutes 注册sign相关路由（新接口）
func (h *signHandler) RegisterRoutes(r router.Router) {
    // 为所有sign路由添加统一中间件
    api := r.Group("/api/signs")
    
    api.POST("", h.CreateSign)
    api.GET("/{id}", h.GetSign)
    api.PUT("/{id}", h.UpdateSign)
    api.DELETE("/{id}", h.DeleteSign)
}

