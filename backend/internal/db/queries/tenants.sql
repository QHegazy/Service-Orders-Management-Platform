-- name: CreateTenant :one
INSERT INTO tenant.tenants (
    tenant_name,
    domain,
    email
) VALUES (
    $1, $2, $3
)
RETURNING id;

-- name: GetTenantByID :one
SELECT * FROM tenant.tenants
WHERE id = $1 AND deleted_at IS NULL;

-- name: UpdateTenant :one
UPDATE tenant.tenants
SET
    tenant_name = COALESCE($2, tenant_name),
    domain = COALESCE($3, domain),
    is_active = COALESCE($4, is_active),
    email = COALESCE($5, email),
    updated_at = timezone('UTC', now())
WHERE id = $1 AND deleted_at IS NULL
RETURNING *;


-- name: SoftDeleteTenant :exec
UPDATE tenant.tenants
SET deleted_at = timezone('UTC', now())
WHERE id = $1 AND deleted_at IS NULL;

-- name: RestoreTenant :one
UPDATE tenant.tenants
SET deleted_at = NULL
WHERE id = $1
RETURNING *;

-- name: DeleteTenant :exec
DELETE FROM tenant.tenants
WHERE id = $1;

-- name: DeactivateTenant :one
UPDATE tenant.tenants
SET is_active = FALSE,
    updated_at = timezone('UTC', now())
WHERE id = $1 AND deleted_at IS NULL
RETURNING *;

-- name: ActivateTenant :one
UPDATE tenant.tenants
SET is_active = TRUE,
    updated_at = timezone('UTC', now())
WHERE id = $1 AND deleted_at IS NULL
RETURNING *;


-- name: ListTenants :many
SELECT * FROM tenant.tenants
WHERE deleted_at IS NULL
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: ListAllTenants :many
SELECT * FROM tenant.tenants
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: GetTenantByName :one
SELECT * FROM tenant.tenants
WHERE tenant_name = $1 AND deleted_at IS NULL;

-- name: GetTenantByDomain :one
SELECT * FROM tenant.tenants
WHERE domain = $1 AND deleted_at IS NULL;


-- name: AddUserToTenant :exec
INSERT INTO tenant.tenant_users (tenant_id, user_id)
VALUES ($1, $2);

-- name: RemoveUserFromTenant :exec
DELETE FROM tenant.tenant_users
WHERE tenant_id = $1 AND user_id = $2;

-- name: ListUsersForTenant :many
SELECT
    u.id,
    u.username,
    u.email,
    u.role,
    u.is_active,
    u.is_verified,
    u.created_at,
    u.updated_at,
    u.deleted_at
FROM
    users AS u
JOIN
    tenant.tenant_users AS tu ON u.id = tu.user_id
WHERE
    tu.tenant_id = $1
ORDER BY
    u.created_at DESC
LIMIT $2 OFFSET $3;
