ALTER TABLE video_prompts RENAME COLUMN prompt TO brief;
ALTER TABLE video_prompts ADD COLUMN script TEXT NOT NULL DEFAULT '';
