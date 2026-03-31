CREATE TABLE IF NOT EXISTS Authentication(
    sessionId INTEGER PRIMARY KEY,
    employeeId INTEGER,
    authToken TEXT,
    datetimeCreated DATETIME,
    lastAccessed DATETIME
)