CREATE TABLE "users" (
  "id" bigserial PRIMARY KEY,
  "username" varchar UNIQUE NOT NULL,
  "email" varchar NOT NULL,
  "hashed_password" varchar NOT NULL,
  "full_name" varchar NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT '0001-01-01'
);

CREATE TABLE "videos" (
  "id" bigserial PRIMARY KEY,
  "title" varchar NOT NULL,
  "stream_url" varchar NOT NULL,
  "description" varchar,
  "thumbnail_url" varchar,
  "created_by" bigserial NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "follows" (
  "following_user_id" bigserial NOT NULL,
  "followed_user_id" bigserial NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE INDEX ON "users" ("id");

CREATE UNIQUE INDEX ON "users" ("email");

CREATE INDEX ON "videos" ("id");

CREATE INDEX "created_at_index" ON "videos" ("created_at");

CREATE INDEX ON "videos" ("stream_url");

CREATE INDEX ON "follows" ("followed_user_id");

CREATE INDEX ON "follows" ("following_user_id");

CREATE UNIQUE INDEX ON "follows" ("followed_user_id", "following_user_id");

ALTER TABLE "videos" ADD FOREIGN KEY ("created_by") REFERENCES "users" ("id");

ALTER TABLE "follows" ADD FOREIGN KEY ("following_user_id") REFERENCES "users" ("id");

ALTER TABLE "follows" ADD FOREIGN KEY ("followed_user_id") REFERENCES "users" ("id");
