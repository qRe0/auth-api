-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS Users
(
    user_id  SERIAL PRIMARY KEY,
    name     VARCHAR(100) UNIQUE NOT NULL,
    phone    VARCHAR(15) UNIQUE  NOT NULL,
    email    VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255)        NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS Users;
-- +goose StatementEnd
