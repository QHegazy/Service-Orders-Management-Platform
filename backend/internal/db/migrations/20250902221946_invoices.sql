-- +goose Up
-- +goose StatementBegin
CREATE SCHEMA IF NOT EXISTS invoice;

CREATE TYPE invoice.invoice_status AS ENUM ('PENDING', 'PAID', 'CANCELLED');
CREATE TYPE invoice.payment_method AS ENUM ('CREDIT_CARD', 'BANK_TRANSFER', 'PAYPAL');
CREATE TYPE invoice.payment_status AS ENUM ('PENDING', 'COMPLETED', 'FAILED');

SELECT 'up SQL query';

CREATE TABLE invoice.invoices (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    ticket_id UUID NOT NULL REFERENCES ticket.tickets(id),  
    amount NUMERIC(10,2) NOT NULL CHECK (amount > 0),
    currency VARCHAR(3) NOT NULL DEFAULT 'USD',
    status invoice.invoice_status DEFAULT 'PENDING',          
    due_date DATE,
    created_at TIMESTAMPTZ DEFAULT timezone('UTC', now()),
    updated_at TIMESTAMPTZ DEFAULT timezone('UTC', now())
);

CREATE TABLE invoice.payments (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    invoice_id UUID NOT NULL REFERENCES invoice.invoices(id),  
    amount NUMERIC(10,2) NOT NULL CHECK (amount > 0),
    payment_date TIMESTAMPTZ DEFAULT (CURRENT_TIMESTAMP AT TIME ZONE 'UTC'),
    method invoice.payment_method DEFAULT 'CREDIT_CARD',
    status invoice.payment_status DEFAULT 'PENDING',
    created_at TIMESTAMPTZ DEFAULT timezone('UTC', now()),
    updated_at TIMESTAMPTZ DEFAULT timezone('UTC', now())
);

CREATE TRIGGER update_invoices_timestamp
BEFORE UPDATE ON invoice.invoices
FOR EACH ROW
EXECUTE FUNCTION updated_at_column();

CREATE TRIGGER update_payments_timestamp
BEFORE UPDATE ON invoice.payments
FOR EACH ROW
EXECUTE FUNCTION updated_at_column();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
DROP TABLE IF EXISTS invoice.payments;
DROP TABLE IF EXISTS invoice.invoices;
DROP TYPE IF EXISTS invoice.invoice_status;
DROP TYPE IF EXISTS invoice.payment_method;
DROP TYPE IF EXISTS invoice.payment_status;
DROP SCHEMA IF EXISTS invoice;
-- +goose StatementEnd
