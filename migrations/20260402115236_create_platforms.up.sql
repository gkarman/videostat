CREATE TABLE platforms
(
    id   SMALLSERIAL PRIMARY KEY,
    name TEXT NOT NULL UNIQUE
);

INSERT INTO platforms (name)
VALUES ('youtube'),
       ('tiktok'),
       ('instagram');