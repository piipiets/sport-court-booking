package request

type CreatePaymentRequest struct {
	BookingID int64   `json:"booking_id" binding:"required"`
	Amount    float64 `json:"amount" binding:"required,gt=0"`
	Method    string  `json:"method" binding:"required,oneof=cash transfer qris"`
}
