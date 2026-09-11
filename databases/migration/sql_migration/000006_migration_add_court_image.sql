-- +migrate Up
-- +migrate StatementBegin

ALTER TABLE courts ADD COLUMN IF NOT EXISTS image_url TEXT;

-- +migrate StatementEnd

-- +migrate Down
-- +migrate StatementBegin

ALTER TABLE courts DROP COLUMN IF EXISTS image_url;

-- +migrate StatementEnd