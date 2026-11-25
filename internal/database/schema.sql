CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    username TEXT NOT NULL COLLATE NOCASE UNIQUE,
    email TEXT NOT NULL COLLATE NOCASE UNIQUE,
    password_hash TEXT NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS sessions (
    id TEXT PRIMARY KEY,
    user_id INTEGER NOT NULL,
    expires_at DATETIME NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,

    FOREIGN KEY (user_id)
        REFERENCES users(id)
        ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS categories (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL COLLATE NOCASE UNIQUE,
    slug TEXT NOT NULL COLLATE NOCASE UNIQUE,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS posts (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    title TEXT NOT NULL
        CHECK (length(trim(title)) BETWEEN 1 AND 150),
    content TEXT NOT NULL
        CHECK (length(trim(content)) BETWEEN 1 AND 10000),
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,

    FOREIGN KEY (user_id)
        REFERENCES users(id)
        ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS post_categories (
    post_id INTEGER NOT NULL,
    category_id INTEGER NOT NULL,

    PRIMARY KEY (post_id, category_id),

    FOREIGN KEY (post_id)
        REFERENCES posts(id)
        ON DELETE CASCADE,
    
    FOREIGN KEY (category_id)
        REFERENCES categories(id)
        ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS comments (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    post_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL,
    content TEXT NOT NULL
        CHECK (length(trim(content)) BETWEEN 1 AND 800),
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,

    FOREIGN KEY (post_id)
        REFERENCES posts(id)
        ON DELETE CASCADE,
    
    FOREIGN KEY (user_id)
        REFERENCES users(id)
        ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS post_reactions (
    user_id INTEGER NOT NULL,
    post_id INTEGER NOT NULL,
    value INTEGER NOT NULL CHECK (value IN (-1, 1)),
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,

    PRIMARY KEY (user_id, post_id),

    FOREIGN KEY (user_id)
        REFERENCES users(id)
        ON DELETE CASCADE,

    FOREIGN KEY (post_id)
        REFERENCES posts(id)
        ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS comment_reactions (
    user_id INTEGER NOT NULL,
    comment_id INTEGER NOT NULL,
    value INTEGER NOT NULL CHECK (value IN (-1, 1)),
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,

    PRIMARY KEY (user_id, comment_id),

    FOREIGN KEY (user_id)
        REFERENCES users(id)
        ON DELETE CASCADE,

    FOREIGN KEY (comment_id)
        REFERENCES comments(id)
        ON DELETE CASCADE
);

-- Existing databases may contain duplicate sessions. Keep the latest created
-- session (rowid breaks ties), then enforce one session per user.
DELETE FROM sessions WHERE rowid NOT IN (
    SELECT rowid FROM sessions AS s
    WHERE s.rowid = (
        SELECT newer.rowid FROM sessions AS newer
        WHERE newer.user_id = s.user_id
        ORDER BY julianday(newer.created_at) DESC, newer.rowid DESC LIMIT 1
    )
);
CREATE UNIQUE INDEX IF NOT EXISTS idx_sessions_unique_user_id ON sessions(user_id);
DROP INDEX IF EXISTS idx_sessions_user_id;

CREATE INDEX IF NOT EXISTS idx_sessions_expires_at
    ON sessions(expires_at);

CREATE INDEX IF NOT EXISTS idx_posts_user_id
    ON posts(user_id);

CREATE INDEX IF NOT EXISTS idx_posts_created_at
    ON posts(created_at DESC);

CREATE INDEX IF NOT EXISTS idx_post_categories_category_id
    ON post_categories(category_id);

CREATE INDEX IF NOT EXISTS idx_comments_post_id
    ON comments(post_id);

CREATE INDEX IF NOT EXISTS idx_comments_user_id
    ON comments(user_id);

CREATE INDEX IF NOT EXISTS idx_post_reactions_post_id
    ON post_reactions(post_id);

CREATE INDEX IF NOT EXISTS idx_post_reactions_user_value
    ON post_reactions(user_id, value);

CREATE INDEX IF NOT EXISTS idx_comment_reactions_comment_id
    ON comment_reactions(comment_id);
