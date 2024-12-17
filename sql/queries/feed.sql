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

-- name: GetFeed :one
SELECT * FROM feeds where feeds.url = $1;

-- name: MarkFeedFetched :exec


UPDATE feeds
SET last_fetched_at = NOW()
WHERE url = $1;

UPDATE feed_follows
SET updated_at = NOW()
WHERE feed_id = $1;



-- name: GetNextFeedToFetch :one
SELECT *
FROM feeds
ORDER BY 
    last_fetched_at ASC NULLS FIRST
LIMIT 1;