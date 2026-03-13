-- Aurora AI Agent - User Table
-- PostgreSQL DDL

CREATE TABLE IF NOT EXISTS "user" (
    id          SERIAL          PRIMARY KEY,
    username    VARCHAR(64)     NOT NULL,
    password    VARCHAR(255)    NOT NULL,
    email       VARCHAR(128)    UNIQUE,
    phone       VARCHAR(20),
    birthday    DATE,
    user_prompt TEXT,
    create_time TIMESTAMPTZ     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    update_time TIMESTAMPTZ     NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT username_not_empty CHECK (username <> ''),
    CONSTRAINT password_not_empty CHECK (password <> '')
);

-- 唯一索引：用户名
CREATE UNIQUE INDEX IF NOT EXISTS idx_user_username ON "user" (username);

-- 普通索引：邮箱（查询登录用）
CREATE INDEX IF NOT EXISTS idx_user_email ON "user" (email);

-- 普通索引：手机号
CREATE INDEX IF NOT EXISTS idx_user_phone ON "user" (phone);

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
