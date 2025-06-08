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


