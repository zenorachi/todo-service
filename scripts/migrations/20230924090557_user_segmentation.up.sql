-- USER'S SESSION --
CREATE TYPE session_type AS (
    refresh_token   VARCHAR(255),
    expires_at      TIMESTAMP
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

-- TASK TYPE --
CREATE TYPE task_type AS ENUM (
    'done', 'not done'
);

-- TO-DO LIST --
CREATE TABLE agenda (
    id              SERIAL PRIMARY KEY,
    user_id         INT NOT NULL,
    title           VARCHAR(255) NOT NULL,
    description     VARCHAR DEFAULT NULL,
    date            TIMESTAMP NOT NULL,
    status          task_type DEFAULT 'not done',
    FOREIGN KEY (user_id) REFERENCES users (id)
);