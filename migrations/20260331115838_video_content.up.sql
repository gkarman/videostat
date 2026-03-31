CREATE TABLE video_content
(
    id           BIGSERIAL PRIMARY KEY,
    account_id   BIGINT    NOT NULL REFERENCES video_accounts (id) ON DELETE CASCADE,
    external_id  TEXT      NOT NULL,
    title        TEXT,
    published_at TIMESTAMP,
    duration_sec INT,
    created_at   TIMESTAMP NOT NULL DEFAULT NOW(),

    UNIQUE (account_id, external_id)
);

CREATE INDEX idx_video_content_account_id ON video_content (account_id);
CREATE INDEX idx_video_content_published_at ON video_content (published_at);
CREATE INDEX idx_video_content_external_id ON video_content (external_id);