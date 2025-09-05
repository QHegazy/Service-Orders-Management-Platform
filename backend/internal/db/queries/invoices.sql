-- name: CreateInvoice :one
INSERT INTO invoice.invoices (
    ticket_id,
    amount,
    currency,
    status,
    due_date
) VALUES (
    $1, $2, $3, $4, $5
)
RETURNING *;

-- name: GetInvoiceByID :one
SELECT * FROM invoice.invoices
WHERE id = $1 AND deleted_at IS NULL;

-- name: UpdateInvoice :one
UPDATE invoice.invoices
SET
    ticket_id = COALESCE($2, ticket_id),
    amount = COALESCE($3, amount),
    currency = COALESCE($4, currency),
    due_date = COALESCE($5, due_date),
    status = COALESCE($6, status),
    updated_at = timezone('UTC', now())
WHERE id = $1 AND deleted_at IS NULL
RETURNING *;


-- name: SoftDeleteInvoice :exec
UPDATE invoice.invoices
SET deleted_at = timezone('UTC', now())
WHERE id = $1 AND deleted_at IS NULL;

-- name: RestoreInvoice :one
UPDATE invoice.invoices
SET deleted_at = NULL
WHERE id = $1
RETURNING *;

-- name: DeleteInvoice :exec
DELETE FROM invoice.invoices
WHERE id = $1;

-- name: ListInvoicesByTenantID :many
SELECT 
    i.id, i.ticket_id, i.amount, i.currency, i.status, i.due_date,
    i.created_at, i.updated_at,
    t.title
FROM invoice.invoices AS i
JOIN ticket.tickets AS t ON i.ticket_id = t.id
WHERE t.tenant_id = $1
  AND i.deleted_at IS NULL
ORDER BY i.created_at DESC
OFFSET $2
LIMIT $3;


-- name: ListAllInvoices :many
SELECT *
FROM invoice.invoices 
ORDER BY created_at DESC
LIMIT $1
OFFSET $2;


-- name: CreatePayment :one
INSERT INTO invoice.payments (
    invoice_id,
    amount,
    payment_date,
    method,
    status
) VALUES (
    $1, $2, $3, $4, $5
)
RETURNING *;

-- name: GetPaymentByID :one
SELECT * FROM invoice.payments
WHERE id = $1;

-- name: ListPaymentsByInvoiceID :many
SELECT * FROM invoice.payments
WHERE invoice_id = $1 
ORDER BY created_at DESC
LIMIT $2 
OFFSET $3;

-- name: ListAllPayments :many
SELECT * FROM invoice.payments
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;
