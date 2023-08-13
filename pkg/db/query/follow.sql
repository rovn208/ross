-- name: GetListFollower :many
SELECT * FROM follows
WHERE following_user_id = $1
LIMIT $2
OFFSET $3;

-- name: GetListFollowing :many
SELECT * FROM follows
WHERE followed_user_id = $1
LIMIT $2
OFFSET $3;

-- name: FollowUser :one
INSERT INTO follows (
    followed_user_id,
    following_user_id
) VALUES (
    $1, $2
) RETURNING *;

-- name: UnfollowUser :exec
DELETE FROM follows
WHERE followed_user_id = $1 AND following_user_id = $2;
