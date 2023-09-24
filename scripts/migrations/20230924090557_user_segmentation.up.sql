-- USER'S SESSION --
CREATE TYPE session_type AS (
    refresh_token   VARCHAR(255),
    expires_at TIMESTAMP
);

-- USERS --
CREATE TABLE users (
    id              SERIAL PRIMARY KEY,
    login           VARCHAR(255) UNIQUE NOT NULL,
    email           VARCHAR(255) UNIQUE NOT NULL,
    password        VARCHAR(255) NOT NULL,
    session         session_type,
    registered_at   TIMESTAMP NOT NULL DEFAULT NOW()
);