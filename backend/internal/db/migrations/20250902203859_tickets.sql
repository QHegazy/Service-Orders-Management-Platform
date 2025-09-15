-- +goose Up
SELECT 'up SQL query';
-- +goose StatementBegin

CREATE SCHEMA IF NOT EXISTS ticket;

CREATE TYPE  ticket.ticket_status AS ENUM ('OPEN', 'IN_PROGRESS', 'RESOLVED', 'CLOSED', 'REOPENED');
CREATE TYPE  ticket.ticket_priority AS ENUM ('LOW', 'MEDIUM', 'HIGH', 'URGENT', 'CRITICAL');

CREATE TABLE ticket.tickets (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID NOT NULL REFERENCES tenant.tenants(id) ON DELETE CASCADE,
    customer_id UUID NOT NULL REFERENCES customers(id),
    assigned_to UUID REFERENCES users(id), 
    title VARCHAR(100) NOT NULL,
    description TEXT,
    status    ticket.ticket_status DEFAULT 'OPEN', 
    priority  ticket.ticket_priority DEFAULT 'MEDIUM',
    created_at TIMESTAMPTZ DEFAULT timezone('UTC', now()),
    updated_at TIMESTAMPTZ DEFAULT timezone('UTC', now()),
    closed_at TIMESTAMPTZ
);

CREATE TABLE ticket.comments (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    ticket_id UUID NOT NULL REFERENCES ticket.tickets(id) ON DELETE CASCADE,
    author_type VARCHAR(20) NOT NULL CHECK (author_type IN ('USER', 'CUSTOMER')),
    author_id UUID NOT NULL,
    comment TEXT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT timezone('UTC', now())
);

CREATE TRIGGER update_tickets_timestamp
BEFORE UPDATE ON ticket.tickets
FOR EACH ROW
EXECUTE FUNCTION updated_at_column();

CREATE TRIGGER update_comments_timestamp
BEFORE UPDATE ON ticket.comments
FOR EACH ROW
EXECUTE FUNCTION updated_at_column();

-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';

DROP TABLE IF EXISTS ticket.comments;
DROP TABLE IF EXISTS ticket.tickets;
DROP TYPE IF EXISTS ticket.ticket_status;
DROP TYPE IF EXISTS ticket.ticket_priority;
DROP SCHEMA IF EXISTS ticket;
-- +goose StatementEnd