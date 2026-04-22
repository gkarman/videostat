ALTER TABLE videos
    ADD COLUMN status TEXT NOT NULL DEFAULT 'created',
    ADD COLUMN error_stage TEXT,
    ADD COLUMN error_message TEXT;

CREATE INDEX idx_videos_status ON videos (status);