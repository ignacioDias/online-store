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
    score SMALLINT NOT NULL CHECK (score BETWEEN 1 AND 5),
    comment TEXT,
    image_url TEXT,
    likes INT NOT NULL DEFAULT 0,
    dislikes INT NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE (user_id, product_id)
);
`
var createShoppingCartTable = `CREATE TABLE IF NOT EXISTS shopping_cart (
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    product_id UUID NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    quantity INT NOT NULL DEFAULT 1 CHECK (quantity > 0),
    added_at TIMESTAMPTZ DEFAULT NOW(),
    PRIMARY KEY (user_id, product_id)
);
`
var createAddressesTable = `CREATE TABLE IF NOT EXISTS addresses (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    street VARCHAR(255) NOT NULL,        -- "Av. Colón 1234"
    apartment VARCHAR(50),               -- "Piso 3, depto B" (opcional)
    city VARCHAR(100) NOT NULL,
    state VARCHAR(100) NOT NULL,         -- provincia
    postal_code VARCHAR(20) NOT NULL,
    country VARCHAR(100) NOT NULL,
    is_default BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ DEFAULT NOW()
);
CREATE UNIQUE INDEX IF NOT EXISTS idx_addresses_one_default_per_user
    ON addresses(user_id)
    WHERE is_default = TRUE;
`

var createPaymentsTable = `CREATE TABLE IF NOT EXISTS payments (
    id UUID PRIMARY KEY,
    order_id UUID NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    provider VARCHAR(50) NOT NULL,              -- 'mercadopago', 'stripe'
    provider_payment_id VARCHAR(255) NOT NULL,  -- el ID que te devuelve el proveedor
    status VARCHAR(50) NOT NULL,                -- 'pending', 'approved', 'rejected', 'refunded'
    amount NUMERIC(10, 2) NOT NULL,
    currency VARCHAR(3) NOT NULL DEFAULT 'ARS',
    payment_method VARCHAR(50),                 -- 'credit_card', 'debit_card', 'cash', etc. (te lo da el proveedor)
    metadata JSONB,                             -- data extra que te devuelva el proveedor
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE (provider, provider_payment_id)
);
`
var createProductsTable = `
CREATE TABLE IF NOT EXISTS products (
    id UUID PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    base_price NUMERIC(10, 2) NOT NULL,
    category_id UUID REFERENCES categories(id),
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS product_variants (
    id UUID PRIMARY KEY,
    product_id UUID NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    sku VARCHAR(100) NOT NULL UNIQUE,      -- código único de esta variante puntual
    price NUMERIC(10, 2),                  -- NULL = usa base_price del producto; si tiene valor, lo sobreescribe
    stock INT NOT NULL DEFAULT 0 CHECK (stock >= 0),
    attributes JSONB NOT NULL,             -- {"talle": "42", "color": "negro"}
    image_url TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW()
);`

var createCategoriesTable = `
CREATE TABLE IF NOT EXISTS categories (
    id UUID PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    slug VARCHAR(100) NOT NULL UNIQUE,     -- "zapatillas-running", para URLs amigables
    parent_id UUID REFERENCES categories(id) ON DELETE CASCADE,  -- NULL = categoría raíz
    created_at TIMESTAMPTZ DEFAULT NOW()
);
`
