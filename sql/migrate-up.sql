CREATE TABLE IF NOT EXISTS users
(
    user_id  bigint PRIMARY KEY,

    username varchar(128) UNIQUE NOT NULL,

    token    char(40) UNIQUE     NOT NULL
);

CREATE TABLE IF NOT EXISTS sessions
(
    token   char(128) PRIMARY KEY,

    user_id bigint NOT NULL,

    created date,

    CONSTRAINT fk
        FOREIGN KEY (user_id) REFERENCES users (user_id)
            ON DELETE CASCADE
            ON UPDATE CASCADE

);
