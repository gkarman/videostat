CREATE TABLE video_content_stats_history
(
    id           BIGSERIAL PRIMARY KEY,
    content_id   BIGINT    NOT NULL REFERENCES video_content (id) ON DELETE CASCADE,

    views        BIGINT,
    likes        BIGINT,
    comments     BIGINT,
    shares       BIGINT,

    collected_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_video_content_stats_content_id
    ON video_content_stats_history (content_id);

CREATE INDEX idx_video_content_stats_time
    ON video_content_stats_history (content_id, collected_at);