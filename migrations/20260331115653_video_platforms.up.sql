CREATE TABLE video_platforms
(
    id   SMALLSERIAL PRIMARY KEY,
    name TEXT NOT NULL UNIQUE
);

INSERT INTO video_platforms (name)
VALUES ('youtube'),
       ('tiktok'),
       ('facebook');