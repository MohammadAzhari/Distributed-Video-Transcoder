-- name: CreateVideo :one
INSERT INTO videos (
  id, filename, status
) VALUES (
  $1, $2, 'new'
)
RETURNING *;


-- name: PublishVideo :one
UPDATE videos
SET status = 'done', worker_ip = $2, scales = $3
WHERE id = $1
RETURNING *;


-- name: GetVideo :one
SELECT * FROM videos
WHERE id = $1 LIMIT 1;