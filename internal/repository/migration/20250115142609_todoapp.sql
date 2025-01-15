-- +goose Up
-- +goose StatementBegin
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    username VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL
);

CREATE TABLE todo_lists (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description VARCHAR(255)
);

CREATE TABLE users_list (
    id SERIAL PRIMARY KEY,
    user_id int REFERENCES users (id) ON DELETE CASCADE NOT NULL,
    list_id int REFERENCES todo_lists (id) ON DELETE CASCADE NOT NULL
);

CREATE TABLE todo_items (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description VARCHAR(255),
    done BOOLEAN NOT NULL DEFAULT false 
);

CREATE TABLE lists_items (
    id SERIAL PRIMARY KEY,
    item_id int REFERENCES todo_items (id) on DELETE CASCADE NOT NULL,
    list_id int REFERENCES todo_lists (id) on DELETE CASCADE NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IS EXISTS lists_items;

DROP TABLE IS EXISTS users_list;

DROP TABLE IS EXISTS todo_lists;

DROP TABLE IS EXISTS users;

DROP TABLE IS EXISTS todo_items;
-- +goose StatementEnd
