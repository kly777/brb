package dto

import "brb/internal/entity"

// SignCreateRequest DTO for creating a sign
type SignCreateRequest struct {
	Signifier string `json:"signifier"`
	Signified string `json:"signified"`
}

// SignUpdateRequest DTO for updating a sign
type SignUpdateRequest struct {
	Signifier string `json:"signifier"`
	Signified string `json:"signified"`
}

// SignResponse DTO for sign responses
type SignResponse struct {
	ID        int64  `json:"id"`
	Signifier string `json:"signifier"`
	Signified string `json:"signified"`
}

// ToEntity converts SignCreateRequest to entity.Sign
func (req *SignCreateRequest) ToEntity() *entity.Sign {
	return &entity.Sign{
		Signifier: req.Signifier,
		Signified: req.Signified,
	}
}

// ToEntity converts SignUpdateRequest to entity.Sign
func (req *SignUpdateRequest) ToEntity(id int64) *entity.Sign {
	return &entity.Sign{
		ID:        id,
		Signifier: req.Signifier,
		Signified: req.Signified,
	}
}

// FromSignEntity converts entity.Sign to SignResponse
func FromSignEntity(sign *entity.Sign) *SignResponse {
	return &SignResponse{
		ID:        sign.ID,
		Signifier: sign.Signifier,
		Signified: sign.Signified,
	}
}

// FromSignEntities converts a slice of entity.Sign to a slice of SignResponse
func FromSignEntities(signs []*entity.Sign) []*SignResponse {
	responses := make([]*SignResponse, len(signs))
	for i, sign := range signs {
		responses[i] = FromSignEntity(sign)
	}
	return responses
}