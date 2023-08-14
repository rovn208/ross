-- name: CreateVideo :one
INSERT INTO videos (
    title,
    stream_url,
    description,
    thumbnail_url,
    created_by
) VALUES (
  $1, $2, $3, $4, $5
) RETURNING *;

-- name: GetVideo :one
SELECT * FROM videos
WHERE id = $1 LIMIT 1;

-- name: DeleteVideo :exec
DELETE FROM videos
WHERE id = $1;

-- name: UpdateVideo :one
UPDATE videos
SET
    title = COALESCE(sqlc.narg(title), title),
    stream_url = COALESCE(sqlc.narg(stream_url), stream_url),
    description = COALESCE(sqlc.narg(description), description),
    thumbnail_url = COALESCE(sqlc.narg(thumbnail_url), thumbnail_url)
WHERE
    id = sqlc.arg(id)
RETURNING *;

-- name: GetListVideo :many
SELECT * FROM videos
ORDER BY created_at DESC
LIMIT $1
OFFSET $2;