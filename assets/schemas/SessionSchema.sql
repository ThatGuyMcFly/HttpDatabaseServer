CREATE TABLE IF NOT EXISTS Session(
    sessionId INTEGER PRIMARY KEY,
    employeeId INTEGER,
    authToken TEXT,
    datetimeCreated DATETIME,
    lastAccessed DATETIME
)