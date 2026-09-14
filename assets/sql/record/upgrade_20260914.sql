BEGIN;

ALTER TABLE points_balance
    ADD CONSTRAINT uk_points_balance_user_id UNIQUE (user_id);

COMMIT;
