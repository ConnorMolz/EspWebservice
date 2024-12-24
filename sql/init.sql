CREATE DATABASE smart_home;

\c smart_home;

CREATE TABLE smart_plants
(
    id         SERIAL PRIMARY KEY,
    sensor_id INT   NOT NULL,
    moist      FLOAT NOT NULL,
    temperature FLOAT NOT NULL,
    humidity   FLOAT NOT NULL,
    date       TEXT  NOT NULL
);