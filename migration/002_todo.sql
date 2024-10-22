-- +goose Up
-- +goose StatementBegin
CREATE TABLE "todo"
(
    "id" UUID PRIMARY KEY DEFAULT GEN_RANDOM_UUID() CHECK ("id" <> '00000000-0000-0000-0000-000000000000'),
    "user_id" UUID REFERENCES "user" ON DELETE SET NULL,
    "subject" VARCHAR(255) NOT NULL,
    "description" VARCHAR(255) NOT NULL DEFAULT '',
    "priority" INT NOT NULL DEFAULT 0,
    "due_date" TIMESTAMPTZ,
    "completed" BOOLEAN NOT NULL DEFAULT FALSE,
    "created_at" TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    "updated_at" TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS "todo";
-- +goose StatementEnd
