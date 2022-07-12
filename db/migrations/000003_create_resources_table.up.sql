CREATE TABLE resources (
   id          INT            AUTO_INCREMENT PRIMARY KEY,
   created_at  DATETIME       NOT NULL,
   updated_at  DATETIME       NOT NULL,
   name        TEXT           NOT NULL,
	amount      DECIMAL(5,2)   NOT NULL,
	measurement VARCHAR(16)    NOT NULL,
   quantity    DECIMAL(5,2)   NOT NULL
);