CREATE TABLE IF NOT EXISTS "users" (
                         "id" bigserial PRIMARY KEY,
                         "username" varchar UNIQUE NOT NULL,
                         "email" varchar NOT NULL,
                         "hash_password" varchar NOT NULL,
                         "password_changed_at" timestamp,
                         "full_name" varchar NOT NULL,
                         "created_at" timestamp DEFAULT (now()),
                         "updated_at" timestamp
);

CREATE TABLE IF NOT EXISTS "videos" (
                          "id" bigserial PRIMARY KEY,
                          "title" varchar NOT NULL,
                          "stream_url" varchar NOT NULL,
                          "description" varchar,
                          "thumbnail_url" varchar,
                          "created_by" bigserial NOT NULL,
                          "created_at" timestamp DEFAULT (now()),
                          "updated_at" timestamp
);

CREATE TABLE IF NOT EXISTS "follows" (
                           "following_user_id" bigserial NOT NULL,
                           "followed_user_id" bigserial NOT NULL,
                           "created_at" timestamp DEFAULT (now()),
                           "updated_at" timestamp
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
