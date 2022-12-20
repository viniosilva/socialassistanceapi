CREATE TABLE resources_to_addresses (
   id          INT            AUTO_INCREMENT PRIMARY KEY,
   created_at  DATETIME       NOT NULL,
   resource_id INT            NOT NULL,
   address_id  INT            NOT NULL,
   quantity    DECIMAL(5,2)   NOT NULL,
   CONSTRAINT resources_history_resources_fk FOREIGN KEY (resource_id)  REFERENCES resources(id),
   CONSTRAINT resources_history_addresses_fk FOREIGN KEY (address_id)   REFERENCES addresses(id)
);