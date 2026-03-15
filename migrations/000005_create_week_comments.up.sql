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

CREATE TABLE IF NOT EXISTS comment_votes(
    comment_id UUID REFERENCES week_comments(id),
    user_id UUID REFERENCES users(id),
    is_upvote BOOLEAN NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    UNIQUE(comment_id, user_id)
);