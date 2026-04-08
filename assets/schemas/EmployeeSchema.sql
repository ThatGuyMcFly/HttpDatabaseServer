PRAGMA foreign_keys = ON;

DROP TABLE Employee;
DROP TABLE Role;
DROP TABLE Password;

CREATE TABLE IF NOT EXISTS Employee(
    employeeId INTEGER PRIMARY KEY,
    firstName TEXT NOT NULL,
    lastName TEXT NOT NULL,
    roleId INTEGER NOT NULL,
    FOREIGN KEY (roleId) REFERENCES Role(roleId)
);

CREATE TABLE IF NOT EXISTS Role(
    roleId INTEGER PRIMARY KEY AUTOINCREMENT,
    roleTitle TEXT UNIQUE NOT NULL
);

CREATE TABLE IF NOT EXISTS Password(
    employeeId INTEGER PRIMARY KEY,
    password TEXT NOT NULL,
    expired INTEGER NOT NULL,
    FOREIGN KEY (employeeId) REFERENCES Employee(employeeId)
);

INSERT INTO Role(roleId, roleTitle) VALUES (1,'Administrator');
INSERT INTO Role(roleId, roleTitle) VALUES (2,'Warehouse');
INSERT INTO Role(roleId, roleTitle) VALUES (3,'SalesFloor');
