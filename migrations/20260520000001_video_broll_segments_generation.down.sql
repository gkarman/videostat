ALTER TABLE video_broll_segments
    DROP COLUMN IF EXISTS generation_status,
    DROP COLUMN IF EXISTS generation_external_id,
    DROP COLUMN IF EXISTS generation_error;
