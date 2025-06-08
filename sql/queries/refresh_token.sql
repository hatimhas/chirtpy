-- name: AddRefToken :one
INSERT INTO refresh_tokens (token, created_at, updated_at, expires_at,revoked_at, user_id)
VALUES (
  $1,                         -- token
  NOW(),                      -- created_at
  NOW(),                      -- updated_at
  NOW() + interval '60 days', -- expires_at
  NULL,                       -- revoked_at
  $2                          -- user_id
)
RETURNING *;


-- name: GetRefreshToken :one
SELECT * FROM refresh_tokens
WHERE token = $1 AND revoked_at IS NULL AND expires_at > NOW();


-- name: RevokeRefreshToken :exec
UPDATE refresh_tokens
SET revoked_at = $1,
    updated_at = $2
WHERE token = $3;
