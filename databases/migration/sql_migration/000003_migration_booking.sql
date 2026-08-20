-- +migrate Up
-- +migrate StatementBegin

CREATE TABLE IF NOT EXISTS bookings (
    id                BIGSERIAL PRIMARY KEY,
    user_id           BIGINT          NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    court_id          BIGINT          NOT NULL REFERENCES courts (id) ON DELETE RESTRICT,
    booking_date      DATE            NOT NULL,
    start_time        TIME            NOT NULL,
    end_time          TIME            NOT NULL,
    status            VARCHAR(20)     NOT NULL DEFAULT 'pending'
                        CHECK (status IN ('pending', 'confirmed', 'cancelled', 'completed')),
    total_price       NUMERIC(12, 2)  NOT NULL CHECK (total_price >= 0),
    payment_deadline  TIMESTAMPTZ     NOT NULL,
    created_at        TIMESTAMPTZ     NOT NULL DEFAULT now(),
 
    CONSTRAINT chk_booking_time_order CHECK (end_time > start_time)
);

-- +migrate StatementEnd

-- +migrate Down
-- +migrate StatementBegin

DROP TABLE IF EXISTS bookings;

-- +migrate StatementEnd