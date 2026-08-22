package response

import "time"

type PaymentResponse struct {
	ID        int64     `json:"id"`
	BookingID int64     `json:"booking_id"`
	Amount    float64   `json:"amount"`
	Method    string    `json:"method"`
	Status    string    `json:"status"`
	PaidAt    time.Time `json:"paid_at"`
}
