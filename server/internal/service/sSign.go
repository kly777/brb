package service

import (
	"brb/internal/entity"
)

// signService 实现handler.signService接口
type signService struct {
	signRepo signRepository
}


type signRepository interface {
	Create(sign *entity.Sign) error
	GetByID(id int64) (*entity.Sign, error)
	Update(sign *entity.Sign) error
	Delete(id int64) error
}

// NewSignService 创建新的SignService实例
func NewSignService(signRepo signRepository) *signService {
	return &signService{
		signRepo: signRepo,
	}
}

// CreateSign creates a new sign
func (s *signService) CreateSign(sign *entity.Sign) error {

	return s.signRepo.Create(sign)
}

// GetSignByID 根据ID获取sign
func (s *signService) GetSignByID(id int64) (*entity.Sign, error) {
	return s.signRepo.GetByID(id)
}


func (s *signService) UpdateSign(sign *entity.Sign) error {

	return s.signRepo.Update(sign)
}


func (s *signService) DeleteSign(id int64) error {
	return s.signRepo.Delete(id)
}