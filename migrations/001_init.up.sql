CREATE TABLE users
(
    id SERIAL PRIMARY KEY,
    name varchar(255) UNIQUE NOT NULL,
    password_hash varchar(255) NOT NULL
);

CREATE TABLE notes (
     id SERIAL PRIMARY KEY,
     user_id INTEGER,
     title varchar(255) NOT NULL,
     text TEXT,
     date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
     done boolean NOT NULL DEFAULT FALSE,
     FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);