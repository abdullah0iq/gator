-- name: CreateFeedFollow :many
WITH inserted_feed_follow AS (
    INSERT INTO feed_follows (id, created_at, updated_at, user_id, feed_id)
    VALUES ($1, NOW(), NOW(), $2, $3)
    RETURNING *
)
SELECT 
    inserted_feed_follow.*, -- All fields from the feed_follows table
    users.name AS user_name, -- The name of the user (from the users table)
    feeds.name AS feed_name  -- The name of the feed (from the feeds table)
FROM inserted_feed_follow
INNER JOIN users ON users.id = inserted_feed_follow.user_id
INNER JOIN feeds ON feeds.url = inserted_feed_follow.feed_id;


-- name: GetFeedFollowsForUser :many
SELECT feed_follows.* , users.name as user_name , feeds.name as feed_name 
from feed_follows 
INNER JOIN users on users.id = feed_follows.user_id
INNER JOIN feeds on feeds.url = feed_follows.feed_id
where users.id = $1;

-- name: UnFollowFeed :exec
DELETE from feed_follows 
where feed_follows.user_id = $1 and feed_follows.feed_id = $2;