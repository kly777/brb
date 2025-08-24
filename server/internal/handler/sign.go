package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"brb/internal/entity"
)

// signHandler 处理sign相关的HTTP请求
type signHandler struct {
	signService SignService
}

type SignService interface {
	CreateSign(sign *entity.Sign) error
	GetSignByID(id int64) (*entity.Sign, error)
	UpdateSign(sign *entity.Sign) error
	DeleteSign(id int64) error
}

// NewSignHandler 创建新的SignHandler
func NewSignHandler(signService SignService) *signHandler {
	return &signHandler{signService: signService}
}

// CreateSign 创建新sign
func (h *signHandler) CreateSign(w http.ResponseWriter, r *http.Request) {
	var sign entity.Sign
	if err := json.NewDecoder(r.Body).Decode(&sign); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.signService.CreateSign(&sign); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(sign)
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

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(sign)
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

	var sign entity.Sign
	if err := json.NewDecoder(r.Body).Decode(&sign); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	sign.ID = id

	if err := h.signService.UpdateSign(&sign); err != nil {
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

// RegisterRoutes 注册sign相关路由
func (h *signHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/signs", h.CreateSign)
	mux.HandleFunc("GET /api/signs/{id}", h.GetSign)
	mux.HandleFunc("PUT /api/signs/{id}", h.UpdateSign)
	mux.HandleFunc("DELETE /api/signs/{id}", h.DeleteSign)
}

