CREATE TABLE users (
    id serial PRIMARY KEY,
    email text UNIQUE NOT NULL,
    password varchar(100) NOT NULL,
    role text NOT NULL
);