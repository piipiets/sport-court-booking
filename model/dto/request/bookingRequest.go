package request

type CreateBookingRequest struct {
	CourtID     int64  `json:"court_id" binding:"required"`
	BookingDate string `json:"booking_date" binding:"required"`
	StartTime   string `json:"start_time" binding:"required"`
	EndTime     string `json:"end_time" binding:"required"`
}

type UpdateBookingStatusRequest struct {
	Status string `json:"status" binding:"required,oneof=confirmed cancelled completed"`
}
