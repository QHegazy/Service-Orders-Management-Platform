-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
CREATE  SCHEMA IF NOT EXISTS tenant;
-- +goose StatementEnd
CREATE TABLE IF NOT EXISTS tenant.tenants (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_name VARCHAR(100) UNIQUE NOT NULL,
    domain VARCHAR(100) UNIQUE NOT NULL,
    logo_url VARCHAR(255),
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMPTZ DEFAULT timezone('UTC', now()),
    updated_at TIMESTAMPTZ DEFAULT timezone('UTC', now()),
    deleted_at TIMESTAMPTZ DEFAULT NULL
);

CREATE TRIGGER update_tenants_timestamp
BEFORE UPDATE ON tenant.tenants
FOR EACH ROW
EXECUTE FUNCTION updated_at_column();

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
DROP TABLE IF EXISTS tenant.tenants;
DROP SCHEMA IF EXISTS tenant;
-- +goose StatementEnd
