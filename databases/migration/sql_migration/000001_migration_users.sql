-- +migrate Up
-- +migrate StatementBegin

CREATE TABLE IF NOT EXISTS users (
    id            BIGSERIAL PRIMARY KEY,
    name          VARCHAR(150)        NOT NULL,
    email         VARCHAR(255)        NOT NULL UNIQUE,
    password      VARCHAR(255)        NOT NULL,
    role          VARCHAR(20)         NOT NULL DEFAULT 'user' CHECK (role IN ('user', 'admin')),
    created_at    TIMESTAMPTZ         NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_users_email ON users (email);

-- +migrate StatementEnd

-- +migrate Down
-- +migrate StatementBegin

DROP TABLE IF EXISTS users;

-- +migrate StatementEnd