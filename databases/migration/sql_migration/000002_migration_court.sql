-- +migrate Up
-- +migrate StatementBegin

CREATE TABLE IF NOT EXISTS courts (
    id              BIGSERIAL PRIMARY KEY,
    name            VARCHAR(150)     NOT NULL,
    type            VARCHAR(50)      NOT NULL,
    price_per_hour  NUMERIC(12, 2)   NOT NULL CHECK (price_per_hour >= 0),
    location        VARCHAR(255)     NOT NULL,
    created_at      TIMESTAMPTZ      NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_users_email ON users (email);

-- +migrate StatementEnd

-- +migrate Down
-- +migrate StatementBegin

DROP TABLE IF EXISTS courts;

-- +migrate StatementEnd