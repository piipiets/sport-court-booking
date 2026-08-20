-- +migrate Up
-- +migrate StatementBegin
CREATE TABLE IF NOT EXISTS reviews (
    id          BIGSERIAL PRIMARY KEY,
    booking_id  BIGINT      NOT NULL REFERENCES bookings (id) ON DELETE CASCADE,
    user_id     BIGINT      NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    rating      SMALLINT    NOT NULL CHECK (rating BETWEEN 1 AND 5),
    comment     TEXT,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),

    -- Satu booking hanya boleh direview sekali
    CONSTRAINT uq_review_per_booking UNIQUE (booking_id)
);

-- +migrate StatementEnd

-- +migrate Down
-- +migrate StatementBegin

DROP TABLE IF EXISTS reviews;

-- +migrate StatementEnd