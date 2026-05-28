ALTER TABLE video_broll_segments
    ADD COLUMN generation_status TEXT NOT NULL DEFAULT 'pending',
    ADD COLUMN generation_external_id TEXT,
    ADD COLUMN generation_error TEXT;
