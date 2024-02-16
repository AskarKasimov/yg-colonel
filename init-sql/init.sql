CREATE TABLE
    IF NOT EXISTS expressions (
        id SERIAL PRIMARY KEY,
        vanilla CHARACTER(256) NOT NULL,
        answer INT
    )