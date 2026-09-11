package entity

import "time"

type Courts struct {
	ID        int64     `db:"id"`
	Name      string    `db:"name"`
	Type      string    `db:"type"`
	Price     float64   `db:"price_per_hour"`
	Location  string    `db:"location"`
	ImageURL  *string   `db:"image_url"`
	CreatedAt time.Time `db:"created_at"`
}
