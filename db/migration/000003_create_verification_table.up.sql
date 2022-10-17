CREATE TABLE "verification" (
    "username" varchar NOT NULL REFERENCES users(username),
    "is_verified" boolean DEFAULT FALSE,
    "verify_key" VARCHAR NOT NULL,
    "verefied_on" timestamptz NOT NULL DEFAULT ('0001-01-01 00:00:00Z'),
    "reset_key" VARCHAR,
    "reset_on" timestamptz NOT NULL DEFAULT ('0001-01-01 00:00:00Z'),
    "is_reset" boolean DEFAULT FALSE
);


