CREATE TABLE actor (
  id SERIAL PRIMARY KEY,
  name TEXT NOT NULL,
  sex TEXT NOT NULL,
  birthday DATE NOT NULL
);

CREATE TABLE film (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    description TEXT,
    release_date DATE NOT NULL,
    rating int NOT NULL
);

CREATE TABLE actor_film (
    id SERIAL PRIMARY KEY,
    actor_id BIGINT NOT NULL REFERENCES actor(id) ON UPDATE CASCADE ON DELETE CASCADE,
    film_id BIGINT NOT NULL REFERENCES film(id) ON UPDATE CASCADE ON DELETE CASCADE
);

CREATE TABLE role (
    id SERIAL PRIMARY KEY,
    role TEXT NOT NULL
);

CREATE TABLE _user (
    id SERIAL PRIMARY KEY,
    username TEXT UNIQUE NOT NULL,
    password TEXT NOT NULL,
    role_id INT NOT NULL REFERENCES role(id) ON UPDATE CASCADE ON DELETE CASCADE
);

INSERT INTO role (role) VALUES
('ADMIN'),
('USER');

INSERT INTO _user (username, password, role_id) VALUES
('admin', '$2a$10$LacZ2mV3TWRotYoceAZAUe/hQSaJ/AfitUA4YnIXtiruzuUgiuCOm', '1'),
('user', '$2a$10$mMMTQDCwRFgUAs.m40rHsOL7mc4TvrMOZGutJsGBD42rEVViQ8I8u', '2');