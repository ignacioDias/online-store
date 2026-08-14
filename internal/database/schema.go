package database

var createSessionsTable = `
CREATE TABLE IF NOT EXISTS sessions (
    id TEXT PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    expires_at TIMESTAMPTZ NOT NULL
);`

var createUsersTable = `CREATE TABLE IF NOT EXISTS users (
    id                UUID PRIMARY KEY,
    username          VARCHAR(50) UNIQUE NOT NULL,
    name              VARCHAR(255) NOT NULL,
    email             VARCHAR(255) UNIQUE NOT NULL,
    hashed_password   TEXT NOT NULL,
    profile_picture   TEXT,                         -- URL o path al archivo
    is_active         BOOLEAN NOT NULL DEFAULT TRUE, -- para soft-disable en vez de borrar
    email_verified    BOOLEAN NOT NULL DEFAULT FALSE,
    created_at        TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at        TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_username ON users(username);
`

var createNotificationsTable = `CREATE TABLE IF NOT EXISTS notifications (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    type VARCHAR(50) NOT NULL,
    title VARCHAR(255) NOT NULL,
    message TEXT NOT NULL,
    metadata JSONB,                   -- data extra según el tipo (ej: order_id)
    is_seen BOOLEAN DEFAULT FALSE,
    seen_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT NOW()
);
CREATE INDEX idx_notifications_user_unseen ON notifications(user_id, is_seen);
`

var createFavoritesTable = `CREATE TABLE IF NOT EXISTS favorites (
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    product_id UUID NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    PRIMARY KEY (user_id, product_id)
);
`

var createReviewsTable = `CREATE TABLE IF NOT EXISTS reviews (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    product_id UUID NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    score SMALLINT NOT NULL CHECK (score BETWEEN 1 AND 10),
    comment TEXT,
    image_url TEXT,
    likes INT NOT NULL DEFAULT 0,
    dislikes INT NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE (user_id, product_id)
);
`
