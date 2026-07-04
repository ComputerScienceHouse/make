CREATE TABLE training_questions (
    id SERIAL PRIMARY KEY,
    training_id INTEGER NOT NULL,
    label TEXT NOT NULL,
    type TEXT NOT NULL,
    answer TEXT NOT NULL,
    required BOOLEAN NOT NULL DEFAULT FALSE,

    FOREIGN KEY (training_id) REFERENCES trainings(id) ON DELETE CASCADE
);

CREATE TABLE question_options (
    id SERIAL PRIMARY KEY,
    question_id INTEGER NOT NULL,
    label TEXT NOT NULL,
    
    FOREIGN KEY (question_id) REFERENCES training_questions(id) ON DELETE CASCADE
);