CREATE TABLE resources_to_families (
   id          INT            AUTO_INCREMENT PRIMARY KEY,
   created_at  DATETIME       NOT NULL,
   resource_id INT            NOT NULL,
   family_id  INT             NOT NULL,
   quantity    DECIMAL(5,2)   NOT NULL,
   CONSTRAINT resources_history_resources_fk FOREIGN KEY (resource_id)  REFERENCES resources(id),
   CONSTRAINT resources_history_families_fk FOREIGN KEY (family_id)   REFERENCES families(id)
);