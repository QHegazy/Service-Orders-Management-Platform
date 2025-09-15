-- name: CreateCustomer :one
INSERT INTO customers (
    last_name,
    first_name,
    username,
    email,
    password
) VALUES (
    $1, $2, $3, $4,$5
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
    password = COALESCE($6, password),
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

-- name: DeleteCustomer :exec
DELETE FROM customers
WHERE id = $1;

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

-- name: GetListOfCustomerByUserId :many 
SELECT
    c.id,
    c.last_name,
    c.first_name,
    c.username,
    c.email,
    c.created_at,
    c.updated_at,
    c.deleted_at
FROM
    customers AS c
JOIN
    ticket.tickets AS t ON c.id = t.customer_id
JOIN
    tenant.tenant_users AS tu ON t.tenant_id = tu.tenant_id
    
WHERE
    tu.user_id = $1
ORDER BY
    c.created_at DESC
LIMIT $2 OFFSET $3;