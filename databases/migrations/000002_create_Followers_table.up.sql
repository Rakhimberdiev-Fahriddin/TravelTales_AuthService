CREATE TABLE IF NOT EXISTS Followers (
    follower_id UUID REFERENCES users(id),
    following_id UUID REFERENCES users(id),
    followed_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (follower_id, following_id)
);

