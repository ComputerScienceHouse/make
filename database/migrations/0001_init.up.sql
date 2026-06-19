CREATE TABLE areas (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    description TEXT,
    ldapgroup TEXT,
    photourl TEXT NOT NULL
);

CREATE TABLE trainings (
    id SERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    description TEXT,
    questions JSONB NOT NULL
);

CREATE TABLE area_trainings (
    area_id INTEGER NOT NULL,
    training_id INTEGER NOT NULL,
    PRIMARY KEY (area_id, training_id),
    FOREIGN KEY (area_id) REFERENCES areas(id) ON DELETE CASCADE,
    FOREIGN KEY (training_id) REFERENCES trainings(id) ON DELETE CASCADE
);

CREATE TABLE user_trainings (
    user_uuid UUID NOT NULL,
    training_id INTEGER NOT NULL,
    completed_at TIMESTAMP NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    PRIMARY KEY (user_uuid, training_id),
    FOREIGN KEY (training_id) REFERENCES trainings(id) ON DELETE CASCADE
);