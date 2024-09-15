-- +goose Up
-- +goose StatementBegin
CREATE TABLE "user"
(
    "id" UUID PRIMARY KEY DEFAULT GEN_RANDOM_UUID() CHECK ("id" <> '00000000-0000-0000-0000-000000000000'),
    "email" TEXT UNIQUE NOT NULL CHECK ("email" <> ''),
    "password_hash" BYTEA,
    "created_at" TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    "updated_at" TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS "user";
-- +goose StatementEnd
