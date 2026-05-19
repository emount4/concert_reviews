BEGIN;

DELETE FROM profile_moderation
WHERE new_value IS NULL;

ALTER TABLE profile_moderation
    DROP CONSTRAINT IF EXISTS profile_moderation_new_value_check;

ALTER TABLE profile_moderation
    ALTER COLUMN new_value SET NOT NULL;

ALTER TABLE profile_moderation
    ADD CONSTRAINT profile_moderation_new_value_check CHECK (new_value !~ '^\s*$');

COMMIT;
