-- name: InsertFeed :one
INSERT INTO feeds (name , url , user_id)
VALUES (
    $1,
    $2,
    $3
)
RETURNING *;

-- name: GetAllFeeds :many
SELECT feeds.name , feeds.url , users.name 
FROM feeds
LEFT JOIN users
ON feeds.user_id = users.id;
