PRAGMA FOREIGN_KEYS = ON;

CREATE TABLE IF NOT EXISTS Items (
    itemId INTEGER PRIMARY KEY AUTOINCREMENT,
    itemName TEXT NOT NULL,
    upc TEXT NOT NULL,
    description TEXT,
    category TEXT,
    price REAL,
    warehouseStock INTEGER,
    salesFloorStock INTEGER
);