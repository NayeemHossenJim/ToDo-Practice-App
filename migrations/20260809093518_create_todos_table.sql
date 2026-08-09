-- +goose Up
CREATE TABLE "todos" (
  "id" bigserial PRIMARY KEY,
  "title" varchar,
  "description" text,
  "completed" boolean DEFAULT false,
  "user_id" bigint NOT NULL,
  "created_at" timestamptz DEFAULT (now()),
  "due_date" date
);

-- +goose Down
DROP TABLE "todos";