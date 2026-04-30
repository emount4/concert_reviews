DROP TRIGGER IF EXISTS trg_users_updated_at ON users;
DROP FUNCTION IF EXISTS update_timestamp();

DROP TABLE IF EXISTS profile_moderation;
DROP TABLE IF EXISTS login_attempts;
DROP TABLE IF EXISTS email_verifications;
DROP TABLE IF EXISTS auth_tokens;
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS roles;
DROP TABLE IF EXISTS cities;

DROP TYPE IF EXISTS moderation_status_enum;