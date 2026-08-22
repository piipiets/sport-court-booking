package entity

import "time"

type Payment struct {
	ID        int64     `db:"id"`
	BookingID int64     `db:"booking_id"`
	Amount    float64   `db:"amount"`
	Method    string    `db:"method"` // cash | transfer | qris
	Status    string    `db:"status"` // unpaid | paid | refunded
	PaidAt    time.Time `db:"paid_at"`
}
