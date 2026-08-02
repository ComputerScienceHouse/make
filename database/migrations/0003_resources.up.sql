CREATE TABLE resources (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    description TEXT NOT NULL,
    url TEXT NOT NULL
);

CREATE TABLE area_resources (
    area_id INTEGER NOT NULL,
    resource_id INTEGER NOT NULL,
    PRIMARY KEY (area_id, resource_id),
    FOREIGN KEY (area_id) REFERENCES areas(id) ON DELETE CASCADE,
    FOREIGN KEY (resource_id) REFERENCES resources(id) ON DELETE CASCADE
);