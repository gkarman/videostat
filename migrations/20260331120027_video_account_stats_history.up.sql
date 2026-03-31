CREATE TABLE video_account_stats_history
(
    id           BIGSERIAL PRIMARY KEY,
    account_id   BIGINT    NOT NULL REFERENCES video_accounts (id) ON DELETE CASCADE,

    followers    BIGINT,
    total_views  BIGINT,

    collected_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_video_account_stats_account_id
    ON video_account_stats_history (account_id);

CREATE INDEX idx_video_account_stats_time
    ON video_account_stats_history (account_id, collected_at);