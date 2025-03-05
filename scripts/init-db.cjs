require('dotenv').config();
const mysql = require('mysql2/promise');

async function initializeDatabase() {
  const db = await mysql.createConnection({
    host: process.env.DB_HOST,
    user: process.env.DB_USER,
    password: process.env.DB_PASSWORD,
    database: process.env.DB_NAME
  });

  // Create tables
  await db.execute(`
    CREATE TABLE IF NOT EXISTS categories (
      catid INT PRIMARY KEY AUTO_INCREMENT,
      name VARCHAR(255) NOT NULL
    );
  `);
  await db.execute(`
    CREATE TABLE IF NOT EXISTS products (
      pid INT PRIMARY KEY AUTO_INCREMENT,
      catid INT NOT NULL,
      name VARCHAR(255) NOT NULL,
      price DECIMAL(10,2) NOT NULL,
      description TEXT,
      image_url TEXT,
      FOREIGN KEY(catid) REFERENCES categories(catid)
    );
  `);

  // Insert sample categories
  await db.execute(`INSERT INTO categories (name) VALUES ('Toys')`);
  await db.execute(`INSERT INTO categories (name) VALUES ('Clothing')`);

  // Insert sample products
  await db.execute(`
    INSERT INTO products (catid, name, price, description, image_url)
    VALUES (1, 'Wireless Mouse', 29.99, 'Ergonomic wireless mouse with long battery life', '')
  `);
  await db.execute(`
    INSERT INTO products (catid, name, price, description, image_url)
    VALUES (1, 'Bluetooth Headphones', 99.99, 'Noise-cancelling Bluetooth headphones', '')
  `);
  await db.execute(`
    INSERT INTO products (catid, name, price, description, image_url)
    VALUES (2, "Men's T-Shirt", 19.99, "100% Cotton crew neck t-shirt", '')
  `);
  await db.execute(`
    INSERT INTO products (catid, name, price, description, image_url)
    VALUES (2, "Women's Dress", 49.99, "Floral print summer dress", '')
  `);

  console.log('Database initialized successfully');
  await db.end();
}

initializeDatabase().catch(err => {
  console.error('Error initializing database:', err);
});
