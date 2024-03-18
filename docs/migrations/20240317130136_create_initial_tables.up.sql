CREATE SEQUENCE roles_role_id_seq AS INTEGER;
ALTER SEQUENCE roles_role_id_seq OWNER TO filmoteca;

CREATE SEQUENCE actors_actor_id_seq AS INTEGER;
ALTER SEQUENCE actors_actor_id_seq OWNER TO filmoteca;

CREATE SEQUENCE movies_movie_id_seq AS INTEGER;
ALTER SEQUENCE movies_movie_id_seq OWNER TO filmoteca;

CREATE SEQUENCE users_user_id_seq AS INTEGER;
ALTER SEQUENCE users_user_id_seq OWNER TO filmoteca;

CREATE TABLE roles (
    role_id SERIAL PRIMARY KEY,
    role_name VARCHAR(50) NOT NULL UNIQUE
);
ALTER SEQUENCE roles_role_id_seq OWNED BY roles.role_id;

CREATE TABLE actors (
    actor_id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    gender VARCHAR(50),
    birthdate DATE
);
ALTER SEQUENCE actors_actor_id_seq OWNED BY actors.actor_id;

CREATE TABLE movies (
    movie_id SERIAL PRIMARY KEY,
    title VARCHAR(150) NOT NULL,
    description TEXT,
    release_date DATE,
    rating NUMERIC(2, 1) CHECK (rating >= 0 AND rating <= 10)
);
ALTER SEQUENCE movies_movie_id_seq OWNED BY movies.movie_id;

CREATE TABLE movies_actors (
    movie_id INTEGER NOT NULL REFERENCES movies ON DELETE CASCADE,
    actor_id INTEGER NOT NULL REFERENCES actors ON DELETE CASCADE,
    PRIMARY KEY (movie_id, actor_id)
);

CREATE TABLE users (
    user_id SERIAL PRIMARY KEY,
    username VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    role_id INTEGER NOT NULL REFERENCES roles
);

INSERT INTO roles (role_name) VALUES ('Admin');
INSERT INTO roles (role_name) VALUES ('User');

ALTER SEQUENCE users_user_id_seq OWNED BY users.user_id;
