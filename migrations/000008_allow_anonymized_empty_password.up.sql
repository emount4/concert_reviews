ALTER TABLE users DROP CONSTRAINT IF EXISTS users_password_hash_check;

ALTER TABLE users
ADD CONSTRAINT users_password_hash_check
CHECK (is_active = false OR password_hash !~ '^\s*$');
