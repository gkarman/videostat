CREATE TABLE video_accounts
(
    id          BIGSERIAL PRIMARY KEY,
    platform_id SMALLINT  NOT NULL REFERENCES video_platforms (id),
    external_id TEXT      NOT NULL,
    title       TEXT      NOT NULL,
    url         TEXT,
    created_at  TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMP NOT NULL DEFAULT NOW(),

    UNIQUE (platform_id, external_id)
);

CREATE INDEX idx_video_accounts_platform ON video_accounts (platform_id);
CREATE INDEX idx_video_accounts_external_id ON video_accounts (external_id);