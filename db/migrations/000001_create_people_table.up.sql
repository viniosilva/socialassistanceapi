CREATE TABLE people (
   id INT AUTO_INCREMENT PRIMARY KEY,
   created_at DATETIME NOT NULL,
   updated_at DATETIME NOT NULL,
   deleted_at DATETIME,
   name VARCHAR(255) NOT NULL
);