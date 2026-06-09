CREATE TABLE user_trainings (
    user_uuid TEXT NOT NULL,
    training_id INTEGER NOT NULL,
    completed_at DATETIME NOT NULL,
    expires_at DATETIME NOT NULL,
    PRIMARY KEY (user_uuid, training_id),
    FOREIGN KEY (training_id) REFERENCES trainings(id)
);