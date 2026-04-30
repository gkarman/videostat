CREATE TABLE video_prompts
(
    id           UUID PRIMARY KEY,
    video_id     UUID        NOT NULL REFERENCES videos (id) ON DELETE CASCADE,
    llm_provider VARCHAR     NOT NULL,
    llm_model    VARCHAR     NOT NULL,
    prompt       TEXT        NOT NULL,
    created_at   TIMESTAMP   NOT NULL DEFAULT NOW()
);
