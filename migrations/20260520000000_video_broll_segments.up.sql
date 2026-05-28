CREATE TABLE video_broll_segments (
    id UUID PRIMARY KEY,
    video_id UUID NOT NULL REFERENCES videos(id) ON DELETE CASCADE,
    position INT NOT NULL,
    start_ms INT NOT NULL,
    end_ms INT NOT NULL,
    text TEXT NOT NULL,
    broll_prompt TEXT NOT NULL,
    broll_url TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_video_broll_segments_video_id ON video_broll_segments (video_id);
