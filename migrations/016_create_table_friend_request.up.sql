create table if not exists friend_request (
    id text primary key,
    from_user_id text not null,
    to_user_id text not null,
    status text not null,
    created_at timestamptz default current_timestamp,
    updated_at timestamptz default current_timestamp
);

create index if not exists idx_from_user_id_friend_request on friend_request(from_user_id);
create index if not exists idx_to_user_id_friend_request on friend_request(to_user_id);
create index if not exists idx_status_friend_request on friend_request("status");