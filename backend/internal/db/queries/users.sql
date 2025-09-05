-- name: GetUserByID :one
SELECT * FROM users
WHERE id = $1
  AND deleted_at IS NULL;

-- name: CreateUser :one
INSERT INTO users (username, email, tenant_id, password, role)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: UpdateUser :one
UPDATE users
SET 
    username    = COALESCE($2, username),
    email       = COALESCE($3, email),
    password    = COALESCE($4, password),
    role        = COALESCE($5, role),
    is_active   = COALESCE($6, is_active),
    is_verified = COALESCE($7, is_verified),
    updated_at  = timezone('UTC', now())
WHERE id = $1
  AND deleted_at IS NULL
RETURNING *;

-- name: SoftDeleteUser :exec
UPDATE users
SET deleted_at = timezone('UTC', now())
WHERE id = $1
  AND deleted_at IS NULL;

-- name: ListUsers :many
SELECT * FROM users
WHERE tenant_id = $1
  AND deleted_at IS NULL
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: ListUsersByRole :many
SELECT * FROM users
WHERE tenant_id = $1
  AND role = $2
  AND deleted_at IS NULL
ORDER BY created_at DESC
LIMIT $3 OFFSET $4;

-- name: ListAllUsers :many
SELECT * FROM users
WHERE tenant_id = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: ListAllUsersByRole :many
SELECT * FROM users
WHERE tenant_id = $1
  AND role = $2
ORDER BY created_at DESC
LIMIT $3 OFFSET $4;

-- name: CountUsers :one
SELECT COUNT(*) FROM users
WHERE tenant_id = $1
  AND deleted_at IS NULL;

-- name: CountUsersByRole :one
SELECT COUNT(*) FROM users
WHERE tenant_id = $1
  AND role = $2
  AND deleted_at IS NULL;

-- name: GetUserByUsername :one
SELECT * FROM users
WHERE username = $1
  AND deleted_at IS NULL;

-- name: VerifyUser :one
UPDATE users
SET is_verified = TRUE,
    updated_at  = timezone('UTC', now())
WHERE id = $1
  AND deleted_at IS NULL
RETURNING *;

-- name: ChangeUserPassword :one
UPDATE users
SET password   = $2,
    updated_at = timezone('UTC', now())
WHERE id = $1
  AND deleted_at IS NULL
RETURNING *;

-- name: DeactivateUser :one
UPDATE users
SET is_active  = FALSE,
    updated_at = timezone('UTC', now())
WHERE id = $1
  AND deleted_at IS NULL
RETURNING *;

-- name: ActivateUser :one
UPDATE users
SET is_active  = TRUE,
    updated_at = timezone('UTC', now())
WHERE id = $1
  AND deleted_at IS NULL
RETURNING *;