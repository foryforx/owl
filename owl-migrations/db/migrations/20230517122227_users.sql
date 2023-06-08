-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';

CREATE TABLE users (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  account_id uuid NOT NULL REFERENCES accounts(id),
  first_name varchar(200) NOT NULL,
  last_name varchar(200) NOT NULL,
  email varchar(200) NOT NULL UNIQUE,
  pwd varchar(500),
  retries int NOT NULL DEFAULT 0,
  is_super_admin boolean NOT NULL DEFAULT false,
  last_login_at timestamp,
  created_at timestamp NOT NULL DEFAULT now(),
  updated_at timestamp NOT NULL DEFAULT now(),
  deleted_at timestamp
);

CREATE INDEX users_account_id_id_deleted_at_idx ON users(account_id, id, deleted_at);
CREATE INDEX users_account_id_email__deleted_at_idx ON users(account_id, email, deleted_at);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
DROP TABLE users;
-- +goose StatementEnd
