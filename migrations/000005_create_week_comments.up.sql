CREATE TABLE IF NOT EXISTS week_comments(
    id UUID PRIMARY KEY,
    week_id UUID REFERENCES weeks(id),
    user_id UUID REFERENCES users(id),
    reply UUID REFERENCES week_comments(id),
    content TEXT NOT NULL,
    upvote INT DEFAULT 0, 
    downvote INT DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);