-- Aurora AI Agent - User Table
-- PostgreSQL DDL

CREATE TABLE IF NOT EXISTS "user" (
    id          SERIAL          PRIMARY KEY,
    username    VARCHAR(64)     NOT NULL,
    password    VARCHAR(255)    NOT NULL,
    email       VARCHAR(128),
    phone       VARCHAR(20),
    birthday    DATE,
    user_prompt TEXT,
    create_time TIMESTAMPTZ     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    update_time TIMESTAMPTZ     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at  TIMESTAMPTZ,

    CONSTRAINT username_not_empty CHECK (username <> ''),
    CONSTRAINT password_not_empty CHECK (password <> '')
);

-- 唯一索引：用户名，仅对未删除账号生效
CREATE UNIQUE INDEX IF NOT EXISTS idx_user_username_active ON "user" (username) WHERE deleted_at IS NULL;

-- 唯一索引：邮箱，仅对未删除账号生效
CREATE UNIQUE INDEX IF NOT EXISTS idx_user_email_active ON "user" (email) WHERE deleted_at IS NULL AND email IS NOT NULL;

-- 普通索引：手机号
CREATE INDEX IF NOT EXISTS idx_user_phone ON "user" (phone);

-- 软删除索引
CREATE INDEX IF NOT EXISTS idx_user_deleted_at ON "user" (deleted_at);

-- 自动更新 update_time 的触发器函数
CREATE OR REPLACE FUNCTION set_update_time()
RETURNS TRIGGER AS $$
BEGIN
    NEW.update_time = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE TRIGGER trg_user_update_time
BEFORE UPDATE ON "user"
FOR EACH ROW EXECUTE FUNCTION set_update_time();
