CREATE TABLE trainings (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    title TEXT NOT NULL,
    description TEXT,
    questions JSON NOT NULL
);

CREATE TABLE area_trainings (
    area_id INTEGER NOT NULL,
    training_id INTEGER NOT NULL,
    PRIMARY KEY (area_id, training_id),
    FOREIGN KEY (area_id) REFERENCES areas(id),
    FOREIGN KEY (training_id) REFERENCES trainings(id)
);