CREATE SCHEMA IF NOT EXISTS storage;

CREATE TABLE storage.notes (
    id SERIAL PRIMARY KEY,
    title VARCHAR(100),
    body TEXT,
    user_id INT NOT NULL,
    created_at BIGINT NOT NULL
)