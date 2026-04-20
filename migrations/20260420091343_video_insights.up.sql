CREATE TABLE video_insights
(
    video_id   UUID PRIMARY KEY REFERENCES videos (id) ON DELETE CASCADE,

    transcript JSONB     NOT NULL,
    hooks      JSONB     NOT NULL,
    prosody    JSONB     NOT NULL,
    structure  JSONB     NOT NULL,

    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);