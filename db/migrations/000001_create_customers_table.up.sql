CREATE TABLE customers (
   id INT AUTO_INCREMENT PRIMARY KEY,
   created_at DATE NOT NULL,
   updated_at DATE NOT NULL,
   deleted_at DATE,
   name VARCHAR(255) NOT NULL
);