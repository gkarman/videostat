ALTER TABLE videos
    DROP COLUMN status,
    DROP COLUMN error_stage,
    DROP COLUMN error_message;

DROP INDEX IF EXISTS idx_videos_status;