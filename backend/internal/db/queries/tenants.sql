-- name: CreateTenant :one
INSERT INTO tenant.tenants (
    tenant_name,
    domain,
    logo_url
) VALUES (
    $1, $2, $3
)
RETURNING *;

-- name: GetTenantByID :one
SELECT * FROM tenant.tenants
WHERE id = $1 AND deleted_at IS NULL;

-- name: UpdateTenant :one
UPDATE tenant.tenants
SET
    tenant_name = COALESCE($2, tenant_name),
    domain = COALESCE($3, domain),
    logo_url = COALESCE($4, logo_url),
    is_active = COALESCE($5, is_active),
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
