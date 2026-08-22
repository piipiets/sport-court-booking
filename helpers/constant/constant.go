package constant

import "errors"

var (
	ErrCourtNotFound         = errors.New("court not found")
	ErrUserNotFound          = errors.New("user not found")
	ErrBookingNotFound       = errors.New("booking not found")
	ErrBookingConflict       = errors.New("booking conflicts with an existing schedule on this court")
	ErrInvalidTimeRange      = errors.New("end_time must be after start_time")
	ErrForbidden             = errors.New("you don't have access to this booking")
	ErrPaymentNotFound       = errors.New("payment not found")
	ErrPaymentAlreadyExists  = errors.New("a payment already exists for this booking")
	ErrPaymentAmountMismatch = errors.New("payment amount does not match booking total price")
)
