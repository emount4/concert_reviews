-- 0. Безопасное создание типа ENUM
DO $$ 
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'moderation_status_enum') THEN
        CREATE TYPE moderation_status_enum AS ENUM ('pending', 'approved', 'rejected');
    END IF;
END $$;

-- 1. Справочник городов
CREATE TABLE IF NOT EXISTS cities (
    city_id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE CHECK (name !~ '^\s*$'),
    slug VARCHAR(100) NOT NULL UNIQUE CHECK (slug !~ '^\s*$' AND slug ~ '^[a-z0-9_-]+$'),
    timezone VARCHAR(50) DEFAULT 'UTC+3' CHECK (timezone ~ '^UTC[+-]\d{1,2}(:\d{2})?$'),
    created_at TIMESTAMP DEFAULT NOW(),
    deleted_at TIMESTAMP DEFAULT NULL
);

-- 2. Роли пользователей
CREATE TABLE IF NOT EXISTS roles (
    role_id SERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL UNIQUE CHECK (name !~ '^\s*$'),
    permissions JSONB DEFAULT '{}'::jsonb,
    created_at TIMESTAMP DEFAULT NOW()
);

-- Наполнение ролей
INSERT INTO roles (name) 
VALUES ('user'), ('admin'), ('super_admin')
ON CONFLICT (name) DO NOTHING;

-- 3. Таблица пользователей
CREATE TABLE IF NOT EXISTS users (
    user_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) NOT NULL UNIQUE CHECK (
        email = lower(email)
        AND email !~ '\\s'
        AND position('@' in email) > 1
    ),
    password_hash VARCHAR(255) NOT NULL CHECK (password_hash !~ '^\s*$'),
    tg_id BIGINT UNIQUE,
    tg_username VARCHAR(100),
    role_id INTEGER DEFAULT 1 REFERENCES roles(role_id) ON DELETE SET DEFAULT,
    username VARCHAR(50) NOT NULL UNIQUE CHECK (username ~ '^[a-zA-Z0-9_]+$' AND LENGTH(username) >= 3),
    bio TEXT CHECK (LENGTH(bio) <= 500),
    avatar_url TEXT CHECK (avatar_url IS NULL OR (avatar_url !~ '^\s*$' AND LENGTH(avatar_url) <= 2048)),
    banner_url TEXT CHECK (banner_url IS NULL OR (banner_url !~ '^\s*$' AND LENGTH(banner_url) <= 2048)),
    is_email_verified BOOLEAN DEFAULT FALSE,
    is_active BOOLEAN DEFAULT TRUE,
    is_banned BOOLEAN DEFAULT FALSE,
    banned_by_user_id UUID REFERENCES users(user_id) ON DELETE SET NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);
 
-- 3.1 Токены (refresh / session)
CREATE TABLE IF NOT EXISTS auth_tokens (
    token_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    token VARCHAR(255) NOT NULL UNIQUE CHECK (token !~ '^\s*$'),
    expires_at TIMESTAMP NOT NULL,
    is_revoked BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT NOW(),
    CHECK (expires_at > created_at)
);

CREATE INDEX IF NOT EXISTS idx_auth_tokens_user_id ON auth_tokens(user_id);
CREATE INDEX IF NOT EXISTS idx_auth_tokens_expires_at ON auth_tokens(expires_at);

-- 4. Подтверждение почты
CREATE TABLE IF NOT EXISTS email_verifications (
    verification_id SERIAL PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    code VARCHAR(64) NOT NULL CHECK (code !~ '^\s*$'),
    type VARCHAR(20),
    expires_at TIMESTAMP NOT NULL CHECK (expires_at > created_at),
    created_at TIMESTAMP DEFAULT NOW()
);

-- 5. Логи входа
CREATE TABLE IF NOT EXISTS login_attempts (
    attempt_id SERIAL PRIMARY KEY,
    identifier VARCHAR(255) NOT NULL CHECK (identifier !~ '^\s*$'),
    success BOOLEAN DEFAULT FALSE,
    ip_address INET,
    user_agent TEXT,
    created_at TIMESTAMP DEFAULT NOW()
);

-- 6. Модерация профиля (ВНИМАТЕЛЬНО, ПРОВЕРЬ ЭТОТ БЛОК)
CREATE TABLE IF NOT EXISTS profile_moderation (
    moderation_id SERIAL PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    field_name VARCHAR(50) NOT NULL CHECK (field_name IN ('username', 'bio', 'avatar_url', 'banner_url')),
    old_value TEXT,
    new_value TEXT NOT NULL CHECK (new_value !~ '^\s*$'),
    status moderation_status_enum DEFAULT 'pending',
    moderated_by_user_id UUID REFERENCES users(user_id) ON DELETE SET NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    CHECK (old_value IS DISTINCT FROM new_value)
);

-- 7. Функция и Триггер
CREATE OR REPLACE FUNCTION update_timestamp() RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS trg_users_updated_at ON users;
CREATE TRIGGER trg_users_updated_at 
    BEFORE UPDATE ON users 
    FOR EACH ROW EXECUTE FUNCTION update_timestamp();