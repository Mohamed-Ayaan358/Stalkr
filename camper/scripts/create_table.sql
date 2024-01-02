CREATE TABLE websites (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255),
    url VARCHAR(255),
    hash VARCHAR(255),
    time VARCHAR(255),
    changed BOOLEAN
);