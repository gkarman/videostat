CREATE TABLE videos
(
    id           UUID PRIMARY KEY,
    blogger_id   UUID NOT NULL,
    external_id  TEXT NOT NULL,
    url          TEXT NOT NULL,
    title        TEXT,
    views        BIGINT,
    likes        BIGINT,
    comments     BIGINT,
    published_at TIMESTAMP,
    created_at   TIMESTAMP NOT NULL,

    CONSTRAINT fk_videos_blogger
        FOREIGN KEY (blogger_id)
            REFERENCES bloggers (id)
            ON DELETE CASCADE,

    CONSTRAINT uniq_video_per_blogger
        UNIQUE (blogger_id, external_id)
);