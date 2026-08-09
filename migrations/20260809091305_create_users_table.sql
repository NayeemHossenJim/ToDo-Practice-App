-- +goose Up
CREATE TABLE "users" (
  "id" bigserial PRIMARY KEY,
  "name" varchar,
  "email" varchar,
  "password" varchar,
  "created_at" timestamptz DEFAULT (now())
);

-- +goose Down
DROP TABLE "users";