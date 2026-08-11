CREATE TABLE "users" (
  "id" bigserial PRIMARY KEY,
  "name" varchar,
  "email" varchar,
  "password" varchar,
  "created_at" timestamptz DEFAULT (now())
);

CREATE TABLE "todos" (
  "id" bigserial PRIMARY KEY,
  "title" varchar,
  "description" text,
  "completed" boolean DEFAULT false,
  "user_id" bigint NOT NULL,
  "created_at" timestamptz DEFAULT (now()),
  "due_date" date
);

CREATE UNIQUE INDEX "idx_users_email" ON "users" ("email");

CREATE INDEX "idx_todos_user_completed" ON "todos" ("user_id", "completed");

CREATE INDEX "idx_todos_user_due_date" ON "todos" ("user_id", "due_date");

ALTER TABLE "todos" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id") DEFERRABLE INITIALLY IMMEDIATE;