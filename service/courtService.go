package service

import (
	"github.com/piipiets/sport-court-booking/model/dto/request"
	"github.com/piipiets/sport-court-booking/model/dto/response"
	"github.com/piipiets/sport-court-booking/model/entity"
	"github.com/piipiets/sport-court-booking/repository"
)

type CourtService interface {
	Create(req request.CourtRequest) error
	GetAll() ([]response.CourtResponse, error)
	GetByID(id int64) (*response.CourtResponse, error)
	Update(id int64, req request.UpdateCourtRequest) error
	Delete(id int64) error
}

type courtService struct {
	courtRepo repository.CourtRepository
}

func NewCourtService(courtRepo repository.CourtRepository) CourtService {
	return &courtService{courtRepo: courtRepo}
}

func (s *courtService) Create(req request.CourtRequest) error {
	court := &entity.Courts{
		Name:     req.Name,
		Type:     req.Type,
		Price:    req.Price,
		Location: req.Location,
	}

	err := s.courtRepo.Create(court)
	if err != nil {
		return err
	}

	return nil
}

func (s *courtService) GetAll() ([]response.CourtResponse, error) {
	courts, err := s.courtRepo.FindAll()
	if err != nil {
		return nil, err
	}

	responses := make([]response.CourtResponse, 0, len(courts))
	for _, c := range courts {
		responses = append(responses, toCourtResponse(c))
	}

	return responses, nil
}

func (s *courtService) GetByID(id int64) (*response.CourtResponse, error) {
	court, err := s.courtRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	resp := toCourtResponse(*court)
	return &resp, nil
}

func (s *courtService) Update(id int64, req request.UpdateCourtRequest) error {
	existing, err := s.courtRepo.FindByID(id)
	if err != nil {
		return err
	}

	if req.Name != nil {
		existing.Name = *req.Name
	}
	if req.Type != nil {
		existing.Type = *req.Type
	}
	if req.Price != nil {
		existing.Price = *req.Price
	}
	if req.Location != nil {
		existing.Location = *req.Location
	}

	err = s.courtRepo.Update(existing)
	if err != nil {
		return err
	}

	return nil
}

func (s *courtService) Delete(id int64) error {
	return s.courtRepo.Delete(id)
}

func toCourtResponse(c entity.Courts) response.CourtResponse {
	return response.CourtResponse{
		ID:       c.ID,
		Name:     c.Name,
		Type:     c.Type,
		Price:    c.Price,
		Location: c.Location,
	}
}
