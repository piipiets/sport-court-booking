package repository

import (
	"database/sql"
	"errors"

	"github.com/piipiets/sport-court-booking/helpers/constant"
	"github.com/piipiets/sport-court-booking/model/entity"
)

type CourtRepository interface {
	Create(court *entity.Courts) error
	FindAll() ([]entity.Courts, error)
	FindByID(id int64) (*entity.Courts, error)
	Update(court *entity.Courts) error
	Delete(id int64) error
}

type courtRepository struct {
	db *sql.DB
}

func NewCourtRepository(db *sql.DB) CourtRepository {
	return &courtRepository{db: db}
}

func (r *courtRepository) Create(court *entity.Courts) error {
	query := `
		INSERT INTO courts (name, type, price_per_hour, location)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at
	`

	err := r.db.QueryRow(
		query,
		court.Name,
		court.Type,
		court.Price,
		court.Location,
	).Scan(&court.ID, &court.CreatedAt)

	if err != nil {
		return err
	}

	return nil
}

func (r *courtRepository) FindAll() ([]entity.Courts, error) {
	query := `
		SELECT id, name, type, price_per_hour, location, created_at
		FROM courts
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	courts := make([]entity.Courts, 0)
	for rows.Next() {
		var c entity.Courts
		if err := rows.Scan(&c.ID, &c.Name, &c.Type, &c.Price, &c.Location, &c.CreatedAt); err != nil {
			return nil, err
		}
		courts = append(courts, c)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return courts, nil
}

func (r *courtRepository) FindByID(id int64) (*entity.Courts, error) {
	query := `
		SELECT id, name, type, price_per_hour, location, created_at
		FROM courts
		WHERE id = $1
	`

	var c entity.Courts
	err := r.db.QueryRow(query, id).Scan(
		&c.ID, &c.Name, &c.Type, &c.Price, &c.Location, &c.CreatedAt,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, constant.ErrCourtNotFound
	}
	if err != nil {
		return nil, err
	}

	return &c, nil
}

func (r *courtRepository) Update(court *entity.Courts) error {
	query := `
		UPDATE courts
		SET name = $1, type = $2, price_per_hour = $3, location = $4
		WHERE id = $5
		RETURNING id, name, type, price_per_hour, location, created_at
	`

	err := r.db.QueryRow(
		query,
		court.Name,
		court.Type,
		court.Price,
		court.Location,
		court.ID,
	).Scan(&court.ID, &court.Name, &court.Type, &court.Price, &court.Location, &court.CreatedAt)

	if errors.Is(err, sql.ErrNoRows) {
		return constant.ErrCourtNotFound
	}
	if err != nil {
		return err
	}

	return nil
}

func (r *courtRepository) Delete(id int64) error {
	query := `DELETE FROM courts WHERE id = $1`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return constant.ErrCourtNotFound
	}

	return nil
}
