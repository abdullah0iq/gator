-- +goose Up
CREATE TABLE users (
        id UUID primary key,
        created_at TIMESTAMP not null ,
        updated_at TIMESTAMP not null,
        name  TEXT unique not null
);
CREATE TABLE feeds (
        name TEXT not null,
        url TEXT unique not null ,
        user_id UUID  not null,
        FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE feeds;
DROP TABLE users;