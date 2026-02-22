CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE users (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  email VARCHAR(254) NOT NULL,
  created_at timestamptz NOT NULL DEFAULT now(),
  updated_at timestamptz,
  deleted_at timestamptz
);

CREATE UNIQUE INDEX IF NOT EXISTS ux_users_email ON users(email) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS ix_users_deleted_at ON users(deleted_at);

CREATE TABLE tokens (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  value VARCHAR(500) NOT NULL,
  created_at timestamptz NOT NULL DEFAULT now(),
  expired_at timestamptz NOT NULL, 
  user_id UUID NOT NULL,

  CONSTRAINT fk_tokens_user FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE mails (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  value TEXT NOT NULL,
  to_email VARCHAR(254) NOT NULL,
  created_at timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS ix_mails_to_email ON mails(to_email);
