-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
CREATE  SCHEMA IF NOT EXISTS tenant;
-- +goose StatementEnd
CREATE TABLE IF NOT EXISTS tenant.tenants (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_name VARCHAR(100) UNIQUE NOT NULL,
    domain VARCHAR(100) UNIQUE NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMPTZ DEFAULT timezone('UTC', now()),
    updated_at TIMESTAMPTZ DEFAULT timezone('UTC', now()),
    deleted_at TIMESTAMPTZ DEFAULT NULL
);

CREATE TABLE IF NOT EXISTS tenant.tenant_users (
    tenant_id UUID NOT NULL REFERENCES tenant.tenants(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    PRIMARY KEY (tenant_id, user_id)
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
