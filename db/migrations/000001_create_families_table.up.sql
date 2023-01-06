CREATE TABLE families (
   id             INT         AUTO_INCREMENT PRIMARY KEY,
   created_at     DATETIME    NOT NULL,
   updated_at     DATETIME    NOT NULL,
   deleted_at     DATETIME,
   country        VARCHAR(2)  NOT NULL,
   state          VARCHAR(2)  NOT NULL,
   city           TEXT        NOT NULL,
   neighborhood   TEXT        NOT NULL,
   street         TEXT        NOT NULL,
   number         VARCHAR(15) NOT NULL,
   complement     TEXT        NOT NULL,
   zipcode        VARCHAR(16) NOT NULL
);