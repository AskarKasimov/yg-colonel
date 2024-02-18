CREATE TABLE expressions (
    id SERIAL PRIMARY KEY,
    vanilla TEXT NOT NULL,
    answer TEXT NOT NULL DEFAULT '',
    progress TEXT NOT NULL DEFAULT 'waiting',
    -- TODO: normalize by separating to another table
    --done processing waiting
    incomingDate TIMESTAMP NOT NULL DEFAULT now()
);
CREATE TABLE workers (
    id SERIAL PRIMARY KEY,
    name TEXT UNIQUE,
    isAlive BOOLEAN NOT NULL DEFAULT true,
    lastHeartbeat TIMESTAMP NOT NULL DEFAULT now()
);
CREATE TABLE workers_and_expressions (
    workerId INT NOT NULL,
    expressionId INT NOT NULL,
    FOREIGN KEY (workerId) REFERENCES workers(id),
    FOREIGN KEY (expressionId) REFERENCES expressions(id)
);