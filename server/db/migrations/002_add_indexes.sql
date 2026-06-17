-- 002_add_indexes.sql: Performance indexes for frequent queries

CREATE INDEX IF NOT EXISTS idx_cart_items_user_id ON cart_items(user_id);

CREATE INDEX IF NOT EXISTS idx_sessions_user_id ON sessions(user_id);

CREATE INDEX IF NOT EXISTS idx_sessions_expires_at ON sessions(expires_at);
