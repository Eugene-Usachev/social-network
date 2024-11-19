CREATE TABLE IF NOT EXISTS posts
(
    id         UUID PRIMARY KEY,
    owner_id   BIGINT,
    text       TEXT,
    survey     TEXT,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    files      LIST<TEXT>,
    likes      BIGINT,
    dislikes   BIGINT,
    PRIMARY KEY ((owner_id), created_at)
) WITH CLUSTERING ORDER BY (created_at DESC);

CREATE INDEX IF NOT EXISTS posts_owner_id_idx ON posts (owner_id);

CREATE TABLE IF NOT EXISTS posts_rates
(
    post_id UUID,
    user_id BIGINT,
    is_like BOOLEAN,
    PRIMARY KEY (post_id, user_id)
)