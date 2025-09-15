-- name: GetUserByID :one
SELECT * FROM users
WHERE id = $1
  AND deleted_at IS NULL;


-- name: CreateUser :one
INSERT INTO users (username, email, password, role)
VALUES ($1, $2, $3, $4)
RETURNING id;

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

-- name: RestoreUser :one
UPDATE users
SET deleted_at = NULL
WHERE id = $1
RETURNING *;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1;

-- name: ListUsers :many
SELECT * FROM users
WHERE deleted_at IS NULL
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: ListAllUsers :many
SELECT * FROM users
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1
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


-- name: GetTenantIdsByUserMail :one
SELECT
    u.id,
    u.username,
    u.password,
    u.role,
    COALESCE(array_agg(tu.tenant_id), '{}') AS tenant_ids
FROM users AS u
LEFT JOIN tenant.tenant_users AS tu ON u.id = tu.user_id
WHERE u.email = $1
GROUP BY u.id, u.username, u.password, u.role;

-- name: GetUsersByRole :many
SELECT 
   id as id,
   username,
   email,
   role,
   is_active,
   is_verified,
   created_at
FROM users
WHERE role = $1 AND deleted_at IS NULL
ORDER BY created_at DESC;