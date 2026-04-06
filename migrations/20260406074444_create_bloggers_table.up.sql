CREATE TABLE bloggers
(
    id          UUID PRIMARY KEY,
    platform_id SMALLINT NOT NULL,
    url         TEXT     NOT NULL,

    CONSTRAINT fk_bloggers_platform
        FOREIGN KEY (platform_id)
            REFERENCES platforms (id)
            ON DELETE RESTRICT
            ON UPDATE CASCADE
);

