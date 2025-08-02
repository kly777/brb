package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"brb/internal/entity"

	"github.com/gorilla/mux"
)

// SignHandler 处理sign相关的HTTP请求
type SignHandler struct {
	signService SignService
}

type SignService interface {
	CreateSign(sign *entity.Sign) error
	GetSignByID(id int64) (*entity.Sign, error)
	UpdateSign(sign *entity.Sign) error
	DeleteSign(id int64) error
}

// NewSignHandler 创建新的SignHandler
func NewSignHandler(signService SignService) *SignHandler {
	return &SignHandler{signService: signService}
}

// CreateSign 创建新sign
func (h *SignHandler) CreateSign(w http.ResponseWriter, r *http.Request) {
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
func (h *SignHandler) GetSign(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
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
func (h *SignHandler) UpdateSign(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
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
func (h *SignHandler) DeleteSign(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
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
func (h *SignHandler) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/signs", h.CreateSign).Methods("POST")
	r.HandleFunc("/signs/{id}", h.GetSign).Methods("GET")
	r.HandleFunc("/signs/{id}", h.UpdateSign).Methods("PUT")
	r.HandleFunc("/signs/{id}", h.DeleteSign).Methods("DELETE")
}
