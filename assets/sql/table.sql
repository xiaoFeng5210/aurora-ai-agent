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
    is_liked     INT NOT NULL DEFAULT 0,   -- 0: 未点赞, 1: 点赞, -1: 点踩
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









-- Aurora AI Agent - RAG Vectorization Status
CREATE TABLE IF NOT EXISTS rag_vectorization (
    id            SERIAL          PRIMARY KEY,
    user_id       INT             NOT NULL,
    file_name     VARCHAR(255)    NOT NULL,
    file_path     VARCHAR(1024)   NOT NULL,
    status        VARCHAR(32)     NOT NULL DEFAULT 'not_vectorized',
    error_message TEXT,
    create_time   TIMESTAMPTZ     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    update_time   TIMESTAMPTZ     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at    TIMESTAMPTZ,

    CONSTRAINT rag_vectorization_file_name_not_empty CHECK (file_name <> ''),
    CONSTRAINT rag_vectorization_file_path_not_empty CHECK (file_path <> ''),
    CONSTRAINT rag_vectorization_status_valid CHECK (status IN ('not_vectorized', 'vectorizing', 'completed', 'failed'))
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_rag_vectorization_user_file_path
    ON rag_vectorization (user_id, file_path);

CREATE INDEX IF NOT EXISTS idx_rag_vectorization_user_id ON rag_vectorization (user_id);
CREATE INDEX IF NOT EXISTS idx_rag_vectorization_status ON rag_vectorization (status);
CREATE INDEX IF NOT EXISTS idx_rag_vectorization_deleted_at ON rag_vectorization (deleted_at);

CREATE OR REPLACE TRIGGER trg_rag_vectorization_update_time
BEFORE UPDATE ON rag_vectorization
FOR EACH ROW EXECUTE FUNCTION set_update_time();









-- Aurora AI Agent - Card
CREATE TABLE IF NOT EXISTS card (
    id             SERIAL       PRIMARY KEY,
    user_id        INT          NOT NULL,
    content        TEXT         NOT NULL,
    tags           TEXT[]       NOT NULL DEFAULT ARRAY[]::TEXT[],
    external_links TEXT[]       NOT NULL DEFAULT ARRAY[]::TEXT[],
    internal_links TEXT[]       NOT NULL DEFAULT ARRAY[]::TEXT[],
    create_time    TIMESTAMPTZ  NOT NULL DEFAULT CURRENT_TIMESTAMP,
    update_time    TIMESTAMPTZ  NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at     TIMESTAMPTZ,

    CONSTRAINT card_content_not_empty CHECK (content <> '')
);

CREATE INDEX IF NOT EXISTS idx_card_user_create_time
    ON card (user_id, create_time DESC, id DESC);

CREATE INDEX IF NOT EXISTS idx_card_tags
    ON card USING GIN (tags);

CREATE OR REPLACE TRIGGER trg_card_update_time
BEFORE UPDATE ON card
FOR EACH ROW EXECUTE FUNCTION set_update_time();
