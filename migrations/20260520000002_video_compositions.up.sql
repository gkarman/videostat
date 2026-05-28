CREATE TABLE video_compositions (
    id UUID PRIMARY KEY,
    video_id UUID NOT NULL REFERENCES videos(id) ON DELETE CASCADE,
    status TEXT NOT NULL DEFAULT 'processing',
    external_id TEXT,
    result_url TEXT,
    error TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_video_compositions_video_id ON video_compositions (video_id);
