CREATE TABLE video_watchers
(
    id         UUID PRIMARY KEY,
    video_id   UUID   NOT NULL REFERENCES videos (id) ON DELETE CASCADE,
    chat_id    BIGINT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_video_watchers_video_id ON video_watchers (video_id);
