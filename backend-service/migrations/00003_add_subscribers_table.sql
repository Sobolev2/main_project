-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS followers (
    id SERIAL PRIMARY KEY,
    follower_id INT NOT NULL,
    user_id INT NOT NULL,
    FOREIGN KEY (follower_id) REFERENCES users(id),
    FOREIGN KEY (user_id) REFERENCES users(id)
);

-- +goose StatementEnd
 
-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS follower;
-- +goose StatementEnd