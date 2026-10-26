-- +goose Up
-- +goose StatementBegin


CREATE TABLE IF NOT EXISTS friends (
    id SERIAL PRIMARY KEY,
    user_first INT NOT NULL,
    user_second INT NOT NULL,
    FOREIGN KEY (user_first) REFERENCES users(id),
    FOREIGN KEY (user_second) REFERENCES users(id)
);
CREATE TABLE IF NOT EXISTS posts (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL,
    content VARCHAR(10000) NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id)
);
CREATE TABLE IF NOT EXISTS comments (
    id SERIAL PRIMARY KEY,
    post_id INT NOT NULL,
    user_id INT NOT NULL,
    content VARCHAR(10000) NOT NULL,
    FOREIGN KEY (post_id) REFERENCES posts(id),
    FOREIGN KEY (user_id) REFERENCES users(id)
);
CREATE TABLE IF NOT EXISTS likes (
    id SERIAL PRIMARY KEY,
    post_id INT NOT NULL,
    user_id INT NOT NULL,
    FOREIGN KEY (post_id) REFERENCES posts(id),
    FOREIGN KEY (user_id) REFERENCES users(id),
    UNIQUE(user_id, post_id)
);
CREATE TABLE IF NOT EXISTS chats (
    id SERIAL PRIMARY KEY,
    user_first INT NOT NULL, 
    user_second INT NOT NULL,
    FOREIGN KEY (user_first) REFERENCES users(id),
    FOREIGN KEY (user_second) REFERENCES users(id)
);
CREATE TABLE IF NOT EXISTS messages (
    id SERIAL PRIMARY KEY,
    chat_id INT NOT NULL,
    sender_id INT NOT NULL,
    content VARCHAR(10000) NOT NULL,
    mark_read BOOLEAN DEFAULT false,
    FOREIGN KEY (chat_id) REFERENCES chats(id),
    FOREIGN KEY (sender_id) REFERENCES users(id)
);


-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS chats, messages, likes, comments, posts, friends;
-- +goose StatementEnd
