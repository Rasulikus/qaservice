-- +goose Up
-- +goose StatementBegin

CREATE TABLE IF NOT EXISTS questions (
    id         SERIAL PRIMARY KEY,
    text       TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS answers (
    id          SERIAL PRIMARY KEY,
    question_id INT NOT NULL REFERENCES questions(id) ON DELETE CASCADE,
    user_id     TEXT NOT NULL,
    text        TEXT NOT NULL,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS answers;
DROP TABLE IF EXISTS questions;

-- +goose StatementEnd
