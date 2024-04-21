CREATE TABLE expressions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    vanilla TEXT NOT NULL,
    answer TEXT NOT NULL DEFAULT '',
    progress TEXT NOT NULL DEFAULT 'waiting',
    -- TODO: normalize by separating to another table
    --done processing waiting
    incomingDate TIMESTAMP NOT NULL DEFAULT now()
);
CREATE TABLE workers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT UNIQUE,
    numberOfGoroutines INT NOT NULL,
    isAlive BOOLEAN NOT NULL DEFAULT true,
    lastHeartbeat TIMESTAMP NOT NULL DEFAULT now()
);
CREATE TABLE workers_and_expressions (
    workerId UUID NOT NULL,
    expressionId UUID NOT NULL,
    FOREIGN KEY (workerId) REFERENCES workers(id),
    FOREIGN KEY (expressionId) REFERENCES expressions(id)
);
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    login TEXT UNIQUE NOT NULL,
    passwordHash TEXT NOT NULL
);
CREATE TABLE users_and_expressions (
    userId UUID NOT NULL,
    expressionId UUID NOT NULL,
    FOREIGN KEY (userId) REFERENCES users(id),
    FOREIGN KEY (expressionId) REFERENCES expressions(id)
);