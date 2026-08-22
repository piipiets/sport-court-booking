package response

import "time"

type BookingResponse struct {
	ID              int64     `json:"id"`
	CourtID         int64     `json:"court_id"`
	BookingDate     string    `json:"booking_date"`
	StartTime       string    `json:"start_time"`
	EndTime         string    `json:"end_time"`
	Status          string    `json:"status"`
	TotalPrice      float64   `json:"total_price"`
	PaymentDeadline time.Time `json:"payment_deadline"`
	CreatedAt       time.Time `json:"created_at"`
}
