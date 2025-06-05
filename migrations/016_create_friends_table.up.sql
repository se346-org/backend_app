CREATE TABLE IF NOT EXISTS friends (
    id VARCHAR(36) PRIMARY KEY,
    user_id VARCHAR(36) NOT NULL,
    friend_id VARCHAR(36) NOT NULL,
    status VARCHAR(20) NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    FOREIGN KEY (user_id) REFERENCES user_info(id),
    FOREIGN KEY (friend_id) REFERENCES user_info(id),
    UNIQUE(user_id, friend_id)
); 