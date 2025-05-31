create table if not exists account (
    id text primary key,
    username text not null unique,
    password text not null,
    created_at timestamptz default current_timestamp,
    updated_at timestamptz default current_timestamp
);

create table if not exists user_info (
    id text primary key,
    account_id text not null unique,
    type text not null,
    email text not null unique,
    fulll_name text not null,
    avatar text,
    created_at timestamptz default current_timestamp,
    updated_at timestamptz default current_timestamp,
    foreign key (account_id) references account(id)
);

create index if not exists idx_user_info_account_id on user_info(account_id);

create table if not exists session (
    session_token text primary key,
    account_id text not null,
    created_at timestamptz default current_timestamp,
    updated_at timestamptz default current_timestamp,
    expired_at timestamptz,
    is_active boolean default true,
    user_agent text,
    ip_address text,
    foreign key (account_id) references account(id)
);

create index if not exists idx_session_account_id on session(account_id);