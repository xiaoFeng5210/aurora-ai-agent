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

-- Aurora AI Agent - Document（对话侧后续可单独建表关联本表 id）
CREATE TABLE IF NOT EXISTS document (
    id          SERIAL          PRIMARY KEY,
    user_id     INT             NOT NULL,
    display_name        VARCHAR(255)    NOT NULL,
    file_name   VARCHAR(512),
    create_time TIMESTAMPTZ     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    update_time TIMESTAMPTZ     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at  TIMESTAMPTZ,

    CONSTRAINT document_display_name_not_empty CHECK (display_name <> ''),
    CONSTRAINT document_file_name_not_empty CHECK (file_name <> '')
);

CREATE INDEX IF NOT EXISTS idx_document_deleted_at ON document (deleted_at);
CREATE INDEX IF NOT EXISTS idx_document_user_id ON document (user_id);
CREATE INDEX IF NOT EXISTS idx_document_create_time ON document (create_time DESC);

CREATE OR REPLACE TRIGGER trg_document_update_time
BEFORE UPDATE ON document
FOR EACH ROW EXECUTE FUNCTION set_update_time();

-- Aurora AI Agent - Messages
CREATE TABLE IF NOT EXISTS messages (
    id           SERIAL PRIMARY KEY,
    message_id   VARCHAR(128) NOT NULL,
    document_id  INT NOT NULL,
    role         VARCHAR(32) NOT NULL,
    content      TEXT NOT NULL DEFAULT '',
    tool_calls   JSONB NOT NULL DEFAULT '[]'::jsonb,
    create_time  TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    update_time  TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at   TIMESTAMPTZ,

    CONSTRAINT messages_message_id_not_empty CHECK (message_id <> ''),
    CONSTRAINT messages_role_valid CHECK (role IN ('user', 'assistant', 'system', 'tool')),
    CONSTRAINT messages_tool_calls_is_array CHECK (jsonb_typeof(tool_calls) = 'array')
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_messages_document_message_id
    ON messages (document_id, message_id);

CREATE INDEX IF NOT EXISTS idx_messages_document_create_time
    ON messages (document_id, create_time DESC, id DESC);

CREATE INDEX IF NOT EXISTS idx_messages_document_role
    ON messages (document_id, role);

CREATE INDEX IF NOT EXISTS idx_messages_deleted_at
    ON messages (deleted_at);

CREATE OR REPLACE TRIGGER trg_messages_update_time
BEFORE UPDATE ON messages
FOR EACH ROW EXECUTE FUNCTION set_update_time();
