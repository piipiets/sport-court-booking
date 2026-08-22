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

type PaymentService interface {
	Create(userID int64, req request.CreatePaymentRequest) error
	GetByBookingID(bookingID int64, requesterID int64, isAdmin bool) (*response.PaymentResponse, error)
	GetAllPaymentsByUserID(isAdmin bool, userID int64) ([]response.PaymentResponse, error)
}

type paymentService struct {
	paymentRepo repository.PaymentRepository
	bookingRepo repository.BookingRepository
}

func NewPaymentService(paymentRepo repository.PaymentRepository, bookingRepo repository.BookingRepository) PaymentService {
	return &paymentService{paymentRepo: paymentRepo, bookingRepo: bookingRepo}
}

func (s *paymentService) Create(userID int64, req request.CreatePaymentRequest) error {
	booking, err := s.bookingRepo.FindByBookingID(req.BookingID)
	if err != nil {
		return err
	}

	if booking.UserID != userID {
		return constant.ErrForbidden
	}

	if req.Amount != booking.TotalPrice {
		return constant.ErrPaymentAmountMismatch
	}

	payment := &entity.Payment{
		BookingID: req.BookingID,
		Amount:    req.Amount,
		Method:    req.Method,
		Status:    "paid",
		PaidAt:    time.Now(),
	}

	paymentExist, err := s.paymentRepo.FindByBookingID(req.BookingID)
	if err != nil && !errors.Is(err, constant.ErrPaymentNotFound) {
		return err
	}
	if paymentExist != nil {
		return constant.ErrPaymentAlreadyExists
	}

	err = s.paymentRepo.Create(payment)
	return err
}

func (s *paymentService) GetByBookingID(bookingID int64, requesterID int64, isAdmin bool) (*response.PaymentResponse, error) {
	booking, err := s.bookingRepo.FindByBookingID(bookingID)
	if err != nil {
		return nil, err
	}

	if !isAdmin && booking.UserID != requesterID {
		return nil, constant.ErrForbidden
	}

	payment, err := s.paymentRepo.FindByBookingID(bookingID)
	if err != nil {
		return nil, err
	}

	resp := toPaymentResponse(*payment)
	return &resp, nil
}

func (s *paymentService) GetAllPaymentsByUserID(isAdmin bool, userID int64) ([]response.PaymentResponse, error) {
	payments, err := s.paymentRepo.GetAllPaymentsByUserId(isAdmin, userID)
	if err != nil {
		return nil, err
	}

	responses := make([]response.PaymentResponse, 0, len(payments))
	for _, b := range payments {
		responses = append(responses, toPaymentResponse(b))
	}

	return responses, nil
}

func toPaymentResponse(p entity.Payment) response.PaymentResponse {
	return response.PaymentResponse{
		ID:        p.ID,
		BookingID: p.BookingID,
		Amount:    p.Amount,
		Method:    p.Method,
		Status:    p.Status,
		PaidAt:    p.PaidAt,
	}
}
