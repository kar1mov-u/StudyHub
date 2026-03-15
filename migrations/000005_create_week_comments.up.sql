CREATE TABLE IF NOT EXISTS week_comments(
    id UUID PRIMARY KEY,
    week_id UUID REFERENCES weeks(id),
    user_id UUID REFERENCES users(id),
    reply UUID REFERENCES week_comments(id),
    content TEXT NOT NULL,
    upvote INT, 
    downvote INT,
    created_at TIMESTAMPTZ DEFAULT NOW()
);