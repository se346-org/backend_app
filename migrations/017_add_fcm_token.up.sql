create table if not exists fcm_token (
    id text primary key,
    user_id text not null,
    token text not null,
    created_at timestamptz default current_timestamp,
    updated_at timestamptz default current_timestamp
);

create index if not exists idx_user_id_fcm_token on fcm_token(user_id);