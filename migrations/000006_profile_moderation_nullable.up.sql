BEGIN;

ALTER TABLE profile_moderation
    ALTER COLUMN new_value DROP NOT NULL;

ALTER TABLE profile_moderation
    DROP CONSTRAINT IF EXISTS profile_moderation_new_value_check;

ALTER TABLE profile_moderation
    ADD CONSTRAINT profile_moderation_new_value_check CHECK (
        (
            field_name = 'username'
            AND new_value IS NOT NULL
            AND new_value !~ '^\s*$'
        )
        OR (
            field_name <> 'username'
            AND (new_value IS NULL OR new_value !~ '^\s*$')
        )
    );

COMMIT;
