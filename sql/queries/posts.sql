-- name: CreatePost :one
INSERT INTO posts (id, created_at, updated_at, title, url, description, published_at, feed_id)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING *;

-- name: GetPostsByUser :many
SELECT P.id, P.created_at, P.updated_at, P.title, P.url, P.description, P.published_at, P.feed_id
FROM posts as P
LEFT JOIN feed_follows as FF on P.feed_id = FF.feed_id 
WHERE FF.user_id = $1
Limit $2;