CREATE INDEX IF NOT EXISTS user_roles_role_id_index ON user_roles (role_id);
CREATE INDEX IF NOT EXISTS role_permissions_permission_id_index ON role_permissions (permission_id);
CREATE INDEX IF NOT EXISTS users_created_at_live_index ON users (created_at) WHERE deleted_at IS NULL;
