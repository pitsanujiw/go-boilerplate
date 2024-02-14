BEGIN;

CREATE TABLE
    IF NOT EXISTS users(
        user_id uuid PRIMARY KEY,
        username VARCHAR (50) UNIQUE NOT NULL,
        password VARCHAR (50) NOT NULL,
        email VARCHAR (300) UNIQUE NOT NULL
    );

COMMIT;