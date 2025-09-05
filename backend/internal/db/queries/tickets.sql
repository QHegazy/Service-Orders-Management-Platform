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
RETURNING *;

-- name: GetTicketByID :one
SELECT * FROM ticket.tickets
WHERE id = $1 AND deleted_at IS NULL;

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
SELECT * FROM ticket.tickets
WHERE customer_id = $1 
ORDER BY created_at DESC
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

-- name: ListTicketsByTenantIDAndStatus :many
SELECT * FROM ticket.tickets
WHERE tenant_id = $1 AND status = $2
ORDER BY created_at DESC
LIMIT $3 OFFSET $4;


-- name: CreateComment :one
INSERT INTO ticket.comments (
    ticket_id,
    comment
) VALUES (
    $1, $2
)
RETURNING *;

-- name: ListCommentsByTicketID :many
SELECT * FROM ticket.comments
WHERE ticket_id = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;
