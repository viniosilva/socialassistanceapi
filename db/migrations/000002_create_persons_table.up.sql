CREATE TABLE persons (
   id          INT            AUTO_INCREMENT PRIMARY KEY,
   created_at  DATETIME       NOT NULL,
   updated_at  DATETIME       NOT NULL,
   deleted_at  DATETIME,
   family_id  INT             NOT NULL,
   name        VARCHAR(255)   NOT NULL,
   CONSTRAINT persons_families_fk  FOREIGN KEY (family_id)   REFERENCES families (id)
);