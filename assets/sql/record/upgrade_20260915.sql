-- Aurora AI Agent
-- Card module upgrade script: add is_content_visible column
-- Generated at: 2026-09-15 00:00:00 Asia/Shanghai
--
-- Production execution example:
--   psql "host=<host> port=<port> user=<user> dbname=<dbname> password=<password> sslmode=disable" -f assets/sql/record/upgrade_20260915.sql

BEGIN;

ALTER TABLE "card"
    ADD COLUMN IF NOT EXISTS is_content_visible BOOLEAN NOT NULL DEFAULT TRUE;

COMMIT;
