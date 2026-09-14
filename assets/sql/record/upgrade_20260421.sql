



-- 添加消息点赞和点踩字段
BEGIN;

ALTER TABLE "messages"
    ADD COLUMN IF NOT EXISTS is_liked INT NOT NULL DEFAULT 0;
COMMIT;
