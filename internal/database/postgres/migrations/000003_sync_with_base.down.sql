DROP TABLE IF EXISTS password_reset_tokens;
CREATE TABLE IF NOT EXISTS forgot_passwords (
    id         UUID         PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id    UUID         NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token      VARCHAR(255) NOT NULL,
    expires_at TIMESTAMP    NOT NULL,
    created_at TIMESTAMP    NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP    NOT NULL DEFAULT NOW()
);

DELETE FROM email_verifications;
DROP INDEX IF EXISTS email_verification_user_id_used_at_index;
ALTER TABLE email_verifications DROP CONSTRAINT IF EXISTS email_verification_token_unique;
ALTER TABLE email_verifications DROP COLUMN IF EXISTS used_at;
ALTER TABLE email_verifications RENAME COLUMN expired_at TO expires_in;
UPDATE email_verifications SET created_at = NOW() WHERE created_at IS NULL;
UPDATE email_verifications SET updated_at = NOW() WHERE updated_at IS NULL;
ALTER TABLE email_verifications
    ALTER COLUMN created_at SET NOT NULL,
    ALTER COLUMN updated_at SET NOT NULL;

ALTER TABLE user_roles DROP COLUMN IF EXISTS assigned_at;

ALTER TABLE permissions DROP CONSTRAINT IF EXISTS permissions_name_unique;
UPDATE permissions SET name = substring(name FROM length("group") + 2)
WHERE left(name, length("group") + 1) = "group" || ' ';
ALTER TABLE permissions ALTER COLUMN "group" TYPE VARCHAR(255);
ALTER TABLE permissions RENAME COLUMN "group" TO modules;
ALTER TABLE permissions ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMP DEFAULT NULL;

ALTER TABLE roles DROP CONSTRAINT IF EXISTS roles_name_unique;
ALTER TABLE roles ALTER COLUMN name TYPE VARCHAR(255);
ALTER TABLE roles ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMP DEFAULT NULL;

DROP INDEX IF EXISTS users_email_unique_live;
DROP INDEX IF EXISTS users_email_deleted_at_status_index;
ALTER TABLE users
    DROP COLUMN IF EXISTS remark,
    ALTER COLUMN status DROP DEFAULT,
    ALTER COLUMN status TYPE VARCHAR(255) USING status::text,
    ALTER COLUMN status DROP NOT NULL,
    ALTER COLUMN name DROP NOT NULL,
    ALTER COLUMN email DROP NOT NULL,
    ALTER COLUMN password DROP NOT NULL;
ALTER TABLE users ADD CONSTRAINT users_email_key UNIQUE (email);

DROP TYPE IF EXISTS user_status;
