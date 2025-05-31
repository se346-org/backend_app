ALTER TABLE conversation ADD COLUMN last_message_id text;
CREATE INDEX idx_conversation_last_message_id ON conversation (last_message_id);
