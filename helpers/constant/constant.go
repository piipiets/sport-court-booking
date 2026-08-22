package constant

import "errors"

var (
	ErrCourtNotFound    = errors.New("court not found")
	ErrUserNotFound     = errors.New("user not found")
	ErrBookingNotFound  = errors.New("booking not found")
	ErrBookingConflict  = errors.New("booking conflicts with an existing schedule on this court")
	ErrInvalidTimeRange = errors.New("end_time must be after start_time")
	ErrForbidden        = errors.New("you don't have access to this booking")
)
