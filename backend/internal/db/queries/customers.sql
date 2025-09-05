-- name: CreateCustomer :one
INSERT INTO customers (
    last_name,
    first_name,
    username,
    email
) VALUES (
    $1, $2, $3, $4
)
RETURNING *;

-- name: GetCustomerByID :one
SELECT * FROM customers
WHERE id = $1 AND deleted_at IS NULL;

-- name: UpdateCustomer :one
UPDATE customers
SET
    last_name = COALESCE($2, last_name),
    first_name = COALESCE($3, first_name),
    username = COALESCE($4, username),
    email = COALESCE($5, email),
    updated_at = timezone('UTC', now())
WHERE id = $1 AND deleted_at IS NULL
RETURNING *;

-- name: SoftDeleteCustomer :exec
UPDATE customers
SET deleted_at = timezone('UTC', now())
WHERE id = $1 AND deleted_at IS NULL;

-- name: RestoreCustomer :one
UPDATE customers
SET deleted_at = NULL
WHERE id = $1
RETURNING *;


-- name: ListAllCustomers :many
SELECT * FROM customers
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: GetCustomerByUsername :one
SELECT * FROM customers
WHERE username = $1 AND deleted_at IS NULL;

-- name: GetCustomerByEmail :one
SELECT * FROM customers
WHERE email = $1 AND deleted_at IS NULL;
