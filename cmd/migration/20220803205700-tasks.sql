
-- +migrate Up
CREATE TABLE IF NOT EXISTS tasks
(
    id BIGSERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    action_time timestamp NOT NULL,
    is_finished BOOLEAN DEFAULT FALSE,
    created_at timestamp,
    updated_at timestamp
);

-- +migrate Down
DROP TABLE IF EXISTS tasks;
