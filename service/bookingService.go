package service

import (
	"errors"
	"time"

	"github.com/piipiets/sport-court-booking/helpers/constant"
	"github.com/piipiets/sport-court-booking/model/dto/request"
	"github.com/piipiets/sport-court-booking/model/dto/response"
	"github.com/piipiets/sport-court-booking/model/entity"
	"github.com/piipiets/sport-court-booking/repository"
)

const (
	dateFormat            = "2006-01-02"
	timeFormat            = "15:04"
	paymentDeadlineWindow = 15 * time.Minute
)

type BookingService interface {
	Create(userID int64, req request.CreateBookingRequest) error
	GetByUserID(userID int64) ([]response.BookingResponse, error)
	GetByID(id int64, requesterID int64, isAdmin bool) (*response.BookingResponse, error)
	UpdateStatus(id int64, req request.UpdateBookingStatusRequest) error
}

type bookingService struct {
	bookingRepo repository.BookingRepository
	courtRepo   repository.CourtRepository
}

func NewBookingService(bookingRepo repository.BookingRepository, courtRepo repository.CourtRepository) BookingService {
	return &bookingService{bookingRepo: bookingRepo, courtRepo: courtRepo}
}

func (s *bookingService) Create(userID int64, req request.CreateBookingRequest) error {
	bookingDate, err := time.Parse(dateFormat, req.BookingDate)
	if err != nil {
		return errors.New("invalid booking_date format, expected YYYY-MM-DD")
	}

	startTime, err := time.Parse(timeFormat, req.StartTime)
	if err != nil {
		return errors.New("invalid start_time format, expected HH:MM")
	}

	endTime, err := time.Parse(timeFormat, req.EndTime)
	if err != nil {
		return errors.New("invalid end_time format, expected HH:MM")
	}

	if !endTime.After(startTime) {
		return constant.ErrInvalidTimeRange
	}

	court, err := s.courtRepo.FindByID(req.CourtID)
	if err != nil {
		return err
	}

	durationHours := endTime.Sub(startTime).Hours()
	totalPrice := court.Price * durationHours

	booking := &entity.Bookings{
		UserID:          userID,
		CourtID:         req.CourtID,
		BookingDate:     bookingDate,
		StartTime:       req.StartTime,
		EndTime:         req.EndTime,
		Status:          "pending",
		TotalPrice:      totalPrice,
		PaymentDeadline: time.Now().Add(paymentDeadlineWindow),
	}

	err = s.bookingRepo.Create(booking)
	return err
}

func (s *bookingService) GetByUserID(userID int64) ([]response.BookingResponse, error) {
	bookings, err := s.bookingRepo.FindAllByUserID(userID)
	if err != nil {
		return nil, err
	}

	responses := make([]response.BookingResponse, 0, len(bookings))
	for _, b := range bookings {
		responses = append(responses, toBookingResponse(b))
	}

	return responses, nil
}

func (s *bookingService) GetByID(id int64, requesterID int64, isAdmin bool) (*response.BookingResponse, error) {
	booking, err := s.bookingRepo.FindByBookingID(id)
	if err != nil {
		return nil, err
	}

	if !isAdmin && booking.UserID != requesterID {
		return nil, constant.ErrForbidden
	}

	resp := toBookingResponse(*booking)
	return &resp, nil
}

func (s *bookingService) UpdateStatus(id int64, req request.UpdateBookingStatusRequest) error {
	return s.bookingRepo.UpdateStatusBooking(id, req.Status)
}

func toBookingResponse(b entity.Bookings) response.BookingResponse {
	return response.BookingResponse{
		ID:              b.ID,
		CourtID:         b.CourtID,
		BookingDate:     b.BookingDate.Format(dateFormat),
		StartTime:       b.StartTime,
		EndTime:         b.EndTime,
		Status:          b.Status,
		TotalPrice:      b.TotalPrice,
		PaymentDeadline: b.PaymentDeadline,
		CreatedAt:       b.CreatedAt,
	}
}
