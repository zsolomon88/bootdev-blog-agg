-- name: CreateFeed :one
INSERT INTO feeds (id, created_at, updated_at, name, url, user_id)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetAllFeeds :many
SELECT * FROM feeds;

-- name: UpdateFetchTime :one
UPDATE feeds
SET fetched_at = $2
WHERE id = $1
RETURNING *;

-- name: FetchNextFeeds :many
SELECT * FROM feeds
ORDER BY fetched_at NULLS FIRST
LIMIT $1;
