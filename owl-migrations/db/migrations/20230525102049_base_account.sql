-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query- Base accounts';
INSERT INTO accounts(id, name) VALUES ('6f532e2d-a3b3-4f83-82c7-9fe0f6c93fa8', 'Base');
INSERT INTO users(id, account_id, first_name, last_name, email, pwd) 
  VALUES ('6f532e2d-a3b3-4f83-82c7-9fe0f6c93fa9', '6f532e2d-a3b3-4f83-82c7-9fe0f6c93fa8', 'Super', 'Admin', 'karuppaiah.al@gmail.com', '$2a$14$ItL4NbZrUoIRO/xaWv64Eequfgif1F.2pA9nuFt5R2bTP/B11q0UC');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query - Base accounts';
DELETE FROM users WHERE id = '6f532e2d-a3b3-4f83-82c7-9fe0f6c93fa9';
DELETE FROM accounts WHERE id = '6f532e2d-a3b3-4f83-82c7-9fe0f6c93fa8';
-- +goose StatementEnd
