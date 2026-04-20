CREATE TABLE video_analysis
(
    id          UUID PRIMARY KEY,
    video_id    UUID      NOT NULL REFERENCES videos (id) ON DELETE CASCADE,
    provider    TEXT      NOT NULL, -- assemblyai / whisper / deepgram
    raw_payload JSONB     NOT NULL,
    created_at  TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_video_analysis_video_id ON video_analysis (video_id);
CREATE INDEX idx_video_analysis_provider ON video_analysis (provider);