package repository

import (
	"database/sql"
	"errors"

	"github.com/piipiets/sport-court-booking/helpers/constant"
	"github.com/piipiets/sport-court-booking/model/entity"
)

type PaymentRepository interface {
	Create(payment *entity.Payment) error
	FindByBookingID(bookingID int64) (*entity.Payment, error)
	GetAllPaymentsByUserId(isAdmin bool, userId int64) ([]entity.Payment, error)
}

type paymentRepository struct {
	db *sql.DB
}

func NewPaymentRepository(db *sql.DB) PaymentRepository {
	return &paymentRepository{db: db}
}

func (r *paymentRepository) Create(payment *entity.Payment) error {
	query := `
		INSERT INTO payments (booking_id, amount, method, status, paid_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`

	err := r.db.QueryRow(
		query,
		payment.BookingID,
		payment.Amount,
		payment.Method,
		payment.Status,
		payment.PaidAt,
	).Scan(&payment.ID)

	if err != nil {
		return err
	}

	return nil
}

func (r *paymentRepository) FindByBookingID(bookingID int64) (*entity.Payment, error) {
	query := `
		SELECT id, booking_id, amount, method, status, paid_at
		FROM payments
		WHERE booking_id = $1
	`

	var p entity.Payment
	err := r.db.QueryRow(query, bookingID).Scan(
		&p.ID, &p.BookingID, &p.Amount, &p.Method, &p.Status, &p.PaidAt,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, constant.ErrPaymentNotFound
	}
	if err != nil {
		return nil, err
	}

	return &p, nil
}

func (r *paymentRepository) GetAllPaymentsByUserId(isAdmin bool, userId int64) ([]entity.Payment, error) {
	var query string

	if isAdmin {
		query = `
			SELECT id, booking_id, amount, method, status, paid_at
			FROM payments	
		`
	} else {
		query = `
			SELECT p.id, p.booking_id, p.amount, p.method, p.status, p.paid_at
			FROM payments p
			join bookings b ON b.id = p.booking_id 
			join users u on b.user_id = u.id
			where u.id = $1 
		`
	}

	var rows *sql.Rows
	var err error
	if isAdmin {
		rows, err = r.db.Query(query)
	} else {
		rows, err = r.db.Query(query, userId)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	payments := make([]entity.Payment, 0)
	for rows.Next() {
		var p entity.Payment
		if err := rows.Scan(
			&p.ID, &p.BookingID, &p.Amount, &p.Method, &p.Status, &p.PaidAt,
		); err != nil {
			return nil, err
		}
		payments = append(payments, p)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return payments, nil
}
