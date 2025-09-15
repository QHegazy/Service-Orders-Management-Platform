-- name: CreateTicket :one
INSERT INTO ticket.tickets (
    tenant_id,
    customer_id,
    assigned_to,
    title,
    description,
    status,
    priority
) VALUES (
    $1, $2, $3, $4, $5, $6, $7
)
RETURNING id;



-- name: GetTicketByID :one
SELECT * FROM ticket.tickets
WHERE id = $1;

-- name: UpdateTicket :one
UPDATE ticket.tickets
SET
    assigned_to = COALESCE($2, assigned_to),
    title = COALESCE($3, title),
    description = COALESCE($4, description),
    status = COALESCE($5, status),
    priority = COALESCE($6, priority),
    updated_at = timezone('UTC', now()),
    closed_at = CASE WHEN $5 = 'RESOLVED' OR $5 = 'CLOSED' THEN timezone('UTC', now()) ELSE closed_at END
WHERE id = $1 
RETURNING *;



-- name: DeleteTicket :exec
DELETE FROM ticket.tickets
WHERE id = $1;

-- name: ListTicketsByTenantID :many
SELECT * FROM ticket.tickets
WHERE tenant_id = $1 
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: ListTicketsByCustomerID :many
-- name: ListTicketsByCustomerID :many
SELECT 
    t.id,
    t.tenant_id,
    t.customer_id,
    t.created_at,
    t.closed_at,
    t.status,
    t.priority,
    u.username AS assigned_username
FROM ticket.tickets AS t
JOIN users AS u ON t.assigned_to = u.id
WHERE t.customer_id = $1
ORDER BY t.created_at DESC
LIMIT $2 OFFSET $3;
-- name: ListTicketsByAssignedTo :many
SELECT * FROM ticket.tickets
WHERE assigned_to = $1 
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: ListTicketsByStatus :many
SELECT * FROM ticket.tickets
WHERE tenant_id = $1 AND status = $2 
ORDER BY created_at DESC
LIMIT $3;

-- name: ListTicketsByPriority :many
SELECT * FROM ticket.tickets
WHERE tenant_id = $1 AND priority = $2
ORDER BY created_at DESC
LIMIT $3;



-- name: CreateComment :one
INSERT INTO ticket.comments (
    ticket_id,
    author_type,
    author_id,
    comment
) VALUES (
    $1, $2, $3, $4
)
RETURNING id;

-- name: ListCommentsByTicketID :many
SELECT 
    c.*,
    u.username
FROM ticket.comments c
LEFT JOIN users u ON c.user_id = u.id
WHERE c.ticket_id = $1
ORDER BY c.created_at DESC
LIMIT $2 OFFSET $3;


-- name: ListTicketsByUserId :many 
SELECT
    t.id,
    t.tenant_id,
    t.customer_id,
    t.assigned_to,
    t.title,
    t.description,
    t.status,
    t.priority,
    t.created_at,
    t.updated_at,
    t.closed_at,
    c.first_name AS customer_first_name,
    c.last_name AS customer_last_name,
    u.username AS assigned_username
FROM
    ticket.tickets AS t
JOIN
    tenant.tenant_users AS tu ON t.tenant_id = tu.tenant_id
LEFT JOIN
    users AS u ON t.assigned_to = u.id
JOIN  
    customers AS c ON t.customer_id = c.id
WHERE
    tu.user_id = $1
ORDER BY
    t.created_at DESC
LIMIT $2 OFFSET $3;



-- name: ListTicketsByCustomerId :many
SELECT
    t.id,
    t.tenant_id,
    t.customer_id,
    t.assigned_to,
    t.title,
    t.description,
    t.status,
    t.priority,
    t.created_at,
    t.updated_at,
    t.closed_at,
    c.username AS customer_username,
    u.username AS assigned_username
FROM
    ticket.tickets AS t
LEFT JOIN
    users AS u ON t.assigned_to = u.id
LEFT JOIN
    customers AS c ON t.customer_id = c.id
WHERE
    t.customer_id = $1
ORDER BY
    t.created_at DESC
LIMIT $2 OFFSET $3;