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
      thumbnail_url TEXT,
      FOREIGN KEY(catid) REFERENCES categories(catid)
    );
  `);
  await db.execute(`
    CREATE TABLE IF NOT EXISTS users (
      id INT PRIMARY KEY AUTO_INCREMENT,
      email TEXT NOT NULL,
      password TEXT,
      admin BOOL
    );
  `);

  // Insert categories
  await db.execute(`INSERT INTO categories (name) VALUES ('Clothing')`);
  await db.execute(`INSERT INTO categories (name) VALUES ('Tools')`);
  await db.execute(`INSERT INTO categories (name) VALUES ('Toys')`);
  await db.execute(`INSERT INTO categories (name) VALUES ('Beauty')`);
  await db.execute(`INSERT INTO categories (name) VALUES ('Pets')`);

  // Insert users
  await db.execute(`
    INSERT INTO users (id, email, password, admin) VALUES (1, "1155229013@link.cuhk.edu.hk", "1155229013", 1)
    `);

  // await db.execute(`
  //   //   INSERT INTO users (id, email, password, admin_flag)
  //   //   VALUES (2, "", "1155229013", 0)
  //   // `);  
  // // Insert products
  // await db.execute(`
  //   INSERT INTO products (catid, name, price, description, image_url)
  //   VALUES (1, "Women's Plus Pleated Midi Dress", 34.98, 'Material: 100% Polyester', '/images/dress.jpg')
  // `);
  // await db.execute(`
  //   INSERT INTO products (catid, name, price, description, image_url)
  //   VALUES (2, 'VQJTCVLY Cordless Drill', 35.49, '21 Voltage & 2 Variable Speeds', '/images/drill.jpg')
  // `);
  // await db.execute(`
  //   INSERT INTO products (catid, name, price, description, image_url)
  //   VALUES (3, 'Hot Wheels Set of 8 Basic Toy Cars & Trucks', 8.88, "It's an instant collection with a set of 8 Hot Wheels, including 1 exclusive vehicle!", '/images/toy.jpg')
  // `);
  // await db.execute(`
  //   INSERT INTO products (catid, name, price, description, image_url)
  //   VALUES (4, 'Maybelline Super Stay Teddy Tint, Long Lasting Matte Lip Stain', 9.97, "Meet Super Stay Teddy Tint, Maybelline's teddy-soft Lip tint that lasts. Now you can tint Lips in teddy-soft color for a plush, light feel that lasts all day. This no transfer Lipcolor lasts up to 12 hours", '/images/lip.jpg')
  // `);
  // await db.execute(`
  //   INSERT INTO products (catid, name, price, description, image_url)
  //   VALUES (5, 'Meow Mix Original Choice Dry Cat Food, 16 Pound Bag', 16.98, 'Contains one (1) 16-pound bag of Meow Mix Original Choice Dry Cat Food, now with a new look', '/images/cat_food.jpg')
  // `);

  console.log('Database initialized successfully');
  await db.end();
}

initializeDatabase().catch(err => {
  console.error('Error initializing database:', err);
});
