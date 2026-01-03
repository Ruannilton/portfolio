-- +goose Up
-- +goose StatementBegin

CREATE TABLE users (
    id UUID PRIMARY KEY,
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    email VARCHAR(255) NOT NULL,
    password_hash VARCHAR(255), -- Nullable (Pode ser nulo se for OAuth)
    provider VARCHAR(50) NOT NULL, -- Ex: 'local', 'google', 'github'
    provider_id VARCHAR(255),   -- Nullable (ID externo do Google/Github)
    reset_token VARCHAR(255),   -- Nullable
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Índices para performance
-- 1. Garante unicidade e velocidade no login
CREATE UNIQUE INDEX idx_users_email ON users(email);

-- 2. Acelera o login via OAuth (buscar user pelo ID do provider)
CREATE INDEX idx_users_provider_id ON users(provider_id) WHERE provider_id IS NOT NULL;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_users_provider_id;
DROP INDEX IF EXISTS idx_users_email;
DROP TABLE IF EXISTS users;
-- +goose StatementEnd
