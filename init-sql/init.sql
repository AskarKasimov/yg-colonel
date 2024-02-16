CREATE TABLE expressions (
    id SERIAL PRIMARY KEY,
    vanilla TEXT NOT NULL,
    answer TEXT NOT NULL DEFAULT '',
    progress TEXT NOT NULL DEFAULT 'waiting'
);
ALTER TABLE expressions
ADD COLUMN incomingDate TIMESTAMP;
ALTER TABLE expressions
ALTER COLUMN incomingDate
SET DEFAULT now();