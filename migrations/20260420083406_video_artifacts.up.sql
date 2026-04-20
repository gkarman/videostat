CREATE TABLE video_artifacts
(
    id         UUID PRIMARY KEY,
    video_id   UUID      NOT NULL REFERENCES videos (id) ON DELETE CASCADE,

    type       TEXT      NOT NULL, -- audio | video | thumbnail
    url        TEXT      NOT NULL, -- S3

    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_video_artifacts_video_id ON video_artifacts (video_id);
CREATE INDEX idx_video_artifacts_video_type ON video_artifacts (video_id, type);