create table if not exists user_online (
    id text PRIMARY KEY,
    user_id text NOT NULL,
    connection_id text NOT NULL,
    created_at timestamptz NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_user_id FOREIGN KEY (user_id) REFERENCES user_info (id)
);
