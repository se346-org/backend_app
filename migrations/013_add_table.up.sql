create table if not exists seen_message (
    id text primary key,
    message_id text not null,
    user_id text not null,
    conversation_id text not null,
    created_at timestamptz default current_timestamp,
    updated_at timestamptz default current_timestamp,
    foreign key (message_id) references message(id),
    foreign key (user_id) references user_info(id),
    foreign key (conversation_id) references conversation(id),
    unique (user_id, conversation_id)
);

create index if not exists idx_seen_message_user_id on seen_message(user_id);
create index if not exists idx_seen_message_conversation_id on seen_message(conversation_id);