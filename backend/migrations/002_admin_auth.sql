-- Add user role + admin invite tokens.

ALTER TABLE users
    ADD COLUMN IF NOT EXISTS role VARCHAR(20) NOT NULL DEFAULT 'user';

CREATE TABLE IF NOT EXISTS admin_invites (
    id UUID PRIMARY KEY,
    token_hash TEXT UNIQUE NOT NULL,
    role VARCHAR(20) NOT NULL DEFAULT 'admin',
    expires_at TIMESTAMP NOT NULL,
    used_at TIMESTAMP NULL,
    used_by UUID NULL REFERENCES users(id) ON DELETE SET NULL,
    created_by UUID REFERENCES users(id) ON DELETE SET NULL,
    created_at TIMESTAMP NOT NULL
);

CREATE INDEX IF NOT EXISTS admin_invites_expires_at_idx ON admin_invites(expires_at);
CREATE INDEX IF NOT EXISTS admin_invites_used_at_idx ON admin_invites(used_at);
