-- 创建已验证订单表
CREATE TABLE IF NOT EXISTS verified_orders (
    id SERIAL PRIMARY KEY,
    order_id INTEGER NOT NULL,
    invoice VARCHAR(255) NOT NULL UNIQUE,
    user_id INTEGER,
    username VARCHAR(255),
    email VARCHAR(255) NOT NULL,
    total_price DECIMAL(10,2) NOT NULL,
    currency VARCHAR(10) NOT NULL,
    status VARCHAR(50) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

-- 创建已验证订单产品表
CREATE TABLE IF NOT EXISTS verified_order_products (
    id SERIAL PRIMARY KEY,
    verified_order_id BIGINT UNSIGNED NOT NULL,
    product_id INTEGER NOT NULL,
    quantity INTEGER NOT NULL,
    price DECIMAL(10,2) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (verified_order_id) REFERENCES verified_orders(id)
);

-- 创建索引
CREATE INDEX idx_verified_orders_invoice ON verified_orders(invoice);
CREATE INDEX idx_verified_orders_user_id ON verified_orders(user_id);
CREATE INDEX idx_verified_order_products_order_id ON verified_order_products(verified_order_id);
