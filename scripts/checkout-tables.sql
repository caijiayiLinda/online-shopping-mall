-- Combined SQL for checkout tables initialization and modifications
use shopping_mall;

-- Create orders table with all columns (including those that were previously added via ALTER)
CREATE TABLE IF NOT EXISTS orders (
    id INT AUTO_INCREMENT PRIMARY KEY,
    currency VARCHAR(3) NOT NULL,
    merchant_email VARCHAR(255) NOT NULL,
    salt VARCHAR(255) NOT NULL,
    total_price DECIMAL(10,2) NOT NULL,
    user_id INT NULL,  -- Nullable for guest users
    username VARCHAR(255) NOT NULL,
    digest VARCHAR(255) NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'pending',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_user_id (user_id),
    INDEX idx_status (status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- Create order_products table
CREATE TABLE IF NOT EXISTS order_products (
    id INT AUTO_INCREMENT PRIMARY KEY,
    order_id INT NOT NULL,
    product_id INT NOT NULL,
    quantity INT NOT NULL,
    price DECIMAL(10,2) NOT NULL,
    FOREIGN KEY (order_id) REFERENCES orders(id) ON DELETE CASCADE,
    INDEX idx_order_id (order_id),
    INDEX idx_product_id (product_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- Sample data for testing
INSERT INTO orders (currency, merchant_email, salt, total_price, user_id, username, digest, status, created_at)
VALUES 
    ('USD', 'customer1@example.com', 'salt123', 99.99, 1, 'john_doe', 'digest123', 'completed', NOW()),
    ('USD', 'customer2@example.com', 'salt456', 49.99, NULL, 'guest', 'digest456', 'paid', NOW());

INSERT INTO order_products (order_id, product_id, quantity, price)
VALUES
    (1, 101, 2, 29.99),
    (1, 102, 1, 39.99),
    (2, 103, 3, 16.66);
