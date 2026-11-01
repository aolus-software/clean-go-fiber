DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'user_status') THEN
        CREATE TYPE user_status AS ENUM ('active', 'inactive', 'suspended', 'blocked');
    END IF;
END
$$;

DELETE FROM users WHERE email IS NULL OR password IS NULL;
UPDATE users SET name = '' WHERE name IS NULL;
UPDATE users SET status = 'active' WHERE status IS NULL OR status NOT IN ('active', 'inactive', 'suspended', 'blocked');

ALTER TABLE users DROP CONSTRAINT IF EXISTS users_email_key;
ALTER TABLE users
    ALTER COLUMN name SET NOT NULL,
    ALTER COLUMN email SET NOT NULL,
    ALTER COLUMN password SET NOT NULL,
    ALTER COLUMN status DROP DEFAULT,
    ALTER COLUMN status TYPE user_status USING status::user_status,
    ALTER COLUMN status SET DEFAULT 'active',
    ALTER COLUMN status SET NOT NULL,
    ADD COLUMN IF NOT EXISTS remark VARCHAR(255) DEFAULT NULL;

CREATE INDEX IF NOT EXISTS users_email_deleted_at_status_index ON users (email, deleted_at, status);
CREATE UNIQUE INDEX IF NOT EXISTS users_email_unique_live ON users (email) WHERE deleted_at IS NULL;

DELETE FROM roles WHERE deleted_at IS NOT NULL;
ALTER TABLE roles DROP COLUMN IF EXISTS deleted_at;
UPDATE roles SET name = 'admin'
WHERE name = 'superadmin' AND NOT EXISTS (SELECT 1 FROM roles WHERE name = 'admin');
UPDATE roles r SET name = left(r.name, 90) || '-' || left(r.id::text, 8)
WHERE EXISTS (
    SELECT 1 FROM roles o
    WHERE o.name = r.name AND (o.created_at, o.id) < (r.created_at, r.id)
);
ALTER TABLE roles ALTER COLUMN name TYPE VARCHAR(100);
ALTER TABLE roles ADD CONSTRAINT roles_name_unique UNIQUE (name);

DELETE FROM permissions WHERE deleted_at IS NOT NULL;
ALTER TABLE permissions DROP COLUMN IF EXISTS deleted_at;
DELETE FROM permissions
WHERE modules IN ('users', 'roles', 'permissions') AND name IN ('read', 'create', 'update', 'delete');
UPDATE permissions SET name = modules || ' ' || name;
UPDATE permissions p SET name = left(p.name, 240) || '-' || left(p.id::text, 8)
WHERE EXISTS (
    SELECT 1 FROM permissions o
    WHERE o.name = p.name AND (o.created_at, o.id) < (p.created_at, p.id)
);
ALTER TABLE permissions RENAME COLUMN modules TO "group";
ALTER TABLE permissions ALTER COLUMN "group" TYPE VARCHAR(100);
ALTER TABLE permissions ADD CONSTRAINT permissions_name_unique UNIQUE (name);

ALTER TABLE user_roles ADD COLUMN IF NOT EXISTS assigned_at TIMESTAMP NOT NULL DEFAULT NOW();

DELETE FROM email_verifications;
ALTER TABLE email_verifications RENAME COLUMN expires_in TO expired_at;
ALTER TABLE email_verifications
    ADD COLUMN used_at TIMESTAMP DEFAULT NULL,
    ALTER COLUMN created_at DROP NOT NULL,
    ALTER COLUMN updated_at DROP NOT NULL;
ALTER TABLE email_verifications ADD CONSTRAINT email_verification_token_unique UNIQUE (token);
CREATE INDEX email_verification_user_id_used_at_index ON email_verifications (user_id, used_at);

DROP TABLE IF EXISTS forgot_passwords;
CREATE TABLE password_reset_tokens (
    id         UUID         PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id    UUID         NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token      VARCHAR(255) NOT NULL,
    expired_at TIMESTAMP    NOT NULL,
    used_at    TIMESTAMP    DEFAULT NULL,
    created_at TIMESTAMP    DEFAULT NOW(),
    updated_at TIMESTAMP    DEFAULT NOW(),
    CONSTRAINT password_reset_token_token_unique UNIQUE (token)
);
CREATE INDEX password_reset_token_user_id_used_at_index ON password_reset_tokens (user_id, used_at);
