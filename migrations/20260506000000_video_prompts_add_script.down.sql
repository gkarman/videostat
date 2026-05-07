ALTER TABLE video_prompts DROP COLUMN script;
ALTER TABLE video_prompts RENAME COLUMN brief TO prompt;
