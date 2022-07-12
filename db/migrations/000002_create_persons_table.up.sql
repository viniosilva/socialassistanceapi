CREATE TABLE persons (
   id          INT            AUTO_INCREMENT PRIMARY KEY,
   created_at  DATETIME       NOT NULL,
   updated_at  DATETIME       NOT NULL,
   deleted_at  DATETIME,
   address_id  INT            NOT NULL,
   name        VARCHAR(255)   NOT NULL,
   CONSTRAINT persons_addresses_fk  FOREIGN KEY (address_id)   REFERENCES addresses (id)
);