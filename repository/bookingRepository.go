package repository

import (
	"database/sql"
	"errors"

	"github.com/piipiets/sport-court-booking/helpers/constant"
	"github.com/piipiets/sport-court-booking/model/entity"
)

type BookingRepository interface {
	Create(booking *entity.Bookings) error
	FindByBookingID(id int64) (*entity.Bookings, error)
	FindAllByUserID(userID int64) ([]entity.Bookings, error)
	UpdateStatusBooking(id int64, status string) error
}

type bookingRepository struct {
	db *sql.DB
}

func NewBookingRepository(db *sql.DB) BookingRepository {
	return &bookingRepository{db: db}
}

func (r *bookingRepository) Create(booking *entity.Bookings) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// lock transaction to prevent race condition
	if _, err := tx.Exec(`SELECT pg_advisory_xact_lock($1)`, booking.CourtID); err != nil {
		return err
	}

	// check for overlapping booking
	var conflictID int64
	err = tx.QueryRow(`
		SELECT id FROM bookings
		WHERE court_id = $1 AND booking_date = $2 AND status != 'cancelled'
		AND start_time < $4 AND end_time > $3
		LIMIT 1
	`, booking.CourtID, booking.BookingDate, booking.StartTime, booking.EndTime).Scan(&conflictID)

	if err == nil {
		return constant.ErrBookingConflict
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return err
	}

	insertQuery := `
		INSERT INTO bookings (user_id, court_id, booking_date, start_time, end_time, status, total_price, payment_deadline)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, created_at
	`
	err = tx.QueryRow(
		insertQuery,
		booking.UserID,
		booking.CourtID,
		booking.BookingDate,
		booking.StartTime,
		booking.EndTime,
		booking.Status,
		booking.TotalPrice,
		booking.PaymentDeadline,
	).Scan(&booking.ID, &booking.CreatedAt)

	if err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (r *bookingRepository) FindByBookingID(id int64) (*entity.Bookings, error) {
	query := `
		SELECT id, user_id, court_id, booking_date, start_time, end_time, status, total_price, payment_deadline, created_at
		FROM bookings
		WHERE id = $1
	`

	var b entity.Bookings
	err := r.db.QueryRow(query, id).Scan(
		&b.ID, &b.UserID, &b.CourtID, &b.BookingDate, &b.StartTime, &b.EndTime,
		&b.Status, &b.TotalPrice, &b.PaymentDeadline, &b.CreatedAt,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, constant.ErrBookingNotFound
	}
	if err != nil {
		return nil, err
	}

	return &b, nil
}

func (r *bookingRepository) FindAllByUserID(userID int64) ([]entity.Bookings, error) {
	query := `
		SELECT id, user_id, court_id, booking_date, start_time, end_time, status, total_price, payment_deadline, created_at
		FROM bookings
		WHERE user_id = $1
		ORDER BY booking_date DESC, start_time DESC
	`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	bookings := make([]entity.Bookings, 0)
	for rows.Next() {
		var b entity.Bookings
		if err := rows.Scan(
			&b.ID, &b.UserID, &b.CourtID, &b.BookingDate, &b.StartTime, &b.EndTime,
			&b.Status, &b.TotalPrice, &b.PaymentDeadline, &b.CreatedAt,
		); err != nil {
			return nil, err
		}
		bookings = append(bookings, b)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return bookings, nil
}

func (r *bookingRepository) UpdateStatusBooking(id int64, status string) error {
	query := `UPDATE bookings SET status = $1 WHERE id = $2`

	result, err := r.db.Exec(query, status, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return constant.ErrBookingNotFound
	}

	return nil
}
