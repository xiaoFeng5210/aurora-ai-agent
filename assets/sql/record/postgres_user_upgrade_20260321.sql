-- Aurora AI Agent
-- User module upgrade script
-- Generated at: 2026-03-21 00:00:00 Asia/Shanghai
--
-- Production execution example:
--   psql "host=<host> port=<port> user=<user> dbname=<dbname> password=<password> sslmode=disable" -f database/postgres_user_upgrade_20260321.sql

BEGIN;

ALTER TABLE "user"
    ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMPTZ;

ALTER TABLE "user"
    DROP CONSTRAINT IF EXISTS user_email_key;

DROP INDEX IF EXISTS idx_user_username;
DROP INDEX IF EXISTS idx_user_email;

CREATE UNIQUE INDEX IF NOT EXISTS idx_user_username_active
    ON "user" (username)
    WHERE deleted_at IS NULL;

CREATE UNIQUE INDEX IF NOT EXISTS idx_user_email_active
    ON "user" (email)
    WHERE deleted_at IS NULL AND email IS NOT NULL;

CREATE INDEX IF NOT EXISTS idx_user_phone ON "user" (phone);
CREATE INDEX IF NOT EXISTS idx_user_deleted_at ON "user" (deleted_at);

COMMIT;
