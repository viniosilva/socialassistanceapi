CREATE TABLE addresses (
   id INT AUTO_INCREMENT PRIMARY KEY,
   created_at DATETIME NOT NULL,
   updated_at DATETIME NOT NULL,
   deleted_at DATETIME,
   country VARCHAR(2) NOT NULL,
   state VARCHAR(2) NOT NULL,
   city VARCHAR(128) NOT NULL,
   neighborhood VARCHAR(128) NOT NULL,
   street VARCHAR(128) NOT NULL,
   number VARCHAR(15) NOT NULL,
   complement VARCHAR(128) NOT NULL,
   zipcode VARCHAR(16) NOT NULL
);