create table if not exists conversation (
    id text primary key,
    user_id text not null,
    title text,
    created_at timestamptz default current_timestamp,
    updated_at timestamptz default current_timestamp,
    deleted_at timestamptz,
    avatar text,
    foreign key (user_id) references user_info(id)
);

create index if not exists idx_updated_at_conversation on conversation(updated_at);
create index if not exists idx_user_id_conversation on conversation(user_id);

create table if not exists conversation_member (
    id text primary key,
    conversation_id text not null,
    user_id text not null,
    role text not null,
    created_at timestamptz default current_timestamp,
    updated_at timestamptz default current_timestamp,
    deleted_at timestamptz,
    foreign key (conversation_id) references conversation(id),
    foreign key (user_id) references user_info(id)
);

create index if not exists idx_conversation_id_conversation_member on conversation_member(conversation_id);
create index if not exists idx_user_id_conversation_member on conversation_member(user_id);

create table if not exists message (
    id text primary key,
    conversation_id text not null,
    user_id text not null,
    type text not null,
    body text,
    created_at timestamptz default current_timestamp,
    updated_at timestamptz default current_timestamp,
    deleted_at timestamptz,
    reply_to text,
    foreign key (conversation_id) references conversation(id),
    foreign key (user_id) references user_info(id)
);

create index if not exists idx_conversation_id_message on message(conversation_id);
create index if not exists idx_user_id_message on message(user_id);
create index if not exists idx_reply_to_message on message(reply_to);
create index if not exists idx_id_created_at_message on message(id, created_at);

create table if not exists contact (
    id text primary key,
    user_id text not null,
    friend_id text not null,
    created_at timestamptz default current_timestamp,
    updated_at timestamptz default current_timestamp,
    foreign key (user_id) references user_info(id),
    foreign key (friend_id) references user_info(id)
);

create index if not exists idx_user_id_contact on contact(user_id);
create index if not exists idx_friend_id_contact on contact(friend_id);