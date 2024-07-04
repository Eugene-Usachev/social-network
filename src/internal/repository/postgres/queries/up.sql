CREATE TABLE IF NOT EXISTS users (
    id               BIGSERIAL       PRIMARY KEY,
    name             VARCHAR(255)    NOT NULL,
    second_name      VARCHAR(255)    NOT NULL,
    email            VARCHAR(255)    NOT NULL UNIQUE,
    password         VARCHAR(255)    NOT NULL,
    friends          BIGINT []       default ARRAY[]::BIGINT[],
    subscribers      BIGINT []       default ARRAY[]::BIGINT[],
    avatar           VARCHAR(64)     default '',
    birthday         VARCHAR(32)     default '',
    -- 0 - not specified, 1 - male, 2 - female
    gender           SMALLINT        default 0,
    -- info is a jsonb field because info can be changed in the future and it's easier to keep it without a schema
    info             jsonb           default '{}',
    registered_at    TIMESTAMP       NOT NULL DEFAULT NOW()
)