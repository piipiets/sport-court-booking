-- +migrate Up
-- +migrate StatementBegin

CREATE TABLE IF NOT EXISTS payments (
    id          BIGSERIAL PRIMARY KEY,
    booking_id  BIGINT          NOT NULL UNIQUE REFERENCES bookings (id) ON DELETE CASCADE,
    amount      NUMERIC(12, 2)  NOT NULL CHECK (amount >= 0),
    method      VARCHAR(20)     NOT NULL CHECK (method IN ('cash', 'transfer', 'qris')),
    status      VARCHAR(20)     NOT NULL DEFAULT 'unpaid'
                  CHECK (status IN ('unpaid', 'paid', 'refunded')),
    paid_at     TIMESTAMPTZ
);

-- +migrate StatementEnd

-- +migrate Down
-- +migrate StatementBegin

DROP TABLE IF EXISTS payments;

-- +migrate StatementEnd