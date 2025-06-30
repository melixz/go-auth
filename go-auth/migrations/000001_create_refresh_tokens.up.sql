CREATE TABLE IF NOT EXISTS refresh_tokens (
                                              id SERIAL PRIMARY KEY,
                                              user_id UUID NOT NULL,
                                              bcrypt_hash TEXT NOT NULL,
                                              user_agent TEXT NOT NULL,
                                              ip INET NOT NULL,
                                              created_at TIMESTAMPTZ DEFAULT now(),
    revoked BOOLEAN DEFAULT FALSE,
    used BOOLEAN DEFAULT FALSE
    );
CREATE INDEX IF NOT EXISTS idx_refresh_tokens_user_id ON refresh_tokens(user_id);