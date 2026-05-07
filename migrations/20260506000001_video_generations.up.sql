CREATE TABLE video_generations
(
    id          UUID PRIMARY KEY,
    video_id    UUID        NOT NULL REFERENCES videos (id) ON DELETE CASCADE,
    platform    VARCHAR     NOT NULL,
    external_id VARCHAR     NOT NULL,
    status      VARCHAR     NOT NULL DEFAULT 'pending',
    s3_url      TEXT        NOT NULL DEFAULT '',
    created_at  TIMESTAMP   NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMP   NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_video_generations_video_id ON video_generations (video_id);
CREATE INDEX idx_video_generations_status   ON video_generations (status);
