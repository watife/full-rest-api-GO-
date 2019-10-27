CREATE TABLE inbox (
    email text UNIQUE NOT NULL,
    send int NOT NULL,
    user_id int NOT NULL
)