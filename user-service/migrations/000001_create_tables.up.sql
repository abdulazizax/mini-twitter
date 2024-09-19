-- Users table
CREATE TABLE Users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    phone_number VARCHAR(50) UNIQUE DEFAULT NULL,
    bio TEXT DEFAULT '',
    username VARCHAR(50) UNIQUE NOT NULL,
    profile_picture TEXT DEFAULT 'https://cdn.pixabay.com/photo/2015/10/05/22/37/blank-profile-picture-973460_1280.png',
    role VARCHAR(50) DEFAULT 'user',
    refresh_token VARCHAR(1000) DEFAULT '',
    is_active BOOLEAN DEFAULT FALSE,
    first_name VARCHAR(50) DEFAULT '',
    last_name VARCHAR(50) DEFAULT '',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP DEFAULT NULL
);


-- Tweets table
CREATE TABLE Tweets (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    username VARCHAR(50) REFERENCES Users(username) ON DELETE CASCADE,
    tweet_serial INT NOT NULL,
    content TEXT NOT NULL,
    media VARCHAR[] DEFAULT '{}',
    comments_count INT DEFAULT 0,
    views_count INT DEFAULT 0,
    repost_count INT DEFAULT 0,
    shares_count INT DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP DEFAULT NULL,
    UNIQUE(username, tweet_serial)
);


-- Comments table
CREATE TABLE Comments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    username VARCHAR(50) REFERENCES Users(username) ON DELETE CASCADE,
    tweet_id UUID REFERENCES Tweets(id) ON DELETE CASCADE,
    comment_serial INT NOT NULL,
    content TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP DEFAULT NULL,
    UNIQUE(tweet_id, comment_serial)
);


-- Likes table
CREATE TABLE Likes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    username VARCHAR(50) REFERENCES Users(username) ON DELETE CASCADE,
    target_id UUID NOT NULL,
    target_type VARCHAR(10) CHECK (target_type IN ('tweet', 'comment')),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(username, target_id, target_type)
);


-- Followers table
CREATE TABLE Followers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    follower_id UUID REFERENCES Users(id) ON DELETE CASCADE,
    followed_id UUID REFERENCES Users(id) ON DELETE CASCADE,
    status BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(follower_id, followed_id)
);
