ALTER TABLE persons
   ADD address_id INT NOT NULL;

ALTER TABLE persons
   ADD CONSTRAINT addresses_fk
   FOREIGN KEY (address_id)
   REFERENCES addresses (id)