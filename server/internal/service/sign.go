package service

import (
	"brb/internal/entity"
)

// SignService 实现handler.SignService接口
type SignService struct {
	signRepo SignRepository
}


type SignRepository interface {
	Create(sign *entity.Sign) error
	GetByID(id int64) (*entity.Sign, error)
	Update(sign *entity.Sign) error
	Delete(id int64) error
}
// NewSignService 创建新的SignService实例
func NewSignService(signRepo SignRepository) *SignService {
	return &SignService{
		signRepo: signRepo,
	}
}

// CreateSign creates a new sign
func (s *SignService) CreateSign(sign *entity.Sign) error {

	return s.signRepo.Create(sign)
}

// GetSignByID 根据ID获取sign
func (s *SignService) GetSignByID(id int64) (*entity.Sign, error) {
	return s.signRepo.GetByID(id)
}


func (s *SignService) UpdateSign(sign *entity.Sign) error {

	return s.signRepo.Update(sign)
}


func (s *SignService) DeleteSign(id int64) error {
	return s.signRepo.Delete(id)
}