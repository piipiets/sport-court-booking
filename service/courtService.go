package service

import (
	"fmt"
	"io"
	"path/filepath"

	"github.com/piipiets/sport-court-booking/helpers/constant"
	"github.com/piipiets/sport-court-booking/model/dto/request"
	"github.com/piipiets/sport-court-booking/model/dto/response"
	"github.com/piipiets/sport-court-booking/model/entity"
	"github.com/piipiets/sport-court-booking/repository"
	"github.com/piipiets/sport-court-booking/storage"
)

type CourtService interface {
	Create(req request.CourtRequest, image *UploadImage) error
	GetAll() ([]response.CourtResponse, error)
	GetByID(id int64) (*response.CourtResponse, error)
	GetImage(id int64) ([]byte, string, error)
	Update(id int64, req request.UpdateCourtRequest, image *UploadImage) error
	Delete(id int64) error
}

type UploadImage struct {
	File     io.Reader
	FileName string
	Size     int64
}

type courtService struct {
	courtRepo repository.CourtRepository
	storage   storage.Storage
}

func NewCourtService(courtRepo repository.CourtRepository, fileStorage storage.Storage) CourtService {
	return &courtService{courtRepo: courtRepo, storage: fileStorage}
}

func (s *courtService) Create(req request.CourtRequest, image *UploadImage) error {
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

	if image != nil {
		_, err = s.uploadCourtImage(court.ID, image)
		if err != nil {
			return err
		}
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

func (s *courtService) GetImage(id int64) ([]byte, string, error) {
	court, err := s.courtRepo.FindByID(id)
	if err != nil {
		return nil, "", err
	}

	if court.ImageURL == nil || *court.ImageURL == "" {
		return nil, "", constant.ErrCourtImageNotFound
	}

	data, contentType, err := s.storage.Download(*court.ImageURL)
	if err != nil {
		return nil, "", err
	}

	return data, contentType, nil
}

func (s *courtService) Update(id int64, req request.UpdateCourtRequest, image *UploadImage) error {
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

	var oldImageURL *string
	var newImageURL *string
	if image != nil {
		oldImageURL = existing.ImageURL
		newImageURL, err = s.uploadCourtImage(id, image)
		if err != nil {
			return err
		}
		existing.ImageURL = newImageURL
	}

	err = s.courtRepo.Update(existing)
	if err != nil {
		return err
	}

	if oldImageURL != nil && newImageURL != nil && *oldImageURL != *newImageURL {
		_ = s.storage.Delete(*oldImageURL)
	}

	return nil
}

func (s *courtService) Delete(id int64) error {
	court, err := s.courtRepo.FindByID(id)
	if err != nil {
		return err
	}

	err = s.courtRepo.Delete(id)
	if err != nil {
		return err
	}

	if court.ImageURL != nil {
		_ = s.storage.Delete(*court.ImageURL)
	}

	return nil
}

func (s *courtService) uploadCourtImage(courtID int64, image *UploadImage) (*string, error) {
	fileName := filepath.Base(image.FileName)
	relativePath := fmt.Sprintf("courts/%d/%s", courtID, fileName)
	publicURL, err := s.storage.Upload(relativePath, image.File, storage.ContentType(relativePath))
	if err != nil {
		return nil, err
	}

	err = s.courtRepo.UpdateImageURL(courtID, &publicURL)
	if err != nil {
		return nil, err
	}

	return &publicURL, nil
}

func toCourtResponse(c entity.Courts) response.CourtResponse {
	return response.CourtResponse{
		ID:       c.ID,
		Name:     c.Name,
		Type:     c.Type,
		Price:    c.Price,
		Location: c.Location,
		ImageURL: c.ImageURL,
	}
}
