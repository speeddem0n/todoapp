CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    username VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL
);

CREATE TABLE todo_lists (
    id SERIAL PRIMARY KEY,,
    title VARCHAR(255) NOT NULL,
    description VARCHAR(255)
);

CREATE TABLE users_list (
    id SERIAL PRIMARY KEY,,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (list_id) REFERENCES todo_lists(id) ON DELETE CASCADE
);

CREATE TABLE todo_items (
    id SERIAL PRIMARY KEY,,
    title VARCHAR(255) NOT NULL,
    description VARCHAR(255),
    done BOOLEAN NOT NULL DEFAULT false 
);

CREATE TABLE lists_items (
    id SERIAL PRIMARY KEY,,
    FOREIGN KEY (item_id) REFERENCES todo_items (id) on DELETE CASCADE,
    FOREIGN KEY (list_id) REFERENCES todo_lists (id) on DELETE CASCADE
);