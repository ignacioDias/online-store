package database

var createSessionsTable = `
CREATE TABLE IF NOT EXISTS sessions (
    id TEXT PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
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
     is_seen BOOLEAN NOT NULL DEFAULT FALSE,
     seen_at TIMESTAMPTZ,
     created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_notifications_user_unseen ON notifications(user_id, is_seen);
`

var createFavoritesTable = `CREATE TABLE IF NOT EXISTS favorites (
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    product_id UUID NOT NULL REFERENCES products(id) ON DELETE CASCADE,
     created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
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
     likes INT NOT NULL DEFAULT 0 CHECK (likes >= 0),
     dislikes INT NOT NULL DEFAULT 0 CHECK (dislikes >= 0),
     created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (user_id, product_id)
);
`
var createShoppingCartTable = `CREATE TABLE IF NOT EXISTS shopping_cart (
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    product_id UUID NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    quantity INT NOT NULL DEFAULT 1 CHECK (quantity > 0),
     added_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
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
     created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
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
     amount NUMERIC(10, 2) NOT NULL CHECK (amount >= 0),
    currency VARCHAR(3) NOT NULL DEFAULT 'ARS',
    payment_method VARCHAR(50),                 -- 'credit_card', 'debit_card', 'cash', etc. (te lo da el proveedor)
    metadata JSONB,                             -- data extra que te devuelva el proveedor
     created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
     updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (provider, provider_payment_id)
);
`
var createCategoriesTable = `
CREATE TABLE IF NOT EXISTS categories (
    id UUID PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    slug VARCHAR(100) NOT NULL UNIQUE,     -- "zapatillas-running", para URLs amigables
    parent_id UUID REFERENCES categories(id) ON DELETE CASCADE,  -- NULL = categoría raíz
     created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
`

var createProductsTable = `
CREATE TABLE IF NOT EXISTS products (
    id UUID PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
     base_price NUMERIC(10, 2) NOT NULL CHECK (base_price >= 0),
     category_id UUID REFERENCES categories(id),
     created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS product_variants (
    id UUID PRIMARY KEY,
    product_id UUID NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    sku VARCHAR(100) NOT NULL UNIQUE,      -- código único de esta variante puntual
     price NUMERIC(10, 2) CHECK (price IS NULL OR price >= 0), -- NULL = usa base_price del producto; si tiene valor, lo sobreescribe
    stock INT NOT NULL DEFAULT 0 CHECK (stock >= 0),
    attributes JSONB NOT NULL,             -- {"talle": "42", "color": "negro"}
    image_url TEXT,
     created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);`

var createCouponsTable = `
CREATE TABLE IF NOT EXISTS coupons (
    id UUID PRIMARY KEY,
    code VARCHAR(50) NOT NULL UNIQUE,           -- "SAVE20", lo que tipea el usuario
    discount_type VARCHAR(20) NOT NULL,          -- 'percentage' o 'fixed_amount'
     discount_value NUMERIC(10, 2) NOT NULL CHECK (discount_value >= 0), -- 20 (=20%) o 500 (=$500)
     min_purchase_amount NUMERIC(10, 2) CHECK (min_purchase_amount IS NULL OR min_purchase_amount >= 0), -- compra mínima para aplicar, opcional
     max_discount_amount NUMERIC(10, 2) CHECK (max_discount_amount IS NULL OR max_discount_amount >= 0), -- tope del descuento (útil con %), opcional
     usage_limit INT CHECK (usage_limit IS NULL OR usage_limit > 0),     -- cuántas veces se puede usar en total, NULL = ilimitado
     usage_count INT NOT NULL DEFAULT 0 CHECK (usage_count >= 0),       -- cuántas veces ya se usó
     usage_limit_per_user INT NOT NULL DEFAULT 1 CHECK (usage_limit_per_user > 0), -- veces que un mismo usuario puede usarlo
    starts_at TIMESTAMPTZ NOT NULL,
    expires_at TIMESTAMPTZ NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
     created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
     CHECK (discount_type IN ('percentage', 'fixed_amount')),
     CHECK (discount_type <> 'percentage' OR discount_value <= 100),
     CHECK (expires_at > starts_at)
);`

var createInventoryTable = `
CREATE TABLE IF NOT EXISTS inventory (
    id UUID PRIMARY KEY,
    variant_id UUID NOT NULL REFERENCES product_variants(id) ON DELETE CASCADE,
    quantity INT NOT NULL DEFAULT 0 CHECK (quantity >= 0),
    reserved_quantity INT NOT NULL DEFAULT 0 CHECK (reserved_quantity >= 0),
     updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    UNIQUE (variant_id)
);

CREATE TABLE IF NOT EXISTS inventory_movements (
    id UUID PRIMARY KEY,
    variant_id UUID NOT NULL REFERENCES product_variants(id) ON DELETE CASCADE,
    type VARCHAR(20) NOT NULL,          -- 'sale', 'restock', 'return', 'adjustment'
    quantity INT NOT NULL,              -- positivo (entra) o negativo (sale)
    order_id INT REFERENCES orders(id), -- si el movimiento viene de una venta
    reason TEXT,                        -- para ajustes manuales, ej: "producto dañado"
     created_at TIMESTAMP NOT NULL DEFAULT NOW()
);
`
