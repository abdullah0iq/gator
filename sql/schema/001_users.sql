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
        last_fetched_at TIMESTAMP,
        user_id UUID  not null,
        FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE feed_follows (
    id UUID PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    user_id UUID NOT NULL,
    feed_id TEXT NOT NULL,
    
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (feed_id) REFERENCES feeds(url),
    CONSTRAINT unique_user_feed UNIQUE (user_id, feed_id)
);
-- +goose Down
DROP TABLE feed_follows;
DROP TABLE feeds;
DROP TABLE users;