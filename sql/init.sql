CREATE DATABASE smart_home;

\c smart_home;

CREATE TABLE smart_plants
(
    id         SERIAL PRIMARY KEY,
    moist      FLOAT NOT NULL,
    temperature FLOAT NOT NULL,
    humidity   FLOAT NOT NULL,
    date       TEXT  NOT NULL
);