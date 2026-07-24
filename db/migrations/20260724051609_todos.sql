-- +goose Up
-- +goose StatementBegin
CREATE TABLE "todos" (
  "id" integer PRIMARY KEY,
  "title" varchar,
  "description" text,
  "completed" boolean DEFAULT false,
  "user_id" integer NOT NULL,
  "created_at" timestamptz,
  "due_date" date
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS todos;
-- +goose StatementEnd
