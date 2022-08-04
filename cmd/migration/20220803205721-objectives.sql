
-- +migrate Up
CREATE TABLE IF NOT EXISTS objectives
(
    id BIGSERIAL PRIMARY KEY,
    task_id BIGINT NOT NULL,
    objective_name VARCHAR(255) NOT NULL,
    is_finished BOOLEAN DEFAULT FALSE,

    FOREIGN KEY ("task_id") REFERENCES "tasks" ("id")
);

-- +migrate Down
DROP TABLE IF EXISTS objectives;
