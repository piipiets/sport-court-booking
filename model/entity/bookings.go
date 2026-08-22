package entity

import "time"

type Bookings struct {
	ID              int64     `db:"id"`
	UserID          int64     `db:"user_id"`
	CourtID         int64     `db:"court_id"`
	BookingDate     time.Time `db:"booking_date"`
	StartTime       string    `db:"start_time"` // format "HH:MM:SS", Postgres TIME
	EndTime         string    `db:"end_time"`
	Status          string    `db:"status"` // pending | confirmed | cancelled | completed
	TotalPrice      float64   `db:"total_price"`
	PaymentDeadline time.Time `db:"payment_deadline"`
	CreatedAt       time.Time `db:"created_at"`
}
