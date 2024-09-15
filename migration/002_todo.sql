-- +goose Up
-- +goose StatementBegin
CREATE TABLE "todo"
(
    "id" UUID PRIMARY KEY DEFAULT GEN_RANDOM_UUID() CHECK ("id" <> '00000000-0000-0000-0000-000000000000'),
    "user_id" UUID REFERENCES "user" ON DELETE CASCADE,
    "description" VARCHAR(255) NOT NULL,
    "completed" BOOLEAN NOT NULL DEFAULT FALSE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS "todo";
-- +goose StatementEnd
