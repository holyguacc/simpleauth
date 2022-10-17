CREATE TABLE "posts" (
    "id" VARCHAR PRIMARY KEY,
    "title" VARCHAR NOT NULL,
    "post_description" VARCHAR NOT NULL,
    "author_name" VARCHAR NOT NULL,
    "post_date" timestamptz NOT NULL DEFAULT (now())
);