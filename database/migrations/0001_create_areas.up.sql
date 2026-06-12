CREATE TABLE areas (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    description TEXT,
    ldapgroup TEXT,
    photourl TEXT NOT NULL
);