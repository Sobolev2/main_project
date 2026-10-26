-- +goose Up
-- +goose StatementBegin


CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    user_name VARCHAR(255) UNIQUE NOT NULL, CHECK (char_length(user_name) >= 3),
    first_name VARCHAR(255), CHECK (char_length(first_name) >= 2),
    last_name VARCHAR(255), CHECK (char_length(last_name) >= 3),
    password VARCHAR(255) NOT NULL CHECK (char_length(password) >= 8)
);



-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS users;
-- +goose StatementEnd
