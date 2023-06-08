-- +goose Up
-- +goose StatementBegin
SELECT 'Accounts table creation';
CREATE TABLE accounts (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  name varchar(500) NOT NULL UNIQUE,
  created_at timestamp not NULL DEFAULT now(),
  updated_at timestamp not NULL DEFAULT now(),
  deleted_at timestamp
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'Accounts table drop';
DROP TABLE accounts;
-- +goose StatementEnd
