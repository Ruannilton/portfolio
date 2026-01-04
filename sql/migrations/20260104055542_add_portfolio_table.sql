-- +goose Up
-- +goose StatementBegin
CREATE TABLE profiles (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL,
    
    -- Dados Profissionais Básicos
    headline VARCHAR(255),
    bio TEXT,
    seniority VARCHAR(50),         -- Ex: 'JUNIOR', 'SENIOR'
    years_of_experience INT DEFAULT 0,
    open_to_work BOOLEAN DEFAULT TRUE,
    
    -- Dados de Contratação
    salary_expectation NUMERIC(15, 2), -- Suporta valores altos com 2 casas decimais
    currency CHAR(3) DEFAULT 'BRL',    -- Ex: 'BRL', 'USD'
    contract_type VARCHAR(50),         -- Ex: 'PJ', 'CLT'
    location VARCHAR(50),              -- Enum: 'ON_SITE', 'REMOTE', 'HYBRID', 'ANY'
    remote_only BOOLEAN DEFAULT FALSE, -- Flag booleana extra da struct
    
    -- Dados Complexos (JSONB)
    skills JSONB DEFAULT '[]'::jsonb,       -- Array de strings
    social_links JSONB DEFAULT '{}'::jsonb, -- Objeto
    experiences JSONB DEFAULT '[]'::jsonb,  -- Array de objetos
    projects JSONB DEFAULT '[]'::jsonb,     -- Array de objetos
    educations JSONB DEFAULT '[]'::jsonb,   -- Array de objetos

    -- Metadados
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    -- Constraints
    CONSTRAINT fk_user FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT uq_user_profile UNIQUE(user_id) -- Garante 1 perfil por usuário
);

-- Índices Estratégicos

-- 1. Índice GIN para buscar dentro do array de Skills (Performance Alta)
-- Permite queries como: WHERE skills ?& array['Go', 'Docker']
CREATE INDEX idx_profiles_skills ON profiles USING GIN (skills);

-- 2. Filtros comuns de recrutadores
CREATE INDEX idx_profiles_seniority ON profiles(seniority);
CREATE INDEX idx_profiles_location ON profiles(location);
CREATE INDEX idx_profiles_open_to_work ON profiles(open_to_work) WHERE open_to_work IS TRUE;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_profiles_open_to_work;
DROP INDEX IF EXISTS idx_profiles_location;
DROP INDEX IF EXISTS idx_profiles_seniority;
DROP INDEX IF EXISTS idx_profiles_skills;

DROP TABLE IF EXISTS profiles;
-- +goose StatementEnd
